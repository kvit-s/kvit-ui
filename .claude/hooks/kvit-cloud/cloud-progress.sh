#!/bin/bash
# PostToolUse and Notification hook, registered in .claude/settings.json. In a
# cloud session Kvit Works is driving, it reports progress on the record
# branch through cloud-record.sh: after a tool call at most once a minute,
# with what the session's log gained since the last report, and at once when
# a permission prompt is waiting on claude.ai. The PostToolUse registration is
# asynchronous, so the agent never waits for the push.
[ "$CLAUDE_CODE_REMOTE" = "true" ] || exit 0
hooks=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) || exit 0
cd "$CLAUDE_PROJECT_DIR" || exit 0
. "$hooks"/cloud-hook-common.sh
input=$(cat)
kvit_ensure_state "$(jq -r '.transcript_path // ""' <<<"$input")"
[ -f "$kvit_state" ] || exit 0
# Nothing is reported between a turn's reported end and the session going on
# with it: the step the machine's own Stop hook adds would otherwise tell the
# application the turn is running again, with no end to follow.
kvit_turn_open || exit 0

state=running
if [ "$(jq -r .hook_event_name <<<"$input")" = Notification ]; then
    [ "$(jq -r .notification_type <<<"$input")" = permission_prompt ] || exit 0
    state=waiting
else
    last=$(stat -c %Y .git/kvit-progress 2>/dev/null || echo 0)
    [ $(( $(date +%s) - last )) -ge 60 ] || exit 0
fi
touch .git/kvit-progress
exec bash "$hooks"/cloud-record.sh "$state" "" \
    "$(jq -r '.transcript_path // ""' <<<"$input")"
