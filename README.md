# kvit-ui

kvit-ui: the component library and design values the Kvit desktop apps draw
with, built on the unison toolkit.

It builds with cgo off and cross-compiles to Windows, macOS and Linux.
Building needs Go 1.27; with an older Go installed, the `toolchain` line in
`go.mod` makes Go fetch 1.27 by itself.

```sh
./build.sh --test        # build, check formatting, vet, run the headless tests
./build.sh --win         # start the gallery on the Windows desktop (from WSL)
./build.sh --help        # every option
```

**What it holds:** four themes (light, dark, sepia and high contrast), an
interface size from 10 to 24 px that every size and spacing follows, and
document typography; settings shared across the apps; the Phosphor icons,
named by what they mean; a text package that shapes text (kerning, emoji
sequences, other scripts); the desktop's tray icon, notifications and file
dialog; and the components, from labels and buttons to tables, charts, menus,
trees and split views. Kvit Notes, Kvit Works, Kvit Cash, Kvit Hub and the
kvit-term terminal are built on it.

The gallery (`cmd/kvit-ui-gallery`) shows every component with its states in
each theme and interface size, and `agent/skills/kvit-ui/catalog.md`
describes each one with a working code sample.

**Licence:** MPL-2.0.
