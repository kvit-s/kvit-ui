I've written the guidelines into `docs/ui-guidelines/`, from Apple's 75 pages and the Windows 7 guide's 46 pages (read from your epubs). Microsoft's current Fluent pages were not read: this environment's network policy blocks learn.microsoft.com, so the Windows rules rest on the Windows 7 guide alone. If you allow learn.microsoft.com in the environment's network settings, or send those pages as files like the epubs, the Windows-specific rules can be checked against them (Q-39).

I couldn't reach the kvit-cash, kvit-notes and kvit-works repositories from here. The before-and-after examples are the ones your survey named, plus kvit-ui's own strings, which have the same habits ("Clear the filter", "Close the pane", "Select every row the current filter shows").

The files: README.md, extracted.jsonl, rules.md, open-questions.md, library-gaps.md, checklist.md (45 items) and skill-section-draft.md (26 rules).

**Counts**
- **Pages read:** 121 (Apple 75, Windows 7 46). Apple pages fetched 2026-10-23; each page's latest change-log date is in the README.
- **Rules extracted:** 3,746 (Apple 1,267, Windows 7 2,479), of which 3,217 are kept.
- **Left out:** 224 for visual appearance, 251 for other platforms, 54 not applicable.
- **On Microsoft's checklist:** 458 records are marked.
- **Merged rules:** 198, and every kept record ends up in one. 91 have a separate Windows and macOS form.
- **kvit-ui today:** enforces 26, partly covers 92, contradicts 10, doesn't cover 70.
- **Open questions:** 40.

**The ten rules the apps break most, from your survey**
1. WRD-2: labels of one to three words.
2. WRD-3: no articles on labels.
3. WRD-1: a label starts with the verb that names what happens.
4. CON-1: one name per command (Quick Capture, three ways).
5. CON-2: one symbol, one meaning (the X used for five actions).
6. CMD-22: no letters standing in for symbols (B, I, U, <>).
7. WRD-7: no OK or Confirm on a confirming button.
8. KEY-2: standard shortcuts (Ctrl+S, F, N, W).
9. KEY-14: menus show their shortcuts.
10. WRD-6: one capitalisation style per kind of label.

WRD-22 (the permanent paragraph) and WRD-10 (a "Delete…" that moves to the trash) come next.

**Open questions that block the most later work**
- **Q-3:** how a label gets its Windows and macOS capitalisation. I recommend the app writes sentence case and kvit-ui converts it on macOS.
- **Q-1:** whether apps have a menu bar on Windows and Linux.
- **Q-7:** the default button and Enter in Windows dialogs.
- **Q-11:** how far access keys go.
- **Q-5:** when OK, Yes and No are allowed on Windows.
- **Q-2:** how settings apply. Here I recommend against the Windows 7 text: settings that take effect at once on every system, not a modal Options window.
- **Q-40:** what Linux gets. 70 rules depend on it.
- **Q-39:** whether to read the Fluent pages before building anything.

**Library changes I'd make first, in order**
1. **`UI.Form`:** a Windows or macOS form for each screen, which tests can switch so every screen is checked in both.
2. **A `Command` type and shortcut registry:** one place for each command's name, symbol and shortcut, with the standard commands getting their standard keys. This alone fixes most of problems 3 and 5.
3. **A menu bar:** built from the app's commands. kvit-ui has none on any system today.
4. **The dialog by platform:** button order (it currently puts the confirming button on the right everywhere, `feedback.go:541-552`), Enter pressing the default button (today Space alone presses a button, `button.go:75`), no "Confirm" default label (`feedback.go:424`), and focus returning to where it was when a dialog closes.
5. **A title-case converter and cleanup of the library's own strings.**
6. **Two new tests:** one checks each screen's labels, the other checks every screen of an app together, which is the only way to catch one command under several names and one symbol for several actions.

The changes in `docs/ui-guidelines/` will be committed and pushed when this turn ends.
