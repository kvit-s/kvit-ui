// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import QtQuick
import Kvit.Ui

// A run of chrome text.
//
// Not in prd.md §5.4's inventory, and here because every component in it needs
// one. A `Text` in this estate carries four properties that are the same every
// time — the chrome family, native rendering, a token colour, a role size —
// and writing them out per call site is how a stray `font.pixelSize: 11`
// appears in a tree that has a rule against numeric font sizes.
//
// `role` is what a call site sets, and it names what the text is rather than
// how big it is. The seven roles are InterfaceMetrics'.
//
// Native rendering rather than the default distance-field rasterizer: chrome
// text is drawn at 10 to 20 pixels and never scaled or rotated, and the
// distance field trades sharpness at exactly those sizes for a transform
// nothing here performs.
Text {
    id: root

    // "caption" | "small" | "body" | "strong" | "title" | "headline" | "display"
    property string role: "body"
    // Setting this to `true` draws the text in the monospace family, for an
    // identifier that has to line up down a column.
    property bool mono: false
    // Tabular numerals: every digit the same width, so a column of figures
    // lines up and a changing value does not shift the text beside it.
    property bool tabular: false

    color: Theme.textPrimary
    // The resolved families, never the stored preference: an empty
    // `font.family` is matched against nothing and falls back to an
    // arbitrary installed face rather than to the desktop default.
    font.family: root.mono ? Interface.resolvedMonoFamily
                           : Interface.resolvedFontFamily
    font.pixelSize: {
        switch (root.role) {
        case "caption":  return Interface.caption
        case "small":    return Interface.small
        case "strong":   return Interface.strong
        case "title":    return Interface.title
        case "headline": return Interface.headline
        case "display":  return Interface.display
        default:         return Interface.body
        }
    }
    font.features: root.tabular ? ({ "tnum": 1 }) : ({})
    renderType: Text.NativeRendering
    // Eliding by default, because a label that grows past its column pushes
    // whatever is beside it off the screen, and a name cut short with an
    // ellipsis is legible where an overlapped one is not.
    elide: Text.ElideRight
    verticalAlignment: Text.AlignVCenter
}
