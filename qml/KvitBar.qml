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

// One quantity against a scale.
//
// kvit-hub's AxisBar generalised, and it carries three of the nine rules
// prd.md §5.5 attaches to the data-display components.
//
// **A bar with no value is not a bar of zero.** `measured: false` draws the
// track and a tick where the bar would start, rather than a zero-width fill,
// because a zero-width fill is indistinguishable from a measured zero.
//
// **Hatching means an upper bound.** A figure that is "at most this much" is
// drawn hatched, and the hatch survives grayscale — it is geometry rather than
// a lighter shade — which is what makes it legible in a printed screenshot and
// to a reader who cannot separate the two tints.
//
// **The scale is stated, not implied.** `maximum` is required. A bar drawn
// against the largest value in its own list rescales every time the list
// changes, so two screenshots taken a day apart are not comparable and nothing
// on the screen says so.
Item {
    id: root

    property real value: 0
    required property real maximum
    property bool measured: true
    // The value is an upper bound rather than a measurement.
    property bool bounded: false
    property color color: Theme.accent
    property bool wide: false
    // What this bar is, for a screen reader. A bar is a shape with no text in
    // it, so without this it is announced as nothing.
    property string label: ""
    property string unit: ""

    implicitHeight: wide ? Interface.barHeightWide : Interface.barHeight
    implicitWidth: parent ? parent.width : Interface.px(120)

    readonly property real fraction:
        maximum > 0 ? Math.max(0, Math.min(1, value / maximum)) : 0

    Accessible.role: Accessible.ProgressBar
    Accessible.name: root.label
    Accessible.description: !root.measured ? qsTr("not measured")
        : root.bounded ? qsTr("at most %1 of %2 %3").arg(root.value)
                             .arg(root.maximum).arg(root.unit)
        : qsTr("%1 of %2 %3").arg(root.value).arg(root.maximum).arg(root.unit)

    // The track: what the whole scale would look like. Always drawn, so the
    // reader can see how much of the scale a short bar is short of.
    Rectangle {
        anchors.fill: parent
        radius: Interface.radiusBar
        color: Theme.chipBackground
    }

    Rectangle {
        id: fill
        visible: root.measured
        height: parent.height
        width: Math.round(root.width * root.fraction)
        radius: Interface.radiusBar
        color: root.color
    }

    // The hatch. Drawn as a repeater of thin diagonals in the theme's own
    // hatch-stroke colour, which is per theme so the stripes survive every
    // ground a bar can sit on.
    Item {
        id: hatch
        anchors.fill: fill
        visible: root.measured && root.bounded
        clip: true
        Repeater {
            model: Math.ceil(fill.width / Interface.spaceSnug) + 4
            delegate: Rectangle {
                required property int index
                width: Interface.hairline
                // `hatch` by id rather than `parent`: a delegate's parent is
                // null while it is being built, and a binding that reads it
                // then warns before settling.
                height: hatch.height * 3
                x: index * Interface.spaceSnug - hatch.height
                y: -hatch.height
                rotation: 45
                transformOrigin: Item.TopLeft
                color: Theme.hatchAlt
                opacity: 0.7
            }
        }
    }

    // What "not measured" looks like: a tick at the origin. It is the same
    // idea as KvitFigure's em dash — an absence drawn as an absence rather
    // than as a zero.
    Rectangle {
        visible: !root.measured
        x: 0
        width: Interface.spaceTight
        height: parent.height
        radius: Interface.radiusBar
        color: Theme.textFaint
    }
}
