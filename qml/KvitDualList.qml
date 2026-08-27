// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// Two lists with items moving between them: what is available, and what is
// chosen and in what order.
//
// prd.md §5.6 Group B names one caller — configuring which columns a table
// shows and in what order. That is the shape this is for, and it is why the
// chosen side is ordered and the available side is not: column order matters
// and which columns exist does not.
//
// Everything works from the keyboard. A dual list where the only route is
// dragging between the two panes is unusable without a mouse, which is the
// state most implementations of this are in.
Item {
    id: root

    // Strings, or { value, label } objects.
    property var available: []
    property var chosen: []
    property string availableLabel: qsTr("Available")
    property string chosenLabel: qsTr("Shown, in order")

    signal changed(var available, var chosen)

    implicitWidth: Interface.px(480)
    implicitHeight: Interface.px(280)

    function labelOf(item) {
        return item !== null && typeof item === "object" ? item.label : item
    }

    function move(fromList, toList, index) {
        const from = fromList.slice()
        const to = toList.slice()
        to.push(from.splice(index, 1)[0])
        return [from, to]
    }

    RowLayout {
        anchors.fill: parent
        spacing: Interface.columnGap

        // The available side.
        ColumnLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            spacing: Interface.spaceNear

            KvitLabel {
                text: root.availableLabel
                role: "small"
                color: Theme.textMuted
            }
            KvitPanel {
                Layout.fillWidth: true
                Layout.fillHeight: true
                ruleTop: true
                ruleBottom: true
                ruleLeft: true
                ruleRight: true

                ListView {
                    id: leftList
                    anchors.fill: parent
                    clip: true
                    model: root.available
                    boundsBehavior: Flickable.StopAtBounds

                    delegate: KvitRow {
                        id: leftRow
                        required property var modelData
                        required property int index
                        width: leftList.width
                        form: "slim"
                        rule: false
                        current: leftList.currentIndex === leftRow.index
                        label: String(root.labelOf(leftRow.modelData))
                        onActivated: {
                            leftList.currentIndex = leftRow.index
                            const moved = root.move(root.available, root.chosen,
                                                    leftRow.index)
                            root.changed(moved[0], moved[1])
                        }

                        KvitLabel {
                            anchors.fill: parent
                            anchors.leftMargin: Interface.spaceNear
                            text: String(root.labelOf(leftRow.modelData))
                            role: "body"
                        }
                    }
                }
            }
        }

        // The buttons between. Symbols rather than words, because "Add" and
        // "Remove" are ambiguous about direction here — each button says which
        // way with an arrow and says what it does in its accessible label.
        ColumnLayout {
            Layout.alignment: Qt.AlignVCenter
            spacing: Interface.spaceNear

            KvitIconButton {
                symbol: "arrow-right"
                label: qsTr("Show the selected column")
                form: "ordinary"
                enabled: leftList.currentIndex >= 0
                         && leftList.currentIndex < root.available.length
                onClicked: {
                    const moved = root.move(root.available, root.chosen,
                                            leftList.currentIndex)
                    root.changed(moved[0], moved[1])
                }
            }
            KvitIconButton {
                symbol: "arrow-left"
                label: qsTr("Hide the selected column")
                form: "ordinary"
                enabled: rightList.currentIndex >= 0
                         && rightList.currentIndex < root.chosen.length
                onClicked: {
                    const moved = root.move(root.chosen, root.available,
                                            rightList.currentIndex)
                    root.changed(moved[1], moved[0])
                }
            }
            KvitIconButton {
                symbol: "arrow-up"
                label: qsTr("Move the selected column earlier")
                form: "ordinary"
                enabled: rightList.currentIndex > 0
                onClicked: {
                    const list = root.chosen.slice()
                    const at = rightList.currentIndex
                    const item = list.splice(at, 1)[0]
                    list.splice(at - 1, 0, item)
                    rightList.currentIndex = at - 1
                    root.changed(root.available, list)
                }
            }
            KvitIconButton {
                symbol: "arrow-down"
                label: qsTr("Move the selected column later")
                form: "ordinary"
                enabled: rightList.currentIndex >= 0
                         && rightList.currentIndex < root.chosen.length - 1
                onClicked: {
                    const list = root.chosen.slice()
                    const at = rightList.currentIndex
                    const item = list.splice(at, 1)[0]
                    list.splice(at + 1, 0, item)
                    rightList.currentIndex = at + 1
                    root.changed(root.available, list)
                }
            }
        }

        // The chosen side, in order.
        ColumnLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            spacing: Interface.spaceNear

            KvitLabel {
                text: root.chosenLabel
                role: "small"
                color: Theme.textMuted
            }
            KvitPanel {
                Layout.fillWidth: true
                Layout.fillHeight: true
                ruleTop: true
                ruleBottom: true
                ruleLeft: true
                ruleRight: true

                ListView {
                    id: rightList
                    anchors.fill: parent
                    clip: true
                    model: root.chosen
                    boundsBehavior: Flickable.StopAtBounds

                    delegate: KvitRow {
                        id: rightRow
                        required property var modelData
                        required property int index
                        width: rightList.width
                        form: "slim"
                        rule: false
                        current: rightList.currentIndex === rightRow.index
                        label: String(root.labelOf(rightRow.modelData))
                        onActivated: rightList.currentIndex = rightRow.index

                        Row {
                            anchors.fill: parent
                            anchors.leftMargin: Interface.spaceNear
                            spacing: Interface.spaceNear
                            KvitLabel {
                                anchors.verticalCenter: parent.verticalCenter
                                text: String(rightRow.index + 1)
                                role: "caption"
                                color: Theme.textFaint
                                tabular: true
                                elide: Text.ElideNone
                            }
                            KvitLabel {
                                anchors.verticalCenter: parent.verticalCenter
                                text: String(root.labelOf(rightRow.modelData))
                                role: "body"
                                elide: Text.ElideNone
                            }
                        }
                    }
                }
            }
        }
    }
}
