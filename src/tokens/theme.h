// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef THEME_H
#define THEME_H

#include <QColor>
#include <QObject>
#include <QString>
#include <QStringList>
#include <QVariantMap>

class SettingsStore;
class SystemAppearance;

// The theme token object: one C++ QObject
// exposing every color the application draws with, as a Q_PROPERTY per
// token. QML binds to the tokens — no hardcoded colors in new QML —
// and the engine's highlighter reads the inline
// tokens from the same object, so both layers restyle
// from one source. All tokens notify through the single themeChanged
// signal: a theme switch re-evaluates every binding wholesale, which
// is the intended (user-rare) whole-window repaint.
//
// Built-in themes are rows in one static table: "light", "dark",
// "sepia". "system" resolves to light or dark from the OS setting
// (QStyleHints::colorScheme, followed live). The chosen theme id and
// the §10.3 overrides (accent, highlight) live in the settings store.
class Theme : public QObject
{
    Q_OBJECT

    // "light" | "dark" | "sepia" | "highContrast" | "system" (persisted).
    Q_PROPERTY(QString themeId READ themeId WRITE setThemeId NOTIFY themeChanged)
    // Reduced motion (§14.3): one source every animation and the typewriter
    // scroll reads. motionScale is 0 when on (all motion stilled), 1 when off,
    // so a duration written `150 * theme.motionScale` becomes instant.
    //
    // This is the effective value, and writing it is still the way to make an
    // explicit choice — it sets reducedMotionSetting to "on" or "off". What
    // it can no longer be is the whole story, because "follow the system" is
    // now a third state (accessibility.md Finding 7).
    Q_PROPERTY(bool reducedMotion READ reducedMotion WRITE setReducedMotion
                   NOTIFY reducedMotionChanged)
    // "on" | "off" | "system" (persisted). "system" is the default for a new
    // installation; an existing one keeps whatever explicit value it had.
    Q_PROPERTY(QString reducedMotionSetting READ reducedMotionSetting
                   WRITE setReducedMotionSetting NOTIFY reducedMotionChanged)
    Q_PROPERTY(qreal motionScale READ motionScale NOTIFY reducedMotionChanged)
    // What themeId currently renders as: never "system".
    Q_PROPERTY(QString resolvedTheme READ resolvedTheme NOTIFY themeChanged)
    Q_PROPERTY(QStringList availableThemes READ availableThemes CONSTANT)

    // §10.3 customization: "" means the theme's own value.
    Q_PROPERTY(QString accentOverride READ accentOverride
                   WRITE setAccentOverride NOTIFY themeChanged)
    Q_PROPERTY(QString highlightOverride READ highlightOverride
                   WRITE setHighlightOverride NOTIFY themeChanged)

    // Surfaces.
    Q_PROPERTY(QColor windowBackground READ windowBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor panelBackground READ panelBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor listBackground READ listBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor footerBackground READ footerBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor popupBackground READ popupBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor chipBackground READ chipBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor bannerBackground READ bannerBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor codePanelBackground READ codePanelBackground NOTIFY themeChanged)

    // Text.
    Q_PROPERTY(QColor textPrimary READ textPrimary NOTIFY themeChanged)
    Q_PROPERTY(QColor textSecondary READ textSecondary NOTIFY themeChanged)
    Q_PROPERTY(QColor textMuted READ textMuted NOTIFY themeChanged)
    Q_PROPERTY(QColor textFaint READ textFaint NOTIFY themeChanged)
    Q_PROPERTY(QColor textDisabled READ textDisabled NOTIFY themeChanged)
    Q_PROPERTY(QColor bannerText READ bannerText NOTIFY themeChanged)
    // The label colour on an accent-filled control. Derived from the
    // effective accent rather than stored, because the accent is
    // user-overridable (§10.3) and this was not: a pale custom accent gave
    // white text on a pale fill with nothing anywhere to catch it. See
    // labelOn(), and use that directly for a fill that is NOT the accent —
    // a danger button, a tag colour — since this one answers only for the
    // accent (accessibility.md Finding 3).
    Q_PROPERTY(QColor onAccent READ onAccent NOTIFY themeChanged)

    // Lines and glyphs.
    //
    // The two border tokens are not "the same line, lighter and darker".
    // `border` is decorative: the rule between two panels, a separator, a
    // block frame. WCAG asks nothing of it, and it is deliberately below 3:1
    // in every theme. `borderStrong` is the control-boundary token — the
    // outline that says where a button, field, checkbox or swatch is — and
    // it is held to 3:1 against the surface behind it by the contrast floors
    // in tests/test_theme.cpp. A new control's resting outline takes
    // `borderStrong`; using `border` for it is the mistake accessibility.md
    // Finding 3 describes.
    Q_PROPERTY(QColor border READ border NOTIFY themeChanged)
    Q_PROPERTY(QColor borderStrong READ borderStrong NOTIFY themeChanged)
    Q_PROPERTY(QColor quoteBar READ quoteBar NOTIFY themeChanged)
    Q_PROPERTY(QColor mutedGlyph READ mutedGlyph NOTIFY themeChanged)

    // Interactive tints.
    Q_PROPERTY(QColor hoverTint READ hoverTint NOTIFY themeChanged)
    Q_PROPERTY(QColor blockHoverTint READ blockHoverTint NOTIFY themeChanged)
    Q_PROPERTY(QColor focusTint READ focusTint NOTIFY themeChanged)
    Q_PROPERTY(QColor focusRing READ focusRing NOTIFY themeChanged)
    Q_PROPERTY(QColor selectionTint READ selectionTint NOTIFY themeChanged)
    Q_PROPERTY(QColor selectionActiveTint READ selectionActiveTint NOTIFY themeChanged)
    Q_PROPERTY(QColor blockSelectionTint READ blockSelectionTint NOTIFY themeChanged)

    // Accent and semantic colors.
    Q_PROPERTY(QColor accent READ accent NOTIFY themeChanged)
    Q_PROPERTY(QColor danger READ danger NOTIFY themeChanged)
    Q_PROPERTY(QColor dangerBright READ dangerBright NOTIFY themeChanged)
    Q_PROPERTY(QColor success READ success NOTIFY themeChanged)
    Q_PROPERTY(QColor warning READ warning NOTIFY themeChanged)
    Q_PROPERTY(QColor pinColor READ pinColor NOTIFY themeChanged)

    // Inline styling (read by the engine's highlighter).
    Q_PROPERTY(QColor marker READ marker NOTIFY themeChanged)
    Q_PROPERTY(QColor inlineCodeBackground READ inlineCodeBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor highlightBackground READ highlightBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor link READ link NOTIFY themeChanged)
    Q_PROPERTY(QColor searchMatchBackground READ searchMatchBackground NOTIFY themeChanged)
    Q_PROPERTY(QColor searchCurrentBackground READ searchCurrentBackground NOTIFY themeChanged)
    // The wash behind text that differs from another version of the same
    // document — what the backup dialog paints over the part of a stored
    // version the note in front of the reader no longer has. Its own token
    // rather than a reuse of the search tint: the two answer different
    // questions, and a reader who has both in their recent memory should not
    // have to work out which mark they are looking at.
    Q_PROPERTY(QColor changedTextBackground READ changedTextBackground NOTIFY themeChanged)

    // Code-block syntax highlighting: the five
    // token colors the code highlighter paints, read the same way the inline
    // tokens above are. codeType is the "function/type" color.
    Q_PROPERTY(QColor codeKeyword READ codeKeyword NOTIFY themeChanged)
    Q_PROPERTY(QColor codeType READ codeType NOTIFY themeChanged)
    Q_PROPERTY(QColor codeString READ codeString NOTIFY themeChanged)
    Q_PROPERTY(QColor codeComment READ codeComment NOTIFY themeChanged)
    Q_PROPERTY(QColor codeNumber READ codeNumber NOTIFY themeChanged)

    // Callout tip accent: the five other callout types reuse
    // accent/warning/success/danger/textMuted; tip needs its own hue.
    Q_PROPERTY(QColor calloutTip READ calloutTip NOTIFY themeChanged)

    // Portfolio-dashboard vocabulary, for a separate application that links
    // this library and draws project dashboards. The two effort axes own two
    // hues used for nothing else; discovered scope and hygiene signals share
    // violet. They are defined here so that application and the editor stay
    // one token system; the editor itself does not draw with them.
    Q_PROPERTY(QColor axisAttention READ axisAttention NOTIFY themeChanged)
    Q_PROPERTY(QColor axisAttentionText READ axisAttentionText NOTIFY themeChanged)
    Q_PROPERTY(QColor axisAgent READ axisAgent NOTIFY themeChanged)
    Q_PROPERTY(QColor axisAgentText READ axisAgentText NOTIFY themeChanged)
    Q_PROPERTY(QColor scopeDiscovered READ scopeDiscovered NOTIFY themeChanged)
    Q_PROPERTY(QColor signalHard READ signalHard NOTIFY themeChanged)
    Q_PROPERTY(QColor signalSoft READ signalSoft NOTIFY themeChanged)
    Q_PROPERTY(QColor signalHygiene READ signalHygiene NOTIFY themeChanged)
    // The stroke drawn across a hatched (bounded) bar; per theme, so the hatch
    // survives every ground the bar can sit on.
    Q_PROPERTY(QColor hatchAlt READ hatchAlt NOTIFY themeChanged)

    // ── The three chart ramps ──────────────────────────────────────────────
    //
    // prd.md §5.5 calls the absence of these the largest single gap the review
    // found, and the reason is that every saturated hue above is already
    // spoken for: `signalSoft` and `warning` are the same value, so are
    // `signalHygiene` and `scopeDiscovered`, and so are `signalHard` and
    // `danger`. A chart drawn in any of them says something it does not mean,
    // and kvit-cash's spending categories can borrow none of them.
    //
    // Three ramps, because charts ask three different questions. A
    // *categorical* ramp says which series this is, and its steps must be
    // identifiable rather than ordered. A *sequential* ramp says how much, and
    // its steps must be rankable — including in grayscale, which is why it
    // moves through lightness. A *diverging* ramp says how far from a
    // reference and in which direction, so it has a neutral midpoint and two
    // arms.
    //
    // Every step is checked rather than asserted: tests/test_palette.cpp runs
    // each ramp against this theme's surfaces and this theme's reserved hues
    // through KvitPalette, so editing a value below fails the build rather
    // than the eye. tools/design-ramps.py is what produced them and is how to
    // produce new ones if a rule changes.
    //
    // The hues are the same in all four themes, within about 14 degrees. What
    // changes per theme is how light and how saturated a hue has to be to sit
    // on that theme's ground. Four independently chosen palettes would mean
    // that switching from light to dark recolours every category on the
    // screen and the reader has to learn the legend again.
    Q_PROPERTY(QStringList categoricalRamp READ categoricalRamp NOTIFY themeChanged)
    Q_PROPERTY(QStringList sequentialRamp READ sequentialRamp NOTIFY themeChanged)
    Q_PROPERTY(QStringList divergingRamp READ divergingRamp NOTIFY themeChanged)

    // User-data colors (folder and tag palettes): content, not chrome,
    // so identical in every theme.
    Q_PROPERTY(QStringList colorPalette READ colorPalette CONSTANT)
    // Soft span tints offered by the highlight-color picker (§10.3):
    // text sits ON these, so the saturated folder palette above would
    // be wrong. Choice candidates, not chrome — theme-independent.
    Q_PROPERTY(QStringList highlightPalette READ highlightPalette CONSTANT)
    // What each of those swatches is called, in the same order. A screen
    // reader has to say something when the caret reaches a swatch, and a
    // six-digit hex value is not it (accessibility.md Finding 2). The two
    // lists are kept the same length as the palettes above by
    // paletteNamesCoverEveryColor in tests/test_theme.cpp.
    Q_PROPERTY(QStringList colorPaletteNames READ colorPaletteNames CONSTANT)
    Q_PROPERTY(QStringList highlightPaletteNames READ highlightPaletteNames
                   CONSTANT)

public:
    struct Tokens {
        QColor windowBackground, panelBackground, listBackground,
            footerBackground, popupBackground, chipBackground,
            bannerBackground, codePanelBackground;
        // onAccent is the one token no table sets: labelOn(accent) fills it
        // in, so a built-in accent and a user's custom one go through the
        // same rule.
        QColor textPrimary, textSecondary, textMuted, textFaint,
            textDisabled, bannerText, onAccent;
        QColor border, borderStrong, quoteBar, mutedGlyph;
        QColor hoverTint, blockHoverTint, focusTint, selectionTint,
            selectionActiveTint, blockSelectionTint, focusRing;
        QColor accent, danger, dangerBright, success, warning, pinColor;
        QColor marker, inlineCodeBackground, highlightBackground, link,
            searchMatchBackground, searchCurrentBackground,
            changedTextBackground;
        QColor codeKeyword, codeType, codeString, codeComment, codeNumber;
        QColor calloutTip;
        QColor axisAttention, axisAttentionText, axisAgent, axisAgentText,
            scopeDiscovered, signalHard, signalSoft, signalHygiene, hatchAlt;
        // Eight categorical steps, seven sequential, seven diverging. See the
        // property declarations above for what each is for.
        QStringList categoricalRamp, sequentialRamp, divergingRamp;
    };

    explicit Theme(QObject *parent = nullptr);

    // Reads theme.id / theme.accent / theme.highlight and saves them
    // back on change. Optional: without a store the theme is
    // session-scoped (tests).
    void setSettings(SettingsStore *settings);

    QString themeId() const { return m_themeId; }
    void setThemeId(const QString &themeId);
    QString resolvedTheme() const { return m_resolved; }
    QStringList availableThemes() const;

    bool reducedMotion() const;
    void setReducedMotion(bool reduced);
    QString reducedMotionSetting() const { return m_reducedMotionSetting; }
    void setReducedMotionSetting(const QString &mode);
    // The three values reducedMotionSetting accepts, for the settings UI and
    // for rejecting a hand-edited settings file.
    Q_INVOKABLE static QStringList availableReducedMotionSettings();
    qreal motionScale() const { return reducedMotion() ? 0.0 : 1.0; }

    // What the desktop says about high contrast and animation. Optional:
    // without one, "system" resolves to light or dark only and "follow the
    // system" for motion means off. Attached once by ProcessServices, and by
    // a test that wants to drive both from a stand-in.
    void setSystemAppearance(SystemAppearance *appearance);
    // Friendly label for a theme id (the menu uses it); handles camelCase ids.
    Q_INVOKABLE QString displayName(const QString &themeId) const;

    QString accentOverride() const { return m_accentOverride; }
    void setAccentOverride(const QString &hex);
    QString highlightOverride() const { return m_highlightOverride; }
    void setHighlightOverride(const QString &hex);

    QColor windowBackground() const { return m_tokens.windowBackground; }
    QColor panelBackground() const { return m_tokens.panelBackground; }
    QColor listBackground() const { return m_tokens.listBackground; }
    QColor footerBackground() const { return m_tokens.footerBackground; }
    QColor popupBackground() const { return m_tokens.popupBackground; }
    QColor chipBackground() const { return m_tokens.chipBackground; }
    QColor bannerBackground() const { return m_tokens.bannerBackground; }
    QColor codePanelBackground() const { return m_tokens.codePanelBackground; }
    QColor textPrimary() const { return m_tokens.textPrimary; }
    QColor textSecondary() const { return m_tokens.textSecondary; }
    QColor textMuted() const { return m_tokens.textMuted; }
    QColor textFaint() const { return m_tokens.textFaint; }
    QColor textDisabled() const { return m_tokens.textDisabled; }
    QColor bannerText() const { return m_tokens.bannerText; }
    QColor onAccent() const { return m_tokens.onAccent; }
    QColor border() const { return m_tokens.border; }
    QColor borderStrong() const { return m_tokens.borderStrong; }
    QColor quoteBar() const { return m_tokens.quoteBar; }
    QColor mutedGlyph() const { return m_tokens.mutedGlyph; }
    QColor hoverTint() const { return m_tokens.hoverTint; }
    QColor blockHoverTint() const { return m_tokens.blockHoverTint; }
    QColor focusTint() const { return m_tokens.focusTint; }
    QColor focusRing() const { return m_tokens.focusRing; }
    QColor selectionTint() const { return m_tokens.selectionTint; }
    QColor selectionActiveTint() const { return m_tokens.selectionActiveTint; }
    QColor blockSelectionTint() const { return m_tokens.blockSelectionTint; }
    QColor accent() const { return m_tokens.accent; }
    QColor danger() const { return m_tokens.danger; }
    QColor dangerBright() const { return m_tokens.dangerBright; }
    QColor success() const { return m_tokens.success; }
    QColor warning() const { return m_tokens.warning; }
    QColor pinColor() const { return m_tokens.pinColor; }
    QColor marker() const { return m_tokens.marker; }
    QColor inlineCodeBackground() const { return m_tokens.inlineCodeBackground; }
    QColor highlightBackground() const { return m_tokens.highlightBackground; }
    QColor link() const { return m_tokens.link; }
    QColor searchMatchBackground() const { return m_tokens.searchMatchBackground; }
    QColor searchCurrentBackground() const { return m_tokens.searchCurrentBackground; }
    QColor changedTextBackground() const { return m_tokens.changedTextBackground; }
    QColor codeKeyword() const { return m_tokens.codeKeyword; }
    QColor codeType() const { return m_tokens.codeType; }
    QColor codeString() const { return m_tokens.codeString; }
    QColor codeComment() const { return m_tokens.codeComment; }
    QColor codeNumber() const { return m_tokens.codeNumber; }
    QColor calloutTip() const { return m_tokens.calloutTip; }
    QColor axisAttention() const { return m_tokens.axisAttention; }
    QColor axisAttentionText() const { return m_tokens.axisAttentionText; }
    QColor axisAgent() const { return m_tokens.axisAgent; }
    QColor axisAgentText() const { return m_tokens.axisAgentText; }
    QColor scopeDiscovered() const { return m_tokens.scopeDiscovered; }
    QColor signalHard() const { return m_tokens.signalHard; }
    QColor signalSoft() const { return m_tokens.signalSoft; }
    QColor signalHygiene() const { return m_tokens.signalHygiene; }
    QColor hatchAlt() const { return m_tokens.hatchAlt; }

    QStringList categoricalRamp() const { return m_tokens.categoricalRamp; }
    QStringList sequentialRamp() const { return m_tokens.sequentialRamp; }
    QStringList divergingRamp() const { return m_tokens.divergingRamp; }

    // One categorical step by index, wrapping past the eighth.
    //
    // Wrapping rather than running out: a chart handed nine series has to draw
    // the ninth as something, and a repeat is a smaller lie than an
    // unannounced grey. A caller with more series than steps should group the
    // tail into "other" instead, and the gallery's chart pages say so.
    Q_INVOKABLE QColor categorical(int index) const;

    // Every colour in this theme that already means something, for the
    // validator and for anything else that needs to stay off them.
    Q_INVOKABLE QStringList reservedHues() const;

    QStringList colorPalette() const;
    QStringList highlightPalette() const;
    QStringList colorPaletteNames() const;
    QStringList highlightPaletteNames() const;

    // The name for one colour value: a palette entry's own name, one of the
    // two greys the text-colour picker adds, "Theme default" for the empty
    // string, or a fallback naming the hex for anything chosen through the
    // system colour dialog. This is what a swatch publishes as its accessible
    // name, and it is here rather than in QML so the folder, tag, text and
    // highlight pickers all say the same words for the same colour.
    Q_INVOKABLE QString colorName(const QString &value) const;

    // The label colour to draw on a filled control: near-black or white,
    // whichever contrasts more with the fill. This is what onAccent is
    // computed with, and what a call site drawing a label on any other fill
    // should ask — a delete button on `danger`, a chip on a tag's own colour.
    Q_INVOKABLE static QColor labelOn(const QColor &fill);

    // WCAG 2.1 relative contrast between two opaque colours, from 1.0
    // (identical) to 21.0 (black on white). Exposed to QML so the settings
    // dialog can warn about a custom accent that disappears into the page,
    // and used by the contrast floors in tests/test_theme.cpp.
    Q_INVOKABLE static qreal contrastRatio(const QColor &a, const QColor &b);

    // The unmodified token table of a built-in theme (tests compare
    // against these; the settings dialog previews swatches from them).
    // The accent override does not ride on this, so `onAccent` here is the
    // value derived from the theme's own accent.
    static const Tokens &tokensFor(const QString &resolvedTheme);

    // Swatch colors for the settings dialog's theme cards ("system"
    // previews as its current resolution).
    Q_INVOKABLE QVariantMap themePreview(const QString &themeId) const;

signals:
    void themeChanged();
    void reducedMotionChanged();

private:
    void refresh();
    QString resolveSystem() const;

    SettingsStore *m_settings = nullptr;
    SystemAppearance *m_systemAppearance = nullptr;
    QString m_themeId = QStringLiteral("light");
    QString m_reducedMotionSetting = QStringLiteral("system");
    QString m_resolved = QStringLiteral("light");
    QString m_accentOverride;
    QString m_highlightOverride;
    Tokens m_tokens;
};

#endif // THEME_H
