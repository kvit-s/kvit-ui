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
import QtQml.Models
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
// `KvitCell` is a Loader over a small component per kind rather than one
// delegate with a branch per kind in it. A branching delegate evaluates every
// branch's bindings for every cell.
//
// And the delegate is one component per column rather than one for the table,
// through a DelegateChooser. TableView keeps its pool of recycled items per
// delegate and picks from it by how near the new index is to the old one, not
// by column, so a single delegate hands most cells an item that was last used
// in a different column — a different kind, and a full rebuild of the cell's
// content. Measured on the benchmark, that is the difference between eleven
// milliseconds to relay out a page and thirty. See the note on TableCell.
//
// ── What the reader gets ───────────────────────────────────────────────────
//
// Column resize, reorder, sort and a per-column menu through the header,
// multi-select through a selection model, keyboard navigation, copy of a
// selection, and a saved column set. The empty state is KvitEmptyState rather
// than a blank grid, which is the same rule the charts follow: nothing to show
// is a thing to say, not an absence to leave.
//
// ── The header is one focus stop, not one per column ───────────────────────
//
// Everything the header offers a pointer — sort, resize, the column menu — is
// on the keyboard too, and the arrangement that makes that possible is that
// the header row takes focus as a single thing and moves a cursor between
// columns with the arrow keys. Putting every header cell in the tab order
// instead would mean a twelve-column table costs twelve presses of Tab to
// walk past, and the cells are recycled delegates, so the thing holding focus
// would be destroyed and rebuilt underneath the reader as the table scrolls
// sideways.
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

    // Which column the rows are ordered by, and in which direction. -1 is no
    // ordering at all, which is the state a table starts in.
    //
    // The header writes both of these itself when the reader sorts, and then
    // says so through `sortRequested`: the model is the application's, so the
    // library cannot reorder anything, and the reader still has to see the
    // indicator move on the press rather than a frame after the rows change.
    // An application that keeps the sort in a saved view sets them back on
    // load.
    property int sortColumn: -1
    property bool sortAscending: true

    // The select-displayed control that a `Check` column draws in its header,
    // and its third state for "some of the rows shown".
    //
    // The state lives with the application rather than with the checkbox,
    // because what "all of them" means is the application's question: a
    // selection keyed by row would empty itself the moment a filter changed.
    property bool headerChecked: false
    property bool headerPartial: false

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
    onHiddenColumnsChanged: {
        table.forceLayout()
        if (internal.headerColumn >= 0 && root.isColumnHidden(internal.headerColumn))
            internal.headerColumn = root.stepHeaderColumn(internal.headerColumn, 1)
    }

    // A row was opened. `rowPressed` is one press, which is how a ledger row
    // opens its record; `rowActivated` is the double press, kept because it
    // is what every existing caller listens to.
    signal rowActivated(int row)
    signal rowPressed(int row)

    // The reader asked for a different order. The table has already moved its
    // own indicator to match, so a handler that cannot sort should put
    // `sortColumn` and `sortAscending` back rather than ignore this.
    signal sortRequested(int column, bool ascending)

    // The select-displayed control in a Check column's header was toggled.
    signal headerToggled(bool checked)
    // One row's box in a Check column was toggled, with the value it now
    // shows.
    signal cellToggled(int row, int column, bool checked)

    implicitWidth: Interface.px(600)
    implicitHeight: Interface.px(400)

    Accessible.role: Accessible.Table
    Accessible.name: qsTr("Table")

    // View state that is nobody else's business: the header's keyboard
    // cursor, which column the open menu belongs to, and which row the
    // pointer is over. Held here rather than as properties of the table so
    // that the component's published interface stays the fifteen things above.
    QtObject {
        id: internal
        property int headerColumn: -1
        property int menuColumn: -1
        // Hover is a property of the row rather than of the cell under the
        // pointer: a row that says "I open something" one cell at a time says
        // it about the cell.
        property int hoveredRow: -1
        // Whether the keyboard is in the table at all. The cell cursor keeps
        // its place when focus leaves, so a cell that disclosed its full
        // value on being current alone would leave a tooltip floating over a
        // surface the reader has already moved away from.
        property bool keyboardHeld: table.activeFocus
    }

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

        // One stop in the tab order for the whole header. See the note at the
        // top of the file for why it is not one per column.
        activeFocusOnTab: true
        onActiveFocusChanged: {
            if (header.activeFocus && internal.headerColumn < 0)
                internal.headerColumn = root.stepHeaderColumn(-1, 1)
        }

        Keys.onPressed: function (event) {
            const column = internal.headerColumn
            // Width first, because Ctrl+Left is a resize and a bare Left is a
            // move, and testing the bare key first would swallow both.
            if (event.modifiers & Qt.ControlModifier) {
                if (event.key === Qt.Key_Left) {
                    root.resizeColumnBy(column, -Interface.px(16))
                    event.accepted = true
                    return
                }
                if (event.key === Qt.Key_Right) {
                    root.resizeColumnBy(column, Interface.px(16))
                    event.accepted = true
                    return
                }
            }
            // Alt+Down is the chord every desktop uses for "open this
            // control's list", and Key_Menu is the dedicated key for it.
            if ((event.modifiers & Qt.AltModifier) && event.key === Qt.Key_Down) {
                root.openColumnMenu(column)
                event.accepted = true
                return
            }
            switch (event.key) {
            case Qt.Key_Left:
                internal.headerColumn = root.stepHeaderColumn(column, -1)
                event.accepted = true
                break
            case Qt.Key_Right:
                internal.headerColumn = root.stepHeaderColumn(column, 1)
                event.accepted = true
                break
            case Qt.Key_Home:
                internal.headerColumn = root.stepHeaderColumn(-1, 1)
                event.accepted = true
                break
            case Qt.Key_End:
                internal.headerColumn = root.stepHeaderColumn(-1, -1)
                event.accepted = true
                break
            case Qt.Key_Return:
            case Qt.Key_Enter:
            case Qt.Key_Space:
                root.activateHeaderColumn(column)
                event.accepted = true
                break
            case Qt.Key_Menu:
                root.openColumnMenu(column)
                event.accepted = true
                break
            case Qt.Key_Down:
                // Out of the header and into the rows, which is where the
                // reader was heading.
                table.forceActiveFocus(Qt.TabFocusReason)
                event.accepted = true
                break
            default:
                break
            }
        }

        delegate: Rectangle {
            id: headerCell
            required property int index

            readonly property string kindName:
                root.model ? root.model.columnKindName(headerCell.index) : "Text"
            readonly property bool sortable:
                root.model ? root.model.columnSortable(headerCell.index) : false
            readonly property bool sorted: root.sortColumn === headerCell.index
            readonly property bool cursorHere:
                header.activeFocus && internal.headerColumn === headerCell.index

            implicitHeight: Interface.rowHeightCompact
            color: Theme.panelBackground

            Accessible.role: Accessible.ColumnHeader
            Accessible.name: root.headerName(headerCell.index)
            Accessible.description: headerCell.kindName === "Check"
                ? qsTr("Enter selects or clears every row shown; the menu key opens this column's options")
                : headerCell.sortable
                  ? qsTr("Enter sorts this column; the menu key opens its options")
                  : qsTr("This column does not sort; the menu key opens its options")
            Accessible.onPressAction: root.activateHeaderColumn(headerCell.index)

            // The select-displayed control, in the one column that draws
            // boxes. It stays bound to `headerChecked` rather than keeping
            // its own state: the reader's answer leaves as a signal and comes
            // back as the property, so the box can never say "all" while the
            // application believes something else.
            KvitCheck {
                id: selectDisplayed
                visible: headerCell.kindName === "Check"
                anchors.left: parent.left
                anchors.leftMargin: Interface.spaceTight
                anchors.verticalCenter: parent.verticalCenter
                height: Interface.rowHeightCompact
                text: ""
                // The header row owns the keyboard: Enter on this column
                // toggles the box. A control inside a recycled delegate must
                // not take focus of its own, because the item holding it is
                // destroyed when the column scrolls out from under it.
                focusPolicy: Qt.NoFocus
                checked: root.headerChecked
                partial: root.headerPartial
                Accessible.name: qsTr("Select the rows shown")
                onToggled: {
                    const wanted = selectDisplayed.checked
                    selectDisplayed.checked = Qt.binding(() => root.headerChecked)
                    root.headerToggled(wanted)
                }

                KvitTooltip {
                    text: qsTr("Select every row the current filter shows")
                    visible: selectDisplayed.hovered
                }
            }

            KvitLabel {
                id: title
                visible: headerCell.kindName !== "Check"
                anchors.left: parent.left
                anchors.leftMargin: Interface.spaceNear
                anchors.right: headerMarks.left
                anchors.rightMargin: Interface.spaceTight
                anchors.verticalCenter: parent.verticalCenter
                text: root.model ? root.model.columnTitle(headerCell.index) : ""
                role: "small"
                // The sorted column is named in the reader's own text colour
                // while every other column stays context. That is the second
                // channel beside the indicator's hue, and the indicator's
                // shape is the third.
                color: headerCell.sorted ? Theme.textPrimary : Theme.textMuted
            }

            Row {
                id: headerMarks
                anchors.right: parent.right
                anchors.rightMargin: Interface.spaceTight
                anchors.verticalCenter: parent.verticalCenter
                spacing: Interface.spaceTight

                // Which way the column is sorted, as two different glyphs
                // rather than one glyph in two colours: the direction has to
                // survive the grayscale reading of a screenshot.
                KvitIcon {
                    visible: headerCell.sorted
                    name: root.sortAscending ? "sort-ascending" : "sort-descending"
                    color: Theme.accent
                    implicitWidth: Interface.iconSizeSmall
                    implicitHeight: Interface.iconSizeSmall
                    anchors.verticalCenter: parent.verticalCenter
                }

                KvitIconButton {
                    id: columnMenuButton
                    // Shown on hover and while the keyboard cursor is on this
                    // column. A button on every header the whole time turns a
                    // twelve-column header into a row of dots.
                    visible: headerHover.hovered || headerCell.cursorHere
                    symbol: "more-vertical"
                    label: qsTr("Column options")
                    width: Interface.rowHeightCompact
                    height: Interface.rowHeightCompact
                    // The header row owns the keyboard; a button inside a
                    // recycled delegate must not take focus, because the item
                    // holding it is destroyed when the column scrolls out.
                    focusPolicy: Qt.NoFocus
                    onClicked: root.openColumnMenu(headerCell.index)
                    anchors.verticalCenter: parent.verticalCenter
                }
            }

            // The keyboard cursor, drawn inside the cell so it is not clipped
            // by the header's own `clip`.
            Rectangle {
                anchors.fill: parent
                anchors.margins: Interface.focusRingWidth
                visible: headerCell.cursorHere
                color: "transparent"
                radius: Interface.radiusControl
                border.width: Interface.focusRingWidth
                border.color: Theme.focusRing
            }

            KvitDivider {
                anchors.bottom: parent.bottom
                anchors.left: parent.left
                anchors.right: parent.right
                color: Theme.borderStrong
            }

            HoverHandler { id: headerHover }
            TapHandler {
                acceptedButtons: Qt.LeftButton
                onSingleTapped: {
                    internal.headerColumn = headerCell.index
                    root.activateHeaderColumn(headerCell.index)
                }
            }
            TapHandler {
                acceptedButtons: Qt.RightButton
                onSingleTapped: root.openColumnMenu(headerCell.index)
            }
        }
    }

    // One menu for the whole header rather than one per header cell. The
    // cells are recycled delegates: a menu built inside one is destroyed
    // while it is open the first time the reader scrolls sideways.
    KvitMenu {
        id: columnMenu
        title: internal.menuColumn >= 0 && root.model
               ? root.model.columnTitle(internal.menuColumn) : ""

        KvitMenuItem {
            text: qsTr("Sort ascending")
            symbol: "sort-ascending"
            enabled: internal.menuColumn >= 0 && root.model
                     ? root.model.columnSortable(internal.menuColumn) : false
            onTriggered: root.sortBy(internal.menuColumn, true)
        }
        KvitMenuItem {
            text: qsTr("Sort descending")
            symbol: "sort-descending"
            enabled: internal.menuColumn >= 0 && root.model
                     ? root.model.columnSortable(internal.menuColumn) : false
            onTriggered: root.sortBy(internal.menuColumn, false)
        }
        KvitMenuItem {
            text: qsTr("Clear the sort")
            symbol: "caret-sort"
            enabled: root.sortColumn >= 0
            onTriggered: root.sortBy(-1, true)
        }
        KvitMenuItem {
            text: qsTr("Hide this column")
            symbol: "eye-off"
            enabled: internal.menuColumn >= 0
            onTriggered: root.hideColumn(internal.menuColumn)
        }
        KvitMenuItem {
            text: qsTr("Show every column")
            symbol: "eye"
            enabled: root.hiddenColumns.length > 0
            onTriggered: root.hiddenColumns = []
        }
    }

    // One cell, and the reason it is an inline component instantiated once
    // per column rather than a single delegate.
    //
    // TableView recycles a cell that scrolls out of view, and it hands the
    // recycled item to whichever cell it builds next — it matches on how near
    // the new index is to the old one, not on the column. Scrolling a page
    // therefore gives almost every cell an item that was last used in a
    // different column, and a cell whose column changed has changed kind, so
    // KvitCell's loader tears its content down and builds another. On a
    // thirteen-column ledger that is three hundred component trees per page,
    // and it takes the relayout of one page from eleven milliseconds to
    // thirty — measured on the benchmark model, one page at a time, taking
    // the middle page of forty.
    //
    // The pool is kept per delegate, so a DelegateChooser that answers with a
    // different component object per column gives each column a pool of its
    // own and a recycled item always comes back to the column it left. The
    // components are identical; it is their being separate objects that does
    // the work.
    component TableCell: Rectangle {
        id: cell

        // The table this cell belongs to, and its private view state. Passed
        // in rather than read from the enclosing scope, because an inline
        // component has a scope of its own and the ids out here are not in
        // it.
        required property var host
        required property var chrome

        required property int row
        required property int column
        required property bool selected
        required property bool current
        required property var display
        required property bool measured
        required property var mark
        // `var` rather than the types they hold, because a model that
        // answers nothing for a role gives the delegate `undefined`, and
        // a typed required property refuses it out loud.
        required property var marks
        required property var checked
        required property var unit
        required property var fullText

        // Guarded, because a delegate's bindings are evaluated once while it
        // is being incubated and once more while it is being torn down, and
        // at neither moment does it have the two objects passed in below.
        readonly property bool rowHovered:
            cell.chrome ? cell.chrome.hoveredRow === cell.row : false

        implicitHeight: Interface.rowHeightSlim
        color: cell.selected ? Theme.selectionTint
             : cell.current ? Theme.focusTint
             : cell.rowHovered ? Theme.hoverTint
             : "transparent"

        KvitCell {
            anchors.fill: parent
            // The enumeration key rather than the enumeration. A Q_ENUM
            // arrives in JavaScript as an integer, so `String()` on it
            // yields "1" and every column in the estate drew as text.
            kind: cell.host && cell.host.model
                  ? cell.host.model.columnKindName(cell.column) : "Text"
            value: cell.display
            measured: cell.measured
            mark: cell.mark === undefined ? "neutral" : String(cell.mark)
            marks: cell.marks === undefined || cell.marks === null
                   ? [] : cell.marks
            checked: cell.checked === true
            unit: cell.unit === undefined || cell.unit === null
                  ? "" : String(cell.unit)
            fullValue: cell.fullText === undefined || cell.fullText === null
                       ? "" : String(cell.fullText)
            selected: cell.selected
            // Where the keyboard is, so the cell discloses a value too long
            // for its column to a reader who never touches the pointer. D43
            // asks for the whole name on hover and by keyboard, and a cell
            // takes no focus of its own to answer the second half with.
            current: cell.current && cell.chrome
                     ? cell.chrome.keyboardHeld : false
            // The row rather than the cell. A reader moving along a row of
            // twelve columns is told "42.10" and nothing about which record
            // that belongs to, so every cell says which row it is in after
            // it has said what it holds.
            Accessible.description: cell.host && cell.host.model
                                    ? cell.host.model.rowName(cell.row) : ""
            onToggled: state => cell.host.cellToggled(cell.row, cell.column, state)
        }

        KvitDivider {
            anchors.bottom: parent.bottom
            anchors.left: parent.left
            anchors.right: parent.right
        }

        // The row says it opens something before it is pressed: the pointer
        // becomes a hand and the row is underlined, which is the same pair of
        // signals a link carries everywhere else in the estate. The tint
        // alone said only "the pointer is here".
        Rectangle {
            anchors.bottom: parent.bottom
            anchors.left: parent.left
            anchors.right: parent.right
            height: Interface.focusRingWidth
            visible: cell.rowHovered
            color: Theme.link
        }

        // A delegate that is recycled under a stationary pointer changes
        // which row it draws without the hover handler hearing anything, so
        // scrolling with the wheel would leave the highlight on the row that
        // used to be there.
        onRowChanged: {
            if (cell.chrome && hover.hovered)
                cell.chrome.hoveredRow = cell.row
        }

        HoverHandler {
            id: hover
            cursorShape: Qt.PointingHandCursor
            onHoveredChanged: {
                if (!cell.chrome)
                    return
                if (hover.hovered)
                    cell.chrome.hoveredRow = cell.row
                else if (cell.chrome.hoveredRow === cell.row)
                    cell.chrome.hoveredRow = -1
            }
        }
        TapHandler {
            // One press opens the record. The double press is kept because
            // every caller written against this listens for it, and because a
            // reader who double-clicks a row means the same thing twice
            // rather than something else.
            onSingleTapped: cell.host.rowPressed(cell.row)
            onDoubleTapped: cell.host.rowActivated(cell.row)
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
            const explicit = table.explicitColumnWidth(column)
            if (explicit >= 0)
                return explicit
            return root.model ? Interface.px(root.model.columnWidth(column))
                              : Interface.px(120)
        }

        // One delegate per column rather than one for the whole table. See
        // the note on TableCell for why that is what makes a table of mixed
        // kinds hold its frame budget.
        delegate: DelegateChooser {
            DelegateChoice { column: 0;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 1;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 2;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 3;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 4;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 5;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 6;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 7;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 8;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 9;  delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 10; delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 11; delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 12; delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 13; delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 14; delegate: TableCell { host: root; chrome: internal } }
            DelegateChoice { column: 15; delegate: TableCell { host: root; chrome: internal } }
            // Every column past the sixteenth shares one pool and pays the
            // rebuild the fifteen above avoid. Sixteen is where the estate's
            // widest table sits; a table wider than that still draws
            // correctly and only scrolls less cheaply.
            DelegateChoice { delegate: TableCell { host: root; chrome: internal } }
        }

        // The keyboard half of what the pointer can do to a row: copy a
        // selection, open the current row, and tick its box where the table
        // has one.
        Keys.onPressed: function (event) {
            if (event.matches(StandardKey.Copy)) {
                root.copySelection()
                event.accepted = true
                return
            }
            const current = selectionModel.currentIndex
            if (!current.valid)
                return
            if (event.key === Qt.Key_Return || event.key === Qt.Key_Enter) {
                root.rowPressed(current.row)
                root.rowActivated(current.row)
                event.accepted = true
            } else if (event.key === Qt.Key_Space) {
                const box = root.checkColumn()
                if (box < 0)
                    return
                root.cellToggled(current.row, box,
                                 !root.model.cellChecked(current.row, box))
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

    // Put the keyboard on one row: the cell cursor moves there, the row is
    // scrolled into view and the table takes focus. A surface raised from a
    // row uses this to hand focus back to the row that raised it when it
    // closes, which is the only way the reader keeps their place.
    function focusRow(row) {
        if (!root.model || row < 0 || row >= root.rowCount)
            return
        table.positionViewAtRow(row, TableView.Contain)
        selectionModel.setCurrentIndex(root.model.index(row, 0),
                                       ItemSelectionModel.NoUpdate)
        table.forceActiveFocus()
    }

    // Whether the reader has hidden this column.
    function isColumnHidden(column) {
        for (let i = 0; i < root.hiddenColumns.length; ++i) {
            if (root.hiddenColumns[i] === column)
                return true
        }
        return false
    }

    // The next column the header cursor should land on, skipping hidden ones
    // and wrapping at both ends. `from` of -1 with a step of 1 is the first
    // column and with a step of -1 is the last, which is what Home and End
    // want.
    function stepHeaderColumn(from, step) {
        const count = table.columns
        if (count <= 0)
            return -1
        let column = from
        for (let i = 0; i < count; ++i) {
            column += step
            if (column < 0)
                column = count - 1
            else if (column >= count)
                column = 0
            if (!root.isColumnHidden(column))
                return column
        }
        return from
    }

    // What a screen reader says about a header: the column, and the order it
    // is holding the rows in.
    function headerName(column) {
        if (!root.model || column < 0)
            return ""
        const title = root.model.columnTitle(column)
        if (root.sortColumn !== column)
            return title
        return root.sortAscending ? qsTr("%1, sorted ascending").arg(title)
                                  : qsTr("%1, sorted descending").arg(title)
    }

    // Activating a header sorts on it, and activating the column it is
    // already sorted on turns the order around. That is the one gesture, and
    // it is the same whether it arrived from a press or from Enter.
    function activateHeaderColumn(column) {
        if (column < 0 || !root.model)
            return
        // A Check column's header is the select-displayed control rather than
        // a sort, so activating it does what pressing the box does.
        if (root.model.columnKindName(column) === "Check") {
            root.headerToggled(!root.headerChecked)
            return
        }
        if (!root.model.columnSortable(column))
            return
        root.sortBy(column, root.sortColumn === column ? !root.sortAscending : true)
    }

    function sortBy(column, ascending) {
        root.sortColumn = column
        root.sortAscending = ascending
        root.sortRequested(column, ascending)
    }

    function hideColumn(column) {
        if (column < 0 || root.isColumnHidden(column))
            return
        root.hiddenColumns = root.hiddenColumns.concat([column])
    }

    // Keyboard resize, in steps. A floor rather than none, because a column
    // narrowed to nothing is a column the reader cannot find again to widen.
    function resizeColumnBy(column, amount) {
        if (column < 0 || root.isColumnHidden(column))
            return
        const floor = Interface.px(48)
        table.setColumnWidth(column,
                             Math.max(floor, table.columnWidth(column) + amount))
        table.forceLayout()
    }

    // Which column, if any, draws boxes. -1 when the table has none, which is
    // every table but a ledger with a selection in it.
    function checkColumn() {
        if (!root.model)
            return -1
        for (let i = 0; i < table.columns; ++i) {
            if (root.model.columnKindName(i) === "Check")
                return i
        }
        return -1
    }

    // Where the menu for a column opens: under the header, at the column's
    // left edge, pulled back inside the table when the column is the last one.
    function openColumnMenu(column) {
        if (column < 0)
            return
        internal.menuColumn = column
        internal.headerColumn = column
        let left = -table.contentX
        for (let i = 0; i < column; ++i)
            left += table.columnWidth(i)
        columnMenu.x = Math.max(0, Math.min(root.width - columnMenu.width, left))
        columnMenu.y = header.height
        columnMenu.open()
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
            // The cell's unit goes with it. A Money cell draws its currency
            // code beside the amount rather than inside it, so a copy that
            // took only the display value would give a column of bare
            // numbers in a workspace holding several currencies.
            const shown = String(root.model.data(index))
            const unit = root.model.data(index, TableModelBase.UnitRole)
            byRow[index.row][index.column] =
                unit === undefined || unit === null || String(unit) === ""
                    ? shown : shown + " " + String(unit)
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
