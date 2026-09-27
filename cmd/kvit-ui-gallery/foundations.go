package main

import (
	"fmt"
	"reflect"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
)

// colorGroups are the token groups as the Tokens struct lists them, with the
// number of fields in each.
var colorGroups = []struct {
	name  string
	count int
}{
	{"Surfaces", 8}, {"Text", 7}, {"Lines and glyphs", 4}, {"Interactive tints", 7},
	{"Accent and meaning", 6}, {"Inline styling", 9}, {"Code", 5}, {"Callout", 1}, {"Dashboard", 9},
}

// drawFoundations draws the foundations page, the design values themselves,
// with its top-left at (x, y) in a column width wide, and returns its height.
func drawFoundations(p painter, x, y, width float32) float32 {
	ui := p.ui
	m := ui.Interface
	t := ui.Theme.Tokens()
	top := y
	gap := float32(m.SpaceLoose())
	primary, secondary, muted := t.TextPrimary, t.TextSecondary, t.TextMuted

	heading := func(title, about string) {
		_, h := p.text(span(title, ui.Chrome(m.Title(), text.Bold, primary)), x, y, width)
		y += h + float32(m.SpaceTight())
		if about != "" {
			_, h = p.text(span(about, ui.Chrome(m.Small(), text.Regular, secondary)), x, y, width)
			y += h
		}
		y += float32(m.Space())
	}

	_, h := p.text(span("Foundations", ui.Chrome(m.Headline(), text.Bold, primary)), x, y, width)
	y += h + float32(m.SpaceTight())
	_, h = p.text(span(fmt.Sprintf("The design values every component draws with, in the %s theme at interface size %d px. "+
		"Nothing on this page is written as a literal: each colour, size and gap comes from the tokens package.",
		tokens.DisplayName(ui.Theme.ResolvedTheme()), m.FontSize()),
		ui.Chrome(m.Small(), text.Regular, secondary)), x, y, width)
	y += h + gap*2

	// Colours, group by group, in the order the Tokens struct lists them.
	heading("Colours", "Each theme's tokens. Border is decorative and deliberately below 3:1; BorderStrong marks where a control is.")
	v := reflect.ValueOf(t)
	typ := v.Type()
	colorType := reflect.TypeOf(palette.Color{})
	field := 0
	sw, sh := float32(m.Px(28)), float32(m.Px(20))
	cell := float32(m.Px(176))
	for _, g := range colorGroups {
		_, h := p.text(span(g.name, ui.Chrome(m.Caption(), text.Semibold, muted)), x, y, width)
		y += h + float32(m.SpaceSnug())
		cx, rowH := x, float32(0)
		for n := 0; n < g.count; field++ {
			if typ.Field(field).Type != colorType {
				continue
			}
			n++
			c := v.Field(field).Interface().(palette.Color)
			if cx+cell > x+width {
				cx, y = x, y+rowH+float32(m.SpaceNear())
				rowH = 0
			}
			r := geom.NewRect(cx, y, sw, sh)
			p.round(r, float32(m.RadiusChip()), c)
			p.outline(r, float32(m.RadiusChip()), float32(m.Hairline()), t.Border)
			_, h1 := p.text(span(typ.Field(field).Name, ui.Chrome(m.Caption(), text.Regular, primary)), cx+sw+float32(m.SpaceNear()), y, 0)
			_, h2 := p.text(span(c.Hex(), ui.Mono(m.Caption(), muted)), cx+sw+float32(m.SpaceNear()), y+h1, 0)
			rowH = max(rowH, sh, h1+h2)
			cx += cell
		}
		y += rowH + gap
	}

	// The three chart ramps.
	heading("Chart ramps", "Categorical says which series, sequential how much, diverging how far and which way. Each is checked against this theme's surfaces.")
	for _, ramp := range []struct {
		name  string
		steps []palette.Color
	}{{"Categorical", t.CategoricalRamp}, {"Sequential", t.SequentialRamp}, {"Diverging", t.DivergingRamp}} {
		_, h := p.text(span(ramp.name, ui.Chrome(m.Caption(), text.Semibold, muted)), x, y, 0)
		y += h + float32(m.SpaceSnug())
		bw := float32(m.Px(56))
		for i, c := range ramp.steps {
			p.round(geom.NewRect(x+float32(i)*(bw+float32(m.SpaceSnug())), y, bw, float32(m.BarHeightWide()*3)), float32(m.RadiusBar()), c)
		}
		y += float32(m.BarHeightWide()*3) + gap
	}

	// The chrome's seven type roles, then the document's scale.
	heading("Type roles", fmt.Sprintf("The chrome's seven roles in %s, derived from the interface size.", ui.Fonts.ResolveFamily(m.FontFamily())))
	roles := []struct {
		name string
		size int
	}{{"Display", m.Display()}, {"Headline", m.Headline()}, {"Title", m.Title()}, {"Strong", m.Strong()},
		{"Body", m.Body()}, {"Small", m.Small()}, {"Caption", m.Caption()}}
	label := float32(m.Px(120))
	for _, r := range roles {
		_, h1 := p.text(span(fmt.Sprintf("%s %d px", r.name, r.size), ui.Chrome(m.Caption(), text.Regular, muted)), x, y+float32(m.SpaceTight()), 0)
		_, h2 := p.text(span("Reconcile 1,284 transactions before Friday", ui.Chrome(r.size, text.Regular, primary)), x+label, y, width-label)
		y += max(h1, h2) + float32(m.SpaceSnug())
	}
	y += gap
	ty := ui.Typography
	heading("Document type", fmt.Sprintf("Body text and headings at base %d px, line height %.1f; separate from the interface size.", ty.BaseSize(), ty.LineHeight()))
	for _, r := range []struct {
		name string
		role tokens.FontRole
		w    int
	}{{"Heading 1", tokens.Heading1, text.Bold}, {"Heading 2", tokens.Heading2, text.Bold}, {"Heading 3", tokens.Heading3, text.Medium},
		{"Heading 4", tokens.Heading4, text.Medium}, {"Body", tokens.Body, text.Regular}, {"Code", tokens.Mono, text.Regular}} {
		st := ui.Chrome(ty.SizeForRole(r.role), r.w, primary)
		st.Family = ty.FontFamily()
		if r.role == tokens.Mono {
			st.Family = ty.MonoFamily()
		}
		_, h1 := p.text(span(fmt.Sprintf("%s %d px", r.name, ty.SizeForRole(r.role)), ui.Chrome(m.Caption(), text.Regular, muted)), x, y+float32(m.SpaceTight()), 0)
		l := ui.Fonts.Layout(span("A note about the quarter", st), text.Options{MaxWidth: width - label, LineHeight: float32(ty.LineHeight())})
		if p.gc != nil {
			l.Draw(p.gc, x+label, y)
		}
		_, h2 := l.Size()
		y += max(h1, h2)
	}
	y += gap

	// Spacing, heights and radii as bars and shapes.
	heading("Spacing and geometry", "The named gaps, row and control heights, and corner radii, all scaled by the one interface size.")
	bar := func(name string, value int, horizontal bool) {
		_, hl := p.text(span(fmt.Sprintf("%s %d", name, value), ui.Chrome(m.Caption(), text.Regular, muted)), x, y, 0)
		if horizontal {
			p.fill(geom.NewRect(x+label, y+float32(m.SpaceTight()), float32(value), float32(m.BarHeightWide())), t.Accent)
			y += max(hl, float32(m.BarHeightWide())) + float32(m.SpaceSnug())
		} else {
			p.fill(geom.NewRect(x+label, y, float32(m.Px(40)), float32(value)), t.PanelBackground)
			p.outline(geom.NewRect(x+label, y, float32(m.Px(40)), float32(value)), 0, float32(m.Hairline()), t.BorderStrong)
			y += max(hl, float32(value)) + float32(m.SpaceSnug())
		}
	}
	for _, s := range []struct {
		n string
		v int
	}{{"SpaceTight", m.SpaceTight()}, {"SpaceSnug", m.SpaceSnug()}, {"SpaceNear", m.SpaceNear()}, {"Space", m.Space()},
		{"SpaceWide", m.SpaceWide()}, {"SpaceLoose", m.SpaceLoose()}, {"ViewMargin", m.ViewMargin()}} {
		bar(s.n, s.v, true)
	}
	y += float32(m.Space())
	for _, s := range []struct {
		n string
		v int
	}{{"RowHeight", m.RowHeight()}, {"RowHeightSlim", m.RowHeightSlim()}, {"ControlHeight", m.ControlHeight()}, {"ChipHeight", m.ChipHeight()}} {
		bar(s.n, s.v, false)
	}
	rx := x + label
	for _, r := range []struct {
		n string
		v int
	}{{"Bar", m.RadiusBar()}, {"Chip", m.RadiusChip()}, {"Control", m.RadiusControl()}, {"Card", m.RadiusCard()}, {"Pill", m.RadiusPill()}} {
		box := geom.NewRect(rx, y, float32(m.Px(56)), float32(m.Px(32)))
		p.round(box, float32(r.v), t.ChipBackground)
		p.outline(box, float32(r.v), float32(m.Hairline()), t.BorderStrong)
		p.text(span(fmt.Sprintf("%s %d", r.n, r.v), ui.Chrome(m.Caption(), text.Regular, muted)), rx, y+box.Height+float32(m.SpaceTight()), 0)
		rx += box.Width + float32(m.SpaceLoose())
	}
	p.text(span("Radii", ui.Chrome(m.Caption(), text.Regular, muted)), x, y, 0)
	y += float32(m.Px(32)) + float32(m.Caption()) + gap*2

	// Every meaning name with the glyph it draws.
	heading("Icons", fmt.Sprintf("The %d meaning names components are written in, from the embedded Phosphor font.", len(icons.MeaningNames())))
	cx, cellW := x, float32(m.Px(150))
	rowH := float32(0)
	for _, name := range icons.MeaningNames() {
		if cx+cellW > x+width {
			cx, y, rowH = x, y+rowH+float32(m.SpaceSnug()), 0
		}
		g, _ := icons.Glyph(name)
		_, hi := p.text(span(string(g), ui.Icon(m.IconSize(), primary)), cx, y, 0)
		_, hn := p.text(span(name, ui.Chrome(m.Caption(), text.Regular, secondary)), cx+float32(m.IconSize()+m.SpaceNear()), y+float32(m.SpaceTight()), 0)
		rowH = max(rowH, hi, hn)
		cx += cellW
	}
	y += rowH + gap*2

	// Text the text package shapes.
	heading("Text", "Shaped with HarfBuzz: kerning, emoji sequences, other scripts and mixed styles on one line.")
	body := m.Title()
	for _, line := range [][]text.Span{
		span("AVATAR To Wa Ty — kerning pulls these pairs together", ui.Chrome(body, text.Regular, primary)),
		span("👍 👍🏽 🇺🇸 🇩🇪 👨‍👩‍👧 1️⃣ ❤️", ui.Chrome(body, text.Regular, primary)),
		{
			{Text: "Plain, ", Style: ui.Chrome(body, text.Regular, primary)},
			{Text: "bold", Style: ui.Chrome(body, text.Bold, primary)},
			{Text: ", ", Style: ui.Chrome(body, text.Regular, primary)},
			{Text: "italic", Style: func() text.Style { s := ui.Chrome(body, text.Regular, primary); s.Italic = true; return s }()},
			{Text: ", ", Style: ui.Chrome(body, text.Regular, primary)},
			{Text: "code", Style: func() text.Style {
				s := ui.Mono(m.Strong(), primary)
				s.Background = kvitui.TextColor(t.InlineCodeBackground)
				return s
			}()},
			{Text: " and a ", Style: ui.Chrome(body, text.Regular, primary)},
			{Text: "link", Style: func() text.Style { s := ui.Chrome(body, text.Regular, t.Link); s.Underline = true; return s }()},
		},
		span("Ελληνικά · Русский · العربية · עברית", ui.Chrome(body, text.Regular, primary)),
	} {
		_, h := p.text(line, x, y, width)
		y += h + float32(m.SpaceNear())
	}
	return y - top
}
