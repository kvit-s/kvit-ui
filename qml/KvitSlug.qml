// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// An identifier: a reference, a short hash, a key, a ticket number.
//
// Monospace, because the reader's task with an identifier is comparison rather
// than reading — is this the same one as that — and two proportional strings
// of the same length are different widths, so the eye cannot line them up.
//
// It elides in the middle rather than at the end. The end of a hash is what
// distinguishes it from its neighbours; an identifier cut short at the right
// is an identifier that matches everything.
Rectangle {
    id: root

    property string text: ""
    // A flat run of monospace text, for a slug inside a table cell where a
    // ground on every row would read as a column of boxes.
    property bool ground: true

    implicitHeight: root.ground ? Interface.chipHeight : label.implicitHeight
    implicitWidth: Math.min(label.implicitWidth + (root.ground ? Interface.spaceNear * 2 : 0),
                            parent ? parent.width : label.implicitWidth)
    radius: Interface.radiusBar
    color: root.ground ? Theme.inlineCodeBackground : "transparent"

    Accessible.role: Accessible.StaticText
    Accessible.name: root.text

    KvitLabel {
        id: label
        anchors.centerIn: parent
        width: parent.width - (root.ground ? Interface.spaceNear * 2 : 0)
        text: root.text
        role: "small"
        mono: true
        color: Theme.textSecondary
        elide: Text.ElideMiddle
        horizontalAlignment: Text.AlignLeft
    }
}
