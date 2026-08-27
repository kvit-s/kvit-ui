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
         : hover.hovered ? Theme.hoverTint
         : "transparent"

    Accessible.role: Accessible.ListItem
    Accessible.name: root.label
    Accessible.selectable: true
    Accessible.selected: root.selected

    Item {
        id: layout
        anchors.fill: parent
        anchors.leftMargin: root.form === "sub" ? Interface.px(34)
                                                : Interface.spaceTight
        anchors.rightMargin: Interface.spaceTight
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

    HoverHandler { id: hover }
    TapHandler { onTapped: root.activated() }
}
