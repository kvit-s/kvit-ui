// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>

#include "palette.h"
#include "theme.h"

// The chart ramps, checked rather than asserted.
//
// prd.md §5.5 asks for three ramps that survive four theme surfaces and three
// colour-vision deficiencies, and the reason that is a test rather than a
// review is that a ramp reviewed by eye stays right until the next person
// edits one hex value. Every step of all three ramps, in all four themes, goes
// through the same arithmetic here that tools/design-ramps.py used to find
// them.
//
// The failure message names the rule, the step and the number, so a broken
// ramp says what is wrong with it rather than that something is.
class TestPalette : public QObject
{
    Q_OBJECT

private slots:
    void testOklchRoundTrips();
    void testDeltaEIsZeroForIdenticalColours();
    void testDeficiencySimulationCollapsesTheAxisItShould();
    void testCategoricalRampHolds_data();
    void testCategoricalRampHolds();
    void testSequentialRampHolds_data();
    void testSequentialRampHolds();
    void testDivergingRampHolds_data();
    void testDivergingRampHolds();
    void testRampsShareTheirHuesAcrossThemes();
    void testEveryThemeHasEveryRamp();

private:
    // Every ground a chart mark can be drawn on in a theme. The list matches
    // what tools/design-ramps.py searched against; if a theme gains a surface
    // token, it belongs in both.
    static QList<QColor> surfacesOf(const QString &theme);
    static QList<QColor> reservedOf(const QString &theme);
    static QList<QColor> rampOf(const QStringList &hexes);
    static QString describe(const QList<KvitPalette::Finding> &findings);
};

QList<QColor> TestPalette::surfacesOf(const QString &theme)
{
    const Theme::Tokens &t = Theme::tokensFor(theme);
    return { t.windowBackground, t.panelBackground, t.listBackground,
             t.chipBackground };
}

QList<QColor> TestPalette::reservedOf(const QString &theme)
{
    // Theme::reservedHues() is a member because it reads the *effective*
    // tokens, including a user's accent override. A test wants the built-in
    // table, so it builds a Theme, points it at the theme under test and asks
    // there — which also means the list the validator uses is the one the
    // application would use.
    Theme built;
    built.setThemeId(theme);
    QList<QColor> hues;
    const QStringList names = built.reservedHues();
    for (const QString &name : names)
        hues.append(QColor(name));
    return hues;
}

QList<QColor> TestPalette::rampOf(const QStringList &hexes)
{
    QList<QColor> ramp;
    for (const QString &hex : hexes)
        ramp.append(QColor(hex));
    return ramp;
}

QString TestPalette::describe(const QList<KvitPalette::Finding> &findings)
{
    QStringList lines;
    for (const KvitPalette::Finding &finding : findings)
        lines.append(QStringLiteral("  [%1] %2").arg(finding.rule, finding.detail));
    return lines.join(QLatin1Char('\n'));
}

void TestPalette::testOklchRoundTrips()
{
    // Any colour converted to perceptual coordinates and back is the colour it
    // started as. If this drifts, every threshold in the validator is being
    // applied to something slightly other than the ramp.
    const QStringList samples = {
        "#000000", "#ffffff", "#7f7f7f", "#5c9fe0", "#e06c60", "#5abd82",
        "#e0a34c", "#a37fd4", "#1e1e1e", "#f6efdf",
    };
    for (const QString &hex : samples) {
        const QColor original(hex);
        const QColor back = KvitPalette::fromOklch(KvitPalette::toOklch(original));
        QVERIFY2(qAbs(original.red() - back.red()) <= 1
                     && qAbs(original.green() - back.green()) <= 1
                     && qAbs(original.blue() - back.blue()) <= 1,
                 qPrintable(QStringLiteral("%1 round-tripped to %2")
                                .arg(hex, back.name())));
    }

    // A grey has no hue to speak of, and its chroma must say so: the chroma
    // floor in the validator is what stops a step that reads as grey being
    // used as a category, and it can only do that if grey measures as zero.
    QVERIFY(KvitPalette::toOklch(QColor("#7f7f7f")).c < 0.005);
}

void TestPalette::testDeltaEIsZeroForIdenticalColours()
{
    QCOMPARE(KvitPalette::deltaE(QColor("#5c9fe0"), QColor("#5c9fe0")), 0.0);
    // Black to white is the largest distance in the space; anything much less
    // than 100 here means the scale factor is wrong and every threshold is
    // being read on the wrong scale.
    QVERIFY(KvitPalette::deltaE(QColor("#000000"), QColor("#ffffff")) > 95.0);
}

void TestPalette::testDeficiencySimulationCollapsesTheAxisItShould()
{
    // The check on the simulation itself: a deuteranope loses the red-green
    // axis, so a saturated red and a saturated green must come out much closer
    // together than they went in. Without this, a validator could be running a
    // near-identity matrix and reporting that everything is fine.
    const QColor red("#c0392b");
    const QColor green("#1e874b");
    const qreal ordinary = KvitPalette::deltaE(red, green);
    const qreal deuteranopic = KvitPalette::deltaE(
        KvitPalette::simulate(red, KvitPalette::Deficiency::Deuteranopia),
        KvitPalette::simulate(green, KvitPalette::Deficiency::Deuteranopia));
    QVERIFY2(deuteranopic < ordinary * 0.6,
             qPrintable(QStringLiteral("red/green were ΔE %1 apart and are ΔE "
                                       "%2 under deuteranopia")
                            .arg(ordinary).arg(deuteranopic)));

    // And the blue-yellow axis survives it, which is why a ramp that leans on
    // that axis works for the readers this is about.
    const qreal blueYellow = KvitPalette::deltaE(
        KvitPalette::simulate(QColor("#2970c8"),
                              KvitPalette::Deficiency::Deuteranopia),
        KvitPalette::simulate(QColor("#e0a34c"),
                              KvitPalette::Deficiency::Deuteranopia));
    QVERIFY(blueYellow > 30.0);
}

void TestPalette::testCategoricalRampHolds_data()
{
    QTest::addColumn<QString>("theme");
    for (const char *theme : { "light", "dark", "sepia", "highContrast" })
        QTest::newRow(theme) << QString::fromLatin1(theme);
}

void TestPalette::testCategoricalRampHolds()
{
    QFETCH(QString, theme);
    const Theme::Tokens &tokens = Theme::tokensFor(theme);

    QCOMPARE(tokens.categoricalRamp.size(), 8);
    const QList<KvitPalette::Finding> findings = KvitPalette::validateCategorical(
        rampOf(tokens.categoricalRamp), surfacesOf(theme), reservedOf(theme));
    QVERIFY2(findings.isEmpty(),
             qPrintable(QStringLiteral("the %1 categorical ramp fails:\n%2")
                            .arg(theme, describe(findings))));
}

void TestPalette::testSequentialRampHolds_data()
{
    testCategoricalRampHolds_data();
}

void TestPalette::testSequentialRampHolds()
{
    QFETCH(QString, theme);
    const Theme::Tokens &tokens = Theme::tokensFor(theme);

    QCOMPARE(tokens.sequentialRamp.size(), 7);
    const QList<KvitPalette::Finding> findings = KvitPalette::validateSequential(
        rampOf(tokens.sequentialRamp), surfacesOf(theme));
    QVERIFY2(findings.isEmpty(),
             qPrintable(QStringLiteral("the %1 sequential ramp fails:\n%2")
                            .arg(theme, describe(findings))));
}

void TestPalette::testDivergingRampHolds_data()
{
    testCategoricalRampHolds_data();
}

void TestPalette::testDivergingRampHolds()
{
    QFETCH(QString, theme);
    const Theme::Tokens &tokens = Theme::tokensFor(theme);

    QCOMPARE(tokens.divergingRamp.size(), 7);
    const QList<KvitPalette::Finding> findings = KvitPalette::validateDiverging(
        rampOf(tokens.divergingRamp), surfacesOf(theme));
    QVERIFY2(findings.isEmpty(),
             qPrintable(QStringLiteral("the %1 diverging ramp fails:\n%2")
                            .arg(theme, describe(findings))));
}

void TestPalette::testRampsShareTheirHuesAcrossThemes()
{
    // The property that makes a category recognisable across a theme switch:
    // step 3 is the same colour in light and in dark, allowing for how much
    // lighter and how much more saturated it has to be to sit on a dark
    // ground.
    //
    // Thirty degrees: each theme may sit up to fourteen either side of the
    // shared spoke (HUE_TOLERANCE in tools/design-ramps.py), so two themes can
    // be twice that apart, plus a degree or two from writing the answer out as
    // a six-digit hex value. A wider drift means somebody has replaced one
    // theme's ramp with a separate palette, which is the thing this
    // arrangement exists to prevent.
    const QStringList themes = { "light", "dark", "sepia", "highContrast" };
    for (int step = 0; step < 8; ++step) {
        qreal reference = -1.0;
        for (const QString &theme : themes) {
            const QColor color(Theme::tokensFor(theme).categoricalRamp.at(step));
            const qreal hue = KvitPalette::toOklch(color).h;
            if (reference < 0.0) {
                reference = hue;
                continue;
            }
            qreal apart = qAbs(hue - reference);
            apart = qMin(apart, 360.0 - apart);
            QVERIFY2(apart <= 30.0,
                     qPrintable(QStringLiteral("categorical step %1 is hue %2 "
                                               "in %3 against %4 elsewhere, a "
                                               "drift of %5 degrees")
                                    .arg(step).arg(hue, 0, 'f', 0)
                                    .arg(theme).arg(reference, 0, 'f', 0)
                                    .arg(apart, 0, 'f', 0)));
        }
    }
}

void TestPalette::testEveryThemeHasEveryRamp()
{
    // A theme added to the table without ramps would leave a chart in it drawn
    // in whatever categorical() falls back to, which is one muted grey for
    // every series.
    for (const char *theme : { "light", "dark", "sepia", "highContrast" }) {
        // By value: tokensFor takes a QString reference, so binding its
        // result to a reference while the argument is a temporary built here
        // is a dangling reference the moment the statement ends.
        const QString id = QString::fromLatin1(theme);
        const Theme::Tokens &tokens = Theme::tokensFor(id);
        QVERIFY(!tokens.categoricalRamp.isEmpty());
        QVERIFY(!tokens.sequentialRamp.isEmpty());
        QVERIFY(!tokens.divergingRamp.isEmpty());
    }

    // And the wrap: nine series on an eight-step ramp reuse the first colour
    // rather than running out.
    Theme built;
    built.setThemeId(QStringLiteral("dark"));
    QCOMPARE(built.categorical(8), built.categorical(0));
    QCOMPARE(built.categorical(-1), built.categorical(7));
}

QTEST_MAIN(TestPalette)
#include "test_palette.moc"
