// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// The side pane: detail about the one thing the reader has selected, beside
// the list they selected it from.
//
// The same width wherever it appears, which is `Interface.paneWidth`. That is
// a decision rather than a coincidence: a pane that is 320 pixels on one
// screen and 420 on another makes the list beside it reflow when the reader
// moves between them, and the reflow reads as the application losing its
// place.
//
// It slides rather than appearing, and the slide follows the reduced-motion
// setting. A pane that appears instantly next to a list makes the list jump
// sideways with no indication of what happened.
KvitPanel {
    id: root

    property string title: ""
    // What the close control is called, both on hover and to a screen
    // reader. A pane holding one record says "Close record" rather than
    // "Close the pane", because the reader is closing the record and the
    // pane is the furniture it arrived in.
    property string closeLabel: qsTr("Close the pane")
    property bool open: true
    default property alias content: bodySlot.data

    signal closeRequested()

    width: Interface.paneWidth
    ruleLeft: true
    visible: x < (parent ? parent.width : 0)

    x: parent ? (root.open ? parent.width - width : parent.width) : 0
    Behavior on x {
        NumberAnimation {
            duration: 160 * Theme.motionScale
            easing.type: Easing.OutCubic
        }
    }

    Accessible.role: Accessible.Pane
    Accessible.name: root.title

    RowLayout {
        id: head
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
        height: Interface.breadcrumbHeight
        anchors.leftMargin: Interface.viewMargin
        anchors.rightMargin: Interface.spaceNear
        spacing: Interface.space

        KvitLabel {
            Layout.fillWidth: true
            text: root.title
            role: "strong"
            font.bold: true
        }
        KvitIconButton {
            symbol: "close"
            label: root.closeLabel
            onClicked: root.closeRequested()
        }
    }

    KvitDivider {
        id: rule
        anchors.top: head.bottom
        anchors.left: parent.left
        anchors.right: parent.right
    }

    Item {
        id: bodySlot
        anchors.top: rule.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.bottom: parent.bottom
    }
}
