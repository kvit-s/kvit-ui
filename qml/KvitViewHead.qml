// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// The strip at the top of a view: what this view is, how much is in it, and
// the controls that act on the whole of it.
//
// Distinct from KvitHeader, which belongs to the window. This one belongs to
// whatever is currently being shown, and it is where a filter, a sort control
// and a count go — so a reader looking for "how do I narrow this down" looks
// in the same place on every screen.
Item {
    id: root

    property string title: ""
    // How many things the view holds, and the word for one of them. A count
    // with no noun beside it does not say what was counted.
    property int count: -1
    property string counted: qsTr("item")
    property string subtitle: ""
    default property alias controls: controlSlot.data

    implicitHeight: Interface.rowHeight
    implicitWidth: Interface.px(600)

    Accessible.role: Accessible.Heading
    Accessible.name: root.title

    RowLayout {
        anchors.fill: parent
        spacing: Interface.columnGap

        ColumnLayout {
            Layout.fillWidth: true
            spacing: Interface.spaceTight

            RowLayout {
                spacing: Interface.space
                KvitLabel {
                    text: root.title
                    role: "title"
                    font.bold: true
                }
                KvitLabel {
                    visible: root.count >= 0
                    text: qsTr("%n %1(s)", "", root.count).arg(root.counted)
                    role: "small"
                    color: Theme.textFaint
                    tabular: true
                }
            }
            KvitLabel {
                Layout.fillWidth: true
                visible: root.subtitle !== ""
                text: root.subtitle
                role: "small"
                color: Theme.textMuted
                // A subtitle is a sentence rather than a name, so it wraps
                // where a row label elides: cutting the end off an
                // explanation loses the half that explains.
                wrapMode: Text.WordWrap
                elide: Text.ElideNone
            }
        }

        Row {
            id: controlSlot
            Layout.alignment: Qt.AlignVCenter
            spacing: Interface.space
        }
    }
}
