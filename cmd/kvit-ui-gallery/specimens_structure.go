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
