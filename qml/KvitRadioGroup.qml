// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// One choice from a few, all visible at once, each with room for a sentence
// saying what it means.
//
// The alternative shapes and when each is right: KvitSegmented for two to four
// short options where the labels are self-explanatory, KvitSelect for a list
// too long to lay out, and this for a handful of options where the choice
// needs explaining — "Follow the system" beside "Always dark" beside "Always
// light", each with a line under it.
//
// Arrow keys move within the group and Tab leaves it, which is what a screen
// reader user expects of a radio group and what a column of separate controls
// does not do.
Column {
    id: root

    // A list of { value, label, detail } objects.
    property var options: []
    property var current: undefined
    property string label: ""

    signal chosen(var value)

    spacing: Interface.spaceNear

    Accessible.role: Accessible.Grouping
    Accessible.name: root.label

    ButtonGroup { id: exclusive }

    Repeater {
        model: root.options
        delegate: RadioButton {
            id: option
            required property var modelData
            required property int index

            width: root.width
            checked: root.current !== undefined
                     && modelData.value === root.current
            ButtonGroup.group: exclusive
            spacing: Interface.space

            Accessible.role: Accessible.RadioButton
            Accessible.name: modelData.detail !== undefined && modelData.detail !== ""
                ? qsTr("%1. %2").arg(modelData.label).arg(modelData.detail)
                : modelData.label

            onClicked: root.chosen(modelData.value)

            indicator: Rectangle {
                implicitWidth: Interface.px(16)
                implicitHeight: Interface.px(16)
                x: option.leftPadding
                y: Interface.spaceTight
                radius: width / 2
                color: "transparent"
                border.width: Interface.hairline
                border.color: option.enabled ? Theme.borderStrong : Theme.border

                Rectangle {
                    anchors.centerIn: parent
                    visible: option.checked
                    width: parent.width - Interface.spaceNear
                    height: width
                    radius: width / 2
                    color: Theme.accent
                }
                Rectangle {
                    anchors.fill: parent
                    anchors.margins: -Interface.focusRingWidth
                    visible: option.visualFocus
                    color: "transparent"
                    radius: parent.radius + Interface.focusRingWidth
                    border.width: Interface.focusRingWidth
                    border.color: Theme.focusRing
                }
            }

            contentItem: Column {
                x: option.indicator.width + option.spacing
                spacing: Interface.spaceTight

                KvitLabel {
                    text: option.modelData.label
                    role: "body"
                    color: option.enabled ? Theme.textPrimary : Theme.textDisabled
                    elide: Text.ElideNone
                }
                KvitLabel {
                    width: root.width - option.indicator.width - option.spacing
                    visible: option.modelData.detail !== undefined
                             && option.modelData.detail !== ""
                    text: option.modelData.detail === undefined
                          ? "" : option.modelData.detail
                    role: "small"
                    color: Theme.textMuted
                    wrapMode: Text.WordWrap
                    elide: Text.ElideNone
                }
            }
        }
    }
}
