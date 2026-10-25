package kvitui

// History is the reader's own way back through the places they went in a
// window (LAY-10). The app says when the reader goes somewhere that counts as
// a step, and Back takes them to the place before. Choosing the place
// already shown, switching between peers in one navigation element, moving
// through a list and its details, and opening a dialog or a popup are not
// steps. Back is disabled when there is nowhere to go back to, and is never
// the only way to a place.
type History struct {
	ui      *UI
	current any
	before  []any
	// OnGo shows a place the reader went back to.
	OnGo func(place any)
	back *Command
}

// NewHistory returns a history that starts at a place.
func NewHistory(ui *UI, start any, onGo func(place any)) *History {
	h := &History{ui: ui, current: start, OnGo: onGo}
	h.back = &Command{Std: StdBack, Enabled: h.CanGoBack, Run: h.Back}
	return h
}

// Visit records that the reader went to a place: the place they were at
// becomes the one Back goes to. Going to the place already shown is no step.
func (h *History) Visit(place any) {
	if place == h.current {
		return
	}
	h.before = append(h.before, h.current)
	h.current = place
	h.ui.Refresh()
}

// Current is the place the reader is at.
func (h *History) Current() any { return h.current }

// CanGoBack reports whether there is a place to go back to.
func (h *History) CanGoBack() bool { return len(h.before) > 0 }

// Back takes the reader to the place before, through OnGo.
func (h *History) Back() {
	if len(h.before) == 0 {
		return
	}
	h.current = h.before[len(h.before)-1]
	h.before = h.before[:len(h.before)-1]
	if h.OnGo != nil {
		h.OnGo(h.current)
	}
	h.ui.Refresh()
}

// Command is the standard Back command for the history: Alt+Left on the
// Windows form and Cmd-[ on the macOS form, disabled when there is nowhere
// to go back to. Put it on the window's menu bar or toolbar, and its button
// at the top leading corner.
func (h *History) Command() *Command { return h.back }

// NewBackButton returns the Back button for a history, a symbol button for
// the top leading corner of the window's header.
func NewBackButton(ui *UI, h *History) *IconButton {
	return NewCommandIconButton(ui, h.back)
}
