// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import QtQuick.Controls
import Kvit.Ui

// A choice from a list that is too long to lay out: a font family, a currency,
// a folder.
//
// For a short list of mutually exclusive options that are worth seeing at
// once, KvitSegmented is the right component — a three-item dropdown hides two
// of the three answers behind a click for no reason.
ComboBox {
    id: root

    property string label: ""

    implicitHeight: Interface.controlHeight
    implicitWidth: Interface.px(180)
    font.pixelSize: Interface.body

    Accessible.role: Accessible.ComboBox
    Accessible.name: root.label

    delegate: ItemDelegate {
        id: option
        required property var model
        required property int index
        width: root.width
        height: Interface.rowHeightSlim
        highlighted: root.highlightedIndex === option.index

        background: Rectangle {
            color: option.highlighted ? Theme.hoverTint : "transparent"
        }
        contentItem: KvitLabel {
            text: option.model[root.textRole] !== undefined
                  ? option.model[root.textRole] : option.model
            role: "body"
            leftPadding: Interface.spaceNear
        }
    }

    indicator: KvitIcon {
        x: root.width - width - Interface.spaceNear
        y: (root.height - height) / 2
        name: "chevron-down"
        color: Theme.textMuted
        // 13 px rather than 18: an indicator drawn at the icon default
        // overpowers the words beside it.
        implicitWidth: Interface.iconSizeSmall
        implicitHeight: Interface.iconSizeSmall
    }

    contentItem: KvitLabel {
        leftPadding: Interface.spaceNear
        rightPadding: root.indicator.width + Interface.spaceLoose
        text: root.displayText
        role: "body"
        color: root.enabled ? Theme.textPrimary : Theme.textDisabled
    }

    background: Rectangle {
        radius: Interface.radiusControl
        color: root.enabled ? Theme.popupBackground : Theme.chipBackground
        border.width: Interface.hairline
        border.color: root.activeFocus ? Theme.focusRing : Theme.borderStrong

        Rectangle {
            anchors.fill: parent
            anchors.margins: -Interface.focusRingWidth
            visible: root.visualFocus
            color: "transparent"
            radius: parent.radius + Interface.focusRingWidth
            border.width: Interface.focusRingWidth
            border.color: Theme.focusRing
        }
    }

    popup: Popup {
        y: root.height
        width: root.width
        implicitHeight: Math.min(contentItem.implicitHeight + Interface.spaceSnug * 2,
                                 Interface.px(320))
        padding: Interface.spaceSnug

        contentItem: ListView {
            clip: true
            implicitHeight: contentHeight
            model: root.delegateModel
            currentIndex: root.highlightedIndex
            boundsBehavior: Flickable.StopAtBounds
        }

        background: Rectangle {
            radius: Interface.radiusControl
            color: Theme.popupBackground
            border.width: Interface.hairline
            border.color: Theme.borderStrong
        }
    }
}
