package main

import (
	"fmt"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
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

// gallery is the window and the state it shows.
type gallery struct {
	ui     *kvitui.UI
	wnd    *unison.Window
	page   string
	header *unison.Panel
	side   *unison.Panel
	body   *unison.Panel
	scroll *unison.ScrollPanel
	status *unison.Panel
	// segments are the theme choices' rectangles in the header, for clicks.
	segments map[string]geom.Rect
	// rows are the sidebar's page rows as last drawn, for clicks.
	rows  map[string]geom.Rect
	minus geom.Rect
	plus  geom.Rect
}

func newGallery(ui *kvitui.UI, page string, firstFrame func()) (*gallery, error) {
	if !hasPage(page) {
		return nil, fmt.Errorf("no page %q; the pages are the ones listed without grey in the sidebar", page)
	}
	g := &gallery{ui: ui, segments: map[string]geom.Rect{}, rows: map[string]geom.Rect{}}
	wnd, err := unison.NewWindow("kvit-ui gallery")
	if err != nil {
		return nil, err
	}
	g.wnd = wnd
	content := wnd.Content()
	content.SetLayout(&unison.FlexLayout{Columns: 1})
	drawn := false
	content.DrawCallback = func(gc *unison.Canvas, r geom.Rect) {
		painter{gc, ui}.fill(r, ui.Theme.Tokens().WindowBackground)
		if !drawn {
			drawn = true
			if firstFrame != nil {
				unison.InvokeTask(firstFrame)
			}
		}
	}

	g.header = unison.NewPanel()
	g.header.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := geom.NewSize(0, float32(ui.Interface.HeaderHeight()))
		return s, s, geom.NewSize(unison.DefaultMaxSize, s.Height)
	})
	g.header.DrawCallback = g.drawHeader
	g.header.MouseDownCallback = g.headerClick
	g.header.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	content.AddChild(g.header)

	middle := unison.NewPanel()
	middle.SetLayout(&unison.FlexLayout{Columns: 2})
	middle.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
	content.AddChild(middle)

	g.side = unison.NewPanel()
	g.side.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := geom.NewSize(float32(ui.Interface.SidebarWidth()), 0)
		return s, s, geom.NewSize(s.Width, unison.DefaultMaxSize)
	})
	g.side.DrawCallback = g.drawSidebar
	g.side.MouseDownCallback = func(where geom.Point, _, _ int, _ mod.Modifiers) bool {
		for name, r := range g.rows {
			if where.In(r) && hasPage(name) {
				g.setPage(name)
				return true
			}
		}
		return false
	}
	g.side.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Fill, VGrab: true})
	middle.AddChild(g.side)

	g.body = unison.NewPanel()
	g.body.SetLayout(&unison.FlexLayout{Columns: 1})
	g.body.SetBorder(scrollStrip{ui})
	g.body.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		painter{gc, ui}.fill(g.body.ContentRect(true), ui.Theme.Tokens().WindowBackground)
	}
	g.scroll = unison.NewScrollPanel()
	g.scroll.SetContent(g.body, behavior.Fill, behavior.Unmodified)
	g.scroll.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
	middle.AddChild(g.scroll)

	g.status = unison.NewPanel()
	g.status.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := geom.NewSize(0, float32(ui.Interface.StatusBarHeight()))
		return s, s, geom.NewSize(unison.DefaultMaxSize, s.Height)
	})
	g.status.DrawCallback = g.drawStatus
	g.status.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	content.AddChild(g.status)

	g.setPage(page)
	wnd.KeyDownCallback = g.keyDown
	ui.OnChanged(func() {
		content.MarkForLayoutRecursively()
		wnd.MarkForRedraw()
	})
	wnd.SetContentRect(geom.NewRect(0, 0, windowWidth, windowHeight))
	wnd.ToFront()
	return g, nil
}

// scrollStrip keeps the strip a vertical scroll bar is drawn in clear of the
// page whether or not the bar is showing, as the Qt region does, so the page
// does not move sideways when it grows past the fold.
type scrollStrip struct{ ui *kvitui.UI }

func (s scrollStrip) Insets() geom.Insets {
	return geom.Insets{Right: float32(s.ui.Interface.SpaceWide())}
}
func (s scrollStrip) Draw(*unison.Canvas, geom.Rect) {}

// setPage shows a page, rebuilt from scratch, scrolled to its top.
func (g *gallery) setPage(name string) {
	g.page = name
	g.body.RemoveAllChildren()
	pg := buildPage(g.ui, name)
	pg.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	g.body.AddChild(pg)
	g.body.MarkForLayoutRecursively()
	g.scroll.SetPosition(0, 0)
	g.wnd.MarkForRedraw()
}

// pageHeight is the height the current page needs at a width.
func (g *gallery) pageHeight(width float32) float32 {
	_, pref, _ := g.body.Sizes(geom.NewSize(width, 0))
	return pref.Height
}

var themeLabels = []struct{ id, label string }{
	{tokens.Light, "Light"}, {tokens.Dark, "Dark"}, {tokens.Sepia, "Sepia"}, {tokens.HighContrast, "Contrast"},
}

func (g *gallery) drawHeader(gc *unison.Canvas, r geom.Rect) {
	ui := g.ui
	m := ui.Interface
	t := ui.Theme.Tokens()
	p := painter{gc, ui}
	b := g.header.ContentRect(false)
	p.fill(b, t.PanelBackground)
	p.fill(geom.NewRect(0, b.Height-float32(m.Hairline()), b.Width, float32(m.Hairline())), t.Border)
	l := ui.Fonts.Layout(span("kvit-ui", ui.Chrome(m.Headline(), text.Bold, t.TextPrimary)), text.Options{})
	_, lh := l.Size()
	l.Draw(gc, float32(m.ViewMargin()), (b.Height-lh)/2)

	// The size stepper at the right, then the theme choices left of it.
	ch := float32(m.ControlHeight())
	cy := (b.Height - ch) / 2
	x := b.Width - float32(m.ViewMargin())
	g.plus = geom.NewRect(x-ch, cy, ch, ch)
	x -= ch
	sizeLabel := ui.Fonts.Layout(span(fmt.Sprintf("%d px", m.FontSize()), ui.Chrome(m.Body(), text.Regular, t.TextSecondary)), text.Options{})
	sw, sh := sizeLabel.Size()
	x -= sw + float32(m.Space())*2
	sizeLabel.Draw(gc, x+float32(m.Space()), (b.Height-sh)/2)
	g.minus = geom.NewRect(x-ch, cy, ch, ch)
	for _, btn := range []struct {
		r    geom.Rect
		icon string
	}{{g.minus, "minus"}, {g.plus, "plus"}} {
		p.round(btn.r, float32(m.RadiusControl()), t.WindowBackground)
		p.outline(btn.r, float32(m.RadiusControl()), float32(m.Hairline()), t.BorderStrong)
		glyph, _ := icons.Glyph(btn.icon)
		il := ui.Fonts.Layout(span(string(glyph), ui.Icon(m.IconSizeSmall(), t.TextPrimary)), text.Options{})
		iw, ih := il.Size()
		il.Draw(gc, btn.r.X+(btn.r.Width-iw)/2, btn.r.Y+(btn.r.Height-ih)/2)
	}
	x = g.minus.X - float32(m.SpaceLoose())
	var widths []float32
	var layouts []*text.Layout
	total := float32(0)
	for _, tl := range themeLabels {
		weight := text.Regular
		if tl.id == ui.Theme.ThemeID() {
			weight = text.Semibold
		}
		ll := ui.Fonts.Layout(span(tl.label, ui.Chrome(m.Body(), weight, t.TextPrimary)), text.Options{})
		w, _ := ll.Size()
		w += float32(m.SpaceLoose()) * 2
		widths, layouts = append(widths, w), append(layouts, ll)
		total += w
	}
	x -= total
	group := geom.NewRect(x, cy, total, ch)
	p.round(group, float32(m.RadiusControl()), t.WindowBackground)
	for i, tl := range themeLabels {
		seg := geom.NewRect(x, cy, widths[i], ch)
		g.segments[tl.id] = seg
		if tl.id == ui.Theme.ThemeID() {
			p.round(seg, float32(m.RadiusControl()), t.ChipBackground)
			p.outline(seg, float32(m.RadiusControl()), float32(m.FocusRingWidth()/2+1), t.BorderStrong)
		}
		_, lh := layouts[i].Size()
		layouts[i].Draw(gc, x+float32(m.SpaceLoose()), cy+(ch-lh)/2)
		x += widths[i]
	}
	p.outline(group, float32(m.RadiusControl()), float32(m.Hairline()), t.BorderStrong)
}

func (g *gallery) headerClick(where geom.Point, button, clickCount int, mods mod.Modifiers) bool {
	for id, r := range g.segments {
		if where.In(r) {
			g.ui.Theme.SetThemeID(id)
			return true
		}
	}
	switch {
	case where.In(g.minus):
		g.ui.Interface.SetFontSize(g.ui.Interface.FontSize() - 1)
	case where.In(g.plus):
		g.ui.Interface.SetFontSize(g.ui.Interface.FontSize() + 1)
	default:
		return false
	}
	return true
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

func (g *gallery) drawSidebar(gc *unison.Canvas, r geom.Rect) {
	ui := g.ui
	m := ui.Interface
	t := ui.Theme.Tokens()
	p := painter{gc, ui}
	b := g.side.ContentRect(false)
	p.fill(b, t.ListBackground)
	p.fill(geom.NewRect(b.Width-float32(m.Hairline()), 0, float32(m.Hairline()), b.Height), t.Border)
	y := float32(m.Space())
	for _, gr := range groups {
		_, h := p.text(span(gr.name, ui.Chrome(m.Small(), text.Semibold, t.TextMuted)), float32(m.Space()), y, 0)
		y += h + float32(m.SpaceSnug())
		for _, name := range gr.pages {
			row := geom.NewRect(0, y, b.Width-float32(m.Hairline()), float32(m.RowHeightSlim()))
			g.rows[name] = row
			color := t.TextDisabled
			if hasPage(name) {
				color = t.TextPrimary
			}
			if name == g.page {
				p.fill(row, t.SelectionTint)
			}
			l := ui.Fonts.Layout(span(name, ui.Chrome(m.Body(), text.Regular, color)), text.Options{})
			_, lh := l.Size()
			l.Draw(gc, float32(m.SpaceLoose()+m.SpaceTight()), y+(row.Height-lh)/2)
			y += row.Height
			if y > b.Height {
				return
			}
		}
		y += float32(m.SpaceNear())
	}
}

func (g *gallery) drawStatus(gc *unison.Canvas, r geom.Rect) {
	ui := g.ui
	m := ui.Interface
	t := ui.Theme.Tokens()
	p := painter{gc, ui}
	b := g.status.ContentRect(false)
	p.fill(b, t.FooterBackground)
	p.fill(geom.NewRect(0, 0, b.Width, float32(m.Hairline())), t.Border)
	st := ui.Chrome(m.Caption(), text.Regular, t.TextMuted)
	built := len(catalog)
	left := ui.Fonts.Layout(span(fmt.Sprintf("%s  ·  %d of 74 components built", g.page, built), st), text.Options{})
	_, lh := left.Size()
	left.Draw(gc, float32(m.Space()), (b.Height-lh)/2)
	right := ui.Fonts.Layout(span(fmt.Sprintf("%s  ·  %d px", tokens.DisplayName(ui.Theme.ResolvedTheme()), m.FontSize()), st), text.Options{})
	rw, rh := right.Size()
	right.Draw(gc, b.Width-rw-float32(m.Space()), (b.Height-rh)/2)
}
