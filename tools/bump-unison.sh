#!/usr/bin/env bash
# Moves every Kvit Go repository on this machine to one unison version, then
# runs tools/check-all.sh. All the repositories use the same version, so
# upgrade them together with this script rather than one at a time. unison is
# a 0.x series whose releases can change its API: read its release notes
# first.
#
#   tools/bump-unison.sh v0.109.0
set -euo pipefail
version=${1:?usage: bump-unison.sh VERSION}
for mod in ~/kvit-*/go.mod; do
    grep -q '^module github.com/kvit-s/' "$mod" || continue
    dir=$(dirname "$mod")
    grep -q 'github.com/richardwilkes/unison ' "$mod" || continue
    echo "== $dir"
    (cd "$dir" && go get "github.com/richardwilkes/unison@$version" && go mod tidy)
done
exec "$(dirname "$0")/check-all.sh"
