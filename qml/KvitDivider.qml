// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A rule between two things.
//
// One design pixel, which stays one device pixel at any interface size
// because Interface.hairline never rounds to nothing, and `border` rather than
// `borderStrong` because this is decoration rather than the boundary of a
// control. Using `borderStrong` here is the mistake accessibility.md Finding 3
// describes from the other direction: it makes a separator look like the edge
// of something clickable.
Rectangle {
    id: root

    property bool vertical: false
    // A rule that stops short of the edges, for a divider inside a padded
    // container where a full-bleed line would cut the padding in half.
    property int inset: 0

    implicitWidth: vertical ? Interface.hairline : 0
    implicitHeight: vertical ? 0 : Interface.hairline
    width: vertical ? implicitWidth : undefined
    height: vertical ? undefined : implicitHeight

    color: Theme.border

    anchors.leftMargin: vertical ? 0 : inset
    anchors.rightMargin: vertical ? 0 : inset
    anchors.topMargin: vertical ? inset : 0
    anchors.bottomMargin: vertical ? inset : 0
}
