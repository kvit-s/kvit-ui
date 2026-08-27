// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// One tab in a row of them: a different view of the same thing.
//
// The selected tab is marked by an underline as well as by colour and weight,
// which is the second channel the high-contrast pass looks for. Selection
// shown by colour alone is the most common place this rule gets broken,
// because a coloured tab label looks obviously different to whoever drew it.
AbstractButton {
    id: root

    property bool selected: false
    // A count beside the label — how many things this tab holds.
    property int count: -1

    implicitHeight: Interface.tabHeight
    implicitWidth: content.implicitWidth + Interface.spaceLoose * 2

    Accessible.role: Accessible.PageTab
    Accessible.name: root.count >= 0
        ? qsTr("%1, %n item(s)", "", root.count).arg(root.text)
        : root.text
    Accessible.selected: root.selected
    Accessible.onPressAction: root.clicked()

    background: Rectangle {
        color: root.hovered && !root.selected ? Theme.hoverTint : "transparent"

        // The underline: two design pixels, so it survives at the smallest
        // interface size and does not read as a hairline rule.
        Rectangle {
            anchors.bottom: parent.bottom
            anchors.left: parent.left
            anchors.right: parent.right
            height: Interface.spaceTight
            visible: root.selected
            color: Theme.accent
        }

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.visualFocus
            color: "transparent"
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    contentItem: Row {
        id: content
        anchors.centerIn: parent
        spacing: Interface.spaceNear

        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            text: root.text
            role: "body"
            font.bold: root.selected
            color: root.selected ? Theme.textPrimary : Theme.textMuted
            elide: Text.ElideNone
        }
        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            visible: root.count >= 0
            text: String(root.count)
            role: "caption"
            color: Theme.textFaint
            tabular: true
            elide: Text.ElideNone
        }
    }
}
