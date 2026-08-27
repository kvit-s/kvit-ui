// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// A card carrying one figure, what it is, how it changed and its recent shape.
//
// prd.md §5.5 notes this is the shape of six of kvit-cash's seven dashboard
// widgets, which is why it is a component rather than four components
// assembled per screen.
//
// The four parts are optional and the order is fixed. A dashboard of tiles is
// read by scanning down one column at a time — every headline, then every
// change — and a tile that puts its parts in a different order breaks that
// scan for every tile beside it.
KvitCard {
    id: root

    property string label: ""
    property string value: ""
    property string unit: ""
    property bool measured: true
    property real change: 0
    property bool hasChange: false
    // "up" | "down" | "none": which direction of change is the good one here.
    property string goodDirection: "none"
    // The recent shape, with null for a period that was not measured.
    property var history: []
    property string caption: ""

    implicitWidth: Interface.px(200)
    implicitHeight: Interface.px(96)

    Accessible.role: Accessible.StaticText
    Accessible.name: root.measured
        ? qsTr("%1: %2 %3").arg(root.label).arg(root.value).arg(root.unit)
        : qsTr("%1: not measured").arg(root.label)

    ColumnLayout {
        anchors.fill: parent
        spacing: Interface.spaceNear

        KvitLabel {
            Layout.fillWidth: true
            text: root.label
            role: "small"
            color: Theme.textMuted
        }

        RowLayout {
            Layout.fillWidth: true
            spacing: Interface.space

            KvitFigure {
                value: root.value
                unit: root.unit
                measured: root.measured
                role: "display"
            }
            Item { Layout.fillWidth: true }
            KvitDelta {
                visible: root.hasChange
                change: root.change
                goodDirection: root.goodDirection
                measured: root.measured
            }
        }

        KvitSpark {
            Layout.fillWidth: true
            visible: root.history.length > 0
            values: root.history
            label: qsTr("%1 over time").arg(root.label)
            implicitHeight: Interface.rowHeightCompact
        }

        KvitLabel {
            Layout.fillWidth: true
            visible: root.caption !== ""
            text: root.caption
            role: "caption"
            color: Theme.textFaint
        }
    }
}
