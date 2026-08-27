// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "tablemodelbase.h"

namespace KvitUi {

TableModelBase::TableModelBase(QObject *parent)
    : QAbstractTableModel(parent)
{
}

QHash<int, QByteArray> TableModelBase::roleNames() const
{
    return {
        { DisplayRole, "display" },
        { MeasuredRole, "measured" },
        { SortRole, "sortValue" },
        { MarkRole, "mark" },
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

}   // namespace KvitUi
