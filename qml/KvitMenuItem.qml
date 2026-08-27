// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// One line of a menu.
//
// It carries its shortcut on the right, which is how a menu teaches: a reader
// who opens a menu for the same command four times is being shown, four times,
// that there is a faster way.
//
// A destructive item is separated by more than colour — it takes the danger
// colour *and* sits below a divider its caller places — because a red menu
// item in a list of black ones is the most common colour-only distinction in
// any application, and the consequence of misreading it is the worst.
MenuItem {
    id: root

    property string symbol: ""
    property string shortcut: ""
    property bool destructive: false

    implicitHeight: Interface.rowHeightSlim
    implicitWidth: Math.max(Interface.px(180),
                            label.implicitWidth + hint.implicitWidth
                            + Interface.px(72))

    Accessible.role: Accessible.MenuItem
    Accessible.name: root.text
    Accessible.description: root.shortcut

    background: Rectangle {
        color: root.highlighted ? Theme.hoverTint : "transparent"
    }

    contentItem: Item {
        KvitIcon {
            id: mark
            anchors.left: parent.left
            anchors.leftMargin: Interface.space
            anchors.verticalCenter: parent.verticalCenter
            visible: root.symbol !== ""
            name: root.symbol === "" ? "dot" : root.symbol
            color: root.destructive ? Theme.danger : Theme.textMuted
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }
        KvitLabel {
            id: label
            anchors.left: parent.left
            anchors.leftMargin: Interface.px(30)
            anchors.right: hint.left
            anchors.rightMargin: Interface.spaceLoose
            anchors.verticalCenter: parent.verticalCenter
            text: root.text
            role: "body"
            color: !root.enabled ? Theme.textDisabled
                 : root.destructive ? Theme.danger
                 : Theme.textPrimary
        }
        KvitLabel {
            id: hint
            anchors.right: parent.right
            anchors.rightMargin: Interface.space
            anchors.verticalCenter: parent.verticalCenter
            text: root.shortcut
            role: "small"
            color: Theme.textFaint
            elide: Text.ElideNone
        }
    }
}
