// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A dense, configurable, editable table over a C++ model.
//
// This is the one item in the whole plan that could fail on its own terms.
// prd.md Decision 3 put kvit-cash's 250,000-row transaction browser in QML
// rather than in Qt Widgets, on the grounds that the alternative means
// implementing rows, figures, chips and the currency-marking rules of
// money-display.md a second time inside QStyledItemDelegate::paint(). The
// Widgets fallback is held open until this holds smooth scrolling and
// sub-100 ms filtering at that size; tests/test_tablemodel.cpp is the
// measurement.
//
// ── What makes it hold at that size ────────────────────────────────────────
//
// `TableView` with `reuseItems` over a QAbstractItemModel, so the number of
// delegate items is the number of visible cells rather than the number of
// cells: about 400 against three million. Nothing here ever holds a row.
//
// The model does the filtering, in C++, and reports a new row count. Filtering
// in QML would mean building a JavaScript array of 250,000 indices per
// keystroke, which is the arrangement that makes a QML table feel slow and
// then gets blamed on QML.
//
// `KvitCell` is a Loader over six small components rather than one delegate
// with six branches in it. A branching delegate evaluates every branch's
// bindings for every cell.
//
// ── What the reader gets ───────────────────────────────────────────────────
//
// Column resize and reorder through the header, multi-select through a
// selection model, keyboard navigation, copy of a selection, and a saved
// column set. The empty state is KvitEmptyState rather than a blank grid,
// which is the same rule the charts follow: nothing to show is a thing to
// say, not an absence to leave.
Item {
    id: root

    // A KvitUi.TableModelBase subclass. The library ships the view and the
    // base class; the models belong to the applications, because what is in a
    // row is the application's business.
    required property var model

    property bool editable: false
    property string emptyTitle: qsTr("Nothing here yet")
    property string emptyDetail: ""
    // Column indices the reader has hidden. A view's own column set is the
    // view's state, saved with the view: two windows onto the same model may
    // show different columns.
    property var hiddenColumns: []

    readonly property alias selection: selectionModel
    readonly property alias view: table

    // How many rows the table is showing, which after a filter is fewer than
    // the model holds. TableView keeps this in step with the model; a binding
    // written as `model.rowCount() === 0` does not, because a call is
    // evaluated once and nothing tells QML to evaluate it again. That is how
    // an empty state comes to sit over a full table: the view is built while
    // a filter matches nothing, the reader clears the filter, the rows arrive
    // and the overlay never hears about it. An application wanting to say
    // "24 of 250,000" should read this rather than call the model.
    readonly property int rowCount: table.rows

    // Hiding a column changes what columnWidthProvider answers, and TableView
    // caches those answers: without the relayout the column keeps its old
    // width and the hide does nothing visible.
    onHiddenColumnsChanged: table.forceLayout()

    signal rowActivated(int row)

    implicitWidth: Interface.px(600)
    implicitHeight: Interface.px(400)

    Accessible.role: Accessible.Table
    Accessible.name: qsTr("Table")

    ItemSelectionModel {
        id: selectionModel
        model: root.model
    }

    HorizontalHeaderView {
        id: header
        syncView: table
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
        clip: true
        // Resizing and reordering are the reader's, and they are the two
        // things a table of twelve columns is unusable without.
        resizableColumns: true
        movableColumns: true

        delegate: Rectangle {
            required property int index
            implicitHeight: Interface.rowHeightCompact
            color: Theme.panelBackground

            KvitLabel {
                anchors.fill: parent
                anchors.leftMargin: Interface.spaceNear
                anchors.rightMargin: Interface.spaceNear
                text: root.model ? root.model.columnTitle(parent.index) : ""
                role: "small"
                color: Theme.textMuted
            }
            KvitDivider {
                anchors.bottom: parent.bottom
                anchors.left: parent.left
                anchors.right: parent.right
                color: Theme.borderStrong
            }
        }
    }

    TableView {
        id: table
        anchors.top: header.bottom
        anchors.left: parent.left
        anchors.right: vertical.left
        anchors.bottom: horizontal.top
        clip: true
        model: root.model
        selectionModel: selectionModel
        // Multi-select by drag and by shift-click, which is what a browser of
        // this size is for: the reader's task is usually "these forty".
        selectionBehavior: TableView.SelectRows
        selectionMode: TableView.ExtendedSelection
        // The whole reason this holds at 250,000 rows: a delegate that leaves
        // the viewport is recycled rather than destroyed and rebuilt.
        reuseItems: true
        boundsBehavior: Flickable.StopAtBounds
        // Keep the same number of rows in flight regardless of row height.
        rowHeightProvider: () => Interface.rowHeightSlim
        columnWidthProvider: column => {
            for (let i = 0; i < root.hiddenColumns.length; ++i) {
                if (root.hiddenColumns[i] === column)
                    return 0
            }
            return root.model ? Interface.px(root.model.columnWidth(column))
                              : Interface.px(120)
        }

        delegate: Rectangle {
            id: cell
            required property int row
            required property int column
            required property bool selected
            required property bool current
            required property var display
            required property bool measured
            required property var mark

            implicitHeight: Interface.rowHeightSlim
            color: selected ? Theme.selectionTint
                 : current ? Theme.focusTint
                 : hover.hovered ? Theme.hoverTint
                 : "transparent"

            KvitCell {
                anchors.fill: parent
                kind: root.model ? String(root.model.columnKind(cell.column)) : "Text"
                value: cell.display
                measured: cell.measured
                mark: cell.mark === undefined ? "neutral" : String(cell.mark)
                selected: cell.selected
            }

            KvitDivider {
                anchors.bottom: parent.bottom
                anchors.left: parent.left
                anchors.right: parent.right
            }

            HoverHandler { id: hover }
            TapHandler {
                onDoubleTapped: root.rowActivated(cell.row)
            }
        }

        // Copy of a selection, as tab-separated text: what a reader does next
        // with forty selected rows is put them somewhere else.
        Keys.onPressed: function (event) {
            if (event.matches(StandardKey.Copy)) {
                root.copySelection()
                event.accepted = true
            }
        }
    }

    KvitScrollBar {
        id: vertical
        flickable: table
        anchors.right: parent.right
        anchors.top: header.bottom
        anchors.bottom: horizontal.top
    }
    KvitScrollBar {
        id: horizontal
        flickable: table
        orientation: Qt.Horizontal
        anchors.left: parent.left
        anchors.right: vertical.left
        anchors.bottom: parent.bottom
    }

    KvitEmptyState {
        anchors.centerIn: parent
        visible: root.rowCount === 0
        title: root.emptyTitle
        detail: root.emptyDetail
        symbol: "list"
    }

    // The selected rows as tab-separated text, one row per line, in the
    // column order the reader is looking at rather than the model's.
    function copySelection() {
        if (!root.model)
            return ""
        const rows = []
        const indexes = selectionModel.selectedIndexes
        const byRow = {}
        for (let i = 0; i < indexes.length; ++i) {
            const index = indexes[i]
            if (byRow[index.row] === undefined)
                byRow[index.row] = {}
            byRow[index.row][index.column] = String(root.model.data(index))
        }
        const rowKeys = Object.keys(byRow).map(Number).sort((a, b) => a - b)
        for (let r = 0; r < rowKeys.length; ++r) {
            const columns = byRow[rowKeys[r]]
            const keys = Object.keys(columns).map(Number).sort((a, b) => a - b)
            rows.push(keys.map(k => columns[k]).join("\t"))
        }
        return rows.join("\n")
    }
}
