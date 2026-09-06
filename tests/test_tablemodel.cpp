// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QtTest/QtTest>
#include <QAbstractItemModelTester>
#include <QElapsedTimer>
#include <QLocale>
#include <QQmlComponent>
#include <QQmlEngine>
#include <QQuickItem>
#include <QQuickWindow>

#include <algorithm>

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
    void testTheColumnKindArrivesAsAWordRatherThanANumber();
    void testAMoneyColumnIsFormattedByTheModel();
    void testARowCanBeInSeveralStatesAtOnce();
    void testTicksSurviveAFilterChange();
    void testTheWholeValueIsAnsweredWhereTheColumnShortensIt();
    void testEveryCellAnswersWithoutMaterialisingTheTable();
    void testAnUnmeasuredValueIsDistinguishableFromZero();
    void testFilteringStaysUnderTheBudget();
    void testEditingSurvivesAFilterChange();
    void testTheViewScrollsSmoothly();
    void testTheEmptyStateFollowsTheRowCount();
    void testTickingABoxDoesNotAlsoOpenTheRow();

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
    //
    // The budget is spent against the middle page rather than against the
    // average of all of them. Each page is timed on its own and the median is
    // what has to fit, because this gate runs on a build machine directly
    // after a suite that writes a 131 MB workspace and a build that uses
    // every core: a page that waits on the machine rather than on the table
    // has been measured at four times the cost of its neighbours, and one
    // such page moves an average of forty enough to fail the whole run. It
    // does not move the middle one. A table that really costs more than a
    // frame to relay out costs it on every page, so the regression this
    // exists to catch still fails it. The worst page is printed beside the
    // median rather than asserted on.
    static constexpr int ScrollBudgetMsPerPage = 16;
    static constexpr int PagesScrolled = 40;
    // A screen of 1200x700 at the default row height holds about 23 rows of
    // 13 columns, so roughly 300 cells are on screen and a little over that
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
    // Twelve columns of value, which is the width the benchmark is stated in,
    // and one more holding the box the reader ticks.
    QCOMPARE(KvitUi::BenchmarkTableModel::ValueColumnCount, 12);
    QCOMPARE(model.columnCount(), 13);
}

void TestTableModel::testTheColumnKindArrivesAsAWordRatherThanANumber()
{
    // What KvitCell is given for `kind`, and the reason it exists.
    //
    // columnKind() answers with a Q_ENUM, and a Q_ENUM returned through a
    // Q_INVOKABLE reaches JavaScript as an integer — so the view's
    // `String(model.columnKind(column))` yielded "1" for every column, no
    // case in KvitCell's switch matched, and every column in every table in
    // the estate drew as plain text. The word is read out of the
    // enumeration's own metadata, so a kind added to the enumeration cannot
    // go missing from the name.
    KvitUi::BenchmarkTableModel model;
    QCOMPARE(model.columnKindName(0), QStringLiteral("Date"));
    QCOMPARE(model.columnKindName(1), QStringLiteral("Slug"));
    QCOMPARE(model.columnKindName(2), QStringLiteral("Text"));
    QCOMPARE(model.columnKindName(4), QStringLiteral("Chip"));
    QCOMPARE(model.columnKindName(6), QStringLiteral("Money"));
    QCOMPARE(model.columnKindName(7), QStringLiteral("Figure"));
    QCOMPARE(model.columnKindName(9), QStringLiteral("Marks"));
    QCOMPARE(model.columnKindName(12), QStringLiteral("Check"));

    // A column outside the table is text rather than nothing: an empty kind
    // draws nothing at all, and a column of nothing is harder to notice than
    // a column of text.
    QCOMPARE(model.columnKindName(-1), QStringLiteral("Text"));
    QCOMPARE(model.columnKindName(999), QStringLiteral("Text"));
}

void TestTableModel::testAMoneyColumnIsFormattedByTheModel()
{
    // The Money cell draws the string it is given and does no arithmetic, so
    // the model has to hand it one. Two facts about money make that
    // necessary. An amount is a signed 64-bit count of minor units and a
    // JavaScript number carries 53 bits of integer, so a cell converting the
    // value would round a large amount and never say it had; and how many
    // minor digits an amount has belongs to its currency, which is knowledge
    // a cell does not have.
    KvitUi::BenchmarkTableModel model;
    const QVariant amount = model.data(model.index(3, 6));
    QCOMPARE(amount.metaType().id(), QMetaType::QString);
    QCOMPARE(amount.toString(), QLocale::system().toString(-498.89, 'f', 2));

    // And the sort still runs on the number, so the column orders by value
    // rather than by the digits of its formatting.
    const QVariant sortValue =
        model.data(model.index(3, 6), KvitUi::TableModelBase::SortRole);
    QCOMPARE(sortValue.metaType().id(), QMetaType::Double);
    QCOMPARE(sortValue.toDouble(), -498.89);
}

void TestTableModel::testARowCanBeInSeveralStatesAtOnce()
{
    // Why the Marks kind exists: a transaction is regularly settled *and*
    // unreviewed, and one chip has to drop one of the two. Each mark carries
    // a tone, the shape that says the same thing without colour, and the word
    // a reader hears.
    KvitUi::BenchmarkTableModel model;
    const int statusColumn = 9;

    const QVariant one = model.data(model.index(1, statusColumn),
                                    KvitUi::TableModelBase::MarksRole);
    const QVariantList oneMark = one.toList();
    QCOMPARE(oneMark.size(), 1);
    const QVariantMap first = oneMark.at(0).toMap();
    QCOMPARE(first.value(QStringLiteral("label")).toString(),
             QStringLiteral("Pending"));
    QVERIFY(!first.value(QStringLiteral("tone")).toString().isEmpty());
    // The second channel beside the hue. Without it two states differ only in
    // colour, which is no difference at all for about one man in twelve and
    // none in a grayscale screenshot.
    QVERIFY(!first.value(QStringLiteral("shape")).toString().isEmpty());

    // Row 7 is both settled and unreviewed.
    const QVariantList two = model.data(model.index(7, statusColumn),
                                        KvitUi::TableModelBase::MarksRole).toList();
    QCOMPARE(two.size(), 2);
    QCOMPARE(two.at(1).toMap().value(QStringLiteral("label")).toString(),
             QStringLiteral("Not reviewed"));

    // Only the column whose kind is Marks answers with any, because the kind
    // belongs to the column.
    QVERIFY(!model.data(model.index(7, 2),
                        KvitUi::TableModelBase::MarksRole).isValid());
}

void TestTableModel::testTicksSurviveAFilterChange()
{
    // The selection a Check column holds is keyed by the underlying row, not
    // by the row on screen. A reader ticks some rows, narrows the filter, and
    // expects the ticks to still be there when it widens again — which is the
    // whole reason the model owns the answer rather than the view.
    KvitUi::BenchmarkTableModel model;
    model.setTotalRows(5000);
    const int boxColumn = 12;
    QCOMPARE(model.columnKindName(boxColumn), QStringLiteral("Check"));
    // And the column of boxes does not sort: its header is the
    // select-displayed control instead.
    QVERIFY(!model.columnSortable(boxColumn));
    QVERIFY(model.columnSortable(0));

    model.setChecked(3, true);
    model.setChecked(9, true);
    QVERIFY(model.cellChecked(3, boxColumn));
    QVERIFY(model.cellChecked(9, boxColumn));
    QVERIFY(!model.cellChecked(4, boxColumn));
    QVERIFY(model.someShownChecked());
    QVERIFY(!model.allShownChecked());

    model.setFilter(QStringLiteral("Harlow"));
    const int matched = model.rowCount();
    QVERIFY(matched > 2);
    model.setEveryShownChecked(true);
    QVERIFY(model.allShownChecked());
    QVERIFY(model.cellChecked(0, boxColumn));

    // Widening the filter brings back both what was ticked before it and what
    // was ticked under it.
    model.setFilter(QString());
    QVERIFY(model.cellChecked(3, boxColumn));
    QVERIFY(model.cellChecked(9, boxColumn));
    QVERIFY(model.someShownChecked());
    QVERIFY(!model.allShownChecked());
}

void TestTableModel::testTheWholeValueIsAnsweredWhereTheColumnShortensIt()
{
    // A column too narrow for its value has to have somewhere to send the
    // reader for the rest of it, and FullTextRole is that somewhere: the view
    // binds it into KvitCell's `fullValue`, which is both what the cell
    // discloses under the pointer and the keyboard cursor and what a screen
    // reader is told the cell holds.
    //
    // The role reaches nothing unless a model answers it, and a model that
    // invents a role name of its own instead — `payeeFull` rather than
    // `fullText` — leaves the two halves of the feature built and never
    // joined, with nothing failing anywhere. So the one model this library
    // ships answers it, and this is what holds the two names together.
    KvitUi::BenchmarkTableModel model;
    model.setTotalRows(200);
    const int payeeColumn = 3;

    QCOMPARE(model.data(model.index(0, payeeColumn)).toString(),
             QStringLiteral("Ashford & Co"));
    QCOMPARE(model.data(model.index(0, payeeColumn),
                        KvitUi::TableModelBase::FullTextRole).toString(),
             QStringLiteral("Ashford & Co (Holdings) Limited"));

    // Only the column that shortens its value answers, in the same way only
    // the Marks column answers with marks: what a cell holds belongs to its
    // column.
    QVERIFY(!model.data(model.index(0, 2),
                        KvitUi::TableModelBase::FullTextRole).isValid());

    // And the name the view binds it through is the one TableModelBase
    // publishes, which is the half a second name would break.
    QCOMPARE(model.roleNames().value(KvitUi::TableModelBase::FullTextRole),
             QByteArray("fullText"));
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
    QList<qint64> pageCost;
    pageCost.reserve(PagesScrolled);
    timer.start();
    for (int i = 0; i < PagesScrolled; ++i) {
        const qint64 from = timer.nsecsElapsed();
        view->setProperty("contentY", page * i);
        // forceLayout is the delegate work: bind, position and reuse every
        // cell the new viewport needs. It is synchronous, which is what makes
        // it measurable.
        QMetaObject::invokeMethod(view, "forceLayout");
        pageCost.append(timer.nsecsElapsed() - from);
    }
    const qint64 elapsed = timer.elapsed();

    QList<qint64> sorted = pageCost;
    std::sort(sorted.begin(), sorted.end());
    const qreal perPage = qreal(sorted.at(sorted.size() / 2)) / 1'000'000.0;
    const qreal worstPage = qreal(sorted.last()) / 1'000'000.0;

    auto *pool = view->property("contentItem").value<QQuickItem *>();
    QVERIFY(pool);
    const int live = pool->childItems().size();
    qInfo("scrolling 250,000 rows: %.2f ms for the middle page (budget %d ms), "
          "%.2f ms for the worst, %lld ms for all %d, "
          "%d live delegate items afterwards",
          perPage, ScrollBudgetMsPerPage, worstPage,
          static_cast<long long>(elapsed), PagesScrolled, live);

    QVERIFY2(perPage < ScrollBudgetMsPerPage,
             qPrintable(QStringLiteral("relaying out a page of a 250,000-row "
                                       "table cost %1 ms against a %2 ms "
                                       "budget (the middle of %3 pages; the "
                                       "worst cost %4 ms and all of them %5 ms)")
                            .arg(perPage, 0, 'f', 2)
                            .arg(ScrollBudgetMsPerPage)
                            .arg(PagesScrolled)
                            .arg(worstPage, 0, 'f', 2)
                            .arg(elapsed)));

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

namespace {

// The KvitEmptyState inside a KvitTable. It has no id a test can reach and no
// object name, so it is found by type: a QML-defined type's class name is its
// file name with a suffix on it.
QQuickItem *findEmptyState(QQuickItem *item)
{
    const auto children = item->childItems();
    for (QQuickItem *child : children) {
        if (QByteArray(child->metaObject()->className()).startsWith("KvitEmptyState"))
            return child;
        if (QQuickItem *found = findEmptyState(child))
            return found;
    }
    return nullptr;
}

// The first item of a given QML type under an item. A QML-defined type's
// class name is its file name with a generated suffix on it, which is why
// this matches on a prefix.
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

// The cell delegate at one row and column. A TableCell publishes both as
// required properties, which is what tells them apart: searching by type
// would find whichever delegate the view happens to have built first.
QQuickItem *findCell(QQuickItem *item, int row, int column)
{
    const auto children = item->childItems();
    for (QQuickItem *child : children) {
        const QVariant childRow = child->property("row");
        const QVariant childColumn = child->property("column");
        if (childRow.isValid() && childColumn.isValid()
            && childRow.toInt() == row && childColumn.toInt() == column) {
            return child;
        }
        if (QQuickItem *found = findCell(child, row, column))
            return found;
    }
    return nullptr;
}

// A QML-declared signal, as the QMetaMethod QSignalSpy wants. Naming it
// through SIGNAL() would mean writing the parameter list a QML signal
// declaration does not spell the same way.
QMetaMethod signalNamed(QObject *object, const char *name)
{
    const QMetaObject *meta = object->metaObject();
    for (int i = 0; i < meta->methodCount(); ++i) {
        const QMetaMethod method = meta->method(i);
        if (method.methodType() == QMetaMethod::Signal
            && method.name() == QByteArray(name)) {
            return method;
        }
    }
    return {};
}

}   // namespace

void TestTableModel::testTheEmptyStateFollowsTheRowCount()
{
    // The overlay that says "nothing here yet" has to disappear when rows
    // arrive. Written as a binding to `model.rowCount()` it does not: a
    // function call is evaluated when the binding is set up and nothing tells
    // QML to evaluate it again, so the table shows the message over a full
    // grid for as long as the view lives. That is what this case holds, and
    // it is worth a windowed test rather than a reading of the property
    // because what the reader sees is the item's visibility.
    if (!m_haveWindow)
        QSKIP("no scene graph on this platform; the view half cannot run");

    QQmlEngine engine;
    QQmlComponent component(&engine);
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    width: 800; height: 400\n"
        "    property alias table: table\n"
        "    property alias model: model\n"
        "    BenchmarkTableModel { id: model; totalRows: 400 }\n"
        "    KvitTable { id: table; anchors.fill: parent; model: model }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/emptystate.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    QQuickWindow window;
    auto *content = qobject_cast<QQuickItem *>(holder.data());
    QVERIFY(content);
    content->setParentItem(window.contentItem());
    window.resize(800, 400);
    window.show();
    QVERIFY(QTest::qWaitForWindowExposed(&window));

    auto *table = content->property("table").value<QQuickItem *>();
    QVERIFY(table);
    auto *model = content->property("model").value<QObject *>();
    QVERIFY(model);
    QQuickItem *empty = findEmptyState(table);
    QVERIFY2(empty, "KvitTable has no KvitEmptyState in it");

    // Four hundred rows, so no message.
    QTRY_COMPARE(table->property("rowCount").toInt(), 400);
    QVERIFY(!empty->isVisible());

    // A filter nothing matches, so the message.
    QVERIFY(model->setProperty("filter", QStringLiteral("zzzzz")));
    QTRY_COMPARE(table->property("rowCount").toInt(), 0);
    QTRY_VERIFY(empty->isVisible());

    // Clearing it puts the rows back, and this is the assertion the bug
    // failed: the count changes underneath a view that was built empty.
    QVERIFY(model->setProperty("filter", QString()));
    QTRY_COMPARE(table->property("rowCount").toInt(), 400);
    QTRY_VERIFY(!empty->isVisible());

    // And a filter that matches some of the rows leaves the message off while
    // the count follows the model down.
    QVERIFY(model->setProperty("filter", QStringLiteral("Harlow")));
    const int matched = model->property("matchedRows").toInt();
    QVERIFY(matched > 0 && matched < 400);
    QTRY_COMPARE(table->property("rowCount").toInt(), matched);
    QVERIFY(!empty->isVisible());
}

void TestTableModel::testTickingABoxDoesNotAlsoOpenTheRow()
{
    // A press on the box in a Check column is a request to tick that row and
    // nothing else.
    //
    // Every cell in the table carries the TapHandler that opens a record, and
    // a checkbox is a child of one of those cells. A TapHandler at its default
    // gesture policy takes a passive grab rather than an exclusive one, so it
    // still fires when the control underneath it has accepted the press --
    // which means one press on the box would tick the row and open it at the
    // same time, and the reader who ticked forty rows would have opened forty
    // records on the way.
    if (!m_haveWindow)
        QSKIP("no scene graph on this platform; the view half cannot run");

    QQmlEngine engine;
    QQmlComponent component(&engine);
    // Wide enough for all thirteen columns, so the box column is on screen
    // without the test having to scroll to it: the widths in
    // BenchmarkTableModel::columnWidth add up to 1,526.
    component.setData(
        "import QtQuick\n"
        "import Kvit.Ui\n"
        "Item {\n"
        "    width: 1600; height: 400\n"
        "    property alias table: table\n"
        "    property alias model: model\n"
        "    BenchmarkTableModel { id: model; totalRows: 200 }\n"
        "    KvitTable { id: table; anchors.fill: parent; model: model }\n"
        "}\n",
        QUrl(QStringLiteral("qrc:/test/checkpress.qml")));
    QVERIFY2(component.isReady(), qPrintable(component.errorString()));

    QScopedPointer<QObject> holder(component.create());
    QVERIFY2(!holder.isNull(), qPrintable(component.errorString()));

    QQuickWindow window;
    auto *content = qobject_cast<QQuickItem *>(holder.data());
    QVERIFY(content);
    content->setParentItem(window.contentItem());
    window.resize(1600, 400);
    window.show();
    QVERIFY(QTest::qWaitForWindowExposed(&window));

    auto *table = content->property("table").value<QQuickItem *>();
    QVERIFY(table);
    auto *view = table->property("view").value<QQuickItem *>();
    QVERIFY(view);
    // One layout, so the first page's delegates exist to be found and pressed.
    QMetaObject::invokeMethod(view, "forceLayout");

    // Searched under the view rather than under the table, so the header's
    // own select-displayed box is not what gets pressed.
    QQuickItem *box = findByType(view, "KvitCheck");
    QVERIFY2(box, "no checkbox in the table's Check column");

    QSignalSpy toggled(table, signalNamed(table, "cellToggled"));
    QSignalSpy pressed(table, signalNamed(table, "rowPressed"));
    QSignalSpy activated(table, signalNamed(table, "rowActivated"));
    QVERIFY(toggled.isValid());
    QVERIFY(pressed.isValid());
    QVERIFY(activated.isValid());

    const QPointF centre =
        box->mapToScene(QPointF(box->width() / 2, box->height() / 2));
    QTest::mouseClick(&window, Qt::LeftButton, Qt::NoModifier, centre.toPoint());

    QTRY_COMPARE(toggled.count(), 1);
    QCOMPARE(toggled.at(0).at(2).toBool(), true);
    QCOMPARE(pressed.count(), 0);
    QCOMPARE(activated.count(), 0);

    // And a press on the first cell of the same row still opens it, so what
    // is being checked above is the box rather than a table that stopped
    // answering the pointer.
    QQuickItem *first = findCell(view, 0, 0);
    QVERIFY2(first, "the table has no delegate for row 0, column 0");
    const QPointF elsewhere =
        first->mapToScene(QPointF(first->width() / 2, first->height() / 2));
    QTest::mouseClick(&window, Qt::LeftButton, Qt::NoModifier,
                      elsewhere.toPoint());
    QTRY_COMPARE(pressed.count(), 1);
    QCOMPARE(pressed.at(0).at(0).toInt(), 0);
    QCOMPARE(toggled.count(), 1);
}

QTEST_MAIN(TestTableModel)
#include "test_tablemodel.moc"
