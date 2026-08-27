// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// What happened to something, in order.
//
// Two callers name it in prd.md §5.6: kvit-cash's per-account update history
// and kvit-notes-pro's session history.
//
// Newest first, which is the order both of those want: a reader opening a
// history is asking what happened recently, and making them scroll to the
// bottom to find out is the most common thing a history gets wrong.
//
// Each entry carries who did it as well as what and when. In an estate where
// an agent and a person both change the same things, "who" is the column that
// makes a history worth reading at all.
Column {
    id: root

    // A list of { when, what, who, detail, tone } objects, newest first.
    property var entries: []
    property string label: ""

    spacing: 0

    Accessible.role: Accessible.List
    Accessible.name: root.label

    Repeater {
        model: root.entries
        delegate: Item {
            id: entry
            required property var modelData
            required property int index

            width: root.width
            implicitHeight: Math.max(Interface.rowHeightSlim,
                                     body.implicitHeight + Interface.space)

            readonly property color toneColor: {
                switch (modelData.tone) {
                case "success": return Theme.success
                case "warning": return Theme.warning
                case "danger":  return Theme.danger
                default:        return Theme.textMuted
                }
            }

            Accessible.role: Accessible.ListItem
            Accessible.name: qsTr("%1, %2, %3")
                .arg(modelData.when === undefined ? "" : modelData.when)
                .arg(modelData.what === undefined ? "" : modelData.what)
                .arg(modelData.who === undefined ? "" : modelData.who)

            // The spine: a rule down the left with a dot per entry. It stops
            // at the last entry rather than running past it, because a line
            // continuing past the oldest event says there is more below.
            Rectangle {
                x: Interface.spaceNear
                y: entry.index === 0 ? entry.height / 2 : 0
                width: Interface.hairline
                height: entry.index === root.entries.length - 1
                        ? entry.height / 2 : entry.height
                color: Theme.border
            }

            KvitDot {
                x: Interface.spaceSnug
                y: entry.height / 2 - height / 2
                width: Interface.spaceNear
                height: Interface.spaceNear
                color: entry.toneColor
                // The shape carries the tone as well as the colour, which is
                // what lets a reader who cannot separate the hues still see
                // that one entry is unlike its neighbours.
                shape: entry.modelData.tone === "danger" ? "diamond"
                     : entry.modelData.tone === "warning" ? "square"
                     : "circle"
                label: entry.modelData.what === undefined
                       ? "" : entry.modelData.what
            }

            Column {
                id: body
                x: Interface.px(28)
                y: Interface.spaceSnug
                width: parent.width - x - Interface.spaceNear
                spacing: Interface.spaceTight

                Row {
                    spacing: Interface.space
                    KvitLabel {
                        text: entry.modelData.what === undefined
                              ? "" : entry.modelData.what
                        role: "body"
                        elide: Text.ElideNone
                    }
                    KvitLabel {
                        text: entry.modelData.who === undefined
                              ? "" : entry.modelData.who
                        role: "small"
                        color: Theme.textMuted
                        elide: Text.ElideNone
                    }
                    KvitLabel {
                        text: entry.modelData.when === undefined
                              ? "" : entry.modelData.when
                        role: "small"
                        color: Theme.textFaint
                        tabular: true
                        elide: Text.ElideNone
                    }
                }
                KvitLabel {
                    width: parent.width
                    visible: entry.modelData.detail !== undefined
                             && entry.modelData.detail !== ""
                    text: entry.modelData.detail === undefined
                          ? "" : entry.modelData.detail
                    role: "small"
                    color: Theme.textMuted
                    wrapMode: Text.WordWrap
                    elide: Text.ElideNone
                }
            }
        }
    }

    KvitEmptyState {
        width: root.width
        visible: root.entries.length === 0
        title: qsTr("Nothing has happened yet")
        detail: qsTr("Changes will be listed here, newest first, with who made them.")
        symbol: "clock"
    }
}
