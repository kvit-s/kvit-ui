// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A button with words on it, in three forms.
//
//   primary   filled with the accent. One per screen region: it is the
//             action the screen is for, and a screen with three primary
//             buttons has told the reader nothing about which to press.
//   ordinary  outlined. Everything else that is an action.
//   quiet     no ground until hover. For a dense strip where a row of
//             outlines would read as a fence.
//
// `danger` is separate from the form rather than a fourth form, because a
// destructive action can be any of the three: a primary Delete in a
// confirmation dialog, an ordinary one in a settings row, a quiet one in a
// context menu.
//
// The label colour on a filled button comes from Theme.labelOn rather than
// from onAccent, because onAccent answers only for the accent and a filled
// danger button needs the label that contrasts with its own fill
// (accessibility.md Finding 3).
AbstractButton {
    id: root

    // "primary" | "ordinary" | "quiet"
    property string form: "ordinary"
    property bool danger: false
    property string symbol: ""
    // A button that is doing something. It stays enabled so its accessible
    // name can say what is happening, and the label changes rather than a
    // spinner appearing, because a spinner says "wait" without saying for what.
    property bool busy: false
    property string busyText: qsTr("Working…")
    // One sentence saying what the words on the button cannot: why it is
    // disabled, or what pressing it opens.
    //
    // A button names itself from its own text, which is the right default and
    // is all most buttons need. The case this is for is the one where the
    // reason a control is in the state it is in lives somewhere the reader
    // cannot see — a Pull button that is disabled because the remote has not
    // been fetched, an Archive button that is disabled because the branch has
    // unpushed work. Without somewhere to put that sentence, a screen offers
    // a grey control and no account of it.
    //
    // It is shown as the tooltip and announced as the accessible description,
    // so a pointer reader and a screen reader are told the same thing. It is
    // never where the label belongs: a button whose purpose is only in its
    // explanation cannot be used without hovering it, which rules out
    // everybody on a touch screen or a keyboard.
    //
    // A disabled button still shows it. Qt stops sending hover events to a
    // disabled item, so the words reach a pointer reader through the button's
    // own HoverHandler below rather than through `hovered`.
    property string explanation: ""

    readonly property color fill: root.danger ? Theme.danger : Theme.accent

    implicitHeight: Math.max(Interface.controlHeight, contentRow.implicitHeight)
    // Measure the words and symbol, not the content item's assigned width.
    // Control stretches its content item across the available rectangle; a
    // Row used directly therefore reports the right natural width but lays
    // its children out at the stretched item's left edge.
    implicitWidth: contentRow.implicitWidth + Interface.spaceLoose * 2
    padding: Interface.space
    font.pixelSize: Interface.body

    Accessible.role: Accessible.Button
    Accessible.name: root.busy ? root.busyText : root.text
    Accessible.description: root.explanation
    Accessible.onPressAction: root.clicked()

    // A HoverHandler rather than the control's own `hovered`, for the one
    // case that matters: `hovered` is false on a disabled button, and a
    // disabled button is exactly where the explanation is worth reading. A
    // handler keeps receiving the pointer either way.
    HoverHandler {
        id: explanationHover
        enabled: root.explanation !== ""
    }

    KvitTooltip {
        objectName: "tooltip"
        text: root.explanation
        visible: root.explanation !== ""
                 && (explanationHover.hovered || root.visualFocus)
    }

    background: Rectangle {
        radius: Interface.radiusControl
        color: {
            if (!root.enabled)
                return root.form === "primary" ? Theme.chipBackground : "transparent"
            if (root.form === "primary")
                return root.pressed ? Qt.darker(root.fill, 1.15) : root.fill
            if (root.pressed)
                return Theme.selectionTint
            if (root.hovered)
                return Theme.hoverTint
            return "transparent"
        }
        border.width: root.form === "ordinary" ? Interface.hairline : 0
        border.color: root.enabled
            ? (root.danger ? Theme.danger : Theme.borderStrong)
            : Theme.border

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.visualFocus
            color: "transparent"
            radius: parent.radius + Interface.focusRingWidth
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    contentItem: Item {
        id: content
        implicitWidth: contentRow.implicitWidth
        implicitHeight: contentRow.implicitHeight

        readonly property color foreground: {
            if (!root.enabled)
                return Theme.textDisabled
            if (root.form === "primary")
                return Theme.labelOn(root.fill)
            return root.danger ? Theme.danger : Theme.textPrimary
        }

        Row {
            id: contentRow
            objectName: "contentRow"
            anchors.horizontalCenter: parent.horizontalCenter
            // Keep this explicit. Native-rendered text can otherwise settle
            // on a different baseline on Windows than it does on Linux.
            anchors.verticalCenter: parent.verticalCenter
            spacing: Interface.spaceNear

            KvitIcon {
                anchors.verticalCenter: parent.verticalCenter
                visible: root.symbol !== ""
                name: root.symbol === "" ? "dot" : root.symbol
                color: content.foreground
                implicitWidth: Interface.iconSizeSmall
                implicitHeight: Interface.iconSizeSmall
            }
            KvitLabel {
                anchors.verticalCenter: parent.verticalCenter
                text: root.busy ? root.busyText : root.text
                role: "body"
                color: content.foreground
                elide: Text.ElideNone
            }
        }
    }
}
