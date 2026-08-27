// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Kvit.Ui

// A modal surface that has to be answered before anything else happens.
//
// `KvitDialog` is one of the two names prd.md §5.4 says the estate already
// uses, and this is that component made shared.
//
// Modality is expensive and worth spending only where the answer really does
// have to come first: a destructive action, a choice the next step depends on.
// A confirmation of something reversible belongs in KvitConfirmInPlace or in a
// toast with an undo, both of which let the reader carry on.
//
// The button order is fixed and the confirming button is on the right, which
// is what Qt does on every platform this ships to. A dialog that puts them the
// other way round is a dialog somebody dismisses by muscle memory and destroys
// something.
Dialog {
    id: root

    property string detail: ""
    // Words rather than "OK": a button that says what it does can be read
    // without reading the sentence above it, which is what a reader in a hurry
    // actually does.
    property string confirmText: qsTr("Confirm")
    property string cancelText: qsTr("Cancel")
    property bool destructive: false
    default property alias content: bodySlot.data

    modal: true
    anchors.centerIn: parent ? Overlay.overlay : undefined
    padding: Interface.viewMargin
    closePolicy: Popup.CloseOnEscape
    implicitWidth: Interface.px(420)

    background: Rectangle {
        radius: Interface.radiusCard
        color: Theme.popupBackground
        border.width: Interface.hairline
        border.color: Theme.borderStrong

        // On the background rather than on the Dialog, because a Dialog is a
        // Popup and a Popup is not an Item, which is what Qt's Accessible
        // attached property requires.
        Accessible.role: Accessible.Dialog
        Accessible.name: root.title
        Accessible.description: root.detail
    }

    header: Item {
        implicitHeight: Interface.rowHeightSlim
        KvitLabel {
            anchors.fill: parent
            anchors.leftMargin: Interface.viewMargin
            anchors.rightMargin: Interface.viewMargin
            anchors.topMargin: Interface.viewMargin
            text: root.title
            role: "title"
            font.bold: true
        }
    }

    contentItem: ColumnLayout {
        spacing: Interface.spaceLoose

        KvitLabel {
            Layout.fillWidth: true
            visible: root.detail !== ""
            text: root.detail
            role: "body"
            color: Theme.textSecondary
            wrapMode: Text.WordWrap
            elide: Text.ElideNone
        }
        Item {
            id: bodySlot
            Layout.fillWidth: true
            implicitHeight: childrenRect.height
        }
    }

    footer: Item {
        implicitHeight: Interface.rowHeightSlim + Interface.viewMargin

        Row {
            anchors.right: parent.right
            anchors.rightMargin: Interface.viewMargin
            anchors.verticalCenter: parent.verticalCenter
            spacing: Interface.space

            KvitButton {
                text: root.cancelText
                form: "ordinary"
                onClicked: root.reject()
            }
            KvitButton {
                text: root.confirmText
                form: "primary"
                danger: root.destructive
                // The keyboard's default. A destructive dialog deliberately
                // does not take one: Return should not delete anything, and a
                // reader who hits it expecting to dismiss the dialog would.
                focus: !root.destructive
                onClicked: root.accept()
            }
        }
    }
}
