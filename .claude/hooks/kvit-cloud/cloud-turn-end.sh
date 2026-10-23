#!/bin/bash
# Stop and StopFailure hook, registered in .claude/settings.json. It acts only
# in a cloud session Kvit Works is driving: one whose current turn is in
# .git/kvit-cloud.json, written by cloud-turn-start.sh or, for the message the
# session was created with, rebuilt here from the session's log.
#
# It does at the end of a cloud turn what the application does at the end of a
# local one: commits what the turn changed, with the message the application
# writes, kvit-agent(session=<conversation>,turn=<n>): settle turn, so the log
# of who changed which files is rebuilt from git the same way for both. Then it
# pushes the branch and reports the turn's end with the reply.
#
# The cloud machine has a Stop hook of its own, ~/.claude/stop-hook-git-check.sh,
# which runs at the same time as this one and keeps the turn going when it finds
# uncommitted changes or unpushed commits. It usually looks before this commits,
# so a turn that changed files gets one more step, and this hook runs a second
# time. That second run finds nothing new and reports nothing, so the reply the
# application reads is the one recorded the first time.
#
# A push that fails is shown to the agent once, by exiting with status 2, which
# keeps the turn going; the second time (stop_hook_active) the failure is
# reported instead of blocking again.
#
# A project may have two tools of its own that settle things before anything
# is committed, and each is run only when the project has it.
# tools/prove-fix.sh, run in place here, takes the fix's files out of the tree
# while it builds the code before the fix, and names its run in .proof/swapped
# until it has put them back; a turn that ends then would commit the old files
# as the turn's work. The agent is told once, and the second time the files are
# put back here. Then tools/proof_gate.py decides whether the reply needs
# pictures of the fix. It runs from here rather than as a Stop hook of its own
# because Claude Code starts every Stop hook at once, and this one would report
# the turn ended while that one sent the agent back; when it refuses, its
# answer is passed on and nothing is committed or reported.
[ "$CLAUDE_CODE_REMOTE" = "true" ] || exit 0
hooks=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) || exit 0
cd "$CLAUDE_PROJECT_DIR" || exit 0
. "$hooks"/cloud-hook-common.sh
input=$(cat)
log=$(jq -r '.transcript_path // ""' <<<"$input")
kvit_ensure_state "$log"
[ -f "$kvit_state" ] || exit 0
branch=$(jq -r .branch "$kvit_state")
conversation=$(jq -r .conversation "$kvit_state")
turn=$(jq -r .turn "$kvit_state")
reply=.git/kvit-reply.md
record="$hooks/cloud-record.sh"

if [ "$(jq -r .hook_event_name <<<"$input")" = StopFailure ]; then
    jq -r '.last_assistant_message // .error // ""' <<<"$input" > "$reply"
    exec bash "$record" failed "$reply"
fi
# cloud-turn-start.sh has already told the agent why nothing can be pushed;
# the application is told too, so it does not wait for a report that would
# never come.
kvit_fetch_work "$branch" || true
if ! kvit_on_track "$branch"; then
    current=$(git branch --show-current)
    echo "The session is on ${current:-a detached HEAD}, whose history is not" \
         "the track's; nothing was committed or pushed." > "$reply"
    exec bash "$record" failed "$reply" "$log"
fi

if [ -f .proof/swapped ]; then
    if [ "$(jq -r .stop_hook_active <<<"$input")" != true ]; then
        echo "tools/prove-fix.sh has the fix's files out of the tree (it is still" \
             "running, or it was stopped). Let it finish, or run" \
             "tools/prove-fix.sh --restore, before ending the turn." >&2
        exit 2
    fi
    "$CLAUDE_PROJECT_DIR/tools/prove-fix.sh" --restore >&2 || exit 0
fi
gate="$CLAUDE_PROJECT_DIR/tools/proof_gate.py"
if [ -f "$gate" ] && command -v python3 >/dev/null; then
    verdict=$(PROOF_GATE_CLOUD_HOOK=1 python3 "$gate" <<<"$input")
    if [ -n "$verdict" ] \
       && [ "$(jq -r '.decision // ""' <<<"$verdict" 2>/dev/null)" = block ]; then
        printf '%s\n' "$verdict"
        exit 0
    fi
fi

git add -A
if ! git diff --cached --quiet; then
    git commit -q -m "kvit-agent(session=$conversation,turn=$turn): settle turn"
fi
if ! error=$(git push -q origin "HEAD:refs/heads/$branch" 2>&1); then
    if [ "$(jq -r .stop_hook_active <<<"$input")" != true ]; then
        echo "The turn's commits could not be pushed to $branch: $error" >&2
        exit 2
    fi
    echo "Pushing to $branch failed: $error" > "$reply"
    exec bash "$record" failed "$reply"
fi

head=$(git rev-parse HEAD)
git update-ref "refs/remotes/origin/$branch" "$head"
# A turn whose end was reported already ends again in two ways: the machine's
# own Stop hook kept it going for one more step, or the session went on with
# it when a background command finished (cloud-turn-start.sh). The first is
# reported only when that step changed files; the second always is. Either
# way the report holds only what the agent said since, which the application
# adds to the conversation as a further reply.
if ! kvit_turn_open && [ "$(jq -r '.endedHead // ""' "$kvit_state")" = "$head" ]; then
    exit 0
fi
# Every reply of the turn since the last report, rather than the last message
# only (cloud-reply.py); the last message alone where that cannot run.
mark=$(jq -r '.replyMark // ""' "$kvit_state")
collected=""
if command -v python3 >/dev/null; then
    collected=$(python3 "$hooks/cloud-reply.py" "$log" "$mark" \
                    <<<"$input" 2>/dev/null) || collected=""
fi
if [ -n "$(jq -r '.reply // ""' <<<"$collected" 2>/dev/null)" ]; then
    jq -r .reply <<<"$collected" > "$reply"
    mark=$(jq -r '.last // ""' <<<"$collected")
else
    jq -r '.last_assistant_message // ""' <<<"$input" > "$reply"
fi
bash "$record" ended "$reply" "$log" || exit 0
jq --arg h "$head" --arg m "$mark" '.endedHead = $h | .replyMark = $m | del(.continued)' \
    "$kvit_state" > "$kvit_state.new" && mv "$kvit_state.new" "$kvit_state"
