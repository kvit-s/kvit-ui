# Merged rules

The rules Kvit apps follow for wording, commands, keyboard, pointer, windows and dialogs, messages, controls, feedback, layout and consistency, merged from Apple's Human Interface Guidelines, Microsoft's Windows 7 UX guide and Microsoft's current Fluent and WinUI guidance. README.md says where they come from and how they were made.

How to read a rule:

- **Windows** and **macOS** give the two forms where the platforms differ. The app uses the form for the system it runs on; Linux uses the Windows form (the owner's decision, open-questions.md Q-40). Where no forms are given, the rule is the same everywhere.
- **Sources** lists the extracted records (ids in extracted.jsonl, `apple-<page>-<n>` or `win7-<page>-<n>`), grouped by page, each number the record's counter on that page (the record's url in extracted.jsonl links its section). A `*` marks a record on Microsoft's own checklist (top-violations).
- **Enforce** says where the rule should be held: a component, a test reading the Go source, a test reading one running screen's controls with their roles and names, a test reading every screen of an app together, or a reviewer.
- **kvit-ui now** says whether the library enforces the rule today, partly does, contradicts it, or does not cover it, with the files that do it (line numbers in rules not yet changed are as of commit 0a36b77).

Ids are stable. WRD-5 was folded into WRD-6 and WIN-21 into WIN-11; their numbers are not reused. Where the Windows 7 guide and the current guidance disagree on behaviour, the Windows form follows the current guidance.

199 rules. kvit-ui enforces 73, partly covers 97, contradicts none and does not cover 29. 94 have a Windows and a macOS form.

## Contents

- [Wording](#wording) (WRD, 28)
- [Commands: menus, toolbars, buttons](#commands-menus-toolbars-buttons) (CMD, 25)
- [Keyboard](#keyboard) (KEY, 15)
- [Pointer](#pointer) (PTR, 15)
- [Windows and dialogs](#windows-and-dialogs) (WIN, 24)
- [Messages](#messages) (MSG, 17)
- [Controls](#controls) (CTL, 26)
- [Feedback](#feedback) (FDB, 10)
- [Layout](#layout) (LAY, 10)
- [Accessibility](#accessibility) (A11Y, 8)
- [Consistency](#consistency) (CON, 8)
- [Help](#help) (HLP, 8)
- [Other](#other) (OTH, 5)

## Wording

### WRD-1

A button or menu item is a command: it starts with a specific verb that names what happens (connect, install; not configure, manage or set), named for what the reader wants done, not the technology. Its words tell it apart from its neighbours: paired opposites (Add and Remove, Show and Hide), a modifier for similar commands (Insert row above, Insert row below), the object it acts on (Open Windows Defender), singular or plural by how many things it acts on.

- **Windows:** Standard labels without a verb are accepted: OK, Yes, No, About, Advanced, Back, Details, Forward, Full screen, Less, More, New, Next, Personalize, Previous, Properties, Settings. A noun stands alone where the verb would be Create, Show, View or Manage (Tools, Page, Status bar). A verb that repeats the menu title is dropped (Insert > Table). When every answer in a set starts with the same words, they are dropped from all. Background work is Suspend/Resume, anything else Enable/Disable.
- **macOS:** The verb stays where Windows drops it (Show Toolbar, Enter Full Screen); verbless labels only for the standard noun commands (Settings…, About <App>, Bold).
- **Sources** (32 records):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 18; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 8; [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 1; [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 11; [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 13
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 201, 191, 135; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 77; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 80, 83, 84, 85, 86, 87; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 37; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 84, 85, 86, 87, 88, 89, 90; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 30; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 74, 75; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 50; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 36, 37; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 39
  - Fluent: [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 14; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 26
- **Enforce:**
  - Test, one screen: Every Button, IconButton and MenuItem name starts with a word from an allowed verb list or the platform form's verbless list (above); the lists live in uitest and grow by review.
  - Test, all screens of an app: Across all screens, two commands with the same verb and object but different results are flagged (Insert row above vs Insert row).
  - Review: The verb is the specific one for the result.
- **kvit-ui now:** partly. `uitest.Labels` fails a Button, IconButton or menu line that does not start with a word from `uitest.Verbs` or `uitest.Verbless`, standard commands and switches excepted (`uitest/labels.go`); whether the verb is the specific one is the reviewer's.

### WRD-2

Keep a command label to one to three words, using only the words the command needs, and button text within about 26 characters; a word more is better than a separate explanation. When specific answers cannot be that short, offer them as radio buttons, each with one line of explanation, committed by one button with a specific verb (WIN-2), never as command links (decided, Q-10, Q-51); Yes and No only where no specific answer is possible (WRD-7).

- **Windows:** A toolbar button's label is preferably one word.
- **Sources** (12 records):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 18; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 4; [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 18; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 7
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 114; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 75; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 63; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 38
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 13; [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 11; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 17; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 25
- **Enforce:**
  - Test, one screen: Fail a Button or MenuItem label over three words and an IconButton Label over three words; fail a Button label over 26 characters; allow a list of exceptions with a reason.
  - Review: A label that fits but is padded; a toolbar label longer than one word on the Windows form.
- **kvit-ui now:** enforces it. `uitest.Labels` fails a command over three words and a Button over 26 characters (`uitest/labels.go`); the library's own labels keep within it ("Move up", "Select shown rows").

### WRD-3

Leave articles and possessives off labels, menu items and tooltips.

- **Sources** (3 records):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 4; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 13; [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 22
- **Enforce:**
  - Test, one screen: Fail a command label that contains a, an, the, my or your as a whole word.
- **kvit-ui now:** enforces it. The library's own labels carry no articles (Clear, Close, Dismiss, Clear sort, Hide column, Select shown rows, Move up; `searchfield.go`, `pane.go`, `feedback.go`, `table.go`, `duallist.go`), and `uitest.Labels` fails a command with a, an, the, my or your (`uitest/labels.go`).

### WRD-4

A command that needs more input before it acts (a dialog that asks, a chooser, a confirmation) ends with an ellipsis (…); one that acts at once, or opens a submenu, has none.

- **Windows:** A command whose whole job is to show a window has no ellipsis, even if that window holds settings: About, Advanced, Help, Options, Properties, Settings, Sharing settings. A command followed only by a confirmation has one. An icon-only toolbar button that asks for more has the ellipsis in its tooltip (Print…); toolbar commands should mostly act at once. Help links never end in an ellipsis.
- **macOS:** Ellipsis character U+2026, not three full stops. Settings… and other commands that open a window to collect input take it.
- **Sources** (22 records, 9 on Microsoft's checklist):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 18; [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 7; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 18
  - Windows 7: [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 76\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 21\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 70\*, 71, 72, 89; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 9\*, 8\*, 45, 46; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 40; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 27\*; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 37\*; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 77\*, 78\*, 92, 79; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 25; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 36
- **Enforce:**
  - Component: MenuItem, Button and IconButton gain `Asks bool`; the component appends "…" (U+2026) on both forms, and to an IconButton's tooltip on the Windows form. A second flag, `OpensWindow`, appends it on the macOS form only (Settings…). A MenuItem with a submenu never gets one.
  - Test, Go source: Fail a label literal ending in three full stops or in "…" once `Asks` exists.
  - Test, one screen: Fail a MenuItem with a submenu whose name ends in "…".
  - Review: The command really asks for more.
- **kvit-ui now:** enforces it. `Asks` on Command, Button, IconButton and MenuItem adds "…", and Command's `OpensWindow` adds it on the macOS form only (`command.go`, `button.go`, `iconbutton.go`, `menu.go`); `uitest.CheckRules` flags three full stops and an ellipsis typed into a command's words (`uitest/wording.go`), and `uitest.Labels` a submenu line ending in "…".

### WRD-6

Capitalise each kind of text in the platform's style, and give every label of one kind the same style across the app. The app writes every label once, in sentence case, and kvit-ui converts it to title case where the macOS form asks for it, with an exception list; labels never contain user data (decided, Q-3). Capitalise feature names sparingly and consistently; generic words (toolbar, menu, button, icon) stay lower case.

- **Windows:** Sentence case for all text, titles included: window and dialog titles, a guided flow's name, menu titles and items, buttons, tabs, column headings, field, check box, radio, list and slider labels, list options and data, prompts, links, tooltips, toolbar labels, main instructions, supplemental text and progress text. Title case nowhere. Titles have no ending punctuation.
- **macOS:** Title case: menu items and menu titles, buttons, tabs, segment labels, column headings, panel titles, notification titles and buttons. Sentence case: tooltips, alert body text, box titles, field and slider labels, notification body. Alert buttons are title case (decided, Q-4).
- **Sources** (79 records, 35 on Microsoft's checklist):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 3; [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 13; [tab-views](https://developer.apple.com/design/human-interface-guidelines/tab-views) 4; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 8; [segmented-controls](https://developer.apple.com/design/human-interface-guidelines/segmented-controls) 8; [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 6; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 3; [panels](https://developer.apple.com/design/human-interface-guidelines/panels) 5; [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 8, 16, 9, 11; [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 21, 12; [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 24; [boxes](https://developer.apple.com/design/human-interface-guidelines/boxes) 4; [sliders](https://developer.apple.com/design/human-interface-guidelines/sliders) 7; [combo-boxes](https://developer.apple.com/design/human-interface-guidelines/combo-boxes) 3; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 11
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 84, 186\*, 192, 199, 202\*; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 30\*, 35, 38; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 30; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 46\*; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 42\*, 47\*, 50\*; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 29\*; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 67\*, 73\*, 85; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 45\*, 51\*, 55; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 39\*, 49\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 72\*, 78\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 79\*, 84\*, 89; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 16\*, 17\*, 19\*, 18\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 58, 73\*; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 39\*; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 33; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 32\*, 53, 33\*, 35\*, 58; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 30\*, 42, 37; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 28; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 81\*; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 33; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 25; [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 19; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 32; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 10\*; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 49\*, 54\*; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 28\*; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 41\*, 44\*; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 43; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 25, 30
  - Fluent: [signature-experiences-typography](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/typography) 4; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 19
- **Enforce:**
  - Component: Labels pass through one function that gives the platform form (library-gaps.md, Platform knowledge). Button, IconButton, MenuItem, Tab, Segmented, table column titles, Dialog buttons and window titles call it: sentence case on the Windows form, title case where the macOS form asks for it, skipping the words on its exception list.
  - Test, one screen: On a headless run with each platform form, compare every name with the converted form.
  - Test, all screens of an app: Every label of one kind (menu item, button, tab, column heading, tooltip) across all screens of an app follows the form's style.
  - Test, Go source: Fail a label built from user data (a label string assembled with a variable) passed to a component that converts case.
- **kvit-ui now:** partly. `TitleCase` converts a sentence-case label on the macOS form, with `TitleCaseKept` for exceptions (`wording.go`); Button, MenuItem, menu titles, Tab, Segmented, table column headings and Dialog's buttons go through it, and `uitest.Labels` compares every such label with its converted form in each form (`uitest/labels.go`). Nothing flags a label built from the reader's data.

### WRD-7

A confirming button names the action it confirms (Insert, Import), never a generic word such as Confirm; what people must notice about the action goes on the button itself (Format disk, not OK under a paragraph about formatting). When the dialog gives a reason not to go ahead, the button adds anyway (Install anyway). Now and later appear only to tell doing it at once from doing it later.

- **Windows:** Specific answers are preferred. A message that only needs dismissing has a single Close button, never OK. OK, Cancel and Apply appear only on a properties dialog of one object, spelled OK (never Ok or Okay). Yes and No only where no specific answer is possible, always as a pair, never replaced by OK and Cancel. When an error ends the program the button is Exit program.
- **macOS:** Never OK, Yes or No on a confirming button; an informing alert has a single OK (decided, Q-5).
- **Sources** (42 records, 14 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 22, 19; [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 17
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 94, 103, 104, 106, 107, 112, 113\*, 118, 33, 102\*, 111, 115, 117, 119\*, 120; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 39; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 86, 87\*; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 26\*, 54, 25\*, 27; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 48; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 36\*, 118\*, 112\*, 34\*; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 35\*, 44\*, 40, 41; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 76, 77\*, 8, 64\*
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 12, 30; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 24; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 11
- **Enforce:**
  - Component: Dialog drops its "Confirm" default and requires ConfirmText. Dialog gains `Kind` (Ask, Inform, Properties): Inform shows one button, Close on the Windows form and OK on the macOS form; Properties allows OK, Cancel and Apply; `YesNo` is allowed only on the Windows form and only with a Title ending in "?".
  - Test, one screen: Per form: fail a dialog button named Confirm, Ok or Okay; on macOS fail OK, Yes or No except a one-button Inform OK; on Windows fail OK outside a Properties dialog, Yes without No, Yes beside OK.
  - Review: The verb is specific; anyway where a reason not to go ahead is given; Yes and No only where no specific answer exists.
- **kvit-ui now:** enforces it. Dialog has no "Confirm" default and warns on a dialog that asks without ConfirmText; `DialogInform` (`NewMessage`) has one button, Close on the Windows form and OK on the macOS form, and `DialogProperties` OK, Cancel and Apply (`dialog.go`); `uitest.Labels` fails Confirm, Ok, Okay, and OK, Yes or No where the form does not allow them (`uitest/labels.go`).

### WRD-8

The button that abandons the action is Cancel, and Cancel leaves no trace. Once something cannot be undone it becomes Close (the work is done, or the window only shows things) or Stop (work is running and what is done is kept); Escape presses it.

- **Windows:** Cancel turns into Close once a step has run or applied changes that cannot be undone (and back to Cancel after Apply), and into Stop once going back is no longer possible mid-operation. A modeless dialog has Close, never Cancel. A finished flow, sheet or popover ends with Close (decided, Q-6). Yes never pairs with Cancel alone (Yes/No/Cancel is allowed). When several pending actions make it unclear what Cancel would cancel, or only the dialog itself can be declined, a specific negative answer replaces it (Don't run, Not now).
- **macOS:** Cancel is never renamed; a specific negative answer (Don't Save) sits beside Cancel, not in its place. A finished flow, sheet or popover ends with Done (decided, Q-6).
- **Sources** (23 records, 13 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 20, 30; [sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) 1
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 109, 110\*, 124, 130, 126, 108\*, 123; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 33\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 35\*, 38\*; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 36\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 39\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 28\*, 29, 50\*; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 78\*, 79\*, 80\*; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 17\*
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 20
- **Enforce:**
  - Component: Dialog keeps "Cancel" as its default and gains `Done()`, which turns Cancel into Close on the Windows form and Done on the macOS form; a progress dialog gets a halt button whose words follow `Stoppable` (Cancel, then Stop on the Windows form).
  - Test, one screen: Fail a dialog whose abandoning button is not Cancel, Close or Stop, except a `Don't …` answer on the Windows form; fail Yes beside Cancel without No.
  - Review: Cancel really leaves no trace.
- **kvit-ui now:** enforces it. CancelText is Cancel unless set, and `Dialog.Done` turns it into Close on the Windows form and Done on the macOS form (`dialog.go`); Progress's halt button takes the words the app gives, Cancel or Stop (`progress.go`); `uitest.Labels` fails another abandoning word and Yes beside Cancel without No.

### WRD-9

A New command names the kind of thing it makes.

- **Sources** (1 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 17
- **Enforce:**
  - Review: "New" alone where the app makes more than one kind of thing.
- **kvit-ui now:** does not cover it. Not covered.

### WRD-10

Delete removes without touching the clipboard and matches the Delete key; Cut copies then removes. Do not call removal Erase or Clear, and do not call a move to the trash Delete without saying so.

- **Windows:** A move to the Recycle Bin (or trash) is Delete…, and its confirmation or undo says where the item went (decided, Q-14).
- **macOS:** A move to the trash is Move to Trash… (decided, Q-14).
- **Sources** (3 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 34, 35; [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 10
- **Enforce:**
  - Review: The label says what really happens to the item.
- **kvit-ui now:** does not cover it. Not covered.

### WRD-11

Undo and Redo name the action they reverse.

- **Sources** (2 records):
  - Apple: [undo-and-redo](https://developer.apple.com/design/human-interface-guidelines/undo-and-redo) 1; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 33
- **Enforce:**
  - Component: An undo stack in kvit-ui (FDB-8) supplies the names; the Edit menu reads them.
- **kvit-ui now:** enforces it. `UndoStack` names Undo and Redo for the edit they reverse, "Undo delete transaction", and feeds the Edit menu, a Toast and a ConfirmInPlace (`undo.go`).

### WRD-12

A command that switches a state says unmistakably what it does in each state.

- **Windows:** A menu item keeps its name whatever the state and shows the state with a tick (on or off) or a bullet (one of several modes): "Status bar" with a tick, not "Show status bar" turning into "Hide status bar". Items named after objects (recent files, open windows) may change. A tray menu item stays enabled and shows its state with a tick.
- **macOS:** The words change with the state: Show Sidebar / Hide Sidebar; a tick only for attributes (CMD-15).
- **Sources** (7 records):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 21, 22; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 41; [panels](https://developer.apple.com/design/human-interface-guidelines/panels) 9
  - Windows 7: [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 38; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 8; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 30
- **Enforce:**
  - Component: MenuItem gains `Toggle{Name, On, Off string}`: the Windows form shows Name with `Checked`; the macOS form (the system menu) shows On or Off. `platform.TrayItem` gains `Checked`.
  - Review: A label that could be read as a status.
- **kvit-ui now:** enforces it. Command and MenuItem take `TurnOn` and `TurnOff`: the Windows form keeps the name and ticks it, the macOS form shows the words for the state (`command.go`, `menu.go`); a tray menu line takes `Checked`, a tick with its words kept (`platform/tray.go`).

### WRD-13

A tooltip is short, concerns only its control, does not repeat a visible label, and never says how to operate the control (right-click to…). Each control has one tooltip string: its verb-led name, plus an optional Explanation line (decided, Q-23).

- **Windows:** A tooltip on an unlabelled control is its label: the command name in sentence case, plus the shortcut and any default in parentheses (Bold (Ctrl+B)), and an ellipsis when the command asks for more. A descriptive tip on an item (infotip) is complete sentences in the present tense, with full stops, in short parts. A notification-area icon's tip gives the program's name and, if useful, a short status, without version numbers. The Start menu description is at most 25 words, starts with a specific verb, lists the main tasks in the words people search for, without catch-alls or marketing, in sentence case.
- **macOS:** A fragment starting with a verb that says what the control does: no articles, no full stop, sentence case, within about 60 characters.
- **Sources** (38 records, 4 on Microsoft's checklist):
  - Apple: [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 19, 20, 21, 22, 24, 25, 26, 18
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 12, 29, 31, 33, 34, 36, 38, 39, 40, 41, 42, 43, 44, 45, 23, 24, 26, 27\*; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 16; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 74\*, 77\*; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 53, 54, 55, 56, 57; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 50\*; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 26
  - Fluent: [ui-controls-tooltips](https://learn.microsoft.com/en-us/windows/apps/design/controls/tooltips) 11, 15
- **Enforce:**
  - Component: IconButton's tip already starts with its Label (`iconbutton.go:69-74`); on the Windows form it appends the command's shortcut in parentheses, which needs a `Key` on IconButton (or one command table that menu and button share).
  - Test, one screen: Windows form: an IconButton tooltip starts with its name and carries the shortcut when the same command has one in a menu. macOS form: fail a tooltip over 60 characters or ending in a full stop. Both: fail a Button Explanation equal to the button's words.
  - Review: The tooltip describes this control only.
- **kvit-ui now:** enforces it. A control made from a Command adds its shortcut in parentheses to its tooltip on the Windows form, "Bold (Ctrl+B)", and none on the macOS form (`command.go`, `iconbutton.go`, `button.go`); `uitest.Labels` fails a macOS-form tooltip over 60 characters or ending in a full stop, and an Explanation repeating a button's words (`uitest/labels.go`).

### WRD-14

Speak to the reader as you; never call them the user; say you can, not the app lets you. The app never speaks of itself as I, me or we. Write a problem message in a form that suits every system, leaving out the subject where the active voice would blame (Couldn't save the file), and never make the reader the culprit (decided, Q-15).

- **Windows:** An option the reader picks may speak as the reader in the first person, and is then exempt from WRD-3 (Print the photos on my camera).
- **Sources** (17 records, 1 on Microsoft's checklist):
  - Apple: [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 3, 4; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 14
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 181, 180; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 57\*; [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 23, 2, 22, 21, 24; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 70; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 74
  - Fluent: [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 2, 9, 11, 10
- **Enforce:**
  - Test, Go source: Flag the words user, our, I, me, we and us in string literals passed to labels, Detail and Text fields; flag lets you and allows you.
  - Review: Tone; blame in problem messages.
- **kvit-ui now:** does not cover it. Not covered.

### WRD-15

Use plain words in the reader's terms: what they want to do, not the technology behind it. One clear word where one will do (lock, not lock down; because, not since), short product names without trademarks or versions, accurate statements; no jargon or abbreviation without a definition where the reader first meets it, no abbreviations that look alike or in text that will be translated, no idioms, slang, jokes or exclamations such as Oops. Short fragments are fine where they read naturally.

- **Sources** (49 records, 8 on Microsoft's checklist):
  - Apple: [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 2, 5, 6, 7; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 19, 5, 3, 1
  - Windows 7: [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 20, 1, 4, 5, 12, 13, 14, 15, 16, 32, 34, 35, 8\*, 9, 10; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 23; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 15\*, 17\*; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 62, 64, 55\*, 65; [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 17\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 60\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 85\*, 13\*; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 52, 51; [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 44; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 195; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 37
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 8; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 1, 4, 8, 12, 13, 14, 21, 28, 30
- **Enforce:**
  - Test, Go source: Flag "Oops", "Whoops", "!" at the end of a label or title.
  - Review: Jargon, technology words and tone.
- **kvit-ui now:** partly. `uitest.CheckRules` flags "Oops" in a string (`uitest/wording.go`); "Whoops" and a label ending in "!" are not flagged.

### WRD-16

Write without gender where it is irrelevant, never use a disability as a metaphor, and never word text as if every reader can see, hear or use a mouse.

- **Sources** (5 records):
  - Apple: [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 10, 19, 20
  - Windows 7: [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 64
  - Fluent: [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 29
- **Enforce:**
  - Review: Pronouns, metaphors and assumed abilities in interface text.
- **kvit-ui now:** does not cover it. Not covered.

### WRD-17

Name the action in words that fit the input: never tap on the desktop; press a key, type text, hold a key kept down, use the arrow keys; select an item in a list or a check box and clear it to uncheck. Call them keyboard shortcuts, never accelerators, mnemonics or hot keys; tab alone means a tab, otherwise the Tab key or a tab stop.

- **Windows:** Select is the input-neutral verb for buttons and commands as well as list items and check boxes (Select Save); click is allowed but not required; never choose or pick. Enter for a value that can also be chosen from a list.
- **macOS:** Choose for menu items (choose File > Save), click for buttons.
- **Sources** (19 records):
  - Apple: [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 5; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 15
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 54; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 94; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 96; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 32; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 71, 72, 73, 74, 75, 76, 77; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 64; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 37; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 54; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 55; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 43
  - Fluent: [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 6
- **Enforce:**
  - Test, Go source: Flag "tap", "hit", "hot key", "accelerator" in string literals.
  - Review: Verbs in help and instructions match the form.
- **kvit-ui now:** partly. `uitest.CheckRules` flags "tap" in a string (`uitest/wording.go`); "hit", "hot key" and "accelerator" are not flagged.

### WRD-18

Link text says what is behind it, in words that tell the screen's links apart; never "click here". A help link is the question its topic answers (How do I add exceptions?), the whole sentence is the link. A link that opens the web says so and carries the external-link symbol. A generic More information is acceptable only right after an instruction already on screen. An inline link keeps its sentence's capitalisation, with the full stop outside it. Links only navigate (CTL-26).

- **Sources** (23 records, 8 on Microsoft's checklist):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 9
  - Windows 7: [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 28, 29\*, 31, 32\*, 34\*, 35\*, 36\*; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 20, 29, 30, 31, 34, 37, 38, 41, 28, 35\*, 36; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 44; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 48\*, 128\*
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 30
- **Enforce:**
  - Test, one screen: Fail a Link named "here", "click here", "more", "link" or "Learn more"; fail a Link ending in a full stop; fail two Links on one screen with the same name and different targets.
  - Component: Link gains a web target that sets the "external" symbol (`arrow-square-out`, `icons/catalog_gen.go:1585`) and the domain as its tooltip.
  - Review: The link predicts its target.
- **kvit-ui now:** partly. `uitest.Labels` fails a Link named here, click here, more, link, learn more or read more, and one ending in a full stop, and `uitest.CheckRules` flags "click here" (`uitest/labels.go`, `uitest/wording.go`); nothing marks a web link or compares two links' targets.

### WRD-19

Every field, list, tree, slider and group of controls has a visible label: a noun or short noun phrase naming what it holds, not how to operate it, with no colon or other ending punctuation, saying when several may be chosen or how to separate them; instructions go above the label in full sentences and may be punctuated. A control that completes a check box or radio option uses that option's label. A prompt (placeholder) names the purpose or gives an example of the format; it stays while the field is empty, focused or not, goes at the first character and returns when the field is emptied. A prompt may stand in for the label only where the purpose is obvious from context (a search field, a toolbar field, a sign-in pair), and never holds information needed while typing.

- **Windows:** The prompt is a sentence-case fragment without ending punctuation or ellipsis.
- **Sources** (54 records, 16 on Microsoft's checklist):
  - Apple: [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 3; [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 2, 3; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 23; [combo-boxes](https://developer.apple.com/design/human-interface-guidelines/combo-boxes) 2; [pop-up-buttons](https://developer.apple.com/design/human-interface-guidelines/pop-up-buttons) 4
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 22\*; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 40, 44, 51; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 36\*, 40\*, 45\*, 37, 46\*, 52, 54, 55; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 45\*, 47\*, 51, 57; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 30\*, 61\*, 62\*, 63\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 79\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 59; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 39, 43, 45; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 23; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 25, 59\*, 68; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 27\*, 31\*, 35, 28, 45; [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 8\*; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 157
  - Fluent: [ui-controls-labels](https://learn.microsoft.com/en-us/windows/apps/design/controls/labels) 1, 3, 4, 5; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 5, 6; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 28; [ui-controls-number-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/number-box) 3, 4; [ui-controls-password-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/password-box) 2, 3; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 4
- **Enforce:**
  - Component: A form-row helper draws the label and sets `LabeledBy`. SearchField's default prompt loses its ellipsis on the Windows form.
  - Test, one screen: Fail a label ending in a colon; fail a Field whose only name is its placeholder unless it is a SearchField, sits in a toolbar or is listed as obvious from context with a reason; fail a placeholder ending in punctuation on the Windows form.
- **kvit-ui now:** partly. FormRow draws a field's label, without a colon, and names the field by it (`formrow.go`); `uitest.Labels` fails a label ending in a colon and, on the Windows form, a placeholder ending in punctuation (`uitest/labels.go`); a field named only by its placeholder is not flagged.

### WRD-20

Name what a group of things is with a specific noun: column headings, tabs, segments, panel titles, settings sections, sidebar groups, submenus, and menu titles (one word each in a menu bar). No catch-all such as General, Advanced or Settings unless the page really holds common or expert-only items that fit nowhere else; a settings section is one word where possible and never contains the word settings (Accounts, not Account settings); a tab for a document or view carries its name. No colon or other ending punctuation. A group label says what the choice is for, and that at least one is needed when it is, not how to choose.

- **Windows:** A menu title may be a verb (Insert, Format), and is a noun when the verb would be Create, Show, View or Manage. List groups by date use Today, Yesterday, Last week, Earlier this year, A long time ago.
- **Sources** (39 records, 8 on Microsoft's checklist):
  - Apple: [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 6, 7; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 3; [tab-views](https://developer.apple.com/design/human-interface-guidelines/tab-views) 3; [segmented-controls](https://developer.apple.com/design/human-interface-guidelines/segmented-controls) 8; [sidebars](https://developer.apple.com/design/human-interface-guidelines/sidebars) 8; [boxes](https://developer.apple.com/design/human-interface-guidelines/boxes) 3; [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 6
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 42, 44, 41; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 39, 38; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 82, 14, 83; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 49, 46\*; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 21\*, 24, 20, 28\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 23, 24, 32\*, 31; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 86\*, 87\*, 88\*; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 42; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 8, 11\*, 13; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 81
  - Fluent: [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 10; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 29; [ui-controls-labels](https://learn.microsoft.com/en-us/windows/apps/design/controls/labels) 2; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 9; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 3
- **Enforce:**
  - Test, one screen: Fail a Tab, column title, segment, group label or section heading ending in punctuation, and a Tab, column title or segment over three words; fail a Tab named General, Advanced or Settings unless listed with a reason; fail a section heading containing "settings".
  - Review: Nouns, specific to the content.
- **kvit-ui now:** partly. `uitest.Labels` fails a tab or segment ending in punctuation or over three words (`uitest/labels.go`); column titles and section headings are not read.

### WRD-21

A multi-step flow is named by its task (Connect to a network), never with the word wizard; every page has a main instruction (WRD-26); one advancing word is used throughout, and the flow can always be left (WIN-22).

- **Windows:** Next only while Back or Cancel can still undo the step. The step that commits has a specific verb answering the page's instruction (Print, Connect); Finish only when that verb would be generic (Save, Choose) or the flow only changes settings. Pages after the work is done end with Close, never Done or Cancel. Back is a button labelled Back beside Next among the bottom buttons.
- **macOS:** Continue (or Next) advances throughout and Done ends the flow.
- **Sources** (15 records, 12 on Microsoft's checklist):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 12
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 127\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 31\*, 38\*, 63\*, 75, 74\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 91\*, 93\*, 96\*, 97\*; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 82\*, 83\*, 84\*
  - Fluent: [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 3
- **Enforce:**
  - Component: A flow component (none exists) takes the commit verb from the app and supplies Back/Next/Finish/Close or Continue/Done by form, with Back beside Next on the Windows form.
  - Test, all screens of an app: Across an app's flows, the advancing button has one name; on the Windows form no button is named Done.
- **kvit-ui now:** does not cover it. No wizard or flow component.

### WRD-22

A setting has a plain label and at most one short line saying what it does; anything longer goes behind a Hint. Explain only the options that need it, never restating a label; a short note may sit in parentheses after the label or below the control; say when a choice becomes the default for next time. Settings text describes; it does not order people about. Skippable detail goes in a footnote or behind Help.

- **Windows:** On the Settings page each setting is a row: its label and an optional one-line description at the leading side, its control at the trailing side.
- **Sources** (23 records, 2 on Microsoft's checklist):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 21; [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 23
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 133, 134, 141\*, 151; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 52; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 56; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 86, 88; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 58; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 44\*; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 46; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 73, 71; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 46; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 34, 33; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 39, 40; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 36, 37
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 18
- **Enforce:**
  - Test, one screen: Fail a static Label under a control longer than about 20 words; fail a RadioOption Detail or Button Explanation that repeats its label.
  - Component: Hint (`popover.go:159-213`) is the place for a longer explanation; RadioGroup's per-option Detail (`radiogroup.go:19`) is the one-line note.
- **kvit-ui now:** enforces it. SettingRow has one line for a setting's description (`settings.go`) and Hint holds a longer one (`popover.go`); `uitest.Labels` fails a label over 20 words under a control, and a radio option's detail or a button's explanation repeating its label (`uitest/labels.go`).

### WRD-23

Call things by what the reader sees, never by implementation words (popover, pane, scene, modal, list view, text box). Help names a control by its exact label, without colon or access-key mark.

- **Windows:** Help adds the plain kind after the label: the Password box, the Programs list, the Underline check box, the View tab, the Screen resolution slider, the Network icon in the notification area. Settings is the Settings page; a flow may be called the <Name> wizard. A group is named only to tell apart options with the same name (Under Effects, select Hidden). The area beside the clock is the notification area, never the system tray.
- **macOS:** Help refers to the control by its label or to the task, without the kind.
- **Sources** (38 records):
  - Apple: [windows](https://developer.apple.com/design/human-interface-guidelines/windows) 5; [panels](https://developer.apple.com/design/human-interface-guidelines/panels) 10; [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 12
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 45; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 53; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 31; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 16, 15; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 57, 58; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 80, 82; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 90, 91; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 59, 60; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 35, 34, 36; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 57, 54; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 33; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 58, 59; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 47, 48; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 42; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 28; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 179, 204, 205, 206; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 40, 41; [ctrl-spin-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-spin-controls) 20; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 92; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 49; [vis-color](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-color) 19
- **Enforce:**
  - Test, Go source: Flag popover, modal, list box, text box, list view and system tray in user-visible string literals.
  - Review: Implementation words in visible text.
- **kvit-ui now:** does not cover it. Not covered.

### WRD-24

Write shortcuts in the platform's form.

- **Windows:** Key names as printed (Ctrl, Alt, Shift, Enter, Page Up, Tab), letters in capitals, joined with + and no spaces (Ctrl+Shift+N), modifiers in the order Windows logo, Ctrl, Alt, Shift. A shifted punctuation key is written with Shift and the shifted character (Ctrl+Shift+?). In running text, punctuation and arrows are spelled out (Ctrl+plus sign, left arrow); never glyphs or all capitals.
- **macOS:** Glyphs ⌃⌥⇧⌘ in that order, and the shifted symbol itself (⌘?).
- **Sources** (8 records, 3 on Microsoft's checklist):
  - Apple: [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 11, 12
  - Windows 7: [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 20\*; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 67, 68\*, 69, 70; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 36\*
- **Enforce:**
  - Component: Menu shortcut text comes from unison's `KeyBinding.String()` (`menu.go:257`); kvit-ui wraps it in one function that, on the Windows form, renames Return to Enter and PageUp to Page Up and puts the logo key first.
  - Test, one screen: Every menu shortcut text on the Windows form matches `^(Win\+)?(Ctrl\+)?(Alt\+)?(Shift\+)?\S+$` with Windows key names.
- **kvit-ui now:** enforces it. `UI.ShortcutText` writes a shortcut in the form's way: "Ctrl+Shift+S" with Windows' key names (Enter, Del, PgUp, Plus) on the Windows form, "⇧⌘S" with the Mac's symbols on the macOS form (`shortcut.go`); menus and tooltips use it.

### WRD-25

A progress description says exactly what is happening, never just Working… or Please wait: an -ing phrase naming the operation, ending in an ellipsis (Printing your pictures…, Copying…). Loading… alone is allowed only beside a bar where what is loading is obvious from the screen (Fluent's own examples show it). Work the reader did not start says why they wait (Please wait while …). Time remaining runs from the largest unit down, two units at most, no more precise than the estimate (hours and minutes, no seconds), singular for one (1 minute). Each progress detail has a label and a unit (Speed: 2 MB/s).

- **Sources** (17 records, 3 on Microsoft's checklist):
  - Apple: [progress-indicators](https://developer.apple.com/design/human-interface-guidelines/progress-indicators) 7
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 194, 79, 80\*, 83; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 55; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 26; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 34, 35, 38, 47, 49, 44; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 54\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 23\*
  - Fluent: [ui-controls-progress-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/progress-controls) 4, 7
- **Enforce:**
  - Component: Button.BusyText loses its "Working…" default and becomes required with Busy; Progress requires Label and gains `Remaining`, written by one library function using `UI.CountPhrase` (`numbers.go:21-37`).
  - Test, one screen: Fail a busy text or progress label of Working… or bare Please wait; fail Loading… on a busy button.
  - Review: Loading… only where the object is obvious from the screen.
- **kvit-ui now:** enforces it. Button has no busy default and warns once about a busy button without BusyText (`button.go`); Progress shows no percentage unless asked (`progress.go`); `uitest.Labels` fails Working…, Please wait and Loading… on a busy button (`uitest/labels.go`).

### WRD-26

Each dialog, page or message leads with one main instruction: what the reader should do or decide, as an imperative or a specific question in terms of their goal, in one sentence, essentials first. Leave it out where the purpose is obvious (Open, Save). Supplemental text only adds context, in short sentences, with detail behind Help; a dialog the app opened by itself says why. The buttons answer the instruction directly. Say each thing once across title, instruction, text and buttons, keeping the full words in the instruction and on the controls; a control label that only repeats the instruction is dropped. A confirmation asks Are you sure you want to… only for the direct result of a request and Do you want to… for a side effect, says permanently when it cannot be undone, and never writes Warning or Caution.

- **Sources** (50 records, 15 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 177\*, 178, 187\*, 188, 189, 190, 197; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 27\*, 28, 31, 32, 36; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 3; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 22, 54\*, 82; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 40\*, 41, 48, 49, 50; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 47\*; [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 33\*, 36; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 10, 11, 71\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 15; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 46\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 15\*, 26\*, 27\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 39; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 1, 3, 4, 5\*, 6, 48\*, 49\*, 51, 52, 60; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 40; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 43
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 2, 6, 10; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 5, 23
- **Enforce:**
  - Component: Dialog's Title is the main instruction and Detail the supplemental text (`feedback.go:385-386`); its doc says so, and a separate window's title (WIN-12) is never repeated as the instruction.
  - Test, one screen: Fail a dialog whose Title ends in a full stop, holds more than one sentence, or is Warning, Error or Caution; fail a dialog whose Detail begins with its Title's words.
  - Review: Says each thing once; written from the reader's goal.
- **kvit-ui now:** partly. `uitest.Labels` fails a dialog title ending in a full stop, a title of Error, Warning or Caution, and a detail starting with the title's words (`uitest/labels.go`); a title of more than one sentence is not flagged.

### WRD-27

Labels, buttons, menu items, tabs, titles, list items, prompts and main instructions end without a full stop; a question ends with a question mark wherever it is. Complete sentences in supplemental text and explanations end with a full stop (tooltips: WRD-13); status text is short fragments without one. One space between sentences. Double quotation marks around object names such as file names, not around paths, web addresses or names obvious from context. A colon introduces a list; use the serial comma; write and, or and number in words, not &, +, / or #. Emphasis only for what must not be missed: semibold rarely, never capitals; no italics in interface text (they read poorly for people with dyslexia), so a word mentioned as a word goes in quotation marks (decided, Q-51).

- **Sources** (44 records, 14 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 193\*, 203; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 31; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 63, 74\*; [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 11; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 73, 79\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 80; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 55\*; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 29; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 25\*, 28\*, 29\*, 64\*; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 46; [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 13\*, 14; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 9\*, 12\*, 15, 17, 20\*, 21, 22, 24, 26, 30, 31, 37, 54\*, 75; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 43; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 29; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 25; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 34, 42; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 24, 30; [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 20\*
  - Fluent: [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 34; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 17, 18; [signature-experiences-typography](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/typography) 9
- **Enforce:**
  - Test, one screen: Fail a Button, MenuItem, Tab, column title or control label name ending in a full stop; fail two spaces after a full stop; fail an object name in single quotation marks.
  - Test, Go source: Flag " & ", " / " and "#" inside label literals ("&" before a letter is the access-key mark, "&&" a literal ampersand).
  - Review: Emphasis and quotation.
  - Test, Go source: Flag `Italic: true` in a text.Style outside the gallery (`text/layout.go:35`).
- **kvit-ui now:** partly. `uitest.Labels` fails a command, check box or switch ending in a full stop and a command with two spaces (`uitest/labels.go`); tabs, column titles and quoted object names are not read for it.

### WRD-28

Be polite, calm and positive; never blame, talk down to or alarm the reader. In messages say problem, not error or failure; unable to or can't, not failed to; incorrect or not valid, not illegal, invalid or bad; stop, not abort, kill or terminate; serious, not fatal or catastrophic. Say please only when asking the reader to put up with something, sorry only for serious problems such as lost data, now or immediately only when it is urgent. Use the active voice and the present tense, contractions where they read naturally, but not with names and never forced to save space, and say what the reader can do rather than what they cannot.

- **Sources** (21 records, 4 on Microsoft's checklist):
  - Windows 7: [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 33; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 56\*, 60, 61; [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 3, 6, 7, 19, 25, 26, 27, 28\*, 29, 30, 31; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 14\*, 113\*
  - Fluent: [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 7, 15, 16, 20
- **Enforce:**
  - Test, Go source: Flag error, failure, failed, illegal, invalid, abort, kill, terminate, fatal, catastrophic and sorry in string literals passed to Dialog, Notice, Toast and Field.Error.
  - Review: Tone, blame, needless apology.
- **kvit-ui now:** does not cover it. Not checked. The library's own messages avoid the words ("Not a number", "Must be at least N", `numberfield.go:109-113`).

### WRD-29

Labels of the options in one group (check boxes, radio buttons, list options, slider ends, answers) are parallel in form and similar in length, each distinct, and state only what differs: shared words go in the group label. Choosing an option turns something on (Don't show again is the accepted exception). A choice meaning none or all is in parentheses ((None)); a strongly recommended one adds (recommended) to its own label, an expert one (advanced); a unit goes in parentheses after the label (Size (MB)). A check box label is at most two lines, a radio button label one line, and each slider end label one word where possible (decided, Q-51). Answers too long for buttons are radio buttons, each with one line of explanation, never command links (decided, Q-10). A label is never half a sentence completed by a control (Remind me [5] minutes before), and text is never assembled from pieces at run time.

- **Sources** (56 records, 1 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 136, 137, 200; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 33, 34, 35, 38, 40; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 45; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 12, 14\*; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 51, 53; [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 18, 37; [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 22, 23; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 64; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 27, 52, 53, 56; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 27, 28, 29, 32; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 7, 42, 43, 65, 66, 67, 69, 70; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 18, 39, 40, 41, 44; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 30, 31, 32, 35, 37; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 28, 29, 31, 32, 37, 38, 41
  - Fluent: [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 7, 8; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 32, 33; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 10
- **Enforce:**
  - Component: RadioOption, Select options and Check gain `Recommended bool`; the component adds "(recommended)" in the app's language, so the app never types it.
  - Test, one screen: Fail two options with the same name in one radio group, select or list; fail a check box or radio name starting Don't or Do not other than Don't show again; fail options that all start with the same word as the group label.
  - Test, one screen: Fail a radio option name that wraps to a second line or a check box name over two lines at the window's default width.
  - Review: Parallel form; no sentence built around a control.
- **kvit-ui now:** partly. `uitest.Labels` fails two radio options with one name and a check box or switch named with Don't or Do not other than a don't-show-again option (`uitest/labels.go`); wrapping option names are not measured.

## Commands: menus, toolbars, buttons

### CMD-1

Every command the app has is on a visible command surface (the menu bar, or where the platform allows, a toolbar and its More menu), including ones also on buttons and context menus, never only in a context menu, so every task can be found on screen rather than remembered.

- **Windows:** A menu bar (Kvit-drawn, in the main window) is for an app with many commands in several categories, such as an editing app (Notepad). Another main window may use a toolbar (command bar) with a More menu as its command surface instead. Either way every command is reachable from the menu bar, or from the toolbar and its More menu. Offer each command on as many surfaces as fit: menu bar or toolbar, and the context menu (decided, Q-1).
- **macOS:** Required: the system menu bar, always present.
- **Sources** (33 records, 2 on Microsoft's checklist):
  - Apple: [designing-for-macos](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) 3; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 44, 69; [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 31; [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 4; [dock-menus](https://developer.apple.com/design/human-interface-guidelines/dock-menus) 2; [split-views](https://developer.apple.com/design/human-interface-guidelines/split-views) 7
  - Windows 7: [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 14; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 9, 22, 8, 34, 10, 35; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 3, 33, 76, 1, 2, 9, 11, 12, 35, 57\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 77\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 30
  - Fluent: [ui-controls-menus-and-context-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus-and-context-menus) 8, 9; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 5, 6; [ui-controls-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus) 8; [ui-controls-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/commanding) 1; [basics-app-silhouette](https://learn.microsoft.com/en-us/windows/apps/design/basics/app-silhouette) 3
- **Enforce:**
  - Component: A Command type (name, symbol, shortcut, enabled, run) that a MenuItem, an IconButton, a toolbar button and a shortcut are made from; the menu bar is built from the app's commands, Kvit-drawn in the window on Windows and Linux, the system's on macOS.
  - Test, all screens of an app: Every command reachable from a button or context menu appears in the menu bar, or, on the Windows form without one, on the toolbar or in its More menu.
- **kvit-ui now:** enforces it. `Command` is made once and shown by MenuItem, Button, IconButton, Toolbar and MenuBar (`command.go`, `menubar.go`, `toolbar.go`); `Window.OffSurface` lists a command on no menu bar, toolbar or sidebar, and `uitest.App` fails it (`uitest/app.go`).

### CMD-2

Put the platform's standard menus in its standard order, leave out only the ones the app has no use for, and keep to at most ten menus, none holding only one or two items.

- **Windows:** File, Edit, View, the app's own menus, Tools, Help; no app menu and no Window menu; Settings is the last item of Tools when the window has no sidebar (WIN-16). A program that does not create or view documents may name and group its menus by its readers' tasks (Play, Library, Rip) instead. unison's per-window bar has File, Edit, Window, Help, with Exit in File and Preferences in Edit (unison `menu_bar.go:84-136`), so kvit-ui builds its own on Windows and Linux.
- **macOS:** App, File, Edit, Format, View, app menus, Window, Help.
- **Sources** (13 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 1, 2, 43, 15, 37, 40, 48, 38
  - Windows 7: [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 10, 17, 18, 15, 16
- **Enforce:**
  - Component: A `kvitui.MenuBar` with a platform form: on macOS unison's `InsertStdMenus` (unison `menu_bar.go:53-63`) with the app's menus between View and Window; on Windows and Linux File, Edit, View, app menus, Tools, Help, drawn by Kvit.
  - Test, all screens of an app: Fail a menu bar of more than ten menus or with a menu of one or two items.
- **kvit-ui now:** enforces it. MenuBar shows the standard menus in the form's order, File, Edit, Format, View, the app's menus, Tools and Help on the Windows form, the app menu, File, Edit, Format, View, the app's menus, Window and Help on the macOS form, and leaves an empty one out (`menubar.go`); `uitest.App` fails more than ten menus and an app menu of one or two lines (`uitest/app.go`).

### CMD-3

About, the settings command and Quit or Exit sit where the platform puts them, with the same short app name; About holds the app's own facts, such as its version.

- **Windows:** No app menu. Settings (named Settings, not Options) is the last item of the navigation sidebar, or the last item of the toolbar's More menu or on the toolbar; with a menu bar and no sidebar, the last item of Tools (WIN-16). About is a section at the bottom of the Settings page, not a Help menu item: the app's name, icon and version, and links to its website, a place to report problems or ask for features, the third-party notices and the legal terms. Exit is the last item of File, its Alt+F4 usually not shown.
- **macOS:** The app menu holds About, Settings…, app-level configuration, Services, Hide, Hide Others, Show All and Quit, in that order.
- **Sources** (11 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 9, 10, 11, 13, 14, 12; [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 9
  - Windows 7: [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 3; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 71
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 4, 25
- **Enforce:**
  - Component: MenuBar places the settings command and Quit/Exit by platform; the macOS label is "Settings…" (unison still says "Preferences…", `menu_bar.go:132`, so kvit-ui passes its own item), the Windows and Linux labels are Settings and Exit. The Settings page (WIN-16) ends with an About section on Windows and Linux.
- **kvit-ui now:** enforces it. MenuBar puts Settings last in Tools on the Windows form, unless the sidebar holds it at its foot (`Sidebar.AddFooter`, `NewSidebarCommand`), and Settings… and Quit <App> in the app menu on the macOS form, renaming unison's "Preferences…"; Exit is last in File on the Windows form (`menubar.go`, `sidebar.go`). `NewAboutSection` is the About section that ends the settings page (`settings.go`).

### CMD-4

The File menu holds the document commands in the platform's standard order, as the app needs them.

- **Windows:** New, Open…, Close; Save, Save as…; Send to; Print…, Print preview, Page setup; the numbered recent files inline; Exit (Alt+F4). Save as… saves under a name the reader gives, where the Mac offers Duplicate.
- **macOS:** New, Open, Open Recent, Close, Save, Duplicate, Rename…, Move To…, Export As…, Print…; Duplicate rather than Save As.
- **Sources** (17 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 16, 19, 20, 22, 21, 26, 27, 28, 29; [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 1, 2; [printing](https://developer.apple.com/design/human-interface-guidelines/printing) 1; [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 11; [windows](https://developer.apple.com/design/human-interface-guidelines/windows) 3
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 49; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 50; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 70
- **Enforce:**
  - Component: MenuBar offers the standard items with their standard shortcuts.
  - Review: Order and wording of app-specific items.
- **kvit-ui now:** enforces it. MenuBar puts the standard File commands the app has in the form's order with their standard keys (`menubar.go`, `command.go`).

### CMD-5

The Edit menu starts with Undo and Redo, then Cut, Copy and Paste, Select all and Delete, then Find, in the platform's order; app-specific edit commands sit beside the standard ones they relate to.

- **Windows:** Undo, Redo (Ctrl+Y and Ctrl+Shift+Z); Cut, Copy, Paste; Select all; Delete, Rename…; Find… (Ctrl+F), Find next (F3), Replace… (Ctrl+H), Go to… (Ctrl+G).
- **macOS:** Undo, Redo; Cut, Copy, Paste, Delete, Select All; Find; ending with the system's Dictation and Emoji items.
- **Sources** (11 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 32, 31, 36; [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 14, 1, 4, 5; [undo-and-redo](https://developer.apple.com/design/human-interface-guidelines/undo-and-redo) 9; [image-wells](https://developer.apple.com/design/human-interface-guidelines/image-wells) 2
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 50; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 72
- **Enforce:**
  - Component: MenuBar's Edit menu in the platform's order; Field's context menu gains Undo and Redo.
- **kvit-ui now:** enforces it. MenuBar's Edit menu takes the standard commands in the form's order (`menubar.go`); Field's context menu has Undo and Redo, from the window's undo manager (`field.go`).

### CMD-6

The View menu holds the show and hide commands for the window's parts, zoom and full screen.

- **Windows:** Check items for the toolbars, panes and the status bar; then Zoom in (Ctrl++), Zoom out (Ctrl+-), Text size; then Full screen (F11) and Refresh (F5).
- **macOS:** Show/Hide items and Enter Full Screen; window switching goes in the Window menu.
- **Sources** (11 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 40, 42, 39, 52; [sidebars](https://developer.apple.com/design/human-interface-guidelines/sidebars) 5; [going-full-screen](https://developer.apple.com/design/human-interface-guidelines/going-full-screen) 15
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 36, 52; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 13, 73; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 10
- **Enforce:**
  - Component: MenuBar's View menu; Window offers Show/Hide Sidebar as a command (a check item on Windows and Linux).
- **kvit-ui now:** partly. MenuBar's View menu puts the app's own lines, such as the window's parts, before Zoom, Full screen and Refresh (`menubar.go`); Window has `FullScreenCommand` (`windowstate.go`), but no built-in Show or Hide sidebar command.

### CMD-7

On a Mac every app has a Window menu with Minimize, Zoom, the panels' show and hide items and the list of open windows.

- **Windows:** No Window menu: the standard Windows menus have none, and the taskbar switches windows.
- **macOS:** As stated.
- **Sources** (6 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 47, 48, 49, 50, 53; [panels](https://developer.apple.com/design/human-interface-guidelines/panels) 7
- **Enforce:**
  - Component: MenuBar uses unison's `NewWindowMenu` (`menu_bar.go:138`).
- **kvit-ui now:** enforces it. MenuBar always has a Window menu on the macOS form: unison's `NewWindowMenu` in the system's bar on a Mac, and Minimize and Zoom in the drawn preview elsewhere (`menubar.go`). The system's bar is built but has not been run on a Mac.

### CMD-8

The Help menu is last; its first item opens the app's help, set apart from release notes and the like, and it stays short.

- **Windows:** <App> help (F1) first; links to report a problem, give feedback or read the terms may follow, below the help. About is not in the Help menu: it is a section at the bottom of the Settings page (CMD-3). Product registration, if any, is a Help command and a link in About, never a step at first start.
- **macOS:** First item "<App> Help"; About is in the app menu.
- **Sources** (11 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 54, 55, 56
  - Windows 7: [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 8; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 54; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 75; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 16, 40
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 24, 32; [in-app-help-guidelines-for-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/guidelines-for-app-help) 5
- **Enforce:**
  - Component: MenuBar's Help menu.
- **kvit-ui now:** enforces it. Help is the last menu, its first line the standard Help command, "<App> help" on F1 on the Windows form and Cmd-? on the macOS form (`menubar.go`, `command.go`).

### CMD-9

A main-menu item or toolbar button that cannot act now stays in place, disabled, and a menu whose items are all disabled still opens; only a menu that exists for one mode or selection drops the items that do not apply. Where working out whether an item can act would make the app noticeably slow, it stays enabled and says why when chosen.

- **Sources** (13 records):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 5, 6, 19, 23; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 3; [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 3; [printing](https://developer.apple.com/design/human-interface-guidelines/printing) 3
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 18, 17; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 31, 32, 30; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 24
- **Enforce:**
  - Component: MenuItem.Disabled (exists); the Command type computes it, with a flag for commands too slow to check that report the reason on use.
  - Review: An item that appears and disappears with the selection.
- **kvit-ui now:** enforces it. A menu line, button or symbol button made from a Command reads whether it can act each time it is shown, and its key does nothing while it cannot (`command.go`, `commandkeys.go`); a command too slow to check leaves `Enabled` nil and says why when run.

### CMD-10

Every object and region with commands of its own has a context menu holding the commands for what was clicked, in at most about four groups and 15 lines, leaving out those that do not apply (the platform form says whether the clipboard commands stay); commands that everyone must find are also visible in the window or behind a visible menu button.

- **Windows:** The context menu holds every command that acts on the item, since it is the route that works with every input. The common commands (Cut, Copy, Paste, Delete, Rename, Share) come first, drawn as a row of symbols with tooltips above the other commands, as in Windows 11; the item's other commands follow. The default command (what double-click or Enter does) is bold; a checked line keeps its tick and is not bold (decided, Q-33). Commands that do not apply are left out, not disabled. In a text field's menu, Cut and Copy appear only with a selection, Paste only with text on the clipboard, Undo only after a change and Select all only when there is text; it adds spelling suggestions on a misspelt word, and Bold, Italic and Underline in rich text. Selecting text may also show a small bar of these commands beside it, which takes no focus, starts with the most important few and goes away if unused. A list whose view, sorting or grouping can change offers View, Sort by and Group by submenus.
- **macOS:** Only Cut, Copy and Paste stay when they do not apply; no default item is marked.
- **Sources** (45 records):
  - Apple: [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 1, 2, 7, 10, 8, 3, 12; [token-fields](https://developer.apple.com/design/human-interface-guidelines/token-fields) 1
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 18, 27, 17; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 6, 36, 37, 39, 41, 40, 42, 4, 34; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 23; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 23; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 35
  - Fluent: [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 21; [ui-controls-menus-and-context-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus-and-context-menus) 4, 5, 6; [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 1, 2, 5, 7, 17; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 25; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 8; [ui-controls-command-bar-flyout](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar-flyout) 1, 2, 3, 5, 6, 13, 14, 15; [ui-controls-dialogs-and-flyouts-flyouts](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/flyouts) 3; [ui-controls-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/commanding) 5; [ui-controls-nested-ui](https://learn.microsoft.com/en-us/windows/apps/design/controls/nested-ui) 5
- **Enforce:**
  - Component: MenuItem gains Default, drawn bold by the Kvit menu on Windows and Linux; a Checked line keeps its tick and loses the bold it has today (`menu.go:236-239`); on Windows and Linux the Kvit menu draws the common commands as a leading row of symbols, and Field's menu leaves out a command that cannot act instead of disabling it.
  - Test, one screen: Fail a context menu over 15 lines or 4 groups, or, on the Windows form, one with a disabled item or with common commands not at its head.
  - Review: Relevance of the items; every command for an item in its menu; comparable objects all having one.
- **kvit-ui now:** partly. A context menu leaves out a line that cannot act, keeping Cut, Copy and Paste on the macOS form; `Default` is bold on the Windows form; a ticked line is no longer bold (`menu.go`); Field's menu is built from the standard commands (`field.go`). `uitest.Labels` fails an open context menu over 15 lines or four groups and a disabled line on the Windows form (`uitest/labels.go`). There is no leading row of common commands.

### CMD-11

On a Mac a context menu shows no shortcuts; the main menus teach them.

- **Windows:** Shown, as in the other menus (the Windows 7 records say nothing against it; Kvit's field menu shows them, `field.go:211`).
- **macOS:** No shortcuts in context menus.
- **Sources** (1 records):
  - Apple: [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 9
- **Enforce:**
  - Component: ShowMenu for a context menu drops `Key` on macOS.
- **kvit-ui now:** enforces it. A context menu shows no shortcuts on the macOS form (`menu.go`, `menuLine`).

### CMD-12

Submenus only where they clearly help, one level deep, about five items at most; three or more items sharing a leading word become one submenu.

- **Windows:** No submenus in a context menu at all. Windows 7 allows a rare second level in the main menus; Kvit keeps one level, which satisfies both.
- **macOS:** A submenu in a context menu where it clearly helps, as in other menus.
- **Sources** (13 records):
  - Apple: [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 5; [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 18, 16, 17, 20, 14, 15
  - Windows 7: [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 29, 38, 26, 27, 28; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 21
- **Enforce:**
  - Test, one screen: Fail a submenu inside a submenu; fail a submenu over 7 items; on the Windows form, fail a submenu in a context menu.
  - Review: Shared leading words.
- **kvit-ui now:** enforces it. `uitest.Labels` fails a submenu inside a submenu, a submenu over seven lines, and on the Windows form a submenu in a context menu (`uitest/labels.go`).

### CMD-13

Order a menu by its natural sequence where it has one (steps of a task, numbers in numeric order), otherwise most used first; group related items between separators, at most seven to a group (a submenu counts as one) and about 25 to a level; a rare command stays in its group, and a command sits in one menu only unless it truly belongs in two.

- **Sources** (13 records):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 11, 12, 13; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 45, 46
  - Windows 7: [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 19, 22, 21, 24, 25, 23, 47; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 29
- **Enforce:**
  - Test, one screen: Fail a menu group over seven items or a menu level over 25.
  - Review: Order, grouping and commands repeated across menus.
- **kvit-ui now:** enforces it. `uitest.Labels` fails a menu group over seven lines and a level over 25 (`uitest/labels.go`).

### CMD-14

Give a menu item a symbol only when it clarifies the item, mainly for the most used commands and those with a standard, well-known symbol, and then the same symbol everywhere the command appears; never an invented, cryptic one.

- **Sources** (7 records):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 9; [pull-down-buttons](https://developer.apple.com/design/human-interface-guidelines/pull-down-buttons) 10
  - Windows 7: [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 48
  - Fluent: [ui-controls-command-bar-flyout](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar-flyout) 9, 10; [ui-controls-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus) 4, 5
- **Enforce:**
  - Review: Decorative symbols.
- **kvit-ui now:** partly. Symbols are dropped on macOS menus (`menu.go:119-148`), which suits this rule there.

### CMD-15

Mark an option that is on; mark the chosen one of exclusive options in the platform's way; where attributes combine, one item clears them all.

- **Windows:** A tick for an independent on/off option, its words not bold (decided, Q-33); a bullet for the chosen one of a set of mutually exclusive options (at least two in the set); no mark when the option is on for only part of the selection.
- **macOS:** A tick for both; a dash when the option is on for only part of the selection.
- **Sources** (7 records):
  - Apple: [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 24, 25
  - Windows 7: [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 44, 46, 43, 45
  - Fluent: [ui-controls-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus) 3
- **Enforce:**
  - Component: MenuItem.Checked (exists) loses its bold and keeps the tick, gains a Radio form drawn as a bullet by the Kvit menu, and a Mixed state (no mark on Windows and Linux, a dash in the macOS menu); Select marks its current option as Radio.
- **kvit-ui now:** enforces it. A ticked line keeps its tick without bold; `Radio` draws the chosen one of a set with a bullet on the Windows form and a tick on the macOS form; `Mixed` draws nothing on the Windows form and a dash on the macOS form (`menu.go`); Select marks its options as a set (`select.go`).

### CMD-16

Holding a modifier never silently changes what a dialog button does. In a Mac menu a modifier may reveal an alternate command, which then also exists another way, uses one modifier, and never appears in a context or tray menu.

- **Windows:** No alternates revealed by a modifier, in menus or dialogs: they cannot be discovered (XP's Yes to All turning into No to All). Offer the variant as its own item.
- **macOS:** Alternates in the menu bar's menus, as stated.
- **Sources** (4 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 57, 58, 59
  - Windows 7: [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 8
- **Enforce:**
  - Review: Alternate items; dialog buttons that read the modifier keys.
- **kvit-ui now:** enforces it. MenuItem has no alternate (`menu.go:19-55`), so the Windows form holds and the Mac form is not offered.

### CMD-17

A menu holds only commands and options; a notification or confirmation goes in a dialog or popover. A menu button shows a down-arrow and lists at least three commands that refine its own action; variants of one command go in a split button whose main part runs a variant marked in its menu (the platform form says which); a feature that is on or off and comes in styles is a toggle split button, its main part turning it on or off and its menu picking the style; one or two commands are buttons, and a view's main actions are visible, not in a menu. Simple choices (check or radio items, a short list) may sit in a menu rather than a dialog.

- **Windows:** A split button's main part runs the variant last chosen (decided, Q-46).
- **macOS:** The main part runs the likeliest variant, or the last one chosen when people repeat it.
- **Sources** (27 records):
  - Apple: [pull-down-buttons](https://developer.apple.com/design/human-interface-guidelines/pull-down-buttons) 1, 3, 4, 5, 6, 7
  - Windows 7: [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 13, 26, 27, 28, 12; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 15, 100; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 9; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 7
  - Fluent: [ui-controls-menus-and-context-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus-and-context-menus) 1, 2, 3; [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 12; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 6, 7, 28, 29, 32; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 4; [ui-controls-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus) 1, 2
- **Enforce:**
  - Component: A MenuButton with a down-arrow, and a SplitButton: a Button and a chevron opening ShowMenu, the item the main part runs marked as default, with a toggle form.
  - Test, one screen: Fail a menu opened from a button with fewer than three items.
  - Review: Notifications or confirmations put in a menu.
- **kvit-ui now:** enforces it. MenuButton opens a menu under a down-arrow, and SplitButton's main part runs the variant last chosen on the Windows form and the likeliest on the macOS form; `NewToggleSplitButton` makes the toggle form, whose main part turns the feature on and off and whose menu picks the style, marked as the one in use, choosing a style turning the feature on (`menubutton.go`); `uitest.Labels` fails a menu button of fewer than three lines (`uitest/labels.go`). A notification or confirmation put in a menu is the reviewer's.

### CMD-18

Toolbars, menu bars and status bars belong to main windows; dialogs have none. The toolbar sits at the top of the window and holds the commands for the main tasks, mostly single-click buttons that act at once, each recognisable at a glance and shown once, in at most three groups; a command needed on several pages sits once on the toolbar rather than on each page; property buttons show and change the state of the selection; commands or a whole toolbar for one mode or selection show only while it is active, but a view's commands are never all hidden; a button's words never change while the app runs, only its symbol may show state.

- **Windows:** Commands start at the leading edge; the More button, Help, and the standard Settings and View commands sit at the trailing end. Beside a menu bar a toolbar button takes the immediate form of a command (Print prints at once) and leaves the asking form (Print…) to the menu; a toolbar standing in for the menu bar takes the asking form.
- **macOS:** At most one primary action, at the trailing end; navigation and critical actions get their own group.
- **Sources** (42 records):
  - Apple: [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 14, 15, 16, 19, 23, 24, 26, 27, 1, 20, 21, 22, 25
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 14, 19; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 7; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 2, 4, 5, 7, 19, 24, 25, 27, 40, 20, 21, 1, 3, 13, 14, 23, 28, 30, 45; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 42; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 3
  - Fluent: [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 1, 4, 15, 16; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 13
- **Enforce:**
  - Component: A Toolbar component (none today; ViewHead holds view controls) with leading, centre and trailing groups, placed by platform; Dialog takes no toolbar.
  - Test, one screen: Fail a toolbar that repeats a command, a Toolbar or StatusBar inside a Dialog, or a view whose toolbar and More button are both hidden.
  - Review: Choice of commands.
- **kvit-ui now:** partly. Toolbar has leading, centre and trailing groups, and the window answers its commands' keys while it is shown (`toolbar.go`); `uitest.Labels` fails a toolbar repeating a command and a toolbar or status bar inside a dialog (`uitest/labels.go`); a view whose toolbar and More button are both hidden is not flagged.

### CMD-19

When a toolbar does not fit, its commands, added in order of importance, move from the least important into a More menu that exists only then, before any other content on the bar (a title, a status) is cut; every command stays reachable.

- **Sources** (10 records):
  - Apple: [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 2, 3, 4, 9; [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 17
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 15, 32, 33
  - Fluent: [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 6, 13
- **Enforce:**
  - Component: Toolbar overflow (as StatusBar already does for facts, `statusbar.go:38-45`).
  - Test, one screen: At the narrowest window, every command is still on screen or in the overflow menu.
- **kvit-ui now:** enforces it. A toolbar too narrow for its commands moves the words into the tooltips first and then the least important commands into a More menu that exists only then (`toolbar.go`).

### CMD-20

Let people choose which toolbar items show when there are many or the app is used for long sessions; customisation starts locked and can be reset to the original.

- **Windows:** No Customize command beside Settings: customisation is reached from Settings. A customisable toolbar starts locked; its context menu lists the toolbars as check items, Lock the toolbars, and Customize…; the customisation window has Reset. The settings command may itself be a toolbar item (CMD-3).
- **macOS:** View > Customize Toolbar…, the system's customisation sheet; the settings command is not a toolbar item.
- **Sources** (12 records):
  - Apple: [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 5; [designing-for-macos](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) 6; [printing](https://developer.apple.com/design/human-interface-guidelines/printing) 2; [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 10
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 64, 65, 67, 16, 62, 66, 68
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 13
- **Enforce:**
  - Review: Toolbar customisation offered where it pays.
- **kvit-ui now:** does not cover it. Not covered.

### CMD-21

No on-screen buttons for Cut, Copy and Paste; undo and redo buttons only where they are really needed, with the standard symbols.

- **Sources** (3 records):
  - Apple: [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 9; [undo-and-redo](https://developer.apple.com/design/human-interface-guidelines/undo-and-redo) 5, 6
- **Enforce:**
  - Test, one screen: Fail an IconButton or Button named Cut, Copy or Paste on a screen that has an Edit menu.
- **kvit-ui now:** enforces it. `uitest.Labels` fails a Cut, Copy or Paste button on a screen with an Edit menu (`uitest/labels.go`).

### CMD-22

Show text on a button when a word or two says more than a symbol; use a symbol for a familiar action; never type letters or punctuation as a stand-in for a symbol.

- **Windows:** On a wide window a toolbar button shows its words beside its symbol; on a narrow one they move to its tooltip.
- **macOS:** As stated.
- **Sources** (17 records):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 8, 12; [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 16; [icons](https://developer.apple.com/design/human-interface-guidelines/icons) 7, 13
  - Windows 7: [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 47, 24; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 17; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 76; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 6, 11, 26, 43; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 34
  - Fluent: [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 10; [ui-controls-icons](https://learn.microsoft.com/en-us/windows/apps/design/controls/icons) 2, 3
- **Enforce:**
  - Test, one screen: Fail a Button whose words are one or two characters, or only punctuation (B, I, <>, +); those are IconButtons with a Label.
  - Test, Go source: Flag a Button text starting with "+ " or "- ".
- **kvit-ui now:** enforces it. The catalogue has bold, italic, underline, strikethrough, code and heading (`icons/catalog_gen.go`); `uitest.Labels` fails a button of one or two characters or punctuation, and `uitest.CheckRules` a button's words starting "+ " or "- " (`uitest/labels.go`, `uitest/wording.go`).

### CMD-23

A destructive command in a menu sits last, after a separator, marked as destructive, and asks before running unless it can be undone (MSG-3).

- **Sources** (4 records):
  - Apple: [pull-down-buttons](https://developer.apple.com/design/human-interface-guidelines/pull-down-buttons) 8; [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 18
  - Windows 7: [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 10; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 31
- **Enforce:**
  - Component: ShowMenu inserts the separator before a trailing Danger item itself instead of leaving it to the caller.
  - Test, one screen: Fail a Danger item that is not last or has no separator before it.
- **kvit-ui now:** enforces it. A menu puts a separator before a destructive line at its end by itself (`menu.go`, `tidySeparators`); `uitest.Labels` fails a destructive line that is not last (`uitest/labels.go`).

### CMD-24

A tray icon (Windows, Linux) or Dock menu (macOS) serves an app that works in the background: the icon shows its status and is not a launcher, and its menu lists the open windows and a few actions useful with no window open, each also in the menu bar or a window.

- **Windows:** The menu opens at the icon: Open (the default, bold) and Run; a separator; ticked toggles (Suspend and Resume for background work, Minimize to tray, notification opt-ins, Display icon in notification area); a separator; Settings and Exit. Exit quits for this session and removes the icon; no About unless the app has no window. A click on the icon may open a light flyout: a summary matching the tooltip, the common tasks and settings applied at once, closed by a click outside. Minimise to the tray only when the reader opts in, for a long-running single-instance app, and only through Minimize, never Close. A hidden icon comes back only from Settings.
- **macOS:** The Dock menu; a menu bar extra, where there is one, opens its menu on any click (`platform/tray.go:73-75`).
- **Sources** (16 records):
  - Apple: [dock-menus](https://developer.apple.com/design/human-interface-guidelines/dock-menus) 3, 4, 1; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 66, 67, 68
  - Windows 7: [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 5, 33, 39, 41, 47, 51, 52, 34, 42, 43
- **Enforce:**
  - Component: platform.TrayItem gains Default and Checked; the tray menu's order is review.
  - Review: Tray and Dock menu contents; minimise-to-tray offered only as an opt-in.
- **kvit-ui now:** partly. A tray menu line takes `Checked` and `Default`, the default drawn in bold on Windows (`platform/tray.go`, `platform/tray_windows.go`); OnClick exists except on macOS; the menu's contents are the reviewer's.

### CMD-25

Offer printing as deep as the kind of app needs: Print… in the File menu only when there is something to print, through the system's print dialog, showing only the options that apply; a document editor also offers a preview. Deferred: printing is left out until an app needs it (decided, Q-36).

- **Windows:** Print… opens the dialog; a toolbar Print button beside a menu bar prints one copy at once to the default printer, named in its tooltip; without a menu bar Print is a split button offering Print…, Print preview and Page setup. Print preview shows margins, page breaks, orientation, headers and footers, pages through the document and prints from it. Within a session the next job keeps the printer and paper size but not the copies or page range.
- **macOS:** Print… opens the system print panel, whose preview serves; no separate preview command is needed.
- **Sources** (50 records):
  - Apple: [printing](https://developer.apple.com/design/human-interface-guidelines/printing) 1, 4, 5, 6, 7, 8, 9, 10, 11; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 29, 30
  - Windows 7: [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 1, 6, 7, 8, 9, 10, 33, 34, 38, 39, 41, 28, 29, 2, 3, 4, 5, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 27, 31, 32, 35, 36, 43; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 51
- **Enforce:**
  - Review: Print offered where people expect it.
- **kvit-ui now:** does not cover it. No printing in kvit-ui; deferred until an app needs it (decided, Q-36).

## Keyboard

### KEY-1

Everything can be done with the keyboard alone, everything the keyboard does the pointer can do too (typing aside), and the result is the same either way.

- **Sources** (19 records, 2 on Microsoft's checklist):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 23, 25, 22; [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 1; [designing-for-macos](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) 5; [pointing-devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) 3; [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 15
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 46, 45; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 10; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 19; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 5\*, 31; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 5, 44\*; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 17
  - Fluent: [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 16; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 20; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 4
- **Enforce:**
  - Test, one screen: Every control with a Press action is focusable and reachable by Tab (extends CheckNamed).
  - Review: Pointer-only features.
- **kvit-ui now:** enforces it. The gallery's tests tab through every page in both forms and fail a control that takes the keyboard and is never reached (`cmd/kvit-ui-gallery/gallery_test.go`).

### KEY-2

Every standard command the app has answers to the platform's standard shortcut for it, wherever the app's window has focus; when the command does not apply at the moment, the key does nothing, silently.

- **Windows:** Ctrl+N, O, S, P, Z, Ctrl+Y and Ctrl+Shift+Z for Redo, Ctrl+X, C, V, A, F; F3 Find next, Ctrl+H Replace, Ctrl+G Go to, F1 Help (contextual help optional), F2 Rename the selection, F5 Refresh, Delete deletes the selected items, F11 Full screen, Alt+Enter Properties of the focused object, Shift+F10 and the Menu key open the focused item's context menu, Ctrl+W closes the tab or document window, Alt+F4 closes the window; Ctrl+Q is never bound to an app command (decided, Q-19). With document tabs: Ctrl+T opens a new tab, Ctrl+W or Ctrl+F4 closes the current tab, Ctrl+Shift+T reopens closed tabs, most recent first, Ctrl+1 to Ctrl+8 go to the tab in that place and Ctrl+9 to the last. A shortcut acts when its key goes down, not on release.
- **macOS:** Cmd-N, O, S, W, P, Q, Z, Shift-Cmd-Z, X, C, V, A, F, G, comma, ?.
- **Sources** (19 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 5; [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 5; [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 1; [undo-and-redo](https://developer.apple.com/design/human-interface-guidelines/undo-and-redo) 9; [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 3
  - Windows 7: [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 15; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 58; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 6, 49; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 9
  - Fluent: [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 8, 17, 18, 19, 20; [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 6, 8, 10; [ui-controls-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/commanding) 3
- **Enforce:**
  - Component: A Command type knows the standard commands (New, Open, Save, Close, Close tab, Reopen closed tab, Find, Undo, Redo, Rename, Delete, Properties, Help, Settings…) and their keys per platform, registers them on the window, acts on key down, and calls nothing when the command is disabled.
  - Test, all screens of an app: Every app that has Save answers to its shortcut, and so on for each standard command it has.
- **kvit-ui now:** enforces it. A Command with `Std` has its platform's keys, and a window answers them wherever its focus is, leaving typing keys to a focused field and doing Cut, Copy, Paste, Select all, Undo and Redo in a focused control that can (`command.go`, `commandkeys.go`); `uitest.App` fails a standard command made by hand (`uitest/app.go`).

### KEY-3

The command modifier is Command on a Mac and Control elsewhere; Shift adds a related variant; the left and right modifier keys, and the keypad keys, mean the same as their twins.

- **Windows:** Control. Alt+letter and Alt+digit belong to access keys (KEY-15) and are never shortcuts; never Ctrl+Alt, which some layouts read as AltGr for typing characters.
- **macOS:** Command; Option for rare commands; Control avoided.
- **Sources** (9 records):
  - Apple: [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 9, 14, 10
  - Windows 7: [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 61, 64, 57, 55; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 68, 66
- **Enforce:**
  - Component: Shortcuts are declared with unison's `mod.OSMenuCommand()`.
  - Test, Go source: Grep KeyBinding literals: fail a hard-coded `mod.Command`/`mod.Control`, an Alt+letter or Alt+digit binding, and Ctrl+Alt on Windows and Linux.
- **kvit-ui now:** enforces it. A shortcut is written once with `kvitui.Keys` and `ModPrimary`, Ctrl on the Windows form and Command on the macOS form, and keypad keys match their twins (`shortcut.go`); a window refuses Alt+letter and Ctrl+Alt on the Windows form (`commandkeys.go`); `uitest.CheckRules` flags a key binding written with mod.Control or mod.Command (`uitest/wording.go`).

### KEY-4

An app's shortcuts act only while its window has focus: never register a system-wide hotkey, and never take a shortcut or gesture the system or its accessibility features own.

- **Windows:** Every Windows-logo-key combination, including ones the system does not use yet; the Application (menu) key is never a modifier.
- **macOS:** Cmd-Space, Cmd-Tab, Ctrl-F1 to F7, Shift-Cmd-3/4, Opt-Cmd-Esc, Opt-Cmd-D, Ctrl-Space and Cmd-Left/Right for input sources.
- **Sources** (10 records, 3 on Microsoft's checklist):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 24; [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 6; [pointing-devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) 2, 24
  - Windows 7: [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 58, 59, 46\*, 60; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 60\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 79\*
- **Enforce:**
  - Test, Go source: A table of reserved keys per platform; fail a declared shortcut on one, and fail any call to a system-wide hotkey API.
- **kvit-ui now:** enforces it. A window refuses, with a warning, a shortcut on a key the system owns in either form (`commandkeys.go`, `systemKey`); nothing registers a system-wide hotkey.

### KEY-5

Give new shortcuts only to the few frequent app commands: a simple combination (not three keys at once, not keys far apart) on a letter that recalls the command's name, free rather than borrowed from a known meaning, never on punctuation that sits on different keys in other layouts (@ $ { } [ ] \ ~ | ^ ' < >). A shortcut means one thing across the whole app, and an app with many shortcuts lets people reassign them.

- **Windows:** Free for frequent commands: Ctrl+J, K, L, M, R, T, Ctrl+digit, F7–F9, F12, Shift+F2–F5, F7–F9, F11, F12, and Alt+function keys except F4 (not Ctrl+G, which is Go to, nor Ctrl+Q, never bound to an app command; decided, Q-19). In an app with document tabs, Ctrl+T, Ctrl+W, Ctrl+Shift+T and Ctrl+1 to Ctrl+9 belong to the tabs (KEY-2) and are not free. For rarer commands: Ctrl+function key (not F4 or F6) and Ctrl+Shift+letter or digit. Commands on the selection take function keys (F2 renames); commands on the whole document take Ctrl combinations (Ctrl+S).
- **macOS:** The general rule, with the modifiers of KEY-3.
- **Sources** (28 records, 6 on Microsoft's checklist):
  - Apple: [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 3, 4, 7; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 6
  - Windows 7: [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 50, 51, 48, 53, 54, 62, 52, 63, 43\*, 45\*; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 69, 64, 65, 67, 62, 63, 56\*, 59\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 76\*, 78\*
  - Fluent: [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 16; [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 3, 18; [ui-controls-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/commanding) 2
- **Enforce:**
  - Test, all screens of an app: No two commands in an app share a shortcut; a standard key is not bound to another command, and an app with document tabs binds none of the tab keys to anything else.
  - Test, Go source: No shortcut on layout-dependent punctuation or a three-modifier chord.
  - Review: Letters recall the command; reassignment when an app has many shortcuts.
- **kvit-ui now:** partly. A window refuses a second command on a key and another command on a standard key, and never binds Ctrl+Q on the Windows form (`commandkeys.go`); `uitest.App` fails one key for two commands across an app's screens (`uitest/app.go`); punctuation that moves between layouts is not flagged.

### KEY-6

Return or Enter presses the dialog's default button, when it has one, which is drawn as the default; the default is never an action that loses what cannot easily be got back.

- **Windows:** A dialog need not have a default button: give one only when one answer is clearly the likely and safe one. Without one, Enter presses only the focused button. A focused button takes Enter and Space; a focused control that uses Enter itself keeps it. A destructive button may be the default only when its action is easy to undo. A window that appears unprompted has no default button and no access keys, and its initial focus is off the commit buttons. A disclosure opens with Space, a drop-down arrow with Space on its own or Alt+Down inside another control (Windows 7 guide). Enter presses a focused button on this form only (decided, Q-7).
- **macOS:** Return presses only the default button, never a focused one (decided, Q-7); a destructive button is never the default; a dialog that must be read has no default.
- **Sources** (27 records, 4 on Microsoft's checklist):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 14, 15; [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 29, 24, 28
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 168, 169\*; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 31; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 11, 29\*, 30; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 1, 10, 17\*; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 27; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 11, 12, 21; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 26; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 43\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 1
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 25, 26, 28, 16; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 26; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 6
- **Enforce:**
  - Component: Dialog gains an optional default button. On Windows and Linux there is none unless the app marks one (the likely and safe answer, never a Destructive confirm), and Enter or Space presses the focused button; on macOS the confirm button unless Destructive, then none. Button answers Return when focused on Windows and Linux, and when it is the default on macOS (today Space only). Disclosure and SectionHeading drop Return on Windows and Linux.
  - Test, one screen: In every open dialog, Return presses the default the platform form names, or the focused button on Windows and Linux, or nothing.
- **kvit-ui now:** enforces it. Dialog's default is the action on the macOS form unless destructive, and on the Windows form only where `DefaultConfirm` says; Return presses it; a focused Button or IconButton takes Enter on the Windows form; Disclosure and SectionHeading take Space alone on the Windows form (`dialog.go`, `control.go`, `button.go`, `iconbutton.go`).

### KEY-7

Escape cancels a dialog, popover or menu, and abandons a drag, move, resize or split in progress, putting things back; on a Mac Command-period cancels too.

- **Windows:** Escape.
- **macOS:** Escape and Command-period.
- **Sources** (5 records):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 31
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 21; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 50; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 2
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 22
- **Enforce:**
  - Component: Popup OnEscape (exists); add Cmd-. on macOS; Table column drag, SplitView handle and Dialog move go back to where they started on Escape.
  - Test, one screen: Every popup closes on Escape; Escape during a drag leaves the layout as it was.
- **kvit-ui now:** enforces it. Escape closes a dialog, popover, menu and the type-ahead's list, Command-period does on the macOS form, and Escape abandons a table column's move or resize, a split handle's drag and a dialog's move, putting them back (`layer.go`, `typeahead.go`, `windowkeys.go`, `table.go`, `splitview.go`, `dialogmove.go`).

### KEY-8

Focus starts on the control people will most likely use first, never a Help link (in a dialog, the first control of its content, else its default button), and then moves only when the reader acts: when the focused item goes away it moves to a neighbour, a closed dialog or popover gives it back to where it was, and after an input error it goes to the first field in error, scrolled into view with its text selected. While a popover is open, the keyboard stays in it until it closes; Tab closes an open menu (decided, Q-49). Moving focus never changes a value or raises an error, a window never takes focus it was not given, and keyboard focus is always visible.

- **Sources** (23 records, 4 on Microsoft's checklist):
  - Apple: [focus-and-selection](https://developer.apple.com/design/human-interface-guidelines/focus-and-selection) 3, 4
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 43, 44; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 11, 12\*, 18, 9; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 164, 48\*, 47; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 31; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 49\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 69\*; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 15; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 13; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 32, 33
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 27; [ui-controls-item-containers-templates](https://learn.microsoft.com/en-us/windows/apps/design/controls/item-containers-templates) 6; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 30; [ui-controls-command-bar-flyout](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar-flyout) 16; [ui-controls-menus](https://learn.microsoft.com/en-us/windows/apps/design/controls/menus) 6
- **Enforce:**
  - Component: Dialog focuses the first field of its body when it has one, else the confirm button, and finish returns focus to the opener, as Popover does. A form helper moves focus to the first Field with an Error.
  - Test, one screen: After closing each dialog, focus is on the opener; Tab inside an open Popover stays in it; tabbing through every control of a screen changes no value and sets no Error.
- **kvit-ui now:** partly. Dialog puts the keyboard in the first field of its content, else on its default button, and gives it back to the opener; Tab stays inside an open dialog or popover (`dialog.go`, `layer.go`); nothing moves the keyboard to the first field in error.

### KEY-9

Tab moves through controls in reading order; inside a group (radio group, tabs, segmented, list) the arrow keys move and Tab leaves; in a list, tree or table, typing the start of a name jumps to it. Ctrl+Tab and Ctrl+Shift+Tab switch tabs on every system (decided, Q-18).

- **Windows:** A frequently used control may come earlier in the order, and Help links come after the commit buttons. Ctrl+Tab or Ctrl+Page Down goes to the next tab and Ctrl+Shift+Tab or Ctrl+Page Up to the previous, from anywhere in the tabbed view, wrapping at the ends. In a document tab strip the arrows move focus only, without wrapping, and Space or Enter selects; in a view selector row the arrows choose as they move. Radio group: Tab lands on the chosen option, or on the first without choosing it; the arrows choose as they move and stop at the ends (decided, Q-49); Up and Down go by order, Left and Right by position; Ctrl+arrow moves without choosing and Space chooses. Navigation sidebar: the arrows move as in a tree; Space or Enter selects the item and opens or closes its group without moving focus. List: Ctrl+arrow moves focus without selecting, Ctrl+Space toggles the focused item, Shift+arrow extends from the anchor. A row holding buttons: Tab goes through its buttons and then out of the list; Down goes to the next row. A select, open or closed, jumps to the typed letters, ignoring case.
- **macOS:** The general rule.
- **Sources** (32 records, 5 on Microsoft's checklist):
  - Apple: [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 8; [focus-and-selection](https://developer.apple.com/design/human-interface-guidelines/focus-and-selection) 1, 7
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 49\*; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 30; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 33, 14\*, 15\*; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 13; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 24; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 13; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 71\*, 72\*
  - Fluent: [ui-controls-navigationview](https://learn.microsoft.com/en-us/windows/apps/design/controls/navigationview) 22, 23; [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 5, 6, 7; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 21; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 14, 15, 16, 17, 18, 19; [ui-controls-itemsview](https://learn.microsoft.com/en-us/windows/apps/design/controls/itemsview) 6, 9; [ui-controls-selector-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/selector-bar) 10; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 17; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 5; [ui-controls-nested-ui](https://learn.microsoft.com/en-us/windows/apps/design/controls/nested-ui) 9, 11
- **Enforce:**
  - Component: Tab, Segmented and SidebarItem get roving focus with arrow keys, as RadioGroup has; on Windows and Linux RadioGroup stops at its ends and moves without choosing on Ctrl+arrow; the tab row answers to Ctrl+Tab and Ctrl+Shift+Tab on every system, and Ctrl+Page Up/Down on Windows and Linux; Table gets Shift+arrow, Ctrl+arrow and Ctrl+Space; Table, Tree, DualList and Select get prefix type-select.
  - Test, one screen: Tab order matches top-to-bottom, leading-to-trailing order of the controls, Help links after the commit buttons.
- **kvit-ui now:** enforces it. A row of tabs, a segmented control, a radio group and the sidebar are each one stop with the arrows inside; a row of tabs goes round at its ends, and a radio group stops at its ends on the Windows form, where Ctrl and an arrow move without choosing and Space chooses; Ctrl+Tab and Ctrl+Shift+Tab choose tabs on both forms; typing the start of a name jumps in Table, Tree, DualList and Select (`roving.go`, `tab.go`, `segmented.go`, `sidebar.go`, `radiogroup.go`, `windowkeys.go`, `typeselect.go`); in a Table, Shift+arrow extends the selection from the anchor, and on the Windows form Ctrl+arrow moves the keyboard without selecting and Ctrl+Space adds the row to the selection or takes it out (`tablebody.go`). `uitest.Labels` reads the Tab stops in the order Tab visits them and fails a step back up a column or back along a row, a help link after the buttons excepted, comparing each stop where it is seen and leaving out one scrolled out of sight (`uitest/labels.go`).

### KEY-10

Which controls take Tab focus follows the platform: on Windows and Linux every interactive control, read-only fields included; on a Mac it follows the system's full-keyboard-access setting, text fields and lists only when it is off (decided, Q-16, Q-17).

- **Windows:** Every interactive control is a Tab stop, read-only fields included, and a group that acts as one control is one stop (Microsoft's checklist). (Decided, Q-17.)
- **macOS:** Follows the system's full-keyboard-access setting: every control when it is on, text fields and lists only when it is off.
- **Sources** (4 records, 2 on Microsoft's checklist):
  - Apple: [focus-and-selection](https://developer.apple.com/design/human-interface-guidelines/focus-and-selection) 5
  - Windows 7: [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 31; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 70\*; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 13\*
- **Enforce:**
  - Component: Control asks the platform for the full-keyboard-access setting on macOS.
  - Review.
- **kvit-ui now:** enforces it. Every control is a Tab stop on the Windows form; on the macOS form `UI.FullKeyboardAccess`, read from the Mac's setting, keeps Tab to text fields and lists while it is off (`windowkeys.go`, `platform/keyboard_darwin.go`).

### KEY-11

A scrolling view answers to the wheel, the trackpad and the standard keys (arrows, Page Up and Down, Home, End, Space), paging with a small overlap. The wheel scrolls the innermost scrollable view under the pointer, vertically, without a click and without moving focus, faster as it turns faster; tilting scrolls across.

- **Windows:** Ctrl+wheel zooms a zoomable view, up to zoom in; a zoomable view without scroll bars zooms with the wheel alone; a zoomable picture opens fitted to its view; over a tab row the wheel may switch tabs; a wheel-button click may start autoscroll.
- **macOS:** The general rule; the trackpad scrolls both ways and pinch zooms.
- **Sources** (16 records):
  - Apple: [scroll-views](https://developer.apple.com/design/human-interface-guidelines/scroll-views) 1, 5, 2
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 37, 32, 33, 34, 35, 28, 29, 30, 31, 38
  - Fluent: [ui-controls-scroll-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/scroll-controls) 7, 12, 13
- **Enforce:**
  - Component: Region, Table, Tree, DualList and TypeAhead through `wheelScroll` (exists); add Space to Region; Ctrl+wheel zoom where a view zooms.
  - Test, one screen: Page Down in a Region moves by its height less one line; the wheel over a nested list scrolls the list and leaves focus where it was.
- **kvit-ui now:** partly. Region pages on Space and Shift+Space as well as Page Up and Down (`region.go`); there is no Ctrl+wheel zoom.

### KEY-12

Modifiers keep their usual meaning during direct manipulation: the copy modifier while dragging, Shift to keep proportions or constrain, Shift-click for a larger step. No interaction needs two modifiers held together, a drag started by a double-click, or two mouse buttons at once.

- **Windows:** Ctrl copies while dragging; Ctrl+arrow moves or resizes in finer steps; Shift-click opens a link in a new window and Ctrl-click in a new tab; Alt is never a pointer modifier, since it belongs to access keys.
- **macOS:** Option copies while dragging.
- **Sources** (9 records):
  - Apple: [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 8; [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 19; [pointing-devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) 5; [steppers](https://developer.apple.com/design/human-interface-guidelines/steppers) 3
  - Windows 7: [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 56, 65; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 25, 26, 24
- **Enforce:**
  - Review: Drag and resize behaviour.
- **kvit-ui now:** does not cover it. Not covered.

### KEY-13

Going to a search area puts the keyboard in its search field, with its earlier text selected so typing replaces it; Find focuses the window's search field; where the window has a separate find bar, Find opens it (decided, Q-20).

- **Windows:** Ctrl+F and Ctrl+E both go to the window's search field when there is one; where there is a separate find bar, Ctrl+F opens it and Ctrl+E goes to the search field.
- **macOS:** Cmd-F.
- **Sources** (3 records):
  - Apple: [search-fields](https://developer.apple.com/design/human-interface-guidelines/search-fields) 20
  - Windows 7: [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 10, 9
- **Enforce:**
  - Component: SearchField registers itself as the window's Find target, answering Ctrl+F and Ctrl+E on Windows and Linux and Cmd-F on macOS; a find bar, where an app has one, takes Find instead.
  - Test, all screens of an app: Every app screen with a SearchField answers to Find.
- **kvit-ui now:** enforces it. Find's keys, Ctrl+F and Ctrl+E on the Windows form and Cmd-F on the macOS form, put the keyboard in the window's search field with its text selected when no command has them (`windowkeys.go`).

### KEY-14

A shortcut is taught where its command appears: every menu-bar item with a shortcut shows it; context menus show none. A control's tooltip is one string, its verb-led name, with the shortcut added on the Windows form (decided, Q-23).

- **Windows:** All the shortcuts are also listed together on one help page. Every control with a shortcut shows it in its tooltip, in parentheses after the name ("Save (Ctrl+S)"), and a shortcut for a command on a list item is taught in the tooltip of the item or its hover button.
- **macOS:** The general rule (CMD-11); tooltips carry no shortcut.
- **Sources** (6 records):
  - Apple: [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 5; [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 5
  - Windows 7: [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 47; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 61
  - Fluent: [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 9; [ui-controls-tooltips](https://learn.microsoft.com/en-us/windows/apps/design/controls/tooltips) 14
- **Enforce:**
  - Component: With the Command type, a MenuItem shows its command's key, and on Windows and Linux the tooltip of an IconButton, Button or Field adds it in parentheses after the name.
  - Test, all screens of an app: Every command with a shortcut shows it in its menu-bar item, and on Windows and Linux in its toolbar tooltip.
- **kvit-ui now:** enforces it. A menu line made from a Command shows its key, a context menu on the macOS form shows none, and a symbol button's tooltip adds it on the Windows form (`menu.go`, `command.go`); a Field's tooltip writes its shortcut in parentheses on the Windows form only (`field.go`).

### KEY-15

On Windows and Linux every interactive control has an access key: Alt and a letter underlined in its words or its label, unique in the window, which operates the control at once (only focuses it when shared). Every menu and menu item has one, unique within its level. Access keys and shortcuts are separate things (decided, Q-11; the Windows 7 form stands, Q-39).

- **Windows:** None on OK, Cancel and Close (Enter and Escape serve them; a button that means OK or Cancel under other words does get one), links, tabs (Ctrl+Tab), a plain Help button (F1), an ellipsis-only browse button, controls without words, toolbar items, group titles (unless keys run short) or static text. Commit buttons are assigned first with the standard letters (Y, N, the n of Don't, A for Apply), then the first letter or an early, wide, distinctive one (w, m, capitals; x in Exit), avoiding i, l and letters with or beside descenders. About 20 at most per dialog, leaving room for translation. Underlines show while Alt is held, or always when the system asks; Alt alone shows and focuses a hidden menu bar.
- **macOS:** None: macOS has no access keys, and the words are drawn without the marker, as Kvit's menus already do (`menu.go:21-24`).
- **Sources** (92 records, 18 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 51\*, 52\*, 53, 54, 55\*, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 29, 43; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 45; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 41\*; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 9; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 48\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 49\*, 56, 85; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 48\*; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 27; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 47\*, 74\*, 75\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 48, 37; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 2, 32, 40, 41, 42, 43, 44, 45, 46, 47\*, 48\*, 49, 50, 51\*, 52, 53, 54, 55, 56, 57, 58; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 34\*; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 40; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 3, 4, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29\*, 30, 31, 32\*, 34, 35, 36, 37, 38, 39, 40, 41; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 27\*; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 24, 31; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 62\*; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 36; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 27, 41; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 50, 51, 52, 53, 54, 55; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 35
- **Enforce:**
  - Component: The `&` marker of MenuItem (`AccessText`, `menu.go:57-74`) is shared by Button, Check, Switch, RadioGroup options, Disclosure, ChipButton and a visible field or list label; the window routes Alt+letter to the control (controls already ignore keys with Alt, `control.go:153-156`); Dialog's Cancel and close button take none; platform reads the always-underline setting on Windows.
  - Test, one screen: On Windows and Linux, in each window: every eligible control has an access key, access keys are unique among enabled controls, and Link, Tab, IconButton, Cancel and Close have none.
- **kvit-ui now:** partly. Every Button, Check, Switch, radio option, Disclosure and FormRow label has an access key on the Windows form, chosen with "&" or given the first free letter, unique in the window, underlined while Alt is held, and Alt with it operates the control; menu lines get one each; OK, Cancel, Close, links, tabs and symbol buttons take none (`accesskey.go`, `menu.go`); `uitest.Labels` fails two controls sharing a chosen key (`uitest/labels.go`). Where Windows' "Underline access keys" setting is on, `UI.UnderlineAccessKeys` underlines them at all times, on controls and the menu bar's titles (`platform/keyboard_windows.go`, `accesskey.go`, `menubar.go`).

## Pointer

### PTR-1

Clicks mean what they mean everywhere: a primary click selects or activates when the button is released over the target, so moving off abandons it; a secondary click opens the context menu; the wheel scrolls. A double-click is told from two clicks by the system's double-click speed, a double-click where only single clicks make sense counts as one, and clicks made while the app is busy count once.

- **Windows:** Only a customisable toolbar has a context menu, opened on right mouse down; a status-bar area answers only to a single left click and opens its menu on press. Tray icon: acts on release; a left click shows the likeliest view (a flyout, the settings or the main window), a right click the menu with its default command in bold, a double click runs that default. Timeouts scale from the system's double-click speed.
- **macOS:** The status-item menu opens on any click (`platform/tray.go:73-75`).
- **Sources** (25 records):
  - Apple: [pointing-devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) 23, 1; [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 2
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 42, 41; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 32, 25, 27, 28, 31; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 19, 23, 18, 17, 22, 8, 9, 10; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 28; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 22, 23; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 13
  - Fluent: [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 27; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 7; [ui-controls-command-bar-flyout](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar-flyout) 4
- **Enforce:**
  - Component: SetContextMenu (exists); `platform.TrayItem` gains Default, drawn bold; tooltip delay read from the system on Windows.
  - Review: Custom click meanings.
- **kvit-ui now:** partly. Controls act on release over their box (`control.go:122-151`) and Table rows on release over the same row (`tablebody.go:334-351`); right-click, the Menu key and Shift+F10 open the context menu (`window.go:103-109`, `menu.go:714-756`); tray lines have no default (`platform/tray.go:49-60`); the tooltip delay is a fixed 500 ms (`control.go:213`).

### PTR-2

The pointer shape says what will happen: arrow, I-beam over text, pointing hand over links, resize arrows over edges, not-allowed where a drop is refused. The pointer is never hidden, except in full-screen presentation, where it comes back as soon as the mouse moves.

- **Sources** (10 records, 2 on Microsoft's checklist):
  - Apple: [pointing-devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) 25, 26; [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 29
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 39, 3, 4, 7\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 81\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 39; [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 3
- **Enforce:**
  - Component: Components set the cursor (Link, SplitView, Table do).
  - Review.
- **kvit-ui now:** partly. Pointing, Arrow and Resize are used (`link.go:48`, `splitview.go:58-60`); no not-allowed or copy pointer; nothing hides the pointer.

### PTR-3

Support drag and drop where people will try it. Basic tasks never need a drag, double-click, right-click, modifier-click or hover: each has a visible command as well. Buttons that appear only while the pointer is over a row are a shortcut: the same commands are in the row's context menu, reachable by keyboard.

- **Sources** (13 records):
  - Apple: [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 1, 2; [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 19; [split-views](https://developer.apple.com/design/human-interface-guidelines/split-views) 2
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 1, 27, 47; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 29
  - Fluent: [ui-controls-listview-and-gridview](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-and-gridview) 14; [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 4, 11; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 3; [ui-controls-nested-ui](https://learn.microsoft.com/en-us/windows/apps/design/controls/nested-ui) 7
- **Enforce:**
  - Review: Pointer-only features: drags, double-clicks, right-clicks and hover reveals without a visible command; every hover button on a row is also in the row's context menu.
- **kvit-ui now:** partly. DualList is keyboard and buttons only by design (`duallist.go:11-17`); Table drags columns, not rows (`table.go:870-899`).

### PTR-4

A drop within a container moves; a drop into another container copies; when unsure, take the outcome that loses nothing.

- **Sources** (1 records):
  - Apple: [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 3
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. No drag and drop of content.

### PTR-5

While dragging, show what is dragged and how many, keep the pointer captured even outside the object, highlight only the target under the pointer that accepts it, and show a refused drop going back.

- **Sources** (13 records):
  - Apple: [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 10, 12, 13, 14, 15, 28, 20, 21, 16, 9; [scroll-views](https://developer.apple.com/design/human-interface-guidelines/scroll-views) 8
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 20; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 14
- **Enforce:**
  - Component: A drag-and-drop helper in kvit-ui when the first app needs one.
  - Review.
- **kvit-ui now:** partly. EmptyState `Dashed` draws a drop target (`emptystate.go:30-34`); nothing else.

### PTR-6

After a drop, the dropped items are selected where they landed and the drop can be undone; a drop that cannot be undone asks first.

- **Sources** (5 records):
  - Apple: [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 23, 5, 7, 6
  - Windows 7: [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 12
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### PTR-7

Several selected items drag as one; one press-and-drag both selects and drags; dragging from an inactive window leaves its selection alone.

- **Sources** (4 records):
  - Apple: [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 4, 30, 26, 27
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### PTR-8

Fields accept typed, pasted and dropped input; a drop offers and takes the richest format both sides understand.

- **Sources** (7 records):
  - Apple: [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 1, 8; [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 8, 17, 18, 22, 25
- **Enforce:**
  - Review.
- **kvit-ui now:** partly. Fields paste through unison's editor (`field.go:19-21`); dropping is not handled.

### PTR-9

Clicking a column heading sorts by it and clicking again reverses; columns can be resized by dragging their dividers.

- **Sources** (8 records):
  - Apple: [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 12, 13; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 5, 6; [column-views](https://developer.apple.com/design/human-interface-guidelines/column-views) 4
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 9, 33, 35
- **Enforce:**
  - Component: Table (exists).
  - Test, one screen.
- **kvit-ui now:** enforces it. Header press sorts or reverses (`table.go:459-477`); edge drag resizes (`table.go:858-861`).

### PTR-10

Only a node that has, or may have, children shows a disclosure triangle. Clicking it opens that node only; with the modifier it opens all below; a click on a selected name in an editable tree renames, a double click opens.

- **Sources** (6 records):
  - Apple: [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 7, 10
  - Windows 7: [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 14
  - Fluent: [ui-controls-tree-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tree-view) 5, 6, 9
- **Enforce:**
  - Component: Tree.
- **kvit-ui now:** partly. Tree draws the chevron only on nodes with children (`tree.go:389-396`) and expands and collapses by click and keys (`tree.go:299-366`); no expand-all.

### PTR-11

Controls that hide themselves come back when the pointer moves over their area.

- **Sources** (2 records):
  - Apple: [pointing-devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) 4; [going-full-screen](https://developer.apple.com/design/human-interface-guidelines/going-full-screen) 9
- **Enforce:**
  - Review.
- **kvit-ui now:** enforces it. The sidebar rail opens on hover or keyboard focus (`window.go:132-161`).

### PTR-12

Hovering over text that is cut short shows the whole text.

- **Sources** (5 records):
  - Apple: [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 12; [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 14; [typography](https://developer.apple.com/design/human-interface-guidelines/typography) 20
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 13; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 29
- **Enforce:**
  - Component: Label shows its full text as a tooltip when elided, as table cells already do.
  - Test, one screen: Every elided Label has a tooltip with the full text.
- **kvit-ui now:** enforces it. A Label cut short shows its whole text as its tooltip (`label.go`); table cells and FloatingView titles disclose cut text as before.

### PTR-13

A click outside an open popover closes it and goes no further, except on a control that opens another popover, which opens in the same click (decided, Q-43); so that click never presses a button.

- **Sources** (3 records):
  - Apple: [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 9
  - Fluent: [ui-controls-dialogs-and-flyouts-flyouts](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/flyouts) 7, 8
- **Enforce:**
  - Component: Window's popup layer: a press outside a popover closes it and is consumed, unless it lands on a control that opens a popover.
  - Test, one screen: With a popover open, a click on a button closes the popover without pressing the button; a click on another popover's opener opens it.
- **kvit-ui now:** enforces it. A press outside a popup closes it and goes no further, except on a control marked with `OpensPopup`, which opens its own in the same press; Select, MenuButton, SplitButton, Hint, the toolbar's More and the menu bar are marked (`layer.go`).

### PTR-14

The whole visible target is clickable, its symbol and its words alike, and it sits near where it is used; a button only the pointer uses (one shown on hover) is at least 20 by 20 px at 100 % scale; when marks are too small to hit, the whole plot is the target and moving across it reads values.

- **Sources** (5 records):
  - Apple: [charts](https://developer.apple.com/design/human-interface-guidelines/charts) 13
  - Windows 7: [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 14; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 5
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 19; [ui-controls-collection-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/collection-commanding) 12
- **Enforce:**
  - Component: Control takes presses over its whole box (exists); Trend's crosshair (exists).
  - Test, one screen: Every IconButton is at least 20 px square at the default interface size.
- **kvit-ui now:** partly. `uitest.Labels` fails a symbol button smaller than 20 design pixels square (`uitest/labels.go`); a tag's remove button is smaller, inside a 16-pixel tag.

### PTR-15

In a list, table or tree that allows several, a click selects one item, Shift-click extends to a range and the toggle modifier adds or removes one item. A double-click selects the item and runs its default command, which is also on a button or in the context menu.

- **Windows:** Ctrl-click toggles; Shift and Ctrl select the same way with the right button; a Shift- or Ctrl-double-click adjusts the selection and runs the default on all of it. In a single-selection list Ctrl-click or Ctrl+Space clears the selection. The keyboard matches the pointer: Shift+arrow extends from the anchor, Ctrl+Space toggles the focused item, Ctrl+arrow moves focus without selecting (KEY-9). In text a double click selects a word and a triple click the paragraph.
- **macOS:** Command-click toggles.
- **Sources** (18 records):
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 15, 14, 7; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 30, 29, 22; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 21, 20; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 11, 14, 12, 13
  - Fluent: [ui-controls-listview-and-gridview](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-and-gridview) 9, 11; [ui-controls-itemsview](https://learn.microsoft.com/en-us/windows/apps/design/controls/itemsview) 5, 8, 10; [ui-controls-selection-modes](https://learn.microsoft.com/en-us/windows/apps/design/controls/selection-modes) 1
- **Enforce:**
  - Component: Table (exists, through unison's `DiscontiguousSelectionDown`); Tree when an app needs several.
  - Review: Every OnRowActivated or OnActivate has the same command on a button or in the row context menu.
- **kvit-ui now:** partly. Table: Shift-click extends to a range and the toggle modifier's click adds or removes a row, a double press activates, Shift+arrow extends from the anchor, and on the Windows form Ctrl+arrow moves the keyboard without selecting and Ctrl+Space toggles the row (`tablebody.go`); Tree selects one node and runs OnActivate on a double-click of a leaf (`tree.go`); no row context menu built in (survey §12).

## Windows and dialogs

### WIN-1

Dialog buttons follow the platform's order and placement.

- **Windows:** Commit buttons form one right-aligned row at the bottom, even a lone button, never centred. Left to right: the main action, an optional second action (a neutral alternative, such as Try it), then the safe button (Cancel, Close or a specific negative) at the right end; a properties dialog adds Apply after Cancel, and Help, where a dialog has one, comes last.
- **macOS:** Default button at the trailing end, Cancel just before it (leading side of the pair).
- **Sources** (14 records, 4 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 23, 24, 25
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 98, 99\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 44\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 73\*; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 16\*
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 15, 24, 29; [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 9; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 10, 20
- **Enforce:**
  - Component: Dialog lays out its buttons by platform: on macOS Cancel then the confirming button at the trailing end; on Windows and Linux the confirming button first and the safe button at the right end.
  - Test, one screen: On each platform form, the confirming button is where the platform puts it.
- **kvit-ui now:** enforces it. Dialog lays out its buttons by form: the action, the second answer and Cancel from the left on the Windows form; the second answer apart, then Cancel and the action at the trailing end on the macOS form (`dialog.go`).

### WIN-2

An alert asks one question, with a specific title and at most three buttons; a routine or informational message that can recur offers a don't-show-again option, while an error or a necessary warning never does: if it seems to need one, it should not be shown.

- **Windows:** A dialog always has body text under its title (every dialog has a title, decided, Q-42). One or two buttons for a decision; a third only sparingly, its action clearly different from the safe one; further choices become checkboxes or radio buttons committed by one button. The suppression option is worded "Don't show this <item> again" (Don't show this reminder again), or "Don't show this message again" for a dialog in general.
- **macOS:** Explanatory text only when it adds something. The alert's suppression checkbox, "Do not show this message again" unless the app words it.
- **Sources** (25 records, 11 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 7, 13, 36, 17, 35
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 17, 140\*, 139\*, 142\*, 143, 144\*, 145, 146; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 23; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 42\*; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 32\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 43\*, 42\*, 45\*, 46\*; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 20\*
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 7, 14, 21; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 12
- **Enforce:**
  - Component: Dialog gains a suppression option with the platform's wording, refused on a Destructive dialog.
  - Test, one screen: Fail a Dialog with more than three buttons, or, on the Windows form, one with no body text.
  - Review: Errors and warnings offering suppression; one question per alert.
- **kvit-ui now:** partly. Dialog's `Suppressible` offers "Don't show this message again" ("Do not show…" on the macOS form), never on a destructive dialog (`dialog.go`); `uitest.Labels` fails a dialog of more than three buttons (`uitest/labels.go`); a dialog with no body text is not flagged.

### WIN-3

Every dialog has a safe button that leaves without acting (Cancel, Close, or a specific one such as Not now), and a dialog with one button has only that one; Cancel is never the default; a sheet or step with a finishing button (Done on macOS, Close on Windows and Linux, decided, Q-6) also has Cancel or Back.

- **Sources** (7 records, 1 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 27, 28; [sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) 6, 7
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 138\*
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 11; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 21
- **Enforce:**
  - Component: Dialog refuses an empty CancelText (log a warning).
  - Test, one screen: Fail a Dialog without a safe button.
- **kvit-ui now:** enforces it. Every dialog that asks has Cancel unless set otherwise, and one with none is warned about (`dialog.go`); `uitest.Labels` fails a dialog without a safe button (`uitest/labels.go`).

### WIN-4

Mark a confirming button as destructive only when the destruction is a side effect the reader may not expect; a deliberate Empty Trash or Overwrite is not marked.

- **Sources** (3 records):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 26, 38, 39
- **Enforce:**
  - Review: Danger used for chosen deletions.
- **kvit-ui now:** partly. Dialog.Destructive exists (`feedback.go:390-392`); when to use it is not said.

### WIN-5

Use a modal dialog only for a short task that must be answered first; long work, tools used alongside the window and find-and-replace go in a window, a panel or a popover, and a little input an action needs before it runs may be asked in a popover beside its command.

- **Sources** (22 records, 4 on Microsoft's checklist):
  - Apple: [modality](https://developer.apple.com/design/human-interface-guidelines/modality) 1, 2, 3, 5; [sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) 3, 5, 15; [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 1, 2
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 1, 3, 10, 12, 16, 21\*, 24\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 82\*, 83\*
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 1; [ui-controls-dialogs-and-flyouts-index](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/) 2, 6; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 13
- **Enforce:**
  - Review: Dialogs used for long work, for tools used alongside the window, or where a popover would do.
- **kvit-ui now:** enforces it. Dialog's own doc says it is "expensive, and worth it only where the answer really does come first" (catalog KvitDialog).

### WIN-6

A modal view has a title naming its task and an obvious way out inside it, with nothing else looking like the close control; the close control means Cancel or Close, never OK.

- **Windows:** Every dialog has the title bar's close symbol (decided, Q-8) (Kvit draws it beside the title), never disabled except in a progress dialog that must finish to avoid losing data. A dialog with settings also has Cancel in its body, a modeless one Close. The title is the main instruction, a short statement of what to do or decide (a progress dialog's names the operation); it does not name the command, feature or program.
- **macOS:** No dialog, alert or sheet has a close symbol (decided, Q-8); the way out is Cancel or Escape.
- **Sources** (17 records):
  - Apple: [modality](https://developer.apple.com/design/human-interface-guidelines/modality) 6, 8, 4
  - Windows 7: [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 6; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 26, 43, 44, 45, 122, 182, 125, 183, 185; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 29; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 52; [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 42
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 4
- **Enforce:**
  - Component: Dialog (title, Cancel, close, Escape); its close symbol is drawn on Windows and Linux only.
  - Test, one screen: Every dialog has a non-empty Title on every system (decided, Q-42).
- **kvit-ui now:** enforces it. Dialog draws its close symbol on the Windows form only, warns about a dialog without a title, and Escape cancels it (`dialog.go`); `uitest.Labels` fails a dialog without a title (`uitest/labels.go`).

### WIN-7

One modal at a time: close a dialog before opening another, except that an alert may sit over a dialog, and on Windows a helper dialog over its owner; at most one popover, none opened from inside another.

- **Windows:** A dialog may open a helper over itself (a chooser, a settings window it owns); the owner stays open underneath and the helper closes first.
- **macOS:** As stated.
- **Sources** (8 records):
  - Apple: [modality](https://developer.apple.com/design/human-interface-guidelines/modality) 9, 10; [sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) 4; [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 7, 8
  - Windows 7: [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 30, 26; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 28
- **Enforce:**
  - Component: Window.Show refuses a second modal popup other than an alert (or, on Windows and Linux, a helper owned by the top dialog); Popover closes any open popover first.
- **kvit-ui now:** partly. A dialog is not opened, with a warning, while another is open, unless it informs, is an `Alert` (a question about something outside the dialog under it, such as quitting) or, on the Windows form, is a `Helper` (`dialog.go`); a popover does not close another before opening.

### WIN-8

Cancel discards a dialog's or flow's input without asking; closing it would throw away significant effort, or by a route people may take for OK (the close control after many changes), first says so and offers to keep it.

- **Sources** (5 records, 1 on Microsoft's checklist):
  - Apple: [modality](https://developer.apple.com/design/human-interface-guidelines/modality) 7; [progress-indicators](https://developer.apple.com/design/human-interface-guidelines/progress-indicators) 10
  - Windows 7: [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 5; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 35\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 57
- **Enforce:**
  - Component: Dialog gains `OnDiscard`, asked before reject from the close control or Escape when the content holds effortful changes.
  - Review: What counts as significant effort.
- **kvit-ui now:** enforces it. Dialog's `OnDiscard` is asked before the close symbol or Escape throws away what was put in; Cancel discards without asking (`dialog.go`).

### WIN-9

A dialog that belongs to one window blocks only that window and stays with it; the app's other windows stay usable, and bringing a window forward raises only it and the windows it owns, which restore with it.

- **Sources** (4 records):
  - Apple: [sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) 14, 13
  - Windows 7: [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 29, 31
- **Enforce:**
  - Component: Dialog's scrim covers its own window (`layer.go:80-99`).
- **kvit-ui now:** enforces it. The modal scrim is per window (`layer.go:80-99`).

### WIN-10

A popover closes on a click outside it, on Escape and when its window is resized, and the click that closes it does not reach the control under it except on a control that opens another popover, which opens in the same click (decided, Q-43); a single-choice popover closes when the choice is made; Cancel and a finishing button appear only when they clarify whether changes are kept; a tray icon's popup applies its settings at once; a popover never stands in for a tooltip or a context menu.

- **Windows:** A popover with commit and Cancel buttons discards its changes on a click outside; one without them keeps what was entered. The finishing button is Close (decided, Q-6).
- **macOS:** A click outside keeps what was entered. The finishing button is Done (decided, Q-6).
- **Sources** (18 records):
  - Apple: [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 3, 4, 5, 6, 10, 11, 15, 16
  - Windows 7: [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 29
  - Fluent: [ui-controls-color-picker](https://learn.microsoft.com/en-us/windows/apps/design/controls/color-picker) 6, 7; [ui-controls-dialogs-and-flyouts-index](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/) 1; [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 14; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 12; [ui-controls-command-bar-flyout](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar-flyout) 7; [ui-controls-dialogs-and-flyouts-flyouts](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/flyouts) 1, 4, 6
- **Enforce:**
  - Component: Popover (exists) closes on a resize of its window and keeps the closing press from what is under it.
  - Review: Popovers used as tooltips or context menus.
- **kvit-ui now:** partly. Popover closes on Escape and an outside press (`popover.go:114`), but the press still reaches what is under it (`layer.go:31-34`), and a resize does not close it.

### WIN-11

Every main window can be resized, moved and minimised in the system's frame, which shows whether it is active and holds no app controls: Back, the sidebar toggle, search and tabs sit in the window's own header below it (decided, Q-41). It opens at a size that shows its typical content, has a minimum size where its content stops working (Kvit's floor is Px(880) by Px(600); the Windows 7 800×600 floor is not followed, decided, Q-25), drops optional parts below that rather than squeezing them, and its layout works at every size above it.

- **Windows:** Every resizable window can be maximised and restored, to its largest useful size when it gains nothing from more; its default layout aims at 1024×768.
- **macOS:** Zoom and the system's full-screen mode (WIN-15).
- **Sources** (43 records, 5 on Microsoft's checklist):
  - Apple: [designing-for-macos](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) 2; [windows](https://developer.apple.com/design/human-interface-guidelines/windows) 1, 4, 9; [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 10
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 23, 32\*, 19\*, 20\*, 25, 29\*; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 7, 8, 10; [win-window-frames](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-frames) 4, 2; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 59, 60; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 40\*, 173; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 39; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 17; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 44; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 38
  - Fluent: [ui-controls-navigationview](https://learn.microsoft.com/en-us/windows/apps/design/controls/navigationview) 25, 26; [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 1, 3, 5, 6, 7, 9, 10, 12, 20, 21; [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 15; [layout-responsive-design](https://learn.microsoft.com/en-us/windows/apps/design/layout/responsive-design) 7; [ui-controls-title-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/title-bar) 1, 3, 5, 7; [basics-app-silhouette](https://learn.microsoft.com/en-us/windows/apps/design/basics/app-silhouette) 6
- **Enforce:**
  - Component: Window (exists, system frame, minimum Px(880)×Px(600)); full screen (WIN-15).
  - Test, one screen: Screenshots at the minimum and design sizes (exists for pages).
- **kvit-ui now:** partly. System frame, with nothing of the app's in it; minimum size Px(880)×Px(600) (`window.go:88-92`), as decided; no full screen (survey §9).

### WIN-12

A window's title names its content and tells it apart from the app's other windows; a new document is Untitled until named.

- **Windows:** A document window: "document name - program name", with an asterisk before the name while it has unsaved changes (*Report - Kvit Cash, WIN-20); another main window: the program name; a dialog: its main instruction (WIN-6); a properties window "<name> Properties" ("<first>, … Properties", or Selection Properties, for several). A window on the taskbar puts what tells it apart first (66% Complete, 3 Reminders). A window that shows its name prominently inside may hide the caption, but still sets the title.
- **macOS:** The document or place, without the app's name.
- **Sources** (16 records):
  - Apple: [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 10, 11, 12; [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 19
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 57; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 184, 46; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 53; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 77, 78, 81; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 45, 46, 47; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 22
  - Fluent: [ui-controls-title-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/title-bar) 2
- **Enforce:**
  - Test, all screens of an app: No two open windows have the same title.
  - Review.
- **kvit-ui now:** partly. `Window.SetDocument` titles the window "document - app" on the Windows form and with the document alone on the macOS form (`windowstate.go`); two windows with one title are not flagged.

### WIN-13

Open a new window only when seeing two things at once helps; offer Open in New Window where that is useful.

- **Sources** (2 records):
  - Apple: [windows](https://developer.apple.com/design/human-interface-guidelines/windows) 2, 3
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### WIN-14

Bring back the windows, their sizes, positions, views, scroll positions, expanded tree nodes and the last settings pane when the app starts again.

- **Sources** (14 records):
  - Apple: [launching](https://developer.apple.com/design/human-interface-guidelines/launching) 4; [multitasking](https://developer.apple.com/design/human-interface-guidelines/multitasking) 1; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 8; [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 15
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 39, 171; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 12; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 16; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 18; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 19; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 34; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 61, 69; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 17
- **Enforce:**
  - Component: Window saves and restores its frame in settings (`settings/store.go`); Tree and SplitView report state to keep.
  - Review.
- **kvit-ui now:** partly. `Window.RememberFrame` keeps the window's place and size in the settings file and puts it back fitted to the screen's usable area (`windowstate.go`); Dialog and SplitView hand their positions to the app; trees and scroll positions are not kept.

### WIN-15

Offer full screen where people view, present or work in depth; it starts and ends only on the reader's request, never changes the display mode, and keeps the task's controls one pointer move away.

- **Windows:** Full screen is a toggle in View (and on a toolbar) shown checked while on; F11 enters and leaves it, and Esc also leaves when it has no other use. A window that gains nothing from full screen only maximises.
- **macOS:** The system's full-screen mode, Ctrl-Cmd-F.
- **Sources** (19 records):
  - Apple: [going-full-screen](https://developer.apple.com/design/human-interface-guidelines/going-full-screen) 1, 2, 3, 4, 5, 7, 8, 10, 12, 13, 14; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 51; [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 20
  - Windows 7: [win-window-frames](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-frames) 5, 6, 8, 1, 7; [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 19
- **Enforce:**
  - Component: Window.FullScreen through unison.
  - Review.
- **kvit-ui now:** partly. `Window.SetFullScreen` and `FullScreenCommand`, the standard command on F11 and Ctrl-Cmd-F, make the window cover the screen and put it back (`windowstate.go`); unison has no full-screen mode, so the window's frame stays drawn.

### WIN-16

Settings hold app-wide options people seldom change, and every app setting can be reached from them; everyday commands stay in the working view, options for one view live in that view and options for one document in the File menu; settings that are complex or break things when wrong offer Restore defaults.

- **Windows:** Settings (never Options or Preferences) is a page in the app's window, filling it as one scrolling column rather than a dialog; each change applies at once, with no OK, Cancel or Apply (WIN-23). Its command is the last item of the navigation sidebar, below footer items such as Help; or the last item of the toolbar's More menu, or on the toolbar; with a menu bar and no sidebar, the last item of Tools. About is a section at its bottom (CMD-3). There is no separate Customize command; an object's attributes are its Properties (WIN-17). Every option of a tray app not on its tray menu is in Settings, including Display icon in notification area.
- **macOS:** Settings… in the app menu: a window with a toolbar of panes, changes taking effect at once, no minimise or zoom, titled with the pane's name.
- **Sources** (33 records):
  - Apple: [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 6, 7, 12, 13, 14; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 12
  - Windows 7: [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 65, 8, 12, 73, 83, 2, 55, 67; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 48, 49; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 53, 88; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 49; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 74, 91
  - Fluent: [ui-controls-navigationview](https://learn.microsoft.com/en-us/windows/apps/design/controls/navigationview) 7, 4; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 2, 3, 5, 8, 11, 12, 14, 41, 1; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 1
- **Enforce:**
  - Component: A Settings component with a platform form (a page in the app's window, ending with About, on Windows and Linux; a window of panes on macOS), applying each change at once on every platform.
  - Review: Everyday commands kept off the settings page.
- **kvit-ui now:** partly. SettingRow, Settings at the sidebar's foot and `NewAboutSection` build a settings page that applies at once (`settings.go`, `sidebar.go`); the same `SettingsPane`s make both forms: on the Windows form `NewSettingsPage` (or `SettingsWindow.Page`) is one scrolling column with each pane a section under its name and About last, and on the macOS form `SettingsWindow` is a window of panes, titled with the pane shown, as tall as it and not resizable; `SettingsWindow.Command` opens the window on the macOS form and shows the page on the Windows form (`settingswindow.go`). unison cannot take the minimise button off the window, which the macOS form does not have.

### WIN-17

An inspector that follows the selection is a pane or panel with direct-adjustment controls that apply at once; panels are small, never cover what they adjust, and come and go with the app; an object's own properties window holds mostly its settings, shows its name first, and for a container shows the container's own properties.

- **Windows:** The inspector pane is titled Properties; an object's properties window is modeless; a modeless tool window may dock into the main window.
- **macOS:** An inspector panel or pane.
- **Sources** (14 records):
  - Apple: [panels](https://developer.apple.com/design/human-interface-guidelines/panels) 2, 3, 4, 6, 8, 14
  - Windows 7: [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 4, 11, 21, 82, 5, 6; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 27, 2
- **Enforce:**
  - Component: Pane (exists).
  - Review.
- **kvit-ui now:** enforces it. Pane is the fixed side pane of detail about the selection (catalog KvitPane).

### WIN-18

Open and save files through the system's dialogs, which bring name, folder, format and search; start in a sensible folder.

- **Sources** (11 records):
  - Apple: [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 15, 18, 20, 21, 3, 16, 17; [searching](https://developer.apple.com/design/human-interface-guidelines/searching) 10; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 24, 25
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 6
- **Enforce:**
  - Component: platform.FileDialog (exists for open); add save.
  - Test, Go source: Flag a hand-made file chooser.
- **kvit-ui now:** partly. `platform/filedialog.go` opens files with the system dialog; there is no save dialog.

### WIN-19

Never lose work: notes-like data (kvit-notes, kvit-works) saves automatically; kvit-cash asks before closing a window or quitting with unsaved changes, with Save, Don't save and Cancel, where Cancel abandons the close or quit, and has Save on its standard shortcut (decided, Q-12).

- **Sources** (8 records):
  - Apple: [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 4, 23; [the-menu-bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) 23; [multitasking](https://developer.apple.com/design/human-interface-guidelines/multitasking) 1
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 129; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 21
  - Fluent: [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 25; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 16
- **Enforce:**
  - Component: Window gains `AllowClose` wired to unison's `AllowCloseCallback`, for autosaving and asking apps alike; MenuBar's Quit goes through it.
  - Review: Which data saves automatically and which asks.
- **kvit-ui now:** enforces it. `Window.AllowClose` is asked before the window closes, by its close button, Close, Exit or Quit (`window.go`).

### WIN-20

Mark a document with unsaved changes through Window.Edited, and clear the mark when saved.

- **Windows:** An asterisk before the document name in the window title (*Report - Kvit Cash) (decided, Q-13).
- **macOS:** The system's dot in the close button.
- **Sources** (3 records):
  - Apple: [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 24, 25, 26
- **Enforce:**
  - Component: Window.Edited: the title's asterisk on Windows and Linux, the close button's dot on macOS.
- **kvit-ui now:** partly. `Window.SetEdited` puts an asterisk before the document in the title on the Windows form, and on the macOS form asks AppKit for the dot in the close button, finding the window by its title since unison keeps the system's window to itself (`windowstate.go`, `platform/document_darwin.go`); the dot is compiled for macOS but has not been run on a Mac.

### WIN-22

A guided flow can always be skipped, cancelled or left, and started again; Back returns to the page last viewed, stays available unless repeating a step would cause harm, never closes the flow and keeps everything entered; moving forward never needs Back, and going Back after committing lets the reader review, committing again only if something changed.

- **Sources** (19 records, 4 on Microsoft's checklist):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 6; [sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) 2; [onboarding](https://developer.apple.com/design/human-interface-guidelines/onboarding) 5
  - Windows 7: [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 8, 34, 36, 64\*, 65\*, 66, 67, 6, 7, 25, 26, 27, 28, 29; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 92\*, 98\*
- **Enforce:**
  - Component: A Flow component, when one exists, keeps each page's input and its own history for Back.
  - Review.
- **kvit-ui now:** partly. Spotlight dismisses on Escape and outside press (`spotlight.go:33`, `spotlight.go:81`); there is no flow component.

### WIN-23

The confirming button applies the dialog's changes, does its task and closes; Cancel, Escape and the close control discard every change, nested dialogs' included, and leave no trace; settings, a modeless window and an inspector apply each change at once and have no OK or Cancel (a modeless window has Close); Apply is left out, not disabled, wherever the platform form does not give it.

- **Windows:** OK, Cancel and Apply remain only for a properties dialog of one object, Apply enabled once something changed so the changes can be seen before closing. A properties dialog opened from another keeps its changes on OK but they take effect only when the owner is confirmed, and the owner's Cancel undoes them. A dialog of input, progress and results pages offers Cancel on the first two and Close on the last.
- **macOS:** Settings windows and inspectors apply each change at once; there is no Apply.
- **Sources** (20 records, 5 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 22, 25, 34, 95, 116, 121\*, 128\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 7, 27, 46, 47, 51, 52, 53, 54\*, 58; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 37\*; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 81\*
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 10; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 2
- **Enforce:**
  - Component: Dialog keeps confirm, Cancel, close and Escape to these meanings (exists); a properties form of Dialog adds Apply on Windows and Linux, enabled once something changed.
  - Test, one screen: Fail an Apply button outside a properties dialog, a disabled Apply, or OK or Cancel on a settings page.
  - Review: Cancel leaving no side effect; modeless windows using Close.
- **kvit-ui now:** enforces it. `DialogProperties` adds Apply on the Windows form only, disabled until `MarkChanged` says something changed and again after it is pressed, which leaves the dialog open; the macOS form leaves it out (`dialog.go`). Settings apply at once and have no OK, Cancel or Apply (`settingswindow.go`). `uitest.Labels` fails an Apply outside a properties dialog or on the macOS form, and OK, Cancel or Apply in a region holding setting rows (`uitest/labels.go`). Whether Cancel leaves no trace, and whether a modeless window uses Close, is the reviewer's.

### WIN-24

Use a guided flow only for one task whose steps must come in order and that a dialog, a pane or a single page cannot carry. Keep only pages that ask for real decisions (no welcome page), branch little, put the steps on pages and only helpers such as choosers in dialogs, commit once on a page whose button names the task, and show a summary only for risky or complex input and a follow-up only for likely next tasks. One input step with progress and results is a dialog of up to three pages without Back and Next; a one-page flow is a dialog.

- **Sources** (32 records, 12 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 7, 29, 30, 31, 32; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 1, 2, 3, 5, 12\*, 13, 20, 21, 22, 23, 24, 30, 32, 42, 43, 47\*, 49\*, 51\*, 57\*, 60\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 1; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 90\*, 99\*, 100\*, 102\*, 103\*, 105\*
- **Enforce:**
  - Component: A Flow component, when one exists, renames the commit page's button to the task and offers no welcome page.
  - Review: Whether a flow is needed; page count, branches, summary and follow-up pages.
- **kvit-ui now:** does not cover it. No flow component in kvit-ui; Dialog is a single page (`feedback.go:368-575`).

### WIN-25

A window or dialog opens wholly on screen, inside the area the taskbar or Dock leaves, shrunk to fit a smaller screen: a dialog centred over its owner and always above it, a window opened from an object or the tray next to it without covering it, a new main window cascaded from the top leading corner. A window that appears unasked opens away from the pointer and takes no keyboard focus, and no window stays always on top.

- **Windows:** A window needing attention flashes its taskbar button three times and stays highlighted, with no sound; a progress window goes to the bottom trailing corner; a stand-alone utility opens centred; a saved position that no longer fits uses the last monitor, shrunk and shifted, else the default placement.
- **macOS:** The system places windows and draws attention through the Dock icon.
- **Sources** (24 records, 6 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 35\*, 37, 38, 41; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 27, 28; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 18, 19\*, 21, 22, 23, 24, 25\*, 26, 27, 28\*, 30, 35; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 5\*, 6\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 17; [inter-keyboard](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-keyboard) 7; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 15, 30
- **Enforce:**
  - Component: Window opens at the design size only when it fits the screen's work area, otherwise shrunk, and cascaded; a restored frame (WIN-14) is clamped to the monitors; Window.RequestAttention flashes the taskbar button instead of taking focus.
  - Test, all screens of an app: On a 1024×768 screen every window opens wholly visible.
- **kvit-ui now:** partly. Dialog opens centred in its window and stays inside it when moved (`dialogmove.go:41-62`), in the window's popup layer above its content (`layer.go:58-68`); NewWindow sets Px(1440)×Px(960) with no fitting to the screen (`window.go:117`).

## Messages

### MSG-1

Show a message only when the reader will act on it or must know it now; otherwise fix the problem silently, give the nearest result, or write it to a log. Background work the reader did not start needs no message. What is shown goes where its urgency puts it: inline, then a notice, then a dialog, escalating only as a condition becomes urgent, once per condition; information shown every few minutes gets its own place in the window. A notice (an in-window bar) carries long-lived app-wide status without interrupting: it takes its own space and covers nothing, has at most one action, closes when its condition is resolved or by the reader (no close control while it must stay until resolved), is never shown and withdrawn in rapid turns, and is announced to screen readers when it appears. A notice never confirms an action that changed nothing; success the reader just caused is a toast, and a success notice is only for background work that ended while the reader was elsewhere (decided, Q-47). Messages make no sound unless the reader must respond at once.

- **Windows:** A dialog that needs attention while the app is in the background flashes the app's taskbar button three times and leaves it highlighted, without taking focus or making a sound.
- **macOS:** A background app bounces its Dock icon once for a dialog that needs attention (critical: until the app is active).
- **Sources** (62 records, 10 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 1, 2, 6; [feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) 1, 4; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 20; [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 13; [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 12
  - Windows 7: [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 5, 6, 24, 22; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 20, 5, 8\*, 163; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 7, 18, 19\*, 3\*, 4, 14, 15; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 2, 3\*, 7, 8, 11, 12, 13, 15, 17, 53\*, 21, 30; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 114\*, 116\*, 66\*, 67\*, 108\*; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 14, 3; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 3, 7
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 1, 2, 3, 5, 6, 7, 8, 10, 11, 13, 14, 16, 23; [ui-controls-dialogs-and-flyouts-index](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/) 4, 9; [ui-controls-progress-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/progress-controls) 10; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 9; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 9
- **Enforce:**
  - Review: A message the reader would dismiss without acting; a problem the app could fix itself; administrator or debug detail on screen; a warning before the condition is near; a sound with a message.
  - Component: Window gains RequestAttention, mapped to FlashWindowEx (Windows), the urgency hint (Linux) and requestUserAttention (macOS), called by Dialog.Open when the window is not active.
  - Test, Go source: Fail a Go source that plays a system sound or beep alongside a Dialog, Notice or Toast.
  - Component: Notice gains Shown, which announces it to a screen reader as Toast's does.
  - Review: A notice toggled by a flickering condition; a notice confirming a no-change action; a notice laid over content.
- **kvit-ui now:** partly. Notice, Toast and Dialog are separated by meaning (`feedback.go:218-223`, catalog KvitNotice); Notice has one Action (`feedback.go:231-235`), a close control only when Dismissible (`feedback.go:236-237`, `feedback.go:275`) and its own room in the layout (`feedback.go:281-337`); it is not announced when shown, only named (`feedback.go:356-367`), unlike Toast (`feedback.go:117-120`); nothing stops rapid toggling. No library message plays a sound; nothing asks for attention when a dialog opens in an inactive window (no FlashWindow or equivalent anywhere in kvit-ui).

### MSG-2

Show no dialog when the app starts; news found at start goes in the window.

- **Sources** (3 records):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 5; [onboarding](https://developer.apple.com/design/human-interface-guidelines/onboarding) 12, 13
- **Enforce:**
  - Test, one screen: No dialog is open on the first frame of any app screen.
- **kvit-ui now:** does not cover it. Not checked.

### MSG-3

Do not ask before an action that can be undone, is frequent, or whose loss is the obvious point; take the likely safe choice and make it easy to reverse. Ask only when going on needs a choice with no good default, the consequences may be unknown to the reader, the action is probably a mistake in its context, or it is a rare action that loses data for good, has security consequences, spends money (unless the reader has turned that confirmation off), or commits the reader, such as submitting a sign-up. A confirmation people always answer the same way is removed, not given a hide option; it never teaches or promotes.

- **Sources** (38 records, 10 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 3, 4; [feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) 5, 6; [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 8; [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 6
  - Windows 7: [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 2\*, 3, 6\*, 7\*, 8, 9, 13, 16\*, 17\*, 18, 19, 20, 33\*, 4\*, 11; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 66; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 4\*, 96; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 6\*, 8, 9, 10; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 117\*
  - Fluent: [ui-controls-dialogs-and-flyouts-index](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/) 3, 5, 10; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 14, 15, 17, 18, 19; [ui-controls-selection-modes](https://learn.microsoft.com/en-us/windows/apps/design/controls/selection-modes) 6
- **Enforce:**
  - Component: ConfirmInPlace and Toast with Undo (exist).
  - Review: A confirmation before an undoable or frequent action; a question that is almost always answered the same way; a question carrying a tip or an advertisement.
- **kvit-ui now:** enforces it. ConfirmInPlace is for undoable actions (`confirminplace.go:14-23`); Dialog's own doc sends reversible confirmations to a toast with an undo (`feedback.go:368-371`).

### MSG-4

A message's title (its main instruction) says what happened, or asks the one question, specifically, in at most two lines; never a bare Error or a code; calm, and never making the reader the subject of a mistake. Every dialog has a title on every system (decided, Q-42), and no caption names the command, feature or program (decided, Q-9). Problem messages are worded to suit every system ("Couldn't save the file"), never with "we" (decided, Q-15). An error code appears only when the message cannot be specific or the code helps a search, in the details, as "Error code: 0x80070005".

- **Sources** (30 records, 3 on Microsoft's checklist):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 8, 9, 10, 11
  - Windows 7: [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 66, 48, 50, 51, 52, 9, 18, 49, 58\*, 68, 69, 72, 75; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 42, 43, 1, 46, 44, 47; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 111\*, 110\*; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 34; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 45; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 23
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 5, 3
- **Enforce:**
  - Test, one screen: Fail a Dialog or Notice title of "Error", "Warning", a code alone, or over about 100 characters, and one that starts with "You " in a Notice or Dialog of danger tone.
  - Review: A confirmation's title asks more than one question; an error code outside the details; "we" in a message.
- **kvit-ui now:** partly. `uitest.Labels` fails a dialog or notice titled Error, Warning, Caution or Alert, or over 100 characters (`uitest/labels.go`); a title of a code alone or starting "You " is not flagged.

### MSG-5

Explanatory text adds what the title does not: for a problem, the likely cause, what happens next, and what the reader can do, in two or three steps at most (longer procedures go to help); for a warning, also what happens if nothing is done; for a question, what is needed to choose each answer. Say only what is likely true, and plainly say when the cause is unknown. One message per detectable cause. Never tell the reader to contact support, and an administrator only when only one can fix it. Name a file by its short name in the title and give the full path here. A details disclosure only when it reveals something new. The text does not describe the buttons; when it must name one, it uses the button's exact words.

- **Sources** (36 records):
  - Apple: [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 13, 14, 15, 16
  - Windows 7: [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 12, 13, 17, 37; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 10, 19, 23, 33, 34, 40, 59, 70, 71, 80, 83, 84, 76, 79, 81; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 15, 37, 52, 53; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 196, 198; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 41; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 56, 57; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 27
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 9; [style-writing-style](https://learn.microsoft.com/en-us/windows/apps/design/style/writing-style) 3, 22
- **Enforce:**
  - Review: Text that repeats the title; a fix that is a guess; "contact technical support"; a full path in a title; details that reword the text.
- **kvit-ui now:** does not cover it. Not checked; Dialog and Notice have a free Detail (`feedback.go:385-386`, `feedback.go:227-228`) and no details disclosure.

### MSG-6

Report an input problem next to its cause, saying how to put it right: a character or length problem as it is typed, a problem with the whole value (including a required field left empty) when the field is left, one message at a time, at the field only when that field alone is the cause. Warn as it happens of a probably unintended state, such as Caps Lock on in a password field or typing past the length limit. Problems found when a form is submitted are marked at each field in error, in words, and the form stays open to be corrected and submitted again; a problem tied to no field (a lost connection) gets its own message, never a dialog for a field's problem; a guided flow that cannot finish ends on a page saying why and what to do.

- **Sources** (50 records, 2 on Microsoft's checklist):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 18, 24; [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 9; [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 9, 10; [feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) 8; [progress-indicators](https://developer.apple.com/design/human-interface-guidelines/progress-indicators) 4
  - Windows 7: [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 40; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 27, 26, 29, 32, 77\*, 78; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 32, 33, 34, 35; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 1, 12, 4, 8, 10, 13, 17, 18, 19, 20, 21, 2, 5, 9, 11, 16, 26, 28, 29, 31; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 44, 61, 46; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 11, 160, 162; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 5; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 109\*; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 22
  - Fluent: [ui-controls-dialogs-and-flyouts-dialogs](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/dialogs) 17; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 14; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 8
- **Enforce:**
  - Component: Field.Error (exists); NumberField checks the range on leaving the field, not on every keystroke; Field gains a Caps Lock note for masked fields and a note at its length limit.
  - Review: The text says how to fix it; one field blamed for a problem of two; a dialog listing problems found on submit.
- **kvit-ui now:** partly. Field shows Error under the field (`field.go:227-228`, `field.go:386-395`); NumberField runs its whole check, range included, on every change (`numberfield.go:46`, `numberfield.go:101-120`), so "Must be at least 100" shows while 150 is still being typed; "Not a number" (`numberfield.go:109`) says what is wrong, not how to fix it; no Caps Lock or length-limit note.

### MSG-7

Confirm success only for significant tasks, by a toast when the reader just caused it (decided, Q-47); report ordinary actions only when they fail. Background work longer than about two minutes whose end matters more than its progress reports its end by a notification, and a guided flow ends on a completion page. A completion sound only for rare, important work long enough for people to walk away.

- **Windows:** Windows may give on-demand success feedback through a tray icon shown only while the work runs (MSG-17).
- **Sources** (15 records, 6 on Microsoft's checklist):
  - Apple: [feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) 7; [multitasking](https://developer.apple.com/design/human-interface-guidelines/multitasking) 5
  - Windows 7: [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 13, 12; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 3, 14; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 59\*, 56\*, 58\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 107\*, 104\*, 106\*; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 77
  - Fluent: [ui-controls-progress-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/progress-controls) 12; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 10
- **Enforce:**
  - Review: Toasts for routine success; a long export with no word when it ends.
- **kvit-ui now:** partly. Toast exists for acknowledgements (`feedback.go:67-91`); Tray.Notify exists for the end of background work (`platform/tray.go:269-283`).

### MSG-8

Nothing the reader needs closes itself on a timer: a toast closes after 4 s only when it acknowledges something nothing later needs, and a field's message stays until the value is fixed. Tooltip timing comes from the system where it has a setting (Windows), Kvit's own formula elsewhere; a tooltip is never the only place a control's meaning lives (HLP-1) (decided, Q-21).

- **Sources** (8 records):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 27
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 161; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 16, 17; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 28; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 14, 15
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 4
- **Enforce:**
  - Component: Toast keeps its rule that one with an action stays and closes after 4 s otherwise (exists); Field.Error stays until the application clears it (exists); tooltips on Windows take the system's hover and show times.
- **kvit-ui now:** partly. A Toast with an action stays; one without closes after 4 s (`feedback.go:67-91`, `feedback.go:117-131`); Field's message stays while Error is set (`field.go:386-395`).

### MSG-9

An empty view says what would be there and gives the action that fills it; it is never the only place for information needed later.

- **Sources** (4 records):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 16, 17
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 11; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 14
- **Enforce:**
  - Component: EmptyState (exists).
  - Test, one screen.
- **kvit-ui now:** enforces it. EmptyState has Title, Detail and Action (`emptystate.go:22-55`).

### MSG-10

A desktop notification is for something that happened in the background, not the result of what the reader just did: one fact, once, silently, never private data, never continuously changing information such as progress, never an advertisement on any system (decided, Q-31), and none while the app's window is in front. Notify only what may change what the reader does. Assume it may go unseen: nothing, first-run information included, depends on it, and what must not be ignored is not a notification.

- **Windows:** An ignorable failure of background work unrelated to the current task (a sync that failed) may be an action-failure notification with the error symbol, and a problem that matters while other programs are in use is shown on the tray icon (MSG-17). Windows holds notifications back for a new system's first hours.
- **macOS:** No errors in notifications; a failure is shown in the window and, if needed, as a badge (MSG-14).
- **Sources** (22 records, 1 on Microsoft's checklist):
  - Apple: [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 1, 2, 3, 4, 5, 6
  - Windows 7: [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 4, 5; [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 27; [vis-std-icons](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-std-icons) 9; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 11, 1, 3, 4, 7, 8, 14, 23, 16, 49; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 78; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 65\*
- **Enforce:**
  - Component: Tray.Notify skips when a Kvit window is active, or the app decides (`platform/tray.go:269-283`); the Windows backend sets NIIF_NOSOUND.
  - Review: A notification for the result of the reader's own command; progress in a notification.
- **kvit-ui now:** enforces it. `UI.Notify` posts nothing while one of the app's windows is in front (`notify.go`); the Windows balloon is posted with NIIF_NOSOUND and Linux's with suppress-sound (`platform/tray_windows.go`, `platform/tray_linux.go`).

### MSG-11

A notification's title is short and specific without the app's name; the body is complete sentences; the system shortens it, not the app.

- **Windows:** Title in sentence case, within 48 English characters (Windows cuts it at 63) and a body within 200 (cut at 255), so translations fit.
- **macOS:** Title in title case.
- **Sources** (13 records):
  - Apple: [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 7, 8, 9, 10, 11, 12
  - Windows 7: [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 39, 40, 44, 37, 38, 42; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 24
- **Enforce:**
  - Component: Tray.Notify cases nothing itself; a test of the app's notification strings checks the case for the platform and the Windows lengths.
  - Review.
- **kvit-ui now:** partly. The Windows backend cuts the title at 63 and the body at 255 UTF-16 units without an ellipsis (`platform/tray_windows.go:373-380`, `platform/tray_windows.go:132-134`); nothing checks case or length.

### MSG-12

Give each notification an honest kind and urgency, and the timing that goes with it: success, failure and information once, when the work ends, and dropped if it cannot be shown then; a problem the reader must solve when it starts, then at most daily, withdrawn once fixed; an optional task at most once a day, three times in all. Hold back a failure that may fix itself within seconds, and avoid the burst at logon. People choose which kinds they get. Never advertise features or offers in a notification, on any system (decided, Q-31).

- **Windows:** Informational (FYI) notifications stay off until the reader turns them on, from the tray icon's menu.
- **Sources** (21 records, 3 on Microsoft's checklist):
  - Apple: [managing-notifications](https://developer.apple.com/design/human-interface-guidelines/managing-notifications) 1, 4, 5, 6, 7, 8, 9
  - Windows 7: [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 26\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 68\*; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 9\*, 10, 17, 18, 19, 20, 21, 35, 15, 50; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 40, 6
- **Enforce:**
  - Component: Notification gains a Kind (success, failure, event, task, information) and an Urgency mapped per platform, and Tray gains Withdraw(id) for a notification that no longer applies.
  - Review: A notification repeated more often than its kind allows; an advertisement.
- **kvit-ui now:** partly. A notification has a `Kind` (information, success, failure, problem, task), shown with the matching symbol on Windows and urgency on Linux, and `Tray.Withdraw` takes one that no longer applies off the screen by its ID: the balloon on Windows, `CloseNotification` on Linux, and Notification Centre's delivered and pending lists on macOS (`platform/tray.go`, `platform/tray_windows.go`, `platform/tray_linux.go`, `platform/tray_darwin.go`). Informational notifications are not held back until the reader turns them on, and nothing limits how often a kind is posted.

### MSG-13

Clicking a notification opens the window where the reader acts or reads more, and shows more than the notification said. Notification buttons, where the platform has them, offer quick common tasks, never Open, and avoid destructive ones.

- **Windows:** No buttons: the notification is clicked as a whole and opens a window for the action, never performing it; a clickable one ends its body with "Click to <do something>." or "Click for more information."
- **macOS:** Buttons for quick tasks may act without opening the app.
- **Sources** (9 records):
  - Apple: [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 15, 16, 17, 18, 13
  - Windows 7: [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 25, 26, 48, 27
- **Enforce:**
  - Component: Tray.OnNotificationClick (exists) brings the app's window forward; buttons only in the macOS backend.
  - Review: A click that rephrases the notification.
- **kvit-ui now:** partly. A click comes back to OnNotificationClick with the notification's ID (`platform/tray.go:76-78`); notifications have no buttons (`platform/tray.go:62-69`); the Linux backend offers only the default action, labelled Open (`platform/tray_linux.go:385-386`).

### MSG-14

A badge on the app's icon counts only unread items, falls as they are read, and its news is also in the window. Inside the window a badge on a control (a sidebar item, a tab, a button) draws attention there without interrupting: a count of items needing attention, or a dot or symbol for new or changed status. It is announced by the control it sits on, has a stated rule for when it appears and goes, and clears when seen. Never a badge for a critical error or anything needing action now (a dialog), for an alert that cannot be dismissed (a notice), or as a plain symbol. One kind of badge throughout a sidebar; a collapsed parent item shows the sum of its children's counts.

- **Sources** (16 records):
  - Apple: [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 20, 21, 22, 23
  - Fluent: [ui-controls-info-badge](https://learn.microsoft.com/en-us/windows/apps/design/controls/info-badge) 1, 2, 3, 4, 5, 6, 7, 8, 9, 13, 14, 15
- **Enforce:**
  - Component: Badge (exists) gains a dot form; SidebarItem's count (exists) keeps folding into its name.
  - Review: A badge that never clears; a badge for a critical error; mixed badge kinds in one sidebar.
- **kvit-ui now:** partly. Badge counts, hides at zero and caps at 99+ (`badge.go:24-48`, `badge.go:59-71`); SidebarItem.Count draws one and puts the count in the item's name (`sidebar.go:93-109`, `sidebar.go:144-150`, `sidebar.go:222-232`); Tab shows a count in words and name (`tab.go:22-24`, `tab.go:111-122`); no dot or symbol badge, no nested sidebar items, no app-icon badge.

### MSG-15

Ask for access to private data at the first use that needs it, saying why; do not ask for what the app can find itself.

- **Sources** (3 records):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 8; [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 2
  - Windows 7: [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 18
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### MSG-16

Choose a message's type by what the reader must do: an error for a problem that has happened and blocks them, a warning for a condition that may cause trouble later and lets them go on, information otherwise; success the reader just caused is a toast, and an in-window notice is a success only for background work that ended while the reader was elsewhere (decided, Q-47). Its symbol and words agree with that type: never an error softened into a warning, a question with the warning symbol only when its answer risks losing data, money, access or much time, no symbol on routine questions, minor input problems or ordinary information. An error shown in place beside a field does carry the error symbol so it is not overlooked; a minor problem at the end of a flow is one more task, without error words.

- **Windows:** A notification's symbol follows its kind: error for an action failure, warning for a system event, the subject's own symbol otherwise, never the information symbol (MSG-10).
- **Sources** (38 records, 23 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 86\*; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 1\*, 2, 11\*; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 1\*, 20\*, 36\*, 38\*; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 23\*, 24\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 62; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 115\*, 120\*, 121\*, 122\*, 123\*, 124\*; [vis-std-icons](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-std-icons) 1, 2\*, 3, 4\*, 5\*, 6\*, 8\*, 10, 11, 12\*, 13\*, 14, 15, 16; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 28, 29, 30, 33; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 22\*
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 18, 19
- **Enforce:**
  - Component: Dialog gains a Tone (none by default) that draws the symbol beside the title; Field's message gains the small error symbol; Toast and Notice keep Tone and drop the info symbol where the tone is neutral.
  - Test, one screen: Fail a Dialog with the warning tone whose buttons do not include a Destructive confirm or a Cancel, and a Field in error without its symbol.
  - Review: Type chosen by severity instead of by what happened.
- **kvit-ui now:** partly. Toast and Notice pick a symbol by Tone (`feedback.go:46-57`), and show the info symbol even for neutral and info tones; Dialog has no symbol or tone (`feedback.go:380-400`); Field's message is a caption in danger ink with no symbol (`field.go:227-228`); the Windows notification balloon always carries the information symbol (`platform/tray_windows.go:439`).

### MSG-17

A tray icon (notification area on Windows and Linux, menu bar extra on macOS) is for what has no window, or status that matters while people work in other programs. An app has at most one; it shows only while it has something to show (a running task or event, a notification's source, or a background feature people turned on) and goes when that ends, unless the reader keeps it. Status shows on the icon by a variation or one standard overlay and in its tooltip, changing calmly, never animated for ongoing activity, never drawn disabled. An app with a window shows its own status in the window.

- **Windows:** On Windows notifications come from the tray icon, so an app with no permanent icon shows one only for the task or event behind its notification.
- **macOS:** A menu bar extra only for a feature people use when the app is not in front; people can remove it; notifications come from the app, not the extra.
- **Sources** (16 records):
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 14; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 36; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 1, 4, 8, 9, 11, 12, 13, 18, 20, 21, 22, 24; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 1, 2
- **Enforce:**
  - Component: Tray gains SetOverlay (warning, error, disabled, offline) and its docs state when to Show and Hide.
  - Review: A tray icon kept around in case a notification comes; a second icon; status only on the icon of an app whose window is open.
- **kvit-ui now:** partly. Tray has Show and Hide (`platform/tray.go:176-192`), one per app by its doc (`platform/tray.go:70-71`), a tooltip (`platform/tray.go:208-216`) and SetIcon (`platform/tray.go:200`); no overlay; Notify posts nothing unless the icon shows (`platform/tray.go:276`).

## Controls

### CTL-1

One primary button per view, for the action the view is for; alternatives offered together are the same size and differ only in style.

- **Sources** (3 records):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 4, 5, 6
- **Enforce:**
  - Test, one screen: Fail a region with two or more ButtonPrimary buttons.
- **kvit-ui now:** enforces it. `uitest.Labels` fails a region, dialog, popover, pane or card with two enabled primary buttons (`uitest/labels.go`).

### CTL-2

Choose the control by the choice: a checkbox for one on/off option whose off state is obvious, or for picking several; radio buttons (stacked) for one of two to five options worth the space, a select beyond that or to play a choice down; segmented for two to five short options. A single on/off choice is one checkbox or switch, never two radio buttons and never a slider; a three-valued choice is a radio group, never a checkbox's mixed state. Evenly spaced steps of a range, four or more, are a slider; fewer are radio buttons. Values from data, numbers and choices changed often while working go in a list, select or number field, not radio buttons or checkboxes. A group holds at most ten checkboxes; when at least one must be ticked, an error says so.

- **Windows:** An on/off setting that applies at once is a switch, never a checkbox; a checkbox is for a choice committed later with its form (OK, Next, Submit), or for picking several. A toggle button only where a button fits, such as a toolbar's Bold.
- **macOS:** A switch is for a setting with large scope that acts at once; otherwise an on/off option is a checkbox.
- **Sources** (84 records):
  - Apple: [toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) 1, 10, 14, 18, 19, 20, 12; [pop-up-buttons](https://developer.apple.com/design/human-interface-guidelines/pop-up-buttons) 1, 5, 2; [pickers](https://developer.apple.com/design/human-interface-guidelines/pickers) 1, 2; [pull-down-buttons](https://developer.apple.com/design/human-interface-guidelines/pull-down-buttons) 2
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 4, 5, 7, 10, 1, 2, 3, 6, 8, 26; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 4, 10, 7, 3, 5, 12, 2, 1; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 1, 7, 5, 10, 13; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 3, 12, 4, 5, 1, 6; [ctrl-spin-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-spin-controls) 1; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 4; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 13; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 27; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 25; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 2, 5
  - Fluent: [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 1, 2, 3, 5, 12, 16; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 6, 7, 34, 35, 36; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 8; [ui-controls-toggles](https://learn.microsoft.com/en-us/windows/apps/design/controls/toggles) 1, 3, 4, 5; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 1, 3, 4, 8, 9, 10, 14; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 2, 3; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 2, 3, 4, 5, 6, 7, 8; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 19, 20; [ui-controls-text-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-controls) 6
- **Enforce:**
  - Component: RadioGroup states the limit, two to five, in its doc and in the gallery (`radiogroup.go:22-28`); Switch's doc draws the line between a switch and a checkbox (`switch.go:16-20`).
  - Test, one screen: Fail a RadioGroup with more than five options or fewer than two.
  - Review: Program option or data; a checkbox whose off state is unclear; on Windows and Linux, a checkbox that applies at once.
- **kvit-ui now:** partly. `uitest.Labels` fails a radio group of fewer than two or more than five options (`uitest/labels.go`); the docs of Segmented, RadioGroup and Select draw the line between them. A Switch with `SmallScope` is a switch on the Windows form and a checkbox on the macOS form, acting at once on both (`switch.go`).

### CTL-3

A segmented control either selects or acts, never both; its segments are all text or all symbols, equal in width, five to seven at most, each with a tooltip when symbols.

- **Sources** (10 records):
  - Apple: [segmented-controls](https://developer.apple.com/design/human-interface-guidelines/segmented-controls) 1, 2, 3, 4, 5, 6, 9, 11, 12, 13
- **Enforce:**
  - Component: Segmented.
  - Test, one screen: Fail a Segmented mixing symbol-only and text segments or with more than seven.
- **kvit-ui now:** partly. `uitest.Labels` fails a segmented control of more than seven segments (`uitest/labels.go`); mixing symbol-only and worded segments is not flagged.

### CTL-4

A choice starts on its safest value, or where safety is not at stake on the one most people want, and shows the value in effect, never a recommended one that is not. Options go in a logical order (related together, the default and the common first), alphabetical only from about twelve; All and None come first, and an explicit (None) replaces a blank item. A select holds values, never commands, shows the chosen value, lists the likely values and puts rare extra choices behind Custom…; it takes typing only when people may need a value it does not list, takes a typed value on Enter or when focus leaves, and puts back the last valid value, saying why, when the typed one is not valid. A radio group starts with one option chosen, unless the choice must be explicit (safety, law, no value set yet), and then committing without one is an error; once chosen it cannot be cleared again, so a none answer is its own option.

- **Sources** (40 records, 6 on Microsoft's checklist):
  - Apple: [pop-up-buttons](https://developer.apple.com/design/human-interface-guidelines/pop-up-buttons) 3, 6; [combo-boxes](https://developer.apple.com/design/human-interface-guidelines/combo-boxes) 1, 4, 5; [pickers](https://developer.apple.com/design/human-interface-guidelines/pickers) 3; [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 4, 7; [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 19
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 13; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 13, 23\*, 14, 24, 25; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 71, 72\*; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 38\*, 25, 26, 28, 6; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 31\*; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 26\*, 16, 17, 21, 2, 8, 9, 25; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 25\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 78
  - Fluent: [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 5, 6, 7, 20, 22, 23; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 1
- **Enforce:**
  - Component: Select: an explicit None option and a blank shown value when several selected objects differ; RadioGroup: no option chosen only when the caller says the choice must be explicit.
  - Test, one screen: Fail a Select or RadioGroup option whose words are empty.
  - Review: Which value is safest and which is likeliest; option order.
- **kvit-ui now:** partly. Select shows its value, starts on the first option (`select.go:36-38`), shows a Placeholder when nothing is chosen (`select.go:21-28`) and opens with the current option ticked (`select.go:93-99`); there is no typed (combo) form and no mixed value.

### CTL-5

A checkbox that governs others shows a mixed state when they differ; a click on a mixed checkbox turns all on and the next turns all off, never back to the mix, and people never set the mix directly. In a tree with checkboxes, a parent's box governs its children: ticking it ticks all below, including children loaded later, and ticking every child ticks the parent. A control that belongs to an option sits to its right or indented under it and shares its label, with no colon after the label; typing or pasting into it chooses the option. Keep choices at one level where possible; nested ones are disabled until their parent is chosen. A group's title holds no control except a checkbox that turns the whole group on or off.

- **Windows:** The dependent control shares its option's access key (Windows 7 guide; Fluent's keyboard pages were not read).
- **macOS:** As stated.
- **Sources** (28 records):
  - Apple: [toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) 16, 17, 13, 15
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 22, 16, 21, 15, 23, 25; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 20, 17, 19, 21, 22; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 6; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 24; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 10; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 36, 39; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 34
  - Fluent: [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 4, 14, 15; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 21; [ui-controls-tree-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tree-view) 11, 13, 14
- **Enforce:**
  - Component: Check: a press on a mixed box checks all and the next clears all (exists); a dependent-control helper that ticks the box on input; Tree: a parent's box governs its children.
  - Review: Nesting and group titles.
- **kvit-ui now:** partly. A press on a mixed Check checks the whole and the next press clears it, as the rule says (`check.go:49-51`); nothing links a dependent field to its box, and Tree has no checkboxes.

### CTL-6

Make on and off differ by more than colour, and say what a toggle turns on; a switch's label is one or two words, preferably a noun naming what it controls (Wi-Fi).

- **Windows:** A switch shows On or Off beside it (decided, Q-44).
- **macOS:** As stated.
- **Sources** (7 records):
  - Apple: [toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) 2, 3, 4
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 32
  - Fluent: [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 6; [ui-controls-toggles](https://learn.microsoft.com/en-us/windows/apps/design/controls/toggles) 2, 6
- **Enforce:**
  - Component: Switch, Check, checkable buttons (exist); Switch shows On or Off beside it on the Windows form.
- **kvit-ui now:** enforces it. A switch says On or Off beside it on the Windows form (`switch.go`); its knob moves as its track fills.

### CTL-7

Prefer a control that accepts only valid values (select, list, slider, stepper) to free text; a text field is for values limited only by length or kind of character. Use a single-line field, with no scroll bar, for short entries and a wrapping text area of set height, which scrolls rather than grows, for long ones; plain text unless formatting is needed, and a rich text editor comes with working formatting commands. A field accepts the common ways a value is written, limits its length to the longest valid form, states the format when it cannot accept them all, and uses the unit that spares people decimals, labelled; long text with a limit shows a live count of what remains. Pasted text is fitted to the field: a single-line field drops line breaks, and a paste past the limit is reported, not silently cut. Number fields take only numbers and show them formatted, and an emptied one holds no value, never zero; a special case such as Unlimited is its own option, never a magic number. On focus a field selects its whole value when people usually replace it. Spell checking is on where people type words and sentences, even in a language the app does not support, and off for names, numbers and codes; it corrects only when sure and otherwise underlines the word and offers alternatives in its context menu. Secrets are masked and a password field opens empty; a reveal button, clear of other controls, shows the password only while held, or a Show password checkbox switches it. A new password is typed twice; signing in asks once.

- **Windows:** A focused single-line field holding text shows a clear button (decided, Q-44).
- **macOS:** As stated.
- **Sources** (66 records, 5 on Microsoft's checklist):
  - Apple: [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 1, 4, 11, 13; [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 5, 6, 10; [text-views](https://developer.apple.com/design/human-interface-guidelines/text-views) 1; [labels](https://developer.apple.com/design/human-interface-guidelines/labels) 1; [image-wells](https://developer.apple.com/design/human-interface-guidelines/image-wells) 1
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 159\*; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 1\*, 10, 11, 14, 17, 15, 16, 18, 19, 22, 26, 5, 20, 27; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 32\*; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 19\*; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 25\*; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 24
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 37; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 1, 2, 4, 8, 11, 15, 17, 20, 21, 22, 23, 24, 26; [ui-controls-number-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/number-box) 1, 2, 5, 10, 11, 13; [ui-controls-password-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/password-box) 1, 5, 6, 7, 8, 9, 10; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 3, 18; [ui-controls-text-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-controls) 5, 7, 8, 15, 16, 17, 18, 19
- **Enforce:**
  - Component: Field: a length limit with a remaining count, paste fitting, a masked mode with reveal, select-all on focus and spell checking set per field; NumberField and MoneyField (exist). A clear button on a focused single-line Field holding text, on the Windows form.
  - Review: A constrained control wherever the valid values can be listed; spell checking off for names, numbers and codes.
- **kvit-ui now:** partly. Field has `Masked` for a secret, `MaxLength` with text past it reported as the field's error rather than cut, and a clear button while it holds the keyboard and text on the Windows form (`field.go`); NumberField checks as the reader types. There is no remaining count, paste fitting, select-all on focus or spell checking.

### CTL-8

A control that cannot act now is removed when people cannot make it work or would not expect it; disabled, with its label and the buttons that serve it, when they expect it and can see why; and otherwise left enabled, saying what is wrong when it is used. A disabled control shows the state that actually applies, and a button that acts on a selection waits for one. A submit button (Save, Next, Continue) may stay disabled until every required field is filled, with the required fields marked and the reason shown; otherwise it stays enabled and pressing it reports what is missing beside the field.

- **Sources** (25 records, 7 on Microsoft's checklist):
  - Apple: [entering-data](https://developer.apple.com/design/human-interface-guidelines/entering-data) 11
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 152\*, 153\*, 154\*, 155\*, 156, 97; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 33\*; [ctrl-spin-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-spin-controls) 12; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 18; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 20; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 7; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 32; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 24; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 18; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 19; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 7; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 22; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 16\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 33\*
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 28, 40; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 11; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 7; [ui-controls-selection-modes](https://learn.microsoft.com/en-us/windows/apps/design/controls/selection-modes) 3
- **Enforce:**
  - Component: Button.Explanation on a disabled button (exists); a field, select or list that is disabled takes its visible label with it, which needs a label linked to the control.
  - Test, one screen: Fail a disabled Button with no Explanation.
  - Review: Remove, disable or explain, by the three tests.
- **kvit-ui now:** enforces it. A disabled button's Explanation, or its command's, is shown and read (`button.go`); `uitest.Labels` fails a disabled button without one (`uitest/labels.go`).

### CTL-9

A stepper is for small numbers changed by small amounts: each press moves by the smallest common step, holding a button repeats it, it never yields an invalid value, special values are words people can step to (Never), and the value can also be typed. Its buttons sit beside the value; stepping stops at the limits, and the button that would pass one is disabled. A slider is for a quantity judged by feel or effect (volume, brightness), oriented as the quantity usually is (a vertical one has its largest value at the top in every language), long enough to hit every value, non-linear over a very wide range or snapping only to a few representative values, with tick marks where it snaps and a value label with its unit, or a sample of the effect, when the value matters; whole units snap to whole steps and the last step lands on the maximum. It applies while dragged and labels its ends, one word each (decided, Q-51); where the slider's label sits is appearance, not a rule (decided, Q-44). An exact, known number is a number field, a value that cannot change is text, and a slider never shows progress.

- **Windows:** The Up and Down keys step a focused number field or stepper, Page Up and Page Down by a large step. A slider gets a number field beside it only when both instant feedback and exact entry are needed; a number field alone where space is tight or the keyboard is likely.
- **macOS:** A slider over a wide range comes with a number field.
- **Sources** (72 records):
  - Apple: [steppers](https://developer.apple.com/design/human-interface-guidelines/steppers) 1, 2; [sliders](https://developer.apple.com/design/human-interface-guidelines/sliders) 2, 3, 5, 6, 8, 9, 10
  - Windows 7: [ctrl-spin-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-spin-controls) 4, 5, 6, 9, 13, 14, 15, 16, 17, 18, 11, 2, 3; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 1, 3, 4, 6, 7, 8, 10, 11, 14, 15, 20, 19, 12, 13; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 3, 12, 23; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 8
  - Fluent: [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 5; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 12; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 1, 4, 5, 6, 7, 8, 9, 10, 13, 15, 16, 19, 22, 23, 24, 25, 26, 27, 31, 35, 36, 37; [ui-controls-number-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/number-box) 6, 7, 8, 9; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 22; [ui-controls-text-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-controls) 9, 10, 11
- **Enforce:**
  - Component: Stepper: a typed value, Up and Down keys, words for special values, repeat while held; Slider (exists).
  - Review: Slider, stepper or number field.
- **kvit-ui now:** partly. Stepper is a figure with minus and plus buttons beside it, clamped at the ends with the button at an end disabled, as the rule says (`stepper.go:12-15`, `stepper.go:61-69`, `stepper.go:72-75`); its value cannot be typed, it takes no arrow keys and a held button does not repeat (`control.go:153-154`); Slider draws its value beside it (`slider.go:17-21`) and moves with the arrow keys (`slider.go:125-144`).

### CTL-10

Static tabs hold peer pages of one thing; a control on a page changes only that page, and switching tabs only moves: it applies nothing, shows no error, and the last tab open is not a setting. They carry text, never nest, never scroll sideways and are not the steps of a task; beyond a few static pages, use a navigation sidebar rather than a tab row. A window with one page has no tab; a tab that does not apply and that nobody would look for is removed, and one people might look for stays, its controls disabled and the reason said on the page. Pages run general first, advanced last. A short row of view selectors may switch between a few views or filters of the same data (All, Unread, Flagged), always with one chosen; many views switch with a select or split button, and a choice that is not navigation is radio buttons. The tab row is visible unless the app alone switches pages. Document tabs, which people open, close, reorder and drag between windows, are a separate case: one is always active, and closing it activates a neighbour; each closable tab has its own close button; a new-tab button and Ctrl+T open one, Ctrl+W closes one, Ctrl+Shift+T reopens the last closed and Ctrl+1 to Ctrl+9 go to a tab (Command on macOS); tabs reorder by dragging and tear off into a new window; a strip that overflows scrolls with buttons at its ends; closing the last tab may close the window.

- **Windows:** Up to seven static tabs in one row that fits even when translated. A narrow inspector chooses its page with a select.
- **macOS:** At most six tabs.
- **Sources** (60 records, 5 on Microsoft's checklist):
  - Apple: [tab-views](https://developer.apple.com/design/human-interface-guidelines/tab-views) 1, 2, 5, 6, 7, 8
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 26; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 5, 3, 4, 7, 13, 14, 15, 19, 22\*, 26, 23\*, 6, 1, 2, 8, 9, 11, 12; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 20; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 39, 38, 19, 20, 36\*, 37, 13, 35, 40\*, 41, 43, 63, 33, 34; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 89\*; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 5
  - Fluent: [ui-controls-tab-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tab-view) 1, 2, 3, 4, 9, 11, 12, 13, 14; [ui-controls-pivot](https://learn.microsoft.com/en-us/windows/apps/design/controls/pivot) 3; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 11, 12, 13; [ui-controls-selector-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/selector-bar) 1, 2, 4, 6, 7
- **Enforce:**
  - Component: Tab carries text and a count, no symbol (`tab.go:16-29`); a document tab strip (close buttons, new tab, reorder, tear-off, the tab shortcuts).
  - Test, one screen: Fail a tab row longer than the platform allows (seven on Windows and Linux, six on macOS), a row of one tab, or a tab row inside a tab page.
- **kvit-ui now:** partly. `uitest.Labels` fails a row of one tab and a row longer than seven on the Windows form or six on the macOS form (`uitest/labels.go`); a tab row inside a tab page is not flagged.

### CTL-11

Use a list or table for text items, a preview list where pictures choose better than words, and a tree only for real hierarchy: broad rather than deep (at most four levels, the most used in the first two), ordered so people can predict where things are, with search or a filter when it is large, a root node only when the whole tree needs commands, headings when it has columns, and detail for a separate task shown beside it. The items of one list share one structure and one click behaviour. Columns start with the name, then what tells items apart, long text last; the most useful show by default and people can add, remove and reorder them. Sort by default as the data suggests: names alphabetically, numbers by value, dates by time. A long grouped list keeps the current group's heading at the top while it scrolls, and a collection switched between list and grid keeps its selection.

- **Windows:** A tree only when seeing the hierarchy matters; otherwise people drill down level by level through lists. No column browser.
- **macOS:** A column browser for deep hierarchies that need no sorting.
- **Sources** (74 records):
  - Apple: [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 1, 15, 4, 9; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 1, 2, 4; [column-views](https://developer.apple.com/design/human-interface-guidelines/column-views) 1, 2, 3; [collections](https://developer.apple.com/design/human-interface-guidelines/collections) 1, 2
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 89; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 8, 30, 31, 36, 1, 2, 3; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 11, 5, 6, 9, 31, 13, 15, 18, 25, 26, 27, 28, 29, 30, 3, 1, 2, 4, 7, 8; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 14, 2, 3, 9, 20; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 22, 23
  - Fluent: [ui-controls-listview-and-gridview](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-and-gridview) 1, 2, 3, 4; [ui-controls-item-containers-templates](https://learn.microsoft.com/en-us/windows/apps/design/controls/item-containers-templates) 1, 2; [ui-controls-tree-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tree-view) 2, 3, 4, 8; [ui-controls-annotated-scrollbar](https://learn.microsoft.com/en-us/windows/apps/design/controls/annotated-scrollbar) 12; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 2, 11; [ui-controls-listview-filtering](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-filtering) 1; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 25; [ui-controls-itemsview](https://learn.microsoft.com/en-us/windows/apps/design/controls/itemsview) 1, 2; [ui-controls-lists](https://learn.microsoft.com/en-us/windows/apps/design/controls/lists) 1, 3, 4, 6, 7, 10, 11, 12; [ui-controls-selector-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/selector-bar) 5; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 21
- **Enforce:**
  - Component: Table: reader-owned columns (exist); Tree.
  - Review: List, tree or browser; column order and default sort.
- **kvit-ui now:** partly. Table, Tree, Row and Card exist; Table columns are the reader's to sort, resize, move and hide (`table.go:62`, `table.go:595-605`); there is no column browser, preview list or tree search.

### CTL-12

Text cut short ends in an ellipsis, never a bare clip, and where there is room for more lines it wraps instead. In a narrow column, shorten text in the middle when its start and end both identify it, such as file names.

- **Sources** (8 records, 1 on Microsoft's checklist):
  - Apple: [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 5; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 12
  - Windows 7: [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 25; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 22\*; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 47
  - Fluent: [signature-experiences-typography](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/typography) 5, 12, 13
- **Enforce:**
  - Component: Cell for path and name kinds; Slug elides in the middle (`marks.go:279-303`).
- **kvit-ui now:** partly. Slug elides in the middle; table cells elide at the end.

### CTL-13

Let people reorder, add and remove items where order matters, by dragging with commands kept for the keyboard, and change values such as ratings in place. To build a short list people must review, use a list builder (available and chosen lists with Add, Remove, Move Up and Move Down), or, where space is short or values come from outside the list, an add and remove list whose Add opens a chooser. A row has one primary action, what pressing it does, at its leading side; buttons for other actions sit at its trailing end only when every input must reach them at all times, and otherwise on the page the row opens; either way they are also in the row's context menu, reachable by keyboard.

- **Sources** (18 records):
  - Apple: [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 2; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 11; [rating-indicators](https://developer.apple.com/design/human-interface-guidelines/rating-indicators) 1; [collections](https://developer.apple.com/design/human-interface-guidelines/collections) 3, 5
  - Windows 7: [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 16, 17, 18, 19, 23, 24, 15
  - Fluent: [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 2; [ui-controls-nested-ui](https://learn.microsoft.com/en-us/windows/apps/design/controls/nested-ui) 1, 2, 3, 4, 8
- **Enforce:**
  - Component: DualList (exists); an add and remove list; UI.SetContextMenu for row actions (exists).
  - Review: Row actions: shown in the row or on its page, and in its context menu.
- **kvit-ui now:** partly. DualList moves items between its lists and reorders them by keyboard and buttons (`duallist.go:11-17`), not by dragging; there is no add and remove list; ListRow marks a pressable row with a chevron at its trailing edge (`listrow.go:34-39`), and UI.SetContextMenu gives a row a context menu that also opens from the keyboard (`menu.go:691-696`, `menu.go:714`).

### CTL-14

Choosing a value in a select, radio group, checkbox or list only sets it: it never runs a command, opens a window or shows or hides other controls. A row that navigates (a sidebar, a master list) stays highlighted while its content shows beside it, and the path to the detail stays highlighted in every pane; a selected item inside a collapsed group marks the group until it opens; a row that picks an option shows a tick instead. Selection covers the whole row and stays visible when the list loses focus; a list nothing acts on selects nothing, and a collection of destinations may open an item on one click instead. A click outside a list never clears its selection; a selection belongs to its own view, and its commands touch nothing elsewhere. Where choosing several is common or losing the set would hurt, rows have checkboxes, so one click or Space changes one item and the arrow keys move without changing the set, with Select all and Clear all on long lists, and Select all only where several may be chosen.

- **Sources** (44 records):
  - Apple: [lists-and-tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) 3; [split-views](https://developer.apple.com/design/human-interface-guidelines/split-views) 1
  - Windows 7: [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 19; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 5, 6, 10, 23, 34; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 11, 37, 12, 31, 33; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 10, 12, 33; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 15; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 18
  - Fluent: [ui-controls-navigationview](https://learn.microsoft.com/en-us/windows/apps/design/controls/navigationview) 15, 20; [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 13; [ui-controls-listview-and-gridview](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-and-gridview) 7, 8, 10, 12, 13; [ui-controls-item-containers-templates](https://learn.microsoft.com/en-us/windows/apps/design/controls/item-containers-templates) 7; [ui-controls-tree-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/tree-view) 10, 12; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 16, 19; [ui-controls-itemsview](https://learn.microsoft.com/en-us/windows/apps/design/controls/itemsview) 3, 4, 7, 11; [ui-controls-lists](https://learn.microsoft.com/en-us/windows/apps/design/controls/lists) 2, 8; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 16; [ui-controls-selection-modes](https://learn.microsoft.com/en-us/windows/apps/design/controls/selection-modes) 2, 4, 7; [ui-controls-list-details](https://learn.microsoft.com/en-us/windows/apps/design/controls/list-details) 1, 3; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 10
- **Enforce:**
  - Component: Sidebar, Table (exist); the docs of Select, RadioGroup and Check say OnChange only sets.
  - Review: Choosing a value only sets it.
- **kvit-ui now:** partly. SidebarItem.Selected and Table's current row (`sidebar.go`, `table.go`); Table has a check column toggled by Space (`tablebody.go:415-424`) and a header box that selects every shown row (`table.go:843`); nothing stops an OnChange from showing or hiding controls.

### CTL-15

A search field says what it searches, updates results as the reader types, starts with the widest scope, offers recent searches and narrows by scope or token; once text is typed it has a clear button, and clearing it restores the original view; one place searches everything. A filter over a list keeps the items that contain the typed text anywhere; a search across the app ranks matches at the start of a word first (decided, Q-45). A search that finds nothing says No results where the results would be.

- **Windows:** App-wide search is a box centred in the window's own header, below the system's title bar (decided, Q-41); a search of one view sits at that view's top trailing corner. The clear button also stops a slow search; many results are ranked by relevance, and the results are one step in the Back history.
- **macOS:** Window search at the trailing end of the toolbar; sidebar filter at the top of the sidebar.
- **Sources** (38 records):
  - Apple: [search-fields](https://developer.apple.com/design/human-interface-guidelines/search-fields) 1, 2, 3, 4, 5, 6, 7, 8, 9, 17, 18, 19, 21; [searching](https://developer.apple.com/design/human-interface-guidelines/searching) 1, 2, 3, 4, 6; [outline-views](https://developer.apple.com/design/human-interface-guidelines/outline-views) 13; [token-fields](https://developer.apple.com/design/human-interface-guidelines/token-fields) 3
  - Windows 7: [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 4, 13, 14, 15, 17, 16, 6, 2, 5, 11, 12, 18
  - Fluent: [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 18; [ui-controls-auto-suggest-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/auto-suggest-box) 2, 4; [ui-controls-listview-filtering](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-filtering) 3, 4; [ui-controls-text-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-controls) 3
- **Enforce:**
  - Component: SearchField: the placeholder is set per use; clear button and Escape (exist).
  - Test, one screen: Fail a SearchField left at the bare default placeholder on an app screen.
- **kvit-ui now:** enforces it. `uitest.Labels` fails a search field left at its "Filter…" placeholder (`uitest/labels.go`).

### CTL-16

Show the common options and hide advanced ones behind a labelled disclosure beside what it expands, at most one advanced-options disclosure per view; several independent expandable sections may stack in one view, and opening one closes no other unless the reader chose that. The same control closes it again, it goes away when there is nothing more to show, it pushes what follows down rather than covering it, and a closed section's header sums up the values chosen inside it. Never hide what people need: an option everyone understands is not advanced because its technology is. Bulky detail or a separate task opens in its own window, and hiding items to narrow a list is a filter, not a disclosure.

- **Windows:** Label it More options and Fewer options (or commands, details); in a message, Show details and Hide details.
- **macOS:** A disclosure triangle or button beside a label naming what it reveals.
- **Sources** (46 records, 8 on Microsoft's checklist):
  - Apple: [disclosure-controls](https://developer.apple.com/design/human-interface-guidelines/disclosure-controls) 1, 2, 3, 4; [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 5
  - Windows 7: [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 23\*, 3, 6, 9, 17, 20, 1, 2, 10, 16, 31\*, 32\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 25; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 147\*, 148, 149; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 24, 25; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 39\*; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 36; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 7; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 11; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 45; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 53\*, 54\*, 55\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 33; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 15
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 20, 31; [ui-controls-listview-filtering](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-filtering) 2; [ui-controls-expander](https://learn.microsoft.com/en-us/windows/apps/design/controls/expander) 1, 2, 3, 4, 5, 7, 8, 9, 10, 11
- **Enforce:**
  - Component: Disclosure (exists): its Title is the caller's, so the Windows words are written per platform by the app.
  - Review.
- **kvit-ui now:** partly. Disclosure has a titled trigger that opens and closes the same section and grows open, pushing what follows down (`disclosure.go:16-20`, `disclosure.go:23-47`); a shared Group makes an accordion that closes the others, which is the caller's choice, not the reader's (`disclosure.go:21-22`, `disclosure.go:28-29`); its header shows a count, not a summary of the values; it supplies no More/Fewer wording per platform.

### CTL-17

Show the sidebar by default, let it be hidden and shown by a command and shortcut, keep it to two levels, and group long lists under headings. On a narrow window it shrinks to a rail of symbols, each with a tooltip, when every section has a clear symbol, and otherwise to one button that opens it over the content; the content stays visible. A group with no page of its own only opens and closes and is never shown selected; a group with a page opens it on a click, and its arrow only opens and closes.

- **Sources** (29 records):
  - Apple: [sidebars](https://developer.apple.com/design/human-interface-guidelines/sidebars) 2, 3, 5, 6, 7, 12; [split-views](https://developer.apple.com/design/human-interface-guidelines/split-views) 5, 6, 7
  - Fluent: [ui-controls-navigationview](https://learn.microsoft.com/en-us/windows/apps/design/controls/navigationview) 1, 2, 3, 5, 8, 10, 11, 12, 13, 18, 19, 21; [ui-controls-split-view](https://learn.microsoft.com/en-us/windows/apps/design/controls/split-view) 1, 2, 3; [layout-responsive-design](https://learn.microsoft.com/en-us/windows/apps/design/layout/responsive-design) 10; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 15; [ui-controls-selector-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/selector-bar) 3; [basics-app-silhouette](https://learn.microsoft.com/en-us/windows/apps/design/basics/app-silhouette) 1, 4
- **Enforce:**
  - Component: Window: a Show/Hide Sidebar command and shortcut over SidebarVisible; sidebar groups.
  - Review.
- **kvit-ui now:** partly. The sidebar collapses to a rail below WidthLaptop (`window.go:163-167`, `tokens/interface.go:196`), every item has a symbol for it (`sidebar.go:77-82`), the rail names an item in its tooltip (`sidebar.go:131-138`) and opens over the body on hover or keyboard focus (`window.go:34-36`, `window.go:170-173`); Window.SidebarVisible hides it as the reader's choice (`window.go:31-33`), but no built-in command or shortcut toggles it, and SidebarItem has no groups (`sidebar.go:83-114`).

### CTL-18

Content text people may want to reuse (messages, addresses, identifiers) can be selected and copied; control labels cannot; read-only text never looks disabled. A field is read-only only while it is for a time not editable; text that is never editable is plain, selectable text, not a field.

- **Sources** (10 records):
  - Apple: [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 6, 7; [labels](https://developer.apple.com/design/human-interface-guidelines/labels) 4; [text-views](https://developer.apple.com/design/human-interface-guidelines/text-views) 4
  - Windows 7: [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 11; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 6, 30; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 4
  - Fluent: [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 9; [ui-controls-text-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-controls) 2
- **Enforce:**
  - Component: Label gains `Selectable` for content text.
- **kvit-ui now:** partly. A Label with `Selectable` is selected whole by a press and copied by Copy, Ctrl+C or its context menu (`label.go`); part of its text cannot be selected.

### CTL-19

A field for values typed again and again (names, addresses, file names) completes them from a known set as people type, never secrets such as passwords, PINs or card numbers; moving through the suggestions puts each one's text, worded as in the list, into the field. A token field turns typed text into a token on comma or Return, each token has its own menu, and suggestions wait for a pause.

- **Sources** (9 records):
  - Apple: [token-fields](https://developer.apple.com/design/human-interface-guidelines/token-fields) 1, 2, 3
  - Windows 7: [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 8; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 35; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 10
  - Fluent: [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 3; [ui-controls-auto-suggest-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/auto-suggest-box) 1, 3
- **Enforce:**
  - Component: TypeAhead (exists); a token field.
  - Review: No completion on secret fields.
- **kvit-ui now:** partly. TypeAhead suggests matches as the reader types (`typeahead.go:47-56`, `typeahead.go:133-159`); there is no token field.

### CTL-20

Show a picker in place, beside the field it edits or in a popover, with predictably ordered values. Readers pick colours only for their own categories (labels, tags), and only from the chart ramp's named steps; nowhere else does the app offer a colour choice (decided, Q-35). A file or folder field has a Browse… button beside it that opens the system's file dialog, and any field whose valid values can be looked up may have a button beside it that opens a chooser.

- **Sources** (18 records):
  - Apple: [pickers](https://developer.apple.com/design/human-interface-guidelines/pickers) 3, 4, 5, 7; [color](https://developer.apple.com/design/human-interface-guidelines/color) 7; [color-wells](https://developer.apple.com/design/human-interface-guidelines/color-wells) 1, 2
  - Windows 7: [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 21, 2; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 14; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 11
  - Fluent: [ui-controls-color-picker](https://learn.microsoft.com/en-us/windows/apps/design/controls/color-picker) 1, 4, 5, 8, 3; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 16; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 23
- **Enforce:**
  - Component: A path field with Browse… over platform.FileDialog; a colour choice offering the categorical ramp's steps (`tokens/theme.go:183-189`).
  - Review: A colour choice anywhere but the reader's own categories.
- **kvit-ui now:** partly. No date or colour picker; the categorical ramp's steps exist as theme values (`tokens/theme.go:183-189`, `data.go:32-33`) but nothing offers them to the reader; platform.FileDialog opens the system's file dialog with a title and filters (`platform/filedialog.go:12-56`), but no field pairs it with a Browse… button.

### CTL-21

Never nest scrolling in the same direction; prefer vertical scrolling to horizontal, and scroll only along the axes the content overflows (long text wraps and scrolls vertically); show that content continues; a scroll bar appears only when the content does not fit and never covers a control; bring into view whatever a command acts on or changes. A very long collection may have an annotated scroll bar: short category labels where they help (the first always kept when labels collide), a hover that previews the place, a click that jumps there and a drag that scrolls; short collections and short views keep the plain bar.

- **Sources** (23 records):
  - Apple: [scroll-views](https://developer.apple.com/design/human-interface-guidelines/scroll-views) 3, 4, 6, 7, 9, 10, 11
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 21; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 34
  - Fluent: [ui-controls-annotated-scrollbar](https://learn.microsoft.com/en-us/windows/apps/design/controls/annotated-scrollbar) 1, 2, 6, 7, 8, 9, 10, 11; [ui-controls-scroll-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/scroll-controls) 1, 3, 4, 5, 8, 11
- **Enforce:**
  - Test, one screen: Fail a vertical Region inside a vertical Region.
  - Component: ScrollBar in its own strip, hidden when everything fits (exists).
- **kvit-ui now:** partly. Region scrolls; ScrollBar keeps a strip of its own rather than floating over the content, and is hidden when everything fits (`scrollbar.go:9-15`); nothing checks nesting; there is no annotated scroll bar.

### CTL-22

A control looks like a control and is never mistaken for the content it acts on; a text link is underlined at rest, so it is recognisable without hovering and without colour.

- **Sources** (15 records, 5 on Microsoft's checklist):
  - Apple: [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 6
  - Windows 7: [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 9\*, 7\*, 10; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 5\*; [ctrl-tabs](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tabs) 10; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 29; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 8; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 80\*; [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 16; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 19; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 21; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 6\*
  - Fluent: [ui-controls-hyperlinks](https://learn.microsoft.com/en-us/windows/apps/design/controls/hyperlinks) 11; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 9
- **Enforce:**
  - Component: Link: an underline at rest.
  - Review.
- **kvit-ui now:** partly. BorderStrong outlines controls at 3:1 (SKILL.md, Lines); a Link turns the accent colour and underlines only on hover or focus (`link.go:12-15`).

### CTL-23

A gauge or level gives its value and its limits in words; a level that is not progress, such as disk space used, is a gauge, never a progress bar; ratings use stars. A rating the reader cannot change is shown read-only; one may show the average until the reader sets their own, told apart by a caption.

- **Sources** (7 records):
  - Apple: [gauges](https://developer.apple.com/design/human-interface-guidelines/gauges) 1, 3; [rating-indicators](https://developer.apple.com/design/human-interface-guidelines/rating-indicators) 2
  - Windows 7: [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 42
  - Fluent: [ui-controls-rating](https://learn.microsoft.com/en-us/windows/apps/design/controls/rating) 1, 2, 3
- **Enforce:**
  - Component: Gauge (exists).
- **kvit-ui now:** partly. Gauge exists (catalog KvitGauge).

### CTL-24

Never play audio or video without visible controls to stop it; the playback position is a horizontal slider, not a progress bar; media with speech offers captions or an alternative audio track.

- **Sources** (4 records):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 29
  - Windows 7: [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 4
  - Fluent: [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 17; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 7
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. No media.

### CTL-25

When a question's answers need more words than a button holds, each answer is a radio button with a short main line and, where needed, a line saying why one would pick it, and a specific commit button acts on the choice; there are no command links on any system (decided, Q-10). Offer two to five, the likeliest or the do-everything answer first; Cancel stays a separate button, never one of the choices. If only one answer would remain, add one worded by how it differs, or make the window a plain message. Up to three short answers are instead buttons with specific verbs.

- **Sources** (30 records, 12 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 131, 132; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 28, 30, 29; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 40\*, 41\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 39\*, 40\*, 41\*, 94\*, 95\*; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 11; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 1, 2, 3, 4\*, 6, 7, 8, 9\*, 10, 11, 13\*, 15\*, 16\*, 18, 14, 26; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 7
- **Enforce:**
  - Component: RadioGroup with a Detail line per option (`radiogroup.go:14-20`) and a Dialog form whose commit button is worded by the action, with Cancel kept.
  - Test, one screen: Fail a dialog's radio choice of fewer than two or more than five answers, or one whose commit button is OK, Yes or Next.
  - Review: Answers too long for buttons are radio buttons with a specific commit button; no command links.
- **kvit-ui now:** partly. RadioGroup gives each option a line under it, read with its name (`radiogroup.go:14-20`, `radiogroup.go:150-153`, `radiogroup.go:242-243`); Dialog has one confirm and one cancel button (`feedback.go:387-389`, `feedback.go:424`) and no form that holds a radio group.

### CTL-26

A button starts an action at once, in one to three words; a link only goes somewhere (another page, window, help topic or website) and commits nothing. Every action, however light or wherever it stands, is a button; a destructive, irreversible or committing action is always a button. A secondary command unrelated to the window's purpose is a quiet button. The words people must not miss are on the control, not in text around it. Every link has a tooltip saying where it leads (for a website, its domain), and a web link opens in the default browser. A link is words, or a picture whose name says where it goes; separate links never touch; related ones are gathered in one place; links show that they were visited.

- **Sources** (31 records):
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 1, 16; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 4, 1, 2, 3, 5, 6; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 1, 2, 21, 5, 6, 15, 16, 19, 22, 32
  - Fluent: [ui-controls-hyperlinks](https://learn.microsoft.com/en-us/windows/apps/design/controls/hyperlinks) 1, 2, 3, 4, 6, 7, 9; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 38, 39; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 1, 2; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 1; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 17
- **Enforce:**
  - Component: Button and Link (exist); Link.Explanation is the tooltip saying where the link leads, for a website its domain.
  - Test, one screen: Fail a Button whose words are five or more words or end with a full stop; fail a Link with no name or no Explanation.
  - Review: Button or link: a link that runs an action.
- **kvit-ui now:** partly. Button and Link are separate controls with roles Button and Link (`button.go:202-219`, `link.go:12-15`); Link's doc still calls it a destination or action (`link.go:12`); Link.Explanation is its tooltip and description (`link.go:24-27`, `link.go:49`), but nothing requires one; nothing limits a Button's words, which never elide (`button.go:119-123`); Link has no visited state.

## Feedback

### FDB-1

Every control shows that it is pressed while the button or key is held.

- **Sources** (2 records):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 3
  - Fluent: [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 14
- **Enforce:**
  - Component: control (exists).
- **kvit-ui now:** enforces it. `control.go` tracks Pressed for every control.

### FDB-2

Show progress for the computer's work, never for steps the reader takes in the window, for anything longer than a moment, in the same place for the same kind of work (modeless progress in the window's status bar); a task whose progress is obvious needs none. While the app waits on a step the reader completes elsewhere (signing in through a browser), an indeterminate indicator may show. Work of unknown length that blocks the window, or a large part of it, shows an indeterminate ring where its result will appear; other work shows a bar. A list being filled shows one bar at its top, never one per item. Only one indicator at a time, never a busy pointer and a bar together. Start indeterminate and switch to determinate as soon as the total is known, never back. The bar is honest: it moves while work goes on and stops when none is done, only ever goes forward, does not race early or sit at 99 percent, and is full only when the work is done. A problem or a pause shows on the bar itself, by more than colour.

- **Windows:** Wait of one to five seconds with nothing else showing it: the busy pointer (or the working-in-background pointer when the app still responds); longer: a progress bar or progress dialog with Cancel or Stop. Windows marks a blocked bar red and a paused one yellow; Kvit adds a word and a symbol (A11Y-2).
- **macOS:** The app sets no wait pointer (the system shows its own only when the app stops responding); progress is always an indicator in the window.
- **Sources** (55 records, 8 on Microsoft's checklist):
  - Apple: [loading](https://developer.apple.com/design/human-interface-guidelines/loading) 5, 6; [progress-indicators](https://developer.apple.com/design/human-interface-guidelines/progress-indicators) 1, 2, 3, 5, 6, 8, 13, 14
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 67, 76; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 1, 27, 2, 8, 18\*, 21, 22, 25, 29, 31, 5, 6, 7, 20\*, 24\*, 26, 28, 32, 48, 50; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 16; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 40, 41, 42\*; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 12; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 52, 53\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 60\*, 56\*, 58\*; [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 3; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 17; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 3
  - Fluent: [ui-controls-progress-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/progress-controls) 1, 2, 3, 5, 8, 9, 14, 15, 16, 19
- **Enforce:**
  - Component: Progress gains Paused and Problem states drawn with a symbol and words, and an indeterminate ring form; a busy-pointer helper sets the platform's wait cursor on Windows and Linux, refusing while a Progress shows.
  - Review: A bar that jumps to the end and waits; a bar for the reader's own steps; an indicator on each list item.
- **kvit-ui now:** partly. Progress has determinate and indeterminate bar forms with a label (`progress.go:14-31`); no ring form; no busy pointer anywhere (no wait cursor is set in kvit-ui); Button.Busy changes its words to "Working…" (`button.go:51-56`); no paused or problem state.

### FDB-3

A long operation can be cancelled when stopping is harmless, and paused as well when cancelling would throw work away; cancelling that loses progress asks first.

- **Windows:** The halt button reads Cancel when stopping leaves no trace and Stop when partial results remain, changing from Cancel to Stop once undoing is no longer possible; Pause becomes Resume while paused.
- **Sources** (7 records, 1 on Microsoft's checklist):
  - Apple: [progress-indicators](https://developer.apple.com/design/human-interface-guidelines/progress-indicators) 9, 10
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 71\*, 72; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 41; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 37; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 52
- **Enforce:**
  - Component: Progress gains an optional halt button whose words are Cancel or Stop as the application says, and Pause/Resume.
- **kvit-ui now:** enforces it. Progress has a halt button whose words the app gives, Cancel or Stop, and Pause, which becomes Resume (`progress.go`).

### FDB-4

Show a usable window within a couple of seconds; load in the background with placeholders, never an empty window. A preview page appears within about a second, quick and adequate before slow and exact.

- **Sources** (9 records):
  - Apple: [loading](https://developer.apple.com/design/human-interface-guidelines/loading) 1, 2, 3; [launching](https://developer.apple.com/design/human-interface-guidelines/launching) 1, 3; [onboarding](https://developer.apple.com/design/human-interface-guidelines/onboarding) 8, 9
  - Windows 7: [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 40; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 16
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered; `--win-smoke` measures first-frame time (CLAUDE.md).

### FDB-5

Put status next to what it describes, and only status people may act on. Background work of minor interest shows its status as text, not an indicator; brief status may also sit in the toolbar. A list with multiple selection says beneath it how many are selected (or what they add up to). The status bar holds only brief facts: mode indicators shown disabled when they do not apply, other items removed; important ones in words, not a symbol alone; changes infrequent and calm. Nothing essential lives only at the window's bottom edge.

- **Sources** (27 records):
  - Apple: [feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) 3; [windows](https://developer.apple.com/design/human-interface-guidelines/windows) 10, 11; [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 19; [sidebars](https://developer.apple.com/design/human-interface-guidelines/sidebars) 13
  - Windows 7: [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 22; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 35; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 32; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 2; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 6, 8, 12, 13, 17, 18, 19, 4, 5, 7, 9, 11, 14; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 9; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 6; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 2
  - Fluent: [ui-controls-progress-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/progress-controls) 11; [basics-commanding-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/commanding-basics) 11
- **Enforce:**
  - Component: StatusBar (facts only); Table and List gain a selection count line.
  - Review: A status item nobody would act on; a symbol-only fact.
- **kvit-ui now:** partly. StatusBar holds Activity and Facts, overflowing into "N more" (`statusbar.go:33-78`, `statusbar.go:288`); no disabled mode indicator; no selection count under Table or List.

### FDB-6

Animate only to tell the reader something, briefly, without ever blocking input, and drop movement when reduced motion is asked for.

- **Sources** (16 records):
  - Apple: [motion](https://developer.apple.com/design/human-interface-guidelines/motion) 1, 3, 4, 5, 6, 7, 2; [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 31
  - Windows 7: [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 11, 12; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 8, 38; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 23
  - Fluent: [signature-experiences-motion](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/motion) 4; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 17; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 13
- **Enforce:**
  - Component: Theme reduced motion (exists).
  - Review.
- **kvit-ui now:** enforces it. Reduced motion is honoured (`window.go:243-252`, `confirminplace.go:84-95`, `progress.go:89-93`).

### FDB-7

Every edit can be undone back to the last open or save, small related edits undo as one step, and undo shows the change it reverted. Notes-like data (kvit-notes, kvit-works) saves itself; kvit-cash asks on closing with unsaved changes (Save, Don't save, Cancel) on every system and has Save on its standard shortcut (decided, Q-12).

- **Sources** (7 records):
  - Apple: [edit-menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) 8; [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 7; [undo-and-redo](https://developer.apple.com/design/human-interface-guidelines/undo-and-redo) 2, 3, 4; [drag-and-drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) 5
  - Windows 7: [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 2
- **Enforce:**
  - Component: An UndoStack in kvit-ui with named entries, feeding the Edit menu, ConfirmInPlace and Toast.
  - Component: Window gains AllowClose, a hook that can hold a close to ask about unsaved changes.
- **kvit-ui now:** enforces it. `UndoStack` keeps named edits for the Edit menu, a Toast and a ConfirmInPlace (`undo.go`), and `Window.AllowClose` can hold a close to ask about unsaved changes (`window.go`).

### FDB-8

Mark items being inserted, removed or reordered so the change can be followed.

- **Sources** (3 records):
  - Apple: [collections](https://developer.apple.com/design/human-interface-guidelines/collections) 6
  - Fluent: [ui-controls-listview-filtering](https://learn.microsoft.com/en-us/windows/apps/design/controls/listview-filtering) 5; [ui-controls-selection-modes](https://learn.microsoft.com/en-us/windows/apps/design/controls/selection-modes) 5
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### FDB-9

Work the reader started keeps going: while the window is in the background, after the window that started it closes, and past recoverable problems, which are collected and shown together at the end. Ask for everything long work needs before it starts, and offer to apply an answer to all items when a bulk operation must ask. Long work shows its progress without blocking the app, disabling only what needs its result; never freeze waiting on a device or the network. Where possible show the result at once and finish the rest in the background. When the window loses focus, only what scores or times the reader pauses, resuming where it was; everything else keeps running (decided, Q-38).

- **Windows:** A modeless progress window can be minimised and has its taskbar button.
- **Sources** (18 records):
  - Apple: [multitasking](https://developer.apple.com/design/human-interface-guidelines/multitasking) 2, 3, 4, 9; [going-full-screen](https://developer.apple.com/design/human-interface-guidelines/going-full-screen) 6
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 68, 69, 70, 73, 74, 75, 66; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 30; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 34, 35; [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 37
  - Fluent: [ui-controls-progress-controls](https://learn.microsoft.com/en-us/windows/apps/design/controls/progress-controls) 6, 13
- **Enforce:**
  - Review: A modal dialog over long work; a recoverable error that ends the task; a question per item.
- **kvit-ui now:** does not cover it. Not covered; the library has no task or job runner.

### FDB-10

Progress details only when the reader can use them: what is happening in words (WRD-25), time remaining once the estimate is steady, refreshed at least every five seconds, rather than elapsed time; no percentage text beside a bar that shows it; extra details such as a transfer rate only when people can act on them and left up long enough to read; identifiers and support detail go to a log. A wait of one to two minutes shows something useful; a longer one says how long and suggests doing something else.

- **Sources** (13 records, 5 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 81, 82; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 9\*, 10, 16, 23\*, 36, 37, 51; [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 24; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 55\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 57\*, 59\*
- **Enforce:**
  - Component: Progress drops its percentage text by default and gains Remaining, shown only once the application says the estimate is steady.
  - Review: File names flickering past; a debug identifier on a progress page.
- **kvit-ui now:** enforces it. Progress shows no percentage unless asked, and says the time left, "About 3 minutes left" (`progress.go`).

## Layout

### LAY-1

Put the most important things at the top and the leading edge, in reading order, and keep them there at every size: when the window grows, content stays anchored at the top leading corner, never re-centred. Each window has one obvious place to look first, and emphasis goes only there. Everything a task needs is in one place, its required steps on the main path and optional ones off it, and each page serves one purpose.

- **Sources** (19 records):
  - Apple: [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 1; [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 6, 7; [typography](https://developer.apple.com/design/human-interface-guidelines/typography) 23; [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 4, 17
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 4, 30, 45, 46, 48, 2, 41, 42, 47; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 8, 7; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 48; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 14
- **Enforce:**
  - Review.
- **kvit-ui now:** partly. The shell keeps header, sidebar and body in fixed places (`window.go:14-28`); nothing checks the content.

### LAY-2

Keep related things together and set apart from others by space, a container or a separator, in at most about seven groups per window or page and never a group inside a group; two radio groups never sit side by side; align on shared edges and indent only what is subordinate.

- **Sources** (29 records):
  - Apple: [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 3, 4; [boxes](https://developer.apple.com/design/human-interface-guidelines/boxes) 1, 2; [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 19; [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 1
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 11, 9, 10, 44, 50, 53; [ctrl-group-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-group-boxes) 5, 1, 2, 3, 4; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 101; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 11, 12; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 25; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 16
  - Fluent: [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 11; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 16, 26; [basics-content-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/content-basics) 7; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 12; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 13; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 13
- **Enforce:**
  - Review.
- **kvit-ui now:** enforces it. Spacing scale and Region, Card, Divider (SKILL.md).

### LAY-3

Lay out by the space the window has, never by which system or screen it is, designing for a few width classes rather than each screen; every function stays reachable at every size. A window fits the screen's work area, has a minimum size at which it still works, and shows more when enlarged: a narrow window shows the least metadata and a list-and-details view one pane at a time, list first; a wide one adds columns and shows both side by side. A long form may take columns on a wide window but is split into pages rather than crowding one. Kvit's window floor is Px(880) by Px(600) on every system; the Windows 7 800 × 600 screen floor is not followed (decided, Q-25).

- **Sources** (33 records, 9 on Microsoft's checklist):
  - Apple: [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 15, 16, 17; [designing-for-macos](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) 1; [sidebars](https://developer.apple.com/design/human-interface-guidelines/sidebars) 7
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 174, 175, 172\*; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 12, 13, 14\*, 16\*, 1\*, 2\*, 17; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 11\*, 12\*, 1\*, 2\*; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 33
  - Fluent: [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 19; [layout-responsive-design](https://learn.microsoft.com/en-us/windows/apps/design/layout/responsive-design) 1, 4, 5, 6, 8, 9; [layout-index](https://learn.microsoft.com/en-us/windows/apps/design/layout/) 1; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 15, 16; [ui-controls-list-details](https://learn.microsoft.com/en-us/windows/apps/design/controls/list-details) 2; [layout-screen-sizes-and-breakpoints-for-responsive-design](https://learn.microsoft.com/en-us/windows/apps/design/layout/screen-sizes-and-breakpoints-for-responsive-design) 1, 2
- **Enforce:**
  - Component: WidthFloor and HeightFloor are Px(880) and Px(600) (`tokens/interface.go:195`, `tokens/interface.go:203`), set as the window's minimum (`window.go:88-92`).
  - Test, all screens of an app: Open every screen at WidthFloor by HeightFloor and fail any control out of reach or cut off.
- **kvit-ui now:** partly. Breakpoints are Measures (`tokens/interface.go:195-197`) and the window's minimum is Px(880) by Px(600) (`tokens/interface.go:195`, `tokens/interface.go:203`, set at `window.go:88-92`), as decided; nothing shows a list and its details one at a time on a narrow window.

### LAY-4

Mirror the layout for right-to-left languages: leading and trailing swap, direction-bound symbols mirror, numbers and logos do not, nor do a media seek bar and playback controls. Text aligns by its own direction, so a block of more than two or three lines in the other direction keeps its own alignment; a string running against the interface gets direction marks so the punctuation at its ends stays put, and paired brackets show correctly. People can type in either direction, and stored text keeps its direction.

- **Sources** (36 records):
  - Apple: [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 2; [right-to-left](https://developer.apple.com/design/human-interface-guidelines/right-to-left) 1, 2, 3, 5, 6, 7, 8, 9, 11, 12, 13, 16, 17, 18; [keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) 13
  - Windows 7: [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 9
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 22; [usability-index](https://learn.microsoft.com/en-us/windows/apps/design/usability/) 2; [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 2; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 12, 20, 21; [ui-controls-title-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/title-bar) 6; [globalizing-design-for-bidi-text](https://learn.microsoft.com/en-us/windows/apps/design/globalizing/design-for-bidi-text) 1, 2, 4, 5, 7, 9, 10, 11, 12, 13, 14, 15
- **Enforce:**
  - Component: Layouts by leading and trailing.
- **kvit-ui now:** partly. The text package orders right-to-left text (CLAUDE.md); layouts are left-to-right; Field's typed text is drawn by unison's own engine (`field.go:29-30`), which does not order right-to-left text.

### LAY-5

Buttons that act on one view sit inside or under it; the buttons that finish a task sit at the bottom trailing corner, in the same place on every page of a flow; switches, checkboxes, radio buttons, path bars and help buttons belong in the content, not in a toolbar or status bar.

- **Windows:** A single button on a page in the window (not a dialog) aligns at the leading edge with the page's content; a lone button in a popover, notification or list row sits at its trailing end; a dialog's buttons stay at the bottom trailing corner.
- **macOS:** As stated.
- **Sources** (18 records):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 20, 21, 22, 23, 28, 30, 33; [path-controls](https://developer.apple.com/design/human-interface-guidelines/path-controls) 1; [toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) 9
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 43, 52; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 9; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 45; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 20, 21
  - Fluent: [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 22, 23; [ui-controls-selection-modes](https://learn.microsoft.com/en-us/windows/apps/design/controls/selection-modes) 8
- **Enforce:**
  - Review.
- **kvit-ui now:** partly. Dialog puts its buttons at the bottom right (`feedback.go:541-552`); nothing checks app screens.

### LAY-6

A field's width matches what it expects; stack fields and radio buttons in a column, related ones the same width, each label clearly with its field and on its text baseline when beside it. In a settings or expander row the label is at the leading side and the control at the trailing edge. Numbers, money and times align to the trailing edge, in fields and columns alike; text, identifiers and dates to the leading edge; a column heading aligns like its data.

- **Windows:** A field's label sits above it, aligned with its leading edge; a list's label sits above it too.
- **macOS:** As stated.
- **Sources** (37 records, 1 on Microsoft's checklist):
  - Apple: [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 5, 6, 7; [toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) 21
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 51, 54, 12, 13, 37; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 32, 48, 43; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 24, 25, 50, 41\*, 43; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 38, 32; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 15; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 14; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 40, 50; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 60; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 42; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 26; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 40, 61
  - Fluent: [basics-content-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/content-basics) 6; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 7, 13; [ui-controls-slider](https://learn.microsoft.com/en-us/windows/apps/design/controls/slider) 29; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 24; [ui-controls-password-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/password-box) 4; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 10; [ui-controls-list-details](https://learn.microsoft.com/en-us/windows/apps/design/controls/list-details) 4; [signature-experiences-typography](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/typography) 2
- **Enforce:**
  - Component: A form row that draws a field's label (above it on Windows and Linux) and links it to the field.
  - Review.
- **kvit-ui now:** enforces it. FormRow puts a field's label above it on the Windows form and beside it on the macOS form (`formrow.go`); SettingRow puts the label at the leading side and the control at the trailing edge (`settings.go`).

### LAY-7

Content needed now and then goes in a popover beside its control, sized to its content; a dialog about one object opens near it.

- **Windows:** A dialog about one object opens beside it, offset down and to the right so the object stays visible.
- **macOS:** Content about one object is a popover pointing at it.
- **Sources** (7 records, 1 on Microsoft's checklist):
  - Apple: [popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) 2, 3, 10
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 36\*
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 9; [ui-controls-dialogs-and-flyouts-index](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/) 7, 8
- **Enforce:**
  - Component: Popover (exists); Dialog opens beside an anchor it is given.
- **kvit-ui now:** partly. Popover is anchored to a control (catalog KvitPopover); Dialog opens in the middle of the window unless MoveTo places it (`feedback.go:462-469`, `dialogmove.go:16-21`).

### LAY-8

Show a chart only when it makes a point, chosen by the question it answers, with its takeaway in words, honest axes and round ticks; show the data as a table when people only need the values.

- **Sources** (22 records):
  - Apple: [charts](https://developer.apple.com/design/human-interface-guidelines/charts) 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 17, 18, 26, 27; [charting-data](https://developer.apple.com/design/human-interface-guidelines/charting-data) 1, 2, 3, 6, 7, 8, 10
- **Enforce:**
  - Component: Data components.
- **kvit-ui now:** enforces it. Data components (catalog Data); an empty chart is an EmptyState (SKILL.md).

### LAY-9

Size controls, columns, lists and windows for their usual content: nothing people need is cut short where there is room, a button's words are never cut (it grows, wraps or its row stacks), a slightly larger control is better than a scroll bar, and running text keeps to 50 to 60 characters a line, never under 20. Leave a third more room (up to double for short words) for translated text, none for numbers, and size line and row heights for non-Latin scripts too. A window opens at the size its typical content needs, not near full screen; a dialog or a page of a flow never scrolls as a whole, and a flow keeps one page size. A list shows whole rows, three to about twenty.

- **Sources** (51 records, 10 on Microsoft's checklist):
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 18, 176\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 14, 15, 16, 18, 19; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 9, 15, 11\*; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 8, 35\*, 36\*, 31, 38, 39; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 28, 29\*; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 9\*, 42; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 7\*, 8\*, 9\*; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 35, 37, 36; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 28, 41; [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 12\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 87; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 41, 42, 43; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 33, 14, 34
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 21; [ui-controls-checkbox](https://learn.microsoft.com/en-us/windows/apps/design/controls/checkbox) 10; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 15, 33; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 18, 19; [ui-controls-text-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/text-box) 19; [ui-controls-combo-box](https://learn.microsoft.com/en-us/windows/apps/design/controls/combo-box) 13, 15; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 11; [ui-controls-lists](https://learn.microsoft.com/en-us/windows/apps/design/controls/lists) 14, 15, 16; [signature-experiences-typography](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/typography) 10, 11
- **Enforce:**
  - Component: Table: column widths from the data, with a minimum and a maximum; Label wraps instead of eliding where the layout has room.
  - Test, all screens of an app: Render every screen at its default size in English and in a pseudo-translation a third longer; fail any Label, Button or cell cut short while its row has room, and any dialog that scrolls.
- **kvit-ui now:** partly. Label is cut short with … rather than wrapped unless Wrap is set (`label.go:12-15`); Button words never elide (`button.go:119-123`); table columns default to Px(120), shrink to Px(48) at least and have no maximum (`table.go:236-239`, `table.go:496`); a cut-short cell shows its whole value on hover (`tablebody.go:200-252`); the window opens at Px(1440) by Px(960) (`window.go:117`, `tokens/interface.go:197`).

### LAY-10

Choose a flat or hierarchical structure by the content: fewer than eight pages viewed in any order, with no parent and child, are peers; more pages, a set order or a parent and child make a hierarchy, and peer links stay within one branch. Top-level sections go in a sidebar, or a row along the top for five or fewer, and the content is headed by the selected section's name. Back follows the reader's own history within the window, sits at the top leading corner, is disabled when there is nothing to go back to, and is never the only way to a place. Choosing the section already shown, switching between peers in one navigation element, moving through a list-and-details page and opening a dialog or popup are not history steps; Back first closes an open dialog, popup or a sidebar open over the content. A hierarchy deeper than two levels shows a breadcrumb: the path from the root, every crumb but the last navigates, the last does nothing, and crumbs that do not fit collapse into an ellipsis that lists them; two levels need only Back. Related pages link to each other directly, not only through their parent.

- **Sources** (29 records):
  - Fluent: [ui-controls-navigationview](https://learn.microsoft.com/en-us/windows/apps/design/controls/navigationview) 6, 14, 16, 17, 2, 3, 8, 13; [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 16; [ui-controls-buttons](https://learn.microsoft.com/en-us/windows/apps/design/controls/buttons) 4, 24; [ui-controls-breadcrumbbar](https://learn.microsoft.com/en-us/windows/apps/design/controls/breadcrumbbar) 1, 2, 4, 5, 6; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 5, 6, 7, 8, 9, 14, 18, 19, 20, 21, 15; [basics-app-silhouette](https://learn.microsoft.com/en-us/windows/apps/design/basics/app-silhouette) 1; [layout-responsive-design](https://learn.microsoft.com/en-us/windows/apps/design/layout/responsive-design) 10
- **Enforce:**
  - Component: Sidebar and Breadcrumb (exist); a back history per window with a Back button at the top leading corner, disabled when empty; a breadcrumb ellipsis that lists the hidden crumbs.
  - Test, one screen: Fail a Breadcrumb whose last crumb can be followed, and a Back button enabled with an empty history.
  - Review: Flat or hierarchical; which moves are history steps; direct links between related pages.
- **kvit-ui now:** enforces it. `History` follows the reader's own steps, and its Back button and standard command (Alt+Left, Cmd-[) are disabled with nowhere to go back to (`history.go`); the breadcrumb's "…" lists the places it folds (`breadcrumb.go`).

## Accessibility

### A11Y-1

Every element has a name a screen reader says, and the name is the words on screen, so voice control can use them. A control without words is named by what it does, in at most 150 characters, without "Click to": the role is spoken already.

- **Sources** (16 records, 2 on Microsoft's checklist):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 12, 20; [icons](https://developer.apple.com/design/human-interface-guidelines/icons) 9
  - Windows 7: [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 61, 3, 14, 20; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 11; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 28\*; [ctrl-spin-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-spin-controls) 19\*; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 31, 32, 33; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 26
  - Fluent: [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 2, 14
- **Enforce:**
  - Test, one screen: CheckNamed (exists); add: a Button's name equals its visible words, and no name is over 150 characters or starts with Click.
- **kvit-ui now:** enforces it. `uitest.CheckNamed`; IconButton's Label is its name and tooltip, a Button's name is its words, and `uitest.Labels` fails a name starting with Click (`uitest/labels.go`).

### A11Y-2

No distinction rests on colour alone, nor on bold or italic alone: the words make sense in plain type. Colour pairs, decorative ones included, stay distinguishable under common colour-vision deficiencies.

- **Sources** (17 records):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 10, 11; [color](https://developer.apple.com/design/human-interface-guidelines/color) 8; [toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) 4; [charts](https://developer.apple.com/design/human-interface-guidelines/charts) 19, 20
  - Windows 7: [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 17; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 14; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 4, 19; [vis-color](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-color) 1, 4
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 20; [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 13, 14; [signature-experiences-color](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/color) 10; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 12
- **Enforce:**
  - Component: Existing second channels.
  - Test, Go source: palette tests (exist), including colour-vision simulation (`palette/palette_test.go:42`).
- **kvit-ui now:** enforces it. Library rule (CLAUDE.md, The library's rules; SKILL.md); a checked menu line is a tick as well as bold (`menu.go:40-42`); it is to keep the tick and lose the bold, which marks a context menu's default item (decided, Q-33).

### A11Y-3

Text and controls meet contrast minimums in every theme: text at least 5:1 wherever the palette tests show it costs little, else 4.5:1, and 3:1 for large text (18 pt, or 14 pt bold). High contrast is Kvit's own high-contrast table, switched on by the desktop's high-contrast setting; the title bar follows high contrast and light or dark (decided, Q-50).

- **Sources** (19 records):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 6, 7, 8; [color](https://developer.apple.com/design/human-interface-guidelines/color) 9; [dark-mode](https://developer.apple.com/design/human-interface-guidelines/dark-mode) 3, 6
  - Windows 7: [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 44; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 6; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 8; [vis-color](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-color) 3, 10
  - Fluent: [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 4; [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 9, 12; [signature-experiences-color](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/color) 8; [style-acrylic](https://learn.microsoft.com/en-us/windows/apps/design/style/acrylic) 5, 6; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 9, 11
- **Enforce:**
  - Test, Go source: Theme contrast tests (exist); raise the text floor to 5:1 for each pair the palette tests show can take it.
  - Component: The desktop's high-contrast setting switches to Kvit's high-contrast table (exists).
- **kvit-ui now:** partly. Every token pair meets WCAG AA, 4.5:1 for text (`tokens/theme_test.go:168-200`), with a stricter floor for the high-contrast table (`tokens/theme_test.go:278`); the desktop's high contrast switches to Kvit's own high-contrast table (`tokens/theme_test.go:351`), as decided; no pair is held to 5:1 yet.

### A11Y-4

The interface works at every interface size: headings stay above body text, layouts reflow rather than clip, and information-carrying symbols grow with text.

- **Windows:** Also test at 100, 125 and 150 percent display scaling for clipped controls and text and stretched images. The interface, title bar included, also follows the system's text-size setting.
- **macOS:** The general rule.
- **Sources** (24 records, 3 on Microsoft's checklist):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 3; [typography](https://developer.apple.com/design/human-interface-guidelines/typography) 4, 7, 9, 16, 17, 18, 19, 21, 22, 23; [layout](https://developer.apple.com/design/human-interface-guidelines/layout) 11, 13, 16; [text-views](https://developer.apple.com/design/human-interface-guidelines/text-views) 3
  - Windows 7: [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 21\*; [win-window-mgt](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-window-mgt) 3\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 3\*; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 7; [vis-fonts](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-fonts) 5
  - Fluent: [basics-titlebar-design](https://learn.microsoft.com/en-us/windows/apps/design/basics/titlebar-design) 13; [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 11; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 1, 10
- **Enforce:**
  - Test, one screen: TestPagesFollowTheInterfaceSize (exists); add the 125 and 150 percent display scales to the screenshot set.
- **kvit-ui now:** partly. Sizes are Measures read at layout time (CLAUDE.md); Windows' text-size setting is not read (nothing in `platform` asks for it).

### A11Y-5

Never carry meaning by sound or animation alone. Nothing flashes between 2 and 55 times a second; the app's sounds stay between 500 and 3000 Hz and follow the system's volume.

- **Sources** (9 records):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 13, 15; [feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) 2; [notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) 14; [charts](https://developer.apple.com/design/human-interface-guidelines/charts) 16
  - Windows 7: [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 27, 26, 37, 25
- **Enforce:**
  - Review.
- **kvit-ui now:** enforces it. Toast and ConfirmInPlace announce themselves (`feedback.go:120`, `confirminplace.go:122-131`); the library plays no sound and nothing it draws flashes.

### A11Y-6

Charts are readable by a screen reader: an overview of kind and axes, each mark named context first, series named by meaning, a sensible keyboard route through values.

- **Sources** (11 records):
  - Apple: [charting-data](https://developer.apple.com/design/human-interface-guidelines/charting-data) 4, 9; [charts](https://developer.apple.com/design/human-interface-guidelines/charts) 11, 14, 15, 22, 23, 24, 25, 28, 30
- **Enforce:**
  - Component: Data components (exist).
- **kvit-ui now:** enforces it. Data components give screen-reader phrasing (`data.go:146-148`, `data.go:606-612`).

### A11Y-7

Treat accessibility as a requirement from the start: run the automated checks on every build, and test the core tasks by hand with a screen reader, keyboard only and the system's accommodations, and with people who rely on them.

- **Sources** (16 records):
  - Apple: [accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) 1, 22, 26; [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 17, 18; [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 13
  - Windows 7: [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 9; [vis-color](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-color) 2
  - Fluent: [usability-index](https://learn.microsoft.com/en-us/windows/apps/design/usability/) 1; [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 1, 3, 4, 16, 17; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 3, 15
- **Enforce:**
  - Test, one screen: uitest's screen-reader view (exists).
  - Review: Core tasks tried with a screen reader and keyboard only before release.
- **kvit-ui now:** partly. uitest gives a screen reader's view (CLAUDE.md) and its checks run in `./build.sh --test`; no testing with people is arranged.

### A11Y-8

What the screen shows, the screen reader gets: text is drawn as text (never inside a picture, built from characters or drawn over a progress bar); each label is tied to its control and never wraps around it as a sentence; a group shown by placement is declared as a group; each option in a group reports whether it is chosen and its place (2 of 5); a notice or badge that appears or changes is announced, a badge as part of the focusable control it sits on; and every focus change raises the event screen readers follow.

- **Sources** (14 records):
  - Windows 7: [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 62, 21, 15; [progress-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/progress-bars) 17; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 56, 55, 57
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 15; [ui-controls-info-badge](https://learn.microsoft.com/en-us/windows/apps/design/controls/info-badge) 10, 11; [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 6; [ui-controls-radio-button](https://learn.microsoft.com/en-us/windows/apps/design/controls/radio-button) 21; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 11; [accessibility-accessibility-overview](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/accessibility-overview) 8
- **Enforce:**
  - Component: Field and Select gain a visible label that sets LabeledBy (today Label is only the accessible name, `field.go:34-35`); containers that group by placement take a Group role. Notice announces itself when shown or changed; a badge's count is in its control's name and announced when it changes; RadioGroup options report their position.
  - Test, one screen: Every visible label beside a control is that control's name or LabeledBy; groups are nodes in the tree.
  - Review: Words in pictures; text-character drawings.
- **kvit-ui now:** partly. FormRow names its field by the label the reader sees (`formrow.go`); no component sets LabeledBy; Notice is a Group named by its text.

## Consistency

### CON-1

One command or feature has one name, symbol and shortcut, spelled and capitalised the same, everywhere it appears: menu, button, tooltip, tray, shortcut list, help. Two different commands never share a name. A window a command opens is titled with the command's words; on Windows a dialog's title is its main instruction (WIN-6), which uses the command's verb where it names the action.

- **Sources** (11 records, 1 on Microsoft's checklist):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 2, 10; [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 2
  - Windows 7: [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 20; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 34\*; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 43, 17; [text-style-tone](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-style-tone) 17; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 55; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 8
  - Fluent: [ui-controls-commanding](https://learn.microsoft.com/en-us/windows/apps/design/controls/commanding) 6
- **Enforce:**
  - Component: A Command type gives one name, symbol and shortcut to every place it is shown (menu line, Button, IconButton tooltip, tray item, shortcut list) and is the title of the dialog it opens.
  - Test, all screens of an app: Across all screens of an app, collect (handler or Command, name): fail one command with two names, one name on two different commands, and one name capitalised two ways.
  - Review: A Dialog's Title uses the words of the command that opened it; on Windows, as the main instruction.
- **kvit-ui now:** enforces it. A Command gives one name, symbol and shortcut to every place it is shown (`command.go`); `uitest.App` fails one command written two ways, or with two symbols or two sets of keys, across an app's screens in one form (`uitest/app.go`).

### CON-2

One symbol has one meaning in an app, and a standard action uses its standard symbol.

- **Sources** (24 records):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 10; [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 8; [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 13; [icons](https://developer.apple.com/design/human-interface-guidelines/icons) 11, 12, 14, 15, 16, 17, 18, 19; [rating-indicators](https://developer.apple.com/design/human-interface-guidelines/rating-indicators) 2; [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 13
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 87, 88; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 21; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 15, 22, 24, 26; [vis-std-icons](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-std-icons) 17, 18; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 12
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 17
- **Enforce:**
  - Component: The icon catalogue lists each name's meaning (`icons/catalog_gen.go`); a symbol used for a meaning not listed logs a warning.
  - Test, all screens of an app: Across all screens of an app, collect (symbol, name); fail a symbol with two unrelated names.
- **kvit-ui now:** partly. `uitest.App` fails one symbol used under two names across an app's screens (`uitest/app.go`); a symbol used for a meaning the catalogue does not list is not warned about.

### CON-3

A component behaves the same everywhere; a custom-drawn control (such as the Kvit menu) answers to the same keys and clicks as the system's.

- **Sources** (12 records):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 10, 11; [menus](https://developer.apple.com/design/human-interface-guidelines/menus) 26; [color-wells](https://developer.apple.com/design/human-interface-guidelines/color-wells) 1; [context-menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) 3; [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 8, 21
  - Windows 7: [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 21; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 33; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 11
  - Fluent: [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 5, 15
- **Enforce:**
  - Component: Behaviour lives in components, not apps (SKILL.md, workflow step 3).
- **kvit-ui now:** enforces it. The Kvit menu takes arrows, Home, End, Return, Escape and access keys (`menu.go:489-564`).

### CON-4

Feel native on each platform. Keep content and controls in the same places across platforms, sizes and similar windows. Do similar things the same way: a recurring field keeps its label and place, a similar choice uses the kind of control its neighbours use, required fields are marked one way, with an asterisk after the label on every system (decided, Q-51), and similar objects all have descriptive tips or none do.

- **Sources** (22 records):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 14, 16, 22; [toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) 25; [search-fields](https://developer.apple.com/design/human-interface-guidelines/search-fields) 16
  - Windows 7: [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 5; [ctrl-drop](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-drop) 7; [ctrl-radio-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-radio-buttons) 6, 8; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 9; [ctrl-spin-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-spin-controls) 7, 10; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 158; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 17, 16; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 15; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 30; [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 1
  - Fluent: [ui-controls-command-bar](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar) 8; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 1; [ui-controls-forms](https://learn.microsoft.com/en-us/windows/apps/design/controls/forms) 6; [design-principles](https://learn.microsoft.com/en-us/windows/apps/design/design-principles) 4
- **Enforce:**
  - Component: Platform forms (library-gaps.md, Platform knowledge); Field gets one Required flag, drawn as an asterisk after the label on every system and given to screen readers as required, so an app cannot mix two markings.
  - Test, all screens of an app: Across all screens of an app, collect (field label, position in its form, control kind): fail one label at two positions or as two kinds of control.
  - Review: Similar windows share one layout; choices on one page use matching control kinds.
- **kvit-ui now:** partly. Field's `Required` is drawn as an asterisk after its FormRow label on every system and told to a screen reader (`field.go`, `formrow.go`); one label at two positions across an app is not flagged.

### CON-5

Follow the desktop's settings (appearance, contrast, motion, scrolling, hover and message timing) by default and do not copy them into app settings.

- **Windows:** Tooltip delay and display time are multiples of the system's hover time and message duration (SPI_GETMOUSEHOVERTIME, SPI_GETMESSAGEDURATION), not fixed times.
- **Sources** (14 records):
  - Apple: [dark-mode](https://developer.apple.com/design/human-interface-guidelines/dark-mode) 1, 2; [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 4, 5
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 18; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 10, 16; [vis-color](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-color) 11; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 36
  - Fluent: [style-mica](https://learn.microsoft.com/en-us/windows/apps/design/style/mica) 2; [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 8; [design-principles](https://learn.microsoft.com/en-us/windows/apps/design/design-principles) 3; [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 22; [signature-experiences-color](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/color) 1
- **Enforce:**
  - Component: Theme follows the desktop (exists). control's tooltip timing (`control.go:213`, `control.go:228-232`) and partTip's (`parttip.go:28-49`) read new platform.HoverTime and platform.MessageDuration on Windows, keeping today's values where the system has no such setting.
  - Test, Go source: No fixed time.Duration literal in tooltip timing outside the platform package.
- **kvit-ui now:** partly. Theme follows the desktop's dark mode, contrast and motion (`platform/appearance_*`); the wheel follows the system's lines per notch (`wheel.go:41`, `platform/wheel_windows.go`). Tooltip timing is fixed: 500 ms delay (`control.go:213`) and at least 3 s plus 60 ms a character (`control.go:228-232`), so on Windows it contradicts the Windows form.

### CON-6

A colour has one meaning across the app, one colour marks what can be acted on, and a colour given to data stays with that data when the theme changes. Readers may colour their own categories (labels, tags) from the chart ramp's named steps, and choose colours nowhere else (decided, Q-35).

- **Sources** (12 records, 1 on Microsoft's checklist):
  - Apple: [color](https://developer.apple.com/design/human-interface-guidelines/color) 1, 12, 10; [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 24
  - Windows 7: [vis-color](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-color) 18, 16, 17, 14; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 24\*; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 19; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 16
  - Fluent: [signature-experiences-color](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/color) 5
- **Enforce:**
  - Test, Go source: Ink names and rules_test (exist).
  - Test, Go source: palette tests: categorical step i keeps its hue family in all four colour tables.
  - Test, Go source: tokens tests: Accent and Link keep one hue family in all four colour tables.
  - Component: A category colour picker offers only the chart ramp's named steps (palette categorical steps); no free colour choice elsewhere.
- **kvit-ui now:** partly. Inks are named by meaning (SKILL.md, What each colour means). Controls are marked with Accent and links with Link (`link.go:66-70`, `check.go:98`, `switch.go:134`), one hue family in each table (`tokens/tables_gen.go:35,44,96,105,157,166,218,227`); no test checks that. A series keeps its categorical step in every theme (`tokens/theme.go:183-188`), and the four tables keep the same hue order with different lightness (`tokens/tables_gen.go:65,126,187,248`); no test checks the hue order. Readers cannot yet colour their own categories: there is no picker limited to the chart ramp's named steps (decided, Q-35).

### CON-7

Charts with the same purpose use the same type and wording across the app; familiar chart forms first.

- **Sources** (4 records):
  - Apple: [charting-data](https://developer.apple.com/design/human-interface-guidelines/charting-data) 5, 11, 12; [charts](https://developer.apple.com/design/human-interface-guidelines/charts) 29
- **Enforce:**
  - Review.
- **kvit-ui now:** enforces it. Fixed data vocabulary (catalog Data).

### CON-8

File names are shown without extensions unless the reader chose otherwise.

- **Sources** (3 records):
  - Apple: [file-management](https://developer.apple.com/design/human-interface-guidelines/file-management) 5; [icons](https://developer.apple.com/design/human-interface-guidelines/icons) 27
  - Windows 7: [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 26
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

## Help

### HLP-1

Every symbol-only control shows a short tooltip on hover and on keyboard focus: one string, its verb-led name plus the optional Explanation line (decided, Q-23). It stays while the pointer moves onto it, goes when pointer and focus have left, and Escape dismisses it without moving focus when no popup is open; with a popup open, Escape closes the popup (decided, Q-49).

- **Windows:** Ctrl also dismisses an open tooltip (Windows 11). The tooltip adds the control's shortcut in parentheses after its name (decided, Q-23).
- **Sources** (18 records, 4 on Microsoft's checklist):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 9; [segmented-controls](https://developer.apple.com/design/human-interface-guidelines/segmented-controls) 12; [sliders](https://developer.apple.com/design/human-interface-guidelines/sliders) 10
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 1, 2\*, 15, 28\*; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 18\*; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 49\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 38, 82; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 26; [ctrl-status-bars](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-status-bars) 20, 24
  - Fluent: [ui-controls-tooltips](https://learn.microsoft.com/en-us/windows/apps/design/controls/tooltips) 1, 2, 5; [ui-controls-command-bar-flyout](https://learn.microsoft.com/en-us/windows/apps/design/controls/command-bar-flyout) 8
- **Enforce:**
  - Component: IconButton (exists). control keeps its tooltip while the pointer crosses to it, and the window's popupKeys hides a shown tooltip on Escape, and on Windows and Linux on Ctrl, before the key goes on to anything else. On Windows and Linux the tooltip adds the Command's shortcut in parentheses after the name.
  - Test, one screen: Every IconButton has a tooltip; with no popup open, a shown tooltip goes on Escape (and Ctrl on Windows) with focus unmoved, and stays while the pointer moves from the control onto it; on Windows a control with a shortcut shows it in parentheses.
- **kvit-ui now:** partly. A symbol button's tooltip is its name, with its command's shortcut in parentheses on the Windows form (`iconbutton.go`, `command.go`); `uitest.Labels` fails a symbol button without one (`uitest/labels.go`); Escape and Ctrl do not hide a shown tooltip.

### HLP-2

A tooltip is short, static, supplemental text beside what it explains, shown only when the reader hovers, focuses or presses, never on its own (UI.ShowTooltip is for the gallery and tests only, decided Q-52): never over it, nor over the next item in its row or column. Its text wraps, never cut short. Nothing in it can be clicked or looks clickable, and it holds no error or changing status (a tray icon's tip excepted) and nothing essential or worth printing. A control that needs a long tooltip needs redesigning; a longer explanation goes in a Hint the reader opens. A link's tooltip says where it leads (CTL-26); any further explanation goes beside the link or in that tooltip, not both.

- **Windows:** Objects such as files and list items may also show a longer infotip describing them on hover; it is supplemental, never essential, and where such tips could distract the reader can turn them off.
- **Sources** (35 records, 4 on Microsoft's checklist):
  - Apple: [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 23
  - Windows 7: [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 3, 5, 6, 7, 21, 19\*, 20\*, 4, 9, 8, 35; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 25; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 46, 44; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 6; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 51\*, 52\*; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 22; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 4
  - Fluent: [ui-controls-annotated-scrollbar](https://learn.microsoft.com/en-us/windows/apps/design/controls/annotated-scrollbar) 5; [in-app-help-in-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/in-app-help) 13; [ui-controls-tooltips](https://learn.microsoft.com/en-us/windows/apps/design/controls/tooltips) 3, 4, 6, 7, 8, 9, 10, 12, 13, 16, 17, 18; [ui-controls-dialogs-and-flyouts-flyouts](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/flyouts) 2
- **Enforce:**
  - Component: Hint (exists). Tooltip: newTooltip takes words only, so nothing clickable can go in it (exists). PlaceBeside (`layer.go:212-246`) learns whether the anchor sits in a row (siblings across) and then places below, keeping trailing-side placement for a column.
  - Component: Explanation and IconButton labels over about 60 characters log a warning (WRD-13's length).
  - Test, one screen: On a screen with a row of IconButtons, hover each and assert its tooltip does not intersect the next button.
  - Component: UI.ShowTooltip's doc says it is for the gallery and tests only and apps must not call it.
  - Test, Go source: No call to UI.ShowTooltip outside the gallery and tests.
  - Review: No errors, status or essential facts in tooltips.
- **kvit-ui now:** partly. Hint exists (`popover.go:160-213`); a tooltip is passive words only, never focused (`layer.go:27-30`, `control.go:389-420`), and wraps at Px(280), never cut (`control.go:403`), but its length is unlimited. UI.ShowTooltip lets an app show one unasked and its doc does not forbid apps calling it (`popover.go:11-16`); only the gallery calls it (`cmd/kvit-ui-gallery/specimens_feedback.go:14`). Placement takes the trailing side, level with the top, first (`layer.go:218-241`), which in a row of icon buttons covers the next button: it contradicts the row case. partTip's hover card serves as the infotip (`parttip.go:51-67`).

### HLP-3

Put what people need in the window itself: explain with clear labels first, then a short line of text, and a help link last; a page or message is complete without help. Help in the window concerns what the reader is doing now, matches the input they use, explains what a control does in this app, and closes in one step.

- **Sources** (29 records, 2 on Microsoft's checklist):
  - Apple: [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 1, 2, 3, 4, 6, 7, 8, 9
  - Windows 7: [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 4, 30, 1, 3; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 68, 70\*; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 44; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 68; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 16; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 165; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 125\*
  - Fluent: [in-app-help-guidelines-for-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/guidelines-for-app-help) 1, 2; [in-app-help-external-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/external-help) 3; [in-app-help-in-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/in-app-help) 1, 2, 4, 8, 11, 12; [in-app-help-instructional-ui](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/instructional-ui) 6
- **Enforce:**
  - Component: Hint holds optional extra explanation, never the only place of a needed fact (doc comment on Hint).
  - Review: Labels, then static text, then help; messages and pages complete without help.
- **kvit-ui now:** partly. Hint closes on Escape or an outside click (`popover.go:159-213`); ShowTooltip's doc says a tooltip may not be the only place a meaning lives (`popover.go:11-15`). Nothing checks that a needed fact is in the window.

### HLP-4

Teaching tips use a Tip (anchored, non-blocking, dismissable); Spotlight stays for tours (decided, Q-48). Tips are few and short, shown near what they explain, only to people who have not used the feature, at most one at a time and spread across sessions, never selling anything. A tip points out a feature that is not obvious, or a new or changed one; never the basics, what people work out alone, or a complex feature, which help covers. It holds nothing needed that is not also elsewhere, no error or status, and no video, audio or large picture. It covers and blocks nothing and can be ignored; once shown it is not shown again, and a setting can bring tips back.

- **Windows:** A tip about one element points at it from above where it fits; a tip about nothing in particular sits centred at the bottom of the window, pointing at nothing. It has a close button (an X by its title, or a worded Close at the bottom) and at most one action button. In a scrolling area it has no close button and closes on a scroll or any other use. A tip whose close button cannot be shown is not shown.
- **Sources** (33 records):
  - Apple: [offering-help](https://developer.apple.com/design/human-interface-guidelines/offering-help) 10, 11, 12, 13, 14, 17; [onboarding](https://developer.apple.com/design/human-interface-guidelines/onboarding) 2, 3
  - Windows 7: [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 6
  - Fluent: [ui-controls-infobar](https://learn.microsoft.com/en-us/windows/apps/design/controls/infobar) 12; [in-app-help-guidelines-for-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/guidelines-for-app-help) 7; [ui-controls-dialogs-and-flyouts-teaching-tip](https://learn.microsoft.com/en-us/windows/apps/design/controls/dialogs-and-flyouts/teaching-tip) 1, 2, 3, 4, 5, 8, 9, 10, 11, 13; [in-app-help-in-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/in-app-help) 9, 10; [in-app-help-instructional-ui](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/instructional-ui) 1, 2, 3, 4, 5, 8, 9, 10, 11, 12
- **Enforce:**
  - Component: New Tip: a passive card that points at its target (or sits centred at the window's bottom), with a close button and at most one action, closing on scroll inside a scroll panel; it never shades the window, and records by id that it was shown.
  - Test, one screen: With a Tip shown, a control outside it is still pressable; a Tip inside a scroll panel closes on scroll.
  - Review: What tips teach and how often they appear.
- **kvit-ui now:** enforces it. Tip is a card beside what it explains, or at the bottom of the window, with a close button and at most one action; it blocks nothing, closes when what it points at moves, and is not shown again once closed, until `ResetTips` (`tip.go`).

### HLP-5

First run asks only what is required and has no safe default; everything else starts from good defaults that can be changed in settings, and a feature is set up in context when first used. An introduction is short, lets people try the thing, can be skipped, stays reachable later, explains the app not the system, goes from simple personal choices to technical ones, and shows one surface at a time.

- **Sources** (21 records):
  - Apple: [onboarding](https://developer.apple.com/design/human-interface-guidelines/onboarding) 1, 4, 6, 7, 10, 11; [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 9; [charting-data](https://developer.apple.com/design/human-interface-guidelines/charting-data) 6
  - Windows 7: [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 13, 2, 4, 9, 11, 12, 5, 15, 14, 16, 7, 25; [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 36
- **Enforce:**
  - Review: First-run screens hold only required choices; setup of secondary features waits for their first use.
- **kvit-ui now:** does not cover it. Not covered: kvit-ui has no first-run or setup component (survey §11; Spotlight is only a highlight, catalog KvitSpotlight).

### HLP-6

Help is reached from where it is needed, the same way throughout the app, and opens the one topic for that place, leaving the window open. Add it only where a topic says something the window cannot: never a search, a list of links or an FAQ, and never a generic Help or Troubleshoot button. Help topics are HTML pages shipped with the app and opened in the system browser at the topic (decided, Q-34); help on the web is only for detailed tutorials and topics too large for the app, and a link to it says it goes online.

- **Windows:** A text link worded as the question it answers, a complete sentence ending only in a question mark, saying when it goes online; at the bottom of the content area, or under the group it concerns; never a button, a disclosure, a "Learn more" link or a title-bar question mark. Where no topic exists for the place, a link to a help home organised by categories is acceptable. (decided, Q-34)
- **macOS:** The round question-mark button, uncaptioned, one per window, opening the topic for that window.
- **Sources** (39 records, 10 on Microsoft's checklist):
  - Apple: [buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) 24, 25, 26, 27, 29; [alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) 37
  - Windows 7: [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 166, 167\*; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 69\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 69; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 12, 19, 10\*, 11\*, 33, 13, 18; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 45\*, 41, 46, 47; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 51; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 127\*, 129\*, 130\*, 131\*, 126\*; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 17, 18
  - Fluent: [in-app-help-guidelines-for-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/guidelines-for-app-help) 4, 8, 9; [usability-index](https://learn.microsoft.com/en-us/windows/apps/design/usability/) 3; [in-app-help-external-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/external-help) 1, 2, 4, 5; [in-app-help-in-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/in-app-help) 6, 3
- **Enforce:**
  - Component: New HelpLink: takes a topic id and words; on Windows and Linux it is a Link (`link.go`) placed at the bottom of its group, on macOS the round question-mark IconButton; Dialog (`feedback.go:382-400`) gets a Help slot that places it per platform. A topic id opens the shipped HTML page at its anchor in the system browser.
  - Test, Go source: HelpLink words end in "?" with no other closing punctuation and do not start with "Learn more", "Get help" or "Troubleshoot"; no Button labelled Help or Troubleshoot.
  - Test, all screens of an app: Every HelpLink's topic id resolves to an existing topic anchor; on macOS at most one help button per window.
- **kvit-ui now:** enforces it. HelpLink is a question-worded link on the Windows form and the round help button on the macOS form, opening `UI.HelpRoot`'s HTML page at its anchor (`helplink.go`); `Dialog.Help` puts it at the bottom of the content on the Windows form, reached by Tab after the buttons, and at the leading end of the button row on the macOS form (`dialog.go`). `uitest.Labels` fails help words that are not a question ending in a question mark or that start "Learn more", "Get help" or "Troubleshoot", a Button labelled Help or Troubleshoot, and more than one help button to a window or dialog on the macOS form; `uitest.App` fails a topic whose page or anchor is not in `UI.HelpRoot` (`uitest/labels.go`, `uitest/app.go`). Whether a topic says something the window cannot is the reviewer's.

### HLP-7

When text points to a setting, give a link that opens it rather than describing the path.

- **Sources** (3 records):
  - Apple: [writing](https://developer.apple.com/design/human-interface-guidelines/writing) 22
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 23; [in-app-help-instructional-ui](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/instructional-ui) 13
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### HLP-8

Help topics are HTML pages shipped with the app (HLP-6). They cover the top tasks in normal use, the way most people do them, not every feature or what the window already says. Each topic leads to something people can do and is written for scanning: numbered steps, tables for reference, subheadings, short lists. It names interface items by their exact words in bold, and messages by their main instruction. Topics are found through contents, an index and search under everyday words, are as accessible as the app, and one topic lists every keyboard shortcut and accessibility option. Tutorials cover only advanced work, are kept apart from task topics, and are linked from a list of tutorials and from the topics they relate to.

- **Windows:** Help names the kind of control with its words (the **Export** button, the **File** menu) (decided, Q-34) and uses Windows terms: select for buttons, commands and list items whatever the input (click, double-click, right-click and drag only where the action needs a mouse); select and clear a check box; context menu (not shortcut menu); unavailable (not dimmed or greyed); access keys (not mnemonics); pointer (not cursor); the menu bar; a tip (not tool tip); a message (not error message). Apart from select, the terms come from the Windows 7 guide.
- **macOS:** Help names an item by its words without naming its kind of control (decided, Q-34). Apple's help-writing terms were not among the sources; use Apple's own terms for the same items.
- **Sources** (89 records):
  - Windows 7: [inter-accessibility](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-accessibility) 12, 13, 63; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 60, 59; [winenv-help](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-help) 21, 2, 6, 20, 23, 26, 25, 7, 9, 24, 5, 38; [mess-error](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-error) 88, 89; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 56, 57; [mess-notif](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-notif) 53, 56; [ctrl-balloons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-balloons) 32, 34, 33; [exper-printing](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-printing) 45, 46; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 207; [mess-warn](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-warn) 42; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 48, 46, 47; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 93, 92; [ctrl-list-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-boxes) 62, 61; [ctrl-progressive-disclosure-controls](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-progressive-disclosure-controls) 38; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 97, 90, 91, 93, 95, 96; [ctrl-command-buttons](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-buttons) 55, 53; [ctrl-tree-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tree-views) 51, 47, 48, 50, 52; [ctrl-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-links) 44, 42; [ctrl-sliders](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-sliders) 35, 34; [cmd-menus](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-menus) 101, 93, 94, 95, 97, 98, 99, 100; [ctrl-command-links](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-command-links) 47, 44, 45, 46; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 55; [ctrl-tooltips-and-infotips](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-tooltips-and-infotips) 62, 63; [inter-mouse](https://learn.microsoft.com/en-us/windows/win32/uxguide/inter-mouse) 48, 49, 50, 51, 52, 53, 54, 55, 56; [ctrl-search-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-search-boxes) 21, 22
  - Fluent: [in-app-help-guidelines-for-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/guidelines-for-app-help) 3, 6; [in-app-help-external-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/external-help) 6, 7, 8, 9; [in-app-help-in-app-help](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/in-app-help) 5, 7
- **Enforce:**
  - Component: A keyboard-shortcuts topic generated from the app's Commands (CON-1), so it lists every shortcut the menus have.
  - Test, all screens of an app: The shortcuts topic lists every key binding found across the app's menus.
  - Review: Topic choice, structure, wording and terms.
- **kvit-ui now:** does not cover it. No help system, topics or shortcut list (survey §11).

## Other

### OTH-1

Examples, pictures and features suit everyone: varied people, everyday settings, no gendered figures for unspecified people, no assumptions about shared experiences. Colours and symbols are checked for what they mean in the cultures the app serves.

- **Sources** (10 records):
  - Apple: [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 1, 11, 12, 13, 14, 15, 16; [icons](https://developer.apple.com/design/human-interface-guidelines/icons) 6
  - Fluent: [signature-experiences-color](https://learn.microsoft.com/en-us/windows/apps/design/signature-experiences/color) 7; [iconography-index](https://learn.microsoft.com/en-us/windows/apps/design/iconography/) 6
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### OTH-2

Numbers, dates, times and money follow the reader's locale, with the long date form where the detail helps (a deadline) and the short one where space is tight; descriptions use real numbers and spelled-out dates and units.

- **Sources** (6 records):
  - Apple: [inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) 22; [text-fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) 12; [right-to-left](https://developer.apple.com/design/human-interface-guidelines/right-to-left) 4
  - Windows 7: [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 39, 38; [ctrl-text-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-text-boxes) 4
- **Enforce:**
  - Component: ui.Number, CountPhrase (exist); a UI.Date(t, Long|Short) from platform.Locale, used by the Date cell.
  - Test, Go source: No time.Format with a literal layout outside the date helper.
- **kvit-ui now:** partly. Locale grouping (`numbers.go`); TypeAhead counts with `CountPhrase` (`typeahead.go`); a Date cell is always drawn as ISO 2006-01-02 whatever the locale (`cell.go`).

### OTH-3

Collect only what is needed and keep private things off screens others may see.

- **Sources** (2 records):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 9; [searching](https://developer.apple.com/design/human-interface-guidelines/searching) 5
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not covered.

### OTH-4

Choose defaults that are safe, secure and private first and otherwise suit most people. Keep settings few, about four or five where possible: only the properties people need, meaningful bundles (a theme) rather than many single settings, and nothing that can be detected. Anything that shares information, changes the system or desktop, adds software or accepts legal terms is opt-in, with nothing preselected.

- **Sources** (24 records, 9 on Microsoft's checklist):
  - Apple: [settings](https://developer.apple.com/design/human-interface-guidelines/settings) 1, 2
  - Windows 7: [exper-first-exper](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-first-exper) 20\*, 21; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 12, 5, 6, 10\*, 18\*; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 59\*, 61; [top-violations](https://learn.microsoft.com/en-us/windows/win32/uxguide/top-violations) 84\*, 101\*; [win-dialog-box](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-dialog-box) 13, 170\*; [ctrl-check-boxes](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-check-boxes) 27\*; [ctrl-list-views](https://learn.microsoft.com/en-us/windows/win32/uxguide/ctrl-list-views) 25; [mess-confirm](https://learn.microsoft.com/en-us/windows/win32/uxguide/mess-confirm) 14; [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 50\*; [cmd-toolbars](https://learn.microsoft.com/en-us/windows/win32/uxguide/cmd-toolbars) 63; [winenv-notification](https://learn.microsoft.com/en-us/windows/win32/uxguide/winenv-notification) 7, 50
  - Fluent: [app-settings-guidelines-for-app-settings](https://learn.microsoft.com/en-us/windows/apps/design/app-settings/guidelines-for-app-settings) 9, 27
- **Enforce:**
  - Review: Defaults, the settings list and opt-ins.
- **kvit-ui now:** partly. Theme follows the desktop (CON-5); of its own the library keeps only theme, interface size and typography in the settings file, and apps add their own keys (`settings/store.go:1-4`). Opt-ins are app matters, not covered.

### OTH-5

Design judgment: serve the product's purpose, make the few common features excellent, let people move freely, and let nothing slow the task. Start with the simplest design that works and remove every element and word the main tasks do not need; fix a confusing feature by redesigning it, not by putting a guided flow in front of it; a simple task that needs a deeply nested or very wide structure is redesigned.

- **Sources** (28 records):
  - Apple: [design-principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) 1, 2, 3, 5, 20, 21, 23, 24, 25, 26; [designing-for-macos](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) 4
  - Windows 7: [win-wizards](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-wizards) 4; [vis-layout](https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-layout) 18; [how-to-design-desktop-ux](https://learn.microsoft.com/en-us/windows/win32/uxguide/how-to-design-desktop-ux) 2, 4, 9, 1, 3, 20; [win-property-win](https://learn.microsoft.com/en-us/windows/win32/uxguide/win-property-win) 62; [text-ui](https://learn.microsoft.com/en-us/windows/win32/uxguide/text-ui) 2
  - Fluent: [accessibility-designing-inclusive-software](https://learn.microsoft.com/en-us/windows/apps/design/accessibility/designing-inclusive-software) 2, 7; [basics-navigation-basics](https://learn.microsoft.com/en-us/windows/apps/design/basics/navigation-basics) 2, 4; [in-app-help-instructional-ui](https://learn.microsoft.com/en-us/windows/apps/design/in-app-help/instructional-ui) 7; [design-principles](https://learn.microsoft.com/en-us/windows/apps/design/design-principles) 1, 2
- **Enforce:**
  - Review.
- **kvit-ui now:** does not cover it. Not a library matter.
