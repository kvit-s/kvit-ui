// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// One value as it stands, and the value something proposes to replace it
// with.
//
// This is the shape every screen that asks the reader to approve a change is
// built from, and there are four of them in kvit-cash alone: the import
// preview, the quarantine group where an incoming record sits beside the
// stored one it may duplicate, the plain-language diff shown before a recipe
// update is adopted, and the review of a proposal an agent has made. Composed
// by hand each time, the four drift — in which of the two values is
// emphasised, in what is drawn where a record is new and has no old value,
// and in what a screen reader is told.
//
// The old value is on the left in the muted colour and the new one is on the
// right in the ordinary text colour, with an arrow between them. Three
// channels say which is which: position, weight of colour, and the direction
// the arrow points. A reader who cannot separate the two colours still has
// the other two.
//
// Both halves are a KvitFigure, so the em dash that stands for a value nobody
// measured arrives here without being restated. That is what `before` is for
// on a record being added: `beforeMeasured: false` draws the dash, which says
// there was nothing here rather than that there was a zero.
Row {
    id: root

    // The values, already formatted. This component does not know what a
    // currency is; `money-display.md` says how kvit-cash formats one.
    property string before: ""
    property string after: ""
    // Drawn once after each value, in the muted colour, as KvitFigure does.
    property string unit: ""
    // False draws the em dash on that side: a record being added has no
    // before, and one being deleted has no after.
    property bool beforeMeasured: true
    property bool afterMeasured: true
    // What the pair is a value of — "Amount", "Category", "Date". Optional;
    // a table whose column headings already say it should leave it empty.
    property string label: ""
    // "body" | "small" | "caption" | "strong" — passed to both figures, so
    // the two are always the same size. A pair drawn at two sizes reads as a
    // heading and a value rather than as a before and an after.
    property string role: "body"

    // True when the two values are the same, which a preview listing every
    // field of a record needs in order to draw the unchanged ones quietly.
    // Not a colour change: the arrow is what goes away, because an arrow
    // between two identical values invites the reader to look for the
    // difference.
    readonly property bool unchanged:
        root.beforeMeasured === root.afterMeasured && root.before === root.after

    spacing: Interface.spaceNear

    Accessible.role: Accessible.StaticText
    Accessible.name: {
        const was = root.beforeMeasured ? qsTr("was %1 %2").arg(root.before)
                                                           .arg(root.unit)
                                        : qsTr("was not set")
        const now = root.afterMeasured ? qsTr("now %1 %2").arg(root.after)
                                                          .arg(root.unit)
                                       : qsTr("now not set")
        if (root.unchanged)
            return root.label !== ""
                ? qsTr("%1, unchanged, %2").arg(root.label).arg(was)
                : qsTr("unchanged, %1").arg(was)
        return root.label !== ""
            ? qsTr("%1, %2, %3").arg(root.label).arg(was).arg(now)
            : qsTr("%1, %2").arg(was).arg(now)
    }

    KvitLabel {
        anchors.verticalCenter: parent.verticalCenter
        visible: root.label !== ""
        text: root.label
        role: "small"
        color: Theme.textMuted
        elide: Text.ElideNone
    }

    KvitFigure {
        anchors.verticalCenter: parent.verticalCenter
        value: root.before
        unit: root.unit
        measured: root.beforeMeasured
        role: root.role
        color: Theme.textMuted
    }

    KvitIcon {
        anchors.verticalCenter: parent.verticalCenter
        visible: !root.unchanged
        name: "arrow-right"
        color: Theme.textFaint
        implicitWidth: Interface.iconSizeSmall
        implicitHeight: Interface.iconSizeSmall
    }

    KvitFigure {
        anchors.verticalCenter: parent.verticalCenter
        visible: !root.unchanged
        value: root.after
        unit: root.unit
        measured: root.afterMeasured
        role: root.role
        color: Theme.textPrimary
    }
}
