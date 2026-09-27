# kvit-ui-go

This repository is the Go version of kvit-ui, the component library and
design values that every Kvit desktop app draws with. It is built on the
unison toolkit (`github.com/richardwilkes/unison`). The Qt/QML version it
replaces is in `~/kvit-ui` and is the specification: `PARITY.md` lists
everything that version does and marks what exists here.

This repository is one step of moving all the Kvit desktop apps from Qt to Go.
The plan is `~/kvit-shirei/go-ui-plan.md`; read its sections 3 and 4 before
changing how this repository is laid out or built.

## Recording the migration

The migration is recorded in `~/kvit-shirei/migration-log.md` for a later
blog post. Append a dated entry there, newest last, when you:
- finish a step of the plan;
- make a decision the plan does not cover;
- find something that works differently from what the plan expects;
- measure something.

Give the command behind every number, and save screenshots under
`~/kvit-shirei/migration-log/<date>/`.

## What is where

| Package | What it holds |
|---|---|
| `palette` | Colours and the colour science the design values are checked with: OKLab/OKLCH, perceptual distance, colour-vision simulation, WCAG contrast, chart-ramp rules |
| `tokens` | `Theme` (four colour tables, overrides, reduced motion, following the desktop), `Interface` (everything derived from the 10–24 px interface size), `Typography` (document text). `tables_gen.go` is generated |
| `settings` | The settings file (`ui.json`), in the Qt library's format and keys |
| `platform` | What the desktop says about dark mode, high contrast and reduced motion |
| `icons` | The embedded Phosphor font and its names. `catalog_gen.go` is generated |
| `text` | Font discovery and fallback, shaping, line breaking, caret and hit testing, drawing through unison's canvas |
| root (`kvitui`) | `UI`, which ties the above together and applies the theme to unison; the components will live here |
| `cmd/kvit-ui-gallery` | The gallery and its screenshot mode |
| `tools/import-qt-*` | Generators reading the Qt library's sources; their tests fail when the Qt source has changed since the last run |
| `third_party/typesetting` | go-text v0.3.5 with one patch (see below) |

The two generated files are rebuilt with:

```sh
go run ./tools/import-qt-theme ~/kvit-ui/src/tokens/theme.cpp > tokens/tables_gen.go
go run ./tools/import-qt-icons ~/kvit-ui/src/qml/iconcatalog.cpp > icons/catalog_gen.go
```

**go-text is a patched copy.** `third_party/typesetting` is go-text v0.3.5
with one addition, a font-loader option in its font finder, used through a
`replace` line in `go.mod`. `KVIT-PATCH.md` there says what, why and when
it goes away. Every app built on kvit-ui-go needs the same line:

    replace github.com/go-text/typesetting => ../kvit-ui-go/third_party/typesetting

Without it the app still compiles against go-text's own release, which
lacks `SetFaceLoader`, so the build fails and the missing line is visible.

**Draw text with the `text` package, never with unison's `Text` or
`Label`.** unison's own text draws one glyph per character, without
kerning, emoji sequences or right-to-left ordering.

**The Qt gallery's screenshots are the reference for parity.** They are in
`~/kvit-qt-reference/kvit-ui-0a0b210/`, outside every repository, because
they cannot be made again once Qt is removed. `./build.sh --shots` stacks
each Go screenshot above the Qt one with the same name, in
`build/shots/compare/`.

## Building and checking

```sh
./build.sh               # build every package and build/kvit-ui-gallery
./build.sh --test        # also gofmt check, go vet and the headless tests
./build.sh --cross       # also the gallery for Windows, macOS (both) and Linux
./build.sh --win         # the gallery for Windows onto D:, started on the Windows desktop
./build.sh --win-smoke   # the same, closing after 6 s, with first-frame time and memory
./build.sh --shots       # the screenshot set, compared with the Qt one
tools/check-all.sh       # ./build.sh --test in every Kvit Go repository
```

The git hook in `.githooks/pre-commit` runs `./build.sh --test`
(`git config core.hooksPath .githooks` enables it in a fresh clone).

- **Tests run headless.** They use `unison.StartHeadless`, which runs the
  real event loop against an in-memory screen, with input injection,
  screenshots and screen-reader assertions. Tests never open a window on the
  desktop.
- **Looking at something in a real window** is done on Windows with
  `./build.sh --win`, since this machine is WSL on Windows and Windows is the
  main platform.
- **cgo is always off,** and every build must stay cross-compilable.

## Conventions

- **Module path.** It is `github.com/kvit-s/kvit-ui`, the repository's final
  name. At the switch this repository's history moves into `~/kvit-ui`, so
  import paths never mention `-go`. Apps use it through a `replace` line in
  their `go.mod`, pointing at `../kvit-ui-go`.
- **unison version.** One version is used across all the Kvit Go
  repositories. Change it only with `tools/bump-unison.sh <version>`, which
  moves every repository and then runs `tools/check-all.sh`. unison is a 0.x
  series, so read its release notes first.
- **unison stays unmodified.** Anything unison lacks is written here or in
  the app, never patched into unison. Bugs in unison go upstream as issues
  with a test that reproduces them.
- **History.** Commit on `main` and keep it linear, with no branches and no
  merge commits. The switch replays this history commit by commit.
- **The library's rules,** as in the Qt version:
  - no colour literal, numeric font size or unnamed geometry value outside
    the design values;
  - no distinction resting on hue alone;
  - every interactive component gives screen readers a role and a name.

  Each rule is enforced by a test once the code it governs exists.
