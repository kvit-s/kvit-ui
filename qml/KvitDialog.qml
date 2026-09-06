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

    // The title, with the way out at the right end of the same row — where
    // KvitPane puts its close control, so a reader who has learnt one surface
    // already knows where the other one's is.
    //
    // It is a third route out rather than a replacement for either of the
    // other two: `closePolicy` is untouched, so Escape still dismisses the
    // dialog, and the cancel button is still there. The reason to have all
    // three is that they are found by different readers — the pointer finds
    // the symbol, the keyboard finds Escape, and somebody reading the
    // sentence finds the button that says what declining does.
    //
    // The header is fourteen pixels taller than it was, and the title's right
    // margin is ten narrower. That is a deliberate change to what every
    // dialog in the estate draws, so here is the arithmetic behind it.
    //
    // The header used to be `Interface.rowHeightSlim` — thirty pixels at the
    // default interface size — with the title held sixteen below the top by
    // `viewMargin`, which left fourteen. `Interface.controlHeight` is
    // twenty-eight, so a close control in the old header either loses the
    // margin above the title, moving every dialog title up twelve pixels, or
    // is squeezed into fourteen pixels and becomes a hit target half the size
    // of every other icon button. Keeping the margin and growing the header
    // to `viewMargin + controlHeight` is the third option, and it is the one
    // that leaves the title where a reader is used to finding it.
    //
    // The right margin is `spaceNear` rather than `viewMargin` because the
    // icon button carries its own padding around the symbol; an equal margin
    // would set the symbol further from the edge than the title is from the
    // other one. These are the two margins KvitPane's head already uses, so
    // the close control sits in the same place on both surfaces.
    header: Item {
        implicitHeight: Interface.viewMargin + Interface.controlHeight

        RowLayout {
            anchors.fill: parent
            anchors.leftMargin: Interface.viewMargin
            anchors.rightMargin: Interface.spaceNear
            anchors.topMargin: Interface.viewMargin
            spacing: Interface.space

            KvitLabel {
                Layout.fillWidth: true
                Layout.alignment: Qt.AlignVCenter
                text: root.title
                role: "title"
                font.bold: true
            }
            KvitIconButton {
                Layout.alignment: Qt.AlignVCenter
                symbol: "close"
                // One word, and it is both the tooltip a pointer reader sees
                // and the name a screen reader says, which is what
                // KvitIconButton's `label` is for.
                label: qsTr("Close")
                // The same answer Escape and the cancel button give: a
                // dismissal is a decline, not an unstated agreement.
                onClicked: root.reject()
            }
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
