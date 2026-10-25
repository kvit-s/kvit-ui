package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// FormRow is a field with its visible label (LAY-6, WRD-19): above the field
// on the Windows form, which Linux takes, and beside it on the macOS form.
// The label names the field for a screen reader, so the two cannot differ;
// it ends without a colon, and with an asterisk when the field must be
// filled, the one marking on every system (CON-4). On the Windows form the
// label holds the field's access key: Alt and its letter put the keyboard in
// the field (KEY-15).
type FormRow struct {
	unison.Panel
	ui *UI
	// Label names the field, in sentence case, with an "&" before the
	// access key's letter where the app chooses one.
	Label string
	// Required marks a field that must be filled.
	Required bool
	// LabelWidth, on the macOS form, is the width of the label's column, so
	// the fields of a form line up; nil takes the label's own width.
	LabelWidth Measure

	label *formLabel
	field unison.Paneler
}

// formLabel is a form row's label, which holds the field's access key.
type formLabel struct {
	unison.Panel
	row *FormRow
}

// NewFormRow returns a field with its label.
func NewFormRow(ui *UI, label string, field unison.Paneler) *FormRow {
	r := &FormRow{ui: ui, Label: label, field: field}
	r.Self = r
	r.label = &formLabel{row: r}
	r.label.Self = r.label
	r.label.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		w, h := r.layout().Size()
		size := geom.NewSize(w, h)
		return size, size, size
	})
	r.label.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		l := r.layout()
		b := r.label.ContentRect(false)
		w, h := l.Size()
		x := b.X
		if ui.form == FormMac {
			// Beside the field, the label is set against it.
			x = b.Right() - w
		}
		l.Draw(gc, x, b.Y+(b.Height-h)/2)
	}
	r.AddChild(r.label)
	r.AddChild(field)
	r.SetLayout(syncing{Layout: formRowLayout{r}, sync: r.sync})
	r.sync()
	return r
}

// words are what the label says: its words without the access-key marker,
// and the asterisk of a required field.
func (r *FormRow) words() string {
	w := plainWords(r.Label)
	if r.Required {
		w += " *"
	}
	return w
}

func (r *FormRow) layout() *text.Layout {
	ui := r.ui
	st := ui.Chrome(ui.Size(RoleBody), text.Regular, ui.Theme.Tokens().TextSecondary)
	return ui.Fonts.Layout(ui.accessSpans(r.label.AsPanel(), r.words(), st), text.Options{})
}

// sync names the field by the label, and tells a screen reader it is
// required.
func (r *FormRow) sync() {
	name := plainWords(r.Label)
	switch f := r.field.AsPanel().Self.(type) {
	case *Field:
		f.Label, f.Required = name, r.Required
	case *SearchField:
		f.Label, f.Required = name, r.Required
	case *NumberField:
		f.Label, f.Required = name, r.Required
	case *MoneyField:
		f.Label, f.Required = name, r.Required
	case *TextArea:
		f.Label, f.Required = name, r.Required
	case *TypeAhead:
		f.Label, f.Required = name, r.Required
	case *Select:
		f.Label = name
	}
}

// accessWords are the label's words for the field's access key.
func (l *formLabel) accessWords() string { return l.row.Label }

// accessPress puts the keyboard in the field.
func (l *formLabel) accessPress() {
	f := l.row.field.AsPanel()
	if fo, ok := f.Self.(interface{ Focus() }); ok {
		fo.Focus()
		return
	}
	if c := f.FirstFocusableChild(); c != nil {
		c.RequestFocus()
	}
}

// ProvideAccessibility leaves the label out: the field it names is named by
// its words.
func (l *formLabel) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	b.Node().Role = role.None
}

type formRowLayout struct{ r *FormRow }

// arrange places the label and the field and returns the row's height:
// the label above the field on the Windows form, and on the macOS form
// beside it, in a column of LabelWidth, centred on the field's first line.
func (l formRowLayout) arrange(box geom.Rect, set bool) float32 {
	r, m := l.r, l.r.ui.Interface
	_, lp, _ := r.label.Sizes(geom.Size{})
	if r.ui.form == FormMac {
		lw := lp.Width
		if r.LabelWidth != nil {
			lw = float32(r.LabelWidth.Of(r.ui))
		}
		gap := float32(m.Space())
		fw := max(0, box.Width-lw-gap)
		_, fp, _ := r.field.AsPanel().Sizes(geom.NewSize(fw, 0))
		if set {
			first := float32(m.ControlHeight())
			r.label.SetFrameRect(geom.NewRect(box.X, box.Y+(first-lp.Height)/2, lw, lp.Height))
			r.field.AsPanel().SetFrameRect(geom.NewRect(box.X+lw+gap, box.Y, fw, fp.Height))
		}
		return max(fp.Height, lp.Height)
	}
	gap := float32(m.SpaceTight())
	_, fp, _ := r.field.AsPanel().Sizes(geom.NewSize(box.Width, 0))
	if set {
		r.label.SetFrameRect(geom.NewRect(box.X, box.Y, min(lp.Width, box.Width), lp.Height))
		r.field.AsPanel().SetFrameRect(geom.NewRect(box.X, box.Y+lp.Height+gap, box.Width, fp.Height))
	}
	return lp.Height + gap + fp.Height
}

func (l formRowLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.r.ui.Interface
	w := float32(m.Px(280))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := l.arrange(geom.NewRect(0, 0, w, 0), false)
	return geom.NewSize(float32(m.Px(120)), h), geom.NewSize(float32(m.Px(280)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l formRowLayout) PerformLayout(target *unison.Panel) {
	l.arrange(target.ContentRect(false), true)
}

// ProvideAccessibility leaves the row out: the field is what a reader
// meets.
func (r *FormRow) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	b.Node().Role = role.None
}
