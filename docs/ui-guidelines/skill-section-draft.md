# Draft: a section for agent/skills/kvit-ui/SKILL.md

This is a draft for the owner, not yet in the skill. It condenses the wording, symbol, command, menu and shortcut rules of rules.md into what an agent needs while building a screen. Each rule names its rule ids so the reasoning can be looked up.

Assumptions it makes, which change it if decided otherwise:

- Q-3: the app writes every label once, in sentence case, and kvit-ui gives the macOS title case. If both forms are supplied instead, rule 6 changes.
- Q-5, Q-7, Q-11: the Windows forms for OK and Yes/No, the default button and access keys, as recommended in open-questions.md.
- The library changes in library-gaps.md steps 1 to 5 exist (`ui.Form`, `Command`, `MenuBar`, Dialog by form). Until they do, rules 15 to 20 are things the agent does by hand.

The before-and-after pairs come from the survey of kvit-cash, kvit-notes and kvit-works of 2026-10-22 and from kvit-ui's own strings and gallery. The application repositories were not available where this draft was written, so the app pairs are the ones the survey named, with its file:line.

---

## Words, symbols, commands and shortcuts

Interface text is not prose. The writing rules for documents and reports ask for whole sentences, spelled-out meaning and no compressed labels; they do not apply to a button, a menu item, a tab or a tooltip. Interface text has its own conventions, shared by Apple and Microsoft: a command starts with its verb, takes one to three words, drops articles, and is named the same everywhere it appears. A sentence on a button is a defect, however clear the sentence is.

The rules below hold on every system unless they say "Windows" or "macOS". kvit-ui knows which system it runs on (`ui.Form`; Linux takes the Windows form) and does the platform differences itself where it can: button order, capitalisation, Ctrl or Cmd, the menu bar.

### Labels

**1. A command is a verb that names what happens.** Buttons and menu items start with a specific verb: Import, Archive, Reconcile, not Manage, Configure or Do it. A statement or a description is not a command. (WRD-1)

- "The data is right" (kvit-cash `app/health/card.go:167`) → "Mark correct"
- "Rows shown" → "Select all"

**2. One to three words.** Use only the words the command needs. When a command seems to need a sentence, the sentence belongs in its Explanation, its dialog, or nowhere. (WRD-2)

- "Take the new default" (kvit-cash `app/configuration/editor.go:116`) → "Reset to default…"
- "Save to the workspace" (kvit-cash `app/configuration/editor.go:157`) → "Save"
- "Move the selected column earlier" (`duallist.go:56`) → "Move up"

**3. No articles, no "my" or "your".** Not on labels, menu items, tabs or tooltips. (WRD-3)

- "Close the preview" (kvit-cash `app/preview/preview.go:56`) → "Back" (it sits on a back arrow)
- "Use the shipped default" (kvit-cash `app/configuration/editor.go:153`) → the same command as rule 2's "Reset to default…"
- "Clear the filter" (`searchfield.go:56`) → "Clear"
- "Dismiss this notice" (`feedback.go:261`) → "Dismiss"

**4. The label says exactly what happens.** Read the handler before naming it. A label that names a different action from the one that runs is the worst kind of wrong. (WRD-1, WRD-10)

- "Add the rows shown" that selects every row (kvit-cash `app/transactions/transactions.go:112`) → "Select all"
- "Delete…" that moves a folder to the trash (kvit-notes `app/actions.go:108`) → "Move to Trash…" on macOS, "Delete…" with the dialog saying it goes to the Recycle Bin on Windows (Q-14)

**5. Ellipsis when the command asks for more.** A command that opens a dialog, a chooser or a confirmation before acting ends with "…" (one character). One that acts at once has none. On Windows, a command whose whole job is to show a window (Options, Properties, About) has none. Set `Asks` and let the component add it. (WRD-4)

- "Export the ledger" opening a dialog (gallery `dialogMoved`) → "Export…"

**6. Write sentence case; kvit-ui does the rest.** Write every label once, in sentence case: "Select all", "Don't save". On macOS the component turns menu items, buttons, tabs and column headings into title case; on Windows and Linux they stay as written. Never write title case yourself, and never mix the two in one menu. (WRD-6)

- kvit-notes File menu (`app/menus.go:39-56`), "Quick Capture…" beside "Quick capture note…" → every item in sentence case

**7. Confirming buttons name the action; never Confirm.** The button that carries out a dialog says what it does. OK appears only on a settings window (Windows) or alone on a message that only informs (macOS). Yes and No appear only together, answering a question that really is yes or no, on Windows. The abandoning button is Cancel; once the work is done it becomes Close. (WRD-7, WRD-8)

- "OK" on Insert Link (kvit-notes `app/links.go:260`) → "Insert"
- "OK" on Confirm import (kvit-notes `app/exchange.go:299`) → "Import"
- Dialog's default ConfirmText "Confirm" (`feedback.go:424`) → always set it: "Delete", "Import", "Send"
- "Delete them" (gallery `dialogDestructive`) → "Delete"

**8. Busy words say what is happening.** Not "Working…", "Loading…" or "Please wait": an -ing phrase naming the operation. (WRD-25)

- Button BusyText "Working…" (`button.go:73`) → "Importing statement…"

**9. A search field says what it searches.** (CTL-15)

- "Filter…" (`searchfield.go:38`) → "Search transactions"

**10. Tooltips are short and add something.** A symbol-only button's Label is its tooltip and its screen-reader name: the command, verb first. Explanation adds a sentence only when the name cannot say it: why the control is disabled, what it opens. On Windows kvit-ui adds the shortcut, "Bold (Ctrl+B)". Never repeat a visible label, never explain how to operate the control. (WRD-13, HLP-2)

- "Select every row the current filter shows" (`table.go:843`) → "Select shown rows"

### Symbols or words

**11. A symbol for a familiar action, words for the rest.** Bold, italic, search, delete, close, back, add: symbols, as IconButtons with a Label. An action people would have to guess from a symbol: words. (CMD-22)

**12. Never letters or punctuation standing in for a symbol.** A text button reading "B", "<>" or "+ Insert" is a symbol drawn in the wrong component. Use an IconButton with the meaning name and a Label in words. If the symbol you need is not in the catalogue, say so; do not type it. (CMD-22)

- "B", "I", "U", "S" (kvit-notes `app/toolbar.go:105-146`, `editor/formatbar.go:59-71`) → IconButtons named "Bold", "Italic", "Underline", "Strikethrough"
- "<>" → IconButton named "Code"; "H" → IconButton named "Heading"
- "+ Insert" → "Insert…" with the `plus` symbol before the words

**13. One symbol, one meaning, across the whole app.** Before using a symbol, check what else in the app uses it. The close symbol closes; it does not archive, withdraw, remove or cancel. (CON-2)

- X for archive (kvit-works `ui/shell/overviewcard.go:61`) → `archive`, "Archive"
- X for withdraw (kvit-works `ui/conversation/view.go:1242`) → `undo`, "Withdraw"
- X labelled "Close agent" that opens a browser (kvit-works `ui/shell/detailheader.go:184`) → `external`, "Open in browser"

### Commands

**14. One command, one name, everywhere.** A command has the same words in the menu bar, on its button, in its tooltip, in the tray menu and in any dialog it opens. Define it once as a `Command` and make every control from it. (CON-1)

- "Quick Capture…" (kvit-notes `app/tray.go:52`), "Quick capture note…" (`app/menus.go:51`), "Quick capture" (`app/capture.go:25`) → one command, "Quick capture…"

**15. Every command is in the menu bar.** Buttons, toolbars and context menus are shortcuts to commands that also live in a menu. A command that exists only on a button cannot be found from the keyboard or by a reader who missed the button. (CMD-1)

**16. Standard commands in the standard menus.** Undo, Redo, Cut, Copy, Paste, Select all and Find go in Edit; New, Open, Save, Close and Print in File; show and hide commands in View. Settings and About go where the platform puts them, which `MenuBar` does for you: Tools > Options and Help > About on Windows, the app menu on macOS. Exit on Windows, Quit on macOS. (CMD-2 to CMD-8)

**17. Menus: most used first, grouped, short, destructive last.** Separate groups with separators; at most one level of submenu, none in a context menu on Windows; a destructive item goes last, after a separator. A context menu holds only what applies to what was clicked. (CMD-10, CMD-12, CMD-13, CMD-23)

**18. A toggle says what it does.** On macOS a menu item swaps its words ("Show sidebar" / "Hide sidebar"); on Windows it keeps its name and shows a tick. Give MenuItem both words and let it choose. (WRD-12)

### Shortcuts

**19. Standard commands get the standard shortcut.** If the app has Save, Ctrl+S (Cmd-S) saves; Find, New, Close, Undo, Redo, Copy, Paste, Select all likewise. Use the standard `Command` and the key comes with it. A shortcut that does nothing when its command cannot act stays silent. (KEY-2)

- kvit-cash Save with no Ctrl+S → the standard Save command
- kvit-works without Ctrl+F, Ctrl+N, Ctrl+W → Find focuses the SearchField; New and Close as standard commands

**20. Menus show their shortcuts.** A menu item whose command has a key shows it; that is how people learn it. Context menus on macOS show none. (KEY-14, CMD-11)

- kvit-cash: none of 23 menu items shows a shortcut; kvit-notes: 4 of 109 → every item with a key shows it

**21. Few new shortcuts, and never the system's.** Give a new shortcut only to a frequent command, on a letter that recalls its name, built on the command modifier (`mod.OSMenuCommand()`, never `mod.Control` or `mod.Command` by hand). Never Alt+letter on Windows (access keys own it), never Ctrl+Alt, never a key the system or another command already uses. (KEY-3, KEY-4, KEY-5)

**22. Access keys on Windows and Linux.** Write labels with "&" before the access letter ("&Import"); kvit-ui underlines it on Windows and Linux and drops it on macOS. Not on OK, Cancel, Close, links or tabs. Two controls in one window never share a letter. (KEY-15)

### Dialogs and messages

**23. kvit-ui orders dialog buttons; do not.** Give Dialog its ConfirmText and CancelText and let it place them: confirming first on Windows and Linux, last on macOS. Never build a row of dialog buttons yourself. (WIN-1)

**24. Ask only when it matters.** No dialog before an action that can be undone: do it, and show ConfirmInPlace or a toast with Undo. Ask before a rare action that loses data for good, and then say what will be lost. (MSG-3)

**25. A message says what happened and what to do.** The title is the problem or the question, specific, never "Error". The detail adds the cause and the fix. An input error sits under its field and says how to put it right. No "Oops", no blame. (MSG-4, MSG-5, MSG-6, WRD-28)

**26. Explanations stay short.** A setting gets a label and at most one short line. A longer explanation goes behind a Hint. A paragraph that is always on screen is read once and then ignored. (WRD-22)

- kvit-notes `app/settings.go:269`, a permanent 50-word paragraph under one setting → one line, the rest in a Hint
