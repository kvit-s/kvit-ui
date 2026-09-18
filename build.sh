#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
set -e

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BUILD_DIR="$PROJECT_DIR/build"

# Settings the five Kvit repositories agree on: the compiler cache, the linker,
# how debug information is stored, where tests write their evidence. They live
# in one place, ~/kvit-build, so that changing one of those decisions is one
# edit rather than five; kvit-works/build-infra.md is the design they belong to.
#
# Absent, this builds exactly as it did before, which is what keeps the
# repository buildable for continuous integration and on Windows.
KVIT_BUILD_ENV="${KVIT_BUILD_ENV:-$HOME/kvit-build/kvit-build-env.sh}"
KVIT_CMAKE_SHARED_ARGS=()
if [ -f "$KVIT_BUILD_ENV" ]; then
    . "$KVIT_BUILD_ENV"
    kvit_build_env "$PROJECT_DIR"
fi

# Find Qt. The newest 6.x kit under ~/Qt unless told otherwise, which is what
# kvit-notes' script does and what every machine in this estate has.
if [ -d "$HOME/Qt" ]; then
    QT_VERSION=$(ls "$HOME/Qt" | grep -E '^6\.' | sort -V | tail -1)
    if [ -n "$QT_VERSION" ]; then
        QT_PATH="$HOME/Qt/$QT_VERSION/gcc_64"
    fi
fi

CLEAN=0
RUN_TESTS=0
RUN_GALLERY=0
HEADLESS=0
SHOW_SHOTS=0
SHOT_ONLY=0

for arg in "$@"; do
    case $arg in
        --clean) CLEAN=1 ;;
        --test) RUN_TESTS=1 ;;
        --gallery) RUN_GALLERY=1 ;;
        --headless) HEADLESS=1 ;;
        --shots) SHOW_SHOTS=1 ;;
        --shots-only) SHOT_ONLY=1; SHOW_SHOTS=1 ;;
        --qt=*) QT_PATH="${arg#*=}" ;;
        --help)
            echo "Usage: ./build.sh [options]"
            echo "Options:"
            echo "  --clean       Clean the build directory first"
            echo "  --test        Run ctest after building"
            echo "  --gallery     Launch the gallery after building"
            echo "  --headless    Run without a display (offscreen platform)"
            echo "  --shots       Write the gallery screenshot set after tests"
            echo "  --shots-only  Write the screenshot set and nothing else"
            echo "  --qt=PATH     Use this Qt installation"
            exit 0
            ;;
    esac
done

if [ "$CLEAN" -eq 1 ] && [ -d "$BUILD_DIR" ]; then
    echo "Cleaning the build directory..."
    rm -rf "$BUILD_DIR"
fi

mkdir -p "$BUILD_DIR"

(
    cd "$BUILD_DIR"

    # Shared libraries for a local build.
    #
    # Every test executable links the token layer and the QML module, and a
    # static link copies what it pulls into each of them. It is not the
    # default in CMakeLists.txt because packaging has no flag that marks it as
    # packaging: a release preset sets only CMAKE_BUILD_TYPE, and a default of
    # ON would reach those and produce an application whose libraries were
    # never installed beside it. Asking for it here keeps it to the builds run
    # from this script, which is where the disk is.
    CMAKE_ARGS="-DKVIT_UI_SHARED_LIBS=ON"
    if [ -n "$QT_PATH" ]; then
        echo "Using Qt from: $QT_PATH"
        CMAKE_ARGS="$CMAKE_ARGS -DCMAKE_PREFIX_PATH=$QT_PATH"
    fi

    cmake .. $CMAKE_ARGS "${KVIT_CMAKE_SHARED_ARGS[@]}"
    make -j"$(nproc)"

    echo ""
    echo "Build complete."

    if [ "$RUN_TESTS" -eq 1 ]; then
        echo ""
        if [ "$HEADLESS" -eq 1 ]; then
            echo "Running tests (headless)..."
            export QT_QPA_PLATFORM=offscreen
        else
            echo "Running tests..."
        fi
        ctest --output-on-failure
    fi

    if [ "$SHOW_SHOTS" -eq 1 ]; then
        # The screenshot set: one image per component per theme, under fixed
        # names. What gets reviewed after a token change is the diff against
        # the previous run rather than the whole set again.
        SHOT_DIR="$BUILD_DIR/screenshots"
        mkdir -p "$SHOT_DIR"
        echo ""
        echo "Writing the gallery screenshot set to $SHOT_DIR..."
        QT_QPA_PLATFORM=${QT_QPA_PLATFORM:-offscreen} \
            "$BUILD_DIR/gallery/kvit-ui-gallery" --shots "$SHOT_DIR"
        echo "$(ls "$SHOT_DIR" | wc -l) images"
    fi

    if [ "$RUN_GALLERY" -eq 1 ]; then
        echo ""
        "$BUILD_DIR/gallery/kvit-ui-gallery"
    fi
)
