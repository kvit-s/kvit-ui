// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// One place in the sidebar.
//
// kvit-hub's SidebarItem generalised. The selected item is marked by a bar
// down its leading edge as well as by a tint, which is the second channel: a
// tint alone is a two-percent lightness difference in the high-contrast theme
// and invisible in a screenshot.
//
// `symbol` is required rather than optional because of the rail. When the
// sidebar collapses the label is gone and the symbol is the whole item, and an
// item that only has words vanishes.
AbstractButton {
    id: root

    required property string symbol
    property bool selected: false
    property bool collapsed: false
    property int count: -1
    // The word for one of the things counted. Empty gives the generic
    // "Review, 11 items"; a caller that knows the noun sets it and the item
    // is announced as "Review, 1 decision" or "Review, 11 decisions".
    property string counted: ""
    // The word for several of them. English adds an s and that is what this
    // falls back to; a noun that does not pluralise that way says so here.
    // Two slots rather than one already-inflected word, because a single word
    // is right at one count and wrong at every other.
    property string countedPlural: ""
    // Whether the count stays visible once the sidebar has collapsed to the
    // rail. Off by default, because that is what the rail drew before this
    // existed and a rail is deliberately a strip of symbols; an application
    // whose rail has to keep a backlog in front of the reader turns it on,
    // and the badge then sits on the symbol's upper corner rather than on
    // the label's line.
    property bool countInRail: false
    // The largest count the badge writes out; past it the badge says "99+".
    // The cap is right where the number is only a signal and wrong where the
    // place it points at states the true figure, because the two then
    // disagree about the same thing. An item whose destination shows the real
    // number raises this to it.
    property int countMax: 99

    implicitHeight: Interface.rowHeightSlim
    implicitWidth: Interface.sidebarWidth

    Accessible.role: Accessible.ListItem
    Accessible.name: {
        if (root.count < 0)
            return root.text
        // Grouped by the reader's locale, for the reason KvitBadge groups its
        // own: a count read out digit by digit is not a number anybody hears.
        const grouped = Number(root.count).toLocaleString(Qt.locale(), 'f', 0)
        if (root.counted !== "") {
            if (root.count === 1) {
                return qsTr("%1, %2 %3",
                            "a label, a count and the singular of what it counts",
                            root.count)
                    .arg(root.text).arg(grouped).arg(root.counted)
            }
            if (root.countedPlural !== "") {
                return qsTr("%1, %2 %3",
                            "a label, a count and the plural of what it counts",
                            root.count)
                    .arg(root.text).arg(grouped).arg(root.countedPlural)
            }
            // No plural was given, so English suffixes an s. The whole phrase
            // is the translatable unit rather than the suffix on its own,
            // which is nothing a translator can act on.
            return qsTr("%1, %2 %3s",
                        "a label, a count and a noun pluralised by suffixing s",
                        root.count)
                .arg(root.text).arg(grouped).arg(root.counted)
        }
        if (root.count === 1)
            return qsTr("%1, %2 item", "a label and a count of one thing",
                        root.count).arg(root.text).arg(grouped)
        return qsTr("%1, %2 items", "a label and a count of several things",
                    root.count).arg(root.text).arg(grouped)
    }
    Accessible.selected: root.selected
    Accessible.onPressAction: root.clicked()

    // The tooltip is what makes the rail usable: collapsed, the label is the
    // only thing that says where this goes.
    //
    // KvitTooltip rather than the attached `ToolTip.text`, which instantiates
    // the platform style's own and arrives as a yellow box belonging to no
    // theme here.
    KvitTooltip {
        text: root.text
        visible: root.collapsed && root.hovered
        delay: 400
    }

    background: Rectangle {
        color: root.selected ? Theme.selectionTint
             : root.hovered ? Theme.hoverTint
             : "transparent"

        Rectangle {
            anchors.left: parent.left
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            width: Interface.spaceTight
            visible: root.selected
            color: Theme.accent
        }

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.visualFocus
            color: "transparent"
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    contentItem: Item {
        anchors.fill: parent

        KvitIcon {
            id: mark
            anchors.left: parent.left
            anchors.leftMargin: root.collapsed
                ? (root.width - width) / 2 : Interface.spaceLoose
            anchors.verticalCenter: parent.verticalCenter
            name: root.symbol
            color: root.selected ? Theme.accent : Theme.textMuted
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }

        KvitLabel {
            anchors.left: mark.right
            anchors.leftMargin: Interface.space
            anchors.right: badge.left
            anchors.rightMargin: Interface.spaceNear
            anchors.verticalCenter: parent.verticalCenter
            visible: !root.collapsed
            text: root.text
            role: "body"
            font.bold: root.selected
            color: root.selected ? Theme.textPrimary : Theme.textSecondary
        }

        // The count, and where it goes once the label is gone.
        //
        // Hidden in the rail unless the caller asks for it, which is the
        // state every sidebar in the estate was drawing before `countInRail`
        // existed. A caller that turns it on gets the badge lifted off the
        // centre line and sitting on the symbol's upper corner, which is
        // where a count on an icon goes everywhere else — a backlog that
        // disappears when the sidebar collapses is a backlog the reader stops
        // acting on.
        KvitBadge {
            id: badge
            anchors.right: parent.right
            anchors.rightMargin: root.collapsed
                ? Interface.spaceTight : Interface.space
            anchors.verticalCenter: parent.verticalCenter
            // Lifted off the centre line in the rail, so it sits on the
            // symbol's upper corner rather than across the middle of it.
            // One anchor set in both states: switching an anchor line on and
            // off is how a component ends up unanchored at a size nobody
            // tested.
            anchors.verticalCenterOffset: root.collapsed ? -Interface.spaceNear : 0
            visible: root.count > 0
                     && (!root.collapsed || root.countInRail)
            count: root.count
            max: root.countMax
            // The item states the count in its own accessible name, so the
            // badge inside it says the same words rather than falling back to
            // "items" beside them. Both slots go across: a badge given only
            // the singular would inflect it on its own and disagree with the
            // item it sits in.
            counted: root.counted
            countedPlural: root.countedPlural
            tone: "neutral"
        }
    }
}
