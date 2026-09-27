# How to draw a view

You are drawing one static HTML page that will be rendered to a 1440×960 PNG and judged as
an image by the owner. This file is the contract: the skeleton to copy, the classes that
already exist, and the things that will make the drawing wrong. Read it in full before you
write anything; it is short on purpose.

Four other files in this directory matter to you:

| File | What it is |
|---|---|
| `tokens.css` | the running application's colours, type scale, density and existing components |
| `frame.css` | the window, the global chrome, the caption strip, and the pieces tokens.css lacks |
| `../views.md` | what your view contains, in reading order, and which data state it is in |
| `../concept-model.md` | the only words allowed on the drawing |

`visual-language.md` in this directory is the full design note. You do not need to re-read
it if you follow this file, but its "Rules a drawing agent must not break" section is the
authority if anything here is ambiguous.

`data-pack.json` (400 KB) holds real content from the data repository. Inspect it with
`python3 -c` or `jq`; never read it whole into context.

---

## The skeleton

Copy this exactly. Everything between `<div class="stage">` and its closing tag is yours.

```html
<!doctype html>
<meta charset="utf-8">
<title>view-id — one clause saying what it shows</title>
<link rel="stylesheet" href="tokens.css">
<link rel="stylesheet" href="frame.css">
<style>
  /* Only what this view needs. Layout, column widths, one-off spacing.
     Never a colour literal — use the tokens. */
</style>

<body>
<div class="app">

  <div class="chrome">
    <span class="back is-disabled">‹</span>          <!-- drop is-disabled on entered views -->
    <span class="wordmark">kvit-hub</span>
    <span class="places">
      <span class="tab is-current">Today</span>
      <span class="tab">Projects</span>
      <span class="tab">Decisions <span class="badge">5</span></span>
    </span>
    <span class="tools">
      <span class="btn btn--quiet">?</span>
      <span class="btn btn--quiet">↻</span>
      <span class="btn btn--quiet">⚙</span>
    </span>
  </div>

  <!-- Entered views only. Today, Projects and Decisions omit this element. -->
  <div class="crumbs">
    <span class="link">Projects</span><span class="sep">›</span>
    <span class="link">kvit-editor</span><span class="sep">›</span>
    <span class="here">public open-source launch</span>
  </div>

  <div class="stage">
    …your view…
  </div>
</div>

<div class="caption">
  <span class="id">view-id</span>
  <span class="state">Real.</span>
  <span class="said">One sentence naming what the figures and strings came from.</span>
</div>
</body>
```

Three rules on the skeleton itself.

**`is-current` moves to whichever of the three places your view lives under.** Today,
Projects and Decisions are the only top-level places; a project page, a milestone page, the
work view, the evidence trail and the document surface are all *entered*, so they carry
`.crumbs` and an enabled back control. A top-level view omits `.crumbs` entirely and leaves
`is-disabled` on back, because there is nothing behind it.

**The badge on Decisions is 5.** That is the number of rows a decision settles; the fifteen
informational rows are deliberately not in it. If your view has a reason to show the other
number, show it inside the view; the badge stays at 5.

**The caption's `state` is one of exactly three words**, matching what `views.md` says about
your view: `Real.`, `Real with specimen prose.` or `Constructed.` The `id` is your filename
without `.html`, and it must be the id `views.md` uses. The `said` sentence names the source
in one line: which day, which record, which tool.

---

## The vertical budget

The page is exactly the screenshot and it never scrolls. Work out whether your content fits
*before* you write it, because a region that overflows clips silently: the renderer passes
`--hide-scrollbars`, so nothing appears to tell you it happened.

```
960  page
-34  caption
-52  chrome
-34  crumbs        (entered views only)
────
874  a top-level view's stage
840  an entered view's stage
```

Horizontally the stage is 1408 px wide inside its 16 px margins.

Row heights, and how many fit:

| Class | Height | Fits in 874 | What it is |
|---|---|---|---|
| `.row` | 56 px | 15 | tokens.css's ledger row: name over description, chips and figures beside |
| `.row--sub` | 48 px | n/a | an expanded line under a row |
| `.row--slim` | 30 px | 26 | frame.css's one-line row, for a list at rest |
| `.row--compact` | 24 px | n/a | a disclosure or triage line |

`projects-base` uses `.row--slim` because 26 projects at the 56 px height need 1456 px and
it has 833. The drawing makes that a deliberate position: at rest the row carries one line,
and a layer that adds content takes it back to the full height. If your view's rows carry
two lines of content, use `.row` and accept that fewer of them fit rather than shrinking the
row to squeeze more in.

Anything you cannot fit either goes in a `.region` (which owns the leftover height and
scrolls) or is cut. Cutting and saying so in the caption is better than clipping.

---

## What already exists

### From `tokens.css`

Do not redefine these; they are the application's own components.

- `.row` `.row--sub` `.row--compact` `.row--dim` `.row-title` `.row-sub` `.disclose`: rows.
  `.is-hover` and `.is-focus` paint the two states a static page cannot produce.
- `.figure` `.figure--attention` `.figure--agent` `.figure--large` `.figure .unit`: a
  measured quantity. Tabular numerals; whole hours at 10 h and above, one decimal below.
- `.dash`: not measured. An empty `<span class="dash"></span>` renders the em dash on its own.
- `.bar` `.bar--wide` `.bar--agent` with `.spent` and `.remaining` children, driven by
  `--spent` and `--remaining` as percentages of a scale shared down the column.
  `.axis-dot` / `.axis-dot--agent` label which axis a bar belongs to.
- `.scope` with `.closed` and `.discovered`: the item-scope bar that a percentage never
  travels without.
- `.chip` `.chip--closed` `.chip--active` (`--fill`) `.chip--planned` `.chip--unestimated`,
  `.sev-hard`, `.sev-soft` or `.sev-hygiene` on the border, and `.pip` for ready-to-close.
- `.tag` `.tag--hard` `.tag--soft` `.tag--hygiene`: a kind tag. `.count` and `.count--zero`: a
  severity count. `.confidence` `.confidence--bounded` `.confidence--unmeasured`.
  `.bucket` `.bucket--agent` `.bucket--none`: a frozen estimate.
- `.slug`: an identifier, monospace, never leading a row a person reads.
- `.btn` `.btn--primary` `.btn--quiet`, `.tab` `.tab .badge`, `.card` `.panel` `.list`.
- `.stack` `.hstack` `.grow` `.nowrap`, `.muted` `.faint` `.secondary` `.label`.

### From `frame.css`

- **window**: `.app`, `.chrome` (`.back` `.wordmark` `.places` `.tools`), `.crumbs`,
  `.stage`, `.caption`. `.stage--flush` drops the margins; `.stage--split` puts a `.pane`
  down the right with `.stage-main` beside it.
- **`.region`** takes the leftover height and scrolls its own overflow. `.region--fixed`
  for a block that should keep its natural height.
- **`.view-head`** is the 40 px strip at the top of a view: `.counts` for what it is counting
  (wrap each number in `<strong>`), then the view's controls.
- **`.heading`** is a section heading with an optional `.count-of` and a hairline reaching
  the right edge. `.heading--strong` for the rule under a column-header row.
- **`.row--slim`** with `.row-name`, `.row-what` and `.row-when`: the one-line row. Lay the
  slots out with your own grid so the columns line up.
- **`.phrase`** carries health and absence, in words. Five levels, each with a shape as well as a
  hue, so none of them rests on colour:

  | Class | Marker | Means |
  |---|---|---|
  | `.phrase--act` | tall solid bar, bold text, red | he can settle this today |
  | `.phrase--decide` | short solid bar, amber | a decision is wanted |
  | `.phrase--tidy` | small square, violet | record cleanup |
  | `.phrase--state` | none, capitalised `.lead` | the object is absent or only partly there |
  | `.phrase` | none | everything else, including ordinary sequencing |

  A `<span class="detail">` inside any of them carries the qualifying clause. If you print
  amber text of your own, use `--signal-soft-text` rather than `--signal-soft`: the latter is
  a bar hue and reads at about 2.5:1 on the light and sepia grounds. `--signal-hard` and
  `--signal-hygiene` are legible as text on every ground and need no partner.
- **`.figure-block`** is a figure with its label above it. `.axis-line` with `.axis-name` puts
  one axis's name, dot and bar on one line.
- **`.spark`** is an inline activity strip. One `<i>` per day with `style="height:N%"`; a day
  that was never measured takes `.is-none` and draws a baseline tick instead of a bar, so
  nothing and zero do not look the same. `.spark--agent` for the other axis; never both in
  one strip.
- **`.disclosure`**: `.trigger` (with a `.disclose` triangle) over a `.body`, with `.source`
  for provenance. Add `.is-open` when you are drawing it open.
- **`.toast`**: `.said` (the action in plain words), `.reason.is-focused` (an empty field
  already focused, with a `.caret`), `.moves`. Position it inside a `position: relative`
  parent.
- **`.pane`**: the side pane, built from `.pane-head` (with `.scope-of`), `.pane-body`, `.pane-foot`
  containing `.ask`. Always `--pane-width`, wherever it appears.
- **helpers**: `.right` `.tabular` `.baseline` `.pad` `.pad-y` `.divider` `.divider--strong`.

---

## Adding a view without touching `frame.css`

Write your layout in the page's own `<style>` block, scoped under a class you put on the
container. `projects-base` does exactly this:

```html
<style>
  .plist .row {
    display: grid;
    grid-template-columns: 160px minmax(0, 1fr) 520px 88px;
    align-items: baseline;
    column-gap: var(--column-gap);
  }
</style>
…
<div class="region plist"> …rows… </div>
```

Change `frame.css` only if you need a piece no view has needed yet, and say so when you
report back, because every other drawing already loaded it and a change to a shared class
silently restyles twenty-seven pages.

---

## Rendering, and checking

```bash
./render.sh projects-base.html      # → renders/projects-base.png, 1440×960
./render.sh                         # every *.html in this directory
```

Then **open the PNG and look at it**. The drawing's only output is what it draws; nothing
else tells you whether it is right, and the first render is never the last one.

Check one other theme before you finish. Copy the file with a theme class on `<body>`,
render, look, delete the copy:

```bash
sed 's|<body>|<body class="theme-contrast">|' myview.html > _t.html
./render.sh _t.html && rm _t.html renders/_t.png
```

`theme-contrast` is the accessibility pass, on a pure black ground with white text, and it
is what catches a distinction that turned out to rest on colour. `theme-light` catches a hue that is
illegible on white. Leave no `_`-prefixed files behind: `./render.sh` with no arguments
renders everything in the directory.

---

## Rules that must not be broken

1. **No network requests of any kind.** No web fonts, no CDN, no remote images, no `fetch`.
   Everything inline or local.
2. **No literal colours.** Every colour comes from a token variable. A hex is wrong in three
   of the four themes and silently breaks the high-contrast check.
3. **No distinction resting on colour alone.** Give it a shape, a border style or a word as
   well, then render `theme-contrast` and confirm it survives.
4. **Red only for what can be acted on today.** A milestone whose tasks are all finished and
   whose definition of done is unconfirmed is red; a stalled blocker is red. Ordinary
   sequencing is muted ink, always. `--danger` is for the application failing and is never a
   data state.
5. **Not measured is a dash, never a zero and never an empty cell.** Zero, not-measured and
   unknown are three different states. An unknown figure draws no ink inside its track; the
   track still renders.
6. **The two effort quantities never stack, sum or share a scale.** They are *my time* and
   *agent time*, on separate rows with separate scales, my time first where both appear.
7. **No dates on forward-looking figures.** Everything forward is effort ("8 h to go",
   "2–3.5 working days"). Only historical facts carry dates.
8. **No project-level percentage complete.** A milestone has a denominator; a project does
   not. Project effort is absolute measured hours.
9. **Titles for people, identifiers on demand.** Lead a row with the title; the slug follows
   in `.slug`. A cross-project reference is project + milestone chip + title, never
   `project:milestone-slug`.
10. **No abbreviation the drawing never explains.** It is judged as an image, so a tooltip is
    not an escape hatch. Anything essential is visible.
11. **Match the density.** Margins 16, column gaps 14, rows at their table height above. A
    drawing airier than the application reads as a proposal to make it airier.

---

## The words

Every string on the drawing comes from the approved vocabulary in `../concept-model.md`.
The ones that catch people out:

| Say | Never say |
|---|---|
| Projects (the screen) | the ledger |
| task | item |
| acceptance criteria | `done_when` |
| definition of done | gate |
| decision, settled, decide | ruling, verdict, delta, proposal sheet, review form |
| evidence, evidence missing | anchor, dangling anchor |
| my time, agent time | attention, agent-wall |
| estimate, remaining estimate | bucket XS–XL, midpoint, remaining effort |
| unplanned | discovered, added later |
| fixing earlier work | rework link, `fixes:` |
| unconfirmed | marked guess, `est_by: agent` |
| finished work not yet confirmed | gate-audit |
| needs planning · no milestone | planning-wanted · task-without-milestone |
| recent activity | momentum |
| inbox item | unruled capture |
| refresh from the repository | re-grounding |
| tracked / partly tracked / not tracked | full / partial / none |
| "direction not reviewed since 2026-06-02" | stale, needs review |
| blocked, blocked by | waiting on |

Never on any surface at all: locus, capped focus, engaged time, block, day window, the 05:00
boundary, fact table, `stats.jsonl`, derived, writer, corrector, applier, channel, registry,
alias, path window, projection, consolidation, backfill.

---

## Glyphs

Only these are installed on the rendering machine. Anything else draws a tofu box, and you
will not notice unless you look at the PNG.

Available: `‹ › ← → ▸ ▾ ▴ ◂ ▪ ▫ ■ □ ● ○ ↻ ⟳ ⚙ ✓ ✕ ⚠ ⚑ ↗ ≤ — · ⋯ ✱`

Missing: `⛭ ⭯ ⮜ ⟨`, and every emoji with a variation selector (`⚙️` renders as a colour
emoji rather than a glyph in the interface's ink). Inline SVG is fine if you need a shape that is not here.
