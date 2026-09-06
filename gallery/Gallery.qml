// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component.
pragma ComponentBehavior: Bound

import QtQuick
import QtQuick.Layouts
import Kvit.Ui
import Kvit.Gallery

// The gallery window: a list of components down the left, one component's
// page in the middle, and the theme and interface size along the top.
//
// The controls are at the top rather than in a settings dialog because they
// are the point. What the gallery is for is looking at the same component in
// four themes, with alignment-sensitive controls also captured at the minimum
// and maximum interface sizes; a control that takes two clicks to reach makes
// that comparison something a reviewer stops doing.
KvitWindow {
    id: window

    visible: true
    title: qsTr("kvit-ui gallery")

    // Where the screenshot run writes, and which page to open. Set from
    // main.cpp as initial properties rather than installed as context
    // properties: a context property is a name a static analyser cannot see,
    // so every use of one is reported as unqualified access and the lint
    // category that catches a real typo becomes noise.
    property string shotDirectory: ""
    property string startingPage: ""

    // Which component's page is showing.
    property int current: 0
    property string filter: ""

    readonly property var currentEntry:
        Catalog.components[window.current]

    // How many components there are, and that number written the way the
    // reader writes numbers.
    readonly property int componentTally: Catalog.components.length
    readonly property string groupedTally:
        Number(window.componentTally).toLocaleString(Qt.locale(), 'f', 0)

    readonly property var shown: {
        const wanted = window.filter.trim().toLowerCase()
        if (wanted === "")
            return Catalog.components
        const found = []
        for (let i = 0; i < Catalog.components.length; ++i) {
            const entry = Catalog.components[i]
            if (entry.name.toLowerCase().indexOf(wanted) >= 0
                || entry.group.toLowerCase().indexOf(wanted) >= 0)
                found.push(entry)
        }
        return found
    }

    // Open one component's page by name.
    //
    // Not `show`. A Window already has a `show()` slot, taking no arguments
    // and doing something else entirely, and the call site does not read as
    // ambiguous: `window.show(name)` compiles, runs, re-shows the window that
    // is already visible, and discards the name with one line on the console
    // that says nothing about pages. Every click in the component list did
    // that instead of changing the page, and so did `--page`, while the
    // screenshot run kept working because it assigns `current` itself.
    function showPage(name: string) {
        for (let i = 0; i < Catalog.components.length; ++i) {
            if (Catalog.components[i].name === name) {
                window.current = i
                return
            }
        }
        console.warn("no gallery page for", name)
    }

    Component.onCompleted: {
        if (window.startingPage !== "")
            window.showPage(window.startingPage)
        if (window.shotDirectory !== "")
            shotStart.start()
    }

    header: KvitHeader {
        anchors.fill: parent
        wordmark: qsTr("kvit-ui")

        RowLayout {
            anchors.fill: parent
            spacing: Interface.columnGap

            Item { Layout.fillWidth: true }

            KvitSegmented {
                label: qsTr("Theme")
                current: Theme.themeId
                options: [{ value: "light", label: qsTr("Light") },
                          { value: "dark", label: qsTr("Dark") },
                          { value: "sepia", label: qsTr("Sepia") },
                          { value: "highContrast", label: qsTr("Contrast") }]
                onChosen: value => Theme.themeId = value
            }

            KvitStepper {
                label: qsTr("Interface size")
                from: Interface.minFontSize
                to: Interface.maxFontSize
                value: Interface.fontSize
                unit: qsTr("px")
                onValueModified: v => Interface.fontSize = v
            }
        }
    }

    sidebar: ColumnLayout {
        anchors.fill: parent
        spacing: 0

        KvitSearchField {
            Layout.fillWidth: true
            Layout.margins: Interface.space
            placeholderText: qsTr("Filter components")
            matches: window.shown.length
            matchedNoun: qsTr("component")
            onTextChanged: window.filter = text
        }

        KvitRegion {
            Layout.fillWidth: true
            Layout.fillHeight: true
            padding: 0

            Column {
                width: parent.width

                Repeater {
                    model: window.shown

                    delegate: Column {
                        id: listing
                        required property var modelData
                        required property int index
                        width: parent.width

                        readonly property bool startsGroup:
                            index === 0
                            || window.shown[index - 1].group !== modelData.group

                        KvitSectionHeading {
                            width: listing.width
                            visible: listing.startsGroup
                            height: visible ? implicitHeight : 0
                            text: listing.modelData.group
                        }

                        KvitRow {
                            width: listing.width
                            form: "slim"
                            rule: false
                            label: listing.modelData.name
                            selected: window.currentEntry !== undefined
                                      && window.currentEntry.name
                                         === listing.modelData.name
                            interactive: true
                            onActivated: window.showPage(listing.modelData.name)

                            KvitLabel {
                                anchors.fill: parent
                                anchors.leftMargin: Interface.spaceLoose
                                text: listing.modelData.name
                                role: "body"
                            }
                        }
                    }
                }
            }
        }
    }

    body: KvitRegion {
        id: pageRegion
        anchors.fill: parent

        GalleryPage {
            id: page
            width: parent.width
            entry: window.currentEntry
        }
    }

    statusBar: KvitStatusBar {
        anchors.fill: parent
        activity: shots.running
            ? qsTr("Writing screenshots: %1").arg(shots.progress) : ""
        // Two source strings picked by the count rather than one with `%n` in
        // it, and the number substituted through %1 so it is grouped the way
        // the reader groups digits. This is the same shape every counted
        // string in the library now takes, and the gallery is where somebody
        // looks to see how the library writes one: "%n component(s)" put a
        // form field on the shop window.
        facts: [window.componentTally === 1
                    ? qsTr("%1 component", "a count of one component",
                           window.componentTally).arg(window.groupedTally)
                    : qsTr("%1 components", "a count of several components",
                           window.componentTally).arg(window.groupedTally),
                Theme.displayName(Theme.resolvedTheme),
                qsTr("%1 px").arg(Interface.fontSize)]
    }

    // The screenshot run does not begin in Component.onCompleted.
    //
    // grabToImage schedules its work for after the next frame is rendered, and
    // at construction time there has been no frame: the callback is queued
    // behind a render that has not happened, so the run stalls before its
    // first image and the directory stays empty. One short timer puts the
    // start after the window's first frame.
    Timer {
        id: shotStart
        interval: 100
        onTriggered: shots.start()
    }

    // The screenshot run.
    //
    // A fixed, named set of images — one per component per theme — so that
    // what gets reviewed after a token change is the diff against the
    // previous run rather than sixty images looked at again from scratch.
    // prd.md §12 names screenshot review as the thing that does not scale, and
    // a stable file name per specimen is the whole mitigation.
    QtObject {
        id: shots

        property bool running: false
        property string progress: ""
        property int themeIndex: 0
        property int pageIndex: 0
        property int sizeIndex: 0
        property int baselineSize: Interface.fontSize
        readonly property var themes: ["light", "dark", "sepia", "highContrast"]

        function sizesForPage() {
            const requested = Catalog.components[pageIndex].shotSizes
            if (requested === undefined)
                return [baselineSize]

            const resolved = []
            for (let i = 0; i < requested.length; ++i) {
                let size = baselineSize
                if (requested[i] === "minimum")
                    size = Interface.minFontSize
                else if (requested[i] === "maximum")
                    size = Interface.maxFontSize
                if (resolved.indexOf(size) < 0)
                    resolved.push(size)
            }
            return resolved
        }

        function start() {
            // Still every animation for the duration of the run.
            //
            // A screenshot is taken one frame after the page changes, and
            // anything with a transition is then captured mid-flight: a
            // KvitPopover fading in over 100 ms is grabbed at opacity zero,
            // which is why its page was a blank card. Reduced motion sets
            // Theme.motionScale to 0, and every duration in the library is
            // written as `n * Theme.motionScale`, so the whole set becomes
            // deterministic rather than a race against the frame clock.
            //
            // It also happens to be a state worth having screenshots of.
            Theme.reducedMotion = true
            running = true
            baselineSize = Interface.fontSize
            themeIndex = 0
            pageIndex = 0
            sizeIndex = 0
            step()
        }

        function step() {
            if (themeIndex >= themes.length) {
                running = false
                Qt.exit(0)
                return
            }
            Theme.themeId = themes[themeIndex]
            const sizes = sizesForPage()
            Interface.fontSize = sizes[sizeIndex]
            window.current = pageIndex
            progress = themes[themeIndex] + " / "
                       + Interface.fontSize + " px / "
                       + Catalog.components[pageIndex].name
            // One turn of the event loop before grabbing, so the page has
            // built and its bindings have settled. Without it the first
            // component of each theme is captured mid-layout.
            Qt.callLater(shots.capture)
        }

        function capture() {
            const name = Catalog.components[pageIndex].name
            const sizePart = Interface.fontSize === baselineSize
                ? "" : "-" + Interface.fontSize + "px"
            const path = window.shotDirectory + "/" + themes[themeIndex]
                         + sizePart + "-" + name
                         + ".png"
            // grabToImage on the content item, not grabWindow on the Window:
            // the latter is a C++ method that QML cannot call, so a run
            // written that way fails on the first image.
            //
            // It is asynchronous — the grab happens after the next frame is
            // rendered — which is what the callback is for. Advancing without
            // waiting captures the previous page.
            window.contentItem.grabToImage(function (result) {
                result.saveToFile(path)
                shots.advance()
            })
        }

        function advance() {
            sizeIndex += 1
            if (sizeIndex >= sizesForPage().length) {
                sizeIndex = 0
                pageIndex += 1
                if (pageIndex >= Catalog.components.length) {
                    pageIndex = 0
                    themeIndex += 1
                }
            }
            Qt.callLater(shots.step)
        }
    }
}
