# Plan: from the guidelines to the library and the apps

What comes after the guidelines in this folder, in the order to do it. Each step says what it builds, which rules it serves (rules.md), what is done when it ends, and roughly how big it is. Steps 1 to 5 are the foundation the rest stands on; take them in order. Steps 6 to 9 can overlap once step 3 is in.

The owner's decisions of 2026-10-23 (open-questions.md) are settled and the plan assumes them. library-gaps.md has the file:line detail behind each step.

## Done

- **The rules**: 199 merged rules from Apple's guidelines, the Windows 7 guide and Microsoft's current Fluent and WinUI guidance, with every open question decided.
- **The skill**: `agent/skills/kvit-ui/SKILL.md` has "Words, symbols, commands and shortcuts", 28 rules for agents building a screen, written for the library as it is today.

## What gets built, and what does not

kvit-ui adds a component when a rule cannot be kept without one or when a screen needs it, not because a guideline describes it. Most guideline controls already exist under a Kvit name (InfoBar is Notice, InfoBadge is Badge, Expander is Disclosure, NavigationView is Sidebar, SelectorBar is Segmented, NumberBox is NumberField, AutoSuggestBox is TypeAhead, BreadcrumbBar is Breadcrumb, Flyout is Popover), or are a field or two on one.

- **New, because the rules need them:** MenuBar, Toolbar with a More menu, split and menu buttons, the settings pieces (setting row, Settings at the foot of the sidebar, About section), form row, Tip, back history. All built on the Command type.
- **Extensions of existing components:** listed under steps 4, 5 and 9.
- **Only when an app asks:** date picker, colour chooser limited to the chart ramp, token field, document tabs, rating, annotated scroll bar, list-and-details helper.
- **Not built:** command links, a column browser, flip view, pips pager, pivot, semantic zoom, media playback, person picture, a Kvit-drawn title bar.

## 1. The platform form

`UI.Form`, `FormWindows` or `FormMac`, set from the system (Linux takes `FormWindows`), and settable in `uitest.Open` so any screen can be drawn headless in either form. Components read it at layout and draw time, as they read colours and sizes. The gallery's tests run every page in both forms.

- Rules: all 94 with a Windows and a macOS form.
- Done when: the gallery's screenshot and accessibility tests pass in both forms, and the screenshot set has a macOS-form pass.
- Size: small.

## 2. Commands and shortcuts

A `Command` (name, symbol, shortcut, enabled, run, optional explanation) from which MenuItem, Button, IconButton and toolbar items are made, so a command has one name, one symbol and one key everywhere and is disabled everywhere at once. The window keeps a shortcut registry that runs a command's key wherever the window has focus, refuses two commands on one key and keys the system owns, and stays silent when the command cannot act. A table of standard commands carries their keys per form: New, Open, Save, Close, Print, Undo, Redo (Ctrl+Y and Ctrl+Shift+Z), Cut, Copy, Paste, Delete, Select all, Find (Ctrl+F and Ctrl+E), Settings, Help, Full screen, Rename, Refresh, and the document-tab keys. Ctrl+Q is never bound.

- Rules: CON-1, CON-2, KEY-2, KEY-3, KEY-4, KEY-5, KEY-14, CMD-9.
- Done when: the gallery builds its own menus and shortcuts from commands, and an app that declares Save gets Ctrl+S (Cmd-S) without choosing it.
- Size: medium. It changes how every app wires its commands, so do it in kvit-ui first and move one app as the proof.

## 3. Command surfaces: menu bar and toolbar

`MenuBar` built from commands and the standard menus per form: the system's bar on macOS (App, File, Edit, View, Window, Help, with kvit-ui's own "Settings…"), a Kvit-drawn bar with access keys on Windows and Linux (File, Edit, View, Tools, Help). `Toolbar` with leading, centre and trailing groups and a More menu that takes what does not fit. Split button and menu button. MenuItem gains both toggle words, a bullet item, a bold default item, and context menus without shortcuts on macOS.

- Rules: CMD-1 to CMD-20, CMD-23, WRD-12, KEY-15 (menus).
- Done when: a gallery page shows a window with each surface in both forms, and a test fails a command that is on no visible surface.
- Size: large. Split it: the menu bar first, then the toolbar and buttons.

## 4. Dialog by platform

Button order by form (main action first and safe button last on Windows and Linux; default last on macOS); an optional default, Return pressing it on macOS and Enter pressing the focused button on Windows; no "Confirm" default; a Close-only form for messages (OK on macOS); a safe button in every dialog; no close symbol on macOS; Command-period on macOS; focus starting in the first field and returning to the opener; Tab kept inside; a warning before typed input is discarded; one modal at a time except an alert over a dialog; a don't-show-again option for informing dialogs. Button and IconButton press on Enter when focused on the Windows form.

- Rules: WIN-1, WIN-2, WIN-3, WIN-6, WIN-7, WIN-8, KEY-6, KEY-7, KEY-8, WRD-7, WRD-8.
- Done when: the dialog specimens pass a test of order, keys and focus in both forms, and WIN-1, WIN-8 and WRD-7 no longer read "contradicts" in rules.md.
- Size: medium.

## 5. Wording in the components

- A title-case converter for the macOS form, with an exception list, called by every component that shows a label.
- `Asks` on Button, IconButton and MenuItem, adding "…" per form.
- The library's own strings rewritten: no articles ("Clear", "Close", "Dismiss", "Clear sort", "Move up"), BusyText required with Busy, ShowPercent off by default, SearchField's placeholder set per use.
- The formatting symbols in the catalogue: bold, italic, underline, strikethrough, code, heading.
- The shortcut in a control's tooltip on the Windows form.

- Rules: WRD-3, WRD-4, WRD-6, WRD-13, WRD-25, CMD-22, CTL-15, FDB-10.
- Done when: the remaining wording contradictions in rules.md are gone and the gallery draws every label through the converter.
- Size: medium.

## 6. Tests

- `uitest.CheckLabels(window, form)`, reading one screen's controls: verb first, at most three words and 26 characters, no articles, no colon, no OK, Yes, No or Confirm outside the allowed cases, the form's case, no letters standing in for symbols, short tooltips, a titled dialog with a safe button, one primary button per region, no paragraph under a setting, a tooltip on cut-short text.
- `uitest.CheckApp(screens)`, reading every screen of an app together: one name per command, one meaning per symbol, one command per key, standard commands on their standard keys, every command on a visible surface, every shortcut shown in its menu.
- `uitest.CheckRules` extended to the Go source: "Oops", "tap", "click here", hand-typed "...", `mod.Control` or `mod.Command` in a key binding, "+ " at the start of a label.
- All of them run in both forms, in kvit-ui's gallery and in each app's test suite.

- Rules: the "test" entries in rules.md (66 one-screen, 18 whole-app, 24 source).
- Done when: the gallery passes them, and each app runs them in its tests, allowed to fail at first with a list of known failures that only shrinks.
- Size: medium. Writing them after steps 3 to 5 means they test the new behaviour instead of failing on what is about to change.

## 7. Keyboard

Access keys on every eligible control on Windows and Linux, unique in the window, underlined while Alt is held, with a uniqueness test. Arrow-key focus inside Tab, Segmented and the sidebar; radio groups that stop at their ends on Windows. Ctrl+Tab between tabs on every system. Type-ahead jumping in Table, Tree and lists. Escape closing TypeAhead's list and abandoning drags. Space paging a Region. Tab focus following the Mac's full-keyboard-access setting.

- Rules: KEY-1, KEY-7, KEY-9, KEY-10, KEY-11, KEY-15.
- Done when: every gallery page can be worked from the keyboard alone in both forms, checked by a test that tabs through each page.
- Size: medium.

## 8. Review the apps

With steps 1 to 6 in, review kvit-cash, kvit-notes and kvit-works screen by screen: the tests from step 6, and a reviewer agent applying `checklist.md` to each screenshot with its list of controls. Fix screens through the components, never by hand-drawing. Start with the survey's findings: labels over three words, the letter buttons in kvit-notes, the X symbol in kvit-works, OK on Insert link and Confirm import, the File menu's mixed case, the 50-word settings paragraph, the missing shortcuts.

- Needs: a session with the three app repositories attached (they were not reachable where these documents were written).
- Done when: each app's known-failure list from step 6 is empty.
- Size: large; one app per piece of work.

## 9. The rest of the library

In the order the apps will meet them:

- Settings: setting rows, Settings at the foot of the sidebar, an About section (WIN-16, WRD-22, CMD-3).
- Form row: a visible label above its field on Windows, tied by LabeledBy, with a required mark (WRD-19, LAY-6, CON-4).
- Field: a masked mode, a clear button on Windows, a length limit (CTL-7). Switch: On or Off beside it on Windows (CTL-6). Label: selectable content text, the full text as a tooltip when cut short (CTL-18, PTR-12).
- Window: AllowClose for unsaved work, Edited for the unsaved mark, its size and place restored and fitted to the screen, full screen (WIN-14, WIN-15, WIN-19, WIN-20).
- Tip for teaching tips (HLP-4); back history and a breadcrumb whose folded crumbs open a list (LAY-10).
- Progress with Cancel or Stop and time remaining (FDB-3, FDB-10); notifications that stay quiet while the app is in front and carry an urgency (MSG-10, MSG-12).
- An undo stack with named entries for the Edit menu, ConfirmInPlace and toasts (FDB-7, WRD-11).
- Help links to shipped HTML topics (HLP-6).

Each item is small to medium and gets its gallery page, catalogue entry and tests as CLAUDE.md requires.

## Alongside

- Use `checklist.md` with a reviewer agent on every new screen from now on, before the tests exist.
- Keep the skill in step: when a step lands, remove its "until kvit-ui has…" note from SKILL.md and give the rule its library form.
- Keep rules.md in step: when a step lands, update the "kvit-ui now" of the rules it serves, so the counts in the README show progress.
