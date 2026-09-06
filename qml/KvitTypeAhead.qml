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

// A field that offers matches as the reader types.
//
// prd.md §5.6 Group B: kvit-cash's category and tag pickers need one. The
// difference from KvitSelect is what the reader knows — a select is for
// choosing from a list they will read, and this is for a list too long to
// read where they already know roughly what they want.
//
// `allowNew` is the property that matters. A tag picker should let the reader
// make a tag that does not exist yet; a category picker usually should not,
// because a typo would silently create a second category next to the right
// one. It has no default: the caller has to decide.
Item {
    id: root

    // Everything that could be matched: strings, or { value, label } objects.
    property var source: []
    property string text: ""
    property string label: ""
    property string placeholder: ""
    required property bool allowNew
    property int maximumSuggestions: 8

    signal chosen(var value)

    implicitHeight: Interface.controlHeight
    implicitWidth: Interface.px(220)

    function labelOf(item) {
        return item !== null && typeof item === "object" ? item.label : item
    }
    function valueOf(item) {
        return item !== null && typeof item === "object" ? item.value : item
    }

    readonly property var matches: {
        const needle = root.text.trim().toLowerCase()
        if (needle === "")
            return []
        const found = []
        for (let i = 0; i < root.source.length
                        && found.length < root.maximumSuggestions; ++i) {
            if (String(root.labelOf(root.source[i])).toLowerCase()
                    .indexOf(needle) >= 0)
                found.push(root.source[i])
        }
        return found
    }

    KvitField {
        id: field
        anchors.fill: parent
        text: root.text
        label: root.label
        placeholderText: root.placeholder
        onTextChanged: root.text = text

        Accessible.role: Accessible.EditableText
        Accessible.name: root.label
        // What a screen reader has no other way to learn: that a list appeared
        // under the field and how long it is.
        // Two source strings picked by the count rather than one with `%n`
        // in it: with no translator installed Qt substitutes the number and
        // chooses no plural form, so one match is announced as
        // "1 suggestion(s)".
        Accessible.description: {
            const many = root.matches.length
            if (many === 0)
                return ""
            return many === 1 ? qsTr("%n suggestion", "", many)
                              : qsTr("%n suggestions", "", many)
        }

        Keys.onDownPressed: if (suggestions.opened) list.incrementCurrentIndex()
        Keys.onUpPressed: if (suggestions.opened) list.decrementCurrentIndex()
        Keys.onReturnPressed: root.accept()
        Keys.onEnterPressed: root.accept()
    }

    function accept() {
        if (list.currentIndex >= 0 && list.currentIndex < root.matches.length) {
            const picked = root.matches[list.currentIndex]
            root.text = String(root.labelOf(picked))
            root.chosen(root.valueOf(picked))
        } else if (root.allowNew && root.text.trim() !== "") {
            root.chosen(root.text.trim())
        }
        suggestions.close()
    }

    KvitPopover {
        id: suggestions
        y: root.height
        width: root.width
        padding: Interface.spaceSnug
        // Opened by typing rather than by a click, and it takes no focus: the
        // caret must stay in the field or the next keystroke goes nowhere.
        focus: false
        visible: root.matches.length > 0 || (root.allowNew && root.text !== "")
        implicitHeight: Math.min(list.contentHeight + newRow.height
                                 + Interface.space, Interface.px(240))

        Column {
            width: parent.width

            ListView {
                id: list
                width: parent.width
                height: Math.min(contentHeight, Interface.px(200))
                model: root.matches
                currentIndex: 0
                clip: true
                boundsBehavior: Flickable.StopAtBounds

                delegate: KvitRow {
                    required property var modelData
                    required property int index
                    width: list.width
                    form: "slim"
                    rule: false
                    current: index === list.currentIndex
                    label: String(root.labelOf(modelData))
                    // The list moves its own cursor with the arrow keys, so
                    // the row is not in the tab order and says here that it
                    // acts.
                    interactive: true
                    onActivated: {
                        root.text = String(root.labelOf(modelData))
                        root.chosen(root.valueOf(modelData))
                        suggestions.close()
                    }

                    KvitLabel {
                        anchors.fill: parent
                        anchors.leftMargin: Interface.spaceNear
                        text: String(root.labelOf(parent.modelData))
                        role: "body"
                    }
                }
            }

            // The "make a new one" row, when the caller allows it. It says
            // what it will make, in quotes, because a row reading "Create" at
            // the bottom of a list of matches is ambiguous about which of them
            // it applies to.
            KvitRow {
                id: newRow
                width: parent.width
                visible: root.allowNew && root.text.trim() !== ""
                form: "slim"
                rule: false
                height: visible ? implicitHeight : 0
                label: qsTr("Create \"%1\"").arg(root.text.trim())
                interactive: true
                onActivated: {
                    root.chosen(root.text.trim())
                    suggestions.close()
                }

                Row {
                    anchors.fill: parent
                    anchors.leftMargin: Interface.spaceNear
                    spacing: Interface.spaceNear
                    KvitIcon {
                        anchors.verticalCenter: parent.verticalCenter
                        name: "plus"
                        color: Theme.accent
                        implicitWidth: Interface.iconSizeSmall
                        implicitHeight: Interface.iconSizeSmall
                    }
                    KvitLabel {
                        anchors.verticalCenter: parent.verticalCenter
                        text: qsTr("Create \"%1\"").arg(root.text.trim())
                        role: "body"
                        color: Theme.accent
                        elide: Text.ElideNone
                    }
                }
            }
        }
    }
}
