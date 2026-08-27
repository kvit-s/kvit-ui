// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "benchmarktablemodel.h"

#include <QDate>

namespace {

// Twelve columns, which is the width prd.md's benchmark asks for and roughly
// what a transaction browser shows: when, what, where, how much, and the
// several small marks that say what kind of thing it was.
const char *kTitles[] = {
    "Date", "Reference", "Description", "Payee", "Category", "Account",
    "Amount", "Balance", "Currency", "Status", "Tags", "Note",
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

    if (role == SortRole) {
        switch (column) {
        case 0: return QDate(2020, 1, 1).addDays(row % 2200);
        case 6: return (row * 37 % 100000) / 100.0 - 500.0;
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
    case 3:  return QStringLiteral("%1 & Co")
                        .arg(QLatin1String(kPlaces[row % countOf(kPlaces)]));
    case 4:  return QLatin1String(kCategories[row % countOf(kCategories)]);
    case 5:  return QLatin1String(kAccounts[row % countOf(kAccounts)]);
    case 6:  return (row * 37 % 100000) / 100.0 - 500.0;
    case 7:  return (row % 40 == 0) ? QVariant()
                                    : QVariant((row * 977 % 4000000) / 100.0);
    case 8:  return QStringLiteral("GBP");
    case 9:  return QLatin1String(kStatuses[row % countOf(kStatuses)]);
    case 10: return QStringLiteral("#%1")
                        .arg(QString::fromLatin1(kCategories[(row / 3)
                                 % countOf(kCategories)]).toLower());
    case 11: return row % 11 == 0 ? QStringLiteral("Checked against statement")
                                  : QString();
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
    case 9:  return Chip;
    case 10: return Chip;
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
    if (!m_filter.isEmpty())
        setFilter(m_filter);
    endResetModel();
    emit totalRowsChanged();
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
    emit filterChanged();
}

}   // namespace KvitUi
