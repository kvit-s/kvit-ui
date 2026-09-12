// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
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
//
// `dashed` turns the same thing into the target a file is dropped onto. The
// drop itself is a DropArea the application wraps around this, since what a
// dropped file means is the application's business; what belongs here is the
// chrome, and specifically the action button. A drop target is reachable by
// pointer only — there is no keyboard gesture for dropping a file — so one
// that offers no other way in is unusable without a mouse. Getting the button
// from the same component as the dashed edge is what stops that being a thing
// each caller has to remember.
//
// It comes in two sizes, because "nothing here" is said in two places. In the
// middle of an empty pane it is the block above: a symbol, a sentence, a line
// under it and the way to fill it. Where seven or eight sections are stacked
// down one column and any of them may be empty, that block is several times
// taller than the rows it stands in for, and a caller that reaches for a bare
// Text instead has drawn a control by hand. `form: "compact"` is the same
// thing on one line, at the height of a slim row.
Item {
    id: root

    // "full" | "compact". The block for an empty pane, or the single line for
    // an empty section in a stack of them.
    property string form: "full"
    // What would be here. One short sentence, in the reader's terms.
    property string title: ""
    // Why it is not, or what to do about it. Optional.
    property string detail: ""
    property string symbol: ""
    // The action that would fill it, if there is one. Drawn as a button in
    // the full form and as a link in the compact one, which is the same
    // action reachable the same ways: a KvitButton is taller than the line it
    // would sit on, and a section's empty line that grew to fit one would be
    // taller than the rows it replaces.
    property string action: ""
    // Draw a dashed outline around the whole thing: this is a region
    // something can be put into, rather than a region that happens to be
    // empty. Qt Quick has no dashed border, so it is a Canvas.
    property bool dashed: false

    signal actioned()

    readonly property bool compact: root.form === "compact"

    implicitWidth: root.compact ? line.implicitWidth : body.implicitWidth
    implicitHeight: root.compact ? line.implicitHeight : body.implicitHeight

    Accessible.role: Accessible.StaticText
    Accessible.name: root.detail !== ""
        ? qsTr("%1. %2").arg(root.title).arg(root.detail) : root.title

    Canvas {
        id: edge
        anchors.fill: parent
        visible: root.dashed
        renderStrategy: Canvas.Cooperative

        // A Canvas keeps its first frame until something asks for another, so
        // without these it draws the border once at the size it was born at
        // and never again.
        Connections {
            target: root
            function onWidthChanged() { edge.requestPaint() }
            function onHeightChanged() { edge.requestPaint() }
        }
        Connections {
            target: Theme
            function onThemeChanged() { edge.requestPaint() }
        }

        onPaint: {
            const ctx = getContext("2d")
            ctx.reset()
            const line = Math.max(1, Interface.hairline)
            // Inset by half the stroke, or the outer half of it falls off the
            // edge of the canvas and the border reads as thinner on screen
            // than the value asks for.
            const inset = line / 2
            ctx.lineWidth = line
            ctx.strokeStyle = Theme.borderStrong
            ctx.setLineDash([Interface.px(5), Interface.px(4)])
            ctx.beginPath()
            ctx.roundedRect(inset, inset,
                            Math.max(0, width - line),
                            Math.max(0, height - line),
                            Interface.radiusCard, Interface.radiusCard)
            ctx.stroke()
        }
    }

    // The one line, for a section in a stack of them.
    //
    // The words are the same words: a reader who has learnt what "No matched
    // transactions yet" means in the middle of a pane reads the same sentence
    // at the top of an empty section rather than a different one.
    //
    // It starts where the rows start, which is the one place it can start.
    // The full form centres itself because it is the whole of an empty pane
    // and there is nothing beside it to line up with; this form stands in for
    // a row in a column of sections that each may be empty, and a line
    // centred among left-aligned rows reads as a different kind of thing from
    // the rows it replaces — a notice about the section rather than the
    // section's contents. The inset is `Interface.space`, which is where a
    // heading's own label starts.
    Item {
        id: line
        visible: root.compact
        anchors.fill: parent
        implicitWidth: lineRow.implicitWidth + Interface.space * 2
        implicitHeight: Math.max(Interface.rowHeightSlim,
                                 lineRow.implicitHeight + Interface.spaceSnug * 2)

        // A layout rather than a Row, because what has to survive a narrow
        // column is the sentence saying what would be here. The title and the
        // reason both give way and elide; the symbol and the action keep
        // their size, since a symbol at half width is a smudge and a
        // half-drawn action is a control nobody can read.
        RowLayout {
            id: lineRow
            objectName: "compactLine"
            anchors.left: parent.left
            anchors.leftMargin: Interface.space
            anchors.verticalCenter: parent.verticalCenter
            // The right inset is kept even though nothing is anchored to it:
            // it is what the sentence elides into, so the words stop short of
            // the edge rather than running up against it.
            width: Math.min(implicitWidth,
                            Math.max(0, parent.width - Interface.space * 2))
            spacing: Interface.spaceNear

            KvitIcon {
                Layout.alignment: Qt.AlignVCenter
                visible: root.symbol !== ""
                name: root.symbol === "" ? "info" : root.symbol
                color: Theme.textFaint
                implicitWidth: Interface.iconSizeSmall
                implicitHeight: Interface.iconSizeSmall
            }
            KvitLabel {
                objectName: "compactTitle"
                Layout.alignment: Qt.AlignVCenter
                Layout.fillWidth: true
                Layout.minimumWidth: 0
                text: root.title
                role: "small"
                color: Theme.textSecondary
            }
            KvitLabel {
                Layout.alignment: Qt.AlignVCenter
                visible: root.detail !== "" && root.title !== ""
                text: "·"
                role: "small"
                color: Theme.textFaint
                elide: Text.ElideNone
            }
            KvitLabel {
                objectName: "compactDetail"
                Layout.alignment: Qt.AlignVCenter
                Layout.fillWidth: true
                Layout.minimumWidth: 0
                visible: root.detail !== ""
                text: root.detail
                role: "small"
                color: Theme.textMuted
            }
            KvitLink {
                Layout.alignment: Qt.AlignVCenter
                Layout.leftMargin: Interface.spaceSnug
                visible: root.action !== ""
                text: root.action
                role: "small"
                elide: Text.ElideNone
                onActivated: root.actioned()
            }
        }
    }

    Column {
        id: body
        visible: !root.compact

        // Width from the parent and height from the content. The other way
        // round — filling the parent while the parent sizes itself from this
        // — is the binding loop CLAUDE.md names first.
        width: root.width

        spacing: Interface.space
        padding: Interface.viewMargin

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
}
