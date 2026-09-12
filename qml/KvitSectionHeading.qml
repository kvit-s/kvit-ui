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
//
// A collapsible heading opens and closes from the keyboard as well as the
// pointer: it joins the tab order, draws the ring every other control draws,
// and answers Return, Enter, Space and the press action a screen reader
// offers. The hoisted action beside it is a control of its own — a KvitLink
// showing its words, or a KvitIconButton showing its symbol where the caller
// named one — and either consumes its own keys, so the key that runs the
// action never also collapses the group behind it. The heading answers a key
// only while the heading itself holds the keyboard, which is the same rule
// KvitRow states at length.
Item {
    id: root

    property string text: ""
    property int count: -1
    // What the count counts, singular. "" leaves the number bare.
    property string counted: ""
    // The word for several of them. English adds an s and that is what this
    // falls back to; a noun that does not pluralise that way says so here.
    property string countedPlural: ""
    // A count the caller has already written out, for a heading whose count is
    // not a number.
    //
    // kvit-notes-pro's Changes heading reads "1 · +0 −0": one changed file,
    // and the lines added and removed across it. There is no integer that
    // says that, and the service the screen reads it from has already
    // composed the phrase, so `count` and `counted` have nothing to work
    // with. Anything given here is drawn as written — no grouping, no
    // pluralising, no locale — because the caller has already decided all
    // three.
    //
    // It sits beside the name rather than at the right end, which is where an
    // integer count goes. The two are in different places because they are
    // different things to a reader: a tally of how many rows are under the
    // heading is something to glance at, and it is the same glance in every
    // heading down the column, whereas a phrase the caller composed is part
    // of what the group is and reads as an extension of its name. Setting
    // this leaves the right end empty, so a heading cannot show two counts.
    property string countText: ""
    // The kind of thing the group holds, said beside the name where the name
    // alone does not say it.
    property string kind: ""
    // The hoisted action, true for every row in the group. Empty for none.
    property string action: ""
    // Draw that action as a symbol instead of as its words.
    //
    // A heading whose action is "Open the diff for all of them" spends more of
    // the bar on the action than on the name of the group, and a column of
    // eight headings each ending in a different sentence is a column a reader
    // has to read rather than scan. A symbol is the same action in the width
    // of one glyph. The words do not go anywhere: they stay in `action`, which
    // becomes the button's accessible name and its tooltip, so the action is
    // still reachable and still says what it does.
    //
    // Only for an action a symbol can actually carry — a diff, a plus, a
    // terminal. An action with no obvious symbol keeps its words, which is
    // what leaving this empty does.
    property string actionSymbol: ""
    property bool strong: false
    property bool expanded: true
    property bool collapsible: false

    signal toggled()
    signal actioned()

    // A heading that can be opened is something to arrive at. One that cannot
    // is a line of text, and putting it in the tab order would give a reader
    // a stop that does nothing.
    activeFocusOnTab: root.collapsible

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
        // A phrase the caller wrote is the count, and it is drawn beside the
        // name. Nothing goes at the right end as well.
        if (root.countText !== "")
            return ""
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
            anchors.top: parent.top
            color: root.strong ? Theme.borderStrong : Theme.border
        }

        // Drawn inside the bar rather than outside it, because a heading sits
        // flush against the rows above and below and a ring drawn outside
        // would be painted over by whichever of them is drawn last.
        Rectangle {
            objectName: "focusRing"
            anchors.fill: parent
            anchors.margins: Interface.focusRingWidth
            visible: root.activeFocus
            color: "transparent"
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
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
    // Whether the group is open, which the chevron says to everybody else.
    // Empty on a heading that does not collapse: there is no state to report
    // and a reader told "collapsed" about a fixed heading would go looking
    // for the way to open it.
    Accessible.description: !root.collapsible ? ""
                          : root.expanded ? qsTr("Expanded") : qsTr("Collapsed")
    Accessible.onPressAction: {
        if (root.collapsible)
            root.toggled()
    }

    // The same three keys KvitRow answers, under the same rule: only while
    // the heading itself has the keyboard, so that Space on the focused
    // action runs the action and leaves the group open.
    Keys.onPressed: event => {
        if (!root.collapsible || !root.activeFocus)
            return
        switch (event.key) {
        case Qt.Key_Return:
        case Qt.Key_Enter:
        case Qt.Key_Space:
            root.toggled()
            event.accepted = true
            break
        default:
            break
        }
    }

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
            objectName: "countText"
            visible: root.countText !== ""
            text: root.countText
            role: "small"
            color: Theme.textFaint
            tabular: true
            elide: Text.ElideNone
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
            objectName: "action"
            Layout.leftMargin: Interface.stackGap
            visible: root.action !== "" && root.actionSymbol === ""
            text: root.action
            role: "small"
            onActivated: root.actioned()
        }
        KvitIconButton {
            objectName: "actionButton"
            Layout.leftMargin: Interface.stackGap
            // The bar is one compact row tall and the button draws its focus
            // ring outside its own ground, so anything taller than the row
            // less a ring on each side has its ring cut off by the rows above
            // and below.
            Layout.preferredWidth: Interface.rowHeightCompact
                                   - Interface.focusRingWidth * 2
            Layout.preferredHeight: Interface.rowHeightCompact
                                    - Interface.focusRingWidth * 2
            visible: root.action !== "" && root.actionSymbol !== ""
            // `dot` while there is no symbol to draw: this button is hidden
            // then, and a KvitIcon given a name it does not know draws its
            // marked placeholder and writes a warning whether or not anybody
            // can see it.
            symbol: root.actionSymbol === "" ? "dot" : root.actionSymbol
            label: root.action
            dense: true
            onClicked: root.actioned()
        }
    }
}
