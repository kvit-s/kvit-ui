package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
	"github.com/richardwilkes/unison/enums/slant"
	"github.com/richardwilkes/unison/enums/spacing"
	"github.com/richardwilkes/unison/enums/weight"
)

// Field is a single line of text the reader types.
//
// The editing is unison's field: the caret, selection, the clipboard, undo
// and what a screen reader is told about the text. The Kvit part is how it
// looks and what it says around the text: an outline in the strong border
// colour, since a field outlined in the ordinary border colour is below 3:1
// in every theme and its edges cannot be found; the focus ring whenever it
// holds the focus, since typing goes there; a placeholder in the faint text
// colour; and an error that is a message under the field as well as a red
// outline, never the outline alone. A field that turns red and says nothing
// has told a reader who cannot see the red nothing at all.
//
// The typed text is drawn by unison's own text engine rather than the Kvit
// text layer, so it can differ slightly from the labels around it.
type Field struct {
	unison.Panel
	ui *UI
	// Label names the field for a screen reader.
	Label string
	// Placeholder is shown, faint, while the field is empty.
	Placeholder string
	// Error is a message under the field. Non-empty is what "in error"
	// means; there is no separate flag, so the two cannot disagree.
	Error string
	// Shortcut names the keys that reach the field, where something does. It
	// is shown with the label as the field's tooltip, so a route to the field
	// is written down somewhere; "" shows no tooltip.
	Shortcut string
	// OnChange runs after every change to the text.
	OnChange func(text string)

	edit    *unison.Field
	message *Label
	// padLeft and padRight are the space between the outline and the text,
	// which a search field widens for its symbol and its clear button.
	padLeft, padRight func() float32
	// decorate draws over the field, as a search field draws its symbol.
	decorate func(gc *unison.Canvas, box geom.Rect)
	// place lays out anything else the field holds, given the box the text
	// field takes, as a search field places its clear button.
	place func(box geom.Rect)
	font  fieldFontKey
}

// fieldFontKey is the family and pixel size the field's unison font was made
// for.
type fieldFontKey struct {
	family string
	px     int
}

// unisonFont is a unison font of a family at an em size in pixels, as the
// rest of the chrome is sized. unison sizes a font by the height of its
// capital letters rather than its em, so the cap height is found from the
// line height: unison's at a trial size against the text layer's per em.
func unisonFont(ui *UI, family string, px int) unison.Font {
	face := unison.FontFaceDescriptor{Family: family, Weight: weight.Regular, Spacing: spacing.Standard, Slant: slant.Upright}.Face()
	if face == nil {
		return unison.FieldFont
	}
	const trial = 10
	probe := face.Font(trial)
	_, big := ui.Fonts.Layout([]text.Span{{Text: "Hg", Style: text.Style{Family: family, Size: 1000}}}, text.Options{}).Size()
	if probe.LineHeight() <= 0 || big <= 0 {
		return probe
	}
	em := probe.LineHeight() / (big / 1000)
	return face.Font(trial * float32(px) / em)
}

// NewField returns an empty field.
func NewField(ui *UI) *Field {
	f := &Field{ui: ui}
	f.Self = f
	f.initField(ui)
	return f
}

func (f *Field) initField(ui *UI) {
	f.padLeft = func() float32 { return float32(ui.Interface.SpaceNear()) }
	f.padRight = f.padLeft
	e := unison.NewField()
	f.edit = e
	unison.UninstallFocusBorders(e, e)
	e.SetBorder(fieldPadding{f})
	e.ClientData()[ringOwnerKey] = f
	e.DrawCallback = f.draw
	e.GainedFocusCallback = func() {
		e.DefaultFocusGained()
		if w := e.Window(); w != nil {
			ui.watchWindow(w)
			w.MarkForRedraw()
		}
	}
	e.LostFocusCallback = func() {
		e.DefaultFocusLost()
		if w := e.Window(); w != nil {
			w.MarkForRedraw()
		}
	}
	e.FrameChangeCallback = func() {
		if w := e.Window(); w != nil {
			ui.watchWindow(w)
		}
	}
	e.ModifiedCallback = func(_, after *unison.FieldState) {
		if f.OnChange != nil {
			f.OnChange(after.Text)
		}
		f.MarkForLayoutAndRedraw()
	}
	e.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		e.Tooltip = nil
		if f.Shortcut != "" {
			say := f.Shortcut
			if f.Label != "" {
				say = f.Label + " · " + f.Shortcut
			}
			e.Tooltip = newTooltip(ui, say, "")
		}
		return e.RectToRoot(e.ContentRect(true))
	}
	e.Accessibility.Callback = func(n *accessibility.Node) {
		n.Name = f.Label
		n.Description = f.Error
		n.Placeholder = f.Placeholder
		n.Invalid = f.Error != ""
	}
	f.message = NewLabel(ui, "")
	f.message.Role, f.message.Ink, f.message.Wrap = RoleCaption, InkDanger, true
	f.AddChild(e)
	f.AddChild(f.message)
	// The field's own panel says nothing; the text field inside it is what a
	// screen reader meets, under the field's label.
	f.Accessibility.Role = role.None
	f.SetLayout(fieldLayout{f})
}

// Text is what is typed in the field.
func (f *Field) Text() string { return f.edit.Text() }

// SetText replaces what is typed in the field.
func (f *Field) SetText(s string) { f.edit.SetText(s) }

// Edit is unison's field inside, for what the Kvit field does not cover.
func (f *Field) Edit() *unison.Field { return f.edit }

// Focus puts the caret in the field.
func (f *Field) Focus() {
	if f.Window() != nil {
		f.edit.RequestFocus()
		return
	}
	unison.InvokeTask(func() {
		if f.Window() != nil {
			f.edit.RequestFocus()
		}
	})
}

// style keeps unison's field in the current theme and interface size.
func (f *Field) style() {
	ui, t := f.ui, f.ui.Theme.Tokens()
	key := fieldFontKey{ui.Interface.FontFamily(), ui.Size(RoleBody)}
	if key != f.font {
		f.font = key
		f.edit.Font = unisonFont(ui, key.family, key.px)
	}
	// The ground and the outline are drawn here, rounded; unison's own
	// square ground is left transparent.
	f.edit.EditableInk, f.edit.BackgroundInk, f.edit.ErrorInk = unison.Transparent, unison.Transparent, unison.Transparent
	f.edit.OnEditableInk, f.edit.OnErrorInk = Color(t.TextPrimary), Color(t.TextPrimary)
	f.edit.OnBackgroundInk = Color(t.TextDisabled)
	f.edit.SelectionInk, f.edit.OnSelectionInk = Color(t.SelectionActiveTint), Color(t.TextPrimary)
}

func (f *Field) draw(gc *unison.Canvas, dirty geom.Rect) {
	ui, t, m := f.ui, f.ui.Theme.Tokens(), f.ui.Interface
	f.style()
	p := painterFor(gc, ui)
	box := f.edit.ContentRect(true)
	radius := float32(m.RadiusControl())
	ground := t.PopupBackground
	if !f.edit.Enabled() {
		ground = t.ChipBackground
	}
	p.round(box, radius, ground)
	edge := t.BorderStrong
	switch {
	case f.Error != "":
		edge = t.Danger
	case f.edit.Focused():
		edge = t.FocusRing
	}
	p.outline(box, radius, float32(m.Hairline()), edge)
	if !f.edit.Enabled() {
		// unison greys a disabled field's text a second time, which leaves it
		// fainter than the disabled colour; it is drawn here instead.
		f.edit.OnBackgroundInk = unison.Transparent
	}
	gc.Save()
	f.edit.DefaultDraw(gc, dirty)
	gc.Restore()
	if !f.edit.Enabled() && f.edit.Text() != "" {
		inner := f.edit.ContentRect(false)
		l := ui.Fonts.Layout([]text.Span{{Text: f.edit.Text(), Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.TextDisabled)}},
			text.Options{MaxWidth: inner.Width, Elide: true})
		_, h := l.Size()
		l.Draw(gc, inner.X, box.Y+(box.Height-h)/2)
	}
	if f.edit.Text() == "" && f.Placeholder != "" {
		inner := f.edit.ContentRect(false)
		l := ui.Fonts.Layout([]text.Span{{Text: f.Placeholder, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.TextFaint)}},
			text.Options{MaxWidth: inner.Width, Elide: true})
		_, h := l.Size()
		l.Draw(gc, inner.X, box.Y+(box.Height-h)/2)
	}
	if f.decorate != nil {
		f.decorate(gc, box)
	}
}

// focusRing puts the ring around the field whenever it holds the focus,
// however the focus arrived: the caret is there, and so is what is typed.
func (f *Field) focusRing() (*unison.Panel, float32, bool) {
	return f.edit.AsPanel(), float32(f.ui.Interface.RadiusControl()), f.edit.Focused()
}

// fieldPadding centres the line of text on the field's height, a near space
// in from the outline at each end.
type fieldPadding struct{ f *Field }

func (b fieldPadding) Insets() geom.Insets {
	f := b.f
	f.style()
	v := max(0, (float32(f.ui.Interface.ControlHeight())-f.edit.Font.LineHeight())/2)
	return geom.Insets{Top: v, Bottom: v, Left: f.padLeft(), Right: f.padRight()}
}

func (b fieldPadding) Draw(*unison.Canvas, geom.Rect) {}

// fieldLayout is the field, a control tall, and the error message under it
// when there is one, as wide as the field.
type fieldLayout struct{ f *Field }

func (l fieldLayout) message(width float32) float32 {
	f := l.f
	f.message.Text = f.Error
	f.message.Hidden = f.Error == ""
	if f.message.Hidden {
		return 0
	}
	_, p, _ := f.message.Sizes(geom.NewSize(width, 0))
	return float32(f.ui.Interface.SpaceTight()) + p.Height
}

func (l fieldLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.f.ui.Interface
	w := float32(m.Px(200))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := float32(m.ControlHeight()) + l.message(w)
	return geom.NewSize(float32(m.Px(60)), h), geom.NewSize(float32(m.Px(200)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l fieldLayout) PerformLayout(target *unison.Panel) {
	f := l.f
	r := target.ContentRect(false)
	ch := float32(f.ui.Interface.ControlHeight())
	f.edit.SetFrameRect(geom.NewRect(r.X, r.Y, r.Width, ch))
	if mh := l.message(r.Width); mh > 0 {
		gap := float32(f.ui.Interface.SpaceTight())
		f.message.SetFrameRect(geom.NewRect(r.X, r.Y+ch+gap, r.Width, mh-gap))
	}
	if f.place != nil {
		f.place(f.edit.FrameRect())
	}
}

// SetEnabled enables or disables the field, the text field inside included.
func (f *Field) SetEnabled(enabled bool) {
	f.Panel.SetEnabled(enabled)
	f.edit.SetEnabled(enabled)
}
