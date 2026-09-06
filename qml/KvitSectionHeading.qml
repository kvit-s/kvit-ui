// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// A group heading: a solid bar carrying a disclosure chevron, the group's
// name, what it holds, and the one action that takes the whole group.
//
// kvit-hub's SectionHeading, with its reasoning intact.
//
// It is a filled bar rather than a label with a rule running off to the right,
// because a rule drawn between the words and the action reads as a divider
// between two unrelated things when it is meant to bind them into one heading.
//
// The action is hoisted here on purpose, and it has to be true for every row
// under the heading. If it is not true for all of them the group is wrong and
// wants splitting — which is the test that keeps a "hand all four to an agent"
// button from appearing over a group where it only applies to three.
//
// The count carries the word for what it counts, because a bare number at the
// end of a heading does not say what was counted. It is drawn whenever the
// caller gives one, a group of one included: a heading that reads "2 accounts"
// at two and nothing at all at one moves the count in and out of the bar as
// the group changes size, and a reader who has learnt to look at the end of
// the heading finds an empty space there instead of "1 account".
//
// The noun arrives in two slots, singular and plural, for the reason
// KvitViewHead and KvitBadge take two: one already-inflected word is right at
// exactly one count and wrong at every other.
Item {
    id: root

    property string text: ""
    property int count: -1
    // What the count counts, singular. "" leaves the number bare.
    property string counted: ""
    // The word for several of them. English adds an s and that is what this
    // falls back to; a noun that does not pluralise that way says so here.
    property string countedPlural: ""
    // The kind of thing the group holds, said beside the name where the name
    // alone does not say it.
    property string kind: ""
    // The hoisted action, true for every row in the group. Empty for none.
    property string action: ""
    property bool strong: false
    property bool expanded: true
    property bool collapsible: false

    signal toggled()
    signal actioned()

    // "1 account", "1,200 accounts", or the grouped number on its own where
    // the caller named nothing. Empty below zero, which is a caller saying it
    // has no count to give.
    //
    // The digits are grouped by the reader's locale, because a heading reading
    // "1200 accounts" sits in the same window as a ledger reading "1,200" and
    // a reader who notices the difference has to work out whether the two are
    // the same number. The number goes in through %1 rather than through Qt's
    // %n, which substitutes the bare integer and would undo the grouping,
    // while the count is still handed to qsTr so the choice of plural form
    // stays Qt's and a translator gets every form the language has.
    readonly property string countPhrase: {
        if (root.count < 0)
            return ""
        const grouped = Number(root.count).toLocaleString(Qt.locale(), 'f', 0)
        if (root.counted === "")
            return grouped
        if (root.count === 1) {
            return qsTr("%1 %2", "a count and the singular of what it counts",
                        root.count).arg(grouped).arg(root.counted)
        }
        if (root.countedPlural !== "") {
            return qsTr("%1 %2", "a count and the plural of what it counts",
                        root.count).arg(grouped).arg(root.countedPlural)
        }
        // No plural was given, so English suffixes an s. The whole phrase is
        // the translatable unit rather than the suffix on its own: a language
        // that pluralises some other way can rewrite this form, and there is
        // nothing a translator can do with a lone "s".
        return qsTr("%1 %2s", "a count and a noun pluralised by suffixing s",
                    root.count).arg(grouped).arg(root.counted)
    }

    implicitHeight: Interface.rowHeightCompact
    implicitWidth: Interface.px(400)

    Rectangle {
        anchors.fill: parent
        color: root.collapsible && headHover.hovered ? Theme.hoverTint
                                                     : Theme.panelBackground

        KvitDivider {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.bottom: parent.bottom
            color: root.strong ? Theme.borderStrong : Theme.border
        }

        HoverHandler {
            id: headHover
            enabled: root.collapsible
            cursorShape: Qt.PointingHandCursor
        }
        TapHandler {
            enabled: root.collapsible
            onTapped: root.toggled()
        }
    }

    Accessible.role: root.collapsible ? Accessible.Heading : Accessible.StaticText
    Accessible.name: root.text

    // A layout rather than a row of natural widths: the name, what the group
    // holds and the count are of unknown length together, and left to
    // themselves they run under the action at the right edge. The middle piece
    // gives way first, since it is the one a reader can lose without losing
    // which group this is.
    RowLayout {
        anchors.left: parent.left
        anchors.leftMargin: Interface.space
        anchors.right: parent.right
        anchors.rightMargin: Interface.space
        anchors.verticalCenter: parent.verticalCenter
        spacing: Interface.stackGap

        KvitIcon {
            visible: root.collapsible
            name: root.expanded ? "chevron-down" : "chevron-right"
            color: Theme.textFaint
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall
        }
        KvitLabel {
            text: root.text
            role: "small"
            font.bold: root.strong
        }
        KvitLabel {
            Layout.fillWidth: true
            visible: root.kind !== ""
            text: root.kind
            role: "small"
            color: Theme.textMuted
        }
        Item { Layout.fillWidth: root.kind === "" }
        KvitLabel {
            visible: root.countPhrase !== ""
            text: root.countPhrase
            role: "small"
            color: Theme.textFaint
            tabular: true
        }
        KvitLink {
            Layout.leftMargin: Interface.stackGap
            visible: root.action !== ""
            text: root.action
            role: "small"
            onActivated: root.actioned()
        }
    }
}
