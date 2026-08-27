// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A measured value.
//
// Three rules travel with it, and they are the reason this is a component
// rather than a Text with a number in it.
//
// Tabular numerals, so a column of figures lines up and a value that changes
// does not shift the label beside it. Proportional digits make a right-aligned
// column ragged and a live figure jitter.
//
// The unit in a muted colour and at the smaller role, because the unit is the
// same on every row and the number is not: setting them alike makes the reader
// re-read the unit each time to find the number.
//
// And an em dash where a value was not measured, rather than a zero. This is
// the rule money-display.md describes without naming the component: a balance
// nobody has computed is not a balance of zero, an effort nobody recorded is
// not zero effort, and drawing either as 0 states something false. The dash
// says "not measured" and nothing else does.
Row {
    id: root

    property string value: ""
    property string unit: ""
    property bool measured: true
    // "body" | "small" | "caption" | "strong" | "title" | "headline" | "display"
    property string role: "body"
    property color color: Theme.textPrimary
    // A figure that is an upper bound rather than a measurement — the same
    // distinction hatching carries on a bar. The tilde is a second channel
    // beside the hatch, for a figure that has no bar beside it.
    property bool bounded: false

    spacing: Interface.spaceSnug

    Accessible.role: Accessible.StaticText
    Accessible.name: root.measured
        ? (root.bounded ? qsTr("at most %1 %2").arg(root.value).arg(root.unit)
                        : qsTr("%1 %2").arg(root.value).arg(root.unit))
        : qsTr("not measured")

    KvitLabel {
        anchors.verticalCenter: parent.verticalCenter
        // The em dash, not a hyphen and not "n/a": it is one glyph wide at
        // every size, it reads as an absence rather than as a minus sign, and
        // it needs no translation.
        text: root.measured ? (root.bounded ? "~" + root.value : root.value)
                            : "—"
        role: root.role
        color: root.measured ? root.color : Theme.textFaint
        tabular: true
        elide: Text.ElideNone
    }
    KvitLabel {
        anchors.verticalCenter: parent.verticalCenter
        visible: root.unit !== "" && root.measured
        text: root.unit
        role: root.role === "caption" ? "caption" : "small"
        color: Theme.textMuted
        elide: Text.ElideNone
    }
}
