// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_BENCHMARKTABLEMODEL_H
#define KVIT_UI_BENCHMARKTABLEMODEL_H

#include <QList>
#include <QSet>
#include <QString>
#include <QtQml/qqmlregistration.h>

#include "tablemodelbase.h"

namespace KvitUi {

// A synthetic table with as many rows as anything in the estate will ever
// have, for the one measurement in this plan that can come back the wrong way.
//
// prd.md Decision 3 puts kvit-cash's transaction browser in QML on KvitTable
// rather than in Qt Widgets, and holds the Widgets fallback open until a
// 250,000-row browser is shown to scroll smoothly and filter in under 100 ms
// with twelve configurable columns, multi-select and inline editing. This is
// the model that benchmark runs against, and it is the only model this library
// ships: the real ones belong to the applications.
//
// Rows are generated rather than stored, so the model costs nothing to build
// and the benchmark measures the view rather than the machine's memory
// bandwidth. Every value is derived from the row index by arithmetic, which
// also makes the assertions in the benchmark exact: row 173,404's amount is
// whatever the formula says it is, whatever order the view asked for it in.
class BenchmarkTableModel : public TableModelBase
{
    Q_OBJECT
    QML_ELEMENT

    Q_PROPERTY(int totalRows READ totalRows WRITE setTotalRows NOTIFY totalRowsChanged)
    // The rows the current filter admits. Setting `filter` recomputes it,
    // which is the operation the 100 ms budget is about.
    Q_PROPERTY(QString filter READ filter WRITE setFilter NOTIFY filterChanged)
    Q_PROPERTY(int matchedRows READ matchedRows NOTIFY filterChanged)
    // What the select-displayed box in the header of the column of boxes has
    // to show. Two properties rather than one, because a box has three states
    // and the third one — some of the rows shown, not all — is the one a
    // reader has to be able to tell from the other two.
    Q_PROPERTY(bool allShownChecked READ allShownChecked NOTIFY checkedChanged)
    Q_PROPERTY(bool someShownChecked READ someShownChecked NOTIFY checkedChanged)

public:
    static constexpr int DefaultTotalRows = 250000;
    // Twelve columns of value and one of boxes. The benchmark's claim is
    // about the twelve; the thirteenth is here so the Check kind is exercised
    // by the only model this library ships.
    static constexpr int ValueColumnCount = 12;
    static constexpr int ColumnCount = ValueColumnCount + 1;

    explicit BenchmarkTableModel(QObject *parent = nullptr);

    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    int columnCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QVariant headerData(int section, Qt::Orientation orientation,
                        int role = Qt::DisplayRole) const override;
    Qt::ItemFlags flags(const QModelIndex &index) const override;
    bool setData(const QModelIndex &index, const QVariant &value,
                 int role = Qt::EditRole) override;

    int columnWidth(int column) const override;
    CellKind columnKind(int column) const override;
    bool columnSortable(int column) const override;

    // Ticking one row, and ticking every row the filter is showing.
    //
    // The model owns the answer rather than the view, because a selection
    // outlives the rows on screen: a reader ticks forty rows, narrows the
    // filter, and expects the forty to still be ticked when it widens again.
    // The set is keyed by generated row, so a filter change moves no ticks.
    Q_INVOKABLE void setChecked(int row, bool checked);
    Q_INVOKABLE void setEveryShownChecked(bool checked);
    bool allShownChecked() const;
    bool someShownChecked() const;

    int totalRows() const { return m_totalRows; }
    void setTotalRows(int rows);
    QString filter() const { return m_filter; }
    void setFilter(const QString &text);
    int matchedRows() const { return int(m_matched.size()); }

signals:
    void totalRowsChanged();
    void filterChanged();
    void checkedChanged();

private:
    // The generated row behind a visible row, which is the identity when no
    // filter is set.
    int sourceRow(int row) const;
    QString descriptionFor(int row) const;

    int m_totalRows = DefaultTotalRows;
    QString m_filter;
    QList<int> m_matched;
    // An edit overrides the generated value for one cell. A hash rather than a
    // copy of the table, because the benchmark edits a handful of cells and
    // materialising 250,000 rows to allow it would be measuring the wrong
    // thing.
    QHash<int, QString> m_edits;
    // Ticked rows, by generated row rather than by visible row.
    QSet<int> m_checked;
    // How many of the rows the current filter shows are ticked. Counted when
    // the filter changes and adjusted by one on each tick, rather than
    // recomputed whenever the header asks: the header asks on every repaint,
    // and walking 250,000 rows there is the shape that makes a table feel
    // slow for a reason nobody can find.
    int m_checkedShown = 0;

    void recountCheckedShown();
};

}   // namespace KvitUi

#endif // KVIT_UI_BENCHMARKTABLEMODEL_H
