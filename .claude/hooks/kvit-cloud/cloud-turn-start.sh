#!/bin/bash
# UserPromptSubmit hook, registered in .claude/settings.json. It acts only in a
# Claude Code cloud session that Kvit Works is driving, and exits at once
# everywhere else.
#
# The application ends every message it sends to a cloud session with three
# lines naming the track's GitHub branch, the local conversation and the turn
# (remote-agent-1.md, section 5):
#
#   Kvit-Branch: claude/kvit-<track>
#   Kvit-Conversation: <the conversation's folder name>
#   Kvit-Turn: <n>
#
# Before the agent reads the message, this brings the clone up to what the
# application pushed to that branch, and keeps the three values in
# .git/kvit-cloud.json for the hook that ends the turn. A message typed on
# claude.ai has no such lines: in a session the application started it counts
# as the next turn, and in any other session nothing happens.
#
# The session has to have started on that branch, or, when `claude --cloud`
# started it, on the branch of its own that command makes from it, whatever
# its name (kvit_on_track in cloud-hook-common.sh). Either way the clone is
# brought up to the work branch by a fast-forward and nothing is checked out:
# a session whose history is not the track's has nothing to do with it.
[ "$CLAUDE_CODE_REMOTE" = "true" ] || exit 0
hooks=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) || exit 0
cd "$CLAUDE_PROJECT_DIR" || exit 0

. "$hooks"/cloud-hook-common.sh
prompt=$(jq -r '.prompt // ""')
# A background command the agent started has finished and the session goes on
# with the turn that started it. That is not a message: the turn keeps its
# number, the tree is left as the agent has it, and the application hears
# that the turn is running again. The end of what follows is reported as a
# further end of the same turn.
if kvit_is_notice "$prompt"; then
    [ -f "$kvit_state" ] || exit 0
    jq '.continued = true' "$kvit_state" > "$kvit_state.new" \
        && mv "$kvit_state.new" "$kvit_state"
    bash "$hooks"/cloud-record.sh running || true
    exit 0
fi
kvit_fields "$prompt"
if [ -n "$branch" ]; then
    source=application
else
    [ -f "$kvit_state" ] || exit 0
    branch=$(jq -r .branch "$kvit_state")
    conversation=$(jq -r .conversation "$kvit_state")
    turn=$(( $(jq -r .turn "$kvit_state") + 1 ))
    source=claude.ai
fi
kvit_write_state "$(kvit_message "$prompt")" "$source"

if ! kvit_fetch_work "$branch"; then
    echo "Fetching $branch from GitHub failed, so the tree may be older than" \
         "the track."
fi
if ! kvit_on_track "$branch"; then
    current=$(git branch --show-current)
    echo "This session is on ${current:-a detached HEAD}, whose history is not" \
         "that of the track's branch, $branch, which the message names. The" \
         "session has to have been started from that branch. Say so in your" \
         "reply and change nothing."
    exit 0
fi
if ! git merge -q --ff-only "origin/$branch"; then
    echo "This clone has commits on $branch that GitHub does not. Merge" \
         "origin/$branch into it before anything else, and say so in your reply."
fi
# The session's own branch follows the work branch, where its commits are
# pushed, so that a check comparing the branch with its upstream, as the
# machine's own Stop hook may, finds them pushed. What that hook compares
# against has not been seen.
git branch -q --set-upstream-to="origin/$branch" 2>/dev/null || true
git submodule update -q --init --recursive --filter=blob:none || true
bash "$hooks"/cloud-record.sh running || true
echo "The tree is the track as it stood when this message was sent, at" \
     "$(git rev-parse --short HEAD). Do not commit, switch branches, rebase or" \
     "push: when the turn ends, what it changed is committed and pushed for you."
