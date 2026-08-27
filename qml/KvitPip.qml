// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// A row of dots standing for a small count: three of five days recorded, two
// of four checks passed.
//
// For counts a reader can take in at a glance without counting, which in
// practice is up to about seven. Past that the row stops being faster to read
// than the number and KvitFigure is the right component.
//
// Filled and empty rather than a filled prefix on a bar, because a pip row
// says "how many of these" while a bar says "how far along", and a reader
// asked the first question should not have to work out the denominator.
Row {
    id: root

    property int filled: 0
    property int total: 5
    property color color: Theme.accent
    property string label: ""

    spacing: Interface.spaceTight

    Accessible.role: Accessible.ProgressBar
    Accessible.name: root.label !== "" ? root.label
                                       : qsTr("%1 of %2").arg(root.filled).arg(root.total)

    Repeater {
        model: Math.max(0, root.total)
        delegate: KvitDot {
            required property int index
            color: index < root.filled ? root.color : Theme.border
            hollow: index >= root.filled
            width: Interface.spaceSnug
            height: Interface.spaceSnug
        }
    }
}
