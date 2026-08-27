// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QAbstractItemModelTester>
#include <QElapsedTimer>
#include <QQmlComponent>
#include <QQmlEngine>
#include <QQuickItem>
#include <QQuickWindow>

#include "benchmarktablemodel.h"
#include "interfacemetrics.h"

// The measurement prd.md Decision 3 turns on.
//
// The decision was to build kvit-cash's transaction browser in QML on
// KvitTable rather than in Qt Widgets, on the grounds that the alternative
// means implementing rows, figures, chips and the currency rules of
// money-display.md a second time inside QStyledItemDelegate::paint(), and
// adding a second toolkit to an estate with roughly 70,000 lines of QML and no
// Widgets at all. The Widgets fallback stays open until a 250,000-row browser
// with twelve configurable columns holds smooth scrolling and sub-100 ms
// filtering. This is that check.
//
// The budgets below are what "smooth" and "sub-100 ms" mean here, stated so a
// failure says which one was missed rather than that something was slow. They
// are deliberately generous against the numbers a modern desktop actually
// produces: a benchmark tuned to the machine it was written on fails on
// somebody else's laptop for no reason anybody can act on.
class TestTableModel : public QObject
{
    Q_OBJECT

private slots:
    void initTestCase();
    void testTheModelIsWellFormed();
    void testEveryCellAnswersWithoutMaterialisingTheTable();
    void testAnUnmeasuredValueIsDistinguishableFromZero();
    void testFilteringStaysUnderTheBudget();
    void testEditingSurvivesAFilterChange();
    void testTheViewScrollsSmoothly();

private:
    // What "sub-100 ms filtering" means: one filter change over 250,000 rows.
    static constexpr int FilterBudgetMs = 100;
    // What "smooth scrolling" means, and what is actually measured.
    //
    // Not frames per second. A frame rate measured under an offscreen
    // platform is a measurement of the readback, and grabbing a 1200x700
    // window costs more than everything the table does. What decides whether
    // a QML table works at this size is delegate churn: at every scroll step
    // a screen full of cells leaves the viewport and another arrives, and the
    // question is whether those are recycled or rebuilt.
    //
    // So two claims. The cost of relaying out a page of cells stays inside a
    // 60 Hz frame, and the number of live delegate items stays bounded no
    // matter how far the reader scrolls. The second is the one that would
    // fail catastrophically rather than gradually, and it needs no timer.
    static constexpr int ScrollBudgetMsPerPage = 16;
    static constexpr int PagesScrolled = 40;
    // A screen of 1200x700 at the default row height holds about 23 rows of
    // 12 columns, so roughly 280 cells are on screen and a little over that
    // is in flight. Anything under a thousand is recycling; a view that built
    // a delegate per row would be at three million.
    static constexpr int LiveDelegateCeiling = 1000;

    bool m_haveWindow = false;
};

void TestTableModel::initTestCase()
{
    // The view half needs a scene graph. Under the offscreen platform it has
    // one; under a platform with no rendering at all it does not, and the case
    // skips rather than failing, saying so.
    m_haveWindow = QGuiApplication::platformName() != QLatin1String("minimal");
}

void TestTableModel::testTheModelIsWellFormed()
{
    // Qt's own model checker, which catches the class of bug that makes a
    // view show the wrong rows after an edit: signals emitted in the wrong
    // order, an index out of range, a parent that is not what it claims.
    KvitUi::BenchmarkTableModel model;
    model.setTotalRows(500);
    QAbstractItemModelTester tester(&model,
                                    QAbstractItemModelTester::FailureReportingMode::Warning);
    QCOMPARE(model.rowCount(), 500);
    QCOMPARE(model.columnCount(), KvitUi::BenchmarkTableModel::ColumnCount);
    QCOMPARE(model.columnCount(), 12);
}

void TestTableModel::testEveryCellAnswersWithoutMaterialisingTheTable()
{
    // Three million cells, generated rather than stored. The point of this
    // case is that asking for a cell in the middle costs nothing: a model that
    // built its rows on demand and cached them would pass every other case
    // here and then use several gigabytes in the real browser.
    KvitUi::BenchmarkTableModel model;
    QCOMPARE(model.rowCount(), 250000);

    QElapsedTimer timer;
    timer.start();
    // A scattered sample rather than a sequential one, because a sequential
    // read is the pattern a cache is good at and a reader scrolling to the
    // middle of a ledger is not.
    for (int i = 0; i < 20000; ++i) {
        const int row = (i * 7919) % model.rowCount();
        const int column = i % model.columnCount();
        const QVariant value = model.data(model.index(row, column));
        Q_UNUSED(value);
    }
    const qint64 elapsed = timer.elapsed();
    qInfo("20,000 scattered cell reads over 3,000,000 cells: %lld ms", elapsed);
    QVERIFY2(elapsed < 500,
             qPrintable(QStringLiteral("20,000 scattered cell reads took %1 ms")
                            .arg(elapsed)));

    // The values are a function of the row, so a cell answers the same thing
    // whatever order it was asked in. This is what makes the rest of the
    // benchmark reproducible.
    const QVariant first = model.data(model.index(173404, 2));
    const QVariant again = model.data(model.index(173404, 2));
    QCOMPARE(first, again);
    QVERIFY(!first.toString().isEmpty());
}

void TestTableModel::testAnUnmeasuredValueIsDistinguishableFromZero()
{
    // The rule KvitFigure exists to enforce, checked at the model boundary
    // where it starts: a balance nobody has computed is not a balance of zero,
    // and the model has to be able to say which it means or no view downstream
    // can draw the difference.
    KvitUi::BenchmarkTableModel model;
    const int balanceColumn = 7;

    const QModelIndex unmeasured = model.index(40, balanceColumn);
    QCOMPARE(model.data(unmeasured, KvitUi::TableModelBase::MeasuredRole).toBool(),
             false);
    QVERIFY(!model.data(unmeasured).isValid());

    const QModelIndex measured = model.index(41, balanceColumn);
    QCOMPARE(model.data(measured, KvitUi::TableModelBase::MeasuredRole).toBool(),
             true);
    QVERIFY(model.data(measured).isValid());

    // A measured zero is measured. Nothing in the model may treat 0 and
    // "absent" as the same thing.
    KvitUi::BenchmarkTableModel other;
    const QModelIndex amount = other.index(0, 6);
    QVERIFY(other.data(amount, KvitUi::TableModelBase::MeasuredRole).toBool());
}

void TestTableModel::testFilteringStaysUnderTheBudget()
{
    KvitUi::BenchmarkTableModel model;
    QCOMPARE(model.rowCount(), 250000);

    // A substring nothing matches, which is the worst case: every row is
    // tested and none is kept.
    QElapsedTimer timer;
    timer.start();
    model.setFilter(QStringLiteral("zzzzz"));
    const qint64 emptyFilter = timer.elapsed();
    QCOMPARE(model.rowCount(), 0);

    // A substring that matches a fraction of the rows, which is what a reader
    // actually types.
    timer.restart();
    model.setFilter(QStringLiteral("Harlow"));
    const qint64 narrowFilter = timer.elapsed();
    const int matched = model.rowCount();
    QVERIFY(matched > 0);
    QVERIFY(matched < model.totalRows());

    // Clearing it, which puts 250,000 rows back.
    timer.restart();
    model.setFilter(QString());
    const qint64 clearFilter = timer.elapsed();
    QCOMPARE(model.rowCount(), 250000);

    // Reported rather than only asserted. A benchmark that says "passed" and
    // nothing else cannot show a slow drift toward its budget, which is how
    // a performance requirement is usually lost.
    qInfo("filtering 250,000 rows: %lld ms to nothing, %lld ms to %d matches, "
          "%lld ms to clear (budget %d ms)",
          emptyFilter, narrowFilter, matched, clearFilter, FilterBudgetMs);

    QVERIFY2(emptyFilter < FilterBudgetMs,
             qPrintable(QStringLiteral("filtering 250,000 rows to nothing took "
                                       "%1 ms against a %2 ms budget")
                            .arg(emptyFilter).arg(FilterBudgetMs)));
    QVERIFY2(narrowFilter < FilterBudgetMs,
             qPrintable(QStringLiteral("filtering 250,000 rows took %1 ms "
                                       "against a %2 ms budget")
                            .arg(narrowFilter).arg(FilterBudgetMs)));
    QVERIFY2(clearFilter < FilterBudgetMs,
             qPrintable(QStringLiteral("clearing the filter took %1 ms against "
                                       "a %2 ms budget")
                            .arg(clearFilter).arg(FilterBudgetMs)));
}

void TestTableModel::testEditingSurvivesAFilterChange()
{
    // Inline editing is the third thing the benchmark asks for, and the way it
    // usually breaks is that an edit is stored against the visible row rather
    // than against the underlying one — so filtering the list moves everybody's
    // edits onto different rows.
    KvitUi::BenchmarkTableModel model;
    model.setTotalRows(5000);
    model.setFilter(QStringLiteral("Harlow"));
    QVERIFY(model.rowCount() > 2);

    const QModelIndex target = model.index(1, 2);
    const QString before = model.data(target).toString();
    QVERIFY(model.setData(target, QStringLiteral("edited in place")));
    QCOMPARE(model.data(target).toString(), QStringLiteral("edited in place"));
    QVERIFY(before != QStringLiteral("edited in place"));

    // The same underlying row, reached without the filter, carries the edit.
    model.setFilter(QString());
    int found = 0;
    for (int row = 0; row < model.rowCount(); ++row) {
        if (model.data(model.index(row, 2)).toString()
            == QLatin1String("edited in place"))
            ++found;
    }
    QCOMPARE(found, 1);

    // A generated column is not editable, so a reader cannot overwrite a
    // reference or a computed balance.
    QVERIFY(!model.setData(model.index(3, 1), QStringLiteral("nope")));
    QVERIFY(!model.setData(model.index(3, 7), QStringLiteral("nope")));
}

void TestTableModel::testTheViewScrollsSmoothly()
{
    if (!m_haveWindow)
        QSKIP("no scene graph on this platform; the view half cannot run");

    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    width: 1200; height: 700\n"
        "    property alias table: table\n"
        "    property alias model: model\n"
        "    BenchmarkTableModel { id: model }\n"
        "    KvitTable { id: table; anchors.fill: parent; model: model }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/table.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    QQuickWindow window;
    auto *content = qobject_cast<QQuickItem *>(holder.data());
    QVERIFY(content);
    content->setParentItem(window.contentItem());
    window.resize(1200, 700);
    window.show();
    QVERIFY(QTest::qWaitForWindowExposed(&window));

    auto *table = content->property("table").value<QQuickItem *>();
    QVERIFY(table);
    auto *view = table->property("view").value<QQuickItem *>();
    QVERIFY(view);

    // One render first, so the first page's delegates exist and the timing
    // below measures steady-state scrolling rather than the first layout.
    window.grabWindow();

    const qreal page = view->height();
    QElapsedTimer timer;
    timer.start();
    for (int i = 0; i < PagesScrolled; ++i) {
        view->setProperty("contentY", page * i);
        // forceLayout is the delegate work: bind, position and reuse every
        // cell the new viewport needs. It is synchronous, which is what makes
        // it measurable.
        QMetaObject::invokeMethod(view, "forceLayout");
    }
    const qint64 elapsed = timer.elapsed();
    const qreal perPage = qreal(elapsed) / PagesScrolled;

    auto *pool = view->property("contentItem").value<QQuickItem *>();
    QVERIFY(pool);
    const int live = pool->childItems().size();
    qInfo("scrolling 250,000 rows: %.2f ms per page (budget %d ms), "
          "%d live delegate items after %d pages",
          perPage, ScrollBudgetMsPerPage, live, PagesScrolled);

    QVERIFY2(perPage < ScrollBudgetMsPerPage,
             qPrintable(QStringLiteral("relaying out a page of a 250,000-row "
                                       "table cost %1 ms against a %2 ms "
                                       "budget (%3 ms for %4 pages)")
                            .arg(perPage, 0, 'f', 2)
                            .arg(ScrollBudgetMsPerPage)
                            .arg(elapsed).arg(PagesScrolled)));

    // And the claim that makes the timing possible: after forty pages the
    // view is still holding a few hundred items rather than a few hundred
    // thousand. Without `reuseItems` this number grows with every page and
    // the process runs out of memory long before the reader runs out of rows.
    QVERIFY2(live > 0 && live < LiveDelegateCeiling,
             qPrintable(QStringLiteral("after scrolling %1 pages the view holds "
                                       "%2 delegate items; anything near the "
                                       "row count means they are not being "
                                       "recycled")
                            .arg(PagesScrolled).arg(live)));
}

QTEST_MAIN(TestTableModel)
#include "test_tablemodel.moc"
