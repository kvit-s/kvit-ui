package uitest_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// A session opens a component in the theme and size asked for, and changes
// either as a reader would.
func TestASessionOpensInTheThemeAndSizeAskedFor(t *testing.T) {
	var save *kvitui.Button
	s := uitest.Open(t, uitest.Options{Theme: tokens.Dark, InterfaceSize: 20}, func(ui *kvitui.UI) unison.Paneler {
		save = kvitui.NewButton(ui, "Save")
		return kvitui.Left(save)
	})
	if n := s.Node(save); n == nil || n.Role != role.Button || n.Name != "Save" {
		t.Fatalf("the button's node is %+v", n)
	}
	ground := func() uint32 {
		r, g, b, _ := s.Capture().At(2, 2).RGBA()
		return (r>>8)<<16 | (g>>8)<<8 | b>>8
	}
	dark := ground()
	s.SetTheme(tokens.Light)
	if light := ground(); light == dark {
		t.Error("changing the theme left the window's ground as it was")
	}
	var h float32
	s.Do(func() { h = save.FrameRect().Height })
	if want := float32(s.UI.Interface.ControlHeight()); h != want {
		t.Errorf("at 20 px the button is %v tall, want %v", h, want)
	}
	s.SetInterfaceSize(12)
	s.Do(func() { h = save.FrameRect().Height })
	if want := float32(s.UI.Interface.ControlHeight()); h != want {
		t.Errorf("at 12 px the button is %v tall, want %v", h, want)
	}
}

// The naming check finds a control that says nothing to a screen reader.
func TestTheNamingCheckFindsAControlWithNoName(t *testing.T) {
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		return kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewIconButton(ui, "trash", "Delete"), kvitui.NewIconButton(ui, "pencil", ""))
	})
	if got := uitest.Unnamed(s.Tree()); len(got) != 1 {
		t.Errorf("found %d unnamed controls, want the one: %v", len(got), got)
	}
	if n := s.Named("Delete"); n == nil || n.Role != role.Button {
		t.Errorf("the named button's node is %+v", n)
	}
}
