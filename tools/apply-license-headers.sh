#!/bin/bash
# Apply (or verify) the MPL-2.0 Exhibit A header on every first-party source
# file: src/, qml/, gallery/, tests/, tools/. resources/ is never touched — the
# vendored font and its stylesheet carry their own licence.
# Idempotent - files already carrying the header are left alone.
#
#   tools/apply-license-headers.sh          # add missing headers
#   tools/apply-license-headers.sh --check  # exit 1 if any file lacks it
set -e
cd "$(dirname "$0")/.."

CHECK=0
[ "$1" = "--check" ] && CHECK=1

MARKER="subject to the terms of the Mozilla Public"

slash_header() {
    cat <<'EOF'
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
EOF
}

hash_header() {
    cat <<'EOF'
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
EOF
}

MISSING=0
process() {
    local f="$1" style="$2"
    if grep -q "$MARKER" "$f"; then
        return
    fi
    if [ "$CHECK" -eq 1 ]; then
        echo "missing header: $f"
        MISSING=1
        return
    fi
    local tmp mode
    mode=$(stat -c%a "$f")
    tmp=$(mktemp)
    if [ "$style" = hash ] && head -1 "$f" | grep -q '^#!'; then
        # Keep the shebang on line 1.
        { head -1 "$f"; hash_header; tail -n +2 "$f"; } > "$tmp"
    elif [ "$style" = hash ]; then
        { hash_header; cat "$f"; } > "$tmp"
    else
        { slash_header; cat "$f"; } > "$tmp"
    fi
    mv "$tmp" "$f"
    chmod "$mode" "$f"
    echo "added: $f"
}

while IFS= read -r -d '' f; do
    case "$f" in
        *.cpp|*.h|*.hpp|*.qml|*.js) process "$f" slash ;;
        *.py|*.sh) process "$f" hash ;;
    esac
done < <(find src qml gallery tests tools -type f \
             \( -name '*.cpp' -o -name '*.h' -o -name '*.hpp' \
                -o -name '*.qml' -o -name '*.js' \
                -o -name '*.py' -o -name '*.sh' \) -print0)

if [ "$CHECK" -eq 1 ]; then
    [ "$MISSING" -eq 0 ] && echo "all first-party sources carry the header"
    exit "$MISSING"
fi
