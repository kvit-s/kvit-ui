#!/bin/bash
# PreToolUse hook on AskUserQuestion, registered in .claude/settings.json. It
# acts only in a cloud session Kvit Works is driving (one with
# .git/kvit-cloud.json, or whose log has the application's Kvit- lines), where
# nobody is watching claude.ai to answer a question asked mid-turn and the
# turn would wait there indefinitely.
#
# The question is kept rather than lost: its questions and options are saved
# in .git/kvit-questions.json, which cloud-record.sh puts in the turn's report
# as questions.json when the turn ends. The application shows them as the
# buttons a local turn's question gets, and the pick is sent as the next
# message. The tool itself is refused, so the agent ends its turn instead of
# waiting. In a session the owner drives from the web page the tool is left
# alone.
[ "$CLAUDE_CODE_REMOTE" = "true" ] || exit 0
hooks=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd) || exit 0
cd "$CLAUDE_PROJECT_DIR" || exit 0
. "$hooks"/cloud-hook-common.sh
input=$(cat)
kvit_ensure_state "$(jq -r '.transcript_path // ""' <<<"$input")"
[ -f "$kvit_state" ] || exit 0
jq -c '.tool_input.questions // []' <<<"$input" > .git/kvit-questions.json
cat <<'EOF'
{"hookSpecificOutput": {"hookEventName": "PreToolUse", "permissionDecision": "deny", "permissionDecisionReason": "Kvit Works shows this question to the person, with its options, when your turn ends, and their answer comes as your next message. End your turn now; there is no need to write the question out again."}}
EOF
