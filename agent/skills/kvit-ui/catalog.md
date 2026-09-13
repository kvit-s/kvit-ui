# The kvit-ui component catalogue

GENERATED FILE — do not edit. Written by `kvit-ui-gallery --catalog`
from the same list the gallery renders and the test suite compiles, so
what is described here is what the repository contains.

Every component is in the `Kvit.Ui` module: one `import Kvit.Ui` reaches
all of them, the token singletons `Theme`, `Interface` and `Typography`,
and the icon font.

Each entry gives what the component is for, the properties it declares, and
a working sample. The samples are compiled by `tests/test_gallery`, so one
that does not work stops the build.

## What there is

74 components, grouped by what they are for.

**Foundation** — `KvitLabel`, `KvitIcon`, `KvitIconButton`, `KvitLink`
**Structure** — `KvitHeader`, `KvitSidebar`, `KvitSidebarItem`, `KvitBreadcrumb`, `KvitRegion`, `KvitViewHead`, `KvitStatusBar`, `KvitWindow`
**Content** — `KvitSectionHeading`, `KvitRow`, `KvitSlimRow`, `KvitCard`, `KvitPanel`, `KvitPane`, `KvitDivider`, `KvitDisclosure`, `KvitEmptyState`
**Marks** — `KvitChip`, `KvitTag`, `KvitBadge`, `KvitSlug`, `KvitDot`, `KvitSignal`, `KvitPip`
**Quantities** — `KvitFigure`, `KvitBeforeAfter`
**Controls** — `KvitButton`, `KvitChipButton`, `KvitStepper`, `KvitField`, `KvitTextArea`, `KvitSearchField`, `KvitCheck`, `KvitSelect`, `KvitTab`
**Feedback** — `KvitTooltip`, `KvitPopover`, `KvitHint`, `KvitHoverCard`, `KvitToast`, `KvitNotice`, `KvitDialog`
**Data** — `KvitBar`, `KvitStackedBar`, `KvitSpark`, `KvitTrend`, `KvitDistribution`, `KvitGauge`, `KvitDelta`, `KvitStatTile`, `KvitFigureBlock`, `KvitCell`, `KvitTable`
**Flow** — `KvitScrollBar`, `KvitMenu`, `KvitMenuItem`, `KvitTree`, `KvitSwitch`, `KvitRadioGroup`, `KvitProgress`, `KvitSlider`, `KvitSplitView`, `KvitSegmented`, `KvitTypeAhead`, `KvitConfirmInPlace`, `KvitTimeline`, `KvitNumberField`, `KvitMoneyField`, `KvitDualList`, `KvitSpotlight`

## Symbols

`KvitIcon` and every component that takes a `symbol` accept these names.
They say what a symbol means rather than what it looks like, so the drawing
can change without touching a call site. A name outside this list draws a
marked placeholder and fails `tests/test_components`.

`plus`, `messages-square`, `message-square`, `send`, `chevron-right`, `chevron-down`, `zoom-in`, `zoom-out`, `pencil`, `check`, `archive`, `rotate-ccw`, `chevron-left`, `chevron-up`, `close`, `sidebar`, `list`, `search`, `note`, `tag`, `caret-sort`, `sort-ascending`, `sort-descending`, `dot`, `dot-outline`, `circle`, `square`, `minus`, `info`, `warning`, `error`, `success`, `question`, `settings`, `filter`, `columns`, `calendar`, `clock`, `user`, `robot`, `folder`, `file`, `link`, `external`, `copy`, `trash`, `undo`, `redo`, `refresh`, `more`, `more-vertical`, `drag`, `pin`, `star`, `eye`, `eye-off`, `lock`, `arrow-up`, `arrow-down`, `arrow-left`, `arrow-right`, `trend-up`, `trend-down`, `chart`, `wallet`, `coins`, `bank`, `receipt`, `repeat`, `split`, `merge`, `play`, `pause`, `stop`, `attachment`, `diff`, `git`, `agent`, `ask`, `ask-in`, `rename`, `terminal`, `chat`

## The components

### KvitLabel

A run of chrome text at one of the seven type roles. Every other component uses it, which is what keeps the chrome family, the rendering mode and the eliding rule in one place.

| Property | Type | |
|---|---|---|
| `role` | string | "caption" \| "small" \| "body" \| "strong" \| "title" \| "headline" \| "display" |
| `mono` | bool | Setting this to `true` draws the text in the monospace family, for an identifier that has to line up down a column. |
| `tabular` | bool | Tabular numerals: every digit the same width, so a column of figures lines up and a changing value does not shift the text beside it. |

*The seven roles*

```qml
Column {
    spacing: Interface.spaceSnug
    KvitLabel { text: "display — a page title"; role: "display" }
    KvitLabel { text: "headline — a pane title"; role: "headline" }
    KvitLabel { text: "title — a section heading"; role: "title" }
    KvitLabel { text: "strong — a name"; role: "strong" }
    KvitLabel { text: "body — row text and prose"; role: "body" }
    KvitLabel { text: "small — chip labels and sub-lines"; role: "small" }
    KvitLabel { text: "caption — kind tags and counts"; role: "caption" }
}
```

*Monospace and tabular numerals*

```qml
Column {
    spacing: Interface.spaceSnug
    KvitLabel { text: "a1b2c3d4"; mono: true }
    KvitLabel { text: "1,234.56"; tabular: true }
}
```

### KvitIcon

A symbol asked for by what it means rather than by what it looks like. The font ships with the module, so importing Kvit.Ui is enough.

| Property | Type | |
|---|---|---|
| `name` | string | **required.** What the symbol means. |
| `color` | color |  |
| `designSize` | readonly int | Call sites give an icon a size by setting width and height, and the glyph is drawn at whichever is smaller, so a 13-pixel combo indicator stays 13 pixels and an unconstrained icon stays 18. |
| `glyph` | readonly string |  |
| `recognized` | readonly bool |  |

*At the sizes the chrome uses*

```qml
Row {
    spacing: Interface.space
    KvitIcon { name: "search"; width: Interface.iconSizeSmall; height: width }
    KvitIcon { name: "chevron-right" }
    KvitIcon { name: "trash"; color: Theme.danger }
    KvitIcon { name: "success"; color: Theme.success }
}
```

*An unrecognised name is visible, not blank*

```qml
KvitIcon { name: "not-a-symbol" }
```

### KvitIconButton

A button whose whole label is a symbol. A real AbstractButton, so it takes tab focus and a screen reader is told it is there; `label` fills both the accessible name and the tooltip shown on pointer hover or keyboard focus, and `explanation` is the second sentence beside it. `dense` draws the symbol at the 13 every symbol beside words in this library is drawn at, rather than at 18.

| Property | Type | |
|---|---|---|
| `symbol` | string | **required.** What the symbol means: an IconCatalog meaning name. |
| `label` | string | **required.** What the button does, in words. |
| `form` | string | "ordinary" \| "quiet". |
| `dense` | bool | Draw the symbol at 13 rather than at 18. |
| `tooltipEnabled` | bool | A control that opens a labelled explanation already presents these words in that surface. |
| `explanation` | string | One sentence saying more than the label can: why the button is disabled, or what pressing it opens. |
| `iconColor` | color |  |

*Quiet, ordinary and checked*

```qml
Row {
    spacing: Interface.space
    KvitIconButton { symbol: "pencil"; label: "Edit" }
    KvitIconButton { symbol: "trash"; label: "Delete"; form: "ordinary" }
    KvitIconButton { symbol: "pin"; label: "Pin"; checked: true }
    KvitIconButton { symbol: "copy"; label: "Copy"; enabled: false }
    KvitIconButton {
        symbol: "settings"; label: "Settings"
        Component.onCompleted: if (Window.window)
            forceActiveFocus(Qt.TabFocusReason)
    }
}
```

*Two symbol sizes*

```qml
// A symbol standing on its own is drawn at 18. `dense` draws it at 13,
// which is the size every symbol beside words in this library is drawn
// at — a link's, a select's indicator, a section heading's chevron.
// A strip of buttons across a header or hoisted onto a heading bar is
// that second case, and at 18 they read as larger than the row they
// sit in.
Column {
    spacing: Interface.space
    Row {
        spacing: Interface.space
        KvitIconButton { symbol: "search"; label: "Search" }
        KvitIconButton { symbol: "plus"; label: "New track" }
        KvitIconButton { symbol: "git"; label: "Branch" }
    }
    Row {
        spacing: Interface.space
        KvitIconButton { symbol: "search"; label: "Search"; dense: true }
        KvitIconButton { symbol: "plus"; label: "New track"; dense: true }
        KvitIconButton { symbol: "git"; label: "Branch"; dense: true }
    }
}
```

*A second sentence, where the name is not enough*

```qml
// `label` names the button and is what a screen reader is told.
// `explanation` is the sentence beside it, for what the name cannot
// hold: what pressing this changes, or why it cannot be pressed. Both
// go to the tooltip and to the accessible description, in that order,
// so a pointer reader and a screen reader are told the same thing.
Row {
    spacing: Interface.space
    KvitIconButton {
        symbol: "rename"; label: "Rename"
        explanation: "Renames the track and its branch. Its history "
                     + "and folder stay unchanged."
        Component.onCompleted: if (Window.window)
            forceActiveFocus(Qt.TabFocusReason)
    }
    KvitIconButton {
        symbol: "archive"; label: "Archive"; enabled: false
        explanation: "The branch has work that has not been pushed."
    }
}
```

### KvitLink

An inline destination or action with link semantics, natural width and optional symbol. Hover and keyboard focus both add accent and an underline; no chevron, button ground or border is invented. `explanation` says in a sentence where following it goes, shown as the tooltip and announced as the accessible description.

| Property | Type | |
|---|---|---|
| `symbol` | string |  |
| `role` | string |  |
| `elide` | int |  |
| `explanation` | string | One sentence saying where following this goes, beside its own words. |

*Plain, symbolic and keyboard-focused*

```qml
Row {
    spacing: Interface.spaceLoose
    KvitLink { text: "Privacy policy" }
    KvitLink { text: "Open calendar"; symbol: "calendar" }
    KvitLink {
        text: "Keyboard focus"
        Component.onCompleted: if (Window.window)
            forceActiveFocus(Qt.TabFocusReason)
    }
}
```

*Elided in a narrow column*

```qml
KvitLink {
    width: Interface.px(120)
    text: "A destination whose full name does not fit here"
}
```

### KvitHeader

The strip across the top of the window: the wordmark, navigation in the middle, actions on the right, in that fixed order on every screen.

| Property | Type | |
|---|---|---|
| `wordmark` | string |  |
| `actions` | alias |  |

*With navigation and actions*

```qml
KvitHeader {
    width: parent.width
    wordmark: "kvit"
    Row {
        spacing: Interface.space
        anchors.verticalCenter: parent.verticalCenter
        KvitTab { text: "Today"; selected: true }
        KvitTab { text: "Projects"; count: 26 }
        KvitTab { text: "Decisions"; count: 4 }
    }
    actions: [
        Component { KvitIconButton { symbol: "settings"; label: "Settings" } }
    ]
}
```

### KvitSidebar

The list of places down the left edge. Below the laptop breakpoint it collapses to its rail and every item becomes its symbol alone.

| Property | Type | |
|---|---|---|
| `collapsed` | bool |  |

*Expanded and collapsed*

```qml
Row {
    spacing: Interface.columnGap
    KvitSidebar {
        width: Interface.sidebarWidth
        KvitSidebarItem { text: "Today"; symbol: "calendar"; selected: true }
        KvitSidebarItem { text: "Projects"; symbol: "folder"; count: 26 }
        KvitSidebarItem { text: "Decisions"; symbol: "question"; count: 4 }
    }
    KvitSidebar {
        width: Interface.railWidth
        collapsed: true
        KvitSidebarItem { text: "Today"; symbol: "calendar"; selected: true }
        KvitSidebarItem { text: "Projects"; symbol: "folder" }
        KvitSidebarItem { text: "Decisions"; symbol: "question" }
    }
}
```

### KvitSidebarItem

One place in the sidebar. The selected item takes a bar down its leading edge as well as a tint, so the selection is not resting on colour. A count is hidden once the sidebar collapses to the rail unless `countInRail` asks for it, and `countMax` raises the badge's cap where the place the item points at states the true number.

| Property | Type | |
|---|---|---|
| `symbol` | string | **required.**  |
| `selected` | bool |  |
| `collapsed` | bool |  |
| `count` | int |  |
| `counted` | string | The word for one of the things counted. |
| `countedPlural` | string | The word for several of them. |
| `countInRail` | bool | Whether the count stays visible once the sidebar has collapsed to the rail. |
| `countMax` | int | The largest count the badge writes out; past it the badge says "99+". |

*Selected, hovered and counted*

```qml
Column {
    width: Interface.sidebarWidth
    KvitSidebarItem {
        text: "Inbox"; symbol: "note"; selected: true
        count: 12; counted: "message"
    }
    KvitSidebarItem { text: "Archive"; symbol: "archive" }
}
```

*The rail, with and without its counts*

```qml
Row {
    spacing: Interface.columnGap

    // The default: a rail is a strip of symbols, and the count stays behind
    // with the label.
    KvitSidebar {
        width: Interface.railWidth
        collapsed: true
        KvitSidebarItem { text: "Review"; symbol: "question"; count: 214 }
        KvitSidebarItem { text: "Accounts"; symbol: "bank" }
    }

    // A rail that has to keep the backlog in front of the reader. The badge
    // moves onto the symbol's upper corner, and `countMax` is raised past
    // KvitBadge's default of 99 because the screen this points at states the
    // true number — a badge reading "99+" beside it disagrees with it.
    KvitSidebar {
        width: Interface.railWidth
        collapsed: true
        KvitSidebarItem {
            text: "Review"; symbol: "question"
            count: 214; countInRail: true; countMax: 999
        }
        KvitSidebarItem { text: "Accounts"; symbol: "bank" }
    }
}
```

### KvitBreadcrumb

Where the reader is and the way back. The last crumb is the current place and is deliberately not a link; a long trail elides from the middle, keeping the section and the current place.

| Property | Type | |
|---|---|---|
| `trail` | var | A list of { label, id } objects, root first. |
| `maximumVisible` | int |  |
| `shown` | readonly var |  |

*A short trail and a long one*

```qml
Column {
    spacing: Interface.space
    KvitBreadcrumb {
        trail: [{ label: "Projects", id: "p" },
                { label: "kvit-ui", id: "u" }]
    }
    KvitBreadcrumb {
        trail: [{ label: "Projects", id: "p" }, { label: "kvit", id: "k" },
                { label: "Wave 1", id: "w" }, { label: "Tokens", id: "t" },
                { label: "Theme", id: "h" }]
    }
}
```

### KvitRegion

A body that takes the height left over and scrolls what does not fit. The scroll bar sits beside the content rather than over it, so nothing is ever hidden behind it.

| Property | Type | |
|---|---|---|
| `contentHeight` | alias |  |
| `contentWidth` | alias |  |
| `padding` | int | Padding inside the scrolled area. |
| `horizontal` | bool |  |

*Scrolling a column of rows*

```qml
KvitRegion {
    width: 480; height: Interface.px(140)
    Column {
        id: rows
        width: 460
        Repeater {
            model: 12
            KvitSlimRow { width: rows.width; name: "Row " + (index + 1) }
        }
    }
}
```

### KvitViewHead

The strip at the top of a view: what the view is, how much is in it, and the controls that act on all of it.

| Property | Type | |
|---|---|---|
| `title` | string |  |
| `count` | int | How many things the view holds, and the word for one of them. |
| `counted` | string |  |
| `countedPlural` | string | The word for several of them. |
| `subtitle` | string |  |
| `countPhrase` | readonly string | "1 transaction", "250,000 transactions". |

*Title, count and controls*

```qml
KvitViewHead {
    width: parent.width
    title: "Transactions"
    count: 1284
    counted: "transaction"
    subtitle: "Everything since the account was opened"
    KvitSearchField { width: Interface.px(180) }
    KvitButton { text: "Export"; form: "ordinary" }
}
```

### KvitStatusBar

The strip along the bottom: what is happening on the left, standing facts on the right. Facts come in three shapes — plain strings, named groups whose facts open what they name, and whole controls at the end. What the bar has no room for goes into a menu behind a control saying how many there are, rather than being cut off the end of the list. `groupsAt: "left"` puts the groups before the activity, for a bar whose left end is the list of what is waiting.

| Property | Type | |
|---|---|---|
| `activity` | string | What is happening now. |
| `facts` | var | Standing facts, right aligned. |
| `groups` | var | Groups of facts that do something. |
| `controls` | alias | One or two controls, after the facts. |
| `groupsAt` | string | "right" \| "left": which end of the bar the groups sit at. |
| `shownGroups` | int | How many of the groups are drawn on the bar. |
| `groupsFirst` | readonly bool | The groups come before the activity rather than after it. |
| `activityFloor` | readonly int | How much of the bar the activity keeps before the groups start taking room from it: about eight words at the default interface size. |
| `hiddenFacts` | readonly var | The facts that did not fit, flattened into what the menu draws. |

*Working, with two facts*

```qml
KvitStatusBar {
    width: parent.width
    activity: "Reindexing 4 of 26 projects"
    facts: ["1,284 notes", "last synced 14:02"]
}
```

*Grouped facts that open what they name*

```qml
// Each fact is a control: it takes tab focus, says its own name to a
// screen reader and answers Return and Space. The bar keeps none of the
// caller's state — it is handed words and hands back which one was
// pressed.
KvitStatusBar {
    id: bar
    width: parent.width
    property string lastPressed: ""

    activity: "Indexing 4 of 26 working copies"
    groups: [
        {
            "label": "Waiting on you",
            "facts": [
                { "text": "3 reviews", "symbol": "question" },
                { "text": "1 conflict", "symbol": "warning" }
            ]
        },
        {
            "label": "Running",
            "facts": [
                { "text": "2 agents", "symbol": "robot" }
            ]
        }
    ]
    onFactActivated: (group, fact) =>
        bar.lastPressed = group + "/" + fact
}
```

*Narrow: what does not fit is in the menu, not gone*

```qml
// The same two groups in a bar too narrow for them. The groups that fit
// are drawn left to right; the rest are behind the count of what is not
// there, which is a control — tab reaches it, Return opens the menu,
// and the arrow keys move through it. Nothing is dropped, and nothing is
// silent about it.
Column {
    spacing: Interface.columnGap
    width: parent.width

    KvitStatusBar {
        width: Math.min(parent.width, Interface.px(420))
        activity: "Indexing"
        groups: [
            {
                "label": "Waiting on you",
                "facts": [
                    { "text": "3 reviews", "symbol": "question" },
                    { "text": "1 conflict", "symbol": "warning" }
                ]
            },
            {
                "label": "Running",
                "facts": [{ "text": "2 agents", "symbol": "robot" }]
            }
        ]
    }

    // The resting height is unchanged: a bar holding only text is as tall as
    // the text, and only a control in the slot at the end makes it grow.
    KvitStatusBar {
        width: Math.min(parent.width, Interface.px(320))
        facts: ["1,284 notes", "last synced 14:02"]
    }
}
```

*Groups first, for a bar whose left end is the work*

```qml
// A bar whose left end is the work rather than a sentence about it: what is
// waiting and what is running, each item something to press, with the
// standing facts at the right. `groupsAt: "left"` is the whole of the
// difference — the groups come before the activity, and the plain facts and
// the controls stay where they were.
KvitStatusBar {
    width: parent.width
    groupsAt: "left"
    activity: "Fetching origin"
    groups: [
        {
            "label": "Waiting on you",
            "facts": [
                { "text": "3 reviews", "symbol": "question" },
                { "text": "1 conflict", "symbol": "warning" }
            ]
        },
        {
            "label": "Running",
            "facts": [{ "text": "2 agents", "symbol": "robot" }]
        }
    ]
    facts: ["7 changes", "main, 2 ahead"]
}
```

### KvitWindow

The application shell: header, optional sidebar, body and status bar, with the sidebar collapsing to its rail below the laptop breakpoint. Not shown here because it is a Window; see the code sample.

| Property | Type | |
|---|---|---|
| `header` | alias |  |
| `sidebar` | alias |  |
| `body` | alias |  |
| `statusBar` | alias |  |
| `sidebarVisible` | bool | False hides the sidebar entirely, which is different from collapsing it: a hidden sidebar is a reader's choice and a collapsed one is the window being narrow. |
| `hasHeader` | bool | Whether the window has a header and a status bar at all. |
| `hasStatusBar` | bool |  |
| `sidebarCollapsed` | readonly bool |  |
| `narrow` | readonly bool |  |

*The whole shell*

```qml
// Source only: a Window has no place inside another window's item tree, so
// this is the one page that shows a sample without drawing it. It is compiled
// and run by tests/test_gallery like every other sample.
//
// Each of the four slots takes one item, and each fills the slot it is in.
KvitWindow {
    id: window
    title: "kvit-cash"

    header: KvitHeader {
        anchors.fill: parent
        wordmark: "kvit"
        Row {
            anchors.verticalCenter: parent.verticalCenter
            spacing: Interface.space
            KvitTab { text: "Accounts"; selected: true }
            KvitTab { text: "Budget" }
        }
    }

    sidebar: KvitSidebar {
        anchors.fill: parent
        // The window says when it is narrow; the sidebar draws itself as a
        // rail when it is.
        collapsed: window.sidebarCollapsed
        KvitSidebarItem { text: "Everyday"; symbol: "wallet"; selected: true }
        KvitSidebarItem { text: "Savings"; symbol: "bank" }
        KvitSidebarItem { text: "Budget"; symbol: "chart-line" }
    }

    body: KvitRegion {
        anchors.fill: parent
        Column {
            id: rows
            width: parent.width
            KvitSlimRow { width: rows.width; name: "Checking" }
            KvitSlimRow { width: rows.width; name: "Joint checking" }
            KvitSlimRow { width: rows.width; name: "Emergency fund" }
        }
    }

    statusBar: KvitStatusBar {
        anchors.fill: parent
        activity: "Matching 4 of 26 statements"
        facts: ["1,284 transactions", "last synced 14:02"]
    }
}
```

### KvitSectionHeading

A group heading: a filled bar with a disclosure chevron, the name, what the group holds, the count with the word for what was counted, and the one action that applies to every row under it. A count that is not a number goes in `countText` and is drawn beside the name as written; `actionSymbol` draws the action as a symbol and keeps its words as the button's name and tooltip, with `actionExplanation` saying in a sentence what running it does. A collapsible heading joins the tab order and opens on Return, Enter or Space; the hoisted action is a control of its own, so running it never also collapses the group.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `count` | int |  |
| `counted` | string | What the count counts, singular. |
| `countedPlural` | string | The word for several of them. |
| `countText` | string | A count the caller has already written out, for a heading whose count is not a number. |
| `kind` | string | The kind of thing the group holds, said beside the name where the name alone does not say it. |
| `action` | string | The hoisted action, true for every row in the group. |
| `actionSymbol` | string | Draw that action as a symbol instead of as its words. |
| `actionExplanation` | string | One sentence saying what running the action does, beside its words. |
| `strong` | bool |  |
| `expanded` | bool |  |
| `collapsible` | bool |  |
| `countPhrase` | readonly string | "1 account", "1,200 accounts", or the grouped number on its own where the caller named nothing. |

*Collapsible, counted, with an action*

```qml
Column {
    width: parent.width
    spacing: Interface.space
    KvitSectionHeading {
        width: parent.width
        text: "Waiting on me"; counted: "project"; count: 4
        action: "Hand all to an agent"; collapsible: true
    }
    // A group of one is still counted. The count is drawn wherever the caller
    // gives one, so it does not come and go as the group changes size.
    KvitSectionHeading {
        width: parent.width
        text: "Waiting on Sam"; counted: "project"; count: 1
    }
    KvitSectionHeading {
        width: parent.width
        text: "Archived"; kind: "closed last quarter"
        collapsible: true; expanded: false; strong: true
    }
}
```

*Reached by tab, opened by Return*

```qml
Column {
    width: parent.width
    spacing: Interface.space

    // The ring says the heading has the keyboard; Return, Enter and Space
    // open and close it. The action beside it is a KvitLink, which consumes
    // its own key, so tabbing on to it and pressing Space runs the action
    // and leaves the group where it was.
    KvitSectionHeading {
        width: parent.width
        text: "Waiting on me"; counted: "project"; count: 4
        action: "Hand all to an agent"; collapsible: true
        Component.onCompleted: if (Window.window)
            forceActiveFocus(Qt.TabFocusReason)
    }
    KvitSectionHeading {
        width: parent.width
        text: "Archived"; kind: "closed last quarter"
        collapsible: true; expanded: false
    }
}
```

*A written count, and an action drawn as a symbol*

```qml
// Some counts are not numbers. A Changes heading reads "1 · +0 −0" —
// one changed file, and the lines added and removed across it — which
// the service composes and hands over already written, and which no
// integer expresses. `countText` is drawn exactly as given, beside the
// name rather than at the right end where a number goes.
//
// `actionSymbol` draws the hoisted action as a symbol. The words stay
// in `action` and become the button's accessible name and its tooltip,
// so nothing is lost by a reader who cannot see the symbol; what is
// gained is a column of headings a reader can scan rather than read.
Column {
    width: parent.width
    spacing: Interface.space
    KvitSectionHeading {
        width: parent.width
        text: "Changes"; countText: "1 · +0 −0"
        action: "Open the diff"; actionSymbol: "diff"
        collapsible: true
    }
    KvitSectionHeading {
        width: parent.width
        text: "Agents"; countText: "2 running"
        action: "Start an agent"; actionSymbol: "plus"
        collapsible: true
    }
    KvitSectionHeading {
        width: parent.width
        text: "Terminal"
        action: "Open a terminal"; actionSymbol: "terminal"
    }
}
```

*What the action does, in a sentence*

```qml
// An action drawn as a symbol has nowhere to say what it does beyond its
// own name, and "Start an agent" leaves out the part a reader wants
// before pressing it. `actionExplanation` is the sentence beside the
// name, on the button rather than on the heading: a heading is not
// something a reader presses, and a sentence announced there would be
// read out on arriving at the group instead of on reaching the action.
//
// Hover either button, or tab to it, to read it. It is the accessible
// description too, which a screen reader reads after the name.
Column {
    width: parent.width
    spacing: Interface.space
    KvitSectionHeading {
        width: parent.width
        text: "Agents"; countText: "2 running"
        action: "Start an agent"; actionSymbol: "plus"
        actionExplanation: "Creates an agent at this project’s root. An agent may change files. Nothing runs until you send its first message."
        collapsible: true
    }
    KvitSectionHeading {
        width: parent.width
        text: "Changes"; countText: "1 · +0 −0"
        action: "Review all"
        actionExplanation: "Opens every changed file as one diff."
        collapsible: true
    }
}
```

### KvitRow

A list row at one of four heights, chosen by what the row carries rather than by how many rows a view wants to fit. Hover and keyboard focus are separate tints, because they are different rows. A row that does something when it is pressed says so first: the hover tint and the tap both follow `interactive`, which is on for a row that takes `activeFocusOnTab` and set by hand for a row inside a list that carries its own cursor.

| Property | Type | |
|---|---|---|
| `form` | string | "full" \| "sub" \| "slim" \| "compact" |
| `selected` | bool |  |
| `current` | bool |  |
| `rule` | bool |  |
| `label` | string | What a screen reader says when the caret reaches this row. |
| `interactive` | bool | Whether pressing this row does something. |
| `hovered` | readonly bool |  |

*Pressable and static*

```qml
Column {
    id: rows
    width: parent.width
    property int opened: 0

    // A row that opens something is reachable by the keyboard, and takes the
    // hover tint and the tap with it.
    KvitRow {
        width: rows.width; label: "Groceries"
        activeFocusOnTab: true
        onActivated: rows.opened++
        KvitLabel { anchors.centerIn: parent; text: "opens the record" }
    }

    // A field name beside its value declares nothing, draws no hover tint
    // and answers no tap.
    KvitRow {
        width: rows.width; label: "Amount"
        KvitLabel { anchors.centerIn: parent; text: "layout only" }
    }
}
```

*Opened from the keyboard, without opening it twice*

```qml
Column {
    id: rows
    width: parent.width
    property int opened: 0

    // Tab reaches this row, the ring says which row the keyboard is on, and
    // Return, Enter or Space opens it. A screen reader's press action does
    // the same one thing.
    KvitRow {
        width: rows.width; form: "sub"; label: "Groceries"
        activeFocusOnTab: true
        onActivated: rows.opened++
        Component.onCompleted: if (Window.window)
            forceActiveFocus(Qt.TabFocusReason)
        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            text: "keyboard focus — Return opens it"
        }
    }

    // A row is usually a container, and the key may be meant for something
    // inside it. The row answers only while the row itself has the keyboard,
    // so Space on this button presses the button and leaves the record shut.
    KvitRow {
        width: rows.width; form: "sub"; label: "Rent"
        activeFocusOnTab: true
        onActivated: rows.opened++
        KvitLabel {
            anchors.verticalCenter: parent.verticalCenter
            text: "the button takes its own Space"
        }
        KvitButton {
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            text: "Split"; form: "quiet"
        }
    }
}
```

*The four heights, and the three states*

```qml
Column {
    width: parent.width
    KvitRow { width: parent.width; form: "full"; label: "Full"
        KvitLabel { anchors.centerIn: parent; text: "full — 56" } }
    KvitRow { width: parent.width; form: "sub"; label: "Sub"
        KvitLabel { anchors.centerIn: parent; text: "sub — 48" } }
    KvitRow { width: parent.width; form: "slim"; selected: true; label: "Slim"
        KvitLabel { anchors.centerIn: parent; text: "slim — 30, selected" } }
    KvitRow { width: parent.width; form: "compact"; current: true; label: "Compact"
        KvitLabel { anchors.centerIn: parent; text: "compact — 24, keyboard focus" } }
}
```

### KvitSlimRow

The row a reader sees most: a name, what it is, one phrase about it and one figure, right aligned, in that order on every screen.

| Property | Type | |
|---|---|---|
| `name` | string |  |
| `kind` | string | What kind of thing this is, said next to the name where the name does not say it. |
| `phrase` | string | One phrase about its state. |
| `figure` | string |  |
| `unit` | string |  |
| `measured` | bool | False draws an em dash rather than the figure, which is what an unmeasured value looks like. |
| `symbol` | string |  |

*Measured and unmeasured*

```qml
Column {
    width: parent.width
    KvitSlimRow { width: parent.width; symbol: "folder"; name: "kvit-ui"
        kind: "library"; phrase: "waiting on review"; figure: "3.5"; unit: "d" }
    KvitSlimRow { width: parent.width; symbol: "folder"; name: "kvit-cash"
        kind: "application"; phrase: "not started"; measured: false }
}
```

### KvitCard

A bounded block of content on a surface. An interactive card takes the control-boundary token and the hover tint, so a card that responds to a click looks like it will before it is clicked.

| Property | Type | |
|---|---|---|
| `interactive` | bool |  |
| `selected` | bool |  |
| `hovered` | alias |  |
| `padding` | int |  |

*Static and interactive*

```qml
Row {
    spacing: Interface.columnGap
    KvitCard { KvitLabel { text: "A card" } }
    KvitCard { interactive: true; KvitLabel { text: "Interactive" } }
    KvitCard { selected: true; KvitLabel { text: "Selected" } }
}
```

### KvitPanel

A region of the window with its own ground: a sidebar, a toolbar strip. Structural, where a card is content — which is why it has no radius.

| Property | Type | |
|---|---|---|
| `ruleTop` | bool | Which edges carry a rule. |
| `ruleBottom` | bool |  |
| `ruleLeft` | bool |  |
| `ruleRight` | bool |  |

*With rules on two edges*

```qml
KvitPanel {
    width: parent.width; height: Interface.px(60)
    ruleTop: true; ruleBottom: true
    KvitLabel { anchors.centerIn: parent; text: "A panel" }
}
```

### KvitPane

The side pane: detail about the one thing selected, the same width wherever it appears so the list beside it never reflows.

| Property | Type | |
|---|---|---|
| `title` | string |  |
| `closeLabel` | string | What the close control is called, both on hover and to a screen reader. |
| `open` | bool |  |

*Open*

```qml
Item {
    width: parent.width; height: Interface.px(160)
    KvitPane {
        height: parent.height
        title: "Transaction"
        KvitLabel { anchors.centerIn: parent; text: "Detail goes here" }
    }
}
```

### KvitDivider

A rule between two things. One design pixel, in the decorative border token rather than the control-boundary one.

| Property | Type | |
|---|---|---|
| `vertical` | bool |  |
| `inset` | int | A rule that stops short of the edges, for a divider inside a padded container where a full-bleed line would cut the padding in half. |

*Horizontal and vertical*

```qml
Column {
    width: parent.width
    spacing: Interface.space
    KvitDivider { width: parent.width }
    Row {
        height: Interface.px(24); spacing: Interface.space
        KvitLabel { text: "left"; height: parent.height }
        KvitDivider { vertical: true; height: parent.height }
        KvitLabel { text: "right"; height: parent.height }
    }
}
```

### KvitDisclosure

A trigger with a body under it. Takes a `group` so several sections can behave as an accordion, which is what most of the eighteen hand-rolled versions in the estate actually are.

| Property | Type | |
|---|---|---|
| `title` | string |  |
| `expanded` | bool |  |
| `group` | string | Sections sharing a group string behave as an accordion: opening one closes the rest. |
| `count` | int | How many things are inside, drawn beside the title. |
| `bodyHeight` | readonly real | How tall the body wants to be, taken from the preferred height of whatever the caller put inside it. |

*Open and closed*

```qml
Column {
    width: parent.width
    spacing: Interface.space
    KvitDisclosure {
        width: parent.width; title: "What changed"; count: 3; expanded: true
        KvitLabel { text: "Three files were rewritten." }
    }
    KvitDisclosure { width: parent.width; title: "What did not"; count: 12 }
}
```

### KvitEmptyState

What a view says when it has nothing to show: what would be here, why it is not, and the action that would fill it. Also the answer for a chart with no data, in place of an axis drawn around zeros. The compact form is the same sentence on one line, at the height of a slim row and starting where the rows start, for a section in a stack of sections that may each be empty.

| Property | Type | |
|---|---|---|
| `form` | string | "full" \| "compact". |
| `title` | string | What would be here. |
| `detail` | string | Why it is not, or what to do about it. |
| `symbol` | string |  |
| `action` | string | The action that would fill it, if there is one. |
| `dashed` | bool | Draw a dashed outline around the whole thing: this is a region something can be put into, rather than a region that happens to be empty. |
| `compact` | readonly bool |  |

*With an action*

```qml
KvitEmptyState {
    width: parent.width
    symbol: "wallet"
    title: "No transactions yet"
    detail: "Import a statement or add one by hand, and it will appear here."
    action: "Import a statement"
}
```

*Dashed, as a drop target*

```qml
KvitEmptyState {
    width: parent.width
    dashed: true
    symbol: "file-arrow-down"
    title: "Drop a statement here"
    detail: "CSV, OFX and QIF. The file is read on this machine and nothing is sent anywhere."
    action: "Choose a file"
}
```

*Compact, one line per empty section*

```qml
// A column of sections, each of which may have nothing in it. The full
// block is several times taller than the rows it stands in for, so a stack
// of them uses the one-line form; the words are the same words.
Column {
    width: parent.width
    spacing: Interface.spaceSnug

    KvitSectionHeading { width: parent.width; text: "Changes"; count: 0
        counted: "change" }
    KvitEmptyState {
        width: parent.width
        form: "compact"
        title: "No changes"
    }

    KvitSectionHeading { width: parent.width; text: "Agents"; count: 0
        counted: "agent" }
    KvitEmptyState {
        width: parent.width
        form: "compact"
        symbol: "robot"
        title: "No agents"
        detail: "none started here yet"
        action: "Start one"
    }
}
```

### KvitChip

A small labelled mark saying what kind of thing this is or what state it is in. The tone names a meaning rather than a colour, and every tone carries its outline as well as its tint so the distinction is not resting on hue. It is a mark and not a control: a chip that opens something when it is pressed is KvitChipButton, drawn from this same tone table.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `tone` | string | "neutral" \| "accent" \| "success" \| "warning" \| "danger" \| "info" |
| `strong` | bool | A filled chip rather than a tinted one, for the one chip on a row that is the point of the row. |
| `symbol` | string |  |
| `explanation` | string | One sentence saying what the word on the chip means and what the reader is expected to do about it. |
| `toneColor` | readonly color |  |

*Every tone, tinted and filled*

```qml
Column {
    spacing: Interface.space
    Row {
        spacing: Interface.spaceNear
        KvitChip { text: "neutral" }
        KvitChip { text: "accent"; tone: "accent" }
        KvitChip { text: "success"; tone: "success" }
        KvitChip { text: "warning"; tone: "warning" }
        KvitChip { text: "danger"; tone: "danger" }
        KvitChip { text: "info"; tone: "info" }
    }
    Row {
        spacing: Interface.spaceNear
        KvitChip { text: "settled"; tone: "success"; strong: true }
        KvitChip { text: "disputed"; tone: "danger"; strong: true; symbol: "warning" }
    }
}
```

*A state whose word is not enough to act on*

```qml
// A chip that names a kind of thing is finished at its word. A chip that
// names a state somebody has to act on is not, and `explanation` is where
// the rest of it goes: shown in the tooltip and announced as the
// accessible description, so both readers are told the same thing.
Row {
    spacing: Interface.spaceNear
    KvitChip {
        text: "stale"
        tone: "warning"
        explanation: "The note changed since this was staged. Update redoes "
                     + "it against the note as it stands now."
    }
    KvitChip {
        text: "read only"
        explanation: "Another session holds the write lease on this copy, "
                     + "so this turn only reads."
    }
}
```

### KvitTag

A label a person put there, in a colour they chose. Theme-independent, because the reader picked that red and would not expect it to change; the label colour is derived from the fill rather than taken from a token.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `tint` | string | A value from Theme.colorPalette. |
| `removable` | bool |  |
| `tinted` | readonly bool |  |
| `tintColor` | readonly color |  |

*Tinted, plain and removable*

```qml
Row {
    spacing: Interface.spaceNear
    KvitTag { text: "groceries"; tint: Theme.colorPalette[0] }
    KvitTag { text: "transport"; tint: Theme.colorPalette[2] }
    KvitTag { text: "untinted" }
    KvitTag { text: "removable"; tint: Theme.colorPalette[4]; removable: true }
}
```

### KvitBadge

A count attached to something else. Caps rather than growing wide, and hides at zero — a badge showing nought says look here about nothing. The number is drawn and announced with the reader's own digit grouping, and the noun beside it comes from the caller in two slots, a singular and a plural.

| Property | Type | |
|---|---|---|
| `count` | int |  |
| `max` | int |  |
| `counted` | string | The word for one of the things counted. |
| `countedPlural` | string | The word for several of them. |
| `tone` | string | "neutral" \| "accent" \| "danger" |
| `groupedCount` | readonly string | The count with the reader's own digit grouping. |
| `toneColor` | readonly color |  |

*Counts, capped, and hidden at zero*

```qml
Row {
    spacing: Interface.space
    KvitBadge { count: 3 }
    KvitBadge { count: 42; tone: "neutral" }
    KvitBadge { count: 1204; tone: "danger" }
    KvitBadge { count: 0 }
}
```

*What the count counts*

```qml
// Same three pills, three different sentences for a screen reader:
// "1 decision", "214 decisions", "214 entries". One already inflected
// word would be right at one count and wrong at every other.
Row {
    spacing: Interface.space
    KvitBadge { count: 1; counted: "decision" }
    KvitBadge { count: 214; counted: "decision"; max: 999 }
    KvitBadge {
        count: 214; max: 999; tone: "neutral"
        counted: "entry"; countedPlural: "entries"
    }
}
```

### KvitSlug

An identifier: a reference, a hash, a key. Monospace, because the task is comparison rather than reading, and elided from the middle because the end is what distinguishes one from its neighbours.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `ground` | bool | A flat run of monospace text, for a slug inside a table cell where a ground on every row would read as a column of boxes. |

*With and without a ground*

```qml
Column {
    spacing: Interface.spaceNear
    KvitSlug { text: "TX-00173404" }
    KvitSlug { text: "9f2c1ab4e77d0031"; ground: false }
}
```

### KvitDot

A small filled circle standing for one thing's state. The shape is the second channel: a level that differs only by hue says nothing to a reader who cannot separate red from amber.

| Property | Type | |
|---|---|---|
| `color` | color |  |
| `shape` | string | "circle" \| "square" \| "diamond" |
| `hollow` | bool |  |
| `label` | string | What the dot means, for a screen reader. |

*Three shapes, three levels*

```qml
Row {
    spacing: Interface.space
    KvitDot { color: Theme.success; label: "healthy" }
    KvitDot { color: Theme.warning; shape: "square"; label: "slipping" }
    KvitDot { color: Theme.danger; shape: "diamond"; label: "stalled" }
    KvitDot { color: Theme.textMuted; hollow: true; label: "not measured" }
}
```

### KvitSignal

A mark saying what state something is in and how many things are in it, for a list where one row may have several of each. The number is drawn only past one, because a column of marks all reading `1` says nothing the mark did not already say. The colour is the caller’s, since which states exist is an application’s own question; the shape and the hollow form are the second channel, so states told apart by hue alone are not.

| Property | Type | |
|---|---|---|
| `count` | int |  |
| `max` | int |  |
| `color` | color |  |
| `shape` | string | "square" \| "circle" |
| `hollow` | bool | A mark drawn as an outline rather than filled, which is the third distinction a caller has without reaching for another hue. |
| `label` | string | **required.** What the mark means, in words, for a screen reader. |

*One of each, and several of one*

```qml
// At one the mark is the whole statement. Past one the count goes
// inside it, and the mark grows sideways to hold the digits rather
// than growing taller and pushing the row apart.
Column {
    spacing: Interface.space
    Row {
        spacing: Interface.space
        KvitSignal { color: Theme.accent; label: "1 needs you" }
        KvitSignal { color: Theme.success; shape: "circle"
            label: "1 running" }
        KvitSignal { color: Theme.danger; hollow: true; label: "1 failed" }
    }
    Row {
        spacing: Interface.space
        KvitSignal { count: 4; color: Theme.accent; label: "4 need you" }
        KvitSignal { count: 12; color: Theme.success; shape: "circle"
            label: "12 running" }
        KvitSignal { count: 240; color: Theme.danger; hollow: true
            label: "240 failed" }
    }
}
```

*Beside the words it marks*

```qml
// Where these are actually drawn: at the head of a row, in front of
// what the row is about. The mark is the height of a caption plus a
// margin, which is what lets it sit in a line of text without
// setting the line height.
Column {
    width: parent.width
    KvitRow {
        width: parent.width
        label: "Dialog Scout, 2 need you"
        Row {
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
            spacing: Interface.space
            KvitSignal { count: 2; color: Theme.accent
                label: "2 need you" }
            KvitLabel { text: "Dialog Scout"; role: "small" }
        }
    }
    KvitRow {
        width: parent.width
        label: "Nightly check, running"
        Row {
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
            spacing: Interface.space
            KvitSignal { color: Theme.success; shape: "circle"
                label: "1 running" }
            KvitLabel { text: "Nightly check"; role: "small" }
        }
    }
}
```

### KvitPip

A row of dots standing for a small count — three of five days recorded. For counts a reader takes in without counting; past about seven, the figure is faster.

| Property | Type | |
|---|---|---|
| `filled` | int |  |
| `total` | int |  |
| `color` | color |  |
| `label` | string |  |

*Three of five, and none of four*

```qml
Column {
    spacing: Interface.space
    KvitPip { filled: 3; total: 5; label: "3 of 5 days recorded" }
    KvitPip { filled: 0; total: 4; color: Theme.warning; label: "no checks passed" }
}
```

### KvitFigure

A measured value: tabular numerals, the unit in muted colour at the smaller role, and an em dash where nothing was measured rather than a zero. A balance nobody computed is not a balance of zero.

| Property | Type | |
|---|---|---|
| `value` | string |  |
| `unit` | string |  |
| `measured` | bool |  |
| `role` | string | "body" \| "small" \| "caption" \| "strong" \| "title" \| "headline" \| "display" |
| `color` | color |  |
| `bounded` | bool | A figure that is an upper bound rather than a measurement — the same distinction hatching carries on a bar. |

*Measured, unmeasured and bounded*

```qml
Column {
    spacing: Interface.spaceNear
    KvitFigure { value: "1,284.50"; unit: "GBP" }
    KvitFigure { value: "0"; unit: "d"; }
    KvitFigure { measured: false }
    KvitFigure { value: "12"; unit: "d"; bounded: true }
    KvitFigure { value: "94"; unit: "%"; role: "display" }
}
```

### KvitBeforeAfter

One value as it stands and the value something proposes to replace it with. Position, colour weight and the arrow all say which is which, so no reader depends on separating the two colours. A record being added has no before, and the em dash says so.

| Property | Type | |
|---|---|---|
| `before` | string | The values, already formatted. |
| `after` | string |  |
| `unit` | string | Drawn once after each value, in the muted colour, as KvitFigure does. |
| `beforeMeasured` | bool | False draws the em dash on that side: a record being added has no before, and one being deleted has no after. |
| `afterMeasured` | bool |  |
| `label` | string | What the pair is a value of — "Amount", "Category", "Date". |
| `role` | string | "body" \| "small" \| "caption" \| "strong" — passed to both figures, so the two are always the same size. |
| `unchanged` | readonly bool | True when the two values are the same, which a preview listing every field of a record needs in order to draw the unchanged ones quietly. |

*Changed, added, removed and unchanged*

```qml
Column {
    spacing: Interface.spaceNear
    KvitBeforeAfter { label: "Amount"; before: "42.00"; after: "44.50"; unit: "GBP" }
    KvitBeforeAfter { label: "Category"; before: ""; beforeMeasured: false; after: "Groceries" }
    KvitBeforeAfter { label: "Payee"; before: "TESCO 4471"; after: ""; afterMeasured: false }
    KvitBeforeAfter { label: "Date"; before: "2026-08-14"; after: "2026-08-14" }
    KvitBeforeAfter { before: "1,284.50"; after: "1,301.75"; unit: "GBP"; role: "strong" }
}
```

### KvitButton

A button with words on it, in three forms. One primary per screen region: it is the action the screen is for. `danger` is separate from the form, because a destructive action can be any of the three. `explanation` is one sentence saying what the words cannot — why it is disabled, or what pressing it opens — and it is read on a disabled button too.

| Property | Type | |
|---|---|---|
| `form` | string | "primary" \| "ordinary" \| "quiet" |
| `danger` | bool |  |
| `symbol` | string |  |
| `busy` | bool | A button that is doing something. |
| `busyText` | string |  |
| `explanation` | string | One sentence saying what the words on the button cannot: why it is disabled, or what pressing it opens. |
| `fill` | readonly color |  |

*Text, icon plus text, and busy text in every form*

```qml
Column {
    spacing: Interface.space
    Row {
        spacing: Interface.space
        KvitButton { text: "Save"; form: "primary" }
        KvitButton { text: "Save"; form: "ordinary" }
        KvitButton { text: "Save"; form: "quiet" }
    }
    Row {
        spacing: Interface.space
        KvitButton { text: "Schedule"; symbol: "calendar"; form: "primary" }
        KvitButton { text: "Schedule"; symbol: "calendar"; form: "ordinary" }
        KvitButton { text: "Schedule"; symbol: "calendar"; form: "quiet" }
    }
    Row {
        spacing: Interface.space
        KvitButton { text: "Save"; busy: true; form: "primary" }
        KvitButton { text: "Save"; busy: true; form: "ordinary" }
        KvitButton { text: "Save"; busy: true; form: "quiet" }
    }
}
```

*Destructive and disabled*

```qml
Row {
    spacing: Interface.space
    KvitButton { text: "Delete"; form: "primary"; danger: true }
    KvitButton { text: "Delete"; form: "ordinary"; danger: true; symbol: "trash" }
    KvitButton { text: "Disabled"; form: "ordinary"; enabled: false }
}
```

*A button that stays on, beside the same button off*

```qml
Row {
    spacing: Interface.space
    KvitButton { text: "Select region"; checkable: true; checked: true }
    KvitButton { text: "Select region"; checkable: true }
}
```

*Why a button cannot be pressed*

```qml
// The reason a control is in the state it is in often lives somewhere
// the reader cannot see. A disabled button with nothing to say about
// itself is a grey rectangle and no account of it; `explanation` is
// where that sentence goes, shown on hover and announced as the
// accessible description. It is read on a disabled button too, which
// is the case it exists for.
Row {
    spacing: Interface.space
    KvitButton {
        text: "Pull"
        explanation: "Brings the 2 commits on the remote into this branch."
    }
    KvitButton {
        text: "Push"; enabled: false
        explanation: "Nothing here has been committed yet."
    }
    KvitButton {
        text: "Archive"; enabled: false; danger: true
        explanation: "The branch has work that has not been pushed."
    }
}
```

### KvitChipButton

KvitChip's twin for a fact that opens something. Drawn from the same tone table, so a row mixing facts that act with facts that do not reads as one row; what separates them is what a real AbstractButton has anyway — a ground that changes under the pointer, a hand cursor, a focus ring and a button role. A chip that cannot be pressed keeps its place in the tab order and says why, a chip that can may say where it goes, and the chip whose destination is already open is drawn as the current one.

| Property | Type | |
|---|---|---|
| `tone` | string | "neutral" \| "accent" \| "success" \| "warning" \| "danger" \| "info" |
| `strong` | bool | A filled chip rather than a tinted one, for the one chip on a row that is the point of the row. |
| `symbol` | string |  |
| `trailingSymbol` | string | A symbol after the label, where the caller knows what pressing this leads to: a chevron for a list that opens beneath, an out-arrow for something that opens elsewhere. |
| `unavailableReason` | string | Why this chip cannot be pressed, in the reader's own terms. |
| `unavailable` | readonly bool |  |
| `explanation` | string | One sentence saying what pressing this opens, where the label alone does not say it. |
| `current` | bool | The chip whose destination is already open. |
| `toneColor` | readonly color |  |
| `lit` | readonly bool | Under the pointer or under the keyboard. |

*Every tone, tinted, filled and keyboard-focused*

```qml
Column {
    id: chips
    spacing: Interface.space
    property string lastPressed: ""

    Row {
        spacing: Interface.spaceNear
        KvitChipButton { text: "neutral"; onActivated: chips.lastPressed = text }
        KvitChipButton { text: "accent"; tone: "accent" }
        KvitChipButton { text: "success"; tone: "success" }
        KvitChipButton { text: "warning"; tone: "warning" }
        KvitChipButton { text: "danger"; tone: "danger" }
        KvitChipButton { text: "info"; tone: "info" }
    }
    Row {
        spacing: Interface.spaceNear
        KvitChipButton { text: "settled"; tone: "success"; strong: true }
        KvitChipButton {
            text: "disputed"; tone: "danger"; strong: true; symbol: "warning"
        }
        KvitChipButton {
            text: "keyboard focus"
            Component.onCompleted: if (Window.window)
                forceActiveFocus(Qt.TabFocusReason)
        }
    }
}
```

*What pressing it leads to, where the caller knows*

```qml
// No chevron is drawn for a chip that acts. Whether pressing it discloses
// a list under the row, leaves the view or opens another window belongs to
// the destination rather than to the mark, which is why the symbol is the
// caller's to name — the same reason KvitLink draws none.
Row {
    spacing: Interface.spaceNear
    KvitChipButton { text: "3 changes"; symbol: "file" }
    KvitChipButton { text: "1 ahead"; trailingSymbol: "chevron-right" }
    KvitChipButton { text: "2 behind"; trailingSymbol: "chevron-right" }
    KvitChipButton { text: "Open on the hub"; trailingSymbol: "external" }
}
```

*Unavailable, with the reason attached*

```qml
// `enabled: false` would take this out of the tab order and stop its
// hover events, making the one chip with something to explain the one chip
// a reader cannot reach to hear it. A chip with a reason keeps its place,
// still shows the words in its tooltip and its accessible description, and
// emits nothing when it is pressed. The fill going away rather than
// changing hue is what a reader who cannot separate the tones still sees.
Row {
    spacing: Interface.spaceNear
    KvitChipButton { text: "1 ahead" }
    KvitChipButton {
        text: "2 behind"
        unavailableReason: "The other branch has not been fetched yet."
    }
    KvitChipButton {
        text: "in sync"; tone: "success"
        unavailableReason: "There is nothing on either side to compare."
    }
}
```

*Elided in a narrow column, at any interface size*

```qml
// The label gives way and the symbols keep their size: a symbol at half
// width is a smudge, and the trailing one is what says where pressing this
// goes. The whole label is in the tooltip once it no longer fits.
Column {
    spacing: Interface.spaceSnug
    KvitChipButton {
        width: Interface.px(150)
        symbol: "warning"
        text: "A fact whose whole phrase does not fit in this column"
        trailingSymbol: "chevron-right"
    }
    KvitChipButton {
        width: Interface.px(90)
        text: "A fact whose whole phrase does not fit in this column"
    }
}
```

*The one that is already open*

```qml
// A row of chips is a row of places to go, and one of them is where the
// reader already is. Drawn as the current one, it stops reading as
// somewhere to go: the selection tint under it, the accent on its edge,
// and the label in the reader's own colour and weight — the weight
// because the tints are two percent of lightness in the high-contrast
// theme and nothing at all in a grayscale screenshot.
//
// This is not `tone: "accent"`, which says the chip is a different kind
// of fact from the ones beside it rather than the same fact in a
// different state.
Row {
    spacing: Interface.spaceNear
    KvitChipButton { text: "3 changes"; symbol: "diff" }
    KvitChipButton {
        text: "1 ahead"; current: true
        trailingSymbol: "chevron-right"
    }
    KvitChipButton { text: "2 behind"; trailingSymbol: "chevron-right" }
}
```

*Where pressing one goes*

```qml
// "1 ahead" is a count. What it opens is the one commit, and a reader
// who has not pressed one of these before has no way to know that.
// `explanation` is that sentence; `unavailableReason` is the other one,
// about a chip that cannot be pressed at all, and a chip moving between
// the two states does not have to swap one string into the other's
// property.
Row {
    spacing: Interface.spaceNear
    KvitChipButton {
        text: "1 ahead"
        trailingSymbol: "chevron-right"
        explanation: "Opens the one commit this branch has and the "
                     + "remote does not."
        Component.onCompleted: if (Window.window)
            forceActiveFocus(Qt.TabFocusReason)
    }
    KvitChipButton {
        text: "2 behind"
        unavailableReason: "The remote has not been fetched yet."
    }
}
```

### KvitStepper

A number with a minus and a plus beside it, for a small range a reader adjusts by one or two. The range comes from whatever owns it rather than being repeated at the call site.

| Property | Type | |
|---|---|---|
| `value` | int |  |
| `from` | int |  |
| `to` | int |  |
| `step` | int |  |
| `label` | string |  |
| `unit` | string | The unit, drawn after the number in muted colour. |

*The interface-size row, pointed at the right setting*

```qml
KvitStepper {
    label: "Interface size"
    from: Interface.minFontSize
    to: Interface.maxFontSize
    value: Interface.fontSize
    unit: "px"
    onValueModified: v => Interface.fontSize = v
}
```

### KvitField

A single line of text. The outline is the control-boundary token, so its edges are visible. An error is a message and a border together, never a border alone.

| Property | Type | |
|---|---|---|
| `error` | string | A message under the field. |
| `label` | string |  |

*Resting, filled and in error*

```qml
Column {
    spacing: Interface.spaceLoose
    KvitField { label: "Payee"; placeholderText: "Who was paid" }
    KvitField { label: "Reference"; text: "TX-00173404" }
    KvitField { label: "Amount"; text: "twelve"; error: "Not a number" }
    KvitField { label: "Locked"; text: "read only"; enabled: false }
}
```

### KvitTextArea

Several lines of text, typed or read. KvitField's outline and error rule in the `field` form, no ground at all in the `plain` one, which is for a document filling a pane. `underlay` takes items positioned in the text's own coordinates, so a wash over a marked passage does not mean replacing the background.

| Property | Type | |
|---|---|---|
| `form` | string | "field" \| "plain". |
| `error` | string | A message under the box. |
| `label` | string |  |
| `mono` | bool |  |
| `underlay` | alias | Items drawn behind the text, in the text's own coordinates. |

*Resting, filled, monospace and in error*

```qml
Column {
    spacing: Interface.spaceLoose
    KvitTextArea { width: Interface.px(280); label: "Message"; placeholderText: "What happened" }
    KvitTextArea {
        width: Interface.px(280)
        label: "Source"
        mono: true
        readOnly: true
        text: "# Notes\n\nThe file as it stands."
    }
    KvitTextArea { width: Interface.px(280); label: "Summary"; text: "..."; error: "Say what changed" }
}
```

*Plain: a document filling a pane, with no ground of its own*

```qml
KvitTextArea {
    width: Interface.px(280)
    form: "plain"
    mono: true
    readOnly: true
    label: "Source"
    text: "# Notes\n\nA document has no edges to find:\nthe pane is its edge."
}
```

*A wash behind a marked passage, drawn through `underlay`*

```qml
KvitTextArea {
    id: marked
    width: Interface.px(280)
    label: "Reviewed source"
    mono: true
    readOnly: true
    text: "one\ntwo\nthree"
    underlay: Rectangle {
        // `length` is the document’s own character count, so this
        // is measured again whenever the text changes and never
        // asked before the characters are there. A rect taken from
        // the method call alone is measured once and never again.
        readonly property rect box: marked.length > 7
            ? marked.positionToRectangle(4) : Qt.rect(0, 0, 0, 0)
        x: box.x
        y: box.y
        width: Interface.px(40)
        height: box.height
        color: Theme.selectionTint
    }
}
```

### KvitSearchField

A field that filters something. Escape clears rather than reverting, and it announces its result count — filtering is the one interaction whose whole outcome happens somewhere else on the screen.

| Property | Type | |
|---|---|---|
| `matches` | int | How many things the filter left. |
| `matchedNoun` | string |  |
| `matchedNounPlural` | string | The word for several of them. |
| `matchPhrase` | readonly string | "1 transaction", "250,000 transactions" — what a screen reader is told the filter left. |

*Empty and filtering*

```qml
Column {
    spacing: Interface.spaceLoose
    KvitSearchField { width: Interface.px(220) }
    KvitSearchField {
        width: Interface.px(220); text: "harlow"
        matches: 47; matchedNoun: "transaction"
    }
    // What a screen reader is told is a sentence: the digits grouped by the
    // reader's locale and the plural a word, so "250,000 entries" rather than
    // "250000 entry(s)". A noun that does not take an s says its own plural.
    KvitSearchField {
        width: Interface.px(220); text: "2026"
        matches: 250000
        matchedNoun: "entry"; matchedNounPlural: "entries"
    }
}
```

### KvitCheck

A checkbox in three states. The third — partial — is what a parent row shows when some of its children are checked; drawing that as unchecked loses the information and drawing it as checked is a lie.

| Property | Type | |
|---|---|---|
| `partial` | bool | Set this instead of `checked` when the state is "some of them". |

*Off, on, partial and disabled*

```qml
Column {
    spacing: Interface.spaceNear
    KvitCheck { text: "Include archived" }
    KvitCheck { text: "Include drafts"; checked: true }
    KvitCheck { text: "Some of these"; partial: true }
    KvitCheck { text: "Not available"; enabled: false }
}
```

### KvitSelect

A choice from a list too long to lay out. For two to four self-evident options, KvitSegmented is right instead — a three-item dropdown hides two of the three answers for no reason.

| Property | Type | |
|---|---|---|
| `label` | string |  |

*A currency picker*

```qml
KvitSelect {
    label: "Currency"
    model: ["GBP", "EUR", "USD", "JPY", "CHF", "SEK"]
}
```

### KvitTab

One tab in a row of them. The selected tab is marked by an underline as well as by colour and weight — selection shown by colour alone is the most common place the rule gets broken. `explanation` is one sentence saying what the view behind the tab shows, which two or three words cannot.

| Property | Type | |
|---|---|---|
| `selected` | bool |  |
| `count` | int | A count beside the label — how many things this tab holds. |
| `explanation` | string | One sentence saying what this view shows, which two or three words on a tab cannot. |
| `groupedCount` | readonly string | The count with the reader's own digit grouping, said and drawn the same way. |

*Selected, counted and plain*

```qml
Row {
    KvitTab { text: "All"; selected: true; count: 1284 }
    KvitTab { text: "Uncategorised"; count: 47 }
    KvitTab { text: "Disputed" }
}
```

*Each tab saying what its view shows*

```qml
Row {
    KvitTab {
        text: "Changes"; selected: true
        explanation: "What this version changed."
    }
    KvitTab {
        text: "Source"
        explanation: "The file’s own text."
    }
    KvitTab {
        text: "History"
        explanation: "Every version and who wrote it."
    }
}
```

### KvitTooltip

A short label next to a control after a pause. Never the only place a control's meaning lives — a control explained only by its tooltip is unusable on a keyboard, a screen reader and a touch screen.

*On a button*

```qml
Item {
    id: stage
    width: parent.width; height: Interface.px(80)
    KvitButton {
        id: action
        x: Interface.px(60); y: Interface.px(40)
        text: "Reconcile"
        // KvitTooltip rather than the attached `ToolTip.text`: the attached
        // one is the platform style's, which is a yellow box belonging to no
        // theme here.
        KvitTooltip {
            id: hint
            text: "Match these against the statement"
            delay: 0
            Component.onCompleted: if (stage.Window.window) hint.visible = true
        }
    }
}
```

### KvitPopover

A small surface anchored to a control holding something to act on. Takes focus, closes on Escape and on a click outside — which is what separates it from a hover card.

| Property | Type | |
|---|---|---|
| `title` | string |  |

*Open, holding a form*

```qml
Item {
    id: stage
    width: parent.width; height: Interface.px(150)
    KvitPopover {
        id: filter
        parent: stage
        title: "Filter"
        width: Interface.px(240)
        // A Popup opens into a window, and `visible: true` on one with no
        // parent does nothing at all.
        Component.onCompleted: if (stage.Window.window) filter.open()
        Column {
            spacing: Interface.spaceNear
            KvitCheck { text: "Settled"; checked: true }
            KvitCheck { text: "Pending" }
            KvitCheck { text: "Disputed" }
        }
    }
}
```

### KvitHint

An information-icon trigger for an explanation too long for a tooltip. Click or keyboard activation keeps its KvitPopover open for reading, and Escape or an outside click dismisses it.

| Property | Type | |
|---|---|---|
| `label` | string | **required.**  |
| `text` | string | **required.**  |
| `popoverWidth` | int |  |
| `opened` | readonly alias |  |

*Open for a longer explanation*

```qml
Item {
    width: parent.width; height: Interface.px(130)
    KvitHint {
        label: "About automatic matching"
        text: "Automatic matching compares the date, amount and reference. It never changes the imported statement."
        Component.onCompleted: if (Window.window) open()
    }
}
```

### KvitHoverCard

More about the thing under the pointer. Read-only, always: a surface that appears on hover is unreachable by keyboard and by touch, so nothing inside one may be the only route to an action.

| Property | Type | |
|---|---|---|
| `shown` | bool |  |
| `padding` | int |  |

*Showing what a row could not fit*

```qml
KvitHoverCard {
    shown: true
    Column {
        spacing: Interface.spaceTight
        KvitLabel { text: "Payment to Harlow depot"; role: "strong" }
        KvitLabel { text: "17 March 2026 at 09:14"; role: "small"
                    color: Theme.textMuted }
        KvitFigure { value: "1,284.50"; unit: "GBP" }
    }
}
```

### KvitToast

A message that appears, says one thing and goes away. For confirming what already happened. A toast with an action stays until dismissed — an undo that times out mid-read is an undo the reader cannot use.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `tone` | string | "info" \| "success" \| "warning" \| "danger" |
| `action` | string |  |
| `shown` | bool |  |
| `timeout` | int | How long it stays, in milliseconds. |
| `toneColor` | readonly color |  |

*Four tones, and one with an undo*

```qml
Column {
    spacing: Interface.space
    KvitToast { shown: true; text: "Copied to the clipboard" }
    KvitToast { shown: true; tone: "success"; text: "Statement imported" }
    KvitToast { shown: true; tone: "warning"; text: "Two rows could not be matched" }
    KvitToast { shown: true; tone: "danger"; text: "The file could not be read" }
    KvitToast { shown: true; tone: "success"; text: "40 transactions archived"
                action: "Undo" }
}
```

### KvitNotice

A message that stays: a condition that is still true, where a toast is an acknowledgement of something that finished. Dismissible only when dismissing it is meaningful.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `detail` | string |  |
| `tone` | string | "info" \| "warning" \| "danger" \| "success" |
| `action` | string |  |
| `dismissible` | bool |  |
| `toneColor` | readonly color |  |

*A warning that can be dismissed, and an error that cannot*

```qml
Column {
    width: parent.width
    spacing: Interface.space
    KvitNotice {
        width: parent.width; tone: "warning"
        text: "Your licence expires in three days"
        detail: "Renew before 30 March to keep syncing."
        action: "Renew"; dismissible: true
    }
    KvitNotice {
        width: parent.width; tone: "danger"
        text: "This vault could not be saved"
        detail: "The disk is full. Nothing has been lost; the changes are still held."
    }
}
```

### KvitDialog

A modal surface that has to be answered first. Expensive, and worth it only where the answer really does come first; the confirming button is on the right, and a destructive dialog takes no keyboard default. There are three ways out, because different readers find different ones: the close symbol at the right of the title, the cancel button, and Escape.

| Property | Type | |
|---|---|---|
| `detail` | string |  |
| `confirmText` | string | Words rather than "OK": a button that says what it does can be read without reading the sentence above it, which is what a reader in a hurry actually does. |
| `cancelText` | string |  |
| `destructive` | bool |  |

*A destructive confirmation*

```qml
// Shown here without its modality so it sits on the page. A real one is
// modal and centres itself on the window, which is what a dialog that has
// to be answered first should do and is not something a gallery card can
// contain.
Item {
    id: stage
    width: parent.width; height: Interface.px(210)
    KvitDialog {
        id: confirm
        parent: stage
        modal: false
        anchors.centerIn: stage
        title: "Delete 40 transactions?"
        detail: "They will be removed from every report. This cannot be undone."
        confirmText: "Delete them"
        destructive: true
        Component.onCompleted: if (stage.Window.window) confirm.open()
    }
}
```

### KvitBar

One quantity against a stated scale. A bar with no value draws a tick rather than a zero-width fill; a bounded figure is hatched, and the hatch survives grayscale; and `maximum` is required, because a bar drawn against its own list rescales invisibly.

| Property | Type | |
|---|---|---|
| `value` | real |  |
| `maximum` | real | **required.**  |
| `measured` | bool |  |
| `bounded` | bool | The value is an upper bound rather than a measurement. |
| `color` | color |  |
| `wide` | bool |  |
| `label` | string | What this bar is, for a screen reader. |
| `unit` | string |  |
| `fraction` | readonly real |  |

*Measured, bounded and unmeasured*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceLoose
    KvitBar { width: parent.width; value: 62; maximum: 100
              label: "Attention"; unit: "h" }
    KvitBar { width: parent.width; value: 38; maximum: 100; wide: true
              bounded: true; color: Theme.axisAgent; label: "Agent"; unit: "h" }
    KvitBar { width: parent.width; maximum: 100; measured: false
              label: "Unrecorded" }
}
```

### KvitStackedBar

Several quantities adding to one total. The two-pixel gap is the reason this is a component: adjacent fills that touch read as one fill with a colour change in it.

| Property | Type | |
|---|---|---|
| `segments` | var | A list of { value, color, label } objects. |
| `maximum` | real | The scale. |
| `wide` | bool |  |
| `label` | string |  |
| `total` | readonly real |  |
| `scale` | readonly real |  |

*A four-way breakdown*

```qml
KvitStackedBar {
    width: parent.width
    wide: true
    label: "Where the month went"
    segments: [
        { value: 420, label: "rent", color: Theme.categorical(0) },
        { value: 180, label: "groceries", color: Theme.categorical(1) },
        { value: 95, label: "transport", color: Theme.categorical(2) },
        { value: 60, label: "everything else", color: Theme.categorical(3) }
    ]
}
```

### KvitSpark

The shape of a series, small enough to sit in a row. A period that was never measured is a baseline tick rather than a zero-height bar: a missing week and a quiet week are different facts.

| Property | Type | |
|---|---|---|
| `values` | var | Numbers, with null for a period that was not measured. |
| `color` | color |  |
| `label` | string |  |
| `maximum` | real | The scale. |
| `scale` | readonly real |  |

*With a hole in the middle*

```qml
KvitSpark {
    width: Interface.px(160)
    label: "Notes written"
    values: [3, 5, 8, 6, null, null, 9, 12, 7, 4, 6, 11]
}
```

### KvitTrend

A series with a value axis and a hover crosshair — read for values, where a spark is read for shape. A gap in the data draws as a gap: interpolating over a hole asserts values nobody measured. A second series is dashed as well as differently coloured, and both are named in the key and in the crosshair.

| Property | Type | |
|---|---|---|
| `points` | var | A list of { x, y } points, y null where the period was not measured. |
| `minimumY` | real |  |
| `maximumY` | real |  |
| `color` | color |  |
| `label` | string |  |
| `unit` | string |  |
| `gridlines` | int | How many horizontal gridlines to draw, including the two extremes. |
| `secondPoints` | var | A second series over the same periods, drawn on the same axis: assets against liabilities, spending against income, this year against last. |
| `secondColor` | color | The second series is dashed as well as differently coloured. |
| `secondLabel` | string | What each line is. |
| `hasSecond` | readonly bool |  |
| `periods` | readonly int | How many periods the x axis spans. |
| `hovered` | readonly int | Where the pointer is, as an index into `points`, or -1. |
| `axisWidth` | readonly int | The plot rectangle, published so a caller can draw over it — an annotation marking where an account's history begins, a band behind a period, a threshold rule. |
| `plotWidth` | readonly real |  |
| `plotTop` | readonly real |  |
| `plotHeight` | readonly real |  |

*A series, and the empty state*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceLoose
    KvitTrend {
        width: parent.width
        label: "Balance"; unit: "GBP"
        minimumY: 0; maximumY: 2000
        points: [{ y: 400 }, { y: 620 }, { y: 580 }, { y: 900 },
                 { y: null }, { y: 1400 }, { y: 1250 }, { y: 1700 }]
    }
    KvitTrend { width: parent.width; height: Interface.px(90); label: "Savings" }
}
```

*Two series, with an annotation drawn over the plot*

```qml
KvitTrend {
    id: worth
    width: parent.width
    height: Interface.px(150)
    label: "Assets"; secondLabel: "Liabilities"; unit: "GBP"
    minimumY: 0; maximumY: 2000
    points: [{ y: 400 }, { y: 620 }, { y: 580 }, { y: 900 },
             { y: 1150 }, { y: 1400 }, { y: 1250 }, { y: 1700 }]
    secondPoints: [{ y: null }, { y: null }, { y: 300 }, { y: 340 },
                   { y: 320 }, { y: 290 }, { y: 260 }, { y: 240 }]

    // Where the second account's history begins. The plot rectangle is
    // published, so an overlay lines up with the same arithmetic.
    Rectangle {
        x: worth.axisWidth + worth.plotWidth * 2 / 7
        y: worth.plotTop
        width: Interface.hairline
        height: worth.plotHeight
        color: Theme.marker
    }
    KvitLabel {
        x: worth.axisWidth + worth.plotWidth * 2 / 7 + Interface.spaceSnug
        y: worth.plotTop + worth.plotHeight - height
        text: "history starts here"
        role: "caption"
        color: Theme.textFaint
    }
}
```

### KvitDistribution

How a set of values is spread. An average of four days and an average made of one twenty-day outlier are the same number and different situations.

| Property | Type | |
|---|---|---|
| `minimum` | real |  |
| `lowerQuartile` | real |  |
| `median` | real |  |
| `upperQuartile` | real |  |
| `maximum` | real |  |
| `scaleMinimum` | real | The scale the box is drawn against, which is usually wider than this one distribution so several rows are comparable. |
| `scaleMaximum` | real |  |
| `measured` | bool |  |
| `color` | color |  |
| `label` | string |  |
| `unit` | string |  |

*Two rows on one scale*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceLoose
    KvitDistribution {
        width: parent.width; label: "Time to settle, card"
        scaleMinimum: 0; scaleMaximum: 30; unit: "d"
        minimum: 1; lowerQuartile: 2; median: 3; upperQuartile: 5; maximum: 21
    }
    KvitDistribution {
        width: parent.width; label: "Time to settle, transfer"
        scaleMinimum: 0; scaleMaximum: 30; unit: "d"
        minimum: 1; lowerQuartile: 1; median: 2; upperQuartile: 2; maximum: 4
    }
}
```

### KvitGauge

How much of an allowance is used, with the target marked and a pace mark saying where an even rate would have reached. Sixty percent spent is fine on day eighteen and a problem on day six.

| Property | Type | |
|---|---|---|
| `value` | real |  |
| `allowance` | real | **required.**  |
| `pace` | real | Where an even rate would have reached by now, 0 to 1. |
| `measured` | bool |  |
| `label` | string |  |
| `unit` | string |  |
| `fraction` | readonly real |  |
| `over` | readonly bool |  |
| `toneColor` | readonly color | Behind pace, at pace, or ahead of it. |

*Ahead of pace, behind it, and over*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceLoose
    KvitGauge { width: parent.width; value: 240; allowance: 600; pace: 0.6
                label: "Groceries"; unit: "GBP" }
    KvitGauge { width: parent.width; value: 480; allowance: 600; pace: 0.6
                label: "Leisure"; unit: "GBP" }
    KvitGauge { width: parent.width; value: 720; allowance: 600; pace: 0.6
                label: "Transport"; unit: "GBP" }
}
```

### KvitDelta

How much something changed and in which direction, with an arrow as well as a colour. Whether up is good is the caller's to say: a rise in spending and a rise in savings are the same arrow and opposite colours.

| Property | Type | |
|---|---|---|
| `change` | real |  |
| `unit` | string |  |
| `measured` | bool |  |
| `goodDirection` | string | "up" \| "down" \| "none". |
| `precision` | int | Round to this many decimal places when drawing. |
| `rising` | readonly bool |  |
| `flat` | readonly bool |  |
| `toneColor` | readonly color |  |

*Up, down and unchanged*

```qml
Row {
    spacing: Interface.spaceLoose
    KvitDelta { change: 12.4; unit: "%"; precision: 1; goodDirection: "up" }
    KvitDelta { change: -3.2; unit: "%"; precision: 1; goodDirection: "up" }
    KvitDelta { change: 18; unit: "GBP"; goodDirection: "down" }
    KvitDelta { change: 0 }
    KvitDelta { measured: false }
}
```

### KvitStatTile

A card carrying one figure, what it is, how it changed and its recent shape — the shape of six of kvit-cash's seven dashboard widgets. The order is fixed so a row of tiles can be scanned one part at a time.

| Property | Type | |
|---|---|---|
| `label` | string |  |
| `value` | string |  |
| `unit` | string |  |
| `measured` | bool |  |
| `change` | real |  |
| `hasChange` | bool |  |
| `goodDirection` | string | "up" \| "down" \| "none": which direction of change is the good one here. |
| `history` | var | The recent shape, with null for a period that was not measured. |
| `caption` | string |  |

*A dashboard row*

```qml
Row {
    spacing: Interface.columnGap
    KvitStatTile {
        label: "Balance"; value: "4,182.30"; unit: "GBP"
        hasChange: true; change: 240.10; goodDirection: "up"
        history: [3200, 3400, 3390, 3800, 4000, 4182]
    }
    KvitStatTile {
        label: "Spent this month"; value: "812.40"; unit: "GBP"
        hasChange: true; change: 96.20; goodDirection: "down"
        caption: "12 days remaining"
    }
    KvitStatTile { label: "Uncategorised"; measured: false }
}
```

### KvitFigureBlock

A figure with its name under it: one number a reader is meant to take away. The number is above and larger, because a row of these is read across the numbers.

| Property | Type | |
|---|---|---|
| `value` | string |  |
| `unit` | string |  |
| `label` | string |  |
| `measured` | bool |  |
| `bounded` | bool |  |
| `color` | color |  |
| `role` | string | "display" \| "headline" \| "title" — how much of the screen this figure is meant to command. |

*A row of three*

```qml
Row {
    spacing: Interface.px(40)
    KvitFigureBlock { label: "Transactions"; value: "1,284" }
    KvitFigureBlock { label: "Categorised"; value: "97"; unit: "%" }
    KvitFigureBlock { label: "Reconciled"; measured: false }
}
```

### KvitCell

One cell of a table, drawn according to what kind of value its column holds. The kind comes from the column, so every cell in it aligns the same way and says the same thing about a missing value. A cell formats nothing: Money and Figure are given the string they draw.

| Property | Type | |
|---|---|---|
| `kind` | string | "Text" \| "Figure" \| "Chip" \| "Slug" \| "Date" \| "Money" \| "Marks" \| "Check" |
| `value` | var |  |
| `measured` | bool |  |
| `mark` | string | For a Chip cell: which tone. |
| `unit` | string |  |
| `fullValue` | string | The whole value where `value` is the shortened one a narrow column has room for, which a model supplies through TableModelBase's FullTextRole under the role name `fullText`. |
| `selected` | bool |  |
| `current` | bool | Whether the keyboard cursor is on this cell. |
| `marks` | var | For a Marks cell: every state this row is in, as { tone, shape, label } objects, one dot each. |
| `checked` | bool | For a Check cell: whether this row is picked. |
| `numeric` | readonly bool |  |
| `valueText` | readonly string | The value as the string it is drawn as, with nothing in place of a value that is not there. |
| `negative` | readonly bool | Whether a money value is money leaving, read from the sign already written into it. |

*The six kinds a column of values can be*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceTight
    KvitCell { width: parent.width; kind: "Text"; value: "Payment to Harlow depot" }
    KvitCell { width: parent.width; kind: "Date"; value: new Date(2026, 7, 30) }
    KvitCell { width: parent.width; kind: "Slug"; value: "TX-00173404" }
    KvitCell { width: parent.width; kind: "Figure"; value: "1284.5"; unit: "GBP" }
    // Money takes an amount that has already been formatted, exactly as
    // KvitBeforeAfter takes its `before` and `after`. The cell does no
    // arithmetic: an amount is a signed count of the minor unit of its own
    // currency, and both halves of that are lost the moment a cell turns one
    // into a JavaScript number and prints two decimal places after it.
    KvitCell { width: parent.width; kind: "Money"; value: "-42.90"; unit: "GBP" }
    KvitCell { width: parent.width; kind: "Chip"; value: "Settled"; mark: "success" }
    KvitCell { width: parent.width; kind: "Figure"; measured: false }
}
```

*Several states at once, and a row the reader picks*

```qml
Column {
    id: sample
    width: parent.width
    spacing: Interface.spaceTight

    property bool picked: true

    // Every state a row is in at once, one dot each. A transaction is
    // regularly settled and unreviewed together, and one chip has to drop one
    // of the two. Each mark carries a tone, the shape that says the same
    // thing without colour, and the word that is both the tooltip and what a
    // screen reader says.
    KvitCell {
        width: parent.width
        kind: "Marks"
        marks: [
            { "tone": "success", "shape": "circle", "label": "Settled" },
            { "tone": "warning", "shape": "diamond", "label": "Not reviewed" },
            { "tone": "neutral", "shape": "square", "label": "Has an attachment" }
        ]
    }

    // The box does not tick itself. Whatever owns the selection does, and the
    // cell redraws from it, so a request that is refused leaves no tick
    // behind.
    KvitCell {
        width: parent.width
        kind: "Check"
        checked: sample.picked
        onToggled: wanted => sample.picked = wanted
    }
}
```

*A value the column cut short*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceTight

    // A value wider than the column it is in. The model answers the whole of
    // it through TableModelBase's FullTextRole, under the role name
    // `fullText`, and the cell discloses it under the pointer and under the
    // keyboard cursor. `current` is how the view says where that cursor is,
    // because a cell is not a control and takes no focus of its own.
    KvitCell {
        width: Interface.px(150)
        kind: "Text"
        value: "Harlow depot retainer"
        fullValue: "Harlow depot — quarterly maintenance retainer"
    }
}
```

### KvitTable

A dense, configurable, editable table over a C++ model. Columns sort, resize and open a menu from their own header, by keyboard as well as pointer; a row opens its record on one press; and a Check column draws a box per row with a select-displayed control in its header. The model says what each column is through `columnKindName()` and answers the roles KvitUi.TableModelBase names: `marks` for a Marks column, `checked` for a Check column, `unit` for the currency or unit a Money or Figure cell draws beside its value, and `fullText` for the whole of a value the column cut short, which the cell then discloses on hover and under the keyboard cursor. Holds smooth scrolling and sub-100 ms filtering at 250,000 rows with twelve columns of value, which is the measurement prd.md Decision 3 turns on.

| Property | Type | |
|---|---|---|
| `model` | var | **required.** A KvitUi.TableModelBase subclass. |
| `editable` | bool |  |
| `emptyTitle` | string |  |
| `emptyDetail` | string |  |
| `hiddenColumns` | var | Column indices the reader has hidden. |
| `sortColumn` | int | Which column the rows are ordered by, and in which direction. |
| `sortAscending` | bool |  |
| `headerChecked` | bool | The select-displayed control that a `Check` column draws in its header, and its third state for "some of the rows shown". |
| `headerPartial` | bool |  |
| `selection` | readonly alias |  |
| `view` | readonly alias |  |
| `rowCount` | readonly int | How many rows the table is showing, which after a filter is fewer than the model holds. |

*Two hundred and fifty thousand rows*

```qml
Item {
    width: parent.width; height: Interface.px(260)
    BenchmarkTableModel { id: rows }
    KvitTable { anchors.fill: parent; model: rows }
}
```

*Sorting, a column menu and a column of boxes*

```qml
Item {
    width: parent.width; height: Interface.px(260)
    BenchmarkTableModel { id: rows; totalRows: 400 }
    KvitTable {
        anchors.fill: parent
        model: rows
        // Six of the thirteen columns, so the boxes at the right-hand end are
        // on screen without scrolling.
        hiddenColumns: [1, 2, 5, 8, 10, 11]
        // The table draws the indicator and says what was asked for; putting
        // the rows in that order is the application's, because the model is.
        // Nothing here sorts, so the arrow stays where it was set.
        sortColumn: 0
        sortAscending: false
        // The column of boxes, and the select-displayed control in its
        // header. The selection is the model's and is keyed by the underlying
        // row, so narrowing the filter does not lose it.
        headerChecked: rows.allShownChecked
        headerPartial: rows.someShownChecked
        onHeaderToggled: wanted => rows.setEveryShownChecked(wanted)
        onCellToggled: (row, column, wanted) => rows.setChecked(row, wanted)
    }
}
```

*Nothing to show*

```qml
Item {
    width: parent.width; height: Interface.px(200)
    BenchmarkTableModel { id: rows; totalRows: 0 }
    KvitTable {
        anchors.fill: parent
        model: rows
        emptyTitle: qsTr("No transactions")
        emptyDetail: qsTr("Nothing matches the current filter.")
    }
}
```

### KvitScrollBar

A scroll bar, occupying its own strip rather than floating over content. Thirty-two files in the estate have a private version; an overlay bar hides the right-hand column of a table and the last character of every elided label.

| Property | Type | |
|---|---|---|
| `flickable` | Flickable | **required.**  |
| `orientation` | int |  |
| `isVertical` | readonly bool |  |
| `span` | readonly real | Guarded even though `flickable` is required, because a scroll bar outlives its flickable by one event loop turn when a view is torn down, and an unguarded binding re-evaluates in that gap and warns. |
| `position` | readonly real |  |

*Beside a scrolling column*

```qml
KvitRegion {
    width: 480; height: Interface.px(120)
    Column {
        id: rows
        width: 460
        Repeater { model: 10; KvitSlimRow { width: rows.width
                   name: "Row " + (index + 1) } }
    }
}
```

### KvitMenu

A list of commands. Eighteen private versions in the estate, and what they mostly get wrong is the same two things: no keyboard route in, and no separator before the destructive item.

*Opened, with shortcuts and a destructive item*

```qml
Item {
    id: stage
    width: parent.width; height: Interface.px(170)
    KvitMenu {
        id: menu
        // Only once the item is in a window: popping a menu on an item
        // with no window is a crash rather than a no-op.
        Component.onCompleted: if (stage.Window.window) menu.popup(stage, 0, 0)
        KvitMenuItem { text: "Open"; symbol: "file"; shortcut: "Ctrl+O" }
        KvitMenuItem { text: "Duplicate"; symbol: "copy"; shortcut: "Ctrl+D" }
        KvitMenuItem { text: "Archive"; symbol: "archive" }
        MenuSeparator {}
        KvitMenuItem { text: "Delete"; symbol: "trash"; destructive: true }
    }
}
```

### KvitMenuItem

One line of a menu, carrying its shortcut on the right — which is how a menu teaches a faster route to a reader who keeps using it.

| Property | Type | |
|---|---|---|
| `symbol` | string |  |
| `shortcut` | string |  |
| `destructive` | bool |  |
| `explanation` | string | One sentence saying what choosing this does, for an entry whose words cannot say it: that applying a staged change writes it into the note as one undo step, that discarding it tells the agent the change was rejected. |

*Ordinary, disabled and destructive*

```qml
Item {
    id: stage
    width: parent.width; height: Interface.px(110)
    KvitMenu {
        id: menu
        // Only once the item is in a window: popping a menu on an item
        // with no window is a crash rather than a no-op.
        Component.onCompleted: if (stage.Window.window) menu.popup(stage, 0, 0)
        KvitMenuItem { text: "Reconcile"; symbol: "check"; shortcut: "Ctrl+R" }
        KvitMenuItem { text: "Split"; symbol: "split"; enabled: false }
        KvitMenuItem { text: "Delete"; symbol: "trash"; destructive: true }
    }
}
```

*The chosen one of a set, and an entry with a sentence*

```qml
// The chosen entry of a set is drawn as chosen by the component: a tick
// in the same gutter a symbol would use, the label in bold beside it, and
// the checked state announced rather than spelled into the label. An
// entry whose words cannot say what choosing it does takes a sentence.
Item {
    id: stage
    width: parent.width; height: Interface.px(150)
    KvitMenu {
        id: menu
        // Only once the item is in a window: popping a menu on an item
        // with no window is a crash rather than a no-op.
        Component.onCompleted: if (stage.Window.window) menu.popup(stage, 0, 0)
        KvitMenuItem {
            text: "Apply"
            symbol: "check"
            explanation: "Writes the change into the note as one undo step."
        }
        MenuSeparator {}
        KvitMenuItem { text: "Fast"; checkable: true; checked: true }
        KvitMenuItem { text: "Careful"; checkable: true }
    }
}
```

### KvitTree

A nested list the reader can open and close. Twelve private versions; what a shared one has to get right is the keyboard, because depth without Left and Right is a wall.

| Property | Type | |
|---|---|---|
| `label` | string | What a screen reader calls this tree. |
| `nodes` | var | A hierarchy written out here, for a tree whose shape is fixed: account groups, a category hierarchy, the sections of a settings page. |

*A small hierarchy*

```qml
Item {
    width: parent.width; height: Interface.px(250)
    KvitTree {
        anchors.fill: parent
        label: "Accounts"
        nodes: [
            { label: "Everyday", children: [
                "Checking", "Joint checking", "Cash"] },
            { label: "Savings", children: [
                "Emergency fund",
                { label: "Certificates", children: ["18 months", "3 years"] }] },
            "Credit card"
        ]
        Component.onCompleted: expandRecursively(-1, 1)
    }
}
```

### KvitSwitch

An option that takes effect the moment it moves — where a checkbox is a value in a form that takes effect on submit. The knob moves and the track fills, so the state is not resting on hue.

*On, off and disabled*

```qml
Column {
    spacing: Interface.spaceNear
    KvitSwitch { text: "Follow the system theme"; checked: true }
    KvitSwitch { text: "Reduce motion" }
    KvitSwitch { text: "Sync over cellular"; enabled: false }
}
```

### KvitRadioGroup

One choice from a handful where the choice needs explaining. Arrow keys move within the group and Tab leaves it, which a column of separate controls does not do.

| Property | Type | |
|---|---|---|
| `options` | var | A list of { value, label, detail, objectName } objects. |
| `current` | var |  |
| `label` | string |  |

*Three options with detail*

```qml
KvitRadioGroup {
    width: parent.width
    label: "Appearance"
    current: "system"
    options: [
        { value: "system", label: "Follow the system",
          detail: "Light or dark, whichever the desktop is set to." },
        { value: "light", label: "Always light" },
        { value: "dark", label: "Always dark" }
    ]
}
```

*An option a test can find by name*

```qml
KvitRadioGroup {
    width: parent.width
    label: "Sessions run in"
    current: "worktree"
    options: [
        { value: "worktree", label: "Separate working copies",
          objectName: "copyModeWorktree",
          detail: "Its own checkout per session. Works on any disk." },
        { value: "project", label: "The project folder itself",
          objectName: "copyModeProject",
          detail: "No copy; one writer at a time, and never confined." }
    ]
}
```

### KvitProgress

How far through something the application is. An unknown total is a moving band rather than a bar creeping toward the end, and it carries a label saying how far through *what*.

| Property | Type | |
|---|---|---|
| `value` | real |  |
| `maximum` | real |  |
| `determinate` | bool |  |
| `label` | string |  |
| `showPercent` | bool |  |
| `fraction` | readonly real |  |

*Determinate and indeterminate*

```qml
Column {
    width: parent.width
    spacing: Interface.spaceLoose
    KvitProgress { width: parent.width; label: "Importing statement"
                   value: 0.47 }
    KvitProgress { width: parent.width; label: "Reconciling"
                   determinate: false }
}
```

### KvitSlider

A value chosen by dragging, for a feel rather than a number. Where the exact value matters, KvitStepper or KvitNumberField are right — hitting a particular number on a slider is hard and it cannot be typed.

| Property | Type | |
|---|---|---|
| `label` | string |  |
| `unit` | string |  |
| `precision` | int |  |

*With its value shown*

```qml
KvitSlider { width: Interface.px(200); label: "Opacity"; value: 0.6
             unit: "%"; from: 0; to: 1 }
```

### KvitSplitView

Two regions the reader can resize. The handle is a wide invisible strip with a hairline down the middle, so the target is comfortable and the rule is still thin; it also moves with the arrow keys.

*Two panes*

```qml
KvitSplitView {
    width: parent.width; height: Interface.px(120)
    KvitPanel { SplitView.preferredWidth: Interface.px(160)
                KvitLabel { anchors.centerIn: parent; text: "left" } }
    KvitPanel { SplitView.fillWidth: true
                KvitLabel { anchors.centerIn: parent; text: "right" } }
}
```

### KvitSegmented

One choice from two to five short options, all visible — the shape kvit-cash's dashboard period control needs, where a control governing seven widgets should say what they are showing without being opened.

| Property | Type | |
|---|---|---|
| `options` | var | A list of strings, or of { value, label } objects. |
| `current` | var |  |
| `label` | string |  |

*A period control*

```qml
Column {
    spacing: Interface.space
    KvitSegmented {
        label: "Period"; current: "month"
        options: [{ value: "week", label: "Week" },
                  { value: "month", label: "Month" },
                  { value: "quarter", label: "Quarter" },
                  { value: "year", label: "Year" }]
    }
    KvitSegmented { options: ["All", "Mine"]; current: "All" }
}
```

### KvitTypeAhead

A field offering matches as the reader types, for a list too long to read. `allowNew` has no default: a tag picker should let the reader invent one, a category picker should not.

| Property | Type | |
|---|---|---|
| `source` | var | Everything that could be matched: strings, or { value, label } objects. |
| `text` | string |  |
| `label` | string |  |
| `placeholder` | string |  |
| `allowNew` | bool | **required.**  |
| `maximumSuggestions` | int |  |
| `matches` | readonly var |  |

*A category picker that will not invent categories*

```qml
KvitTypeAhead {
    width: Interface.px(240)
    label: "Category"
    placeholder: "Start typing"
    allowNew: false
    source: ["Groceries", "Transport", "Utilities", "Rent",
             "Leisure", "Health", "Savings", "Income"]
}
```

### KvitConfirmInPlace

A strip saying what just happened, with the undo inside it. Cheaper than a dialog before every action when the reader is doing the same thing forty times — and only honest when the action really is reversible.

| Property | Type | |
|---|---|---|
| `text` | string |  |
| `undoText` | string | The words on the undo. |
| `shown` | bool |  |
| `affected` | int | What was affected, so the reader can check the count before undoing. |
| `affectedPhrase` | readonly string | "1 item", "250,000 items" — built once and used by both the drawn label and the announcement, so the two cannot disagree about the same number. |

*After a bulk edit*

```qml
KvitConfirmInPlace {
    width: parent.width
    shown: true
    text: "Recategorised as Groceries"
    affected: 40
}
```

### KvitTimeline

What happened to something, newest first, with who did it. In an estate where an agent and a person change the same things, who is the column that makes a history worth reading.

| Property | Type | |
|---|---|---|
| `entries` | var | A list of { when, what, who, detail, tone } objects, newest first. |
| `label` | string |  |

*An account's recent history*

```qml
KvitTimeline {
    width: parent.width
    label: "Account history"
    entries: [
        { when: "14:02", what: "Statement imported", who: "agent",
          detail: "412 transactions, 8 unmatched", tone: "success" },
        { when: "11:20", what: "Two rows disputed", who: "you",
          tone: "warning" },
        { when: "Yesterday", what: "Account opened", who: "you" }
    ]
}
```

### KvitNumberField

A number the reader types, right-aligned in tabular numerals. It validates and says why rather than refusing keystrokes — a field that ignores a key gives no reason, and the usual cause is a decimal separator the reader's locale writes differently.

| Property | Type | |
|---|---|---|
| `minimum` | real |  |
| `maximum` | real |  |
| `decimals` | int | 0 for an integer field. |
| `value` | readonly real | The parsed value, or NaN while the text is not a number. |
| `valid` | readonly bool |  |
| `normalised` | readonly string | The text with the locale's decimal separator turned into a point, so "3,50" from a German keyboard parses the same as "3.50". |

*Integer and decimal, valid and out of range*

```qml
Column {
    spacing: Interface.spaceLoose
    KvitNumberField { label: "Days"; text: "14"; minimum: 1; maximum: 365 }
    KvitNumberField { label: "Rate"; text: "4.25"; decimals: 2 }
    KvitNumberField { label: "Days"; text: "999"; minimum: 1; maximum: 365 }
}
```

### KvitMoneyField

An amount of money, in minor units. The reader types 12.34 and the field reports 1234, an integer — money in floating point drifts by a penny somewhere nobody can find. The decimal count comes from the currency.

| Property | Type | |
|---|---|---|
| `currency` | string | An ISO 4217 code. |
| `minorDigits` | int | How many decimal places this currency has. |
| `minorUnits` | readonly real | The amount as an integer number of minor units — pence, cents, sen. |

*Sterling, yen and dinar*

```qml
Column {
    spacing: Interface.spaceLoose
    KvitMoneyField { label: "Amount"; text: "1284.50"; currency: "GBP" }
    KvitMoneyField { label: "Amount"; text: "4200"; currency: "JPY"
                     minorDigits: 0 }
    KvitMoneyField { label: "Amount"; text: "18.750"; currency: "BHD"
                     minorDigits: 3 }
}
```

### KvitDualList

Two lists with items moving between them: what is available, and what is chosen and in what order. Everything works from the keyboard, which most implementations of this shape do not.

| Property | Type | |
|---|---|---|
| `available` | var | Strings, or { value, label } objects. |
| `chosen` | var |  |
| `availableLabel` | string |  |
| `chosenLabel` | string |  |

*Choosing table columns*

```qml
KvitDualList {
    width: parent.width; height: Interface.px(220)
    available: ["Payee", "Account", "Tags", "Note"]
    chosen: ["Date", "Description", "Amount", "Balance"]
}
```

### KvitSpotlight

Darken everything except one region and say something about it — the primitive under a guided tour, and deliberately only the primitive: what drives the stepping differs between the two applications that want one.

| Property | Type | |
|---|---|---|
| `target` | Item | The item to leave lit, in this item's coordinates. |
| `title` | string |  |
| `detail` | string |  |
| `shown` | bool |  |
| `padding` | int |  |
| `hole` | readonly rect |  |

*Focusing a button*

```qml
Item {
    width: parent.width; height: Interface.px(160)
    KvitButton { id: target; text: "Import a statement"; form: "primary"
                 x: Interface.px(40); y: Interface.px(20) }
    KvitSpotlight {
        shown: true
        target: target
        title: "Start here"
        detail: "Import a statement and the dashboard fills itself in."
    }
}
```

