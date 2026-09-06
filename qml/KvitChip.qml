// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// The lazy parts below are nested components, and Bound is what makes them
// read this file's ids legally rather than by accident: without it a nested
// component resolves an outer id at run time through the object hierarchy.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// A small labelled mark: what kind of thing this is, what state it is in, one
// word about it.
//
// The `tone` names a meaning rather than a colour, which is what lets the same
// call site be right in all four themes. A chip with `tone: "danger"` is red
// in three of them and bright coral in high contrast, and a chip written with
// a hex value is wrong in three.
//
// Every tone carries a second channel besides its hue: the outline is drawn in
// the tone as well as the fill being tinted with it, and a `strong` chip is
// filled rather than tinted. That is the rule the high-contrast pass in the
// gallery checks, and it is why a reader who cannot separate the red tone from
// the amber one can still separate a warning from a failure by where the
// weight is.
Rectangle {
    id: root

    property string text: ""
    // "neutral" | "accent" | "success" | "warning" | "danger" | "info"
    property string tone: "neutral"
    // A filled chip rather than a tinted one, for the one chip on a row that
    // is the point of the row.
    property bool strong: false
    property string symbol: ""

    readonly property color toneColor: {
        switch (tone) {
        case "accent":  return Theme.accent
        case "success": return Theme.success
        case "warning": return Theme.warning
        case "danger":  return Theme.danger
        case "info":    return Theme.link
        default:        return Theme.textMuted
        }
    }

    implicitHeight: Interface.chipHeight
    implicitWidth: content.width + Interface.spaceNear * 2
    radius: Interface.radiusChip

    color: root.strong ? root.toneColor
         : root.tone === "neutral" ? Theme.chipBackground
         : Qt.alpha(root.toneColor, 0.16)
    border.width: root.strong ? 0 : Interface.hairline
    border.color: root.tone === "neutral" ? Theme.border : root.toneColor

    Accessible.role: Accessible.StaticText
    Accessible.name: root.text

    Row {
        id: content
        anchors.centerIn: parent
        spacing: Interface.spaceSnug

        // Built only when the chip has a symbol. A KvitIcon carries a font
        // loader and two items of its own, and a table draws several hundred
        // chips at a time; most of them have no symbol at all.
        //
        // `visible` as well as `active`, because a Row lays out every visible
        // child and would keep a gap where the absent icon would have been.
        Loader {
            anchors.verticalCenter: parent.verticalCenter
            active: root.symbol !== ""
            visible: active
            sourceComponent: KvitIcon {
                name: root.symbol
                color: root.strong ? Theme.labelOn(root.toneColor) : root.toneColor
                implicitWidth: Interface.caption
                implicitHeight: Interface.caption
            }
        }
        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            text: root.text
            role: "caption"
            // labelOn rather than onAccent: onAccent answers only for the
            // accent, and a filled danger chip needs the label that contrasts
            // with *its* fill (accessibility.md Finding 3).
            color: root.strong ? Theme.labelOn(root.toneColor)
                 : root.tone === "neutral" ? Theme.textSecondary
                 : root.toneColor
            elide: Text.ElideNone
        }
    }
}
