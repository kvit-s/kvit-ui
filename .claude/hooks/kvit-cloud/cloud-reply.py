#!/usr/bin/env python3
"""cloud-reply.py <session log> [<after uuid>] -- the reply a cloud turn gave.

Called by cloud-turn-end.sh with the hook's own input on standard input; never
registered as a hook itself. It prints one JSON object:

    {"reply": "<every reply of the turn, in order>", "last": "<uuid>"}

Claude Code hands the end-of-turn hook `last_assistant_message`, which is the
last message only. A turn that was sent back once -- by the pictures check,
tools/proof_gate.py, or by the machine's own check for uncommitted work --
wrote its report first and something shorter after, so the last message alone
loses the report. This reads the session's own log instead and takes every
message of the agent's since the turn's message that is a reply rather than
work: a message of text alone, or text beside a question. A message that calls
any other tool is the agent narrating its work, and a subagent's messages are
never the reply. That is the rule the application applies to a local turn
(agent/stream/claudestream.go).

The turn starts at the person's message, or at the notice Claude Code gives
the session when a background command finishes; a check that sent the agent
back is not a new message. `<after uuid>` is where the previous report of the
same turn stopped, so a turn reported again says only what came since, and
"last" is where this one stops, for the next.

Claude Code writes the log asynchronously, so the turn's last message may not
be in it yet when the hook runs. The log is read again for up to two seconds
until it is; after that, the hook's `last_assistant_message` is added at the
end. KVIT_CLOUD_REPLY_WAIT_MS changes the two seconds; the cloud turn tests
set it to 0, since the logs they write never hold the last message.
"""

from __future__ import annotations

import json
import os
import sys
import time

HOOK_BOUNCE_PREFIXES = ("Stop hook feedback",)


def entries(path: str) -> list[dict]:
    rows = []
    try:
        with open(path, encoding="utf-8") as log:
            for line in log:
                try:
                    row = json.loads(line)
                except ValueError:
                    continue
                if isinstance(row, dict):
                    rows.append(row)
    except OSError:
        pass
    return rows


def user_text(row: dict) -> str | None:
    """The words of a message the person or the session sent, or None when
    the row is a tool's result or something that is not a message."""
    if row.get("type") != "user" or row.get("isSidechain") or row.get("isMeta"):
        return None
    content = (row.get("message") or {}).get("content")
    if isinstance(content, str):
        return content
    if isinstance(content, list):
        texts = [block.get("text", "") for block in content
                 if isinstance(block, dict) and block.get("type") == "text"]
        if texts and not any(isinstance(block, dict)
                             and block.get("type") == "tool_result"
                             for block in content):
            return "\n".join(texts)
    return None


def collect(rows: list[dict], after: str) -> tuple[list[str], str]:
    start = 0
    for index, row in enumerate(rows):
        text = user_text(row)
        if text is not None and not text.lstrip().startswith(HOOK_BOUNCE_PREFIXES):
            start = index + 1
        if after and row.get("uuid") == after:
            start = max(start, index + 1)
    # message id -> [text, at work]; ids kept in the order first seen
    order: list[str] = []
    messages: dict[str, list] = {}
    for position, row in enumerate(rows[start:]):
        if row.get("type") != "assistant" or row.get("isSidechain"):
            continue
        message = row.get("message") or {}
        key = message.get("id") or f"row-{position}"
        if key not in messages:
            messages[key] = ["", False]
            order.append(key)
        for block in message.get("content") or []:
            if not isinstance(block, dict):
                continue
            if block.get("type") == "text":
                text = block.get("text", "")
                messages[key][0] += ("\n\n" if messages[key][0] else "") + text
            elif block.get("type") == "tool_use" and block.get("name") != "AskUserQuestion":
                messages[key][1] = True
    replies = [messages[key][0].strip() for key in order
               if not messages[key][1] and messages[key][0].strip()]
    last = rows[-1].get("uuid", "") if rows else ""
    return replies, last or after


def main() -> int:
    if len(sys.argv) < 2:
        print(__doc__, file=sys.stderr)
        return 2
    path = sys.argv[1]
    after = sys.argv[2] if len(sys.argv) > 2 else ""
    try:
        hook = json.load(sys.stdin)
    except ValueError:
        hook = {}
    final = (hook.get("last_assistant_message") or "").strip()

    replies, last = collect(entries(path), after)
    try:
        wait_ms = int(os.environ.get("KVIT_CLOUD_REPLY_WAIT_MS", "2000"))
    except ValueError:
        wait_ms = 2000
    for _ in range(max(0, wait_ms) // 100):
        if not final or (replies and replies[-1] == final):
            break
        time.sleep(0.1)
        replies, last = collect(entries(path), after)
    if final and (not replies or replies[-1] != final):
        replies.append(final)
    print(json.dumps({"reply": "\n\n".join(replies), "last": last}))
    return 0


if __name__ == "__main__":
    sys.exit(main())
