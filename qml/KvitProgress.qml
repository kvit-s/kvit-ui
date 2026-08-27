// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// How far through something the application is.
//
// Two states, and which one is showing is the most useful thing on the bar.
// `determinate` means the total is known and the fill says how much is done;
// otherwise it is a moving band that says only that work is happening. Drawing
// an unknown total as a determinate bar creeping toward the end is the pattern
// that makes people distrust progress bars.
//
// It carries a label saying what is happening. "47%" answers how far and not
// how far through what, and a reader who left the screen and came back has no
// way to find out.
Item {
    id: root

    property real value: 0
    property real maximum: 1
    property bool determinate: true
    property string label: ""
    property bool showPercent: true

    implicitHeight: Interface.rowHeightCompact
    implicitWidth: Interface.px(200)

    readonly property real fraction:
        maximum > 0 ? Math.max(0, Math.min(1, value / maximum)) : 0

    Accessible.role: Accessible.ProgressBar
    Accessible.name: root.label
    Accessible.description: root.determinate
        ? qsTr("%1 percent").arg(Math.round(root.fraction * 100))
        : qsTr("in progress")

    Row {
        id: text
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
        spacing: Interface.space

        KvitLabel {
            width: parent.width - percent.width - parent.spacing
            text: root.label
            role: "small"
            color: Theme.textSecondary
        }
        KvitLabel {
            id: percent
            visible: root.determinate && root.showPercent
            text: Math.round(root.fraction * 100) + "%"
            role: "small"
            color: Theme.textMuted
            tabular: true
            elide: Text.ElideNone
        }
    }

    Rectangle {
        id: track
        anchors.bottom: parent.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        height: Interface.barHeight
        radius: Interface.radiusBar
        color: Theme.chipBackground
        clip: true

        Rectangle {
            visible: root.determinate
            width: Math.round(track.width * root.fraction)
            height: parent.height
            radius: parent.radius
            color: Theme.accent
            Behavior on width {
                NumberAnimation { duration: 160 * Theme.motionScale }
            }
        }

        // The indeterminate band. With motion stilled it does not animate:
        // instead the whole track takes a muted fill, so a reader who has
        // asked for no motion still sees that something is happening rather
        // than an empty bar that looks stalled.
        Rectangle {
            id: band
            visible: !root.determinate
            width: Theme.motionScale > 0 ? track.width * 0.3 : track.width
            height: parent.height
            radius: parent.radius
            color: Theme.motionScale > 0 ? Theme.accent
                                         : Qt.alpha(Theme.accent, 0.4)

            XAnimator on x {
                running: !root.determinate && Theme.motionScale > 0
                loops: Animation.Infinite
                from: -band.width
                to: track.width
                duration: 1200
            }
        }
    }
}
