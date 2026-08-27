// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// What a view says when it has nothing to show.
//
// It says what will appear here and, where there is one, offers the action
// that would make something appear. Both halves matter. A blank region tells
// the reader nothing — not whether the data is still loading, not whether the
// filter is too narrow, not whether they have simply not made anything yet —
// and "No results" tells them only that the application knows.
//
// It is also the answer for a chart with no data. Drawing an axis around a set
// of zeros is a chart that says the values are zero, which is a different
// claim from having no values, and prd.md §5.5 calls that out as one of the
// nine rules the data-display components carry.
Column {
    id: root

    // What would be here. One short sentence, in the reader's terms.
    property string title: ""
    // Why it is not, or what to do about it. Optional.
    property string detail: ""
    property string symbol: ""
    // The action that would fill it, if there is one.
    property string action: ""

    signal actioned()

    spacing: Interface.space
    padding: Interface.viewMargin

    Accessible.role: Accessible.StaticText
    Accessible.name: root.detail !== ""
        ? qsTr("%1. %2").arg(root.title).arg(root.detail) : root.title

    KvitIcon {
        anchors.horizontalCenter: parent.horizontalCenter
        visible: root.symbol !== ""
        name: root.symbol === "" ? "info" : root.symbol
        color: Theme.textFaint
        implicitWidth: Interface.px(28)
        implicitHeight: Interface.px(28)
    }
    KvitLabel {
        anchors.horizontalCenter: parent.horizontalCenter
        text: root.title
        role: "body"
        color: Theme.textSecondary
        horizontalAlignment: Text.AlignHCenter
        elide: Text.ElideNone
    }
    KvitLabel {
        anchors.horizontalCenter: parent.horizontalCenter
        width: Math.min(root.parent ? root.parent.width * 0.7 : Interface.px(320),
                        Interface.px(320))
        visible: root.detail !== ""
        text: root.detail
        role: "small"
        color: Theme.textMuted
        wrapMode: Text.WordWrap
        horizontalAlignment: Text.AlignHCenter
        elide: Text.ElideNone
    }
    KvitButton {
        anchors.horizontalCenter: parent.horizontalCenter
        visible: root.action !== ""
        text: root.action
        form: "ordinary"
        onClicked: root.actioned()
    }
}
