# Working in kvit-ui

This repository is the shared interface layer for four Qt desktop
applications: the note editor kvit-notes, the agent-facing editor
kvit-notes-pro, the portfolio dashboard kvit-hub, and the personal-finance
application kvit-cash. Each consumes it as a pinned git submodule.

`prd.md` says what the system is and why. `plan.md` says in what order it gets
built. `agent/skills/kvit-ui/` is the skill an application's agent loads; if
you are changing how the vocabulary is described, that is the file to change,
and `catalog.md` beside it is generated rather than edited.

## What must not appear in this tree

Each of these is held by a test or the lint gate. They are listed here because
the failure they prevent is one nobody notices by looking at their own screen.

**A colour literal in QML.** Every colour comes from `Theme`. A literal is
correct in the theme it was picked in and wrong in the other three.

**A numeric font size.** Every size comes from `Interface`: `caption`,
`small`, `body`, `strong`, `title`, `headline`, `display`.

**An unnamed spacing or geometry value.** `Interface` names the spacing scale,
the row heights, the control heights, the radii and the reflow breakpoints.
`Interface.px(n)` is for a one-off; a `px()` with a literal that a second view
will also need is a token waiting to be added.

**A distinction resting on hue alone.** Every state carries a second channel:
a shape, an underline, a weight, a hatch, a position. The gallery's
high-contrast pass is where this is checked.

**A zero standing in for a value nobody measured.** `KvitFigure` draws an em
dash, `KvitBar` a tick at the origin, a chart with no data a `KvitEmptyState`.

**A `Rectangle` with a `MouseArea` used as a control.** It reaches no
assistive technology at all. Use the components; anything whose whole label is
a symbol also takes a `label` in words, which fills both the tooltip and the
accessible name.

**A hand-drawn symbol.** `KvitIcon`, by meaning name. An unknown name draws a
marked placeholder and fails `tests/test_components`.

## Three shapes that fail without saying so

All three were hit while building this, and none of them produces an error
where the mistake is.

**A container sized from its children while a child fills the container.**
`implicitHeight: holder.childrenRect.height` with `holder { anchors.fill:
parent }` is a cycle. Anchor the holder's width to the parent and take its
height from `childrenRect`, and nothing loops.

**A non-URL passed to `Qt.createQmlObject`.** The third argument has to be a
real URL. A plain string is not rejected; it hangs inside the creation.

**An empty `font.family`.** It does not mean "the platform default". Qt
matches it against nothing and falls back to whichever installed face its font
matching lands on, which on a Linux desktop with the usual DejaVu set is
`DejaVu Math TeX Gyre`, a serif maths face. Bind to
`Interface.resolvedFontFamily` and `Interface.resolvedMonoFamily`, never to
the stored preferences `fontFamily` and `monoFamily`, which are empty when the
reader has not chosen one. This one shipped: every label in the library was
drawn in a serif until somebody looked at a screenshot, because a serif
interface looks like a decision rather than a defect.

## Generated files

Four things are generated and checked by `ctest`. Editing one by hand fails
the build; regenerating is always the fix.

| File | Written by |
|---|---|
| `ux/tokens.css` | `cmake --build <dir> --target tokens-css` |
| `src/qml/iconcatalog.cpp` | `python3 tools/generate-icon-catalog.py` |
| `agent/skills/kvit-ui/catalog.md` | `kvit-ui-gallery --catalog <path>` |
| the chart ramps in `theme.cpp` | `python3 tools/design-ramps.py`, by hand |

The ramps are the odd one out: the search prints values and somebody pastes
them in, because it takes a minute and its answer only changes when a rule
does. `tests/test_palette` is what holds the pasted values honest.

## Building and looking

```
./build.sh --test              # build and run the suite
./build.sh --gallery           # build and open the gallery
./build.sh --shots-only        # write the screenshot set
KVIT_UI_BUILD_DIR=build tools/run-qmllint.sh
```

The gallery is the reference and the review surface. `--page <Component>`
opens one, `--theme` and `--interface-size` set the conditions, `--shots`
writes the whole set. What gets reviewed after a token change is the diff
against the previous run rather than the whole set again.

## Adding a component

1. Write it in `qml/`, and add it to the right group in `qml/CMakeLists.txt`.
2. Add an entry to `gallery/Catalog.qml`: what it is for, and at least one
   working sample per state worth seeing. The samples are the same strings the
   gallery renders and `tests/test_gallery` compiles, so a broken one stops the
   build.
3. Regenerate the skill catalogue.
4. Run the suite and look at the four screenshots.

`tests/test_gallery` fails on a component with no page, so step 2 is not
optional.

## The duplication window

`Theme`, `Typography` and `InterfaceMetrics` exist in this repository *and* in
kvit-notes' `src/platform/` until Wave 2 moves that application onto this one.
The six source pairs in kvit-notes are frozen for the duration: a change made
there has to be re-made here by hand, the two have already diverged on the
type scale, and no automatic check can reconcile them. If kvit-notes needs a
token added meanwhile, add it here and cherry-pick the same hunk there.
