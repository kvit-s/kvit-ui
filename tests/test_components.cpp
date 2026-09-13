// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QAccessible>
#include <QDirIterator>
#include <QLocale>
#include <QScopeGuard>
#include <QQmlComponent>
#include <QRegularExpression>
#include <QQmlEngine>
#include <QQmlExpression>
#include <QQmlContext>
#include <QQuickWindow>

#include <functional>
#include <QScreen>
#include <QQuickItem>
#include <QQuickWindow>
#include <QSignalSpy>

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
    void testACountReadsAsASentenceRatherThanAFormField();
    void testASuggestionDrawsItsOwnLabel();
    void testTheRailShowsACountOnlyWhenAskedTo();
    void testACountedNameInflectsTheCallersNoun();
    void testAGroupHeadingCountsItsGroupTheWayEverythingElseDoes();
    void testAProgressBarCannotTrailTheNumberBesideIt();
    void testARowSaysWhetherItIsPressableBeforeItIsPressed();
    void testARowAnnouncesItselfAsWhatItActuallyIs();
    void testAShortenedValueIsDisclosedToTheKeyboardAsWellAsThePointer();
    void testARowOpensFromTheKeyboardAndOnlyOnce();
    void testASectionHeadingOpensFromTheKeyboardWithoutSwallowingItsAction();
    void testAnIconButtonDrawsTheSymbolSizeItWasAskedFor();
    void testAControlCarriesOneSentenceOfExplanation();
    void testAChipThatActsIsAControlAndSaysWhyItCannot();
    void testAChipCanSayItIsTheOneAlreadyOpen();
    void testAHeadingTakesAWrittenCountAndASymbolicAction();
    void testAChipsLabelGivesWayRatherThanClipping_data();
    void testAChipsLabelGivesWayRatherThanClipping();
    void testAnEmptySectionSaysSoOnOneLine_data();
    void testAnEmptySectionSaysSoOnOneLine();
    void testTheStatusBarKeepsWhatDoesNotFitReachable();
    void testTheStatusBarCanPutItsGroupsFirst();
    void testTheStatusBarSurvivesItsGroupsBeingReplaced();

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
    // A mark whose whole statement is a colour and a shape has nothing else
    // to tell a screen reader, which is why the words are required.
    if (name == QLatin1String("KvitSignal.qml"))
        return { { QStringLiteral("label"), QStringLiteral("2 running") } };
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

    QObject *hintTooltip = trigger->findChild<QObject *>(QStringLiteral("tooltip"));
    QVERIFY(hintTooltip);
    trigger->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(trigger->hasActiveFocus());
    QVERIFY(!hintTooltip->property("visible").toBool());
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_VERIFY(hint->property("opened").toBool());

    auto *hintTitle = hint->findChild<QQuickItem *>(QStringLiteral("title"));
    auto *hintDetail = hint->findChild<QQuickItem *>(QStringLiteral("detail"));
    QVERIFY(hintTitle);
    QVERIFY(hintDetail);
    QTRY_VERIFY(hintTitle->isVisible());
    QTRY_VERIFY(hintDetail->isVisible());
    QVERIFY(hintTitle->width() > 0);
    QVERIFY(hintTitle->height() > 0);
    QVERIFY(hintDetail->width() > 0);
    QVERIFY(hintDetail->height() > 0);
    QVERIFY(hintDetail->mapToItem(hint, QPointF()).y()
            > hintTitle->mapToItem(hint, QPointF(0, hintTitle->height())).y());
    QCOMPARE(hintTitle->property("text").toString(),
             QStringLiteral("About automatic matching"));
    QCOMPARE(hintDetail->property("text").toString(),
             QStringLiteral("Automatic matching compares the date, amount and reference."));

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

// What a suggestion in a type-ahead draws.
//
// The list's delegate is a KvitRow, and KvitRow's default property reparents
// whatever a caller puts inside it into its own inner layout. The label inside
// the delegate read `parent.modelData`, which is that layout rather than the
// delegate, so every suggestion in an open list drew the word `undefined`.
// The delegate's `label`, which is what a screen reader is given, was right the
// whole time, so nothing that read the accessible name could see it.
void TestComponents::testASuggestionDrawsItsOwnLabel()
{
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    width: 400; height: 300\n"
        "    property alias picker: picker\n"
        "    KvitTypeAhead { id: picker; allowNew: false; width: 200\n"
        "                    label: \"Category\"\n"
        "                    source: [{ value: \"1\", label: \"Groceries\" },\n"
        "                             { value: \"2\", label: \"Green fees\" },\n"
        "                             { value: \"3\", label: \"Rent\" }] }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/suggestion.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    QObject *picker = holder->property("picker").value<QObject *>();
    QVERIFY(picker);

    // The list's rows exist only once the popup is drawn, so this needs a
    // window. The gate runs offscreen, which has one.
    if (QGuiApplication::platformName() == QLatin1String("minimal"))
        return;
    QQuickWindow window;
    qobject_cast<QQuickItem *>(holder.data())->setParentItem(window.contentItem());
    window.resize(400, 300);
    window.show();
    QVERIFY(QTest::qWaitForWindowExposed(&window));

    QVERIFY(picker->setProperty("text", QStringLiteral("gr")));
    const QVariantList matched = picker->property("matches").toList();
    QCOMPARE(matched.size(), 2);

    // The popup builds its rows when it opens, so the list exists only after
    // the event loop has run.
    QObject *popup = nullptr;
    const QList<QObject *> children = picker->findChildren<QObject *>();
    for (QObject *child : children) {
        if (QString::fromLatin1(child->metaObject()->className())
                .startsWith(QStringLiteral("KvitPopover")))
            popup = child;
    }
    QVERIFY2(popup, "the type-ahead holds no suggestion popup");
    QTRY_VERIFY2(popup->property("visible").toBool(),
                 "the suggestion popup did not open with two matches");
    QTest::qWait(50);

    // Every piece of text the open list draws. Two of them are the matched
    // labels; none of them is the word an undefined value prints as. Read
    // through the property rather than through QQuickText, which is private
    // Qt: what is being checked is what a component put on the screen, and
    // every drawn label answers to `text`.
    // Walked down the drawing tree rather than through QObject parentage: a
    // popup's content is reparented into the window's overlay, and a view's
    // delegates are not QObject children of the view that made them.
    QStringList drawn;
    std::function<void(QQuickItem *)> walk = [&](QQuickItem *item) {
        if (!item || !item->isVisible())
            return;
        const QVariant text = item->property("text");
        if (text.isValid() && !text.toString().isEmpty())
            drawn.append(text.toString());
        const QList<QQuickItem *> below = item->childItems();
        for (QQuickItem *child : below)
            walk(child);
    };
    walk(window.contentItem());
    QVERIFY2(!drawn.contains(QStringLiteral("undefined")),
             qPrintable(QStringLiteral("the suggestion list drew: %1")
                            .arg(drawn.join(QStringLiteral(" | ")))));
    QVERIFY2(drawn.contains(QStringLiteral("Groceries")),
             qPrintable(QStringLiteral("the suggestion list drew: %1")
                            .arg(drawn.join(QStringLiteral(" | ")))));
    QVERIFY2(drawn.contains(QStringLiteral("Green fees")),
             qPrintable(QStringLiteral("the suggestion list drew: %1")
                            .arg(drawn.join(QStringLiteral(" | ")))));
}

void TestComponents::testACountReadsAsASentenceRatherThanAFormField()
{
    // What a reader is told about how much a view holds, and how much a
    // filter left.
    //
    // Two things go wrong with Qt's %n on its own, and both of them reached a
    // screen. With no translator installed Qt substitutes the number into the
    // one source string and does not choose a plural form, so "%n result(s)"
    // is read out as "1 result(s)"; and %n writes the bare integer, so a
    // quarter of a million matches arrives as "250000" and is read digit by
    // digit. Neither shows up in a screenshot, because both of these strings
    // exist only for a screen reader.

    // A fixed locale, so what follows is a fact about the components rather
    // than a property of the machine the suite runs on: the C locale groups
    // nothing, and a build run under LC_ALL=C would otherwise fail here for a
    // reason nobody could act on. Restored at the end, because the default
    // locale is process-wide and the cases after this one share the process.
    const QLocale previous;
    QLocale::setDefault(QLocale(QLocale::English, QLocale::UnitedStates));
    const auto restore = qScopeGuard([previous] { QLocale::setDefault(previous); });

    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    property alias head: head\n"
        "    property alias field: field\n"
        "    KvitViewHead { id: head; title: \"Transactions\"\n"
        "                   counted: \"transaction\" }\n"
        "    KvitSearchField { id: field; matchedNoun: \"transaction\" }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/counted.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    QObject *head = holder->property("head").value<QObject *>();
    QObject *field = holder->property("field").value<QObject *>();
    QVERIFY(head);
    QVERIFY(field);

    const QString grouped = QStringLiteral("250,000");

    QVERIFY(head->setProperty("count", 1));
    QCOMPARE(head->property("countPhrase").toString(),
             QStringLiteral("1 transaction"));
    QVERIFY(head->setProperty("count", 250000));
    QCOMPARE(head->property("countPhrase").toString(),
             grouped + QStringLiteral(" transactions"));

    // A noun that does not pluralise by suffixing an s says its own plural,
    // and the component uses it rather than guessing.
    QVERIFY(head->setProperty("countedPlural", QStringLiteral("entries")));
    QCOMPARE(head->property("countPhrase").toString(),
             grouped + QStringLiteral(" entries"));

    // Below zero says nothing at all, which is right for a view whose
    // contents are not a countable list.
    QVERIFY(head->setProperty("count", -1));
    QCOMPARE(head->property("countPhrase").toString(), QString());

    // The filter field says the same thing in the same shape, which is the
    // half that was left standing when the view head was fixed.
    QVERIFY(field->setProperty("matches", 1));
    QCOMPARE(field->property("matchPhrase").toString(),
             QStringLiteral("1 transaction"));
    QVERIFY(field->setProperty("matches", 250000));
    QCOMPARE(field->property("matchPhrase").toString(),
             grouped + QStringLiteral(" transactions"));
    QVERIFY(field->setProperty("matchedNounPlural", QStringLiteral("entries")));
    QCOMPARE(field->property("matchPhrase").toString(),
             grouped + QStringLiteral(" entries"));
    QVERIFY(field->setProperty("matches", -1));
    QCOMPARE(field->property("matchPhrase").toString(), QString());
}

namespace {

// The KvitBadge inside a KvitSidebarItem. It has no object name a test can
// reach, so it is found by type: a QML-defined type's class name is its file
// name with a suffix on it.
QQuickItem *findBadge(QQuickItem *item)
{
    const auto children = item->childItems();
    for (QQuickItem *child : children) {
        if (QByteArray(child->metaObject()->className()).startsWith("KvitBadge"))
            return child;
        if (QQuickItem *found = findBadge(child))
            return found;
    }
    return nullptr;
}

}   // namespace

void TestComponents::testTheRailShowsACountOnlyWhenAskedTo()
{
    // What a collapsed sidebar draws belongs to the estate rather than to one
    // application.
    //
    // kvit-cash's shell keeps every item's count in the rail; the other three
    // applications draw a strip of symbols there and were drawing it that way
    // before this property existed. So the default is what they already had
    // and an application opts in, rather than every rail in the estate
    // changing because one screen asked for it.
    //
    // The cap goes the same way. KvitBadge writes "99+" past `max`, which is
    // right where the number is only a signal that there is a backlog and
    // wrong where the place the item points at states the true figure: the
    // two then disagree about the same thing in the same window.
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "KvitSidebarItem { text: \"Review\"; symbol: \"question\"; count: 214 }\n",
        QUrl(QStringLiteral("qrc:/test/rail-count.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    auto *item = qobject_cast<QQuickItem *>(holder.data());
    QVERIFY(item);
    item->setWidth(240);
    item->setHeight(30);
    QCoreApplication::processEvents();

    QQuickItem *badge = findBadge(item);
    QVERIFY2(badge, "KvitSidebarItem has no KvitBadge in it");

    // Expanded, the count sits at the end of the label's line.
    QVERIFY(badge->isVisible());

    // Collapsed to the rail it is gone, which is what every sidebar in the
    // estate drew before kvit-cash needed otherwise.
    QVERIFY(item->setProperty("collapsed", true));
    QCoreApplication::processEvents();
    QVERIFY2(!badge->isVisible(),
             "a collapsed sidebar showed its badge without being asked to");

    QVERIFY(item->setProperty("countInRail", true));
    QCoreApplication::processEvents();
    QVERIFY(badge->isVisible());

    // And the cap is the caller's.
    QCOMPARE(badge->property("max").toInt(), 99);
    QVERIFY(item->setProperty("countMax", 999));
    QCoreApplication::processEvents();
    QCOMPARE(badge->property("max").toInt(), 999);
}

namespace {

// The first descendant item whose type name starts with `prefix`. A
// QML-defined type's class name is its file name with a suffix on it, which is
// how a component with no object name is found from C++.
QQuickItem *findByType(QQuickItem *item, const char *prefix)
{
    const auto children = item->childItems();
    for (QQuickItem *child : children) {
        if (QByteArray(child->metaObject()->className()).startsWith(prefix))
            return child;
        if (QQuickItem *found = findByType(child, prefix))
            return found;
    }
    return nullptr;
}

}   // namespace

void TestComponents::testACountedNameInflectsTheCallersNoun()
{
    // What a sidebar item and a badge say about a count, at one and at many.
    //
    // Neither component can know what a count counts, so the noun comes from
    // the caller. One already-inflected word is right at exactly one count and
    // wrong at every other: a caller that hands over "decisions" makes a
    // workspace with a single decision waiting announce its sidebar item as
    // "Review, 1 decisions". So both take the singular and, where English's
    // suffixed s is wrong, the plural as well — the same two slots KvitViewHead
    // takes — and this is what walks that branch. Nothing did before, and the
    // fault is invisible in a screenshot because these strings exist only for
    // a screen reader.
    //
    // A fixed locale, for the reason the count-phrase case above fixes one:
    // the C locale groups nothing, and the grouping is half of what is
    // asserted here.
    const QLocale previous;
    QLocale::setDefault(QLocale(QLocale::English, QLocale::UnitedStates));
    const auto restore = qScopeGuard([previous] { QLocale::setDefault(previous); });

    QQmlEngine engine;
    QQmlComponent component(&engine);
    // The names are read through properties of the holder, because an
    // attached property is not a property of the item a test can ask for by
    // name.
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    property alias item: item\n"
        "    property alias badge: badge\n"
        "    readonly property string itemName: item.Accessible.name\n"
        "    readonly property string badgeName: badge.Accessible.name\n"
        "    KvitSidebarItem {\n"
        "        id: item\n"
        "        text: \"Review\"; symbol: \"question\"\n"
        "        count: 1; counted: \"decision\"\n"
        "    }\n"
        "    KvitBadge { id: badge; count: 1; counted: \"decision\"; max: 250000 }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/counted-name.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    auto *item = holder->property("item").value<QQuickItem *>();
    auto *badge = holder->property("badge").value<QQuickItem *>();
    QVERIFY(item);
    QVERIFY(badge);
    QCoreApplication::processEvents();

    // One of them, said as one of them.
    QCOMPARE(holder->property("itemName").toString(),
             QStringLiteral("Review, 1 decision"));
    QCOMPARE(holder->property("badgeName").toString(),
             QStringLiteral("1 decision"));

    // Several, with the digits grouped the way the reader's locale groups
    // them: a six-figure count run together is read out digit by digit.
    QVERIFY(item->setProperty("count", 250000));
    QVERIFY(badge->setProperty("count", 250000));
    QCoreApplication::processEvents();
    QCOMPARE(holder->property("itemName").toString(),
             QStringLiteral("Review, 250,000 decisions"));
    QCOMPARE(holder->property("badgeName").toString(),
             QStringLiteral("250,000 decisions"));

    // A noun English does not pluralise by suffixing an s says its own plural,
    // and both components use it rather than guessing.
    QVERIFY(item->setProperty("counted", QStringLiteral("entry")));
    QVERIFY(item->setProperty("countedPlural", QStringLiteral("entries")));
    QVERIFY(badge->setProperty("counted", QStringLiteral("entry")));
    QVERIFY(badge->setProperty("countedPlural", QStringLiteral("entries")));
    QCoreApplication::processEvents();
    QCOMPARE(holder->property("itemName").toString(),
             QStringLiteral("Review, 250,000 entries"));
    QCOMPARE(holder->property("badgeName").toString(),
             QStringLiteral("250,000 entries"));

    // The item hands both slots to the badge inside it, so the two do not
    // disagree about the same count in the same place.
    QQuickItem *inner = findBadge(item);
    QVERIFY2(inner, "KvitSidebarItem has no KvitBadge in it");
    QCOMPARE(inner->property("counted").toString(), QStringLiteral("entry"));
    QCOMPARE(inner->property("countedPlural").toString(),
             QStringLiteral("entries"));

    // With no noun at all it falls back to the generic word, still inflected
    // and still grouped.
    QVERIFY(badge->setProperty("counted", QString()));
    QVERIFY(badge->setProperty("countedPlural", QString()));
    QCoreApplication::processEvents();
    QCOMPARE(holder->property("badgeName").toString(),
             QStringLiteral("250,000 items"));
    QVERIFY(badge->setProperty("count", 1));
    QCoreApplication::processEvents();
    QCOMPARE(holder->property("badgeName").toString(),
             QStringLiteral("1 item"));

    // And what the badge draws is the same number written the same way. The
    // pill and the line under it state one figure, and a reader who sees
    // "250000" above "250,000 records included" has to work out whether they
    // are the same thing.
    QVERIFY(badge->setProperty("count", 250000));
    QCoreApplication::processEvents();
    QQuickItem *label = findByType(badge, "KvitLabel");
    QVERIFY2(label, "KvitBadge draws no label");
    QCOMPARE(label->property("text").toString(), QStringLiteral("250,000"));

    // The cap is unchanged: past `max` the pill is too narrow for the number
    // and says so, while the announcement above still carries the real figure.
    QVERIFY(badge->setProperty("max", 99));
    QCoreApplication::processEvents();
    QCOMPARE(label->property("text").toString(), QStringLiteral("99+"));
    QCOMPARE(holder->property("badgeName").toString(),
             QStringLiteral("250,000 items"));
}

namespace {

// Every descendant whose type name starts with `prefix`, in tree order.
// findByType above answers with the first one, and a section heading draws
// four labels of which the one under test is the third.
void collectByType(QQuickItem *item, const char *prefix,
                   QList<QQuickItem *> &found)
{
    const auto children = item->childItems();
    for (QQuickItem *child : children) {
        if (QByteArray(child->metaObject()->className()).startsWith(prefix))
            found.append(child);
        collectByType(child, prefix, found);
    }
}

}   // namespace

void TestComponents::testAGroupHeadingCountsItsGroupTheWayEverythingElseDoes()
{
    // What a group heading says at the end of its bar.
    //
    // The heading was the last component still counting its group the way the
    // library counted before any of this was settled, and it had all three of
    // the faults the rest were repaired for: the plural was an "s" suffixed to
    // the caller's noun with no singular form to fall back to, the number was
    // written out bare so a four-figure group read "1200 accounts" beside a
    // ledger reading "1,200", and the whole phrase was hidden below two, so a
    // heading that says "2 accounts" for two said nothing at all for one.
    //
    // kvit-cash draws its sidebar's account group with this component, so all
    // three reached a screen.
    //
    // A fixed locale, for the reason the count-phrase case above fixes one:
    // the grouping is half of what is asserted here, and the C locale groups
    // nothing.
    const QLocale previous;
    QLocale::setDefault(QLocale(QLocale::English, QLocale::UnitedStates));
    const auto restore = qScopeGuard([previous] { QLocale::setDefault(previous); });

    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "KvitSectionHeading {\n"
        "    width: 400; height: 30\n"
        "    text: \"Accounts\"; counted: \"account\"; count: 1200\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/section-count.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    auto *heading = qobject_cast<QQuickItem *>(holder.data());
    QVERIFY(heading);
    QCoreApplication::processEvents();

    QCOMPARE(heading->property("countPhrase").toString(),
             QStringLiteral("1,200 accounts"));

    // Drawn, and not only computed. The label is held on to, so what follows
    // asks the same item what it says and whether it is on screen.
    QList<QQuickItem *> labels;
    collectByType(heading, "KvitLabel", labels);
    QQuickItem *drawn = nullptr;
    for (QQuickItem *label : labels) {
        if (label->property("text").toString()
                == QStringLiteral("1,200 accounts"))
            drawn = label;
    }
    QVERIFY2(drawn, "the heading computed a count phrase it does not draw");
    QVERIFY(drawn->isVisible());

    // One account is one account. This is the case the heading used to draw
    // nothing for.
    QVERIFY(heading->setProperty("count", 1));
    QCoreApplication::processEvents();
    QCOMPARE(heading->property("countPhrase").toString(),
             QStringLiteral("1 account"));
    QCOMPARE(drawn->property("text").toString(), QStringLiteral("1 account"));
    QVERIFY2(drawn->isVisible(),
             "a group of one drew no count where a group of two draws one");

    // A noun English does not pluralise by suffixing an s says its own plural.
    QVERIFY(heading->setProperty("counted", QStringLiteral("entry")));
    QVERIFY(heading->setProperty("countedPlural", QStringLiteral("entries")));
    QVERIFY(heading->setProperty("count", 1200));
    QCoreApplication::processEvents();
    QCOMPARE(heading->property("countPhrase").toString(),
             QStringLiteral("1,200 entries"));

    // With no noun at all the number stands on its own, and is still grouped.
    QVERIFY(heading->setProperty("counted", QString()));
    QVERIFY(heading->setProperty("countedPlural", QString()));
    QCoreApplication::processEvents();
    QCOMPARE(heading->property("countPhrase").toString(),
             QStringLiteral("1,200"));
    QCOMPARE(drawn->property("text").toString(), QStringLiteral("1,200"));

    // Below zero is a caller saying it has no count to give, and it is the
    // only case where the heading draws none.
    QVERIFY(heading->setProperty("count", -1));
    QCoreApplication::processEvents();
    QCOMPARE(heading->property("countPhrase").toString(), QString());
    QVERIFY2(!drawn->isVisible(),
             "a heading with no count to give drew one anyway");
}

void TestComponents::testAProgressBarCannotTrailTheNumberBesideIt()
{
    // What the bar draws while the work reporting to it holds the thread.
    //
    // The fill slides to each new value over 160 ms. Work that runs on the
    // interface thread and reports its own progress — a commit of a quarter of
    // a million records, which yields a few milliseconds between batches — is
    // work that leaves the animation almost no time to run, so each slide
    // covers a fraction of its step, the next one starts from where that one
    // stopped, and the shortfall compounds. Measured on the real commit: the
    // label read 42 per cent while the bar drew about 5.
    //
    // A starved event loop is exactly what this case is: values are set one
    // after another with nothing processed in between, so no animation
    // advances at all. The bar has to stay within one report of the truth
    // anyway.
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    property alias bar: bar\n"
        "    readonly property real motion: Theme.motionScale\n"
        "    KvitProgress {\n"
        "        id: bar\n"
        "        width: 400; height: 30\n"
        "        label: \"Committing\"; maximum: 250000\n"
        "    }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/progress-honest.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    auto *bar = holder->property("bar").value<QQuickItem *>();
    QVERIFY(bar);
    QCoreApplication::processEvents();

    // The track is the first rectangle under the bar and the fill the first
    // one under the track; neither has a name, for the reason nothing in this
    // library has one.
    QQuickItem *track = findByType(bar, "QQuickRectangle");
    QVERIFY2(track, "KvitProgress draws no track");
    QQuickItem *fill = findByType(track, "QQuickRectangle");
    QVERIFY2(fill, "KvitProgress draws no fill");

    const qreal full = track->width();
    QVERIFY2(full > 0, "the progress track has no width");
    const qreal oneReport = full / 20.0;

    // Twenty reports of five per cent each, with the event loop never given a
    // turn — a commit publishing its progress from the thread that would draw
    // it.
    for (int i = 1; i <= 20; ++i) {
        QVERIFY(bar->setProperty("value", 12500.0 * i));
        const qreal honest = qRound(full * i / 20.0);
        const qreal shown = fill->width();
        const QString said =
            QStringLiteral("at %1 per cent of %2 pixels the bar drew %3")
                .arg(i * 5).arg(full).arg(shown);
        QVERIFY2(shown <= honest + 0.5, qPrintable(said));
        QVERIFY2(shown >= honest - oneReport - 0.5, qPrintable(said));
    }

    // And a bar whose value moves while the application is idle still slides
    // rather than stepping, which is the half of this that must not have been
    // paid for the other half. Not asserted where the reader has asked for no
    // motion: there every duration in the library is zero and there is
    // nothing to slide.
    if (holder->property("motion").toReal() <= 0)
        return;

    QObject *slide = nullptr;
    const auto belongings = fill->children();
    for (QObject *child : belongings) {
        if (QByteArray(child->metaObject()->className())
                .startsWith("QQuickNumberAnimation"))
            slide = child;
    }
    QVERIFY2(slide, "KvitProgress no longer animates its fill at all");

    // Back to the start. The bar is behind after the run above, so this
    // report is one it jumps to, and a jump leaves it exactly where it was
    // asked to be with nothing running — a settled bar, with no event loop
    // needed to settle it.
    QVERIFY(bar->setProperty("value", 0.0));
    QCOMPARE(fill->width(), 0.0);
    QCOMPARE(slide->property("running").toBool(), false);

    // From there, the next report slides: the animation is running and the
    // width has not yet moved.
    QVERIFY(bar->setProperty("value", 125000.0));
    QVERIFY2(slide->property("running").toBool(),
             "a bar that was up to date stepped to its new value");
    QCOMPARE(fill->width(), 0.0);
}

void TestComponents::testARowSaysWhetherItIsPressableBeforeItIsPressed()
{
    // A row that opens something says so before it is pressed, and a row that
    // is only a layout says nothing.
    //
    // The same component builds both: a list of accounts whose rows open a
    // ledger, and a pane whose rows are a field name beside its value. When
    // every one of them tints under the pointer and answers a tap, the tint
    // stops carrying any information and fifteen field rows in a record pane
    // read as fifteen things to click.
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    id: holder\n"
        "    width: 240; height: 80\n"
        "    Column {\n"
        "        anchors.fill: parent\n"
        "        KvitRow {\n"
        "            id: plain\n"
        "            objectName: \"plain\"\n"
        "            width: 240; height: 40; label: \"Amount\"\n"
        "        }\n"
        "        KvitRow {\n"
        "            id: acting\n"
        "            objectName: \"acting\"\n"
        "            width: 240; height: 40; label: \"Groceries\"\n"
        "            activeFocusOnTab: true\n"
        "        }\n"
        "    }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/row-pressable.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));
    auto *content = qobject_cast<QQuickItem *>(holder.data());
    QVERIFY(content);

    auto *plain = content->findChild<QQuickItem *>(QStringLiteral("plain"));
    auto *acting = content->findChild<QQuickItem *>(QStringLiteral("acting"));
    QVERIFY(plain);
    QVERIFY(acting);

    // A row reachable by the keyboard is a row that does something, so it
    // needs nothing else said about it. A row that is not is static until it
    // says otherwise.
    QCOMPARE(plain->property("interactive").toBool(), false);
    QCOMPARE(acting->property("interactive").toBool(), true);

    // The tap follows it. A row inside a list that moves its own cursor with
    // the arrow keys is not in the tab order and sets `interactive` itself,
    // which is what this last step stands for.
    auto tapEnabled = [](QQuickItem *row) {
        const auto handlers = row->findChildren<QObject *>();
        for (QObject *child : handlers) {
            if (QByteArray(child->metaObject()->className())
                    .startsWith("QQuickTapHandler"))
                return child->property("enabled").toBool();
        }
        return false;
    };
    QVERIFY2(!tapEnabled(plain), "a layout row answered a tap");
    QVERIFY2(tapEnabled(acting), "a row that acts refused a tap");
    QVERIFY(plain->setProperty("interactive", true));
    QCoreApplication::processEvents();
    QVERIFY2(tapEnabled(plain), "a row that declared itself refused a tap");
    QVERIFY(plain->setProperty("interactive", false));

    // And the hover tint follows it, which is the half a reader sees. This
    // needs a window and a pointer, so it runs wherever the suite has one —
    // the gate runs offscreen, which does.
    if (QGuiApplication::platformName() == QLatin1String("minimal"))
        return;

    QQuickWindow window;
    content->setParentItem(window.contentItem());
    window.resize(240, 80);
    window.show();
    QVERIFY(QTest::qWaitForWindowExposed(&window));

    const QPoint overPlain =
        window.contentItem()->mapFromItem(plain, QPointF(120, 20)).toPoint();
    const QPoint overActing =
        window.contentItem()->mapFromItem(acting, QPointF(120, 20)).toPoint();

    QTest::mouseMove(&window, overPlain);
    QTRY_COMPARE(plain->property("hovered").toBool(), true);
    QCOMPARE(plain->property("color").value<QColor>().alpha(), 0);

    QSignalSpy tapped(acting, SIGNAL(activated()));
    QTest::mouseMove(&window, overActing);
    QTRY_COMPARE(acting->property("hovered").toBool(), true);
    QVERIFY2(acting->property("color").value<QColor>().alpha() > 0,
             "a row that opens something drew no hover tint");

    QTest::mouseClick(&window, Qt::LeftButton, Qt::NoModifier, overActing);
    QTRY_COMPARE(tapped.count(), 1);

    QSignalSpy ignored(plain, SIGNAL(activated()));
    QTest::mouseClick(&window, Qt::LeftButton, Qt::NoModifier, overPlain);
    QCoreApplication::processEvents();
    QCOMPARE(ignored.count(), 0);
}

void TestComponents::testARowAnnouncesItselfAsWhatItActuallyIs()
{
    // What a screen reader is told about a row, which has to say the same
    // thing the drawing says.
    //
    // A row that does nothing takes no hover tint and answers no tap, so a
    // reader who can see it reads it as a layout. Announcing it as a list
    // item — something to move to, select and open — tells a reader who
    // cannot see it the opposite, and kvit-cash's record pane is twenty-three
    // such rows under one that does open something. The tint and the tap were
    // put behind `interactive`; this is the half that did not follow.
    //
    // The roles are read through properties of the holder, because an
    // attached property is not a property of the item a test can ask for by
    // name.
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    width: 240; height: 160\n"
        "    property alias field: field\n"
        "    readonly property int fieldRole: field.Accessible.role\n"
        "    readonly property bool fieldSelectable: field.Accessible.selectable\n"
        "    readonly property int holdingRole: holding.Accessible.role\n"
        "    readonly property int actingRole: acting.Accessible.role\n"
        "    readonly property bool actingSelectable: acting.Accessible.selectable\n"
        "    readonly property bool pickedSelectable: picked.Accessible.selectable\n"
        "    readonly property bool pickedSelected: picked.Accessible.selected\n"
        "    Column {\n"
        "        anchors.fill: parent\n"
        "        KvitRow { id: field; width: 240; label: \"Amount, 12.50\" }\n"
        "        KvitRow { id: holding; width: 240\n"
        "            KvitLabel { anchors.centerIn: parent; text: \"12.50\" } }\n"
        "        KvitRow { id: acting; width: 240; label: \"Groceries\"\n"
        "            activeFocusOnTab: true }\n"
        "        KvitRow { id: picked; width: 240; label: \"Groceries\"\n"
        "            selected: true }\n"
        "    }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/row-announcement.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));
    QCoreApplication::processEvents();

    // A field name beside its value is the line of text it looks like, and
    // there is nothing to select.
    QCOMPARE(holder->property("fieldRole").toInt(), int(QAccessible::StaticText));
    QCOMPARE(holder->property("fieldSelectable").toBool(), false);

    // A row with no label of its own is the container it is: the labels
    // inside it are what a reader hears, rather than an unnamed line of text
    // between the rows that do say something.
    QCOMPARE(holder->property("holdingRole").toInt(), int(QAccessible::Grouping));

    // A row that opens something is a thing in a list, and can be picked.
    QCOMPARE(holder->property("actingRole").toInt(), int(QAccessible::ListItem));
    QCOMPARE(holder->property("actingSelectable").toBool(), true);

    // Selection is announced wherever it is drawn. A row a caller has drawn
    // as picked says so whether or not pressing it does anything.
    QCOMPARE(holder->property("pickedSelectable").toBool(), true);
    QCOMPARE(holder->property("pickedSelected").toBool(), true);

    // And the announcement follows the property, so a row inside a list that
    // moves its own cursor with the arrow keys — which sets `interactive`
    // rather than joining the tab order — is announced as a list item too.
    auto *field = holder->property("field").value<QQuickItem *>();
    QVERIFY(field);
    QVERIFY(field->setProperty("interactive", true));
    QCoreApplication::processEvents();
    QCOMPARE(holder->property("fieldRole").toInt(), int(QAccessible::ListItem));
    QCOMPARE(holder->property("fieldSelectable").toBool(), true);
}

namespace {

// Whether a tooltip has been built under an item. A ToolTip is a Popup rather
// than a child item — its content is drawn in the window's overlay — so it is
// found among the object's children rather than among its items.
bool hasTooltip(QObject *object)
{
    const auto children = object->findChildren<QObject *>();
    for (QObject *child : children) {
        if (QByteArray(child->metaObject()->className()).startsWith("KvitTooltip"))
            return true;
    }
    return false;
}

}   // namespace

void TestComponents::testAShortenedValueIsDisclosedToTheKeyboardAsWellAsThePointer()
{
    // A value the column is too narrow for stays reachable without a pointer.
    //
    // A cell is not a control and takes no focus of its own, so it cannot
    // answer that on its own: the view that owns the cell cursor says where
    // the cursor is through `current`, and the disclosure opens for hover or
    // for the cursor, which is the pair every real control in the estate uses
    // for the same purpose. Gated on hover alone, a reader who never touches
    // the pointer sees the truncated name and has no way to reach the rest.
    //
    // The tooltip is built when it is wanted rather than with the cell,
    // because a table draws several hundred cells at a time and a tooltip is
    // a popup with a window behind it. So its absence is what this reads
    // first, and it is also what says the loader is doing its job.
    QQmlEngine engine;
    QQmlComponent component(&engine);
    // The accessible name is read through a property of the holder, because
    // an attached property is not a property of the item a test can ask for
    // by name.
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    property alias cell: cell\n"
        "    readonly property string spokenName: cell.Accessible.name\n"
        "    KvitCell {\n"
        "        id: cell\n"
        "        width: 80; height: 30\n"
        "        kind: \"Text\"\n"
        "        value: \"Harlow & Co\"\n"
        "        fullValue: \"Harlow & Co (Holdings) Limited\"\n"
        "    }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/disclosure.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));
    auto *cell = holder->property("cell").value<QQuickItem *>();
    QVERIFY(cell);
    QCoreApplication::processEvents();

    // Nothing built while the cell is neither hovered nor under the cursor.
    QVERIFY2(!hasTooltip(cell),
             "a cell built its tooltip before anything asked for it");

    // The keyboard arrives, and the rest of the value is there.
    QVERIFY(cell->setProperty("current", true));
    QCoreApplication::processEvents();
    QTRY_VERIFY2(hasTooltip(cell),
                 "a cell under the keyboard cursor disclosed nothing");

    // And it goes away again with the cursor, so a cell the reader has moved
    // off does not leave a label floating over the row below it.
    QVERIFY(cell->setProperty("current", false));
    QCoreApplication::processEvents();
    QTRY_VERIFY(!hasTooltip(cell));

    // What a screen reader is told is the whole value rather than the
    // shortened one: a reader who is read a truncated payee is read a
    // different payee.
    QCOMPARE(holder->property("spokenName").toString(),
             QStringLiteral("Harlow & Co (Holdings) Limited"));
}

namespace {

// The accessibility press action, which is what a screen reader offers on
// whatever it is reading. It is a different route in from a key press: the
// caller has named the thing it means rather than arrived at it by tabbing.
bool pressThroughAccessibility(QQuickItem *item)
{
    QAccessibleInterface *accessible = QAccessible::queryAccessibleInterface(item);
    if (!accessible)
        return false;
    QAccessibleActionInterface *actions = accessible->actionInterface();
    if (!actions)
        return false;
    if (!actions->actionNames().contains(QAccessibleActionInterface::pressAction()))
        return false;
    actions->doAction(QAccessibleActionInterface::pressAction());
    return true;
}

// Every item under `root` with this object name, in the order they are drawn.
//
// The visual tree rather than the object tree: a Repeater's delegates are
// parented into the item they appear in but are owned by the delegate model,
// so findChildren() does not reach them.
void collectNamed(QQuickItem *item, const QString &name,
                  QList<QQuickItem *> &found)
{
    const QList<QQuickItem *> children = item->childItems();
    for (QQuickItem *child : children) {
        if (child->objectName() == name)
            found.append(child);
        collectNamed(child, name, found);
    }
}

QList<QQuickItem *> itemsNamed(QQuickItem *root, const QString &name)
{
    QList<QQuickItem *> found;
    collectNamed(root, name, found);
    return found;
}

}   // namespace

void TestComponents::testARowOpensFromTheKeyboardAndOnlyOnce()
{
    // A row that takes the keyboard and does nothing with it.
    //
    // KvitRow drew a focus ring, announced itself as a list item and answered
    // a tap, and had no key handler at all: Tab reached it, the ring said it
    // was there, and Return did nothing. kvit-notes-pro hand-wrote its own
    // key handling on a private copy of this control for exactly that reason,
    // and that copy is what its migration onto this library deletes.
    //
    // The part worth a test is not the key handler. It is that a row is
    // usually a container: a Space pressed while a button inside it has the
    // keyboard has to press the button and leave the record shut, or one
    // keystroke does two things and only one of them is visible.
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
            width: 320
            height: 240
            property int opened: 0
            property int split: 0
            Column {
                anchors.fill: parent
                KvitRow {
                    objectName: "acting"
                    width: 320; height: 40; label: "Groceries"
                    activeFocusOnTab: true
                    onActivated: window.opened += 1
                }
                KvitRow {
                    objectName: "holding"
                    width: 320; height: 40; label: "Rent"
                    activeFocusOnTab: true
                    onActivated: window.opened += 1
                    KvitButton {
                        objectName: "inside"
                        anchors.verticalCenter: parent.verticalCenter
                        text: "Split"
                        onClicked: window.split += 1
                    }
                }
                KvitRow {
                    objectName: "layout"
                    width: 320; height: 40; label: "Amount"
                    onActivated: window.opened += 1
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/row-keyboard.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto *acting = window->findChild<QQuickItem *>(QStringLiteral("acting"));
    auto *holding = window->findChild<QQuickItem *>(QStringLiteral("holding"));
    auto *layout = window->findChild<QQuickItem *>(QStringLiteral("layout"));
    auto *inside = window->findChild<QQuickItem *>(QStringLiteral("inside"));
    QVERIFY(acting);
    QVERIFY(holding);
    QVERIFY(layout);
    QVERIFY(inside);

    // The ring follows the keyboard as well as the caller's own cursor. A row
    // that answers Return without saying it has the keyboard is a row nobody
    // can tell is about to open.
    auto *ring = acting->findChild<QQuickItem *>(QStringLiteral("focusRing"));
    QVERIFY(ring);
    QVERIFY(!ring->isVisible());

    acting->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(acting->hasActiveFocus());
    QTRY_VERIFY(ring->isVisible());

    QSignalSpy opened(acting, SIGNAL(activated()));
    QTest::keyClick(window, Qt::Key_Return);
    QTRY_COMPARE(opened.count(), 1);
    QTest::keyClick(window, Qt::Key_Enter);
    QTRY_COMPARE(opened.count(), 2);
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_COMPARE(opened.count(), 3);

    QVERIFY2(pressThroughAccessibility(acting),
             "a row announced as a list item offered no press action");
    QTRY_COMPARE(opened.count(), 4);
    QCOMPARE(window->property("opened").toInt(), 4);

    // The pointer is the fourth route and produces one activation like the
    // rest. It needs a real pointer, so it runs wherever the suite has one —
    // the gate runs offscreen, which does.
    if (QGuiApplication::platformName() != QLatin1String("minimal")) {
        const QPoint overActing =
            window->contentItem()->mapFromItem(acting, QPointF(160, 20)).toPoint();
        QTest::mouseClick(window, Qt::LeftButton, Qt::NoModifier, overActing);
        QTRY_COMPARE(opened.count(), 5);
    }

    // The container case. The button takes the Space it was sent and the row
    // behind it stays shut.
    QSignalSpy holdingOpened(holding, SIGNAL(activated()));
    inside->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(inside->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_COMPARE(window->property("split").toInt(), 1);
    QCoreApplication::processEvents();
    QCOMPARE(holdingOpened.count(), 0);

    // Return is the key nothing in Qt Quick Controls consumes for a button,
    // so it is the one that would reach the row by propagation. It does not,
    // because the row answers only while the row itself has the keyboard.
    QTest::keyClick(window, Qt::Key_Return);
    QCoreApplication::processEvents();
    QCOMPARE(holdingOpened.count(), 0);
    QCOMPARE(window->property("split").toInt(), 1);

    // And a row that is only a layout ignores all four routes, for the same
    // reason it draws no hover tint: a record pane's field rows are not
    // things to open.
    QSignalSpy ignored(layout, SIGNAL(activated()));
    layout->forceActiveFocus(Qt::OtherFocusReason);
    QTRY_VERIFY(layout->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Return);
    QTest::keyClick(window, Qt::Key_Space);
    QCoreApplication::processEvents();
    QCOMPARE(ignored.count(), 0);
    // The press action included. Qt offers one wherever a component declares
    // a handler and gives no way to withdraw it, so what a layout row
    // promises is unavoidable and what it does when asked is not.
    pressThroughAccessibility(layout);
    QCoreApplication::processEvents();
    QCOMPARE(ignored.count(), 0);
}

void TestComponents::testASectionHeadingOpensFromTheKeyboardWithoutSwallowingItsAction()
{
    // The same gap as the row's, in the component every section of a stacked
    // overview is built from: the chevron said the group could be opened and
    // only the pointer could open it.
    //
    // The half that was already right is the action. It is a separate control
    // rather than a second meaning for pressing the bar, so the test that
    // matters is that running it does not also collapse the group behind it.
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
            width: 420
            height: 200
            Column {
                anchors.fill: parent
                KvitSectionHeading {
                    objectName: "group"
                    width: 420
                    text: "Waiting on me"
                    counted: "project"; count: 4
                    action: "Hand all to an agent"
                    collapsible: true
                }
                KvitSectionHeading {
                    objectName: "fixed"
                    width: 420
                    text: "Everything else"
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/heading-keyboard.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto *heading = window->findChild<QQuickItem *>(QStringLiteral("group"));
    auto *fixed = window->findChild<QQuickItem *>(QStringLiteral("fixed"));
    QVERIFY(heading);
    QVERIFY(fixed);

    // A heading that can be opened is a stop worth having; one that cannot is
    // a line of text, and a stop there does nothing.
    QCOMPARE(heading->property("activeFocusOnTab").toBool(), true);
    QCOMPARE(fixed->property("activeFocusOnTab").toBool(), false);

    auto *ring = heading->findChild<QQuickItem *>(QStringLiteral("focusRing"));
    QVERIFY(ring);
    QVERIFY(!ring->isVisible());

    QSignalSpy toggled(heading, SIGNAL(toggled()));
    QSignalSpy actioned(heading, SIGNAL(actioned()));

    heading->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(heading->hasActiveFocus());
    QTRY_VERIFY(ring->isVisible());

    QTest::keyClick(window, Qt::Key_Return);
    QTRY_COMPARE(toggled.count(), 1);
    QTest::keyClick(window, Qt::Key_Enter);
    QTRY_COMPARE(toggled.count(), 2);
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_COMPARE(toggled.count(), 3);
    QVERIFY2(pressThroughAccessibility(heading),
             "a heading announced as a heading offered no press action");
    QTRY_COMPARE(toggled.count(), 4);

    // The action is its own stop and its own key. Space on it runs the action
    // and leaves the group at four.
    auto *action = heading->findChild<QQuickItem *>(QStringLiteral("action"));
    QVERIFY(action);
    action->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(action->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_COMPARE(actioned.count(), 1);
    QCoreApplication::processEvents();
    QCOMPARE(toggled.count(), 4);

    QTest::keyClick(window, Qt::Key_Return);
    QTRY_COMPARE(actioned.count(), 2);
    QCoreApplication::processEvents();
    QCOMPARE(toggled.count(), 4);

    // A heading that does not collapse has nothing to open, whichever way it
    // is asked.
    QSignalSpy fixedToggled(fixed, SIGNAL(toggled()));
    fixed->forceActiveFocus(Qt::OtherFocusReason);
    QTRY_VERIFY(fixed->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Space);
    QCoreApplication::processEvents();
    QCOMPARE(fixedToggled.count(), 0);
    pressThroughAccessibility(fixed);
    QCoreApplication::processEvents();
    QCOMPARE(fixedToggled.count(), 0);
}

void TestComponents::testAnIconButtonDrawsTheSymbolSizeItWasAskedFor()
{
    // A strip of icon buttons across a header, or one hoisted onto a heading
    // bar, is drawn beside symbols the rest of the library draws at 13 —
    // KvitLink's, KvitSelect's, KvitSectionHeading's own chevron. At 18 they
    // read as larger than everything around them, and there was no way to ask
    // for the other size: the content item is anchored rather than laid out,
    // so neither the control's padding nor its width and height reach the
    // symbol, and a caller setting the button to 20 pixels got an 18-pixel
    // symbol filling it edge to edge.
    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(
        InterfaceMetrics::DefaultFontSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(R"(
        import QtQuick
        import Kvit.Ui
        Item {
            property alias ordinary: ordinary
            property alias dense: dense
            KvitIconButton { id: ordinary; symbol: "search"; label: "Search" }
            KvitIconButton {
                id: dense
                symbol: "search"; label: "Search"; dense: true
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/icon-button-size.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    auto *metrics = KvitUi::DefaultServices::interfaceMetrics();
    for (const char *name : { "ordinary", "dense" }) {
        auto *button = holder->property(name).value<QQuickItem *>();
        QVERIFY(button);
        auto *symbol = button->findChild<QQuickItem *>(
            QStringLiteral("symbol"));
        QVERIFY2(symbol, "the icon button has no symbol");
        const qreal expected = button->property("dense").toBool()
            ? metrics->iconSizeSmall() : metrics->iconSize();
        QCOMPARE(symbol->width(), expected);
        QCOMPARE(symbol->height(), expected);
    }

    // The box is the caller's, and it is still the box whichever size the
    // symbol is: a heading bar sets it to the row it sits in.
    auto *dense = holder->property("dense").value<QQuickItem *>();
    QCOMPARE(dense->implicitHeight(), qreal(metrics->controlHeight()));
}

void TestComponents::testAControlCarriesOneSentenceOfExplanation()
{
    // A button's own words are its name and are usually all it needs. The
    // case this is for is the one where the reason a control is in the state
    // it is in lives somewhere the reader cannot see: a Pull button disabled
    // because the remote has not been fetched, an Archive button disabled
    // because the branch still has unpushed work. Without somewhere to put
    // that sentence, the screen offers a grey control and no account of it,
    // and a screen reader is told even less.
    //
    // The same string goes to both surfaces on all three controls, which is
    // what stops the words a pointer reader sees drifting from the words a
    // screen reader hears — the defect accessibility.md Finding 1 records.
    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(
        InterfaceMetrics::DefaultFontSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(R"(
        import QtQuick
        import QtQuick.Controls
        import Kvit.Ui
        ApplicationWindow {
            visible: true
            width: 520
            height: 260
            KvitButton {
                objectName: "pull"
                x: 20; y: 20
                text: "Pull"
                explanation: "Brings the 2 commits on the remote into this branch."
            }
            KvitButton {
                objectName: "archive"
                x: 20; y: 80
                text: "Archive"
                enabled: false
                explanation: "The branch has work that has not been pushed."
            }
            KvitIconButton {
                objectName: "rename"
                x: 240; y: 20
                symbol: "rename"
                label: "Rename"
                explanation: "Renames the track and its branch. Its history and folder stay unchanged."
            }
            KvitChipButton {
                objectName: "ahead"
                x: 240; y: 80
                text: "1 ahead"
                explanation: "Opens the one commit this branch has and the remote does not."
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/explanation.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    struct Case
    {
        const char *objectName;
        const char *name;
        const char *explanation;
    };
    const Case cases[] = {
        { "pull", "Pull",
          "Brings the 2 commits on the remote into this branch." },
        { "archive", "Archive",
          "The branch has work that has not been pushed." },
        { "rename", "Rename",
          "Renames the track and its branch. Its history and folder stay "
          "unchanged." },
        { "ahead", "1 ahead",
          "Opens the one commit this branch has and the remote does not." },
    };

    for (const Case &one : cases) {
        auto *control = window->findChild<QQuickItem *>(
            QLatin1String(one.objectName));
        QVERIFY2(control, one.objectName);

        // The name stays the control's own words. An explanation is a second
        // string beside it, never a replacement for it: a control whose
        // purpose is only in its explanation cannot be used without hovering
        // it, which rules out everybody on a touch screen or a keyboard.
        QAccessibleInterface *accessible =
            QAccessible::queryAccessibleInterface(control);
        QVERIFY2(accessible, one.objectName);
        QCOMPARE(accessible->text(QAccessible::Name),
                 QString::fromLatin1(one.name));
        QCOMPARE(accessible->text(QAccessible::Description),
                 QString::fromLatin1(one.explanation));

        QObject *tooltip = control->findChild<QObject *>(
            QStringLiteral("tooltip"));
        QVERIFY2(tooltip, one.objectName);
        QVERIFY2(!tooltip->property("visible").toBool(),
                 "an explanation was showing before anybody asked for it");
    }

    // What a pointer reader sees, on the control where it matters most: the
    // disabled one. Qt stops sending hover events to a disabled item and
    // takes it out of the tab order, so neither `hovered` nor the focus ring
    // can be what opens this.
    if (QGuiApplication::platformName() != QLatin1String("minimal")) {
        auto *archive = window->findChild<QQuickItem *>(
            QStringLiteral("archive"));
        QVERIFY(archive);
        QObject *tooltip = archive->findChild<QObject *>(
            QStringLiteral("tooltip"));
        QVERIFY(tooltip);
        const QPoint over = window->contentItem()
            ->mapFromItem(archive,
                          QPointF(archive->width() / 2, archive->height() / 2))
            .toPoint();
        QTest::mouseMove(window, over);
        QTRY_VERIFY2(tooltip->property("visible").toBool(),
                     "a disabled button never said why it was disabled");
        QCOMPARE(tooltip->property("text").toString(),
                 QStringLiteral("The branch has work that has not been "
                                "pushed."));
        QTest::mouseMove(window, QPoint(2, 2));
        QTRY_VERIFY(!tooltip->property("visible").toBool());
    }

    // And what a keyboard reader sees on the ones it can reach. The icon
    // button shows its label and its explanation together, because its label
    // is a picture everywhere else.
    auto *rename = window->findChild<QQuickItem *>(QStringLiteral("rename"));
    QVERIFY(rename);
    QObject *renameTooltip = rename->findChild<QObject *>(
        QStringLiteral("tooltip"));
    QVERIFY(renameTooltip);
    rename->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(rename->hasActiveFocus());
    QTRY_VERIFY(renameTooltip->property("visible").toBool());
    QCOMPARE(renameTooltip->property("text").toString(),
             QStringLiteral("Rename\nRenames the track and its branch. Its "
                            "history and folder stay unchanged."));

    auto *ahead = window->findChild<QQuickItem *>(QStringLiteral("ahead"));
    QVERIFY(ahead);
    QObject *aheadTooltip = ahead->findChild<QObject *>(
        QStringLiteral("tooltip"));
    QVERIFY(aheadTooltip);
    ahead->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(ahead->hasActiveFocus());
    QTRY_VERIFY(aheadTooltip->property("visible").toBool());
    QCOMPARE(aheadTooltip->property("text").toString(),
             QStringLiteral("Opens the one commit this branch has and the "
                            "remote does not."));

    // A chip that cannot be pressed says that instead. The two are different
    // sentences about different states and neither has to be swapped into the
    // other's property for the right one to be shown.
    ahead->setProperty("unavailableReason",
                       QStringLiteral("The remote has not been fetched yet."));
    QCoreApplication::processEvents();
    QTRY_COMPARE(aheadTooltip->property("text").toString(),
                 QStringLiteral("The remote has not been fetched yet."));
    QAccessibleInterface *aheadAccessible =
        QAccessible::queryAccessibleInterface(ahead);
    QVERIFY(aheadAccessible);
    QCOMPARE(aheadAccessible->text(QAccessible::Description),
             QStringLiteral("The remote has not been fetched yet."));
}

void TestComponents::testAChipThatActsIsAControlAndSaysWhyItCannot()
{
    // The chip a person can press, and the chip that says why they cannot.
    //
    // A consuming application drew the first as a Button with a bordered
    // background of its own and got the alignment wrong; the second is the
    // state it needs and the one `enabled: false` cannot express, because Qt
    // takes a disabled item out of the tab order and stops sending it hover
    // events — which makes the one chip on the row with something to explain
    // the one chip a reader cannot reach to hear it.
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
            width: 420
            height: 160
            property int opened: 0
            property int refused: 0
            Row {
                anchors.centerIn: parent
                spacing: 12
                KvitChipButton {
                    objectName: "open"
                    text: "1 ahead"
                    trailingSymbol: "chevron-right"
                    onActivated: window.opened += 1
                }
                KvitChipButton {
                    objectName: "blocked"
                    text: "2 behind"
                    unavailableReason: "The other branch has not been fetched yet."
                    onActivated: window.refused += 1
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/chip-button.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto *open = window->findChild<QQuickItem *>(QStringLiteral("open"));
    auto *blocked = window->findChild<QQuickItem *>(QStringLiteral("blocked"));
    QVERIFY(open);
    QVERIFY(blocked);

    // It is a button, not a rectangle with a handler on it, which is what
    // reaches assistive technology at all.
    QAccessibleInterface *openAccessible =
        QAccessible::queryAccessibleInterface(open);
    QVERIFY2(openAccessible, "KvitChipButton has no accessible interface");
    QCOMPARE(openAccessible->role(), QAccessible::Button);
    QCOMPARE(openAccessible->text(QAccessible::Name), QStringLiteral("1 ahead"));

    // Both keys, then the pointer.
    open->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(open->hasActiveFocus());
    QTRY_VERIFY2(open->property("visualFocus").toBool(),
                 "a chip with the keyboard drew no focus state");
    QTest::keyClick(window, Qt::Key_Space);
    QTRY_COMPARE(window->property("opened").toInt(), 1);
    QTest::keyClick(window, Qt::Key_Return);
    QTRY_COMPARE(window->property("opened").toInt(), 2);
    QVERIFY(pressThroughAccessibility(open));
    QTRY_COMPARE(window->property("opened").toInt(), 3);

    if (QGuiApplication::platformName() != QLatin1String("minimal")) {
        const QPoint overOpen = window->contentItem()
            ->mapFromItem(open, QPointF(open->width() / 2, open->height() / 2))
            .toPoint();
        QTest::mouseClick(window, Qt::LeftButton, Qt::NoModifier, overOpen);
        QTRY_COMPARE(window->property("opened").toInt(), 4);
    }

    // An available chip whose label fits has nothing to add, so nothing is
    // shown over it.
    QObject *openTooltip = open->findChild<QObject *>(QStringLiteral("tooltip"));
    QVERIFY(openTooltip);
    QVERIFY(!openTooltip->property("visible").toBool());

    // The chip that cannot be pressed is still a place the keyboard stops,
    // still says why in both the surface a pointer reads and the one a screen
    // reader reads, and still emits nothing.
    QCOMPARE(blocked->property("unavailable").toBool(), true);
    QCOMPARE(blocked->property("enabled").toBool(), true);
    QCOMPARE(blocked->property("activeFocusOnTab").toBool(), true);

    QAccessibleInterface *blockedAccessible =
        QAccessible::queryAccessibleInterface(blocked);
    QVERIFY(blockedAccessible);
    QCOMPARE(blockedAccessible->text(QAccessible::Description),
             QStringLiteral("The other branch has not been fetched yet."));

    QObject *blockedTooltip =
        blocked->findChild<QObject *>(QStringLiteral("tooltip"));
    QVERIFY(blockedTooltip);
    blocked->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(blocked->hasActiveFocus());
    QTRY_VERIFY2(blockedTooltip->property("visible").toBool(),
                 "a chip that cannot be pressed did not say why");
    QCOMPARE(blockedTooltip->property("text").toString(),
             QStringLiteral("The other branch has not been fetched yet."));

    QTest::keyClick(window, Qt::Key_Space);
    QTest::keyClick(window, Qt::Key_Return);
    QVERIFY(pressThroughAccessibility(blocked));
    QCoreApplication::processEvents();
    QCOMPARE(window->property("refused").toInt(), 0);
}

void TestComponents::testAChipCanSayItIsTheOneAlreadyOpen()
{
    // A row of chips that open things — "3 changes", "1 ahead", "2 behind" —
    // and the reader is already looking at one of them. Without a way to draw
    // that chip as the current one, pressing it appears to do nothing,
    // because what it opens is already open.
    //
    // `tone: "accent"` is what a caller reaches for instead, and it says
    // something else: that this is a different kind of fact from the ones
    // beside it. It is the same fact, in the state of being the current one.
    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(
        InterfaceMetrics::DefaultFontSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(R"(
        import QtQuick
        import QtQuick.Controls
        import Kvit.Ui
        ApplicationWindow {
            visible: true
            width: 420
            height: 140
            readonly property color currentGround: ahead.background.color
            readonly property color currentEdge: ahead.background.border.color
            readonly property real currentEdgeWidth: ahead.background.border.width
            readonly property color plainGround: changes.background.color
            Row {
                anchors.centerIn: parent
                spacing: 12
                KvitChipButton {
                    id: changes
                    objectName: "changes"; text: "3 changes"
                }
                KvitChipButton {
                    id: ahead
                    objectName: "ahead"; text: "1 ahead"; current: true
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/chip-current.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto *plain = window->findChild<QQuickItem *>(QStringLiteral("changes"));
    auto *current = window->findChild<QQuickItem *>(QStringLiteral("ahead"));
    QVERIFY(plain);
    QVERIFY(current);

    // It is still a chip that acts: the same tab stop, the same keys, the
    // same button role. Being the current one is a state, not a different
    // control.
    QCOMPARE(current->property("activeFocusOnTab").toBool(), true);
    QAccessibleInterface *accessible =
        QAccessible::queryAccessibleInterface(current);
    QVERIFY(accessible);
    QCOMPARE(accessible->role(), QAccessible::Button);
    QVERIFY2(accessible->state().selected,
             "the current chip did not announce itself as the current one");

    // Three channels, because the tints are two percent of lightness in the
    // high-contrast theme and nothing at all in a grayscale screenshot: the
    // ground, the accent on the edge, and the weight of the label.
    QCOMPARE(window->property("currentGround").value<QColor>(),
             KvitUi::DefaultServices::theme()->selectionTint());
    QCOMPARE(window->property("currentEdge").value<QColor>(),
             KvitUi::DefaultServices::theme()->accent());
    QVERIFY(window->property("currentEdgeWidth").toReal() > 0);
    QVERIFY(window->property("currentGround").value<QColor>()
            != window->property("plainGround").value<QColor>());

    auto *label = current->findChild<QQuickItem *>(QStringLiteral("label"));
    auto *plainLabel = plain->findChild<QQuickItem *>(QStringLiteral("label"));
    QVERIFY(label);
    QVERIFY(plainLabel);
    QCOMPARE(label->property("font").value<QFont>().bold(), true);
    QCOMPARE(plainLabel->property("font").value<QFont>().bold(), false);
}

void TestComponents::testAHeadingTakesAWrittenCountAndASymbolicAction()
{
    // Two things a heading is asked for that it could not say.
    //
    // A count that is not a number: kvit-notes-pro's Changes heading reads
    // "1 · +0 −0" — one changed file, and the lines added and removed across
    // it — which arrives from the service already written and which no
    // integer expresses. And an action drawn as a symbol: a column of eight
    // headings each ending in a different sentence is a column a reader has
    // to read rather than scan, and the words are still there as the button's
    // name and its tooltip.
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
            width: 420
            height: 180
            property int opened: 0
            KvitSectionHeading {
                objectName: "changes"
                width: parent.width
                y: 20
                text: "Changes"
                countText: "1 · +0 −0"
                action: "Open the diff"
                actionSymbol: "diff"
                onActioned: window.opened += 1
            }
            KvitSectionHeading {
                objectName: "counted"
                width: parent.width
                y: 80
                text: "Waiting on me"
                count: 4
                counted: "project"
                action: "Hand all to an agent"
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/heading-written-count.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto *changes = window->findChild<QQuickItem *>(QStringLiteral("changes"));
    auto *counted = window->findChild<QQuickItem *>(QStringLiteral("counted"));
    QVERIFY(changes);
    QVERIFY(counted);

    // The phrase is drawn as written: no grouping, no plural, no locale. The
    // caller decided all three before handing it over.
    auto *written = changes->findChild<QQuickItem *>(
        QStringLiteral("countText"));
    QVERIFY(written);
    QVERIFY(written->isVisible());
    QCOMPARE(written->property("text").toString(),
             QString::fromUtf8("1 · +0 −0"));

    // And it leaves the right end empty, so a heading never shows two counts.
    QCOMPARE(changes->property("countPhrase").toString(), QString());
    QCOMPARE(counted->property("countPhrase").toString(),
             QStringLiteral("4 projects"));

    // Beside the name rather than at the right end, which is where the number
    // goes. The two are different things to a reader: a tally to glance at in
    // every heading down a column, against a phrase that is part of what the
    // group is.
    QVERIFY2(written->mapToItem(changes, QPointF()).x()
                 < changes->width() / 2,
             "the written count was drawn at the right end of the bar");

    // The action keeps its words where a reader who cannot see the symbol
    // finds them: the button's accessible name, and its tooltip.
    auto *button = changes->findChild<QQuickItem *>(
        QStringLiteral("actionButton"));
    QVERIFY(button);
    QVERIFY(button->isVisible());
    QAccessibleInterface *accessible =
        QAccessible::queryAccessibleInterface(button);
    QVERIFY(accessible);
    QCOMPARE(accessible->role(), QAccessible::Button);
    QCOMPARE(accessible->text(QAccessible::Name),
             QStringLiteral("Open the diff"));
    QCOMPARE(button->property("symbol").toString(), QStringLiteral("diff"));

    // It fits inside the bar, ring and all. The ring is drawn outside the
    // button's own ground, so a button as tall as the row has its ring cut
    // off by the rows above and below.
    auto *metrics = KvitUi::DefaultServices::interfaceMetrics();
    QVERIFY2(button->height() + 2 * metrics->focusRingWidth()
                 <= changes->height() + 0.5,
             qPrintable(QStringLiteral("a %1-tall action and its ring in a "
                                       "%2-tall bar")
                            .arg(button->height()).arg(changes->height())));

    // A heading given no symbol keeps the link it always had, and a heading
    // given one does not draw both.
    auto *link = changes->findChild<QQuickItem *>(QStringLiteral("action"));
    QVERIFY(link);
    QVERIFY2(!link->isVisible(), "a heading drew its action twice");
    auto *countedLink = counted->findChild<QQuickItem *>(
        QStringLiteral("action"));
    QVERIFY(countedLink);
    QVERIFY(countedLink->isVisible());

    // And the symbol runs the action.
    QVERIFY(pressThroughAccessibility(button));
    QTRY_COMPARE(window->property("opened").toInt(), 1);
}

void TestComponents::testAChipsLabelGivesWayRatherThanClipping_data()
{
    QTest::addColumn<int>("interfaceSize");
    QTest::addColumn<int>("chipWidth");

    for (int size : { InterfaceMetrics::DefaultFontSize,
                      // 200% of the default, which is where a phrase that fit
                      // at rest stops fitting.
                      2 * InterfaceMetrics::DefaultFontSize }) {
        for (int width : { 200, 120, 70 }) {
            QTest::newRow(qPrintable(QStringLiteral("%1px-%2wide")
                                         .arg(size).arg(width)))
                << size << width;
        }
    }
}

void TestComponents::testAChipsLabelGivesWayRatherThanClipping()
{
    // What has to survive a column too narrow for the phrase.
    //
    // The label is the piece that gives way; the symbols keep their size,
    // because a symbol at half width is a smudge and the trailing one is what
    // says where pressing the chip goes. A label wider than the box it is in
    // is the failure this catches: it draws over the border and over the
    // chevron rather than ending in an ellipsis.
    QFETCH(int, interfaceSize);
    QFETCH(int, chipWidth);

    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(interfaceSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\nimport Kvit.Ui\n"
        "KvitChipButton {\n"
        "    symbol: \"warning\"\n"
        "    trailingSymbol: \"chevron-right\"\n"
        "    text: \"A fact whose whole phrase does not fit in this column\"\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/chip-elision.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *chip = qobject_cast<QQuickItem *>(instance.data());
    QVERIFY(chip);
    chip->setWidth(chipWidth);
    QCoreApplication::processEvents();

    auto *row = chip->findChild<QQuickItem *>(QStringLiteral("contentRow"));
    auto *label = chip->findChild<QQuickItem *>(QStringLiteral("label"));
    QVERIFY(row);
    QVERIFY(label);

    // Everything the chip draws stays inside the chip.
    const qreal available =
        chip->width() - chip->property("leftPadding").toReal()
        - chip->property("rightPadding").toReal();
    QVERIFY2(row->width() <= available + 0.5,
             qPrintable(QStringLiteral("the content is %1 wide inside %2")
                            .arg(row->width()).arg(available)));
    QVERIFY2(label->width() <= row->width() + 0.5,
             qPrintable(QStringLiteral("the label is %1 wide inside %2")
                            .arg(label->width()).arg(row->width())));

    // And it is the label that gave way, by eliding rather than by being cut
    // off at the edge.
    if (label->width() + 0.5 < label->implicitWidth()) {
        QVERIFY2(label->property("truncated").toBool(),
                 "a label narrower than its text did not elide");
    }
}

void TestComponents::testAnEmptySectionSaysSoOnOneLine_data()
{
    QTest::addColumn<int>("interfaceSize");
    QTest::addColumn<int>("width");

    for (int size : { InterfaceMetrics::DefaultFontSize,
                      2 * InterfaceMetrics::DefaultFontSize }) {
        for (int width : { 360, 180 }) {
            QTest::newRow(qPrintable(QStringLiteral("%1px-%2wide")
                                         .arg(size).arg(width)))
                << size << width;
        }
    }
}

void TestComponents::testAnEmptySectionSaysSoOnOneLine()
{
    // The compact empty state, which exists because the full one is several
    // times too tall for a column of seven or eight sections that may each be
    // empty. A consuming application drew a bare Text there instead, which is
    // one more hand-drawn control and one more sentence written in a
    // different voice from the block in the middle of an empty pane.
    QFETCH(int, interfaceSize);
    QFETCH(int, width);

    KvitUi::DefaultServices::interfaceMetrics()->setFontSize(interfaceSize);
    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\nimport Kvit.Ui\n"
        "Item {\n"
        "    property alias compact: compact\n"
        "    property alias full: full\n"
        "    readonly property string compactName: compact.Accessible.name\n"
        "    KvitEmptyState {\n"
        "        id: compact\n"
        "        form: \"compact\"\n"
        "        symbol: \"robot\"\n"
        "        title: \"No agents\"\n"
        "        detail: \"none started here yet\"\n"
        "        action: \"Start one\"\n"
        "    }\n"
        "    KvitEmptyState {\n"
        "        id: full\n"
        "        symbol: \"robot\"\n"
        "        title: \"No agents\"\n"
        "        detail: \"none started here yet\"\n"
        "        action: \"Start one\"\n"
        "    }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/empty-compact.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));
    auto *compact = holder->property("compact").value<QQuickItem *>();
    auto *full = holder->property("full").value<QQuickItem *>();
    QVERIFY(compact);
    QVERIFY(full);
    compact->setWidth(width);
    full->setWidth(width);
    QCoreApplication::processEvents();

    // One line, at the height a slim row would have had. The full form is
    // still the block it was, which is the half that must not have moved for
    // the callers already using it.
    const qreal slim =
        KvitUi::DefaultServices::interfaceMetrics()->rowHeightSlim();
    QCOMPARE(compact->implicitHeight(), slim);
    QVERIFY2(full->implicitHeight() > 2 * compact->implicitHeight(),
             qPrintable(QStringLiteral("the full form is %1 tall against the "
                                       "compact form's %2")
                            .arg(full->implicitHeight())
                            .arg(compact->implicitHeight())));

    // Nothing on the line is drawn outside the line.
    auto *title = compact->findChild<QQuickItem *>(QStringLiteral("compactTitle"));
    auto *detail = compact->findChild<QQuickItem *>(QStringLiteral("compactDetail"));
    QVERIFY(title);
    QVERIFY(detail);
    for (QQuickItem *piece : { title, detail }) {
        const QPointF origin = piece->mapToItem(compact, QPointF(0, 0));
        QVERIFY2(origin.x() >= -0.5 && origin.x() + piece->width()
                     <= compact->width() + 0.5,
                 qPrintable(QStringLiteral("a piece runs from %1 to %2 inside "
                                           "a line %3 wide")
                                .arg(origin.x())
                                .arg(origin.x() + piece->width())
                                .arg(compact->width())));
        QVERIFY2(origin.y() >= -0.5 && origin.y() + piece->height()
                     <= compact->implicitHeight() + 0.5,
                 "a piece of the compact line is taller than the line");
    }

    // It starts where the rows it stands in for start. A line centred among
    // left-aligned rows reads as a notice about the section rather than as
    // the section's contents, which is what a caller who reached for a bare
    // Text instead was avoiding. The full form still centres: it is the whole
    // of an empty pane and there is nothing beside it to line up with.
    const qreal inset = KvitUi::DefaultServices::interfaceMetrics()->space();
    auto *lineRow = compact->findChild<QQuickItem *>(
        QStringLiteral("compactLine"));
    QVERIFY(lineRow);
    QCOMPARE(lineRow->mapToItem(compact, QPointF()).x(), inset);
    QVERIFY2(title->mapToItem(compact, QPointF()).x() < compact->width() / 2,
             "the sentence was drawn past the middle of the line");

    // What a screen reader hears is the sentence and its reason, the same in
    // both forms: the compact line is a shorter drawing of the same thing,
    // not a shorter thing.
    QCOMPARE(holder->property("compactName").toString(),
             QStringLiteral("No agents. none started here yet"));
}

namespace {

// What a bar too narrow for its groups does, whichever end they sit at.
//
// Three claims, and the order the groups are drawn in changes none of them,
// which is why every bar that has groups is put through this rather than only
// the one the component was written for. The window has already been narrowed
// by the time this is called.
//
//   • the bar says how many facts are not on it, in a control rather than in
//     a label, so there is a way to reach them
//   • what it holds is the tail of the list — every fact from the first
//     group that did not fit to the end — rather than whichever ones
//     happened to be narrow
//   • a fact activated from the menu reports the same two indices it would
//     have reported from the bar, so the caller cannot tell the difference
void checkWhatDoesNotFitStaysReachable(QQuickWindow *window, QQuickItem *bar)
{
    const int total = bar->property("groups").toList().size();
    QTRY_VERIFY(bar->property("shownGroups").toInt() < total);
    const int shown = bar->property("shownGroups").toInt();
    const QVariantList hidden = bar->property("hiddenFacts").toList();
    QVERIFY2(!hidden.isEmpty(), "a bar too narrow for its groups hid nothing");

    // The tail, group by group: the hidden facts start at the first group
    // that did not fit, run to the last group, and take every fact of each
    // in the order they were given.
    int expectedGroup = shown;
    int expectedFact = 0;
    for (const QVariant &entry : hidden) {
        const QVariantMap fact = entry.toMap();
        if (fact.value(QStringLiteral("fact")).toInt() == 0
            && expectedFact != 0) {
            expectedGroup += 1;
            expectedFact = 0;
        }
        QCOMPARE(fact.value(QStringLiteral("group")).toInt(), expectedGroup);
        QCOMPARE(fact.value(QStringLiteral("fact")).toInt(), expectedFact);
        expectedFact += 1;
    }
    QCOMPARE(expectedGroup, total - 1);

    auto *overflow = bar->findChild<QQuickItem *>(QStringLiteral("overflow"));
    QVERIFY(overflow);
    QTRY_VERIFY2(overflow->isVisible() && overflow->width() > 0,
                 "nothing on the bar said that anything was hidden");
    QCOMPARE(overflow->property("text").toString(),
             QStringLiteral("%1 more").arg(hidden.size()));

    QAccessibleInterface *overflowAccessible =
        QAccessible::queryAccessibleInterface(overflow);
    QVERIFY(overflowAccessible);
    QVERIFY2(!overflowAccessible->text(QAccessible::Description).isEmpty(),
             "the overflow control did not say which groups it holds");

    // A group that is not on the bar is out of the tab order rather than an
    // invisible stop on the way along it.
    const QList<QQuickItem *> wrappers =
        itemsNamed(bar, QStringLiteral("group"));
    QCOMPARE(wrappers.size(), total);
    for (int i = 0; i < wrappers.size(); ++i) {
        QCOMPARE(wrappers.at(i)->isEnabled(), i < shown);
        QCOMPARE(wrappers.at(i)->width() > 0, i < shown);
    }

    // And what is in the menu activates exactly as it would have on the bar.
    QObject *menu = bar->findChild<QObject *>(QStringLiteral("overflowMenu"));
    QVERIFY(menu);
    QTRY_COMPARE(menu->property("count").toInt(), hidden.size());
    QQuickItem *entry = nullptr;
    QVERIFY(QMetaObject::invokeMethod(menu, "itemAt",
                                      Q_RETURN_ARG(QQuickItem *, entry),
                                      Q_ARG(int, 0)));
    QVERIFY(entry);
    const QVariantMap firstHidden = hidden.at(0).toMap();
    QCOMPARE(entry->property("text").toString(),
             firstHidden.value(QStringLiteral("label")).toString());

    const int before = window->property("presses").toInt();
    QVERIFY(QMetaObject::invokeMethod(entry, "clicked"));
    QTRY_COMPARE(window->property("presses").toInt(), before + 1);
    QCOMPARE(window->property("lastGroup").toInt(),
             firstHidden.value(QStringLiteral("group")).toInt());
    QCOMPARE(window->property("lastFact").toInt(),
             firstHidden.value(QStringLiteral("fact")).toInt());
}

// The facts on a bar, in the order the keyboard walks them, checked against
// the order they are drawn in. A reader tabbing along a bar and a reader
// looking at it have to be walking the same list.
void checkTheKeyboardFollowsTheDrawing(QQuickWindow *window, QQuickItem *bar,
                                       const QStringList &expected)
{
    const QList<QQuickItem *> facts = itemsNamed(bar, QStringLiteral("fact"));
    QCOMPARE(facts.size(), expected.size());

    QStringList spoken;
    qreal previous = -1;
    for (QQuickItem *fact : facts) {
        QAccessibleInterface *accessible =
            QAccessible::queryAccessibleInterface(fact);
        QVERIFY2(accessible, "a status-bar fact has no accessible interface");
        QCOMPARE(accessible->role(), QAccessible::Link);
        spoken.append(accessible->text(QAccessible::Name));
        const qreal x = fact->mapToItem(bar, QPointF(0, 0)).x();
        QVERIFY2(x > previous, "the facts are not laid out left to right");
        previous = x;
    }
    QCOMPARE(spoken, expected);

    facts.first()->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(facts.first()->hasActiveFocus());
    for (int i = 1; i < facts.size(); ++i) {
        QTest::keyClick(window, Qt::Key_Tab);
        QTRY_VERIFY2(facts.at(i)->hasActiveFocus(),
                     "tab did not reach the facts in the order they are drawn");
    }
}

}   // namespace

void TestComponents::testTheStatusBarKeepsWhatDoesNotFitReachable()
{
    // What a bar does when it runs out of room.
    //
    // The consuming application's own bar sliced each of its two groups to
    // the first two items, so a third thing waiting for somebody was not on
    // the screen and nothing said it existed. This bar keeps the tail in a
    // menu behind a control that says how many there are.
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
            width: 900
            height: 120
            property int lastGroup: -1
            property int lastFact: -1
            property int presses: 0
            KvitStatusBar {
                objectName: "bar"
                anchors.left: parent.left
                anchors.right: parent.right
                anchors.bottom: parent.bottom
                activity: "Indexing 4 of 26 working copies"
                groups: [
                    {
                        "label": "Waiting on you",
                        "facts": [
                            { "text": "3 reviews", "symbol": "question" },
                            { "text": "1 conflict", "symbol": "warning" }
                        ]
                    },
                    {
                        "label": "Running",
                        "facts": [{ "text": "2 agents", "symbol": "robot" }]
                    }
                ]
                onFactActivated: (group, fact) => {
                    window.lastGroup = group
                    window.lastFact = fact
                    window.presses += 1
                }
            }
            KvitStatusBar {
                objectName: "resting"
                width: 400
                facts: ["1,284 notes", "last synced 14:02"]
            }
            // The same bar with its groups at the other end. It sits above
            // the bottom of the window rather than on it, because a bar drawn
            // over the first would take the pointer presses meant for that
            // one.
            KvitStatusBar {
                objectName: "leading"
                anchors.left: parent.left
                anchors.right: parent.right
                y: 40
                groupsAt: "left"
                activity: "Indexing 4 of 26 working copies"
                groups: [
                    {
                        "label": "Waiting on you",
                        "facts": [
                            { "text": "3 reviews", "symbol": "question" },
                            { "text": "1 conflict", "symbol": "warning" }
                        ]
                    },
                    {
                        "label": "Running",
                        "facts": [{ "text": "2 agents", "symbol": "robot" }]
                    }
                ]
                onFactActivated: (group, fact) => {
                    window.lastGroup = group
                    window.lastFact = fact
                    window.presses += 1
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/status-groups.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto *bar = window->findChild<QQuickItem *>(QStringLiteral("bar"));
    auto *resting = window->findChild<QQuickItem *>(QStringLiteral("resting"));
    auto *leading = window->findChild<QQuickItem *>(QStringLiteral("leading"));
    QVERIFY(bar);
    QVERIFY(resting);
    QVERIFY(leading);

    // A bar holding only text is as tall as the text. Only a control in the
    // slot at the end makes it grow, which is what the height has always
    // said and what the groups must not change.
    QCOMPARE(resting->implicitHeight(),
             qreal(KvitUi::DefaultServices::interfaceMetrics()->statusBarHeight()));

    // Wide enough for both groups: nothing is in the menu, at either end.
    QTRY_COMPARE(bar->property("shownGroups").toInt(), 2);
    QCOMPARE(bar->property("hiddenFacts").toList().size(), 0);
    QTRY_COMPARE(leading->property("shownGroups").toInt(), 2);
    QCOMPARE(leading->property("hiddenFacts").toList().size(), 0);

    // The facts are controls, and they run left to right in the order they
    // were given, which is the order the keyboard walks them in.
    const QStringList inOrder{ QStringLiteral("3 reviews"),
                               QStringLiteral("1 conflict"),
                               QStringLiteral("2 agents") };
    checkTheKeyboardFollowsTheDrawing(window, bar, inOrder);
    checkTheKeyboardFollowsTheDrawing(window, leading, inOrder);

    const QList<QQuickItem *> facts = itemsNamed(bar, QStringLiteral("fact"));
    QCOMPARE(facts.size(), 3);

    // Pressed by the keyboard, and by the pointer, with the two indices the
    // caller needs to know which one it was.
    facts.at(1)->forceActiveFocus(Qt::TabFocusReason);
    QTRY_VERIFY(facts.at(1)->hasActiveFocus());
    QTest::keyClick(window, Qt::Key_Return);
    QTRY_COMPARE(window->property("presses").toInt(), 1);
    QCOMPARE(window->property("lastGroup").toInt(), 0);
    QCOMPARE(window->property("lastFact").toInt(), 1);

    if (QGuiApplication::platformName() != QLatin1String("minimal")) {
        const QPoint overFact = window->contentItem()
            ->mapFromItem(facts.at(2), QPointF(facts.at(2)->width() / 2,
                                               facts.at(2)->height() / 2))
            .toPoint();
        QTest::mouseClick(window, Qt::LeftButton, Qt::NoModifier, overFact);
        QTRY_COMPARE(window->property("presses").toInt(), 2);
        QCOMPARE(window->property("lastGroup").toInt(), 1);
        QCOMPARE(window->property("lastFact").toInt(), 0);
    }

    // Now take the room away. What no longer fits is not dropped: the bar
    // says how many there are and keeps them behind a control — at whichever
    // end the groups sit, which is why both bars are put through it.
    window->setWidth(300);
    checkWhatDoesNotFitStaysReachable(window, bar);
    checkWhatDoesNotFitStaysReachable(window, leading);
}

void TestComponents::testTheStatusBarCanPutItsGroupsFirst()
{
    // Which end of the bar the groups sit at.
    //
    // The bar was written for a window whose left end is a sentence about
    // what is running, so the groups sit at the right, beside the standing
    // facts. The editor this layer serves draws the other bar: its left end
    // lists what is waiting for the reader and what is running, each item
    // something to press, with the standing facts at the right. Before
    // `groupsAt` that window could not use this component at all, because
    // taking it meant moving its groups across the screen.
    //
    // Two things are checked here. That `groupsAt: "left"` puts the groups
    // before the activity and hard against the left margin, and that a bar
    // that does not ask for it is laid out where it always was — the rest of
    // what "where it always was" means is pinned by the case above, which is
    // unchanged.
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
            width: 900
            height: 400

            // One set of groups, read by all four bars, so that the only
            // difference between them is the two properties under test.
            readonly property var sample: [
                {
                    "label": "Waiting on you",
                    "facts": [
                        { "text": "3 reviews", "symbol": "question" },
                        { "text": "1 conflict", "symbol": "warning" }
                    ]
                },
                {
                    "label": "Running",
                    "facts": [{ "text": "2 agents", "symbol": "robot" }]
                }
            ]

            Column {
                anchors.fill: parent

                KvitStatusBar {
                    objectName: "rightWithActivity"
                    width: parent.width
                    activity: "Indexing 4 of 26 working copies"
                    groups: window.sample
                    facts: ["7 changes"]
                }
                KvitStatusBar {
                    objectName: "rightBare"
                    width: parent.width
                    groups: window.sample
                }
                KvitStatusBar {
                    objectName: "leftWithActivity"
                    width: parent.width
                    groupsAt: "left"
                    activity: "Indexing 4 of 26 working copies"
                    groups: window.sample
                    facts: ["7 changes"]
                }
                KvitStatusBar {
                    objectName: "leftBare"
                    width: parent.width
                    groupsAt: "left"
                    groups: window.sample
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/status-order.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    const qreal margin = KvitUi::DefaultServices::interfaceMetrics()->space();

    auto bar = [window](const char *name) {
        return window->findChild<QQuickItem *>(QLatin1String(name));
    };
    // Where the groups start, in the bar's own coordinates.
    auto groupsStart = [](QQuickItem *item) {
        const QList<QQuickItem *> wrappers =
            itemsNamed(item, QStringLiteral("group"));
        return wrappers.isEmpty()
            ? -1.0 : wrappers.first()->mapToItem(item, QPointF()).x();
    };
    // The activity the bar is actually drawing. There are two labels for it,
    // one at each end of the groups, and the one the order did not ask for is
    // invisible: a layout leaves an invisible child out, so it takes no room
    // and nothing on the bar moves because of it.
    auto activity = [](QQuickItem *item) {
        QQuickItem *drawn = nullptr;
        int visible = 0;
        const QList<QQuickItem *> labels =
            itemsNamed(item, QStringLiteral("activity"));
        for (QQuickItem *label : labels) {
            if (!label->isVisible())
                continue;
            visible += 1;
            drawn = label;
        }
        return visible == 1 ? drawn : nullptr;
    };

    for (const char *name : { "rightWithActivity", "rightBare",
                              "leftWithActivity", "leftBare" }) {
        QVERIFY2(bar(name), name);
        QTRY_COMPARE(bar(name)->property("shownGroups").toInt(), 2);
    }

    // The default. The activity is at the left margin and the groups are past
    // the middle of the bar, beside the standing facts at the right.
    QQuickItem *right = bar("rightWithActivity");
    QQuickItem *rightActivity = activity(right);
    QVERIFY2(rightActivity, "the default bar drew no activity, or drew two");
    QCOMPARE(rightActivity->mapToItem(right, QPointF()).x(), margin);
    QVERIFY2(groupsStart(right) > rightActivity->mapToItem(right, QPointF()).x(),
             "the groups were drawn before the activity by default");
    QVERIFY2(groupsStart(right) > right->width() / 2,
             "the groups were not at the right-hand end of the bar");

    // And with nothing happening, they stay at the right-hand end rather than
    // falling back to the left: the room the activity would have taken is
    // held open.
    QVERIFY2(!activity(bar("rightBare")),
             "a bar with no activity drew one anyway");
    QVERIFY2(groupsStart(bar("rightBare")) > bar("rightBare")->width() / 2,
             "a bar with no activity moved its groups to the left");

    // The new order. The groups are hard against the left margin, and the
    // activity is drawn after them rather than before.
    QQuickItem *left = bar("leftWithActivity");
    QQuickItem *leftActivity = activity(left);
    QVERIFY2(leftActivity, "the bar drew no activity, or drew two");
    QCOMPARE(groupsStart(left), margin);
    const QList<QQuickItem *> leftFacts =
        itemsNamed(left, QStringLiteral("fact"));
    QCOMPARE(leftFacts.size(), 3);
    const QQuickItem *lastFact = leftFacts.last();
    QVERIFY2(leftActivity->mapToItem(left, QPointF()).x()
                 > lastFact->mapToItem(left, QPointF()).x(),
             "the activity was drawn before the groups");

    // With nothing happening they are still at the left margin, which is the
    // case the consuming window is actually in most of the time.
    QVERIFY2(!activity(bar("leftBare")),
             "a bar with no activity drew one anyway");
    QCOMPARE(groupsStart(bar("leftBare")), margin);

    // How much room the groups get does not depend on which end they sit at.
    // The measurement counts what the activity, the facts and the controls
    // take, and none of those changed size by moving.
    window->setWidth(420);
    QTRY_VERIFY(right->property("shownGroups").toInt() < 2);
    QTRY_COMPARE(left->property("shownGroups").toInt(),
                 right->property("shownGroups").toInt());
    QCOMPARE(left->property("hiddenFacts").toList().size(),
             right->property("hiddenFacts").toList().size());

    // Narrower still, until neither has room for anything: still the same
    // answer at both ends.
    window->setWidth(200);
    QTRY_COMPARE(left->property("shownGroups").toInt(),
                 right->property("shownGroups").toInt());
}

void TestComponents::testTheStatusBarSurvivesItsGroupsBeingReplaced()
{
    // A bar whose groups change while it is on the screen.
    //
    // The gallery gives its specimens one array and never changes it, so
    // nothing here had ever replaced a bar's groups. A consuming application
    // replaces them whenever the work the bar reports changes, which for a
    // status bar is constantly, and every replacement wrote one warning per
    // fact:
    //
    //     TypeError: Cannot read property 'verticalCenter' of null
    //
    // The bar drew correctly before and after, so nothing on the screen said
    // anything was wrong. What it cost was the consuming application, whose
    // test gate fails any run that produced a QML warning: adopting this
    // component took that suite from three failing registrations to eight.
    //
    // The cause is that a Repeater unparents a delegate before destroying it,
    // and a delegate anchored to `parent` re-evaluates that binding in the
    // window where `parent` is null. So the check here is not that the bar
    // still draws -- it always did -- but that the console stays clean.
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
            width: 900
            height: 200

            property var reported: [
                {
                    "label": "Waiting on you",
                    "facts": [
                        { "text": "3 reviews", "symbol": "question" },
                        { "text": "1 conflict", "symbol": "warning" }
                    ]
                },
                {
                    "label": "Running",
                    "facts": [{ "text": "2 agents", "symbol": "robot" }]
                }
            ]
            property var plain: ["1,284 notes", "last synced 14:02"]
            property int lastGroup: -1
            property int lastFact: -1
            property int presses: 0

            Column {
                anchors.fill: parent

                KvitStatusBar {
                    objectName: "trailing"
                    width: parent.width
                    activity: "Indexing 4 of 26 working copies"
                    groups: window.reported
                    facts: window.plain
                    onFactActivated: (group, fact) => {
                        window.lastGroup = group
                        window.lastFact = fact
                        window.presses += 1
                    }
                }
                KvitStatusBar {
                    objectName: "leading"
                    width: parent.width
                    groupsAt: "left"
                    activity: "Indexing 4 of 26 working copies"
                    groups: window.reported
                    facts: window.plain
                    onFactActivated: (group, fact) => {
                        window.lastGroup = group
                        window.lastFact = fact
                        window.presses += 1
                    }
                }
            }
        }
    )", QUrl(QStringLiteral("qrc:/test/status-replaced.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> instance(component.create());
    QVERIFY2(!instance.isNull(), qPrintable(component.errorString()));
    auto *window = qobject_cast<QQuickWindow *>(instance.data());
    QVERIFY(window);
    window->show();
    window->requestActivate();
    QVERIFY(QTest::qWaitForWindowExposed(window));

    auto bar = [window](const char *name) {
        return window->findChild<QQuickItem *>(QLatin1String(name));
    };
    QVERIFY(bar("trailing"));
    QVERIFY(bar("leading"));
    QTRY_COMPARE(bar("trailing")->property("shownGroups").toInt(), 2);
    QTRY_COMPARE(bar("leading")->property("shownGroups").toInt(), 2);

    // Four replacements, because they tear down different things: a group
    // dropped, a fact dropped from inside a group that stays, a group added
    // back, and everything taken away at once. The plain facts beside the
    // groups are replaced with them, since those are a repeater too.
    const QList<QPair<QString, QString>> rounds{
        { QStringLiteral("one group instead of two"),
          QStringLiteral(R"([{ "label": "Running",
                               "facts": [{ "text": "2 agents",
                                           "symbol": "robot" }] }])") },
        { QStringLiteral("a fact dropped from a group that stays"),
          QStringLiteral(R"([{ "label": "Running",
                               "facts": [{ "text": "1 agent",
                                           "symbol": "robot" },
                                         { "text": "4 queued",
                                           "symbol": "clock" }] }])") },
        { QStringLiteral("three groups where there was one"),
          QStringLiteral(R"([{ "label": "Waiting on you",
                               "facts": [{ "text": "9 reviews",
                                           "symbol": "question" }] },
                             { "label": "Running",
                               "facts": [{ "text": "2 agents",
                                           "symbol": "robot" }] },
                             { "label": "Failed",
                               "facts": [{ "text": "1 run",
                                           "symbol": "warning" }] }])") },
        { QStringLiteral("nothing left to report"), QStringLiteral("[]") },
    };

    int round = 0;
    for (const auto &[what, groups] : rounds) {
        g_warnings.clear();
        QQmlExpression assign(
            qmlContext(window), window,
            QStringLiteral("reported = %1; plain = [\"%2 notes\"]")
                .arg(groups).arg(1000 + round));
        assign.evaluate();
        QVERIFY2(!assign.hasError(), qPrintable(assign.error().toString()));
        QCoreApplication::processEvents();
        QTest::qWait(50);
        QCoreApplication::processEvents();
        QVERIFY2(g_warnings.isEmpty(),
                 qPrintable(QStringLiteral("replacing the groups with %1 wrote:\n  %2")
                                .arg(what, g_warnings.join(QStringLiteral("\n  ")))));
        ++round;
    }

    // The reason for the anchor has not gone away: a fact still sits on the
    // bar's vertical centre, at either end.
    g_warnings.clear();
    QQmlExpression restore(
        qmlContext(window), window,
        QStringLiteral(R"(reported = [{ "label": "Running",
                                        "facts": [{ "text": "2 agents",
                                                    "symbol": "robot" }] }])"));
    restore.evaluate();
    QVERIFY2(!restore.hasError(), qPrintable(restore.error().toString()));
    QCoreApplication::processEvents();
    for (const char *name : { "trailing", "leading" }) {
        QQuickItem *item = bar(name);
        QTRY_COMPARE(itemsNamed(item, QStringLiteral("fact")).size(), 1);
        QQuickItem *fact = itemsNamed(item, QStringLiteral("fact")).first();
        const qreal centre =
            fact->mapToItem(item, QPointF(0, fact->height() / 2)).y();
        QVERIFY2(qAbs(centre - item->height() / 2) <= 1.0,
                 qPrintable(QStringLiteral("%1: the fact sits at %2 on a bar "
                                           "%3 tall")
                                .arg(QLatin1String(name)).arg(centre)
                                .arg(item->height())));
    }
    QVERIFY2(g_warnings.isEmpty(),
             qPrintable(g_warnings.join(QStringLiteral("\n  "))));

    // And the overflow still behaves when the room runs out, on a bar whose
    // model was replaced rather than set once.
    g_warnings.clear();
    QQmlExpression crowd(
        qmlContext(window), window,
        QStringLiteral(R"(reported = [{ "label": "Waiting on you",
                                        "facts": [{ "text": "3 reviews",
                                                    "symbol": "question" },
                                                  { "text": "1 conflict",
                                                    "symbol": "warning" }] },
                                      { "label": "Running",
                                        "facts": [{ "text": "2 agents",
                                                    "symbol": "robot" }] }])"));
    crowd.evaluate();
    QVERIFY2(!crowd.hasError(), qPrintable(crowd.error().toString()));
    QCoreApplication::processEvents();
    QTRY_COMPARE(bar("trailing")->property("shownGroups").toInt(), 2);
    window->setWidth(300);
    checkWhatDoesNotFitStaysReachable(window, bar("trailing"));
    checkWhatDoesNotFitStaysReachable(window, bar("leading"));
    QVERIFY2(g_warnings.isEmpty(),
             qPrintable(QStringLiteral("narrowing a replaced bar wrote:\n  %1")
                            .arg(g_warnings.join(QStringLiteral("\n  ")))));
}

QTEST_MAIN(TestComponents)
#include "test_components.moc"
