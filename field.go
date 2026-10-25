package kvitui

import (
	"strings"

	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/pathop"
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
	// Shortcut names the keys that reach the field, where something does. On
	// the Windows form it is shown in parentheses after the label as the
	// field's tooltip, "Search (Ctrl+F)", so a route to the field is written
	// down somewhere; the macOS form puts no shortcut in a tooltip (KEY-14).
	// "" shows no tooltip.
	Shortcut string
	// ReadOnly lets the text be selected and copied but not changed.
	ReadOnly bool
	// Required marks a field that must be filled: told to a screen reader,
	// and drawn as an asterisk after the label of the FormRow holding it
	// (CON-4).
	Required bool
	// Masked hides what is typed behind dots, for a secret; the text can be
	// neither copied nor cut (CTL-7).
	Masked bool
	// MaxLength is the most characters the field takes, 0 for no limit.
	// Text past it, typed or pasted, is kept and reported as the field's
	// error, rather than cut without a word (CTL-7).
	MaxLength int
	// OnChange runs after every change to the text.
	OnChange func(text string)

	edit    *unison.Field
	message *Label
	// clear empties a plain single-line field: shown on the Windows form
	// while the field holds the keyboard and some text (CTL-7, decided
	// Q-44). nil for the fields that have their own.
	clear *IconButton
	// overLimit is the error the length limit set, kept apart from the
	// application's own.
	overLimit string
	// padLeft and padRight are the space between the outline and the text,
	// which a search field widens for its symbol and its clear button.
	padLeft, padRight func() float32
	// decorate draws over the field, as a search field draws its symbol.
	decorate func(gc *unison.Canvas, box geom.Rect)
	// place lays out anything else the field holds, given the box the text
	// field takes, as a search field places its clear button.
	place func(box geom.Rect)
	font  fieldFontKey
	// area is the TextArea this field is, for one of several lines.
	area *TextArea
	// spoken, when set, is what a screen reader calls the field instead of
	// its label, as a money field adds its currency.
	spoken func() string
	// changed, when set, runs after every change to the text, before
	// OnChange, for a field that checks what is typed.
	changed func()
	// describe, when set, is what a screen reader is told about the field
	// while it has no error, as a type-ahead says how many suggestions it
	// has.
	describe func() string
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
	// unison knows families by their own names, not by generic ones such as
	// "monospace", so it is given the family the text layer resolves to.
	family = ui.Fonts.ResolveFamily(family)
	face := unison.FontFaceDescriptor{Family: family, Weight: weight.Regular, Spacing: spacing.Standard, Slant: slant.Upright}.Face()
	if face == nil {
		return unison.FieldFont
	}
	const trial = 10
	// probeSize is a size to measure a line at, large so rounding is small;
	// it is never drawn.
	const probeSize = 1000
	probe := face.Font(trial)
	_, big := ui.Fonts.Layout([]text.Span{{Text: "Hg", Style: text.Style{Family: family, Size: probeSize}}}, text.Options{}).Size()
	if probe.LineHeight() <= 0 || big <= 0 {
		return probe
	}
	em := probe.LineHeight() / (big / probeSize)
	return face.Font(trial * float32(px) / em)
}

// NewField returns an empty field.
func NewField(ui *UI) *Field {
	f := &Field{ui: ui}
	f.Self = f
	f.initField(ui, false)
	f.clear = NewIconButton(ui, "close", "Clear")
	// Not a stop in the Tab order: the keyboard empties a field another way.
	f.clear.SetFocusable(false)
	f.clear.OnClick = func() {
		f.SetText("")
		f.edit.RequestFocus()
	}
	f.AddChildAtIndex(f.clear, 0)
	f.padRight = func() float32 {
		if f.clearShown() {
			return float32(ui.Interface.ControlHeight())
		}
		return float32(ui.Interface.SpaceNear())
	}
	f.place = func(box geom.Rect) {
		f.clear.Hidden = !f.clearShown()
		c := float32(ui.Interface.ControlHeight())
		f.clear.SetFrameRect(geom.NewRect(box.Right()-c, box.Y+(box.Height-c)/2, c, c))
	}
	return f
}

// clearShown reports whether the clear button shows: on the Windows form,
// while the field holds the keyboard and some text it may change.
func (f *Field) clearShown() bool {
	return f.clear != nil && f.ui.form == FormWindows && f.edit.Focused() && f.edit.Text() != "" &&
		!f.ReadOnly && !f.Masked && f.edit.Enabled()
}

func (f *Field) initField(ui *UI, multi bool) {
	f.padLeft = func() float32 { return float32(ui.Interface.SpaceNear()) }
	f.padRight = f.padLeft
	e := unison.NewField()
	if multi {
		e = unison.NewMultiLineField()
	}
	f.edit = e
	unison.UninstallFocusBorders(e, e)
	e.SetBorder(fieldPadding{f})
	e.ClientData()[ringOwnerKey] = f
	e.DrawCallback = f.draw
	e.GainedFocusCallback = func() {
		e.DefaultFocusGained()
		f.MarkForLayoutAndRedraw()
		if w := e.Window(); w != nil {
			ui.watchWindow(w)
			w.MarkForRedraw()
		}
	}
	e.LostFocusCallback = func() {
		e.DefaultFocusLost()
		f.MarkForLayoutAndRedraw()
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
		f.checkLength()
		if f.changed != nil {
			f.changed()
		}
		if f.OnChange != nil {
			f.OnChange(after.Text)
		}
		f.MarkForLayoutAndRedraw()
	}
	e.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		e.Tooltip = nil
		if f.Shortcut != "" && ui.form == FormWindows {
			say := f.Shortcut
			if f.Label != "" {
				say = f.Label + " (" + f.Shortcut + ")"
			}
			e.Tooltip = newTooltip(ui, say, "")
		}
		return e.RectToRoot(e.ContentRect(true))
	}
	e.Accessibility.Callback = func(n *accessibility.Node) {
		n.Name = f.Label
		if f.spoken != nil {
			n.Name = f.spoken()
		}
		n.Description = f.Error
		if f.Error == "" && f.describe != nil {
			n.Description = f.describe()
		}
		if f.Required {
			n.Description = joinWords("Required", n.Description)
		}
		n.Placeholder = f.Placeholder
		n.Invalid = f.Error != ""
		if f.ReadOnly {
			n.ReadOnly = true
			n.Actions = n.Actions.Without(accessibility.SetValue, accessibility.ReplaceText)
		}
	}
	// unison's field has no read-only mode, so a read-only field drops what
	// would change its text: typed characters, the keys that delete or break
	// lines, and cutting and pasting. Moving, selecting and copying still work.
	typed := e.RuneTypedCallback
	e.RuneTypedCallback = func(ch rune) bool { return f.ReadOnly || typed(ch) }
	keyDown := e.KeyDownCallback
	e.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if f.ReadOnly {
			switch key {
			case unison.KeyBackspace, unison.KeyDelete, unison.KeyReturn, unison.KeyNumPadEnter:
				return true
			}
		}
		if (key == unison.KeyBackspace || key == unison.KeyDelete) && deletesWord(mods) {
			deleteWord(e, key == unison.KeyDelete)
			return true
		}
		if (key == unison.KeyLeft || key == unison.KeyRight) && ctrlMovesByWord(mods) {
			mods = mods&^mod.Control | mod.Option
		}
		return keyDown(key, mods, repeat)
	}
	for _, id := range []int{unison.CutItemID, unison.PasteItemID, unison.DeleteItemID} {
		can, do := e.InstallCmdHandlers(id, nil, nil)
		e.InstallCmdHandlers(id, func(v any) bool { return !f.ReadOnly && can(v) }, do)
	}
	// A secret is neither copied nor cut out of a masked field.
	for _, id := range []int{unison.CutItemID, unison.CopyItemID} {
		can, do := e.InstallCmdHandlers(id, nil, nil)
		e.InstallCmdHandlers(id, func(v any) bool { return !f.Masked && can(v) }, do)
	}
	// The editing commands as a Kvit menu, as every other menu is, rather
	// than unison's own, made of the standard commands so each has its
	// standard name and key. The menu leaves out what cannot act (CMD-10).
	e.ContextMenuCallback = nil
	ui.SetContextMenu(e, func(geom.Point) (string, []MenuItem) {
		line := func(std StdCommand, can func() bool, do func()) MenuItem {
			return MenuItem{Command: &Command{Std: std, Enabled: can, Run: do}}
		}
		never := func() bool { return false }
		cut, paste := e.CanCut, e.CanPaste
		if f.ReadOnly {
			cut, paste = never, never
		}
		undo, redo := never, never
		var doUndo, doRedo func()
		if m := unison.UndoManagerFor(e); m != nil && !f.ReadOnly {
			undo, redo, doUndo, doRedo = m.CanUndo, m.CanRedo, m.Undo, m.Redo
		}
		return "", []MenuItem{
			line(StdUndo, undo, doUndo),
			line(StdRedo, redo, doRedo),
			{Separator: true},
			line(StdCut, cut, e.Cut),
			line(StdCopy, e.CanCopy, e.Copy),
			line(StdPaste, paste, e.Paste),
			{Separator: true},
			line(StdSelectAll, e.CanSelectAll, e.SelectAll),
		}
	})
	f.message = NewLabel(ui, "")
	f.message.Role, f.message.Ink, f.message.Wrap = RoleCaption, InkDanger, true
	f.AddChild(e)
	f.AddChild(f.message)
	// The field's own panel says nothing; the text field inside it is what a
	// screen reader meets, under the field's label.
	f.Accessibility.Role = role.None
	f.SetLayout(fieldLayout{f})
}

// checkLength reports text past MaxLength as the field's error, and takes
// the report away once the text fits.
func (f *Field) checkLength() {
	n := len([]rune(f.edit.Text()))
	over := ""
	if f.MaxLength > 0 && n > f.MaxLength {
		over = f.ui.CountPhrase(n-f.MaxLength, "character", "") + " over the limit of " + f.ui.Number(f.MaxLength)
	}
	if f.Error == f.overLimit {
		f.Error = over
	}
	f.overLimit = over
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
	family := ui.Interface.FontFamily()
	if f.area != nil && f.area.Mono {
		family = ui.Interface.MonoFamily()
	}
	key := fieldFontKey{family, ui.Size(RoleBody)}
	if key != f.font {
		f.font = key
		f.edit.Font = unisonFont(ui, key.family, key.px)
	}
	f.edit.ObscurementRune = 0
	if f.Masked {
		f.edit.ObscurementRune = '•'
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
	radius := f.radius()
	// A plain text area is a document filling a pane, with no ground and no
	// outline unless it is in error.
	plain := f.area != nil && f.area.Plain
	if !plain {
		ground := t.PopupBackground
		if !f.edit.Enabled() {
			ground = t.ChipBackground
		}
		p.round(box, radius, ground)
	}
	edge := t.BorderStrong
	switch {
	case f.Error != "":
		edge = t.Danger
	case f.edit.Focused():
		edge = t.FocusRing
	}
	if !plain || f.Error != "" {
		p.outline(box, radius, float32(m.Hairline()), edge)
	}
	if f.area != nil && f.area.Underlay != nil {
		e := f.edit
		lineHeight := e.Font.LineHeight()
		f.area.Underlay(gc, func(index int) geom.Rect {
			pt := e.FromSelectionIndex(index)
			return geom.NewRect(pt.X, pt.Y, 0, lineHeight)
		})
	}
	if !f.edit.Enabled() {
		// unison greys a disabled field's text a second time, which leaves it
		// fainter than the disabled colour; it is drawn here instead.
		f.edit.OnBackgroundInk = unison.Transparent
	}
	gc.Save()
	f.edit.DefaultDraw(gc, dirty)
	gc.Restore()
	// Text drawn by the text layer: a disabled field's words, and the
	// placeholder. One line is centred on the field; several start at the
	// top and wrap.
	own := func(words string, ink palette.Color) {
		inner := f.edit.ContentRect(false)
		st := ui.Chrome(ui.Size(RoleBody), text.Regular, ink)
		if f.area != nil && f.area.Mono {
			st = ui.Mono(ui.Size(RoleBody), ink)
		}
		if f.area != nil {
			gc.Save()
			gc.ClipRect(inner, pathop.Intersect, false)
			ui.Fonts.Layout([]text.Span{{Text: words, Style: st}}, text.Options{MaxWidth: inner.Width}).Draw(gc, inner.X, inner.Y)
			gc.Restore()
			return
		}
		l := ui.Fonts.Layout([]text.Span{{Text: words, Style: st}}, text.Options{MaxWidth: inner.Width, Elide: true})
		_, h := l.Size()
		l.Draw(gc, inner.X, box.Y+(box.Height-h)/2)
	}
	if !f.edit.Enabled() && f.edit.Text() != "" {
		shown := f.edit.Text()
		if f.Masked {
			shown = strings.Repeat("•", len([]rune(shown)))
		}
		own(shown, t.TextDisabled)
	}
	if f.edit.Text() == "" && f.Placeholder != "" {
		own(f.Placeholder, t.TextFaint)
	}
	if f.decorate != nil {
		f.decorate(gc, box)
	}
}

// focusRing puts the ring around the field whenever it holds the focus,
// however the focus arrived: the caret is there, and so is what is typed.
func (f *Field) focusRing() (*unison.Panel, float32, bool) {
	return f.edit.AsPanel(), f.radius(), f.edit.Focused()
}

// radius is the corner radius of the field's box: none for a plain text area.
func (f *Field) radius() float32 {
	if f.area != nil && f.area.Plain {
		return 0
	}
	return float32(f.ui.Interface.RadiusControl())
}

// fieldPadding centres the line of text on the field's height, a near space
// in from the outline at each end.
type fieldPadding struct{ f *Field }

func (b fieldPadding) Insets() geom.Insets {
	f := b.f
	f.style()
	if f.area != nil {
		return geom.NewUniformInsets(float32(f.ui.Interface.SpaceNear()))
	}
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
	if l.f.area != nil {
		// Three rows to start; a text area takes more height when given it.
		msg := l.message(w)
		return geom.NewSize(float32(m.Px(60)), float32(m.ControlHeight())+msg),
			geom.NewSize(float32(m.Px(200)), float32(3*m.RowHeight())+msg),
			geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
	}
	h := float32(m.ControlHeight()) + l.message(w)
	return geom.NewSize(float32(m.Px(60)), h), geom.NewSize(float32(m.Px(200)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l fieldLayout) PerformLayout(target *unison.Panel) {
	f := l.f
	r := target.ContentRect(false)
	ch := float32(f.ui.Interface.ControlHeight())
	if f.area != nil {
		ch = max(0, r.Height-l.message(r.Width))
	}
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
