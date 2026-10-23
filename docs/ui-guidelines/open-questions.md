# Open questions

Decisions the owner has to make before the rules can be built into kvit-ui, its tests and the skill. Each question names the rules it touches, the options, what each option changes, and a recommendation.

The first eight block the most later work: each one decides how a component behaves on every platform, so a component change made before it is decided may have to be undone. The rest can be decided as their rules come up.

## Decisions that block the most work

### Q-40. Linux

The owner applies the Windows form on Linux until Linux is decided. That choice matters for 70 rules, listed in the last section. In most of them GTK and Qt do what Windows does (access keys, Enter on the focused button, Ctrl+Tab, every control a Tab stop). They part ways in a handful of visible places:

- WIN-1: GNOME puts the confirming button last, as macOS does; KDE puts it first, as Windows does.
- WRD-6: GNOME and KDE use title case on buttons and menus in many apps; the Windows form is sentence case.
- WIN-16, CMD-1, CMD-3: GNOME apps use a header bar with a primary menu (Preferences, About) and live preferences, not a menu bar with Tools > Options and OK and Cancel.
- KEY-2, KEY-5: Linux desktops use Ctrl+Shift+Z for Redo and Ctrl+Q for Quit.
- MSG-17, CMD-24: GNOME shows tray icons only with an extension.

Options:

1. Windows form on Linux, as now. One code path for two systems; Kvit looks like a Windows app on a Linux desktop.
2. A third, Linux form for the five places above, following GNOME. More strings and layouts to test; Kvit looks at home on GNOME and less so on KDE.
3. Windows form, plus the two Linux key habits (Ctrl+Shift+Z as a second Redo, Ctrl+Q kept free for Quit).

Recommendation: 3. It costs two key bindings and settles the most noticed difference without a third form.

### Q-1. A menu bar on Windows and Linux

Rules: CMD-1, CMD-2, CMD-6, CMD-18. kvit-ui has no menu bar on any system (`menu.go`, survey). On macOS the system menu bar is required. The Windows 7 guide requires a menu bar for a main window with many commands, but lets a simple app hide it (Alt brings it back) or go without one when it would hold three menus or fewer.

Options:

1. Every main window has a visible Kvit-drawn menu bar on Windows and Linux. Every command is findable and has an access key; the window gains a row.
2. As 1, but a simple app may start with the bar hidden; Alt and View > Menu bar show it.
3. Also allow no menu bar when it would hold three menus or fewer, with toolbar menu buttons instead.

What it changes: whether kvit-ui needs a Kvit-drawn menu bar (options 1 and 2) or a toolbar that must reach every command (option 3); whether "every command is in the menu bar" (CMD-1) can be a test on Windows.

Recommendation: 2. One menu bar built from the app's commands serves keyboard and screen-reader users on every system, and small apps can still look uncluttered.

### Q-3. How a label gets its Windows and macOS capitalisation

Rules: WRD-6 and every rule that names a label. The Windows form is sentence case for everything but window titles (Don't save, Select all); the macOS form is title case for menu items, buttons, tabs and column headings (Don't Save, Select All). The library's strings are sentence case today.

Options:

1. **The app writes sentence case; the component converts to title case on macOS.** One string per label. Conversion follows Apple's rule: capitalise every word except articles, coordinating conjunctions and prepositions of three letters or fewer, and always the first and last word. Cost: a converter in kvit-ui and one call in every component that shows a label (Button, IconButton's tooltip, MenuItem, Tab, Segmented, column titles, Dialog's buttons). Failure cases: proper nouns and product names written lower case on purpose (iCloud, macOS) need an exception list; words the rule cannot classify (Set Up as a verb vs setup as a noun, "Move To…" where To is part of a phrasal verb) come out wrong; mixed strings with a user's data in them (Open "my notes") must not be converted, so data has to be passed apart from the label.
2. **The app supplies both forms.** Every label has two strings. Cost: twice the strings and translations; nothing checks that the two forms say the same thing, and they drift.
3. **The app writes title case; the component converts to sentence case on Windows and Linux.** Lowering is easier to get wrong: proper nouns and acronyms (Wi-Fi, PDF, Kvit Notes) must be kept, and the app would need to mark them.
4. **One form everywhere.** No platform behaviour; breaks the owner's decision.

Recommendation: 1, with a short exception list in kvit-ui and labels that never contain user data (WRD-29 already forbids building text from pieces). A test renders every screen in both forms and compares each name with the converter's output, so the exceptions show up.

### Q-7. The default button in Windows dialogs

Rules: KEY-6, WIN-3, MSG-3, KEY-15. Apple: a destructive dialog has no default, and a dialog that must be read has none. Windows 7: every dialog has a default, the safe answer in a risky question, and Enter presses whichever command button has focus. kvit-ui: Buttons answer to Space only (`button.go:75`); Dialog has no default and only focuses a button (`feedback.go:470-479`).

Options:

1. Platform forms. Windows and Linux: a default in every dialog, Cancel when the dialog is destructive, Enter presses the focused command button. macOS: Return presses the default; none in a destructive dialog.
2. Apple's rule everywhere: Return does nothing in a destructive dialog on every system.

Recommendation: 1. Dialog already focuses Cancel in a destructive dialog, so on Windows "Enter presses the focused button" gives the safe default with no new layout.

### Q-5. OK, Yes and No on Windows

Rules: WRD-1, WRD-7, WRD-8, WIN-2. Apple: never OK, Yes or No on a confirming button, except OK in an alert that only informs. Windows 7: OK commits a settings or properties window; Yes and No answer a yes-or-no question, always as a pair, when specific answers would be long or the reader must stop and read; a message with no action of its own has a single Close. Apple's own pages also disagree on the one-button alert: OK (`apple-alerts-19`) or Done (`apple-alerts-30`).

Options:

1. Platform forms as in WRD-7. The default is still a specific verb on both systems, so most dialogs read the same.
2. Specific verbs on every system; Yes, No and OK never.

For the macOS one-button alert: OK (alerts-19, the common Mac practice) or Done (alerts-30).

Recommendation: 1, and OK for the macOS informing alert. Either way Dialog loses its "Confirm" default (`feedback.go:424`).

### Q-2. Settings: where it lives and how it applies

Rules: WIN-16, WIN-23, CMD-3. Windows 7: Tools > Options opens a modal window that commits on OK, with Cancel. macOS: Settings… in the app menu opens a window whose changes apply at once. Current Windows apps mostly apply settings at once, but the only Windows source here is the Windows 7 guide.

Options:

1. Platform forms: modal Options with OK and Cancel on Windows and Linux, live Settings on macOS.
2. Settings apply at once everywhere, in a window with Close.

What it changes: whether a settings component needs OK, Cancel and Apply semantics (WIN-23) and a way to undo every change on Cancel.

Recommendation: 2, against the Windows 7 text, for one reason: a live settings window is how Windows 10 and 11 apps behave and how Kvit's own theme and size settings already work, and building transactional settings for one platform is a large cost for a dated convention. Revisit when the Fluent pages can be read (Q-39).

### Q-11. How far access keys go

Rules: KEY-15, KEY-3, KEY-6. Today only Kvit menus have access keys, and only inside an open menu (`menu.go:525-564`).

Options:

1. The full Windows form on Windows and Linux: every eligible control has an Alt+letter key, underlined when Alt is held or always if the system asks, unique in the window, checked by a test.
2. Menus only, plus Alt opening the menu bar.

Recommendation: 1. Access keys are on Microsoft's checklist several times, GTK and Qt have the same mnemonics, and they are the main keyboard route the apps lack. It needs a window-wide registry and a uniqueness test.

### Q-39. Microsoft's current guidance was not read

The Fluent and WinUI pages under learn.microsoft.com could not be fetched from this environment (the network policy refused the host). The Windows forms rest on the Windows 7 guide alone. Where current Windows practice is known to differ (live settings, "Learn more" links, select rather than click in help), the rule says so.

Options:

1. Treat the Windows 7 forms as final.
2. Fetch the Fluent pages (allow learn.microsoft.com in the environment's network settings, or supply them as files as was done for the Windows 7 guide) and rerun the merge for the rules marked "Fluent not consulted" before the library work starts.

Recommendation: 2, for the WRD, WIN and KEY rules only; the rest of the Windows 7 behaviour guidance is unlikely to change.

## Wording

### Q-4. Alert button capitalisation on macOS

Rules: WRD-6. Apple's Alerts page asks for sentence-case button titles (`apple-alerts-21`); its own examples and the Buttons page use title case (View All, Add to Cart; `apple-buttons-13`).

Options: title case (matches the Mac's own alerts and the Buttons page); sentence case (follows the Alerts page literally).

Recommendation: title case, as every Mac app's alerts are.

### Q-6. Done and Close at the end of flows, sheets and popovers

Rules: WRD-21, WRD-8, WIN-3, WIN-10. Windows 7 forbids Done (it is not an imperative) and ends a finished flow with Close; Apple ends flows and sheets with Done.

Options: platform forms (Done on macOS, Close on Windows and Linux); Done everywhere except Windows guided flows.

Recommendation: platform forms; it is one label table.

### Q-14. Delete or Move to Trash

Rules: WRD-10. kvit-notes calls a command "Delete…" that moves a folder to the trash. Apple's Finder calls it Move to Trash and gives no warning (`apple-feedback-06`); Windows calls the place the Recycle Bin, and the Windows 7 pages read here do not name the command.

Options: Move to Trash on macOS and Delete on Windows and Linux, with the dialog saying where it goes; Move to Trash everywhere; Delete everywhere.

Recommendation: the first. Windows Explorer's own command is Delete, and the confirmation (or the undo toast) says it went to the Recycle Bin.

### Q-15. "We", "you" and blame in messages

Rules: WRD-14, MSG-4. Apple: never use we (`apple-writing-14`) and keep we for the maker (`apple-inclusion-04`), on two pages. Windows 7: no you or your in error messages, so the reader is never the culprit.

Options: you in general, never the reader as the subject of a mistake, we never (one voice everywhere); Windows strictly (no you in any error message).

Recommendation: the first. It removes blame, which is what both rules are for.

### Q-22. A colon after a control's label

Rules: WRD-19, WRD-20. Windows ends a label beside or above its control with a colon (Name:). Apple asks for one on a slider's label and a combo box's label but shows none elsewhere.

Options: colon on the Windows form only; colon everywhere; no colons.

Recommendation: colon on the Windows form only. The form-row helper sets LabeledBy either way, so screen readers do not depend on it.

### Q-23. One tooltip text for both systems

Rules: WRD-13, HLP-1, KEY-14. Windows: the tooltip of an unlabelled control is its name plus the shortcut, Bold (Ctrl+B). Apple: a fragment saying what the control does.

Options: one string, the verb-led name, with the shortcut added on the Windows form; two strings per control.

Recommendation: one string. A verb-led name already says what the control does, and the shortcut in the tooltip is the Windows answer to the missing-shortcuts problem.

### Q-24. Toggle wording in menus

Rules: WRD-12, CMD-15, CTL-16. Windows keeps a toggle's name and ticks it (Status bar ✓); macOS swaps the words (Show Sidebar / Hide Sidebar).

Options: platform forms; ticks everywhere; swapped words everywhere.

Recommendation: platform forms, with MenuItem taking both words and choosing per system.

## Keyboard

### Q-16. What takes Tab focus on a Mac

Rules: KEY-10. On a Mac only text fields and lists take Tab focus unless the reader turned on full keyboard access; kvit-ui makes every control a Tab stop on every system (`control.go:85`).

Options: follow the Mac setting (read it through `platform`, as dark mode is); every control a Tab stop everywhere.

Recommendation: follow the Mac setting. Readers who use the keyboard on a Mac have turned it on, and the rest expect Tab to move between fields.

### Q-17. Read-only fields as Tab stops on Windows

Rules: KEY-10. Microsoft's checklist makes every read-only field a Tab stop; the text-box page makes a single-line one a stop only when people may scroll or copy it.

Recommendation: every read-only field, as the checklist says.

### Q-18. Ctrl+Tab between tabs on a Mac

Rules: KEY-9. Apple names no key for switching tabs; many Mac apps use Ctrl-Tab.

Recommendation: Ctrl+Tab and Ctrl+Shift+Tab on every system.

### Q-19. Contradictions in the Windows shortcut lists

Rules: KEY-2, KEY-5. The Windows 7 keyboard page lists Ctrl+G both as free and as Go to; it gives Redo as Ctrl+Y where Linux uses Ctrl+Shift+Z; Ctrl+Q is free on Windows and Quit on Linux.

Recommendation: Ctrl+G is Go to; Redo is both Ctrl+Y and Ctrl+Shift+Z on Windows and Linux; Ctrl+Q is never bound to an app command (so Linux can use it to quit).

### Q-20. The search key on Windows

Rules: KEY-13. The Windows 7 guide gives Ctrl+E for the window's search box and Ctrl+F for Find; most Windows and Linux apps now use Ctrl+F for both.

Recommendation: both keys go to the window's search field when there is one; Ctrl+F opens a find bar where there is a separate one. This answers kvit-works' missing Ctrl+F.

## Windows and dialogs

### Q-8. A close symbol on dialogs on macOS

Rules: WIN-6, WIN-3. KvitDialog draws a close symbol beside its title on every system (`feedback.go:428`). Windows dialogs always have one; Mac alerts and sheets never do.

Options: keep it everywhere (one look); drop it on macOS.

Recommendation: drop it on macOS. Which ways out a dialog has is behaviour, not appearance.

### Q-9. The command's name over a Windows message

Rules: MSG-4. A Windows message box has a title bar naming the command it came from; the main instruction is the question. KvitDialog is drawn in the window without a title bar, and its Title is the question.

Options: no caption on any system; a small caption naming the command on Windows and Linux.

Recommendation: no caption. The window's own title already names the app, and a second title line costs room for little.

### Q-10. Command links

Rules: CTL-25, WRD-2, WRD-14, WRD-29, CON-4. Windows 7 answers a question whose answers need many words with command links; macOS has no such control.

Options: a CommandLink component on Windows and Linux, radio buttons with a specific commit button on macOS, built from one list of answers; no command links anywhere (radio buttons and one verb button everywhere); command links everywhere.

Recommendation: no command links for now. No Kvit screen in the survey needs one, the wording rules apply either way, and the component can come when an app needs it.

### Q-12. Save automatically or ask on close

Rules: WIN-19, FDB-7. Apple: save automatically; no explicit save needed. Windows 7: closing with unsaved changes asks Save, Don't save, Cancel. kvit-cash has no Ctrl+S today.

Options: autosave on every system, Save kept as a command that saves at once; autosave on macOS and ask-on-close on Windows and Linux; ask on close everywhere.

Recommendation: autosave everywhere for notes-like data (kvit-notes, kvit-works); for kvit-cash, where a save commits a ledger, ask on close on every system and give Save its standard shortcut. This is an app decision as much as a library one; the library needs Window.AllowClose either way.

### Q-13. Marking unsaved changes on Windows and Linux

Rules: WIN-20, WIN-12. The Windows 7 guide sets no mark. Many Windows editors put an asterisk before the document name.

Options: asterisk in the window title (*Report - Kvit Cash); a mark beside the name in the header; no mark.

Recommendation: the asterisk, set through Window.Edited, which on macOS sets the system's dot.

### Q-25. Window floor and Windows' 800 × 600 minimum

Rules: LAY-3, WIN-11, A11Y-4. kvit-ui's window minimum is Px(880) by Px(600) (`tokens/interface.go:195`, `tokens/interface.go:203`), which does not fit an 800 × 600 screen less the taskbar, the Windows 7 floor.

Options: lower the floor to about Px(800) by Px(560) and let the body reflow below; keep it.

Recommendation: keep it. 800 × 600 screens are gone from the machines Kvit targets, and the floor already scales with the interface size. Record the decision so the rule's Windows form can say so.

### Q-28. A choice that reveals more controls

Rules: CTL-14, CTL-16. The Windows list, radio and check box pages say choosing a value only sets it; the progressive-disclosure page allows a choice to reveal follow-up controls.

Recommendation: allow it only for controls placed directly after the one that revealed them, announced to screen readers.

### Q-29. Pages that must not scroll at large interface sizes

Rules: CTL-10, LAY-9, A11Y-4. Windows: a dialog or tab page never scrolls. Kvit's interface size runs to 24 px, where some pages cannot fit.

Recommendation: a page fits without scrolling at the default size and scrolls rather than clips above it.

### Q-30. A disabled Next button

Rules: CTL-8. Windows keeps commit buttons enabled and reports what is missing; Apple suggests keeping Continue unavailable until required fields are filled.

Recommendation: the platform split, with the disabled button's Explanation saying why on macOS.

## Controls

### Q-26. Stepper arrows on Windows

Rules: CTL-9. Windows puts the spin arrows inside the number box; KvitStepper has minus and plus buttons beside the value.

Recommendation: keep Kvit's look and take the behaviour: a typed value, Up and Down keys, wrapping where sensible.

### Q-27. A click on a mixed checkbox on macOS

Rules: CTL-5. Windows cycles all on, all off, the original mix; Apple does not say; KvitCheck turns everything on and drops the mix (`check.go:49-51`).

Recommendation: cycle on every system; it loses nothing.

### Q-33. Bold default item in context menus

Rules: CMD-10, CMD-15. Windows marks a context menu's default item in bold. KvitMenu draws a checked line in bold as well as with a tick (`menu.go:236-239`).

Recommendation: bold for the default item on Windows and Linux; a checked line keeps its tick and loses the bold.

## Messages and feedback

### Q-21. Things that close on a timer

Rules: MSG-8, HLP-1, MSG-6. Apple: nothing the reader needs closes on a timer. KvitToast without an action closes after 4 s (`feedback.go:67-91`); Windows closes a tooltip after 5 s and an input balloon after about 10 s; Kvit's tooltip stays at least 3 s plus 60 ms a character (`control.go:228-232`).

Options: keep the 4 s toast for acknowledgements only (nothing needed later) and keep field messages until fixed; or make every toast stay until dismissed.

Recommendation: the first, with tooltip timing taken from the system where it has a setting (Windows) and the current formula elsewhere.

### Q-31. Advertising and errors in desktop notifications

Rules: MSG-10, MSG-12, MSG-16, MSG-17. Windows 7 never allows advertising in a notification; Apple allows marketing with consent. Windows allows an ignorable background failure as a notification; Apple keeps errors out of notifications.

Recommendation: no advertising on any system; errors by platform form.

### Q-32. A busy pointer

Rules: FDB-2. Windows uses the busy pointer for waits of one to five seconds; kvit-ui sets none, and Button.Busy changes the words instead.

Recommendation: a busy-pointer helper for blocking waits on Windows and Linux only, after checking unison's cursor set.

## Help, consistency and the rest

### Q-34. Help topics

Rules: HLP-6, HLP-8, CMD-8. There is no help in kvit-ui or the apps. Where topics live: HTML pages shipped with the app and opened in the browser at the topic; an in-app viewer; online pages. How links read: as the question the topic answers (Windows 7), or "Learn more" (current Windows Settings, unsourced). How help names controls: with the kind (the Back up daily check box, Windows) or without (macOS).

Recommendation: shipped HTML opened at the topic; question links; the kind named on the Windows form only.

### Q-35. Readers choosing the colours of their own data

Rules: CON-6. The Windows 7 colour page lets people pick colours for their own categories; Kvit's chart colours are a validated ramp.

Recommendation: readers may pick from the ramp's steps, named, for their own categories (labels, tags); nowhere else.

### Q-36. Printed output

Rules: CMD-25. No Kvit app prints today and kvit-ui has no printing.

Recommendation: leave printing out until an app needs it; then widen CMD-25 with what printed output must be.

### Q-37. Kvit's own appearance settings

Rules: CON-5, A11Y-4. Apple: no app-specific appearance setting (`apple-dark-mode-01`) and chrome may stay its size when text grows (`apple-typography-09`). Kvit has four themes with an override and scales the whole interface.

Recommendation: keep Kvit's settings (the owner's visual decisions), following the desktop by default as it does now.

### Q-38. When to pause on losing focus

Rules: FDB-9. Apple pauses things that need attention when the reader switches away, but keeps media playing. The pages do not draw the line.

Recommendation: pause only what scores or times the reader; everything else keeps running.

## Contradictions settled in the rules

No decision is needed unless the owner disagrees.

- Confirming a destructive menu item (`apple-pull-down-buttons-08`) against no confirmation for undoable actions (`apple-alerts-03`): MSG-3 asks only when the action cannot be undone.
- Hiding or disabling unavailable context-menu items (`apple-context-menus-07`): CMD-10 has a form per platform.
- App-menu order, most general first (`apple-the-menu-bar-46`), whose next sentence says the opposite: CMD-13 orders by use.
- "Avoid more than six tabs" followed by "six or more tabs" (`apple-tab-views-06`): CTL-10 allows six on macOS.
- A text Close button (`apple-sheets-01`) against the standard symbol (`apple-toolbars-13`): WRD-8 and WIN-6 give each place its form.
- Nothing essential at the window's bottom edge (`apple-layout-19`) against dialog buttons at the bottom: the rule is about main windows, and dialogs keep their buttons (LAY-5).
- A safe default first (`win7-exper-first-exper-20`): OTH-4 takes it on every system.
- Windows text-box prompts as labels against a visible label (WRD-19): a prompt may stand in only in a search field or a toolbar field.

## Where the Windows form matters on Linux

Seventy rules. For each, what Linux users see with the Windows form:

- A11Y-4: Linux desktops also scale at 125 and 150 percent, so the scaling test applies.
- CMD-1, CMD-2, CMD-3: a menu bar with File, Edit, View, Tools, Help, Options in Tools and About in Help, where GNOME uses a header-bar menu.
- CMD-10, CMD-11, CMD-12, CMD-15, CMD-16: Kvit-drawn context menus with a bold default, shortcuts shown, no submenus, bullets for exclusive options, no alternates.
- CMD-24, MSG-17: tray menus over StatusNotifierItem, which GNOME shows only with an extension; nothing may live only in the tray.
- CON-5: Linux has no hover-time setting, so tooltip timing keeps its defaults.
- CTL-2, CTL-5, CTL-8, CTL-9, CTL-10, CTL-11, CTL-15, CTL-16, CTL-25: up to seven radio buttons, mixed checkboxes cycling, commit buttons left enabled, stepper keys and wrap, up to seven tabs, no column browser, search at the top trailing corner, More options / Fewer options, command links.
- FDB-2, FDB-3: busy pointer from the cursor theme; Cancel or Stop on halt buttons.
- HLP-2, HLP-6, HLP-8: infotips; a question-worded help link where GNOME and KDE use a Help button; Windows help terms.
- KEY-2, KEY-3, KEY-4, KEY-5, KEY-6, KEY-9, KEY-10, KEY-11, KEY-12, KEY-13, KEY-15: Windows shortcuts (Q-19 for Redo and Quit), no Ctrl+Alt, Super reserved like the Windows key, Enter on the focused button, Ctrl+Tab, every control a Tab stop, Ctrl+wheel zoom (the middle button pastes on Linux, so no autoscroll there), no Alt on pointer gestures (window managers use Alt+drag), Ctrl+F and Ctrl+E (Q-20), Alt access keys.
- LAY-3, LAY-6, LAY-7: the 800 × 600 floor (Q-25), labels before fields, a dialog about an object opening beside it.
- MSG-1, MSG-6, MSG-10, MSG-11, MSG-13: attention through the urgency hint or xdg-activation, inline field messages, background failures as notifications, Windows title and body lengths, no notification buttons.
- PTR-1, PTR-15: tray hosts differ on what a left click does; Ctrl-click toggles a selection.
- WIN-1, WIN-6, WIN-7, WIN-11, WIN-12, WIN-15, WIN-16, WIN-19, WIN-23, WIN-25: confirming button first, a close symbol on every dialog titled by its command, helper dialogs over their owner, maximise and the size floor, "document - program" titles, F11 for full screen, Options with OK and Cancel, Save / Don't save / Cancel on close, OK/Cancel/Apply, placement best effort on Wayland.
- WRD-1, WRD-4, WRD-6, WRD-7, WRD-8, WRD-12, WRD-13, WRD-17, WRD-19, WRD-20, WRD-21, WRD-23, WRD-24: verbless standard labels, no ellipsis on Options, sentence case, OK and Yes/No, Cancel becoming Close or Stop, ticked toggles, shortcuts in tooltips, "click" for menu items, colons after labels, Close at the end of flows, "notification area" (which Linux calls the system tray), Ctrl+Shift+N key names.
