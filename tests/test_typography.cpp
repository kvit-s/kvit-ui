// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QSignalSpy>
#include <QTemporaryDir>

#include "typography.h"
#include "settingsstore.h"

// Typography settings: the frozen ratio type scale, clamped setters, and
// persistence.
//
// kvit-notes' copy of this file also drives BlockEditorEngine, checking that
// a line-height or mono-family change reaches the rendered document. That
// class is the editor's document renderer and stays there; what is left here
// is every case that is about the type scale itself.
class TestTypography : public QObject
{
    Q_OBJECT

private slots:
    void testDefaultsAndTheScaleTheyProduce();
    void testScaleDerivesFromBase_data();
    void testScaleDerivesFromBase();
    void testClamps();
    void testSignalDiscipline();
    void testPersistsThroughSettings();
    void testCorruptSettingsClampOnLoad();
    void testMonospaceFamilies();
    void testResetToDefaults();
};

void TestTypography::testDefaultsAndTheScaleTheyProduce()
{
    // What a reader who has never opened the settings dialog is given, and
    // the document that follows from it.
    Typography t;
    QCOMPARE(t.baseSize(), Typography::DefaultBaseSize);
    QCOMPARE(t.baseSize(), 14);
    QCOMPARE(t.lineHeight(), Typography::DefaultLineHeight);
    QCOMPARE(t.lineHeight(), 1.3);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading1)), 30);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading2)), 22);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading3)), 19);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading4)), 16);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Mono)), 12);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Body)), 14);
    // Everything that is not a heading and not code is body text, and asks
    // for it by name rather than by being absent from a switch.
    QCOMPARE(t.bodySize(), 14);
    QCOMPARE(t.monoSize(), 12);
    QCOMPARE(t.paragraphSpacing(), Typography::DefaultParagraphSpacing);
    QCOMPARE(t.maxContentWidth(), 0);
    QCOMPARE(t.monoFamily(), QString("monospace"));
    QCOMPARE(t.fontFamily(), QString());

    // The scale's shape is fixed independently of where the base sits: the
    // ratios are written against 15, so a base of 15 still reproduces the
    // sizes they were measured from, to the pixel.
    t.setBaseSize(15);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading1)), 32);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading2)), 24);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading3)), 20);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading4)), 17);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Mono)), 13);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Body)), 15);
}

void TestTypography::testScaleDerivesFromBase_data()
{
    QTest::addColumn<int>("base");
    QTest::addColumn<int>("h1");
    QTest::addColumn<int>("code");
    // qRound(base * ratio): the scale stays coherent at any base.
    QTest::newRow("base 12") << 12 << 26 << 10;  // 25.6, 10.4
    QTest::newRow("base 18") << 18 << 38 << 16;  // 38.4, 15.6
    QTest::newRow("base 20") << 20 << 43 << 17;  // 42.67, 17.33
}

void TestTypography::testScaleDerivesFromBase()
{
    QFETCH(int, base);
    QFETCH(int, h1);
    QFETCH(int, code);
    Typography t;
    t.setBaseSize(base);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Heading1)), h1);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Mono)), code);
    QCOMPARE(t.sizeForRole(int(Typography::FontRole::Body)), base);
    // Order always holds: H1 > H2 > H3 > H4 >= body > code... at least
    // never inverted.
    QVERIFY(t.sizeForRole(int(Typography::FontRole::Heading1))
            > t.sizeForRole(int(Typography::FontRole::Heading2)));
    QVERIFY(t.sizeForRole(int(Typography::FontRole::Heading2))
            > t.sizeForRole(int(Typography::FontRole::Heading3)));
    QVERIFY(t.sizeForRole(int(Typography::FontRole::Heading3))
            > t.sizeForRole(int(Typography::FontRole::Heading4)));
    QVERIFY(t.sizeForRole(int(Typography::FontRole::Heading4))
            >= t.sizeForRole(int(Typography::FontRole::Body)));
    QVERIFY(t.sizeForRole(int(Typography::FontRole::Body))
            > t.sizeForRole(int(Typography::FontRole::Mono)));
}

void TestTypography::testClamps()
{
    Typography t;
    t.setBaseSize(5);
    QCOMPARE(t.baseSize(), Typography::MinBaseSize);
    t.setBaseSize(99);
    QCOMPARE(t.baseSize(), Typography::MaxBaseSize);
    t.setLineHeight(0.3);
    QCOMPARE(t.lineHeight(), Typography::MinLineHeight);
    t.setLineHeight(5.0);
    QCOMPARE(t.lineHeight(), Typography::MaxLineHeight);
    t.setParagraphSpacing(-4);
    QCOMPARE(t.paragraphSpacing(), 0);
    t.setParagraphSpacing(400);
    QCOMPARE(t.paragraphSpacing(), Typography::MaxParagraphSpacing);
    t.setMaxContentWidth(50);   // too narrow to be usable
    QCOMPARE(t.maxContentWidth(), Typography::MinContentWidth);
    t.setMaxContentWidth(0);    // 0 = off is always legal
    QCOMPARE(t.maxContentWidth(), 0);
}

void TestTypography::testSignalDiscipline()
{
    Typography t;
    QSignalSpy spy(&t, &Typography::typographyChanged);
    t.setBaseSize(18);
    QCOMPARE(spy.count(), 1);
    t.setBaseSize(18);          // unchanged: silent
    QCOMPARE(spy.count(), 1);
    t.setBaseSize(99);          // clamps to 28: one change
    t.setBaseSize(40);          // clamps to 28 again: silent
    QCOMPARE(spy.count(), 2);
}

void TestTypography::testPersistsThroughSettings()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    const QString path = dir.filePath("settings.json");

    {
        SettingsStore store;
        QVERIFY(store.open(path));
        Typography t;
        t.setSettings(&store);
        t.setFontFamily("Serif Family");
        t.setBaseSize(18);
        t.setLineHeight(1.4);
        t.setParagraphSpacing(14);
        t.setMaxContentWidth(700);
        t.setMonoFamily("Mono Family");
        store.flush();
    }

    SettingsStore reopened;
    QVERIFY(reopened.open(path));
    Typography t;
    t.setSettings(&reopened);
    QCOMPARE(t.fontFamily(), QString("Serif Family"));
    QCOMPARE(t.baseSize(), 18);
    QCOMPARE(t.lineHeight(), 1.4);
    QCOMPARE(t.paragraphSpacing(), 14);
    QCOMPARE(t.maxContentWidth(), 700);
    QCOMPARE(t.monoFamily(), QString("Mono Family"));
}

void TestTypography::testCorruptSettingsClampOnLoad()
{
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    SettingsStore store;
    QVERIFY(store.open(dir.filePath("settings.json")));
    store.setValue("typography.fontSize", 400);
    store.setValue("typography.lineHeight", -3);
    store.setValue("typography.maxContentWidth", 10);

    Typography t;
    t.setSettings(&store);
    QCOMPARE(t.baseSize(), Typography::MaxBaseSize);
    QCOMPARE(t.lineHeight(), Typography::MinLineHeight);
    QCOMPARE(t.maxContentWidth(), Typography::MinContentWidth);
}

void TestTypography::testMonospaceFamilies()
{
    Typography t;
    const QStringList families = t.monospaceFamilies();
    QVERIFY(!families.isEmpty());
    // The generic alias is always offered (it is the default value).
    QVERIFY(families.contains("monospace"));
}

void TestTypography::testResetToDefaults()
{
    // Against a freshly constructed object rather than against a list of
    // numbers: what the button promises is the state the application starts
    // in, and asserting the two agree is what catches one of the two places
    // the defaults are written being changed without the other.
    const Typography fresh;
    Typography t;
    t.setBaseSize(20);
    t.setLineHeight(1.6);
    t.setParagraphSpacing(20);
    t.setMaxContentWidth(600);
    t.resetToDefaults();
    QCOMPARE(t.baseSize(), fresh.baseSize());
    QCOMPARE(t.lineHeight(), fresh.lineHeight());
    QCOMPARE(t.baseSize(), Typography::DefaultBaseSize);
    QCOMPARE(t.lineHeight(), Typography::DefaultLineHeight);
    QCOMPARE(t.paragraphSpacing(), Typography::DefaultParagraphSpacing);
    QCOMPARE(t.maxContentWidth(), 0);
    QCOMPARE(t.monoFamily(), QString("monospace"));
    QCOMPARE(t.fontFamily(), fresh.fontFamily());
}

QTEST_MAIN(TestTypography)
#include "test_typography.moc"
