// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A region of the window with its own ground: a sidebar, a toolbar strip, the
// area behind a group of controls.
//
// A panel is structural and a card is content — the distinction is worth
// keeping because they take different grounds and only one of them has a
// corner radius. A panel is part of the window's furniture and its edges are
// the window's edges; a card is an object sitting on a surface.
Rectangle {
    id: root

    // Which edges carry a rule. A panel between two others usually wants one
    // on the side it abuts and none anywhere else.
    property bool ruleTop: false
    property bool ruleBottom: false
    property bool ruleLeft: false
    property bool ruleRight: false

    color: Theme.panelBackground

    KvitDivider { visible: root.ruleTop; anchors.top: parent.top; anchors.left: parent.left; anchors.right: parent.right }
    KvitDivider { visible: root.ruleBottom; anchors.bottom: parent.bottom; anchors.left: parent.left; anchors.right: parent.right }
    KvitDivider { visible: root.ruleLeft; vertical: true; anchors.left: parent.left; anchors.top: parent.top; anchors.bottom: parent.bottom }
    KvitDivider { visible: root.ruleRight; vertical: true; anchors.right: parent.right; anchors.top: parent.top; anchors.bottom: parent.bottom }
}
