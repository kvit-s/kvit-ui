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

## Five shapes that fail without saying so

All five were hit while building this, and none of them produces an error
where the mistake is.

**A container sized from its children while a child fills the container.**
`implicitHeight: holder.childrenRect.height` with `holder { anchors.fill:
parent }` is a cycle. Anchor the holder's width to the parent and take its
height from `childrenRect`, and nothing loops.

**A non-URL passed to `Qt.createQmlObject`.** The third argument has to be a
real URL. A plain string is not rejected; it hangs inside the creation.

**A `contentItem` that positions itself.** A `Control` places its content item
at `(leftPadding, topPadding)` and overwrites any `x` or `y` it was given, so a
content item that indents itself past an indicator lands underneath it instead.
Put the offset in the control's `leftPadding` and leave the content item alone.
`KvitCheck` and `KvitSwitch` reach the same result by giving their label its own
`leftPadding`, which works because a `KvitLabel` has no ground; the control's
padding is the version to copy, because it is what makes `availableWidth` the
real content width.

**A QML function whose name a base type already uses.** `Window` has a
`show()` slot taking no arguments, so a `window.show(name)` written for a
navigation function that was never declared compiles, runs, re-shows the
window and drops the name, saying only `Too many arguments, ignoring 1` on the
console. The gallery's component list did nothing when clicked for as long as
that line stood. Name a function for what it does to the content —
`showPage`, not `show` — and a test that calls it by name from C++ is what
holds it.

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

## Who consumes this, and what that costs

kvit-notes takes `Theme`, `Typography`, `InterfaceMetrics`, `SettingsStore`,
`SystemAppearance` and `PerfLog` from here, as a pinned submodule at
`third_party/kvit-ui` linked as `kvit-ui-tokens`. It no longer has copies of
its own; the duplication window that Wave 2 existed to close is closed, and
the freeze on those files in kvit-notes is over. kvit-notes-pro pins the
same commit beside its `core/` one and points `KVIT_UI_ROOT` at it before
adding the editor, so the two cannot drift by nested-submodule accident.

Two things follow for work done here.

**A token change reaches an application only when somebody bumps its pin.**
There is no automatic propagation and no check that the four are on the same
commit. An application can sit on an older kvit-ui indefinitely, which is the
arrangement working as intended rather than a problem to fix.

**kvit-notes imports `Kvit`, not `Kvit.Ui`.** The singletons keep their old
QML names on its own module, backed by the C++ objects from here. It cannot
import both, because `KvitDialog` exists in both and every dialog would become
ambiguous. Wave 5 is what moves it onto the library's components; until then,
renaming a singleton or changing a property name here breaks it silently at
the QML layer, where nothing in this repository's suite will see it.
