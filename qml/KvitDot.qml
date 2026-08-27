// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A small filled circle standing for one thing's state: a health level, a
// severity, whether something is running.
//
// Drawn as a rounded rectangle rather than as an icon glyph, which is the one
// deliberate exception to "no component draws its own shape". A circle six
// pixels across is exact as geometry and approximate as a rasterized glyph,
// and a row of them at that size is where the difference shows.
//
// `shape` is the second channel. A dot that says "stalled" only by being red
// says nothing to a reader who cannot separate red from amber, so a level
// carries a shape as well: a circle, a square or a diamond. The gallery's
// high-contrast pass is what checks that a screen using dots set the shape.
Item {
    id: root

    property color color: Theme.textMuted
    // "circle" | "square" | "diamond"
    property string shape: "circle"
    property bool hollow: false
    // What the dot means, for a screen reader. A shape has no text in it.
    property string label: ""

    implicitWidth: Interface.spaceNear
    implicitHeight: Interface.spaceNear

    Accessible.role: Accessible.Indicator
    Accessible.name: root.label
    Accessible.ignored: root.label === ""

    Rectangle {
        anchors.centerIn: parent
        width: root.width
        height: root.height
        rotation: root.shape === "diamond" ? 45 : 0
        radius: root.shape === "circle" ? width / 2 : Interface.hairline
        color: root.hollow ? "transparent" : root.color
        border.width: root.hollow ? Interface.hairline : 0
        border.color: root.color
    }
}
