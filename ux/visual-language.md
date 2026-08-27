# Visual language for the mockups

## What this is for

Step 3 of the design programme in `ux/` produces static HTML mockups, rendered to PNG with
headless Chromium, which the owner judges as images. A mockup proposes a change to the
portfolio dashboard, a Qt/QML desktop application that shows where every recorded project
stands. For the owner to judge a proposal rather than a picture, the mockups have to look
like the application they are changing: same colours, same type sizes, same row heights,
same treatment of a measured figure and of a missing one. Otherwise every difference on
screen is ambiguous between "this is the proposal" and "this is the mockup being sloppy".

`tokens.css` is that shared appearance, expressed as CSS custom properties and a small set
of base element styles. Every mockup imports it and adds only what its own layout needs.
This file explains what is in it, where each value came from, and the rules a drawing agent
must not break.

Two terms recur and are worth stating once. **Attention** is the owner's own engaged time;
**agent-wall** is wall-clock time of agent runs. Every effort figure in the system is one of
the two, they are never summed, and each owns a colour used for nothing else.

## Where the values come from

The dashboard links `kvit-core` from the owner's document editor (checkout at
`~/kvit-notes`) and shares its theme token system, so the two applications restyle from one
source and look like one product. The colour tables are four C++ structs in
`/home/sk/kvit-notes/src/platform/theme.cpp` — `lightTokens()`, `darkTokens()`,
`sepiaTokens()`, `highContrastTokens()` — each filling the same `Theme::Tokens` struct
declared in `theme.h`. The portfolio-specific colours (the two axis hues, discovery violet,
the three signal severities, the hatch stroke) are defined there too, one value per theme,
even though the editor itself never draws with them.

Type and spacing are not tokens in the application. Type is one base point size with integer
offsets, from `/home/sk/kvit-notes/src/platform/typography.h`: `baseSize` defaults to 15 pt
and every QML call site asks for `typography.baseSize - 4` or similar. Spacing, row heights
and radii are literals in the dashboard's QML under
`/home/sk/kvit-hub/dashboard/app/qml/`. The values in `tokens.css` were read out of those
files, so a token there names a number that is really on screen.

## Token map

Every colour token keeps the application's name, kebab-cased. `Theme::windowBackground`
becomes `--window-background`, `Theme::axisAttentionText` becomes `--axis-attention-text`.
The table below gives the ones a mockup actually reaches for; the stylesheet carries the full
set, including the five code-highlighting colours a mockup of the document panel would need.

| Application token | CSS variable | dark (default) | light |
|---|---|---|---|
| `windowBackground` | `--window-background` | `#1e1e1e` | `#ffffff` |
| `panelBackground` | `--panel-background` | `#252526` | `#f4f4f4` |
| `listBackground` | `--list-background` | `#212122` | `#fafafa` |
| `popupBackground` | `--popup-background` | `#2d2d30` | `#ffffff` |
| `chipBackground` | `--chip-background` | `#37373a` | `#f2f2f0` |
| `textPrimary` | `--text-primary` | `#e8e8e8` | `#1a1a1a` |
| `textSecondary` | `--text-secondary` | `#c8c8c8` | `#555555` |
| `textMuted` | `--text-muted` | `#a8a8a8` | `#666666` |
| `textFaint` | `--text-faint` | `#848484` | `#999999` |
| `border` | `--border` | `#3c3c3c` | `#dddddd` |
| `borderStrong` | `--border-strong` | `#5a5a5a` | `#b0b0b0` |
| `hoverTint` | `--hover-tint` | `#333336` | `#ebebeb` |
| `focusTint` | `--focus-tint` | `#263544` | `#eaf2fb` |
| `focusRing` | `--focus-ring` | `#58a6ff` | `#1f6feb` |
| `selectionTint` | `--selection-tint` | `#2d4356` | `#dce8f5` |
| `accent` | `--accent` | `#5c9fe0` | `#4a90d9` |
| `link` | `--link` | `#6fb1ff` | `#2970c8` |
| `danger` | `--danger` | `#e06c60` | `#b3261e` |
| `success` | `--success` | `#5abd82` | `#27ae60` |
| `warning` | `--warning` | `#e0a34c` | `#f39c12` |
| `axisAttention` | `--axis-attention` | `#d9a04c` | `#d99a3d` |
| `axisAttentionText` | `--axis-attention-text` | `#e6b877` | `#b06a10` |
| `axisAgent` | `--axis-agent` | `#4aa3a3` | `#4aa3a3` |
| `axisAgentText` | `--axis-agent-text` | `#7cc7c4` | `#1f7a7a` |
| `scopeDiscovered` | `--scope-discovered` | `#a37fd4` | `#8a5cc0` |
| `signalHard` | `--signal-hard` | `#e06c60` | `#c0392b` |
| `signalSoft` | `--signal-soft` | `#e0a34c` | `#d99a3d` |
| `signalHygiene` | `--signal-hygiene` | `#a37fd4` | `#8a5cc0` |
| `hatchAlt` | `--hatch-alt` | `#1e1e1e` | `#ffffff` |

Two pairs deserve a note. Each axis has a **bar hue** and a **text hue**: the bar hue is what
a filled bar is painted with, and the text hue is a shifted version legible as small type on
the theme's background. Attention text in light is `#b06a10`, much darker than the bar's
`#d99a3d`; in dark it is `#e6b877`, lighter. A figure printed in the bar hue would be
unreadable in one theme or the other, so `.figure--attention` and `.figure--agent` use the
text hue and `.bar` uses the bar hue. And `--hatch-alt` is the stroke colour of the hatch
that marks an upper bound; it is the theme's own background colour, which is what keeps the
hatch visible on every ground a bar can sit on.

Sepia and high contrast are in the stylesheet as `.theme-sepia` and `.theme-contrast` with
their full tables. High contrast is the accessibility check: pure black ground, white body
text, white structural borders, a yellow focus ring. Rendering a mockup with
`class="theme-contrast"` on `<body>` is the cheapest way to prove that nothing in it depends
on colour alone.

## The default theme

`:root` carries **dark**, and a mockup renders dark unless it says otherwise. The reason is
what the owner's own installation is set to: `theme.id` is `"dark"` in
`~/.config/kvit/kvit-hub-dashboard/settings.json`, and the dashboard's own screenshot pass
renders dark whenever `--theme` is not forced. A mockup should look like the window he has
open.

Light is a class away. Put `theme-light`, `theme-dark`, `theme-sepia` or `theme-contrast` on
`<body>` and the whole page restyles, because everything downstream reads variables.

## Type

The application's scale is one base size with integer point offsets. The stylesheet converts
each step at 96 dpi (points × 4/3) and gives it a role name, so a mockup asks for a role
rather than doing arithmetic.

| Application | CSS variable | px | Where it is used |
|---|---|---|---|
| `baseSize + 4` | `--type-title` | 25 | page title, the load-failure message |
| `baseSize + 2` | `--type-page` | 23 | the shell wordmark |
| `baseSize + 1` | `--type-heading` | 21 | section heading |
| `baseSize` | `--type-name` | 20 | project name, popover title |
| `baseSize - 2` | `--type-body` | 17 | prose, tab labels |
| `baseSize - 3` | `--type-row` | 16 | the dominant row text, and the page default |
| `baseSize - 4` | `--type-secondary` | 15 | descriptions, figures, column headings |
| `baseSize - 5` | `--type-small` | 13 | chip labels, the axis legend, sub-lines |
| `baseSize - 6` | `--type-micro` | 12 | kind tags, counts |

Families are local, since no mockup may fetch anything over the network. `--font-ui` is
`Ubuntu, "Ubuntu Sans", "DejaVu Sans", system-ui, sans-serif`, because Ubuntu is the family
the owner has selected in his settings and it is installed on this machine. `--font-mono` is
`"DejaVu Sans Mono", "Ubuntu Mono", ui-monospace, monospace`, which is what the application's
default `monospace` resolves to here.

Monospace carries identity, and nothing else: a milestone slug, a task slug, an evidence
anchor, a bucket letter. Prose, titles and figures are always the UI family. The current
interface leans on adjacent steps of this scale — `--type-secondary` and `--type-row` differ
by one pixel — which is part of why the owner reads the surfaces as walls of text. Prefer two
steps of separation between a title and the data under it.

## Density

One reader, one data set, a desktop window, no mobile target. The application opens at
1400×900 and refuses to go below 1280×600; the mockups render at 1440×960.

| Measure | CSS variable | px | Source |
|---|---|---|---|
| view outer margin | `--view-margin` | 16 | `LedgerView.qml` `anchors.margins` |
| gap between columns | `--column-gap` | 14 | `LedgerView.qml` `columnSpacing` |
| gap between stacked blocks | `--stack-gap` | 7 | `LedgerView.qml` `ColumnLayout.spacing` |
| shell header | `--header-height` | 52 | `Main.qml` header |
| breadcrumb strip | `--breadcrumb-height` | 34 | `Main.qml` |
| project row | `--row-height` | 56 | `LedgerView.qml` delegate |
| expanded milestone line | `--row-height-sub` | 48 | `LedgerView.qml` delegate |
| disclosure or triage line | `--row-height-compact` | 24 | `LedgerView.qml` history line |
| tab | `--tab-height` | 30 | `Main.qml` tab delegate |
| milestone chip | `--chip-height` | 17 | `MilestoneChip.qml` |
| kind tag, count chip | `--tag-height` | 16 | `LedgerView.qml` `KindTag`, `SignalChip.qml` |
| confidence tag, bucket chip | `--pill-height` | 15 | `ConfidenceTag.qml`, `BucketChip.qml` |
| compact axis bar | `--bar-height` | 7 | `AxisBar.qml`, `EffortJourney.qml` |
| full axis bar | `--bar-height-wide` | 9 | `EffortJourney.qml` |

Radii are `--radius-bar` 2, `--radius-chip` 3, `--radius-control` 4, `--radius-card` 6,
`--radius-pill` 8. Every separator is one pixel of `--border`; the rule under a column header
row is `--border-strong`.

A project row is two lines at 56 px: the name at `--type-name` bold, and a one-line
description in `--text-muted` at `--type-secondary` beneath it. That shape is a settled
change (L1 in `dashboard/new-ui-v2.md`) and mockups should keep it.

## How a figure is rendered

A figure is a measured quantity in hours, on one of the two axes. It is drawn as a graphic
sized to be read, with a small number as its label — never as a full-strength printed number
with a bar too small to carry it.

The graphic is `.bar`: a hairline track that is always drawn, a solid segment for hours spent
in the axis's bar hue, and a lighter outlined segment for hours remaining, both positioned on
one continuous scale shared by every row in that column. Set `--spent` and `--remaining` as
percentages of that shared scale. The two axes get their own row and their own scale, and
they never stack or sum.

The number is `.figure`: `--type-secondary`, tabular numerals so a column aligns, tinted with
the axis's text hue when it belongs to an axis. Precision is whole hours at 10 h and above
and one decimal below, matching `fmt.js`. Tenths of an hour on a 200-hour project are digit
noise, and the exact value belongs in a tooltip.

An upper bound adds the hatch. Apply `.is-bounded` to the spent segment and prefix the label
with `≤`; the two always travel together, so a hatch with no `≤` or a `≤` with no hatch is a
drawing error. Note that at task and milestone level the `≤` is being retired: the owner
ruled on 2026-07-29 that shared evidence is divided among the items it witnesses, giving each
an approximate figure printed with no marker at all. The application has not caught up yet.
For a mockup, follow what the mockup is proposing and say which convention it uses.

A percentage never travels alone. `.scope` draws the item scope as one bar — closed solid,
open as outlined ground, the discovered share hatched violet across the tail — and the
percentage sits beside it with its closed, open and discovered counts. Discovery enlarges the
denominator, so the number can move down, and the counts are what make that legible.

## How "not measured" is rendered

An em dash in `--text-faint`, with no unit. The `.dash` class styles it, and an empty
`<span class="dash"></span>` renders the dash on its own.

Zero, none-measured and unknown are three different states and the interface keeps them
apart. A zero-length bar reads as "measured: zero", so an unknown figure draws no ink inside
the track at all — the empty track stays, because a missing bar and a missing column look the
same. `≤ 0 h` never prints: an upper bound of zero means nothing was measured, so it renders
as the dash. Terms whose value is zero and carries no information (`+0 discovered`,
`rework 0`) are omitted rather than printed.

Where a whole object is absent, the surface says so in words on the object itself. A project
with no record shows "Activity measured; no project record", and the reason — which the
registry already holds in a `note` field — is one gesture away. A blank cell means "nothing
here", which is true of some objects and false of others.

## What red is reserved for

Red (`--signal-hard`) marks something the reader can act on today. That is the whole of its
meaning, and spending it anywhere else is the specific failure the v2 review identified: on
the interface as it stands, ordinary sequencing renders in the alarm colour on every blocked
row, which leaves nothing distinct for the one thing on screen that does want action.

In practice: a milestone whose items are all closed and whose gate needs a ruling is red. A
blocker that has stalled is red. A wait whose blocker is progressing is muted ink and reads
"after M1 open-core split". Within-project sequencing is muted always. Amber
(`--signal-soft`) means a decision is wanted; violet (`--signal-hygiene`) means record
cleanup, and it is the same violet as discovered scope.

`--danger` is a separate token for the application failing — the data root could not be read,
a tool exited non-zero. It is not a data state and a mockup should not use it to colour a
figure.

## What an interactive affordance looks like

The application distinguishes four states, and a static mockup cannot hover, so the
stylesheet exposes each as a class you paint deliberately.

- **A link** is `--link` coloured with no underline. `.is-hover` on it turns it `--accent`
  and underlines it. Every drill-in in the application is a link of this kind, including the
  project name in a ledger row.
- **A row** is transparent, `--hover-tint` under the pointer (`.is-hover`), `--focus-tint`
  when it holds keyboard focus (`.is-focus`). A row that expands carries `▸` / `▾` at its
  left, in a fixed 14 px slot.
- **A button** is `.btn`: 26 px, one-pixel `--border-strong`, `--footer-background` fill.
  `.btn--primary` fills with `--accent` and takes `--on-accent` text. `.btn--quiet` drops the
  border and is what the header tools (`?`, `↻`, `⛭`) look like. Keyboard focus is a two-pixel
  `--focus-ring` outline, which `.is-keyboard-focus` paints.
- **A tab** is a pill at `--radius-control`; the current one fills with `--selection-tint` and
  goes bold. A count riding on a tab is a `--signal-soft` badge.

Anything the reader can press should look pressable at a glance in the PNG. A mockup that
depends on a hover state to be understood is a mockup the owner cannot judge, so paint the
state and label it.

## Rules a drawing agent must not break

1. **No network requests of any kind.** No web fonts, no CDN stylesheets, no remote images,
   no `fetch`. Chromium renders these files offline. Use the two font stacks in the sheet,
   and inline any graphic as SVG or a data URI.
2. **No literal colours.** Every colour comes from a variable. A hex code in a mockup is a
   colour that will be wrong in three of the four themes, and it silently breaks the
   high-contrast check.
3. **No distinction resting on colour alone.** The application encodes measurement confidence
   as border style (solid exact, dashed bounded, dotted unmeasured), milestone status as fill
   style (filled closed, partly filled active, dashed planned, dotted unestimated), and an
   upper bound as a hatch. Colour is a second channel on top. If a mockup introduces a
   distinction, give it a shape, a style or a word as well, then render it with
   `class="theme-contrast"` and confirm it survives.
4. **Red only for what can be acted on today.** See above.
5. **Unknown is a dash, never a zero and never an empty cell.**
6. **The two axes never stack or sum.** Attention and agent-wall get separate rows and
   separate scales.
7. **No dates on forward-looking figures.** Everything forward is effort ("8 h to go",
   "2–3.5 working days"), because converting effort to a date needs a model of how the
   owner's attention is allocated across projects and that model does not exist. Only
   historical facts carry dates.
8. **No project-level percentage complete.** A milestone has an agreed scope, so its percent
   is defined. A project spans closed milestones that predate estimates and future ones not
   yet planned, so it has no denominator; project effort appears as absolute measured hours.
9. **Titles for people, identifiers on demand.** A row a person reads leads with the title.
   The slug follows in `.slug`, monospace and secondary. A cross-project reference renders as
   project plus milestone chip plus title, and never as `project:milestone-slug`.
10. **No abbreviation the mockup never explains.** "aw", "a/g", "disc", "conf", "P1" either
    get spelled out or get a legend. The mockup is judged as an image, so a tooltip is not
    available as an escape hatch — anything essential has to be visible.
11. **Match the density.** Rows at 56 px, margins at 16, column gaps at 14. A mockup drawn
    airier than the application will read as a proposal to make it airier, whether or not that
    was the intent.

## Rendering a mockup

`render.sh` in this directory drives headless Chromium over every `*.html` here and writes
PNGs to `renders/`:

```bash
./render.sh                # every mockup
./render.sh today-a.html   # just these
```

The viewport is 1440×960, overridable with the `WIDTH` and `HEIGHT` environment variables. A
page taller than the viewport is captured whole, so a long surface renders in full, and a
mockup meant to be judged as one screen should fit 960 px. Then open the PNG and look at it:
the mockup's only output is what it draws, so nothing else says whether it is right.

To check a mockup in another theme, change the `<body>` class and render again. The
high-contrast pass is worth doing once per mockup, because it is what catches a distinction
that turned out to rest on colour.
