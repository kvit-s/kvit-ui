---
name: kvit-ui
description: The shared interface vocabulary for the kvit desktop applications, written in Go on unison: how to describe a screen, which component to reach for, what each colour means and what it is reserved for, how to word labels, menus, dialogs and messages and when to use a symbol or a shortcut, and the workflow for getting a screen from a description to working Go. Use whenever building, changing or reviewing a screen in kvit-notes, kvit-notes-pro, kvit-hub or kvit-cash, or when writing anything in Go that draws.
---

# Building a screen in the kvit estate

The kvit desktop applications share one interface layer, `kvit-ui` (module
`github.com/kvit-s/kvit-ui`, imported as `kvitui`), built on the unison
toolkit. It has three parts: the design values (colours, type sizes, spacing
and density as named values), a vocabulary of components built from them,
and a drawing layer for proposing a screen before building it.

One `import kvitui "github.com/kvit-s/kvit-ui"` reaches every component, the
`UI` that holds the theme (`ui.Theme`), the interface size (`ui.Interface`)
and the fonts, and every colour and size by name (`kvitui.InkTextMuted`,
`kvitui.SizeSpace`).

`catalog.md` beside this file is the reference. It says how many components
there are and lists every one of them: what it is for, the fields and
constructor it has, and a working sample. It is written by
`kvit-ui-gallery --catalog`, from the same list the gallery draws, so it
describes what actually exists.

## Ask for a screen in the vocabulary, not in measurements

The single most useful thing to get right. A request written in what things
*mean* is buildable; a request written in values is not.

> a slim row with the project name, a health phrase, and the attention figure
> right aligned

is buildable. It names components (`SlimRow`, `Figure`) and meanings (a name,
a phrase, a figure), and it draws correctly in all four themes.

> a 30-pixel row with #e6b877 text on the right

is not. It names values instead of meanings: the hex is wrong in three of the
four themes, and 30 pixels stops being right the moment somebody changes the
interface size. A screen written that way passes review on the machine it was
drawn on and fails everywhere else.

This applies to describing an existing screen as much as to asking for a new
one. "The stalled projects are red" is a description that survives a theme
change; "the stalled projects are #c0392b" is not.

## The rules that must not be broken

Each of them is enforced by a test, and each exists because breaking it
produced a specific defect.

**No colour literal.** Every colour comes from the theme, as an `Ink`
(`kvitui.InkDanger`) read when the component draws. A literal is correct in
the theme it was picked in and wrong in the other three, and the
high-contrast pass is what catches it. If a screen needs a hue that is not a
design value, it needs one: say so rather than writing the value.
`rules_test.go` fails on a colour literal outside the design values.

**No numeric font size.** Every size is a type role: `RoleCaption`,
`RoleSmall`, `RoleBody`, `RoleStrong`, `RoleTitle`, `RoleHeadline`,
`RoleDisplay`, turned into pixels by `ui.Size(role)`. A numeric size does not
move when the reader changes the interface size, so the text stays put while
everything around it grows.

**No unnamed spacing or geometry.** The interface names the spacing scale
(`SizeSpaceTight`, `SizeSpaceSnug`, `SizeSpaceNear`, `SizeSpace`,
`SizeSpaceWide`, `SizeSpaceLoose`), the row heights, the control heights, the
radii and the reflow breakpoints, each a `Measure` read at layout time.
`kvitui.Px(n)` is for a measurement that has no name yet; a `Px` with a
literal in it that a second view will also need is a design value waiting to
be added. A size read once when a panel is built keeps the size it had then,
and `TestPagesFollowTheInterfaceSize` fails.

**Every distinction carries a second channel besides hue.** About one man in
twelve cannot separate red from green, and every screenshot is judged in
grayscale by somebody eventually. So a selected tab has an underline as well
as a colour, a health level has a shape as well as a hue, a bounded figure is
hatched as well as tinted, a link is underlined on hover, and a filled chip
is filled as well as coloured. When adding a state, say what its second
channel is.

**A value that was not measured is not zero.** `Figure` draws an em dash,
`Bar` draws a tick at the origin, `Spark` draws a baseline tick, a chart
takes `kvitui.NotMeasured` for a hole, and a chart with no data shows an
`EmptyState` rather than an axis drawn around zeros. A balance nobody has
computed and a balance of zero are different facts, and drawing them the
same way states something false.

**A count is handed over as a number and a noun, never as a finished
phrase.** The components that say how many of something there is —
`ViewHead`, `SectionHeading`, `Badge`, `SidebarItem`, `SearchField` — take the
number in `Count` (or `Matches`) and the word for one of the things in
`Counted` (or `MatchedNoun`), plus `CountedPlural` where English's suffixed s
is wrong. The component groups the digits the way the reader's locale groups
them (`ui.Number`) and picks the form, so one account reads "1 account" and
twelve hundred read "1,200 accounts".

**A control is a control, not a panel with a mouse handler.** Use `Button`,
`IconButton`, `Check`, `Switch` and the rest. A bare panel with a mouse
handler reaches no assistive technology: a screen reader is not told what it
is, the keyboard cannot get to it, and Space does nothing. Anything whose
whole label is a symbol also needs a `Label` in words, and that one string
becomes both the tooltip and the name a screen reader says. The gallery's
`TestEveryControlSaysWhatItIs` fails on a reachable control with no role or
no name.

**A symbol is asked for by meaning.** `kvitui.NewIcon(ui, "chevron-right")`,
never a literal `›`. The catalogue lists the names. An unknown name draws a
marked placeholder and logs a warning, which fails the gallery's screenshot
test.

## What each colour means

The names say what a colour is *for*. Using one for something else is how two
unrelated things end up looking related. Each is a field of `tokens.Tokens`
(`ui.Theme.Tokens().Border`) and an `Ink` (`kvitui.InkBorder`).

**Surfaces.** `WindowBackground` the window, `PanelBackground` a region with
its own ground, `ListBackground` a card or a list, `PopupBackground` anything
floating, `ChipBackground` a small mark, `FooterBackground` the status bar.

**Text.** `TextPrimary` what the reader is reading, `TextSecondary`
supporting text, `TextMuted` context, `TextFaint` a count or a timestamp,
`TextDisabled` a control that cannot be used.

**Lines.** `Border` is decorative: a rule between two panels, a separator. It
is deliberately below 3:1 in every theme. `BorderStrong` is the
control-boundary colour: the outline that says where a button, field,
checkbox or card-you-can-click is, held to 3:1 by the contrast floors in the
theme test. A control outlined in `Border` has invisible edges; a separator
drawn in `BorderStrong` looks clickable.

**Meaning.** `Accent` is progress and selection and never a target. `Link` is
a target. `Success` is finished and nothing else. `Warning` is something to
look at. `Danger` is stalled or destructive, and appears only where something
can be settled today. `FocusRing` is the keyboard, and nothing else uses it.

**Versions of the same text.** `ChangedTextBackground` marks text that
differs from another version without saying in which direction;
`AddedTextBackground` and `RemovedTextBackground` are the two directions, for
a diff. They are not `Success` and `Danger` — a removed line is not an error.
All three are close in luminance so text stays equally legible on any of
them, which means they separate by hue alone and are never the only mark: a
diff draws the `+` or `-` in the gutter and tints the line behind it.

**The portfolio vocabulary**, used by kvit-hub and available to anything
else: `AxisAttention` is the reader's own time and `AxisAgent` is agent time,
everywhere, and the two never share a scale or a total. Each has a more
legible `…Text` partner for small type. `ScopeDiscovered` is work that turned
up rather than being planned. `SignalHard`, `SignalSoft` and `SignalHygiene`
are the three severities.

**Charts.** Three ramps, per theme, all checked by `palette/palette_test.go`
against that theme's surfaces and against every hue that already means
something. `kvitui.InkCategorical(i)` for series identity, the sequential
ramp for magnitude, the diverging ramp for distance from a reference. Never a
meaning colour for a chart series: a bar in `Danger` says "stalled" to anyone
who has read the conventions, whatever the key claims.

## Words, symbols, commands and shortcuts

Interface text is not prose. The writing rules for documents and reports ask
for whole sentences, spelled-out meaning and no compressed labels; they do
not apply to a button, a menu item, a tab or a tooltip. Interface text has
its own conventions, shared by Apple and Microsoft: a command starts with its
verb, takes one to three words, drops articles, and is named the same
everywhere it appears. A sentence on a button is a defect, however clear the
sentence is.

The rules below hold on every system unless they say Windows or macOS; Linux
takes the Windows form. Each names its rules in
`docs/ui-guidelines/rules.md`, where the reasoning and sources are. Some lean
on library support that is still being built (`docs/ui-guidelines/library-gaps.md`);
where they do, the rule says what to do with the library as it is today.

### Labels

**1. A command is a verb that names what happens.** Buttons and menu items
start with a specific verb: Import, Archive, Reconcile, not Manage,
Configure or Do it. A statement or a description is not a command. (WRD-1)

- "The data is right" → "Mark correct"

**2. One to three words, within 26 characters.** Use only the words the
command needs. When a command seems to need a sentence, the sentence belongs
in its Explanation, its dialog, or nowhere. A toolbar label is preferably
one word. (WRD-2)

- "Take the new default" → "Reset to default…"
- "Save to the workspace" → "Save"
- "Move the selected column earlier" → "Move up"

**3. No articles, no "my" or "your".** Not on labels, menu items, tabs or
tooltips. (WRD-3)

- "Close the preview" on a back arrow → "Back"
- "Clear the filter" → "Clear"
- "Dismiss this notice" → "Dismiss"

**4. The label says exactly what happens.** Read the handler before naming
it. A label that names a different action from the one that runs is the
worst kind of wrong. (WRD-1, WRD-10)

- "Add the rows shown" that selects every row → "Select all"
- "Delete…" that moves a folder to the trash → "Move to Trash…" on macOS,
  "Delete…" on Windows and Linux with the dialog or undo saying it went to
  the Recycle Bin

**5. Ellipsis when the command asks for more.** A command that opens a
dialog, a chooser or a confirmation before acting ends with "…", the one
character U+2026, never three full stops. One that acts at once has none. On
Windows, a command whose whole job is to show a window (Settings,
Properties, About) has none. Write it into the label yourself for now.
(WRD-4)

- "Export the ledger" opening a dialog → "Export…"

**6. Write sentence case.** Write every label once, in sentence case:
"Select all", "Don't save". That is the Windows and Linux form for all text,
window and dialog titles included. kvit-ui will turn menu items, buttons,
tabs and column headings into title case on macOS; never write title case
yourself, and never mix the two in one menu. (WRD-6)

- "Quick Capture…" beside "Quick capture note…" in one menu → every item in
  sentence case

**7. Confirming buttons name the action; never Confirm.** The button that
carries out a dialog says what it does, so always set Dialog's ConfirmText:
its default, "Confirm", is a placeholder, not a label. A message that only
informs has one button: Close on Windows, OK on macOS. OK, Cancel and Apply
belong only to a properties dialog of one object (Windows). Yes and No
appear only together, on Windows, where no specific answer is possible. The
abandoning button is Cancel; once the work is done it becomes Close. Every
dialog has a safe button that leaves without acting. (WRD-7, WRD-8)

- "OK" on Insert link → "Insert"
- "OK" on Confirm import → "Import"
- "Delete them" → "Delete"

**8. Busy words say what is happening.** Not "Working…", "Loading…" or
"Please wait": an -ing phrase naming the operation. Always set Button's
BusyText. (WRD-25)

- "Working…" → "Importing statement…"

**9. A search field says what it searches.** Set SearchField's Placeholder;
its default, "Filter…", says nothing. (CTL-15)

- "Filter…" → "Search transactions"

**10. Tooltips are short and add something.** A symbol-only button's Label
is its tooltip and its screen-reader name: the command, verb first.
Explanation adds a sentence only when the name cannot say it: why the
control is disabled, what it opens. Never repeat a visible label, never
explain how to operate the control. (WRD-13, HLP-2)

- "Select every row the current filter shows" → "Select shown rows"

### Symbols or words

**11. A symbol for a familiar action, words for the rest.** Bold, italic,
search, delete, close, back, add: symbols, as IconButtons with a Label. An
action people would have to guess from a symbol: words. (CMD-22)

**12. Never letters or punctuation standing in for a symbol.** A text button
reading "B", "<>" or "+ Insert" is a symbol drawn in the wrong component.
Use an IconButton with the meaning name and a Label in words. If the symbol
you need is not in the catalogue (bold, italic, underline, strikethrough,
code and heading are not, yet), say so; do not type it. (CMD-22)

- "B", "I", "U", "S" → IconButtons named "Bold", "Italic", "Underline",
  "Strikethrough"
- "<>" → an IconButton named "Code"; "H" → one named "Heading"
- "+ Insert" → "Insert…" with the `plus` symbol before the words

**13. One symbol, one meaning, across the whole app.** Before using a
symbol, check what else in the app uses it. The close symbol closes; it does
not archive, withdraw, remove or cancel. (CON-2)

- X for archive → `archive`, "Archive"
- X for withdraw → `undo`, "Withdraw"
- X labelled "Close agent" that opens a browser → `external`, "Open in
  browser"

### Commands

**14. One command, one name, everywhere.** A command has the same words in
the menu bar, on its button, in its tooltip, in the tray menu and in any
dialog it opens. Until kvit-ui has a Command type, give each command one
name constant and one handler, and make every control that runs it from
those two. (CON-1)

- "Quick Capture…", "Quick capture note…" and "Quick capture" for one
  command → "Quick capture…" in every place

**15. Every command is on a visible command surface.** On macOS that is the
menu bar. On Windows and Linux it is the menu bar for an app with many
commands, or a toolbar and its More menu for one with few. Buttons and
context menus are extra routes to commands that also live there; a command
reached only from a context menu or a hover button cannot be found. (CMD-1)

**16. Standard commands in the standard menus.** Undo, Redo, Cut, Copy,
Paste, Select all and Find go in Edit; New, Open, Save, Close and Print in
File; show and hide commands in View. On Windows and Linux, Settings is the
last item of the sidebar or the toolbar's More menu, and About is a section
at the bottom of the Settings page; on macOS both are in the app menu. Exit
on Windows, Quit on macOS. (CMD-2 to CMD-8)

**17. Menus: most used first, grouped, short, destructive last.** Separate
groups with separators; at most one level of submenu, none in a context
menu on Windows; a destructive item goes last, after a separator. A context
menu holds only what applies to what was clicked. (CMD-10, CMD-12, CMD-13,
CMD-23)

**18. A toggle says what it does.** On Windows and Linux a menu item keeps
its name and shows a tick ("Sidebar" ✓); on macOS it swaps its words ("Show
sidebar" / "Hide sidebar"). Until MenuItem takes both words, use the tick
(`Checked`). (WRD-12)

### Shortcuts

**19. Standard commands get the standard shortcut.** If the app has Save,
Ctrl+S (Cmd-S) saves; Find, New, Close, Undo, Redo, Copy, Paste and Select
all likewise. Redo is both Ctrl+Y and Ctrl+Shift+Z on Windows and Linux.
Ctrl+F and Ctrl+E go to the window's search field. A shortcut whose command
cannot act does nothing, silently. Until kvit-ui has a shortcut registry,
handle the keys in `Window.OnKeyDown` and set the same key on the command's
MenuItem. (KEY-2, KEY-13)

- Save with no Ctrl+S → Save on Ctrl+S
- a window with search but no Ctrl+F → Ctrl+F focuses its SearchField

**20. Menus show their shortcuts.** A menu item whose command has a key
shows it in `Key`; that is how people learn it. Context menus on macOS show
none. (KEY-14, CMD-11)

- 23 menu items, none showing a shortcut → every item with a key shows it

**21. Few new shortcuts, and never the system's.** Give a new shortcut only
to a frequent command, on a letter that recalls its name, built on the
command modifier (`mod.OSMenuCommand()`, never `mod.Control` or
`mod.Command` by hand). Never Alt+letter on Windows (access keys own it),
never Ctrl+Alt, never Ctrl+Q, never a key the system or another command
already uses. (KEY-3, KEY-4, KEY-5)

**22. Access keys on Windows and Linux.** In a MenuItem's Text, put "&"
before the access letter ("&Import"); kvit-ui underlines it on Windows and
Linux and drops it on macOS. Not on OK, Cancel, Close, links or tabs; two
items at one level never share a letter. Text that comes from the reader's
data goes through `PlainMenuText`. Outside menus kvit-ui does not read "&"
yet, so do not write it there. (KEY-15)

### Dialogs and messages

**23. Let Dialog place its buttons.** Give Dialog its ConfirmText and
CancelText and let it lay them out; never build a row of dialog buttons
yourself. The platform order (the main action first and the safe button
last on Windows and Linux, the default last on macOS) is the Dialog's job.
(WIN-1)

**24. Ask only when it matters.** No dialog before an action that can be
undone: do it, and show ConfirmInPlace or a toast with Undo. Ask before a
rare action that loses data for good, and then say what will be lost.
(MSG-3)

**25. A message says what happened and what to do.** The title is the
problem or the question, specific, never "Error". The detail adds the cause
and the fix. An input error sits under its field and says how to put it
right. No "Oops", no blame, and no "we": write "Couldn't save the file",
which reads right on both systems. (MSG-4, MSG-5, MSG-6, WRD-14, WRD-28)

**26. Explanations stay short.** A setting gets a label and at most one
short line. A longer explanation goes behind a Hint. A paragraph that is
always on screen is read once and then ignored. (WRD-22)

- a permanent 50-word paragraph under one setting → one line, the rest in a
  Hint

**27. Settings apply at once.** An app's settings are a page (on macOS a
window) where each change takes effect as it is made: no OK, Cancel or
Apply. On Windows a setting that applies at once is a switch, not a
checkbox. Each setting is one row: its label, a one-line description where
needed, and the control. (WIN-16, WIN-23, CTL-2, WRD-22)

**28. Links go somewhere; buttons act.** A link opens a page, a window or a
help topic and commits nothing; anything that changes data is a button.
Labels end without a colon. (CTL-26, WRD-19)

## The workflow for a new screen

Four steps, in order. The first two are cheap and the last two are not, which
is the point: the drawing is where a screen gets rejected.

**1. Draw it as HTML.** A static page against `ux/tokens.css` and
`ux/frame.css` in this repository. `tokens.css` is written by
`tools/tokens-to-css` from the same design values the applications draw
with, so the drawing is in the colours the application will actually use.
`ux/PATTERN.md` is the contract for what a drawing may and may not do;
`ux/visual-language.md` is what the values mean. Nothing in a mockup loads
anything over the network.

**2. Render it and put the image in front of the owner.** `ux/render.sh`
writes a PNG at 1440×960. The `kvit-preview` skill wraps this. A screen is
approved by looking at it rather than by reading a description of it.

**3. Build it in Go from the components.** Assemble it; do not redraw it. If
a screen needs something the vocabulary does not have, that is a component
to add to `kvit-ui` rather than a shape to draw in the application — say
so rather than building it there.

**4. Finish with a screenshot run in four themes.** `./build.sh --shots` in
the application, or `./build.sh --shots` here for the components. Look at the
high-contrast pass in particular: it is where a distinction resting on hue
alone becomes obvious.

## Where things live

- `agent/skills/kvit-ui/catalog.md`: every component, written by
  `kvit-ui-gallery --catalog`.
- The repository root: the components, one file each, and `UI`.
- `tokens/`: `Theme`, `Interface` and `Typography`; `palette/`: the colour
  arithmetic the chart ramps are checked with.
- `ux/`: `tokens.css` (generated), `frame.css`, `render.sh`, `PATTERN.md`,
  `visual-language.md`.
- `cmd/kvit-ui-gallery/`: the gallery, which is also the documentation and
  the screenshot run.
- The tests beside each file, and `rules_test.go`: what holds all of the above
  to the rules on this page.
- `docs/ui-guidelines/`: the interface rules behind "Words, symbols, commands
  and shortcuts", merged from Apple's and Microsoft's guidelines, with the
  owner's decisions and a reviewer's checklist.

An application keeps its own mockups, its own data packs and its own
components built on top; only the stylesheets, the render script and the two
prose documents are shared.
