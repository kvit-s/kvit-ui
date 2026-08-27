// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A value chosen from a continuous range by dragging.
//
// Right where the reader is choosing a feel rather than a number — an opacity,
// a threshold they will judge by looking at the result. Wrong where the exact
// value matters, because hitting a particular number on a slider is difficult
// and the reader cannot type one: KvitStepper or KvitNumberField are for that.
//
// The current value is always drawn, because a slider whose position is its
// only output makes the reader estimate from the handle's position.
Slider {
    id: root

    property string label: ""
    property string unit: ""
    property int precision: 0

    implicitHeight: Interface.controlHeight
    implicitWidth: Interface.px(180)

    Accessible.role: Accessible.Slider
    Accessible.name: root.label

    background: Rectangle {
        x: root.leftPadding
        y: root.topPadding + root.availableHeight / 2 - height / 2
        width: root.availableWidth
        height: Interface.barHeight
        radius: Interface.radiusBar
        color: Theme.chipBackground

        Rectangle {
            width: root.visualPosition * parent.width
            height: parent.height
            radius: parent.radius
            color: root.enabled ? Theme.accent : Theme.border
        }
    }

    handle: Rectangle {
        x: root.leftPadding + root.visualPosition
           * (root.availableWidth - width)
        y: root.topPadding + root.availableHeight / 2 - height / 2
        implicitWidth: Interface.px(14)
        implicitHeight: Interface.px(14)
        radius: width / 2
        color: root.pressed ? Theme.hoverTint : Theme.popupBackground
        border.width: Interface.hairline
        border.color: Theme.borderStrong

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.visualFocus
            color: "transparent"
            radius: parent.radius + Interface.focusRingWidth
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    KvitFigure {
        anchors.left: parent.right
        anchors.leftMargin: Interface.space
        anchors.verticalCenter: parent.verticalCenter
        value: root.value.toFixed(root.precision)
        unit: root.unit
        role: "small"
    }
}
