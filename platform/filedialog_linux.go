package platform

import (
	"errors"
	"os"
	"os/exec"
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// openNative runs kdialog in a KDE session and zenity elsewhere, the order
// unison's own picker tries them in; with neither installed, the desktop's
// portal FileChooser, which draws the system's own GTK dialog, except under
// Windows Subsystem for Linux, where the portal is skipped; with no portal
// either, unison's own dialog.
func (d FileDialog) openNative() []string {
	kdialog, _ := exec.LookPath("kdialog")
	zenity, _ := exec.LookPath("zenity")
	switch {
	case kdialog != "" && os.Getenv("KDE_FULL_SESSION") != "":
		return runExternalDialog(exec.Command(kdialog, d.kdialogArgs()...), d.Multiple)
	case zenity != "":
		return runExternalDialog(exec.Command(zenity, d.zenityArgs()...), d.Multiple)
	case kdialog != "":
		return runExternalDialog(exec.Command(kdialog, d.kdialogArgs()...), d.Multiple)
	}
	if paths, ok := d.openWithPortal(); ok {
		return paths
	}
	return d.openWithUnison()
}

// externalDialogStart is how long the stand-in window waits before it
// hides and the dialog program starts: long enough for the event loop to
// have finished bringing the stand-in to the front, since taking it off the
// screen while the system is still activating it asks X for the focus on a
// window that is no longer there.
const externalDialogStart = 20 * time.Millisecond

// runExternalDialog runs a dialog program, zenity or kdialog, while the
// application's windows take no input, and returns the files it printed.
// unison keeps its own picker modal the same way: a window of its own runs
// a modal event loop until the program exits, so the windows keep drawing
// and answering the system while refusing presses and keys. unison marks
// its stand-in to be never shown, through a field this package cannot set;
// this stand-in is a one-pixel undecorated dialog over the active window,
// hidden as soon as it has been brought to the front.
func runExternalDialog(cmd *exec.Cmd, multiple bool) []string {
	owner := unison.ActiveWindow()
	wnd, err := unison.NewWindow("", unison.WindowKindWindowOption(unison.WindowKindDialog),
		unison.UndecoratedWindowOption(), unison.NotResizableWindowOption())
	if err != nil {
		errs.Log(err, "cmd", cmd.String())
		return nil
	}
	if owner != nil {
		at := owner.ContentRect()
		wnd.SetFrameRect(geom.NewRect(at.X, at.Y, 1, 1))
	}
	var paths []string
	unison.InvokeTaskAfter(func() {
		wnd.Hide()
		go func() {
			out, err := cmd.Output()
			var exit *exec.ExitError
			if err != nil && !errors.As(err, &exit) {
				// The program could not be run at all; a cancel is an exit
				// status of 1 and says nothing.
				errs.Log(err, "cmd", cmd.String())
			}
			unison.InvokeTask(func() {
				if err == nil {
					paths = chosenPaths(out, multiple)
				}
				wnd.StopModal(unison.ModalResponseOK)
			})
		}()
	}, externalDialogStart)
	wnd.RunModal()
	return paths
}
