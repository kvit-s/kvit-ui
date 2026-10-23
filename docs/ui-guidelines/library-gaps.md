# What kvit-ui contradicts or lacks

What the library does today against the merged rules (rules.md), and the component changes and tests that would close the gap, in the order to do them. Line numbers are as of commit 0a36b77.

Of the 198 rules, kvit-ui enforces 26, partly covers 92, contradicts 10 and does not cover 70. Many of the 70 are app matters (help topics, printing, onboarding) or reviewer matters (tone, chart choice); the work below is the part the library can carry.

## Platform knowledge

The owner wants each platform's defaults. That needs kvit-ui to know at run time which form it is drawing: the Windows form (also used on Linux, Q-40) or the macOS form.

### What exists today

- `menu.go:100`: on macOS every menu is the system's (`runtime.GOOS == "darwin"`).
- `menu.go:21-24`, `menu.go:128`: the access-key marker is dropped on macOS.
- `fieldword.go:20`: Option deletes a word on macOS, Ctrl elsewhere.
- `wheel.go:41`: wheel easing differs on macOS.
- unison: `mod.OSMenuCommand()` is Command on macOS and Control elsewhere; `KeyBinding.String()` prints ⌃⌥⇧⌘ on macOS and "Ctrl+Shift+" elsewhere (unison `enums/mod/modifiers.go:90-144`); the standard Quit item is Quit on macOS and Exit elsewhere (unison `menu_bar_darwin.go:21-22`, `menu_bar_other.go:16-17`).

Nothing else varies: button order, labels, capitalisation, dialogs and menus are the same on every system.

### Proposal

Add `UI.Form` (`FormWindows` or `FormMac`), set from `runtime.GOOS` (Linux takes `FormWindows`) and settable by `uitest.Open`, so every screen can be rendered headless in both forms. Components read `ui.Form` at layout and draw time, as they read colours and sizes. Tests run each gallery page in both forms.

### Rules that need it

| What changes by platform | Rules | Windows and Linux | macOS |
|---|---|---|---|
| Dialog button order | WIN-1 | confirming button first (left), Cancel after it | Cancel, then the default button at the trailing end |
| Default button and Enter | KEY-6 (Q-7) | a default in every dialog, Cancel when destructive; Enter presses the focused button | Return presses the default; none when destructive |
| Close symbol on dialogs | WIN-6 (Q-8) | always | none on alerts |
| Capitalisation | WRD-6 (Q-3, Q-4) | sentence case except window titles | title case on menus, buttons, tabs, column headings |
| Command modifier | KEY-3, WRD-24 | Ctrl, written Ctrl+Shift+N | Command, written ⇧⌘N |
| Standard shortcuts | KEY-2 | Ctrl+Y redo, F1, F2, F3, F5, F11, Alt+F4 | Shift-Cmd-Z, Cmd-?, Ctrl-Cmd-F, Cmd-comma |
| Cancel key | KEY-7 | Escape | Escape and Command-period |
| Exit or Quit | CMD-3, CMD-4 | Exit, last in File | Quit <App>, last in the app menu |
| Settings and About | CMD-3, WIN-16 (Q-2) | Options last in Tools; About last in Help | Settings… and About in the app menu |
| Menu bar | CMD-1, CMD-2 (Q-1) | Kvit-drawn in the window; File, Edit, View, Tools, Help | the system's; App, File, Edit, View, Window, Help |
| Window menu | CMD-7 | none | required |
| Access keys | KEY-15 (Q-11) | Alt+letter on every control and menu | none |
| Ellipsis | WRD-4 | none on Options, Properties, About | on Settings… |
| OK, Yes, No, Done, Close | WRD-7, WRD-8, WRD-21 (Q-5, Q-6) | OK for settings windows, Yes/No pairs, Close ends flows | specific verbs, Done ends flows |
| Toggle commands | WRD-12 (Q-24) | name kept, ticked | Show/Hide words swap |
| Shortcut in tooltips | WRD-13, KEY-14 (Q-23) | Bold (Ctrl+B) | no shortcut |
| Shortcuts in context menus | CMD-11 | shown | not shown |
| Context-menu default | CMD-10 (Q-33) | first and bold | none |
| Exclusive menu options | CMD-15 | bullet | tick |
| Colon after labels | WRD-19, WRD-20 (Q-22) | Name: | none |
| Tab focus | KEY-10 (Q-16) | every control | text fields and lists, unless full keyboard access |
| Window title | WIN-12 | "document - program" | the document only |
| Unsaved mark | WIN-20 (Q-13) | asterisk in the title | dot in the close button |
| Full screen | WIN-15 | F11, View > Full screen | the system's full screen, Ctrl-Cmd-F |
| Placeholder | WRD-19 | hidden on focus | kept until typing |
| Radio buttons before a select | CTL-2 | up to seven | up to five |
| Tabs in a row | CTL-10 | up to seven | up to six |
| Mixed checkbox click | CTL-5 (Q-27) | cycles back to the mix | not defined |
| Disabled Next | CTL-8 (Q-30) | stays enabled, explains on press | may be unavailable |
| Window attention | MSG-1, WIN-25 | flash the taskbar button | bounce the Dock icon |
| Notifications | MSG-10, MSG-11, MSG-13 | no buttons, sentence-case title | title case, buttons allowed |

### Capitalisation

Four ways a label can get both forms (Q-3 has the full argument):

1. **The app writes sentence case and the component converts to title case on macOS.** One string per label. Cost: a converter (Apple's rule: every word capitalised except articles, coordinating conjunctions and prepositions of three letters or fewer; first and last always) and one call in Button, IconButton's tooltip, MenuItem, Tab, Segmented, table column titles and Dialog's buttons. Fails on: proper nouns written lower case on purpose (iCloud, macOS), words the rule cannot classify (a phrasal verb's particle, "Set Up" vs "setup"), and any label that contains the reader's data. Mitigation: an exception list in kvit-ui, and data passed apart from the label (WRD-29 already forbids building labels from pieces).
2. **The app supplies both forms.** Twice the strings and translations, and nothing keeps them saying the same thing.
3. **The app writes title case and the component lowers it on Windows.** Lowering breaks proper nouns and acronyms (Wi-Fi, PDF, Kvit Notes) unless the app marks them, which is more work than option 1's exceptions.
4. **One form everywhere.** Breaks the owner's platform decision.

Recommended: 1. A test renders each gallery page and each app screen in both forms and checks every name against the converter, so a wrong conversion is seen once and goes on the exception list.

## Where kvit-ui contradicts a rule

| Rule | What kvit-ui does | Change |
|---|---|---|
| WIN-1 | The confirming button is on the right on every system (`feedback.go:371`, laid out right to left at `feedback.go:541-552`). | Dialog lays out by `ui.Form`: confirm first on Windows and Linux, last on macOS. Test: in each form, the confirming button's x against Cancel's. |
| KEY-6 | Buttons press on Space only (`button.go:75`); Dialog has no default button and only focuses Confirm, or Cancel when destructive (`feedback.go:470-479`); `Dialog.Accept` is left to the app (`dialogaccept.go:3-7`). | Button and IconButton also press on Return and NumPadEnter when focused on the Windows form. Dialog gets a default: Return presses it on macOS (none when Destructive); on Windows Enter presses the focused button. A `Default` mark drawn on the button. |
| WRD-7 | Dialog's ConfirmText defaults to "Confirm" (`feedback.go:424`). | No default: an empty ConfirmText logs a warning unless the dialog only informs (then a single Close, or OK on macOS, Q-5). |
| WRD-3 | Library defaults carry articles: "Clear the filter" (`searchfield.go:56`), "Close the pane" (`pane.go:44`), "Dismiss this notice" (`feedback.go:261`), "Clear the sort" (`table.go:597`), "Move this column left" and the rest of the column menu (`table.go:595-605`), "Select the rows shown" (`table.go:965`), "Show the selected column" and three more (`duallist.go:54-57`). | Rewrite: Clear, Close, Dismiss, Clear sort, Move left, Move right, Widen, Narrow, Hide column, Show all columns, Select all, Show, Hide, Move up, Move down. |
| WRD-25 | Button.BusyText defaults to "Working…" (`button.go:73`). | BusyText required with Busy; a busy button without it logs a warning. |
| FDB-10 | Progress shows "N%" beside the bar by default (`progress.go:45`, `ShowPercent: true`). | ShowPercent defaults to false; the label says what is happening (WRD-25). |
| CMD-11 | The field's context menu shows each action's shortcut (`field.go:211`). | ShowMenu drops `Key` for a context menu on macOS. |
| WIN-8 | Every dismissal is a decline (`feedback.go:498-504`); typed input in a dialog is lost without a word. | Dialog gains `OnDiscard func() bool`, asked before a reject by the close symbol or Escape when the content changed. |
| MSG-10 | `Tray.Notify` posts even while a Kvit window is active (`platform/tray.go:272-283`); the Windows balloon has no NIIF_NOSOUND (`platform/tray_windows.go:438-439`), so it plays a sound. | Notify skips when a Kvit window is the active window (or returns false and lets the app decide); add NIIF_NOSOUND. |
| LAY-3 | The window minimum is Px(880) by Px(600) (`tokens/interface.go:195`, `tokens/interface.go:203`, `window.go:88-92`), above Windows' 800 × 600 floor. | Decide Q-25 first. |

## What kvit-ui lacks, in the order to build it

Each step lists the rules it serves. Steps 1 to 3 are foundations the rest stand on.

### 1. The platform form

`UI.Form` and `uitest.Open` with a form, as above. Rules: every rule with a Windows and a macOS form (91).

### 2. Commands and shortcuts

A `Command` type: name, symbol, shortcut (declared with `mod.OSMenuCommand()`), enabled func, run func, and an optional explanation. MenuItem, IconButton, Button and a toolbar item are made from a command, so one command has one name and one symbol everywhere (CON-1, CON-2), its shortcut shows in its menu item and tooltip (KEY-14, WRD-13), and it is disabled in every place at once (CMD-9). The Window keeps a shortcut registry: a command's key runs it wherever the window has focus (KEY-2), and the registry refuses two commands on one key (KEY-5) and keys the system owns (KEY-4).

A table of standard commands with their keys per form (New, Open, Save, Close, Print, Undo, Redo, Cut, Copy, Paste, Delete, Select all, Find, Settings or Options, Help, Full screen, Rename, Refresh): an app that has the command gets the standard key without choosing one (KEY-2). This is the fix for kvit-cash's missing Ctrl+S and kvit-works' missing Ctrl+F, Ctrl+N and Ctrl+W.

Today: MenuItem.Key is display only (`menu.go:26-28`, survey §4); apps catch keys in `Window.OnKeyDown` (`window.go:37-41`); `docs/plugins.md` proposes command routing but none exists.

### 3. The menu bar

`kvitui.MenuBar`, built from the app's commands and the standard menus per form (CMD-1 to CMD-8): the system's bar through unison on macOS (with kvit-ui's own "Settings…" item, since unison's says "Preferences…", unison `menu_bar.go:130-136`), and a Kvit-drawn bar in the window on Windows and Linux, with access keys, Alt to focus it, and the hide-and-Alt route if Q-1 allows it. unison's per-window bar puts Exit in File and Preferences in Edit with Ctrl+Q (unison `menu_bar.go:84-136`), which is not the Windows form, so the Windows bar is kvit-ui's.

Today: no menu bar on any system (survey §4).

### 4. Dialog

- Button order by form (WIN-1), default and Enter (KEY-6), no "Confirm" default (WRD-7), Close-only informing dialogs, Cancel turning into Close or Stop (WRD-8).
- Focus back to the opener on close (KEY-8; today `feedback.go:482-487` only hides; Popover already does it, `popover.go:125-148`).
- Focus starts in the first field when the body has one, not on the confirm button (KEY-8; `feedback.go:470-479`).
- Tab kept inside an open modal (survey §2: no focus trap).
- OnDiscard (WIN-8); a refused Destructive dialog with no Cancel (WIN-3).
- No close symbol on macOS if Q-8 says so; Command-period cancels on macOS (KEY-7).
- One modal at a time except an alert over a dialog (WIN-7; `layer.go:58-68` stacks without a check).
- A don't-show-again option for informing dialogs (WIN-2).

### 5. Wording in the components

- A title-case converter for the macOS form, called by every component that shows a label (WRD-6, Q-3).
- `Asks bool` on Button, IconButton and MenuItem, the component adding "…" per form (WRD-4).
- The library's own strings rewritten without articles (WRD-3, list above), and TypeAhead's count through `UI.CountPhrase` (OTH-2; `typeahead.go:100-102`).
- SearchField's placeholder set per use, "Filter…" kept only as a fallback that a test flags on app screens (CTL-15; `searchfield.go:38`).
- BusyText required with Busy (WRD-25); ShowPercent off by default (FDB-10).
- MenuItem toggle words per form (WRD-12): `On`/`Off` words on macOS, a tick on Windows.

### 6. Keyboard

- Return and NumPadEnter press a focused Button and IconButton on the Windows form (KEY-6; `button.go:75`, `iconbutton.go:59`).
- Access keys on every eligible control, unique per window, underlined while Alt is held (KEY-15, Q-11). Today only menus have them (`menu.go:525-564`).
- Roving focus with arrow keys in Tab, Segmented and the sidebar, as RadioGroup has (KEY-9; `tab.go:42`, `segmented.go:56-60`, `sidebar.go:120-124`).
- Ctrl+Tab and Ctrl+Shift+Tab between tabs (KEY-9).
- Type-select in Table, Tree and lists (KEY-9).
- SearchField registers as the window's Find target for Ctrl+F (and Ctrl+E on Windows, Q-20) (KEY-13).
- Escape closes TypeAhead's list (KEY-7; `typeahead.go:171`), and abandons a column drag, a split drag and a dialog move (`table.go:870-899`, `dialogmove.go:74-119`).
- Space pages a Region (KEY-11; `region.go:133-144`).
- The Mac's full-keyboard-access setting read through `platform` (KEY-10, Q-16).

### 7. Menus

- A separator inserted before a trailing Danger item by ShowMenu itself (CMD-23; today the caller's job, `menu.go:35-39`).
- One level of submenu; none in context menus on the Windows form (CMD-12; `menu.go:50-54` nests without limit).
- A default item (first, bold) for context menus on the Windows form, and Checked drawn with the tick only (CMD-10, Q-33; `menu.go:236-239`).
- A radio (bullet) item for exclusive options (CMD-15; Select uses Checked, `select.go:96`).
- Tray items with a check mark and a default item (WRD-12, CMD-24; `platform/tray.go:49-60`).
- A split button and a menu button component (CMD-17).

### 8. Tests

Each reads what the rule needs:

- **`uitest.CheckLabels(window, form)`, one screen's controls.** Every Button, IconButton, MenuItem, Tab and column name: starts with an allowed verb or a standard verbless label (WRD-1); at most three words, four for a Windows toolbar button (WRD-2); no a, an, the, my, your (WRD-3); not OK, Yes, No or Confirm outside the allowed cases (WRD-7); no ending full stop (WRD-27); in the form's case (WRD-6); a Button whose words are one or two characters or only punctuation fails (CMD-22: the "B", "I", "U", "<>" buttons). Tooltips: at most 60 characters on the macOS form, no full stop (WRD-13). A Dialog has a Title, at most three buttons and a Cancel when Destructive (WIN-2, WIN-3, WIN-6). At most one ButtonPrimary per region (CTL-1). No static label longer than about 20 words under a control (WRD-22: the 50-word paragraph). Every elided Label has a tooltip with its text (PTR-12).
- **`uitest.CheckApp(screens)`, every screen of an app together.** Collects (command handler, name, symbol, shortcut) from every screen and fails one command with two names (CON-1: Quick Capture… / Quick capture note… / Quick capture), one symbol with two unrelated names (CON-2: the X for archive, withdraw, remove, cancel and close), one key for two commands (KEY-5), a standard command without its standard key (KEY-2), a command reachable from a button or context menu but missing from the menu bar (CMD-1), and a menu item with a shortcut that its menu does not show (KEY-14).
- **`uitest.CheckRules` extended, the Go source.** Flags a label literal ending in "..." or "…" once `Asks` exists (WRD-4); "Oops", "tap", "click here" in strings (WRD-15, WRD-17, WRD-18); `mod.Command` or `mod.Control` in a KeyBinding literal (KEY-3); a Button text starting "+ " (CMD-22).
- **Both forms.** Every gallery page and app screen rendered with `FormWindows` and `FormMac`, and the checks above run in each.

### 9. Windows and state

- `Window.AllowClose` wired to unison's `AllowCloseCallback`, and Quit or Exit through it (WIN-19).
- `Window.Edited` for the unsaved mark (WIN-20).
- Window frame saved and restored in settings, fitted to the screen's work area (WIN-14, WIN-25; `window.go:117` sets Px(1440)×Px(960) without fitting).
- Full screen (WIN-15).
- A window title helper giving "document - program" on Windows and the document on macOS (WIN-12).

### 10. Controls and the rest

- A form row: a visible label tied to its field by LabeledBy, with the colon per form (WRD-19, LAY-6; `field.go:34-35`).
- Label `Selectable` for content text (CTL-18), and a tooltip with the full text when elided (PTR-12).
- Symbols for formatting: bold, italic, underline, strikethrough, code, heading (CMD-22; the catalog's symbol list has none of them, though it has `link` and `plus`). kvit-notes' "B", "I", "U", "S", "<>" and "H" buttons need these before they can become IconButtons.
- Check cycling back to the mix (CTL-5, Q-27; `check.go:49-51`).
- Progress with Cancel or Stop and a time-remaining line (FDB-3, FDB-10).
- A Toolbar component with leading, centre and trailing groups and overflow (CMD-18, CMD-19).
- Notifications: an urgency, skipped while a Kvit window is active, silent on Windows (MSG-10, MSG-12).
- An UndoStack with named entries for the Edit menu, ConfirmInPlace and Toast (FDB-7, WRD-11).
- A help link to a topic (HLP-6, Q-34).

## Already in place

The library already holds a good part of the behaviour the sources ask for: one name for a symbol-only button's tooltip and its screen-reader name (`iconbutton.go:23-27`), tooltips on keyboard focus as well as hover (`control.go:164-166`), Escape and an outside click closing popovers and menus, the menu's keyboard, access keys inside menus, the Menu key and Shift+F10 for context menus, a destructive dialog that takes no keyboard default, Explanation on disabled buttons, ConfirmInPlace and an undo in a toast instead of confirmation dialogs, Notice and Toast separated by whether the condition still holds, EmptyState, counts with nouns, a second channel for every colour distinction, contrast floors and interface-size scaling. rules.md marks these "enforces".
