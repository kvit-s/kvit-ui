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
`docs/ui-guidelines/rules.md`, where the reasoning and sources are. The
library follows the system the app runs on (`UI.Form`: `FormWindows`, also
Linux's, or `FormMac`); a test draws either form on any system with
`uitest.Options{Form: kvitui.FormMac}`, and `Session.CheckLabels` and
`uitest.CheckApp` check a screen and an app against these rules in both.

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
dialog, a chooser or a confirmation before acting ends with "…". Set `Asks`
on the Command, Button, IconButton or MenuItem and kvit-ui adds it; never
type it, and never three full stops. One that acts at once has none. A
command whose whole job is to show a window (Settings) sets `OpensWindow`,
which adds "…" on macOS only. (WRD-4)

- "Export the ledger" opening a dialog → "Export" with `Asks`

**6. Write sentence case.** Write every label once, in sentence case:
"Select all", "Don't save". That is the Windows and Linux form for all text,
window and dialog titles included. kvit-ui turns menu items, menu titles,
buttons, tabs, segments and column headings into title case on macOS
(`TitleCase`); never write title case yourself, and never mix the two in one
menu. A name kept lower case on purpose goes in `kvitui.TitleCaseKept`. A
line naming the reader's own thing is `Verbatim`. (WRD-6)

- "Quick Capture…" beside "Quick capture note…" in one menu → every item in
  sentence case

**7. Confirming buttons name the action; never Confirm.** The button that
carries out a dialog says what it does, so always set Dialog's ConfirmText;
there is no default, and a dialog without one is warned about. A message
that only informs is `NewMessage`: one button, Close on Windows, OK on
macOS. OK, Cancel and Apply belong only to `DialogProperties`. Yes and No
appear only together, on Windows, where no specific answer is possible. The
abandoning button is Cancel; once the work is done, `Done()` turns it into
Close (Done on macOS). Every dialog has a safe button that leaves without
acting. (WRD-7, WRD-8)

- "OK" on Insert link → "Insert"
- "OK" on Confirm import → "Import"
- "Delete them" → "Delete"

**8. Busy words say what is happening.** Not "Working…", "Loading…" or
"Please wait": an -ing phrase naming the operation. Always set Button's
BusyText; a busy button without one keeps its words and is warned about.
(WRD-25)

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
Use an IconButton with the meaning name and a Label in words: `bold`,
`italic`, `underline`, `strikethrough`, `code` and `heading` are in the
catalogue. If the symbol you need is not, say so; do not type it. (CMD-22)

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
dialog it opens. Make each command once as a `kvitui.Command` (name, symbol,
shortcut, Enabled, Run) and build every place it appears from it:
`MenuItem{Command: c}`, `NewCommandButton`, `NewCommandIconButton`, a
`Toolbar`, a `MenuBar`. It is then disabled everywhere at once; call
`ui.Refresh()` when something other than a command changes what can run.
(CON-1, CMD-9)

- "Quick Capture…", "Quick capture note…" and "Quick capture" for one
  command → "Quick capture…" in every place

**15. Every command is on a visible command surface.** On macOS that is the
menu bar. On Windows and Linux it is the menu bar for an app with many
commands, or a `Toolbar` and its More menu for one with few. Buttons and
context menus are extra routes to commands that also live there; a command
reached only from a context menu or a hover button cannot be found.
`Window.OffSurface` lists the ones that are. (CMD-1)

**16. Standard commands in the standard menus.** Give a standard command its
`Std` (StdSave, StdUndo, StdFind…) and `MenuBar.Add` puts it in its menu at
its place for the form: Undo, Redo, Cut, Copy, Paste, Select all and Find in
Edit; New, Open, Save, Close and Print in File; the app's own lines go in a
standard menu with `AddTo` or in a menu of their own with `AddMenu`. On
Windows and Linux, Settings is the last item of the sidebar
(`Sidebar.AddFooter(NewSidebarCommand(ui, settings))`) or of Tools, and About
is `NewAboutSection` at the bottom of the Settings page; on macOS both are in
the app menu. Exit on Windows, Quit on macOS. (CMD-2 to CMD-8)

**17. Menus: most used first, grouped, short, destructive last.** Separate
groups with separators; at most one level of submenu, none in a context
menu on Windows; a destructive item goes last, after a separator. A context
menu holds only what applies to what was clicked. (CMD-10, CMD-12, CMD-13,
CMD-23)

**18. A toggle says what it does.** On Windows and Linux a menu item keeps
its name and shows a tick ("Sidebar" ✓); on macOS it swaps its words ("Show
sidebar" / "Hide sidebar"). Give the Command (or MenuItem) its `Checked`
and both words, `TurnOn` and `TurnOff`; kvit-ui shows each form's. One of a
set of options is `Radio`. A feature that is on or off and comes in styles,
such as bullets, is `NewToggleSplitButton`: its main part turns it on and
off and its menu picks the style. (WRD-12, CMD-15, CMD-17)

### Shortcuts

**19. Standard commands get the standard shortcut.** A Command with `Std`
has its platform's keys without choosing them: StdSave is Ctrl+S on
Windows and Linux and Cmd-S on macOS; Redo is both Ctrl+Y and Ctrl+Shift+Z.
A window answers the keys of the commands on its menu bar
(`Window.SetMenuBar`), on a toolbar it shows, and given to
`Window.AddCommands`, wherever the focus is, leaving typing keys to a
focused field. Ctrl+F and Ctrl+E go to the window's search field by
themselves. A shortcut whose command cannot act does nothing, silently.
(KEY-2, KEY-13)

- Save with no Ctrl+S → Save on Ctrl+S
- a window with search but no Ctrl+F → Ctrl+F focuses its SearchField

**20. Menus show their shortcuts.** A menu line made from a Command shows
its key, and a symbol button's tooltip adds it in parentheses on Windows and
Linux ("Bold (Ctrl+B)"); that is how people learn it. Context menus on macOS
show none, by themselves. (KEY-14, CMD-11)

- 23 menu items, none showing a shortcut → every item with a key shows it

**21. Few new shortcuts, and never the system's.** Give a new shortcut only
to a frequent command, on a letter that recalls its name, written once for
both systems: `kvitui.Keys(unison.KeyK, kvitui.ModPrimary|kvitui.ModShift)`,
never `mod.Control` or `mod.Command` by hand. The window refuses, with a
warning, Alt+letter on Windows (access keys own it), Ctrl+Alt, Ctrl+Q, a
key the system owns, and a key another command or a standard one already
has. (KEY-3, KEY-4, KEY-5)

**22. Access keys on Windows and Linux.** Every control with words has one:
kvit-ui gives each the first free letter of its words, unique in the
window, underlines it while Alt is held (always, where Windows'
"Underline access keys" setting is on), and Alt with the letter operates
the control. To choose the letter, put "&" before it ("&Import") in a
Button, Check, Switch, radio option, FormRow label or MenuItem; macOS draws
the words without it. OK, Cancel, Close, links, tabs and symbol buttons take
none. Text that comes from the reader's data goes through `PlainMenuText`.
(KEY-15)

### Dialogs and messages

**23. Let Dialog place its buttons.** Give Dialog its ConfirmText, and
OtherText for a second answer such as "Don't save", and let it lay them out;
never build a row of dialog buttons yourself. Dialog puts the main action
first and the safe button last on Windows and Linux, and the default last
on macOS; Return presses the default (the action, unless destructive) on
macOS, and on Windows only where `DefaultConfirm` says the action is the
likely, safe answer. Set `OnDiscard` where closing would throw away typed
work. (WIN-1, KEY-6, WIN-8)

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

**27. Settings apply at once.** An app's settings are written once, as
`SettingsPane`s given to `NewSettingsWindow`: a page in the app's window on
Windows and Linux (`Page`), a window of panes on macOS, and `Command` shows
the right one. Each change takes effect as it is made: no OK, Cancel or
Apply. On Windows a setting that applies at once is a switch, not a
checkbox. Each setting is a `SettingRow`: its label, a one-line description
where needed, and the control. A form's field is a `FormRow`, its label
above the field on Windows and beside it on macOS, with `Required` for an
asterisk. (WIN-16, WIN-23, CTL-2, WRD-22, LAY-6)

**28. Links go somewhere; buttons act.** A link opens a page, a window or a
help topic and commits nothing; anything that changes data is a button.
Labels end without a colon. Help for one place is a `HelpLink` worded as
the question its topic answers ("How does a budget roll over?"), never a
Help button or "Learn more"; a dialog takes it as `Help`, and kvit-ui
places it for the form. (CTL-26, WRD-19, HLP-6)

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
