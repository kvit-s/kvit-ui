package kvitui_test

import (
	"path/filepath"
	"testing"

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
