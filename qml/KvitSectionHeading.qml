// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// A group heading: a solid bar carrying a disclosure chevron, the group's
// name, what it holds, and the one action that takes the whole group.
//
// kvit-hub's SectionHeading, with its reasoning intact.
//
// It is a filled bar rather than a label with a rule running off to the right,
// because a rule drawn between the words and the action reads as a divider
// between two unrelated things when it is meant to bind them into one heading.
//
// The action is hoisted here on purpose, and it has to be true for every row
// under the heading. If it is not true for all of them the group is wrong and
// wants splitting — which is the test that keeps a "hand all four to an agent"
// button from appearing over a group where it only applies to three.
//
// The count carries the word for what it counts, because a bare number at the
// end of a heading does not say what was counted, and a group of one is not
// counted at all: the single row under the heading is the count.
Item {
    id: root

    property string text: ""
    property int count: -1
    // What the count counts, singular. "" leaves the number bare.
    property string counted: ""
    // The kind of thing the group holds, said beside the name where the name
    // alone does not say it.
    property string kind: ""
    // The hoisted action, true for every row in the group. Empty for none.
    property string action: ""
    property bool strong: false
    property bool expanded: true
    property bool collapsible: false

    signal toggled()
    signal actioned()

    implicitHeight: Interface.rowHeightCompact
    implicitWidth: Interface.px(400)

    Rectangle {
        anchors.fill: parent
        color: root.collapsible && headHover.hovered ? Theme.hoverTint
                                                     : Theme.panelBackground

        KvitDivider {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.bottom: parent.bottom
            color: root.strong ? Theme.borderStrong : Theme.border
        }

        HoverHandler {
            id: headHover
            enabled: root.collapsible
            cursorShape: Qt.PointingHandCursor
        }
        TapHandler {
            enabled: root.collapsible
            onTapped: root.toggled()
        }
    }

    Accessible.role: root.collapsible ? Accessible.Heading : Accessible.StaticText
    Accessible.name: root.text

    // A layout rather than a row of natural widths: the name, what the group
    // holds and the count are of unknown length together, and left to
    // themselves they run under the action at the right edge. The middle piece
    // gives way first, since it is the one a reader can lose without losing
    // which group this is.
    RowLayout {
        anchors.left: parent.left
        anchors.leftMargin: Interface.space
        anchors.right: parent.right
        anchors.rightMargin: Interface.space
        anchors.verticalCenter: parent.verticalCenter
        spacing: Interface.stackGap

        KvitIcon {
            visible: root.collapsible
            name: root.expanded ? "chevron-down" : "chevron-right"
            color: Theme.textFaint
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }
        KvitLabel {
            text: root.text
            role: "small"
            font.bold: root.strong
        }
        KvitLabel {
            Layout.fillWidth: true
            visible: root.kind !== ""
            text: root.kind
            role: "small"
            color: Theme.textMuted
        }
        Item { Layout.fillWidth: root.kind === "" }
        KvitLabel {
            visible: root.count > 1
            text: root.counted !== ""
                  ? qsTr("%1 %2s").arg(root.count).arg(root.counted)
                  : String(root.count)
            role: "small"
            color: Theme.textFaint
            tabular: true
        }
        KvitLink {
            Layout.leftMargin: Interface.stackGap
            visible: root.action !== ""
            text: root.action
            role: "small"
            onActivated: root.actioned()
        }
    }
}
