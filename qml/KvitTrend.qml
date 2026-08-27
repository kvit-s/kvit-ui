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

    implicitHeight: Interface.px(120)
    implicitWidth: Interface.px(320)

    // Where the pointer is, as an index into `points`, or −1.
    readonly property int hovered: hover.hovered && root.points.length > 0
        ? Math.max(0, Math.min(root.points.length - 1,
                               Math.round(hover.point.position.x / plotWidth
                                          * (root.points.length - 1))))
        : -1

    readonly property int axisWidth: Interface.px(36)
    readonly property real plotWidth: Math.max(1, width - axisWidth)

    Accessible.role: Accessible.Chart
    Accessible.name: root.label
    Accessible.description: qsTr("%n point(s), %1 to %2 %3", "", root.points.length)
        .arg(root.minimumY).arg(root.maximumY).arg(root.unit)

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
            y: Math.round(root.height * fraction)
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
        width: root.plotWidth
        height: root.height
        renderStrategy: Canvas.Cooperative

        // Repaint whenever anything it draws from changes. Without the
        // explicit connections a Canvas keeps its first frame for ever, which
        // is the failure that makes a hand-drawn chart look like it works
        // until the data updates.
        Connections {
            target: root
            function onPointsChanged() { plot.requestPaint() }
            function onColorChanged() { plot.requestPaint() }
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
            const count = root.points.length
            if (count === 0)
                return

            const span = root.maximumY - root.minimumY
            const toY = v => height - (span > 0 ? (v - root.minimumY) / span : 0)
                                       * height
            const toX = i => count > 1 ? (i / (count - 1)) * width : width / 2

            ctx.lineWidth = Math.max(1, Interface.px(2))
            ctx.strokeStyle = root.color
            ctx.lineJoin = "round"
            ctx.lineCap = "round"

            // A run of measured points is one stroke; a null breaks the run,
            // so a gap in the data draws as a gap in the line rather than as a
            // straight segment across it. Interpolating over a hole is the
            // chart asserting values nobody measured.
            let drawing = false
            ctx.beginPath()
            for (let i = 0; i < count; ++i) {
                const p = root.points[i]
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
    }

    // The crosshair: a rule at the hovered point and the value beside it.
    Rectangle {
        visible: root.hovered >= 0
        x: root.axisWidth + (root.points.length > 1
            ? root.hovered / (root.points.length - 1) * root.plotWidth
            : root.plotWidth / 2)
        width: Interface.hairline
        height: root.height
        color: Theme.textMuted
    }

    KvitLabel {
        visible: root.hovered >= 0 && root.points.length > 0
        anchors.top: parent.top
        x: Math.min(root.width - width,
                    root.axisWidth + (root.points.length > 1
                        ? root.hovered / (root.points.length - 1) * root.plotWidth
                        : 0) + Interface.spaceSnug)
        text: {
            if (root.hovered < 0 || root.hovered >= root.points.length)
                return ""
            const p = root.points[root.hovered]
            if (p === null || p.y === null || p.y === undefined)
                return qsTr("not measured")
            return qsTr("%1 %2").arg(p.y).arg(root.unit)
        }
        role: "caption"
        tabular: true
        elide: Text.ElideNone
    }

    HoverHandler { id: hover }

    KvitEmptyState {
        anchors.centerIn: parent
        visible: root.points.length === 0
        title: qsTr("Nothing recorded yet")
        detail: qsTr("%1 will appear here once there is something to plot.")
                    .arg(root.label !== "" ? root.label : qsTr("The series"))
    }
}
