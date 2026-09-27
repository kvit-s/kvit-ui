package main

// The Structure group's specimens.

import (
	"fmt"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

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

func sidebarItemStates(ui *kvitui.UI) unison.Paneler {
	inbox := kvitui.NewSidebarItem(ui, "Inbox", "note")
	inbox.Selected = true
	inbox.Count, inbox.Counted = 12, "message"
	archive := kvitui.NewSidebarItem(ui, "Archive", "archive")
	return kvitui.Width(ui, kvitui.SizeSidebarWidth, kvitui.Column(ui, kvitui.Px(0), inbox, archive))
}

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

func regionScrolling(ui *kvitui.UI) unison.Paneler {
	rows := make([]unison.Paneler, 12)
	for i := range rows {
		rows[i] = kvitui.NewSlimRow(ui, fmt.Sprintf("Row %d", i+1))
	}
	region := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), rows...))
	return kvitui.Sized(ui, kvitui.Px(480), kvitui.Px(140), region)
}

func viewHeadWithControls(ui *kvitui.UI) unison.Paneler {
	head := kvitui.NewViewHead(ui, "Transactions",
		kvitui.Width(ui, kvitui.Px(180), kvitui.NewSearchField(ui)), kvitui.NewButton(ui, "Export"))
	head.Count, head.Counted = 1284, "transaction"
	head.Subtitle = "Everything since the account was opened"
	return kvitui.FullWidth(head)
}

func statusBarWorking(ui *kvitui.UI) unison.Paneler {
	bar := kvitui.NewStatusBar(ui)
	bar.Activity = "Reindexing 4 of 26 projects"
	bar.Facts = []string{"1,284 notes", "last synced 14:02"}
	return kvitui.FullWidth(bar)
}

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
