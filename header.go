package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Header is the strip across the top of the window: what the application
// is, where the reader is in it, and the actions that apply everywhere. It
// has three places in a fixed order, the wordmark on the left, navigation in
// the middle and actions on the right, because a reader who has learned
// where the back control is on one screen should find it in the same place
// on every screen and in every application.
type Header struct {
	Panel
	navigation, actions *unison.Panel
}

// NewHeader returns a header with a wordmark ("" for none), navigation in
// the middle (nil for none) and actions at the right, a near space apart.
func NewHeader(ui *UI, wordmark string, navigation unison.Paneler, actions ...unison.Paneler) *Header {
	h := &Header{}
	h.ui = ui
	h.Self = h
	h.DrawCallback = h.draw
	h.RuleBottom = true
	h.SetBorder(sides{ui, SizeViewMargin})
	row := &spaced{FlexLayout: unison.FlexLayout{VAlign: align.Middle}, ui: ui, gap: SizeColumnGap, horizontal: true}
	h.SetLayout(AtLeast(ui, SizeHeaderHeight, row))
	if wordmark != "" {
		mark := NewLabel(ui, wordmark)
		mark.Role, mark.Weight = RoleHeadline, text.Bold
		mark.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
		h.AddChild(mark)
	}
	// Navigation takes the width left over, and its contents sit at its left,
	// centred on the header's height.
	h.navigation = unison.NewPanel()
	h.navigation.SetLayout(&unison.FlexLayout{Columns: 1, HAlign: align.Start, VAlign: align.Middle})
	h.navigation.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
	if navigation != nil {
		// As wide as the slot, so navigation can put things at its right
		// end; a row of tabs still starts at the left.
		navigation.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Middle, HGrab: true, VGrab: true})
		h.navigation.AddChild(navigation)
	}
	h.AddChild(h.navigation)
	h.actions = Row(ui, SizeSpaceNear, actions...)
	h.actions.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
	h.AddChild(h.actions)
	row.Columns = len(h.Children())
	return h
}

// ProvideAccessibility names the header, which holds the application's
// top-level navigation.
func (h *Header) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.TabList
	}
	n.Name = "Application header"
}

// sides is a border that is only space at the left and right, a Measure
// wide.
type sides struct {
	ui *UI
	m  Measure
}

func (s sides) Insets() geom.Insets {
	w := float32(s.m.Of(s.ui))
	return geom.Insets{Left: w, Right: w}
}
func (s sides) Draw(*unison.Canvas, geom.Rect) {}
