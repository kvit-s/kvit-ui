// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_TABLEMODELBASE_H
#define KVIT_UI_TABLEMODELBASE_H

#include <QAbstractTableModel>
#include <QStringList>
#include <QtQml/qqmlregistration.h>

namespace KvitUi {

// What KvitTable expects of a model, and the small amount of it that is worth
// writing once.
//
// The library ships the view and this base; the models belong to the
// applications, because what is in a row is the application's business.
// kvit-cash's transaction model reads a ledger, kvit-hub's project model reads
// a portfolio, and neither has anything to say to the other.
//
// What is shared is the shape of a column. A table that lets a reader resize,
// reorder, hide and sort columns needs somewhere to put a column's title, its
// preferred width, whether it is sortable and what kind of cell draws it, and
// every application would otherwise invent that description again. The kind is
// what KvitCell switches on, so `Figure` gets tabular numerals and an em dash
// for an unmeasured value while `Text` gets neither.
class TableModelBase : public QAbstractTableModel
{
    Q_OBJECT
    QML_ELEMENT
    QML_UNCREATABLE("TableModelBase is a base class for application models")

public:
    // How a cell in this column is drawn. The names are what a column
    // description says, not what a delegate looks like, so a change to how a
    // figure is drawn does not reach a model.
    enum CellKind {
        Text,     // ordinary text, left aligned
        Figure,   // a measured value: tabular numerals, right aligned
        Chip,     // a small labelled mark
        Slug,     // a monospace identifier
        Date,     // a date, formatted by the view's locale
        Money,    // a figure with a currency, following money-display.md
    };
    Q_ENUM(CellKind)

    // The roles a KvitTable delegate reads. A model may add its own above
    // UserRole; these are the ones the view itself uses.
    enum Role {
        DisplayRole = Qt::DisplayRole,
        // What the cell says when it has no value. Distinguishing "not
        // measured" from "zero" is the rule KvitFigure exists to enforce, and
        // it can only be enforced if the model can say which it means.
        MeasuredRole = Qt::UserRole + 1,
        // The value to sort and filter on, when that differs from what is
        // shown: a date shown as "yesterday", a figure shown rounded.
        SortRole,
        // For a Chip cell, which of the theme's marks to draw it as.
        MarkRole,
    };
    Q_ENUM(Role)

    explicit TableModelBase(QObject *parent = nullptr);

    QHash<int, QByteArray> roleNames() const override;

    // The columns, in model order. A view's own column order, widths and
    // hidden set are the view's state and are saved with the view, not here:
    // two windows onto the same model may show different columns.
    Q_INVOKABLE virtual QString columnTitle(int column) const;
    Q_INVOKABLE virtual int columnWidth(int column) const;
    Q_INVOKABLE virtual CellKind columnKind(int column) const;
    Q_INVOKABLE virtual bool columnSortable(int column) const;
};

}   // namespace KvitUi

#endif // KVIT_UI_TABLEMODELBASE_H
