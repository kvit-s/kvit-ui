// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A figure with its name under it: one number a reader is meant to take away.
//
// `frame.css` already defines this shape, which is why it is in the inventory
// — the drawings have been using it for a while and the applications have been
// redrawing it per screen.
//
// The number is above the label rather than beside it, and it is the larger of
// the two. Both are deliberate: a row of these is read across the numbers, and
// a label set at the same size as its figure makes the reader do the work of
// deciding which is which on every one.
Column {
    id: root

    property string value: ""
    property string unit: ""
    property string label: ""
    property bool measured: true
    property bool bounded: false
    property color color: Theme.textPrimary
    // "display" | "headline" | "title" — how much of the screen this figure
    // is meant to command.
    property string role: "headline"

    spacing: Interface.spaceTight

    Accessible.role: Accessible.StaticText
    Accessible.name: root.measured
        ? qsTr("%1: %2 %3").arg(root.label).arg(root.value).arg(root.unit)
        : qsTr("%1: not measured").arg(root.label)

    KvitFigure {
        value: root.value
        unit: root.unit
        measured: root.measured
        bounded: root.bounded
        color: root.color
        role: root.role
    }
    KvitLabel {
        text: root.label
        role: "small"
        color: Theme.textMuted
        elide: Text.ElideNone
    }
}
