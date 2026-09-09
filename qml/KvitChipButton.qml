// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import Kvit.Ui

// KvitChip's twin for a fact that opens something when it is pressed.
//
// Two components rather than one property on KvitChip, for two reasons. A
// KvitChip is a Rectangle with two items in it and a table draws several
// hundred of them at once; a Control carries a background, a content item and
// the whole AbstractButton machinery, and paying for that on every chip in a
// grid to serve the handful that act is the wrong trade. And the library's
// rule is that a Rectangle with a MouseArea in it is not a control — the way
// to make a chip act is to make it a real AbstractButton, which is what Qt
// publishes with a button role, puts in the tab order and hands Space to,
// rather than to bolt a handler onto the mark.
//
// At rest it is drawn exactly as KvitChip is, from the same tone table, so a
// row mixing facts that act with facts that do not reads as one row of facts.
// What separates them is the behaviour a control has anyway: the ground
// changes under the pointer, the cursor becomes a hand, the focus ring
// appears when the keyboard reaches it, and a screen reader is told it is a
// button. No chevron is drawn for it. Whether pressing this discloses a list,
// leaves the view or opens another window belongs to the destination rather
// than to the mark, which is the same reason KvitLink draws none;
// `trailingSymbol` is there for a caller that knows which of those it is.
AbstractButton {
    id: root

    // "neutral" | "accent" | "success" | "warning" | "danger" | "info"
    property string tone: "neutral"
    // A filled chip rather than a tinted one, for the one chip on a row that
    // is the point of the row.
    property bool strong: false
    property string symbol: ""
    // A symbol after the label, where the caller knows what pressing this
    // leads to: a chevron for a list that opens beneath, an out-arrow for
    // something that opens elsewhere.
    property string trailingSymbol: ""

    // Why this chip cannot be pressed, in the reader's own terms. Empty means
    // it can be.
    //
    // This rather than `enabled: false`, because Qt takes a disabled item out
    // of the tab order and stops sending it hover events, which would make
    // the one chip on the row that has something to explain the one chip a
    // reader cannot reach to hear the explanation. A chip with a reason keeps
    // its place in the tab order, still shows its tooltip, still announces
    // itself, and emits nothing when it is pressed.
    property string unavailableReason: ""
    readonly property bool unavailable: root.unavailableReason !== ""

    readonly property color toneColor: {
        switch (tone) {
        case "accent":  return Theme.accent
        case "success": return Theme.success
        case "warning": return Theme.warning
        case "danger":  return Theme.danger
        case "info":    return Theme.link
        default:        return Theme.textMuted
        }
    }

    // Under the pointer or under the keyboard. The two draw the same ground
    // because they mean the same thing about this chip — it is the one that
    // would act — and they are told apart by the ring, which only the
    // keyboard draws.
    readonly property bool lit: root.hovered || root.visualFocus

    signal activated()

    implicitHeight: Math.max(Interface.chipHeight, contentRow.implicitHeight)
    implicitWidth: contentRow.implicitWidth + root.leftPadding + root.rightPadding
    // The chip's own inset, in the control's padding rather than in the
    // content item's position: a Control places its content item at
    // (leftPadding, topPadding) and overwrites any x it was given, and it is
    // the padding that makes `availableWidth` the width the label really has
    // to elide into.
    leftPadding: Interface.spaceNear
    rightPadding: Interface.spaceNear
    topPadding: 0
    bottomPadding: 0
    hoverEnabled: true
    activeFocusOnTab: true

    Accessible.role: Accessible.Button
    Accessible.name: root.text
    Accessible.description: root.unavailableReason
    Accessible.onPressAction: root.clicked()

    onClicked: {
        if (!root.unavailable)
            root.activated()
    }

    // AbstractButton answers Space on its own. Return and Enter are added for
    // the reason KvitLink adds them: what these chips do is open something,
    // and a reader who has learnt that Return opens the row does not learn a
    // different key for the fact sitting on it.
    Keys.onReturnPressed: event => {
        root.clicked()
        event.accepted = true
    }
    Keys.onEnterPressed: event => {
        root.clicked()
        event.accepted = true
    }

    // Two things need saying and neither is always true, so one surface says
    // whichever applies: why the chip cannot be pressed, or the whole label
    // when the column was too narrow to draw it. An available chip whose
    // label fits has nothing to add and shows nothing.
    KvitTooltip {
        objectName: "tooltip"
        text: root.unavailable ? root.unavailableReason : root.text
        visible: (root.hovered || root.visualFocus)
                 && (root.unavailable || label.truncated)
    }

    HoverHandler {
        cursorShape: root.unavailable ? Qt.ArrowCursor : Qt.PointingHandCursor
    }

    background: Rectangle {
        radius: Interface.radiusChip

        // The unavailable state is the fill going away, not the fill changing
        // hue: a tinted or filled shape becomes an empty outline, which is a
        // difference a reader who cannot separate the tones still sees. The
        // words are in the tooltip and in the accessible description.
        color: {
            if (root.unavailable)
                return "transparent"
            if (root.strong)
                return root.pressed ? Qt.darker(root.toneColor, 1.15)
                                    : root.toneColor
            if (root.tone === "neutral")
                return (root.pressed || root.lit) ? Theme.hoverTint
                                                  : Theme.chipBackground
            return Qt.alpha(root.toneColor,
                            (root.pressed || root.lit) ? 0.28 : 0.16)
        }
        border.width: (root.strong && !root.unavailable) ? 0 : Interface.hairline
        border.color: root.unavailable ? Theme.border
                    : root.tone === "neutral"
                        ? (root.lit ? Theme.borderStrong : Theme.border)
                        : root.toneColor

        // Outside the fill rather than a border colour change, because a ring
        // drawn as a border disappears on a chip that already has one.
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
        implicitWidth: contentRow.implicitWidth
        implicitHeight: contentRow.implicitHeight

        // A layout rather than a Row, because the label is the piece that
        // gives way: at a width narrower than the chip wants — a narrow
        // column, or the same column at twice the interface size — the two
        // symbols keep their size and the words elide between them, which is
        // the one arrangement that stays readable rather than losing the
        // symbol that says what pressing it does.
        RowLayout {
            id: contentRow
            objectName: "contentRow"
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.verticalCenter: parent.verticalCenter
            width: Math.min(implicitWidth, parent.width)
            spacing: Interface.spaceSnug

            KvitIcon {
                Layout.alignment: Qt.AlignVCenter
                visible: root.symbol !== ""
                name: root.symbol === "" ? "dot" : root.symbol
                color: label.color
                implicitWidth: Interface.caption
                implicitHeight: Interface.caption
            }
            KvitLabel {
                id: label
                objectName: "label"
                Layout.alignment: Qt.AlignVCenter
                Layout.fillWidth: true
                Layout.minimumWidth: 0
                text: root.text
                role: "caption"
                // labelOn rather than onAccent: onAccent answers only for the
                // accent, and a filled danger chip needs the label that
                // contrasts with *its* fill (accessibility.md Finding 3).
                color: root.unavailable ? Theme.textDisabled
                     : root.strong ? Theme.labelOn(root.toneColor)
                     : root.tone === "neutral" ? Theme.textSecondary
                                               : root.toneColor
                elide: Text.ElideRight
            }
            KvitIcon {
                Layout.alignment: Qt.AlignVCenter
                visible: root.trailingSymbol !== ""
                name: root.trailingSymbol === "" ? "dot" : root.trailingSymbol
                color: label.color
                implicitWidth: Interface.caption
                implicitHeight: Interface.caption
            }
        }
    }
}
