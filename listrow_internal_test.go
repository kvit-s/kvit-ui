package kvitui

import (
	"runtime"
	"testing"
	"weak"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// The window's pointer watch, which a row with HoverWhole is registered with
// once it is laid out in a window, does not keep a row nothing else holds:
// an application such as a conversation's transcript makes and drops rows
// as its turns are drawn. The row is registered without being put in the
// window, since unison keeps a panel taken out of an open window for a
// while, a row or not. This test is inside the package to reach the row's
// own registration, and so starts the headless screen itself rather than
// through uitest, which imports the package.
func TestAListRowIsNotKeptByThePointerWatch(t *testing.T) {
	ui, err := New(Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 400, Height: 300})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	// The UI outlives the row, as it does in an application.
	defer runtime.KeepAlive(ui)
	var w *unison.Window
	screen.Do(func() {
		if w, err = unison.NewWindow("Rows"); err != nil {
			t.Error(err)
			return
		}
		w.SetContentRect(geom.NewRect(0, 0, 400, 300))
		w.ToFront()
	})
	screen.Sync()
	if w == nil {
		t.Fatal("no window")
	}
	var gone weak.Pointer[ListRow]
	registered := false
	screen.Do(func() {
		row := NewListRow(ui, NewLabel(ui, "Ran a command"))
		row.HoverWhole = true
		row.watchPointerIn(w)
		listeners, _ := w.ClientData()[pointerWatchedKey].(pointerListeners)
		for owner := range listeners {
			registered = registered || owner.Value() == row.AsPanel()
		}
		gone = weak.Make(row)
	})
	if !registered {
		t.Fatal("the row is not registered with the window's pointer watch")
	}
	screen.MouseMove(geom.NewPoint(20, 20), 0)
	screen.Sync()
	for range 10 {
		runtime.GC()
		if gone.Value() == nil {
			return
		}
		screen.Sync()
	}
	t.Error("a row nothing else holds is kept by the window's pointer watch")
}
