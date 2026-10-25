package main

import (
	"fmt"
	"strings"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

// The gallery's window: 1440 × 960.
const (
	windowWidth  = 1440
	windowHeight = 960
)

// group is one section of the sidebar: the component groups of the
// library's catalogue. A component without a page yet is listed greyed.
type group struct {
	name  string
	pages []string
}

var groups = []group{
	{"Foundation", []string{"Foundations", "KvitLabel", "KvitIcon", "KvitIconButton", "KvitLink"}},
	{"Structure", []string{"KvitHeader", "KvitSidebar", "KvitSidebarItem", "KvitBreadcrumb", "KvitRegion", "KvitViewHead", "KvitStatusBar", "KvitWindow", "KvitHistory"}},
	{"Content", []string{"KvitSectionHeading", "KvitRow", "KvitSlimRow", "KvitCard", "KvitPanel", "KvitPane", "KvitDivider", "KvitDisclosure", "KvitEmptyState"}},
	{"Marks", []string{"KvitChip", "KvitTag", "KvitBadge", "KvitSlug", "KvitDot", "KvitSignal", "KvitPip"}},
	{"Quantities", []string{"KvitFigure", "KvitBeforeAfter"}},
	{"Controls", []string{"KvitButton", "KvitChipButton", "KvitStepper", "KvitField", "KvitTextArea", "KvitSearchField", "KvitCheck", "KvitSelect", "KvitTab", "KvitFormRow", "KvitSettingRow", "KvitSettingsWindow"}},
	{"Commands", []string{"KvitCommand", "KvitMenuBar", "KvitToolbar", "KvitMenuButton", "KvitSplitButton", "KvitUndoStack"}},
	{"Feedback", []string{"KvitTooltip", "KvitPopover", "KvitHint", "KvitHoverCard", "KvitToast", "KvitNotice", "KvitDialog", "KvitTip", "KvitHelpLink", "KvitFloatingView"}},
	{"Data", []string{"KvitBar", "KvitStackedBar", "KvitSpark", "KvitNetFlow", "KvitTrend", "KvitDistribution", "KvitGauge", "KvitDelta", "KvitStatTile", "KvitFigureBlock", "KvitCell", "KvitTable"}},
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
// library's own components: a Window holding a Header with the theme and
// size controls, a sidebar of section headings and rows under a filter
// field, the page in a Region, and a StatusBar.
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
	form     *kvitui.Segmented
	size     *kvitui.Stepper
	status   *kvitui.StatusBar
	activity string // what the screenshot run is doing, for the status bar
}

func newGallery(ui *kvitui.UI, page string, firstFrame func()) (*gallery, error) {
	if !hasPage(page) {
		return nil, fmt.Errorf("no page %q; the pages are the ones listed without grey in the sidebar", page)
	}
	// The app menu on the macOS form, and Help, are named after it.
	ui.AppName = "kvit-ui gallery"
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
	g.form = kvitui.NewSegmented(ui, "Form",
		kvitui.Option{Value: kvitui.FormWindows.String(), Label: "Windows"},
		kvitui.Option{Value: kvitui.FormMac.String(), Label: "macOS"})
	g.form.OnChoose = func(v string) { ui.SetForm(formNamed(v)) }
	g.size = kvitui.NewStepper(ui, "Interface size", tokens.MinInterfaceSize, tokens.MaxInterfaceSize)
	g.size.Unit = "px"
	g.size.Follow = ui.Interface.FontSize
	g.size.OnChange = ui.Interface.SetFontSize
	room := kvitui.FullWidth(unison.NewPanel())
	wnd.SetHeader(kvitui.NewHeader(ui, "kvit-ui", kvitui.FullWidth(kvitui.Row(ui, kvitui.SizeColumnGap, room, g.form, g.theme, g.size))))

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
	wnd.SetMenuBar(g.menuBar())

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
	g.form.Current = ui.Form().String()
	g.status.Activity = g.activity
	g.status.Facts = []string{
		fmt.Sprintf("%d components", len(catalog)),
		ui.Form().String() + " form",
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

// formNamed is the form with a name as Form.String gives it.
func formNamed(name string) kvitui.Form {
	if name == kvitui.FormMac.String() {
		return kvitui.FormMac
	}
	return kvitui.FormWindows
}

var themeLabels = []struct{ id, label string }{
	{tokens.Light, "Light"}, {tokens.Dark, "Dark"}, {tokens.Sepia, "Sepia"}, {tokens.HighContrast, "Contrast"},
}

// menuBar is the gallery's commands: a theme on Ctrl+1 to Ctrl+4, the
// interface size on the zoom keys, the form, and Find, which goes to the
// sidebar's filter.
func (g *gallery) menuBar() *kvitui.MenuBar {
	ui := g.ui
	var themes []kvitui.MenuItem
	for i, tl := range themeLabels {
		themes = append(themes, kvitui.MenuItem{Radio: true, Command: &kvitui.Command{
			Name:     tl.label,
			Shortcut: kvitui.Keys(unison.Key1+unison.KeyCode(i), kvitui.ModPrimary),
			Checked:  func() bool { return ui.Theme.ThemeID() == tl.id },
			Run:      func() { ui.Theme.SetThemeID(tl.id) },
		}})
	}
	var forms []kvitui.MenuItem
	for _, f := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
		forms = append(forms, kvitui.MenuItem{Radio: true, Command: &kvitui.Command{
			Name:    f.String() + " form",
			Checked: func() bool { return ui.Form() == f },
			Run:     func() { ui.SetForm(f) },
		}})
	}
	size := ui.Interface
	return kvitui.NewMenuBar(ui).
		Add(&kvitui.Command{Std: kvitui.StdQuit, Run: func() { unison.AttemptQuit() }},
			&kvitui.Command{Std: kvitui.StdFind, Run: func() { g.search.Focus() }},
			&kvitui.Command{Std: kvitui.StdZoomIn, Run: func() { size.SetFontSize(size.FontSize() + 1) },
				Enabled: func() bool { return size.FontSize() < tokens.MaxInterfaceSize }},
			&kvitui.Command{Std: kvitui.StdZoomOut, Run: func() { size.SetFontSize(size.FontSize() - 1) },
				Enabled: func() bool { return size.FontSize() > tokens.MinInterfaceSize }},
			&kvitui.Command{Std: kvitui.StdActualSize, Run: func() { size.SetFontSize(tokens.DefaultInterfaceSize) }}).
		AddTo(kvitui.MenuView,
			kvitui.MenuItem{Text: "&Theme", Items: themes},
			kvitui.MenuItem{Text: "&Form", Items: forms})
}
