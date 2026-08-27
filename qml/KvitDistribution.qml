// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// How a set of values is spread: where the bulk of them sit and how far the
// tails reach.
//
// The reason a chart needs this and not just an average: an average of four
// days and an average of four days made of one twenty-day outlier and three
// same-day items are the same number and completely different situations.
//
// Drawn as a box with whiskers rather than as a histogram, because a row in a
// list has room for one and not the other, and because the five numbers a box
// carries — the extremes, the quartiles and the median — are exactly what a
// reader asks of a spread.
Item {
    id: root

    property real minimum: 0
    property real lowerQuartile: 0
    property real median: 0
    property real upperQuartile: 0
    property real maximum: 0
    // The scale the box is drawn against, which is usually wider than this
    // one distribution so several rows are comparable.
    property real scaleMinimum: 0
    property real scaleMaximum: 1
    property bool measured: true
    property color color: Theme.accent
    property string label: ""
    property string unit: ""

    implicitHeight: Interface.barHeightWide
    // The natural width, which is what implicitWidth means: how wide this
    // wants to be when nothing constrains it. A caller that wants it to fill
    // something assigns `width`.
    //
    // Writing `parent.width` here instead says "fill my parent", and that is
    // unusable inside anything that sizes itself to its children — a Column, a
    // Row, an Item measured by childrenRect — because the child then asks the
    // parent how wide it is while the parent is asking the child. Every one of
    // these collapsed to zero width in a bare Column.
    implicitWidth: Interface.px(160)

    function at(value) {
        const span = scaleMaximum - scaleMinimum
        if (span <= 0)
            return 0
        return Math.round(root.width
                          * Math.max(0, Math.min(1, (value - scaleMinimum) / span)))
    }

    Accessible.role: Accessible.Chart
    Accessible.name: root.label
    Accessible.description: !root.measured ? qsTr("not measured")
        : qsTr("median %1 %2, middle half %3 to %4, range %5 to %6")
              .arg(root.median).arg(root.unit)
              .arg(root.lowerQuartile).arg(root.upperQuartile)
              .arg(root.minimum).arg(root.maximum)

    // The whisker: the full range, as a thin line.
    Rectangle {
        visible: root.measured
        anchors.verticalCenter: parent.verticalCenter
        x: root.at(root.minimum)
        width: Math.max(Interface.hairline,
                        root.at(root.maximum) - root.at(root.minimum))
        height: Interface.hairline
        color: Theme.textFaint
    }

    // The box: the middle half.
    Rectangle {
        visible: root.measured
        x: root.at(root.lowerQuartile)
        width: Math.max(Interface.spaceTight,
                        root.at(root.upperQuartile) - root.at(root.lowerQuartile))
        height: parent.height
        radius: Interface.radiusBar
        color: Qt.alpha(root.color, 0.35)
        border.width: Interface.hairline
        border.color: root.color
    }

    // The median: a solid rule across the box. It is drawn in the text
    // colour rather than in the series colour so it stays visible against a
    // tinted box in every theme.
    Rectangle {
        visible: root.measured
        x: root.at(root.median)
        width: Interface.spaceTight
        height: parent.height
        color: Theme.textPrimary
    }

    KvitLabel {
        anchors.centerIn: parent
        visible: !root.measured
        text: "—"
        role: "caption"
        color: Theme.textFaint
        elide: Text.ElideNone
    }
}
