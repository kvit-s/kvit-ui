# Kvit Works cloud hooks

Kvit Works writes this folder when it adds a project's cloud hooks or
updates them, and replaces any change made here by hand at the next update.

A Kvit Works cloud agent runs as a Claude Code cloud session on Anthropic's
machines. These scripts, registered in `../settings.json`, do the cloud side
of each of its turns. When a message arrives they bring the session's clone
up to the track's branch, `claude/kvit-<track>`. When the turn ends they
commit what the agent changed, push it to that branch, and report the turn
on `claude/kvit-<track>-record`, which Kvit Works reads to show the reply.
Each script exits at once outside a cloud session, where
`CLAUDE_CODE_REMOTE` is not `true`.

`version` is the version of this set of scripts.
