// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// A message that stays: a banner across the top of a view saying something is
// wrong, out of date, or waiting.
//
// prd.md §5.6 lists this in Group B — absent everywhere and named in a written
// flow — and the reason it is separate from KvitToast is the one thing a
// reader needs to know about a message: whether it will still be there when
// they look back. A toast is an acknowledgement of something that finished. A
// notice is a condition that is still true.
//
// So a notice is dismissible only when dismissing it is meaningful. "Your
// licence expires in three days" can be dismissed; "this file could not be
// saved" cannot, because the condition outlives the reader's acknowledgement
// of it and hiding it would be the application forgetting.
Rectangle {
    id: root

    property string text: ""
    property string detail: ""
    // "info" | "warning" | "danger" | "success"
    property string tone: "info"
    property string action: ""
    property bool dismissible: false

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

    implicitHeight: Math.max(Interface.rowHeightSlim + Interface.space,
                             layout.implicitHeight + Interface.space)
    implicitWidth: parent ? parent.width : Interface.px(600)

    radius: Interface.radiusControl
    color: Qt.alpha(root.toneColor, 0.12)
    border.width: Interface.hairline
    border.color: root.toneColor

    Accessible.role: Accessible.AlertMessage
    Accessible.name: root.detail !== ""
        ? qsTr("%1. %2").arg(root.text).arg(root.detail) : root.text

    RowLayout {
        id: layout
        anchors.fill: parent
        anchors.leftMargin: Interface.spaceLoose
        anchors.rightMargin: Interface.spaceNear
        spacing: Interface.space

        KvitIcon {
            Layout.alignment: Qt.AlignTop
            Layout.topMargin: Interface.spaceNear
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

        ColumnLayout {
            Layout.fillWidth: true
            spacing: Interface.spaceTight
            KvitLabel {
                Layout.fillWidth: true
                text: root.text
                role: "body"
                wrapMode: Text.WordWrap
                elide: Text.ElideNone
            }
            KvitLabel {
                Layout.fillWidth: true
                visible: root.detail !== ""
                text: root.detail
                role: "small"
                color: Theme.textMuted
                wrapMode: Text.WordWrap
                elide: Text.ElideNone
            }
        }

        KvitButton {
            Layout.alignment: Qt.AlignVCenter
            visible: root.action !== ""
            text: root.action
            form: "ordinary"
            onClicked: root.actioned()
        }
        KvitIconButton {
            Layout.alignment: Qt.AlignVCenter
            visible: root.dismissible
            symbol: "close"
            label: qsTr("Dismiss this notice")
            onClicked: root.dismissed()
        }
    }
}
