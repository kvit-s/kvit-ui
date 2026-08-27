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

// One choice from two to five short options, all visible.
//
// prd.md §5.6 puts this in Group B — absent everywhere and named in a written
// flow. kvit-cash's dashboard needs one: a period control (week, month,
// quarter, year) governing seven widgets at once. That is the shape this is
// for, and it is why the options are visible rather than behind a dropdown —
// a control that governs the whole screen should say what the screen is
// currently showing without being opened.
//
// Past about five options it stops working: the segments get too narrow for
// their labels and KvitSelect is the right component.
Rectangle {
    id: root

    // A list of strings, or of { value, label } objects.
    property var options: []
    property var current: undefined
    property string label: ""

    signal chosen(var value)

    implicitHeight: Interface.controlHeight
    implicitWidth: row.implicitWidth + Interface.hairline * 2

    radius: Interface.radiusControl
    color: Theme.chipBackground
    border.width: Interface.hairline
    border.color: Theme.borderStrong

    Accessible.role: Accessible.PageTabList
    Accessible.name: root.label

    function valueOf(option) {
        return option !== null && typeof option === "object" ? option.value : option
    }
    function labelOf(option) {
        return option !== null && typeof option === "object" ? option.label : option
    }

    Row {
        id: row
        anchors.fill: parent
        anchors.margins: Interface.hairline

        Repeater {
            model: root.options
            delegate: AbstractButton {
                id: segment
                required property var modelData
                required property int index

                readonly property bool selected:
                    root.valueOf(modelData) === root.current

                // `row` by id rather than `parent`, which is null while the
                // delegate is being built.
                height: row.height
                width: Math.max(implicitContentWidth + Interface.spaceLoose * 2,
                                Interface.px(56))

                Accessible.role: Accessible.PageTab
                Accessible.name: root.labelOf(modelData)
                Accessible.selected: selected
                Accessible.onPressAction: segment.clicked()

                onClicked: root.chosen(root.valueOf(modelData))

                background: Rectangle {
                    radius: Interface.radiusBar
                    // The selected segment is filled and its label goes bold,
                    // which is the second channel: a segmented control marked
                    // only by a tint is unreadable in high contrast, where
                    // every tint is nearly the same.
                    color: segment.selected ? Theme.popupBackground
                         : segment.hovered ? Theme.hoverTint
                         : "transparent"
                    border.width: segment.selected ? Interface.hairline : 0
                    border.color: Theme.borderStrong

                    Rectangle {
                        anchors.fill: parent
                        anchors.margins: -Interface.focusRingWidth
                        visible: segment.visualFocus
                        color: "transparent"
                        radius: parent.radius + Interface.focusRingWidth
                        border.width: Interface.focusRingWidth
                        border.color: Theme.focusRing
                    }
                }

                contentItem: KvitLabel {
                    text: root.labelOf(segment.modelData)
                    role: "body"
                    font.bold: segment.selected
                    color: segment.selected ? Theme.textPrimary : Theme.textMuted
                    horizontalAlignment: Text.AlignHCenter
                    elide: Text.ElideNone
                }
            }
        }
    }
}
