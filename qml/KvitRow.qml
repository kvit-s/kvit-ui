// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A row in a list, at one of four heights.
//
// kvit-hub's ListRow, generalised. The height is chosen by what the row
// carries rather than by how many rows a view wants to fit: 26 projects at the
// full height need 1,456 pixels against the 874 a stage has, so a list at rest
// uses the slim height and a layer that adds content spends the description
// rather than the row height. Shrinking a row below these heights to fit more
// in is what made kvit-hub's previous interface read as busy.
//
//   full     name over description, with chips and figures beside them
//   sub      a line expanded underneath a row, indented under its parent
//   slim     one line: a name, what it is, one phrase and one figure
//   compact  a disclosure or a group heading
//
// Hover and keyboard focus are separate tints, because the row under the
// pointer and the row the keyboard is on are different rows and a reader
// arrowing down a list needs to see which is which.
//
// A row that acts is opened by the pointer, by Return, by Space and by the
// press action a screen reader offers, and all four do the same one thing.
// Only the keyboard half needs explaining. A row is usually a container, and
// the key a person presses on it may be meant for something inside it — a
// KvitButton at the right-hand end, a KvitLink in the middle — so the row
// answers a key only while the row itself is the focused item. A child with
// focus is left to answer for itself, which is what stops one Space from both
// pressing the button and opening the record behind it.
//
// A row that does something when it is pressed says so before it is pressed,
// and a row that does nothing says nothing. Both halves matter: a list whose
// rows open a record and a pane whose rows are a field name beside its value
// are built from this same component, and when every one of them tints under
// the pointer the tint stops meaning anything. So the hover tint, the tap and
// what a screen reader is told all follow `interactive`, which a row that acts
// declares — either by setting it, or by taking `activeFocusOnTab`, since a
// row the pointer can act on and the keyboard cannot is unreachable for half
// the readers anyway.
Rectangle {
    id: root

    // "full" | "sub" | "slim" | "compact"
    property string form: "slim"
    property bool selected: false
    property bool current: false
    property bool rule: true
    // What a screen reader says when the caret reaches this row. A row whose
    // content is several separate labels reads as a jumble without it.
    property string label: ""
    // Whether pressing this row does something. False draws no hover tint and
    // emits no `activated()`, which is what a row that is only a layout —
    // a field name beside its value — should look like and do. It follows
    // `activeFocusOnTab` by default, so a row already reachable by keyboard
    // needs nothing extra; a row inside a list that moves its own cursor with
    // the arrow keys sets this instead.
    property bool interactive: root.activeFocusOnTab

    signal activated()

    readonly property bool hovered: hover.hovered
    default property alias content: layout.data

    implicitHeight: {
        switch (form) {
        case "full":    return Interface.rowHeight
        case "sub":     return Interface.rowHeightSub
        case "compact": return Interface.rowHeightCompact
        default:        return Interface.rowHeightSlim
        }
    }
    implicitWidth: layout.implicitWidth

    color: selected ? Theme.selectionTint
         : current ? Theme.focusTint
         : (root.interactive && hover.hovered) ? Theme.hoverTint
         : "transparent"

    // What a screen reader is told, which follows `interactive` for the same
    // reason the tint does. A list item is something a reader moves to,
    // selects and opens; a record pane whose twenty-three field rows all
    // announce themselves that way offers a reader who cannot see them
    // twenty-three things to try, none of which do anything. A row that only
    // arranges other things is a grouping of what it holds, and one carrying
    // a label of its own reads as that line of text.
    //
    // Selection is announced wherever it exists: a row that acts can be
    // picked, and a row drawn as picked says so whatever else it is.
    Accessible.role: root.interactive ? Accessible.ListItem
                   : root.label !== "" ? Accessible.StaticText
                                       : Accessible.Grouping
    Accessible.name: root.label
    Accessible.selectable: root.interactive || root.selected
    Accessible.selected: root.selected

    Item {
        id: layout
        anchors.fill: parent
        anchors.leftMargin: root.form === "sub" ? Interface.px(34)
                                                : Interface.spaceTight
        anchors.rightMargin: Interface.spaceTight
    }

    // The row the keyboard is on is outlined as well as tinted. Every other
    // control in the estate draws this ring, and a tint on its own is a two
    // percent lightness difference in the high-contrast theme and invisible
    // in a grayscale screenshot.
    //
    // Two ways to be that row, because there are two kinds of list. A list
    // that carries its own cursor sets `current` on the row the cursor is on;
    // a row that joined the tab order is that row when Tab reaches it, and a
    // row that answers Return without saying it has the keyboard is a row
    // nobody can tell is about to open.
    Rectangle {
        objectName: "focusRing"
        anchors.fill: parent
        anchors.margins: Interface.focusRingWidth
        visible: root.current || root.activeFocus
        color: "transparent"
        border.width: Interface.focusRingWidth
        border.color: Theme.focusRing
    }

    // The hairline belongs to the row rather than to the list, so a group that
    // ends mid-list still closes. An expanded sub-line has none: it belongs to
    // the row above it and a rule there would cut it off.
    KvitDivider {
        visible: root.rule && root.form !== "sub" && root.form !== "compact"
        anchors.bottom: parent.bottom
        anchors.left: parent.left
        anchors.right: parent.right
    }

    // The handler stays in place whether or not the row acts, because
    // `hovered` is published and a caller uses it to reveal something inside
    // the row — a disclosure on a truncated value — that is not the row being
    // pressed.
    HoverHandler { id: hover }
    TapHandler {
        enabled: root.interactive
        onTapped: root.activated()
    }

    // Return, Enter and Space, answered only while the row itself holds the
    // keyboard.
    //
    // Qt sends a key to the focused item and then up its parents, so without
    // the `activeFocus` test a row holding a KvitButton would open the record
    // on the same Space that pressed the button — twice for one keystroke,
    // once visibly and once not. AbstractButton accepts Space, so that
    // particular pair would be caught by the event being consumed, but
    // nothing consumes Return, and a row is a container often enough that the
    // rule is worth stating rather than inheriting from which keys Qt Quick
    // Controls happens to handle.
    Keys.onPressed: event => {
        if (!root.interactive || !root.activeFocus)
            return
        switch (event.key) {
        case Qt.Key_Return:
        case Qt.Key_Enter:
        case Qt.Key_Space:
            root.activated()
            event.accepted = true
            break
        default:
            break
        }
    }

    // The same thing again for assistive technology, which does not press
    // keys: a screen reader offers the press action of whatever it is on, and
    // a row announced as a list item with no press action is a row it can
    // read out and cannot open. This is not behind `activeFocus` — the caller
    // has named the row it means rather than arrived at it by tabbing.
    Accessible.onPressAction: {
        if (root.interactive)
            root.activated()
    }
}
