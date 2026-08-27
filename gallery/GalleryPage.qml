// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// One component's page: what it is, what it is for, and each of its states.
Column {
    id: root

    required property var entry

    spacing: Interface.spaceLoose
    width: parent ? parent.width : Interface.px(700)

    KvitViewHead {
        width: parent.width
        title: root.entry.name
        subtitle: root.entry.summary
    }

    Repeater {
        model: root.entry.specimens
        delegate: Specimen {
            required property var modelData
            width: root.width
            caption: modelData.caption
            snippet: modelData.snippet
            extraImports: modelData.extraImports === undefined
                          ? [] : modelData.extraImports
        }
    }
}
