// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// An amount of money.
//
// A number field with three things money needs and a plain number does not.
//
// It works in minor units. The reader types 12.34 and `minorUnits` is 1234,
// an integer, which is what a ledger stores: doing money in floating point
// means 0.1 + 0.2 is not 0.3 and a balance drifts by a penny somewhere nobody
// can find.
//
// The number of decimal places comes from the currency rather than being two.
// Yen has none and Bahraini dinar has three, and a field that assumes two
// silently divides or multiplies those by a hundred.
//
// The currency is beside the field rather than in a column heading, because a
// ledger can hold more than one and a heading that says "Amount (GBP)" over a
// column containing euros is the failure money-display.md is written against.
KvitNumberField {
    id: root

    // An ISO 4217 code. Empty draws no marker, for a field in a context where
    // the currency is already unambiguous.
    property string currency: ""
    // How many decimal places this currency has. Two for most, zero for yen,
    // three for dinar.
    property int minorDigits: 2
    // The amount as an integer number of minor units — pence, cents, sen.
    // NaN while the text is not a number.
    readonly property real minorUnits:
        root.valid ? Math.round(root.value * Math.pow(10, root.minorDigits))
                   : NaN

    decimals: root.minorDigits
    rightPadding: currency !== ""
        ? marker.width + Interface.spaceNear * 2 : Interface.spaceNear

    Accessible.role: Accessible.EditableText
    Accessible.name: root.currency !== ""
        ? qsTr("%1 in %2").arg(root.label).arg(root.currency) : root.label

    KvitLabel {
        id: marker
        anchors.right: parent.right
        anchors.rightMargin: Interface.spaceNear
        anchors.verticalCenter: parent.verticalCenter
        visible: root.currency !== ""
        text: root.currency
        role: "small"
        color: Theme.textMuted
        elide: Text.ElideNone
    }
}
