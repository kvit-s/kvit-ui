#!/usr/bin/env bash
# Render every mockup in this directory to a PNG under renders/.
#
# The mockups are static HTML with no network access of any kind: no web fonts,
# no scripts fetched from anywhere, no images that are not inlined. That is what
# lets them render identically here and on any machine.
#
# Usage:
#   ./render.sh                # render every *.html in this directory
#   ./render.sh today-a.html   # render just these
#
# Viewport is 1440x960, which is the desktop size the design assumes.

set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
out="$here/renders"
mkdir -p "$out"

shell_bin="${HEADLESS_SHELL:-}"
if [[ -z "$shell_bin" ]]; then
  shell_bin="$(ls -d "$HOME"/.cache/ms-playwright/chromium_headless_shell-*/chrome-linux/headless_shell 2>/dev/null | sort -V | tail -1 || true)"
fi
if [[ -z "$shell_bin" || ! -x "$shell_bin" ]]; then
  echo "no headless chromium found; set HEADLESS_SHELL to one" >&2
  exit 1
fi

width="${WIDTH:-1440}"
height="${HEIGHT:-960}"

files=("$@")
if [[ ${#files[@]} -eq 0 ]]; then
  mapfile -t files < <(cd "$here" && ls *.html 2>/dev/null || true)
fi
if [[ ${#files[@]} -eq 0 ]]; then
  echo "nothing to render" >&2
  exit 1
fi

for f in "${files[@]}"; do
  name="$(basename "$f" .html)"
  src="$here/$(basename "$f")"
  [[ -f "$src" ]] || { echo "missing: $src" >&2; exit 1; }
  "$shell_bin" --headless --disable-gpu --no-sandbox --hide-scrollbars \
    --force-device-scale-factor=1 --default-background-color=00000000 \
    --screenshot="$out/$name.png" --window-size="$width,$height" \
    "file://$src" >/dev/null 2>&1
  echo "$out/$name.png"
done
