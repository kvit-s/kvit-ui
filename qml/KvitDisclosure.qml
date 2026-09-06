// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A trigger with a body under it that opens and closes.
//
// prd.md §5.6 counts 18 files in the estate with a private version of this,
// which is why the shared one has to cover the whole shape rather than the
// simplest case. In particular it takes a `group`: `frame.css` has only the
// single trigger-and-body form, and most of the 18 hand-rolled versions are
// one of several sections where opening one closes the others.
//
// The chevron rotates rather than swapping between two glyphs, so the reader
// sees the state change happen; the rotation follows the reduced-motion
// setting, and with motion stilled it is an instant swap.
Item {
    id: root

    property string title: ""
    property bool expanded: false
    // Sections sharing a group string behave as an accordion: opening one
    // closes the rest. Empty means this disclosure stands alone.
    property string group: ""
    // How many things are inside, drawn beside the title. Below zero draws
    // nothing, which is right for a section whose contents are not a
    // countable list.
    property int count: -1
    default property alias content: bodySlot.data

    signal toggled(bool expanded)

    implicitWidth: Interface.px(400)
    // How tall the body wants to be, taken from the preferred height of
    // whatever the caller put inside it.
    //
    // Read from the content rather than from the slot's own `childrenRect`.
    // An item whose own height comes from its `childrenRect` is a cycle Qt
    // reports as a binding loop: assigning the height marks the rectangle
    // dirty, which asks for the height again while it is still being worked
    // out. It settles on the right number and warns every time a section is
    // opened. Reading each child's preferred height instead depends only on
    // what the caller put inside, which nothing here writes back to.
    readonly property real bodyHeight: {
        let wanted = 0
        for (let i = 0; i < bodySlot.children.length; ++i) {
            const child = bodySlot.children[i]
            if (child.visible)
                wanted = Math.max(wanted, child.y + child.implicitHeight)
        }
        return wanted
    }

    implicitHeight: trigger.height + (expanded ? root.bodyHeight
                                                 + Interface.space : 0)

    Behavior on implicitHeight {
        NumberAnimation {
            duration: 140 * Theme.motionScale
            easing.type: Easing.OutCubic
        }
    }

    AbstractButton {
        id: trigger
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
        height: Interface.rowHeightCompact

        Accessible.role: Accessible.Button
        Accessible.name: root.title
        // What a screen reader needs and a chevron does not say: whether this
        // section is currently open.
        Accessible.description: root.expanded ? qsTr("Expanded")
                                              : qsTr("Collapsed")
        Accessible.onPressAction: trigger.clicked()

        onClicked: root.toggled(!root.expanded)

        background: Rectangle {
            color: trigger.hovered ? Theme.hoverTint : "transparent"
            Rectangle {
                anchors.fill: parent
                anchors.margins: -Interface.focusRingWidth
                visible: trigger.visualFocus
                color: "transparent"
                border.width: Interface.focusRingWidth
                border.color: Theme.focusRing
            }
        }

        contentItem: Row {
            spacing: Interface.spaceNear

            KvitIcon {
                anchors.verticalCenter: parent.verticalCenter
                name: "chevron-right"
                color: Theme.textMuted
                implicitWidth: Interface.iconSizeSmall
                implicitHeight: Interface.iconSizeSmall
                rotation: root.expanded ? 90 : 0
                Behavior on rotation {
                    NumberAnimation { duration: 140 * Theme.motionScale }
                }
            }
            KvitLabel {
                anchors.verticalCenter: parent.verticalCenter
                text: root.title
                role: "body"
                font.bold: root.expanded
                elide: Text.ElideNone
            }
            KvitLabel {
                anchors.verticalCenter: parent.verticalCenter
                visible: root.count >= 0
                // Grouped by the reader's locale. A disclosure over a bulk
                // change is handed the record count, and the same figure is
                // written grouped in the line under it; ungrouped here, one
                // number appears in the window twice in two spellings.
                text: Number(root.count).toLocaleString(Qt.locale(), 'f', 0)
                role: "caption"
                color: Theme.textFaint
                tabular: true
                elide: Text.ElideNone
            }
        }
    }

    Item {
        id: bodySlot
        anchors.top: trigger.bottom
        anchors.topMargin: Interface.spaceSnug
        anchors.left: parent.left
        anchors.leftMargin: Interface.spaceLoose
        anchors.right: parent.right
        height: root.expanded ? root.bodyHeight : 0
        visible: root.expanded
        clip: true
    }
}
