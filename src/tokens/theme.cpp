// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "theme.h"

#include <QGuiApplication>
#include <QStyleHints>

#include <cmath>

#include "settingsstore.h"
#include "systemappearance.h"
#include "perflog.h"

// The three built-in token tables (features.md §10.1). The light table
// is the application's pre-Phase-9 appearance verbatim; dark and sepia
// were composed against it and are reviewed through the theme
// storyboards like every visual state.
namespace {

Theme::Tokens lightTokens()
{
    Theme::Tokens t;
    t.windowBackground = QColor("#ffffff");
    t.panelBackground = QColor("#f4f4f4");
    t.listBackground = QColor("#fafafa");
    t.footerBackground = QColor("#f5f5f5");
    t.popupBackground = QColor("#ffffff");
    t.chipBackground = QColor("#f2f2f0");
    t.bannerBackground = QColor("#fdf3d8");
    t.codePanelBackground = QColor("#f6f6f4");
    t.textPrimary = QColor("#1a1a1a");
    t.textSecondary = QColor("#555555");
    t.textMuted = QColor("#666666");
    t.textFaint = QColor("#767676");
    t.textDisabled = QColor("#aaaaaa");
    t.bannerText = QColor("#886c1a");
    t.border = QColor("#dddddd");
    t.borderStrong = QColor("#949494");
    t.quoteBar = QColor("#949494");
    t.mutedGlyph = QColor("#7d909f");
    t.hoverTint = QColor("#ebebeb");
    t.blockHoverTint = QColor("#fafafa");
    t.focusTint = QColor("#eaf2fb");
    t.focusRing = QColor("#1f6feb");   // strong keyboard-focus outline
    t.selectionTint = QColor("#dce8f5");
    t.selectionActiveTint = QColor("#c8dff5");
    t.blockSelectionTint = QColor("#dbe9f9");
    t.accent = QColor("#4a90d9");
    t.danger = QColor("#b3261e");
    t.dangerBright = QColor("#e05c5c");
    t.success = QColor("#1e874b");
    t.warning = QColor("#a66908");
    t.pinColor = QColor("#e0a04c");
    t.marker = QColor("#949494");
    t.inlineCodeBackground = QColor("#f0f0ee");
    t.highlightBackground = QColor("#fdf3a9");
    t.link = QColor("#2970c8");
    t.searchMatchBackground = QColor("#b5dcff");
    t.searchCurrentBackground = QColor("#ffb454");
    t.changedTextBackground = QColor("#e4d7f2");
    t.addedTextBackground = QColor("#d3edd8");
    t.removedTextBackground = QColor("#f6d5d2");
    t.codeKeyword = QColor("#a626a4");
    t.codeType = QColor("#2967f0");
    t.codeString = QColor("#3f7e3e");
    t.codeComment = QColor("#6f7178");
    t.codeNumber = QColor("#986801");
    t.calloutTip = QColor("#2a9d8f");
    t.axisAttention = QColor("#d99a3d");
    t.axisAttentionText = QColor("#b06a10");
    t.axisAgent = QColor("#4aa3a3");
    t.axisAgentText = QColor("#1f7a7a");
    t.scopeDiscovered = QColor("#8a5cc0");
    t.signalHard = QColor("#c0392b");
    t.signalSoft = QColor("#d99a3d");
    t.signalHygiene = QColor("#8a5cc0");
    t.hatchAlt = QColor("#ffffff");
    // The three chart ramps; see theme.h and tools/design-ramps.py.
    t.categoricalRamp = {
        "#971c4b", "#cd7707", "#7a6402", "#093c00",
        "#097f95", "#00459a", "#9977e4", "#530054",
    };
    t.sequentialRamp = {
        "#f5aacd", "#e48bb6", "#cf6e9f", "#b75487",
        "#9b3f70", "#7d2e58", "#5e2142",
    };
    t.divergingRamp = {
        "#227702", "#538647", "#799374", "#9e9e9e",
        "#9d7f9a", "#995e94", "#93398e",
    };

    return t;
}

Theme::Tokens darkTokens()
{
    Theme::Tokens t;
    t.windowBackground = QColor("#1e1e1e");
    t.panelBackground = QColor("#252526");
    t.listBackground = QColor("#212122");
    t.footerBackground = QColor("#2a2a2b");
    t.popupBackground = QColor("#2d2d30");
    t.chipBackground = QColor("#37373a");
    t.bannerBackground = QColor("#3a3320");
    t.codePanelBackground = QColor("#262624");
    t.textPrimary = QColor("#eeeeee");
    t.textSecondary = QColor("#c8c8c8");
    t.textMuted = QColor("#a8a8a8");
    t.textFaint = QColor("#858585");
    t.textDisabled = QColor("#5f5f5f");
    t.bannerText = QColor("#d9c37a");
    t.border = QColor("#3c3c3c");
    t.borderStrong = QColor("#696969");
    t.quoteBar = QColor("#696969");
    t.mutedGlyph = QColor("#6f7d8a");
    t.hoverTint = QColor("#333336");
    t.blockHoverTint = QColor("#232324");
    t.focusTint = QColor("#263544");
    t.focusRing = QColor("#58a6ff");   // strong keyboard-focus outline
    t.selectionTint = QColor("#2d4356");
    t.selectionActiveTint = QColor("#35506b");
    t.blockSelectionTint = QColor("#2a3f52");
    t.accent = QColor("#5c9fe0");
    t.danger = QColor("#e06c60");
    t.dangerBright = QColor("#e05c5c");
    t.success = QColor("#5abd82");
    t.warning = QColor("#e0a34c");
    t.pinColor = QColor("#e0a04c");
    t.marker = QColor("#6f6f6f");
    t.inlineCodeBackground = QColor("#333330");
    t.highlightBackground = QColor("#6b5c17");
    t.link = QColor("#6fb1ff");
    t.searchMatchBackground = QColor("#264f78");
    t.searchCurrentBackground = QColor("#96601f");
    t.changedTextBackground = QColor("#3d3060");
    t.addedTextBackground = QColor("#1d3a26");
    t.removedTextBackground = QColor("#45211f");
    t.codeKeyword = QColor("#c678dd");
    t.codeType = QColor("#61afef");
    t.codeString = QColor("#98c379");
    t.codeComment = QColor("#888c96");
    t.codeNumber = QColor("#d19a66");
    t.calloutTip = QColor("#4db6ac");
    t.axisAttention = QColor("#d9a04c");
    t.axisAttentionText = QColor("#e6b877");
    t.axisAgent = QColor("#4aa3a3");
    t.axisAgentText = QColor("#7cc7c4");
    t.scopeDiscovered = QColor("#a37fd4");
    t.signalHard = QColor("#e06c60");
    t.signalSoft = QColor("#e0a34c");
    t.signalHygiene = QColor("#a37fd4");
    t.hatchAlt = QColor("#1e1e1e");
    // The three chart ramps; see theme.h and tools/design-ramps.py.
    t.categoricalRamp = {
        "#ff8e93", "#ffc163", "#8b8800", "#acffbd",
        "#74deff", "#4081e1", "#ad9efe", "#cc5b9a",
    };
    t.sequentialRamp = {
        "#7c3c5d", "#9d4b75", "#bc5d8d", "#d973a6",
        "#f28dbe", "#feb0d4", "#ffd8e9",
    };
    t.divergingRamp = {
        "#89dd76", "#8ac07f", "#89a384", "#868686",
        "#ae8faa", "#d596cf", "#fd9cf5",
    };

    return t;
}

Theme::Tokens sepiaTokens()
{
    Theme::Tokens t;
    t.windowBackground = QColor("#f6efdf");
    t.panelBackground = QColor("#eee5d0");
    t.listBackground = QColor("#f2ead8");
    t.footerBackground = QColor("#ece2cc");
    t.popupBackground = QColor("#faf4e6");
    t.chipBackground = QColor("#ece3cb");
    t.bannerBackground = QColor("#eadfb2");
    t.codePanelBackground = QColor("#efe7d1");
    t.textPrimary = QColor("#3d3427");
    t.textSecondary = QColor("#5a4f3d");
    t.textMuted = QColor("#6e6250");
    t.textFaint = QColor("#776b59");
    t.textDisabled = QColor("#a99a80");
    t.bannerText = QColor("#776017");
    t.border = QColor("#ded2b8");
    t.borderStrong = QColor("#9c875d");
    t.quoteBar = QColor("#9f8756");
    t.mutedGlyph = QColor("#938166");
    t.hoverTint = QColor("#e8dfc9");
    t.blockHoverTint = QColor("#f0e8d6");
    t.focusTint = QColor("#ece0c4");
    t.focusRing = QColor("#9a6a2b");   // strong keyboard-focus outline
    t.selectionTint = QColor("#e2d4b4");
    t.selectionActiveTint = QColor("#d8c69e");
    t.blockSelectionTint = QColor("#e0d2b2");
    t.accent = QColor("#9a6b2f");
    t.danger = QColor("#a33b30");
    t.dangerBright = QColor("#c25a4e");
    t.success = QColor("#4d7839");
    t.warning = QColor("#91641a");
    t.pinColor = QColor("#b07a2a");
    t.marker = QColor("#9a8863");
    t.inlineCodeBackground = QColor("#ece3cb");
    t.highlightBackground = QColor("#f0dd8c");
    t.link = QColor("#8a5a20");
    t.searchMatchBackground = QColor("#cfd9a8");
    t.searchCurrentBackground = QColor("#e8a94e");
    t.changedTextBackground = QColor("#ded0e0");
    t.addedTextBackground = QColor("#d5e6c9");
    t.removedTextBackground = QColor("#f2d6c8");
    t.codeKeyword = QColor("#9a2f8a");
    t.codeType = QColor("#2f69af");
    t.codeString = QColor("#497236");
    t.codeComment = QColor("#73674e");
    t.codeNumber = QColor("#8a5a20");
    t.calloutTip = QColor("#3a8a7a");
    t.axisAttention = QColor("#c08a2e");
    t.axisAttentionText = QColor("#8a5f14");
    t.axisAgent = QColor("#3f8f8a");
    t.axisAgentText = QColor("#2a6b66");
    t.scopeDiscovered = QColor("#7a4fa8");
    t.signalHard = QColor("#a33b30");
    t.signalSoft = QColor("#b07a1f");
    t.signalHygiene = QColor("#7a4fa8");
    t.hatchAlt = QColor("#f6efdf");
    // The three chart ramps; see theme.h and tools/design-ramps.py.
    t.categoricalRamp = {
        "#b93f67", "#612700", "#838303", "#1c6701",
        "#0088bb", "#2063bf", "#410371", "#892f84",
    };
    t.sequentialRamp = {
        "#f5aacd", "#e48bb6", "#cf6e9f", "#b75487",
        "#9b3f70", "#7d2e58", "#5e2142",
    };
    t.divergingRamp = {
        "#227702", "#538647", "#799374", "#9e9e9e",
        "#9d7f9a", "#995e94", "#93398e",
    };

    return t;
}

// A colour override the painter can actually use, or nothing. The empty
// string is how both overrides say "no override", so an unusable value and
// an absent one collapse to the same state.
QString validColorOrEmpty(const QString &value)
{
    return QColor::isValidColorName(value) ? value : QString();
}

// The darkest label the derived accent text is allowed to be. Pure black on a
// mid accent is harsher than the rest of the interface, and near-black clears
// every default accent by a wide margin anyway.
const QColor kNearBlack("#1a1a1a");

// WCAG 2.1 relative luminance: each channel linearised, then weighted by how
// much the eye takes from it.
qreal relativeLuminance(const QColor &c)
{
    auto channel = [](qreal v) {
        return v <= 0.04045 ? v / 12.92 : std::pow((v + 0.055) / 1.055, 2.4);
    };
    return 0.2126 * channel(c.redF()) + 0.7152 * channel(c.greenF())
         + 0.0722 * channel(c.blueF());
}

const QString kSettingsThemeId = QStringLiteral("theme.id");
const QString kSettingsAccent = QStringLiteral("theme.accent");
const QString kSettingsHighlight = QStringLiteral("theme.highlight");
// The pre-Finding-7 key: a plain bool, present only where somebody had made
// a choice. It is still read, once, to carry an existing installation's
// explicit setting into the three-way one.
const QString kSettingsReducedMotion = QStringLiteral("view.reducedMotion");
const QString kSettingsReducedMotionMode =
    QStringLiteral("view.reducedMotionMode");

} // namespace

Theme::Theme(QObject *parent)
    : QObject(parent)
{
    refresh();

    // "system" follows the OS setting live (§10.1). qGuiApp is null in
    // pure-QObject unit tests; those exercise the fixed themes.
    if (qGuiApp) {
        connect(qGuiApp->styleHints(), &QStyleHints::colorSchemeChanged,
                this, [this]() {
                    if (m_themeId == QLatin1String("system"))
                        refresh();
                });
    }
}

void Theme::setSettings(SettingsStore *settings)
{
    m_settings = settings;
    if (!m_settings)
        return;
    // First start (no persisted choice) follows the OS color scheme; an
    // explicit pick in Settings persists and wins from then on.
    m_themeId = m_settings->value(kSettingsThemeId,
                                  QStringLiteral("system")).toString();
    if (!availableThemes().contains(m_themeId))
        m_themeId = QStringLiteral("system");  // stale or hand-edited
    // Same validation the live setters apply. A hand-edited or corrupted
    // settings file would otherwise install a string QColor cannot parse,
    // and refresh() would paint the accent with an invalid colour instead
    // of falling back to the theme's own token.
    m_accentOverride = validColorOrEmpty(
        m_settings->value(kSettingsAccent).toString());
    m_highlightOverride = validColorOrEmpty(
        m_settings->value(kSettingsHighlight).toString());
    // Reduced motion, in the order that keeps an existing installation's
    // choice: the three-way setting if it has been written, else the old
    // bool — whose presence at all means somebody chose — else "system",
    // which is the default for a new installation.
    const QString previous = m_reducedMotionSetting;
    if (m_settings->contains(kSettingsReducedMotionMode)) {
        m_reducedMotionSetting =
            m_settings->value(kSettingsReducedMotionMode).toString();
    } else if (m_settings->contains(kSettingsReducedMotion)) {
        m_reducedMotionSetting =
            m_settings->value(kSettingsReducedMotion).toBool()
                ? QStringLiteral("on") : QStringLiteral("off");
    } else {
        m_reducedMotionSetting = QStringLiteral("system");
    }
    if (!availableReducedMotionSettings().contains(m_reducedMotionSetting))
        m_reducedMotionSetting = QStringLiteral("system");   // hand-edited
    if (previous != m_reducedMotionSetting)
        emit reducedMotionChanged();
    refresh();
}

QStringList Theme::availableReducedMotionSettings()
{
    return { QStringLiteral("on"), QStringLiteral("off"),
             QStringLiteral("system") };
}

bool Theme::reducedMotion() const
{
    if (m_reducedMotionSetting == QLatin1String("on"))
        return true;
    if (m_reducedMotionSetting == QLatin1String("off"))
        return false;
    // "system": what the desktop says, and off where it says nothing.
    return m_systemAppearance && m_systemAppearance->reducedMotion();
}

void Theme::setReducedMotion(bool reduced)
{
    // Writing the effective value is an explicit choice, which is what the
    // View menu's checkable item means when it is toggled.
    setReducedMotionSetting(reduced ? QStringLiteral("on")
                                    : QStringLiteral("off"));
}

void Theme::setReducedMotionSetting(const QString &mode)
{
    if (m_reducedMotionSetting == mode
        || !availableReducedMotionSettings().contains(mode))
        return;
    m_reducedMotionSetting = mode;
    if (m_settings)
        m_settings->setValue(kSettingsReducedMotionMode, mode);
    emit reducedMotionChanged();
}

void Theme::setSystemAppearance(SystemAppearance *appearance)
{
    if (m_systemAppearance == appearance)
        return;
    if (m_systemAppearance)
        disconnect(m_systemAppearance, nullptr, this, nullptr);
    m_systemAppearance = appearance;
    if (m_systemAppearance) {
        connect(m_systemAppearance, &SystemAppearance::changed, this, [this]() {
            // Both halves ride on this one signal: the theme may have to
            // resolve differently, and the effective motion value may have
            // moved. Emitting both is cheaper than working out which.
            emit reducedMotionChanged();
            refresh();
        });
    }
    emit reducedMotionChanged();
    refresh();
}

QString Theme::displayName(const QString &themeId) const
{
    if (themeId == QLatin1String("highContrast"))
        return tr("High contrast");
    if (themeId.isEmpty())
        return themeId;
    return themeId.at(0).toUpper() + themeId.mid(1);
}

QStringList Theme::availableThemes() const
{
    return { QStringLiteral("light"), QStringLiteral("dark"),
             QStringLiteral("sepia"), QStringLiteral("highContrast"),
             QStringLiteral("system") };
}

void Theme::setThemeId(const QString &themeId)
{
    if (m_themeId == themeId || !availableThemes().contains(themeId))
        return;
    PerfLog::ScopedTimer perf(
        QStringLiteral("theme.switch"),
        QVariantMap{{QStringLiteral("theme"), themeId}});
    m_themeId = themeId;
    if (m_settings)
        m_settings->setValue(kSettingsThemeId, m_themeId);
    refresh();
}

void Theme::setAccentOverride(const QString &hex)
{
    const QString value = (QColor::isValidColorName(hex)) ? hex : QString();
    if (m_accentOverride == value)
        return;
    PerfLog::ScopedTimer perf(
        QStringLiteral("theme.switch"),
        QVariantMap{{QStringLiteral("accentOverride"), value}});
    m_accentOverride = value;
    if (m_settings)
        m_settings->setValue(kSettingsAccent, m_accentOverride);
    refresh();
}

void Theme::setHighlightOverride(const QString &hex)
{
    const QString value = (QColor::isValidColorName(hex)) ? hex : QString();
    if (m_highlightOverride == value)
        return;
    PerfLog::ScopedTimer perf(
        QStringLiteral("theme.switch"),
        QVariantMap{{QStringLiteral("highlightOverride"), value}});
    m_highlightOverride = value;
    if (m_settings)
        m_settings->setValue(kSettingsHighlight, m_highlightOverride);
    refresh();
}

QColor Theme::categorical(int index) const
{
    const QStringList &ramp = m_tokens.categoricalRamp;
    if (ramp.isEmpty())
        return m_tokens.textMuted;
    // A negative index is a caller doing its own arithmetic and getting it
    // wrong; wrapping it rather than clamping keeps the series colours stable
    // when it does.
    const int at = ((index % ramp.size()) + ramp.size()) % ramp.size();
    return QColor(ramp.at(at));
}

QStringList Theme::reservedHues() const
{
    // Every colour in this theme that already means something. The chart-ramp
    // validator keeps its steps away from these; the list is here rather than
    // in the test so that a token added to the table above is a token the
    // validator knows about, without anybody remembering to add it twice.
    //
    // Surfaces, text and lines are deliberately absent: a chart step is
    // already held away from a surface by the contrast rule, and no reader
    // mistakes a bar for the colour of the text beside it.
    //
    // `axisAttentionText` and `axisAgentText` are absent for the same reason
    // and one more. They are the legible text partners of the two axis bar
    // hues — the same colour, adjusted until small type set in it clears 4.5:1
    // — and they sit close enough to their bar hues that reserving both closes
    // off enough of the hue wheel to leave no categorical ramp of eight steps
    // at all. They are also never drawn as a mark, so nothing can be confused
    // with them.
    return {
        m_tokens.accent.name(), m_tokens.danger.name(),
        m_tokens.dangerBright.name(), m_tokens.success.name(),
        m_tokens.warning.name(), m_tokens.link.name(),
        m_tokens.pinColor.name(), m_tokens.calloutTip.name(),
        m_tokens.focusRing.name(),
        m_tokens.axisAttention.name(), m_tokens.axisAgent.name(),
        m_tokens.scopeDiscovered.name(), m_tokens.signalHard.name(),
        m_tokens.signalSoft.name(), m_tokens.signalHygiene.name(),
    };
}

QStringList Theme::colorPalette() const
{
    return { QStringLiteral("#e05c5c"), QStringLiteral("#e0a04c"),
             QStringLiteral("#58a866"), QStringLiteral("#4a90d9"),
             QStringLiteral("#9068c8"), QStringLiteral("#d06ca8") };
}

QStringList Theme::highlightPalette() const
{
    return { QStringLiteral("#fdf3a9"), QStringLiteral("#ffd9a8"),
             QStringLiteral("#c9ecc9"), QStringLiteral("#c9e4ff"),
             QStringLiteral("#f2ccf2") };
}

QStringList Theme::colorPaletteNames() const
{
    return { tr("Red"), tr("Amber"), tr("Green"),
             tr("Blue"), tr("Purple"), tr("Pink") };
}

QStringList Theme::highlightPaletteNames() const
{
    return { tr("Yellow"), tr("Peach"), tr("Mint"), tr("Sky"), tr("Lilac") };
}

QString Theme::colorName(const QString &value) const
{
    if (value.trimmed().isEmpty())
        return tr("Theme default");
    // Compared as parsed colours rather than as strings: the same swatch
    // arrives as "#e05c5c" from the palette and as "#E05C5C" from a document
    // that was hand-edited, and a name that depends on the spelling would
    // silently fall through to the hex fallback for the second one.
    const QColor wanted(value);
    if (!wanted.isValid())
        return value;

    const QStringList palettes[] = { colorPalette(), highlightPalette() };
    const QStringList names[] = { colorPaletteNames(), highlightPaletteNames() };
    for (int p = 0; p < 2; ++p) {
        for (int i = 0; i < palettes[p].size(); ++i) {
            if (QColor(palettes[p].at(i)) == wanted)
                return names[p].value(i, value);
        }
    }
    // The two greys the text-colour picker adds to the content palette for
    // prose; they are not folder or tag choices, so they are not in either
    // list above.
    if (wanted == QColor("#333333"))
        return tr("Near black");
    if (wanted == QColor("#888888"))
        return tr("Grey");
    return tr("Custom colour %1").arg(wanted.name());
}

// High-contrast theme (§14.3): pure black ground, white body text (21:1, WCAG
// AAA), bright saturated accents, and strong white borders so every structural
// edge is visible. A bright-yellow focus ring gives the maximum-contrast
// keyboard indicator. Its stricter contrast floor is asserted in the theme test.
Theme::Tokens highContrastTokens()
{
    Theme::Tokens t;
    t.windowBackground = QColor("#000000");
    t.panelBackground = QColor("#000000");
    t.listBackground = QColor("#050505");
    t.footerBackground = QColor("#000000");
    t.popupBackground = QColor("#0a0a0a");
    t.chipBackground = QColor("#1a1a1a");
    t.bannerBackground = QColor("#1a1a00");
    t.codePanelBackground = QColor("#0a0a0a");
    t.textPrimary = QColor("#ffffff");
    t.textSecondary = QColor("#f0f0f0");
    t.textMuted = QColor("#d8d8d8");
    t.textFaint = QColor("#bcbcbc");
    t.textDisabled = QColor("#8a8a8a");
    t.bannerText = QColor("#ffff00");
    t.border = QColor("#ffffff");
    t.borderStrong = QColor("#ffffff");
    t.quoteBar = QColor("#ffffff");
    t.mutedGlyph = QColor("#ffff66");
    t.hoverTint = QColor("#333333");
    t.blockHoverTint = QColor("#1a1a1a");
    t.focusTint = QColor("#00335c");
    t.focusRing = QColor("#ffff00");   // maximum-contrast keyboard outline
    t.selectionTint = QColor("#0055aa");
    t.selectionActiveTint = QColor("#0070d0");
    t.blockSelectionTint = QColor("#004488");
    t.accent = QColor("#33ccff");
    t.danger = QColor("#ff6666");
    t.dangerBright = QColor("#ff8080");
    t.success = QColor("#44ff99");
    t.warning = QColor("#ffbb33");
    t.pinColor = QColor("#ffbb33");
    t.marker = QColor("#ffff66");
    t.inlineCodeBackground = QColor("#222222");
    t.highlightBackground = QColor("#5c5c00");
    t.link = QColor("#66ddff");
    t.searchMatchBackground = QColor("#0055aa");
    t.searchCurrentBackground = QColor("#995a00");
    t.changedTextBackground = QColor("#4b0082");
    t.addedTextBackground = QColor("#00420f");
    t.removedTextBackground = QColor("#5c0000");
    t.codeKeyword = QColor("#ff99ff");
    t.codeType = QColor("#99ddff");
    t.codeString = QColor("#99ff99");
    t.codeComment = QColor("#d0d0d0");
    t.codeNumber = QColor("#ffcc88");
    t.calloutTip = QColor("#33ffdd");
    // The soft-signal amber matches the attention hue deliberately: a soft
    // signal always sits inside a chip, so the two never compete on one mark.
    t.axisAttention = QColor("#ffbb33");
    t.axisAttentionText = QColor("#ffcc66");
    t.axisAgent = QColor("#33ddcc");
    t.axisAgentText = QColor("#77ffe8");
    t.scopeDiscovered = QColor("#cc99ff");
    t.signalHard = QColor("#ff6666");
    t.signalSoft = QColor("#ffbb33");
    t.signalHygiene = QColor("#cc99ff");
    t.hatchAlt = QColor("#000000");
    // The three chart ramps; see theme.h and tools/design-ramps.py.
    t.categoricalRamp = {
        "#c0453f", "#d38406", "#e4ffa6", "#89dd76",
        "#10919d", "#5ea7ff", "#795ec9", "#fec0f4",
    };
    t.sequentialRamp = {
        "#7c3c5d", "#9d4b75", "#bc5d8d", "#d973a6",
        "#f28dbe", "#feb0d4", "#ffd8e9",
    };
    t.divergingRamp = {
        "#89dd76", "#8ac07f", "#89a384", "#868686",
        "#ae8faa", "#d596cf", "#fd9cf5",
    };

    return t;
}

qreal Theme::contrastRatio(const QColor &a, const QColor &b)
{
    const qreal la = relativeLuminance(a);
    const qreal lb = relativeLuminance(b);
    const qreal lighter = qMax(la, lb);
    const qreal darker = qMin(la, lb);
    return (lighter + 0.05) / (darker + 0.05);
}

QColor Theme::labelOn(const QColor &fill)
{
    // Two candidates rather than a computed shade: a label has to be one of
    // the two colours the rest of the interface already uses for text on a
    // fill, and anything between them reads as a third text colour nobody
    // asked for.
    return contrastRatio(kNearBlack, fill) >= contrastRatio(Qt::white, fill)
        ? kNearBlack : QColor(Qt::white);
}

// The token tables carry no onAccent of their own: it is derived here so a
// built-in theme and a custom accent go through exactly one rule. tokensFor()
// is what the settings previews and the theme test read, so the derivation
// has to happen before they see the table rather than only in refresh().
static Theme::Tokens withDerivedLabels(Theme::Tokens t)
{
    t.onAccent = Theme::labelOn(t.accent);
    return t;
}

const Theme::Tokens &Theme::tokensFor(const QString &resolvedTheme)
{
    static const Tokens light = withDerivedLabels(lightTokens());
    static const Tokens dark = withDerivedLabels(darkTokens());
    static const Tokens sepia = withDerivedLabels(sepiaTokens());
    static const Tokens highContrast = withDerivedLabels(highContrastTokens());
    if (resolvedTheme == QLatin1String("dark"))
        return dark;
    if (resolvedTheme == QLatin1String("sepia"))
        return sepia;
    if (resolvedTheme == QLatin1String("highContrast"))
        return highContrast;
    return light;
}

QVariantMap Theme::themePreview(const QString &themeId) const
{
    const QString resolved = (themeId == QLatin1String("system"))
        ? resolveSystem() : themeId;
    const Tokens &t = tokensFor(resolved);
    return { { QStringLiteral("background"), t.windowBackground },
             { QStringLiteral("panel"), t.panelBackground },
             { QStringLiteral("text"), t.textPrimary },
             { QStringLiteral("accent"), t.accent },
             { QStringLiteral("border"), t.border } };
}

QString Theme::resolveSystem() const
{
    // High contrast first: someone who has turned it on system-wide has asked
    // for it in a way that outranks light-or-dark, and the high-contrast
    // theme answers both questions at once (accessibility.md Finding 7).
    // Only under "system" — an explicit theme choice is the user saying what
    // they want from this application in particular, and it wins.
    if (m_systemAppearance && m_systemAppearance->highContrast())
        return QStringLiteral("highContrast");
    if (qGuiApp
        && qGuiApp->styleHints()->colorScheme() == Qt::ColorScheme::Dark)
        return QStringLiteral("dark");
    return QStringLiteral("light");  // Light and Unknown both render light
}

void Theme::refresh()
{
    m_resolved = (m_themeId == QLatin1String("system")) ? resolveSystem()
                                                        : m_themeId;
    m_tokens = tokensFor(m_resolved);

    // §10.3 overrides ride on top of the table; the engine and QML
    // read plain tokens and never know overrides exist.
    if (!m_accentOverride.isEmpty()) {
        const QColor accent(m_accentOverride);
        m_tokens.accent = accent;
        m_tokens.link = accent;
        // The label follows the accent it sits on, whether that accent came
        // from the table or from the user (accessibility.md Finding 3).
        m_tokens.onAccent = labelOn(accent);
    }
    if (!m_highlightOverride.isEmpty())
        m_tokens.highlightBackground = QColor(m_highlightOverride);

    emit themeChanged();
}
