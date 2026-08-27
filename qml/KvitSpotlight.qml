// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// Darken everything except one region, and say something about it.
//
// This is the primitive under a guided tour, and it is deliberately only the
// primitive. prd.md §5.6 calls the tour itself a judgement rather than a yes
// or no: two applications want one — kvit-cash's first-run flow and
// kvit-notes-pro's "guide me" — and what drives the stepping is different in
// each. One reads a fixed script; the other decides the next step from what
// the agent is doing. A shared stepper would have to be either of those and
// would be wrong for the other, so what is shared is focusing a region and
// each application brings its own reason to move.
//
// Escape always closes it. A reader who wants out of a tour and cannot get out
// is a reader who force-quits the application.
Item {
    id: root

    // The item to leave lit, in this item's coordinates. Null darkens
    // everything, which is the state before the first step.
    property Item target: null
    property string title: ""
    property string detail: ""
    property bool shown: false
    property int padding: Interface.space

    signal dismissed()

    anchors.fill: parent
    visible: opacity > 0
    opacity: shown ? 1 : 0
    Behavior on opacity {
        NumberAnimation { duration: 140 * Theme.motionScale }
    }

    Accessible.role: Accessible.Dialog
    Accessible.name: root.title
    Accessible.description: root.detail

    readonly property rect hole: {
        if (!target)
            return Qt.rect(0, 0, 0, 0)
        const at = target.mapToItem(root, 0, 0)
        return Qt.rect(at.x - padding, at.y - padding,
                       target.width + padding * 2,
                       target.height + padding * 2)
    }

    // The shade, as four rectangles around the hole rather than one rectangle
    // with a hole cut in it. Four rectangles need no shader and no mask, and
    // they leave the lit region genuinely untouched — a mask would still
    // composite over it.
    Repeater {
        model: 4
        delegate: Rectangle {
            required property int index
            color: Qt.rgba(0, 0, 0, 0.55)
            x: index === 1 ? root.hole.x + root.hole.width
             : index === 3 ? 0 : 0
            y: index === 0 ? 0
             : index === 2 ? root.hole.y + root.hole.height
             : root.hole.y
            width: index === 0 || index === 2 ? root.width
                 : index === 1 ? Math.max(0, root.width - root.hole.x - root.hole.width)
                 : Math.max(0, root.hole.x)
            height: index === 0 ? Math.max(0, root.hole.y)
                  : index === 2 ? Math.max(0, root.height - root.hole.y - root.hole.height)
                  : root.hole.height
        }
    }

    // The outline around the lit region, so it reads as chosen rather than as
    // a gap in the shade.
    Rectangle {
        visible: root.target !== null
        x: root.hole.x
        y: root.hole.y
        width: root.hole.width
        height: root.hole.height
        color: "transparent"
        radius: Interface.radiusCard
        border.width: Interface.focusRingWidth
        border.color: Theme.focusRing
    }

    KvitCard {
        visible: root.title !== "" || root.detail !== ""
        x: Math.max(Interface.viewMargin,
                    Math.min(root.width - width - Interface.viewMargin,
                             root.hole.x))
        y: root.hole.y + root.hole.height + Interface.space
           + height < root.height
           ? root.hole.y + root.hole.height + Interface.space
           : Math.max(Interface.viewMargin, root.hole.y - height - Interface.space)
        width: Interface.px(280)

        Column {
            spacing: Interface.spaceNear
            width: parent.width - Interface.spaceLoose * 2

            KvitLabel {
                width: parent.width
                text: root.title
                role: "strong"
                font.bold: true
            }
            KvitLabel {
                width: parent.width
                visible: root.detail !== ""
                text: root.detail
                role: "small"
                color: Theme.textSecondary
                wrapMode: Text.WordWrap
                elide: Text.ElideNone
            }
        }
    }

    focus: root.shown
    Keys.onEscapePressed: root.dismissed()

    // A click anywhere in the shade dismisses. A tour that traps the reader
    // until they find its close button is a tour they resent.
    TapHandler { onTapped: root.dismissed() }
}
