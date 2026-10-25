package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// SettingRow is one setting on a settings page: its name and, where needed,
// one short line saying more, at the leading side, and the control that
// changes it at the trailing edge (LAY-6). A setting applies the moment it
// changes, with no OK, Cancel or Apply, so a setting that is on or off is a
// Switch (WIN-16, WIN-23). The description is one line; a longer explanation
// belongs behind a Hint (WRD-22). The row names its control for a screen
// reader when the control has no name of its own.
type SettingRow struct {
	unison.Panel
	ui *UI
	// Label names the setting.
	Label string
	// Description says in one short line what the setting does, "" for none.
	Description string

	name, desc *Label
	control    unison.Paneler
}

// NewSettingRow returns a row for one setting.
func NewSettingRow(ui *UI, label, description string, control unison.Paneler) *SettingRow {
	r := &SettingRow{ui: ui, Label: label, Description: description, control: control}
	r.Self = r
	r.name = NewLabel(ui, "")
	r.desc = NewLabel(ui, "")
	r.desc.Role, r.desc.Ink = RoleSmall, InkTextMuted
	r.SetLayout(syncing{Layout: settingRowLayout{r}, sync: r.sync})
	r.sync()
	return r
}

func (r *SettingRow) sync() {
	r.name.Text, r.desc.Text = r.Label, r.Description
	parts := []unison.Paneler{r.name}
	if r.Description != "" {
		parts = append(parts, r.desc)
	}
	showOnly(r.AsPanel(), append(parts, r.control)...)
	// The control is named by the setting where it has no name of its own.
	switch c := r.control.AsPanel().Self.(type) {
	case *Switch:
		if c.Text == "" {
			c.Name = r.Label
		}
	case *Select:
		if c.Label == "" {
			c.Label = r.Label
		}
	case *Segmented:
		if c.Label == "" {
			c.Label = r.Label
		}
	case *Field:
		if c.Label == "" {
			c.Label = r.Label
		}
	case *Stepper:
		if c.Label == "" {
			c.Label = r.Label
		}
	}
}

type settingRowLayout struct{ r *SettingRow }

// arrange places the name and description at the leading side and the
// control at the trailing edge, centred on the row, and returns the row's
// height.
func (l settingRowLayout) arrange(box geom.Rect, set bool) float32 {
	r, m := l.r, l.r.ui.Interface
	_, cp, _ := r.control.AsPanel().Sizes(geom.Size{})
	gap := float32(m.ColumnGap())
	textWidth := max(0, box.Width-cp.Width-gap)
	_, np, _ := r.name.Sizes(geom.NewSize(textWidth, 0))
	textHeight := np.Height
	var dp geom.Size
	if r.Description != "" {
		_, dp, _ = r.desc.Sizes(geom.NewSize(textWidth, 0))
		textHeight += float32(m.SpaceTight()) + dp.Height
	}
	h := max(float32(m.RowHeight()), textHeight+float32(2*m.SpaceNear()), cp.Height)
	if set {
		top := box.Y + (h-textHeight)/2
		r.name.SetFrameRect(geom.NewRect(box.X, top, textWidth, np.Height))
		if r.Description != "" {
			r.desc.SetFrameRect(geom.NewRect(box.X, top+np.Height+float32(m.SpaceTight()), textWidth, dp.Height))
		}
		r.control.AsPanel().SetFrameRect(geom.NewRect(box.Right()-cp.Width, box.Y+(h-cp.Height)/2, cp.Width, cp.Height))
	}
	return h
}

func (l settingRowLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.r.ui.Interface
	w := float32(m.Px(480))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := l.arrange(geom.NewRect(0, 0, w, 0), false)
	return geom.NewSize(float32(m.Px(240)), h), geom.NewSize(float32(m.Px(480)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l settingRowLayout) PerformLayout(target *unison.Panel) {
	l.arrange(target.ContentRect(false), true)
}

// ProvideAccessibility names the row as a group by its setting.
func (r *SettingRow) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = r.Label
}

// About is what an About section says about the app (CMD-3): its name and
// version, and where to go for more.
type About struct {
	// Name is the app's name; the UI's AppName unless set.
	Name string
	// Version is the version the app is, "2.4.1".
	Version string
	// Details are further lines of facts, such as a build or a copyright.
	Details []string
	// Links are the places to go: the website, a place to report a problem
	// or ask for a feature, the third-party notices, the legal terms.
	Links []AboutLink
}

// AboutLink is one place an About section links to.
type AboutLink struct {
	// Text says where it goes, in sentence case.
	Text string
	// OnActivate follows it.
	OnActivate func()
}

// NewAboutSection returns the About section that ends the settings page on
// the Windows form (CMD-3, WIN-16): a heading naming the app, its version
// and facts, and its links. On the macOS form About is in the app menu, and
// the app shows the same section in a window of its own from there.
func NewAboutSection(ui *UI, a About) *unison.Panel {
	name := a.Name
	if name == "" {
		name = ui.AppName
	}
	heading := NewSectionHeading(ui, "About")
	title := NewLabel(ui, name)
	title.Role, title.Weight = RoleStrong, text.Bold
	parts := []unison.Paneler{FullWidth(heading), title}
	if a.Version != "" {
		version := NewLabel(ui, "Version "+a.Version)
		version.Ink = InkTextSecondary
		parts = append(parts, version)
	}
	for _, d := range a.Details {
		l := NewLabel(ui, d)
		l.Role, l.Ink = RoleSmall, InkTextMuted
		parts = append(parts, l)
	}
	for _, link := range a.Links {
		l := NewLink(ui, link.Text)
		l.OnActivate = link.OnActivate
		parts = append(parts, Left(l))
	}
	section := Column(ui, SizeSpaceNear, parts...)
	section.Accessibility.Role, section.Accessibility.Name = role.Group, "About "+name
	return section
}
