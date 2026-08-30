// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
pragma Singleton

import QtQuick

// Every component in the vocabulary, what it is for, and the states worth
// looking at.
//
// This is one list with three consumers, which is why it is data rather than
// sixty-eight hand-written pages. The gallery renders it. tests/test_gallery
// instantiates every snippet in it and fails on any error, so a sample that
// does not compile stops the build. And `kvit-ui-gallery --catalog` writes the
// vocabulary skill's catalogue from it, so what an agent is told exists is
// what the repository actually contains.
//
// Each snippet is real QML, rendered as written. A gallery whose samples are
// written separately from what it draws is a gallery whose samples go stale,
// and the first person to find out is somebody who copied one.
QtObject {
    readonly property var components: [
        {
            "name": "KvitLabel",
            "group": "Foundation",
            "summary": "A run of chrome text at one of the seven type roles. Every other component uses it, which is what keeps the chrome family, the rendering mode and the eliding rule in one place.",
            "specimens": [
                {
                    "caption": "The seven roles",
                    "snippet": "Column {\n    spacing: Interface.spaceSnug\n    KvitLabel { text: \"display \u2014 a page title\"; role: \"display\" }\n    KvitLabel { text: \"headline \u2014 a pane title\"; role: \"headline\" }\n    KvitLabel { text: \"title \u2014 a section heading\"; role: \"title\" }\n    KvitLabel { text: \"strong \u2014 a name\"; role: \"strong\" }\n    KvitLabel { text: \"body \u2014 row text and prose\"; role: \"body\" }\n    KvitLabel { text: \"small \u2014 chip labels and sub-lines\"; role: \"small\" }\n    KvitLabel { text: \"caption \u2014 kind tags and counts\"; role: \"caption\" }\n}"
                },
                {
                    "caption": "Monospace and tabular numerals",
                    "snippet": "Column {\n    spacing: Interface.spaceSnug\n    KvitLabel { text: \"a1b2c3d4\"; mono: true }\n    KvitLabel { text: \"1,234.56\"; tabular: true }\n}"
                },
            ]
        },
        {
            "name": "KvitIcon",
            "group": "Foundation",
            "summary": "A symbol asked for by what it means rather than by what it looks like. The font ships with the module, so importing Kvit.Ui is enough.",
            "specimens": [
                {
                    "caption": "At the sizes the chrome uses",
                    "snippet": "Row {\n    spacing: Interface.space\n    KvitIcon { name: \"search\"; width: Interface.iconSizeSmall; height: width }\n    KvitIcon { name: \"chevron-right\" }\n    KvitIcon { name: \"trash\"; color: Theme.danger }\n    KvitIcon { name: \"success\"; color: Theme.success }\n}"
                },
                {
                    "caption": "An unrecognised name is visible, not blank",
                    "snippet": "KvitIcon { name: \"not-a-symbol\" }"
                },
            ]
        },
        {
            "name": "KvitIconButton",
            "group": "Foundation",
            "summary": "A button whose whole label is a symbol. A real AbstractButton, so it takes tab focus and a screen reader is told it is there; `label` fills both the tooltip and the accessible name so the two cannot disagree.",
            "specimens": [
                {
                    "caption": "Quiet, ordinary and checked",
                    "snippet": "Row {\n    spacing: Interface.space\n    KvitIconButton { symbol: \"pencil\"; label: \"Edit\" }\n    KvitIconButton { symbol: \"trash\"; label: \"Delete\"; form: \"ordinary\" }\n    KvitIconButton { symbol: \"pin\"; label: \"Pin\"; checked_: true }\n    KvitIconButton { symbol: \"copy\"; label: \"Copy\"; enabled: false }\n}"
                },
            ]
        },
        {
            "name": "KvitHeader",
            "group": "Structure",
            "summary": "The strip across the top of the window: the wordmark, navigation in the middle, actions on the right, in that fixed order on every screen.",
            "specimens": [
                {
                    "caption": "With navigation and actions",
                    "snippet": "KvitHeader {\n    width: parent.width\n    wordmark: \"kvit\"\n    Row {\n        spacing: Interface.space\n        anchors.verticalCenter: parent.verticalCenter\n        KvitTab { text: \"Today\"; selected: true }\n        KvitTab { text: \"Projects\"; count: 26 }\n        KvitTab { text: \"Decisions\"; count: 4 }\n    }\n    actions: [\n        Component { KvitIconButton { symbol: \"settings\"; label: \"Settings\" } }\n    ]\n}"
                },
            ]
        },
        {
            "name": "KvitSidebar",
            "group": "Structure",
            "summary": "The list of places down the left edge. Below the laptop breakpoint it collapses to its rail and every item becomes its symbol alone.",
            "specimens": [
                {
                    "caption": "Expanded and collapsed",
                    "snippet": "Row {\n    spacing: Interface.columnGap\n    KvitSidebar {\n        width: Interface.sidebarWidth\n        KvitSidebarItem { text: \"Today\"; symbol: \"calendar\"; selected: true }\n        KvitSidebarItem { text: \"Projects\"; symbol: \"folder\"; count: 26 }\n        KvitSidebarItem { text: \"Decisions\"; symbol: \"question\"; count: 4 }\n    }\n    KvitSidebar {\n        width: Interface.railWidth\n        collapsed: true\n        KvitSidebarItem { text: \"Today\"; symbol: \"calendar\"; selected: true }\n        KvitSidebarItem { text: \"Projects\"; symbol: \"folder\" }\n        KvitSidebarItem { text: \"Decisions\"; symbol: \"question\" }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitSidebarItem",
            "group": "Structure",
            "summary": "One place in the sidebar. The selected item takes a bar down its leading edge as well as a tint, so the selection is not resting on colour.",
            "specimens": [
                {
                    "caption": "Selected, hovered and counted",
                    "snippet": "Column {\n    width: Interface.sidebarWidth\n    KvitSidebarItem { text: \"Inbox\"; symbol: \"note\"; selected: true; count: 12 }\n    KvitSidebarItem { text: \"Archive\"; symbol: \"archive\" }\n}"
                },
            ]
        },
        {
            "name": "KvitBreadcrumb",
            "group": "Structure",
            "summary": "Where the reader is and the way back. The last crumb is the current place and is deliberately not a link; a long trail elides from the middle, keeping the section and the current place.",
            "specimens": [
                {
                    "caption": "A short trail and a long one",
                    "snippet": "Column {\n    spacing: Interface.space\n    KvitBreadcrumb {\n        trail: [{ label: \"Projects\", id: \"p\" },\n                { label: \"kvit-ui\", id: \"u\" }]\n    }\n    KvitBreadcrumb {\n        trail: [{ label: \"Projects\", id: \"p\" }, { label: \"kvit\", id: \"k\" },\n                { label: \"Wave 1\", id: \"w\" }, { label: \"Tokens\", id: \"t\" },\n                { label: \"Theme\", id: \"h\" }]\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitRegion",
            "group": "Structure",
            "summary": "A body that takes the height left over and scrolls what does not fit. The scroll bar sits beside the content rather than over it, so nothing is ever hidden behind it.",
            "specimens": [
                {
                    "caption": "Scrolling a column of rows",
                    "snippet": "KvitRegion {\n    width: 480; height: Interface.px(140)\n    Column {\n        id: rows\n        width: 460\n        Repeater {\n            model: 12\n            KvitSlimRow { width: rows.width; name: \"Row \" + (index + 1) }\n        }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitViewHead",
            "group": "Structure",
            "summary": "The strip at the top of a view: what the view is, how much is in it, and the controls that act on all of it.",
            "specimens": [
                {
                    "caption": "Title, count and controls",
                    "snippet": "KvitViewHead {\n    width: parent.width\n    title: \"Transactions\"\n    count: 1284\n    counted: \"transaction\"\n    subtitle: \"Everything since the account was opened\"\n    KvitSearchField { width: Interface.px(180) }\n    KvitButton { text: \"Export\"; form: \"ordinary\" }\n}"
                },
            ]
        },
        {
            "name": "KvitStatusBar",
            "group": "Structure",
            "summary": "The strip along the bottom: what is happening on the left, standing facts on the right.",
            "specimens": [
                {
                    "caption": "Working, with two facts",
                    "snippet": "KvitStatusBar {\n    width: parent.width\n    activity: \"Reindexing 4 of 26 projects\"\n    facts: [\"1,284 notes\", \"last synced 14:02\"]\n}"
                },
            ]
        },
        {
            "name": "KvitWindow",
            "group": "Structure",
            "summary": "The application shell: header, optional sidebar, body and status bar, with the sidebar collapsing to its rail below the laptop breakpoint. Not shown here because it is a Window; see the code sample.",
            "specimens": [
                {
                    "caption": "The whole shell",
                    "showRender": false,
                    "snippet": "// Source only: a Window has no place inside another window's item tree, so\n// this is the one page that shows a sample without drawing it. It is compiled\n// and run by tests/test_gallery like every other sample.\n//\n// Each of the four slots takes one item, and each fills the slot it is in.\nKvitWindow {\n    id: window\n    title: \"kvit-cash\"\n\n    header: KvitHeader {\n        anchors.fill: parent\n        wordmark: \"kvit\"\n        Row {\n            anchors.verticalCenter: parent.verticalCenter\n            spacing: Interface.space\n            KvitTab { text: \"Accounts\"; selected: true }\n            KvitTab { text: \"Budget\" }\n        }\n    }\n\n    sidebar: KvitSidebar {\n        anchors.fill: parent\n        // The window says when it is narrow; the sidebar draws itself as a\n        // rail when it is.\n        collapsed: window.sidebarCollapsed\n        KvitSidebarItem { text: \"Everyday\"; symbol: \"wallet\"; selected: true }\n        KvitSidebarItem { text: \"Savings\"; symbol: \"bank\" }\n        KvitSidebarItem { text: \"Budget\"; symbol: \"chart-line\" }\n    }\n\n    body: KvitRegion {\n        anchors.fill: parent\n        Column {\n            id: rows\n            width: parent.width\n            KvitSlimRow { width: rows.width; name: \"Checking\" }\n            KvitSlimRow { width: rows.width; name: \"Joint checking\" }\n            KvitSlimRow { width: rows.width; name: \"Emergency fund\" }\n        }\n    }\n\n    statusBar: KvitStatusBar {\n        anchors.fill: parent\n        activity: \"Matching 4 of 26 statements\"\n        facts: [\"1,284 transactions\", \"last synced 14:02\"]\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitSectionHeading",
            "group": "Content",
            "summary": "A group heading: a filled bar with a disclosure chevron, the name, what the group holds, the count with the word for what was counted, and the one action that applies to every row under it.",
            "specimens": [
                {
                    "caption": "Collapsible, counted, with an action",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.space\n    KvitSectionHeading {\n        width: parent.width\n        text: \"Waiting on me\"; counted: \"project\"; count: 4\n        action: \"Hand all to an agent\"; collapsible: true\n    }\n    KvitSectionHeading {\n        width: parent.width\n        text: \"Archived\"; kind: \"closed last quarter\"\n        collapsible: true; expanded: false; strong: true\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitRow",
            "group": "Content",
            "summary": "A list row at one of four heights, chosen by what the row carries rather than by how many rows a view wants to fit. Hover and keyboard focus are separate tints, because they are different rows.",
            "specimens": [
                {
                    "caption": "The four heights, and the three states",
                    "snippet": "Column {\n    width: parent.width\n    KvitRow { width: parent.width; form: \"full\"; label: \"Full\"\n        KvitLabel { anchors.centerIn: parent; text: \"full \u2014 56\" } }\n    KvitRow { width: parent.width; form: \"sub\"; label: \"Sub\"\n        KvitLabel { anchors.centerIn: parent; text: \"sub \u2014 48\" } }\n    KvitRow { width: parent.width; form: \"slim\"; selected: true; label: \"Slim\"\n        KvitLabel { anchors.centerIn: parent; text: \"slim \u2014 30, selected\" } }\n    KvitRow { width: parent.width; form: \"compact\"; current: true; label: \"Compact\"\n        KvitLabel { anchors.centerIn: parent; text: \"compact \u2014 24, keyboard focus\" } }\n}"
                },
            ]
        },
        {
            "name": "KvitSlimRow",
            "group": "Content",
            "summary": "The row a reader sees most: a name, what it is, one phrase about it and one figure, right aligned, in that order on every screen.",
            "specimens": [
                {
                    "caption": "Measured and unmeasured",
                    "snippet": "Column {\n    width: parent.width\n    KvitSlimRow { width: parent.width; symbol: \"folder\"; name: \"kvit-ui\"\n        kind: \"library\"; phrase: \"waiting on review\"; figure: \"3.5\"; unit: \"d\" }\n    KvitSlimRow { width: parent.width; symbol: \"folder\"; name: \"kvit-cash\"\n        kind: \"application\"; phrase: \"not started\"; measured: false }\n}"
                },
            ]
        },
        {
            "name": "KvitCard",
            "group": "Content",
            "summary": "A bounded block of content on a surface. An interactive card takes the control-boundary token and the hover tint, so a card that responds to a click looks like it will before it is clicked.",
            "specimens": [
                {
                    "caption": "Static and interactive",
                    "snippet": "Row {\n    spacing: Interface.columnGap\n    KvitCard { KvitLabel { text: \"A card\" } }\n    KvitCard { interactive: true; KvitLabel { text: \"Interactive\" } }\n    KvitCard { selected: true; KvitLabel { text: \"Selected\" } }\n}"
                },
            ]
        },
        {
            "name": "KvitPanel",
            "group": "Content",
            "summary": "A region of the window with its own ground: a sidebar, a toolbar strip. Structural, where a card is content \u2014 which is why it has no radius.",
            "specimens": [
                {
                    "caption": "With rules on two edges",
                    "snippet": "KvitPanel {\n    width: parent.width; height: Interface.px(60)\n    ruleTop: true; ruleBottom: true\n    KvitLabel { anchors.centerIn: parent; text: \"A panel\" }\n}"
                },
            ]
        },
        {
            "name": "KvitPane",
            "group": "Content",
            "summary": "The side pane: detail about the one thing selected, the same width wherever it appears so the list beside it never reflows.",
            "specimens": [
                {
                    "caption": "Open",
                    "snippet": "Item {\n    width: parent.width; height: Interface.px(160)\n    KvitPane {\n        height: parent.height\n        title: \"Transaction\"\n        KvitLabel { anchors.centerIn: parent; text: \"Detail goes here\" }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitDivider",
            "group": "Content",
            "summary": "A rule between two things. One design pixel, in the decorative border token rather than the control-boundary one.",
            "specimens": [
                {
                    "caption": "Horizontal and vertical",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.space\n    KvitDivider { width: parent.width }\n    Row {\n        height: Interface.px(24); spacing: Interface.space\n        KvitLabel { text: \"left\"; height: parent.height }\n        KvitDivider { vertical: true; height: parent.height }\n        KvitLabel { text: \"right\"; height: parent.height }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitDisclosure",
            "group": "Content",
            "summary": "A trigger with a body under it. Takes a `group` so several sections can behave as an accordion, which is what most of the eighteen hand-rolled versions in the estate actually are.",
            "specimens": [
                {
                    "caption": "Open and closed",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.space\n    KvitDisclosure {\n        width: parent.width; title: \"What changed\"; count: 3; expanded: true\n        KvitLabel { text: \"Three files were rewritten.\" }\n    }\n    KvitDisclosure { width: parent.width; title: \"What did not\"; count: 12 }\n}"
                },
            ]
        },
        {
            "name": "KvitEmptyState",
            "group": "Content",
            "summary": "What a view says when it has nothing to show: what would be here, why it is not, and the action that would fill it. Also the answer for a chart with no data, in place of an axis drawn around zeros.",
            "specimens": [
                {
                    "caption": "With an action",
                    "snippet": "KvitEmptyState {\n    width: parent.width\n    symbol: \"wallet\"\n    title: \"No transactions yet\"\n    detail: \"Import a statement or add one by hand, and it will appear here.\"\n    action: \"Import a statement\"\n}"
                },
                {
                    "caption": "Dashed, as a drop target",
                    "snippet": "KvitEmptyState {\n    width: parent.width\n    dashed: true\n    symbol: \"file-arrow-down\"\n    title: \"Drop a statement here\"\n    detail: \"CSV, OFX and QIF. The file is read on this machine and nothing is sent anywhere.\"\n    action: \"Choose a file\"\n}"
                },
            ]
        },
        {
            "name": "KvitChip",
            "group": "Marks",
            "summary": "A small labelled mark saying what kind of thing this is or what state it is in. The tone names a meaning rather than a colour, and every tone carries its outline as well as its tint so the distinction is not resting on hue.",
            "specimens": [
                {
                    "caption": "Every tone, tinted and filled",
                    "snippet": "Column {\n    spacing: Interface.space\n    Row {\n        spacing: Interface.spaceNear\n        KvitChip { text: \"neutral\" }\n        KvitChip { text: \"accent\"; tone: \"accent\" }\n        KvitChip { text: \"success\"; tone: \"success\" }\n        KvitChip { text: \"warning\"; tone: \"warning\" }\n        KvitChip { text: \"danger\"; tone: \"danger\" }\n        KvitChip { text: \"info\"; tone: \"info\" }\n    }\n    Row {\n        spacing: Interface.spaceNear\n        KvitChip { text: \"settled\"; tone: \"success\"; strong: true }\n        KvitChip { text: \"disputed\"; tone: \"danger\"; strong: true; symbol: \"warning\" }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitTag",
            "group": "Marks",
            "summary": "A label a person put there, in a colour they chose. Theme-independent, because the reader picked that red and would not expect it to change; the label colour is derived from the fill rather than taken from a token.",
            "specimens": [
                {
                    "caption": "Tinted, plain and removable",
                    "snippet": "Row {\n    spacing: Interface.spaceNear\n    KvitTag { text: \"groceries\"; tint: Theme.colorPalette[0] }\n    KvitTag { text: \"transport\"; tint: Theme.colorPalette[2] }\n    KvitTag { text: \"untinted\" }\n    KvitTag { text: \"removable\"; tint: Theme.colorPalette[4]; removable: true }\n}"
                },
            ]
        },
        {
            "name": "KvitBadge",
            "group": "Marks",
            "summary": "A count attached to something else. Caps rather than growing wide, and hides at zero \u2014 a badge showing nought says look here about nothing.",
            "specimens": [
                {
                    "caption": "Counts, capped, and hidden at zero",
                    "snippet": "Row {\n    spacing: Interface.space\n    KvitBadge { count: 3 }\n    KvitBadge { count: 42; tone: \"neutral\" }\n    KvitBadge { count: 1204; tone: \"danger\" }\n    KvitBadge { count: 0 }\n}"
                },
            ]
        },
        {
            "name": "KvitSlug",
            "group": "Marks",
            "summary": "An identifier: a reference, a hash, a key. Monospace, because the task is comparison rather than reading, and elided from the middle because the end is what distinguishes one from its neighbours.",
            "specimens": [
                {
                    "caption": "With and without a ground",
                    "snippet": "Column {\n    spacing: Interface.spaceNear\n    KvitSlug { text: \"TX-00173404\" }\n    KvitSlug { text: \"9f2c1ab4e77d0031\"; ground: false }\n}"
                },
            ]
        },
        {
            "name": "KvitDot",
            "group": "Marks",
            "summary": "A small filled circle standing for one thing's state. The shape is the second channel: a level that differs only by hue says nothing to a reader who cannot separate red from amber.",
            "specimens": [
                {
                    "caption": "Three shapes, three levels",
                    "snippet": "Row {\n    spacing: Interface.space\n    KvitDot { color: Theme.success; label: \"healthy\" }\n    KvitDot { color: Theme.warning; shape: \"square\"; label: \"slipping\" }\n    KvitDot { color: Theme.danger; shape: \"diamond\"; label: \"stalled\" }\n    KvitDot { color: Theme.textMuted; hollow: true; label: \"not measured\" }\n}"
                },
            ]
        },
        {
            "name": "KvitPip",
            "group": "Marks",
            "summary": "A row of dots standing for a small count \u2014 three of five days recorded. For counts a reader takes in without counting; past about seven, the figure is faster.",
            "specimens": [
                {
                    "caption": "Three of five, and none of four",
                    "snippet": "Column {\n    spacing: Interface.space\n    KvitPip { filled: 3; total: 5; label: \"3 of 5 days recorded\" }\n    KvitPip { filled: 0; total: 4; color: Theme.warning; label: \"no checks passed\" }\n}"
                },
            ]
        },
        {
            "name": "KvitFigure",
            "group": "Quantities",
            "summary": "A measured value: tabular numerals, the unit in muted colour at the smaller role, and an em dash where nothing was measured rather than a zero. A balance nobody computed is not a balance of zero.",
            "specimens": [
                {
                    "caption": "Measured, unmeasured and bounded",
                    "snippet": "Column {\n    spacing: Interface.spaceNear\n    KvitFigure { value: \"1,284.50\"; unit: \"GBP\" }\n    KvitFigure { value: \"0\"; unit: \"d\"; }\n    KvitFigure { measured: false }\n    KvitFigure { value: \"12\"; unit: \"d\"; bounded: true }\n    KvitFigure { value: \"94\"; unit: \"%\"; role: \"display\" }\n}"
                },
            ]
        },
        {
            "name": "KvitBeforeAfter",
            "group": "Quantities",
            "summary": "One value as it stands and the value something proposes to replace it with. Position, colour weight and the arrow all say which is which, so no reader depends on separating the two colours. A record being added has no before, and the em dash says so.",
            "specimens": [
                {
                    "caption": "Changed, added, removed and unchanged",
                    "snippet": "Column {\n    spacing: Interface.spaceNear\n    KvitBeforeAfter { label: \"Amount\"; before: \"42.00\"; after: \"44.50\"; unit: \"GBP\" }\n    KvitBeforeAfter { label: \"Category\"; before: \"\"; beforeMeasured: false; after: \"Groceries\" }\n    KvitBeforeAfter { label: \"Payee\"; before: \"TESCO 4471\"; after: \"\"; afterMeasured: false }\n    KvitBeforeAfter { label: \"Date\"; before: \"2026-08-14\"; after: \"2026-08-14\" }\n    KvitBeforeAfter { before: \"1,284.50\"; after: \"1,301.75\"; unit: \"GBP\"; role: \"strong\" }\n}"
                },
            ]
        },
        {
            "name": "KvitButton",
            "group": "Controls",
            "summary": "A button with words on it, in three forms. One primary per screen region: it is the action the screen is for. `danger` is separate from the form, because a destructive action can be any of the three.",
            "specimens": [
                {
                    "caption": "The three forms, and destructive",
                    "snippet": "Column {\n    spacing: Interface.space\n    Row {\n        spacing: Interface.space\n        KvitButton { text: \"Save\"; form: \"primary\" }\n        KvitButton { text: \"Cancel\"; form: \"ordinary\" }\n        KvitButton { text: \"More\"; form: \"quiet\" }\n        KvitButton { text: \"Disabled\"; form: \"ordinary\"; enabled: false }\n    }\n    Row {\n        spacing: Interface.space\n        KvitButton { text: \"Delete\"; form: \"primary\"; danger: true }\n        KvitButton { text: \"Delete\"; form: \"ordinary\"; danger: true; symbol: \"trash\" }\n        KvitButton { text: \"Save\"; form: \"primary\"; busy: true }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitStepper",
            "group": "Controls",
            "summary": "A number with a minus and a plus beside it, for a small range a reader adjusts by one or two. The range comes from whatever owns it rather than being repeated at the call site.",
            "specimens": [
                {
                    "caption": "The interface-size row, pointed at the right setting",
                    "snippet": "KvitStepper {\n    label: \"Interface size\"\n    from: Interface.minFontSize\n    to: Interface.maxFontSize\n    value: Interface.fontSize\n    unit: \"px\"\n    onValueModified: v => Interface.fontSize = v\n}"
                },
            ]
        },
        {
            "name": "KvitField",
            "group": "Controls",
            "summary": "A single line of text. The outline is the control-boundary token, so its edges are visible. An error is a message and a border together, never a border alone.",
            "specimens": [
                {
                    "caption": "Resting, filled and in error",
                    "snippet": "Column {\n    spacing: Interface.spaceLoose\n    KvitField { label: \"Payee\"; placeholderText: \"Who was paid\" }\n    KvitField { label: \"Reference\"; text: \"TX-00173404\" }\n    KvitField { label: \"Amount\"; text: \"twelve\"; error: \"Not a number\" }\n    KvitField { label: \"Locked\"; text: \"read only\"; enabled: false }\n}"
                },
            ]
        },
        {
            "name": "KvitSearchField",
            "group": "Controls",
            "summary": "A field that filters something. Escape clears rather than reverting, and it announces its result count \u2014 filtering is the one interaction whose whole outcome happens somewhere else on the screen.",
            "specimens": [
                {
                    "caption": "Empty and filtering",
                    "snippet": "Column {\n    spacing: Interface.spaceLoose\n    KvitSearchField { width: Interface.px(220) }\n    KvitSearchField {\n        width: Interface.px(220); text: \"harlow\"\n        matches: 47; matchedNoun: \"transaction\"\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitCheck",
            "group": "Controls",
            "summary": "A checkbox in three states. The third \u2014 partial \u2014 is what a parent row shows when some of its children are checked; drawing that as unchecked loses the information and drawing it as checked is a lie.",
            "specimens": [
                {
                    "caption": "Off, on, partial and disabled",
                    "snippet": "Column {\n    spacing: Interface.spaceNear\n    KvitCheck { text: \"Include archived\" }\n    KvitCheck { text: \"Include drafts\"; checked: true }\n    KvitCheck { text: \"Some of these\"; partial: true }\n    KvitCheck { text: \"Not available\"; enabled: false }\n}"
                },
            ]
        },
        {
            "name": "KvitSelect",
            "group": "Controls",
            "summary": "A choice from a list too long to lay out. For two to four self-evident options, KvitSegmented is right instead \u2014 a three-item dropdown hides two of the three answers for no reason.",
            "specimens": [
                {
                    "caption": "A currency picker",
                    "snippet": "KvitSelect {\n    label: \"Currency\"\n    model: [\"GBP\", \"EUR\", \"USD\", \"JPY\", \"CHF\", \"SEK\"]\n}"
                },
            ]
        },
        {
            "name": "KvitTab",
            "group": "Controls",
            "summary": "One tab in a row of them. The selected tab is marked by an underline as well as by colour and weight \u2014 selection shown by colour alone is the most common place the rule gets broken.",
            "specimens": [
                {
                    "caption": "Selected, counted and plain",
                    "snippet": "Row {\n    KvitTab { text: \"All\"; selected: true; count: 1284 }\n    KvitTab { text: \"Uncategorised\"; count: 47 }\n    KvitTab { text: \"Disputed\" }\n}"
                },
            ]
        },
        {
            "name": "KvitTooltip",
            "group": "Feedback",
            "summary": "A short label next to a control after a pause. Never the only place a control's meaning lives \u2014 a control explained only by its tooltip is unusable on a keyboard, a screen reader and a touch screen.",
            "specimens": [
                {
                    "caption": "On a button",
                    "snippet": "Item {\n    id: stage\n    width: parent.width; height: Interface.px(80)\n    KvitButton {\n        id: action\n        x: Interface.px(60); y: Interface.px(40)\n        text: \"Reconcile\"\n        // KvitTooltip rather than the attached `ToolTip.text`: the attached\n        // one is the platform style's, which is a yellow box belonging to no\n        // theme here.\n        KvitTooltip {\n            id: hint\n            text: \"Match these against the statement\"\n            delay: 0\n            Component.onCompleted: if (stage.Window.window) hint.visible = true\n        }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitPopover",
            "group": "Feedback",
            "summary": "A small surface anchored to a control holding something to act on. Takes focus, closes on Escape and on a click outside \u2014 which is what separates it from a hover card.",
            "specimens": [
                {
                    "caption": "Open, holding a form",
                    "snippet": "Item {\n    id: stage\n    width: parent.width; height: Interface.px(150)\n    KvitPopover {\n        id: filter\n        parent: stage\n        title: \"Filter\"\n        width: Interface.px(240)\n        // A Popup opens into a window, and `visible: true` on one with no\n        // parent does nothing at all.\n        Component.onCompleted: if (stage.Window.window) filter.open()\n        Column {\n            spacing: Interface.spaceNear\n            KvitCheck { text: \"Settled\"; checked: true }\n            KvitCheck { text: \"Pending\" }\n            KvitCheck { text: \"Disputed\" }\n        }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitHoverCard",
            "group": "Feedback",
            "summary": "More about the thing under the pointer. Read-only, always: a surface that appears on hover is unreachable by keyboard and by touch, so nothing inside one may be the only route to an action.",
            "specimens": [
                {
                    "caption": "Showing what a row could not fit",
                    "snippet": "KvitHoverCard {\n    shown: true\n    Column {\n        spacing: Interface.spaceTight\n        KvitLabel { text: \"Payment to Harlow depot\"; role: \"strong\" }\n        KvitLabel { text: \"17 March 2026 at 09:14\"; role: \"small\"\n                    color: Theme.textMuted }\n        KvitFigure { value: \"1,284.50\"; unit: \"GBP\" }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitToast",
            "group": "Feedback",
            "summary": "A message that appears, says one thing and goes away. For confirming what already happened. A toast with an action stays until dismissed \u2014 an undo that times out mid-read is an undo the reader cannot use.",
            "specimens": [
                {
                    "caption": "Four tones, and one with an undo",
                    "snippet": "Column {\n    spacing: Interface.space\n    KvitToast { shown: true; text: \"Copied to the clipboard\" }\n    KvitToast { shown: true; tone: \"success\"; text: \"Statement imported\" }\n    KvitToast { shown: true; tone: \"warning\"; text: \"Two rows could not be matched\" }\n    KvitToast { shown: true; tone: \"danger\"; text: \"The file could not be read\" }\n    KvitToast { shown: true; tone: \"success\"; text: \"40 transactions archived\"\n                action: \"Undo\" }\n}"
                },
            ]
        },
        {
            "name": "KvitNotice",
            "group": "Feedback",
            "summary": "A message that stays: a condition that is still true, where a toast is an acknowledgement of something that finished. Dismissible only when dismissing it is meaningful.",
            "specimens": [
                {
                    "caption": "A warning that can be dismissed, and an error that cannot",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.space\n    KvitNotice {\n        width: parent.width; tone: \"warning\"\n        text: \"Your licence expires in three days\"\n        detail: \"Renew before 30 March to keep syncing.\"\n        action: \"Renew\"; dismissible: true\n    }\n    KvitNotice {\n        width: parent.width; tone: \"danger\"\n        text: \"This vault could not be saved\"\n        detail: \"The disk is full. Nothing has been lost; the changes are still held.\"\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitDialog",
            "group": "Feedback",
            "summary": "A modal surface that has to be answered first. Expensive, and worth it only where the answer really does come first; the confirming button is on the right, and a destructive dialog takes no keyboard default.",
            "specimens": [
                {
                    "caption": "A destructive confirmation",
                    "snippet": "// Shown here without its modality so it sits on the page. A real one is\n// modal and centres itself on the window, which is what a dialog that has\n// to be answered first should do and is not something a gallery card can\n// contain.\nItem {\n    id: stage\n    width: parent.width; height: Interface.px(210)\n    KvitDialog {\n        id: confirm\n        parent: stage\n        modal: false\n        anchors.centerIn: stage\n        title: \"Delete 40 transactions?\"\n        detail: \"They will be removed from every report. This cannot be undone.\"\n        confirmText: \"Delete them\"\n        destructive: true\n        Component.onCompleted: if (stage.Window.window) confirm.open()\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitBar",
            "group": "Data",
            "summary": "One quantity against a stated scale. A bar with no value draws a tick rather than a zero-width fill; a bounded figure is hatched, and the hatch survives grayscale; and `maximum` is required, because a bar drawn against its own list rescales invisibly.",
            "specimens": [
                {
                    "caption": "Measured, bounded and unmeasured",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.spaceLoose\n    KvitBar { width: parent.width; value: 62; maximum: 100\n              label: \"Attention\"; unit: \"h\" }\n    KvitBar { width: parent.width; value: 38; maximum: 100; wide: true\n              bounded: true; color: Theme.axisAgent; label: \"Agent\"; unit: \"h\" }\n    KvitBar { width: parent.width; maximum: 100; measured: false\n              label: \"Unrecorded\" }\n}"
                },
            ]
        },
        {
            "name": "KvitStackedBar",
            "group": "Data",
            "summary": "Several quantities adding to one total. The two-pixel gap is the reason this is a component: adjacent fills that touch read as one fill with a colour change in it.",
            "specimens": [
                {
                    "caption": "A four-way breakdown",
                    "snippet": "KvitStackedBar {\n    width: parent.width\n    wide: true\n    label: \"Where the month went\"\n    segments: [\n        { value: 420, label: \"rent\", color: Theme.categorical(0) },\n        { value: 180, label: \"groceries\", color: Theme.categorical(1) },\n        { value: 95, label: \"transport\", color: Theme.categorical(2) },\n        { value: 60, label: \"everything else\", color: Theme.categorical(3) }\n    ]\n}"
                },
            ]
        },
        {
            "name": "KvitSpark",
            "group": "Data",
            "summary": "The shape of a series, small enough to sit in a row. A period that was never measured is a baseline tick rather than a zero-height bar: a missing week and a quiet week are different facts.",
            "specimens": [
                {
                    "caption": "With a hole in the middle",
                    "snippet": "KvitSpark {\n    width: Interface.px(160)\n    label: \"Notes written\"\n    values: [3, 5, 8, 6, null, null, 9, 12, 7, 4, 6, 11]\n}"
                },
            ]
        },
        {
            "name": "KvitTrend",
            "group": "Data",
            "summary": "A series with a value axis and a hover crosshair \u2014 read for values, where a spark is read for shape. A gap in the data draws as a gap: interpolating over a hole asserts values nobody measured. A second series is dashed as well as differently coloured, and both are named in the key and in the crosshair.",
            "specimens": [
                {
                    "caption": "A series, and the empty state",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.spaceLoose\n    KvitTrend {\n        width: parent.width\n        label: \"Balance\"; unit: \"GBP\"\n        minimumY: 0; maximumY: 2000\n        points: [{ y: 400 }, { y: 620 }, { y: 580 }, { y: 900 },\n                 { y: null }, { y: 1400 }, { y: 1250 }, { y: 1700 }]\n    }\n    KvitTrend { width: parent.width; height: Interface.px(90); label: \"Savings\" }\n}"
                },
                {
                    "caption": "Two series, with an annotation drawn over the plot",
                    "snippet": "KvitTrend {\n    id: worth\n    width: parent.width\n    height: Interface.px(150)\n    label: \"Assets\"; secondLabel: \"Liabilities\"; unit: \"GBP\"\n    minimumY: 0; maximumY: 2000\n    points: [{ y: 400 }, { y: 620 }, { y: 580 }, { y: 900 },\n             { y: 1150 }, { y: 1400 }, { y: 1250 }, { y: 1700 }]\n    secondPoints: [{ y: null }, { y: null }, { y: 300 }, { y: 340 },\n                   { y: 320 }, { y: 290 }, { y: 260 }, { y: 240 }]\n\n    // Where the second account's history begins. The plot rectangle is\n    // published, so an overlay lines up with the same arithmetic.\n    Rectangle {\n        x: worth.axisWidth + worth.plotWidth * 2 / 7\n        y: worth.plotTop\n        width: Interface.hairline\n        height: worth.plotHeight\n        color: Theme.marker\n    }\n    KvitLabel {\n        x: worth.axisWidth + worth.plotWidth * 2 / 7 + Interface.spaceSnug\n        y: worth.plotTop + worth.plotHeight - height\n        text: \"history starts here\"\n        role: \"caption\"\n        color: Theme.textFaint\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitDistribution",
            "group": "Data",
            "summary": "How a set of values is spread. An average of four days and an average made of one twenty-day outlier are the same number and different situations.",
            "specimens": [
                {
                    "caption": "Two rows on one scale",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.spaceLoose\n    KvitDistribution {\n        width: parent.width; label: \"Time to settle, card\"\n        scaleMinimum: 0; scaleMaximum: 30; unit: \"d\"\n        minimum: 1; lowerQuartile: 2; median: 3; upperQuartile: 5; maximum: 21\n    }\n    KvitDistribution {\n        width: parent.width; label: \"Time to settle, transfer\"\n        scaleMinimum: 0; scaleMaximum: 30; unit: \"d\"\n        minimum: 1; lowerQuartile: 1; median: 2; upperQuartile: 2; maximum: 4\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitGauge",
            "group": "Data",
            "summary": "How much of an allowance is used, with the target marked and a pace mark saying where an even rate would have reached. Sixty percent spent is fine on day eighteen and a problem on day six.",
            "specimens": [
                {
                    "caption": "Ahead of pace, behind it, and over",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.spaceLoose\n    KvitGauge { width: parent.width; value: 240; allowance: 600; pace: 0.6\n                label: \"Groceries\"; unit: \"GBP\" }\n    KvitGauge { width: parent.width; value: 480; allowance: 600; pace: 0.6\n                label: \"Leisure\"; unit: \"GBP\" }\n    KvitGauge { width: parent.width; value: 720; allowance: 600; pace: 0.6\n                label: \"Transport\"; unit: \"GBP\" }\n}"
                },
            ]
        },
        {
            "name": "KvitDelta",
            "group": "Data",
            "summary": "How much something changed and in which direction, with an arrow as well as a colour. Whether up is good is the caller's to say: a rise in spending and a rise in savings are the same arrow and opposite colours.",
            "specimens": [
                {
                    "caption": "Up, down and unchanged",
                    "snippet": "Row {\n    spacing: Interface.spaceLoose\n    KvitDelta { change: 12.4; unit: \"%\"; precision: 1; goodDirection: \"up\" }\n    KvitDelta { change: -3.2; unit: \"%\"; precision: 1; goodDirection: \"up\" }\n    KvitDelta { change: 18; unit: \"GBP\"; goodDirection: \"down\" }\n    KvitDelta { change: 0 }\n    KvitDelta { measured: false }\n}"
                },
            ]
        },
        {
            "name": "KvitStatTile",
            "group": "Data",
            "summary": "A card carrying one figure, what it is, how it changed and its recent shape \u2014 the shape of six of kvit-cash's seven dashboard widgets. The order is fixed so a row of tiles can be scanned one part at a time.",
            "specimens": [
                {
                    "caption": "A dashboard row",
                    "snippet": "Row {\n    spacing: Interface.columnGap\n    KvitStatTile {\n        label: \"Balance\"; value: \"4,182.30\"; unit: \"GBP\"\n        hasChange: true; change: 240.10; goodDirection: \"up\"\n        history: [3200, 3400, 3390, 3800, 4000, 4182]\n    }\n    KvitStatTile {\n        label: \"Spent this month\"; value: \"812.40\"; unit: \"GBP\"\n        hasChange: true; change: 96.20; goodDirection: \"down\"\n        caption: \"12 days remaining\"\n    }\n    KvitStatTile { label: \"Uncategorised\"; measured: false }\n}"
                },
            ]
        },
        {
            "name": "KvitFigureBlock",
            "group": "Data",
            "summary": "A figure with its name under it: one number a reader is meant to take away. The number is above and larger, because a row of these is read across the numbers.",
            "specimens": [
                {
                    "caption": "A row of three",
                    "snippet": "Row {\n    spacing: Interface.px(40)\n    KvitFigureBlock { label: \"Transactions\"; value: \"1,284\" }\n    KvitFigureBlock { label: \"Categorised\"; value: \"97\"; unit: \"%\" }\n    KvitFigureBlock { label: \"Reconciled\"; measured: false }\n}"
                },
            ]
        },
        {
            "name": "KvitCell",
            "group": "Data",
            "summary": "One cell of a table, drawn according to what kind of value its column holds. The kind comes from the column so every cell in it aligns the same way and says the same thing about a missing value.",
            "specimens": [
                {
                    "caption": "The six kinds",
                    "snippet": "Column {\n    width: parent.width\n    KvitCell { width: parent.width; kind: \"Text\"; value: \"Payment to Harlow depot\" }\n    KvitCell { width: parent.width; kind: \"Figure\"; value: 1284.5; unit: \"GBP\" }\n    KvitCell { width: parent.width; kind: \"Money\"; value: -42.9; unit: \"GBP\" }\n    KvitCell { width: parent.width; kind: \"Chip\"; value: \"Settled\"; mark: \"success\" }\n    KvitCell { width: parent.width; kind: \"Slug\"; value: \"TX-00173404\" }\n    KvitCell { width: parent.width; kind: \"Figure\"; measured: false }\n}"
                },
            ]
        },
        {
            "name": "KvitTable",
            "group": "Data",
            "summary": "A dense, configurable, editable table over a C++ model. Holds smooth scrolling and sub-100 ms filtering at 250,000 rows with twelve columns, which is the measurement prd.md Decision 3 turns on.",
            "specimens": [
                {
                    "caption": "Two hundred and fifty thousand rows",
                    "snippet": "Item {\n    width: parent.width; height: Interface.px(260)\n    BenchmarkTableModel { id: rows }\n    KvitTable { anchors.fill: parent; model: rows }\n}"
                },
                {
                    "caption": "Nothing to show",
                    "snippet": "Item {\n    width: parent.width; height: Interface.px(200)\n    BenchmarkTableModel { id: rows; totalRows: 0 }\n    KvitTable {\n        anchors.fill: parent\n        model: rows\n        emptyTitle: qsTr(\"No transactions\")\n        emptyDetail: qsTr(\"Nothing matches the current filter.\")\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitScrollBar",
            "group": "Flow",
            "summary": "A scroll bar, occupying its own strip rather than floating over content. Thirty-two files in the estate have a private version; an overlay bar hides the right-hand column of a table and the last character of every elided label.",
            "specimens": [
                {
                    "caption": "Beside a scrolling column",
                    "snippet": "KvitRegion {\n    width: 480; height: Interface.px(120)\n    Column {\n        id: rows\n        width: 460\n        Repeater { model: 10; KvitSlimRow { width: rows.width\n                   name: \"Row \" + (index + 1) } }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitMenu",
            "group": "Flow",
            "summary": "A list of commands. Eighteen private versions in the estate, and what they mostly get wrong is the same two things: no keyboard route in, and no separator before the destructive item.",
            "specimens": [
                {
                    "caption": "Opened, with shortcuts and a destructive item",
                    "snippet": "Item {\n    id: stage\n    width: parent.width; height: Interface.px(170)\n    KvitMenu {\n        id: menu\n        // Only once the item is in a window: popping a menu on an item\n        // with no window is a crash rather than a no-op.\n        Component.onCompleted: if (stage.Window.window) menu.popup(stage, 0, 0)\n        KvitMenuItem { text: \"Open\"; symbol: \"file\"; shortcut: \"Ctrl+O\" }\n        KvitMenuItem { text: \"Duplicate\"; symbol: \"copy\"; shortcut: \"Ctrl+D\" }\n        KvitMenuItem { text: \"Archive\"; symbol: \"archive\" }\n        MenuSeparator {}\n        KvitMenuItem { text: \"Delete\"; symbol: \"trash\"; destructive: true }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitMenuItem",
            "group": "Flow",
            "summary": "One line of a menu, carrying its shortcut on the right \u2014 which is how a menu teaches a faster route to a reader who keeps using it.",
            "specimens": [
                {
                    "caption": "Ordinary, disabled and destructive",
                    "snippet": "Item {\n    id: stage\n    width: parent.width; height: Interface.px(110)\n    KvitMenu {\n        id: menu\n        // Only once the item is in a window: popping a menu on an item\n        // with no window is a crash rather than a no-op.\n        Component.onCompleted: if (stage.Window.window) menu.popup(stage, 0, 0)\n        KvitMenuItem { text: \"Reconcile\"; symbol: \"check\"; shortcut: \"Ctrl+R\" }\n        KvitMenuItem { text: \"Split\"; symbol: \"split\"; enabled: false }\n        KvitMenuItem { text: \"Delete\"; symbol: \"trash\"; destructive: true }\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitTree",
            "group": "Flow",
            "summary": "A nested list the reader can open and close. Twelve private versions; what a shared one has to get right is the keyboard, because depth without Left and Right is a wall.",
            "specimens": [
                {
                    "caption": "A small hierarchy",
                    "snippet": "Item {\n    width: parent.width; height: Interface.px(250)\n    KvitTree {\n        anchors.fill: parent\n        label: \"Accounts\"\n        nodes: [\n            { label: \"Everyday\", children: [\n                \"Checking\", \"Joint checking\", \"Cash\"] },\n            { label: \"Savings\", children: [\n                \"Emergency fund\",\n                { label: \"Certificates\", children: [\"18 months\", \"3 years\"] }] },\n            \"Credit card\"\n        ]\n        Component.onCompleted: expandRecursively(-1, 1)\n    }\n}"
                },
            ]
        },
        {
            "name": "KvitSwitch",
            "group": "Flow",
            "summary": "An option that takes effect the moment it moves \u2014 where a checkbox is a value in a form that takes effect on submit. The knob moves and the track fills, so the state is not resting on hue.",
            "specimens": [
                {
                    "caption": "On, off and disabled",
                    "snippet": "Column {\n    spacing: Interface.spaceNear\n    KvitSwitch { text: \"Follow the system theme\"; checked: true }\n    KvitSwitch { text: \"Reduce motion\" }\n    KvitSwitch { text: \"Sync over cellular\"; enabled: false }\n}"
                },
            ]
        },
        {
            "name": "KvitRadioGroup",
            "group": "Flow",
            "summary": "One choice from a handful where the choice needs explaining. Arrow keys move within the group and Tab leaves it, which a column of separate controls does not do.",
            "specimens": [
                {
                    "caption": "Three options with detail",
                    "snippet": "KvitRadioGroup {\n    width: parent.width\n    label: \"Appearance\"\n    current: \"system\"\n    options: [\n        { value: \"system\", label: \"Follow the system\",\n          detail: \"Light or dark, whichever the desktop is set to.\" },\n        { value: \"light\", label: \"Always light\" },\n        { value: \"dark\", label: \"Always dark\" }\n    ]\n}"
                },
            ]
        },
        {
            "name": "KvitProgress",
            "group": "Flow",
            "summary": "How far through something the application is. An unknown total is a moving band rather than a bar creeping toward the end, and it carries a label saying how far through *what*.",
            "specimens": [
                {
                    "caption": "Determinate and indeterminate",
                    "snippet": "Column {\n    width: parent.width\n    spacing: Interface.spaceLoose\n    KvitProgress { width: parent.width; label: \"Importing statement\"\n                   value: 0.47 }\n    KvitProgress { width: parent.width; label: \"Reconciling\"\n                   determinate: false }\n}"
                },
            ]
        },
        {
            "name": "KvitSlider",
            "group": "Flow",
            "summary": "A value chosen by dragging, for a feel rather than a number. Where the exact value matters, KvitStepper or KvitNumberField are right \u2014 hitting a particular number on a slider is hard and it cannot be typed.",
            "specimens": [
                {
                    "caption": "With its value shown",
                    "snippet": "KvitSlider { width: Interface.px(200); label: \"Opacity\"; value: 0.6\n             unit: \"%\"; from: 0; to: 1 }"
                },
            ]
        },
        {
            "name": "KvitSplitView",
            "group": "Flow",
            "summary": "Two regions the reader can resize. The handle is a wide invisible strip with a hairline down the middle, so the target is comfortable and the rule is still thin; it also moves with the arrow keys.",
            "specimens": [
                {
                    "caption": "Two panes",
                    "snippet": "KvitSplitView {\n    width: parent.width; height: Interface.px(120)\n    KvitPanel { SplitView.preferredWidth: Interface.px(160)\n                KvitLabel { anchors.centerIn: parent; text: \"left\" } }\n    KvitPanel { SplitView.fillWidth: true\n                KvitLabel { anchors.centerIn: parent; text: \"right\" } }\n}"
                },
            ]
        },
        {
            "name": "KvitSegmented",
            "group": "Flow",
            "summary": "One choice from two to five short options, all visible \u2014 the shape kvit-cash's dashboard period control needs, where a control governing seven widgets should say what they are showing without being opened.",
            "specimens": [
                {
                    "caption": "A period control",
                    "snippet": "Column {\n    spacing: Interface.space\n    KvitSegmented {\n        label: \"Period\"; current: \"month\"\n        options: [{ value: \"week\", label: \"Week\" },\n                  { value: \"month\", label: \"Month\" },\n                  { value: \"quarter\", label: \"Quarter\" },\n                  { value: \"year\", label: \"Year\" }]\n    }\n    KvitSegmented { options: [\"All\", \"Mine\"]; current: \"All\" }\n}"
                },
            ]
        },
        {
            "name": "KvitTypeAhead",
            "group": "Flow",
            "summary": "A field offering matches as the reader types, for a list too long to read. `allowNew` has no default: a tag picker should let the reader invent one, a category picker should not.",
            "specimens": [
                {
                    "caption": "A category picker that will not invent categories",
                    "snippet": "KvitTypeAhead {\n    width: Interface.px(240)\n    label: \"Category\"\n    placeholder: \"Start typing\"\n    allowNew: false\n    source: [\"Groceries\", \"Transport\", \"Utilities\", \"Rent\",\n             \"Leisure\", \"Health\", \"Savings\", \"Income\"]\n}"
                },
            ]
        },
        {
            "name": "KvitConfirmInPlace",
            "group": "Flow",
            "summary": "A strip saying what just happened, with the undo inside it. Cheaper than a dialog before every action when the reader is doing the same thing forty times \u2014 and only honest when the action really is reversible.",
            "specimens": [
                {
                    "caption": "After a bulk edit",
                    "snippet": "KvitConfirmInPlace {\n    width: parent.width\n    shown: true\n    text: \"Recategorised as Groceries\"\n    affected: 40\n}"
                },
            ]
        },
        {
            "name": "KvitTimeline",
            "group": "Flow",
            "summary": "What happened to something, newest first, with who did it. In an estate where an agent and a person change the same things, who is the column that makes a history worth reading.",
            "specimens": [
                {
                    "caption": "An account's recent history",
                    "snippet": "KvitTimeline {\n    width: parent.width\n    label: \"Account history\"\n    entries: [\n        { when: \"14:02\", what: \"Statement imported\", who: \"agent\",\n          detail: \"412 transactions, 8 unmatched\", tone: \"success\" },\n        { when: \"11:20\", what: \"Two rows disputed\", who: \"you\",\n          tone: \"warning\" },\n        { when: \"Yesterday\", what: \"Account opened\", who: \"you\" }\n    ]\n}"
                },
            ]
        },
        {
            "name": "KvitNumberField",
            "group": "Flow",
            "summary": "A number the reader types, right-aligned in tabular numerals. It validates and says why rather than refusing keystrokes \u2014 a field that ignores a key gives no reason, and the usual cause is a decimal separator the reader's locale writes differently.",
            "specimens": [
                {
                    "caption": "Integer and decimal, valid and out of range",
                    "snippet": "Column {\n    spacing: Interface.spaceLoose\n    KvitNumberField { label: \"Days\"; text: \"14\"; minimum: 1; maximum: 365 }\n    KvitNumberField { label: \"Rate\"; text: \"4.25\"; decimals: 2 }\n    KvitNumberField { label: \"Days\"; text: \"999\"; minimum: 1; maximum: 365 }\n}"
                },
            ]
        },
        {
            "name": "KvitMoneyField",
            "group": "Flow",
            "summary": "An amount of money, in minor units. The reader types 12.34 and the field reports 1234, an integer \u2014 money in floating point drifts by a penny somewhere nobody can find. The decimal count comes from the currency.",
            "specimens": [
                {
                    "caption": "Sterling, yen and dinar",
                    "snippet": "Column {\n    spacing: Interface.spaceLoose\n    KvitMoneyField { label: \"Amount\"; text: \"1284.50\"; currency: \"GBP\" }\n    KvitMoneyField { label: \"Amount\"; text: \"4200\"; currency: \"JPY\"\n                     minorDigits: 0 }\n    KvitMoneyField { label: \"Amount\"; text: \"18.750\"; currency: \"BHD\"\n                     minorDigits: 3 }\n}"
                },
            ]
        },
        {
            "name": "KvitDualList",
            "group": "Flow",
            "summary": "Two lists with items moving between them: what is available, and what is chosen and in what order. Everything works from the keyboard, which most implementations of this shape do not.",
            "specimens": [
                {
                    "caption": "Choosing table columns",
                    "snippet": "KvitDualList {\n    width: parent.width; height: Interface.px(220)\n    available: [\"Payee\", \"Account\", \"Tags\", \"Note\"]\n    chosen: [\"Date\", \"Description\", \"Amount\", \"Balance\"]\n}"
                },
            ]
        },
        {
            "name": "KvitSpotlight",
            "group": "Flow",
            "summary": "Darken everything except one region and say something about it \u2014 the primitive under a guided tour, and deliberately only the primitive: what drives the stepping differs between the two applications that want one.",
            "specimens": [
                {
                    "caption": "Focusing a button",
                    "snippet": "Item {\n    width: parent.width; height: Interface.px(160)\n    KvitButton { id: target; text: \"Import a statement\"; form: \"primary\"\n                 x: Interface.px(40); y: Interface.px(20) }\n    KvitSpotlight {\n        shown: true\n        target: target\n        title: \"Start here\"\n        detail: \"Import a statement and the dashboard fills itself in.\"\n    }\n}"
                },
            ]
        },
    ]
}
