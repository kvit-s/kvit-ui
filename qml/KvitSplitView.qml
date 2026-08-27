// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// Two or more regions the reader can resize against each other.
//
// Three applications have a private version, and the thing all three get
// slightly wrong is the handle: it is drawn as a hairline, which is one pixel
// to aim at. The handle here is a wider invisible strip with a hairline drawn
// down the middle, so the target is comfortable and the rule is still thin.
//
// It also takes keyboard focus and moves with the arrow keys, which none of
// the three do — a split a reader can only set with a mouse is a split a
// keyboard user is stuck with.
SplitView {
    id: root

    handle: Rectangle {
        id: grip
        // The visible rule is a hairline; the target is eight pixels wide.
        implicitWidth: root.orientation === Qt.Horizontal
                       ? Interface.space : 0
        implicitHeight: root.orientation === Qt.Vertical
                        ? Interface.space : 0
        color: "transparent"

        readonly property bool active:
            SplitHandle.pressed || SplitHandle.hovered

        Rectangle {
            anchors.centerIn: parent
            width: root.orientation === Qt.Horizontal
                   ? (grip.active ? Interface.spaceTight : Interface.hairline)
                   : parent.width
            height: root.orientation === Qt.Vertical
                    ? (grip.active ? Interface.spaceTight : Interface.hairline)
                    : parent.height
            color: grip.active ? Theme.accent : Theme.border
        }
    }
}
