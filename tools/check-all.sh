#!/usr/bin/env bash
# Runs ./build.sh --test in every Kvit Go repository on this machine, so a
# change to kvit-ui-go is checked against every app that uses it before it is
# committed. The repositories checked are the directories ~/kvit-* whose go.mod
# declares a module under github.com/kvit-s/ and requires unison. That finds
# the -go repositories during the migration and the switched ones after it,
# and leaves out Kvit's other Go projects, such as kvit-coder.
#
#   tools/check-all.sh
set -uo pipefail
failed=()
checked=0
for mod in ~/kvit-*/go.mod; do
    grep -q '^module github.com/kvit-s/' "$mod" || continue
    grep -q 'github.com/richardwilkes/unison ' "$mod" || continue
    dir=$(dirname "$mod")
    echo "== $dir"
    checked=$((checked + 1))
    if [ -x "$dir/build.sh" ]; then
        (cd "$dir" && ./build.sh --test) || failed+=("$dir")
    else
        (cd "$dir" && go vet ./... && go test ./...) || failed+=("$dir")
    fi
done
echo "checked $checked repositories"
if [ ${#failed[@]} -gt 0 ]; then
    printf 'failed: %s\n' "${failed[@]}"
    exit 1
fi
