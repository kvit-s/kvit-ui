// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A label a person put there: a tag, a category, a folder.
//
// A chip is chrome saying what something is; a tag is content the reader wrote.
// The difference matters because a tag carries a colour the reader chose, from
// Theme.colorPalette, which is theme-independent — a red tag is the same red in
// dark mode, because the reader picked that red and would not expect it to
// change.
//
// Which means the label colour cannot be a token: it has to contrast with
// whatever the reader picked, and Theme.labelOn is what answers that.
Rectangle {
    id: root

    property string text: ""
    // A value from Theme.colorPalette. Empty takes the neutral chip ground,
    // which is what an untinted tag looks like.
    property string tint: ""
    property bool removable: false

    signal removeRequested()

    readonly property bool tinted: root.tint !== ""
    readonly property color tintColor: root.tinted ? root.tint : Theme.chipBackground

    implicitHeight: Interface.tagHeight
    implicitWidth: row.width + Interface.spaceNear * 2
    radius: Interface.radiusPill

    color: root.tinted ? Qt.alpha(root.tintColor, 0.22) : Theme.chipBackground
    border.width: Interface.hairline
    border.color: root.tinted ? root.tintColor : Theme.border

    Accessible.role: Accessible.StaticText
    // The colour's name, not its hex value: a screen reader saying "#8a5cc0"
    // tells the reader nothing (accessibility.md Finding 2).
    Accessible.name: root.tinted
        ? qsTr("%1, %2").arg(root.text).arg(Theme.colorName(root.tint))
        : root.text

    Row {
        id: row
        anchors.centerIn: parent
        spacing: Interface.spaceSnug

        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            text: root.text
            role: "caption"
            color: root.tinted ? Theme.textPrimary : Theme.textSecondary
            elide: Text.ElideNone
        }
        KvitIconButton {
            anchors.verticalCenter: parent.verticalCenter
            visible: root.removable
            symbol: "close"
            label: qsTr("Remove %1").arg(root.text)
            implicitWidth: Interface.caption
            implicitHeight: Interface.caption
            onClicked: root.removeRequested()
        }
    }
}
