package main

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

// buildComponentPage lays out a component's page as the Qt gallery does: a
// view head with the name and the summary, and for each specimen its caption,
// the live specimen in a frame, and its code.
func buildComponentPage(ui *kvitui.UI, e entry) *unison.Panel {
	head := kvitui.NewViewHead(ui, e.name)
	head.Subtitle = e.summary
	head.Padding = kvitui.Px(0) // the page already has the view margin
	parts := []unison.Paneler{head}
	for _, s := range e.specimens {
		caption := kvitui.NewLabel(ui, s.caption)
		caption.Role, caption.Ink = kvitui.RoleSmall, kvitui.InkTextMuted
		parts = append(parts, kvitui.Column(ui, kvitui.SizeSpaceNear, caption, specimenFrame(ui, s.build(ui)), codeBlock(ui, sourceOf(s.build))))
	}
	page := kvitui.Column(ui, kvitui.SizeSpaceLoose, parts...)
	page.SetBorder(kvitui.Padding(ui, kvitui.SizeViewMargin))
	return page
}

// specimenFrame is the box a live specimen sits in: the specimen at its top
// left, a loose space from each edge, and the box at least a row tall.
func specimenFrame(ui *kvitui.UI, content unison.Paneler) *unison.Panel {
	f := unison.NewPanel()
	f.SetLayout(&unison.FlexLayout{Columns: 1})
	content.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Start, VAlign: align.Start})
	f.AddChild(content)
	f.SetBorder(kvitui.Padding(ui, kvitui.SizeSpaceLoose))
	f.SetLayout(kvitui.AtLeast(ui, kvitui.SizeRowHeight, f.Layout()))
	f.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		t := ui.Theme.Tokens()
		p := painter{gc, ui}
		r := f.ContentRect(true)
		p.round(r, float32(ui.Interface.RadiusCard()), t.ListBackground)
		p.outline(r, float32(ui.Interface.RadiusCard()), float32(ui.Interface.Hairline()), t.Border)
	}
	return f
}

// codeBlock shows a specimen's source in the monospace family.
func codeBlock(ui *kvitui.UI, src string) *unison.Panel {
	c := unison.NewPanel()
	layout := func() *text.Layout {
		return ui.Fonts.Layout(span(src, ui.Mono(ui.Size(kvitui.RoleSmall), ui.Theme.Tokens().TextSecondary)), text.Options{})
	}
	pad := func() float32 { return float32(ui.Interface.Space()) }
	c.SetSizer(func(hint geom.Size) (geom.Size, geom.Size, geom.Size) {
		// As wide as the column; a long line is cut at the edge rather than
		// widening the page.
		_, h := layout().Size()
		pref := geom.NewSize(0, h+2*pad())
		return pref, pref, geom.NewSize(unison.DefaultMaxSize, pref.Height)
	})
	c.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		painter{gc, ui}.round(c.ContentRect(false), float32(ui.Interface.RadiusBar()), ui.Theme.Tokens().CodePanelBackground)
		layout().Draw(gc, pad(), pad())
	}
	return c
}

// drawnPage wraps a page drawn by a painter function, such as Foundations,
// in a panel as tall as it needs.
func drawnPage(ui *kvitui.UI, draw func(p painter, x, y, width float32) float32) *unison.Panel {
	p := unison.NewPanel()
	margin := func() float32 { return float32(ui.Interface.ViewMargin()) }
	p.SetSizer(func(hint geom.Size) (geom.Size, geom.Size, geom.Size) {
		w := max(hint.Width, 400)
		h := draw(painter{nil, ui}, margin(), margin(), w-2*margin()) + 2*margin()
		return geom.NewSize(400, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
	})
	p.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		b := p.ContentRect(false)
		draw(painter{gc, ui}, margin(), margin(), b.Width-2*margin())
	}
	return p
}
