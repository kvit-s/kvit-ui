// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// The list row a reader sees most: one line carrying a name, what it is, one
// phrase about it and one figure, right aligned.
//
// KvitRow is the ground and this is the arrangement on top of it. They are
// separate because the arrangement is what nearly every list wants and the
// ground is what the few that do not can still build on — a row with a chart
// in it takes KvitRow directly.
//
// The order across the row is fixed, and the fixed order is the point: a
// reader scanning a list down the left edge finds the name in the same place
// on every screen in the estate, and scanning down the right edge finds the
// number.
KvitRow {
    id: root

    form: "slim"

    property string name: ""
    // What kind of thing this is, said next to the name where the name does
    // not say it.
    property string kind: ""
    // One phrase about its state. Muted, because it is context rather than
    // identity.
    property string phrase: ""
    property string figure: ""
    property string unit: ""
    // False draws an em dash rather than the figure, which is what an
    // unmeasured value looks like. A zero is a measurement.
    property bool measured: true
    property string symbol: ""

    label: [name, kind, phrase].filter(part => part !== "").join(", ")

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: Interface.spaceNear
        anchors.rightMargin: Interface.spaceNear
        spacing: Interface.space

        KvitIcon {
            visible: root.symbol !== ""
            name: root.symbol === "" ? "dot" : root.symbol
            color: Theme.textMuted
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }
        KvitLabel {
            text: root.name
            role: "body"
            Layout.maximumWidth: root.width * 0.45
        }
        KvitLabel {
            visible: root.kind !== ""
            text: root.kind
            role: "small"
            color: Theme.textFaint
        }
        KvitLabel {
            Layout.fillWidth: true
            text: root.phrase
            role: "small"
            color: Theme.textMuted
        }
        KvitFigure {
            visible: root.figure !== "" || !root.measured
            value: root.figure
            unit: root.unit
            measured: root.measured
        }
    }
}
