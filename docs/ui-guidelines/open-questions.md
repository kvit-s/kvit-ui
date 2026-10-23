# Open questions and the owner's decisions

The decisions the rules needed from the owner. Each question names the rules it touches, the options, what each option changes, the recommendation, and what was decided.

The questions under "Decisions that block the most work" were the ones the library work depended on most.

Microsoft's current guidance (Fluent and WinUI, read on 2026-10-23 from the copy the owner supplied) settled eleven questions and narrowed five; those are at the end, under "Settled by Microsoft's current guidance", with what was decided. Questions Q-41 to Q-52 are new: the current guidance raised them.

## The decisions

The owner decided every open question on 2026-10-23: Q-40, Q-3, Q-7, Q-5, Q-11 and Q-41 by option 1, Q-39 by keeping the Windows 7 forms, and the rest as recommended. rules.md, library-gaps.md, checklist.md and skill-section-draft.md follow these decisions. Each question below keeps its options for the record, with the decision under its heading.

| Question | Decision |
|---|---|
| Q-3. How a label gets its Windows and macOS capitalisation | Option 1: the app writes sentence case; kvit-ui converts to title case on macOS, with an exception list. |
| Q-4. Alert button capitalisation on macOS | Title case. |
| Q-5. OK, Yes and No on Windows | WRD-7's Windows form as written; a single OK on a macOS informing alert. |
| Q-6. Done and Close at the end of flows, sheets and popovers | Platform forms: Done on macOS, Close on Windows and Linux. |
| Q-7. Enter on a focused button | Option 1: Enter presses a focused button on the Windows form only. |
| Q-8. A close symbol on dialogs on macOS | No close symbol on macOS dialogs. |
| Q-10. Command links | No command links on any system. |
| Q-11. How far access keys go | Option 1: full access keys on Windows and Linux, checked by a test. |
| Q-12. Save automatically or ask on close | Autosave for kvit-notes and kvit-works; kvit-cash asks on close and has Save on its shortcut. |
| Q-13. Marking unsaved changes on Windows and Linux | An asterisk in the window title, through Window.Edited. |
| Q-14. Delete or Move to Trash | Move to Trash… on macOS; Delete… on Windows and Linux, the dialog or undo saying where it went. |
| Q-15. "We" in messages, and one string for two systems | Messages in a form that suits both systems; "we" on neither. |
| Q-16. What takes Tab focus on a Mac | Follow the Mac's full-keyboard-access setting. |
| Q-17. Read-only fields as Tab stops on Windows | Every read-only field is a Tab stop on Windows. |
| Q-18. Ctrl+Tab between tabs on a Mac | Ctrl+Tab and Ctrl+Shift+Tab on every system. |
| Q-19. Contradictions in the Windows shortcut lists | Ctrl+G is Go to; Redo on both Ctrl+Y and Ctrl+Shift+Z; Ctrl+Q never bound. |
| Q-20. The search key on Windows | Both Ctrl+F and Ctrl+E go to the search field. |
| Q-21. Things that close on a timer | The 4 s toast for acknowledgements only; tooltip timing from the system where it has a setting. |
| Q-23. One tooltip text for both systems | One string; the shortcut added on the Windows form. |
| Q-24. Toggle wording in menus | Platform forms; MenuItem takes both words. |
| Q-25. Window floor and Windows' 800 × 600 minimum | Keep Kvit's Px(880) by Px(600) floor. |
| Q-31. Advertising and errors in desktop notifications | No advertising anywhere; errors by platform form. |
| Q-33. Bold default item in context menus | Bold default item on Windows and Linux; a checked line keeps only its tick. |
| Q-34. Help topics | Shipped HTML pages opened at the topic; question links; the control's kind named on Windows only. |
| Q-35. Readers choosing the colours of their own data | From the ramp's named steps, for the reader's own categories only. |
| Q-36. Printed output | Left out until an app needs it. |
| Q-38. When to pause on losing focus | Pause only what scores or times the reader. |
| Q-39. Microsoft's keyboard pages were not read | Keep the Windows 7 forms for the keyboard rules. |
| Q-40. Linux | Option 1: Linux takes the Windows form, with no Linux-specific additions. |
| Q-41. A Kvit-drawn title bar on Windows | Option 1: keep the system frame; Back, search and tabs sit in the window's own header. |
| Q-42. A dialog without a title | A title on every dialog, on every system. |
| Q-43. The click that closes a popover | Absorbed, except on a control that opens another popover. |
| Q-44. Small Windows control conventions Kvit does not have | Switch state word and field clear button on the Windows form; slider label position not taken; accordions only by the reader's choice. |
| Q-45. Search: start of words or anywhere | Contains-anywhere for list filters; start-of-word matches first in app search. |
| Q-46. Split buttons | The last variant chosen, on Windows. |
| Q-47. Success notices or toasts | Toasts for what the reader just did; a success Notice only for background work that ended while they were elsewhere. |
| Q-48. Teaching tips | A Tip component; Spotlight for tours. |
| Q-49. Small keyboard differences in the current guidance | Tab keeps closing a menu; radio groups stop at their ends on Windows; Escape closes a tooltip only when no popup is open. |
| Q-50. Contrast and high contrast on Windows | 5:1 where it costs little; Kvit's own high-contrast table kept. |
| Q-51. Windows-only details that could be shared | All four become general rules. |
| Q-52. Tooltips the reader did not ask for | Gallery and tests only. |

## Decisions that block the most work

### Q-40. Linux

**Decided (2026-10-23):** Option 1: Linux takes the Windows form, with no Linux-specific additions.

The owner applies the Windows form on Linux until Linux is decided. That choice matters for the rules listed in the last section. In most of them GTK and Qt do what Windows does (access keys, Enter on the focused button, Ctrl+Tab, every control a Tab stop). They part ways in a handful of visible places:

- WIN-1: GNOME puts the confirming button last, as macOS does; KDE puts it first, as Windows does.
- WRD-6: GNOME and KDE use title case on buttons and menus in many apps; the Windows form is sentence case for everything.
- CMD-1, CMD-3: GNOME apps use a header bar with a primary menu (Preferences, About). The current Windows form (a toolbar with a More menu, Settings as a page applying at once, About inside Settings) is now close to GNOME's; only the names differ (Settings against Preferences).
- KEY-2, KEY-5: Linux desktops use Ctrl+Shift+Z for Redo and Ctrl+Q for Quit.
- MSG-17, CMD-24: GNOME shows tray icons only with an extension.

Options:

1. Windows form on Linux, as now. One code path for two systems; Kvit looks like a Windows app on a Linux desktop.
2. A third, Linux form for the five places above, following GNOME. More strings and layouts to test; Kvit looks at home on GNOME and less so on KDE.
3. Windows form, plus the two Linux key habits (Ctrl+Shift+Z as a second Redo, Ctrl+Q kept free for Quit).

Recommendation: 3. It costs two key bindings and settles the most noticed difference without a third form.

### Q-3. How a label gets its Windows and macOS capitalisation

**Decided (2026-10-23):** Option 1: the app writes sentence case; kvit-ui converts to title case on macOS, with an exception list.

Rules: WRD-6 and every rule that names a label. The Windows form is now sentence case for all text, window and dialog titles included (Don't save, Select all); the macOS form is title case for menu items, buttons, tabs and column headings (Don't Save, Select All). The library's strings are sentence case today.

Options:

1. **The app writes sentence case; the component converts to title case on macOS.** One string per label. Conversion follows Apple's rule: capitalise every word except articles, coordinating conjunctions and prepositions of three letters or fewer, and always the first and last word. Cost: a converter in kvit-ui and one call in every component that shows a label (Button, IconButton's tooltip, MenuItem, Tab, Segmented, column titles, Dialog's buttons). Failure cases: proper nouns and product names written lower case on purpose (iCloud, macOS) need an exception list; words the rule cannot classify (Set Up as a verb vs setup as a noun, "Move To…" where To is part of a phrasal verb) come out wrong; mixed strings with a user's data in them (Open "my notes") must not be converted, so data has to be passed apart from the label.
2. **The app supplies both forms.** Every label has two strings. Cost: twice the strings and translations; nothing checks that the two forms say the same thing, and they drift.
3. **The app writes title case; the component converts to sentence case on Windows and Linux.** Lowering is easier to get wrong: proper nouns and acronyms (Wi-Fi, PDF, Kvit Notes) must be kept, and the app would need to mark them.
4. **One form everywhere.** No platform behaviour; breaks the owner's decision.

Recommendation: 1, with a short exception list in kvit-ui and labels that never contain user data (WRD-29 already forbids building text from pieces). A test renders every screen in both forms and compares each name with the converter's output, so the exceptions show up. The current Windows guidance makes option 1 simpler: Windows takes the string as written, titles included.

### Q-7. Enter on a focused button

**Decided (2026-10-23):** Option 1: Enter presses a focused button on the Windows form only.

Rules: KEY-6, WIN-3, KEY-15. Narrowed by the current guidance: on both systems a dialog now needs no default button, and gets one only when one answer is clearly the likely and safe one; a destructive answer is never the default. What still differs: on Windows, Enter presses whichever button has focus; on macOS, Return presses only the default, whatever has focus. kvit-ui: Buttons answer to Space only (`button.go:75`); Dialog has no default and only focuses a button (`feedback.go:470-479`).

Options:

1. Platform forms: Button and IconButton also press on Enter when focused, on the Windows form only.
2. Enter presses a focused button on every system.

Recommendation: 1. It is one key binding by form, and it keeps the Mac's Return-means-default habit.

### Q-5. OK, Yes and No on Windows

**Decided (2026-10-23):** WRD-7's Windows form as written; a single OK on a macOS informing alert.

Rules: WRD-1, WRD-7, WRD-8. Narrowed by the current guidance: its Dialogs page asks for specific answers and gives Close for a message that only needs dismissing; its Buttons page still lists OK and Yes/No in the button order. WRD-7's Windows form now says: specific answers; Close on a message; OK, Cancel and Apply only on a properties dialog of one object; Yes and No only where no specific answer is possible, as a pair. Apple's own pages disagree on the one-button alert: OK (`apple-alerts-19`) or Done (`apple-alerts-30`).

Options for Windows: the form as now written; or no Yes and No at all (specific verbs always).

For the macOS one-button alert: OK (alerts-19, the common Mac practice) or Done (alerts-30).

Recommendation: the Windows form as written, and OK for the macOS informing alert. Either way Dialog loses its "Confirm" default (`feedback.go:424`).

### Q-11. How far access keys go

**Decided (2026-10-23):** Option 1: full access keys on Windows and Linux, checked by a test.

Rules: KEY-15, KEY-3, KEY-6. Today only Kvit menus have access keys, and only inside an open menu (`menu.go:525-564`). The current guidance's access-key and keyboard pages were not in the copy supplied (Q-39), so this rests on the Windows 7 guide.

Options:

1. The full Windows form on Windows and Linux: every eligible control has an Alt+letter key, underlined when Alt is held or always if the system asks, unique in the window, checked by a test.
2. Menus only, plus Alt opening the menu bar where an app has one.

Recommendation: 1, but read the current access-keys page first. Access keys are on Microsoft's checklist several times, GTK and Qt have the same mnemonics, and they are the main keyboard route the apps lack.

### Q-41. A Kvit-drawn title bar on Windows

**Decided (2026-10-23):** Option 1: keep the system frame; Back, search and tabs sit in the window's own header.

Rules: WIN-11, CTL-15, LAY-10. The current Windows form lets the title bar hold Back, the sidebar toggle, an app-wide search box, document tabs and the account picture, as long as its empty parts still drag the window. kvit-ui uses the system's frame. Putting anything in the title bar means drawing it: the drag area, double-click to maximise, the system window menu on the icon, snap layouts on the maximise button, and the caption buttons.

Options:

1. Keep the system frame; Back, search and tabs sit in the window's own header below it.
2. Draw the title bar on Windows (and Linux), with those controls in it.

What it changes: option 2 is a large library cost and a source of platform bugs (snap layouts, high contrast, accessibility of the caption buttons); option 1 gives up a few pixels and the Windows 11 look.

Recommendation: 1. The rule allows controls in the title bar; it does not require them, and the owner does not take visual appearance.

### Q-39. Microsoft's keyboard pages were not read

**Decided (2026-10-23):** Keep the Windows 7 forms for the keyboard rules.

The copy of the current guidance supplied on 2026-10-23 holds the design pages and the WinUI controls, but not six pages the brief named: keyboard accessibility, keyboard interactions, access keys, keyboard accelerators, focus navigation and guidelines for visual feedback. KEY-3, KEY-4, KEY-10, KEY-13, KEY-15 and the shortcut parts of WRD-13, WRD-24 and KEY-2 still rest on the Windows 7 guide.

Options: treat the Windows 7 forms as final for these; or supply those six pages and rerun the merge for the KEY rules.

Recommendation: supply them before the keyboard work (library-gaps.md step 6) starts; nothing earlier depends on them.

## Wording

### Q-4. Alert button capitalisation on macOS

**Decided (2026-10-23):** Title case.

Rules: WRD-6. Apple's Alerts page asks for sentence-case button titles (`apple-alerts-21`); its own examples and the Buttons page use title case (View All, Add to Cart; `apple-buttons-13`).

Options: title case (matches the Mac's own alerts and the Buttons page); sentence case (follows the Alerts page literally).

Recommendation: title case, as every Mac app's alerts are.

### Q-6. Done and Close at the end of flows, sheets and popovers

**Decided (2026-10-23):** Platform forms: Done on macOS, Close on Windows and Linux.

Rules: WRD-21, WRD-8, WIN-3, WIN-10. Windows 7 forbids Done (it is not an imperative) and ends a finished flow with Close; Apple ends flows and sheets with Done.

Options: platform forms (Done on macOS, Close on Windows and Linux); Done everywhere except Windows guided flows.

Recommendation: platform forms; it is one label table.

### Q-14. Delete or Move to Trash

**Decided (2026-10-23):** Move to Trash… on macOS; Delete… on Windows and Linux, the dialog or undo saying where it went.

Rules: WRD-10. kvit-notes calls a command "Delete…" that moves a folder to the trash. Apple's Finder calls it Move to Trash and gives no warning (`apple-feedback-06`); Windows calls the place the Recycle Bin, and the Windows 7 pages read here do not name the command.

Options: Move to Trash on macOS and Delete on Windows and Linux, with the dialog saying where it goes; Move to Trash everywhere; Delete everywhere.

Recommendation: the first. Windows Explorer's own command is Delete, and the confirmation (or the undo toast) says it went to the Recycle Bin.

### Q-15. "We" in messages, and one string for two systems

**Decided (2026-10-23):** Messages in a form that suits both systems; "we" on neither.

Rules: WRD-14, MSG-4. Settled in substance by the current guidance: on Windows the app may take responsibility in the active voice with "we" (We couldn't save your file); Apple says never "we" (`apple-writing-14`). WRD-14 now has those two forms. What is left: one source string ships on both systems, so a message written with "we" needs a second, macOS string.

Options:

1. Per-platform message strings where the Windows form uses "we".
2. Write problem messages in a form that suits both (Couldn't save the file), and use "we" on neither.

Recommendation: 2. It costs nothing, reads well on both systems, and still never blames the reader.
### Q-23. One tooltip text for both systems

**Decided (2026-10-23):** One string; the shortcut added on the Windows form.

Rules: WRD-13, HLP-1, KEY-14. Windows: the tooltip of an unlabelled control is its name plus the shortcut, Bold (Ctrl+B). Apple: a fragment saying what the control does.

Options: one string, the verb-led name, with the shortcut added on the Windows form; two strings per control.

Recommendation: one string. A verb-led name already says what the control does, and the shortcut in the tooltip is the Windows answer to the missing-shortcuts problem.

### Q-24. Toggle wording in menus

**Decided (2026-10-23):** Platform forms; MenuItem takes both words.

Rules: WRD-12, CMD-15, CTL-16. Windows keeps a toggle's name and ticks it (Status bar ✓); macOS swaps the words (Show Sidebar / Hide Sidebar).

Options: platform forms; ticks everywhere; swapped words everywhere.

Recommendation: platform forms, with MenuItem taking both words and choosing per system.

## Keyboard

### Q-16. What takes Tab focus on a Mac

**Decided (2026-10-23):** Follow the Mac's full-keyboard-access setting.

Rules: KEY-10. On a Mac only text fields and lists take Tab focus unless the reader turned on full keyboard access; kvit-ui makes every control a Tab stop on every system (`control.go:85`).

Options: follow the Mac setting (read it through `platform`, as dark mode is); every control a Tab stop everywhere.

Recommendation: follow the Mac setting. Readers who use the keyboard on a Mac have turned it on, and the rest expect Tab to move between fields.

### Q-17. Read-only fields as Tab stops on Windows

**Decided (2026-10-23):** Every read-only field is a Tab stop on Windows.

Rules: KEY-10. Microsoft's checklist makes every read-only field a Tab stop; the text-box page makes a single-line one a stop only when people may scroll or copy it.

Recommendation: every read-only field, as the checklist says.

### Q-18. Ctrl+Tab between tabs on a Mac

**Decided (2026-10-23):** Ctrl+Tab and Ctrl+Shift+Tab on every system.

Rules: KEY-9. Apple names no key for switching tabs; many Mac apps use Ctrl-Tab.

Recommendation: Ctrl+Tab and Ctrl+Shift+Tab on every system.

### Q-19. Contradictions in the Windows shortcut lists

**Decided (2026-10-23):** Ctrl+G is Go to; Redo on both Ctrl+Y and Ctrl+Shift+Z; Ctrl+Q never bound.

Rules: KEY-2, KEY-5. The Windows 7 keyboard page lists Ctrl+G both as free and as Go to; it gives Redo as Ctrl+Y where Linux uses Ctrl+Shift+Z; Ctrl+Q is free on Windows and Quit on Linux.

Recommendation: Ctrl+G is Go to; Redo is both Ctrl+Y and Ctrl+Shift+Z on Windows and Linux; Ctrl+Q is never bound to an app command (so Linux can use it to quit).

### Q-20. The search key on Windows

**Decided (2026-10-23):** Both Ctrl+F and Ctrl+E go to the search field.

Rules: KEY-13. The Windows 7 guide gives Ctrl+E for the window's search box and Ctrl+F for Find; most Windows and Linux apps now use Ctrl+F for both.

Recommendation: both keys go to the window's search field when there is one; Ctrl+F opens a find bar where there is a separate one. This answers kvit-works' missing Ctrl+F.

### Q-49. Small keyboard differences in the current guidance

**Decided (2026-10-23):** Tab keeps closing a menu; radio groups stop at their ends on Windows; Escape closes a tooltip only when no popup is open.

Rules: KEY-8, KEY-9, HLP-1.

- A Windows menu keeps focus until it is dismissed; KvitMenu closes on Tab (`menu.go:517-518`).
- A Windows radio group stops at its ends; KvitRadioGroup wraps (`radiogroup.go:102-129`).
- Escape (and Ctrl on Windows) dismisses a tooltip; if Escape closes the tooltip first, a dialog under it needs a second Escape.

Recommendation: keep Tab closing a menu (Escape and Tab both leaving is the common desktop behaviour); stop at the ends on the Windows form; let Escape close a tooltip only when no popup is open, so one Escape still closes a dialog.

## Windows and dialogs

### Q-8. A close symbol on dialogs on macOS

**Decided (2026-10-23):** No close symbol on macOS dialogs.

Rules: WIN-6, WIN-3. KvitDialog draws a close symbol beside its title on every system (`feedback.go:428`). Windows dialogs always have one; Mac alerts and sheets never do.

Options: keep it everywhere (one look); drop it on macOS.

Recommendation: drop it on macOS. Which ways out a dialog has is behaviour, not appearance.

### Q-10. Command links

**Decided (2026-10-23):** No command links on any system.

Rules: CTL-25, WRD-2, WRD-14, WRD-29, CON-4. The Windows 7 guide answers a question whose answers need many words with command links. The current guidance has no command link control, and its Dialogs and Buttons pages turn further choices into radio buttons or checkboxes committed by one button (WIN-2). macOS has none either.

Options: drop CTL-25's Windows form and use radio buttons with a specific commit button on every system; or build a CommandLink component for Windows and Linux anyway.

Recommendation: drop it. Neither platform's current guidance has the control, and no Kvit screen in the survey needs one.
### Q-12. Save automatically or ask on close

**Decided (2026-10-23):** Autosave for kvit-notes and kvit-works; kvit-cash asks on close and has Save on its shortcut.

Rules: WIN-19, FDB-7. Apple: save automatically; no explicit save needed. Windows 7: closing with unsaved changes asks Save, Don't save, Cancel. kvit-cash has no Ctrl+S today.

Options: autosave on every system, Save kept as a command that saves at once; autosave on macOS and ask-on-close on Windows and Linux; ask on close everywhere.

Recommendation: autosave everywhere for notes-like data (kvit-notes, kvit-works); for kvit-cash, where a save commits a ledger, ask on close on every system and give Save its standard shortcut. This is an app decision as much as a library one; the library needs Window.AllowClose either way.

### Q-13. Marking unsaved changes on Windows and Linux

**Decided (2026-10-23):** An asterisk in the window title, through Window.Edited.

Rules: WIN-20, WIN-12. The Windows 7 guide sets no mark. Many Windows editors put an asterisk before the document name.

Options: asterisk in the window title (*Report - Kvit Cash); a mark beside the name in the header; no mark.

Recommendation: the asterisk, set through Window.Edited, which on macOS sets the system's dot.

### Q-25. Window floor and Windows' 800 × 600 minimum

**Decided (2026-10-23):** Keep Kvit's Px(880) by Px(600) floor.

Rules: LAY-3, WIN-11, A11Y-4. kvit-ui's window minimum is Px(880) by Px(600) (`tokens/interface.go:195`, `tokens/interface.go:203`), which does not fit an 800 × 600 screen less the taskbar, the Windows 7 floor.

Options: lower the floor to about Px(800) by Px(560) and let the body reflow below; keep it.

Recommendation: keep it. 800 × 600 screens are gone from the machines Kvit targets, and the floor already scales with the interface size. Record the decision so the rule's Windows form can say so.

### Q-42. A dialog without a title

**Decided (2026-10-23):** A title on every dialog, on every system.

Rules: MSG-4, WIN-6. On Windows a simple message, error or question may have no title, its first line saying what happened or asking the question. KvitDialog always has a title row; an empty Title leaves the row empty (`feedback.go:449`).

Options: give Dialog a title-less form on the Windows form; or keep a title on every dialog on every system (allowed by both platforms).

Recommendation: keep a title everywhere. Both forms allow it, and one layout is cheaper.

### Q-43. The click that closes a popover

**Decided (2026-10-23):** Absorbed, except on a control that opens another popover.

Rules: WIN-10, PTR-13. The current guidance says the click that closes a popover goes no further; kvit-ui passes it on to what is under it (`layer.go:31-34`), and PTR-13 wants one click to close one popover and open the next.

Options: absorb the click (the Windows form); pass it on (today, and the macOS behaviour PTR-13 describes); absorb it except on a control that opens another popover.

Recommendation: the last. It keeps a stray click from pressing something by accident and keeps PTR-13's one-click switch.

## Controls

### Q-33. Bold default item in context menus

**Decided (2026-10-23):** Bold default item on Windows and Linux; a checked line keeps only its tick.

Rules: CMD-10, CMD-15. Windows marks a context menu's default item in bold. KvitMenu draws a checked line in bold as well as with a tick (`menu.go:236-239`).

Recommendation: bold for the default item on Windows and Linux; a checked line keeps its tick and loses the bold.

### Q-44. Small Windows control conventions Kvit does not have

**Decided (2026-10-23):** Switch state word and field clear button on the Windows form; slider label position not taken; accordions only by the reader's choice.

Rules: CTL-6, CTL-7, CTL-9, CTL-16. Each is a small behaviour the current Windows controls have and Kvit's do not:

- A switch shows On or Off (or one other word) beside it (CTL-6). KvitSwitch shows no state word.
- A focused single-line field holding text shows a clear button (CTL-7). Only KvitSearchField has one.
- A slider's value label sits below it, centred, with a value tip while dragging (CTL-9). KvitSlider draws the value to the right (`slider.go:20-21`).
- Expandable sections open independently; an accordion (opening one closes the others) only when the reader chooses it (CTL-16). KvitDisclosure's Group makes an accordion when the caller asks (`disclosure.go:21-29`).

Options per item: adopt on the Windows form; adopt everywhere; leave out as appearance.

Recommendation: adopt the switch's state word and the field's clear button on the Windows form (they are behaviour people rely on); leave the slider label's position as appearance; follow the current guidance on accordions everywhere.

### Q-45. Search: start of words or anywhere

**Decided (2026-10-23):** Contains-anywhere for list filters; start-of-word matches first in app search.

Rules: CTL-15. The Windows 7 guide matches search at the start of words; the current guidance's filtering example keeps items containing the typed text anywhere.

Recommendation: contains-anywhere for a filter over a list (what KvitSearchField does), start-of-words ranking first for a search across the app.

### Q-46. Split buttons

**Decided (2026-10-23):** The last variant chosen, on Windows.

Rules: CMD-17. The current guidance makes the chosen menu item the main part's action every time; the rule keeps the likeliest variant, or the last one chosen when people repeat it.

Recommendation: follow the current guidance on Windows (last chosen becomes the main action); there is no Mac counterpart to keep apart.

## Messages and feedback

### Q-21. Things that close on a timer

**Decided (2026-10-23):** The 4 s toast for acknowledgements only; tooltip timing from the system where it has a setting.

Rules: MSG-8, HLP-1, MSG-6. Apple: nothing the reader needs closes on a timer. KvitToast without an action closes after 4 s (`feedback.go:67-91`); Windows closes a tooltip after 5 s and an input balloon after about 10 s; Kvit's tooltip stays at least 3 s plus 60 ms a character (`control.go:228-232`).

Options: keep the 4 s toast for acknowledgements only (nothing needed later) and keep field messages until fixed; or make every toast stay until dismissed.

Recommendation: the first, with tooltip timing taken from the system where it has a setting (Windows) and the current formula elsewhere.

### Q-31. Advertising and errors in desktop notifications

**Decided (2026-10-23):** No advertising anywhere; errors by platform form.

Rules: MSG-10, MSG-12, MSG-16, MSG-17. Windows 7 never allows advertising in a notification; Apple allows marketing with consent. Windows allows an ignorable background failure as a notification; Apple keeps errors out of notifications.

Recommendation: no advertising on any system; errors by platform form.

### Q-47. Success notices or toasts

**Decided (2026-10-23):** Toasts for what the reader just did; a success Notice only for background work that ended while they were elsewhere.

Rules: MSG-1, MSG-16, MSG-7. The current guidance lets an in-window notice bar report success at the end of long or background work. KvitNotice's doc says a notice is a condition still true and a toast acknowledges what finished (`feedback.go:218-223`).

Options: allow a success Notice for the end of background work the reader may have missed; keep success to toasts.

Recommendation: keep toasts for success the reader just caused, and allow a success Notice only for background work that ended while the reader was elsewhere.

### Q-48. Teaching tips

**Decided (2026-10-23):** A Tip component; Spotlight for tours.

Rules: HLP-4. The current guidance's teaching tip points at a control, never blocks, and can be ignored. KvitSpotlight shades the window, takes every press and places its text below.

Options: add a Tip component (anchored, non-blocking, dismissable) and keep Spotlight for tours; or use Spotlight for both.

Recommendation: a Tip component. A blocking spotlight for a one-line tip is the interruption HLP-4 forbids.

## Help, consistency and the rest

### Q-34. Help topics

**Decided (2026-10-23):** Shipped HTML pages opened at the topic; question links; the control's kind named on Windows only.

Rules: HLP-6, HLP-8, CMD-8. There is no help in kvit-ui or the apps. Where topics live: HTML pages shipped with the app and opened in the browser at the topic; an in-app viewer; online pages. How links read: as the question the topic answers (Windows 7), or "Learn more" (current Windows Settings, unsourced). How help names controls: with the kind (the Back up daily check box, Windows) or without (macOS).

Recommendation: shipped HTML opened at the topic; question links; the kind named on the Windows form only.

### Q-35. Readers choosing the colours of their own data

**Decided (2026-10-23):** From the ramp's named steps, for the reader's own categories only.

Rules: CON-6. The Windows 7 colour page lets people pick colours for their own categories; Kvit's chart colours are a validated ramp.

Recommendation: readers may pick from the ramp's steps, named, for their own categories (labels, tags); nowhere else.

### Q-36. Printed output

**Decided (2026-10-23):** Left out until an app needs it.

Rules: CMD-25. No Kvit app prints today and kvit-ui has no printing.

Recommendation: leave printing out until an app needs it; then widen CMD-25 with what printed output must be.

### Q-38. When to pause on losing focus

**Decided (2026-10-23):** Pause only what scores or times the reader.

Rules: FDB-9. Apple pauses things that need attention when the reader switches away, but keeps media playing. The pages do not draw the line.

Recommendation: pause only what scores or times the reader; everything else keeps running.

### Q-50. Contrast and high contrast on Windows

**Decided (2026-10-23):** 5:1 where it costs little; Kvit's own high-contrast table kept.

Rules: A11Y-3. The current guidance asks for 5:1 text contrast (3:1 for large text) and the system's own high-contrast colours. Kvit's tables are held to 4.5:1, and its high-contrast theme is its own fixed table.

Options: raise every table to 5:1 (they are shared across systems) or keep 4.5:1; read the system's high-contrast colours on Windows or keep Kvit's table.

Recommendation: raise to 5:1 where the palette tests show it costs little, and keep Kvit's own high-contrast table: the owner does not take visual appearance, and the theme already follows the desktop's high-contrast switch.

### Q-51. Windows-only details that could be shared

**Decided (2026-10-23):** All four become general rules.

Rules: WRD-2, WRD-27, WRD-29, CON-4. Some current Windows conventions were put in the Windows form only, because Apple's forms were not to change, though nothing on the Mac side argues against them: button text within 26 characters; no italics (they read badly for people with dyslexia); a check box label at most two lines, a radio label one line, a slider end one word; required fields marked with an asterisk after the label.

Recommendation: make all four general rules. One form is cheaper to test, and none conflicts with Apple.

### Q-52. Tooltips the reader did not ask for

**Decided (2026-10-23):** Gallery and tests only.

Rules: HLP-2. The current guidance shows a tooltip only in answer to the reader. `UI.ShowTooltip` (`popover.go:16-37`) lets an app show one unasked; the gallery uses it to show a tooltip permanently (`cmd/kvit-ui-gallery/specimens_feedback.go:14`).

Recommendation: keep ShowTooltip for the gallery and tests, and say in its doc comment that apps must not call it.

## Settled by Microsoft's current guidance

The owner's rule is that where the current guidance and the Windows 7 guide disagree on behaviour, the current guidance wins for Windows. Applying it settled these. The rules already say what is written here.

- **Q-1. A menu bar on Windows and Linux.** A menu bar is for apps with many commands in several categories (editing apps, as Notepad keeps one); other main windows may use a toolbar with a More menu as their command surface. Either way every command is reachable from a visible command surface, never only from a context menu (CMD-1). The macOS menu bar stays required. For the Kvit apps that likely means a menu bar in kvit-notes and toolbars with More menus in kvit-cash and kvit-works, an app decision the library should support both ways.
- **Q-2. Settings.** On Windows and Linux, Settings (not Options) is a page in the app's window that applies each change at once, with no OK, Cancel or Apply; it is the last item of the navigation sidebar, or of the toolbar's More menu, or of Tools where there is a menu bar and no sidebar. About is a section at the bottom of the Settings page, not a Help menu item. OK, Cancel and Apply stay only for a properties dialog of one object (WIN-16, WIN-23, CMD-3, CMD-8, CMD-20). This is what Q-2 recommended.
- **Q-9. The command's name over a Windows message.** No caption: a dialog's title is its main instruction on Windows too (MSG-4, WIN-6).
- **Q-22. A colon after a label.** No colon after a control's label or a group label on any system (WRD-19, WRD-20, CTL-5).
- **Q-26. Stepper on Windows.** Stepping stops at the limits, the button that would pass a limit is disabled, and the buttons sit beside the value, on every system (CTL-9). KvitStepper already does this (`stepper.go:61-69`).
- **Q-27. A click on a mixed checkbox.** All on, then all off, on every system; no cycle back to the mix (CTL-5). KvitCheck already does this (`check.go:49-51`).
- **Q-28. A choice that reveals more controls.** A checkbox may reveal the fields that depend on it, placed directly after it; a list never shows or hides other controls (CTL-14).
- **Q-29. Pages that must not scroll.** A settings page is one column that scrolls when it does not fit; dialogs and flow pages fit without scrolling at the default interface size and scroll rather than clip above it (LAY-9).
- **Q-30. A disabled Next button.** One rule on every system: a submit button may stay disabled until every required field is filled, with the required fields marked and the reason shown; otherwise it stays enabled and reports what is missing beside the field (CTL-8).
- **Q-32. A busy pointer.** Kept for short waits on Windows and Linux; the current guidance adds an indeterminate ring where the result will appear for blocking work of unknown length (FDB-2).
- **Q-37. Kvit's own appearance settings.** The current guidance has apps follow the system's light or dark theme by default and offer their own theme setting (Light, Dark, Use system setting), which is what Kvit does (CON-5).

Narrowed, still open above: Q-3 (Windows is now sentence case everywhere, which simplifies the conversion), Q-5, Q-7, Q-10, Q-15 and Q-39.

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

For each rule with a Windows form, what Linux users see with it. The current Windows guidance moved several rules closer to GNOME (live settings, a toolbar instead of a menu bar, no colons, placeholders kept on focus), so the visible differences are fewer than before.

- A11Y-3, A11Y-4: Windows contrast and scaling tests; Linux desktops also scale at 125 and 150 percent.
- CMD-1, CMD-2, CMD-3, CMD-8, CMD-20: a toolbar with a More menu or a menu bar with File, Edit, View, Tools, Help; Settings as the last item and About inside Settings, where GNOME has Preferences and About in its primary menu.
- CMD-10, CMD-11, CMD-12, CMD-15, CMD-16, CMD-17, CMD-22: Kvit-drawn context menus with common commands first, shortcuts shown, inapplicable items left out, bullets for exclusive options, no modifier alternates; split buttons; toolbar words beside symbols on wide windows.
- CMD-24, MSG-17: tray menus over StatusNotifierItem, which GNOME shows only with an extension; nothing may live only in the tray.
- CON-4: required fields marked with an asterisk.
- CON-5: Linux has no hover-time setting, so tooltip timing keeps its defaults.
- CTL-2, CTL-10, CTL-11, CTL-15: a switch for any setting that applies at once; a sidebar beyond a few tabs; no column browser, drill-down lists where a tree is not needed; search centred in the window's top area for the app, at a view's top trailing corner for the view.
- FDB-2, FDB-3: busy pointer from the cursor theme; Cancel or Stop on halt buttons.
- HLP-1, HLP-2, HLP-4, HLP-6, HLP-8: Ctrl dismissing a tooltip; infotips; teaching-tip placement; question-worded help links, where GNOME and KDE use a Help button; Windows help terms.
- KEY-2, KEY-3, KEY-4, KEY-5, KEY-6, KEY-9, KEY-10, KEY-11, KEY-12, KEY-13, KEY-14, KEY-15: Windows shortcuts (Q-19 for Redo and Quit), no Ctrl+Alt, Super reserved like the Windows key, document-tab keys, Enter on the focused button (Q-7), Ctrl+Tab, every control a Tab stop, Ctrl+wheel zoom (the middle button pastes on Linux, so no autoscroll there), no Alt on pointer gestures (window managers use Alt+drag), Ctrl+F and Ctrl+E (Q-20), shortcuts in tooltips, Alt access keys.
- LAY-3, LAY-5, LAY-6, LAY-7: the 800 × 600 floor (Q-25); a lone page button at the leading edge; labels above fields; a dialog about an object opening beside it.
- MSG-1, MSG-4, MSG-10, MSG-11, MSG-13: attention through the urgency hint or xdg-activation; a simple message may have no title; background failures as notifications; Windows title and body lengths; no notification buttons.
- PTR-1, PTR-15: tray hosts differ on what a left click does; Ctrl-click toggles a selection.
- WIN-1, WIN-2, WIN-6, WIN-7, WIN-10, WIN-11, WIN-12, WIN-15, WIN-19, WIN-23, WIN-25: the main action first and the safe button last, where GNOME puts the main action last; body text in every dialog; a close symbol on every dialog; helper dialogs over their owner; a popover with commit and Cancel discarding on an outside click; title-bar controls (Q-41); window titles; F11 for full screen; Save, Don't save, Cancel on close; OK, Cancel and Apply only in properties dialogs; placement best effort on Wayland.
- WRD-1, WRD-2, WRD-4, WRD-6, WRD-7, WRD-8, WRD-12, WRD-13, WRD-14, WRD-17, WRD-20, WRD-21, WRD-22, WRD-23, WRD-24, WRD-27, WRD-29: verbless standard labels, one-word toolbar labels, no ellipsis on Settings, sentence case for everything, specific answers with OK only in properties dialogs, Cancel becoming Close or Stop, ticked toggles, shortcuts in tooltips, "we" in problem messages, "select" as the verb, no colons, a labelled Back beside Next, setting rows, "notification area" (which Linux calls the system tray), Ctrl+Shift+N key names, no italics, short option labels.
