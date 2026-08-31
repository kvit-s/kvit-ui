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

// Where the reader is, and the way back.
//
// The last crumb is the current place and is not a link, which is worth
// stating because a breadcrumb whose last item is clickable invites a reader
// to click it and get nothing, and a screen reader announces a link that goes
// nowhere.
//
// Long trails elide from the middle rather than the end. The first crumb says
// which application section this is and the last says where the reader is;
// what a reader can afford to lose is the middle.
Row {
    id: root

    // A list of { label, id } objects, root first.
    property var trail: []
    property int maximumVisible: 4

    signal crumbActivated(var crumb)

    height: Interface.breadcrumbHeight
    spacing: Interface.spaceNear

    Accessible.role: Accessible.Grouping
    Accessible.name: qsTr("Breadcrumb")

    readonly property var shown: {
        if (trail.length <= maximumVisible)
            return trail
        // First, an ellipsis crumb, then the last few. The ellipsis is not
        // activatable — it stands for several places, and picking one of them
        // is what the sidebar is for.
        const tail = trail.slice(trail.length - (maximumVisible - 2))
        return [trail[0], { label: "…", id: "", elided: true }].concat(tail)
    }

    Repeater {
        model: root.shown
        delegate: Row {
            id: crumb
            anchors.verticalCenter: root.verticalCenter
            required property var modelData
            required property int index
            spacing: Interface.spaceNear

            readonly property bool last: index === root.shown.length - 1
            readonly property bool activatable:
                !last && modelData.elided !== true

            KvitIcon {
                anchors.verticalCenter: parent.verticalCenter
                visible: crumb.index > 0
                name: "chevron-right"
                color: Theme.textFaint
                implicitWidth: Interface.caption
                implicitHeight: Interface.caption
            }
            KvitLabel {
                anchors.verticalCenter: parent.verticalCenter
                text: crumb.modelData.label
                role: "body"
                color: crumb.last ? Theme.textPrimary
                     : crumb.activatable && hover.hovered ? Theme.accent
                     : crumb.activatable ? Theme.link
                     : Theme.textFaint
                font.bold: crumb.last
                font.underline: crumb.activatable && hover.hovered
                elide: Text.ElideNone

                Accessible.role: crumb.activatable ? Accessible.Link
                                                    : Accessible.StaticText

                HoverHandler {
                    id: hover
                    enabled: crumb.activatable
                    cursorShape: Qt.PointingHandCursor
                }
                TapHandler {
                    enabled: crumb.activatable
                    onTapped: root.crumbActivated(crumb.modelData)
                }
            }
        }
    }
}
