#!/usr/bin/env bash
# Builds kvit-ui-go: every package, and the gallery into build/.
#
#   ./build.sh               build
#   ./build.sh --test        also check formatting, run go vet and the headless tests
#   ./build.sh --cross       also build the gallery for windows/amd64, darwin/arm64,
#                            darwin/amd64 and linux/amd64 into build/<os>-<arch>/
#   ./build.sh --win         build the gallery for Windows onto D: and start it there
#   ./build.sh --win-smoke   the same, but it closes itself after 6 s; prints when it
#                            drew its first frame and how much memory Windows gave it
#   ./build.sh --run         start the gallery here (needs a display)
#
# Everything builds with cgo off. KVIT_WIN_DIR overrides where Windows builds go
# (default /mnt/d/projects/kvit-ui-go).
set -euo pipefail
cd "$(dirname "$0")"
export CGO_ENABLED=0
test=0 cross=0 win=0 smoke=0 run=0
for a in "$@"; do
    case $a in
        --test) test=1 ;;
        --cross) cross=1 ;;
        --win) win=1 ;;
        --win-smoke) win=1 smoke=1 ;;
        --run) run=1 ;;
        -h|--help) sed -n '2,15p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *) echo "unknown option: $a" >&2; exit 2 ;;
    esac
done

mkdir -p build
go build ./...
go build -o build/kvit-ui-gallery ./cmd/kvit-ui-gallery

if [ $test = 1 ]; then
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        echo "not formatted with gofmt:" >&2
        echo "$unformatted" >&2
        exit 1
    fi
    go vet ./...
    go test ./...
fi

if [ $cross = 1 ]; then
    for target in windows/amd64 darwin/arm64 darwin/amd64 linux/amd64; do
        os=${target%/*} arch=${target#*/} ext=
        [ "$os" = windows ] && ext=.exe
        GOOS=$os GOARCH=$arch go build -o "build/$os-$arch/kvit-ui-gallery$ext" ./cmd/kvit-ui-gallery
    done
fi

if [ $win = 1 ]; then
    dest=${KVIT_WIN_DIR:-/mnt/d/projects/kvit-ui-go}
    mkdir -p "$dest"
    GOOS=windows GOARCH=amd64 go build -o "$dest/kvit-ui-gallery.exe" ./cmd/kvit-ui-gallery
    if [ $smoke = 0 ]; then
        "$dest/kvit-ui-gallery.exe" &
        disown
    else
        "$dest/kvit-ui-gallery.exe" --smoke 6s &
        pid=$!
        sleep 3
        # The working set includes pages shared with system libraries (the
        # OpenGL driver among them); the private working set is the column
        # Task Manager shows as memory.
        powershell.exe -NoProfile -Command '
            $p = Get-Process kvit-ui-gallery
            $c = (Get-Counter "\Process(kvit-ui-gallery)\Working Set - Private").CounterSamples[0].CookedValue
            "Windows memory: working set {0:N0} MB, private working set {1:N0} MB, private bytes {2:N0} MB" -f ($p.WorkingSet64 / 1MB), ($c / 1MB), ($p.PrivateMemorySize64 / 1MB)' | tr -d '\r'
        wait $pid
    fi
fi

if [ $run = 1 ]; then
    build/kvit-ui-gallery
fi
