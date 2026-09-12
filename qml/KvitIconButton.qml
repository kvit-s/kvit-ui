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
// `label` is the one place the words are written. It becomes the accessible
// name and, by default, the tooltip a pointer user reads, so the two cannot
// drift apart — which is the other half of Finding 1, where a control had a
// tooltip saying one thing and an accessible name saying another. A control
// that opens the same label in a persistent surface may suppress the tooltip.
//
// `explanation` is the second string, for a button that has more to say than
// a name: why it cannot be pressed, or what pressing it opens. It goes into
// the same two surfaces, after the label in both.
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
    // Draw the symbol at 13 rather than at 18.
    //
    // Both sizes are the library's: `Interface.iconSize` is what a symbol
    // standing on its own is drawn at, and `Interface.iconSizeSmall` is what a
    // symbol sitting beside words is drawn at, which is what KvitLink,
    // KvitNotice, KvitSelect, KvitTree, KvitMenuItem and KvitSectionHeading
    // all already use. A strip of icon buttons across the top of a pane, or
    // one hoisted onto a heading bar, is the second case: at 18 they read as
    // larger than everything around them, which is what kvit-notes-pro's
    // overview found when half its header had moved to this component and
    // half had not.
    property bool dense: false
    // A control that opens a labelled explanation already presents these words in
    // that surface. It can suppress the duplicate hover surface while retaining this
    // label as its accessible name.
    property bool tooltipEnabled: true
    // One sentence saying more than the label can: why the button is disabled,
    // or what pressing it opens. Optional, and never where the label itself
    // should be — a button whose purpose is only in its explanation is a
    // button nobody can use without hovering it.
    //
    // `label` stays the accessible name and this becomes the accessible
    // description, which is the pair a screen reader reads in that order, and
    // the tooltip shows the same two strings in the same order so a pointer
    // reader and a screen reader are told the same thing.
    property string explanation: ""
    property color iconColor: root.enabled ? Theme.textSecondary
                                           : Theme.textDisabled

    implicitWidth: Interface.controlHeight
    implicitHeight: Interface.controlHeight

    // The tooltip and accessible name normally answer different questions: a
    // pointer user reads the tooltip and a screen reader reads the name. Both
    // come from `label`, so they cannot disagree. tooltipEnabled only removes
    // the duplicate visual surface; the accessible name remains below.
    //
    // An explicit KvitTooltip rather than the attached `ToolTip.text`. The
    // attached property instantiates Qt Quick Controls' own ToolTip, which
    // arrives in the platform style — a yellow box with black text that
    // belongs to no theme in this estate and ignores the interface size. Every
    // tooltip in the library was one of those until this was written.
    KvitTooltip {
        objectName: "tooltip"
        text: root.explanation === "" ? root.label
                                      : root.label + "\n" + root.explanation
        visible: root.tooltipEnabled && (root.hovered || root.visualFocus)
                 && root.label !== ""
    }

    Accessible.role: Accessible.Button
    Accessible.name: label
    Accessible.description: root.explanation
    // A button drawn as on has to say so. Without these a screen reader
    // announces a mode that is running exactly like one that is not.
    Accessible.checkable: root.checkable || root.checked
    Accessible.checked: root.checked
    Accessible.onPressAction: root.clicked()

    background: Rectangle {
        radius: Interface.radiusControl
        color: {
            if (!root.enabled)
                return "transparent"
            if (root.checked)
                return Theme.selectionTint
            if (root.pressed)
                return Theme.selectionTint
            if (root.hovered)
                return Theme.hoverTint
            return "transparent"
        }
        // A button drawn as on carries the accent edge as well as the tint,
        // because a tint on its own is a few percent of lightness in the
        // high-contrast theme and nothing at all in a grayscale capture.
        border.width: root.checked || root.form === "ordinary"
                        ? Interface.hairline : 0
        border.color: root.checked ? Theme.accent : Theme.borderStrong

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

    // The symbol inside a rectangle, rather than being the rectangle.
    //
    // A Control stretches its content item across its available rectangle and
    // overwrites whatever size that item asked for, anchors and all, which is
    // the same rule CLAUDE.md states for a content item's position. A KvitIcon
    // used directly as the content item is therefore drawn at the button's own
    // width — 28 by default — and its glyph at whatever that clamps down to,
    // so there is no size a caller can ask for and no size this component
    // could set. One Item in between is what gives the symbol a size of its
    // own: the Item takes the stretching and the symbol is centred inside it.
    contentItem: Item {
        KvitIcon {
            objectName: "symbol"
            anchors.centerIn: parent
            name: root.symbol
            color: root.checked ? Theme.accent : root.iconColor
            width: root.dense ? Interface.iconSizeSmall : Interface.iconSize
            height: root.dense ? Interface.iconSizeSmall : Interface.iconSize
        }
    }
}
