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

// A series over time, with a value axis and a crosshair.
//
// The difference from KvitSpark is not size: it is that this one can be read
// for values rather than for shape. That means an axis with labelled
// gridlines, a hover crosshair naming the point under the pointer, and a
// stated y range rather than one taken from the data.
//
// It draws with Canvas rather than a Repeater of rectangles. That is the one
// place in this library where a shape is drawn by hand, and the reason is that
// a line is not a stack of rectangles: a polyline of 200 points as 200 rotated
// items costs 200 scene-graph nodes and still renders as a staircase at the
// joins. No charting library is introduced — the estate has zero references to
// QtCharts or QtGraphs, and the conventions that matter here (hatching for an
// upper bound, a shape as well as a hue for every level) are not things a
// charting library exposes.
Item {
    id: root

    // A list of { x, y } points, y null where the period was not measured.
    property var points: []
    property real minimumY: 0
    property real maximumY: 1
    property color color: Theme.accent
    property string label: ""
    property string unit: ""
    // How many horizontal gridlines to draw, including the two extremes.
    property int gridlines: 3

    // A second series over the same periods, drawn on the same axis: assets
    // against liabilities, spending against income, this year against last.
    // Empty for the ordinary one-series chart, which draws exactly as it did
    // before this existed.
    //
    // It is here rather than left to the caller as an overlay because the
    // crosshair cannot be. The readout names the value under the pointer, and
    // an overlay drawn on top of this cannot reach into it — a chart drawing
    // two lines and naming one value is a chart that answers the reader's
    // question with the wrong number half the time.
    //
    // Both series are indexed together: `secondPoints[i]` is the same period
    // as `points[i]`. Where one is shorter than the other it simply stops —
    // an account opened partway through the range is the ordinary case — and
    // the x axis spans whichever runs longest.
    property var secondPoints: []
    // The second series is dashed as well as differently coloured. That is
    // not decoration: the estate's rule is that no distinction rests on hue
    // alone, and two solid lines separated only by colour is that mistake in
    // the one place a reader is most likely to be reading for values.
    property color secondColor: Theme.categoricalRamp.length > 1
        ? Theme.categoricalRamp[1] : Theme.warning
    // What each line is. Drawn as a key above the plot when there are two
    // series, because two unlabelled lines cannot be told apart at all.
    property string secondLabel: ""

    readonly property bool hasSecond: root.secondPoints.length > 0

    // How many periods the x axis spans. Both series are placed on this same
    // grid, so a series that stops early stops early rather than being
    // stretched across the whole width.
    readonly property int periods:
        Math.max(root.points.length, root.secondPoints.length)

    implicitHeight: Interface.px(120)
    implicitWidth: Interface.px(320)

    // Where the pointer is, as an index into `points`, or -1. The pointer
    // position is measured from the left edge of the whole item and the
    // fraction spans the plot only, so the axis has to come off first;
    // without that the crosshair names a period to the left of the one it is
    // drawn over, and names the wrong one everywhere except the far left.
    readonly property int hovered: hover.hovered && root.periods > 0
        ? Math.max(0, Math.min(root.periods - 1,
                               Math.round((hover.point.position.x - axisWidth)
                                          / plotWidth * (root.periods - 1))))
        : -1

    // The plot rectangle, published so a caller can draw over it — an
    // annotation marking where an account's history begins, a band behind a
    // period, a threshold rule. All four are needed: the key above a
    // two-series chart takes height off the top, so a caller that knew only
    // the width would line up horizontally and be wrong vertically.
    readonly property int axisWidth: Interface.px(36)
    readonly property real plotWidth: Math.max(1, width - axisWidth)
    readonly property real plotTop: root.hasSecond
        ? key.implicitHeight + Interface.spaceSnug : 0
    readonly property real plotHeight: Math.max(1, height - plotTop)

    Accessible.role: Accessible.Chart
    Accessible.name: root.label
    // Two source strings on each branch, picked here by the count rather than
    // one with `%n` in it: with no translator installed Qt substitutes the
    // number and chooses no plural form, so a single period is announced as
    // "1 period(s)". The count goes in through a placeholder of its own and
    // is grouped the way the reader groups digits, because %n writes the bare
    // integer and five years of daily readings would be announced as "1826
    // points" beside a window that writes every other figure as "1,826". The
    // number still goes to qsTr as well, so which plural form to use stays
    // Qt's choice and a translator into a language with three of them gets
    // all three.
    Accessible.description: {
        if (root.hasSecond) {
            const grouped =
                Number(root.periods).toLocaleString(Qt.locale(), 'f', 0)
            const form = root.periods === 1
                ? qsTr("%1 period, %2 to %3 %4, two series: %5 and %6",
                       "a count of one period", root.periods)
                : qsTr("%1 periods, %2 to %3 %4, two series: %5 and %6",
                       "a count of several periods", root.periods)
            return form.arg(grouped)
                       .arg(root.minimumY).arg(root.maximumY).arg(root.unit)
                       .arg(root.label).arg(root.secondLabel)
        }
        const many = root.points.length
        const grouped = Number(many).toLocaleString(Qt.locale(), 'f', 0)
        const form = many === 1
            ? qsTr("%1 point, %2 to %3 %4", "a count of one point", many)
            : qsTr("%1 points, %2 to %3 %4", "a count of several points", many)
        return form.arg(grouped)
                   .arg(root.minimumY).arg(root.maximumY).arg(root.unit)
    }

    // The gridlines and their labels. Drawn before the series so the line sits
    // over them.
    Repeater {
        model: Math.max(2, root.gridlines)
        delegate: Item {
            required property int index
            readonly property real fraction:
                index / Math.max(1, Math.max(2, root.gridlines) - 1)
            readonly property real value:
                root.maximumY - (root.maximumY - root.minimumY) * fraction

            // Anchored to the trend itself by id: `parent` is null while a
            // Repeater delegate is being built, and these evaluate then.
            anchors.left: root.left
            anchors.right: root.right
            y: Math.round(root.plotTop + root.plotHeight * fraction)
            height: 1

            Rectangle {
                anchors.left: parent.left
                anchors.leftMargin: root.axisWidth
                anchors.right: parent.right
                height: Interface.hairline
                color: Theme.border
            }
            KvitLabel {
                anchors.right: parent.left
                anchors.rightMargin: -root.axisWidth + Interface.spaceSnug
                anchors.verticalCenter: parent.verticalCenter
                width: root.axisWidth - Interface.spaceSnug
                text: parent.value.toFixed(0)
                role: "caption"
                color: Theme.textFaint
                tabular: true
                horizontalAlignment: Text.AlignRight
            }
        }
    }

    Canvas {
        id: plot
        x: root.axisWidth
        y: root.plotTop
        width: root.plotWidth
        height: root.plotHeight
        renderStrategy: Canvas.Cooperative

        // Repaint whenever anything it draws from changes. Without the
        // explicit connections a Canvas keeps its first frame for ever, which
        // is the failure that makes a hand-drawn chart look like it works
        // until the data updates.
        Connections {
            target: root
            function onPointsChanged() { plot.requestPaint() }
            function onColorChanged() { plot.requestPaint() }
            function onSecondPointsChanged() { plot.requestPaint() }
            function onSecondColorChanged() { plot.requestPaint() }
            function onWidthChanged() { plot.requestPaint() }
            function onHeightChanged() { plot.requestPaint() }
        }
        Connections {
            target: Theme
            function onThemeChanged() { plot.requestPaint() }
        }

        onPaint: {
            const ctx = getContext("2d")
            ctx.reset()
            const count = root.periods
            if (count === 0)
                return

            const span = root.maximumY - root.minimumY
            const toY = v => height - (span > 0 ? (v - root.minimumY) / span : 0)
                                       * height
            const toX = i => count > 1 ? (i / (count - 1)) * width : width / 2

            ctx.lineWidth = Math.max(1, Interface.px(2))
            ctx.lineJoin = "round"
            ctx.lineCap = "round"

            // A run of measured points is one stroke; a null breaks the run,
            // so a gap in the data draws as a gap in the line rather than as a
            // straight segment across it. Interpolating over a hole is the
            // chart asserting values nobody measured.
            const stroke = series => {
                let drawing = false
                ctx.beginPath()
                for (let i = 0; i < count; ++i) {
                    const p = i < series.length ? series[i] : null
                    if (p === null || p.y === null || p.y === undefined) {
                        drawing = false
                        continue
                    }
                    if (!drawing) {
                        ctx.moveTo(toX(i), toY(p.y))
                        drawing = true
                    } else {
                        ctx.lineTo(toX(i), toY(p.y))
                    }
                }
                ctx.stroke()
            }

            // The second series first, so where the two cross the primary
            // one stays on top and stays readable.
            if (root.hasSecond) {
                ctx.strokeStyle = root.secondColor
                ctx.setLineDash([Interface.px(5), Interface.px(4)])
                stroke(root.secondPoints)
                ctx.setLineDash([])
            }

            ctx.strokeStyle = root.color
            stroke(root.points)
        }
    }

    // The crosshair: a rule at the hovered point and the value beside it.
    Rectangle {
        visible: root.hovered >= 0
        x: root.axisWidth + (root.periods > 1
            ? root.hovered / (root.periods - 1) * root.plotWidth
            : root.plotWidth / 2)
        y: root.plotTop
        width: Interface.hairline
        height: root.plotHeight
        color: Theme.textMuted
    }

    // What the hovered period was, one line per series. Both lines are named
    // when there are two, because a bare pair of numbers leaves the reader
    // matching them to lines by guessing.
    Column {
        id: readout
        visible: root.hovered >= 0
        y: root.plotTop
        x: Math.min(root.width - width,
                    root.axisWidth + (root.periods > 1
                        ? root.hovered / (root.periods - 1) * root.plotWidth
                        : 0) + Interface.spaceSnug)
        spacing: Interface.spaceTight

        // The value of one series at the hovered index, as words. An index
        // past the end of a series reads the same as a null: the period is
        // there and this series has nothing for it.
        function valueAt(series) {
            if (root.hovered < 0 || root.hovered >= series.length)
                return qsTr("not measured")
            const p = series[root.hovered]
            if (p === null || p.y === null || p.y === undefined)
                return qsTr("not measured")
            return qsTr("%1 %2").arg(p.y).arg(root.unit)
        }

        KvitLabel {
            text: root.hasSecond && root.label !== ""
                ? qsTr("%1: %2").arg(root.label)
                                .arg(readout.valueAt(root.points))
                : readout.valueAt(root.points)
            role: "caption"
            color: root.hasSecond ? root.color : Theme.textPrimary
            tabular: true
            elide: Text.ElideNone
        }
        KvitLabel {
            visible: root.hasSecond
            text: root.secondLabel !== ""
                ? qsTr("%1: %2").arg(root.secondLabel)
                                .arg(readout.valueAt(root.secondPoints))
                : readout.valueAt(root.secondPoints)
            role: "caption"
            color: root.secondColor
            tabular: true
            elide: Text.ElideNone
        }
    }

    // The key. Only drawn when there are two series, and it takes its room
    // off the top of the plot rather than sitting over the lines.
    //
    // Each entry is a sample of the stroke it stands for, so the dashed line
    // is identified by being dashed and not only by its colour. That is the
    // same rule the second series is drawn under, applied to the thing that
    // explains it.
    Row {
        id: key
        visible: root.hasSecond
        x: root.axisWidth
        spacing: Interface.spaceLoose

        Repeater {
            model: root.hasSecond
                ? [{ "text": root.label, "color": root.color, "dashed": false },
                   { "text": root.secondLabel, "color": root.secondColor,
                     "dashed": true }]
                : []
            delegate: Row {
                id: entry
                required property var modelData
                spacing: Interface.spaceSnug

                // The stroke sample: one dash for the solid series, two for
                // the dashed one, at the same weight the plot draws them.
                Row {
                    anchors.verticalCenter: parent.verticalCenter
                    spacing: entry.modelData.dashed ? Interface.px(3) : 0
                    Repeater {
                        model: entry.modelData.dashed ? 2 : 1
                        delegate: Rectangle {
                            width: entry.modelData.dashed
                                ? Interface.px(6) : Interface.px(15)
                            height: Math.max(1, Interface.px(2))
                            radius: Interface.radiusBar
                            color: entry.modelData.color
                        }
                    }
                }
                KvitLabel {
                    anchors.verticalCenter: parent.verticalCenter
                    text: entry.modelData.text
                    role: "caption"
                    color: Theme.textMuted
                    elide: Text.ElideNone
                }
            }
        }
    }

    HoverHandler { id: hover }

    KvitEmptyState {
        anchors.centerIn: parent
        visible: root.periods === 0
        title: qsTr("Nothing recorded yet")
        detail: qsTr("%1 will appear here once there is something to plot.")
                    .arg(root.label !== "" ? root.label : qsTr("The series"))
    }
}
