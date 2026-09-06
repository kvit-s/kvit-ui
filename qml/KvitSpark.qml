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
    // Two source strings picked by the count rather than one with `%n` in it:
    // with no translator installed Qt substitutes the number and chooses no
    // plural form, so a one-period series announces "1 period(s)". The number
    // goes in through %1 and is grouped the way the reader groups digits,
    // because %n writes the bare integer and a series of a thousand periods
    // is then read out digit by digit. The count still goes to qsTr, so which
    // plural form to use stays Qt's choice and a translator gets every form.
    Accessible.description: {
        const many = root.values.length
        const grouped = Number(many).toLocaleString(Qt.locale(), 'f', 0)
        return many === 1
            ? qsTr("%1 period", "a count of one period", many).arg(grouped)
            : qsTr("%1 periods", "a count of several periods", many)
                  .arg(grouped)
    }

    // The gap between bars, and where each bar starts.
    //
    // Not a Row with fractional widths. Twelve bars across 160 pixels with a
    // one-pixel gap gives each one 12.4167 pixels, a Row lays them out at
    // cumulative fractional positions, and the renderer snaps each bar to
    // device pixels on its own — so the gaps come out one pixel in some
    // places and two in others, which is what a reader sees as an unevenly
    // spaced series.
    //
    // Rounding the *boundaries* instead makes every gap exactly `gap`. The
    // bar widths then vary by at most one pixel, which nobody reads at this
    // size, where uneven gaps are obvious. The slot arithmetic runs over
    // `width + gap` so that the last bar ends exactly at the right edge
    // rather than leaving a trailing gap.
    // Not properties: a reader of the catalogue should see what this
    // component takes, and a bar gap is not that.
    function slotStart(index) {
        const gap = Interface.hairline
        const count = Math.max(1, root.values.length)
        return Math.round(index * (root.width + gap) / count)
    }
    function slotWidth(index) {
        return Math.max(1, root.slotStart(index + 1) - root.slotStart(index)
                            - Interface.hairline)
    }

    Repeater {
        model: root.values
        delegate: Item {
            id: slot
            required property var modelData
            required property int index

            x: root.slotStart(slot.index)
            width: root.slotWidth(slot.index)
            height: root.height

            // A measured value: a bar of its height.
            Rectangle {
                visible: slot.modelData !== null
                anchors.bottom: parent.bottom
                width: parent.width
                height: slot.modelData === null ? 0
                    : Math.max(Interface.hairline,
                               Math.round(root.height
                                          * (slot.modelData / root.scale)))
                radius: Interface.radiusBar
                color: root.color
            }

            // A period nobody measured: a tick on the baseline. Muted and one
            // pixel tall, so it reads as an absence rather than as a very
            // small value.
            Rectangle {
                visible: slot.modelData === null
                anchors.bottom: parent.bottom
                width: parent.width
                height: Interface.hairline
                color: Theme.textFaint
            }
        }
    }
}
