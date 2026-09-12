// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Kvit.Ui

// An inline destination or action announced as a link, never as a button.
//
// The optional symbol is caller-supplied. Links do not grow a disclosure
// chevron of their own: whether navigation discloses, leaves the view or opens
// something external belongs to the destination, not to this primitive.
AbstractButton {
    id: root

    property string symbol: ""
    property string role: "body"
    property int elide: Text.ElideRight
    // One sentence saying where following this goes, beside its own words.
    //
    // The same property KvitButton, KvitIconButton and KvitChipButton carry,
    // for the same reason and drawn the same way: shown as the tooltip and
    // announced as the accessible description, so a pointer reader and a
    // screen reader are told the same thing. A link's words are visible
    // already, so this is never where its purpose belongs — it is for a
    // destination the words alone do not describe.
    property string explanation: ""

    signal activated()

    implicitWidth: contentRow.implicitWidth
    implicitHeight: contentRow.implicitHeight
    padding: 0
    hoverEnabled: true
    activeFocusOnTab: true
    clip: true

    Accessible.role: Accessible.Link
    Accessible.name: root.text
    Accessible.description: root.explanation
    Accessible.onPressAction: root.clicked()

    KvitTooltip {
        objectName: "tooltip"
        text: root.explanation
        visible: root.explanation !== ""
                 && (root.hovered || root.visualFocus)
    }

    onClicked: root.activated()

    // AbstractButton handles Space. Links also follow the browser convention
    // and activate on either form of the Enter key.
    Keys.onReturnPressed: event => {
        root.clicked()
        event.accepted = true
    }
    Keys.onEnterPressed: event => {
        root.clicked()
        event.accepted = true
    }

    background: null

    contentItem: Item {
        implicitWidth: contentRow.implicitWidth
        implicitHeight: contentRow.implicitHeight

        RowLayout {
            id: contentRow
            objectName: "contentRow"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
            width: Math.min(implicitWidth, parent.width)
            spacing: Interface.spaceNear

            KvitIcon {
                Layout.alignment: Qt.AlignVCenter
                visible: root.symbol !== ""
                name: root.symbol === "" ? "dot" : root.symbol
                color: root.hovered || root.visualFocus ? Theme.accent
                                                        : Theme.link
                implicitWidth: Interface.iconSizeSmall
                implicitHeight: Interface.iconSizeSmall
            }
            KvitLabel {
                id: label
                objectName: "label"
                Layout.alignment: Qt.AlignVCenter
                Layout.fillWidth: true
                Layout.minimumWidth: 0
                text: root.text
                role: root.role
                color: root.hovered || root.visualFocus ? Theme.accent
                                                        : Theme.link
                font.underline: root.hovered || root.visualFocus
                elide: root.elide
            }
        }
    }

    HoverHandler {
        cursorShape: Qt.PointingHandCursor
    }
}
