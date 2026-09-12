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
    void testTheOverviewsNamesDrawWhatItMeans();
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

    // And no two meanings quietly point at the same glyph by accident. Two
    // names for one drawing where the estate thinks it has two drawings is a
    // bug the gallery would not catch, because both pages look right on their
    // own.
    //
    // Sharing a drawing on purpose is a different thing and is written down
    // here. A meaning name says what a call site means, and two call sites
    // can mean different things that are drawn the same way: an agent is
    // drawn as a robot and renaming is drawn as a pencil, and a screen about
    // agents that asks for `robot` has written down the drawing instead of
    // what it meant. Each pair below is one of those, with what separates the
    // two meanings; anything not on this list still fails.
    const QList<QPair<QString, QString>> sameDrawingOnPurpose = {
        // What is doing the work, against what the work is being done by.
        { QStringLiteral("agent"), QStringLiteral("robot") },
        // Giving a thing a new name, against editing what is inside it.
        { QStringLiteral("rename"), QStringLiteral("pencil") },
    };

    QSet<QString> allowed;
    for (const auto &pair : sameDrawingOnPurpose) {
        QVERIFY2(names.contains(pair.first) && names.contains(pair.second),
                 qPrintable(QStringLiteral("'%1' and '%2' are written down as "
                                           "sharing a drawing, but one of them "
                                           "is no longer a meaning name")
                                .arg(pair.first, pair.second)));
        QCOMPARE(catalog.glyph(pair.first), catalog.glyph(pair.second));
        allowed.insert(QStringLiteral("%1 and %2").arg(pair.first, pair.second));
        allowed.insert(QStringLiteral("%1 and %2").arg(pair.second, pair.first));
    }

    QHash<QString, QString> byGlyph;
    QStringList shared;
    for (const QString &name : names) {
        const QString glyph = catalog.glyph(name);
        if (byGlyph.contains(glyph)) {
            const QString pair = QStringLiteral("%1 and %2")
                                     .arg(byGlyph.value(glyph), name);
            if (!allowed.contains(pair))
                shared.append(pair);
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

void TestIcons::testTheOverviewsNamesDrawWhatItMeans()
{
    // kvit-notes-pro's overview names its symbols in C++ — `actionIcon` on a
    // section heading, `icon` on a row record — so the name a call site gets
    // is not something the call site can rewrite. Six of these resolved to
    // nothing when the overview was first drawn with KvitIcon, which is a
    // hatched red box on the screen and a console warning that takes the
    // shell's diagnostics gate red.
    //
    // The codepoint rather than only "it resolves to something", because what
    // the overview needs is a particular drawing: it and the web design it was
    // drawn from both show a branching line for `git` and a robot for `agent`,
    // and a name that resolved to some other symbol would pass a test that
    // only asked whether it resolved.
    KvitUi::IconCatalog catalog;
    const QList<QPair<QString, char32_t>> asked = {
        { QStringLiteral("diff"), 0xE27C },      // git-diff
        { QStringLiteral("git"), 0xE278 },       // git-branch
        { QStringLiteral("agent"), 0xE762 },     // robot
        { QStringLiteral("ask"), 0xE176 },       // chat-teardrop-dots
        { QStringLiteral("ask-in"), 0xE026 },    // arrow-bend-up-right
        { QStringLiteral("rename"), 0xE3B4 },    // pencil-simple
        // These two resolved before this, through the escape hatch that takes
        // a Phosphor name directly, and each landed on a glyph nobody chose:
        // `terminal` on the bare prompt at U+E47E rather than the window,
        // `chat` on the square bubble at U+E15C rather than the round one.
        { QStringLiteral("terminal"), 0xEAE8 },  // terminal-window
        { QStringLiteral("chat"), 0xE168 },      // chat-circle
    };
    for (const auto &pair : asked) {
        const QString glyph = catalog.glyph(pair.first);
        QVERIFY2(!glyph.isEmpty(), qPrintable(pair.first));
        QCOMPARE(glyph, QString::fromUcs4(&pair.second, 1));
    }

    // Naming `chat` as a meaning shadows Phosphor's own `chat`, the way
    // `link` and `pencil` already shadow theirs. The square bubble is still
    // reachable, under the name it always had here.
    const char32_t squareBubble = 0xE15C;
    QCOMPARE(catalog.glyph(QStringLiteral("message-square")),
             QString::fromUcs4(&squareBubble, 1));
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
