// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <cmath>
#include <QtTest/QtTest>
#include <QSignalSpy>
#include <QTemporaryDir>

#include "theme.h"
#include "settingsstore.h"
#include "systemappearance.h"

// The theme token object: the four built-in token tables, system resolution,
// the accent and highlight overrides, the contrast floors, and persistence
// through the settings store.
//
// kvit-notes' copy of this file also drives BlockEditorEngine, checking that
// the editor's highlighter takes its inline colours from the theme and
// rehighlights when the theme changes. That class is the editor's document
// renderer and stays there.
class TestTheme : public QObject
{
    Q_OBJECT

private slots:
    void testDefaultIsLightTable();
    void testThemeSwitchSwapsTokens();
    void testInvalidThemeIdRejected();
    void testSystemResolvesToLightOrDark();
    void testTablesAreCompleteAndDistinct();
    void testDarkAndSepiaKeepContrast();
    void testEveryTokenPairMeetsItsFloor_data();
    void testEveryTokenPairMeetsItsFloor();
    void testAccentLabelIsDerivedFromTheAccent();
    void testPaletteNamesCoverEveryColor();
    void testFocusRingIsVisible();
    void testHighContrastMeetsStricterFloor();
    void testReducedMotionScale();
    void testReducedMotionFollowsTheSystem();
    void testHighContrastFollowsTheSystem();
    void testAnExistingMotionChoiceSurvivesTheUpgrade();
    void testAccentOverride();
    void testHighlightOverride();
    void testInvalidOverrideClears();
    void testInvalidPersistedOverrideIgnored();
    void testPersistsThroughSettings();
    void testFirstStartDefaultsToSystem();
    void testStaleSettingsValueFallsBack();

};

void TestTheme::testDefaultIsLightTable()
{
    Theme theme;
    QCOMPARE(theme.themeId(), QString("light"));
    QCOMPARE(theme.resolvedTheme(), QString("light"));

    const Theme::Tokens light = Theme::tokensFor(QStringLiteral("light"));
    QCOMPARE(theme.windowBackground(), light.windowBackground);
    QCOMPARE(theme.textPrimary(), light.textPrimary);
    QCOMPARE(theme.accent(), light.accent);
    QCOMPARE(theme.marker(), light.marker);

    // The light table was the pre-Phase-9 appearance verbatim. It is no
    // longer, in one direction only: the tokens that failed a WCAG floor were
    // darkened to meet it (accessibility.md Finding 3), which is what moved
    // the list marker from #b8b8b8 to a value that clears 3:1 against white.
    // The rest of the engine's documented fallbacks are unchanged.
    QCOMPARE(theme.marker(), QColor("#949494"));
    QCOMPARE(theme.link(), QColor("#2970c8"));
    QCOMPARE(theme.highlightBackground(), QColor("#fdf3a9"));
    QCOMPARE(theme.inlineCodeBackground(), QColor("#f0f0ee"));
    QCOMPARE(theme.searchMatchBackground(), QColor("#b5dcff"));
    QCOMPARE(theme.searchCurrentBackground(), QColor("#ffb454"));
}

void TestTheme::testThemeSwitchSwapsTokens()
{
    Theme theme;
    QSignalSpy spy(&theme, &Theme::themeChanged);

    theme.setThemeId(QStringLiteral("dark"));
    QCOMPARE(spy.count(), 1);
    QCOMPARE(theme.resolvedTheme(), QString("dark"));
    QCOMPARE(theme.windowBackground(),
             Theme::tokensFor(QStringLiteral("dark")).windowBackground);

    theme.setThemeId(QStringLiteral("sepia"));
    QCOMPARE(spy.count(), 2);
    QCOMPARE(theme.windowBackground(),
             Theme::tokensFor(QStringLiteral("sepia")).windowBackground);

    // Setting the current id again is a no-op.
    theme.setThemeId(QStringLiteral("sepia"));
    QCOMPARE(spy.count(), 2);
}

void TestTheme::testInvalidThemeIdRejected()
{
    Theme theme;
    QSignalSpy spy(&theme, &Theme::themeChanged);
    theme.setThemeId(QStringLiteral("neon"));
    QCOMPARE(theme.themeId(), QString("light"));
    QCOMPARE(spy.count(), 0);
}

void TestTheme::testSystemResolvesToLightOrDark()
{
    Theme theme;
    theme.setThemeId(QStringLiteral("system"));
    QCOMPARE(theme.themeId(), QString("system"));
    QVERIFY(theme.resolvedTheme() == QLatin1String("light")
            || theme.resolvedTheme() == QLatin1String("dark"));
    QCOMPARE(theme.windowBackground(),
             Theme::tokensFor(theme.resolvedTheme()).windowBackground);
}

void TestTheme::testTablesAreCompleteAndDistinct()
{
    const QStringList themes{ QStringLiteral("light"), QStringLiteral("dark"),
                              QStringLiteral("sepia"),
                              QStringLiteral("highContrast") };
    for (const QString &id : themes) {
        const Theme::Tokens &t = Theme::tokensFor(id);
        // Spot the full struct through a representative of each group;
        // an unset QColor is invalid.
        const QList<QColor> all{ t.windowBackground, t.panelBackground,
            t.listBackground, t.footerBackground, t.popupBackground,
            t.chipBackground, t.bannerBackground, t.codePanelBackground,
            t.textPrimary, t.textSecondary, t.textMuted, t.textFaint,
            t.textDisabled, t.bannerText, t.onAccent, t.border,
            t.borderStrong, t.quoteBar, t.mutedGlyph, t.hoverTint,
            t.focusTint, t.focusRing, t.selectionTint, t.selectionActiveTint,
            t.blockSelectionTint, t.blockHoverTint, t.accent, t.danger, t.dangerBright,
            t.success, t.warning, t.pinColor, t.marker,
            t.inlineCodeBackground, t.highlightBackground, t.link,
            t.searchMatchBackground, t.searchCurrentBackground,
            t.changedTextBackground,
            t.codeKeyword, t.codeType, t.codeString, t.codeComment,
            t.codeNumber, t.calloutTip,
            t.axisAttention, t.axisAttentionText, t.axisAgent,
            t.axisAgentText, t.scopeDiscovered, t.signalHard, t.signalSoft,
            t.signalHygiene, t.hatchAlt };
        for (const QColor &c : all)
            QVERIFY2(c.isValid(), qPrintable(id + " has an unset token"));

        // The five code-highlight tokens are five
        // distinct colors within each theme, so the token classes are
        // visually separable.
        const QList<QColor> codeTokens{ t.codeKeyword, t.codeType,
            t.codeString, t.codeComment, t.codeNumber };
        for (int i = 0; i < codeTokens.size(); ++i)
            for (int j = i + 1; j < codeTokens.size(); ++j)
                QVERIFY2(codeTokens[i] != codeTokens[j],
                         qPrintable(id + ": code tokens must be distinct"));
    }

    // The three themes are actually different appearances.
    QVERIFY(Theme::tokensFor("light").windowBackground
            != Theme::tokensFor("dark").windowBackground);
    QVERIFY(Theme::tokensFor("light").windowBackground
            != Theme::tokensFor("sepia").windowBackground);
    QVERIFY(Theme::tokensFor("dark").windowBackground
            != Theme::tokensFor("sepia").windowBackground);
}

void TestTheme::testDarkAndSepiaKeepContrast()
{
    // Coarse legibility floor (not a WCAG-grade audit): body
    // text against the editor background keeps a strong luminance gap,
    // and text stays legible over the selection and search tints.
    auto luminance = [](const QColor &c) {
        return 0.2126 * c.redF() + 0.7152 * c.greenF() + 0.0722 * c.blueF();
    };
    for (const QString &id : { QStringLiteral("light"), QStringLiteral("dark"),
                               QStringLiteral("sepia") }) {
        const Theme::Tokens &t = Theme::tokensFor(id);
        QVERIFY2(qAbs(luminance(t.textPrimary)
                      - luminance(t.windowBackground)) > 0.55,
                 qPrintable(id + ": body text vs background"));
        QVERIFY2(qAbs(luminance(t.textPrimary)
                      - luminance(t.selectionTint)) > 0.35,
                 qPrintable(id + ": text vs selection tint"));
        QVERIFY2(qAbs(luminance(t.textPrimary)
                      - luminance(t.searchMatchBackground)) > 0.3,
                 qPrintable(id + ": text vs search tint"));
    }
}

// Every foreground/background pair that actually appears together on screen,
// in every theme, against the WCAG 2.1 level AA floor for what it is: 4.5:1
// where the foreground is text, 3:1 where it is the part of a control that
// says where the control is or what state it is in (accessibility.md
// Finding 3).
//
// This is the check the two coarse floors above do not reach. They compare
// luminance gaps for body text and hold the focus ring to a real ratio;
// everything else — the faint text a date or a count is set in, the code
// theme, a warning, the outline of an unchecked to-do box — went unmeasured,
// and most of it was below the line in the light and sepia themes.
//
// `border` is deliberately absent. It is the decorative rule between two
// panels, which WCAG asks nothing of; `borderStrong` is the control-boundary
// token and is here. See the note on the two in theme.h.
void TestTheme::testEveryTokenPairMeetsItsFloor_data()
{
    QTest::addColumn<QString>("themeId");
    QTest::addColumn<QString>("pairName");
    QTest::addColumn<QColor>("foreground");
    QTest::addColumn<QColor>("background");
    QTest::addColumn<double>("floor");

    struct Pair {
        const char *name;
        QColor Theme::Tokens::*foreground;
        QColor Theme::Tokens::*background;
        double floor;
    };
    // 4.5 = text; 3.0 = a control's own boundary or state mark.
    static const Pair kPairs[] = {
        {"onAccent on accent", &Theme::Tokens::onAccent,
         &Theme::Tokens::accent, 4.5},
        {"warning on window", &Theme::Tokens::warning,
         &Theme::Tokens::windowBackground, 4.5},
        {"success on window", &Theme::Tokens::success,
         &Theme::Tokens::windowBackground, 4.5},
        {"textFaint on window", &Theme::Tokens::textFaint,
         &Theme::Tokens::windowBackground, 4.5},
        {"codeComment on code panel", &Theme::Tokens::codeComment,
         &Theme::Tokens::codePanelBackground, 4.5},
        {"codeString on code panel", &Theme::Tokens::codeString,
         &Theme::Tokens::codePanelBackground, 4.5},
        {"codeType on code panel", &Theme::Tokens::codeType,
         &Theme::Tokens::codePanelBackground, 4.5},
        {"bannerText on banner", &Theme::Tokens::bannerText,
         &Theme::Tokens::bannerBackground, 4.5},
        {"textPrimary on the current match",
         &Theme::Tokens::textPrimary,
         &Theme::Tokens::searchCurrentBackground, 4.5},
        {"textPrimary on active selection", &Theme::Tokens::textPrimary,
         &Theme::Tokens::selectionActiveTint, 4.5},
        {"textPrimary on changed text", &Theme::Tokens::textPrimary,
         &Theme::Tokens::changedTextBackground, 4.5},
        {"borderStrong on window", &Theme::Tokens::borderStrong,
         &Theme::Tokens::windowBackground, 3.0},
        {"mutedGlyph on panel", &Theme::Tokens::mutedGlyph,
         &Theme::Tokens::panelBackground, 3.0},
        {"quoteBar on window", &Theme::Tokens::quoteBar,
         &Theme::Tokens::windowBackground, 3.0},
        {"marker on window", &Theme::Tokens::marker,
         &Theme::Tokens::windowBackground, 3.0},
    };

    for (const QString &id : { QStringLiteral("light"), QStringLiteral("dark"),
                               QStringLiteral("sepia"),
                               QStringLiteral("highContrast") }) {
        const Theme::Tokens &t = Theme::tokensFor(id);
        for (const Pair &p : kPairs) {
            QTest::newRow(qPrintable(id + ": " + QString::fromLatin1(p.name)))
                << id << QString::fromLatin1(p.name)
                << t.*(p.foreground) << t.*(p.background) << p.floor;
        }
    }
}

void TestTheme::testEveryTokenPairMeetsItsFloor()
{
    QFETCH(QString, themeId);
    QFETCH(QString, pairName);
    QFETCH(QColor, foreground);
    QFETCH(QColor, background);
    QFETCH(double, floor);

    const double r = Theme::contrastRatio(foreground, background);
    QVERIFY2(r >= floor,
             qPrintable(themeId + ": " + pairName + " is only "
                        + QString::number(r, 'f', 2) + ":1 (need "
                        + QString::number(floor, 'f', 1) + ":1)"));
}

// The accent label is computed, not stored, so a custom accent gets a label
// that suits it. Without this the label stayed pure white whatever the user
// picked, and a pale accent gave white text on a pale fill.
void TestTheme::testAccentLabelIsDerivedFromTheAccent()
{
    Theme theme;

    // Pale accent → a dark label; dark accent → a light one. Both clear the
    // 4.5:1 text floor, which is the property that matters rather than the
    // exact shade.
    theme.setAccentOverride(QStringLiteral("#ffe9a0"));
    QVERIFY(Theme::contrastRatio(theme.onAccent(), theme.accent()) >= 4.5);
    const QColor onPale = theme.onAccent();

    theme.setAccentOverride(QStringLiteral("#102030"));
    QVERIFY(Theme::contrastRatio(theme.onAccent(), theme.accent()) >= 4.5);
    QVERIFY2(theme.onAccent() != onPale,
             "the label has to move when the accent moves from pale to dark");

    // Clearing the override puts the theme's own accent back, label included.
    theme.setAccentOverride(QString());
    QCOMPARE(theme.onAccent(),
             Theme::tokensFor(QStringLiteral("light")).onAccent);
}

// A swatch is announced by name, so every colour a picker offers has to have
// one. The two lists are parallel to the palettes by position, which is the
// kind of pairing that goes wrong silently when a colour is added.
void TestTheme::testPaletteNamesCoverEveryColor()
{
    Theme theme;
    QCOMPARE(theme.colorPaletteNames().size(), theme.colorPalette().size());
    QCOMPARE(theme.highlightPaletteNames().size(),
             theme.highlightPalette().size());

    const QStringList palette = theme.colorPalette();
    for (int i = 0; i < palette.size(); ++i) {
        QCOMPARE(theme.colorName(palette.at(i)),
                 theme.colorPaletteNames().at(i));
        // Spelling must not decide the answer: the same colour arrives
        // lowercase from the palette and uppercase from a hand-edited note.
        QCOMPARE(theme.colorName(palette.at(i).toUpper()),
                 theme.colorPaletteNames().at(i));
    }
    QVERIFY(!theme.colorName(QString()).isEmpty());     // "Theme default"
    QVERIFY(!theme.colorName(QStringLiteral("#123456")).isEmpty());
}

void TestTheme::testFocusRingIsVisible()
{
    // The keyboard-focus ring (§14.1) must stand out against the editor
    // background: a WCAG non-text contrast floor of 3:1 in every theme.
    auto lin = [](double c) {
        return c <= 0.03928 ? c / 12.92 : std::pow((c + 0.055) / 1.055, 2.4);
    };
    auto relLum = [&](const QColor &c) {
        return 0.2126 * lin(c.redF()) + 0.7152 * lin(c.greenF())
             + 0.0722 * lin(c.blueF());
    };
    auto ratio = [&](const QColor &a, const QColor &b) {
        double la = relLum(a), lb = relLum(b);
        double hi = qMax(la, lb), lo = qMin(la, lb);
        return (hi + 0.05) / (lo + 0.05);
    };
    for (const QString &id : { QStringLiteral("light"), QStringLiteral("dark"),
                               QStringLiteral("sepia") }) {
        const Theme::Tokens &t = Theme::tokensFor(id);
        const double r = ratio(t.focusRing, t.windowBackground);
        QVERIFY2(r >= 3.0,
                 qPrintable(id + ": focus ring vs background is only "
                            + QString::number(r, 'f', 2) + ":1 (need 3:1)"));
    }
}

void TestTheme::testHighContrastMeetsStricterFloor()
{
    // The high-contrast theme holds a stricter floor than the others:
    // WCAG AAA (7:1) for body text and ≥4.5:1 for the focus ring, both vs the
    // editor background.
    auto lin = [](double c) {
        return c <= 0.03928 ? c / 12.92 : std::pow((c + 0.055) / 1.055, 2.4);
    };
    auto relLum = [&](const QColor &c) {
        return 0.2126 * lin(c.redF()) + 0.7152 * lin(c.greenF())
             + 0.0722 * lin(c.blueF());
    };
    auto ratio = [&](const QColor &a, const QColor &b) {
        double hi = qMax(relLum(a), relLum(b)), lo = qMin(relLum(a), relLum(b));
        return (hi + 0.05) / (lo + 0.05);
    };
    const Theme::Tokens t = Theme::tokensFor(QStringLiteral("highContrast"));
    const double text = ratio(t.textPrimary, t.windowBackground);
    QVERIFY2(text >= 7.0,
             qPrintable("high contrast body text is only "
                        + QString::number(text, 'f', 2) + ":1 (need 7:1)"));
    const double ring = ratio(t.focusRing, t.windowBackground);
    QVERIFY2(ring >= 4.5,
             qPrintable("high contrast focus ring is only "
                        + QString::number(ring, 'f', 2) + ":1 (need 4.5:1)"));
    // It is selectable and resolves to its own tokens.
    Theme theme;
    theme.setThemeId(QStringLiteral("highContrast"));
    QCOMPARE(theme.resolvedTheme(), QString("highContrast"));
    QCOMPARE(theme.windowBackground(), QColor("#000000"));
    QCOMPARE(theme.displayName("highContrast"), QString("High contrast"));
}

void TestTheme::testReducedMotionScale()
{
    // Reduced motion (§14.3) is one source: motionScale is 1 normally and 0
    // when on, so every animation multiplying by it stills instantly. It
    // persists through the settings store.
    Theme theme;
    QCOMPARE(theme.reducedMotion(), false);
    QCOMPARE(theme.motionScale(), 1.0);
    QSignalSpy spy(&theme, &Theme::reducedMotionChanged);
    theme.setReducedMotion(true);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(theme.reducedMotion(), true);
    QCOMPARE(theme.motionScale(), 0.0);
    // Setting the same value again is a no-op.
    theme.setReducedMotion(true);
    QCOMPARE(spy.count(), 1);

    // Persistence: a second theme on the same store reads it back.
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore store;
    store.open(dir.filePath(QStringLiteral("s.json")));
    Theme a;
    a.setSettings(&store);
    a.setReducedMotion(true);
    Theme b;
    b.setSettings(&store);
    QCOMPARE(b.reducedMotion(), true);
    QCOMPARE(b.motionScale(), 0.0);
}

// "system" is the third state: the desktop's own reduce-motion preference,
// which is where a person who needs it has usually already said so
// (accessibility.md Finding 7).
void TestTheme::testReducedMotionFollowsTheSystem()
{
    SystemAppearance system;
    system.setOverride(false, /*reducedMotion=*/true);

    Theme theme;
    QCOMPARE(theme.reducedMotionSetting(), QString("system"));
    // Nothing attached: "system" has nothing to follow, so motion stays on.
    QCOMPARE(theme.reducedMotion(), false);

    QSignalSpy spy(&theme, &Theme::reducedMotionChanged);
    theme.setSystemAppearance(&system);
    QVERIFY(spy.count() >= 1);
    QCOMPARE(theme.reducedMotion(), true);
    QCOMPARE(theme.motionScale(), 0.0);

    // An explicit choice outranks the desktop, in both directions.
    theme.setReducedMotionSetting(QStringLiteral("off"));
    QCOMPARE(theme.reducedMotion(), false);
    theme.setReducedMotionSetting(QStringLiteral("on"));
    QCOMPARE(theme.reducedMotion(), true);

    // And going back to "system" starts following it again.
    theme.setReducedMotionSetting(QStringLiteral("system"));
    QCOMPARE(theme.reducedMotion(), true);
    system.setOverride(false, false);
    QCOMPARE(theme.reducedMotion(), false);

    // A value no setting can hold is refused rather than stored.
    theme.setReducedMotionSetting(QStringLiteral("sometimes"));
    QCOMPARE(theme.reducedMotionSetting(), QString("system"));
}

// System-wide high contrast folds into the theme's "system" setting, which
// already means "follow the desktop". An explicit theme choice is the user
// speaking about this application in particular and wins.
void TestTheme::testHighContrastFollowsTheSystem()
{
    SystemAppearance system;
    system.setOverride(/*highContrast=*/true, false);

    Theme theme;
    theme.setThemeId(QStringLiteral("system"));
    theme.setSystemAppearance(&system);
    QCOMPARE(theme.resolvedTheme(), QString("highContrast"));
    QCOMPARE(theme.windowBackground(),
             Theme::tokensFor(QStringLiteral("highContrast")).windowBackground);

    // An explicit theme is not overridden by the desktop.
    theme.setThemeId(QStringLiteral("sepia"));
    QCOMPARE(theme.resolvedTheme(), QString("sepia"));

    // Back to following, and then the desktop turning it off.
    theme.setThemeId(QStringLiteral("system"));
    QCOMPARE(theme.resolvedTheme(), QString("highContrast"));
    system.setOverride(false, false);
    QVERIFY(theme.resolvedTheme() == QLatin1String("light")
            || theme.resolvedTheme() == QLatin1String("dark"));
}

// An installation that had already turned reduced motion on keeps it on. The
// old setting was a plain bool, written only when somebody chose; its
// presence is what distinguishes a choice from a default.
void TestTheme::testAnExistingMotionChoiceSurvivesTheUpgrade()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore store;
    QVERIFY(store.open(dir.filePath(QStringLiteral("s.json"))));
    store.setValue(QStringLiteral("view.reducedMotion"), true);

    Theme upgraded;
    upgraded.setSettings(&store);
    QCOMPARE(upgraded.reducedMotionSetting(), QString("on"));
    QCOMPARE(upgraded.reducedMotion(), true);

    // Someone who had explicitly turned it OFF keeps that too, rather than
    // silently starting to follow a desktop that asks for it.
    SettingsStore off;
    QVERIFY(off.open(dir.filePath(QStringLiteral("off.json"))));
    off.setValue(QStringLiteral("view.reducedMotion"), false);
    SystemAppearance system;
    system.setOverride(false, true);
    Theme keptOff;
    keptOff.setSystemAppearance(&system);
    keptOff.setSettings(&off);
    QCOMPARE(keptOff.reducedMotionSetting(), QString("off"));
    QCOMPARE(keptOff.reducedMotion(), false);

    // A fresh installation has neither key and follows the desktop.
    SettingsStore fresh;
    QVERIFY(fresh.open(dir.filePath(QStringLiteral("fresh.json"))));
    Theme brandNew;
    brandNew.setSystemAppearance(&system);
    brandNew.setSettings(&fresh);
    QCOMPARE(brandNew.reducedMotionSetting(), QString("system"));
    QCOMPARE(brandNew.reducedMotion(), true);
}

void TestTheme::testAccentOverride()
{
    Theme theme;
    QSignalSpy spy(&theme, &Theme::themeChanged);
    theme.setAccentOverride(QStringLiteral("#aa3366"));
    QCOMPARE(spy.count(), 1);
    QCOMPARE(theme.accent(), QColor("#aa3366"));
    QCOMPARE(theme.link(), QColor("#aa3366"));
    // Other tokens are untouched.
    QCOMPARE(theme.textPrimary(),
             Theme::tokensFor(QStringLiteral("light")).textPrimary);

    // The override rides across theme switches...
    theme.setThemeId(QStringLiteral("dark"));
    QCOMPARE(theme.accent(), QColor("#aa3366"));

    // ...and clearing restores the table value.
    theme.setAccentOverride(QString());
    QCOMPARE(theme.accent(), Theme::tokensFor(QStringLiteral("dark")).accent);
    QCOMPARE(theme.link(), Theme::tokensFor(QStringLiteral("dark")).link);
}

void TestTheme::testHighlightOverride()
{
    Theme theme;
    theme.setHighlightOverride(QStringLiteral("#c2f0c2"));
    QCOMPARE(theme.highlightBackground(), QColor("#c2f0c2"));
    QCOMPARE(theme.accent(),
             Theme::tokensFor(QStringLiteral("light")).accent);
    theme.setHighlightOverride(QString());
    QCOMPARE(theme.highlightBackground(),
             Theme::tokensFor(QStringLiteral("light")).highlightBackground);
}

void TestTheme::testInvalidOverrideClears()
{
    Theme theme;
    theme.setAccentOverride(QStringLiteral("#aa3366"));
    theme.setAccentOverride(QStringLiteral("not-a-color"));
    QCOMPARE(theme.accentOverride(), QString());
    QCOMPARE(theme.accent(), Theme::tokensFor(QStringLiteral("light")).accent);
}

// setAccentOverride/setHighlightOverride reject a string QColor cannot
// parse, but setSettings() reads the same two keys straight out of the
// store. A hand-edited or corrupted settings file therefore installs an
// invalid QColor as the accent, which renders as black rather than falling
// back to the theme's own accent.
void TestTheme::testInvalidPersistedOverrideIgnored()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    const QString path = dir.filePath("settings.json");

    {
        SettingsStore store;
        QVERIFY(store.open(path));
        store.setValue(QStringLiteral("theme.accent"),
                       QStringLiteral("not-a-color"));
        store.setValue(QStringLiteral("theme.highlight"),
                       QStringLiteral("#zzzzzz"));
        store.flush();
    }

    SettingsStore store;
    QVERIFY(store.open(path));
    Theme theme;
    theme.setSettings(&store);

    // The junk is discarded and the theme's own tokens stand. Which tokens
    // those are depends on the desktop: the default theme id is "system", and
    // it resolves to dark on a machine set to dark, which is what made this
    // case fail on Windows while passing on a session that reports no
    // preference. The assertion is about the override being ignored, so it
    // reads the tokens of whatever theme resolved.
    QCOMPARE(theme.accentOverride(), QString());
    QCOMPARE(theme.highlightOverride(), QString());
    const Theme::Tokens resolved = Theme::tokensFor(theme.resolvedTheme());
    QVERIFY(theme.accent().isValid());
    QVERIFY(theme.highlightBackground().isValid());
    QCOMPARE(theme.accent(), resolved.accent);
    QCOMPARE(theme.highlightBackground(), resolved.highlightBackground);
}

void TestTheme::testPersistsThroughSettings()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    const QString path = dir.filePath("settings.json");

    {
        SettingsStore store;
        QVERIFY(store.open(path));
        Theme theme;
        theme.setSettings(&store);
        theme.setThemeId(QStringLiteral("sepia"));
        theme.setAccentOverride(QStringLiteral("#aa3366"));
        theme.setHighlightOverride(QStringLiteral("#c2f0c2"));
        store.flush();
    }

    SettingsStore reopened;
    QVERIFY(reopened.open(path));
    Theme theme;
    theme.setSettings(&reopened);
    QCOMPARE(theme.themeId(), QString("sepia"));
    QCOMPARE(theme.accent(), QColor("#aa3366"));
    QCOMPARE(theme.highlightBackground(), QColor("#c2f0c2"));
}

void TestTheme::testFirstStartDefaultsToSystem()
{
    // No persisted choice yet: follow the OS color scheme (which always
    // resolves to a concrete built-in table).
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore store;
    QVERIFY(store.open(dir.filePath("settings.json")));

    Theme theme;
    theme.setSettings(&store);
    QCOMPARE(theme.themeId(), QString("system"));
    QVERIFY(theme.resolvedTheme() == QLatin1String("light")
            || theme.resolvedTheme() == QLatin1String("dark"));
}

void TestTheme::testStaleSettingsValueFallsBack()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore store;
    QVERIFY(store.open(dir.filePath("settings.json")));
    store.setValue("theme.id", "neon");

    Theme theme;
    theme.setSettings(&store);
    QCOMPARE(theme.themeId(), QString("system"));
}

QTEST_MAIN(TestTheme)
#include "test_theme.moc"
