#!/usr/bin/env python3
"""Find past cases where the owner critiqued a screen and an agent answered
with a revised picture, in the session transcripts of the Kvit apps.

Each Kvit project keeps its agent sessions in
<project>/.kvit/sessions/<id>/chat.md. A transcript is a sequence of turns,
each opened by a marker line

    <!--kvit turn=N role=user|assistant|system ...-->

The owner speaks in two places: a user turn, and an answer to a question the
agent asked inside its own turn, written as a quoted line

    > <question> — **<answer>**

A picture the owner pasted is written `![pasted image|320](path)`; its file
is named message-*.png or capture-*.png. A picture an agent made has a
descriptive caption and lives in the session's pictures/ folder.

A case ("candidate") is an owner utterance with a before picture, followed
within a few events by agent pictures (the after), and optionally the
owner's next words (the verdict). The before picture is either

  rule A: a picture the owner pasted in that utterance (any words), or
  rule B: the agent's own picture(s) shown in the agent text just before it,
          when the owner pasted nothing and the words hold a critique word.

Rule A takes every paste because the owner's critiques use too many
different words for a word list; the hand check separates screen critiques
from bug reports. A user turn written by another project's agent (a brief)
is not the owner speaking and is skipped.

Usage, run from anywhere (reads only; writes nothing unless told):

    python3 -I scan.py scan [--rule A|B|AB] [--tsv FILE]
    python3 -I scan.py show ID [ID ...]
    python3 -I scan.py counts
    python3 -I scan.py pastes [--tsv FILE]
    python3 -I scan.py show-paste ID [ID ...]

`pastes` lists the owner's pastes that no agent picture followed: critiques
whose fix was never shown as a picture, and bug reports. They give no
before/after pair, but they still record what the owner found wrong.

ROOT, the folder holding the projects, defaults to the parent of the
kvit-ui checkout this script sits in; override with KVIT_PROJECTS_ROOT.

A project's sessions are read from the first of ROOT/<project> and
/kvit/projects/<project> whose .kvit/sessions folder holds any session
(kvit-ui's checkout under ROOT keeps none; /kvit/projects/kvit-ui does).

kvit-dictate's sessions live only on the owner's Windows PC
(D:/projects/kvit-dictate/.kvit). Copy its `sessions` and `pasted` folders
here first, then name the copy in KVIT_EXTRA_SESSIONS as short=folder pairs,
comma-separated, each folder being a copied `sessions` folder (a `pasted`
folder beside it is used for pictures):

    KVIT_EXTRA_SESSIONS=dictate=~/.cache/kvit-ux-mining/kvit-dictate/sessions

A picture path in a copied transcript points at the PC or at the original
folder, so `show` looks a missing picture up by name in the session's
pictures/ folder and in the pasted/assets folder beside the sessions.
"""

import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.environ.get("KVIT_PROJECTS_ROOT",
                      os.path.dirname(os.path.dirname(HERE)))
PROJECTS = {"works": "kvit-works", "hub": "kvit-hub",
            "notes": "kvit-notes", "cash": "kvit-cash",
            "coder": "kvit-coder", "ui": "kvit-ui"}
PROJECT_PARENTS = [ROOT, "/kvit/projects"]

# Sessions left out: the ones holding this research itself, whose pastes
# are screenshots of documents rather than of an app.
EXCLUDE = {"20261111-4f5afb", "20261124-01d655"}

TURN = re.compile(r"^<!--kvit turn=(\d+) role=(\w+)[^>]*-->\s*$")
IMG = re.compile(r"!\[([^\]]*)\]\(([^)\s]+)\)")
ANSWER = re.compile(r"—\s*\*\*(.+?)\*\*", re.S)

# Critique words for rule B, matched case-insensitively: the list of the
# first pass in agents-ux-mining.md, and words found in the owner's
# critiques that it missed.
CRITIQUE = re.compile(
    r"not good|confusing|not what i|wrong|awkward|\bodd\b|don'?t like|"
    r"should|jargon|intuitive|very bad|not right|mess\b|needs to be fixed|"
    r"looks like|feels like|can'?t tell|\bshows\b|"
    r"\bbad\b|ugly|hard to|don'?t understand|doesn'?t look|no need|"
    r"needs to|too (much|many|big|small|large|wide|narrow|heavy)|padding|"
    r"\bgaps?\b|jarring|clutter|unclear|not clear|weird|strange|"
    r"why (is|does|do)|instead|improve|crap|confus|misalign|alligned|"
    r"aligned|don'?t see|can'?t see|not visible|hidden", re.I)

# A user turn the app writes by itself after a reply: the reply checks, the
# word limit, the reminder to ask what to do next.
APP_MESSAGE = re.compile(
    r"^\s*(Now ask your question again|Your turn changed files|"
    r"The checks found these problems|Your (last )?reply (is|was|has)|"
    r"Your question (is|was|has))", re.I)

# A user turn that is a brief from another project's agent, not the owner.
BRIEF = re.compile(r"^\s*(#+ *Brief\b|\*\*Task:|# .*request\b)|"
                   r"Written \d{4}-\d\d-\d\d by the agent of", re.I | re.M)

# Owner replies that carry no judgment: they only let the agent go on.
EMPTY_REPLY = re.compile(
    r"^\W*(continue|go on|go ahead|ok|okay|yes|proceed|next)\W*$", re.I)

AFTER_WINDOW = 6   # events searched after the critique for agent pictures


def before_picture(alt, path):
    """An agent picture showing only the old state, by caption or name.

    A picture of both states ("before and after") is an after picture."""
    name, cap = os.path.basename(path).lower(), alt.lower()
    if "after" in cap or "after" in name:
        return False
    return (cap.lstrip("*_ ").startswith("before")
            or name.startswith("before") or "-before" in name)


def owner_picture(alt, path):
    name = os.path.basename(path)
    return (alt.lower().startswith(("pasted image", "attached"))
            or "/pasted/" in path
            or name.startswith(("message-", "capture-")))


def parse(path):
    """Return the transcript as a list of events, in order.

    An event is a dict: who ('owner' or 'agent'), turn, line (1-based line
    of its first text), text, pictures (list of (caption, path)).
    """
    with open(path, encoding="utf-8", errors="replace") as f:
        lines = f.read().split("\n")
    events = []
    role, turn = None, 0
    buf, buf_start = [], 0

    def flush_agent(chunk, start):
        """Split an assistant turn into agent text and owner answers."""
        i = 0
        cur, cur_start = [], start
        while i < len(chunk):
            ln = chunk[i]
            if ln.startswith(">"):
                run, run_start = [], start + i
                while i < len(chunk) and chunk[i].startswith(">"):
                    run.append(re.sub(r"^>\s?", "", chunk[i]))
                    i += 1
                joined = "\n".join(run)
                answers = ANSWER.findall(joined)
                if answers and all(APP_MESSAGE.match(a) for a in answers):
                    continue    # the app's own message, shown as an answer
                if not answers:
                    cur.extend(">" + r for r in run)
                    continue
                if cur:
                    events.append(agent_event(cur, cur_start))
                    cur = []
                # Pictures the owner pasted with the answer: inside the
                # quoted answer, or on the lines after the quote.
                pics = [(a, p) for ans in answers for a, p in IMG.findall(ans)
                        if owner_picture(a, p)]
                j = i
                while j < len(chunk) and (not chunk[j].strip()
                                          or IMG.search(chunk[j])):
                    for alt, p in IMG.findall(chunk[j]):
                        if owner_picture(alt, p):
                            pics.append((alt, p))
                    if IMG.search(chunk[j]) and not all(
                            owner_picture(a, p)
                            for a, p in IMG.findall(chunk[j])):
                        break
                    j += 1
                if pics:
                    i = j
                events.append({"who": "owner", "turn": turn,
                               "line": run_start, "text": "\n".join(answers),
                               "question": joined, "pictures": pics,
                               "kind": "answer"})
                cur_start = start + i
                continue
            if not cur:
                cur_start = start + i
            cur.append(ln)
            i += 1
        if cur:
            events.append(agent_event(cur, cur_start))

    def agent_event(chunk, start):
        text = "\n".join(chunk)
        pics = [(a, p) for a, p in IMG.findall(text)
                if not owner_picture(a, p)]
        return {"who": "agent", "turn": turn, "line": start, "text": text,
                "pictures": pics}

    def flush():
        if role == "user" and not APP_MESSAGE.match("\n".join(buf)):
            text = "\n".join(buf)
            events.append({"who": "owner", "turn": turn, "line": buf_start,
                           "text": text, "pictures": IMG.findall(text),
                           "kind": "turn"})
        elif role == "assistant":
            flush_agent(buf, buf_start)

    start = 0
    if lines and lines[0].strip() == "---":     # front matter
        for k in range(1, len(lines)):
            if lines[k].strip() == "---":
                start = k + 1
                break
    for n in range(start, len(lines)):
        m = TURN.match(lines[n])
        if m:
            flush()
            turn, role = int(m.group(1)), m.group(2)
            buf, buf_start = [], n + 2
            continue
        buf.append(lines[n])
    flush()
    return [e for e in events if e["text"].strip() or e["pictures"]]


def session_roots():
    for short, folder in PROJECTS.items():
        for parent in PROJECT_PARENTS:
            base = os.path.join(parent, folder, ".kvit", "sessions")
            if os.path.isdir(base) and os.listdir(base):
                yield short, base
                break
    extra = os.environ.get("KVIT_EXTRA_SESSIONS", "")
    for pair in filter(None, (p.strip() for p in extra.split(","))):
        short, _, folder = pair.partition("=")
        yield short, os.path.expanduser(folder)


def sessions():
    for short, base in session_roots():
        if not os.path.isdir(base):
            continue
        for sid in sorted(os.listdir(base)):
            path = os.path.join(base, sid, "chat.md")
            if sid not in EXCLUDE and os.path.isfile(path):
                yield short, sid, path


def resolve(picture, transcript):
    """Where a picture named in a transcript is on this machine."""
    if os.path.isfile(picture):
        return picture
    name = os.path.basename(picture.replace("\\", "/"))
    sdir = os.path.dirname(transcript)
    for p in (os.path.join(sdir, "pictures", name),
              os.path.join(os.path.dirname(os.path.dirname(sdir)),
                           "pasted", "assets", name)):
        if os.path.isfile(p):
            return p
    return picture + "   (NOT FOUND HERE)"


def critique_text(e):
    # Words of the owner only: drop picture links, which hold file names.
    return IMG.sub(" ", e["text"])


def candidates(events, rules="AB"):
    out = []
    for i, e in enumerate(events):
        if e["who"] != "owner" or BRIEF.search(e["text"]):
            continue
        before = [(a, p) for a, p in e["pictures"] if owner_picture(a, p)]
        rule = "A"
        if not before:
            if not CRITIQUE.search(critique_text(e)):
                continue
            prev = events[i - 1] if i > 0 else None
            if prev and prev["who"] == "agent" and prev["pictures"]:
                before, rule = list(prev["pictures"]), "B"
        if not before or rule not in rules:
            continue
        after, after_idx, j = [], None, i + 1
        while j < len(events) and j <= i + AFTER_WINDOW:
            ev = events[j]
            shown = [(a, p) for a, p in ev["pictures"]
                     if not before_picture(a, p)]
            if ev["who"] == "agent" and shown:
                after.extend(shown)
                after_idx = j
            elif ev["who"] == "owner" and after:
                if not EMPTY_REPLY.match(critique_text(ev)):
                    break
            j += 1
        if not after:
            continue
        verdict = None
        for ev in events[after_idx + 1:]:
            if ev["who"] == "owner" and critique_text(ev).strip():
                verdict = ev
                break
        out.append({"index": i, "rule": rule, "owner": e, "before": before,
                    "after": after, "after_index": after_idx,
                    "verdict": verdict})
    return out


def lone_pastes(events):
    """Owner utterances with a pasted picture and no case: no after."""
    cased = {c["index"] for c in candidates(events)}
    out = []
    for i, e in enumerate(events):
        if (e["who"] != "owner" or i in cased or BRIEF.search(e["text"])
                or not any(owner_picture(a, p) for a, p in e["pictures"])):
            continue
        reply = events[i + 1] if i + 1 < len(events) else None
        if reply and reply["who"] == "owner":
            # The agent wrote only questions: its words are the question.
            reply = ({"text": reply["question"]} if reply.get("question")
                     else None)
        nxt = next((ev for ev in events[i + 1:] if ev["who"] == "owner"
                    and critique_text(ev).strip()), None)
        out.append({"index": i, "owner": e, "reply": reply, "next": nxt,
                    "before": [(a, p) for a, p in e["pictures"]
                               if owner_picture(a, p)]})
    return out


def all_pastes():
    for short, sid, path in sessions():
        events = parse(path)
        for c in lone_pastes(events):
            yield short, sid, path, events, c


def cmd_pastes(args):
    tsv = args[1] if len(args) > 1 and args[0] == "--tsv" else None
    rows = ["\t".join(["id", "pictures", "owner words"])]
    for short, sid, path, events, c in all_pastes():
        rows.append("\t".join([
            case_id(short, sid, c),
            " ".join(os.path.basename(p) for _, p in c["before"]),
            one_line(c["owner"]["text"], 240)]))
    text = "\n".join(rows) + "\n"
    if tsv:
        with open(tsv, "w", encoding="utf-8") as f:
            f.write(text)
        print(f"{len(rows) - 1} pastes written to {tsv}")
    else:
        sys.stdout.write(text)


def cmd_show_paste(args):
    want = set(args)
    for short, sid, path, events, c in all_pastes():
        cid = case_id(short, sid, c)
        if cid not in want:
            continue
        o = c["owner"]
        print(f"=== {cid}  (transcript {path})")
        if o.get("kind") == "answer":
            print("--- agent question the owner answered:")
            print(o["question"][:1200].strip())
        print("--- OWNER WORDS:")
        print(o["text"].strip()[:3000])
        print("--- PASTED pictures:")
        for alt, p in c["before"]:
            print(f"  {resolve(p, path)}   [{alt}]")
        print("--- AGENT REPLY (start):")
        r = c["reply"]
        print(r["text"].strip()[:2500] if r else "  (none)")
        print("--- OWNER'S NEXT WORDS:")
        n = c["next"]
        print(f"  L{n['line']}: " + n["text"].strip()[:1200] if n else "  (none)")
        print()


def case_id(short, sid, c):
    return f"{short}/{sid}/L{c['owner']['line']}"


def one_line(s, n):
    s = re.sub(r"\s+", " ", IMG.sub("[picture]", s or "")).strip()
    return s if len(s) <= n else s[: n - 1] + "…"


def all_cases(rules="AB"):
    for short, sid, path in sessions():
        events = parse(path)
        for c in candidates(events, rules):
            yield short, sid, path, events, c


def cmd_scan(args):
    rules, tsv = "AB", None
    while args:
        a = args.pop(0)
        if a == "--rule":
            rules = args.pop(0)
        elif a == "--tsv":
            tsv = args.pop(0)
    rows = ["\t".join(["id", "rule", "before", "critique", "after",
                       "verdict_line", "verdict"])]
    for short, sid, path, events, c in all_cases(rules):
        v = c["verdict"]
        rows.append("\t".join([
            case_id(short, sid, c), c["rule"],
            " ".join(os.path.basename(p) for _, p in c["before"]),
            one_line(c["owner"]["text"], 200),
            " ".join(os.path.basename(p) for _, p in c["after"]),
            str(v["line"]) if v else "",
            one_line(v["text"], 160) if v else ""]))
    text = "\n".join(rows) + "\n"
    if tsv:
        with open(tsv, "w", encoding="utf-8") as f:
            f.write(text)
        print(f"{len(rows) - 1} cases written to {tsv}")
    else:
        sys.stdout.write(text)


def cmd_counts(args):
    per = {}
    pics_owner = pics_agent = 0
    for short, sid, path in sessions():
        events = parse(path)
        per.setdefault(short, [0, 0, 0, 0])
        per[short][0] += 1
        if any(e["pictures"] for e in events):
            per[short][1] += 1
        for e in events:
            if e["who"] == "owner":
                pics_owner += len(e["pictures"])
            else:
                pics_agent += len(e["pictures"])
        for c in candidates(events):
            per[short][2 if c["rule"] == "A" else 3] += 1
    print("project\tsessions\twith pictures\trule A cases\trule B cases")
    for k, v in per.items():
        print(f"{k}\t" + "\t".join(map(str, v)))
    tot = [sum(v[i] for v in per.values()) for i in range(4)]
    print("total\t" + "\t".join(map(str, tot)))
    print(f"owner pictures {pics_owner}, agent pictures {pics_agent}")


def cmd_show(args):
    want = set(args)
    for short, sid, path, events, c in all_cases():
        cid = case_id(short, sid, c)
        if cid not in want:
            continue
        o = c["owner"]
        print(f"=== {cid}  (rule {c['rule']}, transcript {path})")
        if c["rule"] == "B":
            prev = events[c["index"] - 1]
            print("--- agent text before the critique (end):")
            print(prev["text"][-2500:].strip())
        if o.get("kind") == "answer":
            print("--- agent question the owner answered:")
            print(o["question"][:1500].strip())
        print("--- OWNER CRITIQUE:")
        print(o["text"].strip()[:4000])
        print("--- BEFORE pictures:")
        for alt, p in c["before"]:
            print(f"  {resolve(p, path)}   [{alt}]")
        print("--- AGENT RESPONSE (text up to the after pictures):")
        budget = 6000
        for ev in events[c["index"] + 1: c["after_index"] + 1]:
            who = "AGENT" if ev["who"] == "agent" else "OWNER"
            chunk = ev["text"].strip()
            if ev["who"] == "agent" and len(chunk) > 2500:
                chunk = chunk[:2500] + "\n[… cut …]"
            print(f"[{who} L{ev['line']}] {chunk[:budget]}")
            budget -= len(chunk)
            if budget <= 0:
                print("[… response cut …]")
                break
        print("--- AFTER pictures:")
        for alt, p in c["after"]:
            print(f"  {resolve(p, path)}   [{alt}]")
        v = c["verdict"]
        print("--- OWNER'S NEXT WORDS (verdict candidate):")
        if v:
            if v.get("kind") == "answer":
                print("  question: " + one_line(v["question"], 600))
            print(f"  L{v['line']}: " + v["text"].strip()[:2000])
        else:
            print("  (none)")
        print()


def main():
    cmds = {"scan": cmd_scan, "show": cmd_show, "counts": cmd_counts,
            "pastes": cmd_pastes, "show-paste": cmd_show_paste}
    if len(sys.argv) < 2 or sys.argv[1] not in cmds:
        sys.stdout.write(__doc__)
        sys.exit(2)
    cmds[sys.argv[1]](sys.argv[2:])


if __name__ == "__main__":
    main()
