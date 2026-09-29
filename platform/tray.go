package platform

// The icon in the desktop's notification area, its menu, and the desktop's
// own notifications: the apps' SystemTray (kvit-notes
// src/platform/systemtray.h), which unison has no equivalent of. Each system
// has its own backend: Shell_NotifyIcon on Windows (tray_windows.go),
// StatusNotifierItem and dbusmenu over D-Bus on Linux (tray_linux.go), and
// NSStatusItem with UNUserNotificationCenter on macOS (tray_darwin.go).
// Elsewhere, and on a Linux desktop without a StatusNotifierWatcher (WSLg has
// none), the tray reports that it is unavailable and shows nothing, while
// the app can still ask for everything;  behaves the same way.
//
// Everything here is called on unison's UI thread, and every callback is
// made there: a backend reports from its own thread or from the D-Bus
// connection's, and the Tray hands the report to unison.InvokeTask.

import (
	"image"
	"sync"

	"github.com/richardwilkes/unison"
)

// Authorization is whether the desktop lets the app post notifications,
// in the four states the app knows.
type Authorization int

const (
	// Unsupported means notifications cannot be posted here at all.
	Unsupported Authorization = iota
	// Unknown means the desktop has not been asked yet; only macOS asks.
	Unknown
	// Denied means the person said no.
	Denied
	// Authorized means notifications are shown.
	Authorized
)

func (a Authorization) String() string {
	switch a {
	case Unknown:
		return "Unknown"
	case Denied:
		return "Denied"
	case Authorized:
		return "Authorized"
	}
	return "Unsupported"
}

// TrayItem is one line of the tray icon's menu.
type TrayItem struct {
	// Text is what the line says, as plain text: an ampersand or an
	// underscore is shown as itself on every system.
	Text string
	// Separator makes the line a rule, with no text.
	Separator bool
	// Disabled shows the line greyed out, and choosing it does nothing.
	Disabled bool
	// OnSelect is called when the line is chosen.
	OnSelect func()
}

// Notification is a message posted to the desktop.
type Notification struct {
	Title   string
	Message string
	// ID comes back to Tray.OnNotificationClick when the notification is
	// clicked; "" for one that needs no answer.
	ID string
}

// Tray is the app's icon in the notification area, its menu, and the
// notifications it posts. Make one with NewTray at start; an app has one.
type Tray struct {
	// OnClick is called when the icon is clicked. macOS opens the menu on a
	// click, as menu bar items do there, so it is never called on macOS.
	OnClick func()
	// OnNotificationClick is called with the ID of a notification the
	// person clicked.
	OnNotificationClick func(id string)
	// OnAuthorizationChanged is called after Authorization changes.
	OnAuthorizationChanged func()

	backend trayBackend
	post    func(func())

	mu      sync.Mutex
	visible bool
	closed  bool
	icons   []image.Image
	tooltip string
	menu    []TrayItem
	auth    Authorization
	asking  bool // an authorization request is waiting for its answer
	last    Notification
}

// trayBackend is one system's notification area. Its methods are called on
// the UI thread; it reports what the person does through trayEvents, from
// whatever thread it hears it on.
type trayBackend interface {
	// available reports whether the desktop has a notification area.
	available() bool
	// show puts the icon in the notification area, or brings it up to date.
	show(s trayState)
	// hide takes the icon away.
	hide()
	// notify posts a notification; it is only called while the icon shows
	// and notifications are authorized.
	notify(n Notification)
	// authorization is the desktop's answer when the tray is made.
	authorization() Authorization
	// requestAuthorization asks the desktop, which answers through
	// trayEvents.authorized.
	requestAuthorization()
	// close takes the icon away for good and lets go of everything held.
	close()
}

// trayState is what the icon shows, as a backend is given it.
type trayState struct {
	icons   []image.Image
	tooltip string
	menu    []trayEntry
}

// trayEntry is a menu line as a backend draws it; the callback stays with
// the Tray, which a backend tells the line's index.
type trayEntry struct {
	text      string
	separator bool
	disabled  bool
}

// trayEvents is what a backend reports.
type trayEvents interface {
	chose(index int)
	clicked()
	notificationClicked(id string)
	authorized(a Authorization)
}

// NewTray returns the desktop's notification area. Where there is none,
// Available is false and the rest does nothing visible.
func NewTray() *Tray { return newSystemTray(unison.InvokeTask) }

// newSystemTray returns the desktop's notification area, whose callbacks
// are handed to post.
func newSystemTray(post func(func())) *Tray {
	t := &Tray{post: post}
	t.backend = newTrayBackend(t)
	t.auth = t.backend.authorization()
	return t
}

// NewTestTray returns a tray that says the desktop has a notification area,
// shows nothing, and answers Authorization with auth, for testing an app's
// use of the tray headlessly: Choose, Click and ClickNotification stand in
// for the person. Callbacks are made through unison.InvokeTask, as on a
// desktop.
func NewTestTray(auth Authorization) *Tray {
	return newTray(&testBackend{auth: auth}, unison.InvokeTask)
}

func newTray(b trayBackend, post func(func())) *Tray {
	t := &Tray{backend: b, post: post}
	t.auth = b.authorization()
	return t
}

// Available reports whether the desktop has a notification area. A tray
// that is not available records what it is asked to show and shows none of
// it.
func (t *Tray) Available() bool { return t.backend.available() }

// Show puts the icon in the notification area.
func (t *Tray) Show() {
	t.mu.Lock()
	t.visible = true
	t.mu.Unlock()
	t.update()
}

// Hide takes the icon out of the notification area.
func (t *Tray) Hide() {
	t.mu.Lock()
	was := t.visible && !t.closed
	t.visible = false
	t.mu.Unlock()
	if was && t.Available() {
		t.backend.hide()
	}
}

// Visible reports whether Show was called more recently than Hide.
func (t *Tray) Visible() bool { t.mu.Lock(); defer t.mu.Unlock(); return t.visible }

// SetIcon sets the icon's picture. Give it at several sizes where there are
// several: each system takes the smallest one at least as large as it
// draws the icon, and scales that.
func (t *Tray) SetIcon(sizes ...image.Image) {
	t.mu.Lock()
	t.icons = append([]image.Image(nil), sizes...)
	t.mu.Unlock()
	t.update()
}

// SetTooltip sets what the icon says when the pointer rests on it.
func (t *Tray) SetTooltip(text string) {
	t.mu.Lock()
	t.tooltip = text
	t.mu.Unlock()
	t.update()
}

// Tooltip is what the icon says when the pointer rests on it.
func (t *Tray) Tooltip() string { t.mu.Lock(); defer t.mu.Unlock(); return t.tooltip }

// SetMenu sets the icon's menu, which Windows and Linux open on a right
// click and macOS on any click.
func (t *Tray) SetMenu(items []TrayItem) {
	t.mu.Lock()
	t.menu = append([]TrayItem(nil), items...)
	t.mu.Unlock()
	t.update()
}

// Menu is the icon's menu.
func (t *Tray) Menu() []TrayItem {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]TrayItem(nil), t.menu...)
}

// update gives the backend what the icon shows, when it is shown.
func (t *Tray) update() {
	t.mu.Lock()
	if !t.visible || t.closed {
		t.mu.Unlock()
		return
	}
	s := trayState{icons: t.icons, tooltip: t.tooltip}
	for _, it := range t.menu {
		s.menu = append(s.menu, trayEntry{text: it.Text, separator: it.Separator, disabled: it.Disabled})
	}
	t.mu.Unlock()
	if t.Available() {
		t.backend.show(s)
	}
}

// Authorization is whether notifications are shown.
func (t *Tray) Authorization() Authorization { t.mu.Lock(); defer t.mu.Unlock(); return t.auth }

// RequestAuthorization asks the desktop to let the app post notifications,
// where it asks (macOS); the answer arrives through OnAuthorizationChanged.
// The desktop is asked once: after an answer, or while a question is
// waiting, this does nothing.
func (t *Tray) RequestAuthorization() {
	t.mu.Lock()
	if t.auth != Unknown || t.asking {
		t.mu.Unlock()
		return
	}
	t.asking = true
	t.mu.Unlock()
	t.backend.requestAuthorization()
}

// Notify posts a notification, and reports whether it was handed to the
// desktop: it is when notifications are authorized and the icon shows, as
// in the app. Either way it is kept as LastNotification.
func (t *Tray) Notify(title, message, id string) bool {
	n := Notification{Title: title, Message: message, ID: id}
	t.mu.Lock()
	t.last = n
	post := t.auth == Authorized && t.visible && !t.closed
	t.mu.Unlock()
	if !post || !t.Available() {
		return false
	}
	t.backend.notify(n)
	return true
}

// LastNotification is the notification most recently asked for, whether or
// not the desktop showed it.
func (t *Tray) LastNotification() Notification { t.mu.Lock(); defer t.mu.Unlock(); return t.last }

// Close takes the icon away and lets go of what the backend holds. Call it
// before the app quits: Windows otherwise leaves the icon in the
// notification area until the pointer passes over it.
func (t *Tray) Close() {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return
	}
	t.closed, t.visible = true, false
	t.mu.Unlock()
	t.backend.close()
}

// Choose runs the menu line with this text as though the person had chosen
// it, and reports whether there was one to choose. It is for tests and
// scripted checks.
func (t *Tray) Choose(text string) bool {
	t.mu.Lock()
	i := -1
	for k, it := range t.menu {
		if !it.Separator && it.Text == text {
			i = k
			break
		}
	}
	t.mu.Unlock()
	if i < 0 {
		return false
	}
	t.chose(i)
	return true
}

// Click acts as though the person had clicked the icon; for tests and
// scripted checks.
func (t *Tray) Click() { t.clicked() }

// ClickNotification acts as though the person had clicked the notification
// with this ID; for tests and scripted checks.
func (t *Tray) ClickNotification(id string) { t.notificationClicked(id) }

// chose runs the menu line at index on the UI thread, if it is still there
// and enabled.
func (t *Tray) chose(index int) {
	t.post(func() {
		t.mu.Lock()
		var fn func()
		if index >= 0 && index < len(t.menu) && !t.menu[index].Separator && !t.menu[index].Disabled {
			fn = t.menu[index].OnSelect
		}
		t.mu.Unlock()
		if fn != nil {
			fn()
		}
	})
}

func (t *Tray) clicked() {
	t.post(func() {
		if t.OnClick != nil {
			t.OnClick()
		}
	})
}

func (t *Tray) notificationClicked(id string) {
	t.post(func() {
		if t.OnNotificationClick != nil {
			t.OnNotificationClick(id)
		}
	})
}

// authorized records the desktop's answer to a request, or a change it
// announced, and tells the app when it differs.
func (t *Tray) authorized(a Authorization) {
	t.post(func() {
		t.mu.Lock()
		t.asking = false
		changed := a != t.auth
		t.auth = a
		t.mu.Unlock()
		if changed && t.OnAuthorizationChanged != nil {
			t.OnAuthorizationChanged()
		}
	})
}

// testBackend is NewTestTray's: a notification area that shows nothing and
// keeps what it was asked to post.
type testBackend struct {
	mu     sync.Mutex
	auth   Authorization
	shown  bool
	state  trayState
	posted []Notification
	asked  int
	closed bool
	// answer is what requestAuthorization answers, when answers is set;
	// otherwise the question waits for ever.
	answer  Authorization
	answers bool
	events  trayEvents
}

func (b *testBackend) available() bool { return true }

func (b *testBackend) show(s trayState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.shown, b.state = true, s
}

func (b *testBackend) hide() { b.mu.Lock(); defer b.mu.Unlock(); b.shown = false }

func (b *testBackend) notify(n Notification) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.posted = append(b.posted, n)
}

func (b *testBackend) authorization() Authorization { return b.auth }

func (b *testBackend) requestAuthorization() {
	b.mu.Lock()
	b.asked++
	answer, answers, events := b.answer, b.answers, b.events
	b.mu.Unlock()
	if answers && events != nil {
		events.authorized(answer)
	}
}

func (b *testBackend) close() { b.mu.Lock(); defer b.mu.Unlock(); b.closed, b.shown = true, false }
