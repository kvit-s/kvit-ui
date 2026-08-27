// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QFontInfo>
#include <QGuiApplication>
#include <QMetaProperty>
#include <QSignalSpy>
#include <QTemporaryDir>

#include "interfacemetrics.h"
#include "settingsstore.h"
#include "typography.h"

// The density and spacing tokens, and the claim the whole arrangement rests
// on: one setting moves every type role and every geometry value together, and
// leaves the document type scale alone.
//
// kvit-hub's Tokens.qml is what this replaces — 141 lines of bare integers for
// row heights, chip heights, radii, pane widths and reflow breakpoints, all of
// which stood still while its type setting moved. That is why raising its
// interface size made a row too tight for the text in it, and why the check
// below is over *every* property rather than over a sample.
class TestDensity : public QObject
{
    Q_OBJECT

private slots:
    void testEveryValueIsTheIdentityAtTheDefault();
    void testEveryValueScalesWithOneSetting();
    void testNoDensityValueCollapsesAtTheSmallestSize();
    void testTheTypeScaleIsOrdered();
    void testTheSpacingScaleIsOrdered();
    void testTheDocumentTypeScaleIsUntouched();
    void testEveryKvitHubTokenHasAHome_data();
    void testEveryKvitHubTokenHasAHome();
    void testTheFontFamiliesAreSeparateFromTheDocument();
    void testTheResolvedFamiliesAreDrawable();

private:
    // Every int property InterfaceMetrics publishes except fontSize itself and
    // the two clamps, which are the inputs rather than the derived values.
    static QStringList derivedProperties();
    static int valueOf(const InterfaceMetrics &metrics, const QString &name);
};

QStringList TestDensity::derivedProperties()
{
    static const QSet<QString> inputs = {
        QStringLiteral("fontSize"), QStringLiteral("minFontSize"),
        QStringLiteral("maxFontSize"), QStringLiteral("objectName"),
    };
    QStringList names;
    const QMetaObject *meta = &InterfaceMetrics::staticMetaObject;
    for (int i = 0; i < meta->propertyCount(); ++i) {
        const QMetaProperty property = meta->property(i);
        if (property.metaType().id() != QMetaType::Int)
            continue;
        const QString name = QString::fromLatin1(property.name());
        if (!inputs.contains(name))
            names.append(name);
    }
    return names;
}

int TestDensity::valueOf(const InterfaceMetrics &metrics, const QString &name)
{
    return metrics.property(name.toLatin1().constData()).toInt();
}

void TestDensity::testEveryValueIsTheIdentityAtTheDefault()
{
    // The claim the migration of two applications rests on: at the default
    // base, every named value answers exactly the literal that was written out
    // at the call site before it existed, so converting a pane changes no
    // pixels until somebody moves the setting.
    //
    // A failure here means the default build has silently reflowed, which is
    // the one outcome Wave 2 and Wave 3 cannot absorb: the acceptance test for
    // both is that the screenshots come back identical.
    InterfaceMetrics metrics;
    QCOMPARE(metrics.fontSize(), 12);
    QCOMPARE(metrics.scale(), 1.0);

    // The seven type roles.
    QCOMPARE(metrics.caption(), 10);
    QCOMPARE(metrics.small(), 11);
    QCOMPARE(metrics.body(), 12);
    QCOMPARE(metrics.strong(), 13);
    QCOMPARE(metrics.title(), 15);
    QCOMPARE(metrics.headline(), 17);
    QCOMPARE(metrics.display(), 20);

    // The spacing scale.
    QCOMPARE(metrics.spaceTight(), 2);
    QCOMPARE(metrics.spaceSnug(), 4);
    QCOMPARE(metrics.spaceNear(), 6);
    QCOMPARE(metrics.space(), 8);
    QCOMPARE(metrics.spaceWide(), 10);
    QCOMPARE(metrics.spaceLoose(), 12);

    // Every value kvit-hub's Tokens.qml named, at the value it named.
    QCOMPARE(metrics.viewMargin(), 16);
    QCOMPARE(metrics.columnGap(), 14);
    QCOMPARE(metrics.stackGap(), 7);
    QCOMPARE(metrics.sidebarWidth(), 232);
    QCOMPARE(metrics.railWidth(), 48);
    QCOMPARE(metrics.paneWidth(), 392);
    QCOMPARE(metrics.headerHeight(), 52);
    QCOMPARE(metrics.breadcrumbHeight(), 34);
    QCOMPARE(metrics.rowHeight(), 56);
    QCOMPARE(metrics.rowHeightSub(), 48);
    QCOMPARE(metrics.rowHeightSlim(), 30);
    QCOMPARE(metrics.rowHeightCompact(), 24);
    QCOMPARE(metrics.tabHeight(), 30);
    QCOMPARE(metrics.chipHeight(), 17);
    QCOMPARE(metrics.tagHeight(), 16);
    QCOMPARE(metrics.pillHeight(), 15);
    QCOMPARE(metrics.barHeight(), 7);
    QCOMPARE(metrics.barHeightWide(), 9);
    QCOMPARE(metrics.radiusBar(), 2);
    QCOMPARE(metrics.radiusChip(), 3);
    QCOMPARE(metrics.radiusControl(), 4);
    QCOMPARE(metrics.radiusCard(), 6);
    QCOMPARE(metrics.radiusPill(), 8);
    QCOMPARE(metrics.hairline(), 1);
    QCOMPARE(metrics.widthFloor(), 880);
    QCOMPARE(metrics.widthLaptop(), 1100);
    QCOMPARE(metrics.widthDrawn(), 1440);
}

void TestDensity::testEveryValueScalesWithOneSetting()
{
    // Over every derived property by reflection rather than over the ones
    // somebody remembered to list. A property added to the header without
    // going through px() is exactly what this catches, and it is the mistake
    // that produced kvit-hub's frozen density block in the first place.
    InterfaceMetrics metrics;
    const QStringList names = derivedProperties();
    QVERIFY(names.size() > 30);

    QHash<QString, int> atDefault;
    for (const QString &name : names)
        atDefault.insert(name, valueOf(metrics, name));

    QSignalSpy spy(&metrics, &InterfaceMetrics::changed);
    metrics.setFontSize(24);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(metrics.scale(), 2.0);

    for (const QString &name : names) {
        const int before = atDefault.value(name);
        const int after = valueOf(metrics, name);
        QVERIFY2(after > before,
                 qPrintable(QStringLiteral("%1 stayed at %2 when the interface "
                                           "size doubled — it is a literal "
                                           "rather than a px() value")
                                .arg(name).arg(before)));
        // Doubling the base doubles the value, to within the rounding px()
        // does. A value that moves by some other factor is one that went
        // through its own arithmetic instead.
        QVERIFY2(qAbs(after - before * 2) <= 1,
                 qPrintable(QStringLiteral("%1 went from %2 to %3, which is "
                                           "not the doubling every other value "
                                           "did")
                                .arg(name).arg(before).arg(after)));
    }
}

void TestDensity::testNoDensityValueCollapsesAtTheSmallestSize()
{
    // At the smallest interface size a hairline scales to 0.83 pixels, and a
    // separator rounded to zero is a line that vanishes — which reads as a
    // layout bug rather than as a smaller interface. px() holds every positive
    // value at one or more, and this is over the whole set rather than over
    // the hairline alone, because the radii and the two-pixel spacing step are
    // in the same position.
    InterfaceMetrics metrics;
    metrics.setFontSize(InterfaceMetrics::MinFontSize);
    for (const QString &name : derivedProperties()) {
        QVERIFY2(valueOf(metrics, name) >= 1,
                 qPrintable(QStringLiteral("%1 collapsed to %2 at the smallest "
                                           "interface size")
                                .arg(name).arg(valueOf(metrics, name))));
    }
}

void TestDensity::testTheTypeScaleIsOrdered()
{
    // Seven steps, each strictly larger than the last, at every size in the
    // range. Two roles that collide at some size are two roles a reader cannot
    // tell apart there, which defeats the point of naming them separately.
    InterfaceMetrics metrics;
    for (int size = InterfaceMetrics::MinFontSize;
         size <= InterfaceMetrics::MaxFontSize; ++size) {
        metrics.setFontSize(size);
        const QList<int> scale = { metrics.caption(), metrics.small(),
                                   metrics.body(), metrics.strong(),
                                   metrics.title(), metrics.headline(),
                                   metrics.display() };
        for (int i = 1; i < scale.size(); ++i) {
            QVERIFY2(scale.at(i) > scale.at(i - 1),
                     qPrintable(QStringLiteral("at interface size %1 the type "
                                               "scale reads %2, which is not "
                                               "increasing")
                                    .arg(size)
                                    .arg(QStringList(
                                             std::accumulate(
                                                 scale.begin(), scale.end(),
                                                 QStringList(),
                                                 [](QStringList acc, int v) {
                                                     acc.append(QString::number(v));
                                                     return acc;
                                                 }))
                                             .join(QStringLiteral(", ")))));
        }
    }
}

void TestDensity::testTheSpacingScaleIsOrdered()
{
    InterfaceMetrics metrics;
    for (int size = InterfaceMetrics::MinFontSize;
         size <= InterfaceMetrics::MaxFontSize; ++size) {
        metrics.setFontSize(size);
        const QList<int> steps = { metrics.hairline(), metrics.spaceTight(),
                                   metrics.spaceSnug(), metrics.spaceNear(),
                                   metrics.space(), metrics.spaceWide(),
                                   metrics.spaceLoose() };
        for (int i = 1; i < steps.size(); ++i) {
            // Not strictly increasing, because at the smallest size the floor
            // in px() pulls the bottom of the scale together: a 1 and a 2
            // design pixel both come out at 1. Non-decreasing is what can be
            // asked for, and it is what matters — a gap that is *smaller* than
            // the one below it in the scale would be a scale in name only.
            QVERIFY2(steps.at(i) >= steps.at(i - 1),
                     qPrintable(QStringLiteral("at interface size %1 spacing "
                                               "step %2 is %3, below step %4's "
                                               "%5")
                                    .arg(size).arg(i).arg(steps.at(i))
                                    .arg(i - 1).arg(steps.at(i - 1))));
        }
    }
}

void TestDensity::testTheDocumentTypeScaleIsUntouched()
{
    // Two settings, one store, and neither reaches the other. This is
    // divergence 3.4 from the other side: kvit-hub's "Interface size" row
    // writes `typography.baseSize`, which in kvit-notes is the *document* font
    // size, so moving one silently moves the other in an application that
    // shares the settings file. What Wave 1 owes kvit-hub is a correct thing
    // for that stepper to point at, and this is the check that it is correct.
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore settings;
    QVERIFY(settings.open(dir.filePath(QStringLiteral("settings.json"))));

    InterfaceMetrics metrics;
    Typography typography;
    metrics.setSettings(&settings);
    typography.setSettings(&settings);

    const int documentBefore = typography.baseSize();
    const int documentBodyBefore = typography.bodySize();

    metrics.setFontSize(24);
    QCOMPARE(metrics.body(), 24);
    QCOMPARE(typography.baseSize(), documentBefore);
    QCOMPARE(typography.bodySize(), documentBodyBefore);

    // And the other way: the document setting moves nothing in the chrome.
    const int chromeBefore = metrics.body();
    typography.setBaseSize(20);
    QCOMPARE(typography.baseSize(), 20);
    QCOMPARE(metrics.body(), chromeBefore);

    // They persist under different keys, which is what makes the separation
    // survive a restart rather than only a session.
    settings.flush();
    QCOMPARE(settings.value(QStringLiteral("interface.fontSize")).toInt(), 24);
    QCOMPARE(settings.value(QStringLiteral("typography.fontSize")).toInt(), 20);
}

void TestDensity::testEveryKvitHubTokenHasAHome_data()
{
    // Every name kvit-hub's theme/Tokens.qml defines in its density block,
    // against the property here that replaces it. Wave 3 deletes that file,
    // and it can only do that if every name it defines resolves to something.
    QTest::addColumn<QString>("hubName");
    QTest::addColumn<QString>("property");

    const QList<QPair<QString, QString>> mapping = {
        { "viewMargin", "viewMargin" },
        { "columnGap", "columnGap" },
        { "stackGap", "stackGap" },
        { "sidebarWidth", "sidebarWidth" },
        { "railWidth", "railWidth" },
        { "paneWidth", "paneWidth" },
        { "headerHeight", "headerHeight" },
        { "breadcrumbHeight", "breadcrumbHeight" },
        { "rowHeight", "rowHeight" },
        { "rowHeightSub", "rowHeightSub" },
        { "rowHeightSlim", "rowHeightSlim" },
        { "rowHeightCompact", "rowHeightCompact" },
        { "tabHeight", "tabHeight" },
        { "chipHeight", "chipHeight" },
        { "tagHeight", "tagHeight" },
        { "pillHeight", "pillHeight" },
        { "barHeight", "barHeight" },
        { "barHeightWide", "barHeightWide" },
        { "radiusBar", "radiusBar" },
        { "radiusChip", "radiusChip" },
        { "radiusControl", "radiusControl" },
        { "radiusCard", "radiusCard" },
        { "radiusPill", "radiusPill" },
        { "hairline", "hairline" },
        { "widthFloor", "widthFloor" },
        { "widthLaptop", "widthLaptop" },
        { "widthDrawn", "widthDrawn" },
        // The nine type roles collapse onto seven. `typeSecondary`,
        // `typeRow` and `typeBody` were three names for the two sizes
        // ordinary row text is set at, and `typeName` is what `strong`
        // already meant.
        { "typeTitle", "display" },
        { "typePage", "headline" },
        { "typeHeading", "title" },
        { "typeName", "strong" },
        { "typeBody", "body" },
        { "typeRow", "body" },
        { "typeSecondary", "small" },
        { "typeSmall", "small" },
        { "typeMicro", "caption" },
    };
    for (const auto &pair : mapping)
        QTest::newRow(qPrintable(pair.first)) << pair.first << pair.second;
}

void TestDensity::testEveryKvitHubTokenHasAHome()
{
    QFETCH(QString, property);
    InterfaceMetrics metrics;
    const QVariant value = metrics.property(property.toLatin1().constData());
    QVERIFY2(value.isValid(),
             qPrintable(QStringLiteral("InterfaceMetrics has no '%1', so "
                                       "kvit-hub's token has nowhere to go")
                            .arg(property)));
    QVERIFY(value.toInt() > 0);
}

void TestDensity::testTheFontFamiliesAreSeparateFromTheDocument()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore settings;
    QVERIFY(settings.open(dir.filePath(QStringLiteral("settings.json"))));

    InterfaceMetrics metrics;
    Typography typography;
    metrics.setSettings(&settings);
    typography.setSettings(&settings);

    // A person who sets their notes in a serif has not asked for a serif
    // sidebar, which is why the chrome family is not Typography's.
    typography.setFontFamily(QStringLiteral("Some Serif"));
    QCOMPARE(metrics.fontFamily(), QString());

    metrics.setFontFamily(QStringLiteral("Some Sans"));
    QCOMPARE(typography.fontFamily(), QStringLiteral("Some Serif"));

    // An empty monospace family would leave an identifier in the proportional
    // face, where two references down a column stop lining up.
    metrics.setMonoFamily(QString());
    QCOMPARE(metrics.monoFamily(), QStringLiteral("monospace"));
}

void TestDensity::testTheResolvedFamiliesAreDrawable()
{
    // The bug this exists to prevent, which shipped once and was found by
    // looking at a screenshot.
    //
    // `fontFamily` is empty by default and means "whatever this desktop
    // uses". Assigning that empty string to a QML `font.family` does not mean
    // that: Qt matches it against nothing and falls back to whichever
    // installed face its font matching lands on. On a Linux desktop with the
    // usual DejaVu set, that is `DejaVu Math TeX Gyre`, a serif maths face,
    // and every label in the library was drawn in it.
    //
    // Nothing about that looks like an error. A serif interface is a
    // plausible design, so no test that only checked the code would have
    // caught it. What can be checked is that the resolved family is the one
    // the desktop actually asked for.
    InterfaceMetrics metrics;

    QVERIFY(!metrics.resolvedFontFamily().isEmpty());
    QVERIFY(!metrics.resolvedMonoFamily().isEmpty());

    // With no preference set, the resolved family is the application's own
    // default, and matching it gives the same face the desktop would.
    const QString wanted = QFontInfo(QGuiApplication::font()).family();
    QCOMPARE(QFontInfo(QFont(metrics.resolvedFontFamily())).family(), wanted);

    // And the thing that went wrong: an empty family does *not* resolve to
    // that. If this ever stops being true the resolved accessors become
    // unnecessary, and this case is where that would be noticed.
    QFont empty;
    empty.setFamily(QString());
    QVERIFY2(QFontInfo(empty).family() != wanted
                 || QFontInfo(empty).family().isEmpty(),
             "an empty font family now resolves to the application default; "
             "if that holds on every platform the resolved accessors on "
             "InterfaceMetrics can go away");

    // A preference, once set, is what is drawn with.
    metrics.setFontFamily(QStringLiteral("DejaVu Sans Mono"));
    QCOMPARE(metrics.resolvedFontFamily(), QStringLiteral("DejaVu Sans Mono"));
    // Clearing it goes back to the desktop rather than to nothing.
    metrics.setFontFamily(QString());
    QCOMPARE(metrics.fontFamily(), QString());
    QCOMPARE(metrics.resolvedFontFamily(), QGuiApplication::font().family());

    // The monospace side has to actually be fixed-pitch, or a column of
    // identifiers stops lining up, which is the only reason it is a separate
    // family at all.
    QVERIFY(QFontInfo(QFont(metrics.resolvedMonoFamily())).fixedPitch());
}

QTEST_MAIN(TestDensity)
#include "test_density.moc"
