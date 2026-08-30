// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Write ux/tokens.css from the running token table.
//
//   tokens-to-css <path>            # write the file
//   tokens-to-css --check <path>    # exit 1 if the file is stale
//
// Why this exists. kvit-hub's tokens.css was a hand-copied mirror of the C++
// token table, and by the time this repository was started it had drifted:
// comparing five dark-theme tokens on 2026-08-27, four differed — `textPrimary`
// #eeeeee against #e8e8e8, `borderStrong` #696969 against #5a5a5a, `quoteBar`
// the same pair, and `textFaint` #858585 against #848484. Drawings were being
// judged in colours the application does not draw with, which makes the whole
// draw-then-implement workflow quietly unreliable: a screen approved in the
// drawing comes out different when it is built.
//
// A C++ tool rather than a script that reads theme.cpp, because a script that
// parses C++ is a second implementation of the token table and would drift in
// its own way. This links the real one and asks it.
//
// Every colour comes out by introspection over Theme's own Q_PROPERTY list, so
// a token added to the table appears here without anybody remembering to add
// it. Only the *grouping* is written down, and a token that matches no group
// is reported rather than quietly appended, so a new token gets put somewhere
// deliberately.

#include <QCoreApplication>
#include <QFile>
#include <QMetaObject>
#include <QMetaProperty>
#include <QStringList>
#include <QTextStream>

#include <cstdio>

#include "interfacemetrics.h"
#include "theme.h"

namespace {

// The four themes and the selector each one is written under. Dark is `:root`
// as well as `.theme-dark`, because that is what kvit-hub's mockups assume: a
// drawing with no theme class on its body renders in the theme the owner's
// dashboard is set to.
struct ThemeRow
{
    const char *id;
    const char *selector;
};

const ThemeRow kThemes[] = {
    { "dark", ":root,\n.theme-dark" },
    { "light", ".theme-light" },
    { "sepia", ".theme-sepia" },
    // `.theme-contrast` rather than `.theme-highContrast`: the class name is
    // what 31 existing mockups already write.
    { "highContrast", ".theme-contrast" },
};

// How the tokens are grouped in the output, and in what order. The strings are
// property-name prefixes and exact names; a property matches the first group
// that claims it.
struct Group
{
    const char *heading;
    QStringList members;
};

const QList<Group> &groups()
{
    static const QList<Group> table = {
        { "surfaces",
          { "windowBackground", "panelBackground", "listBackground",
            "footerBackground", "popupBackground", "chipBackground",
            "bannerBackground", "codePanelBackground" } },
        { "text",
          { "textPrimary", "textSecondary", "textMuted", "textFaint",
            "textDisabled", "bannerText", "onAccent" } },
        { "lines and glyphs",
          { "border", "borderStrong", "quoteBar", "mutedGlyph" } },
        { "interactive tints",
          { "hoverTint", "blockHoverTint", "focusTint", "focusRing",
            "selectionTint", "selectionActiveTint", "blockSelectionTint" } },
        { "accent and semantic",
          { "accent", "danger", "dangerBright", "success", "warning",
            "pinColor", "link", "marker", "inlineCodeBackground",
            "highlightBackground", "searchMatchBackground",
            "searchCurrentBackground", "changedTextBackground",
            "addedTextBackground", "removedTextBackground",
            "calloutTip" } },
        { "code, for a mockup that shows a document",
          { "codeKeyword", "codeType", "codeString", "codeComment",
            "codeNumber" } },
        { "portfolio vocabulary",
          { "axisAttention", "axisAttentionText", "axisAgent", "axisAgentText",
            "scopeDiscovered", "signalHard", "signalSoft", "signalHygiene",
            "hatchAlt" } },
    };
    return table;
}

// windowBackground -> --window-background
QString cssName(const QString &property)
{
    QString out = QStringLiteral("--");
    for (const QChar ch : property) {
        if (ch.isUpper()) {
            out.append(QLatin1Char('-'));
            out.append(ch.toLower());
        } else {
            out.append(ch);
        }
    }
    return out;
}

QStringList colorProperties()
{
    QStringList names;
    const QMetaObject *meta = &Theme::staticMetaObject;
    for (int i = 0; i < meta->propertyCount(); ++i) {
        const QMetaProperty property = meta->property(i);
        if (property.metaType().id() == QMetaType::QColor)
            names.append(QString::fromLatin1(property.name()));
    }
    return names;
}

QString hexOf(const Theme &theme, const QString &property)
{
    const QVariant value = theme.property(property.toLatin1().constData());
    return value.value<QColor>().name(QColor::HexRgb);
}

void writeThemeBlock(QTextStream &out, const ThemeRow &row,
                     const QStringList &properties, QStringList *ungrouped)
{
    Theme theme;
    theme.setThemeId(QString::fromLatin1(row.id));

    out << row.selector << " {\n";

    QSet<QString> written;
    bool firstGroup = true;
    for (const Group &group : groups()) {
        QStringList present;
        for (const QString &member : group.members) {
            if (properties.contains(member))
                present.append(member);
        }
        if (present.isEmpty())
            continue;
        if (!firstGroup)
            out << "\n";
        firstGroup = false;
        out << "  /* " << group.heading << " */\n";
        for (const QString &member : present) {
            out << "  " << cssName(member) << ": " << hexOf(theme, member)
                << ";\n";
        }
        written.unite(QSet<QString>(present.begin(), present.end()));
    }

    for (const QString &property : properties) {
        if (!written.contains(property) && !ungrouped->contains(property))
            ungrouped->append(property);
    }

    // The chart ramps, as one variable per step. A CSS list would be shorter
    // and unusable: a drawing needs `var(--categorical-3)` for one bar.
    const auto ramp = [&out](const char *prefix, const QStringList &steps) {
        out << "\n  /* " << prefix << " ramp */\n";
        for (int i = 0; i < steps.size(); ++i) {
            out << "  --" << prefix << "-" << i << ": " << steps.at(i)
                << ";\n";
        }
    };
    ramp("categorical", theme.categoricalRamp());
    ramp("sequential", theme.sequentialRamp());
    ramp("diverging", theme.divergingRamp());

    out << "}\n";
}

// The type and density scale, at the default interface size.
//
// One block rather than four, because these do not vary by theme. They are
// written out at the default base, which is what a drawing is judged at; a
// mockup that wants to see a larger interface changes the variables by hand,
// which is what kvit-hub already does.
void writeScaleBlock(QTextStream &out)
{
    InterfaceMetrics metrics;

    out << ":root {\n";
    out << "  /* Local families only — no @font-face, no network, so a mockup\n"
           "   * renders identically here and on any machine. */\n";
    out << "  --font-ui: Ubuntu, \"Ubuntu Sans\", \"DejaVu Sans\", system-ui, "
           "sans-serif;\n";
    out << "  --font-mono: \"DejaVu Sans Mono\", \"Ubuntu Mono\", "
           "ui-monospace, monospace;\n";

    struct Entry
    {
        const char *name;
        int value;
        const char *note;
    };

    const Entry type[] = {
        { "type-caption", metrics.caption(), "kind tags, counts" },
        { "type-small", metrics.small(), "chip labels, sub-lines" },
        { "type-body", metrics.body(), "row text, prose" },
        { "type-strong", metrics.strong(), "a name, an emphasised row" },
        { "type-title", metrics.title(), "a section heading" },
        { "type-headline", metrics.headline(), "a pane title, the wordmark" },
        { "type-display", metrics.display(), "a page title" },
    };
    out << "\n  /* The seven type roles, in pixels at the default interface\n"
           "   * size of " << metrics.fontSize() << ". */\n";
    for (const Entry &entry : type) {
        out << "  --" << entry.name << ": " << entry.value << "px;"
            << QString(qMax(1, 22 - int(qstrlen(entry.name))
                                   - int(QString::number(entry.value).size())),
                       QLatin1Char(' '))
            << "/* " << entry.note << " */\n";
    }

    const Entry space[] = {
        { "space-tight", metrics.spaceTight(), nullptr },
        { "space-snug", metrics.spaceSnug(), nullptr },
        { "space-near", metrics.spaceNear(), nullptr },
        { "space", metrics.space(), "the ordinary gap" },
        { "space-wide", metrics.spaceWide(), nullptr },
        { "space-loose", metrics.spaceLoose(), nullptr },
    };
    out << "\n  /* The spacing scale. 3 and 5 are deliberately not steps. */\n";
    for (const Entry &entry : space) {
        out << "  --" << entry.name << ": " << entry.value << "px;";
        if (entry.note)
            out << "  /* " << entry.note << " */";
        out << "\n";
    }

    const Entry density[] = {
        { "view-margin", metrics.viewMargin(), "the outer margin of a view" },
        { "column-gap", metrics.columnGap(), nullptr },
        { "stack-gap", metrics.stackGap(), "between stacked blocks" },
        { "sidebar-width", metrics.sidebarWidth(), nullptr },
        { "rail-width", metrics.railWidth(), "the sidebar collapsed" },
        { "pane-width", metrics.paneWidth(), "the side pane" },
        { "header-height", metrics.headerHeight(), nullptr },
        { "breadcrumb-height", metrics.breadcrumbHeight(), nullptr },
        { "status-bar-height", metrics.statusBarHeight(), nullptr },
        { "row-height", metrics.rowHeight(), "name over description" },
        { "row-height-sub", metrics.rowHeightSub(), "an expanded line" },
        { "row-height-slim", metrics.rowHeightSlim(), "a list at rest" },
        { "row-height-compact", metrics.rowHeightCompact(), "a disclosure" },
        { "control-height", metrics.controlHeight(), nullptr },
        { "tab-height", metrics.tabHeight(), nullptr },
        { "chip-height", metrics.chipHeight(), nullptr },
        { "tag-height", metrics.tagHeight(), nullptr },
        { "pill-height", metrics.pillHeight(), nullptr },
        { "bar-height", metrics.barHeight(), "a compact bar" },
        { "bar-height-wide", metrics.barHeightWide(), "a full-size bar" },
        { "icon-size", metrics.iconSize(), nullptr },
        { "icon-size-small", metrics.iconSizeSmall(), nullptr },
        { "radius-bar", metrics.radiusBar(), nullptr },
        { "radius-chip", metrics.radiusChip(), nullptr },
        { "radius-control", metrics.radiusControl(), nullptr },
        { "radius-card", metrics.radiusCard(), nullptr },
        { "radius-pill", metrics.radiusPill(), nullptr },
        { "hairline", metrics.hairline(), nullptr },
        { "focus-ring-width", metrics.focusRingWidth(), nullptr },
        { "width-floor", metrics.widthFloor(), "below this, effort dots collide" },
        { "width-laptop", metrics.widthLaptop(), "the sidebar collapses" },
        { "width-drawn", metrics.widthDrawn(), "every drawing is this wide" },
    };
    // kvit-hub's nine type names, mapped onto the merged seven.
    //
    // Thirty-one mockups are written against these, and they are what Decision
    // 2 renamed. Without the mapping every `font-size: var(--type-row)` in
    // them resolves to nothing, CSS falls back to the inherited size, and the
    // whole page reflows into something that looks plausible and is wrong.
    //
    // The values are the merged scale, so a mockup drawn before the merge
    // renders at the sizes the application will now use, which is the point:
    // the drawing and the screen agree again. What a mockup must not do is
    // keep using these names in new work — they are here so that Wave 3 can
    // move kvit-hub's mockups over deliberately rather than by discovering
    // that they broke.
    const Entry legacy[] = {
        { "type-micro", metrics.caption(), "was baseSize - 6" },
        { "type-small", metrics.small(), "was baseSize - 5" },
        { "type-secondary", metrics.small(), "was baseSize - 4" },
        { "type-row", metrics.body(), "was baseSize - 3" },
        { "type-name", metrics.strong(), "was baseSize" },
        { "type-heading", metrics.title(), "was baseSize + 1" },
        { "type-page", metrics.headline(), "was baseSize + 2" },
    };
    out << "\n  /* kvit-hub's older type names, on the merged scale. New work\n"
           "   * uses the seven roles above; these keep the existing mockups\n"
           "   * rendering until Wave 3 moves them over. */\n";
    for (const Entry &entry : legacy) {
        out << "  --" << entry.name << ": " << entry.value << "px;";
        if (entry.note)
            out << "  /* " << entry.note << " */";
        out << "\n";
    }

    out << "\n  /* Density. One reader, desktop, no mobile target. */\n";
    for (const Entry &entry : density) {
        out << "  --" << entry.name << ": " << entry.value << "px;";
        if (entry.note)
            out << "  /* " << entry.note << " */";
        out << "\n";
    }
    out << "}\n";
}

QString render()
{
    QString text;
    QTextStream out(&text);

    out << "/* kvit-ui design tokens, for the drawing layer.\n"
           " *\n"
           " * GENERATED FILE — do not edit.\n"
           " *\n"
           " * Written by tools/tokens-to-css from the same C++ token table the\n"
           " * applications draw with (src/tokens/theme.cpp and\n"
           " * src/tokens/interfacemetrics.h), so a drawing is judged in the\n"
           " * colours the application will actually use. Editing this file by\n"
           " * hand fails the generator check.\n"
           " *\n"
           " * Values only. The window, the base elements and the component\n"
           " * classes are in frame.css, which is hand-written.\n"
           " *\n"
           " * Themes: :root carries dark. Put .theme-light, .theme-dark,\n"
           " * .theme-sepia or .theme-contrast on <body> to render a mockup in\n"
           " * another one.\n"
           " */\n\n";

    out << "/* ══════════════════════════════════════════════ theme tokens ══════ */\n\n";

    const QStringList properties = colorProperties();
    QStringList ungrouped;
    for (const ThemeRow &row : kThemes) {
        writeThemeBlock(out, row, properties, &ungrouped);
        out << "\n";
    }

    if (!ungrouped.isEmpty()) {
        // Loud rather than silent: a token nobody put in a group would
        // otherwise be missing from every drawing, and the drawing would look
        // right because the browser falls back to `inherit`.
        std::fprintf(stderr,
                     "tokens-to-css: %d colour token(s) belong to no group in "
                     "tools/tokens-to-css.cpp and were left out: %s\n",
                     int(ungrouped.size()),
                     qPrintable(ungrouped.join(QStringLiteral(", "))));
    }

    out << "/* ═══════════════════════════════════ type, density, geometry ══════ */\n\n";
    writeScaleBlock(out);

    out.flush();
    return text;
}

}   // namespace

int main(int argc, char *argv[])
{
    QCoreApplication app(argc, argv);
    QStringList args = app.arguments().mid(1);

    const bool check = args.removeAll(QStringLiteral("--check")) > 0;
    if (args.size() != 1) {
        std::fprintf(stderr, "usage: tokens-to-css [--check] <path>\n");
        return 2;
    }
    const QString path = args.first();
    const QString wanted = render();

    if (check) {
        QFile file(path);
        if (!file.open(QIODevice::ReadOnly | QIODevice::Text)) {
            std::fprintf(stderr, "tokens-to-css: cannot read %s\n",
                         qPrintable(path));
            return 1;
        }
        const QString current = QString::fromUtf8(file.readAll());
        if (current == wanted) {
            std::printf("%s matches the token table\n", qPrintable(path));
            return 0;
        }
        std::fprintf(stderr,
                     "tokens-to-css: %s no longer matches the token table.\n"
                     "Run the generator rather than editing it:\n"
                     "  cmake --build <build> --target tokens-css\n",
                     qPrintable(path));
        return 1;
    }

    QFile file(path);
    if (!file.open(QIODevice::WriteOnly | QIODevice::Truncate | QIODevice::Text)) {
        std::fprintf(stderr, "tokens-to-css: cannot write %s\n",
                     qPrintable(path));
        return 1;
    }
    file.write(wanted.toUtf8());
    std::printf("wrote %s\n", qPrintable(path));
    return 0;
}
