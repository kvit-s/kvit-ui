// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QFontDatabase>
#include <QRawFont>

#include "iconcatalog.h"

// The icon catalogue: what a name resolves to, and what happens when it
// resolves to nothing.
//
// The table itself is generated from the font's published stylesheet by
// tools/generate-icon-catalog.py and checked against the font's cmap, so what
// is left to test here is the layer above it — the meaning names, which are
// written by hand — and the two properties the estate depends on: that every
// meaning name draws something, and that an unknown name does not.
class TestIcons : public QObject
{
    Q_OBJECT

private slots:
    void testEveryMeaningNameResolves();
    void testTheFontHasAGlyphForEveryMeaningName();
    void testPhosphorNamesResolveDirectly();
    void testAnUnknownNameResolvesToNothing();
    void testTheEditorsLiteralsAllHaveSomewhereToGo();
    void testTheCatalogueIsBigEnoughToBeTheWholeFont();
};

void TestIcons::testEveryMeaningNameResolves()
{
    KvitUi::IconCatalog catalog;
    const QStringList names = KvitUi::IconCatalog::meaningNames();
    QVERIFY(names.size() > 40);

    for (const QString &name : names) {
        const QString glyph = catalog.glyph(name);
        QVERIFY2(!glyph.isEmpty(),
                 qPrintable(QStringLiteral("the meaning name '%1' resolves to "
                                           "nothing, so every call site using "
                                           "it draws a placeholder")
                                .arg(name)));
    }

    // And no two meanings quietly point at the same glyph by accident. A
    // duplicate is allowed — `check` and `success` could reasonably share one
    // — but two names for the same drawing where the estate thinks it has two
    // drawings is a bug the gallery would not catch, because both pages look
    // right on their own.
    QHash<QString, QString> byGlyph;
    QStringList shared;
    for (const QString &name : names) {
        const QString glyph = catalog.glyph(name);
        if (byGlyph.contains(glyph)) {
            shared.append(QStringLiteral("%1 and %2")
                              .arg(byGlyph.value(glyph), name));
        }
        byGlyph.insert(glyph, name);
    }
    QVERIFY2(shared.isEmpty(),
             qPrintable(QStringLiteral("these meaning names draw the same "
                                       "glyph: %1")
                            .arg(shared.join(QStringLiteral(", ")))));
}

void TestIcons::testTheFontHasAGlyphForEveryMeaningName()
{
    // The catalogue can only promise that the stylesheet named a codepoint.
    // This asks the shipped font whether it actually has a glyph there, which
    // is the difference between a symbol that draws and one that comes out as
    // the font's fallback box.
    const int id = QFontDatabase::addApplicationFont(
        QStringLiteral(":/qt/qml/Kvit/Ui/fonts/Phosphor.ttf"));
    QVERIFY2(id >= 0, "the icon font is not in the module's resources");

    const QStringList families = QFontDatabase::applicationFontFamilies(id);
    QVERIFY(!families.isEmpty());
    QRawFont face = QRawFont::fromFont(QFont(families.first(), 18));
    QVERIFY(face.isValid());

    KvitUi::IconCatalog catalog;
    const QStringList names = KvitUi::IconCatalog::meaningNames();
    QStringList missing;
    for (const QString &name : names) {
        const QString glyph = catalog.glyph(name);
        if (glyph.isEmpty() || !face.supportsCharacter(glyph.toUcs4().first()))
            missing.append(name);
    }
    QVERIFY2(missing.isEmpty(),
             qPrintable(QStringLiteral("the font has no glyph for: %1")
                            .arg(missing.join(QStringLiteral(", ")))));
}

void TestIcons::testPhosphorNamesResolveDirectly()
{
    // Asking for a Phosphor name works and is the escape hatch for a symbol
    // that has not been given a meaning yet. What a call site *should* write
    // is a meaning name, which is why those are what the vocabulary skill
    // publishes.
    KvitUi::IconCatalog catalog;
    QCOMPARE(catalog.glyph(QStringLiteral("caret-right")),
             catalog.glyph(QStringLiteral("chevron-right")));
    QVERIFY(!catalog.glyph(QStringLiteral("acorn")).isEmpty());
}

void TestIcons::testAnUnknownNameResolvesToNothing()
{
    // Empty rather than a placeholder character. KvitIcon is what turns that
    // into something visible, and tests/test_components.cpp is what holds it
    // to warning; the catalogue's job is only to be honest that it does not
    // know the name.
    KvitUi::IconCatalog catalog;
    QVERIFY(catalog.glyph(QStringLiteral("no-such-symbol")).isEmpty());
    QVERIFY(catalog.glyph(QString()).isEmpty());
    // Near misses, which is what a typo actually looks like.
    QVERIFY(catalog.glyph(QStringLiteral("chevron_right")).isEmpty());
    QVERIFY(catalog.glyph(QStringLiteral("Chevron-Right")).isEmpty());
}

void TestIcons::testTheEditorsLiteralsAllHaveSomewhereToGo()
{
    // kvit-notes passes literal Unicode characters to its IconButton today —
    // `›`, `‹`, `×`, `+`, `▰`, `▤`, `⌕`, `N`, `#` — drawn in whatever font on
    // the machine happens to have them. Wave 5 moves those call sites onto
    // named symbols, and this is the list it moves them to. A name missing
    // here is a call site with nowhere to go.
    KvitUi::IconCatalog catalog;
    const QStringList destinations = {
        QStringLiteral("chevron-right"),   // ›
        QStringLiteral("chevron-left"),    // ‹
        QStringLiteral("close"),           // ×
        QStringLiteral("plus"),            // +
        QStringLiteral("sidebar"),         // ▰
        QStringLiteral("list"),            // ▤
        QStringLiteral("search"),          // ⌕
        QStringLiteral("note"),            // N
        QStringLiteral("tag"),             // #
    };
    for (const QString &name : destinations)
        QVERIFY2(!catalog.glyph(name).isEmpty(), qPrintable(name));

    // And kvit-notes-pro's twelve, unchanged, so its six call sites keep
    // working when it drops WorksIcon.qml.
    const QStringList worksIcon = {
        QStringLiteral("plus"), QStringLiteral("messages-square"),
        QStringLiteral("message-square"), QStringLiteral("send"),
        QStringLiteral("chevron-right"), QStringLiteral("chevron-down"),
        QStringLiteral("zoom-in"), QStringLiteral("zoom-out"),
        QStringLiteral("pencil"), QStringLiteral("check"),
        QStringLiteral("archive"), QStringLiteral("rotate-ccw"),
    };
    for (const QString &name : worksIcon)
        QVERIFY2(!catalog.glyph(name).isEmpty(), qPrintable(name));
}

void TestIcons::testTheCatalogueIsBigEnoughToBeTheWholeFont()
{
    // A generator that silently produced a short table would leave the
    // meaning names working and everything else missing, which is exactly the
    // state kvit-notes-pro's hand-written switch of twelve was in.
    QVERIFY(KvitUi::IconCatalog::glyphNames().size() > 1000);
}

QTEST_MAIN(TestIcons)
#include "test_icons.moc"
