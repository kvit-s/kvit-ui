// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// One state of one component, rendered beside the code that produced it.
//
// The rendering and the code sample are the same string. That is the point:
// a gallery whose samples are written separately from what it draws is a
// gallery whose samples go stale, and the first person to find out is somebody
// who copied one and got a compile error.
//
// The cost is that a broken snippet is a runtime error rather than a compile
// error, which would normally be the wrong trade. It is the right one here
// because tests/test_gallery.cpp instantiates every specimen in the catalogue
// and fails on any error or warning — so a snippet that does not compile stops
// the build, in the same place a QML file would have.
Column {
    id: root

    required property string caption
    required property string snippet
    // Extra imports the snippet needs beyond QtQuick and Kvit.Ui.
    property var extraImports: []
    property bool showSource: true
    // False for the one thing a page cannot draw: a component whose root is a
    // Window, which has no place inside another window's item tree. The
    // snippet is then shown and not built here — tests/test_gallery still
    // compiles and runs it, so it is checked exactly like every other one,
    // and what a reader copies is real QML rather than a sketch of it.
    property bool showRender: true

    // What went wrong, empty when the snippet built. Read by the gallery test.
    property string error: ""

    spacing: Interface.spaceNear
    width: parent ? parent.width : Interface.px(600)

    KvitLabel {
        text: root.caption
        role: "small"
        color: Theme.textMuted
    }

    // The specimen itself, on the list ground so a component with a
    // transparent background is still visible.
    Rectangle {
        width: parent.width
        visible: root.showRender
        height: Math.max(Interface.rowHeight,
                         stage.height + Interface.spaceLoose * 2)
        radius: Interface.radiusCard
        color: Theme.listBackground
        border.width: Interface.hairline
        border.color: Theme.border

        // The stage takes its height from what was built into it and its
        // width from this rectangle. Filling the rectangle instead would make
        // the rectangle's height depend on the stage and the stage's on the
        // rectangle, which hangs during construction rather than warning.
        Item {
            id: stage
            x: Interface.spaceLoose
            y: Interface.spaceLoose
            width: parent.width - Interface.spaceLoose * 2
            height: childrenRect.height
        }
    }

    KvitLabel {
        width: parent.width
        visible: root.error !== ""
        text: root.error
        role: "small"
        color: Theme.danger
        wrapMode: Text.WordWrap
        elide: Text.ElideNone
    }

    // The source, in the monospace family. Selectable, because what a reader
    // does with a sample is copy it.
    Rectangle {
        width: parent.width
        visible: root.showSource
        implicitHeight: source.implicitHeight + Interface.space * 2
        height: implicitHeight
        radius: Interface.radiusBar
        color: Theme.codePanelBackground

        TextEdit {
            id: source
            anchors.fill: parent
            anchors.margins: Interface.space
            text: root.snippet
            readOnly: true
            selectByMouse: true
            wrapMode: TextEdit.NoWrap
            color: Theme.textSecondary
            selectionColor: Theme.selectionActiveTint
            font.family: Interface.resolvedMonoFamily
            font.pixelSize: Interface.small
            renderType: Text.NativeRendering
        }
    }

    Component.onCompleted: build()

    function build() {
        if (!root.showRender)
            return
        const imports = ["import QtQuick", "import QtQuick.Controls",
                         "import QtQuick.Layouts", "import Kvit.Ui"]
            .concat(root.extraImports)
        try {
            // The third argument has to be a real URL. It is only used to
            // name the object in error messages and to resolve relative paths
            // out of the snippet, and a plain string like "specimen: the
            // seven roles" is not one — Qt does not reject it, it hangs
            // inside the creation, which is a long way from where anyone
            // would look.
            Qt.createQmlObject(
                imports.join("\n") + "\n" + root.snippet, stage,
                "qrc:/kvit-ui-gallery/specimen/"
                    + encodeURIComponent(root.caption) + ".qml")
            root.error = ""
        } catch (e) {
            root.error = String(e)
        }
    }
}
