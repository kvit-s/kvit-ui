// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// More about the thing under the pointer, without clicking.
//
// kvit-hub's HoverCard generalised. It is read-only and it must stay that
// way: nothing inside one can be the only route to an action, because a
// surface that appears on hover is unreachable by keyboard and by touch. Where
// an action belongs on the same thing, it goes in a KvitPopover opened by a
// click, or in the row itself.
//
// It shows what a dense row could not fit — the full name behind an elided
// one, the exact figure behind a rounded one, the date behind "3 days ago" —
// which is why the components that round or elide are the ones that most often
// have one attached.
Rectangle {
    id: root

    property bool shown: false
    default property alias content: holder.data
    property int padding: Interface.spaceLoose

    implicitWidth: holder.childrenRect.width + padding * 2
    implicitHeight: holder.childrenRect.height + padding * 2


    visible: opacity > 0
    opacity: shown ? 1 : 0
    Behavior on opacity {
        NumberAnimation { duration: 100 * Theme.motionScale }
    }

    radius: Interface.radiusCard
    color: Theme.popupBackground
    border.width: Interface.hairline
    border.color: Theme.borderStrong

    // Ignored by assistive technology on purpose: everything in here is a
    // second view of something already on the screen, and announcing it again
    // makes a list twice as long to listen to.
    Accessible.ignored: true

    // The content holder takes its height from what is in it and its width
    // from the parent, rather than filling the parent in both directions.
    //
    // `anchors.fill` plus an `implicitHeight` read off `childrenRect` is a
    // cycle: the holder's height comes from this item and this item's from the
    // holder's children. Qt does not warn about it — it hangs during
    // construction, which is a much harder thing to find than a warning.
    Item {
        id: holder
        x: root.padding
        y: root.padding
        width: childrenRect.width
        height: childrenRect.height
    }
}
