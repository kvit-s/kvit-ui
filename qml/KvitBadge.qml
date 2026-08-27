// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A count attached to something else: unread items, pending changes, matches.
//
// A badge is always a number and the number is always small on the screen and
// possibly large in value, so it caps: past `max` it says "99+" rather than
// growing wide enough to shift the thing it is attached to.
//
// Zero hides it rather than drawing a nought. A badge showing 0 is a mark that
// says "look here" about nothing.
Rectangle {
    id: root

    property int count: 0
    property int max: 99
    // "neutral" | "accent" | "danger"
    property string tone: "accent"

    readonly property color toneColor: {
        switch (tone) {
        case "danger": return Theme.danger
        case "neutral": return Theme.textMuted
        default: return Theme.accent
        }
    }

    visible: count > 0
    implicitHeight: Interface.pillHeight
    implicitWidth: Math.max(implicitHeight, label.implicitWidth + Interface.spaceNear)
    radius: height / 2
    color: root.toneColor

    Accessible.role: Accessible.StaticText
    Accessible.name: qsTr("%n item(s)", "", root.count)

    KvitLabel {
        id: label
        anchors.centerIn: parent
        text: root.count > root.max ? root.max + "+" : String(root.count)
        role: "caption"
        color: Theme.labelOn(root.toneColor)
        tabular: true
        elide: Text.ElideNone
    }
}
