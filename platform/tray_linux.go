package platform

// The notification area on Linux, over the session D-Bus, as the own tray
// does it: the icon is a StatusNotifierItem at /StatusNotifierItem,
// registered with the desktop's org.kde.StatusNotifierWatcher; its menu is a
// com.canonical.dbusmenu at /MenuBar, which the desktop draws itself; and a
// notification is org.freedesktop.Notifications' Notify, whose click comes
// back as its ActionInvoked signal with the "default" action.
//
// A session without a StatusNotifierWatcher, or with one that has no host
// drawing the items, has no notification area: GNOME without its
// AppIndicator extension, and WSLg. There the tray is unavailable, as the
// is.

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	sniWatcher     = "org.kde.StatusNotifierWatcher"
	sniWatcherPath = dbus.ObjectPath("/StatusNotifierWatcher")
	sniInterface   = "org.kde.StatusNotifierItem"
	sniPath        = dbus.ObjectPath("/StatusNotifierItem")
	menuInterface  = "com.canonical.dbusmenu"
	menuPath       = dbus.ObjectPath("/MenuBar")
	notifyName     = "org.freedesktop.Notifications"
	notifyPath     = dbus.ObjectPath("/org/freedesktop/Notifications")
	busName        = "org.freedesktop.DBus"
	// busTimeout bounds each call the tray waits for, so a wedged desktop
	// service cannot stall the app.
	busTimeout = 2 * time.Second
)

// sniIconSizes are the sizes the icon is offered at; the desktop picks.
var sniIconSizes = []int{16, 22, 24, 32, 48, 64}

// sniCount numbers the item names a process asks for, as the specification
// has them: org.kde.StatusNotifierItem-<process>-<number>.
var sniCount atomic.Int32

// sniTray is the notification area of a desktop with a StatusNotifierWatcher.
type sniTray struct {
	events trayEvents
	conn   *dbus.Conn
	id     string // the program's name, which the desktop files items under
	auth   Authorization

	exportOnce sync.Once
	exportErr  error
	props      *prop.Properties

	mu       sync.Mutex
	state    trayState
	name     string // the bus name the item is registered under; "" while hidden
	revision uint32 // the menu's, which grows with every change
	posted   map[uint32]string
}

// newTrayBackend connects to the session bus and finds the desktop's
// StatusNotifierWatcher. Without either, there is no notification area.
func newTrayBackend(events trayEvents) trayBackend {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return noTray{}
	}
	if !hasTrayHost(conn) {
		conn.Close()
		return noTray{}
	}
	t := &sniTray{events: events, conn: conn, id: programName(), posted: map[uint32]string{}}
	if nameReachable(conn, notifyName) {
		t.auth = Authorized
	}
	_ = conn.AddMatchSignal(dbus.WithMatchObjectPath(notifyPath), dbus.WithMatchInterface(notifyName))
	_ = conn.AddMatchSignal(dbus.WithMatchSender(busName), dbus.WithMatchInterface(busName),
		dbus.WithMatchMember("NameOwnerChanged"), dbus.WithMatchArg(0, sniWatcher))
	signals := make(chan *dbus.Signal, 16)
	conn.Signal(signals)
	go t.listen(signals)
	return t
}

// hasTrayHost reports whether a StatusNotifierWatcher runs and says a host
// draws its items.
func hasTrayHost(conn *dbus.Conn) bool {
	if !nameOwned(conn, sniWatcher) {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), busTimeout)
	defer cancel()
	var v dbus.Variant
	err := conn.Object(sniWatcher, sniWatcherPath).CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0,
		sniWatcher, "IsStatusNotifierHostRegistered").Store(&v)
	host, _ := v.Value().(bool)
	return err == nil && host
}

func nameOwned(conn *dbus.Conn, name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), busTimeout)
	defer cancel()
	var owned bool
	err := conn.BusObject().CallWithContext(ctx, busName+".NameHasOwner", 0, name).Store(&owned)
	return err == nil && owned
}

// nameReachable reports whether a name has an owner, or the bus would start
// one when it is called.
func nameReachable(conn *dbus.Conn, name string) bool {
	if nameOwned(conn, name) {
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), busTimeout)
	defer cancel()
	var names []string
	if conn.BusObject().CallWithContext(ctx, busName+".ListActivatableNames", 0).Store(&names) != nil {
		return false
	}
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

// programName is the program's file name without its extension, which is
// also the name of its .desktop file.
func programName() string {
	name := filepath.Base(os.Args[0])
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// listen follows the signals the tray asked for: a notification clicked or
// closed, and the watcher starting again, which forgets every item.
func (t *sniTray) listen(signals <-chan *dbus.Signal) {
	for s := range signals {
		switch s.Name {
		case notifyName + ".ActionInvoked":
			var id uint32
			var action string
			if dbus.Store(s.Body, &id, &action) != nil || action != "default" {
				continue
			}
			t.mu.Lock()
			appID, ours := t.posted[id]
			t.mu.Unlock()
			if ours {
				t.events.notificationClicked(appID)
			}
		case notifyName + ".NotificationClosed":
			if len(s.Body) > 0 {
				if id, ok := s.Body[0].(uint32); ok {
					t.mu.Lock()
					delete(t.posted, id)
					t.mu.Unlock()
				}
			}
		case busName + ".NameOwnerChanged":
			var name, from, to string
			if dbus.Store(s.Body, &name, &from, &to) != nil || name != sniWatcher || to == "" {
				continue
			}
			t.mu.Lock()
			registered := t.name
			t.mu.Unlock()
			if registered != "" {
				t.register(registered)
			}
		}
	}
}

func (t *sniTray) available() bool { return true }

func (t *sniTray) authorization() Authorization { return t.auth }

// requestAuthorization has nothing to ask: freedesktop notification servers
// show what they are sent.
func (t *sniTray) requestAuthorization() {}

// sniPixmap is one size of an icon: ARGB32 in network byte order, which is
// A, R, G, B for each pixel, row after row.
type sniPixmap struct {
	Width, Height int32
	Data          []byte
}

// sniToolTip is the item's ToolTip property.
type sniToolTip struct {
	IconName    string
	IconPixmap  []sniPixmap
	Title       string
	Description string
}

func pixmaps(icons []image.Image) []sniPixmap {
	out := []sniPixmap{}
	for _, size := range sniIconSizes {
		img := iconAt(icons, size)
		if img == nil {
			return out
		}
		data := make([]byte, 0, len(img.Pix))
		for i := 0; i+3 < len(img.Pix); i += 4 {
			data = append(data, img.Pix[i+3], img.Pix[i], img.Pix[i+1], img.Pix[i+2])
		}
		out = append(out, sniPixmap{Width: int32(size), Height: int32(size), Data: data})
	}
	return out
}

// title is the item's name for the desktop: the tooltip, or the program's
// name without one.
func (t *sniTray) title(s trayState) string {
	if s.tooltip != "" {
		return s.tooltip
	}
	return t.id
}

// export puts the item and its menu on the bus, the first time the icon is
// shown.
func (t *sniTray) export() error {
	t.exportOnce.Do(func() {
		t.exportErr = t.exportObjects()
	})
	return t.exportErr
}

func (t *sniTray) exportObjects() error {
	item := sniItem{t}
	if err := t.conn.Export(item, sniPath, sniInterface); err != nil {
		return err
	}
	props, err := prop.Export(t.conn, sniPath, prop.Map{sniInterface: {
		"Category":            {Value: "ApplicationStatus"},
		"Id":                  {Value: t.id},
		"Title":               {Value: t.id},
		"Status":              {Value: "Active"},
		"WindowId":            {Value: int32(0)},
		"IconThemePath":       {Value: ""},
		"IconName":            {Value: ""},
		"IconPixmap":          {Value: []sniPixmap{}},
		"OverlayIconName":     {Value: ""},
		"OverlayIconPixmap":   {Value: []sniPixmap{}},
		"AttentionIconName":   {Value: ""},
		"AttentionIconPixmap": {Value: []sniPixmap{}},
		"AttentionMovieName":  {Value: ""},
		"ToolTip":             {Value: sniToolTip{IconPixmap: []sniPixmap{}}},
		"ItemIsMenu":          {Value: false},
		"Menu":                {Value: menuPath},
	}})
	if err != nil {
		return err
	}
	t.props = props
	if err := t.conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: string(sniPath),
		Interfaces: []introspect.Interface{introspect.IntrospectData, prop.IntrospectData, {
			Name:       sniInterface,
			Methods:    introspect.Methods(item),
			Properties: props.Introspection(sniInterface),
			Signals: []introspect.Signal{{Name: "NewTitle"}, {Name: "NewIcon"}, {Name: "NewToolTip"},
				{Name: "NewStatus", Args: []introspect.Arg{{Name: "status", Type: "s"}}}},
		}},
	}), sniPath, "org.freedesktop.DBus.Introspectable"); err != nil {
		return err
	}

	menu := dbusMenu{t}
	if err := t.conn.Export(menu, menuPath, menuInterface); err != nil {
		return err
	}
	menuProps, err := prop.Export(t.conn, menuPath, prop.Map{menuInterface: {
		"Version":       {Value: uint32(3)},
		"TextDirection": {Value: "ltr"},
		"Status":        {Value: "normal"},
		"IconThemePath": {Value: []string{}},
	}})
	if err != nil {
		return err
	}
	return t.conn.Export(introspect.NewIntrospectable(&introspect.Node{
		Name: string(menuPath),
		Interfaces: []introspect.Interface{introspect.IntrospectData, prop.IntrospectData, {
			Name:       menuInterface,
			Methods:    introspect.Methods(menu),
			Properties: menuProps.Introspection(menuInterface),
			Signals: []introspect.Signal{
				{Name: "LayoutUpdated", Args: []introspect.Arg{{Name: "revision", Type: "u"}, {Name: "parent", Type: "i"}}},
				{Name: "ItemsPropertiesUpdated", Args: []introspect.Arg{{Name: "updatedProps", Type: "a(ia{sv})"}, {Name: "removedProps", Type: "a(ias)"}}},
				{Name: "ItemActivationRequested", Args: []introspect.Arg{{Name: "id", Type: "i"}, {Name: "timestamp", Type: "u"}}},
			},
		}},
	}), menuPath, "org.freedesktop.DBus.Introspectable")
}

func (t *sniTray) show(s trayState) {
	if t.export() != nil {
		return
	}
	t.mu.Lock()
	t.state = s
	t.revision++
	revision, name := t.revision, t.name
	t.mu.Unlock()

	title := t.title(s)
	icons := pixmaps(s.icons)
	t.props.SetMust(sniInterface, "Title", title)
	t.props.SetMust(sniInterface, "IconPixmap", icons)
	t.props.SetMust(sniInterface, "ToolTip", sniToolTip{IconPixmap: []sniPixmap{}, Title: title})
	t.props.SetMust(sniInterface, "Status", "Active")
	for _, signal := range []string{"NewTitle", "NewIcon", "NewToolTip"} {
		_ = t.conn.Emit(sniPath, sniInterface+"."+signal)
	}
	_ = t.conn.Emit(sniPath, sniInterface+".NewStatus", "Active")
	_ = t.conn.Emit(menuPath, menuInterface+".LayoutUpdated", revision, int32(0))
	if name == "" {
		t.appear()
	}
}

// appear takes a bus name for the item and registers it with the watcher.
// The name is the specification's, so that releasing it takes the item
// away; when it cannot be had, the connection's own name does.
func (t *sniTray) appear() {
	name := fmt.Sprintf("org.kde.StatusNotifierItem-%d-%d", os.Getpid(), sniCount.Add(1))
	if reply, err := t.conn.RequestName(name, dbus.NameFlagDoNotQueue); err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		name = t.conn.Names()[0]
	}
	t.mu.Lock()
	t.name = name
	t.mu.Unlock()
	t.register(name)
}

func (t *sniTray) register(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), busTimeout)
	defer cancel()
	_ = t.conn.Object(sniWatcher, sniWatcherPath).CallWithContext(ctx, sniWatcher+".RegisterStatusNotifierItem", 0, name).Err
}

// hide releases the item's name, which the watcher follows and takes the
// item away for. An item registered under the connection's own name is
// marked Passive instead, which desktops show out of the way or not at all.
func (t *sniTray) hide() {
	t.mu.Lock()
	name := t.name
	t.name = ""
	t.mu.Unlock()
	if name == "" || t.props == nil {
		return
	}
	t.props.SetMust(sniInterface, "Status", "Passive")
	_ = t.conn.Emit(sniPath, sniInterface+".NewStatus", "Passive")
	if !strings.HasPrefix(name, ":") {
		_, _ = t.conn.ReleaseName(name)
	}
}

// notify sends the notification with a "default" action, which is what a
// click on its body invokes. The call is made off the UI thread, as the
// server may take its time answering.
func (t *sniTray) notify(n Notification) {
	t.mu.Lock()
	app := t.title(t.state)
	t.mu.Unlock()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), busTimeout)
		defer cancel()
		hints := map[string]dbus.Variant{"desktop-entry": dbus.MakeVariant(t.id)}
		var id uint32
		err := t.conn.Object(notifyName, notifyPath).CallWithContext(ctx, notifyName+".Notify", 0,
			app, uint32(0), "", n.Title, n.Message, []string{"default", "Open"}, hints, int32(-1)).Store(&id)
		if err == nil {
			t.mu.Lock()
			t.posted[id] = n.ID
			t.mu.Unlock()
		}
	}()
}

func (t *sniTray) close() {
	t.hide()
	t.conn.Close()
}

// entries is the menu as last shown, and its revision.
func (t *sniTray) entries() ([]trayEntry, uint32) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state.menu, t.revision
}

// sniItem is the org.kde.StatusNotifierItem interface's methods. The
// desktop draws the menu from /MenuBar, so ContextMenu has nothing to do.
type sniItem struct{ t *sniTray }

// Activate is a click on the icon.
func (i sniItem) Activate(_, _ int32) *dbus.Error { i.t.events.clicked(); return nil }

// SecondaryActivate is a middle click, which the app has no use for.
func (i sniItem) SecondaryActivate(_, _ int32) *dbus.Error { return nil }

// ContextMenu is asked for only by a desktop that does not draw the menu.
func (i sniItem) ContextMenu(_, _ int32) *dbus.Error { return nil }

// Scroll is the wheel over the icon.
func (i sniItem) Scroll(_ int32, _ string) *dbus.Error { return nil }

// dbusMenu is the com.canonical.dbusmenu interface's methods: the menu is a
// root, id 0, with one item for each line, id index+1.
type dbusMenu struct{ t *sniTray }

type menuLayout struct {
	ID       int32
	Props    map[string]dbus.Variant
	Children []dbus.Variant
}

type menuProps struct {
	ID    int32
	Props map[string]dbus.Variant
}

type menuEvent struct {
	ID        int32
	EventID   string
	Data      dbus.Variant
	Timestamp uint32
}

// entryProps are a line's properties. dbusmenu marks an access key with an
// underscore, so one in the text is doubled to show as itself.
func entryProps(e *trayEntry, names []string) map[string]dbus.Variant {
	all := map[string]dbus.Variant{"children-display": dbus.MakeVariant("submenu")}
	if e != nil {
		all = map[string]dbus.Variant{"visible": dbus.MakeVariant(true)}
		if e.separator {
			all["type"] = dbus.MakeVariant("separator")
		} else {
			all["label"] = dbus.MakeVariant(strings.ReplaceAll(e.text, "_", "__"))
			all["enabled"] = dbus.MakeVariant(!e.disabled)
		}
	}
	if len(names) == 0 {
		return all
	}
	some := map[string]dbus.Variant{}
	for _, n := range names {
		if v, ok := all[n]; ok {
			some[n] = v
		}
	}
	return some
}

// entry is the line with an id, nil for the root, and whether there is one.
func entry(entries []trayEntry, id int32) (*trayEntry, bool) {
	if id == 0 {
		return nil, true
	}
	if id < 0 || int(id) > len(entries) {
		return nil, false
	}
	return &entries[id-1], true
}

func (m dbusMenu) GetLayout(parent, depth int32, names []string) (uint32, menuLayout, *dbus.Error) {
	entries, revision := m.t.entries()
	e, ok := entry(entries, parent)
	if !ok {
		return 0, menuLayout{}, dbus.MakeFailedError(fmt.Errorf("no menu item %d", parent))
	}
	layout := menuLayout{ID: parent, Props: entryProps(e, names), Children: []dbus.Variant{}}
	if parent == 0 && depth != 0 {
		for i := range entries {
			layout.Children = append(layout.Children, dbus.MakeVariant(menuLayout{
				ID: int32(i + 1), Props: entryProps(&entries[i], names), Children: []dbus.Variant{},
			}))
		}
	}
	return revision, layout, nil
}

func (m dbusMenu) GetGroupProperties(ids []int32, names []string) ([]menuProps, *dbus.Error) {
	entries, _ := m.t.entries()
	if len(ids) == 0 {
		for i := 0; i <= len(entries); i++ {
			ids = append(ids, int32(i))
		}
	}
	out := []menuProps{}
	for _, id := range ids {
		if e, ok := entry(entries, id); ok {
			out = append(out, menuProps{ID: id, Props: entryProps(e, names)})
		}
	}
	return out, nil
}

func (m dbusMenu) GetProperty(id int32, name string) (dbus.Variant, *dbus.Error) {
	entries, _ := m.t.entries()
	e, ok := entry(entries, id)
	if ok {
		if v, ok := entryProps(e, []string{name})[name]; ok {
			return v, nil
		}
	}
	return dbus.Variant{}, dbus.MakeFailedError(fmt.Errorf("no property %q on menu item %d", name, id))
}

// Event is what the person did to a line; "clicked" chooses it.
func (m dbusMenu) Event(id int32, eventID string, _ dbus.Variant, _ uint32) *dbus.Error {
	entries, _ := m.t.entries()
	e, ok := entry(entries, id)
	if !ok {
		return dbus.MakeFailedError(fmt.Errorf("no menu item %d", id))
	}
	if eventID == "clicked" && e != nil && !e.separator && !e.disabled {
		m.t.events.chose(int(id) - 1)
	}
	return nil
}

func (m dbusMenu) EventGroup(events []menuEvent) ([]int32, *dbus.Error) {
	missing := []int32{}
	for _, ev := range events {
		if m.Event(ev.ID, ev.EventID, ev.Data, ev.Timestamp) != nil {
			missing = append(missing, ev.ID)
		}
	}
	return missing, nil
}

// AboutToShow says the menu needs no update before it opens: it is sent
// whole whenever it changes.
func (m dbusMenu) AboutToShow(int32) (bool, *dbus.Error) { return false, nil }

func (m dbusMenu) AboutToShowGroup([]int32) ([]int32, []int32, *dbus.Error) {
	return []int32{}, []int32{}, nil
}
