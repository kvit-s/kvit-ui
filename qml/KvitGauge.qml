// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// How much of an allowance has been used, and whether it is going to be
// enough.
//
// kvit-hub's BudgetGauge generalised, and it is a bar with three things a bar
// does not have: a target marked on the scale, an over-budget state, and a
// pace mark saying where the reader would be if consumption were even.
//
// The pace mark is what makes it a gauge rather than a progress bar. Sixty
// percent of a budget spent is fine on day eighteen of a month and a problem
// on day six, and a bar alone cannot say which.
Item {
    id: root

    property real value: 0
    required property real allowance
    // Where an even rate would have reached by now, 0 to 1. Negative leaves
    // the mark off, for a gauge with no time in it.
    property real pace: -1
    property bool measured: true
    property string label: ""
    property string unit: ""

    implicitHeight: Interface.barHeightWide
    implicitWidth: Interface.px(160)

    readonly property real fraction:
        allowance > 0 ? Math.max(0, Math.min(1, value / allowance)) : 0
    readonly property bool over: allowance > 0 && value > allowance
    // Behind pace, at pace, or ahead of it. The three states each get a
    // colour and the over state also gets the overflow cap drawn, so the
    // distinction is not resting on hue.
    readonly property color toneColor: {
        if (!measured)
            return Theme.textFaint
        if (over)
            return Theme.danger
        if (pace >= 0 && fraction > pace + 0.05)
            return Theme.warning
        return Theme.success
    }

    Accessible.role: Accessible.ProgressBar
    Accessible.name: root.label
    Accessible.description: !root.measured ? qsTr("not measured")
        : root.over ? qsTr("%1 of %2 %3, over by %4").arg(root.value)
                          .arg(root.allowance).arg(root.unit)
                          .arg(root.value - root.allowance)
        : qsTr("%1 of %2 %3").arg(root.value).arg(root.allowance).arg(root.unit)

    Rectangle {
        anchors.fill: parent
        radius: Interface.radiusBar
        color: Theme.chipBackground
    }

    Rectangle {
        visible: root.measured
        height: parent.height
        width: Math.round(root.width * root.fraction)
        radius: Interface.radiusBar
        color: root.toneColor
    }

    // The overflow cap: a solid block at the far end when the value is past
    // the allowance. Without it a gauge at 100% and a gauge at 180% are the
    // same picture.
    Rectangle {
        visible: root.measured && root.over
        anchors.right: parent.right
        height: parent.height
        width: Interface.spaceSnug
        color: Theme.dangerBright
    }

    // The pace mark: a hairline down the bar at where an even rate would be.
    Rectangle {
        visible: root.measured && root.pace >= 0 && root.pace <= 1
        x: Math.round(root.width * root.pace)
        width: Interface.hairline
        height: parent.height + Interface.spaceSnug
        y: -Interface.spaceTight
        color: Theme.textPrimary
    }

    Rectangle {
        visible: !root.measured
        width: Interface.spaceTight
        height: parent.height
        radius: Interface.radiusBar
        color: Theme.textFaint
    }
}
