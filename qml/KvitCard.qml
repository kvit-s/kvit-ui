// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A bounded block of content sitting on a surface: a dashboard widget, a
// summary, a group of related figures.
//
// The outline is `border` rather than `borderStrong` unless the card is
// something a reader can act on, in which case it is a control and takes the
// control-boundary token. `interactive` is what says which, and it also brings
// the hover tint and the focus ring, so a card that responds to a click looks
// like it will before it is clicked.
Rectangle {
    id: root

    property bool interactive: false
    property bool selected: false
    property alias hovered: hover.hovered
    default property alias content: holder.data
    property int padding: Interface.spaceLoose

    signal activated()

    implicitWidth: holder.childrenRect.width + padding * 2
    implicitHeight: holder.height + padding * 2

    radius: Interface.radiusCard
    color: root.selected ? Theme.selectionTint
         : (root.interactive && hover.hovered) ? Theme.hoverTint
         : Theme.listBackground
    border.width: Interface.hairline
    border.color: root.interactive ? Theme.borderStrong : Theme.border

    // The content holder takes its height from what is in it and its width
    // from the parent, rather than filling the parent in both directions.
    //
    // `anchors.fill` plus an `implicitHeight` read off `childrenRect` is a
    // cycle: the holder's height comes from this item and this item's from the
    // holder's children. Qt does not warn about it — it hangs during
    // construction, which is a much harder thing to find than a warning.
    Item {
        id: holder
        x: root.padding
        y: root.padding
        width: root.width - root.padding * 2
        height: childrenRect.height
    }

    HoverHandler {
        id: hover
        enabled: root.interactive
        cursorShape: Qt.PointingHandCursor
    }
    TapHandler {
        enabled: root.interactive
        onTapped: root.activated()
    }
}
