The skill section is in `agent/skills/kvit-ui/SKILL.md` as "Words, symbols, commands and shortcuts", with all 28 rules, just before "The workflow for a new screen".

I changed the draft in three ways, so agents don't reach for library features that aren't built yet:

- **Rules that depend on future library work say what to do today:**
  - Write the ellipsis yourself, since there's no `Asks` yet.
  - Give each command one name constant and one handler, since there's no `Command` type.
  - Handle shortcut keys in `Window.OnKeyDown` and set the same key on the command's `MenuItem`.
  - Use the tick for toggles until `MenuItem` takes both words.
  - Write "&" access keys only in menu items, the only place kvit-ui reads them now.
- **Two defaults agents must always replace:** Dialog's "Confirm" label and SearchField's "Filter…" placeholder.
- **The app examples are before-and-after pairs without file and line numbers**, which will go stale as the screens get fixed. The draft in `docs/ui-guidelines/` keeps the full references.

I also added "how to word labels, menus, dialogs and messages…" to the skill's description so agents load it for wording work. "Where things live" now points to `docs/ui-guidelines/`, and the draft and README there note that the section is in the skill. I ran no tests, as only Markdown changed.
