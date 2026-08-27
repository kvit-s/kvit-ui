#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Static QML lint gate, over the component vocabulary and the gallery.
#
#   tools/run-qmllint.sh              # lint qml/ and gallery/
#   tools/run-qmllint.sh a.qml b.qml  # lint specific files
#
# What this catches that nothing else does: a binding to a property that does
# not exist. QML does not throw on one — it evaluates to `undefined`, which a
# colour property renders as transparent and a numeric one as zero. A component
# with a typo in a token name renders as nothing at all, and without this the
# build says nothing about it.
#
# Every category qmllint knows is on, including `unqualified` and
# `missing-property`, which is only possible because nothing here reaches C++
# through a context property: the token objects are typed singletons in the
# `Kvit.Ui` module, so a static analyser can see them. kvit-notes had to
# disable both categories for years for exactly the opposite reason.
#
# Finding qmllint: $QT_ROOT_DIR/bin/qmllint if that is set (CI exports it via
# install-qt-action), otherwise the newest 6.x kit under ~/Qt, otherwise PATH.
set -e
cd "$(dirname "$0")/.."

if [ -n "$QT_ROOT_DIR" ] && [ -x "$QT_ROOT_DIR/bin/qmllint" ]; then
    QMLLINT="$QT_ROOT_DIR/bin/qmllint"
elif [ -d "$HOME/Qt" ]; then
    QT_VERSION=$(ls "$HOME/Qt" | grep -E '^6\.' | sort -V | tail -1)
    QMLLINT="$HOME/Qt/$QT_VERSION/gcc_64/bin/qmllint"
fi
if [ ! -x "${QMLLINT:-}" ]; then
    QMLLINT=$(command -v qmllint || true)
fi
if [ ! -x "${QMLLINT:-}" ]; then
    echo "qmllint not found; set QT_ROOT_DIR to a Qt installation" >&2
    exit 1
fi

# The module's own type description, so `import Kvit.Ui` resolves. It is
# written next to the built module by qt_add_qml_module, and the
# kvit-ui-qmltypes target in ALL is what guarantees it is current after any
# build. A lint run against a stale description reports errors that are not
# there and misses the ones that are.
BUILD_DIR="${KVIT_UI_BUILD_DIR:-build}"
MODULES="$BUILD_DIR/qml-modules"
if [ ! -f "$MODULES/Kvit/Ui/qmldir" ]; then
    echo "no module description under $MODULES/Kvit/Ui." >&2
    echo "Build first, or set KVIT_UI_BUILD_DIR to a build directory." >&2
    exit 1
fi

IMPORTS=(-I "$MODULES")

files=("$@")
if [ ${#files[@]} -eq 0 ]; then
    mapfile -t files < <(ls qml/*.qml gallery/*.qml)
fi

# Every category is an error rather than a warning, and the tree is clean at
# that setting. Demoting a category to `warning` is the same as having no gate
# for what it catches: a warning nobody fails on is a warning nobody reads.
echo "Linting ${#files[@]} QML files..."
"$QMLLINT" "${IMPORTS[@]}" \
    --unqualified error \
    --missing-property error \
    --unused-imports error \
    --deprecated error \
    "${files[@]}"
echo "qmllint clean"
