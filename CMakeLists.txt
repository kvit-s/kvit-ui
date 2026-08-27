cmake_minimum_required(VERSION 3.21)
project(kvit-ui VERSION 1.0.0 LANGUAGES CXX)

# C++20 matches the two editors and kvit-hub's desktop, so a consumer that
# adds this repository as a submodule never has to reconcile two standards.
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

# Qt Widgets is deliberately absent.
#
# prd.md Decision 3 puts kvit-cash in QML throughout, with its transaction
# browser on KvitTable, and holds a Qt Widgets fallback open only until the
# 250,000-row benchmark in tests/ passes. The QPalette bridge and the C++
# metrics helper that fallback would need are not built here, and leaving
# Widgets out of this line is the visible marker of that. If the benchmark
# ever fails, adding it back is this one word.
find_package(Qt6 REQUIRED COMPONENTS Core Gui Qml Quick QuickControls2 Test QuickTest)

# Where the Kvit.Ui QML module's generated qmldir and .qmltypes live in the
# resource system: `:/qt/qml/Kvit/Ui`, which is on the engine's default import
# path. Unset, every configure prints a QTP0001 warning.
qt_policy(SET QTP0001 NEW)

# How much debug information to compile with.
#
# The default `-g` on a Qt tree writes tens of gigabytes of DWARF into the
# build directory and slows every link, and almost none of it is read: a
# stack trace with function names and line numbers needs `-g1`, and only
# sitting at a debugger inspecting locals needs the rest. This must be set
# before any target is defined, because CMAKE_CXX_FLAGS is captured per
# target at creation time.
option(KVIT_UI_FULL_DEBUG_INFO
    "Compile with full debug info (-g) instead of line tables only (-g1)" OFF)
if(NOT KVIT_UI_FULL_DEBUG_INFO
   AND (CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang"))
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g1")
endif()

# Static or shared internal libraries.
#
# OFF is the shipping arrangement, and it is the default rather than a
# packaging flag because packaging has no flag that marks it as packaging: a
# release preset sets only CMAKE_BUILD_TYPE. A default of ON would reach those
# and produce an application whose libraries were never installed beside it.
# build.sh turns it ON for local development, where a static link would copy
# the whole token layer into each test executable.
option(KVIT_UI_SHARED_LIBS "Build the internal libraries shared" OFF)
if(KVIT_UI_SHARED_LIBS)
    set(KVIT_UI_LIB_TYPE SHARED)
else()
    set(KVIT_UI_LIB_TYPE STATIC)
endif()

option(KVIT_UI_BUILD_GALLERY "Build the gallery application" ON)
option(KVIT_UI_BUILD_TESTS "Build the test suite" ON)

# Warnings are errors on the compilers this is developed against. The token
# layer's switch statements over role enumerations rely on -Werror=switch to
# stop the build when a role is added without a size, which is the mechanism
# Typography's comment describes.
if(CMAKE_CXX_COMPILER_ID MATCHES "GNU|Clang")
    add_compile_options(-Wall -Wextra -Werror=switch)
endif()

add_subdirectory(src/tokens)
add_subdirectory(qml)
add_subdirectory(tools)

# The gallery before the tests: tests/test_gallery.cpp links the gallery's QML
# module, and CMake needs the target to exist before it is named.
if(KVIT_UI_BUILD_GALLERY)
    add_subdirectory(gallery)
endif()

if(KVIT_UI_BUILD_TESTS)
    enable_testing()
    add_subdirectory(tests)
endif()
