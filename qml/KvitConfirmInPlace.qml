// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// A strip that appears where the reader is, saying what just happened and
// offering to undo it.
//
// prd.md §5.6 Group B, and kvit-cash's bulk edit is the flow that names it:
// the undo goes inside the confirmation rather than the confirmation coming
// before the action.
//
// The reason to prefer this over a modal dialog is what each costs when the
// reader is doing the same thing forty times. A dialog before every action
// costs a click each time and stops being read by the third one; a
// confirmation after, with an undo, costs nothing unless something went wrong.
// Which means the undo has to actually work — this shape is only honest for an
// action the application can really reverse.
//
// It stays until it is dismissed or replaced. A timed one would be a toast,
// and an undo that expires while the reader is checking what happened is an
// undo they cannot use.
Rectangle {
    id: root

    property string text: ""
    // The words on the undo. Empty draws no undo at all, which is what a
    // caller whose action has stopped being reversible sets: the strip then
    // states what happened without offering a control that would be refused.
    property string undoText: qsTr("Undo")
    property bool shown: false
    // What was affected, so the reader can check the count before undoing.
    property int affected: -1

    signal undone()
    signal dismissed()

    implicitHeight: shown ? Interface.rowHeightSlim + Interface.space : 0
    implicitWidth: Interface.px(400)

    visible: implicitHeight > 0
    clip: true
    Behavior on implicitHeight {
        NumberAnimation {
            duration: 140 * Theme.motionScale
            easing.type: Easing.OutCubic
        }
    }

    radius: Interface.radiusControl
    color: Qt.alpha(Theme.success, 0.12)
    border.width: Interface.hairline
    border.color: Theme.success

    // "1 item", "250,000 items" — built once and used by both the drawn label
    // and the announcement, so the two cannot disagree about the same number.
    //
    // The digits are grouped by the reader's locale, and the plural is a word
    // rather than a parenthesis: "250000 item(s)" is a form field rather than
    // a sentence. The number goes in through %1 rather than through Qt's %n,
    // which substitutes the bare integer and would lose the grouping; the
    // count is still handed to qsTr, so the choice of plural form stays Qt's
    // and a translator gets every form.
    readonly property string affectedPhrase: {
        if (root.affected < 0)
            return ""
        const grouped = Number(root.affected).toLocaleString(Qt.locale(), 'f', 0)
        return root.affected === 1
            ? qsTr("%1 item", "a count of one thing", root.affected).arg(grouped)
            : qsTr("%1 items", "a count of several things",
                   root.affected).arg(grouped)
    }

    Accessible.role: Accessible.AlertMessage
    Accessible.name: root.affected >= 0
        ? qsTr("%1, %2 affected", "what was done and how much it touched")
              .arg(root.text).arg(root.affectedPhrase)
        : root.text

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: Interface.spaceLoose
        anchors.rightMargin: Interface.spaceNear
        spacing: Interface.space

        KvitIcon {
            name: "success"
            color: Theme.success
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }
        KvitLabel {
            Layout.fillWidth: true
            text: root.affected >= 0
                  ? qsTr("%1 — %2", "what was done and how much it touched")
                        .arg(root.text).arg(root.affectedPhrase)
                  : root.text
            role: "body"
        }
        KvitButton {
            visible: root.undoText !== ""
            text: root.undoText
            form: "quiet"
            onClicked: root.undone()
        }
        KvitIconButton {
            symbol: "close"
            label: qsTr("Dismiss")
            onClicked: root.dismissed()
        }
    }
}
