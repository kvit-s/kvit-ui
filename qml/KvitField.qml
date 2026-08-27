// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A single line of text the reader types.
//
// The outline is `borderStrong`, which is the control-boundary token: a field
// drawn with `border` is below 3:1 in every theme, so its edges are invisible
// and the reader has to guess where to click (accessibility.md Finding 3).
//
// An error is a message and a border together, never a border alone. A field
// that turns red and says nothing has told the reader something is wrong and
// not what, and a reader who cannot see the red has been told nothing at all.
TextField {
    id: root

    // A message under the field. Non-empty is what "in error" means; there is
    // no separate flag, because the two cannot then disagree.
    property string error: ""
    property string label: ""

    implicitHeight: Interface.controlHeight
    implicitWidth: Interface.px(200)
    leftPadding: Interface.spaceNear
    rightPadding: Interface.spaceNear

    color: enabled ? Theme.textPrimary : Theme.textDisabled
    placeholderTextColor: Theme.textFaint
    selectionColor: Theme.selectionActiveTint
    selectedTextColor: Theme.textPrimary
    font.family: Interface.resolvedFontFamily
    font.pixelSize: Interface.body
    renderType: Text.NativeRendering

    Accessible.role: Accessible.EditableText
    Accessible.name: root.label
    Accessible.description: root.error

    background: Rectangle {
        radius: Interface.radiusControl
        color: root.enabled ? Theme.popupBackground : Theme.chipBackground
        border.width: Interface.hairline
        border.color: root.error !== "" ? Theme.danger
                    : root.activeFocus ? Theme.focusRing
                    : Theme.borderStrong

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.activeFocus
            color: "transparent"
            radius: parent.radius + Interface.focusRingWidth
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    KvitLabel {
        anchors.top: parent.bottom
        anchors.topMargin: Interface.spaceTight
        anchors.left: parent.left
        width: parent.width
        visible: root.error !== ""
        text: root.error
        role: "caption"
        color: Theme.danger
        wrapMode: Text.WordWrap
        elide: Text.ElideNone
    }
}
