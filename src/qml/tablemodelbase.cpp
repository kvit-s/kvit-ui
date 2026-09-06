// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "tablemodelbase.h"

#include <QMetaEnum>

namespace KvitUi {

TableModelBase::TableModelBase(QObject *parent)
    : QAbstractTableModel(parent)
{
    // rowName() keeps the sentence it built for the last row it was asked
    // about. Anything that could change what a row holds drops it.
    connect(this, &QAbstractItemModel::dataChanged, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::modelReset, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::layoutChanged, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::rowsInserted, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::rowsRemoved, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::rowsMoved, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::columnsInserted, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::columnsRemoved, this,
            &TableModelBase::forgetNamedRow);
    connect(this, &QAbstractItemModel::columnsMoved, this,
            &TableModelBase::forgetNamedRow);
}

void TableModelBase::forgetNamedRow()
{
    m_namedRow = -1;
    m_namedRowText.clear();
}

QHash<int, QByteArray> TableModelBase::roleNames() const
{
    return {
        { DisplayRole, "display" },
        { MeasuredRole, "measured" },
        { SortRole, "sortValue" },
        { MarkRole, "mark" },
        { MarksRole, "marks" },
        { CheckedRole, "checked" },
        { UnitRole, "unit" },
        { FullTextRole, "fullText" },
    };
}

QString TableModelBase::columnTitle(int column) const
{
    return headerData(column, Qt::Horizontal, Qt::DisplayRole).toString();
}

int TableModelBase::columnWidth(int) const
{
    // A width in design pixels, which the view runs through Interface.px().
    // 120 is what a name or a date needs at the default interface size; a
    // model that knows better says so.
    return 120;
}

TableModelBase::CellKind TableModelBase::columnKind(int) const
{
    return Text;
}

bool TableModelBase::columnSortable(int) const
{
    return true;
}

QString TableModelBase::columnKindName(int column) const
{
    // The keys, converted once. Every cell of every page asks for this — a
    // page of a thirteen-column table is three hundred calls — and building
    // the string each time allocates three hundred times over to say the same
    // eight words. The names still come out of the enumeration's own
    // metadata, so a kind added to it cannot go missing from the list.
    static const QHash<int, QString> names = [] {
        const QMetaEnum kinds = QMetaEnum::fromType<CellKind>();
        QHash<int, QString> built;
        built.reserve(kinds.keyCount());
        for (int i = 0; i < kinds.keyCount(); ++i)
            built.insert(kinds.value(i), QString::fromLatin1(kinds.key(i)));
        return built;
    }();
    // A model that returned a value outside the enumeration gets the default
    // rather than an empty string, because an empty kind draws nothing at all
    // and a column of nothing is harder to notice than a column of text.
    return names.value(int(columnKind(column)), QStringLiteral("Text"));
}

QString TableModelBase::rowName(int row) const
{
    if (row < 0 || row >= rowCount())
        return {};
    // See the note on m_namedRow: every cell of a row asks for this same
    // sentence, and building it reads the whole row.
    if (row == m_namedRow)
        return m_namedRowText;
    QStringList parts;
    const int columns = columnCount();
    for (int column = 0; column < columns; ++column) {
        const CellKind kind = columnKind(column);
        if (kind == Marks || kind == Check)
            continue;
        const QString text = data(index(row, column), DisplayRole).toString();
        if (text.isEmpty())
            continue;
        const QString unit = data(index(row, column), UnitRole).toString();
        parts.append(unit.isEmpty() ? text : text + QLatin1Char(' ') + unit);
    }
    m_namedRow = row;
    m_namedRowText = parts.join(QStringLiteral(", "));
    return m_namedRowText;
}

bool TableModelBase::cellChecked(int row, int column) const
{
    return data(index(row, column), CheckedRole).toBool();
}

}   // namespace KvitUi
