# The kvit-ui component catalogue

GENERATED FILE — do not edit. Written by `kvit-ui-gallery --catalog`
from the same list the gallery draws and the library's own source, so what is
described here is what the repository contains.

Every component is in the `kvitui` package: one
`import kvitui "github.com/kvit-s/kvit-ui"` reaches all of them, the `UI` that
holds the theme, the interface size and the fonts, and every colour and size
by name.

Each entry gives what the component is for, how to make one, its fields and
methods, and a working sample. The samples are the functions the gallery
draws, compiled with it and run by its tests, so one that does not work stops
the build.

## What there is

76 components, grouped by what they are for.

**Foundation** — `KvitLabel`, `KvitIcon`, `KvitIconButton`, `KvitLink`
**Structure** — `KvitHeader`, `KvitSidebar`, `KvitSidebarItem`, `KvitBreadcrumb`, `KvitRegion`, `KvitViewHead`, `KvitStatusBar`, `KvitWindow`
**Content** — `KvitSectionHeading`, `KvitRow`, `KvitSlimRow`, `KvitCard`, `KvitPanel`, `KvitPane`, `KvitDivider`, `KvitDisclosure`, `KvitEmptyState`
**Marks** — `KvitChip`, `KvitTag`, `KvitBadge`, `KvitSlug`, `KvitDot`, `KvitSignal`, `KvitPip`
**Quantities** — `KvitFigure`, `KvitBeforeAfter`
**Controls** — `KvitButton`, `KvitChipButton`, `KvitStepper`, `KvitField`, `KvitTextArea`, `KvitSearchField`, `KvitCheck`, `KvitSelect`, `KvitTab`
**Feedback** — `KvitTooltip`, `KvitPopover`, `KvitHint`, `KvitHoverCard`, `KvitToast`, `KvitNotice`, `KvitDialog`, `KvitFloatingView`
**Data** — `KvitBar`, `KvitStackedBar`, `KvitSpark`, `KvitNetFlow`, `KvitTrend`, `KvitDistribution`, `KvitGauge`, `KvitDelta`, `KvitStatTile`, `KvitFigureBlock`, `KvitCell`, `KvitTable`
**Flow** — `KvitScrollBar`, `KvitMenu`, `KvitMenuItem`, `KvitTree`, `KvitSwitch`, `KvitRadioGroup`, `KvitProgress`, `KvitSlider`, `KvitSplitView`, `KvitSegmented`, `KvitTypeAhead`, `KvitConfirmInPlace`, `KvitTimeline`, `KvitNumberField`, `KvitMoneyField`, `KvitDualList`, `KvitSpotlight`

## Symbols

`NewIcon` and every component that takes a `Symbol` accept these names.
They say what a symbol means rather than what it looks like, so the drawing
can change without touching a call site. A name outside this list draws a
marked placeholder and logs a warning, which fails the gallery's tests.

`plus`, `messages-square`, `message-square`, `send`, `chevron-right`, `chevron-down`, `zoom-in`, `zoom-out`, `pencil`, `check`, `archive`, `rotate-ccw`, `chevron-left`, `chevron-up`, `close`, `sidebar`, `list`, `search`, `note`, `tag`, `caret-sort`, `sort-ascending`, `sort-descending`, `dot`, `dot-outline`, `circle`, `square`, `minus`, `info`, `warning`, `error`, `success`, `question`, `settings`, `filter`, `columns`, `calendar`, `clock`, `user`, `robot`, `folder`, `file`, `link`, `external`, `copy`, `trash`, `undo`, `redo`, `refresh`, `more`, `more-vertical`, `drag`, `pin`, `star`, `eye`, `eye-off`, `lock`, `arrow-up`, `arrow-down`, `arrow-left`, `arrow-right`, `trend-up`, `trend-down`, `chart`, `wallet`, `coins`, `bank`, `receipt`, `repeat`, `split`, `merge`, `play`, `pause`, `stop`, `attachment`, `diff`, `git`, `agent`, `ask`, `ask-in`, `rename`, `terminal`, `chat`

## The components

### KvitLabel

A run of chrome text at one of the seven type roles. Every other component uses it, which is what keeps the chrome family, the colour and the eliding rule in one place.

```go
func NewLabel(ui *UI, s string) *Label
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is what the label says. |
| `Role` | `TypeRole` | Role is what the text is; its size follows from it. |
| `Mono` | `bool` | Mono sets the text in the monospace family, for an identifier that has to line up down a column. |
| `Tabular` | `bool` | Tabular gives every digit the same width, so a column of figures lines up and a changing value does not shift the text beside it. |
| `Weight` | `int` | Weight is text.Regular unless set. |
| `Ink` | `Ink` | Ink is the colour; InkTextPrimary unless set. |
| `Wrap` | `bool` | Wrap lets the text run onto further lines instead of being cut short. |
| `LineHeight` | `float32` | LineHeight multiplies the lines' height, for a paragraph read in a popover; 0 is 1. |

*The seven roles*

```go
func labelRoles(ui *kvitui.UI) unison.Paneler {
    role := func(r kvitui.TypeRole, s string) *kvitui.Label {
        l := kvitui.NewLabel(ui, s)
        l.Role = r
        return l
    }
    return kvitui.Column(ui, kvitui.SizeSpaceSnug,
        role(kvitui.RoleDisplay, "display — a page title"),
        role(kvitui.RoleHeadline, "headline — a pane title"),
        role(kvitui.RoleTitle, "title — a section heading"),
        role(kvitui.RoleStrong, "strong — a name"),
        role(kvitui.RoleBody, "body — row text and prose"),
        role(kvitui.RoleSmall, "small — chip labels and sub-lines"),
        role(kvitui.RoleCaption, "caption — kind tags and counts"),
    )
}
```

*Monospace and tabular numerals*

```go
func labelMonoAndTabular(ui *kvitui.UI) unison.Paneler {
    id := kvitui.NewLabel(ui, "a1b2c3d4")
    id.Mono = true
    amount := kvitui.NewLabel(ui, "1,234.56")
    amount.Tabular = true
    return kvitui.Column(ui, kvitui.SizeSpaceSnug, id, amount)
}
```

### KvitIcon

A symbol asked for by what it means rather than by what it looks like. The font is embedded in the library, so nothing has to be added to an application.

```go
func NewIcon(ui *UI, name string) *Icon
```

| Field | Type | |
|---|---|---|
| `Name` | `string` | Name is the symbol's meaning name. |
| `Ink` | `Ink` | Ink is the colour; InkTextPrimary unless set. |
| `Size` | `Measure` | Size is the glyph's size; SizeIconSize unless set. |
| `Label` | `string` | Label, when set, is what a screen reader is told the symbol means. |

Methods: `Recognized`.

*At the sizes the chrome uses*

```go
func iconSizes(ui *kvitui.UI) unison.Paneler {
    search := kvitui.NewIcon(ui, "search")
    search.Size = kvitui.SizeIconSizeSmall
    trash := kvitui.NewIcon(ui, "trash")
    trash.Ink = kvitui.InkDanger
    success := kvitui.NewIcon(ui, "success")
    success.Ink = kvitui.InkSuccess
    return kvitui.Row(ui, kvitui.SizeSpace, search, kvitui.NewIcon(ui, "chevron-right"), trash, success)
}
```

*An unrecognised name is visible, not blank*

```go
func iconUnknown(ui *kvitui.UI) unison.Paneler {
    return kvitui.NewIcon(ui, "not-a-symbol")
}
```

### KvitIconButton

A button whose whole label is a symbol. A real button, so it takes tab focus and a screen reader is told it is there; `Label` fills both the accessible name and the tooltip shown on pointer hover, and `Explanation` is the second sentence beside it. `Dense` draws the symbol at the 13 every symbol beside words in this library is drawn at, rather than at 18.

```go
func NewIconButton(ui *UI, symbol, label string) *IconButton
```

| Field | Type | |
|---|---|---|
| `Symbol` | `string` | Symbol is the meaning name of the symbol drawn. |
| `Label` | `string` | Label says what the button does. |
| `Explanation` | `string` | Explanation says more than the label can; optional. |
| `Form` | `IconButtonForm` | Form is Quiet or Ordinary. |
| `Dense` | `bool` | Dense draws the symbol at the small icon size, the size every symbol beside words is drawn at, rather than the full one. |
| `Checkable` | `bool` | Checkable makes a press toggle Checked. |
| `Checked` | `bool` | Checked draws the button as on: an accent edge as well as a tint, so the state does not rest on colour alone. |
| `TooltipEnabled` | `bool` | TooltipEnabled shows the tooltip; true unless turned off by a control that already shows the label in a surface of its own. |
| `Size` | `Measure` | Size is the button's square side; a control's height unless set. |
| `OnClick` | `func()` | OnClick runs when the button is pressed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Quiet, ordinary and checked*

```go
func iconButtonForms(ui *kvitui.UI) unison.Paneler {
    remove := kvitui.NewIconButton(ui, "trash", "Delete")
    remove.Form = kvitui.Ordinary
    pin := kvitui.NewIconButton(ui, "pin", "Pin")
    pin.Checked = true
    duplicate := kvitui.NewIconButton(ui, "copy", "Copy")
    duplicate.SetEnabled(false)
    settings := kvitui.NewIconButton(ui, "settings", "Settings")
    settings.Focus()
    return kvitui.Row(ui, kvitui.SizeSpace,
        kvitui.NewIconButton(ui, "pencil", "Edit"), remove, pin, duplicate, settings)
}
```

*Two symbol sizes*

```go
func iconButtonSizes(ui *kvitui.UI) unison.Paneler {
    // A symbol standing on its own is drawn at the full icon size. Dense
    // draws it at the small one, the size every symbol beside words in this
    // library is drawn at: a strip of buttons across a header, or on a
    // heading bar, reads as larger than its row at the full size.
    strip := func(dense bool) *unison.Panel {
        search := kvitui.NewIconButton(ui, "search", "Search")
        add := kvitui.NewIconButton(ui, "plus", "New track")
        branch := kvitui.NewIconButton(ui, "git", "Branch")
        search.Dense, add.Dense, branch.Dense = dense, dense, dense
        return kvitui.Row(ui, kvitui.SizeSpace, search, add, branch)
    }
    return kvitui.Column(ui, kvitui.SizeSpace, strip(false), strip(true))
}
```

*A second sentence, where the name is not enough*

```go
func iconButtonExplanation(ui *kvitui.UI) unison.Paneler {
    // Label names the button and is what a screen reader is told.
    // Explanation is the sentence beside it, for what the name cannot hold:
    // what pressing this changes, or why it cannot be pressed. Both go to
    // the tooltip and to the accessible description, in that order.
    rename := kvitui.NewIconButton(ui, "rename", "Rename")
    rename.Explanation = "Renames the track and its branch. Its history and folder stay unchanged."
    rename.Focus()
    archive := kvitui.NewIconButton(ui, "archive", "Archive")
    archive.Explanation = "The branch has work that has not been pushed."
    archive.SetEnabled(false)
    return kvitui.Row(ui, kvitui.SizeSpace, rename, archive)
}
```

### KvitLink

An inline destination or action with link semantics, natural width and optional symbol. Hover and keyboard focus both add accent and an underline; no chevron, button ground or border is invented. `Explanation` says in a sentence where following it goes, shown as the tooltip and announced as the accessible description.

```go
func NewLink(ui *UI, s string) *Link
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is what the link says. |
| `Symbol` | `string` | Symbol is an optional meaning name drawn before the text. |
| `Role` | `TypeRole` | Role is the text's type role; RoleBody unless set. |
| `Explanation` | `string` | Explanation says in a sentence where following the link goes, for a destination its words alone do not describe. |
| `OnActivate` | `func()` | OnActivate runs when the link is followed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Plain, symbolic and keyboard-focused*

```go
func linkForms(ui *kvitui.UI) unison.Paneler {
    calendar := kvitui.NewLink(ui, "Open calendar")
    calendar.Symbol = "calendar"
    focused := kvitui.NewLink(ui, "Keyboard focus")
    focused.Focus()
    return kvitui.Row(ui, kvitui.SizeSpaceLoose, kvitui.NewLink(ui, "Privacy policy"), calendar, focused)
}
```

*Elided in a narrow column*

```go
func linkElided(ui *kvitui.UI) unison.Paneler {
    return kvitui.Width(ui, kvitui.Px(120), kvitui.NewLink(ui, "A destination whose full name does not fit here"))
}
```

### KvitHeader

The strip across the top of the window: the wordmark, navigation in the middle, actions on the right, in that fixed order on every screen.

```go
func NewHeader(ui *UI, wordmark string, navigation unison.Paneler, actions ...unison.Paneler) *Header
```

| Field | Type | |
|---|---|---|
| `Panel` | embedded | Everything a Panel has. |

*With navigation and actions*

```go
func headerWithNavigation(ui *kvitui.UI) unison.Paneler {
    today := kvitui.NewTab(ui, "Today")
    today.Selected = true
    projects := kvitui.NewTab(ui, "Projects")
    projects.Count = 26
    decisions := kvitui.NewTab(ui, "Decisions")
    decisions.Count = 4
    return kvitui.FullWidth(kvitui.NewHeader(ui, "kvit",
        kvitui.Row(ui, kvitui.SizeSpace, today, projects, decisions),
        kvitui.NewIconButton(ui, "settings", "Settings")))
}
```

### KvitSidebar

The list of places down the left edge. Below the laptop breakpoint it collapses to its rail and every item becomes its symbol alone.

```go
func NewSidebar(ui *UI, items ...unison.Paneler) *Sidebar
```

| Field | Type | |
|---|---|---|
| `Collapsed` | `bool` | Collapsed draws the sidebar as its rail. |

*Expanded and collapsed*

```go
func sidebarExpandedAndCollapsed(ui *kvitui.UI) unison.Paneler {
    today := kvitui.NewSidebarItem(ui, "Today", "calendar")
    today.Selected = true
    projects := kvitui.NewSidebarItem(ui, "Projects", "folder")
    projects.Count = 26
    decisions := kvitui.NewSidebarItem(ui, "Decisions", "question")
    decisions.Count = 4
    expanded := kvitui.NewSidebar(ui, today, projects, decisions)

    railToday := kvitui.NewSidebarItem(ui, "Today", "calendar")
    railToday.Selected = true
    rail := kvitui.NewSidebar(ui, railToday,
        kvitui.NewSidebarItem(ui, "Projects", "folder"),
        kvitui.NewSidebarItem(ui, "Decisions", "question"))
    rail.Collapsed = true
    return kvitui.Row(ui, kvitui.SizeColumnGap, expanded, rail)
}
```

### KvitSidebarItem

One place in the sidebar. The selected item takes a bar down its leading edge as well as a tint, so the selection is not resting on colour. A count is hidden once the sidebar collapses to the rail unless `CountInRail` asks for it, and `CountMax` raises the badge's cap where the place the item points at states the true number.

```go
func NewSidebarItem(ui *UI, label, symbol string) *SidebarItem
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text names the place. |
| `Symbol` | `string` | Symbol is the meaning name of the item's symbol. |
| `Selected` | `bool` | Selected marks the place the window is showing. |
| `Collapsed` | `bool` | Collapsed draws the item as part of the rail; its Sidebar sets it. |
| `Count` | `int` | Count is how many things are waiting there, drawn as a badge; below zero draws none, and zero draws no badge but is still announced. |
| `Counted` | `string` | Counted is the word for one of the things counted, so a screen reader hears "Review, 1 decision"; "" says "items". |
| `CountedPlural` | `string` | CountedPlural is the word for several; "" adds an s to Counted. |
| `CountInRail` | `bool` | CountInRail keeps the badge once the sidebar collapses, on the symbol's upper corner. |
| `CountMax` | `int` | CountMax is the largest count the badge writes out, 99 unless set. |
| `OnClick` | `func()` | OnClick runs when the item is pressed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Name`, `Pressed`.

*Selected, hovered and counted*

```go
func sidebarItemStates(ui *kvitui.UI) unison.Paneler {
    inbox := kvitui.NewSidebarItem(ui, "Inbox", "note")
    inbox.Selected = true
    inbox.Count, inbox.Counted = 12, "message"
    archive := kvitui.NewSidebarItem(ui, "Archive", "archive")
    return kvitui.Width(ui, kvitui.SizeSidebarWidth, kvitui.Column(ui, kvitui.Px(0), inbox, archive))
}
```

*The rail, with and without its counts*

```go
func sidebarItemRail(ui *kvitui.UI) unison.Paneler {
    // The default: a rail is a strip of symbols, and the count stays behind
    // with the label.
    review := kvitui.NewSidebarItem(ui, "Review", "question")
    review.Count = 214
    plain := kvitui.NewSidebar(ui, review, kvitui.NewSidebarItem(ui, "Accounts", "bank"))
    plain.Collapsed = true

    // A rail that has to keep the backlog in front of the reader. The badge
    // moves onto the symbol's upper corner, and CountMax is raised past the
    // badge's default of 99 because the screen this points at states the
    // true number, and a badge reading "99+" beside it disagrees with it.
    backlog := kvitui.NewSidebarItem(ui, "Review", "question")
    backlog.Count, backlog.CountInRail, backlog.CountMax = 214, true, 999
    counted := kvitui.NewSidebar(ui, backlog, kvitui.NewSidebarItem(ui, "Accounts", "bank"))
    counted.Collapsed = true
    return kvitui.Row(ui, kvitui.SizeColumnGap, plain, counted)
}
```

### KvitBreadcrumb

Where the reader is and the way back. The last crumb is the current place and is deliberately not a link; a long trail elides from the middle, keeping the section and the current place.

```go
func NewBreadcrumb(ui *UI, trail ...Crumb) *Breadcrumb
```

| Field | Type | |
|---|---|---|
| `Trail` | `[]Crumb` | Trail is the places from the root to the current one. |
| `MaximumVisible` | `int` | MaximumVisible is how many crumbs are drawn before the middle of the trail is replaced by "…"; 4 unless set. |
| `OnActivate` | `func(Crumb)` | OnActivate runs when a crumb before the current one is followed. |

Methods: `Shown`.

*A short trail and a long one*

```go
func breadcrumbTrails(ui *kvitui.UI) unison.Paneler {
    short := kvitui.NewBreadcrumb(ui,
        kvitui.Crumb{Label: "Projects", ID: "p"},
        kvitui.Crumb{Label: "kvit-ui", ID: "u"})
    long := kvitui.NewBreadcrumb(ui,
        kvitui.Crumb{Label: "Projects", ID: "p"}, kvitui.Crumb{Label: "kvit", ID: "k"},
        kvitui.Crumb{Label: "Wave 1", ID: "w"}, kvitui.Crumb{Label: "Tokens", ID: "t"},
        kvitui.Crumb{Label: "Theme", ID: "h"})
    return kvitui.Column(ui, kvitui.SizeSpace, short, long)
}
```

### KvitRegion

A body that takes the height left over and scrolls what does not fit. The scroll bar sits beside the content rather than over it, so nothing is ever hidden behind it.

```go
func NewRegion(ui *UI, content unison.Paneler) *Region
```

| Field | Type | |
|---|---|---|
| `Padding` | `Measure` | Padding is the space inside the scrolled area, around the content; the view margin unless set, which is what a body directly inside a window wants. |
| `Horizontal` | `bool` | Horizontal lets the content be wider than the region and scroll sideways too. |

Methods: `Position`, `ScrollTo`, `Scrolls`.

*Scrolling a column of rows*

```go
func regionScrolling(ui *kvitui.UI) unison.Paneler {
    rows := make([]unison.Paneler, 12)
    for i := range rows {
        rows[i] = kvitui.NewSlimRow(ui, fmt.Sprintf("Row %d", i+1))
    }
    region := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), rows...))
    return kvitui.Sized(ui, kvitui.Px(480), kvitui.Px(140), region)
}
```

### KvitViewHead

The strip at the top of a view: what the view is, how much is in it, and the controls that act on all of it.

```go
func NewViewHead(ui *UI, title string, controls ...unison.Paneler) *ViewHead
```

| Field | Type | |
|---|---|---|
| `Title` | `string` | Title names the view, and is the heading a screen reader announces. |
| `Count` | `int` | Count is how many things the view holds; below zero shows no count. |
| `Counted` | `string` | Counted is the word for one of the things counted, "item" unless set. |
| `CountedPlural` | `string` | CountedPlural is the word for several; "" adds an s to Counted. |
| `Subtitle` | `string` | Subtitle is a sentence under the title. |
| `Padding` | `Measure` | Padding is how far the content sits in from the left and right edges of the region the head belongs to; the view margin unless set. |

Methods: `CountPhrase`.

*Title, count and controls*

```go
func viewHeadWithControls(ui *kvitui.UI) unison.Paneler {
    head := kvitui.NewViewHead(ui, "Transactions",
        kvitui.Width(ui, kvitui.Px(180), kvitui.NewSearchField(ui)), kvitui.NewButton(ui, "Export"))
    head.Count, head.Counted = 1284, "transaction"
    head.Subtitle = "Everything since the account was opened"
    return kvitui.FullWidth(head)
}
```

### KvitStatusBar

The strip along the bottom: what is happening on the left, standing facts on the right. Facts come in three shapes — plain strings, named groups whose facts open what they name, and whole controls at the end. What the bar has no room for goes into a menu behind a control saying how many there are, rather than being cut off the end of the list. `GroupsFirst` puts the groups before the activity, for a bar whose left end is the list of what is waiting.

```go
func NewStatusBar(ui *UI, controls ...unison.Paneler) *StatusBar
```

| Field | Type | |
|---|---|---|
| `Activity` | `string` | Activity says what is happening now; "" leaves the left blank, which is the resting state. |
| `Facts` | `[]string` | Facts are the standing facts, at the right, separated by dots. |
| `Groups` | `[]StatusGroup` | Groups are the facts that do something. |
| `GroupsFirst` | `bool` | GroupsFirst puts the groups at the left, before the activity, for a bar whose left end is the list of what is waiting rather than a sentence about what is running. |
| `OnFact` | `func(group, fact int)` | OnFact runs when a fact is pressed, on the bar or in the menu, with the index of its group and its index in the group. |

Methods: `Shown`.

*Working, with two facts*

```go
func statusBarWorking(ui *kvitui.UI) unison.Paneler {
    bar := kvitui.NewStatusBar(ui)
    bar.Activity = "Reindexing 4 of 26 projects"
    bar.Facts = []string{"1,284 notes", "last synced 14:02"}
    return kvitui.FullWidth(bar)
}
```

*Grouped facts that open what they name*

```go
func statusBarGroups(ui *kvitui.UI) unison.Paneler {
    // Each fact is a control: it takes tab focus, says its own name to a
    // screen reader and answers Return and Space. The bar keeps none of the
    // caller's state: it is handed words and hands back which one was
    // pressed.
    bar := kvitui.NewStatusBar(ui)
    bar.Activity = "Indexing 4 of 26 working copies"
    bar.Groups = []kvitui.StatusGroup{
        {Label: "Waiting on you", Facts: []kvitui.StatusFact{
            {Text: "3 reviews", Symbol: "question"}, {Text: "1 conflict", Symbol: "warning"}}},
        {Label: "Running", Facts: []kvitui.StatusFact{{Text: "2 agents", Symbol: "robot"}}},
    }
    bar.OnFact = func(group, fact int) {
        bar.Activity = fmt.Sprintf("Pressed fact %d of group %d", fact, group)
        bar.MarkForLayoutAndRedraw()
    }
    return kvitui.FullWidth(bar)
}
```

*Narrow: what does not fit is in the menu, not gone*

```go
func statusBarNarrow(ui *kvitui.UI) unison.Paneler {
    // The same two groups in a bar too narrow for them. The groups that fit
    // are drawn left to right; the rest are behind the count of what is not
    // there, which is a link: Tab reaches it, Return opens the menu, and the
    // arrow keys move through it. Nothing is dropped, and nothing is silent
    // about it.
    narrow := kvitui.NewStatusBar(ui)
    narrow.Activity = "Indexing"
    narrow.Groups = []kvitui.StatusGroup{
        {Label: "Waiting on you", Facts: []kvitui.StatusFact{
            {Text: "3 reviews", Symbol: "question"}, {Text: "1 conflict", Symbol: "warning"}}},
        {Label: "Running", Facts: []kvitui.StatusFact{{Text: "2 agents", Symbol: "robot"}}},
    }
    // The resting height is unchanged: a bar holding only text is as tall as
    // the text, and only a control in the slot at the end makes it grow.
    plain := kvitui.NewStatusBar(ui)
    plain.Facts = []string{"1,284 notes", "last synced 14:02"}
    return kvitui.Column(ui, kvitui.SizeColumnGap,
        kvitui.Width(ui, kvitui.Px(420), narrow), kvitui.Width(ui, kvitui.Px(320), plain))
}
```

*Groups first, for a bar whose left end is the work*

```go
func statusBarGroupsFirst(ui *kvitui.UI) unison.Paneler {
    // A bar whose left end is the work rather than a sentence about it: what
    // is waiting and what is running, each item something to press, with
    // the standing facts at the right. GroupsFirst is the whole of the
    // difference.
    bar := kvitui.NewStatusBar(ui)
    bar.GroupsFirst = true
    bar.Activity = "Fetching origin"
    bar.Groups = []kvitui.StatusGroup{
        {Label: "Waiting on you", Facts: []kvitui.StatusFact{
            {Text: "3 reviews", Symbol: "question"}, {Text: "1 conflict", Symbol: "warning"}}},
        {Label: "Running", Facts: []kvitui.StatusFact{{Text: "2 agents", Symbol: "robot"}}},
    }
    bar.Facts = []string{"7 changes", "main, 2 ahead"}
    return kvitui.FullWidth(bar)
}
```

### KvitWindow

The application shell: header, optional sidebar, body and status bar, with the sidebar collapsing to its rail below the laptop breakpoint and opening over the body while the pointer rests on it. Not shown here because it is a window; see the code sample.

```go
func NewWindow(ui *UI, title string) (*Window, error)
```

| Field | Type | |
|---|---|---|
| `SidebarVisible` | `bool` | SidebarVisible false hides the sidebar entirely, which is a reader's choice, where a collapsed one is the window being narrow. |
| `SidebarExpandsOnHover` | `bool` | SidebarExpandsOnHover opens the rail while the pointer rests on it or the keyboard is inside it; on unless an application has a reason. |
| `OnKeyDown` | `func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool` | OnKeyDown is offered every key press before the focused control, for the application's own shortcuts; true says it used the key. |

Methods: `Narrow`, `Popups`, `SetBody`, `SetHeader`, `SetSidebar`, `SetStatusBar`, `Show`, `Showing`, `SidebarCollapsed`, `SidebarExpanded`, `SidebarNarrow`.

*The whole shell*

```go
func windowShell(ui *kvitui.UI) unison.Paneler {
    // Source only: a window has no place inside another window, so this is
    // the one page that shows a sample without drawing it. The gallery's
    // tests build and run it like every other sample.
    window, err := kvitui.NewWindow(ui, "kvit-cash")
    if err != nil {
        return nil
    }
    accounts := kvitui.NewTab(ui, "Accounts")
    accounts.Selected = true
    window.SetHeader(kvitui.NewHeader(ui, "kvit",
        kvitui.Row(ui, kvitui.SizeSpace, accounts, kvitui.NewTab(ui, "Budget"))))

    // The window says when it is narrow, and a Sidebar in it is drawn as a
    // rail when it is.
    everyday := kvitui.NewSidebarItem(ui, "Everyday", "wallet")
    everyday.Selected = true
    window.SetSidebar(kvitui.NewSidebar(ui, everyday,
        kvitui.NewSidebarItem(ui, "Savings", "bank"), kvitui.NewSidebarItem(ui, "Budget", "chart-line")))

    window.SetBody(kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0),
        kvitui.NewSlimRow(ui, "Checking"), kvitui.NewSlimRow(ui, "Joint checking"), kvitui.NewSlimRow(ui, "Emergency fund"))))

    status := kvitui.NewStatusBar(ui)
    status.Activity = "Matching 4 of 26 statements"
    status.Facts = []string{"1,284 transactions", "last synced 14:02"}
    window.SetStatusBar(status)
    window.ToFront()
    return nil
}
```

### KvitSectionHeading

A group heading: a filled bar with a disclosure chevron, the name, what the group holds, the count with the word for what was counted, and the one action that applies to every row under it. A count that is not a number goes in `CountText` and is drawn beside the name as written; `ActionSymbol` draws the action as a symbol and keeps its words as the button's name and tooltip, with `ActionExplanation` saying in a sentence what running it does. A collapsible heading joins the tab order and opens on Return, Enter or Space; the hoisted action is a control of its own, so running it never also collapses the group.

```go
func NewSectionHeading(ui *UI, label string) *SectionHeading
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text names the group. |
| `Count` | `int` | Count is how many rows the group holds; below zero draws none. |
| `Counted` | `string` | Counted is the word for one of them; "" leaves the number bare. |
| `CountedPlural` | `string` | CountedPlural is the word for several; "" adds an s to Counted. |
| `CountText` | `string` | CountText is a count the caller has already written out, for a count that is not a number, such as "1 · +0 −0". |
| `Kind` | `string` | Kind says what sort of thing the group holds, beside the name. |
| `Action` | `string` | Action is the hoisted action's words; "" for none. |
| `ActionSymbol` | `string` | ActionSymbol draws the action as a symbol, keeping its words as the button's name and tooltip, for an action a symbol can carry. |
| `ActionExplanation` | `string` | ActionExplanation says in a sentence what running the action does. |
| `Strong` | `bool` | Strong draws the name bold and the rule above in the strong border colour. |
| `Collapsible` | `bool` | Collapsible lets the heading open and close its group. |
| `Expanded` | `bool` | Expanded is whether the group is open. |
| `OnToggle` | `func(expanded bool)` | OnToggle runs after the heading opened or closed, with the new state. |
| `OnAction` | `func()` | OnAction runs when the hoisted action is taken. |

Methods: `CountPhrase`, `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Collapsible, counted, with an action*

```go
func sectionHeadingForms(ui *kvitui.UI) unison.Paneler {
    mine := kvitui.NewSectionHeading(ui, "Waiting on me")
    mine.Counted, mine.Count = "project", 4
    mine.Action, mine.Collapsible = "Hand all to an agent", true
    // A group of one is still counted. The count is drawn wherever the
    // caller gives one, so it does not come and go as the group changes size.
    sams := kvitui.NewSectionHeading(ui, "Waiting on Sam")
    sams.Counted, sams.Count = "project", 1
    archived := kvitui.NewSectionHeading(ui, "Archived")
    archived.Kind = "closed last quarter"
    archived.Collapsible, archived.Expanded, archived.Strong = true, false, true
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, mine, sams, archived))
}
```

*Reached by tab, opened by Return*

```go
func sectionHeadingKeyboard(ui *kvitui.UI) unison.Paneler {
    // The ring says the heading has the keyboard; Return, Enter and Space
    // open and close it. The action beside it is a Link, which takes its own
    // keys, so tabbing on to it and pressing Space runs the action and leaves
    // the group where it was.
    mine := kvitui.NewSectionHeading(ui, "Waiting on me")
    mine.Counted, mine.Count = "project", 4
    mine.Action, mine.Collapsible = "Hand all to an agent", true
    mine.Focus()
    archived := kvitui.NewSectionHeading(ui, "Archived")
    archived.Kind = "closed last quarter"
    archived.Collapsible, archived.Expanded = true, false
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, mine, archived))
}
```

*A written count, and an action drawn as a symbol*

```go
func sectionHeadingWrittenCount(ui *kvitui.UI) unison.Paneler {
    // Some counts are not numbers. A Changes heading reads "1 · +0 −0": one
    // changed file, and the lines added and removed across it, which the
    // service hands over already written and which no integer expresses.
    // CountText is drawn exactly as given, beside the name rather than at the
    // right end where a number goes.
    //
    // ActionSymbol draws the hoisted action as a symbol. The words stay in
    // Action and become the button's accessible name and its tooltip.
    changes := kvitui.NewSectionHeading(ui, "Changes")
    changes.CountText, changes.Collapsible = "1 · +0 −0", true
    changes.Action, changes.ActionSymbol = "Open the diff", "diff"
    agents := kvitui.NewSectionHeading(ui, "Agents")
    agents.CountText, agents.Collapsible = "2 running", true
    agents.Action, agents.ActionSymbol = "Start an agent", "plus"
    terminal := kvitui.NewSectionHeading(ui, "Terminal")
    terminal.Action, terminal.ActionSymbol = "Open a terminal", "terminal"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, changes, agents, terminal))
}
```

*What the action does, in a sentence*

```go
func sectionHeadingExplanation(ui *kvitui.UI) unison.Paneler {
    // ActionExplanation is the sentence beside the action's name, on the
    // action rather than on the heading: a heading is not something a reader
    // presses. Hover either action, or tab to it, to read it; it is the
    // accessible description too.
    agents := kvitui.NewSectionHeading(ui, "Agents")
    agents.CountText, agents.Collapsible = "2 running", true
    agents.Action, agents.ActionSymbol = "Start an agent", "plus"
    agents.ActionExplanation = "Creates an agent at this project’s root. An agent may change files. Nothing runs until you send its first message."
    changes := kvitui.NewSectionHeading(ui, "Changes")
    changes.CountText, changes.Collapsible = "1 · +0 −0", true
    changes.Action = "Review all"
    changes.ActionExplanation = "Opens every changed file as one diff."
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, agents, changes))
}
```

### KvitRow

A list row at one of four heights, chosen by what the row carries rather than by how many rows a view wants to fit. Hover and keyboard focus are separate marks, because they are different rows. A row that does something when it is pressed says so first: the hover tint, the press and the chevron at its trailing edge all follow `Interactive`.

*Pressable and static*

```go
func rowPressableAndStatic(ui *kvitui.UI) unison.Paneler {
    // A row that opens something is reachable by the keyboard, and takes the
    // hover tint, the press and the chevron with it.
    opens := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, "opens the record")))
    opens.Interactive, opens.Label = true, "Groceries"
    // A field name beside its value declares nothing, draws no hover tint and
    // answers no press.
    static := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, "layout only")))
    static.Label = "Amount"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), opens, static))
}
```

*Opened from the keyboard, without opening it twice*

```go
func rowFromTheKeyboard(ui *kvitui.UI) unison.Paneler {
    // Tab reaches this row, the ring says which row the keyboard is on, and
    // Return, Enter or Space opens it. A screen reader's press action does
    // the same one thing.
    groceries := kvitui.NewListRow(ui, kvitui.NewLabel(ui, "keyboard focus — Return opens it"))
    groceries.Form, groceries.Interactive, groceries.Label = kvitui.RowSub, true, "Groceries"
    groceries.Focus()
    // A row is usually a container, and the key may be meant for something
    // inside it. The row answers only while the row itself has the keyboard,
    // so Space on this button presses the button and leaves the record shut.
    split := kvitui.NewButton(ui, "Split")
    split.Form = kvitui.ButtonQuiet
    rent := kvitui.NewListRow(ui, kvitui.FullWidth(kvitui.NewLabel(ui, "the button takes its own Space")), split)
    rent.Form, rent.Interactive, rent.Label = kvitui.RowSub, true, "Rent"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), groceries, rent))
}
```

*The four heights, and the three states*

```go
func rowHeights(ui *kvitui.UI) unison.Paneler {
    row := func(form kvitui.RowForm, label, note string) *kvitui.ListRow {
        r := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, note)))
        r.Form, r.Label = form, label
        return r
    }
    slim := row(kvitui.RowSlim, "Slim", "slim — 30, selected")
    slim.Selected = true
    compact := row(kvitui.RowCompact, "Compact", "compact — 24, keyboard focus")
    compact.Current = true
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0),
        row(kvitui.RowFull, "Full", "full — 56"), row(kvitui.RowSub, "Sub", "sub — 48"), slim, compact))
}
```

### KvitSlimRow

The row a reader sees most: a name, what it is, one phrase about it and one figure, right aligned, in that order on every screen.

```go
func NewSlimRow(ui *UI, name string) *SlimRow
```

| Field | Type | |
|---|---|---|
| `ListRow` | embedded | Everything a ListRow has. |
| `Name` | `string` | Name is what the row is about. |
| `Kind` | `string` | Kind says what sort of thing it is, where the name does not. |
| `Phrase` | `string` | Phrase says one thing about its state, muted, because it is context rather than identity. |
| `Figure` | `string` | Figure is the number at the right, as the caller formatted it. |
| `Unit` | `string` | Unit follows the figure. |
| `Measured` | `bool` | Measured is false for a value nobody measured, drawn as "—". |
| `Symbol` | `string` | Symbol is an optional meaning name drawn before the name. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Measured and unmeasured*

```go
func slimRowMeasured(ui *kvitui.UI) unison.Paneler {
    library := kvitui.NewSlimRow(ui, "kvit-ui")
    library.Symbol, library.Kind, library.Phrase = "folder", "library", "waiting on review"
    library.Figure, library.Unit = "3.5", "d"
    app := kvitui.NewSlimRow(ui, "kvit-cash")
    app.Symbol, app.Kind, app.Phrase = "folder", "application", "not started"
    app.Measured = false
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), library, app))
}
```

### KvitCard

A bounded block of content on a surface. An interactive card takes the control-boundary token, the hover tint and the chevron a pressable row carries, so a card that responds to a click looks like it will before it is clicked.

```go
func NewCard(ui *UI, content ...unison.Paneler) *Card
```

| Field | Type | |
|---|---|---|
| `Interactive` | `bool` | Interactive makes the card something to press. |
| `NoOpensMark` | `bool` | NoOpensMark leaves the chevron off a pressable card whose ways in are the rows inside it, so each row carries its own. |
| `OpensLabel` | `string` | OpensLabel says in words what pressing the card opens: the tooltip on the chevron, and the card's accessible description. |
| `Selected` | `bool` | Selected marks a chosen card. |
| `Label` | `string` | Label is what a screen reader calls the card; the words it holds unless set. |
| `Padding` | `Measure` | Padding is the space inside the edge; the loose space unless set. |
| `OnActivate` | `func()` | OnActivate runs when a pressable card is pressed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Static and interactive*

```go
func cardForms(ui *kvitui.UI) unison.Paneler {
    interactive := kvitui.NewCard(ui, kvitui.NewLabel(ui, "Interactive"))
    interactive.Interactive = true
    selected := kvitui.NewCard(ui, kvitui.NewLabel(ui, "Selected"))
    selected.Selected = true
    return kvitui.Row(ui, kvitui.SizeColumnGap, kvitui.NewCard(ui, kvitui.NewLabel(ui, "A card")), interactive, selected)
}
```

### KvitPanel

A region of the window with its own ground: a sidebar, a toolbar strip. Structural, where a card is content — which is why it has no radius.

```go
func NewPanel(ui *UI) *Panel
```

| Field | Type | |
|---|---|---|
| `RuleTop` | `bool` | RuleTop, RuleBottom, RuleLeft and RuleRight draw a hairline along that edge. |
| `RuleBottom` | `bool` | RuleTop, RuleBottom, RuleLeft and RuleRight draw a hairline along that edge. |
| `RuleLeft` | `bool` | RuleTop, RuleBottom, RuleLeft and RuleRight draw a hairline along that edge. |
| `RuleRight` | `bool` | RuleTop, RuleBottom, RuleLeft and RuleRight draw a hairline along that edge. |

*With rules on two edges*

```go
func panelWithRules(ui *kvitui.UI) unison.Paneler {
    p := kvitui.NewPanel(ui)
    p.RuleTop, p.RuleBottom = true, true
    label := kvitui.NewLabel(ui, "A panel")
    label.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Middle, VAlign: align.Middle, HGrab: true, VGrab: true})
    p.AddChild(label)
    p.SetLayout(kvitui.AtLeast(ui, kvitui.Px(60), &unison.FlexLayout{Columns: 1}))
    return kvitui.FullWidth(p)
}
```

### KvitPane

The side pane: detail about the one thing selected, the same width wherever it appears so the list beside it never reflows.

```go
func NewPane(ui *UI, title string, content ...unison.Paneler) *Pane
```

| Field | Type | |
|---|---|---|
| `Panel` | embedded | Everything a Panel has. |
| `Title` | `string` | Title names what the pane shows. |
| `CloseLabel` | `string` | CloseLabel is what the close control is called, on hover and to a screen reader: "Close record" for a pane holding one record, since the reader is closing the record rather than the furniture it arrived in. |
| `Closable` | `bool` | Closable draws the close control; true unless set, and off for a pane the reader cannot dismiss (kvit-cash). |
| `OnClose` | `func()` | OnClose runs when the close control is pressed. |

Methods: `Open`, `SetOpen`.

*Open*

```go
func paneOpen(ui *kvitui.UI) unison.Paneler {
    pane := kvitui.NewPane(ui, "Transaction", kvitui.Centred(kvitui.NewLabel(ui, "Detail goes here")))
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(160), kvitui.WithPane(ui, nil, pane)))
}
```

### KvitDivider

A rule between two things. One design pixel, in the decorative border token rather than the control-boundary one.

```go
func NewDivider(ui *UI) *Divider
```

| Field | Type | |
|---|---|---|
| `Vertical` | `bool` | Vertical draws the rule top to bottom rather than across. |
| `Inset` | `Measure` | Inset stops the rule short of both ends, for a rule inside a padded container where a full-length line would cut the padding in half. |

*Horizontal and vertical*

```go
func dividerBothWays(ui *kvitui.UI) unison.Paneler {
    down := kvitui.NewDivider(ui)
    down.Vertical = true
    down.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Fill, VGrab: true})
    beside := kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewLabel(ui, "left"), down, kvitui.NewLabel(ui, "right"))
    beside.SetLayout(kvitui.AtLeast(ui, kvitui.Px(24), beside.Layout()))
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewDivider(ui), beside))
}
```

### KvitDisclosure

A trigger with a body under it. Takes a `Group` so several sections can behave as an accordion, which is what most of the eighteen hand-rolled versions in the estate actually are.

```go
func NewDisclosure(ui *UI, title string, content ...unison.Paneler) *Disclosure
```

| Field | Type | |
|---|---|---|
| `Title` | `string` | Title names the section. |
| `Group` | `string` | Group makes sections sharing it an accordion; "" stands alone. |
| `Count` | `int` | Count is how many things are inside, beside the title; below zero draws none. |
| `OnToggle` | `func(expanded bool)` | OnToggle runs after the section opened or closed, with the new state. |

Methods: `Expanded`, `SetExpanded`.

*Open and closed*

```go
func disclosureOpenAndClosed(ui *kvitui.UI) unison.Paneler {
    changed := kvitui.NewDisclosure(ui, "What changed", kvitui.NewLabel(ui, "Three files were rewritten."))
    changed.Count = 3
    changed.SetExpanded(true)
    unchanged := kvitui.NewDisclosure(ui, "What did not")
    unchanged.Count = 12
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, changed, unchanged))
}
```

### KvitEmptyState

What a view says when it has nothing to show: what would be here, why it is not, and the action that would fill it. Also the answer for a chart with no data, in place of an axis drawn around zeros. The compact form is the same sentence on one line, at the height of a slim row and starting where the rows start, for a section in a stack of sections that may each be empty.

```go
func NewEmptyState(ui *UI, title string) *EmptyState
```

| Field | Type | |
|---|---|---|
| `Form` | `EmptyForm` | Form is the full block or the compact line. |
| `Title` | `string` | Title says what would be here, in one short sentence. |
| `Detail` | `string` | Detail says why it is not, or what to do about it; optional. |
| `Symbol` | `string` | Symbol is an optional meaning name. |
| `Action` | `string` | Action is the words of the action that would fill it; "" for none. |
| `Dashed` | `bool` | Dashed draws a dashed edge, for a region something can be put into. |
| `OnAction` | `func()` | OnAction runs when the action is taken. |

*With an action*

```go
func emptyStateWithAction(ui *kvitui.UI) unison.Paneler {
    empty := kvitui.NewEmptyState(ui, "No transactions yet")
    empty.Symbol = "wallet"
    empty.Detail = "Import a statement or add one by hand, and it will appear here."
    empty.Action = "Import a statement"
    return kvitui.FullWidth(empty)
}
```

*Dashed, as a drop target*

```go
func emptyStateDropTarget(ui *kvitui.UI) unison.Paneler {
    drop := kvitui.NewEmptyState(ui, "Drop a statement here")
    drop.Dashed, drop.Symbol = true, "file-arrow-down"
    drop.Detail = "CSV, OFX and QIF. The file is read on this machine and nothing is sent anywhere."
    drop.Action = "Choose a file"
    return kvitui.FullWidth(drop)
}
```

*Compact, one line per empty section*

```go
func emptyStateCompact(ui *kvitui.UI) unison.Paneler {
    // A column of sections, each of which may have nothing in it. The full
    // block is several times taller than the rows it stands in for, so a
    // stack of them uses the one-line form; the words are the same words.
    changes := kvitui.NewSectionHeading(ui, "Changes")
    changes.Count, changes.Counted = 0, "change"
    noChanges := kvitui.NewEmptyState(ui, "No changes")
    noChanges.Form = kvitui.EmptyCompact
    agents := kvitui.NewSectionHeading(ui, "Agents")
    agents.Count, agents.Counted = 0, "agent"
    noAgents := kvitui.NewEmptyState(ui, "No agents")
    noAgents.Form, noAgents.Symbol = kvitui.EmptyCompact, "robot"
    noAgents.Detail, noAgents.Action = "none started here yet", "Start one"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceSnug, changes, noChanges, agents, noAgents))
}
```

### KvitChip

A small labelled mark saying what kind of thing this is or what state it is in. The tone names a meaning rather than a colour, and every tone carries its outline as well as its tint so the distinction is not resting on hue. It is a mark and not a control: a chip that opens something when it is pressed is KvitChipButton, drawn from this same tone table.

```go
func NewChip(ui *UI, word string) *Chip
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the chip's word. |
| `Tone` | `Tone` | Tone is what its colour means. |
| `Strong` | `bool` | Strong fills the chip, for the one chip on a row that is its point. |
| `Symbol` | `string` | Symbol is an optional meaning name before the word. |
| `Explanation` | `string` | Explanation says in a sentence what the word means and what to do about it: the tooltip and the accessible description. |

*Every tone, tinted and filled*

```go
func chipTones(ui *kvitui.UI) unison.Paneler {
    chip := func(word string, tone kvitui.Tone) unison.Paneler {
        c := kvitui.NewChip(ui, word)
        c.Tone = tone
        return c
    }
    settled := kvitui.NewChip(ui, "settled")
    settled.Tone, settled.Strong = kvitui.ToneSuccess, true
    disputed := kvitui.NewChip(ui, "disputed")
    disputed.Tone, disputed.Strong, disputed.Symbol = kvitui.ToneDanger, true, "warning"
    return kvitui.Column(ui, kvitui.SizeSpace,
        kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, chip("neutral", kvitui.ToneNeutral), chip("accent", kvitui.ToneAccent),
            chip("success", kvitui.ToneSuccess), chip("warning", kvitui.ToneWarning), chip("danger", kvitui.ToneDanger),
            chip("info", kvitui.ToneInfo))),
        kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, settled, disputed)))
}
```

*A state whose word is not enough to act on*

```go
func chipExplanation(ui *kvitui.UI) unison.Paneler {
    // A chip that names a kind of thing is finished at its word. A chip that
    // names a state somebody has to act on is not, and Explanation is where
    // the rest of it goes: the tooltip and the accessible description.
    stale := kvitui.NewChip(ui, "stale")
    stale.Tone = kvitui.ToneWarning
    stale.Explanation = "The note changed since this was staged. Update redoes it against the note as it stands now."
    readOnly := kvitui.NewChip(ui, "read only")
    readOnly.Explanation = "Another session holds the write lease on this copy, so this turn only reads."
    return kvitui.Row(ui, kvitui.SizeSpaceNear, stale, readOnly)
}
```

### KvitTag

A label a person put there, in a colour they chose. Theme-independent, because the reader picked that red and would not expect it to change; the label colour is derived from the fill rather than taken from a token.

```go
func NewTag(ui *UI, label string) *Tag
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the tag. |
| `Tint` | `palette.Color` | Tint is the colour the reader chose; the zero value is the neutral chip ground. |
| `Tinted` | `bool` | Tinted says Tint is set. |
| `Removable` | `bool` | Removable draws a close control, which calls OnRemove. |
| `OnRemove` | `func()` | OnRemove runs when the close control is pressed. |

Methods: `SetTint`.

*Tinted, plain and removable*

```go
func tagForms(ui *kvitui.UI) unison.Paneler {
    colours := tokens.ColorPalette()
    groceries := kvitui.NewTag(ui, "groceries")
    groceries.SetTint(colours[0])
    transport := kvitui.NewTag(ui, "transport")
    transport.SetTint(colours[2])
    removable := kvitui.NewTag(ui, "removable")
    removable.SetTint(colours[4])
    removable.Removable = true
    return kvitui.Row(ui, kvitui.SizeSpaceNear, groceries, transport, kvitui.NewTag(ui, "untinted"), removable)
}
```

### KvitBadge

A count attached to something else. Caps rather than growing wide, and hides at zero — a badge showing nought says look here about nothing. The number is drawn and announced with the reader's own digit grouping, and the noun beside it comes from the caller in two slots, a singular and a plural.

```go
func NewBadge(ui *UI, count int) *Badge
```

| Field | Type | |
|---|---|---|
| `Count` | `int` | Count is the number; zero or less hides the badge. |
| `Max` | `int` | Max is the largest count written out; 99 unless set. |
| `Counted` | `string` | Counted is the word for one of the things counted, so a screen reader hears "1 decision" or "11 decisions"; "" says "items". |
| `CountedPlural` | `string` | CountedPlural is the word for several; "" adds an s to Counted. |
| `Tone` | `BadgeTone` | Tone is the pill's colour. |

Methods: `Phrase`, `Shows`, `Text`.

*Counts, capped, and hidden at zero*

```go
func badgeCounts(ui *kvitui.UI) unison.Paneler {
    neutral := kvitui.NewBadge(ui, 42)
    neutral.Tone = kvitui.BadgeNeutral
    danger := kvitui.NewBadge(ui, 1204)
    danger.Tone = kvitui.BadgeDanger
    return kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewBadge(ui, 3), neutral, danger, kvitui.NewBadge(ui, 0))
}
```

*What the count counts*

```go
func badgeNouns(ui *kvitui.UI) unison.Paneler {
    // Same three pills, three different sentences for a screen reader:
    // "1 decision", "214 decisions", "214 entries". One already inflected
    // word would be right at one count and wrong at every other.
    one := kvitui.NewBadge(ui, 1)
    one.Counted = "decision"
    many := kvitui.NewBadge(ui, 214)
    many.Counted, many.Max = "decision", 999
    entries := kvitui.NewBadge(ui, 214)
    entries.Max, entries.Tone = 999, kvitui.BadgeNeutral
    entries.Counted, entries.CountedPlural = "entry", "entries"
    return kvitui.Row(ui, kvitui.SizeSpace, one, many, entries)
}
```

### KvitSlug

An identifier: a reference, a hash, a key. Monospace, because the task is comparison rather than reading, and elided from the middle because the end is what distinguishes one from its neighbours.

```go
func NewSlug(ui *UI, id string) *Slug
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the identifier. |
| `Ground` | `bool` | Ground draws the inline-code ground; true unless turned off, for a slug in a table cell where a ground on every row reads as a column of boxes. |

*With and without a ground*

```go
func slugGrounds(ui *kvitui.UI) unison.Paneler {
    bare := kvitui.NewSlug(ui, "9f2c1ab4e77d0031")
    bare.Ground = false
    return kvitui.Column(ui, kvitui.SizeSpaceNear, kvitui.Left(kvitui.NewSlug(ui, "TX-00173404")), kvitui.Left(bare))
}
```

### KvitDot

A small filled circle standing for one thing's state. The shape is the second channel: a level that differs only by hue says nothing to a reader who cannot separate red from amber.

```go
func NewDot(ui *UI) *Dot
```

| Field | Type | |
|---|---|---|
| `Ink` | `Ink` | Ink is the colour; the muted text colour unless set. |
| `Shape` | `Shape` | Shape is circle, square or diamond. |
| `Hollow` | `bool` | Hollow draws the outline only. |
| `Label` | `string` | Label says what the dot means, for a screen reader; "" hides it from one. |
| `Size` | `Measure` | Size is the dot's side; a near space unless set. |

*Three shapes, three levels*

```go
func dotLevels(ui *kvitui.UI) unison.Paneler {
    dot := func(ink kvitui.Ink, shape kvitui.Shape, hollow bool, label string) unison.Paneler {
        d := kvitui.NewDot(ui)
        d.Ink, d.Shape, d.Hollow, d.Label = ink, shape, hollow, label
        return d
    }
    return kvitui.Row(ui, kvitui.SizeSpace,
        dot(kvitui.InkSuccess, kvitui.ShapeCircle, false, "healthy"),
        dot(kvitui.InkWarning, kvitui.ShapeSquare, false, "slipping"),
        dot(kvitui.InkDanger, kvitui.ShapeDiamond, false, "stalled"),
        dot(kvitui.InkTextMuted, kvitui.ShapeCircle, true, "not measured"))
}
```

### KvitSignal

A mark saying what state something is in and how many things are in it, for a list where one row may have several of each. The number is drawn only past one, because a column of marks all reading `1` says nothing the mark did not already say. The colour is the caller’s, since which states exist is an application’s own question; the shape and the hollow form are the second channel, so states told apart by hue alone are not.

```go
func NewSignal(ui *UI, label string) *Signal
```

| Field | Type | |
|---|---|---|
| `Count` | `int` | Count is how many are in the state; zero or less hides the mark. |
| `Max` | `int` | Max is the largest count written out; 99 unless set. |
| `Ink` | `Ink` | Ink is the state's colour, the application's own idea; the accent unless set. |
| `Shape` | `Shape` | Shape is square or circle; a diamond would turn the number with it. |
| `Hollow` | `bool` | Hollow draws the outline only, a third distinction without another hue. |
| `Label` | `string` | Label says what the mark means, "2 agents running": the colour and the number do not say it. |

*One of each, and several of one*

```go
func signalCounts(ui *kvitui.UI) unison.Paneler {
    // At one the mark is the whole statement. Past one the count goes inside
    // it, and the mark grows sideways to hold the digits rather than growing
    // taller and pushing the row apart.
    signal := func(count int, ink kvitui.Ink, shape kvitui.Shape, hollow bool, label string) unison.Paneler {
        s := kvitui.NewSignal(ui, label)
        s.Count, s.Ink, s.Shape, s.Hollow = count, ink, shape, hollow
        return s
    }
    return kvitui.Column(ui, kvitui.SizeSpace,
        kvitui.Left(kvitui.Row(ui, kvitui.SizeSpace,
            signal(1, kvitui.InkAccent, kvitui.ShapeSquare, false, "1 needs you"),
            signal(1, kvitui.InkSuccess, kvitui.ShapeCircle, false, "1 running"),
            signal(1, kvitui.InkDanger, kvitui.ShapeSquare, true, "1 failed"))),
        kvitui.Left(kvitui.Row(ui, kvitui.SizeSpace,
            signal(4, kvitui.InkAccent, kvitui.ShapeSquare, false, "4 need you"),
            signal(12, kvitui.InkSuccess, kvitui.ShapeCircle, false, "12 running"),
            signal(240, kvitui.InkDanger, kvitui.ShapeSquare, true, "240 failed"))))
}
```

*Beside the words it marks*

```go
func signalBesideWords(ui *kvitui.UI) unison.Paneler {
    // Where these are drawn: at the head of a row, in front of what the row
    // is about. The mark is a caption's height plus a margin, which lets it
    // sit in a line of text without setting the line's height.
    scoutSignal := kvitui.NewSignal(ui, "2 need you")
    scoutSignal.Count = 2
    scoutName := kvitui.NewLabel(ui, "Dialog Scout")
    scoutName.Role = kvitui.RoleSmall
    scout := kvitui.NewListRow(ui, scoutSignal, scoutName)
    scout.Label = "Dialog Scout, 2 need you"

    nightlySignal := kvitui.NewSignal(ui, "1 running")
    nightlySignal.Ink, nightlySignal.Shape = kvitui.InkSuccess, kvitui.ShapeCircle
    nightlyName := kvitui.NewLabel(ui, "Nightly check")
    nightlyName.Role = kvitui.RoleSmall
    nightly := kvitui.NewListRow(ui, nightlySignal, nightlyName)
    nightly.Label = "Nightly check, running"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), scout, nightly))
}
```

### KvitPip

A row of dots standing for a small count — three of five days recorded. For counts a reader takes in without counting; past about seven, the figure is faster.

```go
func NewPip(ui *UI, filled, total int) *Pip
```

| Field | Type | |
|---|---|---|
| `Filled` | `int` | Filled of Total dots are filled. |
| `Total` | `int` | Filled of Total dots are filled. |
| `Ink` | `Ink` | Ink is the filled dots' colour; the accent unless set. |
| `Label` | `string` | Label says what is counted; "3 of 5" unless set. |

*Three of five, and none of four*

```go
func pipCounts(ui *kvitui.UI) unison.Paneler {
    days := kvitui.NewPip(ui, 3, 5)
    days.Label = "3 of 5 days recorded"
    checks := kvitui.NewPip(ui, 0, 4)
    checks.Ink, checks.Label = kvitui.InkWarning, "no checks passed"
    return kvitui.Column(ui, kvitui.SizeSpace, kvitui.Left(days), kvitui.Left(checks))
}
```

### KvitFigure

A measured value: tabular numerals, the unit in muted colour at the smaller role, and an em dash where nothing was measured rather than a zero. A balance nobody computed is not a balance of zero.

```go
func NewFigure(ui *UI, value, unit string) *Figure
```

| Field | Type | |
|---|---|---|
| `Value` | `string` | Value is the number as the caller formatted it. |
| `Unit` | `string` | Unit follows the value; "" for none. |
| `Measured` | `bool` | Measured is false for a value nobody measured, which draws "—". |
| `Bounded` | `bool` | Bounded marks an upper bound rather than a measurement, drawn "~12". |
| `Role` | `TypeRole` | Role is the value's type role; RoleBody unless set. |
| `Ink` | `Ink` | Ink is the value's colour; InkTextPrimary unless set. |

Methods: `Phrase`.

*Measured, unmeasured and bounded*

```go
func figureForms(ui *kvitui.UI) unison.Paneler {
    unmeasured := kvitui.NewFigure(ui, "", "")
    unmeasured.Measured = false
    bounded := kvitui.NewFigure(ui, "12", "d")
    bounded.Bounded = true
    large := kvitui.NewFigure(ui, "94", "%")
    large.Role = kvitui.RoleDisplay
    return kvitui.Column(ui, kvitui.SizeSpaceNear,
        kvitui.NewFigure(ui, "1,284.50", "GBP"), kvitui.NewFigure(ui, "0", "d"), unmeasured, bounded, large)
}
```

### KvitBeforeAfter

One value as it stands and the value something proposes to replace it with. Position, colour weight and the arrow all say which is which, so no reader depends on separating the two colours. A record being added has no before, and the em dash says so.

```go
func NewBeforeAfter(ui *UI, label, before, after string) *BeforeAfter
```

| Field | Type | |
|---|---|---|
| `Before` | `string` | Before and After are the values, already formatted. |
| `After` | `string` | Before and After are the values, already formatted. |
| `Unit` | `string` | Unit is drawn after each value. |
| `BeforeMeasured` | `bool` | BeforeMeasured and AfterMeasured false draw an em dash on that side. |
| `AfterMeasured` | `bool` | BeforeMeasured and AfterMeasured false draw an em dash on that side. |
| `Label` | `string` | Label says what the pair is a value of; "" for none, where a heading already says it. |
| `Role` | `TypeRole` | Role is both values' type role, so the two are always the same size. |

Methods: `Unchanged`.

*Changed, added, removed and unchanged*

```go
func beforeAfterForms(ui *kvitui.UI) unison.Paneler {
    amount := kvitui.NewBeforeAfter(ui, "Amount", "42.00", "44.50")
    amount.Unit = "GBP"
    category := kvitui.NewBeforeAfter(ui, "Category", "", "Groceries")
    category.BeforeMeasured = false
    payee := kvitui.NewBeforeAfter(ui, "Payee", "TESCO 4471", "")
    payee.AfterMeasured = false
    date := kvitui.NewBeforeAfter(ui, "Date", "2026-08-14", "2026-08-14")
    balance := kvitui.NewBeforeAfter(ui, "", "1,284.50", "1,301.75")
    balance.Unit, balance.Role = "GBP", kvitui.RoleStrong
    return kvitui.Column(ui, kvitui.SizeSpaceNear,
        kvitui.Left(amount), kvitui.Left(category), kvitui.Left(payee), kvitui.Left(date), kvitui.Left(balance))
}
```

### KvitButton

A button with words on it, in three forms. One primary per screen region: it is the action the screen is for. `Danger` is separate from the form, because a destructive action can be any of the three. `Explanation` is one sentence saying what the words cannot — why it is disabled, or what pressing it opens — and it is read on a disabled button too.

```go
func NewButton(ui *UI, label string) *Button
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is what the button says. |
| `Form` | `ButtonForm` | Form is ordinary, primary or quiet. |
| `Danger` | `bool` | Danger draws the button in the danger colour. |
| `Symbol` | `string` | Symbol is an optional meaning name drawn before the words. |
| `Busy` | `bool` | Busy says the button is doing something: it stays enabled so its name can say what, and the words change to BusyText rather than a spinner appearing, since a spinner says "wait" without saying for what. |
| `BusyText` | `string` | BusyText is what a busy button says; "Working…" unless set. |
| `Explanation` | `string` | Explanation is one sentence the words cannot say: why the button is disabled, or what pressing it opens. |
| `Checkable` | `bool` | Checkable makes a press toggle Checked. |
| `Checked` | `bool` | Checked draws a button that turns a mode on as on: tinted and outlined in the accent, two marks rather than one, since a tint alone is a few percent of lightness in the high-contrast theme. |
| `OnClick` | `func()` | OnClick runs when the button is pressed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Text, icon plus text, and busy text in every form*

```go
func buttonForms(ui *kvitui.UI) unison.Paneler {
    row := func(label, symbol string, busy bool) *unison.Panel {
        var buttons []unison.Paneler
        for _, form := range []kvitui.ButtonForm{kvitui.ButtonPrimary, kvitui.ButtonOrdinary, kvitui.ButtonQuiet} {
            b := kvitui.NewButton(ui, label)
            b.Form, b.Symbol, b.Busy = form, symbol, busy
            buttons = append(buttons, b)
        }
        return kvitui.Row(ui, kvitui.SizeSpace, buttons...)
    }
    return kvitui.Column(ui, kvitui.SizeSpace, row("Save", "", false), row("Schedule", "calendar", false), row("Save", "", true))
}
```

*Destructive and disabled*

```go
func buttonDanger(ui *kvitui.UI) unison.Paneler {
    primary := kvitui.NewButton(ui, "Delete")
    primary.Form, primary.Danger = kvitui.ButtonPrimary, true
    ordinary := kvitui.NewButton(ui, "Delete")
    ordinary.Danger, ordinary.Symbol = true, "trash"
    disabled := kvitui.NewButton(ui, "Disabled")
    disabled.SetEnabled(false)
    return kvitui.Row(ui, kvitui.SizeSpace, primary, ordinary, disabled)
}
```

*A button that stays on, beside the same button off*

```go
func buttonChecked(ui *kvitui.UI) unison.Paneler {
    on := kvitui.NewButton(ui, "Select region")
    on.Checkable, on.Checked = true, true
    off := kvitui.NewButton(ui, "Select region")
    off.Checkable = true
    return kvitui.Row(ui, kvitui.SizeSpace, on, off)
}
```

*Why a button cannot be pressed*

```go
func buttonExplanations(ui *kvitui.UI) unison.Paneler {
    // The reason a control is in the state it is in often lives somewhere the
    // reader cannot see. A disabled button with nothing to say about itself
    // is a grey rectangle and no account of it; Explanation is where that
    // sentence goes, shown on hover and announced as the accessible
    // description. It is read on a disabled button too, which is the case it
    // exists for.
    pull := kvitui.NewButton(ui, "Pull")
    pull.Explanation = "Brings the 2 commits on the remote into this branch."
    push := kvitui.NewButton(ui, "Push")
    push.SetEnabled(false)
    push.Explanation = "Nothing here has been committed yet."
    archive := kvitui.NewButton(ui, "Archive")
    archive.SetEnabled(false)
    archive.Danger = true
    archive.Explanation = "The branch has work that has not been pushed."
    return kvitui.Row(ui, kvitui.SizeSpace, pull, push, archive)
}
```

### KvitChipButton

KvitChip's twin for a fact that opens something. Drawn from the same tone table, so a row mixing facts that act with facts that do not reads as one row; what separates them is what a control has anyway — a ground that changes under the pointer, a hand cursor, a focus ring and a button role. A chip that cannot be pressed keeps its place in the tab order and says why, a chip that can may say where it goes, and the chip whose destination is already open is drawn as the current one.

```go
func NewChipButton(ui *UI, words string) *ChipButton
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the chip's words. |
| `Tone` | `Tone` | Tone is what its colour means. |
| `Strong` | `bool` | Strong fills the chip in its tone. |
| `Symbol` | `string` | Symbol is an optional meaning name before the words; TrailingSymbol one after them, for a caller that knows where pressing leads. |
| `TrailingSymbol` | `string` | Symbol is an optional meaning name before the words; TrailingSymbol one after them, for a caller that knows where pressing leads. |
| `Explanation` | `string` | Explanation says where pressing the chip goes, where the words do not. |
| `UnavailableReason` | `string` | UnavailableReason says why the chip cannot be pressed; "" means it can. |
| `Current` | `bool` | Current marks the chip whose destination is already open. |
| `Selectable` | `bool` | Selectable makes the chip a choice that is on or off; Selected is whether it is on. |
| `Selected` | `bool` | Selectable makes the chip a choice that is on or off; Selected is whether it is on. |
| `OnActivate` | `func()` | OnActivate runs when an available chip is pressed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Every tone, tinted, filled and keyboard-focused*

```go
func chipButtonTones(ui *kvitui.UI) unison.Paneler {
    chip := func(words string, tone kvitui.Tone) *kvitui.ChipButton {
        c := kvitui.NewChipButton(ui, words)
        c.Tone = tone
        return c
    }
    settled := chip("settled", kvitui.ToneSuccess)
    settled.Strong = true
    disputed := chip("disputed", kvitui.ToneDanger)
    disputed.Strong, disputed.Symbol = true, "warning"
    focused := chip("keyboard focus", kvitui.ToneNeutral)
    focused.Focus()
    return kvitui.Column(ui, kvitui.SizeSpace,
        kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, chip("neutral", kvitui.ToneNeutral), chip("accent", kvitui.ToneAccent),
            chip("success", kvitui.ToneSuccess), chip("warning", kvitui.ToneWarning), chip("danger", kvitui.ToneDanger),
            chip("info", kvitui.ToneInfo))),
        kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, settled, disputed, focused)))
}
```

*What pressing it leads to, where the caller knows*

```go
func chipButtonLeadsTo(ui *kvitui.UI) unison.Paneler {
    // No chevron is drawn for a chip that acts. Whether pressing it discloses
    // a list under the row, leaves the view or opens another window belongs
    // to the destination, which is why the symbol is the caller's to name.
    changes := kvitui.NewChipButton(ui, "3 changes")
    changes.Symbol = "file"
    ahead := kvitui.NewChipButton(ui, "1 ahead")
    ahead.TrailingSymbol = "chevron-right"
    behind := kvitui.NewChipButton(ui, "2 behind")
    behind.TrailingSymbol = "chevron-right"
    hub := kvitui.NewChipButton(ui, "Open on the hub")
    hub.TrailingSymbol = "external"
    return kvitui.Row(ui, kvitui.SizeSpaceNear, changes, ahead, behind, hub)
}
```

*Unavailable, with the reason attached*

```go
func chipButtonUnavailable(ui *kvitui.UI) unison.Paneler {
    // A chip with a reason keeps its place in the tab order, still shows the
    // words in its tooltip and its accessible description, and does nothing
    // when pressed. The fill going away rather than changing hue is what a
    // reader who cannot separate the tones still sees.
    behind := kvitui.NewChipButton(ui, "2 behind")
    behind.UnavailableReason = "The other branch has not been fetched yet."
    inSync := kvitui.NewChipButton(ui, "in sync")
    inSync.Tone, inSync.UnavailableReason = kvitui.ToneSuccess, "There is nothing on either side to compare."
    return kvitui.Row(ui, kvitui.SizeSpaceNear, kvitui.NewChipButton(ui, "1 ahead"), behind, inSync)
}
```

*Elided in a narrow column, at any interface size*

```go
func chipButtonElided(ui *kvitui.UI) unison.Paneler {
    // The label gives way and the symbols keep their size: a symbol at half
    // width is a smudge, and the trailing one says where pressing this goes.
    // The whole label is in the tooltip once it no longer fits.
    long := kvitui.NewChipButton(ui, "A fact whose whole phrase does not fit in this column")
    long.Symbol, long.TrailingSymbol = "warning", "chevron-right"
    short := kvitui.NewChipButton(ui, "A fact whose whole phrase does not fit in this column")
    return kvitui.Column(ui, kvitui.SizeSpaceSnug,
        kvitui.Width(ui, kvitui.Px(150), long), kvitui.Width(ui, kvitui.Px(90), short))
}
```

*The one that is already open*

```go
func chipButtonCurrent(ui *kvitui.UI) unison.Paneler {
    // One chip in a row of places to go is where the reader already is:
    // the selection tint under it, the accent on its edge, and bold words,
    // since the tints are nothing in a grayscale screenshot.
    changes := kvitui.NewChipButton(ui, "3 changes")
    changes.Symbol = "diff"
    ahead := kvitui.NewChipButton(ui, "1 ahead")
    ahead.Current, ahead.TrailingSymbol = true, "chevron-right"
    behind := kvitui.NewChipButton(ui, "2 behind")
    behind.TrailingSymbol = "chevron-right"
    return kvitui.Row(ui, kvitui.SizeSpaceNear, changes, ahead, behind)
}
```

*Where pressing one goes*

```go
func chipButtonExplanation(ui *kvitui.UI) unison.Paneler {
    // "1 ahead" is a count; what it opens is the one commit. Explanation is
    // that sentence; UnavailableReason is the other one, about a chip that
    // cannot be pressed at all.
    ahead := kvitui.NewChipButton(ui, "1 ahead")
    ahead.TrailingSymbol = "chevron-right"
    ahead.Explanation = "Opens the one commit this branch has and the remote does not."
    ahead.Focus()
    behind := kvitui.NewChipButton(ui, "2 behind")
    behind.UnavailableReason = "The remote has not been fetched yet."
    return kvitui.Row(ui, kvitui.SizeSpaceNear, ahead, behind)
}
```

### KvitStepper

A number with a minus and a plus beside it, for a small range a reader adjusts by one or two. The range comes from whatever owns it rather than being repeated at the call site.

```go
func NewStepper(ui *UI, label string, from, to int) *Stepper
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names what the number is, such as "Interface size". |
| `Unit` | `string` | Unit follows the number, in the muted colour. |
| `Value` | `int` | Value is the number. |
| `From` | `int` | From and To are the ends of the range; Step is how far a press moves. |
| `To` | `int` | From and To are the ends of the range; Step is how far a press moves. |
| `Step` | `int` | From and To are the ends of the range; Step is how far a press moves. |
| `OnChange` | `func(value int)` | OnChange runs after a press changes the value, with the new value. |
| `Follow` | `func() int` | Follow, when set, gives the value from whatever owns it, read before every layout, so the stepper shows the owner's value when it changes elsewhere, as a binding does in QML. |

*The interface-size row, pointed at the right setting*

```go
func stepperInterfaceSize(ui *kvitui.UI) unison.Paneler {
    size := kvitui.NewStepper(ui, "Interface size", tokens.MinInterfaceSize, tokens.MaxInterfaceSize)
    size.Unit = "px"
    size.Follow = ui.Interface.FontSize
    size.OnChange = ui.Interface.SetFontSize
    return size
}
```

### KvitField

A single line of text. The outline is the control-boundary token, so its edges are visible. An error is a message and a border together, never a border alone.

```go
func NewField(ui *UI) *Field
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the field for a screen reader. |
| `Placeholder` | `string` | Placeholder is shown, faint, while the field is empty. |
| `Error` | `string` | Error is a message under the field. |
| `Shortcut` | `string` | Shortcut names the keys that reach the field, where something does. |
| `ReadOnly` | `bool` | ReadOnly lets the text be selected and copied but not changed. |
| `OnChange` | `func(text string)` | OnChange runs after every change to the text. |

Methods: `Edit`, `Focus`, `SetEnabled`, `SetText`, `Text`.

*Resting, filled and in error*

```go
func fieldStates(ui *kvitui.UI) unison.Paneler {
    payee := kvitui.NewField(ui)
    payee.Label, payee.Placeholder = "Payee", "Who was paid"
    reference := kvitui.NewField(ui)
    reference.Label = "Reference"
    reference.SetText("TX-00173404")
    amount := kvitui.NewField(ui)
    amount.Label, amount.Error = "Amount", "Not a number"
    amount.SetText("twelve")
    locked := kvitui.NewField(ui)
    locked.Label = "Locked"
    locked.SetText("read only")
    locked.SetEnabled(false)
    return kvitui.Width(ui, kvitui.Px(200), kvitui.Column(ui, kvitui.SizeSpaceLoose, payee, reference, amount, locked))
}
```

### KvitTextArea

Several lines of text, typed or read. KvitField's outline and error rule in the field form, no ground at all in the plain one, which is for a document filling a pane. `Underlay` draws behind the words in the text's own coordinates, so a wash over a marked passage does not mean replacing the background.

```go
func NewTextArea(ui *UI) *TextArea
```

| Field | Type | |
|---|---|---|
| `Field` | embedded | Everything a Field has. |
| `Plain` | `bool` | Plain draws a document filling a pane: no ground and no outline, since a whole-pane rectangle drawn as a field says there is something beside it. |
| `Mono` | `bool` | Mono draws the text in the monospace family, for source text and anything where a column has to line up. |
| `Underlay` | `func(gc *unison.Canvas, at func(index int) geom.Rect)` | Underlay draws behind the words, such as a wash behind a marked passage. |

*Resting, filled, monospace and in error*

```go
func textAreaStates(ui *kvitui.UI) unison.Paneler {
    message := kvitui.NewTextArea(ui)
    message.Label, message.Placeholder = "Message", "What happened"
    source := kvitui.NewTextArea(ui)
    source.Label, source.Mono, source.ReadOnly = "Source", true, true
    source.SetText("# Notes\n\nThe file as it stands.")
    summary := kvitui.NewTextArea(ui)
    summary.Label, summary.Error = "Summary", "Say what changed"
    summary.SetText("...")
    return kvitui.Width(ui, kvitui.Px(280), kvitui.Column(ui, kvitui.SizeSpaceLoose, message, source, summary))
}
```

*Plain: a document filling a pane, with no ground of its own*

```go
func textAreaPlain(ui *kvitui.UI) unison.Paneler {
    source := kvitui.NewTextArea(ui)
    source.Label, source.Plain, source.Mono, source.ReadOnly = "Source", true, true, true
    source.SetText("# Notes\n\nA document has no edges to find:\nthe pane is its edge.")
    return kvitui.Width(ui, kvitui.Px(280), source)
}
```

*A wash behind a marked passage, drawn through `Underlay`*

```go
func textAreaUnderlay(ui *kvitui.UI) unison.Paneler {
    marked := kvitui.NewTextArea(ui)
    marked.Label, marked.Mono, marked.ReadOnly = "Reviewed source", true, true
    marked.SetText("one\ntwo\nthree")
    // A wash behind the second line, placed where the text says the fifth
    // character is.
    marked.Underlay = func(gc *unison.Canvas, at func(index int) geom.Rect) {
        box := at(4)
        box.Width = float32(ui.Interface.Px(40))
        gc.DrawRect(box, kvitui.Color(ui.Theme.Tokens().SelectionTint).Paint(gc, box, paintstyle.Fill))
    }
    return kvitui.Width(ui, kvitui.Px(280), marked)
}
```

### KvitSearchField

A field that filters something. Escape clears rather than reverting, and it announces its result count — filtering is the one interaction whose whole outcome happens somewhere else on the screen.

```go
func NewSearchField(ui *UI) *SearchField
```

| Field | Type | |
|---|---|---|
| `Field` | embedded | Everything a Field has. |
| `Matches` | `int` | Matches is how many things the filter left; below zero says nothing, for a filter whose result is not a countable list. |
| `MatchedNoun` | `string` | MatchedNoun is the word for one of them, "result" unless set. |
| `MatchedNounPlural` | `string` | MatchedNounPlural is the word for several; "" adds an s. |

Methods: `MatchPhrase`.

*Empty and filtering*

```go
func searchFieldStates(ui *kvitui.UI) unison.Paneler {
    filtering := kvitui.NewSearchField(ui)
    filtering.SetText("harlow")
    filtering.Matches, filtering.MatchedNoun = 47, "transaction"
    // What a screen reader is told is a sentence: the digits grouped by the
    // reader's locale and the plural a word, so "250,000 entries" rather than
    // "250000 entry(s)". A noun that does not take an s says its own plural.
    year := kvitui.NewSearchField(ui)
    year.SetText("2026")
    year.Matches, year.MatchedNoun, year.MatchedNounPlural = 250000, "entry", "entries"
    return kvitui.Width(ui, kvitui.Px(220), kvitui.Column(ui, kvitui.SizeSpaceLoose, kvitui.NewSearchField(ui), filtering, year))
}
```

### KvitCheck

A checkbox in three states. The third — partial — is what a parent row shows when some of its children are checked; drawing that as unchecked loses the information and drawing it as checked is a lie.

```go
func NewCheck(ui *UI, label string) *Check
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the box's label. |
| `Checked` | `bool` | Checked is the box's state. |
| `Partial` | `bool` | Partial says some of what the box stands for is checked; it is drawn with a dash, and a press checks the whole. |
| `OnChange` | `func(checked bool)` | OnChange runs after a press, with the new state. |
| `Name` | `string` | Name is what a screen reader calls a box with no words beside it, such as the one on each row of a table; the words unless set. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Off, on, partial and disabled*

```go
func checkStates(ui *kvitui.UI) unison.Paneler {
    drafts := kvitui.NewCheck(ui, "Include drafts")
    drafts.Checked = true
    some := kvitui.NewCheck(ui, "Some of these")
    some.Partial = true
    unavailable := kvitui.NewCheck(ui, "Not available")
    unavailable.SetEnabled(false)
    return kvitui.Column(ui, kvitui.SizeSpaceNear,
        kvitui.Left(kvitui.NewCheck(ui, "Include archived")), kvitui.Left(drafts), kvitui.Left(some), kvitui.Left(unavailable))
}
```

### KvitSelect

A choice from a list too long to lay out. For two to four self-evident options, KvitSegmented is right instead — a three-item dropdown hides two of the three answers for no reason.

```go
func NewSelect(ui *UI, label string, options ...Option) *Select
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the choice for a screen reader. |
| `Options` | `[]Option` | Options are the choices. |
| `Current` | `string` | Current is the Value of the chosen option. |
| `OnChoose` | `func(value string)` | OnChoose runs after an option is chosen, with its value. |

Methods: `Choose`, `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*A currency picker*

```go
func selectCurrency(ui *kvitui.UI) unison.Paneler {
    var options []kvitui.Option
    for _, code := range []string{"GBP", "EUR", "USD", "JPY", "CHF", "SEK"} {
        options = append(options, kvitui.Option{Value: code, Label: code})
    }
    return kvitui.NewSelect(ui, "Currency", options...)
}
```

### KvitTab

One tab in a row of them. The selected tab is marked by an underline as well as by colour and weight — selection shown by colour alone is the most common place the rule gets broken. `Explanation` is one sentence saying what the view behind the tab shows, which two or three words cannot.

```go
func NewTab(ui *UI, label string) *Tab
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text names the view. |
| `Selected` | `bool` | Selected marks the tab whose view is showing. |
| `Count` | `int` | Count is how many things the view holds, drawn beside the name; below zero draws none. |
| `Explanation` | `string` | Explanation says in one sentence what the view shows, which two or three words on a tab cannot. |
| `OnClick` | `func()` | OnClick runs when the tab is pressed. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*Selected, counted and plain*

```go
func tabForms(ui *kvitui.UI) unison.Paneler {
    all := kvitui.NewTab(ui, "All")
    all.Selected, all.Count = true, 1284
    uncategorised := kvitui.NewTab(ui, "Uncategorised")
    uncategorised.Count = 47
    return kvitui.Row(ui, kvitui.Px(0), all, uncategorised, kvitui.NewTab(ui, "Disputed"))
}
```

*Each tab saying what its view shows*

```go
func tabExplanations(ui *kvitui.UI) unison.Paneler {
    changes := kvitui.NewTab(ui, "Changes")
    changes.Selected = true
    changes.Explanation = "What this version changed."
    source := kvitui.NewTab(ui, "Source")
    source.Explanation = "The file’s own text."
    history := kvitui.NewTab(ui, "History")
    history.Explanation = "Every version and who wrote it."
    return kvitui.Row(ui, kvitui.Px(0), changes, source, history)
}
```

### KvitTooltip

A short label next to a control after a pause. Never the only place a control's meaning lives — a control explained only by its tooltip is unusable on a keyboard, a screen reader and a touch screen.

```go
func (u *UI) ShowTooltip(anchor unison.Paneler, words string) (hide func())
```

ShowTooltip shows a tooltip beside an anchor until the returned function is called, placed as a control's own tooltip is.

*On a button*

```go
func tooltipOnButton(ui *kvitui.UI) unison.Paneler {
    reconcile := kvitui.NewButton(ui, "Reconcile")
    // Shown here for good; a control shows its own on hover and on keyboard
    // focus.
    ui.ShowTooltip(reconcile, "Match these against the statement")
    return kvitui.FullWidth(kvitui.At(ui, kvitui.Px(60), kvitui.Px(40), kvitui.Px(80), reconcile))
}
```

### KvitPopover

A small surface anchored to a control holding something to act on. Takes focus, closes on Escape and on a click outside — which is what separates it from a hover card.

```go
func NewPopover(ui *UI, title string, content ...unison.Paneler) *Popover
```

| Field | Type | |
|---|---|---|
| `Title` | `string` | Title names the popover for a screen reader. |
| `Width` | `Measure` | Width holds the popover to a width, its height following from it; nil takes the width its content asks for. |
| `OnClose` | `func()` | OnClose runs after the popover closes. |

Methods: `Close`, `Open`, `Opened`.

*Open, holding a form*

```go
func popoverWithForm(ui *kvitui.UI) unison.Paneler {
    settled := kvitui.NewCheck(ui, "Settled")
    settled.Checked = true
    filter := kvitui.NewPopover(ui, "Filter", settled, kvitui.NewCheck(ui, "Pending"), kvitui.NewCheck(ui, "Disputed"))
    filter.Width = kvitui.Px(240)
    stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(150), unison.NewPanel()))
    filter.Open(stage, kvitui.PlaceOver(stage))
    return stage
}
```

### KvitHint

An information-icon trigger for an explanation too long for a tooltip. Click or keyboard activation keeps its KvitPopover open for reading, and Escape or an outside click dismisses it.

```go
func NewHint(ui *UI, label, explanation string) *Hint
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names what is explained; it is the button's name and the popover's title. |
| `Text` | `string` | Text is the explanation. |

Methods: `Close`, `Open`, `Opened`.

*Open for a longer explanation*

```go
func hintOpen(ui *kvitui.UI) unison.Paneler {
    hint := kvitui.NewHint(ui, "About automatic matching",
        "Automatic matching compares the date, amount and reference. It never changes the imported statement.")
    hint.Open()
    return kvitui.FullWidth(kvitui.At(ui, nil, nil, kvitui.Px(130), hint))
}
```

### KvitHoverCard

More about the thing under the pointer. Read-only, always: a surface that appears on hover is unreachable by keyboard and by touch, so nothing inside one may be the only route to an action.

```go
func NewHoverCard(ui *UI, content ...unison.Paneler) *HoverCard
```

*Showing what a row could not fit*

```go
func hoverCardRow(ui *kvitui.UI) unison.Paneler {
    name := kvitui.NewLabel(ui, "Payment to Harlow depot")
    name.Role = kvitui.RoleStrong
    when := kvitui.NewLabel(ui, "17 March 2026 at 09:14")
    when.Role, when.Ink = kvitui.RoleSmall, kvitui.InkTextMuted
    return kvitui.NewHoverCard(ui, name, when, kvitui.NewFigure(ui, "1,284.50", "GBP"))
}
```

### KvitToast

A message that appears, says one thing and goes away. For confirming what already happened. A toast with an action stays until dismissed — an undo that times out mid-read is an undo the reader cannot use.

```go
func NewToast(ui *UI, words string) *Toast
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the message. |
| `Tone` | `Tone` | Tone is info, success, warning or danger. |
| `Action` | `string` | Action is the words of the one action; "" for none. |
| `Timeout` | `time.Duration` | Timeout is how long a toast without an action stays; 4 s unless set. |
| `OnAction` | `func()` | OnAction runs when the action is taken; OnDismiss when the toast goes, by its timeout or its close control. |
| `OnDismiss` | `func()` | OnAction runs when the action is taken; OnDismiss when the toast goes, by its timeout or its close control. |

Methods: `Shown`.

*Four tones, and one with an undo*

```go
func toastTones(ui *kvitui.UI) unison.Paneler {
    toast := func(words string, tone kvitui.Tone, action string) unison.Paneler {
        t := kvitui.NewToast(ui, words)
        t.Tone, t.Action = tone, action
        return kvitui.Left(t)
    }
    return kvitui.Column(ui, kvitui.SizeSpace,
        toast("Copied to the clipboard", kvitui.ToneInfo, ""),
        toast("Statement imported", kvitui.ToneSuccess, ""),
        toast("Two rows could not be matched", kvitui.ToneWarning, ""),
        toast("The file could not be read", kvitui.ToneDanger, ""),
        toast("40 transactions archived", kvitui.ToneSuccess, "Undo"))
}
```

### KvitNotice

A message that stays: a condition that is still true, where a toast is an acknowledgement of something that finished. Dismissible only when dismissing it is meaningful.

```go
func NewNotice(ui *UI, condition string) *Notice
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is the condition; Detail, optional, says more. |
| `Detail` | `string` | Text is the condition; Detail, optional, says more. |
| `Tone` | `Tone` | Tone is info, success, warning or danger. |
| `Action` | `string` | Action is the words of the way out of the condition; "" for none. |
| `ActionHint` | `string` | ActionHint says what pressing the action does, where its words cannot (kvit-cash): the action's tooltip. |
| `Dismissible` | `bool` | Dismissible draws a close control. |
| `OnAction` | `func()` | OnAction and OnDismiss run when the action or the close control is pressed. |
| `OnDismiss` | `func()` | OnAction and OnDismiss run when the action or the close control is pressed. |

*A warning that can be dismissed, and an error that cannot*

```go
func noticeForms(ui *kvitui.UI) unison.Paneler {
    licence := kvitui.NewNotice(ui, "Your licence expires in three days")
    licence.Tone, licence.Detail = kvitui.ToneWarning, "Renew before 30 March to keep syncing."
    licence.Action, licence.Dismissible = "Renew", true
    vault := kvitui.NewNotice(ui, "This vault could not be saved")
    vault.Tone, vault.Detail = kvitui.ToneDanger, "The disk is full. Nothing has been lost; the changes are still held."
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, licence, vault))
}
```

### KvitDialog

A modal surface that has to be answered first. Expensive, and worth it only where the answer really does come first; the confirming button is on the right, and a destructive dialog takes no keyboard default. There are three ways out, because different readers find different ones: the close symbol at the right of the title, the cancel button, and Escape.

```go
func NewDialog(ui *UI, title string, content ...unison.Paneler) *Dialog
```

| Field | Type | |
|---|---|---|
| `Title` | `string` | Title asks the question; Detail says more. |
| `Detail` | `string` | Title asks the question; Detail says more. |
| `ConfirmText` | `string` | ConfirmText and CancelText are the buttons' words; "" leaves a button out, and a dialog with neither has no row of buttons. |
| `CancelText` | `string` | ConfirmText and CancelText are the buttons' words; "" leaves a button out, and a dialog with neither has no row of buttons. |
| `Destructive` | `bool` | Destructive draws the confirming button in the danger colour and takes away the keyboard default. |
| `OnAccept` | `func()` | OnAccept and OnReject run when the dialog is answered. |
| `OnReject` | `func()` | OnAccept and OnReject run when the dialog is answered. |

Methods: `Open`.

*A destructive confirmation*

```go
func dialogDestructive(ui *kvitui.UI) unison.Paneler {
    // Shown here without its modality so it sits on the page. A real one is
    // opened with Open, modal and in the middle of the window.
    confirm := kvitui.NewDialog(ui, "Delete 40 transactions?")
    confirm.Detail = "They will be removed from every report. This cannot be undone."
    confirm.ConfirmText, confirm.Destructive = "Delete them", true
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(210), kvitui.Centred(kvitui.Width(ui, kvitui.Px(420), confirm))))
}
```

### KvitFloatingView

Detail about the one thing selected, held over the list rather than beside it: a centred card with a name and a close control, and the list dimmed behind. Unlike KvitPane it takes no width out of the layout, so the list keeps the whole window at every size; unlike KvitDialog it is closed rather than answered, and it floats over one region rather than over the window. From kvit-cash's copy of the Qt library.

```go
func NewFloatingView(ui *UI, title string, body unison.Paneler) *FloatingView
```

| Field | Type | |
|---|---|---|
| `Title` | `string` | Title names what is in the view. |
| `TitleDetail` | `string` | TitleDetail is the whole of what the title names where the title is shortened, shown under the pointer on the title; "" says the title is the whole of it. |
| `CloseLabel` | `string` | CloseLabel is what the close control is called on hover and to a screen reader: "Close record" rather than "Close the view", since the reader is closing the record. |
| `Closable` | `bool` | Closable draws the close control. |
| `ViewWidth` | `Measure` | ViewWidth is the card's width, at most the region's less the margins; the floating view width unless set. |
| `CloseOnPressOutside` | `bool` | CloseOnPressOutside closes the view on a press on the dimmed area. |
| `OnCloseRequested` | `func()` | OnCloseRequested runs when the reader asks to close the view, by the close control, Escape, or a press outside where that is on. |

Methods: `Close`, `Open`, `Opened`.

*A record held over the list it came from*

```go
func floatingRecord(ui *kvitui.UI) unison.Paneler {
    // Shown over a stand-in list, which is what a floating view is always
    // over: the card is centred on the region it floats in, and the region is
    // the list's, not the window's.
    row := func(name, figure string) unison.Paneler {
        r := kvitui.NewSlimRow(ui, name)
        r.Figure = figure
        return kvitui.FullWidth(r)
    }
    list := kvitui.Column(ui, kvitui.Px(0),
        row("Whole Foods", "-84.10"), row("Rent", "-1,850.00"), row("Payroll", "3,204.55"), row("Con Edison", "-96.22"))
    stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), list))
    detail := func(name, phrase string) unison.Paneler {
        r := kvitui.NewSlimRow(ui, name)
        r.Phrase = phrase
        return kvitui.FullWidth(r)
    }
    body := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.SizeSpace,
        kvitui.FullWidth(kvitui.NewSectionHeading(ui, "Source record")),
        detail("Payee", "Whole Foods"), detail("Category", "Groceries")))
    record := kvitui.NewFloatingView(ui, "Whole Foods · 3 Sep", body)
    record.CloseLabel = "Close record"
    record.OnCloseRequested = record.Close
    whenShown(stage, func() { record.Open(stage) })
    return stage
}
```

### KvitBar

One quantity against a stated scale. A bar with no value draws a tick rather than a zero-width fill; a bounded figure is hatched, and the hatch survives grayscale; and the scale is required, because a bar drawn against its own list rescales invisibly.

```go
func NewBar(ui *UI, value, maximum float64) *Bar
```

| Field | Type | |
|---|---|---|
| `Value` | `float64` | Value against Maximum, the scale. |
| `Maximum` | `float64` | Value against Maximum, the scale. |
| `Measured` | `bool` | Measured is false for a value nobody measured. |
| `Bounded` | `bool` | Bounded marks the value as "at most this much". |
| `Ink` | `Ink` | Ink is the fill; the accent unless set. |
| `Wide` | `bool` | Wide draws the taller form. |
| `Label` | `string` | Label and Unit say what the bar is, for a screen reader. |
| `Unit` | `string` | Label and Unit say what the bar is, for a screen reader. |

*Measured, bounded and unmeasured*

```go
func barForms(ui *kvitui.UI) unison.Paneler {
    attention := kvitui.NewBar(ui, 62, 100)
    attention.Label, attention.Unit = "Attention", "h"
    agent := kvitui.NewBar(ui, 38, 100)
    agent.Wide, agent.Bounded, agent.Ink = true, true, kvitui.InkAxisAgent
    agent.Label, agent.Unit = "Agent", "h"
    unrecorded := kvitui.NewBar(ui, 0, 100)
    unrecorded.Measured, unrecorded.Label = false, "Unrecorded"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose, attention, agent, unrecorded))
}
```

### KvitStackedBar

Several quantities adding to one total. The two-pixel gap is the reason this is a component: adjacent fills that touch read as one fill with a colour change in it.

```go
func NewStackedBar(ui *UI, segments ...Segment) *StackedBar
```

| Field | Type | |
|---|---|---|
| `Segments` | `[]Segment` | Segments are the parts, left to right. |
| `Maximum` | `float64` | Maximum is the scale; the sum of the segments unless set, which is right for a bar that is a breakdown of its own total. |
| `Wide` | `bool` | Wide draws the taller form. |
| `Label` | `string` | Label names the bar for a screen reader. |

*A four-way breakdown*

```go
func stackedBarBreakdown(ui *kvitui.UI) unison.Paneler {
    month := kvitui.NewStackedBar(ui,
        kvitui.Segment{Value: 420, Label: "rent", Ink: kvitui.InkCategorical(0)},
        kvitui.Segment{Value: 180, Label: "groceries", Ink: kvitui.InkCategorical(1)},
        kvitui.Segment{Value: 95, Label: "transport", Ink: kvitui.InkCategorical(2)},
        kvitui.Segment{Value: 60, Label: "everything else", Ink: kvitui.InkCategorical(3)})
    month.Wide, month.Label = true, "Where the month went"
    return kvitui.FullWidth(month)
}
```

### KvitSpark

The shape of a series, small enough to sit in a row. A period that was never measured is a baseline tick rather than a zero-height bar: a missing week and a quiet week are different facts.

```go
func NewSpark(ui *UI, values ...float64) *Spark
```

| Field | Type | |
|---|---|---|
| `Values` | `[]float64` | Values are the periods, oldest first; NotMeasured for a gap. |
| `Ink` | `Ink` | Ink is the bars' colour; the accent unless set. |
| `Maximum` | `float64` | Maximum is the scale; the largest value unless set. |
| `Label` | `string` | Label names the series for a screen reader. |

*With a hole in the middle*

```go
func sparkWithHole(ui *kvitui.UI) unison.Paneler {
    gap := kvitui.NotMeasured
    notes := kvitui.NewSpark(ui, 3, 5, 8, 6, gap, gap, 9, 12, 7, 4, 6, 11)
    notes.Label = "Notes written"
    return kvitui.Width(ui, kvitui.Px(160), notes)
}
```

### KvitNetFlow

Two opposed series over the same periods, on one baseline and one scale, with what they come to drawn as a line. What comes in is drawn up from the baseline and what goes out is drawn down, so which of the two was larger in a period is a glance rather than a comparison between two strips an inch apart. Each column names its own period underneath, and the names thin out to whatever spacing they actually fit in. From kvit-cash's copy of the Qt library.

```go
func NewNetFlow(ui *UI, ins, outs []float64, periods ...string) *NetFlow
```

| Field | Type | |
|---|---|---|
| `Ins` | `[]float64` | Ins and Outs are what came in and went out over each period, indexed together; NotMeasured is a period nobody measured, drawn as a tick on the baseline rather than as a zero. |
| `Outs` | `[]float64` | Ins and Outs are what came in and went out over each period, indexed together; NotMeasured is a period nobody measured, drawn as a tick on the baseline rather than as a zero. |
| `Periods` | `[]string` | Periods name the columns, one each, in the same order. |
| `Maximum` | `float64` | Maximum is the scale both directions are drawn against; 0 takes the largest value present. |
| `RisingLabel` | `string` | RisingLabel, FallingLabel and NetLabel name the parts in a key above the plot; all empty draws no key. |
| `FallingLabel` | `string` | RisingLabel, FallingLabel and NetLabel name the parts in a key above the plot; all empty draws no key. |
| `NetLabel` | `string` | RisingLabel, FallingLabel and NetLabel name the parts in a key above the plot; all empty draws no key. |
| `RisingInk` | `Ink` | RisingInk and FallingInk are the bars' colours: the first two categorical colours unless set, never a meaning colour by default, since success and danger mean finished and stalled everywhere else. |
| `FallingInk` | `Ink` | RisingInk and FallingInk are the bars' colours: the first two categorical colours unless set, never a meaning colour by default, since success and danger mean finished and stalled everywhere else. |
| `NetInk` | `Ink` | RisingInk and FallingInk are the bars' colours: the first two categorical colours unless set, never a meaning colour by default, since success and danger mean finished and stalled everywhere else. |
| `Label` | `string` | Label names the chart for a screen reader. |
| `PlotHeight` | `Measure` | PlotHeight is the plot's height without the key and the names; 64 design pixels unless set. |

*A year of money in and money out, netted off*

```go
func netFlowYear(ui *kvitui.UI) unison.Paneler {
    flow := kvitui.NewNetFlow(ui,
        []float64{0.62, 0.55, 0.71, 0.58, 0.64, 0.60, 0.66, 0.59, 0.63, 0.70, 0.61, 0.65},
        []float64{0.48, 0.72, 0.51, 0.66, 0.44, 0.81, 0.52, 0.47, 0.69, 0.55, 0.90, 0.50},
        "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec")
    flow.Maximum = 1
    flow.RisingLabel, flow.FallingLabel, flow.NetLabel = "Money in", "Money out", "Kept"
    // Money in and money out are good and bad news, so this chart says so
    // with the meaning colours; a chart's default is the categorical pair.
    flow.RisingInk, flow.FallingInk = kvitui.InkSuccess, kvitui.InkDanger
    flow.Label = "Cash flow by month"
    return kvitui.FullWidth(flow)
}
```

*Daily columns, where the names thin out, and a period nobody measured*

```go
func netFlowDays(ui *kvitui.UI) unison.Paneler {
    gap := kvitui.NotMeasured
    var days []string
    for d := 1; d <= 20; d++ {
        days = append(days, strconv.Itoa(d))
    }
    flow := kvitui.NewNetFlow(ui,
        []float64{0.2, 0.0, 0.9, 0.1, gap, 0.3, 0.0, 0.5, 0.2, 0.0, 0.7, 0.1, 0.0, 0.4, 0.0, 0.8, 0.2, 0.0, 0.3, 0.6},
        []float64{0.3, 0.4, 0.2, 0.5, gap, 0.1, 0.6, 0.2, 0.7, 0.3, 0.2, 0.4, 0.5, 0.1, 0.3, 0.2, 0.6, 0.4, 0.2, 0.3},
        days...)
    flow.Maximum, flow.Label = 1, "Cash flow by day"
    return kvitui.FullWidth(flow)
}
```

### KvitTrend

A series with a value axis and a hover crosshair — read for values, where a spark is read for shape. A gap in the data draws as a gap: interpolating over a hole asserts values nobody measured. A second series is dashed as well as differently coloured, and both are named in the key and in the crosshair.

```go
func NewTrend(ui *UI, points ...float64) *Trend
```

| Field | Type | |
|---|---|---|
| `Points` | `[]float64` | Points are the values, oldest first; NotMeasured for a gap. |
| `MinimumY` | `float64` | MinimumY and MaximumY are the scale's ends. |
| `MaximumY` | `float64` | MinimumY and MaximumY are the scale's ends. |
| `Ink` | `Ink` | Ink is the line's colour; the accent unless set. |
| `Label` | `string` | Label and Unit say what the line is. |
| `Unit` | `string` | Label and Unit say what the line is. |
| `Gridlines` | `int` | Gridlines is how many horizontal lines, at least two; 3 unless set. |
| `Second` | `[]float64` | Second is an optional second series, dashed, in SecondInk (the second categorical colour unless set), named SecondLabel. |
| `SecondInk` | `Ink` |  |
| `SecondLabel` | `string` |  |
| `Overlay` | `func(gc *unison.Canvas, plot geom.Rect)` | Overlay draws over the plot after the lines, given the plot's box, so an annotation lines up with the same arithmetic the lines use. |

*A series, and the empty state*

```go
func trendAndEmpty(ui *kvitui.UI) unison.Paneler {
    balance := kvitui.NewTrend(ui, 400, 620, 580, 900, kvitui.NotMeasured, 1400, 1250, 1700)
    balance.Label, balance.Unit, balance.MaximumY = "Balance", "GBP", 2000
    savings := kvitui.NewTrend(ui)
    savings.Label = "Savings"
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
        balance, kvitui.Height(ui, kvitui.Px(90), savings)))
}
```

*Two series, with an annotation drawn over the plot*

```go
func trendTwoSeries(ui *kvitui.UI) unison.Paneler {
    gap := kvitui.NotMeasured
    worth := kvitui.NewTrend(ui, 400, 620, 580, 900, 1150, 1400, 1250, 1700)
    worth.Label, worth.SecondLabel, worth.Unit, worth.MaximumY = "Assets", "Liabilities", "GBP", 2000
    worth.Second = []float64{gap, gap, 300, 340, 320, 290, 260, 240}
    // Where the second account's history begins, drawn over the plot with
    // the same arithmetic the lines use.
    worth.Overlay = func(gc *unison.Canvas, plot geom.Rect) {
        x := plot.X + plot.Width*2/7
        mark := geom.NewRect(x, plot.Y, float32(ui.Interface.Hairline()), plot.Height)
        gc.DrawRect(mark, kvitui.Color(kvitui.InkMarker.Of(ui)).Paint(gc, mark, paintstyle.Fill))
        note := ui.Fonts.Layout(kvitui.Caption(ui, "history starts here", kvitui.InkTextFaint), kvitui.NoOptions)
        _, h := note.Size()
        note.Draw(gc, x+float32(ui.Interface.SpaceSnug()), plot.Bottom()-h)
    }
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(150), worth))
}
```

### KvitDistribution

How a set of values is spread. An average of four days and an average made of one twenty-day outlier are the same number and different situations.

```go
func NewDistribution(ui *UI) *Distribution
```

| Field | Type | |
|---|---|---|
| `Minimum` | `float64` | The five numbers. |
| `LowerQuartile` | `float64` | The five numbers. |
| `Median` | `float64` | The five numbers. |
| `UpperQuartile` | `float64` | The five numbers. |
| `Maximum` | `float64` | The five numbers. |
| `ScaleMinimum` | `float64` | ScaleMinimum and ScaleMaximum are the scale's ends. |
| `ScaleMaximum` | `float64` | ScaleMinimum and ScaleMaximum are the scale's ends. |
| `Measured` | `bool` | Measured is false for a spread nobody measured, drawn as "—". |
| `Ink` | `Ink` | Ink is the box's colour; the accent unless set. |
| `Label` | `string` | Label and Unit say what it is, for a screen reader. |
| `Unit` | `string` | Label and Unit say what it is, for a screen reader. |

*Two rows on one scale*

```go
func distributionRows(ui *kvitui.UI) unison.Paneler {
    card := kvitui.NewDistribution(ui)
    card.Label, card.Unit, card.ScaleMaximum = "Time to settle, card", "d", 30
    card.Minimum, card.LowerQuartile, card.Median, card.UpperQuartile, card.Maximum = 1, 2, 3, 5, 21
    transfer := kvitui.NewDistribution(ui)
    transfer.Label, transfer.Unit, transfer.ScaleMaximum = "Time to settle, transfer", "d", 30
    transfer.Minimum, transfer.LowerQuartile, transfer.Median, transfer.UpperQuartile, transfer.Maximum = 1, 1, 2, 2, 4
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose, card, transfer))
}
```

### KvitGauge

How much of an allowance is used, with the target marked and a pace mark saying where an even rate would have reached. Sixty percent spent is fine on day eighteen and a problem on day six.

```go
func NewGauge(ui *UI, value, allowance float64) *Gauge
```

| Field | Type | |
|---|---|---|
| `Value` | `float64` | Value against Allowance. |
| `Allowance` | `float64` | Value against Allowance. |
| `Pace` | `float64` | Pace is how far through the period it is, 0 to 1, drawn as a tick; below zero draws none. |
| `Measured` | `bool` | Measured is false for a value nobody measured. |
| `Label` | `string` | Label and Unit say what it is, for a screen reader. |
| `Unit` | `string` | Label and Unit say what it is, for a screen reader. |

*Ahead of pace, behind it, and over*

```go
func gaugePace(ui *kvitui.UI) unison.Paneler {
    gauge := func(value float64, label string) unison.Paneler {
        g := kvitui.NewGauge(ui, value, 600)
        g.Pace, g.Label, g.Unit = 0.6, label, "GBP"
        return g
    }
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
        gauge(240, "Groceries"), gauge(480, "Leisure"), gauge(720, "Transport")))
}
```

### KvitDelta

How much something changed and in which direction, with an arrow as well as a colour. Whether up is good is the caller's to say: a rise in spending and a rise in savings are the same arrow and opposite colours.

```go
func NewDelta(ui *UI, change float64) *Delta
```

| Field | Type | |
|---|---|---|
| `Change` | `float64` | Change is the change. |
| `Unit` | `string` | Unit follows it. |
| `Measured` | `bool` | Measured is false for a change nobody measured. |
| `Good` | `Direction` | Good is the direction that is good; Neither draws the muted colour. |
| `Precision` | `int` | Precision is how many decimal places to draw. |

*Up, down and unchanged*

```go
func deltaDirections(ui *kvitui.UI) unison.Paneler {
    delta := func(change float64, unit string, precision int, good kvitui.Direction) *kvitui.Delta {
        d := kvitui.NewDelta(ui, change)
        d.Unit, d.Precision, d.Good = unit, precision, good
        return d
    }
    unmeasured := kvitui.NewDelta(ui, 0)
    unmeasured.Measured = false
    return kvitui.Row(ui, kvitui.SizeSpaceLoose,
        delta(12.4, "%", 1, kvitui.Up), delta(-3.2, "%", 1, kvitui.Up), delta(18, "GBP", 0, kvitui.Down),
        delta(0, "", 0, kvitui.Neither), unmeasured)
}
```

### KvitStatTile

A card carrying one figure, what it is, how it changed and its recent shape — the shape of six of kvit-cash's seven dashboard widgets. The order is fixed so a row of tiles can be scanned one part at a time.

```go
func NewStatTile(ui *UI, label, value, unit string) *StatTile
```

| Field | Type | |
|---|---|---|
| `Card` | embedded | Everything a Card has. |
| `Label` | `string` | Label names the figure; Value and Unit are the figure. |
| `Value` | `string` | Label names the figure; Value and Unit are the figure. |
| `Unit` | `string` | Label names the figure; Value and Unit are the figure. |
| `Measured` | `bool` | Measured is false for a figure nobody measured. |
| `HasChange` | `bool` | HasChange draws Change with the Good direction. |
| `Change` | `float64` |  |
| `Good` | `Direction` |  |
| `History` | `[]float64` | History is the recent shape; NotMeasured for a gap. |
| `Caption` | `string` | Caption is a line under it all. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*A dashboard row*

```go
func statTileRow(ui *kvitui.UI) unison.Paneler {
    balance := kvitui.NewStatTile(ui, "Balance", "4,182.30", "GBP")
    balance.HasChange, balance.Change, balance.Good = true, 240.10, kvitui.Up
    balance.History = []float64{3200, 3400, 3390, 3800, 4000, 4182}
    spent := kvitui.NewStatTile(ui, "Spent this month", "812.40", "GBP")
    spent.HasChange, spent.Change, spent.Good = true, 96.20, kvitui.Down
    spent.Caption = "12 days remaining"
    uncategorised := kvitui.NewStatTile(ui, "Uncategorised", "", "")
    uncategorised.Measured = false
    return kvitui.Row(ui, kvitui.SizeColumnGap, balance, spent, uncategorised)
}
```

### KvitFigureBlock

A figure with its name under it: one number a reader is meant to take away. The number is above and larger, because a row of these is read across the numbers.

```go
func NewFigureBlock(ui *UI, value, unit, label string) *FigureBlock
```

| Field | Type | |
|---|---|---|
| `Figure` | `*Figure` |  |
| `Label` | `string` | Label names the figure. |

*A row of three*

```go
func figureBlockRow(ui *kvitui.UI) unison.Paneler {
    reconciled := kvitui.NewFigureBlock(ui, "", "", "Reconciled")
    reconciled.Figure.Measured = false
    return kvitui.Row(ui, kvitui.Px(40),
        kvitui.NewFigureBlock(ui, "1,284", "", "Transactions"), kvitui.NewFigureBlock(ui, "97", "%", "Categorised"), reconciled)
}
```

### KvitCell

One cell of a table, drawn according to what kind of value its column holds. The kind comes from the column, so every cell in it aligns the same way and says the same thing about a missing value. A cell formats nothing: Money and Figure are given the string they draw.

```go
func NewCell(ui *UI, kind CellKind, v CellValue) *Cell
```

| Field | Type | |
|---|---|---|
| `Kind` | `CellKind` | Kind is how the cell is drawn. |
| `Value` | `CellValue` | Value is what it holds. |
| `Leads` | `bool` | Leads draws the text in the link colour: the row's name, which the reader presses to open the row. |
| `Selected` | `bool` | Selected draws the text in the primary colour. |
| `Current` | `bool` | Current says the keyboard cursor is on the cell, which discloses a value cut short as the pointer does. |
| `OnToggle` | `func(checked bool)` | OnToggle asks for a Check cell's box to change. |

*The six kinds a column of values can be*

```go
func cellKinds(ui *kvitui.UI) unison.Paneler {
    cell := func(kind kvitui.CellKind, v kvitui.CellValue) unison.Paneler {
        return kvitui.FullWidth(kvitui.NewCell(ui, kind, v))
    }
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceTight,
        cell(kvitui.CellText, kvitui.CellValue{Text: "Payment to Harlow depot"}),
        cell(kvitui.CellDate, kvitui.CellValue{Date: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)}),
        cell(kvitui.CellSlug, kvitui.CellValue{Text: "TX-00173404"}),
        cell(kvitui.CellFigure, kvitui.CellValue{Text: "1284.5", Unit: "GBP"}),
        // Money takes an amount that has already been formatted. The cell
        // does no arithmetic: an amount is a count of the minor unit of its
        // own currency, and how many minor digits it has is the currency's.
        cell(kvitui.CellMoney, kvitui.CellValue{Text: "-42.90", Unit: "GBP"}),
        cell(kvitui.CellChip, kvitui.CellValue{Text: "Settled", Tone: kvitui.ToneSuccess}),
        cell(kvitui.CellFigure, kvitui.CellValue{Unmeasured: true})))
}
```

*Several states at once, and a row the reader picks*

```go
func cellStatesAndPick(ui *kvitui.UI) unison.Paneler {
    // Every state a row is in at once, one dot each. Each mark has a tone,
    // the shape that says the same thing without colour, and the word that
    // is both the tooltip and what a screen reader says.
    marks := kvitui.NewCell(ui, kvitui.CellMarks, kvitui.CellValue{Marks: []kvitui.CellMark{
        {Tone: kvitui.ToneSuccess, Shape: kvitui.ShapeCircle, Label: "Settled"},
        {Tone: kvitui.ToneWarning, Shape: kvitui.ShapeDiamond, Label: "Not reviewed"},
        {Tone: kvitui.ToneNeutral, Shape: kvitui.ShapeSquare, Label: "Has an attachment"},
    }})
    // The box does not tick itself. Whatever owns the selection does, and
    // the cell redraws from it, so a request that is refused leaves no tick.
    picked := kvitui.NewCell(ui, kvitui.CellCheck, kvitui.CellValue{Checked: true})
    picked.OnToggle = func(wanted bool) {
        picked.Value.Checked = wanted
        picked.MarkForLayoutAndRedraw()
    }
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceTight, kvitui.FullWidth(marks), kvitui.FullWidth(picked)))
}
```

*A value the column cut short*

```go
func cellCutShort(ui *kvitui.UI) unison.Paneler {
    // A value wider than its column. The model gives the whole of it as
    // FullText, and the cell shows it under the pointer and under the
    // keyboard cursor, which the view says is here through Current.
    return kvitui.Width(ui, kvitui.Px(150), kvitui.NewCell(ui, kvitui.CellText, kvitui.CellValue{
        Text:     "Harlow depot retainer",
        FullText: "Harlow depot — quarterly maintenance retainer",
    }))
}
```

### KvitTable

A dense, configurable table over a model. Columns sort, resize, move and open a menu from their own header, by keyboard as well as pointer; a row opens its record on one press; and a Check column draws a box per row with a box in its header for every row shown. The model says what each column is and answers each cell with its value, its unit, its marks, whether its box is ticked and the whole of a value the column cut short, which the cell then discloses on hover and under the keyboard cursor. Holds smooth scrolling and sub-100 ms filtering at 250,000 rows with twelve columns of value.

```go
func NewTable(ui *UI, model TableModel) *Table
```

| Field | Type | |
|---|---|---|
| `Model` | `TableModel` | Model is what the table shows. |
| `Label` | `string` | Label names the table for a screen reader; "Table" unless set. |
| `EmptyTitle` | `string` | EmptyTitle and EmptyDetail are what the table says when it has no rows, which is a thing to say rather than a blank grid to leave. |
| `EmptyDetail` | `string` | EmptyTitle and EmptyDetail are what the table says when it has no rows, which is a thing to say rather than a blank grid to leave. |
| `HiddenColumns` | `[]int` | HiddenColumns are the model columns the reader has hidden. |
| `Order` | `[]int` | Order is the model columns in the order they are shown; nil is the model's order. |
| `Widths` | `map[int]int` | Widths are the columns the reader has resized, by model column, in design pixels. |
| `SortColumn` | `int` | SortColumn is the column the rows are ordered by, -1 for none, and SortAscending its direction. |
| `SortAscending` | `bool` |  |
| `HeaderChecked` | `bool` | HeaderChecked and HeaderPartial are what the box in a Check column's header shows: every row shown, or some of them. |
| `HeaderPartial` | `bool` | HeaderChecked and HeaderPartial are what the box in a Check column's header shows: every row shown, or some of them. |
| `RowsOpen` | `bool` | RowsOpen draws, in a strip of its own down the right-hand edge, the mark that says a press on a row opens something. |
| `OpensColumn` | `int` | OpensColumn is the column holding the row's name, drawn in the link colour to say the row opens something; -1 for none. |
| `OpensLabel` | `string` | OpensLabel says in words what pressing a row opens, on the strip's tooltip. |
| `OnRowPressed` | `func(row int)` | OnRowPressed runs on one press on a row, which is how a ledger row opens its record; OnRowActivated on a double press and on Return. |
| `OnRowActivated` | `func(row int)` | OnRowPressed runs on one press on a row, which is how a ledger row opens its record; OnRowActivated on a double press and on Return. |
| `OnSort` | `func(column int, ascending bool)` | OnSort runs when the reader asks for a different order. |
| `OnHeaderToggled` | `func(checked bool)` | OnHeaderToggled runs when the box in a Check column's header is pressed, with the state asked for. |
| `OnCellToggled` | `func(row, column int, checked bool)` | OnCellToggled runs when one row's box is pressed, with the state asked for. |

Methods: `CopySelection`, `CurrentRow`, `FocusRow`, `HideColumn`, `IsRowSelected`, `MoveColumnBy`, `Position`, `Refresh`, `Reset`, `ResizeColumnBy`, `RowCount`, `ScrollTo`, `SelectedRows`, `SortBy`.

*Two hundred and fifty thousand rows*

```go
func tableQuarterMillion(ui *kvitui.UI) unison.Paneler {
    rows := kvitui.NewBenchmarkTableModel(kvitui.BenchmarkRows)
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), kvitui.NewTable(ui, rows)))
}
```

*Sorting, a column menu and a column of boxes*

```go
func tableSortingAndBoxes(ui *kvitui.UI) unison.Paneler {
    rows := kvitui.NewBenchmarkTableModel(400)
    table := kvitui.NewTable(ui, rows)
    // Six of the thirteen columns, so the boxes at the right-hand end are on
    // screen without scrolling.
    table.HiddenColumns = []int{1, 2, 5, 8, 10, 11}
    // The table draws the indicator and says what was asked for; putting the
    // rows in that order is the application's, because the model is. Nothing
    // here sorts, so the arrow stays where it was set.
    table.SortColumn, table.SortAscending = 0, false
    // The column of boxes, and the box in its header for every row shown.
    // The ticks are the model's, kept by the underlying row, so narrowing
    // the filter does not lose them.
    ticked := func() {
        table.HeaderChecked, table.HeaderPartial = rows.AllShownChecked(), rows.SomeShownChecked()
        table.Refresh()
    }
    table.OnHeaderToggled = func(wanted bool) { rows.SetEveryShownChecked(wanted); ticked() }
    table.OnCellToggled = func(row, _ int, wanted bool) { rows.SetChecked(row, wanted); ticked() }
    // The mark that says a press on a row opens something, drawn at rest in
    // a strip of its own rather than in a column, so a reader can tell which
    // lists lead anywhere without sweeping the pointer across them.
    table.RowsOpen = true
    table.OnRowPressed = func(row int) { rows.SetChecked(row, true); ticked() }
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), table))
}
```

*Nothing to show*

```go
func tableEmpty(ui *kvitui.UI) unison.Paneler {
    table := kvitui.NewTable(ui, kvitui.NewBenchmarkTableModel(0))
    table.EmptyTitle, table.EmptyDetail = "No transactions", "Nothing matches the current filter."
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(200), table))
}
```

### KvitScrollBar

A scroll bar, occupying its own strip rather than floating over content. Thirty-two files in the estate have a private version; an overlay bar hides the right-hand column of a table and the last character of every elided label.

```go
func NewScrollBar(ui *UI, horizontal bool) *unison.ScrollBar
```

NewScrollBar returns a scroll bar that occupies a strip of its own rather than floating over what it scrolls.

*Beside a scrolling column*

```go
func scrollBarBesideColumn(ui *kvitui.UI) unison.Paneler {
    rows := make([]unison.Paneler, 10)
    for i := range rows {
        rows[i] = kvitui.NewSlimRow(ui, fmt.Sprintf("Row %d", i+1))
    }
    region := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), rows...))
    return kvitui.Sized(ui, kvitui.Px(480), kvitui.Px(120), region)
}
```

### KvitMenu

A list of commands, opened with UI.ShowMenu. It is unison's menu in the Kvit colours and type, as the owner chose: it answers the arrow keys, Return and Escape, and teaches each shortcut by showing it beside its command.

```go
func (u *UI) ShowMenu(anchor unison.Paneler, title string, items []MenuItem) (closeMenu func())
```

ShowMenu opens a menu of items under anchor.

```go
func (u *UI) ShowMenuAt(owner unison.Paneler, part geom.Rect, title string, items []MenuItem) (closeMenu func())
```

ShowMenuAt opens a menu of items under one part of owner, given in owner's own coordinates, such as one column of a table's header, and returns the function that closes it.

*Opened, with shortcuts and a destructive item*

```go
func menuOpened(ui *kvitui.UI) unison.Paneler {
    stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(170), unison.NewPanel()))
    // Opened at the stage's top left once it is in a window, as the Qt page
    // opens its menu.
    whenShown(stage, func() {
        ui.ShowMenuAt(stage, geom.Rect{}, "", []kvitui.MenuItem{
            {Text: "Open", Key: command(unison.KeyO)},
            {Text: "Duplicate", Key: command(unison.KeyD)},
            {Text: "Archive"},
            {Separator: true},
            {Text: "Delete"},
        })
    })
    return stage
}
```

### KvitMenuItem

One line of a menu, carrying its shortcut on the right — which is how a menu teaches a faster route to a reader who keeps using it.

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is what the line says. |
| `Key` | `unison.KeyBinding` | Key is the shortcut shown at the right of the line, which is how a menu teaches that there is a faster way; the zero value shows none. |
| `Checked` | `bool` | Checked marks the entry in use. |
| `Disabled` | `bool` | Disabled draws the line and does nothing when chosen. |
| `Separator` | `bool` | Separator makes the line a divider instead, as goes before a destructive entry at the bottom of a menu. |
| `OnSelect` | `func()` | OnSelect runs when the line is chosen. |

*Ordinary, disabled and destructive*

```go
func menuItemForms(ui *kvitui.UI) unison.Paneler {
    stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(110), unison.NewPanel()))
    whenShown(stage, func() {
        ui.ShowMenuAt(stage, geom.Rect{}, "", []kvitui.MenuItem{
            {Text: "Reconcile", Key: command(unison.KeyR)},
            {Text: "Split", Disabled: true},
            {Text: "Delete"},
        })
    })
    return stage
}
```

### KvitTree

A nested list the reader can open and close. Twelve private versions; what a shared one has to get right is the keyboard, because depth without Left and Right is a wall.

```go
func NewTree(ui *UI, label string, nodes ...TreeNode) *Tree
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the tree for a screen reader. |
| `Model` | `TreeModel` | Model is what the tree shows. |
| `OnChoose` | `func(path []int)` | OnChoose runs when the reader moves to a node. |
| `OnActivate` | `func(path []int)` | OnActivate runs when the reader presses Return on a node with no children, or presses one twice. |

Methods: `Current`, `ExpandTo`, `IsExpanded`, `Refresh`, `SetExpanded`.

*A small hierarchy*

```go
func treeAccounts(ui *kvitui.UI) unison.Paneler {
    leaf := func(label string) kvitui.TreeNode { return kvitui.TreeNode{Label: label} }
    accounts := kvitui.NewTree(ui, "Accounts",
        kvitui.TreeNode{Label: "Everyday", Children: []kvitui.TreeNode{leaf("Checking"), leaf("Joint checking"), leaf("Cash")}},
        kvitui.TreeNode{Label: "Savings", Children: []kvitui.TreeNode{
            leaf("Emergency fund"),
            {Label: "Certificates", Children: []kvitui.TreeNode{leaf("18 months"), leaf("3 years")}}}},
        leaf("Credit card"))
    // The top level open, and everything under it closed.
    accounts.ExpandTo(1)
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(250), accounts))
}
```

### KvitSwitch

An option that takes effect the moment it moves — where a checkbox is a value in a form that takes effect on submit. The knob moves and the track fills, so the state is not resting on hue.

```go
func NewSwitch(ui *UI, label string) *Switch
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text is what the switch turns on. |
| `Checked` | `bool` | Checked is whether it is on. |
| `OnChange` | `func(checked bool)` | OnChange runs after a press, with the new state. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*On, off and disabled*

```go
func switchForms(ui *kvitui.UI) unison.Paneler {
    follow := kvitui.NewSwitch(ui, "Follow the system theme")
    follow.Checked = true
    cellular := kvitui.NewSwitch(ui, "Sync over cellular")
    cellular.SetEnabled(false)
    return kvitui.Column(ui, kvitui.SizeSpaceNear,
        kvitui.Left(follow), kvitui.Left(kvitui.NewSwitch(ui, "Reduce motion")), kvitui.Left(cellular))
}
```

### KvitRadioGroup

One choice from a handful where the choice needs explaining. Arrow keys move within the group and Tab leaves it, which a column of separate controls does not do.

```go
func NewRadioGroup(ui *UI, label string, options ...RadioOption) *RadioGroup
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the group for a screen reader. |
| `Current` | `string` | Current is the chosen option's value. |
| `OnChoose` | `func(value string)` | OnChoose runs when the reader chooses an option. |

Methods: `Choose`.

*Three options with detail*

```go
func radioAppearance(ui *kvitui.UI) unison.Paneler {
    appearance := kvitui.NewRadioGroup(ui, "Appearance",
        kvitui.RadioOption{Value: "system", Label: "Follow the system",
            Detail: "Light or dark, whichever the desktop is set to."},
        kvitui.RadioOption{Value: "light", Label: "Always light"},
        kvitui.RadioOption{Value: "dark", Label: "Always dark"})
    appearance.Current = "system"
    return kvitui.FullWidth(appearance)
}
```

### KvitProgress

How far through something the application is. An unknown total is a moving band rather than a bar creeping toward the end, and it carries a label saying how far through what.

```go
func NewProgress(ui *UI, label string, value float64) *Progress
```

| Field | Type | |
|---|---|---|
| `Value` | `float64` | Value against Maximum, 1 unless set. |
| `Maximum` | `float64` | Value against Maximum, 1 unless set. |
| `Determinate` | `bool` | Determinate says the total is known; false draws the moving band. |
| `Label` | `string` | Label says what is happening. |
| `ShowPercent` | `bool` | ShowPercent draws the share done beside the label. |

*Determinate and indeterminate*

```go
func progressForms(ui *kvitui.UI) unison.Paneler {
    reconciling := kvitui.NewProgress(ui, "Reconciling", 0)
    reconciling.Determinate = false
    return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
        kvitui.FullWidth(kvitui.NewProgress(ui, "Importing statement", 0.47)), kvitui.FullWidth(reconciling)))
}
```

### KvitSlider

A value chosen by dragging, for a feel rather than a number. Where the exact value matters, KvitStepper or KvitNumberField are right — hitting a particular number on a slider is hard and it cannot be typed.

```go
func NewSlider(ui *UI, label string, value float64) *Slider
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the value, for a screen reader. |
| `Unit` | `string` | Unit follows the value. |
| `Precision` | `int` | Precision is how many decimal places the value is drawn with. |
| `Value` | `float64` | Value is the value, From and To the range; 0 to 1 unless set. |
| `From` | `float64` | Value is the value, From and To the range; 0 to 1 unless set. |
| `To` | `float64` | Value is the value, From and To the range; 0 to 1 unless set. |
| `Step` | `float64` | Step is how far an arrow key moves it; a tenth of the range unless set. |
| `OnChange` | `func(value float64)` | OnChange runs as the value moves. |

Methods: `Focus`, `Hovered`, `KeyboardFocus`, `Pressed`.

*With its value shown*

```go
func sliderWithValue(ui *kvitui.UI) unison.Paneler {
    opacity := kvitui.NewSlider(ui, "Opacity", 0.6)
    opacity.Unit = "%"
    return kvitui.Width(ui, kvitui.Px(200), opacity)
}
```

### KvitSplitView

Two regions the reader can resize. The handle is a wide invisible strip with a hairline down the middle, so the target is comfortable and the rule is still thin; it also moves with the arrow keys.

```go
func NewSplitView(ui *UI, panes ...unison.Paneler) *SplitView
```

| Field | Type | |
|---|---|---|
| `Vertical` | `bool` | Vertical stacks the regions top to bottom rather than side by side. |
| `Fill` | `int` | Fill is the region that takes the room the others leave; the last unless set. |

Methods: `SetSize`.

*Two panes*

```go
func splitTwoPanes(ui *kvitui.UI) unison.Paneler {
    pane := func(words string) unison.Paneler {
        p := kvitui.NewPanel(ui)
        p.SetLayout(&unison.FlexLayout{Columns: 1})
        p.AddChild(kvitui.Centred(kvitui.NewLabel(ui, words)))
        return p
    }
    split := kvitui.NewSplitView(ui, pane("left"), pane("right"))
    split.SetSize(0, kvitui.Px(160))
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(120), split))
}
```

### KvitSegmented

One choice from two to five short options, all visible — the shape kvit-cash's dashboard period control needs, where a control governing seven widgets should say what they are showing without being opened.

```go
func NewSegmented(ui *UI, label string, options ...Option) *Segmented
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the choice for a screen reader, such as "Theme". |
| `Current` | `string` | Current is the Value of the chosen option. |
| `OnChoose` | `func(value string)` | OnChoose runs after an option is chosen, with its value. |

Methods: `Choose`.

*A period control*

```go
func segmentedPeriod(ui *kvitui.UI) unison.Paneler {
    period := kvitui.NewSegmented(ui, "Period",
        kvitui.Option{Value: "week", Label: "Week"}, kvitui.Option{Value: "month", Label: "Month"},
        kvitui.Option{Value: "quarter", Label: "Quarter"}, kvitui.Option{Value: "year", Label: "Year"})
    period.Current = "month"
    whose := kvitui.NewSegmented(ui, "",
        kvitui.Option{Value: "All", Label: "All"}, kvitui.Option{Value: "Mine", Label: "Mine"})
    return kvitui.Column(ui, kvitui.SizeSpace, kvitui.Left(period), kvitui.Left(whose))
}
```

### KvitTypeAhead

A field offering matches as the reader types, for a list too long to read. Whether the reader may make a new value has no default: a tag picker should let the reader invent one, a category picker should not.

```go
func NewTypeAhead(ui *UI, label string, allowNew bool, source ...Suggestion) *TypeAhead
```

| Field | Type | |
|---|---|---|
| `Field` | embedded | Everything a Field has. |
| `Source` | `[]Suggestion` | Source is everything that could be matched. |
| `AllowNew` | `bool` | AllowNew offers to make what was typed when nothing matches it exactly. |
| `MaximumSuggestions` | `int` | MaximumSuggestions caps the matches shown; 8 unless set. |
| `OnChoose` | `func(value string)` | OnChoose runs with the value chosen, or the words typed when a new one is made. |

Methods: `Matches`, `SetText`, `SetTrailing`.

*A category picker that will not invent categories*

```go
func typeAheadCategory(ui *kvitui.UI) unison.Paneler {
    var categories []kvitui.Suggestion
    for _, c := range []string{"Groceries", "Transport", "Utilities", "Rent", "Leisure", "Health", "Savings", "Income"} {
        categories = append(categories, kvitui.Suggestion{Value: c})
    }
    // A category picker that will not invent categories: a typo would make
    // a second category beside the right one.
    category := kvitui.NewTypeAhead(ui, "Category", false, categories...)
    category.Placeholder = "Start typing"
    return kvitui.Width(ui, kvitui.Px(240), category)
}
```

### KvitConfirmInPlace

A strip saying what just happened, with the undo inside it. Cheaper than a dialog before every action when the reader is doing the same thing forty times — and right only when the action can really be undone.

```go
func NewConfirmInPlace(ui *UI, text string) *ConfirmInPlace
```

| Field | Type | |
|---|---|---|
| `Text` | `string` | Text says what happened. |
| `UndoText` | `string` | UndoText is the words on the undo; "" draws no undo, for an action that has stopped being reversible. |
| `Shown` | `bool` | Shown opens the strip, growing to its height; false closes it. |
| `Affected` | `int` | Affected is how many things the action touched, so the reader can check the count before undoing; below zero says nothing. |
| `OnUndo` | `func()` | OnUndo and OnDismiss run when the undo or the close button is pressed. |
| `OnDismiss` | `func()` | OnUndo and OnDismiss run when the undo or the close button is pressed. |

*After a bulk edit*

```go
func confirmAfterBulkEdit(ui *kvitui.UI) unison.Paneler {
    done := kvitui.NewConfirmInPlace(ui, "Recategorised as Groceries")
    done.Shown, done.Affected = true, 40
    return kvitui.FullWidth(done)
}
```

### KvitTimeline

What happened to something, newest first, with who did it. In an estate where an agent and a person change the same things, who is the column that makes a history worth reading.

```go
func NewTimeline(ui *UI, label string, entries ...TimelineEntry) *Timeline
```

| Field | Type | |
|---|---|---|
| `Label` | `string` | Label names the history for a screen reader. |
| `Entries` | `[]TimelineEntry` | Entries are the events, newest first. |

*An account's recent history*

```go
func timelineHistory(ui *kvitui.UI) unison.Paneler {
    return kvitui.FullWidth(kvitui.NewTimeline(ui, "Account history",
        kvitui.TimelineEntry{When: "14:02", What: "Statement imported", Who: "agent",
            Detail: "412 transactions, 8 unmatched", Tone: kvitui.ToneSuccess},
        kvitui.TimelineEntry{When: "11:20", What: "Two rows disputed", Who: "you", Tone: kvitui.ToneWarning},
        kvitui.TimelineEntry{When: "Yesterday", What: "Account opened", Who: "you"}))
}
```

### KvitNumberField

A number the reader types, right-aligned in tabular numerals. It validates and says why rather than refusing keystrokes — a field that ignores a key gives no reason, and the usual cause is a decimal separator the reader's locale writes differently.

```go
func NewNumberField(ui *UI) *NumberField
```

| Field | Type | |
|---|---|---|
| `Field` | embedded | Everything a Field has. |
| `Minimum` | `float64` | Minimum and Maximum bound the value; no bound unless set. |
| `Maximum` | `float64` | Minimum and Maximum bound the value; no bound unless set. |
| `Decimals` | `int` | Decimals is how many decimal places the value has; 0 for a whole number. |

Methods: `SetText`, `Valid`, `Value`.

*Integer and decimal, valid and out of range*

```go
func numberFieldForms(ui *kvitui.UI) unison.Paneler {
    field := func(label, text string, minimum, maximum float64, decimals int) unison.Paneler {
        f := kvitui.NewNumberField(ui)
        f.Label, f.Minimum, f.Maximum, f.Decimals = label, minimum, maximum, decimals
        f.SetText(text)
        return kvitui.Left(f)
    }
    unbounded := math.Inf(1)
    return kvitui.Column(ui, kvitui.SizeSpaceLoose,
        field("Days", "14", 1, 365, 0),
        field("Rate", "4.25", -unbounded, unbounded, 2),
        field("Days", "999", 1, 365, 0))
}
```

### KvitMoneyField

An amount of money, in minor units. The reader types 12.34 and the field reports 1234, an integer — money in floating point drifts by a penny somewhere nobody can find. The decimal count comes from the currency.

```go
func NewMoneyField(ui *UI, currency string, minorDigits int) *MoneyField
```

| Field | Type | |
|---|---|---|
| `NumberField` | embedded | Everything a NumberField has. |
| `Currency` | `string` | Currency is the ISO 4217 code shown in the field; "" shows none, for a field where the currency is already clear. |

Methods: `MinorUnits`.

*Sterling, yen and dinar*

```go
func moneyFieldCurrencies(ui *kvitui.UI) unison.Paneler {
    amount := func(text, currency string, minorDigits int) unison.Paneler {
        f := kvitui.NewMoneyField(ui, currency, minorDigits)
        f.Label = "Amount"
        f.SetText(text)
        return kvitui.Left(f)
    }
    return kvitui.Column(ui, kvitui.SizeSpaceLoose,
        amount("1284.50", "GBP", 2), amount("4200", "JPY", 0), amount("18.750", "BHD", 3))
}
```

### KvitDualList

Two lists with items moving between them: what is available, and what is chosen and in what order. Everything works from the keyboard, which most implementations of this shape do not.

```go
func NewDualList(ui *UI, available, chosen []Option) *DualList
```

| Field | Type | |
|---|---|---|
| `Available` | `[]Option` | Available and Chosen are the two lists; Chosen is in order. |
| `Chosen` | `[]Option` | Available and Chosen are the two lists; Chosen is in order. |
| `AvailableLabel` | `string` | AvailableLabel and ChosenLabel head the two lists. |
| `ChosenLabel` | `string` | AvailableLabel and ChosenLabel head the two lists. |
| `OnChange` | `func(available, chosen []Option)` | OnChange runs after an item moves, with both lists as they now are. |

*Choosing table columns*

```go
func dualListColumns(ui *kvitui.UI) unison.Paneler {
    columns := func(names ...string) []kvitui.Option {
        var out []kvitui.Option
        for _, n := range names {
            out = append(out, kvitui.Option{Value: n, Label: n})
        }
        return out
    }
    return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(220), kvitui.NewDualList(ui,
        columns("Payee", "Account", "Tags", "Note"), columns("Date", "Description", "Amount", "Balance"))))
}
```

### KvitSpotlight

Darken everything except one region and say something about it — the primitive under a guided tour, and deliberately only the primitive: what drives the stepping differs between the two applications that want one.

```go
func NewSpotlight(ui *UI, target unison.Paneler, title, detail string) *Spotlight
```

| Field | Type | |
|---|---|---|
| `Target` | `unison.Paneler` | Target is the part left lit; nil darkens everything, the state before the first step. |
| `Title` | `string` | Title and Detail are what is said about it. |
| `Detail` | `string` | Title and Detail are what is said about it. |
| `Padding` | `Measure` | Padding is the room around the target inside the lit hole; a space unless set. |
| `OnDismiss` | `func()` | OnDismiss runs when the reader asks to leave, by Escape or a press. |

Methods: `Close`, `Open`, `Opened`.

*Focusing a button*

```go
func spotlightOnButton(ui *kvitui.UI) unison.Paneler {
    target := kvitui.NewButton(ui, "Import a statement")
    target.Form = kvitui.ButtonPrimary
    stage := kvitui.FullWidth(kvitui.At(ui, kvitui.Px(40), kvitui.Px(20), kvitui.Px(160), target))
    spot := kvitui.NewSpotlight(ui, target, "Start here", "Import a statement and the dashboard fills itself in.")
    spot.OnDismiss = spot.Close
    whenShown(stage, func() { spot.Open(stage) })
    return stage
}
```
