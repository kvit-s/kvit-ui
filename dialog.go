package kvitui

import (
	"log/slog"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// DialogKind is what a dialog is for, which decides its buttons.
type DialogKind int

const (
	// DialogAsk asks a question: a button naming the action (ConfirmText,
	// required, never "Confirm" or "OK"), and the safe button, Cancel.
	DialogAsk DialogKind = iota
	// DialogInform only tells the reader something: one button, Close on
	// the Windows form and OK on the macOS form (WRD-7, decided Q-5).
	DialogInform
	// DialogProperties shows the settings of one object: OK, Cancel and,
	// on the Windows form, Apply, the one place those words belong (WRD-7,
	// WIN-23). Apply is enabled once MarkChanged says something changed;
	// the macOS form leaves it out.
	DialogProperties
)

// Dialog is a surface that has to be answered before anything else happens:
// a destructive action, a choice the next step depends on. Modality is
// expensive, and a confirmation of something reversible belongs in a toast
// with an undo instead.
//
// Its buttons follow the platform (WIN-1): on the Windows form, which Linux
// takes, the action first and the safe button at the right end; on the
// macOS form Cancel and then the action at the trailing end. Every dialog has
// a safe button (WIN-3). Return presses the default button on the macOS
// form, which is the action unless the dialog is destructive; on the Windows
// form Enter presses the focused button, and the action only where
// DefaultConfirm marks it as the likely and safe answer (KEY-6). Escape, and
// Command-period on the macOS form, cancel; on the Windows form so does the
// close symbol beside the title, which the macOS form does not draw (WIN-6).
// The keyboard starts in the first field of the content, or on the default
// button, stays inside while the dialog is open, and goes back to where it
// was when it closes (KEY-8).
//
// A Dialog is a panel; Open shows it modal in a Kvit window, and it can also
// be placed on a page as it is. An open dialog moves when the reader drags
// its title row, as a window does, so it can be pulled off what it covers.
// It opens in the middle of the window unless MoveTo says where, and OnMove
// tells the application where the reader left it, to keep and give back to
// MoveTo the next time.
type Dialog struct {
	unison.Panel
	ui *UI
	// Kind is what the dialog is for: asking (the default), informing, or
	// showing an object's properties.
	Kind DialogKind
	// Title names the question or what happened, on every dialog (WIN-6,
	// decided Q-42); Detail says more.
	Title, Detail string
	// ConfirmText names the action the dialog confirms, "Import" or
	// "Delete": a dialog that asks needs it, and has no "Confirm" to fall
	// back on (WRD-7). A properties dialog's is OK unless set.
	ConfirmText string
	// CancelText is the safe button's words, Cancel unless set; a dialog
	// that asks always has one (WIN-3).
	CancelText string
	// OtherText is an optional second answer, such as "Don't save", beside
	// the action on the Windows form and at the leading end on the macOS
	// form; "" for none.
	OtherText string
	// Destructive draws the confirming button in the danger colour, and
	// takes away the default on both forms.
	Destructive bool
	// DefaultConfirm makes the action the default on the Windows form, which
	// Enter presses wherever the keyboard is in the dialog: only where it is
	// clearly the likely and safe answer. The macOS form's default is the
	// action unless the dialog is destructive.
	DefaultConfirm bool
	// Suppressible offers "Don't show this message again" under the detail,
	// for a routine message that can recur; an error or a destructive
	// question never offers it (WIN-2). SuppressNoun names what is not shown
	// again, "message" unless set.
	Suppressible bool
	SuppressNoun string
	// Helper lets the dialog open over another on the Windows form, as a
	// chooser over the dialog that needs it; an informing dialog may open
	// over another on both forms. Otherwise one dialog is open at a time
	// (WIN-7).
	Helper bool
	// Alert lets a question about something outside the dialog under it,
	// such as quitting while that dialog is open, open over it on both
	// forms, as an alert may sit over a dialog (WIN-7).
	Alert bool
	// OnAccept and OnReject run when the dialog is answered; OnOther when
	// the second answer is chosen; OnApply when a properties dialog's Apply
	// is pressed, which leaves it open.
	OnAccept, OnReject, OnOther, OnApply func()
	// OnDiscard is asked before the close symbol or Escape throws away what
	// the reader put into the dialog: true goes ahead, false keeps it open.
	// The app asks the reader there when the input took effort (WIN-8).
	// Cancel discards without asking.
	OnDiscard func() bool
	// OnMove runs when the reader lets go of the dialog after dragging it,
	// with where its top-left corner is now, as Position gives it.
	OnMove func(x, y int)
	// Help opens the help topic for the dialog, where the topic says what
	// the dialog cannot; nil for none (HLP-6). On the Windows form it is the
	// question-worded link at the bottom of the content, which Tab reaches
	// after the buttons (KEY-9); on the macOS form it is the round help
	// button at the leading end of the button row.
	Help *HelpLink
	// ConfirmKeepsOpen leaves the dialog open when it is confirmed, for an
	// answer the application checks and may refuse, saying why in the
	// dialog; Close takes it away once the answer is taken.
	ConfirmKeepsOpen bool

	title    *Label
	close    *IconButton
	detail   *Label
	body     *unison.Panel
	suppress *Check
	cancel   *Button
	confirm  *Button
	other    *Button
	apply    *Button
	changed  bool // something changed since the dialog opened or Apply was pressed
	hide     func()
	opener   *unison.Panel // what held the keyboard before the dialog opened
	done     bool          // the work is done, and Cancel has become Close or Done
	// Where the top-left corner goes, from the window's top-left corner;
	// nil puts the dialog in the middle.
	x, y Measure
	// A drag of the title row: where the pointer was pressed and where the
	// dialog's top-left corner was then, both in the window's content
	// coordinates, and whether it has moved the dialog yet.
	dragging     bool
	grab, grabAt geom.Point
	moved        bool
}

// NewDialog returns a dialog asking a question, holding content under its
// detail. Give it ConfirmText, the words of the action it asks about.
func NewDialog(ui *UI, title string, content ...unison.Paneler) *Dialog {
	d := &Dialog{ui: ui, Title: title, CancelText: "Cancel", SuppressNoun: "message"}
	d.Self = d
	d.title = NewLabel(ui, "")
	d.title.Role, d.title.Weight = RoleTitle, text.Bold
	d.close = NewIconButton(ui, "close", "Close")
	d.close.OnClick = d.dismiss
	d.detail = NewLabel(ui, "")
	d.detail.Ink, d.detail.Wrap = InkTextSecondary, true
	d.body = unison.NewPanel()
	d.body.SetLayout(&unison.FlexLayout{Columns: 1})
	for _, c := range content {
		if c.AsPanel().LayoutData() == nil {
			c.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		}
		d.body.AddChild(c)
	}
	d.suppress = NewCheck(ui, "")
	d.cancel = NewButton(ui, "")
	d.cancel.OnClick = d.reject
	d.confirm = NewButton(ui, "")
	d.confirm.Form = ButtonPrimary
	d.confirm.OnClick = d.accept
	d.other = NewButton(ui, "")
	d.other.OnClick = d.answerOther
	d.apply = NewButton(ui, "")
	d.apply.OnClick = func() {
		if d.OnApply != nil {
			d.OnApply()
		}
		d.changed = false
		if d.apply.Focused() {
			// The keyboard does not stay on a button that can no longer act.
			d.confirm.RequestFocus()
		}
		d.sync()
		d.MarkForRedraw()
	}
	d.SetLayout(syncing{Layout: dialogLayout{d}, sync: d.sync})
	d.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) { drawSurface(gc, ui, d.ContentRect(true)) }
	d.MouseDownCallback = d.mouseDown
	d.MouseDragCallback = d.mouseDrag
	d.MouseUpCallback = d.mouseUp
	d.KeyDownCallback = d.keyDown
	d.sync()
	return d
}

// NewMessage returns a dialog that only informs: one button, Close on the
// Windows form and OK on the macOS form.
func NewMessage(ui *UI, title, detail string) *Dialog {
	d := NewDialog(ui, title)
	d.Kind, d.Detail = DialogInform, detail
	return d
}

// words are the buttons' words as the dialog's kind and form give them.
func (d *Dialog) words() (confirm, cancel, apply string) {
	mac := d.ui.form == FormMac
	switch d.Kind {
	case DialogInform:
		if mac {
			return "OK", "", ""
		}
		return "Close", "", ""
	case DialogProperties:
		confirm = d.ConfirmText
		if confirm == "" {
			confirm = "OK"
		}
		if !mac {
			apply = "Apply"
		}
		return confirm, d.CancelText, apply
	}
	cancel = d.CancelText
	if d.done {
		cancel = "Close"
		if mac {
			cancel = "Done"
		}
		return "", cancel, ""
	}
	return d.ConfirmText, cancel, ""
}

// sync copies the fields into the parts, keeping as children only the
// parts that have something to show, in the order the form reads them.
func (d *Dialog) sync() {
	d.title.Text, d.detail.Text = d.Title, d.Detail
	confirm, cancel, apply := d.words()
	d.confirm.Text, d.cancel.Text, d.apply.Text = confirm, cancel, apply
	// Apply shows what changed without closing, so there is nothing for it
	// to do until something has (WIN-23).
	d.apply.SetEnabled(d.changed)
	d.apply.Explanation = ""
	if !d.changed {
		d.apply.Explanation = "No changes to apply yet"
	}
	d.other.Text = d.OtherText
	if d.done || d.Kind != DialogAsk {
		d.other.Text = ""
	}
	d.confirm.Danger = d.Destructive && d.Kind == DialogAsk
	noun := d.SuppressNoun
	if noun == "" {
		noun = "message"
	}
	d.suppress.Text = "Don't show this " + noun + " again"
	if d.ui.form == FormMac {
		d.suppress.Text = "Do not show this " + noun + " again"
	}
	parts := []unison.Paneler{d.title}
	if d.ui.form == FormWindows {
		parts = append(parts, d.close)
	}
	if d.Detail != "" {
		parts = append(parts, d.detail)
	}
	parts = append(parts, d.body)
	if d.Suppressible && !d.Destructive && d.Kind != DialogProperties {
		parts = append(parts, d.suppress)
	}
	parts = append(parts, d.buttons()...)
	if d.Help != nil {
		// After the buttons, so Tab reaches it last.
		parts = append(parts, d.Help)
	}
	showOnly(d.AsPanel(), parts...)
}

// buttons are the buttons shown, left to right: on the Windows form the
// action, the second answer, the safe button and Apply; on the macOS form
// the second answer, then the safe button and the action at the trailing
// end (WIN-1).
func (d *Dialog) buttons() []unison.Paneler {
	var out []unison.Paneler
	add := func(b *Button) {
		if b.Text != "" {
			out = append(out, b)
		}
	}
	if d.ui.form == FormMac {
		add(d.other)
		add(d.cancel)
		add(d.apply)
		add(d.confirm)
		return out
	}
	add(d.confirm)
	add(d.other)
	add(d.cancel)
	add(d.apply)
	return out
}

// ButtonWords are the words of the dialog's buttons as they are shown, left
// to right.
func (d *Dialog) ButtonWords() []string {
	d.sync()
	var out []string
	for _, b := range d.buttons() {
		out = append(out, b.(*Button).Words())
	}
	return out
}

// Buttons are the dialog's buttons as they are shown, left to right: the
// same buttons ButtonWords names, as handles the application can enable,
// focus or press, rather than finding by their words.
func (d *Dialog) Buttons() []*Button {
	d.sync()
	var out []*Button
	for _, b := range d.buttons() {
		out = append(out, b.(*Button))
	}
	return out
}

// ConfirmButton is the dialog's action, "Import" or OK: nil when the dialog
// shows no action, as an asking dialog does once Done. An informing dialog's
// only button, Close on the Windows form and OK on the macOS form, is its
// confirming button.
func (d *Dialog) ConfirmButton() *Button {
	d.sync()
	if d.confirm.Text == "" {
		return nil
	}
	return d.confirm
}

// CancelButton is the dialog's safe button, Cancel unless set, Close or Done
// once Done: nil when the dialog shows none, as an informing dialog does.
func (d *Dialog) CancelButton() *Button {
	d.sync()
	if d.cancel.Text == "" {
		return nil
	}
	return d.cancel
}

// OtherButton is the dialog's optional second answer, such as "Don't save":
// nil when the dialog shows none.
func (d *Dialog) OtherButton() *Button {
	d.sync()
	if d.other.Text == "" {
		return nil
	}
	return d.other
}

// ApplyButton is a properties dialog's Apply on the Windows form, which the
// reader presses once MarkChanged says something changed: nil when the dialog
// shows none, as every dialog does on the macOS form.
func (d *Dialog) ApplyButton() *Button {
	d.sync()
	if d.apply.Text == "" {
		return nil
	}
	return d.apply
}

// SafeButton reports whether the dialog has a button that leaves without
// acting (WIN-3).
func (d *Dialog) SafeButton() bool {
	d.sync()
	return d.cancel.Text != "" || d.Kind == DialogInform
}

// hasDefault reports whether the action is the default button, which Return
// presses: on the macOS form unless destructive, on the Windows form only
// where DefaultConfirm says so and it is not destructive (KEY-6).
func (d *Dialog) hasDefault() bool {
	if d.confirm.Text == "" || (d.Destructive && d.Kind == DialogAsk) {
		return false
	}
	return d.ui.form == FormMac || d.DefaultConfirm || d.Kind == DialogInform
}

// MarkChanged tells a properties dialog that the reader has changed
// something, which enables Apply on the Windows form until it is pressed
// (WIN-23). Call it from the fields' and boxes' change callbacks.
func (d *Dialog) MarkChanged() {
	if d.changed {
		return
	}
	d.changed = true
	d.sync()
	d.MarkForRedraw()
}

// Suppressed reports whether the reader ticked "Don't show this message
// again".
func (d *Dialog) Suppressed() bool { return d.Suppressible && d.suppress.Checked }

// Done says the dialog's work has been carried out and cannot be taken back:
// Cancel becomes Close on the Windows form and Done on the macOS form, and the
// action goes (WRD-8).
func (d *Dialog) Done() {
	d.done = true
	d.sync()
	d.MarkForLayoutAndRedraw()
}

// Open shows the dialog modal, in the middle of the window or where MoveTo
// put it, with the keyboard in the first field of its content, or on the
// default button, or on the safe one. A dialog is refused, with a warning,
// while another is open, unless it may sit over it (WIN-7); one that asks
// without ConfirmText is opened and warned about (WRD-7, WIN-3).
func (d *Dialog) Open(window *Window) {
	if d.hide != nil {
		return
	}
	d.sync()
	if d.Kind == DialogAsk && d.ConfirmText == "" && !d.done {
		slog.Warn("kvitui: a dialog that asks has no ConfirmText naming its action", "title", d.Title)
	}
	if d.Kind == DialogAsk && d.CancelText == "" {
		slog.Warn("kvitui: a dialog has no safe button", "title", d.Title)
	}
	if d.Title == "" {
		slog.Warn("kvitui: a dialog has no title")
	}
	over := d.Kind == DialogInform || d.Alert || (d.Helper && d.ui.form == FormWindows)
	for _, p := range window.popups {
		if p.Modal && !over {
			slog.Warn("kvitui: a dialog was not opened over another", "title", d.Title)
			return
		}
	}
	d.opener = window.CurrentFocus()
	d.hide = window.Show(&Popup{Panel: d, Place: d.place, Modal: true, OnEscape: d.dismiss})
	unison.InvokeTask(func() {
		if d.hide == nil {
			return
		}
		if f := firstField(d.body); f != nil {
			f.RequestFocus()
			return
		}
		switch {
		case d.hasDefault():
			d.confirm.Focus()
		case d.cancel.Text != "":
			d.cancel.Focus()
		case d.ui.form == FormWindows:
			d.close.Focus()
		}
	})
}

// firstField is the first control in a panel that takes typing, or nil.
func firstField(p *unison.Panel) *unison.Panel {
	if isTextField(p) && p.Enabled() && !p.Hidden {
		return p
	}
	for _, c := range p.Children() {
		if f := firstField(c); f != nil {
			return f
		}
	}
	return nil
}

// finish takes the dialog away and gives the keyboard back to where it was
// before it opened.
func (d *Dialog) finish() {
	if d.hide == nil {
		return
	}
	d.hide()
	d.hide = nil
	if d.opener != nil && d.opener.Window() != nil {
		d.opener.RequestFocus()
	}
	d.opener = nil
}

func (d *Dialog) accept() {
	if !d.ConfirmKeepsOpen || d.Kind == DialogInform {
		d.finish()
	}
	if d.OnAccept != nil {
		d.OnAccept()
	}
}

func (d *Dialog) answerOther() {
	d.finish()
	if d.OnOther != nil {
		d.OnOther()
	}
}

// A dismissal is a decline, not an agreement left unsaid.
func (d *Dialog) reject() {
	if d.Kind == DialogInform {
		d.accept()
		return
	}
	d.finish()
	if d.OnReject != nil {
		d.OnReject()
	}
}

// dismiss is the close symbol and Escape: a decline, once OnDiscard agrees
// to throw away what the reader put in.
func (d *Dialog) dismiss() {
	if d.OnDiscard != nil && d.Kind != DialogInform && !d.OnDiscard() {
		return
	}
	d.reject()
}

// keyDown presses the default button on Return, from anywhere in the dialog
// a key reaches it from, such as a field.
func (d *Dialog) keyDown(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
	if repeat || mods&(mod.Control|mod.Command|mod.Option) != 0 {
		return false
	}
	if (key == unison.KeyReturn || key == unison.KeyNumPadEnter) && d.hasDefault() && d.confirm.Enabled() {
		d.accept()
		return true
	}
	return false
}

type dialogLayout struct{ d *Dialog }

// arrange places the parts across a width and returns the dialog's height:
// the title a view margin below the top with the close control at the end of
// its row on the Windows form, the detail and the content a loose space
// apart, the don't-show-again box, and the buttons in a row at the bottom
// right.
func (l dialogLayout) arrange(r geom.Rect, set bool) float32 {
	d, m := l.d, l.d.ui.Interface
	margin, gap := float32(m.ViewMargin()), float32(m.Space())
	inner := max(0, r.Width-2*margin)
	put := func(p unison.Paneler, box geom.Rect) {
		if set {
			p.AsPanel().SetFrameRect(box)
		}
	}
	ch := float32(m.ControlHeight())
	closeWidth := float32(0)
	if d.ui.form == FormWindows {
		_, cp, _ := d.close.Sizes(geom.Size{})
		put(d.close, geom.NewRect(r.Right()-margin-cp.Width, r.Y+margin+(ch-cp.Height)/2, cp.Width, cp.Height))
		closeWidth = cp.Width + gap
	}
	_, tp, _ := d.title.Sizes(geom.Size{})
	put(d.title, geom.NewRect(r.X+margin, r.Y+margin+(ch-tp.Height)/2, max(0, inner-closeWidth), tp.Height))
	y := r.Y + margin + ch + margin
	// The detail and the content a loose space apart; the content keeps its
	// place, and the gap before it, even when it is empty.
	loose := float32(m.SpaceLoose())
	if d.Detail != "" {
		_, dp, _ := d.detail.Sizes(geom.NewSize(inner, 0))
		put(d.detail, geom.NewRect(r.X+margin, y, inner, dp.Height))
		y += dp.Height + loose
	}
	_, bp, _ := d.body.Sizes(geom.NewSize(inner, 0))
	put(d.body, geom.NewRect(r.X+margin, y, inner, bp.Height))
	y += bp.Height
	if d.suppress.Parent() != nil {
		_, sp, _ := d.suppress.Sizes(geom.Size{})
		y += loose
		put(d.suppress, geom.NewRect(r.X+margin, y, sp.Width, sp.Height))
		y += sp.Height
	}
	mac := d.ui.form == FormMac
	help := d.Help != nil && d.Help.Parent() == d.AsPanel()
	if help && !mac {
		// The link at the bottom of the content, at the leading edge.
		_, hp, _ := d.Help.Sizes(geom.NewSize(inner, 0))
		y += loose
		put(d.Help, geom.NewRect(r.X+margin, y, min(hp.Width, inner), hp.Height))
		y += hp.Height
	}
	buttons := d.buttons()
	if len(buttons) == 0 {
		return y + margin - r.Y
	}
	row := float32(m.RowHeightSlim()) + margin
	leading := r.X + margin
	if help && mac {
		// The round help button first at the leading end, as the Mac has it.
		_, hp, _ := d.Help.Sizes(geom.Size{})
		put(d.Help, geom.NewRect(leading, y+margin+(row-hp.Height)/2, hp.Width, hp.Height))
		leading += hp.Width + gap
	}
	x := r.Right() - margin
	for i := len(buttons) - 1; i >= 0; i-- {
		b := buttons[i]
		_, bp, _ := b.AsPanel().Sizes(geom.Size{})
		bx := x - bp.Width
		if mac && b == unison.Paneler(d.other) && len(buttons) > 1 {
			// The second answer stands apart at the leading end.
			bx = leading
		} else {
			x = bx - gap
		}
		// Centred in the whole footer, a slim row and a margin tall.
		put(b, geom.NewRect(bx, y+margin+(row-bp.Height)/2, bp.Width, bp.Height))
	}
	return y + margin + row - r.Y
}

func (l dialogLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	w := float32(l.d.ui.Interface.Px(420))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := l.arrange(geom.NewRect(0, 0, w, 0), false)
	return geom.NewSize(w, h), geom.NewSize(float32(l.d.ui.Interface.Px(420)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l dialogLayout) PerformLayout(target *unison.Panel) { l.arrange(target.ContentRect(true), true) }

// ProvideAccessibility describes the dialog by its question and detail.
func (d *Dialog) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Dialog
	}
	n.Name, n.Description = d.Title, d.Detail
	n.Modal = d.hide != nil
}
