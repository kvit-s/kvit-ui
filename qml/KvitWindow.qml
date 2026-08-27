// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Window
import Kvit.Ui

// The application shell: a window with a header, an optional sidebar, a body
// and a status bar.
//
// It exists because those four pieces are assembled slightly differently in
// each application and the differences are all accidents. What it fixes in
// particular is the reflow: the sidebar collapses to its rail below
// `Interface.widthLaptop` and the whole window has a floor at
// `Interface.widthFloor`, so a view never has to work out for itself when it
// is too narrow.
//
// The breakpoints scale with the interface size, which is the reason they are
// tokens rather than literals: at twice the interface size the same screen
// needs twice the width before it stops being cramped.
Window {
    id: root

    property alias header: headerSlot.data
    property alias sidebar: sidebarSlot.data
    property alias body: bodySlot.data
    property alias statusBar: statusSlot.data

    // False hides the sidebar entirely, which is different from collapsing it:
    // a hidden sidebar is a reader's choice and a collapsed one is the window
    // being narrow.
    property bool sidebarVisible: true
    // Whether the window has a header and a status bar at all.
    //
    // Explicit rather than derived from whether the slot has children. A slot
    // sized from `childrenRect` whose child fills the slot is a loop — the
    // child's height comes from the slot and the slot's from the child — and
    // it hangs during construction rather than warning, because neither side
    // ever settles.
    property bool hasHeader: true
    property bool hasStatusBar: true
    readonly property bool sidebarCollapsed: width < Interface.widthLaptop
    readonly property bool narrow: width < Interface.widthFloor

    minimumWidth: Interface.widthFloor
    minimumHeight: Interface.px(600)
    width: Interface.widthDrawn
    height: Interface.px(960)
    color: Theme.windowBackground

    Item {
        id: headerSlot
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
        height: root.hasHeader ? Interface.headerHeight : 0
    }

    KvitPanel {
        id: sidebarPanel
        visible: root.sidebarVisible
        anchors.top: headerSlot.bottom
        anchors.bottom: statusSlot.top
        anchors.left: parent.left
        width: root.sidebarCollapsed ? Interface.railWidth : Interface.sidebarWidth
        ruleRight: true

        // The width change is animated because it is a layout the reader is
        // watching happen, and an instant jump reads as a redraw glitch rather
        // than as the window adapting. It follows the reduced-motion setting:
        // motionScale is 0 when motion is stilled, which makes the duration
        // zero and the change instant (features.md §14.3).
        Behavior on width {
            NumberAnimation {
                duration: 120 * Theme.motionScale
                easing.type: Easing.OutCubic
            }
        }

        Item {
            id: sidebarSlot
            anchors.fill: parent
        }
    }

    Item {
        id: bodySlot
        anchors.top: headerSlot.bottom
        anchors.bottom: statusSlot.top
        anchors.left: sidebarPanel.visible ? sidebarPanel.right : parent.left
        anchors.right: parent.right
    }

    Item {
        id: statusSlot
        anchors.bottom: parent.bottom
        anchors.left: parent.left
        anchors.right: parent.right
        height: root.hasStatusBar ? Interface.statusBarHeight : 0
    }
}
