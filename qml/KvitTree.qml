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

// A nested list the reader can open and close: folders, an outline, a
// hierarchy of accounts.
//
// Twelve files in the estate have a private version. What a shared one has to
// get right, and most hand-rolled ones do not, is the keyboard: Right opens a
// node, Left closes it or moves to the parent, and Home and End go to the ends.
// A tree without those is a tree a keyboard user has to arrow through one
// visible row at a time, and depth becomes a wall.
//
// Indentation is by `Interface.spaceLoose` per level rather than by a literal,
// so a deep tree stays proportionate when the interface size moves.
TreeView {
    id: root

    // What a screen reader calls this tree.
    property string label: ""

    // A hierarchy written out here, for a tree whose shape is fixed: account
    // groups, a category hierarchy, the sections of a settings page. Each
    // entry is a string, or an object with `label` and optional `children`.
    //
    // A TreeView takes a QAbstractItemModel and nothing else, and QML has no
    // tree model of its own, so without this every use of KvitTree — a
    // six-row settings outline included — starts with a C++ class. A tree over
    // something unbounded, kvit-notes' vault being the case that matters,
    // still wants one: set `model` and leave `nodes` alone.
    property var nodes: []

    TreeNodeModel {
        id: declaredNodes
        nodes: root.nodes === undefined ? [] : root.nodes
    }

    model: declaredNodes

    clip: true
    boundsBehavior: Flickable.StopAtBounds
    rowHeightProvider: () => Interface.rowHeightSlim
    // The same recycling KvitTable relies on. A folder tree over a large vault
    // is the other place in this estate where the row count is unbounded.
    reuseItems: true

    Accessible.role: Accessible.Tree
    Accessible.name: root.label

    delegate: Item {
        id: node
        required property TreeView treeView
        required property bool isTreeNode
        required property bool expanded
        required property bool hasChildren
        required property int depth
        required property int row
        required property var display

        implicitWidth: Interface.sidebarWidth
        implicitHeight: Interface.rowHeightSlim

        // `treeView` is null for a moment while a delegate is being torn down
        // or handed back to the reuse pool, and an unguarded read of it there
        // prints a TypeError for every row on screen.
        readonly property bool current: node.treeView
                                        && node.treeView.currentRow === node.row

        Rectangle {
            anchors.fill: parent
            color: node.current ? Theme.selectionTint
                 : hover.hovered ? Theme.hoverTint
                 : "transparent"
            HoverHandler { id: hover }
        }

        Accessible.role: Accessible.TreeItem
        Accessible.name: String(node.display)

        KvitIcon {
            id: chevron
            x: Interface.spaceSnug + node.depth * Interface.spaceLoose
            anchors.verticalCenter: parent.verticalCenter
            visible: node.hasChildren
            name: node.expanded ? "chevron-down" : "chevron-right"
            color: Theme.textMuted
            implicitWidth: Interface.iconSizeSmall
            implicitHeight: Interface.iconSizeSmall

            TapHandler {
                onTapped: node.treeView.toggleExpanded(node.row)
            }
        }

        KvitLabel {
            anchors.left: parent.left
            anchors.leftMargin: Interface.px(24) + node.depth * Interface.spaceLoose
            anchors.right: parent.right
            anchors.rightMargin: Interface.spaceNear
            anchors.verticalCenter: parent.verticalCenter
            text: String(node.display)
            role: "body"
            color: node.current ? Theme.textPrimary : Theme.textSecondary
        }

        TapHandler {
            onTapped: node.treeView.selectionModel.setCurrentIndex(
                node.treeView.index(node.row, 0),
                ItemSelectionModel.ClearAndSelect)
            onDoubleTapped: if (node.hasChildren)
                node.treeView.toggleExpanded(node.row)
        }
    }
}
