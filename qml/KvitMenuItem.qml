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
    // One sentence saying what choosing this does, for an entry whose words
    // cannot say it: that applying a staged change writes it into the note as
    // one undo step, that discarding it tells the agent the change was
    // rejected. Shown as the tooltip and announced as the accessible
    // description, the terms every other acting component in the vocabulary
    // carries it on, and never where the label belongs.
    property string explanation: ""

    implicitHeight: Interface.rowHeightSlim
    implicitWidth: Math.max(Interface.px(180),
                            label.implicitWidth + hint.implicitWidth
                            + Interface.px(72))

    Accessible.role: Accessible.MenuItem
    Accessible.name: root.text
    // The shortcut and the sentence are both descriptions of the same entry,
    // so they are announced together rather than one of them replacing the
    // other.
    Accessible.description: root.explanation === "" ? root.shortcut
                          : root.shortcut === "" ? root.explanation
                          : root.explanation + " " + root.shortcut
    // An entry drawn as the chosen one has to say so. Without these a screen
    // reader announces the profile in use exactly like the four that are not.
    Accessible.checkable: root.checkable || root.checked
    Accessible.checked: root.checked

    KvitTooltip {
        objectName: "tooltip"
        text: root.explanation
        visible: root.explanation !== ""
                 && (root.hovered || root.highlighted)
    }

    background: Rectangle {
        color: root.highlighted ? Theme.hoverTint : "transparent"
    }

    // The check is drawn here rather than left to the host application's
    // style. A MenuItem's own indicator comes from whichever style the
    // application installed, lands in the same left gutter this puts its
    // symbol in, and is therefore a different mark in each of the four
    // applications and sometimes two marks in one place.
    indicator: null

    contentItem: Item {
        KvitIcon {
            id: mark
            anchors.left: parent.left
            anchors.leftMargin: Interface.space
            anchors.verticalCenter: parent.verticalCenter
            visible: root.symbol !== "" || root.checked
            // An entry that is the chosen one is marked as chosen whatever
            // else it would have shown, since the state is the thing the
            // reader is scanning the list for.
            name: root.checked ? "check"
                : root.symbol === "" ? "dot" : root.symbol
            color: root.checked ? Theme.accent
                 : root.destructive ? Theme.danger : Theme.textMuted
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
            // Two channels for the chosen entry rather than the mark alone:
            // the weight as well as the tick, so the entry in use is still
            // the one that stands out in a grayscale capture.
            font.bold: root.checked
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
