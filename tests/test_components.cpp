// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QAccessible>
#include <QDirIterator>
#include <QQmlComponent>
#include <QRegularExpression>
#include <QQmlEngine>
#include <QQuickWindow>
#include <QScreen>
#include <QQuickItem>

#include "interfacemetrics.h"
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
    void testEveryComponentHasANaturalWidth_data();
    void testEveryComponentHasANaturalWidth();
    void testButtonContentIsCentered_data();
    void testButtonContentIsCentered();
    void testTabContentIsCentered_data();
    void testTabContentIsCentered();
    void testLinksHintsAndKeyboardTooltips();
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
    if (name == QLatin1String("KvitHint.qml")) {
        return { { QStringLiteral("label"), QStringLiteral("About this value") },
                 { QStringLiteral("text"),
                   QStringLiteral("A longer explanation of this value.") } };
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

QString centeringError(QQuickItem *control)
{
    QQuickItem *row = control->findChild<QQuickItem *>(
        QStringLiteral("contentRow"));
    if (!row)
        return QStringLiteral("the natural-size content Row was not found");

    const QPointF origin = row->mapToItem(control, QPointF(0, 0));
    const qreal left = origin.x();
    const qreal right = control->width() - origin.x() - row->width();
    const qreal top = origin.y();
    const qreal bottom = control->height() - origin.y() - row->height();
    const QScreen *screen = QGuiApplication::primaryScreen();
    const qreal devicePixelRatio = screen ? screen->devicePixelRatio() : 1.0;

    const qreal horizontalError = qAbs(left - right) * devicePixelRatio;
    const qreal verticalError = qAbs(top - bottom) * devicePixelRatio;
    if (horizontalError <= 1.0 + 0.001 && verticalError <= 1.0 + 0.001)
        return {};

    return QStringLiteral(
               "content insets are L %1 / R %2 and T %3 / B %4 at DPR %5 "
               "(%6 horizontal and %7 vertical device pixels apart)")
        .arg(left).arg(right).arg(top).arg(bottom).arg(devicePixelRatio)
        .arg(horizontalError).arg(verticalError);
}

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

void TestComponents::testEveryComponentHasANaturalWidth_data()
{
    QTest::addColumn<QString>("url");
    for (const QString &url : componentUrls())
        QTest::newRow(qPrintable(url.section(QLatin1Char('/'), -1))) << url;
}

void TestComponents::testEveryComponentHasANaturalWidth()
{
    // No component's own implicit size may be read off its parent.
    //
    // `implicitWidth` means how wide something wants to be when nothing
    // constrains it. Thirteen components had `implicitWidth: parent.width`
    // instead, which says "fill my parent" — and inside anything that sizes
    // itself to its children (a Column, a Row, an Item measured by
    // childrenRect) that is a cycle: the child asks the parent how wide it is
    // while the parent is measuring itself from the child. All thirteen
    // settled at zero and drew nothing. KvitSlug was the one somebody noticed,
    // because its gallery page was the only one whose specimen did not happen
    // to pass an explicit width.
    //
    // Checked in the source rather than by measuring a built component,
    // because zero width is the right answer for several of these when they
    // are empty — a KvitLabel with no text, a KvitFigure with no value — and a
    // measurement cannot tell that apart from the defect. What is wrong is the
    // binding, so that is what is read.
    //
    // Root-level declarations only, at exactly four spaces of indentation.
    // Inside a nested object `parent` means that object's parent, which is
    // ordinary and correct: KvitCheck's tick is sized from the box it sits in.
    QFETCH(QString, url);

    // The url is "qrc" + the resource path, which is how componentUrls()
    // builds it; QFile wants the resource path.
    QFile file(url.mid(3));
    QVERIFY2(file.open(QIODevice::ReadOnly | QIODevice::Text),
             qPrintable(url));
    const QStringList lines =
        QString::fromUtf8(file.readAll()).split(QLatin1Char('\n'));

    static const QRegularExpression implicitSize(
        QStringLiteral("^    implicit(Width|Height):(.*)$"));

    QStringList offenders;
    for (const QString &line : lines) {
        const QRegularExpressionMatch match = implicitSize.match(line);
        if (!match.hasMatch())
            continue;
        if (match.captured(2).contains(QLatin1String("parent")))
            offenders.append(line.trimmed());
    }

    QVERIFY2(offenders.isEmpty(),
             qPrintable(QStringLiteral("%1 sizes itself from its parent, so it "
                                       "collapses inside anything that sizes "
                                       "itself to its children:\n  %2")
                            .arg(url.section(QLatin1Char('/'), -1),
                                 offenders.join(QStringLiteral("\n  ")))));
}

void TestComponents::testButtonContentIsCentered_data()
{
    QTest::addColumn<QString>("form");
    QTest::addColumn<QString>("symbol");
    QTest::addColumn<bool>("busy");
    QTest::addColumn<int>("interfaceSize");

    struct ContentState {
        const char *name;
        const char *symbol;
        bool busy;
    };
    const ContentState states[] = {
        { "text", "", false },
        { "icon-text", "calendar", false },
        { "busy-text", "", true },
    };
    for (int size : { InterfaceMetrics::MinFontSize,
                      InterfaceMetrics::DefaultFontSize,
                      InterfaceMetrics::MaxFontSize }) {
        for (const char *form : { "primary", "ordinary", "quiet" }) {
            for (const ContentState &state : states) {
                const QByteArray name = QByteArray(form) + '-' + state.name
                    + '-' + QByteArray::number(size);
                QTest::newRow(name.constData())
                    << QString::fromLatin1(form)
                    << QString::fromLatin1(state.symbol)
                    << state.busy << size;
            }
        }
    }
}

void TestComponents::testButtonContentIsCentered()
{
    QFETCH(QString, form);
    QFETCH(QString, symbol);
    QFETCH(bool, busy);
    QFETCH(int, interfaceSize);

    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(interfaceSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\nimport Kvit.Ui\n"
        "KvitButton { text: \"Save\"; busyText: \"Working…\" }\n",
        QUrl(QStringLiteral("qrc:/test/centred-button.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    QScopedPointer<QObject> instance(component.createWithInitialProperties({
        { QStringLiteral("form"), form },
        { QStringLiteral("symbol"), symbol },
        { QStringLiteral("busy"), busy },
    }));
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *control = qobject_cast<QQuickItem *>(instance.data());
    QVERIFY(control);
    QCoreApplication::processEvents();

    const QString error = centeringError(control);
    QVERIFY2(error.isEmpty(), qPrintable(error));
}

void TestComponents::testTabContentIsCentered_data()
{
    QTest::addColumn<int>("count");
    QTest::addColumn<int>("interfaceSize");

    for (int size : { InterfaceMetrics::MinFontSize,
                      InterfaceMetrics::DefaultFontSize,
                      InterfaceMetrics::MaxFontSize }) {
        QTest::newRow(qPrintable(QStringLiteral("text-%1").arg(size)))
            << -1 << size;
        QTest::newRow(qPrintable(QStringLiteral("counted-%1").arg(size)))
            << 47 << size;
    }
}

void TestComponents::testTabContentIsCentered()
{
    QFETCH(int, count);
    QFETCH(int, interfaceSize);

    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(interfaceSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\nimport Kvit.Ui\nKvitTab { text: \"Projects\" }\n",
        QUrl(QStringLiteral("qrc:/test/centred-tab.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    QScopedPointer<QObject> instance(component.createWithInitialProperties({
        { QStringLiteral("count"), count },
    }));
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *control = qobject_cast<QQuickItem *>(instance.data());
    QVERIFY(control);
    QCoreApplication::processEvents();

    const QString error = centeringError(control);
    QVERIFY2(error.isEmpty(), qPrintable(error));
}

void TestComponents::testLinksHintsAndKeyboardTooltips()
{
    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(
        InterfaceMetrics::DefaultFontSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(R"(
        import QtQuick
        import QtQuick.Controls
        import Kvit.Ui
        ApplicationWindow {
            id: window
            visible: true
            width: 480
            height: 240
            property int linkActivations: 0

            KvitLink {
                objectName: "link"
                x: 20; y: 20
                text: "Privacy policy"
                onActivated: window.linkActivations += 1
            }
            KvitIconButton {
                objectName: "iconButton"
                x: 20; y: 70
                symbol: "settings"
                label: "Settings"
            }
            KvitHint {
                objectName: "hint"
                x: 420; y: 20
                label: "About automatic matching"
                text: "Automatic matching compares the date, amount and reference."
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/interactive-primitives.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QTRY_VERIFY(window->isVisible());

    auto *link = window->findChild<QQuickItem *>(QStringLiteral("link"));
    auto *iconButton = window->findChild<QQuickItem *>(
        QStringLiteral("iconButton"));
    auto *hint = window->findChild<QQuickItem *>(QStringLiteral("hint"));
    QVERIFY(link);
    QVERIFY(iconButton);
    QVERIFY(hint);

    QAccessibleInterface *linkAccessible =
        QAccessible::queryAccessibleInterface(link);
    QVERIFY2(linkAccessible, "KvitLink has no accessible interface");
    QCOMPARE(linkAccessible->role(), QAccessible::Link);
    QCOMPARE(linkAccessible->text(QAccessible::Name),
             QStringLiteral("Privacy policy"));

    const qreal naturalWidth = link->implicitWidth();
    auto *linkLabel = link->findChild<QQuickItem *>(QStringLiteral("label"));
    QVERIFY(linkLabel);
    QVERIFY(linkLabel->width() + 1.0 >= linkLabel->implicitWidth());
    link->setWidth(60);
    QCoreApplication::processEvents();
    QVERIFY(naturalWidth > link->width());
    QVERIFY(linkLabel->width() <= link->width());

    link->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(link->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Return);
    QTRY_COMPARE(window->property("linkActivations").toInt(), 1);

    QObject *tooltip = iconButton->findChild<QObject *>(
        QStringLiteral("tooltip"));
    QVERIFY(tooltip);
    iconButton->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(iconButton->hasActiveFocus());
    QTRY_VERIFY(tooltip->property("visible").toBool());

    auto *trigger = hint->findChild<QQuickItem *>(QStringLiteral("trigger"));
    QVERIFY(trigger);
    QAccessibleInterface *triggerAccessible =
        QAccessible::queryAccessibleInterface(trigger);
    QVERIFY2(triggerAccessible, "KvitHint's trigger has no accessible interface");
    QCOMPARE(triggerAccessible->role(), QAccessible::Button);
    QCOMPARE(triggerAccessible->text(QAccessible::Name),
             QStringLiteral("About automatic matching"));

    trigger->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(trigger->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_VERIFY(hint->property("opened").toBool());
    QTest::keyClick(window, Qt::Key_Escape);
    QTRY_VERIFY(!hint->property("opened").toBool());
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
