#!/bin/bash
# cloud-record.sh <state> [<reply file>] [<session log>]
#
# Called by the other cloud hooks, never registered itself. It adds one commit
# to the track's record branch, <work branch>-record, and pushes it, without
# touching the agent's working tree or index. <state> is running, waiting,
# ended or failed. The application reads the branch's tip to learn where the
# turn stands (remote-agent-1.md, sections 6 and 7):
#
#   state.json                     {"state", "conversation", "turn", "head", "at"}
#   <conversation>/<turn>/message.md    what the turn was asked
#   <conversation>/<turn>/reply.md      the agent's reply, once there is one
#   <conversation>/<turn>/log.jsonl     the lines the session's own log gained
#                                       since the previous report
#   <conversation>/<turn>/pictures/<run>/...
#                                  each picture the reply links to under
#                                  .proof/runs/, at its path below that
#                                  folder, and its run's proof.json
#   <conversation>/<turn>/questions.json
#                                  in an ended report only: the questions the
#                                  agent asked with AskUserQuestion, which
#                                  ended its turn (cloud-no-questions.sh);
#                                  removed from every other ended report, so
#                                  one never stands for a later reply
#
# Each report is a new commit on the previous one, so the branch is never
# forced. A session whose machine was replaced starts from the branch as
# GitHub has it, fetched into refs/kvit/record below.
[ "$CLAUDE_CODE_REMOTE" = "true" ] || exit 0
cd "$CLAUDE_PROJECT_DIR" || exit 0
state=.git/kvit-cloud.json
[ -f "$state" ] || exit 0

# Progress reports run in the background and can meet the turn's last one.
exec 9>.git/kvit-record.lock
flock 9

branch=$(jq -r .branch "$state")
conversation=$(jq -r .conversation "$state")
turn=$(jq -r .turn "$state")
dir="$conversation/$turn"
export GIT_INDEX_FILE=.git/kvit-record.index

# The record branch continues from GitHub's tip on a machine that has none of
# it yet: a new session on a track that had one before, or a machine that
# replaced a reclaimed one.
if ! git rev-parse -q --verify refs/kvit/record >/dev/null \
   && git fetch -q origin "$branch-record" 2>/dev/null; then
    git update-ref refs/kvit/record FETCH_HEAD
fi
parent=$(git rev-parse -q --verify refs/kvit/record || true)
if [ -n "$parent" ]; then
    git read-tree "$parent"
else
    git read-tree --empty
fi
put() {
    git update-index --add --cacheinfo \
        "100644,$(git hash-object -w "$2"),$1"
}
jq -n --arg s "$1" --arg c "$conversation" --argjson t "$turn" \
      --arg h "$(git rev-parse HEAD)" --arg at "$(date -u +%FT%TZ)" \
      '{state: $s, conversation: $c, turn: $t, head: $h, at: $at}' \
      > .git/kvit-record-state.json
put state.json .git/kvit-record-state.json
jq -r .message "$state" > .git/kvit-record-message.md
put "$dir/message.md" .git/kvit-record-message.md
if [ -n "${2:-}" ] && [ -f "$2" ]; then
    put "$dir/reply.md" "$2"
    # Pictures made by tools/prove-fix.sh exist only on this machine, and the
    # reply's links to them name this machine's paths. They go with the reply,
    # and the application copies them out and points the links at the copies
    # (CloudWatcher). A link that leaves the folder is not followed.
    grep -oE '!\[[^]]*\]\(<?[^)> ]*\.proof/runs/[^)> ]+' "$2" \
        | sed -E 's|.*\.proof/runs/||' | sort -u \
        | while IFS= read -r picture; do
            case "/$picture/" in */../*|*/./*|//*) continue ;; esac
            [ -f ".proof/runs/$picture" ] || continue
            put "$dir/pictures/$picture" ".proof/runs/$picture"
            run="${picture%%/*}"
            if [ -f ".proof/runs/$run/proof.json" ]; then
                put "$dir/pictures/$run/proof.json" ".proof/runs/$run/proof.json"
            fi
        done
fi
if [ "$1" = ended ]; then
    if [ -s .git/kvit-questions.json ] \
       && [ "$(jq 'length' .git/kvit-questions.json 2>/dev/null)" -gt 0 ] 2>/dev/null; then
        put "$dir/questions.json" .git/kvit-questions.json
    else
        git update-index --force-remove "$dir/questions.json"
    fi
    rm -f .git/kvit-questions.json
fi
# Only what the log gained since the last report, in whole lines. A session's
# log was 562 KB after two short turns on 2026-09-24, and a report carrying all
# of it once a minute is several megabytes an hour for a few new lines. The
# application reads every report in order, so the pieces add up to the log.
if [ -n "${3:-}" ] && [ -f "$3" ]; then
    offsets=.git/kvit-record-log-lines
    lines=$(wc -l < "$3")
    previous=0
    if [ -f "$offsets" ] && [ "$(sed -n 1p "$offsets")" = "$3" ]; then
        previous=$(sed -n 2p "$offsets")
    fi
    [ "${previous:-0}" -le "$lines" ] 2>/dev/null || previous=0
    head -n "$lines" "$3" | tail -n +"$((previous + 1))" > .git/kvit-record-log.jsonl
    printf '%s\n%s\n' "$3" "$lines" > "$offsets"
    put "$dir/log.jsonl" .git/kvit-record-log.jsonl
fi

tree=$(git write-tree)
commit=$(git commit-tree "$tree" ${parent:+-p "$parent"} \
             -m "$conversation turn $turn: $1")
git update-ref refs/kvit/record "$commit"
git push -q origin "$commit:refs/heads/$branch-record"
