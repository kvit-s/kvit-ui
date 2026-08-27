# Delivery plan for the shared user-interface system

*Written 2026-08-27, against `prd.md` in this repository. `prd.md` says what the system is
and why; this document says in what order it gets built and by whom. The four decisions in
`prd.md` §11 are settled and this plan assumes them.*

## What this is

Four Qt desktop applications — the note editor **kvit-notes**, the agent-facing editor
**kvit-notes-pro** built on top of it, the portfolio dashboard **kvit-hub/desktop**, and the
personal-finance application **kvit-cash**, which has a requirements document and no code — are
meant to look like one product and currently do not. `prd.md` proposes a shared interface layer
to fix that: a token layer (colours, type sizes, spacing and density as named values), a
vocabulary of reusable drawn components built from those tokens, a gallery that renders every
component in all four themes, and two agent skills that describe the vocabulary and let an agent
look at what it built.

That system lives in this repository, `~/kvit-ui`, and each application consumes it as a pinned
git submodule.

The work splits into five waves. **Wave 1 builds the whole system inside this repository and
changes nothing anywhere else** — it is the only wave planned in detail here, because it is the
only one whose shape is known today. Waves 2 to 5 move each application onto the system, and they
are described as scope and sequence rather than as tasks, since what each of them actually costs
is only knowable once the components exist and the first migration has been done.

## The five waves

| Wave | What happens | Repository touched | Planned in detail |
|---|---|---|---|
| 1 | The system is built: tokens, components, gallery, drawing pipeline, agent skills | `kvit-ui` only | yes, below |
| 2 | kvit-notes gives up the token classes and consumes `kvit-ui` instead | kvit-notes | overview |
| 3 | kvit-hub/desktop moves onto the tokens and the components | kvit-hub | overview |
| 4 | kvit-cash is written on the system from its first screen | kvit-cash | overview |
| 5 | The two editors adopt the components in their own screens | kvit-notes, kvit-notes-pro | overview |

Two ordering facts are worth stating up front, because they are not arbitrary.

**Wave 2 comes before Wave 3 for a build reason.** kvit-hub does not have a theme of its own: it
reaches `Theme` by adding the whole of kvit-notes to its build as a subdirectory
(`desktop/CMakeLists.txt`, `KVIT_EDITOR_ROOT` defaulting to `/home/sk/kvit-notes`) and installing
the object as a lowercase `theme` context property in `app/main.cpp:245`. It also needs kvit-notes
for two unrelated things — the layered graph layout behind the dependency view, and
`DocumentExporter` for the document panel — so it cannot simply stop depending on the editor. If
kvit-hub adopted `kvit-ui` while kvit-notes still owned its own `Theme`, the process would hold
two theme objects with two settings stores and two notions of which theme is current. So
kvit-notes hands over ownership first, in a wave that changes no pixels.

**Wave 4 comes before Wave 5 deliberately**, as `prd.md` §9 argues: kvit-cash is the test of
whether an application that did not grow up with this system can use it, and finding that out on
an empty repository is cheaper than finding it out after rewriting two large ones.

---

# Wave 1 — build the system, inside this repository

## What "inside this repository" means, and the one coordination constraint

Every C++ class, QML file, stylesheet, script and document produced in Wave 1 is created here.
The token classes that exist today in kvit-notes are **copied in, not moved**: at the end of
Wave 1, `Theme`, `Typography` and `InterfaceMetrics` exist in both repositories, and kvit-notes
builds and ships exactly as it does today. Wave 2 is what deletes the kvit-notes copy.

That duplication is the price of keeping Wave 1 self-contained, and it has one requirement
attached: **the three source pairs in kvit-notes are frozen for the duration of Wave 1.** No edits
to `src/platform/theme.{h,cpp}`, `src/platform/typography.{h,cpp}` or
`src/platform/interfacemetrics.{h,cpp}`, and no edits to their three tests. A change made there
during the window has to be re-made here by hand, and the version in this repository diverges
immediately once the type scale and the density tokens are reworked, so no automatic check can
catch it later. If something in kvit-notes genuinely needs a token added mid-wave, add it in
`kvit-ui` and cherry-pick the same hunk into kvit-notes, so the two stay reconcilable.

Nothing else in any other repository is read-only: kvit-hub's `parts/` components, its
`tokens.css` and `frame.css`, its `render.sh` and its `GalleryView.qml` are all **copied** into
this repository as starting material during Wave 1, and the originals stay where they are until
Wave 3 retires them.

## Stage 1.1 — the repository and its build

Nothing here is interesting design work, and getting it wrong is expensive later, so it comes
first.

**Tasks**

- `git init`, MPL-2.0 licence headers matching kvit-notes' convention (its
  `tools/apply-license-headers.sh` is the script to copy), `.gitignore` with `build*/`, and a
  `THIRD-PARTY-NOTICES.md` for the vendored font that arrives in stage 1.6.
- `CMakeLists.txt` at the root: `cmake_minimum_required(VERSION 3.21)`, C++20,
  `CMAKE_AUTOMOC`/`CMAKE_AUTORCC` on, `qt_policy(SET QTP0001 NEW)`, and
  `find_package(Qt6 REQUIRED COMPONENTS Core Gui Quick QuickControls2 Test QuickTest)`. Widgets is
  deliberately absent; see the Widgets bridge section below.
- Two build options, copied from the arrangement kvit-notes uses and from
  `qt-project-setup`'s measurements:
  - `KVIT_UI_FULL_DEBUG_INFO` (default OFF) appending `-g1` on GCC and Clang, placed before any
    `add_subdirectory()`, with `-DKVIT_UI_FULL_DEBUG_INFO=ON` as the escape hatch for anyone who
    needs to sit at a debugger.
  - `KVIT_UI_SHARED_LIBS` (default OFF) selecting `SHARED` or `STATIC` for the internal
    libraries, turned ON by `build.sh` rather than in CMake, so a packaging path that passes only
    `CMAKE_BUILD_TYPE` still gets a static build.
- Three targets, named so that a consumer can read the dependency graph:
  - `kvit-ui-tokens` — the C++ token layer (stage 1.2).
  - `kvit-ui-qml` — the `Kvit.Ui` QML module, holding both the type registrations and the
    component QML.
  - `kvit-ui-gallery` — the gallery application (stage 1.10).
- `CMakePresets.json`, a `build.sh` mirroring kvit-notes' (configure, build, test, `--shots`),
  and `win-build.bat` if Windows parity is wanted from the start.
- `.github/workflows/ci.yml`: configure, build, `ctest`, `run-qmllint.sh`, and the generator check
  from stage 1.5. kvit-notes has a working CI file to start from; kvit-hub has none, which is why
  its gate ends up in `build.sh` in Wave 3.
- `CLAUDE.md` for this repository, stating the rules an agent working here must not break — no
  colour literals, no numeric font sizes, every distinction carrying a second channel besides hue
  — and pointing at the skills in stage 1.12.

**Check.** A clean checkout configures and builds on Linux and Windows, `ctest` runs with zero
tests and passes, and a shared development build links. Building shared once at this point is
worth the minute it costs: a `PRIVATE` dependency of a static library still reaches the consumer's
link line, so an undeclared dependency stays invisible until the first shared build.

## Stage 1.2 — port the token classes

**What moves in, and what comes with it.** The three classes named in `prd.md` §2 do not stand
alone. Their actual dependencies, read out of the includes:

| File | Lines | Depends on |
|---|---|---|
| `theme.{h,cpp}` | 348 + 614 | `SettingsStore`, `SystemAppearance`, `PerfLog`, `QGuiApplication`/`QStyleHints` |
| `typography.{h,cpp}` | — | `SettingsStore`, `PerfLog`, and an `#include "blockkind.h"` that nothing in the file uses |
| `interfacemetrics.{h,cpp}` | 98 + — | `SettingsStore`, `PerfLog` |
| `settingsstore.{h,cpp}` | 83 + 176 | Qt only |
| `systemappearance.{h,cpp}` | 85 + 220 | Qt only, plus a Windows native event filter |
| `perflog.{h,cpp}` | 168 + 442 | Qt only |

So six source pairs come across, not three. `SettingsStore` and `SystemAppearance` are already
platform-level and free of editor concepts. `PerfLog` sits in kvit-notes' `domain/` module today
and is Qt-only; it comes here so the ported code keeps its timing instrumentation rather than
having it stripped out and re-added later.

**Tasks**

- Copy the six pairs into `src/tokens/`, keeping the class names, the file names and the comments.
  The comments in these headers are the design rationale — `interfacemetrics.h` explaining why
  geometry has to travel with type, `theme.h` distinguishing decorative `border` from control
  `borderStrong` — and they are the reason the port is a copy rather than a rewrite.
- Drop the unused `blockkind.h` include from `typography.cpp` after confirming by compilation that
  nothing in the file uses it.
- Copy `tests/test_theme.cpp`, `tests/test_typography.cpp` and `tests/test_interfacemetrics.cpp`,
  including the contrast-floor assertions, which are what hold `borderStrong` at 3:1.
- Register the singletons for QML. kvit-notes uses `QML_FOREIGN` wrappers gathered in one file
  (`src/qml/qmlsingletons.h`) so that the service classes stay free of QML machinery; this
  repository keeps that arrangement with its own smaller list, and the module is `Kvit.Ui` rather
  than `Kvit`. The QML names stay `Theme`, `Interface` and `Typography`, since 3,702 call sites
  across the two editors already use them.
- `qt_add_qml_module(kvit-ui-qml URI Kvit.Ui VERSION 1.0 NO_PLUGIN DEPENDENCIES QtQuick)`, plus a
  `kvit-ui-qmltypes` target in `ALL` so `qmllint` always has a current module description to read.
  kvit-notes learned this the hard way: its lint job named a generated target that had moved, and
  the blocking check failed before running any check at all.

**Check.** The three ported tests pass unchanged. A three-line QML file that imports `Kvit.Ui` and
prints `Theme.accent`, `Interface.body` and `Typography.baseSize` runs under `qmlscene`. `qmllint`
resolves every `Kvit.Ui` type.

## Stage 1.3 — settle the type scale and add the density tokens

This is the first stage that changes behaviour rather than moving code, and it closes divergences
3.1, 3.3, 3.4, 3.5 and 3.6 from `prd.md`.

**The type scale (Decision 2).** `InterfaceMetrics` has five roles derived from a base of 12
(`caption` 10, `small` 11, `body` 12, `strong` 13, `title` 15) in pixels. kvit-hub's `Tokens.qml`
has nine, written as integer offsets from a base of 15 in points (`typeTitle` = base+4 down to
`typeMicro` = base−6). `prd.md` §11 settles that the merged scale is judged from rendered images
rather than decided from the numbers, so:

- Build a scale sheet in the gallery harness that renders both scales side by side at the same
  physical size, at the default base and at the two clamp ends (10 and 24).
- Convert kvit-hub's nine roles to pixels at 96 dpi (a point is 4/3 of a pixel) before comparing,
  so the two columns are the same measurement.
- Put the sheet in front of the owner, agree the merged set — `prd.md` expects roughly seven steps
  once `typeName` and `typeBody` reconcile against `body` and `strong` — and write the answer back
  into `prd.md` §11 as the recorded decision.
- Implement the agreed roles on `InterfaceMetrics`, keeping the existing five names so the editors
  do not have to change, and add the new names for whatever kvit-hub needs above and below them.

**The density tokens.** `Tokens.qml`'s density block is 141 lines of unscaled integers: view
margin 16, column gap 14, stack gap 7, sidebar 232, rail 48, pane 392, header 52, breadcrumb 34,
four row heights (56/48/30/24), tab 30, chip 17, tag 16, pill 15, two bar heights, five radii, a
hairline, and three reflow breakpoints (880/1100/1440). All of it moves into `InterfaceMetrics`
in C++ beside `px()`, so that every one of those values scales with the interface size, which it
does not do today.

**The spacing scale.** The gap `prd.md` §5.3 names: kvit-hub writes 146 bare `spacing:` literals,
concentrated at 8 (44 times), 6 (22), 2 (18), then 5, 4, 3, 1, 10 and 12 — nine unnamed values,
all below `stackGap: 7`, which is its smallest named value. kvit-notes writes 59 more. Name the
1, 2, 4, 6, 8, 10 and 12 pixel steps, decide whether 3 and 5 survive as tokens or get rounded onto
a neighbour, and put them through `px()` like everything else.

**The settings key collision (divergence 3.4).** kvit-hub's settings row labelled "Interface size"
writes `typography.baseSize`, which in kvit-notes is the *document* font size. Both applications
read the same settings store, so one silently moves the other. The fix belongs to kvit-hub and
happens in Wave 3; what Wave 1 owes it is `Interface.fontSize` with its `minFontSize` and
`maxFontSize` clamps exposed, so the stepper has something correct to point at.

**Check.** Every value kvit-hub's `Tokens.qml` names has a home in C++, the scale sheet has been
judged and the answer recorded in `prd.md`, and a unit test asserts that changing
`Interface.fontSize` from 12 to 24 moves every role and every density value and leaves
`Typography.baseSize` alone.

## Stage 1.4 — the three chart ramps

`prd.md` §5.5 calls this the largest single gap the review found, and it is design work rather
than rearrangement, so it is scheduled before anything that would consume it.

The estate has roughly five saturated hues and every one already means something: `signalSoft`
and `warning` are both `#e0a34c`, `signalHygiene` and `scopeDiscovered` are both `#a37fd4`,
`signalHard` and `danger` are both `#e06c60`. Run against the dark surface, the six-hue set fails
three checks — red and green separate by only ΔE 5.7 under deuteranopia against a target of 8,
teal reads as gray at chroma 0.085, and three of the six sit outside the usable lightness band.
kvit-cash needs categorical colour for spending categories and can borrow none of them.

**Tasks**

- Write the palette validator as a runnable script in `tools/`, so a ramp is checked rather than
  asserted: contrast against each of the four theme surfaces, pairwise separation under the three
  common colour-vision deficiencies, chroma floor, lightness band, and no step within a stated
  distance of a reserved semantic hue.
- Design a **categorical ramp of eight steps** that passes against light, dark, sepia and high
  contrast, and shares no step with a reserved hue.
- Design a **sequential ramp** in one hue, light to dark.
- Design a **diverging pair** with a neutral midpoint.
- Add all three to `Theme` as token groups, one per theme table, and add the validator to `ctest`
  so a later edit to a ramp fails the build rather than the eye.

**Check.** The validator passes on all three ramps against all four surfaces, and it is part of
the test suite rather than a script somebody remembers to run.

## Stage 1.5 — the stylesheet generator

`tokens.css` in kvit-hub is a hand-copied mirror of the C++ token table, and it has already
drifted: comparing five dark-theme tokens on 2026-08-27, four differ — `textPrimary` `#eeeeee`
against `#e8e8e8`, `borderStrong` `#696969` against `#5a5a5a`, `quoteBar` the same pair, and
`textFaint` `#858585` against `#848484`. Drawings are being judged in colours the application does
not draw with.

**Tasks**

- A generator in `tools/` that reads the token table and writes `tokens.css` — all four themes, in
  the class structure the existing 782-line file already uses, so the mockups that consume it keep
  working unchanged.
- A check that fails when the committed file differs from what the generator produces, wired into
  this repository's CI in stage 1.1 and into kvit-hub's `build.sh` in Wave 3, since kvit-hub has
  no CI.
- Regenerate and commit, which is what corrects the four drifted values.

**Check.** The generated file matches the committed one, and deliberately editing one hex value in
the committed file fails the build.

## Stage 1.6 — the icon set

Three applications draw a small symbol three different ways today, and one of the three is already
the right answer.

**kvit-notes** has `qml/IconButton.qml`, an accessible button that takes a `glyph` string, and its
call sites pass literal Unicode characters — `›`, `‹`, `×`, `+`, `▰`, `▤`, `⌕`, `N`, `#` — drawn in
whatever font on the machine happens to have them. Three files use it. What a character looks like,
and whether it is present at all, is a property of the reader's system rather than of the design.

**kvit-notes-pro** has `qml/works/WorksIcon.qml`, 78 lines, which introduced the Phosphor icon font
— `resources/fonts/Phosphor.ttf`, 477 KB, MIT-licensed with its licence text beside it, aliased
into the resource system as `qrc:/fonts/Phosphor.ttf`. A call site asks for a symbol by a name that
says what it means (`chevron-right`, `send`, `pencil`, `rotate-ccw`), a table maps that to a
codepoint taken from the font's own stylesheet, and the glyph is drawn at the smaller of the
control's two dimensions with `Text.NativeRendering` so it rasterizes for the pixel grid it lands
on. Twelve names are mapped and six files use it. The header records why it replaced hand-drawn
vector paths: those were drawn on a 24-unit grid and scaled by 18/24 with a fixed stroke width,
which left strokes about 1.7 device pixels wide lying across two rows of pixels.

**kvit-hub** has no icon system at all; each view that needs a symbol draws it with `Canvas`,
`Shape` or a stack of rectangles.

The notes-pro arrangement is the one that generalises, so `KvitIcon` is that component with the
font shipped by the library rather than by one application.

**Tasks**

- Move `Phosphor.ttf` and `Phosphor-LICENSE.txt` into `resources/fonts/` here, and ship the font
  through the `Kvit.Ui` module's own resources so that importing the module is enough — a consuming
  application should not have to add a font alias to its `.qrc`, which is what kvit-notes-pro does
  today. Record the MIT notice in `THIRD-PARTY-NOTICES.md`. The licence is compatible with both the
  MPL-2.0 editor and the proprietary agent build, so nothing here constrains what a consumer ships.
- Port `WorksIcon.qml` as `KvitIcon`, keeping the meaning-name interface, the native rendering, the
  default size of `Interface.px(18)` and the colour defaulting to `Theme.textPrimary`.
- Replace the hand-written `switch` with a name-to-codepoint table generated from the Phosphor
  regular stylesheet, so adding a symbol is a data change and the catalogue in the vocabulary skill
  can list what exists. Keep the kvit meaning names as an alias layer over the Phosphor names,
  since they deliberately differ: what this estate calls `chevron-right` is Phosphor's
  `caret-right`.
- Make an unrecognised name visible rather than blank. `WorksIcon` already computes a `recognized`
  flag and then draws nothing; the shared version fails the gallery check and `qmllint` instead, so
  a typo does not ship as an empty box.
- Name a symbol for each of the nine literal characters kvit-notes uses today, so Wave 5 has
  somewhere to send them.
- Port kvit-notes' `IconButton.qml` as `KvitIconButton`, taking an icon name rather than a glyph
  string. It is the accessible version — a real `AbstractButton`, so Qt publishes it with a role and
  a name, it takes tab focus, and Space or Return activates it (`accessibility.md` Finding 1) —
  and it derives both the tooltip and the accessible name from one `label` property so the two
  cannot drift apart. Most of the estate's small controls are a `Rectangle` with a `MouseArea`,
  which reaches no assistive technology at all, so this is the component that replaces them.
- Decide the weight policy. Phosphor ships six weights at roughly 480 KB each; ship regular alone
  unless a second weight has a caller, and record the decision so it is not reopened per
  application.

**Check.** Every named symbol renders in the gallery at the sizes the chrome actually uses — 13 px
for a combo indicator, 18 px at rest — in all four themes; the codepoint table is generated rather
than hand-written; and an unrecognised name fails a check rather than drawing nothing.

## Stage 1.7 — the primitives

The vocabulary from `prd.md` §5.4: 36 types under the `Kvit` prefix that `KvitDialog` and
`KvitShell` already use, in `qml/` and exported through the `Kvit.Ui` module.

**Starting material.** kvit-hub's `parts/` folder is 18 components already written against tokens
— `ActivityStrip`, `AgentWrote`, `AskButton`, `BlockDot`, `Bucket`, `DayStrip`, `EffortBar`,
`Figure`, `HealthPhrase`, `HistoryButton`, `ListRow`, `PendingRing`, `ScopeBar`,
`SectionHeading`, `SidebarItem`, `SizeDot`, `WhoMark`, `WindowStatusBar`. They are copied in and
sorted: the general ones become the shared implementation under their `Kvit` names (`ListRow` →
`KvitRow`, `SectionHeading` → `KvitSectionHeading`, `Figure` → `KvitFigure`, `SidebarItem` →
`KvitSidebarItem`, `WindowStatusBar` → `KvitStatusBar`), and the portfolio-specific ones
(`HealthPhrase`, `Bucket`, `AgentWrote`, `WhoMark`) stay behind as kvit-hub's own layer and are
not copied.

**Build order.** Structure first, because everything else is placed inside it, then content, then
marks, then controls, then feedback. Everything from the chevron on a disclosure to the tick in a
checkbox draws its symbol with `KvitIcon` from stage 1.6, so no component here draws a shape by
hand:

1. **Window and structure** — `KvitWindow`, `KvitSidebar`, `KvitSidebarItem`, `KvitHeader`,
   `KvitBreadcrumb`, `KvitRegion`, `KvitViewHead`, `KvitStatusBar`.
2. **Content** — `KvitSectionHeading`, `KvitRow`, `KvitSlimRow`, `KvitCard`, `KvitPanel`,
   `KvitPane`, `KvitDivider`, `KvitDisclosure`, `KvitEmptyState`.
3. **Marks** — `KvitChip`, `KvitTag`, `KvitBadge`, `KvitSlug`, `KvitPip`, `KvitDot`.
4. **Quantities** — `KvitFigure`: tabular numerals, the unit in muted colour, and an em dash where
   a value was not measured rather than a zero.
5. **Controls** — `KvitButton` in primary, ordinary and quiet forms, `KvitStepper`, `KvitField`,
   `KvitSearchField`, `KvitCheck`, `KvitSelect`, `KvitTab`.
6. **Feedback** — `KvitToast`, `KvitPopover`, `KvitHoverCard`, `KvitDialog`, `KvitTooltip`.

Each component is added to the gallery in the same commit that adds the component, so stage 1.10 is
assembly rather than a backlog.

**Check.** Every primitive renders in the gallery in all four themes and every state it declares,
and `qmllint` is clean.

## Stage 1.8 — the data-display components

`prd.md` §5.5. These are the expensive ones and they are where each application has independently
invented something, so they come with the nine common rules attached rather than as loose parts.

**Marks** — `KvitBar` (kvit-hub's `AxisBar` generalised), `KvitStackedBar` with a two-pixel surface
gap so adjacent fills never touch, `KvitSpark` drawing a baseline tick for a period that was never
measured rather than a zero-height bar, `KvitTrend` with a value axis and a hover crosshair,
`KvitDistribution`, `KvitGauge` (`BudgetGauge` generalised), `KvitDelta`.

**Tiles** — `KvitStatTile`, which is the shape of six of kvit-cash's seven dashboard widgets, and
`KvitFigureBlock`, which `frame.css` already defines.

**The table** — `KvitTable`, `KvitCell` in its text, figure, chip, slug and date variants, and the
empty state from stage 1.7 covering the no-data case for every mark and the table. `KvitTable` is
a `TableView` with `reuseItems` over a C++ `QAbstractItemModel`, a `HorizontalHeaderView` for
column resize and reorder, a `selectionModel` for multi-select, saved column sets, keyboard
navigation and copy of a selection. The library ships the view, a small model base class, and a
synthetic 250,000-row model used only for the benchmark; the real models belong to the
applications.

No charting library is introduced. The estate has zero references to QtCharts or QtGraphs today,
everything is `Rectangle` and `Repeater`, and the conventions that matter here — hatching that
marks a figure as an upper bound and survives grayscale, confidence carried by border style, a
shape as well as a hue on every health level — are not things a charting library exposes.

**Build `KvitTable` first within this stage,** against the 250,000-row synthetic model from the
first commit. It is the one item in the whole plan that can fail on its own terms, and it is the
go/no-go for Decision 3, which put kvit-cash's transaction browser in QML rather than Qt Widgets.

**Check.** `KvitTable` holds smooth scrolling and sub-100 ms filtering at 250,000 rows with twelve
configurable columns, multi-select and inline editing. Every mark renders in four themes. A chart
with no data shows an empty state describing what will appear rather than an axis drawn around
zeros.

## Stage 1.9 — the components named by a written flow

`prd.md` §5.6 separates HuskarUI's 76 components into what this estate has already built by hand,
what it has never built but has already specified in a written flow, and what has no caller here.
The first two groups are built in Wave 1; the third is not built at all.

**Group A — built by hand between three and thirty-two times already.** A working version can be
lifted from whichever application has the best one rather than designed from nothing: scroll bar
(32 files have a private version), context menu (18), collapsible section (18), tree (12), switch
(10), radio group (10), progress (5), slider (4), split view (3, in three different applications).
The collapsible count is also why `KvitDisclosure` needs a multi-section form and not only the
single trigger-and-body that `frame.css` has today.

**Group B — absent everywhere, and named in a flow that has already been written.** Segmented
control (kvit-cash's dashboard needs one period control governing seven widgets), type-ahead field
(its category and tag pickers), confirm-in-place (its bulk edit puts the undo inside the
confirmation), timeline (per-account update history, and notes-pro's session history), a
persistent notification distinct from the transient `KvitToast`, numeric field with a money
variant carrying minor units and the marking rules from `money-display.md`, and a dual-list
selector for configurable columns.

**One that is a judgement rather than a yes or no.** The spotlight mechanism behind a guided tour
is wanted twice — kvit-cash's first-run flow and notes-pro's "guide me" — and what drives it
differs in each. Build the focus-a-region primitive; leave the stepping to whichever application
wants it.

Group B is built here rather than deferred to Wave 4 despite kvit-cash being the first caller for
all seven, because building them against the written flows costs less than kvit-cash stopping to
build a control mid-screen. Where a flow is ambiguous, build the smaller version and let Wave 4
extend it.

**Check.** Each of the sixteen renders in the gallery, and every hand-rolled version it replaces
is listed in the migration notes for its wave.

## Stage 1.10 — the gallery

kvit-hub's `GalleryView.qml` is 247 lines and already renders its vocabulary in every state across
four themes; its header comment says this is where the rule that no distinction rests on colour
alone gets checked once so every view inherits verified components. That idea moves here and
becomes an application of its own.

**Tasks**

- `kvit-ui-gallery`, a standalone binary rendering every primitive in every state in all four
  themes, with the interface size steppable so the density work in stage 1.3 is visible.
- Following HuskarUI's arrangement, the gallery is also the documentation: one page per component
  showing its states and a working code sample.
- A screenshot run producing a fixed, named set of images, so that what gets reviewed after a
  token change is the diff against the previous run rather than the whole set again.
- A single-component harness that loads one primitive against sample data in a chosen theme, which
  is what the preview skill in stage 1.12 wraps.

**Check.** The gallery covers every component built in stages 1.6 to 1.9, and the high-contrast
pass shows no distinction resting on hue alone.

## Stage 1.11 — the drawing layer

The pipeline that produced the screens the owner has already judged moves here so all four
applications draw against one stylesheet.

**Tasks**

- Copy in `tokens.css` (now generated, per stage 1.5), `frame.css` (665 lines, hand-written, and
  it stays hand-written because it is the window and layout vocabulary rather than a mirror of
  anything), `render.sh` (50 lines, headless Chromium at 1440×960), `visual-language.md` and
  `PATTERN.md`.
- Extend `frame.css` with classes for what has no drawing form today: the trend, the distribution,
  the stat tile and the table.
- Each application keeps its own mockups and its own data pack; only the stylesheets, the render
  script and the two prose documents are shared.

**Check.** A mockup copied from kvit-hub renders here through `render.sh` and is
pixel-indistinguishable from its existing render, apart from the four corrected colours.

## Stage 1.12 — the agent skills

Under `agent/skills/`, mirroring HuskarUI's layout, and installed into each application through
the submodule so the four cannot drift on how the system is described to an agent.

**`kvit-ui` — the vocabulary skill.** How to express intent, which is the rule that a request is
written in the vocabulary rather than in measurements: "a slim row with the project name, a health
phrase, and the attention figure right aligned" is buildable, while "a 30-pixel row with #e6b877
text on the right" is not, because it names values instead of meanings and will be wrong in three
of the four themes. It carries a generated catalogue — every component, its properties, its
states, a working snippet — produced from this repository rather than written by hand, since a
hand-written catalogue goes stale the first time a property is added. It carries the token
reference: what each colour means and what it is reserved for, most of which `visual-language.md`
already records. It carries the rules that must not be broken, which exist today as `PATTERN.md`'s
list. And it carries the workflow: draw the screen as HTML against `tokens.css` and `frame.css`,
render it with `render.sh`, put the image in front of the owner, then implement it in QML from
primitives and finish with a screenshot run in four themes.

**`kvit-preview` — the rendering skill.** Wraps `render.sh` for a drawing, the application's own
screenshot run for its views, and the single-component harness from stage 1.10. The reason this is
a skill rather than a paragraph in a README is stated in kvit-hub's own `CLAUDE.md`: the
application's only output is what it draws, so the test suite says nothing about whether a view is
right, and an agent that cannot see its own output is guessing.

**Check.** An agent given a screen description in the vocabulary, and no colour or pixel value,
produces a drawing, renders it, and implements it from primitives.

## The Qt Widgets bridge: contingency, not scheduled work

`prd.md` Decision 3 puts kvit-cash in QML throughout, with the transaction browser on `KvitTable`,
and holds the Widgets fallback open until the 250,000-row benchmark in stage 1.8 passes. So the
`QPalette` bridge and the C++ metrics helper described in §5.9 are **not built in Wave 1**. If the
benchmark fails, they become a stage of their own before Wave 4, and `kvit-cash/prd.md` §48.1
stands as written instead of being updated.

Keeping Widgets out of `find_package` in stage 1.1 is the visible marker of that decision. Adding
it back is a one-line change if the benchmark goes the other way.

## Wave 1 order and dependencies

```
1.1 build ─┬─ 1.2 token classes ──┬─ 1.3 type scale + density ─┬─ 1.6 icons ─┬─ 1.7 primitives ──┬─ 1.10 gallery
           │                      │                            │             │                   │
           │                      ├─ 1.4 chart ramps ──────────┤             ├─ 1.8 data display ┤
           │                      │                            │             │                   │
           │                      └─ 1.5 css generator ────────┘             └─ 1.9 group A + B ──┤
           │                                                                                     │
           └─ 1.11 drawing layer (any time after 1.5) ───────────────────────────────────────────┴─ 1.12 skills
```

The four hard constraints: the type scale is agreed before any component is drawn, since every
component measures itself against it; the icon set exists before the components, since a button, a
disclosure and a checkbox all ask it for a symbol; the ramps exist before the data-display work,
since those components are the only callers; and the skills come last, because the catalogue is
generated from the finished component set.

`KvitTable` is the one item worth starting early even though its stage sits late — it can be
prototyped against the synthetic model any time after stage 1.2, and it is the piece that could
change Decision 3.

## Wave 1 exit criteria

1. `kvit-ui` builds clean on Linux and Windows, static and shared, with `ctest` green.
2. `Theme`, `Interface` and `Typography` are reachable from QML as `Kvit.Ui` singletons, and the
   three ported test files pass unchanged apart from the type-scale changes agreed in stage 1.3.
3. One interface-size setting moves every type role and every density value; the document type
   setting is untouched by it.
4. The categorical, sequential and diverging ramps pass the validator against all four theme
   surfaces, and the validator runs in `ctest`.
5. `tokens.css` is generated, the committed copy matches, and editing it by hand fails the build.
6. The Phosphor font ships with the module, `KvitIcon` resolves every named symbol from a generated
   table, and no component draws a symbol by hand.
7. The gallery renders every component in four themes, and the high-contrast pass shows no
   distinction resting on hue alone.
8. `KvitTable` meets its 250,000-row benchmark, or it does not and the Widgets contingency is
   scheduled.
9. Both agent skills exist, and the vocabulary skill's catalogue is generated from this repository.
10. kvit-notes, kvit-notes-pro and kvit-hub build and behave exactly as they did before Wave 1
    started. Nothing outside this repository has changed.

---

# Waves 2 to 5 — the applications move onto the system

These are described as scope, order and finishing condition. They are not planned in detail here
on purpose: what each of them costs depends on how the components turn out, and Wave 3 in
particular is the measurement that makes Waves 4 and 5 estimable. Each wave gets its own plan
written at the point it starts.

## Wave 2 — kvit-notes hands over the token classes

**Scope.** Delete `theme.*`, `typography.*` and `interfacemetrics.*` from `src/platform/`, add
`kvit-ui` as a pinned git submodule, link `kvit-ui-tokens`, and change the QML import so the
singletons come from `Kvit.Ui` rather than the `Kvit` module. `SettingsStore`, `SystemAppearance`
and `PerfLog` now have two homes; each one is either taken from `kvit-ui` in both places or kept
in kvit-notes and dropped from `kvit-ui`, decided per class at the time rather than now.

**What does not happen.** No screen changes, no primitives adopted, no hex sweep. This wave is
plumbing, and its value is that the duplication window from Wave 1 closes.

**Finishing condition.** kvit-notes' full test suite passes, its 223 screenshots re-render
identically, and `src/platform/` no longer contains the three classes.

**Also in this wave, or immediately after it:** kvit-notes-pro pins the same `kvit-ui` commit
through its own submodule arrangement. Its `core/` submodule already points at kvit-notes and has
not drifted, so this is a second pin beside an existing one, and the practice of denying agent
edits inside a submodule extends to it unchanged.

## Wave 3 — kvit-hub/desktop migrates

**Scope, in the order the work has to happen.** Replace `add_subdirectory(${KVIT_EDITOR_ROOT})`
with the `kvit-ui` submodule for the interface layer, keeping the kvit-notes dependency only for
the graph layout and `DocumentExporter`. Retire the lowercase `theme`, `typography` and
`appSettings` context properties in `app/main.cpp` in favour of the `Kvit.Ui` singletons. Convert
348 `font.pointSize` sites to `font.pixelSize`. Repoint the "Interface size" stepper in
`SettingsPanel.qml` at `Interface.fontSize`, which is what its label already claims and which ends
the settings collision with kvit-notes. Delete `theme/Tokens.qml` once every name it defines has a
home in `Kvit.Ui`. Sweep the 146 bare spacing literals and the 30 literal radii onto the named
scale. Replace the hand-drawn symbols in its views with `KvitIcon`. Repoint the views at the shared
components, retire the `components/` folder — its five live
members move to the portfolio layer, the six orphans and the duplicate `ScopeBar` are deleted —
and add the generator check to `build.sh`, since this repository has no CI.

**Why it is third and not later.** It is the largest existing body of token-disciplined QML, its
`parts/` folder is where most of the shared components came from, and it is the only application
that can tell you whether the components survive contact with real screens before kvit-cash bets
on them.

**Finishing condition.** `build.sh --shots` output matches the pre-migration screenshots except
where a change was intended, `components/` is gone, and no QML file in the tree contains a colour
literal, a numeric font size or an unnamed spacing value.

## Wave 4 — kvit-cash is built on the system

**Scope.** The finance application's first screens are assembled from components, with its money
figures, account rows and category marks as its own layer built on top. `money-display.md`'s three
kinds of figure — native, converted, same-currency — are implemented against `KvitFigure` rather
than drawn again. The transaction browser is `KvitTable` with a kvit-cash model behind it. Its
spending categories are the first caller of the categorical ramp. Where a Group B control from
stage 1.9 turns out to be smaller than the flow needs, kvit-cash extends it in `kvit-ui` rather
than working around it locally.

**Why it is fourth.** It is the test of whether the system is usable by an application that did
not grow up with it, and an empty repository is the cheapest place to find out.

**Finishing condition.** kvit-cash ships a dashboard with no colour literals and no component of
its own that duplicates a shared one, and `kvit-cash/prd.md` §48.1 has been updated to match
Decision 3.

## Wave 5 — the editors adopt the components

**Scope.** kvit-notes replaces its redrawn rows, chips and section headings with the shared
components, moves its `IconButton` call sites off literal Unicode characters and onto named icons,
and clears its 49 hardcoded hex values across 104 QML files. kvit-notes-pro does the
same in `qml/works/` and `qml/agent/` for the screens it touches, drops its own `WorksIcon.qml` and
its copy of the Phosphor font in favour of the shared `KvitIcon`, and clears the two literal font
sizes in its 276 `font.pixelSize` sites.

**Explicitly not in this wave.** Breaking `WorksShell.qml` (2,481 lines) and `WorksSidebar.qml`
(2,042) into smaller files is separate work under `prd.md` Decision 4, planned and scheduled on its
own. Nothing in this plan waits on it, and notes-pro gets the tokens and components regardless.

**Finishing condition.** Zero hardcoded hex in either editor, and every row or chip that the
shared vocabulary covers is drawn by the shared version.

---

# Risks carried into Wave 1

**`KvitTable` is the item that can fail.** Everything else in Wave 1 is a rearrangement of work
already done and carries schedule risk rather than technical risk. A dense, configurable, editable
250,000-row table in QML is the one piece nobody in this estate has built. Stage 1.8 puts it first
within its own stage and against realistic data for that reason, and the Widgets contingency stays
open until it passes.

**The component set is drawn before kvit-cash exists.** A vocabulary decided against three
applications may not fit the fourth. Wave 4 preceding Wave 5 is the mitigation, and the inventory
in `prd.md` §5.4 is a first cut rather than a specification.

**The duplication window in Wave 1.** Two copies of `Theme`, `Typography` and `InterfaceMetrics`
exist from stage 1.2 until Wave 2 finishes, and no automatic check can reconcile them once stage
1.3 reworks the type scale. The freeze on those six files in kvit-notes is the whole mitigation,
which makes Wave 2 worth starting as soon as Wave 1's token stages are done rather than after all
of Wave 1.

**Screenshot review does not scale.** Four applications, four themes and more than sixty
components is a large number of images for one person to judge. The gallery is useful only if what
gets reviewed is the diff against the previous run.

**The token layer becomes a bottleneck.** Every application waits on `kvit-ui` for a value it
needs. The submodule arrangement lets an application pin an older version, and the application's
own layer stays open so it can build what it needs without asking first.
