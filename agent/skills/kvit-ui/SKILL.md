---
name: kvit-ui
description: The shared interface vocabulary for the kvit desktop applications: how to describe a screen, which component to reach for, what each colour means and what it is reserved for, and the workflow for getting a screen from a description to working QML. Use whenever building, changing or reviewing a screen in kvit-notes, kvit-notes-pro, kvit-hub or kvit-cash, or when writing anything in QML that draws.
---

# Building a screen in the kvit estate

Four Qt desktop applications share one interface layer, `kvit-ui`, which each
of them consumes as a pinned git submodule. It has three parts: a token layer
(colours, type sizes, spacing and density as named values), a vocabulary of
components built from those tokens, and a drawing layer for proposing a screen
before building it.

`import Kvit.Ui` reaches all of it: every component, the token singletons
`Theme`, `Interface` and `Typography`, and the icon font.

`catalog.md` beside this file is the reference. It says how many components
there are and lists every one of them: what it is for, its properties and a
working sample. It is generated from the repository on each build, so it
describes what actually exists.

## Ask for a screen in the vocabulary, not in measurements

The single most useful thing to get right. A request written in what things
*mean* is buildable; a request written in values is not.

> a slim row with the project name, a health phrase, and the attention figure
> right aligned

is buildable. It names components (`KvitSlimRow`, `KvitFigure`) and meanings
(a name, a phrase, a figure), and it renders correctly in all four themes.

> a 30-pixel row with #e6b877 text on the right

is not. It names values instead of meanings: the hex is wrong in three of the
four themes, and 30 pixels stops being right the moment somebody changes the
interface size. A screen written that way passes review on the machine it was
drawn on and fails everywhere else.

This applies to describing an existing screen as much as to asking for a new
one. "The stalled projects are red" is a description that survives a theme
change; "the stalled projects are #c0392b" is not.

## The rules that must not be broken

These are not style preferences. Each of them is enforced by a test or a lint
gate, and each exists because breaking it produced a specific defect.

**No colour literal in QML.** Every colour comes from `Theme`. A literal is
correct in the theme it was picked in and wrong in the other three, and the
high-contrast pass is what catches it. If a screen needs a hue that is not a
token, it needs a token: say so rather than writing the value.

**No numeric font size.** Every size comes from `Interface`: `caption`,
`small`, `body`, `strong`, `title`, `headline`, `display`. A numeric size does
not move when the reader changes the interface size, so the text stays put
while everything around it grows.

**No unnamed spacing or geometry.** `Interface` names the spacing scale
(`spaceTight`, `spaceSnug`, `spaceNear`, `space`, `spaceWide`, `spaceLoose`),
the row heights, the control heights, the radii and the reflow breakpoints.
`Interface.px(n)` is for a measurement that has no name yet; a `px()` with a
literal in it that a second view will also need is a token waiting to be
added.

**Every distinction carries a second channel besides hue.** About one man in
twelve cannot separate red from green, and every screenshot is judged in
grayscale by somebody eventually. So a selected tab has an underline as well
as a colour, a health level has a shape as well as a hue, a bounded figure is
hatched as well as tinted, a link is underlined on hover, and a filled chip is
filled as well as coloured. When adding a state, say what its second channel
is.

**A value that was not measured is not zero.** `KvitFigure` draws an em dash,
`KvitBar` draws a tick at the origin, `KvitSpark` draws a baseline tick, and a
chart with no data shows a `KvitEmptyState` rather than an axis drawn around
zeros. A balance nobody has computed and a balance of zero are different
facts, and drawing them the same way states something false.

**A count is handed over as a number and a noun, never as a finished
phrase.** The components that say how many of something there is —
`KvitViewHead`, `KvitSectionHeading`, `KvitBadge`, `KvitSidebarItem`,
`KvitSearchField` — take the number in `count` (or `matches`) and the word for
one of the things in `counted`, plus `countedPlural` where English's suffixed s
is wrong. The component then groups the digits the way the reader's locale
groups them and picks the form, so one account reads "1 account" and twelve
hundred read "1,200 accounts". A call site that assembles the sentence itself
writes "1 accounts" for a workspace with one account and "1200 accounts" for
one with twelve hundred, in a window whose ledger writes the same number as
"1,200".

**A control is a control, not a rectangle with a mouse area.** Use
`KvitButton`, `KvitIconButton`, `KvitCheck`, `KvitSwitch` and the rest. A
`Rectangle` with a `MouseArea` reaches no assistive technology at all: a
screen reader is not told it is there, the keyboard cannot get to it, and
Space does nothing. Anything whose whole label is a symbol also needs a
`label` in words, and that one string becomes both the tooltip and the
accessible name.

**A symbol is asked for by meaning.** `KvitIcon { name: "chevron-right" }`,
never a literal `›`. The catalogue lists the names. An unknown name draws a
marked placeholder and fails the build rather than drawing nothing.

## What each colour means

The token names say what a colour is *for*. Using one for something else is
how two unrelated things end up looking related.

**Surfaces.** `windowBackground` the window, `panelBackground` a region with
its own ground, `listBackground` a card or a list, `popupBackground` anything
floating, `chipBackground` a small mark, `footerBackground` the status bar.

**Text.** `textPrimary` what the reader is reading, `textSecondary` supporting
text, `textMuted` context, `textFaint` a count or a timestamp, `textDisabled`
a control that cannot be used.

**Lines.** `border` is decorative: a rule between two panels, a separator. It
is deliberately below 3:1 in every theme. `borderStrong` is the
control-boundary token: the outline that says where a button, field, checkbox
or card-you-can-click is, held to 3:1 by the contrast floors in the theme
test. A control outlined in `border` has invisible edges; a separator drawn in
`borderStrong` looks clickable.

**Meaning.** `accent` is progress and selection and never a target. `link` is
a target. `success` is finished and nothing else. `warning` is something to
look at. `danger` is stalled or destructive, and appears only where something
can be settled today. `focusRing` is the keyboard, and nothing else uses it.

**Versions of the same text.** `changedTextBackground` marks text that differs
from another version without saying in which direction; `addedTextBackground`
and `removedTextBackground` are the two directions, for a diff. They are not
`success` and `danger` — a removed line is not an error, and using the colour
that means one everywhere else says it is. All three are close in luminance
so that text stays equally legible on any of them, which means they separate
by hue alone and are never the only mark: a diff draws the `+` or `-` in the
gutter and tints the line behind it.

**The portfolio vocabulary**, used by kvit-hub and available to anything else:
`axisAttention` is the reader's own time and `axisAgent` is agent time,
everywhere, and the two never share a scale or a total. Each has a separate,
more legible `…Text` partner for small type. `scopeDiscovered` is work that
turned up rather than being planned. `signalHard`, `signalSoft` and
`signalHygiene` are the three severities.

**Charts.** Three ramps, per theme, all checked by `tests/test_palette`
against that theme's surfaces and against every hue that already means
something. `Theme.categoricalRamp` and `Theme.categorical(i)` for series
identity, `Theme.sequentialRamp` for magnitude, `Theme.divergingRamp` for
distance from a reference. Never a semantic colour for a chart series: a bar
in `danger` says "stalled" to anyone who has read the conventions, whatever
the legend claims.

## The workflow for a new screen

Four steps, in order. The first two are cheap and the last two are not, which
is the point: the drawing is where a screen gets rejected.

**1. Draw it as HTML.** A static page against `ux/tokens.css` and
`ux/frame.css` in this repository. `tokens.css` is generated from the same C++
token table the applications draw with, so the drawing is in the colours the
application will actually use. `ux/PATTERN.md` is the contract for what a
drawing may and may not do; `ux/visual-language.md` is what the tokens mean.
Nothing in a mockup loads anything over the network.

**2. Render it and put the image in front of the owner.** `ux/render.sh`
writes a PNG at 1440×960. The `kvit-preview` skill wraps this. A screen is
approved by looking at it rather than by reading a description of it.

**3. Implement it in QML from the components.** Assemble it; do not redraw it.
If a screen needs something the vocabulary does not have, that is a component
to add to `kvit-ui` rather than a shape to draw locally — say so rather than
building it in the application.

**4. Finish with a screenshot run in four themes.** `build.sh --shots` in the
application, or `kvit-ui-gallery --shots <dir>` for a component. Look at the
high-contrast pass in particular: it is where a distinction resting on hue
alone becomes obvious.

## Where things live

- `catalog.md`: every component, generated from the repository.
- `qml/`: the components.
- `src/tokens/`: `Theme`, `Typography`, `InterfaceMetrics`, and the colour
  arithmetic the chart ramps are checked with.
- `ux/`: `tokens.css` (generated), `frame.css`, `render.sh`, `PATTERN.md`,
  `visual-language.md`.
- `gallery/`: the gallery, which is also the documentation and the screenshot
  run.
- `tests/`: what holds all of the above to the rules on this page.

An application keeps its own mockups, its own data packs and its own layer of
components built on top; only the stylesheets, the render script and the two
prose documents are shared.
