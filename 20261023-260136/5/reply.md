With every question decided, the documents are ready to build from. I'd go in this order; each step is a separate piece of work:

1. **Put the skill section in.** Merge `skill-section-draft.md` into `agent/skills/kvit-ui/SKILL.md` now. It's a cheap change, and the agents building screens start writing better labels immediately, before any library work lands. Rules 15–20 (commands, menu bar, shortcuts) stay manual until step 3.

2. **Teach kvit-ui which system it is on.** `UI.Form` (Windows or macOS, with Linux taking Windows), switchable in `uitest.Open` so every gallery page renders in both forms. Every later platform-dependent behaviour stands on this.

3. **Commands, shortcuts and command surfaces.** A `Command` type (name, symbol, shortcut, enabled, run) and a shortcut registry on the window. Standard commands get their standard keys automatically, which fixes the missing Ctrl+S, Ctrl+F, Ctrl+N and Ctrl+W. Then `MenuBar` (the system's on macOS, Kvit-drawn on Windows and Linux) and a toolbar with a More menu, both built from commands. This is the largest step, and it removes the one-command-under-several-names problem at its root.

4. **Dialog by platform.** Button order, Enter on a focused button on Windows, an optional default, no "Confirm" default, a safe button in every dialog, focus returned to the opener, no close symbol on macOS, and warning before input is discarded. This clears most of the 9 places where the library contradicts the rules.

5. **Wording in components.** The title-case converter, `Asks` adding the ellipsis, the library's own strings rewritten (no articles, no "Working…"), the formatting symbols (bold, italic and the rest) added to the catalogue, and the percentage hidden by default.

6. **Tests.** `uitest.CheckLabels` reads one screen's controls (verb first, length, no articles, case per form, no letters standing in for symbols). `uitest.CheckApp` reads every screen of an app together (one name per command, one meaning per symbol, standard shortcuts present). Run both in both forms. Adding them after steps 3–5 means they test the new behaviour instead of failing on what is about to change.

7. **Keyboard.** Access keys on every control, arrow-key focus in tabs, segmented controls and the sidebar, Ctrl+Tab, type-ahead search in lists, Escape everywhere.

8. **The rest of library-gaps.md.** Settings page and rows, form rows with visible labels, Tip, back history, window close and restore, notifications, undo stack.

9. **Review the apps.** Use `checklist.md` with a reviewer agent and the new tests on every screen of kvit-cash, kvit-notes and kvit-works. This environment couldn't reach those repositories, so this needs a session where they're attached. Fix the screens through the new components, not by hand.

Two smaller items can happen at any point: after step 1, a reviewer agent can start applying `checklist.md` to new screens; and if Microsoft's six keyboard pages turn up, a quick pass would check the KEY rules against them, though you decided the Windows 7 forms stand.

To stay within reach of one session each, step 3 splits into commands with the registry, then the menu bar with the toolbar, and step 4 is a single session.
