// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// The strip along the bottom of a window: what the application is doing now,
// and one or two standing facts about what is open.
//
// kvit-hub's WindowStatusBar generalised. Two slots rather than a free-form
// row, because the two halves answer different questions and mixing them is
// what makes a status bar unreadable: the left says what is happening and
// changes, the right says what is true and does not.
Rectangle {
    id: root

    // What is happening now. Empty leaves the left side blank, which is the
    // resting state — a status bar that always has something to say trains the
    // reader to stop looking at it.
    property string activity: ""
    // Standing facts, right aligned. Each becomes one label with a separator
    // between.
    property var facts: []

    implicitHeight: Interface.statusBarHeight
    color: Theme.footerBackground

    KvitDivider {
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
    }

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: Interface.space
        anchors.rightMargin: Interface.space
        spacing: Interface.space

        KvitLabel {
            Layout.fillWidth: true
            role: "caption"
            color: Theme.textMuted
            text: root.activity
            visible: text !== ""
        }
        Item { Layout.fillWidth: root.activity === "" }

        Repeater {
            model: root.facts
            delegate: RowLayout {
                required property string modelData
                required property int index
                spacing: Interface.space
                KvitLabel {
                    text: "·"
                    role: "caption"
                    color: Theme.textFaint
                    visible: parent.index > 0
                }
                KvitLabel {
                    text: parent.modelData
                    role: "caption"
                    color: Theme.textMuted
                    tabular: true
                }
            }
        }
    }
}
