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
        Marks,    // several small marks at once, drawn without words
        Check,    // a box the reader ticks
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
        // For a Marks cell, the several marks a row carries at once, as a
        // list of { tone, shape, label }. One mark cannot say that a
        // transaction is both a split and unreviewed, and a row that is both
        // is the ordinary case rather than the exception.
        MarksRole,
        // For a Check cell, whether this row is ticked. The model owns the
        // answer because a selection outlives the rows on screen: a reader
        // ticks forty rows, changes the filter, and expects the forty to
        // still be ticked when the filter comes back.
        CheckedRole,
        // What a Figure or Money cell draws beside the amount: a currency
        // code, or a unit of measure. It is the cell's own unit rather than
        // part of the string, so the amount stays alignable down its decimal
        // separator and a model never has to paste a currency code onto a
        // formatted number.
        UnitRole,
        // The whole value where the cell shows a shortened one. A column
        // narrow enough to elide its text has no other way to give the reader
        // the rest of it, and a name that is only ever seen truncated is a
        // name the reader cannot check.
        FullTextRole,
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

    // The same answer as columnKind(), as the enumeration key: "Text",
    // "Figure", "Chip", "Slug", "Date", "Money", "Marks", "Check".
    //
    // QML needs the word rather than the number. A Q_ENUM returned through a
    // Q_INVOKABLE arrives in JavaScript as an integer, so String() on it
    // yields "1" and every cell falls through to KvitCell's Text default —
    // which is what every table in the estate drew until this existed. The
    // string is read out of the Q_ENUM's own metadata rather than written as
    // a switch here, so a kind added to the enumeration cannot go missing
    // from the name.
    //
    // Not virtual: it is derived from columnKind(), so a model overrides that
    // one and gets this for nothing.
    Q_INVOKABLE QString columnKindName(int column) const;

    // What a screen reader says about a whole row, rather than about the one
    // cell the keyboard happens to be in. A table of a dozen columns read
    // cell by cell tells the reader nothing about which record they are on,
    // so every cell carries this as its description.
    //
    // The default joins what the row's text-bearing columns hold, in column
    // order, which for a ledger is the date, the payee and the amount. Marks
    // and checkboxes are left out: they name themselves, and a row's states
    // are not how the reader identifies it. A model with a better sentence
    // overrides this.
    Q_INVOKABLE virtual QString rowName(int row) const;

    // Whether the box in a Check cell is ticked. The view needs this to
    // answer the space bar on the row the keyboard is on, and reading it
    // through the role from QML would mean naming the role number there.
    Q_INVOKABLE bool cellChecked(int row, int column) const;

private:
    // The last row rowName() was asked for, and the sentence it answered.
    //
    // Every cell of a row asks for the same sentence, because the description
    // belongs to the row and a delegate exists per cell. Answering it afresh
    // each time reads every column of the row, so a thirteen-column table
    // walked one page cost thirteen row reads per row instead of one — about
    // eight thousand data() calls to relayout a page, and a third of the
    // page's whole cost. One row of memory removes twelve of those thirteen.
    //
    // It is dropped whenever anything could have changed what a row holds:
    // the signals connected in the constructor are every one this class can
    // see a subclass emit. A subclass that changes its data without emitting
    // one of them was already breaking every view of it.
    void forgetNamedRow();

    mutable int m_namedRow = -1;
    mutable QString m_namedRowText;
};

}   // namespace KvitUi

#endif // KVIT_UI_TABLEMODELBASE_H
