@echo off
rem This Source Code Form is subject to the terms of the Mozilla Public
rem License, v. 2.0. If a copy of the MPL was not distributed with this
rem file, You can obtain one at https://mozilla.org/MPL/2.0/.
rem
rem Windows parity for build.sh, kept from the start rather than added later:
rem the one thing a cross-platform build breaks on is a path or a compiler
rem flag that was only ever tried on Linux, and finding that out at release
rem time is what makes it expensive.

setlocal
rem Resolve `%~dp0` without its trailing backslash. Passing a quoted path that
rem ends in `\` through CMake's Windows argument parser escapes the closing
rem quote and folds the following `-B`/`-G` arguments into the source path.
for %%I in ("%~dp0.") do set "PROJECT_DIR=%%~fI"
set "BUILD_DIR=%PROJECT_DIR%\build"

if "%QT_ROOT_DIR%"=="" (
    echo Set QT_ROOT_DIR to your Qt installation, for example
    echo   set QT_ROOT_DIR=C:\Qt\6.10.1\msvc2022_64
    exit /b 1
)
set "PATH=%QT_ROOT_DIR%\bin;%PATH%"

if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

cmake -S "%PROJECT_DIR%" -B "%BUILD_DIR%" -G "Visual Studio 17 2022" -A x64 ^
      -DCMAKE_PREFIX_PATH="%QT_ROOT_DIR%"
if errorlevel 1 exit /b 1

cmake --build "%BUILD_DIR%" --config Release
if errorlevel 1 exit /b 1

if "%1"=="--test" (
    ctest --test-dir "%BUILD_DIR%" -C Release --output-on-failure
)
