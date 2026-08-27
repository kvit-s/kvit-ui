// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// One place in the sidebar.
//
// kvit-hub's SidebarItem generalised. The selected item is marked by a bar
// down its leading edge as well as by a tint, which is the second channel: a
// tint alone is a two-percent lightness difference in the high-contrast theme
// and invisible in a screenshot.
//
// `symbol` is required rather than optional because of the rail. When the
// sidebar collapses the label is gone and the symbol is the whole item, and an
// item that only has words vanishes.
AbstractButton {
    id: root

    required property string symbol
    property bool selected: false
    property bool collapsed: false
    property int count: -1

    implicitHeight: Interface.rowHeightSlim
    implicitWidth: Interface.sidebarWidth

    Accessible.role: Accessible.ListItem
    Accessible.name: root.count >= 0
        ? qsTr("%1, %n item(s)", "", root.count).arg(root.text)
        : root.text
    Accessible.selected: root.selected
    Accessible.onPressAction: root.clicked()

    // The tooltip is what makes the rail usable: collapsed, the label is the
    // only thing that says where this goes.
    //
    // KvitTooltip rather than the attached `ToolTip.text`, which instantiates
    // the platform style's own and arrives as a yellow box belonging to no
    // theme here.
    KvitTooltip {
        text: root.text
        visible: root.collapsed && root.hovered
        delay: 400
    }

    background: Rectangle {
        color: root.selected ? Theme.selectionTint
             : root.hovered ? Theme.hoverTint
             : "transparent"

        Rectangle {
            anchors.left: parent.left
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            width: Interface.spaceTight
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

    contentItem: Item {
        anchors.fill: parent

        KvitIcon {
            id: mark
            anchors.left: parent.left
            anchors.leftMargin: root.collapsed
                ? (root.width - width) / 2 : Interface.spaceLoose
            anchors.verticalCenter: parent.verticalCenter
            name: root.symbol
            color: root.selected ? Theme.accent : Theme.textMuted
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }

        KvitLabel {
            anchors.left: mark.right
            anchors.leftMargin: Interface.space
            anchors.right: badge.left
            anchors.rightMargin: Interface.spaceNear
            anchors.verticalCenter: parent.verticalCenter
            visible: !root.collapsed
            text: root.text
            role: "body"
            font.bold: root.selected
            color: root.selected ? Theme.textPrimary : Theme.textSecondary
        }

        KvitBadge {
            id: badge
            anchors.right: parent.right
            anchors.rightMargin: Interface.space
            anchors.verticalCenter: parent.verticalCenter
            visible: !root.collapsed && root.count > 0
            count: root.count
            tone: "neutral"
        }
    }
}
