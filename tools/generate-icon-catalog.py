# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Generate src/qml/iconcatalog.cpp from the icon font itself.
#
#   python3 tools/generate-icon-catalog.py            # rewrite the catalogue
#   python3 tools/generate-icon-catalog.py --check     # fail if it is stale
#
# Why this is generated. kvit-notes-pro's WorksIcon.qml carries a hand-written
# switch of twelve names to twelve codepoints, transcribed from the Phosphor
# stylesheet. Twelve is transcribable; the whole font is 1,200-odd symbols, and
# a shared component has to be able to answer for any of them or every
# consuming application goes back to hand-transcribing the two it needs. It
# also means the vocabulary skill's catalogue lists what actually exists rather
# than what somebody remembered to write down.
#
# Where the names come from. A TrueType font can carry glyph names in its
# `post` table, and Phosphor's does not: its post table is version 3.0, which
# is the version that says "no names here". So the names come from the
# stylesheet Phosphor publishes with the font, vendored beside it as
# resources/fonts/Phosphor-regular.css, and every codepoint the stylesheet
# claims is checked against the font's own `cmap` before it is written out.
# That check is what makes the pair trustworthy: a stylesheet from a different
# release than the .ttf would name codepoints the font has no glyph for, and
# the generator stops rather than producing a catalogue of empty boxes.
#
# The meaning names are the layer above. They are written here by hand and
# deliberately differ from Phosphor's own: what this estate calls
# `chevron-right` is Phosphor's `caret-right`, and a call site should ask for
# the meaning rather than for the drawing, so that the drawing can change
# without touching the call site.

import re
import struct
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
FONT = ROOT / "resources" / "fonts" / "Phosphor.ttf"
STYLESHEET = ROOT / "resources" / "fonts" / "Phosphor-regular.css"
OUT = ROOT / "src" / "qml" / "iconcatalog.cpp"

def read_tables(data):
    """The font's table directory: tag -> (offset, length)."""
    if data[:4] == b"ttcf":
        raise SystemExit("font collections are not supported")
    num_tables = struct.unpack(">H", data[4:6])[0]
    tables = {}
    for i in range(num_tables):
        base = 12 + i * 16
        tag, _checksum, offset, length = struct.unpack(
            ">4sIII", data[base:base + 16])
        tables[tag.decode("latin-1")] = (offset, length)
    return tables


def parse_cmap(data, offset):
    """codepoint -> glyph id, from the best Unicode subtable available."""
    num_subtables = struct.unpack(">H", data[offset + 2:offset + 4])[0]
    best = None
    for i in range(num_subtables):
        base = offset + 4 + i * 8
        platform, encoding, sub_offset = struct.unpack(
            ">HHI", data[base:base + 8])
        # Windows/Unicode-full and Windows/BMP are the two that matter; a
        # Unicode-platform table is equivalent and is taken if it is all there is.
        rank = {(3, 10): 3, (3, 1): 2, (0, 4): 3, (0, 3): 2, (0, 6): 3}.get(
            (platform, encoding), 0)
        if rank and (best is None or rank > best[0]):
            best = (rank, offset + sub_offset)
    if best is None:
        raise SystemExit("no Unicode cmap subtable in the font")

    sub = best[1]
    fmt = struct.unpack(">H", data[sub:sub + 2])[0]
    mapping = {}
    if fmt == 4:
        seg_x2 = struct.unpack(">H", data[sub + 6:sub + 8])[0]
        segments = seg_x2 // 2
        ends = struct.unpack(">%dH" % segments, data[sub + 14:sub + 14 + seg_x2])
        starts_at = sub + 16 + seg_x2
        starts = struct.unpack(">%dH" % segments,
                               data[starts_at:starts_at + seg_x2])
        deltas_at = starts_at + seg_x2
        deltas = struct.unpack(">%dh" % segments,
                               data[deltas_at:deltas_at + seg_x2])
        ranges_at = deltas_at + seg_x2
        ranges = struct.unpack(">%dH" % segments,
                               data[ranges_at:ranges_at + seg_x2])
        for i in range(segments):
            for code in range(starts[i], min(ends[i], 0xFFFF) + 1):
                if ranges[i] == 0:
                    glyph = (code + deltas[i]) & 0xFFFF
                else:
                    at = ranges_at + i * 2 + ranges[i] + (code - starts[i]) * 2
                    if at + 2 > len(data):
                        continue
                    glyph = struct.unpack(">H", data[at:at + 2])[0]
                    if glyph:
                        glyph = (glyph + deltas[i]) & 0xFFFF
                if glyph:
                    mapping[code] = glyph
    elif fmt == 12:
        groups = struct.unpack(">I", data[sub + 12:sub + 16])[0]
        for i in range(groups):
            base = sub + 16 + i * 12
            start, end, glyph = struct.unpack(">III", data[base:base + 12])
            for k in range(end - start + 1):
                mapping[start + k] = glyph + k
    else:
        raise SystemExit("cmap subtable format %d is not supported" % fmt)
    return mapping


def parse_stylesheet(text):
    """Phosphor name -> codepoint, from the `.ph.ph-NAME:before { content }`
    rules the published stylesheet is made of."""
    pattern = re.compile(
        r'\.ph\.ph-([a-z0-9-]+):before\s*\{\s*content:\s*"\\([0-9a-fA-F]{4,6})"')
    return {name: int(code, 16) for name, code in pattern.findall(text)}


# ── the meaning names ──────────────────────────────────────────────────────
#
# Each entry is (the name a call site writes, the Phosphor glyph it draws).
#
# The twelve at the top are kvit-notes-pro's, unchanged, so its six call sites
# keep working when it drops WorksIcon.qml in Wave 5. The block after them is
# one name for each of the nine literal Unicode characters kvit-notes passes to
# IconButton today — `›`, `‹`, `×`, `+`, `▰`, `▤`, `⌕`, `N`, `#` — which is
# what Wave 5 needs somewhere to send them. The rest are what the components in
# this repository ask for: a disclosure's chevron, a checkbox's tick, a
# select's indicator, a sort arrow, a close button on a toast.
MEANINGS = [
    # kvit-notes-pro's twelve
    ("plus", "plus"),
    ("messages-square", "chats"),
    ("message-square", "chat"),
    ("send", "paper-plane-right"),
    ("chevron-right", "caret-right"),
    ("chevron-down", "caret-down"),
    ("zoom-in", "magnifying-glass-plus"),
    ("zoom-out", "magnifying-glass-minus"),
    ("pencil", "pencil-simple"),
    ("check", "check"),
    ("archive", "archive"),
    ("rotate-ccw", "arrow-counter-clockwise"),

    # the nine literals kvit-notes draws today
    ("chevron-left", "caret-left"),          # ‹
    ("chevron-up", "caret-up"),              # ▲ (the pair of chevron-down)
    ("close", "x"),                          # ×
    ("sidebar", "sidebar-simple"),           # ▰ — show or hide the sidebar
    ("list", "list"),                        # ▤ — the note list
    ("search", "magnifying-glass"),          # ⌕
    ("note", "note"),                        # N — a note
    ("tag", "hash"),                         # # — a tag

    # what the components here ask for
    ("caret-sort", "caret-up-down"),
    ("sort-ascending", "sort-ascending"),
    ("sort-descending", "sort-descending"),
    ("dot", "dot"),
    ("dot-outline", "dot-outline"),
    ("circle", "circle"),
    ("square", "square"),
    ("minus", "minus"),
    ("info", "info"),
    ("warning", "warning"),
    ("error", "warning-circle"),
    ("success", "check-circle"),
    ("question", "question"),
    ("settings", "gear"),
    ("filter", "funnel"),
    ("columns", "columns"),
    ("calendar", "calendar-blank"),
    ("clock", "clock"),
    ("user", "user"),
    ("robot", "robot"),
    ("folder", "folder"),
    ("file", "file"),
    ("link", "link-simple"),
    ("external", "arrow-square-out"),
    ("copy", "copy"),
    ("trash", "trash"),
    ("undo", "arrow-arc-left"),
    ("redo", "arrow-arc-right"),
    ("refresh", "arrows-clockwise"),
    ("more", "dots-three"),
    ("more-vertical", "dots-three-vertical"),
    ("drag", "dots-six-vertical"),
    ("pin", "push-pin"),
    ("star", "star"),
    ("eye", "eye"),
    ("eye-off", "eye-slash"),
    ("lock", "lock-simple"),
    ("arrow-up", "arrow-up"),
    ("arrow-down", "arrow-down"),
    ("arrow-left", "arrow-left"),
    ("arrow-right", "arrow-right"),
    ("trend-up", "trend-up"),
    ("trend-down", "trend-down"),
    ("chart", "chart-bar"),
    ("wallet", "wallet"),
    ("coins", "coins"),
    ("bank", "bank"),
    ("receipt", "receipt"),
    ("repeat", "repeat"),
    ("split", "arrows-split"),
    ("merge", "arrows-merge"),
    ("play", "play"),
    ("pause", "pause"),
    ("stop", "stop"),

    # A document kept with a record. kvit-cash marks a transaction row with
    # each state it is in, and four of the five states it draws already had a
    # name here: waiting is `clock`, a split is `split`, a transfer is
    # `repeat`, and one nobody has looked at yet is `dot-outline`. A record
    # with something attached to it had none, and `file`, `link` and `pin`
    # each already mean something else — a document, a target, and something
    # held at the top of a list.
    ("attachment", "paperclip"),

    # The eight kvit-notes-pro's overview asks for. Its section headings and
    # its rows get their icon names from C++ — `actionIcon` on a section and
    # `icon` on a row record — so a call site cannot rewrite them, and six of
    # them resolved to nothing at all: the overview drew a hatched red box
    # where the Changes heading's action should have been.
    #
    # `agent` and `rename` draw glyphs that already have a meaning name here,
    # and that is the two-layer table working rather than a duplicate to
    # remove. A robot is what an agent is drawn as and a pencil is what
    # renaming is drawn as, and a screen about agents that asks for `robot`
    # has written down the drawing instead of what it meant. The deliberate
    # pairs are listed in tests/test_icons.cpp, which still fails on any
    # other two meanings landing on one glyph.
    #
    # `terminal` and `chat` are a decision rather than an addition. Neither
    # had a meaning, so both fell through to the escape hatch that resolves a
    # Phosphor name directly, and both landed on a glyph nobody chose:
    # `terminal` on the bare prompt rather than the window, `chat` on the
    # square bubble rather than the round one. Naming them here is what the
    # meaning layer is for, and it shadows the two Phosphor names the same way
    # `link` already shadows Phosphor's `link` and `pencil` its `pencil`. The
    # square bubble is still reachable, as `message-square`.
    ("diff", "git-diff"),
    ("git", "git-branch"),
    ("agent", "robot"),
    ("ask", "chat-teardrop-dots"),
    ("ask-in", "arrow-bend-up-right"),
    ("rename", "pencil-simple"),
    ("terminal", "terminal-window"),
    ("chat", "chat-circle"),
]

HEADER = """// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// GENERATED FILE — do not edit.
//
// Written by tools/generate-icon-catalog.py from
// resources/fonts/Phosphor-regular.css, checked against the cmap of
// resources/fonts/Phosphor.ttf. Adding a symbol means adding a line to that
// script's MEANINGS list and running it; editing this file by hand fails the
// generator check in CI.
#include "iconcatalog.h"

#include <QHash>

namespace {

// Every glyph the font carries, by its Phosphor name. The names come from the
// published stylesheet and every codepoint below was found in the font's own
// cmap, so this is what the shipped .ttf actually contains rather than what a
// stylesheet from some other release claimed.
struct Glyph
{
    const char *name;
    char32_t codepoint;
};

const Glyph kGlyphs[] = {
"""


def main():
    check = "--check" in sys.argv
    data = FONT.read_bytes()
    tables = read_tables(data)
    if "cmap" not in tables:
        raise SystemExit("the font has no cmap table")
    in_font = parse_cmap(data, tables["cmap"][0])
    declared = parse_stylesheet(STYLESHEET.read_text(encoding="utf-8"))
    if not declared:
        raise SystemExit("no icon rules found in %s" % STYLESHEET)

    # Every name the stylesheet declares, kept only if the font actually has a
    # glyph at that codepoint. Phosphor puts its glyphs in the private use
    # area; a codepoint outside it would be the stylesheet naming an ordinary
    # character, which is not an icon.
    by_name = {}
    absent = []
    for name, code in sorted(declared.items()):
        if not (0xE000 <= code <= 0xF8FF):
            continue
        if code not in in_font:
            absent.append((name, code))
            continue
        by_name[name] = code
    if absent:
        raise SystemExit(
            "%d names in %s point at codepoints the font has no glyph for — "
            "the stylesheet and the .ttf are from different releases. First "
            "few: %s"
            % (len(absent), STYLESHEET.name,
               ", ".join("%s U+%04X" % pair for pair in absent[:5])))

    missing = [(meaning, glyph) for meaning, glyph in MEANINGS
               if glyph not in by_name]
    if missing:
        raise SystemExit(
            "these meaning names point at glyphs the font does not have:\n  "
            + "\n  ".join("%s -> %s" % pair for pair in missing))

    out = [HEADER]
    for name in sorted(by_name):
        out.append('    { "%s", 0x%04x },\n' % (name, by_name[name]))
    out.append("""};

// The meaning names, and the glyph each one draws.
//
// Two layers rather than one because they answer different questions. A call
// site asks for what a symbol means — `chevron-right`, `send`, `trend-up` —
// and the font is asked for what it looks like. Phosphor calls a right-facing
// chevron `caret-right`, and a component written against that name would have
// to change if the drawing ever moved to a different glyph.
struct Meaning
{
    const char *name;
    const char *glyph;
};

const Meaning kMeanings[] = {
""")
    for meaning, glyph in MEANINGS:
        out.append('    { "%s", "%s" },\n' % (meaning, glyph))
    out.append("""};

}   // namespace

namespace KvitUi {

IconCatalog::IconCatalog(QObject *parent)
    : QObject(parent)
{
}

QString IconCatalog::glyph(const QString &name) const
{
    static const QHash<QString, char32_t> byGlyphName = [] {
        QHash<QString, char32_t> map;
        for (const Glyph &g : kGlyphs)
            map.insert(QString::fromLatin1(g.name), g.codepoint);
        return map;
    }();
    static const QHash<QString, QString> byMeaning = [] {
        QHash<QString, QString> map;
        for (const Meaning &m : kMeanings)
            map.insert(QString::fromLatin1(m.name), QString::fromLatin1(m.glyph));
        return map;
    }();

    // A meaning name first, then the Phosphor name it resolves to. Asking for
    // a Phosphor name directly works and is the escape hatch for a symbol that
    // has not been given a meaning yet; the vocabulary skill's catalogue lists
    // the meanings, because those are what a call site should be written in.
    const QString resolved = byMeaning.value(name, name);
    const auto found = byGlyphName.constFind(resolved);
    if (found == byGlyphName.constEnd())
        return QString();
    return QString::fromUcs4(&found.value(), 1);
}

QStringList IconCatalog::meaningNames()
{
    QStringList names;
    names.reserve(int(std::size(kMeanings)));
    for (const Meaning &m : kMeanings)
        names.append(QString::fromLatin1(m.name));
    return names;
}

QStringList IconCatalog::glyphNames()
{
    QStringList names;
    names.reserve(int(std::size(kGlyphs)));
    for (const Glyph &g : kGlyphs)
        names.append(QString::fromLatin1(g.name));
    return names;
}

}   // namespace KvitUi
""")
    text = "".join(out)

    if check:
        current = OUT.read_text(encoding="utf-8") if OUT.exists() else ""
        if current != text:
            print("src/qml/iconcatalog.cpp is stale; run "
                  "tools/generate-icon-catalog.py", file=sys.stderr)
            return 1
        print("iconcatalog.cpp matches the font")
        return 0

    # Generated C++ is repository text, not locale text. On Windows the locale
    # codec is commonly cp1252, which rewrites the em dashes in the header and
    # makes an unchanged UTF-8 file fail `--check`.
    with OUT.open("w", encoding="utf-8", newline="\n") as output:
        output.write(text)
    print("wrote %s: %d glyphs, %d meaning names"
          % (OUT.relative_to(ROOT), len(by_name), len(MEANINGS)))
    return 0


if __name__ == "__main__":
    sys.exit(main())
