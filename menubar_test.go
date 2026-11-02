package kvitui_test

import (
	"image"
	"slices"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
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

// A split button is drawn as one control: its arrow part sits against the
// main part, as narrow as a symbol button may be, and its outline is the
// main part's, which a toggle split button turns to the accent while the
// feature is on. The arrow is drawn before the main part, and follows a
// change at once rather than a frame later.
func TestASplitButtonIsDrawnAsOneControl(t *testing.T) {
	on := false
	bullets := &kvitui.Command{Name: "Bullets", Checked: func() bool { return on }, Run: func() { on = !on }}
	var split *kvitui.SplitButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		split = kvitui.NewToggleSplitButton(ui, bullets,
			&kvitui.Command{Name: "Dot bullets", Run: func() {}},
			&kvitui.Command{Name: "Dash bullets", Run: func() {}})
		return []unison.Paneler{split}
	})
	var main, arrow geom.Rect
	screen.Do(func() {
		main = split.Children()[0].RectToRoot(split.Children()[0].ContentRect(true))
		arrow = split.Children()[1].RectToRoot(split.Children()[1].ContentRect(true))
	})
	if arrow.X != main.Right() || arrow.Width != float32(ui.Interface.Px(20)) {
		t.Fatalf("the main part is %v and the arrow part %v", main, arrow)
	}
	edge := func() unison.Color {
		pr, pg, pb, _ := screen.Capture().At(int(arrow.X+arrow.Width/2), int(arrow.Y)).RGBA()
		return unison.RGB(int(pr>>8), int(pg>>8), int(pb>>8))
	}
	if got, want := edge(), kvitui.Color(ui.Theme.Tokens().BorderStrong); got != want {
		t.Errorf("off, the arrow part's top edge is %v, want the strong border %v", got, want)
	}
	screen.Do(func() {
		on = true
		split.MarkForRedraw()
	})
	screen.Sync()
	if got, want := edge(), kvitui.Color(ui.Theme.Tokens().Accent); got != want {
		t.Errorf("on, the arrow part's top edge is %v, want the accent %v", got, want)
	}
}

// A choice button's button runs the action with the value chosen, and its
// choice changes the value without running it. The pair keeps its width
// whatever is chosen, and a screen reader reads the button and the choice
// with its value.
func TestAChoiceButtonRunsWithTheValueChosen(t *testing.T) {
	var pressed, chosen []string
	var measure *kvitui.ChoiceButton
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		measure = kvitui.NewChoiceButton(ui, "Measure", "How many",
			kvitui.Option{Value: "all", Label: "All"}, kvitui.Option{Value: "5", Label: "5"},
			kvitui.Option{Value: "10", Label: "10"})
		measure.OnPress = func(v string) { pressed = append(pressed, v) }
		measure.OnChoose = func(v string) { chosen = append(chosen, v) }
		return []unison.Paneler{measure}
	})
	button, choice := measure.Children()[0], measure.Children()[1]
	var before geom.Rect
	screen.Do(func() { before = measure.FrameRect() })
	screen.Click(screen.PanelCenter(button))
	if !slices.Equal(pressed, []string{"all"}) {
		t.Fatalf("the button ran with %v", pressed)
	}
	screen.Click(screen.PanelCenter(choice))
	lines := nodes(screen, w, role.MenuItem)
	if len(lines) != 3 || lines[0].Name != "All" || lines[0].Checked != check.On || lines[2].Checked != check.Off {
		t.Fatalf("the choice's lines are %+v", lines)
	}
	screen.Click(lines[2].Bounds.Center())
	if measure.Current != "10" || !slices.Equal(chosen, []string{"10"}) || len(pressed) != 1 {
		t.Errorf("choosing 10 left %q, reported %v and ran %v", measure.Current, chosen, pressed)
	}
	screen.Do(func() {
		if after := measure.FrameRect(); after != before {
			t.Errorf("choosing 10 moved the pair from %v to %v", before, after)
		}
	})
	screen.Click(screen.PanelCenter(button))
	if !slices.Equal(pressed, []string{"all", "10"}) {
		t.Errorf("after 10 was chosen the button ran with %v", pressed)
	}
	// Up changes the value without opening the list.
	screen.Do(func() { choice.RequestFocus() })
	screen.KeyPress(unison.KeyUp, 0)
	if measure.Current != "5" || len(nodes(screen, w, role.MenuItem)) != 0 {
		t.Errorf("Up on the choice left %q", measure.Current)
	}
	if n := screen.AccessibilityNodeFor(measure); n == nil || n.Role != role.Group || n.Name != "Measure" {
		t.Errorf("the pair's node: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(button); n == nil || n.Role != role.Button || n.Name != "Measure" {
		t.Errorf("the button's node: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(choice); n == nil || n.Role != role.ComboBox || n.Name != "How many" || n.Value != "5" {
		t.Errorf("the choice's node: %+v", n)
	}
}

// A choice button is enabled and disabled as one: by its command when it has
// one, and by SetEnabled when it has none.
func TestAChoiceButtonIsDisabledAsOne(t *testing.T) {
	ready := false
	var byCommand, bySetting *kvitui.ChoiceButton
	pressed := 0
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		values := []kvitui.Option{{Value: "all", Label: "All"}, {Value: "5", Label: "5"}}
		byCommand = kvitui.NewChoiceButton(ui, "", "How many", values...)
		byCommand.Command = &kvitui.Command{Name: "Measure", Enabled: func() bool { return ready }, Run: func() { pressed++ }}
		bySetting = kvitui.NewChoiceButton(ui, "Export", "How many", values...)
		bySetting.SetEnabled(false)
		return []unison.Paneler{byCommand, bySetting}
	})
	enabled := func(c *kvitui.ChoiceButton) (button, choice bool) {
		screen.Do(func() { button, choice = c.Children()[0].Enabled(), c.Children()[1].Enabled() })
		return button, choice
	}
	if b, c := enabled(byCommand); b || c {
		t.Errorf("with its command unable to act, the button is enabled %v and the choice %v", b, c)
	}
	if b, c := enabled(bySetting); b || c {
		t.Errorf("disabled, the button is enabled %v and the choice %v", b, c)
	}
	screen.Do(func() {
		ready = true
		byCommand.MarkForRedraw()
		bySetting.SetEnabled(true)
	})
	screen.Sync()
	if b, c := enabled(byCommand); !b || !c {
		t.Errorf("with its command able to act, the button is enabled %v and the choice %v", b, c)
	}
	if b, c := enabled(bySetting); !b || !c {
		t.Errorf("enabled again, the button is enabled %v and the choice %v", b, c)
	}
	screen.Click(screen.PanelCenter(byCommand.Children()[0]))
	if pressed != 1 {
		t.Errorf("pressing the button ran its command %d times", pressed)
	}
}

// On a filled pair, a choice button or a split button, the line between the
// parts is the words' colour at a fifth of its strength, so it parts the two
// without cutting across the fill; the high-contrast theme keeps it at full
// strength. A toggle split button given the filled form stays ordinary,
// since a fill cannot show whether its feature is on.
func TestAFilledPairIsPartedByASoftLine(t *testing.T) {
	var measure *kvitui.ChoiceButton
	var export, bullets *kvitui.SplitButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		measure = kvitui.NewChoiceButton(ui, "Measure", "How many",
			kvitui.Option{Value: "all", Label: "All"}, kvitui.Option{Value: "5", Label: "5"})
		measure.Form = kvitui.ButtonPrimary
		run := func() {}
		export = kvitui.NewSplitButton(ui, &kvitui.Command{Name: "Export as CSV", Run: run},
			&kvitui.Command{Name: "Export as PDF", Run: run})
		export.Form = kvitui.ButtonPrimary
		bullets = kvitui.NewToggleSplitButton(ui, &kvitui.Command{Name: "Bullets", Checked: func() bool { return false }, Run: run},
			&kvitui.Command{Name: "Dot bullets", Run: run})
		bullets.Form = kvitui.ButtonPrimary
		return []unison.Paneler{measure, export, bullets}
	})
	// at is the colour at a point of the screen.
	at := func(img image.Image, x, y int) [3]int {
		pr, pg, pb, _ := img.At(x, y).RGBA()
		return [3]int{int(pr >> 8), int(pg >> 8), int(pb >> 8)}
	}
	// line is the colour across a pair's meeting that differs most from the
	// fill, the colour inside the second part, and the fill and the words'
	// colour.
	line := func(pair unison.Paneler) (got, inside, fill, ink [3]int) {
		var r geom.Rect
		screen.Do(func() {
			part := pair.AsPanel().Children()[1]
			r = part.RectToRoot(part.ContentRect(true))
		})
		img := screen.Capture()
		fr, fg, fb := ui.Theme.Tokens().Accent.RGBA8()
		ir, ig, ib := tokens.LabelOn(ui.Theme.Tokens().Accent).RGBA8()
		fill, ink = [3]int{int(fr), int(fg), int(fb)}, [3]int{int(ir), int(ig), int(ib)}
		far := -1
		for x := int(r.X) - 1; x <= int(r.X)+1; x++ {
			if c := at(img, x, int(r.Y+r.Height/2)); apart(c, fill) > far {
				got, far = c, apart(c, fill)
			}
		}
		// Inside, past the down-arrow and short of the rounded corner.
		return got, at(img, int(r.Right())-3, int(r.Y+r.Height/2)), fill, ink
	}
	for _, pair := range []unison.Paneler{measure, export} {
		got, inside, fill, ink := line(pair)
		if apart(inside, fill) > 2 {
			t.Errorf("%T: the second part is %v inside, want the fill %v", pair, inside, fill)
		}
		if apart(got, fill) < 16 || apart(got, ink) < 16 {
			t.Errorf("%T: the line is %v on the fill %v with words in %v; want a mix of the two", pair, got, fill, ink)
		}
	}
	if _, inside, fill, _ := line(bullets); apart(inside, fill) <= 2 {
		t.Errorf("the toggle split button given the filled form is filled, %v", inside)
	}
	screen.Do(func() { ui.Theme.SetThemeID(tokens.HighContrast) })
	screen.Sync()
	for _, pair := range []unison.Paneler{measure, export} {
		if got, _, _, ink := line(pair); apart(got, ink) > 16 {
			t.Errorf("%T: in the high-contrast theme the line is %v, want the words' colour %v", pair, got, ink)
		}
	}
}

// apart is the largest difference between two colours' channels.
func apart(a, b [3]int) int {
	d := 0
	for i := range a {
		d = max(d, a[i]-b[i], b[i]-a[i])
	}
	return d
}

// placedAt is a layout that puts its one child at a position between
// pixels, as a parent that centres a control often does.
type placedAt struct{ x, y float32 }

func (placedAt) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	size := geom.NewSize(400, 60)
	return size, size, size
}

func (l placedAt) PerformLayout(target *unison.Panel) {
	child := target.Children()[0]
	_, p, _ := child.Sizes(geom.Size{})
	child.SetFrameRect(geom.NewRect(l.x, l.y, p.Width, p.Height))
}

// A pair put between pixels by its parent still draws its parts on whole
// pixels of the window, inside the box it was given, so its outline is one
// row of pixels and the line between the parts one column, and the line
// stops at the outline rather than poking past it.
func TestAPairBetweenPixelsIsDrawnOnWholePixels(t *testing.T) {
	var pairs []unison.Paneler
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		measure := kvitui.NewChoiceButton(ui, "Measure", "How many",
			kvitui.Option{Value: "all", Label: "All"}, kvitui.Option{Value: "5", Label: "5 days"})
		run := func() {}
		export := kvitui.NewSplitButton(ui, &kvitui.Command{Name: "Export as CSV", Run: run},
			&kvitui.Command{Name: "Export as PDF", Run: run})
		pairs = []unison.Paneler{measure, export}
		var out []unison.Paneler
		for _, p := range pairs {
			holder := unison.NewPanel()
			holder.AddChild(p)
			holder.SetLayout(placedAt{4.3, 10.5})
			out = append(out, holder)
		}
		return out
	})
	whole := func(v float32) bool { return v == float32(int(v)) }
	img := screen.Capture()
	at := func(x, y int) [3]int {
		pr, pg, pb, _ := img.At(x, y).RGBA()
		return [3]int{int(pr >> 8), int(pg >> 8), int(pb >> 8)}
	}
	er, eg, eb := ui.Theme.Tokens().BorderStrong.RGBA8()
	edge := [3]int{int(er), int(eg), int(eb)}
	for _, pair := range pairs {
		var box, first, second geom.Rect
		screen.Do(func() {
			kids := pair.AsPanel().Children()
			box = pair.AsPanel().RectToRoot(pair.AsPanel().ContentRect(true))
			first = kids[0].RectToRoot(kids[0].ContentRect(true))
			second = kids[1].RectToRoot(kids[1].ContentRect(true))
		})
		for _, r := range []geom.Rect{first, second} {
			if !whole(r.X) || !whole(r.Y) || !whole(r.Right()) || !whole(r.Bottom()) || !box.Contains(r) {
				t.Errorf("%T in %v has a part at %v, off whole pixels or outside the pair", pair, box, r)
			}
		}
		top, mid := int(first.Y), int(first.X+first.Width/2)
		if got := at(mid, top); apart(got, edge) > 2 {
			t.Errorf("%T: the outline's top row is %v, want the whole edge colour %v", pair, got, edge)
		}
		if got, ground := at(int(second.X), top-1), at(mid, top-1); apart(got, ground) > 2 {
			t.Errorf("%T: above the outline the line's column is %v where the ground is %v", pair, got, ground)
		}
	}
}

// A value shorter than the widest sits in the middle of the room the widest
// takes, as a button's words sit in its box, rather than against the line
// between the parts with a gap before the down-arrow.
func TestAShortValueIsCentredInTheChoice(t *testing.T) {
	var measure *kvitui.ChoiceButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		measure = kvitui.NewChoiceButton(ui, "Measure", "How many",
			kvitui.Option{Value: "all", Label: "All"}, kvitui.Option{Value: "5", Label: "5 days"})
		return []unison.Paneler{measure}
	})
	var r geom.Rect
	screen.Do(func() { r = measure.Children()[1].RectToRoot(measure.Children()[1].ContentRect(true)) })
	m := ui.Interface
	left := r.X + float32(m.Space())
	right := r.Right() - float32(2*m.SpaceSnug()+m.Px(10))
	img := screen.Capture()
	pr, pg, pb, _ := img.At(int(r.X)+2, int(r.Y)+3).RGBA()
	ground := [3]int{int(pr >> 8), int(pg >> 8), int(pb >> 8)}
	// The words' ink: the columns of the room with anything drawn on them.
	inkLeft, inkRight := -1, -1
	for x := int(r.X) + 1; x < int(right); x++ {
		for y := int(r.Y) + 3; y < int(r.Bottom())-3; y++ {
			cr, cg, cb, _ := img.At(x, y).RGBA()
			if apart([3]int{int(cr >> 8), int(cg >> 8), int(cb >> 8)}, ground) > 40 {
				if inkLeft < 0 {
					inkLeft = x
				}
				inkRight = x
				break
			}
		}
	}
	if inkLeft < 0 {
		t.Fatal("no words drawn in the choice")
	}
	before, after := float32(inkLeft)-left, right-float32(inkRight+1)
	if d := before - after; d > 2 || d < -2 {
		t.Errorf("All has %.1f pixels of room before it and %.1f after it", before, after)
	}
}
