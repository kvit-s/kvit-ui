// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A field that filters something.
//
// Separate from KvitField because it behaves differently rather than because
// it looks different: it carries a magnifier so it is findable without reading
// the placeholder, a clear button that appears only when there is something to
// clear, and Escape clears rather than reverting — which is what a reader
// expects of a filter and not of a text field.
//
// It also announces its result count, because filtering is the one interaction
// where the whole outcome happens somewhere else on the screen: a screen reader
// user who types into an unannounced filter gets no feedback at all.
KvitField {
    id: root

    // How many things the filter left. Negative leaves it unsaid, which is
    // right for a filter whose result is not a countable list.
    property int matches: -1
    property string matchedNoun: qsTr("result")
    // The word for several of them. English adds an s and that is what this
    // falls back to; a noun that does not pluralise that way says so here.
    property string matchedNounPlural: ""

    // "1 transaction", "250,000 transactions" — what a screen reader is told
    // the filter left.
    //
    // The digits are grouped by the reader's locale, because a six-figure
    // count run together is read digit by digit, and the plural is a word
    // rather than a parenthesis: "250000 transaction(s)" is a form field and
    // no reader says it that way.
    //
    // The number goes in through %1 rather than through Qt's %n, which is
    // what keeps the grouping — %n substitutes the bare integer. The count is
    // still handed to qsTr, so which plural form to use stays Qt's choice and
    // a translator into a language with three of them gets all three; and the
    // source string is picked here as well, so that with no translator
    // installed the English is right rather than "1 results".
    readonly property string matchPhrase: {
        if (root.matches < 0)
            return ""
        const grouped = Number(root.matches).toLocaleString(Qt.locale(), 'f', 0)
        if (root.matches === 1) {
            return qsTr("%1 %2", "a count and the singular of what it counts",
                        root.matches).arg(grouped).arg(root.matchedNoun)
        }
        if (root.matchedNounPlural !== "") {
            return qsTr("%1 %2", "a count and the plural of what it counts",
                        root.matches).arg(grouped).arg(root.matchedNounPlural)
        }
        // No plural was given, so English suffixes an s. The whole phrase is
        // the translatable unit rather than the suffix on its own: a language
        // that pluralises some other way can rewrite this form, and there is
        // nothing a translator can do with a lone "s".
        return qsTr("%1 %2s", "a count and a noun pluralised by suffixing s",
                    root.matches).arg(grouped).arg(root.matchedNoun)
    }

    leftPadding: Interface.controlHeight
    rightPadding: root.text !== "" ? Interface.controlHeight : Interface.spaceNear
    placeholderText: qsTr("Filter…")

    Accessible.role: Accessible.EditableText
    Accessible.name: root.label !== "" ? root.label : qsTr("Filter")
    Accessible.description: root.matchPhrase

    Keys.onEscapePressed: function (event) {
        if (root.text !== "") {
            root.clear()
            event.accepted = true
        }
    }

    KvitIcon {
        anchors.left: parent.left
        anchors.leftMargin: Interface.spaceNear
        anchors.verticalCenter: parent.verticalCenter
        name: "search"
        color: Theme.textFaint
        implicitWidth: Interface.iconSizeSmall
        implicitHeight: Interface.iconSizeSmall
    }

    KvitIconButton {
        anchors.right: parent.right
        anchors.rightMargin: Interface.spaceTight
        anchors.verticalCenter: parent.verticalCenter
        visible: root.text !== ""
        symbol: "close"
        label: qsTr("Clear the filter")
        implicitWidth: Interface.px(20)
        implicitHeight: Interface.px(20)
        onClicked: root.clear()
    }
}
