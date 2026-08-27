// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A number with a minus and a plus beside it: an interface size, a column
// count, a number of days.
//
// For a small integer range a reader adjusts by one or two, which is exactly
// where a slider is wrong — a slider makes a precise value hard to hit and
// gives no way to read what the current one is without a label beside it.
//
// The range comes from whatever owns it rather than being repeated at the call
// site. kvit-hub's "Interface size" row is the caller this was written for:
// it takes `from` and `to` from Interface.minFontSize and maxFontSize, so the
// clamp is enforced in one place.
Row {
    id: root

    property int value: 0
    property int from: 0
    property int to: 100
    property int step: 1
    property string label: ""
    // The unit, drawn after the number in muted colour.
    property string unit: ""

    signal valueModified(int value)

    spacing: Interface.spaceSnug

    Accessible.role: Accessible.SpinBox
    // The current value goes in the name rather than in a separate value
    // property: Qt's attached Accessible has no numeric value, so a screen
    // reader that is told only "Interface size" learns what the control is and
    // not what it currently says.
    Accessible.name: root.unit !== ""
        ? qsTr("%1, %2 %3").arg(root.label).arg(root.value).arg(root.unit)
        : qsTr("%1, %2").arg(root.label).arg(root.value)

    function apply(next) {
        const bounded = Math.max(root.from, Math.min(root.to, next))
        if (bounded !== root.value)
            root.valueModified(bounded)
    }

    KvitIconButton {
        anchors.verticalCenter: parent.verticalCenter
        symbol: "minus"
        label: qsTr("Decrease %1").arg(root.label)
        form: "ordinary"
        enabled: root.value > root.from
        implicitWidth: Interface.px(22)
        implicitHeight: Interface.px(22)
        onClicked: root.apply(root.value - root.step)
    }

    KvitFigure {
        anchors.verticalCenter: parent.verticalCenter
        // A fixed width, so the row beside it does not shuffle sideways every
        // time the number gains or loses a digit.
        width: Interface.px(34)
        value: String(root.value)
        unit: root.unit
        role: "body"
    }

    KvitIconButton {
        anchors.verticalCenter: parent.verticalCenter
        symbol: "plus"
        label: qsTr("Increase %1").arg(root.label)
        form: "ordinary"
        enabled: root.value < root.to
        implicitWidth: Interface.px(22)
        implicitHeight: Interface.px(22)
        onClicked: root.apply(root.value + root.step)
    }
}
