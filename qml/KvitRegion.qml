// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A body that takes the height left over and scrolls whatever does not fit.
//
// It exists because "the part in the middle" is written differently in every
// view that has one, and the differences are all mistakes: a Flickable with no
// scroll bar, a scroll bar that overlaps the last column of text, content that
// jumps when the bar appears, and a region that grows past the window instead
// of scrolling because nothing constrained its height.
//
// The bar is KvitScrollBar and sits outside the content rather than over it,
// so nothing is ever hidden behind it and the content width does not change
// when it appears.
Item {
    id: root

    default property alias content: holder.data
    property alias contentHeight: flick.contentHeight
    property alias contentWidth: flick.contentWidth
    // Padding inside the scrolled area. The default is the view margin, which
    // is what a body directly inside a window wants; a region nested inside a
    // card that already has padding sets it to zero.
    property int padding: Interface.viewMargin
    property bool horizontal: false

    clip: true

    // The scroll bars' strips are reserved whether or not a bar is showing.
    //
    // Not a cosmetic choice. Sizing the content to the bar's visibility is a
    // cycle — the bar appears, the content narrows, the content gets taller,
    // which is what decided the bar — and a region holding anything whose
    // height depends on its width (wrapped text, a wrapping layout) never
    // settles. Reserving the strip also means the content does not jump
    // sideways the moment a list grows past the fold.
    Flickable {
        id: flick
        anchors.fill: parent
        anchors.rightMargin: Interface.spaceWide
        anchors.bottomMargin: root.horizontal ? Interface.spaceWide : 0

        contentWidth: root.horizontal ? holder.childrenRect.width + root.padding * 2
                                      : width
        contentHeight: holder.childrenRect.height + root.padding * 2
        boundsBehavior: Flickable.StopAtBounds
        // A wheel that keeps travelling after the pointer stops is right on a
        // touch screen and wrong on a desktop, where it makes a list overshoot
        // the row somebody was aiming at.
        flickDeceleration: 5000

        Item {
            id: holder
            x: root.padding
            y: root.padding
            width: flick.width - root.padding * 2
        }
    }

    KvitScrollBar {
        id: vertical
        flickable: flick
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.bottom: parent.bottom
    }

    KvitScrollBar {
        id: horizontal_
        flickable: flick
        orientation: Qt.Horizontal
        anchors.left: parent.left
        anchors.right: vertical.visible ? vertical.left : parent.right
        anchors.bottom: parent.bottom
    }
}
