// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QSignalSpy>
#include <QTemporaryDir>

#include "interfacemetrics.h"
#include "settingsstore.h"

// The chrome's own type scale (accessibility.md Finding 4): one setting, five
// role sizes and a geometry scaler derived from it, persisted under
// `interface.` and clamped so a hand-edited settings file cannot produce an
// unusable layout.
class TestInterfaceMetrics : public QObject
{
    Q_OBJECT

private slots:
    void testDefaultReproducesTheOldLiterals();
    void testRolesAndGeometryScaleTogether();
    void testClamps();
    void testHairlinesNeverRoundToNothing();
    void testPersistsThroughSettings();
    void testResetToDefaults();
};

// The claim the whole migration rests on: at the default base every role
// answers the literal that was written out at the call site before this
// existed, so converting a pane changes no pixels until somebody moves the
// setting. A failure here means the default build has silently reflowed.
void TestInterfaceMetrics::testDefaultReproducesTheOldLiterals()
{
    InterfaceMetrics metrics;
    QCOMPARE(metrics.fontSize(), 12);
    QCOMPARE(metrics.scale(), 1.0);
    QCOMPARE(metrics.caption(), 10);
    QCOMPARE(metrics.small(), 11);
    QCOMPARE(metrics.body(), 12);
    QCOMPARE(metrics.strong(), 13);
    QCOMPARE(metrics.title(), 15);
    // Geometry is the identity at the default too, which is what makes
    // `implicitHeight: Interface.px(28)` a safe rewrite of `28`.
    for (int designPx : { 1, 2, 4, 6, 8, 10, 14, 16, 22, 24, 26, 28, 30, 60,
                          96, 148, 170, 480, 560 })
        QCOMPARE(metrics.px(designPx), designPx);
}

void TestInterfaceMetrics::testRolesAndGeometryScaleTogether()
{
    InterfaceMetrics metrics;
    QSignalSpy spy(&metrics, &InterfaceMetrics::changed);
    metrics.setFontSize(24);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(metrics.scale(), 2.0);
    QCOMPARE(metrics.body(), 24);
    QCOMPARE(metrics.caption(), 20);
    QCOMPARE(metrics.title(), 30);
    // The geometry travels with the type, which is the point: a 24-pixel
    // label inside a 28-pixel button clips, and a 56-pixel button does not.
    QCOMPARE(metrics.px(28), 56);

    // Setting the same size again is a no-op.
    metrics.setFontSize(24);
    QCOMPARE(spy.count(), 1);
}

void TestInterfaceMetrics::testClamps()
{
    InterfaceMetrics metrics;
    metrics.setFontSize(2);
    QCOMPARE(metrics.fontSize(), InterfaceMetrics::MinFontSize);
    metrics.setFontSize(400);
    QCOMPARE(metrics.fontSize(), InterfaceMetrics::MaxFontSize);
}

// A one-pixel rule scaled down and rounded reaches zero, and a separator that
// vanishes at the smallest interface size reads as a layout bug rather than
// as a smaller interface.
void TestInterfaceMetrics::testHairlinesNeverRoundToNothing()
{
    InterfaceMetrics metrics;
    metrics.setFontSize(InterfaceMetrics::MinFontSize);
    QVERIFY(metrics.px(1) >= 1);
    QVERIFY(metrics.px(2) >= 1);
    QCOMPARE(metrics.px(0), 0);
}

void TestInterfaceMetrics::testPersistsThroughSettings()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    const QString path = dir.filePath(QStringLiteral("settings.json"));
    {
        SettingsStore store;
        QVERIFY(store.open(path));
        InterfaceMetrics metrics;
        metrics.setSettings(&store);
        metrics.setFontSize(16);
        store.flush();
    }

    SettingsStore reopened;
    QVERIFY(reopened.open(path));
    InterfaceMetrics metrics;
    metrics.setSettings(&reopened);
    QCOMPARE(metrics.fontSize(), 16);

    // A hand-edited value out of range clamps on load exactly as a live one
    // does, rather than installing a layout nothing fits in.
    SettingsStore corrupt;
    QVERIFY(corrupt.open(dir.filePath(QStringLiteral("corrupt.json"))));
    corrupt.setValue(QStringLiteral("interface.fontSize"), 400);
    InterfaceMetrics clamped;
    clamped.setSettings(&corrupt);
    QCOMPARE(clamped.fontSize(), InterfaceMetrics::MaxFontSize);
}

void TestInterfaceMetrics::testResetToDefaults()
{
    InterfaceMetrics metrics;
    metrics.setFontSize(20);
    metrics.resetToDefaults();
    QCOMPARE(metrics.fontSize(), InterfaceMetrics::DefaultFontSize);
    QCOMPARE(metrics.body(), 12);
}

QTEST_MAIN(TestInterfaceMetrics)
#include "test_interfacemetrics.moc"
