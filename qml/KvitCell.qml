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

// One cell of a table, drawn according to what kind of value it holds.
//
// The kind comes from the column rather than from the value, which is what
// makes a column readable: every cell in it is aligned the same way, set in
// the same face, and says the same thing about a missing value. A cell that
// decided for itself would give a column where some rows are right-aligned and
// some are not, depending on whether that row's value happened to parse as a
// number.
//
// The kinds are TableModelBase::CellKind, so the model says what a column is
// once and every view of it agrees.
Item {
    id: root

    // "Text" | "Figure" | "Chip" | "Slug" | "Date" | "Money"
    property string kind: "Text"
    property var value: undefined
    property bool measured: true
    // For a Chip cell: which tone. The model supplies it through MarkRole.
    property string mark: "neutral"
    property string unit: ""
    property bool selected: false

    implicitHeight: Interface.rowHeightSlim
    implicitWidth: content.implicitWidth

    readonly property bool numeric: kind === "Figure" || kind === "Money"

    Accessible.role: Accessible.Cell
    Accessible.name: root.measured ? String(root.value === undefined ? "" : root.value)
                                   : qsTr("not measured")

    Loader {
        id: content
        anchors.fill: parent
        anchors.leftMargin: Interface.spaceNear
        anchors.rightMargin: Interface.spaceNear
        sourceComponent: {
            switch (root.kind) {
            case "Figure": return figureCell
            case "Money":  return moneyCell
            case "Chip":   return chipCell
            case "Slug":   return slugCell
            case "Date":   return dateCell
            default:       return textCell
            }
        }
    }

    Component {
        id: textCell
        KvitLabel {
            text: root.value === undefined || root.value === null
                  ? "" : String(root.value)
            role: "body"
            color: root.selected ? Theme.textPrimary : Theme.textSecondary
            verticalAlignment: Text.AlignVCenter
        }
    }

    Component {
        id: figureCell
        // Right-aligned, because a column of numbers is compared down its
        // units digit and a left-aligned column of numbers cannot be.
        Item {
            KvitFigure {
                anchors.right: parent.right
                anchors.verticalCenter: parent.verticalCenter
                value: root.value === undefined || root.value === null
                       ? "" : String(root.value)
                unit: root.unit
                measured: root.measured
                role: "body"
            }
        }
    }

    Component {
        id: moneyCell
        // The money variant follows money-display.md: minor units always
        // shown, so 4 and 4.00 do not sit in the same column looking like
        // different precisions, and the currency beside the number rather
        // than in a heading, because a table can hold more than one.
        Item {
            KvitFigure {
                anchors.right: parent.right
                anchors.verticalCenter: parent.verticalCenter
                value: {
                    const v = root.value
                    if (v === undefined || v === null)
                        return ""
                    return Number(v).toFixed(2)
                }
                unit: root.unit
                measured: root.measured
                role: "body"
                color: Number(root.value) < 0 ? Theme.danger : Theme.textPrimary
            }
        }
    }

    Component {
        id: chipCell
        Item {
            KvitChip {
                anchors.left: parent.left
                anchors.verticalCenter: parent.verticalCenter
                text: root.value === undefined || root.value === null
                      ? "" : String(root.value)
                tone: root.mark
            }
        }
    }

    Component {
        id: slugCell
        KvitSlug {
            anchors.verticalCenter: parent.verticalCenter
            text: root.value === undefined || root.value === null
                  ? "" : String(root.value)
            // No ground inside a table: a box on every row of a column reads
            // as a column of boxes rather than as a column of identifiers.
            ground: false
        }
    }

    Component {
        id: dateCell
        KvitLabel {
            text: {
                const v = root.value
                if (v === undefined || v === null)
                    return ""
                return v instanceof Date
                    ? Qt.formatDate(v, Locale.ShortFormat) : String(v)
            }
            role: "body"
            color: Theme.textSecondary
            tabular: true
            verticalAlignment: Text.AlignVCenter
        }
    }
}
