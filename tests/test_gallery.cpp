// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QDirIterator>
#include <QTemporaryDir>
#include <QQmlComponent>
#include <QQmlEngine>
#include <QQuickItem>

#include "catalogwriter.h"
#include "theme.h"
#include "uiservices.h"

// The gallery's catalogue, and the two promises it makes.
//
// Every component in the module has a page — a component nobody put in the
// catalogue is a component that never appears in a screenshot set and never
// gets looked at in four themes, which is the whole mechanism prd.md §7 rests
// on.
//
// And every code sample in it compiles and runs. The samples are the same
// strings the gallery renders, so a sample that does not work is a page that
// draws nothing; and since they are strings rather than QML files, nothing
// else in the build would notice. This is what makes them safe to copy.
class TestGallery : public QObject
{
    Q_OBJECT

private slots:
    void initTestCase();
    void testEveryComponentHasAPage();
    void testEverySnippetBuilds_data();
    void testEverySnippetBuilds();
    void testEveryPageSaysWhatTheComponentIsFor();
    void testAlignmentSensitiveControlsCoverEveryInterfaceExtreme();
    void testTheComponentListNavigates();
    void testTheCatalogueWriterRuns();

private:
    QQmlEngine m_engine;
    QVariantList m_catalog;
    // A sized item to build each specimen inside, so a sample that binds to
    // its parent's width has one.
    QQuickItem *m_stage = nullptr;
};

namespace {

// Every warning QML emits while a specimen is being built.
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

namespace {

QStringList moduleComponents()
{
    QStringList names;
    QDirIterator it(QStringLiteral(":/qt/qml/Kvit/Ui"), { QStringLiteral("*.qml") },
                    QDir::Files);
    while (it.hasNext()) {
        const QString path = it.next();
        names.append(QFileInfo(path).completeBaseName());
    }
    names.sort();
    return names;
}

}   // namespace

void TestGallery::initTestCase()
{
    // Reach into kvit-ui-qml from C++ before touching QML.
    //
    // This is not decoration. With static libraries the linker keeps only the
    // archive members something references, and a QML module's type
    // registration and its resource initialiser are members nothing
    // references by name — so a test that reaches Kvit.Ui only through
    // `import` in a string gets "module is not installed" and an empty
    // resource tree. Naming one symbol from the library pulls it in.
    KvitUi::DefaultServices::theme()->setThemeId(QStringLiteral("dark"));

    g_previous = qInstallMessageHandler(collect);

    QQmlComponent stage(&m_engine);
    stage.setData("import QtQuick\nItem { width: 600; height: 400 }",
                  QUrl(QStringLiteral("qrc:/test/stage.qml")));
    QVERIFY2(stage.isReady(), qPrintable(stage.errorString()));
    m_stage = qobject_cast<QQuickItem *>(stage.create());
    QVERIFY(m_stage);

    QQmlComponent component(&m_engine);
    component.setData("import QtQml\nimport Kvit.Gallery\n"
                      "QtObject { property var all: Catalog.components }",
                      QUrl(QStringLiteral("qrc:/test/catalog.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));
    m_catalog = holder->property("all").toList();
    QVERIFY(!m_catalog.isEmpty());
}

void TestGallery::testEveryComponentHasAPage()
{
    QSet<QString> documented;
    for (const QVariant &entry : m_catalog)
        documented.insert(entry.toMap().value(QStringLiteral("name")).toString());

    QStringList undocumented;
    for (const QString &name : moduleComponents()) {
        if (!documented.contains(name))
            undocumented.append(name);
    }
    QVERIFY2(undocumented.isEmpty(),
             qPrintable(QStringLiteral("these components have no gallery page, "
                                       "so nothing renders them in four "
                                       "themes: %1")
                            .arg(undocumented.join(QStringLiteral(", ")))));

    // And nothing in the catalogue that is not a component, which is what a
    // rename leaves behind.
    const QStringList existing = moduleComponents();
    QStringList stale;
    for (const QString &name : documented) {
        if (!existing.contains(name))
            stale.append(name);
    }
    stale.sort();
    QVERIFY2(stale.isEmpty(),
             qPrintable(QStringLiteral("the catalogue documents components that "
                                       "do not exist: %1")
                            .arg(stale.join(QStringLiteral(", ")))));
}

void TestGallery::testEverySnippetBuilds_data()
{
    QTest::addColumn<QString>("component");
    QTest::addColumn<QString>("caption");
    QTest::addColumn<QString>("snippet");

    for (const QVariant &entry : m_catalog) {
        const QVariantMap map = entry.toMap();
        const QString name = map.value(QStringLiteral("name")).toString();
        const QVariantList specimens =
            map.value(QStringLiteral("specimens")).toList();
        for (const QVariant &specimen : specimens) {
            const QVariantMap one = specimen.toMap();
            const QString caption = one.value(QStringLiteral("caption")).toString();
            QTest::newRow(qPrintable(name + QStringLiteral(" — ") + caption))
                << name << caption
                << one.value(QStringLiteral("snippet")).toString();
        }
    }
}

void TestGallery::testEverySnippetBuilds()
{
    QFETCH(QString, snippet);
    QFETCH(QString, component);
    QFETCH(QString, caption);

    // The same imports the gallery prepends, so what is compiled here is what
    // it renders.
    const QString source = QStringLiteral(
        "import QtQuick\nimport QtQuick.Controls\nimport QtQuick.Layouts\n"
        "import Kvit.Ui\n") + snippet;

    QQmlComponent built(&m_engine);
    built.setData(source.toUtf8(),
                  QUrl(QStringLiteral("qrc:/test/specimen.qml")));
    QVERIFY2(built.isReady(),
             qPrintable(QStringLiteral("%1 — %2 does not compile:\n%3")
                            .arg(component, caption, built.errorString())));

    // Built in two halves with a parent set in between, which is what the
    // gallery does when it puts a specimen on its stage. A sample written to
    // sit inside something — and most are, because that is how a component is
    // actually used — binds to `parent.width`, and creating it in one go with
    // no parent evaluates that binding against null.
    g_warnings.clear();
    QScopedPointer<QObject> instance(built.beginCreate(m_engine.rootContext()));
    QVERIFY2(!instance.isNull(),
             qPrintable(QStringLiteral("%1 — %2 does not run:\n%3")
                            .arg(component, caption, built.errorString())));
    if (auto *item = qobject_cast<QQuickItem *>(instance.data()))
        item->setParentItem(m_stage);
    built.completeCreate();
    QCoreApplication::processEvents();

    // A sample that warns has not worked. The one exception is the icon page's
    // deliberate typo, which exists to show what an unrecognised name does.
    QStringList unexpected;
    for (const QString &warning : g_warnings) {
        if (!warning.contains(QLatin1String("not-a-symbol")))
            unexpected.append(warning);
    }
    QVERIFY2(unexpected.isEmpty(),
             qPrintable(QStringLiteral("%1 — %2 warned:\n  %3")
                            .arg(component, caption,
                                 unexpected.join(QStringLiteral("\n  ")))));
}

void TestGallery::testEveryPageSaysWhatTheComponentIsFor()
{
    // A page with a name and no summary is a page that says a component
    // exists without saying when to reach for it, which is the one thing the
    // vocabulary skill's catalogue has to carry.
    QStringList silent;
    for (const QVariant &entry : m_catalog) {
        const QVariantMap map = entry.toMap();
        const QString summary =
            map.value(QStringLiteral("summary")).toString();
        const QVariantList specimens =
            map.value(QStringLiteral("specimens")).toList();
        if (summary.size() < 40 || specimens.isEmpty())
            silent.append(map.value(QStringLiteral("name")).toString());
    }
    QVERIFY2(silent.isEmpty(),
             qPrintable(QStringLiteral("these pages have no summary or no "
                                       "specimens: %1")
                            .arg(silent.join(QStringLiteral(", ")))));
}

void TestGallery::testAlignmentSensitiveControlsCoverEveryInterfaceExtreme()
{
    // The Windows native-text path is where off-centre button content first
    // showed up. Keep named screenshot variants for both controls whose
    // content is centred by the shared wrapper, so that path is reviewed at
    // the same three sizes the geometry test exercises.
    const QStringList expected = {
        QStringLiteral("minimum"), QStringLiteral("default"),
        QStringLiteral("maximum"),
    };
    for (const QString &name : { QStringLiteral("KvitButton"),
                                 QStringLiteral("KvitTab") }) {
        QVariantMap found;
        for (const QVariant &entry : m_catalog) {
            const QVariantMap candidate = entry.toMap();
            if (candidate.value(QStringLiteral("name")).toString() == name) {
                found = candidate;
                break;
            }
        }
        QVERIFY2(!found.isEmpty(), qPrintable(name));

        QStringList actual;
        for (const QVariant &size :
             found.value(QStringLiteral("shotSizes")).toList()) {
            actual.append(size.toString());
        }
        QCOMPARE(actual, expected);
    }
}

void TestGallery::testTheComponentListNavigates()
{
    // Opening a component's page by name, which is what a click in the
    // component list and `--page` both go through.
    //
    // This is tested from C++ rather than left to the eye because the way it
    // failed says nothing. The call was `window.show(name)`; a Window already
    // has a `show()` slot that takes no arguments, so it re-showed a window
    // that was already visible, discarded the name, and left the first
    // component on screen. The only sign was one line on the console about an
    // argument being ignored. Nothing in the build noticed, because the
    // screenshot run assigns the page index itself and never calls this.
    QQmlComponent component(&m_engine);
    component.setData("import QtQuick\nimport Kvit.Gallery\n"
                      "Gallery { visible: false }",
                      QUrl(QStringLiteral("qrc:/test/gallery.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));
    QScopedPointer<QObject> window(component.create());
    QVERIFY2(!window.isNull(), qPrintable(component.errorString()));

    QCOMPARE(window->property("current").toInt(), 0);

    // The last component in the catalogue, so a page that does not move fails
    // rather than passing on the one already showing.
    const int last = m_catalog.size() - 1;
    const QString wanted =
        m_catalog.at(last).toMap().value(QStringLiteral("name")).toString();
    QVERIFY(QMetaObject::invokeMethod(window.data(), "showPage",
                                      Q_ARG(QString, wanted)));
    QCOMPARE(window->property("current").toInt(), last);
}

void TestGallery::testTheCatalogueWriterRuns()
{
    // The vocabulary skill's catalogue, written the way CI writes it.
    //
    // It is also what keeps this binary linked to the gallery module at all.
    // Built shared, the linker drops a library nothing references by symbol,
    // and `import Kvit.Gallery` then fails with "module is not installed" —
    // which looks like a QML problem and is a link-line one. Naming
    // writeCatalog is what stops that.
    QTemporaryDir dir;
    QVERIFY(dir.isValid());
    const QString path = dir.filePath(QStringLiteral("catalog.md"));
    QVERIFY(KvitUi::writeCatalog(&m_engine, path));

    QFile written(path);
    QVERIFY(written.open(QIODevice::ReadOnly | QIODevice::Text));
    const QString text = QString::fromUtf8(written.readAll());

    // Every component named, so an agent reading it is told about all of them.
    for (const QVariant &entry : m_catalog) {
        const QString name =
            entry.toMap().value(QStringLiteral("name")).toString();
        QVERIFY2(text.contains(QStringLiteral("### ") + name),
                 qPrintable(QStringLiteral("%1 is missing from the catalogue")
                                .arg(name)));
    }
    // And the symbol names, which are the other half of what a call site
    // needs and are not in the component list.
    QVERIFY(text.contains(QStringLiteral("`chevron-right`")));
}

QTEST_MAIN(TestGallery)
#include "test_gallery.moc"
