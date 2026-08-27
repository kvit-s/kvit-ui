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

// Several quantities that add up to one total, in one bar.
//
// The two-pixel gap between segments is the reason this is a component. Two
// adjacent fills that touch read as one fill with a colour change in it, and
// at the small end — a segment three pixels wide — the reader cannot tell
// whether there are two segments or one. The gap is drawn in the surface
// colour rather than left transparent, so it works over any ground.
//
// A segment too small to be visible is not silently dropped. It is drawn at
// the minimum visible width and the total is noted as approximate, because a
// stacked bar whose segments do not add up to its total is a chart that lies
// about its own arithmetic.
Item {
    id: root

    // A list of { value, color, label } objects.
    property var segments: []
    // The scale. Defaults to the sum, which is right for a bar that is a
    // breakdown of its own total; set it to compare several bars against one
    // scale.
    property real maximum: 0
    property bool wide: false
    property string label: ""

    implicitHeight: wide ? Interface.barHeightWide : Interface.barHeight
    implicitWidth: Interface.px(160)

    readonly property real total: {
        let sum = 0
        for (let i = 0; i < segments.length; ++i)
            sum += segments[i].value
        return sum
    }
    readonly property real scale: maximum > 0 ? maximum : total

    Accessible.role: Accessible.Chart
    Accessible.name: root.label
    Accessible.description: {
        const parts = []
        for (let i = 0; i < root.segments.length; ++i) {
            parts.push(qsTr("%1 %2").arg(root.segments[i].value)
                       .arg(root.segments[i].label !== undefined
                            ? root.segments[i].label : ""))
        }
        return parts.join(qsTr(", "))
    }

    Rectangle {
        anchors.fill: parent
        radius: Interface.radiusBar
        color: Theme.chipBackground
    }

    Row {
        id: stack
        anchors.fill: parent
        // The gap between segments, in the surface colour. Two design pixels:
        // one is invisible against a saturated fill and three eats a small
        // segment.
        spacing: Interface.spaceTight

        Repeater {
            model: root.segments
            delegate: Rectangle {
                required property var modelData
                required property int index
                // `stack` by id rather than `parent`, which is null while the
                // delegate is being built.
                height: stack.height
                width: {
                    if (root.scale <= 0)
                        return 0
                    const share = root.width * (modelData.value / root.scale)
                    // A segment that exists is drawn wide enough to be seen.
                    // Rounding it away would make the bar disagree with the
                    // figures beside it.
                    return modelData.value > 0
                        ? Math.max(Interface.spaceTight, Math.round(share))
                        : 0
                }
                radius: Interface.radiusBar
                color: modelData.color !== undefined
                       ? modelData.color : Theme.categorical(index)
            }
        }
    }
}
