// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// The shape of a series, small enough to sit in a row.
//
// It has no axis and no labels, which is the whole point: it answers "is this
// going up" at a glance, and a reader who needs the numbers has KvitTrend.
//
// The rule it carries: a period that was never measured is drawn as a baseline
// tick rather than as a zero-height bar. A missing week and a week with no
// activity are different facts, and drawing both as nothing makes a series
// look continuous when it has a hole in it. A `null` value in `values` is
// what says "not measured".
Item {
    id: root

    // Numbers, with null for a period that was not measured.
    property var values: []
    property color color: Theme.accent
    property string label: ""
    // The scale. Zero uses the largest value present, which is right for a
    // spark whose job is the shape rather than the magnitude.
    property real maximum: 0

    implicitHeight: Interface.rowHeightCompact
    implicitWidth: Interface.px(80)

    readonly property real scale: {
        if (maximum > 0)
            return maximum
        let top = 0
        for (let i = 0; i < values.length; ++i) {
            if (values[i] !== null && values[i] > top)
                top = values[i]
        }
        return top > 0 ? top : 1
    }

    Accessible.role: Accessible.Chart
    Accessible.name: root.label
    Accessible.description: qsTr("%n period(s)", "", root.values.length)

    Row {
        anchors.fill: parent
        spacing: Interface.hairline

        Repeater {
            model: root.values
            delegate: Item {
                required property var modelData
                required property int index
                width: Math.max(Interface.hairline,
                                (root.width - Interface.hairline
                                 * Math.max(0, root.values.length - 1))
                                / Math.max(1, root.values.length))
                height: root.height

                // A measured value: a bar of its height.
                Rectangle {
                    visible: parent.modelData !== null
                    anchors.bottom: parent.bottom
                    width: parent.width
                    height: parent.modelData === null ? 0
                        : Math.max(Interface.hairline,
                                   Math.round(root.height
                                              * (parent.modelData / root.scale)))
                    radius: Interface.radiusBar
                    color: root.color
                }

                // A period nobody measured: a tick on the baseline. Muted and
                // one pixel tall, so it reads as an absence rather than as a
                // very small value.
                Rectangle {
                    visible: parent.modelData === null
                    anchors.bottom: parent.bottom
                    width: parent.width
                    height: Interface.hairline
                    color: Theme.textFaint
                }
            }
        }
    }
}
