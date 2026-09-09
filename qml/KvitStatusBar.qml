// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// The group and fact delegates below are nested components, and Bound is what
// makes them read this file's ids legally rather than by accident: without it
// a nested component resolves an outer id at run time through the object
// hierarchy.
pragma ComponentBehavior: Bound

import QtQml
import QtQuick
import QtQuick.Layouts
import Kvit.Ui

// The strip along the bottom of a window: what the application is doing now,
// and the standing facts about what is open.
//
// kvit-hub's WindowStatusBar generalised. The left says what is happening and
// changes; the right says what is true and does not, and mixing the two is
// what makes a status bar unreadable.
//
// The right-hand side takes facts in three shapes, in increasing order of how
// much they do.
//
//   facts     plain strings, separated by dots. Nothing to press.
//   groups    named groups of facts that open what they name. Each fact is a
//             symbol and a label and is a control: it takes tab focus, says
//             its own name to a screen reader, and answers Return and Space.
//   controls  one or two whole controls at the end, for an action that is not
//             a fact at all — the undo that reverses the last change.
//
// A group is `{ label, facts: [{ text, symbol }] }`, and pressing a fact
// emits `factActivated(group, fact)` with the two indices. The bar holds no
// state of the caller's: it is handed words and hands back which one was
// pressed.
//
// ── What happens when they do not fit ──────────────────────────────────────
//
// A bar is one line at the bottom of a window that can be any width, so the
// facts on it will eventually not fit. The behaviour a consumer reaches for
// on its own is to cut the list — kvit-notes-pro sliced each of its two
// groups to the first two items — and the cost is that the third thing
// waiting for somebody is not on the screen and nothing says it exists.
//
// So the bar shows the groups it has room for, keeps the rest in a menu
// behind a control that says how many there are, and that control is a
// control: reachable by tab, opened by Return, and navigated with the arrow
// keys once open. Nothing is dropped and nothing is silent.
//
// The measurement is deliberate about one thing. A group that does not fit is
// collapsed to no width rather than hidden, because a positioner measures
// only its visible children: a hidden group would report a width of zero, the
// bar would decide it fits after all, and the two answers would alternate.
// Collapsed and disabled, it keeps a width to be measured by and leaves the
// tab order, which is the pair of properties that matter.
Rectangle {
    id: root

    // What is happening now. Empty leaves the left side blank, which is the
    // resting state — a status bar that always has something to say trains the
    // reader to stop looking at it.
    property string activity: ""
    // Standing facts, right aligned. Each becomes one label with a separator
    // between.
    property var facts: []
    // Groups of facts that do something. Each entry is
    // `{ label, facts: [{ text, symbol }] }`; `label` names what the group is
    // about — "Waiting on you", "Running" — and may be empty for a group that
    // needs no heading.
    property var groups: []
    // One or two controls, after the facts. An action that is not a fact —
    // the undo that reverses the last change — needs a real control rather
    // than a sentence telling the reader which key to press.
    property alias controls: controlSlot.data

    // Which fact was pressed, as the index of its group and its index within
    // that group. The caller knows what it put there; the bar does not need
    // to be told what any of it means.
    signal factActivated(int group, int fact)

    // How many of the groups are drawn on the bar. Written by the measurement
    // below and read by the delegates; a caller has no reason to set it.
    property int shownGroups: 0

    // How much of the bar the activity keeps before the groups start taking
    // room from it: about eight words at the default interface size. Below
    // this the activity elides rather than pushing a fact off the bar, since
    // a fact is something to press and the activity is a sentence that says
    // the same thing shortened.
    readonly property int activityFloor: Interface.px(96)

    // The facts that did not fit, flattened into what the menu draws. Each
    // one keeps the two indices it came from, so activating it from the menu
    // is indistinguishable from activating it on the bar.
    readonly property var hiddenFacts: {
        const out = []
        const all = root.groups
        for (let g = root.shownGroups; g < all.length; ++g) {
            const group = all[g]
            const facts = group.facts !== undefined ? group.facts : []
            const groupLabel = group.label !== undefined ? group.label : ""
            for (let f = 0; f < facts.length; ++f) {
                const text = facts[f].text !== undefined ? facts[f].text : ""
                out.push({
                    "group": g,
                    "fact": f,
                    // The group's name travels with the fact, because a menu
                    // item reading "2 behind" out of the bar has lost which
                    // of the two groups it was in.
                    "label": groupLabel === ""
                        ? text
                        : qsTr("%1 — %2", "a status-bar group and one of its facts")
                              .arg(groupLabel).arg(text),
                    "symbol": facts[f].symbol !== undefined ? facts[f].symbol : ""
                })
            }
        }
        return out
    }

    // Text alone fits in `Interface.statusBarHeight`. A control does not:
    // `Interface.controlHeight` is taller than the bar, and a bar that kept
    // its text height would draw the button clipped at both ends. So the bar
    // is as tall as what it holds, and returns to its resting height when
    // the slot is empty or everything in it is hidden — which is why this
    // reads the slots' implicit heights rather than counting their children.
    implicitHeight: Math.max(Interface.statusBarHeight,
                             controlSlot.implicitHeight + Interface.spaceTight,
                             groupRow.implicitHeight + Interface.spaceTight)
    color: Theme.footerBackground

    // Measured once per turn of the event loop rather than on every property
    // that feeds it, because a resize moves several of them at once and the
    // answer is only wanted after the last.
    function scheduleGroupLayout() {
        Qt.callLater(root.relayoutGroups)
    }

    // How many groups fit, left to right, in what is left after the margins,
    // the activity's floor, the plain facts and the controls.
    //
    // Left to right and never a gap: the groups that fit are a prefix of the
    // list, so what the menu holds is the tail. A rule that kept whichever
    // groups happened to be narrow would reorder the bar as the numbers on it
    // changed.
    function relayoutGroups() {
        const total = groupRepeater.count
        if (total === 0) {
            root.shownGroups = 0
            return
        }

        const taken = Interface.space * 2
            + (root.activity === "" ? 0 : root.activityFloor + Interface.space)
            + (factsRow.implicitWidth > 0
                   ? factsRow.implicitWidth + Interface.space : 0)
            + (controlSlot.implicitWidth > 0
                   ? controlSlot.implicitWidth + Interface.space : 0)
        const budget = root.width - taken

        const fitting = function (room) {
            let used = 0
            let count = 0
            for (let i = 0; i < total; ++i) {
                const wrapper = groupRepeater.itemAt(i)
                if (!wrapper)
                    break
                if (used + wrapper.naturalWidth > room)
                    break
                used += wrapper.naturalWidth
                count += 1
            }
            return count
        }

        // Twice, because the control that says what is hidden only exists
        // when something is, and it takes room of its own once it does.
        let count = fitting(budget)
        if (count < total)
            count = fitting(budget - overflow.naturalWidth)
        root.shownGroups = count
    }

    onWidthChanged: root.scheduleGroupLayout()
    onGroupsChanged: root.scheduleGroupLayout()
    onActivityChanged: root.scheduleGroupLayout()
    Component.onCompleted: root.scheduleGroupLayout()

    // Every width above is in interface pixels, so a reader who changes the
    // interface size changes all of them at once.
    Connections {
        target: Interface
        function onChanged() { root.scheduleGroupLayout() }
    }

    KvitDivider {
        anchors.top: parent.top
        anchors.left: parent.left
        anchors.right: parent.right
    }

    // The facts that did not fit, reachable rather than dropped. A Menu is
    // what makes it reachable: Qt gives it arrow-key navigation, Return to
    // activate and Escape to close once it is a real Menu.
    KvitMenu {
        id: overflowMenu
        objectName: "overflowMenu"
        parent: overflowLink
        y: -height

        Instantiator {
            model: root.hiddenFacts
            delegate: KvitMenuItem {
                required property var modelData
                text: modelData.label
                symbol: modelData.symbol
                onTriggered: root.factActivated(modelData.group, modelData.fact)
            }
            onObjectAdded: (index, object) => overflowMenu.insertItem(index, object)
            onObjectRemoved: (index, object) => overflowMenu.removeItem(object)
        }
    }

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: Interface.space
        anchors.rightMargin: Interface.space
        spacing: Interface.space

        KvitLabel {
            Layout.fillWidth: true
            Layout.minimumWidth: 0
            role: "caption"
            color: Theme.textMuted
            text: root.activity
            visible: text !== ""
        }
        Item { Layout.fillWidth: root.activity === "" }

        // The groups, and the control that says what is not among them.
        //
        // No spacing of its own: each wrapper carries its own leading gap, so
        // that a collapsed one takes no room at all. A positioner puts its
        // spacing between every visible child whatever its width, and a bar
        // with three collapsed groups would otherwise carry three gaps to
        // nowhere.
        Row {
            id: groupRow
            Layout.alignment: Qt.AlignVCenter
            visible: root.groups.length > 0
            spacing: 0

            Repeater {
                id: groupRepeater
                model: root.groups

                delegate: Item {
                    id: wrapper
                    objectName: "group"

                    required property int index
                    required property var modelData

                    readonly property int leading:
                        wrapper.index > 0 ? Interface.space : 0
                    readonly property real naturalWidth:
                        content.implicitWidth + wrapper.leading
                    readonly property bool shown:
                        wrapper.index < root.shownGroups

                    width: wrapper.shown ? wrapper.naturalWidth : 0
                    height: content.implicitHeight
                    clip: true
                    opacity: wrapper.shown ? 1 : 0
                    // Out of the tab order while it is not on the bar. A
                    // collapsed group still has its links inside it, and a
                    // reader tabbing along the bar would otherwise stop on
                    // three of them with nothing to see.
                    enabled: wrapper.shown

                    onNaturalWidthChanged: root.scheduleGroupLayout()
                    Component.onCompleted: root.scheduleGroupLayout()

                    Row {
                        id: content
                        x: wrapper.leading
                        spacing: Interface.spaceNear

                        KvitLabel {
                            anchors.verticalCenter: parent.verticalCenter
                            visible: text !== ""
                            text: wrapper.modelData.label !== undefined
                                  ? wrapper.modelData.label : ""
                            role: "caption"
                            color: Theme.textFaint
                            elide: Text.ElideNone
                        }
                        Repeater {
                            model: wrapper.modelData.facts !== undefined
                                   ? wrapper.modelData.facts : []

                            delegate: KvitLink {
                                objectName: "fact"

                                required property var modelData
                                required property int index

                                anchors.verticalCenter: parent.verticalCenter
                                role: "caption"
                                text: modelData.text !== undefined
                                      ? modelData.text : ""
                                symbol: modelData.symbol !== undefined
                                        ? modelData.symbol : ""
                                elide: Text.ElideNone
                                onActivated: root.factActivated(wrapper.index,
                                                                index)
                            }
                        }
                    }
                }
            }

            Item {
                id: overflow

                readonly property int leading:
                    root.shownGroups > 0 ? Interface.space : 0
                readonly property real naturalWidth:
                    overflowLink.implicitWidth + overflow.leading
                readonly property bool shown: root.hiddenFacts.length > 0

                width: overflow.shown ? overflow.naturalWidth : 0
                height: overflowLink.implicitHeight
                clip: true
                opacity: overflow.shown ? 1 : 0
                enabled: overflow.shown

                onNaturalWidthChanged: root.scheduleGroupLayout()

                KvitLink {
                    id: overflowLink
                    objectName: "overflow"
                    x: overflow.leading
                    role: "caption"
                    symbol: "more"
                    // The count is of facts rather than of groups, because a
                    // reader wants to know how many things are not on the bar
                    // and does not know how they were grouped.
                    text: {
                        const hidden = root.hiddenFacts.length
                        return qsTr("%1 more",
                                    "a count of status-bar facts that did not fit",
                                    hidden)
                            .arg(Number(hidden).toLocaleString(Qt.locale(), 'f', 0))
                    }
                    elide: Text.ElideNone
                    // Which groups they came from, for a reader who hears the
                    // control rather than sees the bar it sits on.
                    Accessible.description: {
                        const names = []
                        for (let g = root.shownGroups;
                             g < root.groups.length; ++g) {
                            const label = root.groups[g].label
                            if (label !== undefined && label !== "")
                                names.push(label)
                        }
                        return names.join(qsTr(", "))
                    }
                    onActivated: {
                        if (overflowMenu.opened)
                            overflowMenu.close()
                        else
                            overflowMenu.open()
                    }
                }
            }
        }

        // The plain facts. Wrapped rather than repeated straight into the
        // layout so that the measurement above has one width to read, and
        // hidden altogether when there are none so the bar keeps the gaps it
        // had before groups existed.
        Row {
            id: factsRow
            Layout.alignment: Qt.AlignVCenter
            visible: root.facts.length > 0
            spacing: Interface.space

            onImplicitWidthChanged: root.scheduleGroupLayout()

            Repeater {
                model: root.facts
                delegate: Row {
                    id: fact
                    required property string modelData
                    required property int index
                    spacing: Interface.space

                    KvitLabel {
                        anchors.verticalCenter: parent.verticalCenter
                        text: "·"
                        role: "caption"
                        color: Theme.textFaint
                        visible: fact.index > 0
                        elide: Text.ElideNone
                    }
                    KvitLabel {
                        anchors.verticalCenter: parent.verticalCenter
                        text: fact.modelData
                        role: "caption"
                        color: Theme.textMuted
                        tabular: true
                        elide: Text.ElideNone
                    }
                }
            }
        }

        Row {
            id: controlSlot
            Layout.alignment: Qt.AlignVCenter
            spacing: Interface.spaceTight

            onImplicitWidthChanged: root.scheduleGroupLayout()
        }
    }
}
