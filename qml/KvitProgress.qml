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

        // The fill.
        //
        // It slides to a new value rather than stepping to it, which reads as
        // one bar moving instead of a bar being redrawn — but only while the
        // slide can finish. Long work that reports its own progress reports it
        // from the thread that also drives animations, so between two reports
        // the animation gets whatever few milliseconds that work yields and no
        // more: a 160 ms slide advances a fraction of the way, the next report
        // starts another slide from where that one stopped, and the shortfall
        // compounds. Measured during a commit of 250,000 records, which
        // reports every 2,000: the label read 42 per cent while the bar drew
        // about 5. Over a wait that long the bar is what a reader watches, so
        // the bar is what has to be true.
        //
        // So a step slides only when the bar is where the step before it asked
        // it to be. If the last slide never arrived, this one jumps, and the
        // furthest behind the bar can ever be is a single report. A bar whose
        // value moves while the application is idle is unchanged: there every
        // slide finishes, so every next step slides.
        //
        // The width is written by hand rather than bound, and the animation is
        // one this item starts rather than a Behavior, because a Behavior owns
        // its animation and refuses to have it stopped from outside — and an
        // animation left running would go on writing the width the jump was
        // meant to correct.
        Rectangle {
            id: fill

            // Where the fraction asks the bar to be.
            readonly property real wanted:
                Math.round(track.width * root.fraction)
            // Where the step before this one asked it to be. Below zero until
            // there has been one, so the first value a bar is given is drawn
            // rather than slid to: a bar that appears part of the way through
            // is saying where the work has got to, not travelling there.
            property real previouslyWanted: -1
            // Half a pixel. `wanted` is whole pixels, so a gap smaller than
            // this is the rounding at the end of a slide that finished rather
            // than a slide that did not.
            readonly property real arrived: 0.5

            onWantedChanged: {
                const stalled =
                    Math.abs(fill.width - fill.previouslyWanted) > fill.arrived
                slide.stop()
                fill.previouslyWanted = fill.wanted
                if (stalled || Theme.motionScale <= 0) {
                    fill.width = fill.wanted
                } else {
                    slide.from = fill.width
                    slide.to = fill.wanted
                    slide.start()
                }
            }

            NumberAnimation {
                id: slide
                target: fill
                property: "width"
                duration: 160 * Theme.motionScale
            }

            visible: root.determinate
            height: parent.height
            radius: parent.radius
            color: Theme.accent
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
