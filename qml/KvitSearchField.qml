// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A field that filters something.
//
// Separate from KvitField because it behaves differently rather than because
// it looks different: it carries a magnifier so it is findable without reading
// the placeholder, a clear button that appears only when there is something to
// clear, and Escape clears rather than reverting — which is what a reader
// expects of a filter and not of a text field.
//
// It also announces its result count, because filtering is the one interaction
// where the whole outcome happens somewhere else on the screen: a screen reader
// user who types into an unannounced filter gets no feedback at all.
KvitField {
    id: root

    // How many things the filter left. Negative leaves it unsaid, which is
    // right for a filter whose result is not a countable list.
    property int matches: -1
    property string matchedNoun: qsTr("result")

    leftPadding: Interface.controlHeight
    rightPadding: root.text !== "" ? Interface.controlHeight : Interface.spaceNear
    placeholderText: qsTr("Filter…")

    Accessible.role: Accessible.EditableText
    Accessible.name: root.label !== "" ? root.label : qsTr("Filter")
    Accessible.description: root.matches < 0 ? ""
        : qsTr("%n %1(s)", "", root.matches).arg(root.matchedNoun)

    Keys.onEscapePressed: function (event) {
        if (root.text !== "") {
            root.clear()
            event.accepted = true
        }
    }

    KvitIcon {
        anchors.left: parent.left
        anchors.leftMargin: Interface.spaceNear
        anchors.verticalCenter: parent.verticalCenter
        name: "search"
        color: Theme.textFaint
        implicitWidth: Interface.iconSizeSmall
        implicitHeight: Interface.iconSizeSmall
    }

    KvitIconButton {
        anchors.right: parent.right
        anchors.rightMargin: Interface.spaceTight
        anchors.verticalCenter: parent.verticalCenter
        visible: root.text !== ""
        symbol: "close"
        label: qsTr("Clear the filter")
        implicitWidth: Interface.px(20)
        implicitHeight: Interface.px(20)
        onClicked: root.clear()
    }
}
