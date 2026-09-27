// Package uitest opens Kvit components and screens headlessly for tests: in a
// Kvit window, in a chosen theme and interface size, with a screen reader's
// view of what is shown, and with the library's rules checked on a
// program's own source. The Kvit applications' tests use it in place of a
// window on a desktop, which a build machine does not have.
package uitest

import (
	"image"
	"testing"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
)

// Options say how a Session starts.
type Options struct {
	// Theme is the theme; light unless set.
	Theme string
	// InterfaceSize is the interface size in pixels; the default unless set.
	InterfaceSize int
	// Width and Height are the headless screen's size; 1500 by 1000 unless
	// set.
	Width, Height float32
	// Motion leaves animation on. It is off unless set, so a test sees the
	// state a change ends in at once rather than partway there.
	Motion bool
}

// Session is a headless screen showing one Kvit window.
type Session struct {
	t testing.TB
	// Screen is unison's headless screen: the pointer, the keyboard, and
	// what is drawn.
	Screen *unison.HeadlessScreen
	// UI is the Kvit interface the window is drawn with.
	UI *kvitui.UI
	// Window is the Kvit window whose body is what the test built.
	Window *kvitui.Window
}

// Open starts a headless screen with a Kvit window showing what build
// returns, and stops it when the test ends. What is built is placed as a page
// places its content: in a column a view margin in from the window's edges,
// taking its own height. A screen reader's view is on.
func Open(t testing.TB, opt Options, build func(ui *kvitui.UI) unison.Paneler) *Session {
	t.Helper()
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	theme := opt.Theme
	if theme == "" {
		theme = tokens.Light
	}
	ui.Theme.SetThemeID(theme)
	if opt.InterfaceSize > 0 {
		ui.Interface.SetFontSize(opt.InterfaceSize)
	}
	ui.Theme.SetReducedMotion(!opt.Motion)
	width, height := opt.Width, opt.Height
	if width <= 0 {
		width = 1500
	}
	if height <= 0 {
		height = 1000
	}
	s := &Session{t: t, UI: ui}
	var startErr error
	s.Screen, err = unison.StartHeadless(unison.HeadlessConfig{Width: width, Height: height},
		unison.StartupFinishedCallback(func() {
			if s.Window, startErr = kvitui.NewWindow(ui, t.Name()); startErr != nil {
				return
			}
			page := kvitui.Column(ui, kvitui.SizeSpace, build(ui))
			page.SetBorder(kvitui.Padding(ui, kvitui.SizeViewMargin))
			s.Window.SetBody(page)
			s.Window.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Screen.Stop)
	if startErr != nil {
		t.Fatal(startErr)
	}
	s.Screen.EnableAccessibility()
	s.Screen.Sync()
	return s
}

// Do runs f on the interface's thread and waits for it.
func (s *Session) Do(f func()) { s.Screen.Do(f) }

// Sync lets everything waiting to run, lay out or draw do so.
func (s *Session) Sync() { s.Screen.Sync() }

// SetTheme changes the theme, as a reader choosing one does.
func (s *Session) SetTheme(id string) {
	s.Do(func() { s.UI.Theme.SetThemeID(id) })
	s.Sync()
}

// SetInterfaceSize changes the interface size, as a reader choosing one
// does.
func (s *Session) SetInterfaceSize(px int) {
	s.Do(func() {
		s.UI.Interface.SetFontSize(px)
		s.Window.Content().MarkForLayoutRecursively()
		s.Window.MarkForRedraw()
	})
	s.Sync()
}

// Capture is the window as drawn now.
func (s *Session) Capture() image.Image { return s.Screen.CaptureWindow(s.Window.Window) }

// Tree is what a screen reader is told about the window now.
func (s *Session) Tree() *accessibility.Tree { return s.Screen.AccessibilityTree(s.Window.Window) }

// Node is what a screen reader is told about one panel, or nil.
func (s *Session) Node(p unison.Paneler) *accessibility.Node { return s.Screen.AccessibilityNodeFor(p) }

// Focused is the node a screen reader says holds the keyboard focus, or nil.
func (s *Session) Focused() *accessibility.Node {
	tree := s.Tree()
	if tree == nil {
		return nil
	}
	return tree.Nodes[tree.Focus]
}

// Nodes are the window's nodes of one role.
func (s *Session) Nodes(r role.Enum) []*accessibility.Node {
	var out []*accessibility.Node
	if tree := s.Tree(); tree != nil {
		for _, n := range tree.Nodes {
			if n.Role == r {
				out = append(out, n)
			}
		}
	}
	return out
}

// Named is the first node with a name, or nil.
func (s *Session) Named(name string) *accessibility.Node {
	if tree := s.Tree(); tree != nil {
		for _, n := range tree.Nodes {
			if n.Name == name {
				return n
			}
		}
	}
	return nil
}

// WaitFor runs the interface until cond holds, and fails the test after two
// seconds: long enough for a tooltip's half-second pause, and short enough
// that a wait that never ends says so.
func (s *Session) WaitFor(what string, cond func() bool) {
	s.t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		ok := false
		s.Do(func() { ok = cond() })
		if ok {
			return
		}
		time.Sleep(20 * time.Millisecond)
		s.Sync()
	}
	s.t.Fatalf("waited two seconds for %s", what)
}

// CheckNamed fails the test for every node a reader can reach or press that
// does not say what kind of thing it is or what it is called: a control with
// no name is read as "button" and nothing else. An empty table cell may have
// no name, since its description names its row.
func (s *Session) CheckNamed() {
	s.t.Helper()
	for _, problem := range Unnamed(s.Tree()) {
		s.t.Error(problem)
	}
}

// Unnamed lists the nodes of a tree that a reader can reach or press and
// that have no role or no name.
func Unnamed(tree *accessibility.Tree) []string {
	if tree == nil {
		return nil
	}
	var out []string
	for _, n := range tree.Nodes {
		reachable := n.Focusable || n.Actions.Has(accessibility.Press)
		if !reachable || n.Ignored || n.Disabled || n.Role == role.ScrollBar {
			continue
		}
		if n.Role == role.Auto || n.Role == role.None || n.Role == role.Unknown {
			out = append(out, "a node named "+quote(n.Name)+" has no role")
		}
		if n.Name == "" && len(n.LabeledBy) == 0 && n.Role != role.Cell {
			out = append(out, "a "+n.Role.String()+" node at "+n.Bounds.String()+" has no name")
		}
	}
	return out
}

func quote(s string) string { return "\"" + s + "\"" }
