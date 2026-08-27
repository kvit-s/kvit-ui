// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// The strip across the top of the window: what this application is, where the
// reader is in it, and the actions that apply everywhere.
//
// Three slots in a fixed order — the wordmark on the left, navigation in the
// middle, actions on the right — because a reader who has learned where the
// back control is on one screen should find it in the same place on every
// screen and in every application in the estate.
KvitPanel {
    id: root

    property string wordmark: ""
    default property alias navigation: navigationSlot.data
    property alias actions: actionSlot.data

    implicitHeight: Interface.headerHeight
    ruleBottom: true

    Accessible.role: Accessible.PageTabList
    Accessible.name: qsTr("Application header")

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: Interface.viewMargin
        anchors.rightMargin: Interface.viewMargin
        spacing: Interface.columnGap

        KvitLabel {
            visible: root.wordmark !== ""
            text: root.wordmark
            role: "headline"
            font.bold: true
            elide: Text.ElideNone
        }

        // The navigation slot takes the space left over rather than sizing
        // itself from what is in it. Deriving its width from `childrenRect`
        // while its child fills it is a loop that hangs on construction
        // instead of warning.
        Item {
            id: navigationSlot
            Layout.fillWidth: true
            Layout.fillHeight: true
        }

        Row {
            id: actionSlot
            Layout.alignment: Qt.AlignVCenter
            spacing: Interface.spaceNear
        }
    }
}
