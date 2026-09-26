# kvit-ui-go

The Go version of kvit-ui: the component library and design values the Kvit
desktop apps draw with, built on the unison toolkit. It replaces the Qt/QML
library in `~/kvit-ui`, and takes over that repository's name when it does
everything that one does (`PARITY.md`).

It builds with cgo off and cross-compiles to Windows, macOS and Linux.
Building needs Go 1.27; with an older Go installed, the `toolchain` line in
`go.mod` makes Go fetch 1.27 by itself.

```sh
./build.sh --test        # build, check formatting, vet, run the headless tests
./build.sh --win         # start the gallery on the Windows desktop (from WSL)
./build.sh --help        # every option
```

**Status:** a skeleton. The gallery opens an empty window.

**Licence:** MPL-2.0, as for the Qt version.
