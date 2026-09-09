// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A count attached to something else: unread items, pending changes, matches.
//
// A badge is always a number and the number is always small on the screen and
// possibly large in value, so it caps: past `max` it says "99+" rather than
// growing wide enough to shift the thing it is attached to.
//
// Zero hides it rather than drawing a nought. A badge showing 0 is a mark that
// says "look here" about nothing.
//
// Its gallery page is shot at the maximum interface size as well as the
// default, which is the check that the cap and the pill hold together at
// twice the type size in the high-contrast theme. That is where a badge would
// fail if it were going to: three digits in a pill sized from `pillHeight`,
// with a label whose colour comes from `Theme.labelOn` rather than from a
// fixed foreground. The run of 2026-09-08 draws 214 and 99+ whole in all four
// themes at 24 px, so there is nothing to change here; the shotSizes entry in
// the catalogue is what makes that check happen again rather than once.
Rectangle {
    id: root

    property int count: 0
    property int max: 99
    // The word for one of the things counted. Empty gives the generic
    // "11 items"; a caller that knows the noun sets it and the badge is
    // announced as "1 decision" or "11 decisions".
    property string counted: ""
    // The word for several of them. English adds an s and that is what this
    // falls back to; a noun that does not pluralise that way says so here.
    // Two slots rather than one already-inflected word, because a single word
    // is right at one count and wrong at every other: a caller that hands
    // over "decisions" announces a workspace with one decision waiting as
    // "1 decisions".
    property string countedPlural: ""
    // "neutral" | "accent" | "danger"
    property string tone: "accent"

    // The count with the reader's own digit grouping. Drawn as well as
    // announced: a badge reading 250000 a few pixels above a line reading
    // "250,000 records included" is the same number written two ways in one
    // window, and a reader who notices the difference has to work out whether
    // they are the same number.
    readonly property string groupedCount:
        Number(root.count).toLocaleString(Qt.locale(), 'f', 0)

    readonly property color toneColor: {
        switch (tone) {
        case "danger": return Theme.danger
        case "neutral": return Theme.textMuted
        default: return Theme.accent
        }
    }

    visible: count > 0
    implicitHeight: Interface.pillHeight
    implicitWidth: Math.max(implicitHeight, label.implicitWidth + Interface.spaceNear)
    radius: height / 2
    color: root.toneColor

    Accessible.role: Accessible.StaticText
    // The cap belongs to the drawn badge and not to what is announced: the
    // pill is small enough that "99+" is the only thing that fits, while a
    // screen reader has room for the number and is the one reader who cannot
    // go and look it up somewhere else.
    //
    // The digits are grouped by the reader's locale, because a six-figure
    // count run together is read out digit by digit, and the plural is a word
    // rather than a parenthesis: "250000 item(s)" is a form field rather than
    // a sentence. The number goes in through %1 rather than Qt's %n, which is
    // what keeps the grouping, while the count is still handed to qsTr so the
    // choice of plural form stays Qt's and a translator gets every form.
    Accessible.name: {
        const grouped = root.groupedCount
        if (root.counted !== "") {
            if (root.count === 1) {
                return qsTr("%1 %2", "a count and the singular of what it counts",
                            root.count).arg(grouped).arg(root.counted)
            }
            if (root.countedPlural !== "") {
                return qsTr("%1 %2", "a count and the plural of what it counts",
                            root.count).arg(grouped).arg(root.countedPlural)
            }
            // No plural was given, so English suffixes an s. The whole phrase
            // is the translatable unit rather than the suffix on its own: a
            // language that pluralises some other way can rewrite this form,
            // and there is nothing a translator can do with a lone "s".
            return qsTr("%1 %2s", "a count and a noun pluralised by suffixing s",
                        root.count).arg(grouped).arg(root.counted)
        }
        if (root.count === 1)
            return qsTr("%1 item", "a count of one thing", root.count).arg(grouped)
        return qsTr("%1 items", "a count of several things",
                    root.count).arg(grouped)
    }

    KvitLabel {
        id: label
        anchors.centerIn: parent
        text: root.count > root.max
              ? Number(root.max).toLocaleString(Qt.locale(), 'f', 0) + "+"
              : root.groupedCount
        role: "caption"
        color: Theme.labelOn(root.toneColor)
        tabular: true
        elide: Text.ElideNone
    }
}
