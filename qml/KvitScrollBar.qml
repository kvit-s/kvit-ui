// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A scroll bar.
//
// prd.md §5.6 counts 32 files in this estate with a private version of this,
// which is the highest count of anything in the inventory and says something
// about how much of it is the same twenty lines. What differs between them is
// not design intent: some are two pixels wide and some are ten, some fade and
// some do not, and one draws over the content it is scrolling.
//
// It occupies its own strip rather than floating over the content. An overlay
// bar hides the right-hand column of a table and the last character of every
// elided label, and the alternative — content that reflows when the bar
// appears — is worse.
Rectangle {
    id: root

    required property Flickable flickable
    property int orientation: Qt.Vertical

    readonly property bool isVertical: orientation === Qt.Vertical
    // Guarded even though `flickable` is required, because a scroll bar
    // outlives its flickable by one event loop turn when a view is torn down,
    // and an unguarded binding re-evaluates in that gap and warns.
    readonly property real span: !flickable ? 1.0 : isVertical
        ? (flickable.height / Math.max(1, flickable.contentHeight))
        : (flickable.width / Math.max(1, flickable.contentWidth))
    readonly property real position: !flickable ? 0.0 : isVertical
        ? flickable.visibleArea.yPosition : flickable.visibleArea.xPosition

    // Hidden when everything fits, rather than shown empty. A full-length
    // handle in a permanent track is a scroll bar that says "there is more"
    // when there is not.
    visible: span < 0.999
    implicitWidth: isVertical ? Interface.spaceWide : 0
    implicitHeight: isVertical ? 0 : Interface.spaceWide
    width: isVertical ? implicitWidth : undefined
    height: isVertical ? undefined : implicitHeight

    color: Theme.panelBackground

    Rectangle {
        id: handle
        radius: width / 2
        color: drag.active || hover.hovered ? Theme.borderStrong : Theme.border

        // A minimum length, because a handle sized honestly against 250,000
        // rows is a fraction of a pixel and cannot be grabbed.
        readonly property int minimum: Interface.px(24)

        x: root.isVertical ? Interface.spaceTight
                           : Math.round(root.position * root.width)
        y: root.isVertical ? Math.round(root.position * root.height)
                           : Interface.spaceTight
        width: root.isVertical ? root.width - Interface.spaceTight * 2
                               : Math.max(minimum, Math.round(root.span * root.width))
        height: root.isVertical ? Math.max(minimum, Math.round(root.span * root.height))
                                : root.height - Interface.spaceTight * 2

        HoverHandler { id: hover }
        DragHandler {
            id: drag
            target: null
            onActiveTranslationChanged: {
                if (!active || !root.flickable)
                    return
                if (root.isVertical) {
                    const travel = Math.max(1, root.height - handle.height)
                    root.flickable.contentY = Math.max(0, Math.min(
                        root.flickable.contentHeight - root.flickable.height,
                        root.flickable.contentY
                            + activeTranslation.y / travel
                              * (root.flickable.contentHeight - root.flickable.height)))
                } else {
                    const travelX = Math.max(1, root.width - handle.width)
                    root.flickable.contentX = Math.max(0, Math.min(
                        root.flickable.contentWidth - root.flickable.width,
                        root.flickable.contentX
                            + activeTranslation.x / travelX
                              * (root.flickable.contentWidth - root.flickable.width)))
                }
            }
        }
    }
}
