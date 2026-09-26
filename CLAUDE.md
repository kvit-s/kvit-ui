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

## Building and checking

```sh
./build.sh               # build every package and build/kvit-ui-gallery
./build.sh --test        # also gofmt check, go vet and the headless tests
./build.sh --cross       # also the gallery for Windows, macOS (both) and Linux
./build.sh --win         # the gallery for Windows onto D:, started on the Windows desktop
./build.sh --win-smoke   # the same, closing after 3 s, with first-frame time and memory
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
