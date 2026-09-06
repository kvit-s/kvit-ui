// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "benchmarktablemodel.h"

#include <QDate>
#include <QLocale>
#include <QVariantList>
#include <QVariantMap>

namespace {

// Twelve columns of value, which is the width prd.md's benchmark asks for and
// roughly what a transaction browser shows: when, what, where, how much, and
// the several small marks that say what kind of thing it was. The thirteenth
// holds no value at all — it is the box the reader ticks, and it is here
// because a ledger with bulk actions on it is the surface this table was
// built for.
const char *kTitles[] = {
    "Date", "Reference", "Description", "Payee", "Category", "Account",
    "Amount", "Balance", "Currency", "Status", "Tags", "Note", "",
};

// A few short word lists, so a description reads like text rather than like a
// number. The benchmark's filter is a substring search over these, which is
// the operation a reader actually performs.
const char *kVerbs[] = { "Payment", "Transfer", "Refund", "Deposit",
                         "Withdrawal", "Fee", "Interest", "Adjustment" };
const char *kPlaces[] = { "Ashford", "Barrow", "Colwyn", "Dunmore", "Elmsley",
                          "Fairholt", "Grantham", "Harlow", "Ivybridge",
                          "Jarrow", "Kelmscott", "Lyndhurst" };
const char *kCategories[] = { "Groceries", "Transport", "Utilities", "Rent",
                              "Leisure", "Health", "Savings", "Income" };
const char *kAccounts[] = { "Current", "Savings", "Card", "Joint" };
const char *kStatuses[] = { "Settled", "Pending", "Disputed" };

template <typename T, int N>
constexpr int countOf(T (&)[N]) { return N; }

// The status column, named because it is now referred to twice: once for what
// kind of cell draws it and once for the marks it answers with.
constexpr int kStatusColumn = 9;
// The amount column, which formats its own value; see the note where it does.
constexpr int kAmountColumn = 6;
// The column of boxes. Last rather than first only because every index in
// tests/test_tablemodel.cpp names a column by number, and a ledger that put
// the boxes on the left would renumber all twelve.
constexpr int kSelectColumn = 12;
// The payee column, which is the one that answers with a shortened value and
// the whole of it separately; see the note where it does.
constexpr int kPayeeColumn = 3;

// The amount a row is worth, as a number, kept in one place because the
// column answers with it twice: formatted for the reader, and unformatted for
// the sort.
double amountOf(int row)
{
    return (row * 37 % 100000) / 100.0 - 500.0;
}

// An amount as the string the cell draws.
//
// The formatting is the model's, which is the whole of the Money cell's
// contract: the cell is given a string and draws it. Two things about money
// are lost if the cell does this instead. An amount is a signed 64-bit count
// of minor units and a JavaScript number carries 53 bits of integer, so a
// value large enough to matter is rounded on the way through and nothing
// says so; and how many minor digits an amount has belongs to its currency,
// none for the yen and three for the dinar, which is knowledge a cell does
// not have. Every row here is in sterling, so two.
QString formattedAmount(int row)
{
    return QLocale::system().toString(amountOf(row), 'f', 2);
}

// One mark, in the shape a Marks cell reads: a theme tone, the shape that
// carries the same distinction without colour, and the word a reader hears on
// hover. Three of the estate's tones and all three of KvitDot's shapes, so
// the high-contrast pass has something to fail on if a shape is dropped.
QVariantMap statusMark(int row)
{
    switch (row % countOf(kStatuses)) {
    case 0:
        return { { QStringLiteral("tone"), QStringLiteral("success") },
                 { QStringLiteral("shape"), QStringLiteral("circle") },
                 { QStringLiteral("label"), QStringLiteral("Settled") } };
    case 1:
        return { { QStringLiteral("tone"), QStringLiteral("warning") },
                 { QStringLiteral("shape"), QStringLiteral("diamond") },
                 { QStringLiteral("label"), QStringLiteral("Pending") } };
    default:
        return { { QStringLiteral("tone"), QStringLiteral("danger") },
                 { QStringLiteral("shape"), QStringLiteral("square") },
                 { QStringLiteral("label"), QStringLiteral("Disputed") } };
    }
}

}   // namespace

namespace KvitUi {

BenchmarkTableModel::BenchmarkTableModel(QObject *parent)
    : TableModelBase(parent)
{
}

int BenchmarkTableModel::rowCount(const QModelIndex &parent) const
{
    if (parent.isValid())
        return 0;
    return m_filter.isEmpty() ? m_totalRows : int(m_matched.size());
}

int BenchmarkTableModel::columnCount(const QModelIndex &parent) const
{
    return parent.isValid() ? 0 : ColumnCount;
}

int BenchmarkTableModel::sourceRow(int row) const
{
    if (m_filter.isEmpty())
        return row;
    return row >= 0 && row < m_matched.size() ? m_matched.at(row) : -1;
}

QString BenchmarkTableModel::descriptionFor(int row) const
{
    return QStringLiteral("%1 to %2 depot")
        .arg(QLatin1String(kVerbs[row % countOf(kVerbs)]))
        .arg(QLatin1String(kPlaces[(row / 7) % countOf(kPlaces)]));
}

QVariant BenchmarkTableModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid())
        return {};
    const int row = sourceRow(index.row());
    if (row < 0 || row >= m_totalRows)
        return {};
    const int column = index.column();

    if (role == MeasuredRole) {
        // Every fortieth balance is deliberately unmeasured, so the benchmark
        // exercises the em-dash path rather than only the numeric one. This is
        // the distinction KvitFigure exists for: a balance nobody has computed
        // is not a balance of zero.
        return !(column == 7 && row % 40 == 0);
    }

    if (role == MarksRole) {
        // Only the status column draws marks: the kind belongs to the column,
        // so a model answering marks everywhere would be describing cells no
        // view asked about.
        if (column != kStatusColumn)
            return {};
        QVariantList marks;
        marks.append(statusMark(row));
        // Every seventh row is also unreviewed, which is the whole reason
        // this column is Marks rather than Chip: a row is regularly two
        // things at once and one chip can only say one of them.
        if (row % 7 == 0) {
            marks.append(QVariantMap {
                { QStringLiteral("tone"), QStringLiteral("neutral") },
                { QStringLiteral("shape"), QStringLiteral("square") },
                { QStringLiteral("label"), QStringLiteral("Not reviewed") },
            });
        }
        return marks;
    }

    if (role == CheckedRole) {
        // Only the column of boxes answers, for the same reason only the
        // status column answers with marks: the kind belongs to the column.
        if (column != kSelectColumn)
            return {};
        return m_checked.contains(row);
    }

    if (role == FullTextRole) {
        // The whole name where the column shows a short form of it.
        //
        // This is the shape the role exists for, and the reason the only
        // model this library ships answers it: a view binds `fullText` into
        // KvitCell's `fullValue`, and a model that answers some name of its
        // own instead leaves the disclosure and the accessible name reading
        // the shortened value with nothing anywhere saying so. A trading name
        // in the column and the registered name behind it is the ordinary
        // case in a ledger, and it is what the column is too narrow to show.
        if (column != kPayeeColumn)
            return {};
        return QStringLiteral("%1 & Co (Holdings) Limited")
            .arg(QLatin1String(kPlaces[row % countOf(kPlaces)]));
    }

    if (role == SortRole) {
        switch (column) {
        case 0: return QDate(2020, 1, 1).addDays(row % 2200);
        // The number rather than the string the reader sees, so the column
        // sorts by value and not by the digits of its formatting.
        case kAmountColumn: return amountOf(row);
        case 7: return (row * 977 % 4000000) / 100.0;
        default: break;
        }
        role = DisplayRole;
    }

    if (role != DisplayRole && role != Qt::EditRole)
        return {};

    const auto edited = m_edits.constFind(row * ColumnCount + column);
    if (edited != m_edits.constEnd())
        return edited.value();

    switch (column) {
    case 0:  return QDate(2020, 1, 1).addDays(row % 2200);
    case 1:  return QStringLiteral("TX-%1").arg(row, 8, 10, QLatin1Char('0'));
    case 2:  return descriptionFor(row);
    case kPayeeColumn: return QStringLiteral("%1 & Co")
                        .arg(QLatin1String(kPlaces[row % countOf(kPlaces)]));
    case 4:  return QLatin1String(kCategories[row % countOf(kCategories)]);
    case 5:  return QLatin1String(kAccounts[row % countOf(kAccounts)]);
    case kAmountColumn: return formattedAmount(row);
    case 7:  return (row % 40 == 0) ? QVariant()
                                    : QVariant((row * 977 % 4000000) / 100.0);
    case 8:  return QStringLiteral("GBP");
    case 9:  return QLatin1String(kStatuses[row % countOf(kStatuses)]);
    case 10: return QStringLiteral("#%1")
                        .arg(QString::fromLatin1(kCategories[(row / 3)
                                 % countOf(kCategories)]).toLower());
    case 11: return row % 11 == 0 ? QStringLiteral("Checked against statement")
                                  : QString();
    // The box draws itself from CheckedRole and has no text of its own. An
    // empty string rather than nothing, so a copy of the row keeps its
    // columns lined up.
    case kSelectColumn: return QString();
    default: return {};
    }
}

QVariant BenchmarkTableModel::headerData(int section, Qt::Orientation orientation,
                                         int role) const
{
    if (role != Qt::DisplayRole || orientation != Qt::Horizontal)
        return {};
    if (section < 0 || section >= ColumnCount)
        return {};
    return QLatin1String(kTitles[section]);
}

Qt::ItemFlags BenchmarkTableModel::flags(const QModelIndex &index) const
{
    if (!index.isValid())
        return Qt::NoItemFlags;
    Qt::ItemFlags result = Qt::ItemIsEnabled | Qt::ItemIsSelectable;
    // The two free-text columns are editable, which is what the benchmark's
    // inline-editing requirement needs; a generated reference or balance is
    // not something a reader gets to overwrite.
    if (index.column() == 2 || index.column() == 11)
        result |= Qt::ItemIsEditable;
    return result;
}

bool BenchmarkTableModel::setData(const QModelIndex &index, const QVariant &value,
                                  int role)
{
    if (!index.isValid() || (role != Qt::EditRole && role != DisplayRole))
        return false;
    if (!(flags(index) & Qt::ItemIsEditable))
        return false;
    const int row = sourceRow(index.row());
    if (row < 0)
        return false;
    m_edits.insert(row * ColumnCount + index.column(), value.toString());
    emit dataChanged(index, index, { DisplayRole, Qt::EditRole });
    return true;
}

int BenchmarkTableModel::columnWidth(int column) const
{
    switch (column) {
    case 0:  return 96;    // Date
    case 1:  return 110;   // Reference
    case 2:  return 220;   // Description
    case 3:  return 150;   // Payee
    case 6:  return 96;    // Amount
    case 7:  return 110;   // Balance
    case 8:  return 64;    // Currency
    case 11: return 200;   // Note
    case kSelectColumn: return 40;
    default: return 110;
    }
}

TableModelBase::CellKind BenchmarkTableModel::columnKind(int column) const
{
    switch (column) {
    case 0:  return Date;
    case 1:  return Slug;
    case 4:  return Chip;
    case 6:  return Money;
    case 7:  return Figure;
    // The status column carries marks rather than one chip, because a
    // transaction is regularly settled *and* unreviewed and a single chip
    // has to drop one of the two.
    case kStatusColumn: return Marks;
    case 10: return Chip;
    case kSelectColumn: return Check;
    default: return Text;
    }
}

void BenchmarkTableModel::setTotalRows(int rows)
{
    rows = qMax(0, rows);
    if (m_totalRows == rows)
        return;
    beginResetModel();
    m_totalRows = rows;
    m_edits.clear();
    m_matched.clear();
    m_checked.clear();
    m_checkedShown = 0;
    if (!m_filter.isEmpty())
        setFilter(m_filter);
    endResetModel();
    emit totalRowsChanged();
    emit checkedChanged();
}

bool BenchmarkTableModel::columnSortable(int column) const
{
    // Every column but the boxes. Sorting a column of ticks would order the
    // rows by an answer the reader gave rather than by anything about them,
    // and the header of that column is the select-displayed control instead.
    return column != kSelectColumn;
}

void BenchmarkTableModel::setChecked(int row, bool checked)
{
    const int source = sourceRow(row);
    if (source < 0)
        return;
    if (checked == m_checked.contains(source))
        return;
    if (checked)
        m_checked.insert(source);
    else
        m_checked.remove(source);
    m_checkedShown += checked ? 1 : -1;
    const QModelIndex box = index(row, kSelectColumn);
    emit dataChanged(box, box, { CheckedRole });
    emit checkedChanged();
}

void BenchmarkTableModel::setEveryShownChecked(bool checked)
{
    const int rows = rowCount();
    if (rows == 0)
        return;
    for (int row = 0; row < rows; ++row) {
        const int source = sourceRow(row);
        if (checked)
            m_checked.insert(source);
        else
            m_checked.remove(source);
    }
    m_checkedShown = checked ? rows : 0;
    emit dataChanged(index(0, kSelectColumn), index(rows - 1, kSelectColumn),
                     { CheckedRole });
    emit checkedChanged();
}

bool BenchmarkTableModel::allShownChecked() const
{
    return rowCount() > 0 && m_checkedShown == rowCount();
}

bool BenchmarkTableModel::someShownChecked() const
{
    return m_checkedShown > 0 && !allShownChecked();
}

void BenchmarkTableModel::recountCheckedShown()
{
    if (m_checked.isEmpty()) {
        m_checkedShown = 0;
        return;
    }
    int shown = 0;
    const int rows = rowCount();
    for (int row = 0; row < rows; ++row) {
        if (m_checked.contains(sourceRow(row)))
            ++shown;
    }
    m_checkedShown = shown;
}

void BenchmarkTableModel::setFilter(const QString &text)
{
    beginResetModel();
    m_filter = text;
    m_matched.clear();
    if (!m_filter.isEmpty()) {
        // A linear scan of the generated descriptions, which is what makes
        // this a measurement of the view rather than of an index: 250,000
        // substring tests is the honest cost of filtering an unindexed table,
        // and if the view cannot stay under 100 ms on top of that, an index in
        // the real model would not save it.
        m_matched.reserve(m_totalRows / 8);
        for (int row = 0; row < m_totalRows; ++row) {
            if (descriptionFor(row).contains(m_filter, Qt::CaseInsensitive))
                m_matched.append(row);
        }
    }
    endResetModel();
    // The ticks stay where they were; how many of them are on screen does
    // not, and this is the one place it changes without a tick being given.
    recountCheckedShown();
    emit filterChanged();
    emit checkedChanged();
}

}   // namespace KvitUi
