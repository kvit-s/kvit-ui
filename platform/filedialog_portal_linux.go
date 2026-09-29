package platform

// The portal backend's D-Bus exchange: asking the desktop's FileChooser
// portal to show the dialog and waiting for its Response, while the
// application's windows take no input meanwhile, as runExternalDialog does
// for zenity and kdialog. The dialog is asked for off the UI thread and its
// end posted back with InvokeTask, so the windows keep drawing and answering
// the system while refusing presses and keys.

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

const (
	portalDest        = "org.freedesktop.portal.Desktop"
	portalPath        = dbus.ObjectPath("/org/freedesktop/portal/desktop")
	portalFileChooser = "org.freedesktop.portal.FileChooser"
	portalRequest     = "org.freedesktop.portal.Request"
	// portalResponseOK is the Response its signal reports the reader choosing
	// with; anything else is a cancel.
	portalResponseOK = 0
	// portalDirectoryVersion is the portal version that added choosing
	// folders.
	portalDirectoryVersion = 3
	// portalCallTimeout bounds the calls that set the dialog up; the wait for // the reader itself has none, as a runaway program's has none.
	portalCallTimeout = 2 * time.Second
)

// portalTokenNumber numbers the dialog's handle tokens, which have to be
// valid object path elements.
var portalTokenNumber atomic.Int32

// openWithPortal shows the desktop's portal FileChooser and waits for the
// reader. It reports false when there is no portal to ask, so the caller
// falls back to unison's own dialog.
func (d FileDialog) openWithPortal() ([]string, bool) {
	conn, ok := portalConnection(d.Directories)
	if !ok {
		return nil, false
	}
	owner := unison.ActiveWindow()
	wnd, err := unison.NewWindow("", unison.WindowKindWindowOption(unison.WindowKindDialog), unison.UndecoratedWindowOption(), unison.NotResizableWindowOption())
	if err != nil {
		errs.Log(err, "dialog", d.Title)
		conn.Close()
		return nil, false
	}
	if owner != nil {
		at := owner.ContentRect()
		wnd.SetFrameRect(geom.NewRect(at.X, at.Y, 1, 1))
	}
	var outcome portalOutcome
	unison.InvokeTaskAfter(func() {
		wnd.Hide()
		go func() {
			outcome = d.runPortalDialog(conn)
			unison.InvokeTask(func() {
				wnd.StopModal(unison.ModalResponseOK)
			})
		}()
	}, externalDialogStart)
	wnd.RunModal()
	conn.Close()
	if !outcome.shown {
		return d.openWithUnison(), true
	}
	return outcome.paths, true
}

// portalOutcome is what asking the portal came to: the paths chosen, and
// whether the dialog appeared at all. A cancel appears with no paths.
type portalOutcome struct {
	paths []string
	shown bool
}

// runPortalDialog asks the portal to show the dialog and waits for its
// Response, off the UI thread; the caller posts the modal's end back.
func (d FileDialog) runPortalDialog(conn *dbus.Conn) portalOutcome {
	token := fmt.Sprintf("kvit%d", portalTokenNumber.Add(1))
	options := map[string]dbus.Variant{
		"handle_token": dbus.MakeVariant(token), "modal": dbus.MakeVariant(true), "multiple": dbus.MakeVariant(d.Multiple), "directory": dbus.MakeVariant(d.Directories)}
	if !d.Directories {
		if filters := portalFilters(d.filters()); len(filters) > 0 {
			options["filters"] = dbus.MakeVariant(filters)
		}
	}
	if d.Folder != "" {
		options["current_folder"] = dbus.MakeVariant(portalFolder(d.Folder))
	}
	signals := make(chan *dbus.Signal, 8)
	conn.Signal(signals)
	defer conn.RemoveSignal(signals)
	// The match is added before the call, so a Response arriving while the
	// call returns is still caught; the loop below keeps only this
	// dialog's handle.
	if err := conn.AddMatchSignal(dbus.WithMatchSender(portalDest),
		dbus.WithMatchInterface(portalRequest), dbus.WithMatchMember("Response")); err != nil {
		errs.Log(err, "dialog", d.Title)
		return portalOutcome{shown: true}
	}
	ctx, cancel := context.WithTimeout(context.Background(), portalCallTimeout)
	var handle dbus.ObjectPath
	err := conn.Object(portalDest, portalPath).CallWithContext(ctx,
		portalFileChooser+".OpenFile", 0, "", d.Title, options).Store(&handle)
	cancel()
	if err != nil {
		errs.Log(err, "dialog", d.Title)
		return portalOutcome{}
	}
	for s := range signals {
		if s.Path != handle || s.Name != portalRequest+".Response" {
			continue
		}
		var response uint32
		var results map[string]dbus.Variant
		if dbus.Store(s.Body, &response, &results) != nil || response != portalResponseOK {
			return portalOutcome{shown: true}
		}
		var uris []string
		if v, ok := results["uris"]; ok {
			if v.Store(&uris) != nil {
				errs.Log(errs.New("the file dialog answered without files"), "dialog", d.Title)
				return portalOutcome{shown: true}
			}
		}
		paths := portalURIPaths(uris)
		if !d.Multiple && len(paths) > 1 {
			paths = paths[:1]
		}
		return portalOutcome{paths: paths, shown: true}
	}
	return portalOutcome{shown: true}
}

// isWSL reports whether the program runs under Windows Subsystem for Linux:
// the environment names the distribution, and the kernel names Microsoft.
// There the portal's dialog shows the Windows files while the program works
// with the Linux ones, so the portal is skipped and unison's own dialog,
// which opens where the program is, is used instead.
func isWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != "" {
		return true
	}
	for _, file := range []string{"/proc/sys/kernel/osrelease", "/proc/version"} {
		if release, err := os.ReadFile(file); err == nil {
			lowered := strings.ToLower(string(release))
			if strings.Contains(lowered, "microsoft") || strings.Contains(lowered, "wsl") {
				return true
			}
		}
	}
	return false
}

// portalConnection connects to the session bus and reports whether the
// portal can show this dialog: no bus or portal, a portal from before
// folders could be chosen, or Windows Subsystem for Linux, is none.
func portalConnection(directories bool) (*dbus.Conn, bool) {
	if isWSL() {
		return nil, false
	}
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, false
	}
	if !nameOwned(conn, portalDest) {
		conn.Close()
		return nil, false
	}
	if directories && portalChooserVersion(conn) < portalDirectoryVersion {
		conn.Close()
		return nil, false
	}
	return conn, true
}

// portalChooserVersion is the portal FileChooser's version, 0 when it would
// not say.
func portalChooserVersion(conn *dbus.Conn) uint32 {
	ctx, cancel := context.WithTimeout(context.Background(), busTimeout)
	defer cancel()
	var v dbus.Variant
	if conn.Object(portalDest, portalPath).CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, portalFileChooser, "version").Store(&v) != nil {
		return 0
	}
	version, _ := v.Value().(uint32)
	return version
}
