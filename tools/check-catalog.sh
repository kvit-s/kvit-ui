#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Fail if agent/skills/kvit-ui/catalog.md no longer matches the repository.
#
#   tools/check-catalog.sh <gallery binary> <catalog path>
#
# The catalogue is what an agent is told the vocabulary contains, and it is
# generated so that it cannot drift from what the vocabulary is. A stale one
# claims a component that was renamed still exists, or omits one that was
# added — and in both cases the agent finds out by writing QML that does not
# compile.
set -e

GALLERY="$1"
CATALOG="$2"
if [ -z "$GALLERY" ] || [ -z "$CATALOG" ]; then
    echo "usage: check-catalog.sh <gallery binary> <catalog path>" >&2
    exit 2
fi

CURRENT=$(mktemp)
trap 'rm -f "$CURRENT"' EXIT

"$GALLERY" --catalog "$CURRENT" > /dev/null

if diff -q "$CATALOG" "$CURRENT" > /dev/null; then
    echo "$CATALOG matches the repository"
    exit 0
fi

echo "$CATALOG is stale. Regenerate it:" >&2
echo "  $GALLERY --catalog $CATALOG" >&2
echo "" >&2
diff -u "$CATALOG" "$CURRENT" | head -40 >&2
exit 1
