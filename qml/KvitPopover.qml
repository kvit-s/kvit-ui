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
