# A shared user-interface system for the kvit desktop applications

*Status: proposal, 2026-08-27. The four decisions in section 11 are settled; nothing is built
yet. The measurements are from the working trees on this machine on that date.*

## What this document is

Four desktop applications are being developed in parallel: the note editor **kvit-notes**, the
agent-facing editor **kvit-notes-pro** built on top of it, the portfolio dashboard
**kvit-hub/desktop**, and the personal-finance application **kvit-cash**, which has a product
requirements document and no code. They are all Qt applications, they all draw their own
controls rather than using a stock style, and they are meant to look like one product from one
maker.

Today they do not. Three of them share a colour table and diverge on everything else, and each
new screen costs more than it should because the person or agent building it has to re-derive
decisions that were already made somewhere in one of the other three.

This document proposes a shared user-interface system to fix that: what it contains, where the
code lives, how each application consumes it, and how an agent asked to build a screen is
supposed to express what it is building. It sets out the current state with measurements, names
the specific divergences, proposes a design, and gives a phased delivery plan with acceptance
criteria. The reference point the owner named is
[HuskarUI](https://github.com/mengps/HuskarUI), a Qt/QML component library that pairs a token
system with an agent skill; section 4 says what is worth taking from it.

Two terms are used throughout and are worth fixing now. A **token** is a named design value —
a colour, a type size, a spacing step — that call sites refer to by name instead of writing the
number. A **primitive** is a small reusable visual element built from tokens — a row, a chip, a
button — that screens assemble rather than redraw.

---

## 1. The four applications, and where they stand

| | kvit-notes | kvit-notes-pro | kvit-hub/desktop | kvit-cash |
|---|---|---|---|---|
| Commits | 366 | 142 | (in kvit-hub) | 1 |
| QML files / lines | 104 / 39,408 | 60 own / 18,142 | 86 / 12,676 | none yet |
| C++ files / lines | 577 / 174,974 | own + submodule | 17 / — | none yet |
| Colour access | `Theme.*` singleton | `Theme.*` singleton | `Tokens.*` singleton over `theme` context property | — |
| Hardcoded hex in QML | 49 (39 distinct) | 32 (30 distinct) | **0** | — |
| Type unit | `font.pixelSize` (365 sites) | `font.pixelSize` (276) | `font.pointSize` (348) | — |
| Type source | `Interface.small/body/caption/strong` | same | `Tokens.type*` off `typography.baseSize` | — |
| Literal font sizes | 0 | 2 | 0 | — |
| Bare `spacing:` literals | 59 | 20 | **146** | — |
| Geometry scaling | `Interface.px(n)` | `Interface.px(n)` | unscaled integers | — |
| Component vocabulary | 2 (`KvitDialog`, `KvitShell`) | none | 18 in `parts/`, 22 stranded in `components/` | — |

### kvit-notes

This is the oldest and most disciplined of the three codebases, and it is where the shared
foundation already lives. `src/platform/` holds three C++ objects that QML binds to as
singletons, registered through the `KVIT_QML_SINGLETONS` macro list in `src/qml/qmlsingletons.h`:

- **`Theme`** (`theme.h`, 348 lines; `theme.cpp`, 614) exposes every colour as a `Q_PROPERTY`,
  filled from four static tables — light, dark, sepia, high-contrast. It also carries the
  reduced-motion setting and a `motionScale` that animations multiply their durations by. The
  header's comments state the rules, including the distinction between `border` (decorative,
  deliberately below 3:1) and `borderStrong` (a control boundary, held to 3:1 by contrast floors
  asserted in `tests/test_theme.cpp`).
- **`Typography`** (`typography.h`) owns the *document's* type: family, base size, line height,
  paragraph spacing, maximum content width, and a frozen ratio table so one setting scales a
  document coherently. `DefaultBaseSize` is 14 and the sizes it returns are pixels.
- **`InterfaceMetrics`** (`interfacemetrics.h`), bound in QML as `Interface`, owns the *chrome's*
  type and geometry: five roles (`caption` 10, `small` 11, `body` 12, `strong` 13, `title` 15 at
  the default base of 12) and an `Interface.px(designPx)` function that scales a design pixel
  value. Its header explains why geometry has to travel with type: a 24-pixel label inside a
  28-pixel button clips.

The split between `Typography` and `Interface` came out of `accessibility.md` Finding 4, and it
is the single most reusable decision in the whole estate. Document text and chrome text are two
independent settings because a reader who wants large body text in a dense note list is asking
for something coherent.

The supporting tooling is also further along than elsewhere: 223 screenshots, `record-gallery.sh`,
a `uidriver.cpp` for driving the interface under test, `run-qmllint.sh`, and
`check-accessible-names.py`.

The weaknesses are 49 hardcoded hex values still in QML, 59 bare `spacing:` literals, and the
absence of a component vocabulary. With 104 QML files and only `KvitDialog` and `KvitShell`
named as reusable types, a row or a chip is redrawn wherever it is needed.

### kvit-notes-pro

The build is a superbuild: `core/` is a git submodule pointing at `file:///home/sk/kvit-notes`,
and `term/` at `kvit-term`. As of today the submodule checkout is byte-identical to the
kvit-notes working tree — `diff -rq core/src ~/kvit-notes/src` and the same over `qml/` both
report zero differing files. `CLAUDE.md` forbids editing anything under `core/`, and
`.claude/settings.json` denies it at the tool level.

That matters more than it looks. The mechanism for sharing code between these applications
already exists, is in daily use, and has not drifted.

The application's own interface is `qml/works/` and `qml/agent/`, 60 files and 18,142 lines. Its
token discipline is better than its reputation: 594 `Theme.*` references, 276 `font.pixelSize`
sites of which only two are literals, and 20 bare `spacing:` literals. What is wrong with it is
structural rather than chromatic. `WorksShell.qml` is 2,481 lines and `WorksSidebar.qml` is
2,042; there is no component vocabulary of its own, so layout, state and drawing sit together in
a handful of very large files. That is the shape a codebase takes when there are no primitives to
assemble, and it is the failure mode a shared system is meant to prevent rather than a failure of
care.

### kvit-hub/desktop

The newest of the three, built in milestone `m2-dashboard` and rewritten in `m3-interface-read`,
and the only one with an explicit token file: `app/qml/theme/Tokens.qml`, 141 lines, a QML
singleton with 1,497 references across the tree. It proxies colour from the shared `theme`
object, derives a nine-step type scale from one base size, and adds a density block — row
heights in four variants, chip and tag and pill heights, five radii, a hairline, three reflow
breakpoints. There are **zero** hardcoded hex colours in 12,676 lines of QML, which no other
application in the estate can claim.

It also has the drawing pipeline. `ux/mockups/v2/` holds `tokens.css` (782 lines, four themes, 64
component classes) and `frame.css` (665 lines, 47 window and layout classes), a set of static HTML
mockups, and `render.sh`, which screenshots each one at 1440×960 with headless Chromium. Two prose
documents — `visual-language.md` and `PATTERN.md` — explain what every token means, where it came
from, and the rules a drawing must not break. `GalleryView.qml` renders every vocabulary element
in all four themes, and `build.sh --shots` writes 44 screenshots of the real application.

Its gaps are in geometry and in component layering. Against 38 spacing values taken from tokens
there are 146 bare literals, concentrated at 8 (44 occurrences), 6 (22), 2 (18), then 5, 4, 3, 1,
10 and 12 — nine unnamed values below `stackGap: 7`, which is the smallest token. Radii are 46
tokenised against 30 literal. There are two component folders: `parts/` holds 18 elements written
against `Tokens` and imported by 32 of the 43 views, while `components/` holds 22 written against
the raw `theme` object, imported only by `Frame.qml` (five of them) and the hidden gallery. Six
are referenced by nothing at all, and `ScopeBar` exists in both folders in two versions.

### kvit-cash

One commit: a 2,487-line `prd.md`, thirteen written user flows under `ux/flows/`, and
`ux/money-display.md`, a cross-cutting specification for how an amount reaches the screen. No
code.

Two things in that PRD bear directly on this proposal. The technology decision (§ around line
1965) is **Qt Widgets for the transaction browser and other dense data surfaces, QML for the
dashboard, onboarding, wizards and the configuration editor, composed within one application**.
So the system cannot be QML-only if kvit-cash is to use it. And `money-display.md` specifies three
kinds of figure — native, converted, and same-currency — with rules about when a currency marker
appears. That is the same shape of problem kvit-hub solved with its `figure` element and its
"never a bare zero, never an empty cell" rule, which suggests a measured-quantity primitive
generalises across the two applications.

Being empty makes kvit-cash the useful test. If the system works, its first screen is assembled
from primitives; if it does not, kvit-cash will grow a fourth private vocabulary and the estate
will be worse off than it is now.

---

## 2. What already exists to build on

The estate is not starting from nothing. Six assets are already built, tested and in use:

1. **The colour token table** in `theme.h` / `theme.cpp`: four complete themes, one `Q_PROPERTY`
   per token, all notifying through a single `themeChanged` signal so a theme switch repaints
   wholesale. Contrast floors are asserted in `tests/test_theme.cpp`.
2. **The two-axis type system**: `Typography` for documents, `InterfaceMetrics` for chrome, with
   `Interface.px()` carrying geometry along with type.
3. **The sharing mechanism**: kvit-notes-pro consumes kvit-notes as a pinned git submodule with
   agent edits denied inside it, and has done so without drift.
4. **The drawing and judging pipeline**: `tokens.css` + `frame.css` + `render.sh` +
   `visual-language.md` + `PATTERN.md` in kvit-hub, and `GalleryView.qml` + `build.sh --shots`
   rendering every element in four themes.
5. **The accessibility work**: `accessibility.md` with numbered findings, the reduced-motion
   setting with a `motionScale` every animation reads, and the high-contrast theme used as the
   check that no distinction rests on colour alone.
6. **The icon font**: kvit-notes-pro's `resources/fonts/Phosphor.ttf` (MIT, 477 KB) and the
   78-line `qml/works/WorksIcon.qml` that draws from it. A call site asks for a symbol by a name
   that says what it means and the component resolves the codepoint, which is the same
   name-not-value discipline the colour tokens have.

The proposal below is mostly a matter of moving these into one place, filling the gaps between
them, and writing down how to use them.

---

## 3. Where the applications diverge

Eight divergences, each verifiable in the working trees.

**3.1 Two units for type.** kvit-notes and kvit-notes-pro set `font.pixelSize` (365 and 276 call
sites). kvit-hub sets `font.pointSize` (348 sites). At 96 dpi a point is 4/3 of a pixel, so the
same nominal number is a different size in the two families, and no conversion happens anywhere.

**3.2 Two ways of reaching the theme.** The editors bind to `Theme`, a QML singleton registered
in C++ through `KVIT_QML_SINGLETONS`. kvit-hub binds to `theme`, a lowercase context property
installed in `app/main.cpp:245`. A component written for one will not compile in the other.

**3.3 The chrome metrics are used by two applications and ignored by the third.** kvit-hub does
not reference `Interface` anywhere. Its chrome type comes from a nine-role scale derived from
`typography.baseSize`, and its geometry — row heights, chip heights, radii — is written as
unscaled integers in `Tokens.qml`. So the interface-size setting that works in the editors does
nothing in the dashboard, and the dashboard's geometry does not move when its type does.

**3.4 One settings key, two meanings.** kvit-hub's settings panel has a row labelled "Interface
size" whose stepper writes `typography.baseSize` and prints the value as `%1 pt`
(`app/qml/components/SettingsPanel.qml:182-205`). In kvit-notes the same object is the *document*
font size, in pixels, defaulting to 14; kvit-hub's `Tokens.qml` defaults it to 15 and treats it as
points. Both applications read and write the same settings store. A user who runs both will find
one changing the other in a way neither labels.

**3.5 Two type scales with no mapping.** `Tokens.qml` defines nine roles (`typeTitle` down to
`typeMicro`) as integer offsets from a base. `InterfaceMetrics` defines five (`caption`, `small`,
`body`, `strong`, `title`) as fixed pixel values scaled by a ratio. Nothing states which of the
nine corresponds to which of the five.

**3.6 Density tokens exist in one application only.** Row heights, chip heights, radii, gaps and
reflow breakpoints are named in `Tokens.qml`. The editors have no equivalent; they pass design
pixels through `Interface.px()` at each call site, so the values are scaled but unnamed.

**3.7 Three ways of drawing a symbol.** kvit-notes' `IconButton.qml` takes a `glyph` string and
its three call sites pass literal Unicode characters — `›`, `‹`, `×`, `+`, `▰`, `▤`, `⌕`, `N`, `#`
— rendered in whatever font on the reader's machine happens to have them. kvit-notes-pro asks
`WorksIcon.qml` for a symbol by name and gets a glyph of the Phosphor font, which is bundled with
the application and rasterized for the exact pixel size it is drawn at. kvit-hub has no icon
system: each view that needs a symbol draws it with `Canvas`, `Shape` or a stack of rectangles.

**3.8 Two ways of consuming the shared code.** kvit-notes-pro uses a git submodule.
kvit-hub/desktop uses `add_subdirectory()` against the hardcoded absolute path
`/home/sk/kvit-notes` (`desktop/CMakeLists.txt:66`), so its build depends on a sibling working
tree at a fixed location and pins nothing.

---

## 4. What HuskarUI does, and what to take from it

HuskarUI is an Ant Design component kit for QML: roughly 80 components published as a QML module
named `HuskarUI.Basic`, built with CMake against Qt 6.7 or later. Its relevant structure:

- **`src/cpp/theme/hustheme.h`** — a C++ singleton holding a token table that components query at
  runtime. Light and dark are token variants of the same names, and changing the primary colour
  recalculates the table and pushes updates through bound properties.
- **Component naming** — every type carries a `Hus` prefix and sits in one of eight functional
  categories (General, Layout, Navigation, Data Entry, Data Display, Feedback, Effects, Utilities).
- **`gallery/`** — a showcase application that doubles as the documentation. `Global.qml` holds
  the component metadata, and each component has an example page (`ExpDrawer.qml` and so on)
  showing its states and a working code sample.
- **`agent/skills/`** — two agent skills. `huskarui` answers questions about components,
  properties and examples by querying a JSON metadata file generated from the repository.
  `qmlpreviewer` renders a QML file with `qmlscene` and captures a screenshot so the agent can
  check what it built. The documented intent is that a user says "create a dashboard page with a
  sidebar and top header" and the agent looks up the component specifications itself rather than
  guessing.

**Worth copying:** the C++ token singleton (already present here), the prefixed component
vocabulary in one importable module, the gallery as both showcase and documentation, the
machine-readable component catalogue, and the pairing of a catalogue skill with a preview skill.

**Not worth copying:** the Ant Design visual language and its component set. These applications
have their own drawn language, documented in `visual-language.md`, and it is well developed —
five-level health phrases with a shape as well as a hue, hatching that survives grayscale,
confidence carried by border style rather than colour. Adopting a generic kit would discard work
that the owner has already judged and accepted. The category taxonomy is also aimed at a general
public library and is heavier than four applications need.

**One thing HuskarUI does not solve here:** it is QML-only. kvit-cash needs Qt Widgets for its
dense data surfaces, so the token layer has to serve both.

---

## 5. The proposed system

### 5.1 Shape: three layers

**Layer 0 — tokens.** Colour, type and density as named values, defined once in C++, published to
QML as singletons and to Qt Widgets as a palette and a metrics object. No call site in any
application writes a colour literal or a font size.

**Layer 1 — primitives.** A vocabulary of small visual elements built from Layer 0, published as
one QML module. Rows, headings, buttons, chips, tags, figures, panels, panes, popovers, empty
states. Each is drawn once and looks the same in every application.

**Layer 2 — the application's own vocabulary.** What is specific to one product and does not
belong in a shared library: kvit-hub's axis bars and milestone chips and health phrases,
kvit-notes' block delegates, kvit-cash's money figures and account rows. These are built from
Layers 0 and 1 and live in their own repositories.

The line between Layer 1 and Layer 2 is whether a second application would use it. A row with a
name, a description and a right-aligned figure is Layer 1. A milestone chip whose border recolours
to the worst signal on that milestone is Layer 2.

### 5.2 Where the code lives, and how an application gets it

**Recommendation: a new repository `~/kvit-ui`, consumed by every application as a pinned git
submodule.** kvit-notes moves `theme.*`, `typography.*` and `interfacemetrics.*` out of
`src/platform/` and becomes the library's first consumer.

The reasoning is that today the shared code is inside kvit-notes, which means kvit-cash would have
to build a note editor to obtain a button. kvit-hub/desktop already pays a version of that cost:
it adds the whole of kvit-notes as a subdirectory to get `kvit-core`. Separating the interface
layer from the editor is what lets a fourth and fifth application be cheap.

The cost is real and should be stated plainly. kvit-notes has 366 commits and 174,974 lines of
C++; moving three platform classes out of it touches its build, its tests and its packaging. The
alternative — leave the token classes in kvit-notes and add the QML primitives module alongside
them inside `kvit-core` — is less work now and leaves kvit-cash carrying the editor as a build
dependency. Section 11 puts this to the owner as a decision.

Consumption is by submodule in all four applications, replacing kvit-hub's hardcoded
`add_subdirectory("/home/sk/kvit-notes")`. That gives every application a pinned version it
upgrades deliberately, and it reuses the arrangement kvit-notes-pro has already proven, including
the practice of denying agent edits inside the submodule.

### 5.3 The token layer

Three groups, one C++ source, three published forms.

**Colour** stays as it is: `Theme`, four themes, one property per token, one change signal.
kvit-hub's portfolio-specific colours (the two axis hues, discovery violet, the three signal
severities, the hatch stroke) are already in `theme.cpp` and stay there, since the alternative is
each application inventing hues that fail the high-contrast check.

**Type** resolves divergences 3.1, 3.3, 3.4 and 3.5 by making `InterfaceMetrics` the single source
for chrome type in every application, in pixels, with `Typography` reserved for document text in
the applications that have documents. The five existing roles are extended to cover what kvit-hub
needs — its nine roles collapse onto a scale of roughly seven once `typeName` and `typeBody` are
reconciled against `body` and `strong` — and the exact mapping is worked out during Phase 1 by
rendering both scales side by side rather than decided on paper. kvit-hub's "Interface size"
stepper is repointed at `Interface.fontSize`, which is what its label already claims.

**Density** is the new work. `Tokens.qml`'s density block moves into C++ beside `Interface.px()`
so that every named value scales with the interface size, and it gains the small end that is
currently missing: a spacing scale naming the 1, 2, 4, 6, 8, 10 and 12 pixel steps that kvit-hub
writes 146 times as literals and kvit-notes writes 59 times. Everything a call site can name is a
call site that cannot drift.

**Publication.** To QML as singletons — `Theme`, `Interface`, `Typography` — using the existing
`KVIT_QML_SINGLETONS` registration, which retires kvit-hub's lowercase context properties. To Qt
Widgets as a `QPalette` built from the same tokens plus a small metrics helper, so kvit-cash's
transaction browser and its QML dashboard restyle together. The Widgets bridge is new code, and it
is what makes the system usable by the application that has not been written yet.

### 5.4 The primitive layer

A first inventory, drawn from what the four applications already contain or specify. Names take
the `Kvit` prefix that `KvitDialog` and `KvitShell` already use.

**Window and structure** — `KvitWindow` (the application shell), `KvitSidebar` and
`KvitSidebarItem`, `KvitHeader`, `KvitBreadcrumb`, `KvitRegion` (a body that takes the leftover
height and scrolls its own overflow), `KvitViewHead` (the strip at the top of a view carrying
counts and controls), `KvitStatusBar`.

**Content** — `KvitSectionHeading` (with an optional count and a hairline to the right edge),
`KvitRow` and `KvitSlimRow`, `KvitCard`, `KvitPanel`, `KvitPane` (the side pane), `KvitDivider`,
`KvitDisclosure`, `KvitEmptyState`.

**Marks** — `KvitChip`, `KvitTag`, `KvitBadge`, `KvitSlug` (a monospace identifier), `KvitPip`,
`KvitDot`.

**Symbols** — `KvitIcon` and `KvitIconButton`. The icon font moves into the library and the
component is kvit-notes-pro's `WorksIcon` with its name-to-codepoint table generated from the
font's stylesheet rather than hand-written, so a consuming application gets the font by importing
the module and adds a symbol by naming one. `KvitIconButton` is kvit-notes' `IconButton`, which is
the accessible form: a real `AbstractButton` that Qt publishes with a role and a name, takes tab
focus, and derives its tooltip and its accessible name from one property so the two cannot
drift.

**Quantities** — `KvitFigure` (a measured value: tabular numerals, a unit in muted colour, an em
dash when the value is not measured rather than a zero). The bars, sparks, trends, gauges and
tables that build on it are set out separately in 5.5, because they carry rules of their own.

**Controls** — `KvitButton` in primary, ordinary and quiet forms, `KvitStepper`, `KvitField`,
`KvitSearchField`, `KvitCheck`, `KvitSelect`, `KvitTab`.

**Feedback** — `KvitToast`, `KvitPopover`, `KvitHoverCard`, `KvitDialog`, `KvitTooltip`.

That is 38 types, and 5.5 adds eleven more. The evidence that this is the right cut: 47 of the 64 classes in
kvit-hub's `tokens.css` and most of the 47 in its `frame.css` fall into it, kvit-notes redraws
perhaps a dozen of them repeatedly across 104 files, and `money-display.md` describes `KvitFigure`
without knowing it exists.

The migration path for kvit-hub is short, because its `parts/` folder is already 18 of these
written against tokens; most become the shared implementation rather than being rewritten.

### 5.5 Data-display components: plots, tables and dense surfaces

The vocabulary in 5.4 is chrome. The screens that cost the most to build are the ones that show
data, and those are exactly where each application has independently invented something.

**What exists today.** kvit-hub draws ten of them by hand: `AxisBar` (48 lines), `ScopeBar` (65),
`MomentumSpark` (82), `EffortJourney` (186), `BudgetGauge` (84), `EffortBar` (65),
`ActivityStrip` (65), `DayStrip` (110), `ProjectTable` (308) and `StatsCell` (37). kvit-notes has
the heavy tabular work: `TableBlock.qml` at 1,694 lines, `QueryBlock.qml` at 841,
`KanbanBlock.qml` at 2,627 and `DiagramBlock.qml` at 1,307. kvit-notes-pro adds `DiffBlock.qml`
and `ReviewFileList.qml`. kvit-cash needs a 250,000-row transaction table, a decade-long net worth
trend with per-point historical currency conversion, cash flow against the prior month, spending
by category, budget progress, category and merchant trends, a recurring list, a data-freshness
list, an import preview, and a side-by-side mapping editor.

**No charting library is used anywhere in the estate.** There are zero references to QtCharts or
QtGraphs; everything is `Rectangle` and `Repeater`, with `Canvas` appearing in only two files
(`Hatch.qml` and `DashBorder.qml`). The proposal keeps it that way, for a specific reason: the
conventions in `visual-language.md` — hatching that marks a figure as an upper bound and survives
grayscale, confidence carried by border style, a shape as well as a hue on every health level —
are not things a charting library exposes, and QtCharts would arrive with its own theming to fight.

**The proposed set**, in three groups.

*Marks — what a value is drawn as:*

- `KvitBar` — one value on a scale shared down a column, with a solid spent portion and an
  outlined remaining portion. kvit-hub's `AxisBar` generalised.
- `KvitStackedBar` — parts of a whole, with a two-pixel surface gap between segments so adjacent
  fills never touch.
- `KvitSpark` — an inline per-period strip. A period that was never measured draws a baseline
  tick rather than a zero-height bar, which is kvit-hub's existing `ActivityStrip` rule.
- `KvitTrend` — a time series with a value axis, a hover crosshair and a tooltip. New; kvit-cash's
  net worth chart is the first caller.
- `KvitDistribution` — categories ranked by magnitude. New; spending by category.
- `KvitGauge` — one value against a limit. kvit-hub's `BudgetGauge` generalised; kvit-cash's
  budgets are the second caller.
- `KvitDelta` — a change with its direction and its period, as in "+4.2% against last month".

*Tiles:*

- `KvitStatTile` — a headline number with a label, an optional delta and an optional spark. Six of
  kvit-cash's seven dashboard widgets are this shape.
- `KvitFigureBlock` — a figure with its label above it, which `frame.css` already defines.

*Tables:*

- `KvitTable` — the shared dense table: a C++ `QAbstractItemModel` behind a `TableView` with
  `reuseItems`, a `HorizontalHeaderView` for column resize and reorder, a `selectionModel` for
  multi-select, saved column sets, keyboard navigation, and copy of a selection. This is the
  single most expensive item in the whole proposal and the one that decides whether kvit-cash can
  be QML throughout.
- `KvitCell` in its variants — text, figure, chip, slug, date.
- `KvitEmptyState` from 5.4 covers the no-data case for every mark and table here: a chart with
  nothing in it describes what will appear rather than drawing an axis around zeros, which
  kvit-cash's flow 08 asks for by name.

**The common design.** Nine rules make these one family rather than a collection.

1. The data's job picks the form, and the answer is sometimes a stat tile rather than a chart.
2. One value axis per chart, never two scales on one plot. kvit-hub already holds a version of
   this rule: attention time and agent time never share a scale or a total.
3. Colour is assigned by the job it does — categorical for identity, sequential for magnitude,
   diverging for polarity, status for state — and status hues are reserved rather than reused as
   another series.
4. Every distinction carries a second channel besides hue. The high-contrast theme is the check,
   and it already exists.
5. A measured value and an unmeasured one never look the same. A spark's unmeasured period draws
   a tick; a table cell with no value shows an em dash and never a zero.
6. Figures use tabular numerals so a column of them aligns.
7. Text wears text tokens; the coloured mark beside a label carries identity, and the label
   itself stays in primary, secondary or muted ink.
8. A legend is present from two series upward, and four or fewer are also labelled directly.
9. Any chart with more than one datum ships a hover tooltip.

**The gap this review found: there is no categorical palette, and one cannot be assembled from the
existing tokens.** The dark theme has six distinct saturated hues — accent blue `#5c9fe0`, axis
amber `#d9a04c`, axis teal `#4aa3a3`, discovery violet `#a37fd4`, success green `#5abd82` and
danger red `#e06c60`. Run through a six-check palette validator against the dark surface, that set
fails three checks: red and green separate by only ΔE 5.7 under deuteranopia against a target of
8, teal reads as gray at chroma 0.085, and three of the six sit outside the usable lightness band.

They also collide semantically. `signalSoft` and `warning` are both `#e0a34c`; `signalHygiene` and
`scopeDiscovered` are both `#a37fd4`; `signalHard` and `danger` are both `#e06c60`. So the estate
has roughly five saturated hues and every one of them already means something specific.

kvit-cash needs categorical colour for spending categories and cannot borrow any of these. Phase 1
therefore has to add three things that do not exist today: a categorical ramp of eight steps that
passes the validator against all four theme surfaces and shares no step with a reserved semantic
hue, a sequential ramp in one hue from light to dark, and a diverging pair with a neutral midpoint.
This is new design work rather than a rearrangement of existing values, and it is the largest
single gap the review turned up.

### 5.6 Components to plan for

HuskarUI publishes 76 components. Most of them are Ant Design's set and have no caller here, but
reading that list against the four applications separates cleanly into three groups, and the first
two are worth committing to now rather than discovering one at a time.

**Group A — already built by hand, several times over.** Counting files across kvit-notes,
kvit-notes-pro and kvit-hub that contain a hand-rolled version:

| Component | Files with a private version | Why it recurs |
|---|---|---|
| Scroll bar | 32 | every scrollable surface; kvit-notes has `DocumentScrollBar.qml` |
| Context menu | 18 | `EditorContextMenus`, `BlockMenu`; kvit-cash needs right-click on a transaction row |
| Collapsible section | 18 | kvit-cash's dashboard sections are hideable, and its import preview expands sections by content |
| Tree | 12 | kvit-notes' file tree and outline, notes-pro's session sidebar, kvit-hub's day tree, kvit-cash's category hierarchy and account groups |
| Switch | 10 | every settings surface |
| Radio group | 10 | kvit-cash's mapping editor needs map / create / leave-unmapped per row |
| Progress | 5 | kvit-cash needs cancellable phased import progress that states its checkpoint |
| Slider | 4 | kvit-notes' typography settings |
| Split view | 3, in three different applications | kvit-hub's `DayView.qml` and `DecisionsView.qml`, notes-pro's `WorksShell.qml` |

Nine types, each existing between three and thirty-two times. Extracting them costs less than the
count suggests, because a working version can be lifted rather than designed. The collapsible
count also says that `KvitDisclosure` in 5.4 needs a multi-section form rather than only the
single trigger-and-body it has in `frame.css` today.

**Group B — absent from every application, and named in a flow that has already been written.**

- **Segmented control.** kvit-cash's dashboard needs one period control governing all seven
  widgets (flow 08). Nothing in the estate has one.
- **Type-ahead field.** Flow 05 requires that "the category and tag pickers are type-ahead fields
  rather than menus". kvit-notes' quick switcher and wiki-link menu are two partial versions of
  the same control.
- **Confirm in place.** Flow 05 puts the undo inside the confirmation for a bulk edit that crosses
  the record threshold, and `prd.md` §47 forbids applying an irreversible bulk edit without
  review. That is a different control from a modal dialog, and HuskarUI's `HusPopconfirm` is the
  shape of it.
- **Timeline.** kvit-cash's per-account update history (§20.3) and its undo history; notes-pro's
  session and conversation history.
- **Persistent notification, separate from the transient toast.** kvit-cash's upstream-change
  alert and its data-freshness prompt persist until acted on, while `KvitToast` in 5.4 is
  transient only. HuskarUI splits these as `HusMessage` and `HusNotification`, and the split
  matches what the flows describe.
- **Numeric field.** A money field is a variant of it, carrying minor units, a currency, and the
  marking rules from `money-display.md`.
- **Dual-list selector.** kvit-cash's configurable columns are an available-versus-shown choice.

**Group C — present in HuskarUI, no caller here.** Pagination (these applications virtualize
rather than page), QR code, rating, carousel, image preview, one-time-password input, the acrylic
and liquid-glass effects, the async hasher, and the frameless-window caption bar with its move and
resize handles, unless the owner decides to draw window decoration too. `HusAvatar` is thin cover
for what kvit-hub's `WhoMark` already does.

**One that needs a judgement rather than a yes or no: the guided tour.** HuskarUI has `HusTourStep`
and `HusTourFocus`, which spotlight a region and step a user through it. kvit-cash's first-run flow
explicitly rejects an onboarding sequence that insists on completion. notes-pro's
`ux1/07-guide-me.md` describes something adjacent: an agent watching a chosen window and pointing
at the screen while coaching by voice, never acting. The spotlight mechanism is the same in both
cases and what drives it is not, so the recommendation is to build the focus-a-region primitive
and leave the stepping to whichever application wants it.

Groups A and B add sixteen types, which puts the library at about sixty-three — HuskarUI's 76 less
the ones this estate has no use for.

### 5.7 The application layer

Nothing changes structurally: each application keeps its own views and its own domain vocabulary.
What changes is that a view imports `Kvit.Ui` and assembles primitives, and an application-specific
element is built from primitives rather than from raw rectangles.

kvit-hub's `components/` folder is retired in the process. Its five live members move to `parts/`
on tokens, the six orphans and the duplicate `ScopeBar` are deleted, and the portfolio vocabulary
that remains — axis bars, milestone chips, scope bars, health phrases — becomes kvit-hub's Layer 2.

### 5.8 Naming and module

QML module: `Kvit.Ui`, imported as `import Kvit.Ui`. Component prefix: `Kvit`. Token singletons
keep their current QML names — `Theme`, `Interface`, `Typography` — since 3,702 call sites across
the two editors already use them and renaming buys nothing.

### 5.9 Qt Widgets

`kvit-cash/prd.md` §48.1 currently assigns Qt Widgets to its dense data surfaces. Decision 3 in
section 11 reverses that: the browser is built on `KvitTable`, the estate stays QML-only, and that
is the assumption the rest of this document makes. The Widgets fallback stays open until Phase 2b's
250,000-row benchmark passes.

If that fallback is ever taken, the token layer still reaches Widgets: a `QPalette` built
from the same `Theme` object and applied at application start, plus a small metrics helper
exposing the same type roles and `px()` scaling to C++ call sites. There would be no Widgets
equivalent of the primitive layer and none is proposed, so what a Widgets surface inherits is
colour, type size and row density, and its rows, figures and chips are drawn a second time in
`QStyledItemDelegate::paint()`. That duplication is the cost Decision 3 weighs.

---

## 6. The drawing layer

kvit-hub's mockup pipeline generalises to the estate with one change: `tokens.css` stops being a
hand-copied mirror and becomes generated.

Its header records that its values were read out of `theme.cpp` and the dashboard's QML. That has
already stopped being true. Comparing five dark-theme tokens against the C++ table on 2026-08-27,
four had diverged: `textPrimary` is `#eeeeee` in `theme.cpp` and `#e8e8e8` in the stylesheet,
`borderStrong` is `#696969` against `#5a5a5a`, `quoteBar` the same pair, and `textFaint` `#858585`
against `#848484`. Drawings are being judged in colours the application no longer draws with, and
nothing reports it.

A small generator run from the token source writes `tokens.css` on demand, and a check fails if
the committed file differs from what the generator produces. kvit-notes and kvit-notes-pro have
`.github/workflows/ci.yml` and can gate it there; kvit-hub has no CI, so its gate belongs in
`desktop/build.sh`. Either way this moves the generator from a good idea to Phase 1 work.

`frame.css` stays hand-written, because it is the window and layout vocabulary rather than a
mirror of anything, and it gains the classes for the primitives in 5.4 and 5.5 that it does not
yet cover — the trend, the distribution, the stat tile and the table have no drawing form today.
Both files, plus `render.sh`, `visual-language.md` and `PATTERN.md`, move into `kvit-ui` so all
four applications draw against one stylesheet. Each application keeps its own mockups and its own
data pack.

---

## 7. The gallery as the contract

`GalleryView.qml` in kvit-hub already renders every vocabulary element in every state across four
themes, and its own header comment says this is where the requirement that no distinction rests on
colour alone is checked once, so that every view inherits verified components. That idea is right
and belongs in the shared library.

The proposal is a gallery application in `kvit-ui` that renders every primitive in every state in
all four themes, following HuskarUI's pattern of the gallery being the documentation. It produces
a fixed set of screenshots, and those screenshots are the thing that gets diffed when a token
changes. Each consuming application keeps its own `--shots` run for its own views, which is
already what kvit-hub's `build.sh --shots` does with 44 images and what kvit-notes' 223
screenshots and `record-gallery.sh` do.

The check that matters is the high-contrast theme. It is where a distinction carried only by hue
disappears, and it is cheap to run on every primitive once rather than expensive to discover in a
view later.

---

## 8. The agent skill

The second half of the owner's request: make agent-built interface work less painful by writing
down how to express intent. HuskarUI's answer is two skills, and the same split fits here.

### 8.1 `kvit-ui` — the vocabulary skill

What it contains:

- **How to express intent.** The rule is that a request is written in the vocabulary, not in
  measurements. "A slim row with the project name, a health phrase, and the attention figure right
  aligned" is a buildable instruction. "A 30-pixel row with #e6b877 text on the right" is not,
  because it names values instead of meanings and will be wrong in three of the four themes.
- **The catalogue**, generated rather than written: every primitive, its properties, its states,
  and a working snippet. This is the JSON metadata file HuskarUI generates from its repository,
  and it is what lets an agent look up a component instead of inventing one.
- **The token reference**: what each colour means and what it is reserved for. Much of this is
  written already — `visual-language.md` records that accent is progress and selection and never a
  target, that green is finished and nothing else, that red is stalled and appears only where
  something can be settled today.
- **The rules that must not be broken**, which exist today as `PATTERN.md`'s "Rules that must not
  be broken" and `visual-language.md`'s "Rules a drawing agent must not break": no colour
  literals; every distinction carries a second channel besides hue; a measured figure and an
  unmeasured one never look the same; identifiers are monospace and never lead a row a person
  reads.
- **The workflow**: draw the screen as HTML against `tokens.css` and `frame.css`, render it with
  `render.sh`, put the image in front of the owner, and only then implement it in QML from
  primitives, finishing with a `--shots` run in four themes.

That workflow already exists in kvit-hub and produced the screens the owner judged. Writing it
down as a skill is what makes it available in kvit-cash and kvit-notes-pro without being
re-explained each time.

### 8.2 `kvit-preview` — the rendering skill

HuskarUI's `qmlpreviewer` renders a QML file with `qmlscene` and captures a screenshot so the
agent can look at what it built rather than assert that it works. The equivalent here wraps three
things that already exist: `ux/mockups/v2/render.sh` for a drawing, `build.sh --shots` for an
application's views, and a new single-component harness that loads one primitive against sample
data in a chosen theme.

The reason this belongs in a skill rather than a README is stated in kvit-hub's own `CLAUDE.md`:
the application's only output is what it draws, so the test suite says nothing about whether a
view is right. An agent that cannot see its own output is guessing.

### 8.3 Where the skills live

In `kvit-ui`, under `agent/skills/`, mirroring HuskarUI's layout, and installed into each
application through its submodule so the four applications cannot drift on how they are described
to an agent.

---

## 9. Delivery plan

Each phase has a check that says whether it worked. The phases are ordered so that the estate is
never worse off partway through.

**Phase 0 — decide.** Done, 2026-08-27. All four questions in section 11 are answered there:
`kvit-ui` is a new repository, the unified type scale is settled during Phase 1 by rendering both
scales rather than on paper, kvit-cash is QML throughout, and kvit-notes-pro's restructuring is
handled as separate work outside this plan.

**Phase 1 — unify the token layer.** Reconcile the two type scales, move kvit-hub to
`font.pixelSize` and to `Interface`, repoint its "Interface size" stepper, add the spacing scale
and the density values to `InterfaceMetrics`, design the three chart ramps from 5.5 (categorical,
sequential, diverging), and write the generator that produces `tokens.css` from the token table.

*Check: kvit-hub's 44 screenshots are re-rendered and reviewed for regressions; the interface-size
setting moves the dashboard's chrome; the settings key no longer collides between applications;
the categorical ramp passes the six-check validator against all four theme surfaces; the generated
stylesheet matches the committed one and the build fails if it does not.*

**Phase 2 — extract the primitives.** Create the module, move kvit-hub's 18 `parts/` in as the
starting set, move kvit-notes-pro's icon font and its `WorksIcon` in as `KvitIcon` so that every
symbol the other primitives draw comes from one named set, fill in the rest of the inventory in
5.4, lift the nine Group A components in 5.6 from whichever application has the best existing
version, and build the gallery.

*Check: the gallery renders every primitive in four themes, and the high-contrast pass shows no
distinction resting on hue alone.*

**Phase 2b — the data-display components.** The marks, tiles and table from 5.5. `KvitTable` is
built against 250,000 synthetic rows from the start, since it is also the prototype that settles
kvit-cash's toolkit question, and it is the one item here that could fail on its own terms.

*Check: `KvitTable` holds smooth scrolling and sub-100 ms filtering at 250,000 rows with twelve
configurable columns, multi-select and inline editing; every mark renders in four themes; a chart
with no data shows an empty state rather than an axis around nothing.*

**Phase 3 — kvit-hub migrates.** Repoint its views at `Kvit.Ui`, retire `components/`, delete the
six orphans and the duplicate `ScopeBar`, and sweep the 146 spacing literals onto the new scale.

*Check: `build.sh --shots` output matches the pre-migration screenshots except where a change was
intended; `components/` is gone.*

**Phase 4 — kvit-cash is built on it.** The first screens of the finance application are assembled
from primitives, with its money figures as its own Layer 2. The Group B components in 5.6 are
built here, since kvit-cash is the first caller for all seven and designing them earlier would be
guessing at their requirements.

*Check: kvit-cash ships a dashboard with no colour literals and no primitive of its own that
duplicates a shared one.*

**Phase 5 — the editors migrate.** kvit-notes replaces its redrawn rows and chips with primitives
and clears its 49 hardcoded hex values. kvit-notes-pro adopts the tokens and primitives too, but
only in the ordinary sense of using them where a screen is touched; breaking `WorksShell.qml` and
`WorksSidebar.qml` apart is separate work under Decision 4 and is not part of this phase.

*Check: zero hardcoded hex in kvit-notes; kvit-notes-pro builds against `Kvit.Ui` with its own
duplicated rows and chips replaced.*

Phase 4 before Phase 5 is deliberate. kvit-cash is the honest test of whether the system is usable
by an application that did not grow up with it, and it is cheaper to find that out on an empty
repository than after rewriting two large ones.

---

## 10. Acceptance criteria

The system has worked if, at the end of Phase 5, all of the following hold:

1. No QML file in any of the four applications contains a colour literal.
2. No QML file contains a numeric font size.
3. Every spacing, radius and row-height value is a named token.
4. One interface-size setting moves the chrome of all four applications, and it is separate from
   the document type setting in the applications that have documents.
5. The four themes, including high contrast, render every primitive correctly in the gallery, and
   no distinction anywhere rests on hue alone.
6. Every symbol in every application is a named icon from the shared set, and no view draws one by
   hand or passes a literal character.
7. `tokens.css` is generated from the token source, and a check fails when the committed file
   differs from what the generator produces.
8. An agent given a screen description in the vocabulary can produce a drawing, render it, and
   implement it without being told any colour or pixel value.
9. Charts in every application draw from the same three ramps, and the categorical ramp passes
   the validator against all four surfaces.
10. A new application reaches its first drawn screen without defining a token or a primitive of its
    own.

Criterion 10 is the one the whole exercise is for.

---

## 11. Decisions needed

**Decision 1 — where the shared code lives.** A new `~/kvit-ui` repository with the three platform
classes moved out of kvit-notes, or a new module inside the existing `kvit-core` in kvit-notes.

The first is the better end state: kvit-cash and any later application get the interface layer
without building an editor, and the dependency graph says what it means. The second is
considerably less work now and leaves the editor as a build dependency of every application
forever. *Recommendation: the new repository, on the grounds that the cost is paid once and
kvit-cash is about to be written either way.*

**Decided, 2026-08-27: the new repository.** `~/kvit-ui` exists and this document is its first
file. `theme.*`, `typography.*` and `interfacemetrics.*` move out of kvit-notes into it, and all
four applications consume it as a pinned git submodule.

**Decision 2 — the unified type scale.** kvit-hub has nine roles as offsets from a base;
`InterfaceMetrics` has five as scaled fixed values. They have to become one scale, and the
question is how many steps it has and what each is for. *Recommendation: settle this in Phase 1 by
rendering both scales side by side in the gallery and judging the images, rather than deciding
from the numbers.*

**Decided, 2026-08-27: as recommended.** The scale is not fixed on paper. Phase 1 renders the nine
kvit-hub roles and the five `InterfaceMetrics` roles side by side, the owner judges the images, and
the resulting scale is written back into this section.

**Built, 2026-08-27, pending the owner's judgement of the images.** Seven roles, in pixels, off the
existing base of 12. The five `InterfaceMetrics` names keep their meanings and their values, so
kvit-notes' chrome is pixel-identical after Wave 2 and its call sites need no rename; `headline` and
`display` are new, and cover what kvit-hub had above `title` and kvit-notes did not.

| Role | px at base 12 | What it is for | Replaces |
|---|---|---|---|
| `caption` | 10 | kind tags, counts | `caption`, `typeMicro` |
| `small` | 11 | chip labels, sub-lines | `small`, `typeSmall`, `typeSecondary` |
| `body` | 12 | row text, prose | `body`, `typeRow`, `typeBody` |
| `strong` | 13 | a name, an emphasised row | `strong`, `typeName` |
| `title` | 15 | a section heading | `title`, `typeHeading` |
| `headline` | 17 | a pane title, the wordmark | `typePage` |
| `display` | 20 | a page title | `typeTitle` |

kvit-hub's nine collapse onto seven because `typeSecondary`, `typeRow` and `typeBody` were three
names for the two sizes ordinary row text is set at, and `typeName` is what `strong` already meant.

**What this changes, and what to look at.** kvit-hub's scale was written in points off a base of 15,
which at 96 dpi is 12 to 25 pixels; the merged scale is 10 to 20. Its chrome therefore gets
noticeably smaller, and that is the judgement the images are for. Rendering `conventions.html`
through the shared pipeline before and after: 11.4% of pixels differ, and holding the type sizes at
kvit-hub's old values drops that to 0.83% — all of which is the four drifted colours being
corrected. So the whole visible change is the type scale, and it is one image to look at rather
than a number to argue about. `ux/tokens.css` carries kvit-hub's nine old names mapped onto the
merged scale, so the 31 existing mockups keep rendering until Wave 3 moves them over deliberately.

Changing the answer is one line per role in `src/tokens/interfacemetrics.h`; the density values, the
generated stylesheet and the gallery all follow from it.

**Decision 3 — kvit-cash's toolkit.** `kvit-cash/prd.md` §48.1 assigns Qt Widgets to the
transaction browser and QML to everything else, on the grounds that a 250,000-row browser is
readily done with `QAbstractItemModel` and a virtualized view. That split has a cost the section
does not price: the browser is the one screen that could not use this system's primitives, so
rows, figures, chips and the currency-marking rules of `money-display.md` would be implemented a
second time inside `QStyledItemDelegate::paint()`. It also adds a second toolkit to an estate with
roughly 70,000 lines of QML and no Widgets at all. *Recommendation: build the browser in QML on
`KvitTable`, and treat Phase 2b's 250,000-row benchmark as the go/no-go. The Widgets fallback for
that one screen costs nothing to hold open, whereas starting in Widgets and moving later means
rewriting the screen.*

**Decided, 2026-08-27: as recommended.** kvit-cash is QML throughout, with the transaction browser
on `KvitTable`. Phase 2b's 250,000-row benchmark is the go/no-go, and `kvit-cash/prd.md` §48.1
needs updating to match.

**Decision 4 — how far kvit-notes-pro is restructured.** Breaking a 2,481-line shell file apart is
worthwhile and is not required by anything else in this plan. It can be Phase 5, or it can be
deferred indefinitely while the application still benefits from the shared tokens and primitives.
*Recommendation: defer the decision until Phase 3 is done, when the cost of assembling a screen
from primitives is known rather than estimated.*

**Decided, 2026-08-27: separate work.** Splitting the two large shell files is its own project,
planned and scheduled on its own, and nothing in this plan waits on it. kvit-notes-pro still gets
the shared tokens and primitives in Phase 5.

---

## 12. Risks

**The token layer becomes a bottleneck.** Every application waits on `kvit-ui` for a value it
needs. Mitigated by the submodule arrangement, which lets an application pin an older version, and
by keeping Layer 2 open so an application can build its own element without asking.

**The primitive set is drawn too early.** A vocabulary decided before kvit-cash exists may not fit
kvit-cash. This is the reason Phase 4 comes before Phase 5, and the reason the inventory in 5.4 is
described as a first cut.

**The generated stylesheet drifts anyway.** A generator that is not run is worse than a
hand-written file, because it claims to be current. The four drifted tokens recorded in section 6
are what happens without an enforced check, and the build gate in Phase 1's scope is what makes
the generator real.

**`KvitTable` is the item that can fail.** Everything else in this proposal is a rearrangement of
work already done and carries schedule risk rather than technical risk. A dense, configurable,
editable 250,000-row table in QML is the one piece nobody here has built. Phase 2b puts it early
and against realistic data for that reason, and the Widgets fallback for kvit-cash's browser stays
open until it passes.

**Screenshot review does not scale.** Four applications, four themes and thirty-five primitives is
a large number of images for one person to judge. The gallery's value depends on the diff being
against the previous run rather than a full re-review each time.
