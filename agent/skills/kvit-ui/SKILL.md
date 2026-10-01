---
name: kvit-ui
description: The shared interface vocabulary for the kvit desktop applications, written in Go on unison: how to describe a screen, which component to reach for, what each colour means and what it is reserved for, and the workflow for getting a screen from a description to working Go. Use whenever building, changing or reviewing a screen in kvit-notes, kvit-notes-pro, kvit-hub or kvit-cash, or when writing anything in Go that draws.
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

An application keeps its own mockups, its own data packs and its own
components built on top; only the stylesheets, the render script and the two
prose documents are shared.
