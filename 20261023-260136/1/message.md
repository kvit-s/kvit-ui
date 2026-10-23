# Task: condense Apple's and Microsoft's desktop UI guidelines into rules for the Kvit apps

## Background

Kvit is a set of cross-platform desktop applications written in Go: kvit-cash
(/home/sk/kvit-cash), kvit-notes (/home/sk/kvit-notes) and kvit-works
(/home/sk/kvit-works). All three draw with one component library, kvit-ui
(/home/sk/kvit-ui, module github.com/kvit-s/kvit-ui, built on the unison
toolkit). Read /home/sk/kvit-ui/CLAUDE.md first, then
agent/skills/kvit-ui/SKILL.md (what agents are told when they build a screen)
and agent/skills/kvit-ui/catalog.md (every component with its fields).

AI agents build most of the screens in these apps, and the results are often
clunky. A survey of the three apps on 2026-10-22 found five kinds of problem:

1. Sentences on buttons. About 16 of roughly 300 text-button labels are over
   three words, and many more are padded with articles; 13 of 75 icon-button
   tooltips are over three words.
  - kvit-cash app/configuration/editor.go:116,153,157: "Take the new default"
    (red, with a trash icon), "Use the shipped default", "Save to the
    workspace". The first two open the same dialog.
  - kvit-cash app/health/card.go:167 "The data is right";
    app/transactions/transactions.go:112 "Add the rows shown" (it means
    Select all); app/preview/preview.go:56 "Close the preview" on a back
    arrow.
2. Letters and typed characters instead of icons. kvit-notes
   app/toolbar.go:105-146 and editor/formatbar.go:59-71 have text buttons
   "B", "I", "U", "S", "<>", "H" and "+ Insert".
3. One command under several names, and one icon for several actions.
  - kvit-notes calls one command "Quick Capture…" (app/tray.go:52), "Quick
    capture note…" (app/menus.go:51) and "Quick capture"
    (app/capture.go:25).
  - In kvit-works the X icon means archive (ui/shell/overviewcard.go:61),
    withdraw (ui/conversation/view.go:1242), remove, cancel and close, and
    an X labelled "Close agent" opens a browser
    (ui/shell/detailheader.go:184).
4. Generic or mismatched labels. "OK" confirms Insert Link (kvit-notes
   app/links.go:260) and Confirm import (app/exchange.go:299). The menu item
   "Delete…" (app/actions.go:108) opens a dialog that moves the folder to the
   trash. The File menu (app/menus.go:39-56) mixes Title Case and sentence
   case.
5. Permanent explanations and missing shortcuts.
  - kvit-notes app/settings.go:269 shows a permanent 50-word paragraph under
    one setting.
  - kvit-cash has no Ctrl+S, and none of its 23 menu items shows a shortcut.
  - kvit-works has no Ctrl+F, Ctrl+N or Ctrl+W.
  - In kvit-notes, 4 of 109 menu items show a shortcut.

The kvit-ui skill teaches colour, sizes, screen-reader names and counts, and
says nothing about wording, icon versus text, menus or shortcuts. A likely
cause of problem 1: the writing instructions agents follow tell them to spell
things out and avoid compressed labels. That suits reports, but UI text has
its own conventions (verb first, one to three words, no articles).

This task produces a condensed set of rules, each traced to Apple's or
Microsoft's published guidelines. Later work turns them into:
(a) behaviour built into kvit-ui components;
(b) automated tests;
(c) a section of the kvit-ui skill that agents read while building a screen;
(d) a checklist that a reviewer agent applies to screenshots.
This task writes documents only. Do not change components, tests or the
skills.

## The owner's decisions: apply them, do not reopen them

- The apps behave as close to each platform's defaults as possible. Where
  Apple's and Microsoft's conventions differ, a rule has a Windows form and a
  macOS form, and the app uses the one for the system it runs on. Example:
  both Microsoft guides put the "do it" button leftmost and the safe button
  rightmost, while Apple puts the default button on the trailing side with
  Cancel before it. KvitDialog (catalog entry KvitDialog) currently puts the
  confirming button on the right on every system.
- From the sources, take behaviour, wording, keyboard, commands (toolbars,
  menus, buttons), windows and dialogs, messages, layout of content, feedback
  and consistency.
- Do not take visual appearance. Kvit has its own visual language: four
  themes, its own design values and the Phosphor icon font. So Liquid Glass,
  materials (Mica, Acrylic), SF Symbols specifics, the Aero icon style,
  corner radii and similar rules are left out. Count them, so it is visible
  how many were left out.
- Menus are the system's own on macOS and Kvit-drawn on Windows and Linux
  (already decided; see CLAUDE.md).
- Linux is undecided. Apply the Windows form on Linux, and list in the open
  questions every rule where that choice matters.

## Sources

Fetch every page; do not work from memory. Record the date you fetched each
page and, for Apple pages, the latest date in the page's change log.

Apple Human Interface Guidelines. The website draws its pages with
JavaScript, so a plain fetch returns only a title. The same content is
served as JSON at
  https://developer.apple.com/tutorials/data/design/human-interface-guidelines/<page>.json
for example .../toolbars.json.

- The text is in primaryContentSections[].content, as blocks (heading,
  paragraph, unorderedList, aside, table) made of inline runs. Titles of
  linked pages are in references.
- A rule is normally a paragraph whose first run is bold (type "strong").
- Rules for particular platforms sit under the heading "Platform
  considerations", in subsections named iOS, macOS, visionOS and so on.
  Read the general rules and the macOS subsection, and skip the other
  platforms.

Pages (75):
  design-principles designing-for-macos accessibility color dark-mode icons
  inclusion layout motion right-to-left typography writing charting-data
  drag-and-drop entering-data feedback file-management going-full-screen
  launching loading managing-notifications modality multitasking
  offering-help onboarding printing searching settings undo-and-redo charts
  text-views boxes collections column-views disclosure-controls labels
  lists-and-tables outline-views split-views tab-views buttons context-menus
  dock-menus edit-menus menus pop-up-buttons pull-down-buttons the-menu-bar
  toolbars path-controls search-fields sidebars token-fields alerts panels
  popovers scroll-views sheets windows color-wells combo-boxes image-wells
  pickers segmented-controls sliders steppers text-fields toggles gauges
  progress-indicators rating-indicators notifications focus-and-selection
  keyboards pointing-devices

Windows 7 UX guide (the Win32 guidelines). It is marked "created for Windows
7 and has not been updated". The look it describes is dated, but most of its
behaviour and wording guidance still applies, and it is the most detailed
desktop guidance either company has published.

- Table of contents: https://learn.microsoft.com/en-us/windows/win32/uxguide/toc.json
- Pages: https://learn.microsoft.com/en-us/windows/win32/uxguide/<name>
- The text is inside the page's <main> element.
- The page top-violations ("UX checklist", about 150 items) is Microsoft's
  own selection of the most important rules. Mark every rule that appears
  there.

Pages (46):
  how-to-design-desktop-ux top-violations ctrl-balloons ctrl-check-boxes
  ctrl-command-buttons ctrl-command-links ctrl-drop ctrl-group-boxes
  ctrl-links ctrl-list-boxes ctrl-list-views progress-bars
  ctrl-progressive-disclosure-controls ctrl-radio-buttons ctrl-search-boxes
  ctrl-sliders ctrl-spin-controls ctrl-status-bars ctrl-tabs ctrl-text-boxes
  ctrl-tooltips-and-infotips ctrl-tree-views cmd-menus cmd-toolbars text-ui
  text-style-tone mess-error mess-warn mess-confirm mess-notif inter-keyboard
  inter-mouse inter-accessibility win-window-mgt win-window-frames
  win-dialog-box win-wizards win-property-win vis-layout vis-fonts vis-color
  vis-std-icons exper-first-exper exper-printing winenv-help
  winenv-notification

Current Microsoft guidance (Fluent and WinUI). It calls itself "not a list of
prescriptive rules" and is mostly about appearance; take its behaviour rules.
Where it disagrees with the Windows 7 guide on behaviour, the current
guidance wins for Windows. Pages, under https://learn.microsoft.com/en-us/windows/apps/ :
  design/basics/  design/basics/app-silhouette  design/basics/navigation-basics
  design/basics/commanding-basics  design/basics/content-basics
  design/guidelines-overview  design/style/writing-style  design/iconography/
  design/layout/  design/usability/  design/accessibility/keyboard-accessibility
  develop/input/keyboard-interactions  develop/input/access-keys
  develop/input/keyboard-accelerators  develop/input/focus-navigation
  develop/input/guidelines-for-visualfeedback
  design/controls/dialogs-and-flyouts/dialogs  design/controls/command-bar
  design/controls/menus  design/controls/buttons  design/controls/tooltips
Also read every page under design/controls/ for a control that kvit-ui has
(compare with catalog.md), where the page states behaviour or wording rules.

The reading is about 240,000 words. Split it across parallel subagents, each
taking a group of pages and returning the extraction for those pages in the
format below. A workflow is fine.

## Steps

1. Extraction. Turn every rule on every page into one record (fields below),
   including the rules you leave out, marked with the reason.
2. Merging.
  - The same rule from different sources becomes one merged rule.
  - Where the platforms differ, the merged rule has a Windows form and a
    macOS form.
  - Where one source contradicts itself or leaves a case undefined, add an
    open question instead. Example: Apple's Alerts page says sentence case
    for all button titles, while its Buttons page asks for title case
    ("Add to Cart").
  - Aim for 100 to 200 merged rules.
3. Comparison with kvit-ui. For every merged rule, check what kvit-ui does
   now: it enforces the rule, contradicts it, or does not cover it. Give
   file:line. The components are in the repository root, one file each;
   start with control.go, button.go, iconbutton.go, the dialog files,
   menu.go, the tooltip files, field.go and window.go.
4. Enforcement. Say where each merged rule should be enforced:
  - in a component (name it, and say what changes);
  - by a test (say what the test reads: the Go source, as
    uitest.CheckRules does; the running screen's list of controls with
    their roles and names, as uitest's CheckNamed does; or every screen of
    an app together, which is the only way to catch problem 3 above);
  - by a reviewer's judgment.
5. Platform knowledge. List every rule that needs kvit-ui to know at run time
   which system it is on: button order, capitalisation, Ctrl or Cmd, Exit or
   Quit, where Settings and About go, the menu bar, and any others you find.
   For capitalisation, propose how a label gets its Windows and macOS forms
   (for example, the app writes one form and the component converts it, or
   the app supplies both), with the cost and the failure cases of each.
6. Examples. Look at the apps only enough to pick 20 to 30 before-and-after
   pairs for the skill draft. The full review of the apps is later work.

## Outputs, in /home/sk/kvit-ui/docs/ui-guidelines/

- README.md: what the folder is and why it exists; the sources with fetch
  dates; how the files relate; counts (pages read, rules extracted, rules
  left out by reason, merged rules).
- extracted.jsonl: one JSON object per extracted rule, with these fields:
  - id;
  - source (apple, windows7 or fluent), page, url (with the section anchor
    where the page has one), section;
  - quote: at most 25 words;
  - rule: one sentence in your own words;
  - platforms (macos, windows or all);
  - topic (commands, wording, keyboard, pointer, windows-and-dialogs,
    messages, controls, layout, feedback, accessibility, consistency, help
    or other);
  - kept (true or false), and the reason when false (visual,
    other-platform or not-applicable);
  - in_windows_checklist (true when the rule is on the top-violations page).
- rules.md: the merged rules, grouped by topic. Each rule has:
  - a stable id (a topic prefix and a number, such as CMD-12);
  - the rule in plain words, with an example where one helps;
  - the Windows form and the macOS form, when they differ;
  - its sources (extracted ids and links);
  - where it will be enforced;
  - what kvit-ui does now.
- open-questions.md: every decision the owner has to make. For each, give
  the options, what each one would change, and your recommendation.
- library-gaps.md: what kvit-ui contradicts or lacks, with file:line, the
  component change or test you propose, and an order to do them in.
- checklist.md: 30 to 50 judgment items for a reviewer agent that looks at a
  screenshot together with the screen reader's list of controls. Each item
  names its rule ids and says what a failure looks like on screen.
- skill-section-draft.md: a section for agent/skills/kvit-ui/SKILL.md on
  wording, icons versus text, commands, menus and shortcuts. It has 20 to 30
  rules with before-and-after pairs from the apps, and states that UI text
  follows UI conventions rather than the prose-writing rules.

Apple's and Microsoft's text is copyrighted. Write each rule in your own
words, keep quotes to 25 words or fewer, and keep the downloaded pages
outside the repository (for example in ~/.cache/kvit/ui-guidelines/).

Follow the writing rules in the global and project CLAUDE.md files in every
document.

## When you finish

Report:

- the counts;
- the ten rules the apps break most, judging from the survey above;
- the open questions that block the most later work;
- the library changes you would make first, in order.
