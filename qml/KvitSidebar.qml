// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick

// The list of places down the left edge.
//
// `collapsed` is the state KvitWindow puts it in below the laptop breakpoint:
// the rail, where every item is its symbol alone. That is why KvitSidebarItem
// requires a symbol as well as a label — an item with no symbol disappears
// entirely when the sidebar collapses, and there is nothing at the call site
// to say so.
Column {
    id: root

    property bool collapsed: false
    spacing: 0

    Accessible.role: Accessible.List
    Accessible.name: qsTr("Navigation")

    // Every KvitSidebarItem inside inherits the collapsed state rather than
    // being told individually, which is what stops one item in a rail still
    // drawing its label.
    onCollapsedChanged: propagate()
    Component.onCompleted: propagate()

    function propagate() {
        for (let i = 0; i < children.length; ++i) {
            if (children[i].collapsed !== undefined)
                children[i].collapsed = root.collapsed
        }
    }
}
