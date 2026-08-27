// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QDirIterator>
#include <QQmlComponent>
#include <QQmlEngine>
#include <QQuickItem>

#include "theme.h"
#include "uiservices.h"

// Every component in the module, loaded in all four themes, with no warning.
//
// This is the check that a component referring to a token that does not exist,
// or asking KvitIcon for a symbol with a typo in it, fails the build rather
// than drawing an empty box. It is worth stating why that matters more here
// than in an ordinary tree: the whole output of an interface library is what it
// draws, and a QML binding to a missing property does not throw — it silently
// evaluates to `undefined`, which a colour property renders as transparent and
// a numeric one as zero. A component with a typo in a token name is a
// component that renders as nothing at all, and nothing about the build says
// so.
//
// The component list is not written down. It is read out of the module's
// resource directory, so a component added to qml/CMakeLists.txt is a
// component this covers without anybody remembering to add it here too.
class TestComponents : public QObject
{
    Q_OBJECT

private slots:
    void initTestCase();
    void testEveryComponentLoads_data();
    void testEveryComponentLoads();
    void testAnUnknownIconNameIsLoud();
    void testTheIconFontIsInTheModule();

private:
    static QStringList componentUrls();
    // What a component needs before it can be built at all.
    //
    // A `required property` is how a component says a value is not optional —
    // KvitIcon without a symbol name is nothing, and a scroll bar without
    // something to scroll is nothing — and the price is that this test cannot
    // instantiate one bare. The table is short on purpose: a component that
    // needs three things before it will draw is a component with too many
    // required properties, and having to add a row here is where that gets
    // noticed.
    static QVariantMap requiredFor(const QString &name);
};

QVariantMap TestComponents::requiredFor(const QString &name)
{
    if (name == QLatin1String("KvitIcon.qml"))
        return { { QStringLiteral("name"), QStringLiteral("check") } };
    if (name == QLatin1String("KvitIconButton.qml")) {
        return { { QStringLiteral("symbol"), QStringLiteral("close") },
                 { QStringLiteral("label"), QStringLiteral("Close") } };
    }
    // The two data marks that require a scale. `maximum` is required on
    // purpose: a bar drawn against the largest value in its own list rescales
    // every time the list changes, so two screenshots taken a day apart are
    // not comparable and nothing on the screen says so.
    if (name == QLatin1String("KvitBar.qml"))
        return { { QStringLiteral("maximum"), 100.0 } };
    if (name == QLatin1String("KvitGauge.qml"))
        return { { QStringLiteral("allowance"), 100.0 } };
    // `allowNew` is required rather than defaulted because neither answer is
    // safe as a default: a tag picker should let the reader invent a tag, and
    // a category picker should not, since a typo would silently create a
    // second category beside the right one.
    if (name == QLatin1String("KvitTypeAhead.qml"))
        return { { QStringLiteral("allowNew"), false } };
    if (name == QLatin1String("KvitSidebarItem.qml"))
        return { { QStringLiteral("symbol"), QStringLiteral("folder") } };
    if (name == QLatin1String("KvitTable.qml")) {
        // A table with nothing behind it is a table with nothing to draw, and
        // `model` is required so that a view cannot forget to attach one.
        static QQmlEngine holder;
        QQmlComponent model(&holder);
        model.setData("import Kvit.Ui\nBenchmarkTableModel { totalRows: 50 }",
                      QUrl(QStringLiteral("qrc:/test/model.qml")));
        static QObject *made = model.create();
        return { { QStringLiteral("model"), QVariant::fromValue(made) } };
    }
    if (name == QLatin1String("KvitScrollBar.qml")) {
        // Built here rather than named, because the required value is an
        // object: a bare Flickable is enough for the bindings to resolve.
        static QQmlEngine holder;
        QQmlComponent flickable(&holder);
        flickable.setData("import QtQuick\nFlickable { width: 100; height: 100 }",
                          QUrl(QStringLiteral("qrc:/test/flickable.qml")));
        static QObject *made = flickable.create();
        return { { QStringLiteral("flickable"), QVariant::fromValue(made) } };
    }
    return {};
}

namespace {

// Every warning QML emits while a test is running. A component that loads and
// warns has still failed: the warning is Qt telling us a binding did not
// resolve.
QStringList g_warnings;
QtMessageHandler g_previous = nullptr;

void collect(QtMsgType type, const QMessageLogContext &context,
             const QString &message)
{
    if (type == QtWarningMsg || type == QtCriticalMsg)
        g_warnings.append(message);
    if (g_previous)
        g_previous(type, context, message);
}

}   // namespace

QStringList TestComponents::componentUrls()
{
    QStringList urls;
    QDirIterator it(QStringLiteral(":/qt/qml/Kvit/Ui"), { QStringLiteral("*.qml") },
                    QDir::Files);
    while (it.hasNext())
        urls.append(QStringLiteral("qrc") + it.next());
    urls.sort();
    return urls;
}

void TestComponents::initTestCase()
{
    g_previous = qInstallMessageHandler(collect);
    // A tree with no components in it would make every case below pass by
    // having nothing to check.
    QVERIFY2(componentUrls().size() >= 60,
             qPrintable(QStringLiteral("only %1 components found in the module "
                                       "resources; the QML files are not being "
                                       "compiled in")
                            .arg(componentUrls().size())));
}

void TestComponents::testEveryComponentLoads_data()
{
    QTest::addColumn<QString>("url");
    QTest::addColumn<QString>("theme");
    const QStringList urls = componentUrls();
    for (const QString &url : urls) {
        const QString name = url.section(QLatin1Char('/'), -1);
        for (const char *theme : { "light", "dark", "sepia", "highContrast" }) {
            QTest::newRow(qPrintable(name + QLatin1Char(' ')
                                     + QLatin1String(theme)))
                << url << QString::fromLatin1(theme);
        }
    }
}

void TestComponents::testEveryComponentLoads()
{
    QFETCH(QString, url);
    QFETCH(QString, theme);

    QQmlEngine engine;
    KvitUi::DefaultServices::theme()->setThemeId(theme);

    g_warnings.clear();
    QQmlComponent component(&engine, QUrl(url));
    QVERIFY2(component.isReady(),
             qPrintable(QStringLiteral("%1 did not compile:\n%2")
                            .arg(url, component.errorString())));

    const QString name = url.section(QLatin1Char('/'), -1);
    QScopedPointer<QObject> instance(
        component.createWithInitialProperties(requiredFor(name)));
    QVERIFY2(!instance.isNull(),
             qPrintable(QStringLiteral("%1 did not instantiate:\n%2")
                            .arg(url, component.errorString())));

    // Give it a size and let bindings settle. A component only asked for its
    // implicit size never evaluates the bindings that depend on width, which
    // is where a good half of the arithmetic in these files lives.
    if (auto *item = qobject_cast<QQuickItem *>(instance.data())) {
        item->setWidth(400);
        item->setHeight(120);
    }
    QCoreApplication::processEvents();

    QVERIFY2(g_warnings.isEmpty(),
             qPrintable(QStringLiteral("%1 warned in the %2 theme:\n  %3")
                            .arg(url, theme,
                                 g_warnings.join(QStringLiteral("\n  ")))));
}

void TestComponents::testAnUnknownIconNameIsLoud()
{
    // kvit-notes-pro's WorksIcon computes a `recognized` flag and then draws
    // nothing, so a typo in an icon name ships as an empty box that nobody
    // notices until somebody wonders where the button went. The shared version
    // warns and draws a marked placeholder, and this is what holds it to that.
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\nimport Kvit.Ui\n"
        "KvitIcon { name: \"no-such-symbol-anywhere\" }\n",
        QUrl(QStringLiteral("qrc:/test/unknown-icon.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    g_warnings.clear();
    QScopedPointer<QObject> instance(component.create());
    QVERIFY(!instance.isNull());
    QCOMPARE(instance->property("recognized").toBool(), false);
    QVERIFY2(!g_warnings.isEmpty(),
             "an unrecognised icon name drew nothing and said nothing");
    QVERIFY(g_warnings.join(QString()).contains(
        QStringLiteral("no-such-symbol-anywhere")));
}

void TestComponents::testTheIconFontIsInTheModule()
{
    // The font ships with the module rather than with an application, so
    // importing Kvit.Ui is enough. kvit-notes-pro adds a `qrc:/fonts` alias to
    // its own resource file today, which means a second consumer of KvitIcon
    // would have to know to do the same and would draw blank glyphs if it did
    // not.
    QVERIFY(QFile::exists(
        QStringLiteral(":/qt/qml/Kvit/Ui/fonts/Phosphor.ttf")));
    // And its licence beside it, which is what THIRD-PARTY-NOTICES.md promises
    // a consumer is shipping.
    QVERIFY(QFile::exists(
        QStringLiteral(":/qt/qml/Kvit/Ui/fonts/Phosphor-LICENSE.txt")));
}

QTEST_MAIN(TestComponents)
#include "test_components.moc"
