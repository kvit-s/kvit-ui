package kvitui_test

import (
	"strconv"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/platform"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
	"golang.org/x/text/language"
)

// The view head's count is grouped by the reader's locale and uses a real
// plural.
func TestTheViewHeadCountsInTheReadersLocale(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Locale = language.AmericanEnglish
	h := kvitui.NewViewHead(ui, "Transactions")
	h.Counted = "transaction"
	for _, c := range []struct {
		count  int
		plural string
		want   string
	}{
		{1, "", "1 transaction"},
		{250000, "", "250,000 transactions"},
		{250000, "entries", "250,000 entries"},
		{-1, "", ""},
	} {
		h.Count, h.CountedPlural = c.count, c.plural
		if got := h.CountPhrase(); got != c.want {
			t.Errorf("count %d, plural %q: %q, want %q", c.count, c.plural, got, c.want)
		}
	}
	for tag, want := range map[language.Tag]string{language.German: "250.000", language.Und: "250000"} {
		ui.Locale = tag
		if got := ui.Number(250000); got != want {
			t.Errorf("%v writes %q, want %q", tag, got, want)
		}
	}
}

func TestTheViewHeadIsARowTallWithItsControlsAtTheRight(t *testing.T) {
	var plain, full *kvitui.ViewHead
	var export *kvitui.IconButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		plain = kvitui.NewViewHead(ui, "Accounts")
		export = kvitui.NewIconButton(ui, "export", "Export")
		full = kvitui.NewViewHead(ui, "Transactions", export)
		full.Count, full.Counted = 1284, "transaction"
		full.Subtitle = "Everything since the account was opened, including the transfers between your own accounts."
		for _, h := range []*kvitui.ViewHead{plain, full} {
			h.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true, SizeHint: geom.NewSize(360, 0)})
		}
		return []unison.Paneler{plain, full}
	})
	screen.Do(func() {
		row := float32(ui.Interface.RowHeight())
		if h := plain.FrameRect().Height; h != row {
			t.Errorf("a head with only a title is %.1f tall, want a row, %.1f", h, row)
		}
		// Only the parts with something to say are there for a screen reader.
		if n := len(plain.Children()); n != 1 {
			t.Errorf("a head with only a title has %d children", n)
		}
		if n := len(full.Children()); n != 4 {
			t.Errorf("a head with every part has %d children", n)
		}
		// Narrow, the subtitle runs onto more lines, and the head grows to hold
		// them rather than letting them spill out of it.
		_, narrow, _ := full.Sizes(geom.NewSize(220, 0))
		if narrow.Height <= row {
			t.Errorf("at 220 wide the head is %.1f tall, no taller than a row", narrow.Height)
		}
		head := full.FrameRect()
		b := full.RectFromRoot(export.RectToRoot(export.ContentRect(true)))
		// The head keeps the view margin inside its own edges, as the rest of
		// a region does.
		if right, want := b.Right(), head.Width-float32(ui.Interface.ViewMargin()); right < want-0.5 || right > want+0.5 {
			t.Errorf("the control ends at %.1f, want the view margin in from the right edge, %.1f", right, want)
		}
		if mid := b.CenterY(); mid < head.Height/2-1 || mid > head.Height/2+1 {
			t.Errorf("the control is centred at %.1f, not on the head's middle %.1f", mid, head.Height/2)
		}
	})
	n := screen.AccessibilityNodeFor(full)
	if n == nil || n.Role != role.Heading || n.Name != "Transactions" {
		t.Errorf("the head's node: %+v", n)
	}
}

func TestTheBreadcrumbLinksBackAndCutsTheMiddle(t *testing.T) {
	var short, long *kvitui.Breadcrumb
	var followed []string
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		short = kvitui.NewBreadcrumb(ui, kvitui.Crumb{Label: "Projects", ID: "p"}, kvitui.Crumb{Label: "kvit-ui", ID: "u"})
		short.OnActivate = func(c kvitui.Crumb) { followed = append(followed, c.ID) }
		long = kvitui.NewBreadcrumb(ui,
			kvitui.Crumb{Label: "Projects", ID: "p"}, kvitui.Crumb{Label: "kvit", ID: "k"},
			kvitui.Crumb{Label: "Wave 1", ID: "w"}, kvitui.Crumb{Label: "Tokens", ID: "t"},
			kvitui.Crumb{Label: "Theme", ID: "h"})
		return []unison.Paneler{short, long}
	})
	var labels []string
	for _, c := range long.Shown() {
		labels = append(labels, c.Label)
	}
	if got := strings.Join(labels, " / "); got != "Projects / … / Tokens / Theme" {
		t.Errorf("the long trail shows %q", got)
	}
	screen.Do(func() {
		if h := short.FrameRect().Height; h != float32(ui.Interface.BreadcrumbHeight()) {
			t.Errorf("the breadcrumb is %.1f tall", h)
		}
	})
	// Crumb, chevron, crumb: the first is a link and the last is not.
	kids := short.Children()
	if len(kids) != 3 {
		t.Fatalf("a two-crumb trail has %d children", len(kids))
	}
	first, last := screen.AccessibilityNodeFor(kids[0]), screen.AccessibilityNodeFor(kids[2])
	if first == nil || first.Role != role.Link || last == nil || last.Role != role.Label {
		t.Fatalf("first crumb %+v, last crumb %+v", first, last)
	}
	if !screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: first.ID, Action: accessibility.Press}) {
		t.Fatal("the first crumb did not take a press")
	}
	screen.Sync()
	if len(followed) != 1 || followed[0] != "p" {
		t.Errorf("following the first crumb reported %v", followed)
	}
	// The "…" stands for several places and cannot be followed.
	if n := screen.AccessibilityNodeFor(long.Children()[2]); n == nil || n.Role != role.Label {
		t.Errorf("the elided crumb is %+v", n)
	}
	if n := screen.AccessibilityNodeFor(short); n == nil || n.Role != role.Group || n.Name != "Breadcrumb" {
		t.Errorf("the trail's node: %+v", n)
	}
}

func TestTheHeaderKeepsItsThreePlaces(t *testing.T) {
	var h *kvitui.Header
	var today *kvitui.Tab
	var settings *kvitui.IconButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		today = kvitui.NewTab(ui, "Today")
		settings = kvitui.NewIconButton(ui, "settings", "Settings")
		h = kvitui.NewHeader(ui, "kvit", kvitui.Row(ui, kvitui.SizeSpace, today), settings)
		h.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		return []unison.Paneler{h}
	})
	screen.Do(func() {
		m := ui.Interface
		box := h.FrameRect()
		if box.Height != float32(m.HeaderHeight()) {
			t.Errorf("the header is %.1f tall, want %d", box.Height, m.HeaderHeight())
		}
		in := func(p unison.Paneler) geom.Rect {
			return h.RectFromRoot(p.AsPanel().RectToRoot(p.AsPanel().ContentRect(true)))
		}
		margin := float32(m.ViewMargin())
		mark := h.Children()[0].FrameRect()
		if mark.X != margin {
			t.Errorf("the wordmark starts at %.1f, want the view margin %.1f", mark.X, margin)
		}
		if tab := in(today); tab.X < mark.Right()+float32(m.ColumnGap())-0.5 {
			t.Errorf("the navigation starts at %.1f, inside the wordmark's gap", tab.X)
		}
		if right := in(settings).Right(); right < box.Width-margin-0.5 || right > box.Width-margin+0.5 {
			t.Errorf("the actions end at %.1f, want %.1f", right, box.Width-margin)
		}
	})
	if n := screen.AccessibilityNodeFor(h); n == nil || n.Role != role.TabList || n.Name != "Application header" {
		t.Errorf("the header's node: %+v", n)
	}
}

func TestASidebarItemSaysWhereItGoesAndHowMuchIsThere(t *testing.T) {
	var inbox, review *kvitui.SidebarItem
	var rail *kvitui.Sidebar
	pressed := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Locale = language.AmericanEnglish
		inbox = kvitui.NewSidebarItem(ui, "Inbox", "note")
		inbox.Selected, inbox.Count, inbox.Counted = true, 12, "message"
		inbox.OnClick = func() { pressed++ }
		review = kvitui.NewSidebarItem(ui, "Review", "question")
		review.Count = 214
		rail = kvitui.NewSidebar(ui, review)
		rail.Collapsed = true
		return []unison.Paneler{kvitui.NewSidebar(ui, inbox), rail}
	})
	n := screen.AccessibilityNodeFor(inbox)
	if n == nil || n.Role != role.ListItem || n.Name != "Inbox, 12 messages" || !n.Selected {
		t.Fatalf("the selected item's node: %+v", n)
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: n.ID, Action: accessibility.Press})
	screen.Sync()
	if pressed != 1 {
		t.Errorf("a screen reader's press made %d presses", pressed)
	}
	badgeOf := func(it *kvitui.SidebarItem) *kvitui.Badge { return it.Children()[0].Self.(*kvitui.Badge) }
	screen.Do(func() {
		m := ui.Interface
		if !review.Collapsed {
			t.Error("an item in a collapsed sidebar is not collapsed")
		}
		if w := rail.FrameRect().Width; w != float32(m.RailWidth()) {
			t.Errorf("the rail is %.1f wide, want %d", w, m.RailWidth())
		}
		if h := inbox.FrameRect().Height; h != float32(m.RowHeightSlim()) {
			t.Errorf("an item is %.1f tall, want %d", h, m.RowHeightSlim())
		}
		// The rail keeps its count behind unless asked to show it.
		if !badgeOf(inbox).Shows() || badgeOf(review).Shows() {
			t.Error("the badge shows in the rail, or not in the expanded sidebar")
		}
		review.CountInRail = true
		rail.MarkForLayoutRecursively()
		rail.ValidateLayout()
		if !badgeOf(review).Shows() || badgeOf(review).Text() != "99+" {
			t.Errorf("CountInRail shows %q", badgeOf(review).Text())
		}
	})
	if n := screen.AccessibilityNodeFor(rail); n == nil || n.Role != role.List || n.Name != "Navigation" {
		t.Errorf("the sidebar's node: %+v", n)
	}
}

func TestARegionScrollsBesideItsBarAndAnswersTheWheelAndKeys(t *testing.T) {
	var region *kvitui.Region
	var first, last *kvitui.ListRow
	var rows []*kvitui.SlimRow
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		// The wheel eases in over a few frames; with motion reduced it
		// arrives at once, which is what a test can measure.
		ui.Theme.SetReducedMotion(true)
		first = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "First"))
		first.Interactive = true
		last = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Last"))
		last.Interactive = true
		parts := []unison.Paneler{first}
		for i := range 12 {
			r := kvitui.NewSlimRow(ui, "Row "+strconv.Itoa(i+1))
			rows = append(rows, r)
			parts = append(parts, r)
		}
		parts = append(parts, last)
		region = kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), parts...))
		region.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(480, 140)})
		return []unison.Paneler{region}
	})
	m := ui.Interface
	screen.Do(func() {
		if !region.Scrolls() {
			t.Fatal("fourteen rows in 140 px do not scroll")
		}
		// The bar has a strip of its own at the right, and the rows stop at
		// the padding before it.
		box := region.FrameRect()
		bar := region.Children()[1].FrameRect()
		if bar.X != box.Width-float32(m.SpaceWide()) || bar.Width != float32(m.SpaceWide()) {
			t.Errorf("the bar's strip is %v in a region %.1f wide", bar, box.Width)
		}
		want := box.Width - float32(m.SpaceWide()) - 2*float32(m.ViewMargin())
		if w := rows[0].FrameRect().Width; w != want {
			t.Errorf("a row is %.1f wide, want %.1f", w, want)
		}
	})
	position := func() (y float32) {
		screen.Do(func() { _, y = region.Position() })
		return y
	}
	// One notch travels the desktop's lines per notch, of a slim row each.
	screen.Wheel(screen.PanelCenter(region), geom.NewPoint(0, -1), 0)
	screen.Sync()
	if y, want := position(), float32(platform.WheelScrollLines()*m.RowHeightSlim()); y != want {
		t.Errorf("a notch scrolled %.1f, want %.1f", y, want)
	}
	// Keys nothing inside wanted scroll the region.
	screen.Do(func() { region.ScrollTo(0) })
	screen.Click(screen.PanelCenter(first))
	screen.KeyPress(unison.KeyEnd, 0)
	var end float32
	screen.Do(func() { end = region.Children()[1].Self.(*unison.ScrollBar).MaxValue() })
	if y := position(); y != end || end <= 0 {
		t.Errorf("End scrolled to %.1f, want the bottom, %.1f", y, end)
	}
	screen.KeyPress(unison.KeyHome, 0)
	if y := position(); y != 0 {
		t.Errorf("Home scrolled to %.1f", y)
	}
	// Tab to the last row, below the fold, brings it into view.
	screen.KeyPress(unison.KeyTab, 0)
	screen.Do(func() {
		if !last.Focused() {
			t.Fatal("Tab did not reach the last row")
		}
		_, y := region.Position()
		box := region.RectFromRoot(last.RectToRoot(last.ContentRect(true)))
		if box.Bottom() > region.FrameRect().Height {
			t.Errorf("the focused row ends at %.1f, below the region (scrolled to %.1f)", box.Bottom(), y)
		}
	})
}

func TestAStatusBarKeepsWhatDoesNotFitBehindACount(t *testing.T) {
	groups := []kvitui.StatusGroup{
		{Label: "Waiting on you", Facts: []kvitui.StatusFact{{Text: "3 reviews", Symbol: "question"}, {Text: "1 conflict", Symbol: "warning"}}},
		{Label: "Running", Facts: []kvitui.StatusFact{{Text: "2 agents", Symbol: "robot"}}},
	}
	var wide, narrow, plain *kvitui.StatusBar
	var pressed [][2]int
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		wide = kvitui.NewStatusBar(ui)
		wide.Activity, wide.Groups = "Indexing 4 of 26 working copies", groups
		wide.OnFact = func(g, f int) { pressed = append(pressed, [2]int{g, f}) }
		narrow = kvitui.NewStatusBar(ui)
		narrow.Activity, narrow.Groups = "Indexing", groups
		plain = kvitui.NewStatusBar(ui)
		plain.Facts = []string{"1,284 notes", "last synced 14:02"}
		return []unison.Paneler{
			kvitui.Column(ui, kvitui.SizeSpace, kvitui.Width(ui, kvitui.Px(700), wide),
				kvitui.Width(ui, kvitui.Px(420), narrow), kvitui.Width(ui, kvitui.Px(320), plain)),
		}
	})
	var overflow *kvitui.Link
	screen.Do(func() {
		if wide.Shown() != 2 || narrow.Shown() != 1 {
			t.Errorf("the wide bar shows %d groups and the narrow one %d", wide.Shown(), narrow.Shown())
		}
		if h := plain.FrameRect().Height; h != float32(ui.Interface.StatusBarHeight()) {
			t.Errorf("a bar of plain facts is %.1f tall, want %d", h, ui.Interface.StatusBarHeight())
		}
		for _, c := range narrow.Children() {
			if l, ok := c.Self.(*kvitui.Link); ok && l.Symbol == "more" {
				overflow = l
			}
		}
	})
	if overflow == nil {
		t.Fatal("the narrow bar has no overflow link")
	}
	if n := screen.AccessibilityNodeFor(overflow); n == nil || n.Name != "1 more" || n.Description != "Running" {
		t.Errorf("the overflow link's node: %+v", n)
	}
	// A fact on the bar is a link that reports which fact it is.
	var conflict *kvitui.Link
	screen.Do(func() {
		var find func(p *unison.Panel)
		find = func(p *unison.Panel) {
			for _, c := range p.Children() {
				if l, ok := c.Self.(*kvitui.Link); ok && l.Text == "1 conflict" {
					conflict = l
				}
				find(c)
			}
		}
		find(wide.AsPanel())
	})
	if conflict == nil {
		t.Fatal("the wide bar has no link for the conflict")
	}
	screen.Click(screen.PanelCenter(conflict))
	if len(pressed) != 1 || pressed[0] != [2]int{0, 1} {
		t.Errorf("pressing the conflict reported %v", pressed)
	}
}

func TestTheWindowCollapsesItsSidebarAndOpensTheRailOverTheBody(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Theme.SetReducedMotion(true) // the sidebar's width then changes at once
	var w *kvitui.Window
	var side *kvitui.Sidebar
	var body *kvitui.Region
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1600, Height: 1000},
		unison.StartupFinishedCallback(func() {
			if w, err = kvitui.NewWindow(ui, "shell"); err != nil {
				t.Error(err)
				return
			}
			side = kvitui.NewSidebar(ui, kvitui.NewSidebarItem(ui, "Everyday", "wallet"), kvitui.NewSidebarItem(ui, "Savings", "bank"))
			body = kvitui.NewRegion(ui, kvitui.NewSlimRow(ui, "Checking"))
			w.SetHeader(kvitui.NewHeader(ui, "kvit", nil))
			w.SetSidebar(side)
			w.SetBody(body)
			w.SetStatusBar(kvitui.NewStatusBar(ui))
			w.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()
	m := ui.Interface
	resize := func(width float32) {
		screen.Do(func() {
			w.SetContentRect(geom.NewRect(0, 0, width, 900))
			w.Content().MarkForLayoutRecursively()
			w.ValidateLayout()
		})
		screen.Sync()
	}
	check := func(when string, collapsed bool, sidebar, bodyLeft int) {
		screen.Do(func() {
			if side.Collapsed != collapsed {
				t.Errorf("%s: the sidebar's Collapsed is %v", when, side.Collapsed)
			}
			if got := side.Parent().FrameRect().Width; got != float32(sidebar) {
				t.Errorf("%s: the sidebar is %.1f wide, want %d", when, got, sidebar)
			}
			if got := body.RectToRoot(body.ContentRect(true)).X; got != float32(bodyLeft) {
				t.Errorf("%s: the body starts at %.1f, want %d", when, got, bodyLeft)
			}
		})
	}
	resize(float32(m.WidthDrawn()))
	check("wide", false, m.SidebarWidth(), m.SidebarWidth())
	resize(float32(m.WidthLaptop() - 100))
	check("narrow", true, m.RailWidth(), m.RailWidth())
	// The pointer on the rail opens it over the body, which stays put.
	screen.MouseMove(geom.NewPoint(float32(m.RailWidth())/2, 200), 0)
	screen.Sync()
	check("rail under the pointer", false, m.SidebarWidth(), m.RailWidth())
	screen.MouseMove(geom.NewPoint(600, 200), 0)
	screen.Sync()
	check("pointer gone", true, m.RailWidth(), m.RailWidth())
	// The keyboard inside the rail opens it as well: Tab moves between its
	// items. The focus the window handed its first item when it opened did
	// not open it, as the narrow check above shows.
	screen.KeyPress(unison.KeyTab, 0)
	check("keyboard in the rail", false, m.SidebarWidth(), m.RailWidth())
	// Never smaller than the floor.
	screen.Do(func() { w.SetContentRect(geom.NewRect(0, 0, 400, 300)) })
	screen.Sync()
	screen.Do(func() {
		if r := w.ContentRect(); r.Width < float32(m.WidthFloor()) || r.Height < float32(m.HeightFloor()) {
			t.Errorf("the window was made %v, below the floor", r.Size)
		}
	})
}

func TestTheStatusBarMenuHoldsWhatDidNotFit(t *testing.T) {
	var bar *kvitui.StatusBar
	var pressed [][2]int
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		bar = kvitui.NewStatusBar(ui)
		bar.Activity = "Indexing"
		bar.Groups = []kvitui.StatusGroup{
			{Label: "Waiting on you", Facts: []kvitui.StatusFact{{Text: "3 reviews", Symbol: "question"}, {Text: "1 conflict", Symbol: "warning"}}},
			{Label: "Running", Facts: []kvitui.StatusFact{{Text: "2 agents", Symbol: "robot"}}},
		}
		bar.OnFact = func(g, f int) { pressed = append(pressed, [2]int{g, f}) }
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(420), bar)}
	})
	var more *kvitui.Link
	screen.Do(func() {
		for _, c := range bar.Children() {
			if l, ok := c.Self.(*kvitui.Link); ok && l.Symbol == "more" && !l.Hidden {
				more = l
			}
		}
	})
	if more == nil {
		t.Fatal("no overflow link")
	}
	screen.Click(screen.PanelCenter(more))
	screen.Sync()
	var w *unison.Window
	screen.Do(func() { w = bar.Window() })
	var entry *accessibility.Node
	for _, n := range screen.AccessibilityTree(w).Nodes {
		if n.Role == role.MenuItem && n.Name == "Running — 2 agents" {
			entry = n
		}
	}
	if entry == nil {
		t.Fatal("the menu has no entry for the hidden fact")
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: entry.ID, Action: accessibility.Press})
	screen.Sync()
	if len(pressed) != 1 || pressed[0] != [2]int{1, 0} {
		t.Errorf("choosing the entry reported %v", pressed)
	}
}
