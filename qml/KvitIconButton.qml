// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A button whose whole label is a symbol.
//
// This is kvit-notes' IconButton, taking an icon name rather than a literal
// Unicode character. It is the component that replaces the estate's most
// common accessibility defect: most of the small controls across the four
// applications are a Rectangle with a MouseArea in it, which reaches no
// assistive technology at all — a screen reader is not told there is a control
// there, the keyboard cannot get to it, and Space does nothing.
//
// An AbstractButton is what fixes that, and it fixes it by being one rather
// than by having accessibility bolted on: Qt publishes it with a button role
// and a name, it takes tab focus, and Space or Return activates it
// (accessibility.md Finding 1).
//
// `label` is the one place the words are written. It becomes both the tooltip
// a pointer user reads and the name a screen reader says, so the two cannot
// drift apart — which is the other half of Finding 1, where a control had a
// tooltip saying one thing and an accessible name saying another.
AbstractButton {
    id: root

    // What the symbol means: an IconCatalog meaning name.
    // `symbol` rather than `icon`, throughout the vocabulary: Qt Quick
    // Controls' AbstractButton already has a FINAL `icon` grouped property,
    // and a component that derives from one cannot declare its own. Naming it
    // `icon` here and `symbol` on the buttons would mean a call site has to
    // remember which kind of component it is looking at.
    required property string symbol
    // What the button does, in words. Required, because a button whose entire
    // label is a picture has nothing else to tell a screen reader.
    required property string label

    // "ordinary" | "quiet". A quiet button has no resting ground at all and
    // draws one on hover, which is what a dense toolbar wants; an ordinary one
    // carries its outline the whole time so its edges are findable.
    property string form: "quiet"
    property bool checked_: false
    property color iconColor: root.enabled ? Theme.textSecondary
                                           : Theme.textDisabled

    implicitWidth: Interface.controlHeight
    implicitHeight: Interface.controlHeight

    // Both the tooltip and the accessible name, because they answer different
    // questions: a pointer user reads the tooltip and a screen reader reads
    // the name, and a control that has one and not the other is unusable by
    // half its readers. Both come from `label`, so they cannot disagree.
    //
    // An explicit KvitTooltip rather than the attached `ToolTip.text`. The
    // attached property instantiates Qt Quick Controls' own ToolTip, which
    // arrives in the platform style — a yellow box with black text that
    // belongs to no theme in this estate and ignores the interface size. Every
    // tooltip in the library was one of those until this was written.
    KvitTooltip {
        text: root.label
        visible: root.hovered && root.label !== ""
    }

    Accessible.role: Accessible.Button
    Accessible.name: label
    Accessible.onPressAction: root.clicked()

    background: Rectangle {
        radius: Interface.radiusControl
        color: {
            if (!root.enabled)
                return "transparent"
            if (root.checked_)
                return Theme.selectionTint
            if (root.pressed)
                return Theme.selectionTint
            if (root.hovered)
                return Theme.hoverTint
            return "transparent"
        }
        border.width: root.form === "ordinary" ? Interface.hairline : 0
        border.color: Theme.borderStrong

        // The focus ring is a separate rectangle outside the fill rather than
        // a border colour change, because a ring drawn as a border disappears
        // on a control that already has one, and because hue alone must never
        // be the only signal (accessibility.md Finding 3).
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

    contentItem: KvitIcon {
        name: root.symbol
        color: root.checked_ ? Theme.accent : root.iconColor
        anchors.centerIn: parent
    }
}
