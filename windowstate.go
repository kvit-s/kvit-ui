package kvitui

import (
	"github.com/kvit-s/kvit-ui/platform"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// The window's title, its unsaved-changes mark, what it asks before closing,
// where it was left, and full screen.

// SetDocument titles the window for the document it shows (WIN-12): on the
// Windows form, which Linux takes, "Report - Kvit Cash", the document and
// then the app's name, with an asterisk before the document while it has
// unsaved changes (WIN-20, decided Q-13); on the macOS form the document
// alone. "" titles it with the app's name.
func (w *Window) SetDocument(name string) {
	w.document = name
	w.retitle()
}

// SetEdited marks the document as having unsaved changes, or clears the
// mark once it is saved (WIN-20).
func (w *Window) SetEdited(edited bool) {
	w.edited = edited
	w.retitle()
}

// Edited reports whether the document is marked as having unsaved changes.
func (w *Window) Edited() bool { return w.edited }

// DocumentTitle is the title the window's document gives it in the app's
// form, as SetDocument says.
func (w *Window) DocumentTitle() string {
	app := w.ui.AppName
	doc := w.document
	if w.ui.form == FormMac {
		// The Mac marks unsaved changes with the system's dot in the close
		// button, which retitle asks for; the title stays the document.
		if doc == "" {
			return app
		}
		return doc
	}
	if doc == "" {
		return app
	}
	if w.edited {
		doc = "*" + doc
	}
	if app == "" {
		return doc
	}
	return doc + " - " + app
}

// retitle sets the title the document gives the window, once one was given,
// and on the macOS form the dot in the close button while the document has
// unsaved changes, which a Mac draws and other systems do not have.
func (w *Window) retitle() {
	if w.document != "" || w.ui.AppName != "" {
		title := w.DocumentTitle()
		w.SetTitle(title)
		platform.MarkEdited(title, w.edited && w.ui.form == FormMac)
	}
}

// RememberFrame keeps the window's place and size in the settings file under
// key whenever the reader moves or resizes it, and puts it back where it was
// left when called again at the next start, fitted to the screen's usable
// area so a window from a larger screen still shows whole (WIN-14). With no
// settings file it only fits the window to the screen.
func (w *Window) RememberFrame(key string) {
	store := w.ui.Settings
	if store != nil {
		if v, ok := store.Value(key); ok {
			if f, ok := v.([]any); ok && len(f) == 4 {
				var n [4]float32
				good := true
				for i, x := range f {
					x, ok := x.(float64)
					good = good && ok
					n[i] = float32(x)
				}
				if good && n[2] > 0 && n[3] > 0 {
					w.SetFrameRect(geom.NewRect(n[0], n[1], n[2], n[3]))
				}
			}
		}
	}
	w.FitToScreen()
	if store == nil {
		return
	}
	save := func() {
		if w.fullScreen {
			return
		}
		r := w.FrameRect()
		store.SetValue(key, []float64{float64(r.X), float64(r.Y), float64(r.Width), float64(r.Height)})
	}
	moved, resized := w.MovedCallback, w.ResizedCallback
	w.MovedCallback = func() {
		save()
		if moved != nil {
			moved()
		}
	}
	w.ResizedCallback = func() {
		save()
		if resized != nil {
			resized()
		}
	}
}

// FitToScreen moves and shrinks the window as little as it takes to show
// whole on the usable area of its screen, the area the task bar and the menu
// bar leave.
func (w *Window) FitToScreen() {
	d := w.Display()
	if d == nil {
		d = unison.PrimaryDisplay()
	}
	if d == nil {
		return
	}
	if r := w.FrameRect(); r != d.FitRectOnto(r) {
		w.SetFrameRect(d.FitRectOnto(r))
	}
}

// SetFullScreen fills the screen with the window, or puts it back where it
// was (WIN-15). It starts and ends only when the reader asks, by the
// standard Full screen command (F11 on the Windows form, Ctrl-Cmd-F on the
// macOS form). unison has no full-screen mode of its own, so the window
// covers the whole screen with its frame still drawn; on a Mac the window's
// green button also takes it into the system's full screen.
func (w *Window) SetFullScreen(on bool) {
	if on == w.fullScreen {
		return
	}
	d := w.Display()
	if d == nil {
		d = unison.PrimaryDisplay()
	}
	if on {
		w.beforeFullScreen = w.FrameRect()
		w.fullScreen = true
		if d != nil {
			w.SetFrameRect(d.Frame)
		}
	} else {
		w.fullScreen = false
		w.SetFrameRect(w.beforeFullScreen)
	}
	w.ui.Refresh()
}

// FullScreen reports whether the window fills the screen.
func (w *Window) FullScreen() bool { return w.fullScreen }

// FullScreenCommand is the standard Full screen command for the window: a
// switch, shown ticked on its View menu line while on, on F11 on the Windows
// form and Ctrl-Cmd-F on the macOS form.
func (w *Window) FullScreenCommand() *Command {
	return &Command{Std: StdFullScreen, Checked: w.FullScreen,
		TurnOn: "Enter full screen", TurnOff: "Exit full screen",
		Run: func() { w.SetFullScreen(!w.fullScreen) }}
}
