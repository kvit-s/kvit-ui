package kvitui_test

import (
	"path/filepath"
	"runtime"
	"testing"
	"weak"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/thememode"
)

func TestTheThemeReachesUnison(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 200, Height: 100})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	for _, id := range tokens.BuiltInThemes() {
		screen.Do(func() { ui.Theme.SetThemeID(id) })
		tk := tokens.TokensFor(id)
		var surface, focus unison.Color
		var mode thememode.Enum
		screen.Do(func() {
			surface, focus = unison.ThemeSurface.GetColor(), unison.ThemeFocus.GetColor()
			mode = unison.CurrentThemeMode()
		})
		if surface != kvitui.Color(tk.WindowBackground) || focus != kvitui.Color(tk.Accent) {
			t.Errorf("%s: unison's surface and focus are not the theme's window background and accent", id)
		}
		wantDark := id == tokens.Dark || id == tokens.HighContrast
		if (mode == thememode.Dark) != wantDark {
			t.Errorf("%s: unison is in mode %v", id, mode)
		}
	}
}

func TestChoicesPersistAcrossRuns(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ui.json")
	{
		ui, err := kvitui.New(kvitui.Options{SettingsPath: path, IgnoreDesktop: true})
		if err != nil {
			t.Fatal(err)
		}
		ui.Theme.SetThemeID(tokens.Sepia)
		ui.Interface.SetFontSize(15)
		ui.Typography.SetBaseSize(17)
		if err := ui.Settings.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	ui, err := kvitui.New(kvitui.Options{SettingsPath: path, IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	if ui.Theme.ThemeID() != tokens.Sepia || ui.Interface.FontSize() != 15 || ui.Typography.BaseSize() != 17 {
		t.Errorf("reopened as %s, %d px, base %d", ui.Theme.ThemeID(), ui.Interface.FontSize(), ui.Typography.BaseSize())
	}
}

// A window a Kvit control was shown in, once closed, is not kept by the UI:
// menus, popovers and dialogs are windows too, and a UI that kept each one
// it had drawn a focus ring in held every one of them, with all its panels,
// for the life of the program.
func TestAClosedWindowIsNotKeptByTheUI(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 400, Height: 300})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	// The UI outlives the window, as it does in an application.
	defer runtime.KeepAlive(ui)
	var closed weak.Pointer[unison.Window]
	var watched bool
	screen.Do(func() {
		w, err := unison.NewWindow("Closing")
		if err != nil {
			t.Error(err)
			return
		}
		w.Content().SetLayout(&unison.FlexLayout{Columns: 1})
		w.Content().AddChild(kvitui.NewButton(ui, "Press"))
		w.Pack()
		w.ToFront()
		closed = weak.Make(w)
	})
	screen.Sync()
	screen.Do(func() {
		w := closed.Value()
		// Watched: the window's key presses go through the UI's watch.
		watched = w != nil && w.KeyDownCallback != nil
		w.Dispose()
	})
	screen.Sync()
	if !watched {
		t.Fatal("the window was not watched")
	}
	for range 10 {
		runtime.GC()
		if closed.Value() == nil {
			return
		}
		screen.Sync()
	}
	t.Error("a closed window is still held")
}

// A Kvit window, once closed, is not kept by the UI: the panels it draws
// its components' spills in go with it. Its body is a label rather than a
// control: closing the only window ends the headless session, and a
// control under the pointer has a tooltip task waiting that would then
// never run, holding the control until the process ends.
func TestAClosedKvitWindowIsNotKeptByTheUI(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 800, Height: 600})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	defer runtime.KeepAlive(ui)
	// Nothing eases, so no frame of the sidebar's easing is still waiting
	// to run when the window closes.
	ui.Theme.SetReducedMotion(true)
	var closed weak.Pointer[unison.Window]
	screen.Do(func() {
		w, err := kvitui.NewWindow(ui, "Closing")
		if err != nil {
			t.Error(err)
			return
		}
		w.SetBody(kvitui.NewLabel(ui, "Closing"))
		w.ToFront()
		closed = weak.Make(w.Window)
	})
	screen.Sync()
	screen.Do(func() { closed.Value().AttemptClose() })
	screen.Sync()
	for range 10 {
		runtime.GC()
		if closed.Value() == nil {
			return
		}
		screen.Sync()
	}
	t.Error("a closed Kvit window is still held")
}

// A disclosure that is no longer used is not kept by the UI for its body,
// which clips the spills of what is inside it.
func TestADisclosureNoLongerUsedIsNotKeptByTheUI(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.KeepAlive(ui)
	gone := weak.Make(kvitui.NewDisclosure(ui, "Details", kvitui.NewLabel(ui, "inside")))
	for range 10 {
		runtime.GC()
		if gone.Value() == nil {
			return
		}
	}
	t.Error("a disclosure no longer used is still held")
}
