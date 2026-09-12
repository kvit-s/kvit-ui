// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A mark saying what state something is in, and how many things are in it.
//
// It sits between the two marks that already exist and answers what neither
// does. KvitDot says which state and carries no number, so a list showing
// three agents waiting has to say "waiting" three times or say it once and
// lose the three. KvitBadge says how many and always draws the digits, so a
// row where exactly one thing is running gets a pill reading "1" beside every
// other row's pill reading "1", which is a column of ones that says nothing.
//
// This draws the number only when there is more than one of them. At one the
// mark is the whole statement: something is in this state. Past one the count
// is what the reader needs and the mark is where it goes.
//
// The colour is the caller's. A state here is an application's own idea —
// running, needs you, failed, stalled, queued — and the library has no table
// of them; what it owns is the drawing. `Theme.accent` is the default so a
// caller that has not decided still gets a themed mark rather than a
// transparent one.
//
// `shape` is the second channel, for the same reason KvitDot has one: three
// states told apart by hue alone are three identical squares to a reader who
// cannot separate the hues, and to anybody reading a grayscale screenshot.
// Two filled shapes and a hollow one is three distinctions that survive both.
// No rotated shape, because a diamond holding a number rotates the number
// with it.
//
// The height is the caption type plus a snug margin above and below, which is
// what makes a single digit fit inside a mark small enough to sit in a row of
// text. Two digits make it wider than it is tall, which is the same shape a
// badge takes and reads the same way.
Item {
    id: root

    property int count: 1
    property int max: 99
    property color color: Theme.accent
    // "square" | "circle"
    property string shape: "square"
    // A mark drawn as an outline rather than filled, which is the third
    // distinction a caller has without reaching for another hue.
    property bool hollow: false
    // What the mark means, in words, for a screen reader. A shape and a digit
    // have none: "2 agents running" is what this is, and neither the colour
    // nor the number says it.
    required property string label

    visible: root.count > 0

    implicitHeight: Interface.caption + Interface.spaceSnug
    implicitWidth: Math.max(implicitHeight,
                            digits.implicitWidth + Interface.spaceSnug)

    Accessible.role: Accessible.StaticText
    Accessible.name: root.label

    Rectangle {
        anchors.fill: parent
        radius: root.shape === "circle" ? height / 2 : Interface.radiusBar
        color: root.hollow ? "transparent" : root.color
        border.width: root.hollow ? Interface.hairline : 0
        border.color: root.color
    }

    // Hidden at one rather than absent, so the mark keeps one implicit width
    // to measure whether or not the count is drawn.
    KvitLabel {
        id: digits
        objectName: "count"
        anchors.centerIn: parent
        visible: root.count > 1
        text: root.count > root.max
              ? Number(root.max).toLocaleString(Qt.locale(), 'f', 0) + "+"
              : Number(root.count).toLocaleString(Qt.locale(), 'f', 0)
        role: "caption"
        font.bold: true
        // labelOn rather than a fixed foreground, for the reason KvitBadge
        // gives: the caller's colour is not the accent and only the ground it
        // is actually drawn on says which label contrasts with it. A hollow
        // mark has the page behind it, so its digits take the mark's colour.
        color: root.hollow ? root.color : Theme.labelOn(root.color)
        tabular: true
        elide: Text.ElideNone
        // The mark is announced as a whole by `label`; the digits inside it
        // would otherwise be read out a second time as a bare number.
        Accessible.ignored: true
    }
}
