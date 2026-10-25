package kvitui_test

import (
	"slices"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// barCommands are an app's commands: standard ones, which the bar places by
// itself, and two of the app's own.
type barCommands struct {
	save, undo, cut, find, settings, about, help, quit, sidebar, importStatement *kvitui.Command
	saved, imported                                                              int
	sidebarShown                                                                 bool
}

func newBarCommands() *barCommands {
	c := &barCommands{sidebarShown: true}
	run := func() {}
	c.save = &kvitui.Command{Std: kvitui.StdSave, Run: func() { c.saved++ }}
	c.undo = &kvitui.Command{Std: kvitui.StdUndo, Run: run}
	c.cut = &kvitui.Command{Std: kvitui.StdCut, Run: run}
	c.find = &kvitui.Command{Std: kvitui.StdFind, Run: run}
	c.settings = &kvitui.Command{Std: kvitui.StdSettings, Run: run}
	c.about = &kvitui.Command{Std: kvitui.StdAbout, Run: run}
	c.help = &kvitui.Command{Std: kvitui.StdHelp, Run: run}
	c.quit = &kvitui.Command{Std: kvitui.StdQuit, Run: run}
	c.sidebar = &kvitui.Command{Name: "Sidebar", TurnOn: "Show sidebar", TurnOff: "Hide sidebar",
		Checked: func() bool { return c.sidebarShown }, Run: func() { c.sidebarShown = !c.sidebarShown }}
	c.importStatement = &kvitui.Command{Name: "Import statement", Asks: true,
		Shortcut: kvitui.Keys(unison.KeyI, kvitui.ModPrimary|kvitui.ModShift), Run: func() { c.imported++ }}
	return c
}

func (c *barCommands) bar(ui *kvitui.UI) *kvitui.MenuBar {
	return kvitui.NewMenuBar(ui).
		Add(c.save, c.undo, c.cut, c.find, c.settings, c.about, c.help, c.quit).
		AddTo(kvitui.MenuView, kvitui.MenuItem{Command: c.sidebar}).
		AddMenu("&Statement", kvitui.MenuItem{Command: c.importStatement})
}

// shape is a menu bar's menus and lines as text, a separator as "-".
func shape(ui *kvitui.UI, b *kvitui.MenuBar) map[string][]string {
	out := map[string][]string{}
	for _, m := range b.Menus() {
		title, _, _ := kvitui.AccessText(m.Title)
		var lines []string
		for _, it := range m.Items {
			switch {
			case it.Separator:
				lines = append(lines, "-")
			case it.Command != nil:
				lines = append(lines, it.Command.Label(ui))
			default:
				lines = append(lines, it.Text)
			}
		}
		out[title] = lines
	}
	return out
}

func titles(b *kvitui.MenuBar) []string {
	var out []string
	for _, m := range b.Menus() {
		t, _, _ := kvitui.AccessText(m.Title)
		out = append(out, t)
	}
	return out
}

// The standard commands go in the standard menus at their standard places,
// in each form, without being told where (CMD-2 to CMD-8).
func TestAMenuBarPlacesTheStandardCommands(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.AppName = "Kvit Cash"
	c := newBarCommands()
	b := c.bar(ui)
	if got, want := titles(b), []string{"File", "Edit", "View", "Statement", "Tools", "Help"}; !slices.Equal(got, want) {
		t.Errorf("the Windows form's menus are %v, want %v", got, want)
	}
	s := shape(ui, b)
	for menu, want := range map[string][]string{
		"File":  {"Save", "-", "Exit"},
		"Edit":  {"Undo", "-", "Cut", "-", "Find"},
		"View":  {"Sidebar"},
		"Tools": {"Settings"},
		"Help":  {"Kvit Cash help"},
	} {
		if !slices.Equal(s[menu], want) {
			t.Errorf("the Windows form's %s menu is %v, want %v", menu, s[menu], want)
		}
	}
	ui.SetForm(kvitui.FormMac)
	if got, want := titles(b), []string{"Kvit Cash", "File", "Edit", "View", "Statement", "Window", "Help"}; !slices.Equal(got, want) {
		t.Errorf("the macOS form's menus are %v, want %v", got, want)
	}
	s = shape(ui, b)
	if app := s["Kvit Cash"]; len(app) < 4 || app[0] != "About Kvit Cash" || app[2] != "Settings…" || app[len(app)-1] != "Quit Kvit Cash" {
		t.Errorf("the macOS form's app menu is %v", app)
	}
	if got := s["View"]; !slices.Equal(got, []string{"Hide Sidebar"}) {
		t.Errorf("the macOS form's View menu is %v, want the switch in words", got)
	}
	if got := s["File"]; !slices.Equal(got, []string{"Save"}) {
		t.Errorf("the macOS form's File menu is %v", got)
	}
}

// The Kvit-drawn bar opens a menu on a press, moves to the next menu on
// Right, and is reached from the keyboard by Alt with a title's letter, by
// F10 and by Alt on its own; and the window answers the keys of the commands
// on it.
func TestAMenuBarWorksFromThePointerAndTheKeyboard(t *testing.T) {
	c := newBarCommands()
	var field *kvitui.Field
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		field = kvitui.NewField(ui)
		field.Label = "Payee"
		return []unison.Paneler{kvitui.Left(field)}
	})
	var bar *kvitui.MenuBar
	screen.Do(func() {
		bar = c.bar(ui)
		w.SetMenuBar(bar)
	})
	screen.Sync()
	openMenu := func() string {
		for _, n := range nodes(screen, w, role.Menu) {
			return n.Name
		}
		return ""
	}
	titlesShown := nodes(screen, w, role.MenuItem)
	if len(titlesShown) != 6 {
		t.Fatalf("the bar shows %d titles", len(titlesShown))
	}
	screen.Click(titlesShown[0].Bounds.Center())
	if got := openMenu(); got != "File" {
		t.Errorf("a press on File opened %q", got)
	}
	screen.KeyPress(unison.KeyRight, 0)
	if got := openMenu(); got != "Edit" {
		t.Errorf("Right in File opened %q", got)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if got := openMenu(); got != "" {
		t.Errorf("Escape left %q open", got)
	}
	screen.Click(screen.PanelCenter(field.Edit()))
	screen.KeyPress(unison.KeyS, mod.Option)
	if got := openMenu(); got != "Statement" {
		t.Errorf("Alt+S opened %q", got)
	}
	screen.KeyPress(unison.KeyReturn, 0)
	if c.imported != 1 {
		t.Errorf("Return on the Statement menu's first line imported %d times", c.imported)
	}
	screen.KeyPress(unison.KeyF10, 0)
	screen.KeyPress(unison.KeyDown, 0)
	if got := openMenu(); got != "File" {
		t.Errorf("F10 and Down opened %q", got)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	screen.KeyPress(unison.KeyEscape, 0)
	screen.KeyDown(unison.KeyLOption, mod.Option)
	screen.KeyUp(unison.KeyLOption, 0)
	screen.KeyPress(unison.KeyDown, 0)
	if got := openMenu(); got != "File" {
		t.Errorf("Alt on its own and Down opened %q", got)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Click(screen.PanelCenter(field.Edit()))
	screen.KeyPress(unison.KeyS, mod.Control)
	screen.KeyPress(unison.KeyI, mod.Control|mod.Shift)
	if c.saved != 1 || c.imported != 2 {
		t.Errorf("Ctrl+S saved %d times and Ctrl+Shift+I imported %d times", c.saved, c.imported-1)
	}
}

// A toolbar too narrow for its commands moves the least important into a
// More menu, so every command is still on the bar or in the menu (CMD-19).
func TestAToolbarMovesWhatDoesNotFitIntoMore(t *testing.T) {
	var bar *kvitui.Toolbar
	var cmds []*kvitui.Command
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		bar = kvitui.NewToolbar(ui, "Statement")
		for _, name := range []string{"Import", "Export", "Reconcile", "Split", "Merge", "Archive"} {
			cmds = append(cmds, &kvitui.Command{Name: name, Run: func() {}})
		}
		bar.Add(kvitui.ToolLeading, cmds[:4]...)
		bar.Add(kvitui.ToolTrailing, cmds[4:]...)
		return []unison.Paneler{kvitui.FullWidth(bar)}
	})
	if got := bar.InMore(); len(got) != 0 {
		t.Errorf("a wide toolbar put %d commands in More", len(got))
	}
	screen.Do(func() {
		narrow := float32(ui.Interface.Px(200))
		bar.SetFrameRect(bar.FrameRect().Inset(geom.Insets{Right: bar.FrameRect().Width - narrow}))
		bar.ValidateLayout()
	})
	more := bar.InMore()
	if len(more) == 0 || more[len(more)-1] != cmds[len(cmds)-1] {
		t.Errorf("a narrow toolbar put %v in More, the least important last", more)
	}
	if more[0] == cmds[0] {
		t.Error("the most important command went to More")
	}
	if off := w.OffSurface(); len(off) != 0 {
		t.Errorf("commands on a toolbar were reported off every surface: %v", off)
	}
}

// A command the window runs that is on no menu bar or toolbar is reported
// (CMD-1).
func TestACommandOnNoSurfaceIsFound(t *testing.T) {
	stray := &kvitui.Command{Name: "Recalculate", Run: func() {}}
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		return []unison.Paneler{kvitui.NewCommandButton(ui, stray)}
	})
	var off []*kvitui.Command
	screen.Do(func() { off = w.OffSurface() })
	if len(off) != 1 || off[0] != stray {
		t.Errorf("OffSurface found %v", off)
	}
}

// A split button's main part runs the variant last chosen on the Windows
// form, and the likeliest on the macOS form (CMD-17).
func TestASplitButtonRunsTheVariantTheFormSays(t *testing.T) {
	var ran []string
	variant := func(name string) *kvitui.Command {
		return &kvitui.Command{Name: name, Run: func() { ran = append(ran, name) }}
	}
	var split *kvitui.SplitButton
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		split = kvitui.NewSplitButton(ui, variant("Export as CSV"), variant("Export as PDF"))
		split.MenuLabel = "Choose format"
		return []unison.Paneler{split}
	})
	pick := func(i int) {
		screen.Click(screen.PanelCenter(split.Children()[1]))
		lines := nodes(screen, w, role.MenuItem)
		screen.Click(lines[i].Bounds.Center())
	}
	pick(1)
	screen.Click(screen.PanelCenter(split.Children()[0]))
	if !slices.Equal(ran, []string{"Export as PDF", "Export as PDF"}) {
		t.Errorf("on the Windows form the main part ran %v after PDF was chosen", ran)
	}
	ran = nil
	screen.Do(func() {
		ui.SetForm(kvitui.FormMac)
		split.Current = 0
	})
	pick(1)
	screen.Click(screen.PanelCenter(split.Children()[0]))
	if !slices.Equal(ran, []string{"Export as PDF", "Export as CSV"}) {
		t.Errorf("on the macOS form the main part ran %v after PDF was chosen", ran)
	}
}

// A toggle split button's main part turns the feature on and off, and its
// menu picks the style, the one in use marked; choosing a style turns the
// feature on (CMD-17).
func TestAToggleSplitButtonTurnsOnAndPicksAStyle(t *testing.T) {
	on := false
	style := ""
	bullets := &kvitui.Command{Name: "Bullets", Symbol: "list-bullets",
		Checked: func() bool { return on }, Run: func() { on = !on }}
	styleOf := func(name string) *kvitui.Command {
		return &kvitui.Command{Name: name, Run: func() { style = name }}
	}
	var split *kvitui.SplitButton
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		split = kvitui.NewToggleSplitButton(ui, bullets, styleOf("Dot bullets"), styleOf("Dash bullets"))
		return []unison.Paneler{split}
	})
	screen.Click(screen.PanelCenter(split.Children()[0]))
	if !on {
		t.Fatal("the main part did not turn the bullets on")
	}
	if n := screen.AccessibilityNodeFor(split.Children()[0]); n == nil || n.Checked != check.On {
		t.Errorf("the main part, on, is described as %+v", n)
	}
	screen.Click(screen.PanelCenter(split.Children()[0]))
	if on {
		t.Fatal("the main part did not turn the bullets off")
	}
	screen.Click(screen.PanelCenter(split.Children()[1]))
	lines := nodes(screen, w, role.MenuItem)
	if len(lines) != 2 || lines[0].Checked != check.On || lines[1].Checked != check.Off {
		t.Fatalf("the menu's lines are %+v", lines)
	}
	screen.Click(lines[1].Bounds.Center())
	if style != "Dash bullets" || !on || split.Current != 1 {
		t.Errorf("choosing Dash bullets left the style %q, the bullets on %v and Current %d", style, on, split.Current)
	}
}
