// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A checkbox, in three states.
//
// The third state — partial — is what a parent row shows when some of its
// children are checked and some are not. Drawing that as unchecked loses the
// information that anything under it is selected, and drawing it as checked is
// a lie the reader acts on.
//
// The tick and the dash are drawn with KvitIcon rather than by hand, which is
// what stops the tick being a two-line Canvas path in one application and a
// rotated rectangle in another.
CheckBox {
    id: root

    // Set this instead of `checked` when the state is "some of them".
    property bool partial: false

    implicitHeight: Interface.controlHeight
    spacing: Interface.spaceNear
    font.pixelSize: Interface.body

    Accessible.role: Accessible.CheckBox
    Accessible.name: root.text
    Accessible.checkStateMixed: root.partial

    indicator: Rectangle {
        implicitWidth: Interface.px(16)
        implicitHeight: Interface.px(16)
        x: root.leftPadding
        y: (root.height - height) / 2
        radius: Interface.radiusBar

        color: (root.checked || root.partial) && root.enabled
            ? Theme.accent : "transparent"
        border.width: Interface.hairline
        border.color: root.enabled ? Theme.borderStrong : Theme.border

        KvitIcon {
            anchors.centerIn: parent
            visible: root.checked || root.partial
            name: root.partial ? "minus" : "check"
            color: Theme.labelOn(Theme.accent)
            implicitWidth: parent.width - Interface.spaceSnug
            implicitHeight: parent.height - Interface.spaceSnug
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
