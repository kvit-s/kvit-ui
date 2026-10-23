# Sourced by the other cloud hooks in this folder; not a hook itself.
#
# The turn a cloud session is on is kept in .git/kvit-cloud.json. It is written
# from the three lines that end every message Kvit Works sends:
#
#   Kvit-Branch: claude/kvit-<track>
#   Kvit-Conversation: <the conversation's folder name>
#   Kvit-Turn: <n>
#
# Leading spaces are allowed before each, since a message pasted from a
# terminal arrives indented.

kvit_state=.git/kvit-cloud.json

# Sets branch, conversation and turn from the lines in "$1"; empty when absent.
kvit_fields() {
    branch=$(grep -oP '^\s*Kvit-Branch:\s*\K\S+' <<<"$1" | tail -1)
    conversation=$(grep -oP '^\s*Kvit-Conversation:\s*\K\S+' <<<"$1" | tail -1)
    turn=$(grep -oP '^\s*Kvit-Turn:\s*\K[0-9]+' <<<"$1" | tail -1)
}

# kvit_is_notice <prompt>: whether the prompt is not a message at all but the
# notice Claude Code gives the session when a background command or monitor
# the agent started has finished. It arrives through UserPromptSubmit like a
# message, after the turn that started the command has ended, and the agent
# goes on with that turn's work (2.1.282, 2026-09-24).
kvit_is_notice() {
    [[ "$1" =~ ^[[:space:]]*\<task-notification\> ]]
}

# kvit_turn_open: whether the current turn is still to be reported as ended,
# either because it has not ended yet or because the session went on with it
# after its end was reported.
kvit_turn_open() {
    [ -z "$(jq -r '.endedHead // ""' "$kvit_state")" ] \
        || [ "$(jq -r '.continued // false' "$kvit_state")" = true ]
}

# The message without those three lines.
kvit_message() {
    grep -vP '^\s*Kvit-(Branch|Conversation|Turn):' <<<"$1"
}

# kvit_fetch_work <work branch>: brings refs/remotes/origin/<work branch> up
# to GitHub's, whichever branches the clone was made to follow.
kvit_fetch_work() {
    git fetch -q origin "+refs/heads/$1:refs/remotes/origin/$1"
}

# kvit_on_track <work branch>: whether the clone works on the track, decided
# by history rather than by the branch's name. A session started from the web
# page on the work branch is on it by name. One started by `claude --cloud`
# from a checkout of it is on a branch of its own at the same commit, named by
# Claude Code: the work branch's name with a suffix
# (claude/kvit-cloud-trial1-9e9b37-x9phqr, 2.1.281) or a name made from the
# message (claude/check-table-rendering-asua45, 2.1.282). Its commits still go
# to the work branch, since the push names it, and a session may push to any
# claude/ branch. So the clone is on the track when its HEAD and the work
# branch as last fetched are the same commit or one descends from the other.
# The record branch shares no history with the work branch and never passes.
kvit_on_track() {
    local tip
    tip=$(git rev-parse -q --verify "refs/remotes/origin/$1") || return 1
    git merge-base --is-ancestor "$tip" HEAD 2>/dev/null \
        || git merge-base --is-ancestor HEAD "$tip" 2>/dev/null
}

# kvit_write_state <message> <source>, from branch, conversation and turn.
kvit_write_state() {
    jq -n --arg b "$branch" --arg c "$conversation" --argjson t "${turn:-0}" \
          --arg m "$1" --arg s "$2" \
          '{branch: $b, conversation: $c, turn: $t, message: $m, source: $s}' \
          > "$kvit_state"
}

# kvit_ensure_state <session log>
#
# The message a cloud session is created with does not pass through the
# UserPromptSubmit hook, so cloud-turn-start.sh never sees the first turn. The
# session's own log does: this finds the latest user message there that has
# the three lines, and writes the state from it when the state is missing or
# names an earlier turn. The clone was made at the branch's tip for that first
# turn, so nothing needs bringing up to date.
kvit_ensure_state() {
    local log="$1" text
    [ -n "$log" ] && [ -f "$log" ] || return 0
    text=$(jq -rR 'fromjson? | select(.type == "user") | .message.content
                   | strings | select(test("(^|\n)\\s*Kvit-Branch:")) | @json' \
               "$log" | tail -1 | jq -r '. // empty' 2>/dev/null)
    [ -n "$text" ] || return 0
    kvit_fields "$text"
    [ -n "$branch" ] || return 0
    if [ -f "$kvit_state" ] \
       && [ "$(jq -r .conversation "$kvit_state")" = "$conversation" ] \
       && [ "$(jq -r .turn "$kvit_state")" -ge "${turn:-0}" ]; then
        return 0
    fi
    kvit_write_state "$(kvit_message "$text")" application
}
