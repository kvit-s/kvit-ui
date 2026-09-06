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
    // The word for several of them. English adds an s and that is what this
    // falls back to; a noun that does not pluralise that way says so here.
    property string countedPlural: ""
    property string subtitle: ""

    // "1 transaction", "250,000 transactions".
    //
    // The digits are grouped by the reader's locale, because a six-figure
    // count run together is read digit by digit, and the plural is a word
    // rather than a parenthesis: "250000 transaction(s)" is a form field
    // rather than a sentence, and no reader says it that way.
    //
    // The number goes in through %1 rather than through Qt's %n, which is
    // what keeps the grouping — %n substitutes the bare integer. The count is
    // still handed to qsTr, so which plural form to use stays Qt's choice and
    // a translator into a language with three of them gets all three; and the
    // source string is picked here as well, so that with no translator
    // installed the English is right rather than "1 transactions".
    readonly property string countPhrase: {
        if (root.count < 0)
            return ""
        const grouped = Number(root.count).toLocaleString(Qt.locale(), 'f', 0)
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
                    text: root.countPhrase
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
