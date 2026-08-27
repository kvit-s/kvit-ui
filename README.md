# kvit-ui

The shared interface layer for the kvit desktop applications: the note editor
kvit-notes, the agent-facing editor kvit-notes-pro, the portfolio dashboard
kvit-hub, and the personal-finance application kvit-cash. Each consumes it as a
pinned git submodule.

Three layers.

**Tokens.** Every value the four applications draw with, as named C++
properties: `Theme` for colour across four themes (light, dark, sepia, high
contrast), `InterfaceMetrics` for the chrome's type scale and its geometry, and
`Typography` for document text in the applications that have documents. One
interface-size setting moves every type role and every density value together,
which is what stops a row being too tight for the text in it at any size.

**Components.** Sixty-eight QML types under the `Kvit` prefix, built from the
tokens: the window and its structure, rows and cards, marks and figures,
controls, feedback, the data-display vocabulary (bars, sparks, trends,
distributions, gauges, tiles) and a table that holds 250,000 rows.

**A drawing layer.** A generated stylesheet, a hand-written frame, and a
headless renderer, so a screen can be proposed as an image before any QML is
written, in the colours the application will actually use, because the
stylesheet is generated from the same token table.

Everything is reached through one import:

```qml
import Kvit.Ui
```

## Building

Qt 6.5 or newer, CMake 3.21, a C++20 compiler.

```
./build.sh --test          # build and run the suite
./build.sh --gallery       # build and open the gallery
./build.sh --help
```

Or through the presets: `cmake --preset linux-release`, `ctest --preset linux`.

## The gallery

```
kvit-ui-gallery                                  # browse
kvit-ui-gallery --page KvitSlimRow --theme dark  # one component
kvit-ui-gallery --shots <directory>              # 68 components x 4 themes
```

It is the reference, with one page per component showing its states and a
working code sample, and it is the review surface. The samples are the same strings
it renders and the test suite compiles, so a sample that does not work stops
the build.

## What holds it together

`ctest` runs thirteen checks. The ones worth naming:

- **The chart ramps are checked, not asserted.** All three ramps in all four
  themes, against that theme's surfaces and against every hue that already
  means something, under normal vision and the three common colour-vision
  deficiencies. Editing a ramp fails the build rather than the eye.
- **Every component loads in four themes with no warning.** A binding to a
  property that does not exist evaluates to `undefined` rather than throwing,
  which a colour renders as transparent; this is what catches it.
- **`KvitTable` meets its benchmark.** 250,000 rows, twelve columns, smooth
  scrolling and sub-100 ms filtering. It is what prd.md's Decision 3 turns on.
- **The generated files match their sources.** The drawing stylesheet, the icon
  catalogue and the agent skill's component catalogue.
- **qmllint is clean at full strength**, `unqualified` and `missing-property`
  included.

## Agent skills

`agent/skills/` holds two, installed into each application through the
submodule so the four cannot drift on how the system is described.

`kvit-ui` is the vocabulary: how to describe a screen, which component to
reach for, what each colour is reserved for, and the workflow from a
description to working QML. `kvit-preview` is how to look at what was built.

## Licence

MPL-2.0 (`LICENSE`). One third-party asset, the Phosphor icon font, under MIT;
see `THIRD-PARTY-NOTICES.md`.
