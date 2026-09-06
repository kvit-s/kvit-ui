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
#include <QQuickWindow>
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
    void testTheRailShowsACountOnlyWhenAskedTo();
    void testACountedNameInflectsTheCallersNoun();
    void testAGroupHeadingCountsItsGroupTheWayEverythingElseDoes();
    void testAProgressBarCannotTrailTheNumberBesideIt();
    void testARowSaysWhetherItIsPressableBeforeItIsPressed();
    void testARowAnnouncesItselfAsWhatItActuallyIs();
    void testAShortenedValueIsDisclosedToTheKeyboardAsWellAsThePointer();

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

QTEST_MAIN(TestComponents)
#include "test_components.moc"
