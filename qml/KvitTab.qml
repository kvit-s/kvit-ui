// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// One tab in a row of them: a different view of the same thing.
//
// The selected tab is marked by an underline as well as by colour and weight,
// which is the second channel the high-contrast pass looks for. Selection
// shown by colour alone is the most common place this rule gets broken,
// because a coloured tab label looks obviously different to whoever drew it.
AbstractButton {
    id: root

    property bool selected: false
    // A count beside the label — how many things this tab holds.
    property int count: -1

    implicitHeight: Math.max(Interface.tabHeight, contentRow.implicitHeight)
    implicitWidth: contentRow.implicitWidth + Interface.spaceLoose * 2

    // The count with the reader's own digit grouping, said and drawn the same
    // way. Ungrouped, a five-figure count is read out digit by digit and is
    // spelled differently from the same number everywhere else in the window.
    readonly property string groupedCount:
        Number(root.count).toLocaleString(Qt.locale(), 'f', 0)

    Accessible.role: Accessible.PageTab
    // Two source strings picked here by the count rather than one with `%n` in
    // it: with no translator installed Qt substitutes the number and does not
    // choose a plural form, so a tab holding one thing announces "1 item(s)".
    // The count still goes to qsTr, so which form to use stays Qt's and a
    // translator into a language with three of them gets all three.
    Accessible.name: {
        if (root.count < 0)
            return root.text
        return root.count === 1
            ? qsTr("%1, %2 item", "a label and a count of one thing",
                   root.count).arg(root.text).arg(root.groupedCount)
            : qsTr("%1, %2 items", "a label and a count of several things",
                   root.count).arg(root.text).arg(root.groupedCount)
    }
    Accessible.selected: root.selected
    Accessible.onPressAction: root.clicked()

    background: Rectangle {
        color: root.hovered && !root.selected ? Theme.hoverTint : "transparent"

        // The underline: two design pixels, so it survives at the smallest
        // interface size and does not read as a hairline rule.
        Rectangle {
            anchors.bottom: parent.bottom
            anchors.left: parent.left
            anchors.right: parent.right
            height: Interface.spaceTight
            visible: root.selected
            color: Theme.accent
        }

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.visualFocus
            color: "transparent"
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    contentItem: Item {
        implicitWidth: contentRow.implicitWidth
        implicitHeight: contentRow.implicitHeight

        Row {
            id: contentRow
            objectName: "contentRow"
            anchors.horizontalCenter: parent.horizontalCenter
            // Native text must be centred deliberately on both rendering
            // paths; relying on the Row's assigned height shifts on Windows.
            anchors.verticalCenter: parent.verticalCenter
            spacing: Interface.spaceNear

            KvitLabel {
                anchors.verticalCenter: parent.verticalCenter
                text: root.text
                role: "body"
                font.bold: root.selected
                color: root.selected ? Theme.textPrimary : Theme.textMuted
                elide: Text.ElideNone
            }
            KvitLabel {
                anchors.verticalCenter: parent.verticalCenter
                visible: root.count >= 0
                text: root.groupedCount
                role: "caption"
                color: Theme.textFaint
                tabular: true
                elide: Text.ElideNone
            }
        }
    }
}
