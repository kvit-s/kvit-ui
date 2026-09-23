// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A small surface anchored to a control, holding something the reader can act
// on: a picker, a filter, a short form.
//
// It takes focus and closes on Escape or on a click outside, which is what
// distinguishes it from KvitHoverCard. The distinction matters: a surface a
// reader can type into must be reachable by keyboard and dismissible by
// keyboard, and a surface that only appears on hover can be neither.
Popup {
    id: root

    property string title: ""

    padding: Interface.spaceLoose
    modal: false
    // Dismiss on Escape and on a click outside, which are the two ways a
    // reader expects to get out of one.
    closePolicy: Popup.CloseOnEscape | Popup.CloseOnPressOutsideParent
    focus: true

    enter: Transition {
        NumberAnimation {
            property: "opacity"; from: 0; to: 1
            duration: 100 * Theme.motionScale
        }
    }
    exit: Transition {
        NumberAnimation {
            property: "opacity"; from: 1; to: 0
            duration: 80 * Theme.motionScale
        }
    }

    // Where the reader put the keyboard while this was fading out. A popup
    // that held the keyboard when it began to close takes it back when the
    // fade ends -- Qt does that, not this file -- and gives it to whatever
    // had it before the popup opened, or to the window's content. That
    // takes it from wherever the reader went in the meantime: the field
    // pressed straight after closing, and on a loaded machine, where the
    // fade runs long, a field they have already typed half a word into. So
    // it is given back to them. Only moves made while the fade is still
    // visible count, since Qt's own move comes once it has finished. The
    // window is held from the start of the fade, because by the time
    // `closed` arrives this surface is in none.
    QtObject {
        id: exitState
        property Item keyboard: null
        property var window: null
    }
    onAboutToHide: {
        exitState.keyboard = null
        exitState.window = root.contentItem ? root.contentItem.Window.window : null
    }
    onClosed: {
        var item = exitState.keyboard
        var win = exitState.window
        exitState.keyboard = null
        exitState.window = null
        if (item && win && item.visible && win.activeFocusItem !== item)
            item.forceActiveFocus(Qt.PopupFocusReason)
    }
    Connections {
        target: exitState.window
        enabled: root.visible && !root.opened
        function onActiveFocusItemChanged() {
            var win = exitState.window
            var item = win.activeFocusItem
            if (root.opacity <= 0 || !item || item === win.contentItem
                || item.parent === null)
                return
            // Anything of this surface's own, the frame around its content
            // included, is the keyboard leaving rather than arriving.
            for (var up = item; up; up = up.parent) {
                if (up === root.contentItem || up === root.contentItem.parent)
                    return
            }
            exitState.keyboard = item
        }
    }

    background: Rectangle {
        radius: Interface.radiusCard
        color: Theme.popupBackground
        border.width: Interface.hairline
        border.color: Theme.borderStrong

        // The accessible identity goes on the background rather than on the
        // Popup: a Popup is not an Item, and Qt's Accessible attached
        // property only attaches to one. The background is the item that
        // covers the whole surface, so it is what a screen reader finds when
        // the reader lands inside.
        Accessible.role: Accessible.Dialog
        Accessible.name: root.title
    }
}
