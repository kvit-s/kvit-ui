package platform

// The Linux tray against a private session bus: a dbus-daemon of the test's
// own, with a StatusNotifierWatcher and a notification server played by the
// test, so nothing reaches the desktop the tests run on.

import (
	"bufio"
	"context"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

// privateBus starts a dbus-daemon for the test and points the session bus
// at it, or skips the test where dbus-daemon is not installed.
func privateBus(t *testing.T) {
	t.Helper()
	daemon, err := exec.LookPath("dbus-daemon")
	if err != nil {
		t.Skip("dbus-daemon is not installed")
	}
	dir := t.TempDir()
	config := filepath.Join(dir, "bus.conf")
	if err := os.WriteFile(config, []byte(`<!DOCTYPE busconfig PUBLIC "-//freedesktop//DTD D-BUS Bus Configuration 1.0//EN"
 "http://www.freedesktop.org/standards/dbus/1.0/busconfig.dtd">
<busconfig>
  <type>session</type>
  <listen>unix:dir=`+dir+`</listen>
  <auth>EXTERNAL</auth>
  <policy context="default">
    <allow send_destination="*" eavesdrop="true"/>
    <allow eavesdrop="true"/>
    <allow own="*"/>
  </policy>
</busconfig>
`), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(daemon, "--config-file="+config, "--nofork", "--print-address=1")
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	address, err := bufio.NewReader(out).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", strings.TrimSpace(address))
}

func connect(t *testing.T) *dbus.Conn {
	t.Helper()
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

// fakeWatcher is the desktop's StatusNotifierWatcher.
type fakeWatcher struct {
	mu    sync.Mutex
	items []string
}

func (w *fakeWatcher) RegisterStatusNotifierItem(service string) *dbus.Error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.items = append(w.items, service)
	return nil
}

func (w *fakeWatcher) RegisterStatusNotifierHost(string) *dbus.Error { return nil }

func (w *fakeWatcher) registered() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return slices.Clone(w.items)
}

func startWatcher(t *testing.T) *fakeWatcher {
	t.Helper()
	conn := connect(t)
	w := &fakeWatcher{}
	if err := conn.Export(w, sniWatcherPath, sniWatcher); err != nil {
		t.Fatal(err)
	}
	if _, err := prop.Export(conn, sniWatcherPath, prop.Map{sniWatcher: {
		"IsStatusNotifierHostRegistered": {Value: true},
		"ProtocolVersion":                {Value: int32(0)},
		"RegisteredStatusNotifierItems":  {Value: []string{}},
	}}); err != nil {
		t.Fatal(err)
	}
	if reply, err := conn.RequestName(sniWatcher, dbus.NameFlagDoNotQueue); err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		t.Fatalf("could not own the watcher's name: %v %v", reply, err)
	}
	return w
}

// fakeNotifier is the desktop's notification server.
type fakeNotifier struct {
	conn *dbus.Conn
	mu   sync.Mutex
	got  []notified
}

type notified struct {
	id                 uint32
	app, summary, body string
	actions            []string
	desktopEntry       string
	expire             int32
}

func (f *fakeNotifier) Notify(app string, _ uint32, _, summary, body string, actions []string,
	hints map[string]dbus.Variant, expire int32,
) (uint32, *dbus.Error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := notified{id: uint32(len(f.got) + 7), app: app, summary: summary, body: body, actions: actions, expire: expire}
	if v, ok := hints["desktop-entry"]; ok {
		n.desktopEntry, _ = v.Value().(string)
	}
	f.got = append(f.got, n)
	return n.id, nil
}

func (f *fakeNotifier) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"actions", "body"}, nil
}

func (f *fakeNotifier) CloseNotification(uint32) *dbus.Error { return nil }

func (f *fakeNotifier) GetServerInformation() (string, string, string, string, *dbus.Error) {
	return "fake", "kvit", "1", "1.2", nil
}

func (f *fakeNotifier) received() []notified {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.got)
}

func startNotifier(t *testing.T) *fakeNotifier {
	t.Helper()
	conn := connect(t)
	f := &fakeNotifier{conn: conn}
	if err := conn.Export(f, notifyPath, notifyName); err != nil {
		t.Fatal(err)
	}
	if reply, err := conn.RequestName(notifyName, dbus.NameFlagDoNotQueue); err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		t.Fatalf("could not own the notification server's name: %v %v", reply, err)
	}
	return f
}

// queued is a tray whose callbacks wait for the test to run them, as they
// wait for unison's UI thread in an app.
type queued chan func()

func (q queued) post(f func()) { q <- f }

// next runs the next callback the tray handed over.
func (q queued) next(t *testing.T, what string) {
	t.Helper()
	select {
	case f := <-q:
		f()
	case <-time.After(3 * time.Second):
		t.Fatalf("no callback for %s", what)
	}
}

// eventually waits for a condition the bus answers in its own time.
func eventually(t *testing.T, what string, cond func() bool) {
	t.Helper()
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		if cond() {
			return
		}
	}
	t.Fatalf("timed out waiting for %s", what)
}

// Without a StatusNotifierWatcher, as on WSLg, there is no notification
// area and no notifications.
func TestLinuxHasNoTrayWithoutAWatcher(t *testing.T) {
	privateBus(t)
	startNotifier(t)
	tray := newSystemTray(func(f func()) { f() })
	defer tray.Close()
	if tray.Available() || tray.Authorization() != Unsupported {
		t.Errorf("available %v, notifications %v, with no watcher", tray.Available(), tray.Authorization())
	}
}

// The icon registers with the watcher and answers what a desktop asks of
// it: its properties, its menu, clicks and menu choices, and notifications
// with their clicks. Hiding it releases the name the watcher follows.
func TestLinuxTrayOverDBus(t *testing.T) {
	privateBus(t)
	watcher := startWatcher(t)
	notifier := startNotifier(t)
	q := make(queued, 16)
	tray := newSystemTray(q.post)
	defer tray.Close()
	if !tray.Available() || tray.Authorization() != Authorized {
		t.Fatalf("available %v, notifications %v, with a watcher and a server", tray.Available(), tray.Authorization())
	}

	var ran []string
	var clicked []string
	tray.OnClick = func() { ran = append(ran, "click") }
	tray.OnNotificationClick = func(id string) { clicked = append(clicked, id) }
	tray.SetIcon(square(32, color.NRGBA{R: 0x12, G: 0x34, B: 0x56, A: 0xff}))
	tray.SetTooltip("Kvit Notes")
	tray.SetMenu([]TrayItem{
		{Text: "New Note", OnSelect: func() { ran = append(ran, "new") }},
		{Text: "Quick Capture…"},
		{Separator: true},
		{Text: "Show_Kvit", Disabled: true},
		{Text: "Quit", OnSelect: func() { ran = append(ran, "quit") }},
	})
	tray.Show()

	eventually(t, "the item to register", func() bool { return len(watcher.registered()) == 1 })
	name := watcher.registered()[0]
	if !strings.HasPrefix(name, "org.kde.StatusNotifierItem-") {
		t.Errorf("registered as %q", name)
	}

	desktop := connect(t)
	item := desktop.Object(name, sniPath)
	get := func(obj dbus.BusObject, iface, prop string) any {
		t.Helper()
		v, err := obj.GetProperty(iface + "." + prop)
		if err != nil {
			t.Fatalf("%s: %v", prop, err)
		}
		return v.Value()
	}
	if got := get(item, sniInterface, "Title"); got != "Kvit Notes" {
		t.Errorf("Title %v", got)
	}
	if got := get(item, sniInterface, "Menu"); got != menuPath {
		t.Errorf("Menu %v", got)
	}
	if got := get(item, sniInterface, "Status"); got != "Active" {
		t.Errorf("Status %v", got)
	}
	var pixmaps []sniPixmap
	if err := dbus.Store([]any{get(item, sniInterface, "IconPixmap")}, &pixmaps); err != nil {
		t.Fatal(err)
	}
	if len(pixmaps) != len(sniIconSizes) || pixmaps[0].Width != 16 || len(pixmaps[0].Data) != 16*16*4 {
		t.Fatalf("pixmaps: %d, first %dx%d", len(pixmaps), pixmaps[0].Width, pixmaps[0].Height)
	}
	if px := pixmaps[0].Data[:4]; !slices.Equal(px, []byte{0xff, 0x12, 0x34, 0x56}) {
		t.Errorf("a pixel is %x, want ARGB ff123456", px)
	}
	var tip sniToolTip
	if err := dbus.Store([]any{get(item, sniInterface, "ToolTip")}, &tip); err != nil || tip.Title != "Kvit Notes" {
		t.Errorf("ToolTip %+v %v", tip, err)
	}
	if got := get(desktop.Object(name, menuPath), menuInterface, "Version"); got != uint32(3) {
		t.Errorf("dbusmenu Version %v", got)
	}

	// The menu, as a desktop reads it to draw it.
	var revision uint32
	var layout menuLayout
	if err := desktop.Object(name, menuPath).Call(menuInterface+".GetLayout", 0, int32(0), int32(-1), []string{}).
		Store(&revision, &layout); err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, child := range layout.Children {
		var line menuLayout
		if err := dbus.Store([]any{child.Value()}, &line); err != nil {
			t.Fatal(err)
		}
		switch {
		case line.Props["type"].Value() == "separator":
			lines = append(lines, "---")
		case line.Props["enabled"].Value() == false:
			lines = append(lines, "("+line.Props["label"].Value().(string)+")")
		default:
			lines = append(lines, line.Props["label"].Value().(string))
		}
	}
	if want := []string{"New Note", "Quick Capture…", "---", "(Show__Kvit)", "Quit"}; !slices.Equal(lines, want) {
		t.Errorf("menu lines %q, want %q", lines, want)
	}

	// A click on the icon, then two lines chosen from the menu; the
	// disabled one does nothing.
	if err := item.Call(sniInterface+".Activate", 0, int32(10), int32(10)).Err; err != nil {
		t.Fatal(err)
	}
	q.next(t, "the click")
	for _, id := range []int32{1, 4, 5} {
		if err := desktop.Object(name, menuPath).Call(menuInterface+".Event", 0, id, "clicked",
			dbus.MakeVariant(""), uint32(0)).Err; err != nil {
			t.Fatal(err)
		}
	}
	q.next(t, "New Note")
	q.next(t, "Quit")
	if !slices.Equal(ran, []string{"click", "new", "quit"}) {
		t.Errorf("ran %q", ran)
	}

	// A notification, with a default action, and its click.
	if !tray.Notify("Kvit", "Note captured", "note:one") {
		t.Fatal("the notification was not posted")
	}
	eventually(t, "the notification", func() bool { return len(notifier.received()) == 1 })
	n := notifier.received()[0]
	if n.app != "Kvit Notes" || n.summary != "Kvit" || n.body != "Note captured" || !slices.Contains(n.actions, "default") ||
		n.desktopEntry != programName() || n.expire != -1 {
		t.Errorf("notification as the server got it: %+v", n)
	}
	// Another program's notification clicked is not this app's.
	_ = notifier.conn.Emit(notifyPath, notifyName+".ActionInvoked", uint32(999), "default")
	eventually(t, "the id to be recorded", func() bool {
		tr := tray.backend.(*sniTray)
		tr.mu.Lock()
		defer tr.mu.Unlock()
		return tr.posted[n.id] == "note:one"
	})
	if err := notifier.conn.Emit(notifyPath, notifyName+".ActionInvoked", n.id, "default"); err != nil {
		t.Fatal(err)
	}
	q.next(t, "the notification's click")
	if !slices.Equal(clicked, []string{"note:one"}) {
		t.Errorf("notification clicks %q", clicked)
	}

	// Hiding releases the name, which takes the item off the desktop.
	tray.Hide()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var owned bool
	if err := desktop.BusObject().CallWithContext(ctx, busName+".NameHasOwner", 0, name).Store(&owned); err != nil || owned {
		t.Errorf("the item's name is still owned after Hide (%v)", err)
	}
	// Shown again, it registers again.
	tray.Show()
	eventually(t, "the second registration", func() bool { return len(watcher.registered()) == 2 })
}
