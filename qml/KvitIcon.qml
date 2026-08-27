// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// Small, recolourable symbols. Chrome asks for one by a name that says what it
// means rather than what it looks like, so the drawing can change without
// touching a call site.
//
// This is kvit-notes-pro's WorksIcon with three changes. The font ships with
// this module rather than with an application, so importing Kvit.Ui is enough
// and a consumer never adds a font alias to its own resource file. The
// name-to-codepoint table is generated from the font's published stylesheet
// rather than hand-written, so adding a symbol is a data change and the
// vocabulary skill's catalogue can list what exists. And an unrecognised name
// is visible: WorksIcon computed a `recognized` flag and then drew nothing,
// which means a typo ships as an empty box that nobody notices until somebody
// wonders where the button went.
//
// Each symbol is one glyph of the Phosphor icon font (MIT, shipped in
// resources/fonts beside its licence text). The hand-drawn vector paths this
// replaced were laid out on a 24-unit grid and then scaled by 18/24 with a
// fixed stroke width, which left strokes about 1.7 device pixels wide lying
// across two rows of pixels; a font is rasterized for the size and pixel grid
// it is actually drawn at, so the strokes come out even.
Item {
    id: root

    // What the symbol means. IconCatalog.meaningNames() is the list.
    required property string name
    property color color: Theme.textPrimary

    // Call sites give an icon a size by setting width and height, and the
    // glyph is drawn at whichever is smaller, so a 13-pixel combo indicator
    // stays 13 pixels and an unconstrained icon stays 18.
    readonly property int designSize: Interface.iconSize
    implicitWidth: designSize
    implicitHeight: designSize

    readonly property string glyph: IconCatalog.glyph(root.name)
    readonly property bool recognized: glyph.length > 0

    // One loader per icon costs almost nothing: Qt caches application fonts by
    // URL, so every loader after the first finds the family already registered
    // rather than parsing the file again.
    FontLoader {
        id: face
        // Relative to this file's own URL, so it resolves whether the
        // module was loaded from the binary's resources or from a
        // directory on disk.
        source: Qt.resolvedUrl("fonts/Phosphor.ttf")
    }

    Text {
        anchors.fill: parent
        visible: root.recognized
        text: root.glyph
        color: root.color
        font.family: face.name
        font.pixelSize: Math.max(1, Math.round(Math.min(
            root.width > 0 ? root.width : root.designSize,
            root.height > 0 ? root.height : root.designSize,
            root.designSize)))
        // The em box is square and the glyph fills it, so one centred
        // character is the icon centred in whatever rectangle a control hands
        // its content item.
        horizontalAlignment: Text.AlignHCenter
        verticalAlignment: Text.AlignVCenter
        // Rasterize for this exact pixel size; see the note on KvitLabel.
        renderType: Text.NativeRendering
    }

    // What an unrecognised name draws: a hatched box in the danger colour,
    // and a warning on the console that names the icon so the call site is
    // findable.
    //
    // This is the whole difference between a typo that is found in the first
    // render and one that ships. tests/test_components.cpp fails on the
    // warning, so a bad name cannot reach a release even if nobody looked at
    // the screen.
    Rectangle {
        anchors.fill: parent
        visible: !root.recognized
        color: "transparent"
        border.width: Interface.hairline
        border.color: Theme.danger
        Text {
            anchors.centerIn: parent
            text: "?"
            color: Theme.danger
            font.family: Interface.fontFamily
            font.pixelSize: Math.max(1, Math.round(parent.height * 0.7))
            renderType: Text.NativeRendering
        }
    }

    onRecognizedChanged: if (!recognized && name !== "")
        console.warn("KvitIcon: no symbol named '" + name + "'")
    Component.onCompleted: if (!recognized && name !== "")
        console.warn("KvitIcon: no symbol named '" + name + "'")
}
