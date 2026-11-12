package kvitui_test

import (
	"testing"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// tipWords are the words of the tooltips shown in a window.
func tipWords(screen *unison.HeadlessScreen, w *kvitui.Window) (words []string) {
	screen.Do(func() {
		for _, p := range w.Popups() {
			words = append(words, p.Panel.AsPanel().Accessibility.Name)
		}
	})
	return words
}

// settle lets a pause pass and everything waiting on it run.
func settle(screen *unison.HeadlessScreen, d time.Duration) {
	time.Sleep(d)
	screen.Sync()
}

func TestATooltipOnAnyPanelFollowsThePointer(t *testing.T) {
	var owner, silent *unison.Panel
	var tip, none *kvitui.Tooltip
	entered, left := 0, 0
	words := "Three files changed"
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		owner, silent = box(120, 30), box(120, 30)
		owner.MouseEnterCallback = func(geom.Point, mod.Modifiers) bool { entered++; return false }
		owner.MouseExitCallback = func() bool { left++; return false }
		tip = kvitui.AddTooltip(ui, owner, func() string { return words })
		none = kvitui.AddTooltip(ui, silent, func() string { return "" })
		return []unison.Paneler{kvitui.Left(owner), kvitui.Left(silent)}
	})
	shown := func(tip *kvitui.Tooltip) (s bool) {
		screen.Do(func() { s = tip.Shown() })
		return s
	}
	away := geom.NewPoint(900, 700)
	screen.MouseMove(screen.PanelCenter(owner), 0)
	if shown(tip) || entered != 1 {
		t.Fatalf("as the pointer arrived the tooltip is shown %v and the panel's own callback ran %d times", shown(tip), entered)
	}
	settle(screen, 300*time.Millisecond)
	if shown(tip) {
		t.Error("the tooltip showed before half a second")
	}
	waitFor(t, screen, "the tooltip", func() bool { return tip.Shown() })
	if got := tipWords(screen, w); len(got) != 1 || got[0] != words || tip.Words() != words {
		t.Errorf("the tooltip says %v and answers %q", got, tip.Words())
	}
	screen.MouseMove(away, 0)
	if shown(tip) || len(tipWords(screen, w)) != 0 || left != 1 {
		t.Errorf("after the pointer left the tooltip is shown %v and the panel's own callback ran %d times", shown(tip), left)
	}
	// A pointer that leaves within the half second shows nothing.
	screen.MouseMove(screen.PanelCenter(owner), 0)
	settle(screen, 200*time.Millisecond)
	screen.MouseMove(away, 0)
	settle(screen, 600*time.Millisecond)
	if shown(tip) {
		t.Error("a pointer that left before half a second showed the tooltip")
	}
	// Empty words show nothing.
	screen.MouseMove(screen.PanelCenter(silent), 0)
	settle(screen, 700*time.Millisecond)
	if shown(none) || len(tipWords(screen, w)) != 0 {
		t.Error("a tooltip with no words showed")
	}
	// Hide takes it away until the pointer comes back.
	screen.MouseMove(screen.PanelCenter(owner), 0)
	waitFor(t, screen, "the tooltip again", func() bool { return tip.Shown() })
	screen.Do(func() { tip.Hide() })
	settle(screen, 700*time.Millisecond)
	if shown(tip) {
		t.Error("a hidden tooltip came back with the pointer still there")
	}
}

// A tooltip that stays is shown for as long as the keyboard is there, and
// through a press; one that does not goes after the time it takes to read,
// and on a press.
func TestATooltipStaysOrGoes(t *testing.T) {
	var staying, going *unison.Panel
	var stays, goes *kvitui.Tooltip
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		staying, going = box(120, 30), box(120, 30)
		stays = kvitui.AddTooltip(ui, staying, func() string { return "Stays" })
		stays.Stays, stays.Focus, stays.FocusAtOnce = true, kvitui.TipFocusTold, true
		goes = kvitui.AddTooltip(ui, going, func() string { return "Goes" })
		goes.Focus, goes.FocusAtOnce = kvitui.TipFocusTold, true
		return []unison.Paneler{kvitui.Left(staying), kvitui.Left(going)}
	})
	screen.Do(func() {
		stays.SetFocused(true)
		goes.SetFocused(true)
	})
	screen.Sync()
	if got := tipWords(screen, w); len(got) != 2 {
		t.Fatalf("two tooltips told the keyboard is there, at once, show %v", got)
	}
	settle(screen, 2500*time.Millisecond)
	if got := tipWords(screen, w); len(got) != 2 {
		t.Errorf("after two and a half seconds the tooltips shown are %v", got)
	}
	settle(screen, 1600*time.Millisecond)
	if got := tipWords(screen, w); len(got) != 1 || got[0] != "Stays" {
		t.Errorf("after four seconds the tooltips shown are %v, want the one that stays", got)
	}
	screen.Do(func() {
		stays.SetFocused(false)
		goes.SetFocused(false)
	})
	screen.Sync()
	if got := tipWords(screen, w); len(got) != 0 {
		t.Errorf("with the keyboard gone and the pointer elsewhere the tooltips shown are %v", got)
	}

	// A press on the panel.
	for _, c := range []struct {
		owner *unison.Panel
		tip   *kvitui.Tooltip
		after bool
	}{{staying, stays, true}, {going, goes, false}} {
		at := screen.PanelCenter(c.owner)
		screen.MouseMove(at, 0)
		waitFor(t, screen, "the tooltip under the pointer", func() bool { return c.tip.Shown() })
		screen.Click(at)
		settle(screen, 100*time.Millisecond)
		var shown bool
		screen.Do(func() { shown = c.tip.Shown() })
		if shown != c.after {
			t.Errorf("after a press the tooltip with Stays %v is shown %v", c.tip.Stays, shown)
		}
		screen.MouseMove(geom.NewPoint(900, 700), 0)
	}
}

func TestHowATooltipLearnsOfTheKeyboard(t *testing.T) {
	var before, never, watched *kvitui.Check
	var told, asked *unison.Panel
	var tipNever, tipWatched, tipTold, tipAsked *kvitui.Tooltip
	questions, fromKeyboard := 0, false
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		before = kvitui.NewCheck(ui, "Before")
		never = kvitui.NewCheck(ui, "Never")
		tipNever = kvitui.AddTooltip(ui, never, func() string { return "Never" })
		tipNever.Focus = kvitui.TipFocusNone
		watched = kvitui.NewCheck(ui, "Watched")
		tipWatched = kvitui.AddTooltip(ui, watched, func() string { return "Watched" })
		told = box(120, 30)
		tipTold = kvitui.AddTooltip(ui, told, func() string { return "Told" })
		tipTold.Focus = kvitui.TipFocusTold
		asked = box(120, 30)
		asked.SetFocusable(true)
		tipAsked = kvitui.AddTooltip(ui, asked, func() string { return "Asked" })
		tipAsked.FocusAtOnce = true
		tipAsked.Keyboard = func() bool { questions++; return fromKeyboard }
		return []unison.Paneler{kvitui.Left(before), kvitui.Left(never), kvitui.Left(watched), kvitui.Left(told), kvitui.Left(asked)}
	})
	shown := func(tip *kvitui.Tooltip) (s bool) {
		screen.Do(func() { s = tip.Shown() })
		return s
	}
	// Never: Tab onto it shows nothing.
	screen.Do(func() { before.Focus() })
	screen.KeyPress(unison.KeyTab, 0)
	settle(screen, 700*time.Millisecond)
	if shown(tipNever) || len(tipWords(screen, w)) != 0 {
		t.Error("a tooltip that ignores the keyboard showed for Tab")
	}
	// Watched: Tab onto it shows it after the pause; a press does not.
	screen.KeyPress(unison.KeyTab, 0)
	if shown(tipWatched) {
		t.Error("a watched tooltip showed for Tab without the pause")
	}
	waitFor(t, screen, "the watched tooltip", func() bool { return tipWatched.Shown() })
	screen.Do(func() { before.Focus() })
	screen.Sync()
	if shown(tipWatched) {
		t.Error("the watched tooltip stayed after the keyboard left")
	}
	screen.Click(screen.PanelCenter(watched))
	screen.MouseMove(geom.NewPoint(900, 700), 0)
	settle(screen, 700*time.Millisecond)
	if shown(tipWatched) {
		t.Error("a press that gave the focus showed the watched tooltip")
	}
	// Told: after the pause, or at once; gone when told the keyboard left.
	screen.Do(func() { tipTold.SetFocused(true) })
	screen.Sync()
	if shown(tipTold) {
		t.Error("a told tooltip showed without the pause")
	}
	waitFor(t, screen, "the told tooltip", func() bool { return tipTold.Shown() })
	screen.Do(func() { tipTold.SetFocused(false) })
	if shown(tipTold) {
		t.Error("the told tooltip stayed after the keyboard left")
	}
	screen.Do(func() {
		tipTold.FocusAtOnce = true
		tipTold.SetFocused(true)
	})
	screen.Sync()
	if !shown(tipTold) {
		t.Error("a told tooltip that shows at once did not")
	}
	screen.Do(func() { tipTold.SetFocused(false) })
	// A Keyboard function is asked, and believed.
	screen.Do(func() { asked.RequestFocus() })
	screen.Sync()
	if questions != 1 || shown(tipAsked) {
		t.Errorf("focus the function says is not the keyboard's: asked %d times, shown %v", questions, shown(tipAsked))
	}
	fromKeyboard = true
	screen.Do(func() { before.Focus() })
	screen.Sync()
	screen.Do(func() { asked.RequestFocus() })
	screen.Sync()
	if questions != 2 || !shown(tipAsked) {
		t.Errorf("focus the function says is the keyboard's: asked %d times, shown %v", questions, shown(tipAsked))
	}
}

func TestAControlsTooltipCanBeSet(t *testing.T) {
	var check *kvitui.Check
	var toggle *kvitui.Switch
	var button, added *kvitui.IconButton
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		check = kvitui.NewCheck(ui, "Wrap lines")
		check.SetTooltip(func() string { return "Wraps long lines" })
		toggle = kvitui.NewSwitch(ui, "Sync")
		toggle.SetTooltip(func() string { return "Keeps the folder in step" })
		button = kvitui.NewIconButton(ui, "pencil", "Edit")
		button.SetTooltip(func() string { return "Rename the track" })
		added = kvitui.NewIconButton(ui, "trash", "Delete")
		kvitui.AddTooltip(ui, added, func() string { return "Delete the track for good" })
		return []unison.Paneler{kvitui.Left(check), kvitui.Left(toggle), kvitui.Left(button), kvitui.Left(added)}
	})
	away := geom.NewPoint(900, 700)
	hover := func(p unison.Paneler, want string) {
		t.Helper()
		screen.MouseMove(screen.PanelCenter(p), 0)
		waitFor(t, screen, "the tooltip saying "+want, func() bool { return len(w.Popups()) == 1 })
		if got := tipWords(screen, w); got[0] != want {
			t.Errorf("the tooltip says %q, want %q", got[0], want)
		}
		screen.MouseMove(away, 0)
		screen.Sync()
	}
	hover(check, "Wraps long lines")
	hover(toggle, "Keeps the folder in step")
	hover(button, "Rename the track")
	if got := button.TooltipWords(); got != "Rename the track" {
		t.Errorf("TooltipWords answers %q for a button given other words", got)
	}
	// Nil gives the button back its own.
	screen.Do(func() { button.SetTooltip(nil) })
	hover(button, "Edit")
	if got := button.TooltipWords(); got != "Edit" {
		t.Errorf("TooltipWords answers %q for a button given its own back", got)
	}
	// A tooltip added to a control is the only one it shows.
	screen.MouseMove(screen.PanelCenter(added), 0)
	waitFor(t, screen, "the added tooltip", func() bool { return len(w.Popups()) > 0 })
	settle(screen, 700*time.Millisecond)
	if got := tipWords(screen, w); len(got) != 1 || got[0] != "Delete the track for good" {
		t.Errorf("a button with an added tooltip shows %v", got)
	}
}
