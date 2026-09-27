package main

// The Content group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

func dividerBothWays(ui *kvitui.UI) unison.Paneler {
	down := kvitui.NewDivider(ui)
	down.Vertical = true
	down.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Fill, VGrab: true})
	beside := kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewLabel(ui, "left"), down, kvitui.NewLabel(ui, "right"))
	beside.SetLayout(kvitui.AtLeast(ui, kvitui.Px(24), beside.Layout()))
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewDivider(ui), beside))
}

func panelWithRules(ui *kvitui.UI) unison.Paneler {
	p := kvitui.NewPanel(ui)
	p.RuleTop, p.RuleBottom = true, true
	label := kvitui.NewLabel(ui, "A panel")
	label.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Middle, VAlign: align.Middle, HGrab: true, VGrab: true})
	p.AddChild(label)
	p.SetLayout(kvitui.AtLeast(ui, kvitui.Px(60), &unison.FlexLayout{Columns: 1}))
	return kvitui.FullWidth(p)
}
