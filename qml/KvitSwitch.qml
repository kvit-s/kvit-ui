// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// An option that is on or off and takes effect immediately.
//
// Ten files in the estate have a private version. The one thing they disagree
// about is what the difference between a switch and a checkbox is, and it is
// worth stating: a switch takes effect the moment it moves, and a checkbox is
// a value in a form that takes effect when the form is submitted. A settings
// row that applies at once is a switch; a row of options above a Save button
// is checkboxes.
//
// The knob moves *and* the track fills, so the state is carried by position as
// well as by colour. A switch whose only difference between on and off is the
// track's hue is one a colour-blind reader has to toggle to find out.
Switch {
    id: root

    implicitHeight: Interface.controlHeight
    spacing: Interface.space
    font.pixelSize: Interface.body

    Accessible.role: Accessible.CheckBox
    Accessible.name: root.text
    Accessible.checked: root.checked

    indicator: Rectangle {
        implicitWidth: Interface.px(28)
        implicitHeight: Interface.px(16)
        x: root.leftPadding
        y: (root.height - height) / 2
        radius: height / 2

        color: root.checked && root.enabled ? Theme.accent : Theme.chipBackground
        border.width: Interface.hairline
        border.color: root.enabled ? Theme.borderStrong : Theme.border

        Rectangle {
            id: knob
            width: parent.height - Interface.spaceSnug
            height: width
            radius: width / 2
            y: Interface.spaceTight
            x: root.checked ? parent.width - width - Interface.spaceTight
                            : Interface.spaceTight
            color: root.checked && root.enabled
                   ? Theme.labelOn(Theme.accent) : Theme.textMuted

            Behavior on x {
                NumberAnimation {
                    duration: 120 * Theme.motionScale
                    easing.type: Easing.OutCubic
                }
            }
        }

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

    contentItem: KvitLabel {
        text: root.text
        role: "body"
        color: root.enabled ? Theme.textPrimary : Theme.textDisabled
        leftPadding: root.indicator.width + root.spacing
        elide: Text.ElideNone
    }
}
