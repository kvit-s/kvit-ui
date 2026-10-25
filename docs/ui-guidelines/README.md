# UI guidelines for the Kvit apps

AI agents build most of the screens in kvit-cash, kvit-notes and kvit-works, and the results are often clunky: sentences on buttons, letters typed where symbols belong, one command under several names, generic labels such as OK, permanent paragraphs of explanation, missing shortcuts (survey of 2026-10-22). The kvit-ui skill teaches colour, sizes and screen-reader names, and says nothing about wording, symbols, menus or shortcuts.

This folder condenses Apple's and Microsoft's published desktop guidelines (Apple's Human Interface Guidelines, Microsoft's Windows 7 UX guide, and Microsoft's current Fluent and WinUI guidance) into rules for the Kvit apps, each traced to its source, and compares each rule with what kvit-ui does today. The rules are turned into behaviour in kvit-ui's components, automated tests (`uitest`), a section of the kvit-ui skill and a reviewer's checklist; guidelines-plan.md says what is built and what is left.

The owner's decisions applied throughout:

- The apps behave as close to each platform's defaults as possible. Where Apple and Microsoft differ, a rule has a Windows form and a macOS form, and the app uses the one for the system it runs on.
- Behaviour, wording, keyboard, commands, windows and dialogs, messages, layout of content, feedback and consistency are taken; visual appearance is not (Kvit has its own look). The rules left out for that reason are counted below.
- Menus are the system's own on macOS and Kvit-drawn on Windows and Linux.
- Linux takes the Windows form (decided 2026-10-23, Q-40); open-questions.md lists every rule where that matters.

## The files

| File | What it holds |
|---|---|
| `extracted.jsonl` | Every rule found on every page read, one JSON object per line, including those left out and why. Each record names the merged rules it went into (`merged_into`). |
| `rules.md` | The merged rules, grouped by topic, each with its Windows and macOS forms, its sources, where it should be enforced and what kvit-ui does now. |
| `open-questions.md` | Every question the rules raised, with options, what each changes, the recommendation and the owner's decision of 2026-10-23 (a table of decisions comes first), and the questions the current Microsoft guidance settled. |
| `library-gaps.md` | Where kvit-ui contradicts or lacks a rule, with file:line, the change or test proposed, and the order to do them in. Also what kvit-ui needs to know about the platform at run time, and how a label can get its two capitalisations. |
| `checklist.md` | 47 judgment items for a reviewer agent looking at a screenshot and the screen reader's list of controls. |
| `skill-section-draft.md` | The draft of the section now in `agent/skills/kvit-ui/SKILL.md` ("Words, symbols, commands and shortcuts"): 28 rules on wording, symbols, commands, menus and shortcuts, with before-and-after pairs from the apps and the library. |
| `guidelines-plan.md` | The next steps: what gets built in kvit-ui and in which order, the tests, and the review of the apps. |
| `app-review-brief.md` | The brief for reviewing each app (step 8): what has to happen first, giving the app's agents the kvit-ui skills, adding the checks to the app's tests, moving its commands, the survey's findings still to fix, and when an app is done. |

The extracted records are the evidence; the merged rules are what later work builds on; the open questions say what is not settled; library-gaps, the checklist and the skill draft are the three ways the rules reach code, review and agents.

## Sources

### Apple Human Interface Guidelines

75 pages, fetched on 2026-10-23 as JSON from `https://developer.apple.com/tutorials/data/design/human-interface-guidelines/<page>.json` (the website draws its pages with JavaScript). From each page the general rules and the macOS subsection were taken; rules under other platforms' subsections were recorded and left out. The second column is the latest date in the page's change log; 21 pages have no change log.

| Page | Latest change |
|---|---|
| [Accessibility](https://developer.apple.com/design/human-interface-guidelines/accessibility) | 2025-06-09 |
| [Alerts](https://developer.apple.com/design/human-interface-guidelines/alerts) | 2024-02-02 |
| [Boxes](https://developer.apple.com/design/human-interface-guidelines/boxes) | (no change log) |
| [Buttons](https://developer.apple.com/design/human-interface-guidelines/buttons) | 2025-12-16 |
| [Charting data](https://developer.apple.com/design/human-interface-guidelines/charting-data) | 2022-09-23 |
| [Charts](https://developer.apple.com/design/human-interface-guidelines/charts) | 2022-09-23 |
| [Collections](https://developer.apple.com/design/human-interface-guidelines/collections) | (no change log) |
| [Color](https://developer.apple.com/design/human-interface-guidelines/color) | 2025-12-16 |
| [Color wells](https://developer.apple.com/design/human-interface-guidelines/color-wells) | (no change log) |
| [Column views](https://developer.apple.com/design/human-interface-guidelines/column-views) | (no change log) |
| [Combo boxes](https://developer.apple.com/design/human-interface-guidelines/combo-boxes) | (no change log) |
| [Context menus](https://developer.apple.com/design/human-interface-guidelines/context-menus) | 2023-12-05 |
| [Dark Mode](https://developer.apple.com/design/human-interface-guidelines/dark-mode) | 2024-08-06 |
| [Design principles](https://developer.apple.com/design/human-interface-guidelines/design-principles) | 2026-06-08 |
| [Designing for macOS](https://developer.apple.com/design/human-interface-guidelines/designing-for-macos) | (no change log) |
| [Disclosure controls](https://developer.apple.com/design/human-interface-guidelines/disclosure-controls) | (no change log) |
| [Dock menus](https://developer.apple.com/design/human-interface-guidelines/dock-menus) | (no change log) |
| [Drag and drop](https://developer.apple.com/design/human-interface-guidelines/drag-and-drop) | 2023-10-24 |
| [Edit menus](https://developer.apple.com/design/human-interface-guidelines/edit-menus) | 2023-06-21 |
| [Entering data](https://developer.apple.com/design/human-interface-guidelines/entering-data) | 2023-06-21 |
| [Feedback](https://developer.apple.com/design/human-interface-guidelines/feedback) | (no change log) |
| [File management](https://developer.apple.com/design/human-interface-guidelines/file-management) | 2024-06-10 |
| [Focus and selection](https://developer.apple.com/design/human-interface-guidelines/focus-and-selection) | 2023-10-24 |
| [Gauges](https://developer.apple.com/design/human-interface-guidelines/gauges) | 2022-09-23 |
| [Going full screen](https://developer.apple.com/design/human-interface-guidelines/going-full-screen) | 2025-06-09 |
| [Icons](https://developer.apple.com/design/human-interface-guidelines/icons) | 2025-06-09 |
| [Image wells](https://developer.apple.com/design/human-interface-guidelines/image-wells) | (no change log) |
| [Inclusion](https://developer.apple.com/design/human-interface-guidelines/inclusion) | (no change log) |
| [Keyboards](https://developer.apple.com/design/human-interface-guidelines/keyboards) | 2025-06-09 |
| [Labels](https://developer.apple.com/design/human-interface-guidelines/labels) | 2023-06-05 |
| [Launching](https://developer.apple.com/design/human-interface-guidelines/launching) | 2024-06-10 |
| [Layout](https://developer.apple.com/design/human-interface-guidelines/layout) | 2026-09-09 |
| [Lists and tables](https://developer.apple.com/design/human-interface-guidelines/lists-and-tables) | 2023-06-21 |
| [Loading](https://developer.apple.com/design/human-interface-guidelines/loading) | 2025-06-09 |
| [Managing notifications](https://developer.apple.com/design/human-interface-guidelines/managing-notifications) | (no change log) |
| [Menus](https://developer.apple.com/design/human-interface-guidelines/menus) | 2026-06-08 |
| [Modality](https://developer.apple.com/design/human-interface-guidelines/modality) | 2023-12-05 |
| [Motion](https://developer.apple.com/design/human-interface-guidelines/motion) | 2025-09-09 |
| [Multitasking](https://developer.apple.com/design/human-interface-guidelines/multitasking) | 2025-06-09 |
| [Notifications](https://developer.apple.com/design/human-interface-guidelines/notifications) | 2023-10-24 |
| [Offering help](https://developer.apple.com/design/human-interface-guidelines/offering-help) | 2023-12-05 |
| [Onboarding](https://developer.apple.com/design/human-interface-guidelines/onboarding) | 2024-06-10 |
| [Outline views](https://developer.apple.com/design/human-interface-guidelines/outline-views) | (no change log) |
| [Panels](https://developer.apple.com/design/human-interface-guidelines/panels) | (no change log) |
| [Path controls](https://developer.apple.com/design/human-interface-guidelines/path-controls) | (no change log) |
| [Pickers](https://developer.apple.com/design/human-interface-guidelines/pickers) | 2023-06-05 |
| [Pointing devices](https://developer.apple.com/design/human-interface-guidelines/pointing-devices) | 2023-06-21 |
| [Pop-up buttons](https://developer.apple.com/design/human-interface-guidelines/pop-up-buttons) | 2023-10-24 |
| [Popovers](https://developer.apple.com/design/human-interface-guidelines/popovers) | (no change log) |
| [Printing](https://developer.apple.com/design/human-interface-guidelines/printing) | (no change log) |
| [Progress indicators](https://developer.apple.com/design/human-interface-guidelines/progress-indicators) | 2023-09-12 |
| [Pull-down buttons](https://developer.apple.com/design/human-interface-guidelines/pull-down-buttons) | 2022-09-14 |
| [Rating indicators](https://developer.apple.com/design/human-interface-guidelines/rating-indicators) | 2022-09-23 |
| [Right to left](https://developer.apple.com/design/human-interface-guidelines/right-to-left) | (no change log) |
| [Scroll views](https://developer.apple.com/design/human-interface-guidelines/scroll-views) | 2026-06-08 |
| [Search fields](https://developer.apple.com/design/human-interface-guidelines/search-fields) | 2026-06-08 |
| [Searching](https://developer.apple.com/design/human-interface-guidelines/searching) | 2026-06-08 |
| [Segmented controls](https://developer.apple.com/design/human-interface-guidelines/segmented-controls) | 2023-06-21 |
| [Settings](https://developer.apple.com/design/human-interface-guidelines/settings) | 2024-06-10 |
| [Sheets](https://developer.apple.com/design/human-interface-guidelines/sheets) | 2026-03-24 |
| [Sidebars](https://developer.apple.com/design/human-interface-guidelines/sidebars) | 2026-06-08 |
| [Sliders](https://developer.apple.com/design/human-interface-guidelines/sliders) | 2023-06-21 |
| [Split views](https://developer.apple.com/design/human-interface-guidelines/split-views) | 2025-06-09 |
| [Steppers](https://developer.apple.com/design/human-interface-guidelines/steppers) | (no change log) |
| [Tab views](https://developer.apple.com/design/human-interface-guidelines/tab-views) | 2023-06-05 |
| [Text fields](https://developer.apple.com/design/human-interface-guidelines/text-fields) | 2023-06-05 |
| [Text views](https://developer.apple.com/design/human-interface-guidelines/text-views) | 2023-06-05 |
| [The menu bar](https://developer.apple.com/design/human-interface-guidelines/the-menu-bar) | 2025-06-09 |
| [Toggles](https://developer.apple.com/design/human-interface-guidelines/toggles) | 2024-03-29 |
| [Token fields](https://developer.apple.com/design/human-interface-guidelines/token-fields) | (no change log) |
| [Toolbars](https://developer.apple.com/design/human-interface-guidelines/toolbars) | 2025-12-16 |
| [Typography](https://developer.apple.com/design/human-interface-guidelines/typography) | 2025-12-16 |
| [Undo and redo](https://developer.apple.com/design/human-interface-guidelines/undo-and-redo) | (no change log) |
| [Windows](https://developer.apple.com/design/human-interface-guidelines/windows) | 2025-06-09 |
| [Writing](https://developer.apple.com/design/human-interface-guidelines/writing) | 2025-12-16 |

### Microsoft Windows 7 UX guide

46 pages: the 46 named in the brief, among them the UX checklist (top-violations), Microsoft's own selection of its most important rules. The guide says it "was created for Windows 7 and has not been updated"; its behaviour and wording guidance is taken, its look is not.

learn.microsoft.com could not be reached from the environment this work ran in (its network policy refused the host). The owner supplied the guide on 2026-10-23 as two EPUB files, an offline copy of the MicrosoftDocs/win32 source (CC-BY-4.0), and the pages were read from those. The copy carries no per-page dates. Each record's url is the page's address on learn.microsoft.com with an anchor derived from its heading the way learn.microsoft.com derives them; the anchors were not checked against the live site.

Pages: how-to-design-desktop-ux, top-violations, ctrl-balloons, ctrl-check-boxes, ctrl-command-buttons, ctrl-command-links, ctrl-drop, ctrl-group-boxes, ctrl-links, ctrl-list-boxes, ctrl-list-views, progress-bars, ctrl-progressive-disclosure-controls, ctrl-radio-buttons, ctrl-search-boxes, ctrl-sliders, ctrl-spin-controls, ctrl-status-bars, ctrl-tabs, ctrl-text-boxes, ctrl-tooltips-and-infotips, ctrl-tree-views, cmd-menus, cmd-toolbars, text-ui, text-style-tone, mess-error, mess-warn, mess-confirm, mess-notif, inter-keyboard, inter-mouse, inter-accessibility, win-window-mgt, win-window-frames, win-dialog-box, win-wizards, win-property-win, vis-layout, vis-fonts, vis-color, vis-std-icons, exper-first-exper, exper-printing, winenv-help, winenv-notification. All under `https://learn.microsoft.com/en-us/windows/win32/uxguide/`.

### Microsoft's current guidance (Fluent and WinUI)

learn.microsoft.com could not be reached from this environment, so the owner supplied the current guidance on 2026-10-23 as two EPUB files, an offline copy of the MicrosoftDocs/windows-dev-docs source (CC-BY-4.0): "Windows App Design Guidance (Fluent, Windows 11)" and "WinUI Controls and Patterns Library (Fluent, Windows 11)". The copy carries no per-page dates. Each record's url is the page's address under `https://learn.microsoft.com/en-us/windows/apps/design/`, derived from the file name, with an anchor derived from its heading; neither was checked against the live site.

83 pages were read: the 31 design pages that state behaviour or wording (design principles, guidelines overview, commanding, navigation, content and app silhouette, title bar, layout and responsive design, usability, accessibility, app settings, bidirectional text, in-app help, writing style, and the colour, typography, motion, materials and iconography overviews, whose visual rules were recorded and left out), and the 52 WinUI control pages for the controls kvit-ui has or the rules name. Left unread: app-icon construction, icon fonts, widgets and downloads (appearance or not applicable), and the control pages for media, images, shapes, date and time pickers, flip view, pips pager, swipe, pull to refresh, inverted lists, semantic zoom, two-pane view, item templates and item repeaters.

Six pages the brief named were not in the copy: keyboard accessibility, keyboard interactions, access keys, keyboard accelerators, focus navigation, and guidelines for visual feedback. The keyboard rules that depend on them still rest on the Windows 7 guide (open-questions.md, Q-39).

Where the current guidance and the Windows 7 guide disagree on behaviour, the current guidance wins for the Windows form, as the brief says. It changed the Windows forms of about forty rules and settled eleven open questions; the biggest changes are listed in open-questions.md under "Settled by Microsoft's current guidance".

The downloaded pages and their text conversions are kept outside the repository, in `~/.cache/kvit/ui-guidelines/`.

## How the rules were made

1. **Extraction.** Every guideline on every page became one record: a paragraph led by a bold sentence, a guideline bullet, a guidance sentence, a table of standard items. Each record has a quote of at most 25 words, the rule in our own words, its platform, its topic, and whether it is kept. Six agents read the Apple pages and eight the Windows 7 pages, one group of pages each; a ninth read Microsoft's checklist and marked every record whose rule appears there (`in_windows_checklist`).
2. **Merging.** The Apple records were merged first into draft rules. Each Windows record was then mapped onto a draft rule, noted where Windows does it differently (the Windows form) or contradicts it, or marked new. Six agents, one per family of topics, folded the new Windows records into the rules under a cap that kept the total under 200, wrote the Windows forms and raised open questions. Last, the current Microsoft guidance was read the same way (seven agents) and each record mapped onto the merged rules as agreeing, adding, changing or removing part of a Windows form, answering an open question, or new; six agents, one per family, applied those changes, and one new rule (LAY-10, in-app navigation) came of it. Every kept record ends in at least one merged rule.
3. **Comparison with kvit-ui.** An agent surveyed the library (every component, default string, key binding and platform switch, with file:line); each merged rule's "kvit-ui now" was written against it and checked in the code.
4. **Enforcement, platform knowledge, examples:** library-gaps.md, open-questions.md and skill-section-draft.md.

The application repositories (kvit-cash, kvit-notes, kvit-works) were not available in this environment, so the app examples are the ones the survey of 2026-10-22 named, with its file:line; the rest come from kvit-ui's own strings and gallery.

## Counts

| | Apple | Windows 7 | Fluent and WinUI | Total |
|---|---|---|---|---|
| Pages read | 75 | 46 | 83 | 204 |
| Rules extracted | 1,267 | 2,479 | 1,045 | 4,791 |
| Kept | 881 | 2,336 | 854 | 4,071 |
| Left out: visual appearance | 116 | 108 | 110 | 334 |
| Left out: other platforms (iOS, watchOS, visionOS, tvOS; pen, touch, gamepad) | 245 | 6 | 35 | 286 |
| Left out: not applicable (Apple-only services; Windows-only technology such as UAC, the Ribbon, widgets; XAML API choices) | 25 | 29 | 46 | 100 |
| On Microsoft's Windows 7 checklist | | 458 (131 of them the checklist page's own items) | | |

- Merged rules: 199 (WRD 28, CMD 25, KEY 15, PTR 15, WIN 24, MSG 17, CTL 26, FDB 10, LAY 10, A11Y 8, CON 8, HLP 8, OTH 5). 94 have a separate Windows and macOS form.
- kvit-ui today: enforces 77, partly covers 93, contradicts none, does not cover 29 (as of 2026-10-24, after steps 1 to 7 and 9 of guidelines-plan.md; rules.md has the detail).
- Where each rule should be enforced (a rule can have several): a component 143, a test reading one screen's controls 66, a test reading the Go source 24, a test reading every screen of an app 18, a reviewer 123.
- Questions: 52 in all. 11 were settled by the current Microsoft guidance; the owner decided the other 41 on 2026-10-23.
