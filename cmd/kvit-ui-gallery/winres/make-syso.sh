#!/usr/bin/env bash
# Writes cmd/kvit-ui-gallery/rsrc_windows_amd64.syso: the Windows resources
# of kvit-ui-gallery.exe, which are its icon (winres/kvit-ui-gallery.ico,
# which tools/appicon draws from icons/kvit-ui-gallery.svg) and its name as
# Explorer, the taskbar and Task Manager show it ("Kvit UI Gallery", under
# Properties > Details). kvit-ui has no release version, so the version
# shown there is 0.0.0.0.
#
#   cmd/kvit-ui-gallery/winres/make-syso.sh
#
# The go command links a .syso file in a package's folder into every
# windows/amd64 build of that package, so ./build.sh --win, --cross and `go
# build` get the icon too.
#
# The tool is go-winres (github.com/tc-hib/go-winres), run with `go run` at a
# pinned version, as kvit-notes and kvit-works run it, so it needs nothing
# installed and leaves go.mod alone.
set -euo pipefail
cd "$(dirname "$0")/.."

GO_WINRES=github.com/tc-hib/go-winres@v0.3.3

go run "$GO_WINRES" make \
    --in winres/winres.json \
    --out rsrc \
    --arch amd64
echo "wrote cmd/kvit-ui-gallery/rsrc_windows_amd64.syso"
