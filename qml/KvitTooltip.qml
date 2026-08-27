// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A short label that appears next to a control after a pause.
//
// A tooltip may name a control and may not be the only place its meaning
// lives: a control whose purpose is only in its tooltip is unusable by anyone
// on a keyboard or a screen reader, and unusable on a touch screen. That is
// what KvitIconButton's required `label` is for — it fills the tooltip and the
// accessible name from one string so the two cannot say different things.
ToolTip {
    id: root

    delay: 500
    // A tooltip that vanishes while somebody is reading it is worse than none.
    timeout: Math.max(3000, root.text.length * 60)
    padding: Interface.spaceNear

    contentItem: KvitLabel {
        text: root.text
        role: "small"
        color: Theme.textPrimary
        wrapMode: Text.WordWrap
        elide: Text.ElideNone
        // Wide enough for a sentence, narrow enough that the eye does not
        // have to travel back across the screen to find the next line.
        width: Math.min(implicitWidth, Interface.px(280))
    }

    background: Rectangle {
        radius: Interface.radiusControl
        color: Theme.popupBackground
        border.width: Interface.hairline
        border.color: Theme.borderStrong
    }
}
