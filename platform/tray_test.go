package platform

import (
	"image"
	"image/color"
	"slices"
	"testing"
)

// syncTray is a tray on the test backend whose callbacks run at once, on
// the test's goroutine.
func syncTray(auth Authorization) (*Tray, *testBackend) {
	b := &testBackend{auth: auth}
	t := newTray(b, func(f func()) { f() })
	b.events = t
	return t, b
}

func square(side int, c color.Color) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, side, side))
	for i := range side * side {
		img.Set(i%side, i/side, c)
	}
	return img
}

// The menu's lines run their callbacks by text; separators and disabled
// lines are never run.
func TestTheMenuRunsTheLineChosen(t *testing.T) {
	tray, _ := syncTray(Authorized)
	var ran []string
	tray.SetMenu([]TrayItem{
		{Text: "New Note", OnSelect: func() { ran = append(ran, "new") }},
		{Separator: true},
		{Text: "Greyed", Disabled: true, OnSelect: func() { ran = append(ran, "greyed") }},
		{Text: "Quit", OnSelect: func() { ran = append(ran, "quit") }},
	})
	for _, text := range []string{"New Note", "Quit", "Greyed", "Missing"} {
		tray.Choose(text)
	}
	if !slices.Equal(ran, []string{"new", "quit"}) {
		t.Errorf("ran %q", ran)
	}
	if tray.Choose("Missing") {
		t.Error("Choose should report a line that is not there")
	}
	// A backend reports by index; one past the end, or a separator, does nothing.
	tray.chose(1)
	tray.chose(9)
	if len(ran) != 2 {
		t.Errorf("a separator or a missing index ran something: %q", ran)
	}
}

// The backend is given the icon only while it is shown, and each change
// while it is.
func TestTheBackendSeesTheIconOnlyWhileShown(t *testing.T) {
	tray, b := syncTray(Authorized)
	tray.SetTooltip("Kvit Notes")
	tray.SetMenu([]TrayItem{{Text: "Show Kvit"}, {Separator: true}, {Text: "Quit", Disabled: true}})
	if b.shown {
		t.Fatal("the icon was shown before Show")
	}
	tray.Show()
	if !b.shown || b.state.tooltip != "Kvit Notes" || len(b.state.menu) != 3 {
		t.Fatalf("after Show: shown %v, state %+v", b.shown, b.state)
	}
	if !b.state.menu[1].separator || !b.state.menu[2].disabled || b.state.menu[0].text != "Show Kvit" {
		t.Errorf("the menu as the backend has it: %+v", b.state.menu)
	}
	tray.SetTooltip("Kvit Notes (2)")
	if b.state.tooltip != "Kvit Notes (2)" {
		t.Errorf("a change while shown did not reach the backend: %q", b.state.tooltip)
	}
	tray.Hide()
	if b.shown || tray.Visible() {
		t.Error("Hide left the icon shown")
	}
	tray.Close()
	tray.Show()
	if b.shown || !b.closed {
		t.Error("a closed tray was shown again")
	}
}

// A notification reaches the desktop only when authorized and while the
// icon shows, and is recorded either way.
func TestNotificationsNeedAuthorizationAndTheIcon(t *testing.T) {
	tray, b := syncTray(Authorized)
	if tray.Notify("Kvit", "Before the icon", "") {
		t.Error("posted before Show")
	}
	tray.Show()
	if !tray.Notify("Kvit", "Note captured", "note:one") {
		t.Error("not posted while shown and authorized")
	}
	if got := tray.LastNotification(); got != (Notification{Title: "Kvit", Message: "Note captured", ID: "note:one"}) {
		t.Errorf("last notification %+v", got)
	}
	if len(b.posted) != 1 || b.posted[0].ID != "note:one" {
		t.Errorf("posted %+v", b.posted)
	}

	denied, db := syncTray(Denied)
	denied.Show()
	if denied.Notify("Kvit", "Denied", "x") || len(db.posted) != 0 {
		t.Error("posted while denied")
	}
	if denied.LastNotification().Message != "Denied" {
		t.Error("a notification not posted should still be recorded")
	}
}

// A notification that no longer applies is withdrawn by its ID; one with
// no ID cannot be named, and a closed tray withdraws nothing (MSG-12).
func TestANotificationIsWithdrawnByItsID(t *testing.T) {
	tray, b := syncTray(Authorized)
	tray.Show()
	tray.Post(Notification{Title: "Sync stopped", Message: "The bank asks you to sign in again.", ID: "sync:bank", Kind: NotifyProblem})
	tray.Post(Notification{Title: "Import finished", ID: "import:1", Kind: NotifySuccess})
	tray.Withdraw("sync:bank")
	tray.Withdraw("")
	if !slices.Equal(b.withdrawn, []string{"sync:bank"}) {
		t.Errorf("the backend was asked to withdraw %q", b.withdrawn)
	}
	if len(b.posted) != 1 || b.posted[0].ID != "import:1" {
		t.Errorf("still posted: %+v", b.posted)
	}
	if tray.LastNotification().ID != "import:1" {
		t.Errorf("withdrawing another notification cleared the last one: %+v", tray.LastNotification())
	}
	tray.Withdraw("import:1")
	if tray.LastNotification() != (Notification{}) {
		t.Errorf("the last notification is %+v after it was withdrawn", tray.LastNotification())
	}
	tray.Close()
	tray.Withdraw("import:1")
	if len(b.withdrawn) != 2 {
		t.Errorf("a closed tray withdrew again: %q", b.withdrawn)
	}
}

// Clicks come back to the app's callbacks with the notification's ID.
func TestClicksReachTheCallbacks(t *testing.T) {
	tray, _ := syncTray(Authorized)
	clicks := 0
	var ids []string
	tray.OnClick = func() { clicks++ }
	tray.OnNotificationClick = func(id string) { ids = append(ids, id) }
	tray.Click()
	tray.ClickNotification("note:one")
	tray.ClickNotification("note:two")
	if clicks != 1 || !slices.Equal(ids, []string{"note:one", "note:two"}) {
		t.Errorf("clicks %d, notification ids %q", clicks, ids)
	}
}

// The desktop is asked once: a second request after an answer, or while
// one is waiting, does not ask again.
func TestAuthorizationIsAskedOnce(t *testing.T) {
	for _, answer := range []Authorization{Authorized, Denied} {
		tray, b := syncTray(Unknown)
		b.answer, b.answers = answer, true
		changes := 0
		tray.OnAuthorizationChanged = func() { changes++ }
		tray.RequestAuthorization()
		tray.RequestAuthorization()
		if tray.Authorization() != answer || changes != 1 || b.asked != 1 {
			t.Errorf("%v: authorization %v, %d changes, asked %d times", answer, tray.Authorization(), changes, b.asked)
		}
	}
	// No answer yet: the question is still waiting, so it is not put again.
	tray, b := syncTray(Unknown)
	tray.RequestAuthorization()
	tray.RequestAuthorization()
	if b.asked != 1 || tray.Authorization() != Unknown {
		t.Errorf("asked %d times while waiting", b.asked)
	}
	// Platforms without a prompt start with their answer and never ask.
	plain, pb := syncTray(Authorized)
	plain.RequestAuthorization()
	if pb.asked != 0 {
		t.Error("a platform that answers at once was asked")
	}
}

// The icon is drawn from the smallest picture at least as large as the
// size needed, or the largest there is.
func TestTheIconIsScaledFromTheNearestSize(t *testing.T) {
	red, green, blue := color.NRGBA{255, 0, 0, 255}, color.NRGBA{0, 255, 0, 255}, color.NRGBA{0, 0, 255, 255}
	sizes := []image.Image{square(64, blue), square(16, red), square(32, green)}
	for size, want := range map[int]color.NRGBA{16: red, 20: green, 32: green, 48: blue, 128: blue} {
		img := iconAt(sizes, size)
		if img.Bounds().Dx() != size || img.NRGBAAt(size/2, size/2) != want {
			t.Errorf("at %d: %v px, middle %v, want %v", size, img.Bounds().Dx(), img.NRGBAAt(size/2, size/2), want)
		}
	}
	if iconAt(nil, 16) != nil {
		t.Error("an icon with no picture")
	}
}

// The system's own tray always answers, whatever the desktop has.
func TestTheSystemTrayAnswers(t *testing.T) {
	tray := newSystemTray(func(f func()) { f() })
	defer tray.Close()
	t.Logf("tray available %v, notifications %v", tray.Available(), tray.Authorization())
	if !tray.Available() && tray.Authorization() != Unsupported {
		t.Error("notifications are authorized without a notification area")
	}
}

// A line that turns something on and off keeps its words and is ticked
// while on; the line a click on the icon runs is the default (WRD-12,
// CMD-24).
func TestTheMenuTicksAndMarksTheDefault(t *testing.T) {
	tray, b := syncTray(Authorized)
	tray.SetMenu([]TrayItem{{Text: "Open", Default: true}, {Text: "Pause syncing", Checked: true}})
	tray.Show()
	if len(b.state.menu) != 2 || !b.state.menu[0].isDefault || !b.state.menu[1].checked || b.state.menu[1].disabled {
		t.Errorf("the menu as the backend has it: %+v", b.state.menu)
	}
}
