// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A number the reader types.
//
// Right-aligned and in tabular numerals, for the same reason KvitFigure is: a
// column of these is compared down the units digit.
//
// It validates as the reader types and reports through KvitField's `error`
// rather than silently refusing keystrokes. Refusing them is worse than it
// sounds: a field that ignores a key gives no reason, and the most common
// cause is a decimal separator the reader's locale writes differently from the
// one the field expects.
KvitField {
    id: root

    property real minimum: -Infinity
    property real maximum: Infinity
    // 0 for an integer field.
    property int decimals: 0
    // The parsed value, or NaN while the text is not a number.
    readonly property real value: Number(root.normalised)
    readonly property bool valid: !isNaN(value)
                                  && value >= minimum && value <= maximum

    // The text with the locale's decimal separator turned into a point, so
    // "3,50" from a German keyboard parses the same as "3.50".
    readonly property string normalised:
        root.text.trim().replace(Qt.locale().decimalPoint, ".")

    horizontalAlignment: TextInput.AlignRight
    font.features: ({ "tnum": 1 })
    inputMethodHints: root.decimals > 0 ? Qt.ImhFormattedNumbersOnly
                                        : Qt.ImhDigitsOnly

    error: {
        if (root.text.trim() === "")
            return ""
        if (isNaN(Number(root.normalised)))
            return qsTr("Not a number")
        if (root.value < root.minimum)
            return qsTr("Must be at least %1").arg(root.minimum)
        if (root.value > root.maximum)
            return qsTr("Must be at most %1").arg(root.maximum)
        return ""
    }

    // Round to the field's precision when the reader leaves it, rather than
    // as they type: reformatting under the caret moves it, which is the thing
    // that makes a "helpful" numeric field infuriating.
    onActiveFocusChanged: {
        if (!activeFocus && root.valid)
            root.text = root.value.toFixed(root.decimals)
    }
}
