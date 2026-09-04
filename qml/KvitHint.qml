// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// A persistent explanation reached through the established information icon.
// Use this for more than a short tooltip: one click or keyboard press opens the
// labelled explanation, keeps the words present while they are read, and
// KvitPopover supplies Escape and outside-click dismissal. There is deliberately
// no tooltip before it: the title and detail are one piece of information, reached
// through one action.
Item {
    id: root

    required property string label
    required property string text
    property int popoverWidth: Interface.px(280)
    readonly property alias opened: explanation.opened

    implicitWidth: trigger.implicitWidth
    implicitHeight: trigger.implicitHeight

    function open() {
        explanation.open()
    }

    function close() {
        explanation.close()
    }

    function shouldOpenLeft(): bool {
        const window = root.Window.window
        if (!window)
            return false
        const origin = root.mapToItem(window.contentItem, 0, 0)
        return origin.x + explanation.width > window.width
    }

    KvitIconButton {
        id: trigger
        objectName: "trigger"
        anchors.centerIn: parent
        symbol: "info"
        label: root.label
        tooltipEnabled: false
        onClicked: explanation.opened ? explanation.close() : explanation.open()
    }

    KvitPopover {
        id: explanation
        objectName: "popover"
        parent: root
        x: root.shouldOpenLeft() ? root.width - width : 0
        y: root.height + Interface.spaceNear
        width: root.popoverWidth
        title: root.label

        contentItem: ColumnLayout {
            spacing: Interface.spaceNear

            KvitLabel {
                objectName: "title"
                Layout.fillWidth: true
                text: root.label
                role: "body"
                font.bold: true
                color: Theme.textPrimary
            }
            KvitLabel {
                objectName: "detail"
                Layout.fillWidth: true
                text: root.text
                role: "small"
                color: Theme.textSecondary
                wrapMode: Text.WordWrap
                elide: Text.ElideNone
                lineHeightMode: Text.ProportionalHeight
                lineHeight: 1.3
            }
        }
    }
}
