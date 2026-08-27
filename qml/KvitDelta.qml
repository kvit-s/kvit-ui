// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// How much something changed, and in which direction.
//
// Direction is carried by an arrow as well as by colour, which is the rule
// that makes this a component rather than a coloured number: green-up and
// red-down is the single most common place a screen rests a whole meaning on
// hue, and it is the pair a deuteranope cannot separate.
//
// Whether up is good is the caller's to say. A rise in spending and a rise in
// savings are the same arrow and opposite colours, so `goodDirection` is
// explicit and there is no default that quietly assumes finance or quietly
// assumes the reverse.
Row {
    id: root

    property real change: 0
    property string unit: ""
    property bool measured: true
    // "up" | "down" | "none". "none" draws the arrow in the neutral colour,
    // for a change that is neither good nor bad.
    property string goodDirection: "none"
    // Round to this many decimal places when drawing.
    property int precision: 0

    spacing: Interface.spaceSnug

    readonly property bool rising: change > 0
    readonly property bool flat: change === 0
    readonly property color toneColor: {
        if (!measured || flat || goodDirection === "none")
            return Theme.textMuted
        const good = rising ? goodDirection === "up" : goodDirection === "down"
        return good ? Theme.success : Theme.danger
    }

    Accessible.role: Accessible.StaticText
    Accessible.name: !root.measured ? qsTr("not measured")
        : root.flat ? qsTr("unchanged")
        : root.rising ? qsTr("up %1 %2").arg(Math.abs(root.change).toFixed(root.precision)).arg(root.unit)
        : qsTr("down %1 %2").arg(Math.abs(root.change).toFixed(root.precision)).arg(root.unit)

    KvitIcon {
        anchors.verticalCenter: parent.verticalCenter
        visible: root.measured && !root.flat
        name: root.rising ? "arrow-up" : "arrow-down"
        color: root.toneColor
        implicitWidth: Interface.caption
        implicitHeight: Interface.caption
    }
    KvitFigure {
        anchors.verticalCenter: parent.verticalCenter
        value: root.flat ? qsTr("no change")
                         : Math.abs(root.change).toFixed(root.precision)
        unit: root.flat ? "" : root.unit
        measured: root.measured
        color: root.toneColor
        role: "small"
    }
}
