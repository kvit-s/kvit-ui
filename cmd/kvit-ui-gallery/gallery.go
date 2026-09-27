package main

import (
	"fmt"
	"strings"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
)

// The gallery's window, as the Qt gallery draws it: 1440 × 960.
const (
	windowWidth  = 1440
	windowHeight = 960
)

// group is one section of the sidebar: the component groups of the Qt
// library's catalogue. A component without a page yet is listed greyed.
type group struct {
	name  string
	pages []string
}

var groups = []group{
	{"Foundation", []string{"Foundations", "KvitLabel", "KvitIcon", "KvitIconButton", "KvitLink"}},
	{"Structure", []string{"KvitHeader", "KvitSidebar", "KvitSidebarItem", "KvitBreadcrumb", "KvitRegion", "KvitViewHead", "KvitStatusBar", "KvitWindow"}},
	{"Content", []string{"KvitSectionHeading", "KvitRow", "KvitSlimRow", "KvitCard", "KvitPanel", "KvitPane", "KvitDivider", "KvitDisclosure", "KvitEmptyState"}},
	{"Marks", []string{"KvitChip", "KvitTag", "KvitBadge", "KvitSlug", "KvitDot", "KvitSignal", "KvitPip"}},
	{"Quantities", []string{"KvitFigure", "KvitBeforeAfter"}},
	{"Controls", []string{"KvitButton", "KvitChipButton", "KvitStepper", "KvitField", "KvitTextArea", "KvitSearchField", "KvitCheck", "KvitSelect", "KvitTab"}},
	{"Feedback", []string{"KvitTooltip", "KvitPopover", "KvitHint", "KvitHoverCard", "KvitToast", "KvitNotice", "KvitDialog"}},
	{"Data", []string{"KvitBar", "KvitStackedBar", "KvitSpark", "KvitTrend", "KvitDistribution", "KvitGauge", "KvitDelta", "KvitStatTile", "KvitFigureBlock", "KvitCell", "KvitTable"}},
	{"Flow", []string{"KvitScrollBar", "KvitMenu", "KvitMenuItem", "KvitTree", "KvitSwitch", "KvitRadioGroup", "KvitProgress", "KvitSlider", "KvitSplitView", "KvitSegmented", "KvitTypeAhead", "KvitConfirmInPlace", "KvitTimeline", "KvitNumberField", "KvitMoneyField", "KvitDualList", "KvitSpotlight"}},
}

// pageNames lists the pages that exist: the foundations, then every
// component in the catalogue.
func pageNames() []string {
	names := []string{"Foundations"}
	for _, e := range catalog {
		names = append(names, e.name)
	}
	return names
}

func hasPage(name string) bool {
	_, ok := entryNamed(name)
	return ok || name == "Foundations"
}

// buildPage builds a page's panel.
func buildPage(ui *kvitui.UI, name string) *unison.Panel {
	if e, ok := entryNamed(name); ok {
		return buildComponentPage(ui, e)
	}
	return drawnPage(ui, drawFoundations)
}

// gallery is the window and the state it shows. It is built from the
// library's own components, as the Qt gallery is: a Window holding a Header
// with the theme and size controls, a sidebar of section headings and rows
// under a filter field, the page in a Region, and a StatusBar.
type gallery struct {
	ui       *kvitui.UI
	wnd      *kvitui.Window
	page     string
	pages    *unison.Panel // holds the page being shown, inside the region
	region   *kvitui.Region
	list     *unison.Panel // the sidebar's groups and rows
	rows     map[string]*kvitui.ListRow
	search   *kvitui.SearchField
	theme    *kvitui.Segmented
	size     *kvitui.Stepper
	status   *kvitui.StatusBar
	activity string // what the screenshot run is doing, for the status bar
}

func newGallery(ui *kvitui.UI, page string, firstFrame func()) (*gallery, error) {
	if !hasPage(page) {
		return nil, fmt.Errorf("no page %q; the pages are the ones listed without grey in the sidebar", page)
	}
	wnd, err := kvitui.NewWindow(ui, "kvit-ui gallery")
	if err != nil {
		return nil, err
	}
	g := &gallery{ui: ui, wnd: wnd, rows: map[string]*kvitui.ListRow{}}

	var options []kvitui.Option
	for _, tl := range themeLabels {
		options = append(options, kvitui.Option{Value: tl.id, Label: tl.label})
	}
	g.theme = kvitui.NewSegmented(ui, "Theme", options...)
	g.theme.OnChoose = func(v string) { ui.Theme.SetThemeID(v) }
	g.size = kvitui.NewStepper(ui, "Interface size", tokens.MinInterfaceSize, tokens.MaxInterfaceSize)
	g.size.Unit = "px"
	g.size.Follow = ui.Interface.FontSize
	g.size.OnChange = ui.Interface.SetFontSize
	room := kvitui.FullWidth(unison.NewPanel())
	wnd.SetHeader(kvitui.NewHeader(ui, "kvit-ui", kvitui.FullWidth(kvitui.Row(ui, kvitui.SizeColumnGap, room, g.theme, g.size))))

	g.search = kvitui.NewSearchField(ui)
	g.search.Placeholder, g.search.MatchedNoun = "Filter components", "component"
	g.search.OnChange = func(string) { g.buildList() }
	searchBox := kvitui.Column(ui, kvitui.Px(0), g.search)
	searchBox.SetBorder(kvitui.Padding(ui, kvitui.SizeSpace))
	g.list = kvitui.Column(ui, kvitui.Px(0))
	listRegion := kvitui.NewRegion(ui, g.list)
	listRegion.Padding = kvitui.Px(0)
	listRegion.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
	wnd.SetSidebar(kvitui.Column(ui, kvitui.Px(0), searchBox, listRegion))

	g.pages = unison.NewPanel()
	g.pages.SetLayout(&unison.FlexLayout{Columns: 1})
	g.region = kvitui.NewRegion(ui, g.pages)
	wnd.SetBody(g.region)

	g.status = kvitui.NewStatusBar(ui)
	wnd.SetStatusBar(g.status)
	wnd.OnKeyDown = g.keyDown

	ui.OnChanged(g.sync)
	g.buildList()
	g.setPage(page)
	g.sync()
	if firstFrame != nil {
		content := wnd.Content()
		draw := content.DrawCallback
		drawn := false
		content.DrawCallback = func(gc *unison.Canvas, r geom.Rect) {
			draw(gc, r)
			if !drawn {
				drawn = true
				unison.InvokeTask(firstFrame)
			}
		}
	}
	wnd.ToFront()
	return g, nil
}

// sync shows the current theme and size in the header and the status bar.
func (g *gallery) sync() {
	ui := g.ui
	g.theme.Current = ui.Theme.ThemeID()
	g.status.Activity = g.activity
	g.status.Facts = []string{
		fmt.Sprintf("%d of 74 components", len(catalog)),
		tokens.DisplayName(ui.Theme.ResolvedTheme()),
		fmt.Sprintf("%d px", ui.Interface.FontSize()),
	}
	g.wnd.Content().MarkForLayoutRecursively()
	g.wnd.MarkForRedraw()
}

// buildList fills the sidebar with the groups and pages the filter matches:
// a section heading for each group, and a row for each page. A component
// without a page yet is listed faint, and cannot be opened.
func (g *gallery) buildList() {
	ui := g.ui
	g.list.RemoveAllChildren()
	g.rows = map[string]*kvitui.ListRow{}
	wanted := strings.ToLower(strings.TrimSpace(g.search.Text()))
	shown := 0
	for _, gr := range groups {
		var pages []string
		for _, name := range gr.pages {
			if wanted == "" || strings.Contains(strings.ToLower(name), wanted) || strings.Contains(strings.ToLower(gr.name), wanted) {
				pages = append(pages, name)
			}
		}
		if len(pages) == 0 {
			continue
		}
		g.list.AddChild(kvitui.FullWidth(kvitui.NewSectionHeading(ui, gr.name)))
		for _, name := range pages {
			label := kvitui.NewLabel(ui, name)
			label.SetBorder(kvitui.Insets(ui, nil, kvitui.SizeSpaceLoose, nil, nil))
			if !hasPage(name) {
				label.Ink = kvitui.InkTextDisabled
			}
			row := kvitui.NewListRow(ui, kvitui.FullWidth(label))
			row.Rule, row.Label, row.Selected = false, name, name == g.page
			// A list of pages is navigation, like a sidebar, and draws no
			// chevron on every line.
			row.Interactive, row.NoOpensMark = hasPage(name), true
			row.OnActivate = func() { g.setPage(name) }
			g.rows[name] = row
			g.list.AddChild(kvitui.FullWidth(row))
			shown++
		}
	}
	g.search.Matches = shown
	g.list.MarkForLayoutAndRedraw()
}

// setPage shows a page, rebuilt from scratch, scrolled to its top.
func (g *gallery) setPage(name string) {
	g.page = name
	for n, row := range g.rows {
		row.Selected = n == name
		row.MarkForRedraw()
	}
	g.pages.RemoveAllChildren()
	pg := buildPage(g.ui, name)
	pg.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	g.pages.AddChild(pg)
	g.pages.MarkForLayoutRecursively()
	g.region.ScrollTo(0)
	g.wnd.MarkForRedraw()
}

// pageHeight is the height the body needs to show the current page in full
// at a body width: the page at the width the region gives it, and the
// region's padding above and below.
func (g *gallery) pageHeight(width float32) float32 {
	m := g.ui.Interface
	inner := width - float32(m.SpaceWide()) - 2*float32(m.ViewMargin())
	_, pref, _ := g.pages.Sizes(geom.NewSize(inner, 0))
	return pref.Height + 2*float32(m.ViewMargin())
}

var themeLabels = []struct{ id, label string }{
	{tokens.Light, "Light"}, {tokens.Dark, "Dark"}, {tokens.Sepia, "Sepia"}, {tokens.HighContrast, "Contrast"},
}

// keyDown: Ctrl+1 to Ctrl+4 choose a theme, Ctrl+plus and Ctrl+minus change
// the interface size, Ctrl+0 resets it.
func (g *gallery) keyDown(key unison.KeyCode, mods mod.Modifiers, _ bool) bool {
	if !mods.CommandDown() {
		return false
	}
	switch key {
	case unison.Key1, unison.Key2, unison.Key3, unison.Key4:
		g.ui.Theme.SetThemeID(themeLabels[int(key-unison.Key1)].id)
	case unison.KeyEqual, unison.KeyNumPadAdd:
		g.ui.Interface.SetFontSize(g.ui.Interface.FontSize() + 1)
	case unison.KeyMinus, unison.KeyNumPadSubtract:
		g.ui.Interface.SetFontSize(g.ui.Interface.FontSize() - 1)
	case unison.Key0:
		g.ui.Interface.SetFontSize(tokens.DefaultInterfaceSize)
	default:
		return false
	}
	return true
}
