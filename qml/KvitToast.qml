// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A message that appears, says one thing, and goes away.
//
// For confirming something that already happened — saved, copied, three rows
// archived — and for the undo that goes with it. Not for anything the reader
// has to act on: a toast that carries the only route to a decision is a
// decision that disappears after four seconds.
//
// A toast with an action stays until it is dismissed. An undo that times out
// while the reader is still reading what happened is an undo they cannot use,
// and the whole reason to offer undo rather than a confirmation dialog is that
// it does not interrupt.
Rectangle {
    id: root

    property string text: ""
    // "info" | "success" | "warning" | "danger"
    property string tone: "info"
    property string action: ""
    property bool shown: false
    // How long it stays, in milliseconds. Ignored when there is an action.
    property int timeout: 4000

    signal actioned()
    signal dismissed()

    readonly property color toneColor: {
        switch (tone) {
        case "success": return Theme.success
        case "warning": return Theme.warning
        case "danger":  return Theme.danger
        default:        return Theme.accent
        }
    }

    implicitHeight: Interface.rowHeightSlim + Interface.space
    implicitWidth: row.width + Interface.spaceLoose * 2

    visible: opacity > 0
    opacity: shown ? 1 : 0
    Behavior on opacity {
        NumberAnimation { duration: 140 * Theme.motionScale }
    }

    radius: Interface.radiusCard
    color: Theme.popupBackground
    border.width: Interface.hairline
    border.color: Theme.borderStrong

    // A live region, so a screen reader says what happened. A toast that is
    // only visual tells a screen reader user nothing at all — and "your change
    // was saved" is exactly the kind of thing they otherwise have to infer.
    Accessible.role: Accessible.AlertMessage
    Accessible.name: root.text
    Accessible.description: root.action

    // The tone stripe down the leading edge: the second channel beside the
    // symbol's colour, so a warning and a confirmation differ by more than hue.
    Rectangle {
        anchors.left: parent.left
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        width: Interface.spaceTight
        radius: parent.radius
        color: root.toneColor
    }

    Timer {
        running: root.shown && root.action === ""
        interval: root.timeout
        onTriggered: root.dismissed()
    }

    Row {
        id: row
        anchors.centerIn: parent
        spacing: Interface.space

        KvitIcon {
            anchors.verticalCenter: parent.verticalCenter
            name: {
                switch (root.tone) {
                case "success": return "success"
                case "warning": return "warning"
                case "danger":  return "error"
                default:        return "info"
                }
            }
            color: root.toneColor
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }
        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            text: root.text
            role: "body"
            elide: Text.ElideNone
        }
        KvitButton {
            anchors.verticalCenter: parent.verticalCenter
            visible: root.action !== ""
            text: root.action
            form: "quiet"
            onClicked: root.actioned()
        }
        KvitIconButton {
            anchors.verticalCenter: parent.verticalCenter
            visible: root.action !== ""
            symbol: "close"
            label: qsTr("Dismiss")
            implicitWidth: Interface.px(20)
            implicitHeight: Interface.px(20)
            onClicked: root.dismissed()
        }
    }
}
