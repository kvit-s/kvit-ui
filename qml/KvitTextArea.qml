// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// Several lines of text the reader types, or reads and selects.
//
// KvitField is one line, and a message box, a note, a commit description and
// a file's own source are all several. Before this existed each application
// drew its own styled `TextArea`, which is how three of them ended up with a
// field outline below 3:1 and one of them with none at all.
//
// The outline and the error rule are KvitField's, for the same reasons: the
// boundary token is `borderStrong`, because a control the reader has to click
// into needs findable edges, and an error is a message and a border together,
// never a border alone.
//
// Two things separate this from a styled TextArea.
//
// `underlay` is where a caller puts items drawn behind the words: a wash
// behind a marked passage, a band on a line somebody is discussing. They are
// positioned in the text's own coordinates, which is the space
// `positionToRectangle` answers in, so a caller can place one over characters
// without measuring anything else. Without somewhere to put them the caller's
// only option is to replace `background`, which is the component's own
// drawing and the whole of what a shared component exists to keep.
//
// `mono` draws the text in the monospace family, for source text and anything
// else where a column has to line up.
TextArea {
    id: root

    // "field" | "plain". A field is a box the reader types into and carries
    // its outline the whole time so its edges are findable. A plain one is a
    // document filling a pane: it has no ground and no outline, because a
    // whole-pane rectangle drawn as a field says there is something beside it
    // when there is not. Both draw the focus ring, and both take the error
    // message.
    property string form: "field"
    // A message under the box. Non-empty is what "in error" means; there is
    // no separate flag, because the two cannot then disagree.
    property string error: ""
    property string label: ""
    property bool mono: false

    // Items drawn behind the text, in the text's own coordinates.
    property alias underlay: underlayHolder.data

    // Three rows to start. A box holding a whole document sets its own height
    // from its content, and one that grows as the reader types sets it from
    // `contentHeight`; neither belongs in the default, because a box that
    // resizes itself inside a layout that sizes itself from the box is the
    // cycle that has no error message.
    implicitHeight: Interface.rowHeight * 3
    implicitWidth: Interface.px(200)
    leftPadding: Interface.spaceNear
    rightPadding: Interface.spaceNear
    topPadding: Interface.spaceNear
    bottomPadding: Interface.spaceNear

    color: enabled ? Theme.textPrimary : Theme.textDisabled
    placeholderTextColor: Theme.textFaint
    selectionColor: Theme.selectionActiveTint
    selectedTextColor: Theme.textPrimary
    font.family: root.mono ? Interface.resolvedMonoFamily
                           : Interface.resolvedFontFamily
    font.pixelSize: Interface.body
    renderType: Text.NativeRendering
    selectByMouse: true

    Accessible.role: Accessible.EditableText
    Accessible.name: root.label
    Accessible.description: root.error
    Accessible.multiLine: true
    Accessible.readOnly: root.readOnly

    background: Rectangle {
        readonly property bool plain: root.form === "plain"
        radius: plain ? 0 : Interface.radiusControl
        color: plain ? "transparent"
             : root.enabled ? Theme.popupBackground : Theme.chipBackground
        border.width: !plain || root.error !== "" ? Interface.hairline : 0
        border.color: root.error !== "" ? Theme.danger
                    : root.activeFocus ? Theme.focusRing
                    : Theme.borderStrong

        // The washes go in the background rather than over the text, because
        // the background is the one thing a Control draws before its words.
        // It fills the control, so an item placed at the rectangle
        // `positionToRectangle` returns lands on those characters.
        Item {
            id: underlayHolder
            anchors.fill: parent
        }

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.activeFocus
            color: "transparent"
            radius: parent.radius + Interface.focusRingWidth
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    KvitLabel {
        anchors.top: parent.bottom
        anchors.topMargin: Interface.spaceTight
        anchors.left: parent.left
        width: parent.width
        visible: root.error !== ""
        text: root.error
        role: "caption"
        color: Theme.danger
        wrapMode: Text.WordWrap
        elide: Text.ElideNone
    }
}
