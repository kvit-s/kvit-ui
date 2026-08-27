// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A list of commands, opened by a right-click or by a button.
//
// Eighteen files in the estate have a private version, and what they mostly
// get wrong is the same two things: no keyboard route in, and no separator
// before the destructive item at the bottom.
//
// Qt's Menu handles the keyboard once it is a real Menu rather than a
// Rectangle full of MouseAreas — arrow keys move, Return activates, Escape
// closes, and the menu takes focus when it opens. That is most of why this
// wraps Menu rather than building one.
Menu {
    id: root

    implicitWidth: Interface.px(200)
    padding: Interface.spaceSnug
    overlap: 0

    delegate: KvitMenuItem {}

    background: Rectangle {
        radius: Interface.radiusControl
        color: Theme.popupBackground
        border.width: Interface.hairline
        border.color: Theme.borderStrong

        Accessible.role: Accessible.PopupMenu
        Accessible.name: root.title
    }
}
