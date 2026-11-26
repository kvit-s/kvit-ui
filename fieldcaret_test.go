package kvitui_test

import (
	"image"
	"sync"
	"testing"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// restAfter and blinkEvery stand in for the ten seconds a caret blinks with
// its field left alone and unison's blink of 560 ms, so the tests take
// seconds rather than minutes. The rest is long against the blink, so a
// caret that stopped at once is told from one that waited.
const (
	restAfter  = 500 * time.Millisecond
	blinkEvery = 25 * time.Millisecond
)

// caretWatch keeps the time of every frame a field is drawn in. unison
// redraws a field for every blink of its caret, so a field drawn in no frame
// for many blinks has a caret that is not blinking.
type caretWatch struct {
	lock   sync.Mutex
	frames []time.Time
}

// watchCaret makes a field blink every blinkEvery and keeps its frames. It
// is called as the field is made, before it is first drawn: a field drawn
// with the keyboard has a blink waiting at the rate it had then.
func watchCaret(e *unison.Field) *caretWatch {
	w := &caretWatch{}
	e.BlinkRate = blinkEvery
	draw := e.DrawCallback
	e.DrawCallback = func(gc *unison.Canvas, dirty geom.Rect) {
		w.lock.Lock()
		w.frames = append(w.frames, time.Now())
		w.lock.Unlock()
		draw(gc, dirty)
	}
	return w
}

// drawnBetween is how many frames the field was drawn in from one time up to
// another, and when the last of them was.
func (w *caretWatch) drawnBetween(from, to time.Time) (n int, last time.Time) {
	w.lock.Lock()
	defer w.lock.Unlock()
	for _, at := range w.frames {
		if !at.Before(from) && at.Before(to) {
			n++
			last = at
		}
	}
	return n, last
}

// drawnWithin is how many frames the field is drawn in over the next d.
func (w *caretWatch) drawnWithin(d time.Duration) int {
	from := time.Now()
	time.Sleep(d)
	n, _ := w.drawnBetween(from, time.Now())
	return n
}

// rests gives an input, then fails the test unless the caret blinked until
// restAfter had passed and then stopped, with nothing more drawn for a dozen
// blinks. The input's own frames are few; the blinks over restAfter are
// many. Whether it blinked is judged from the frames' times rather than by
// watching after the input returns, since a headless screen's input waits
// for the window to have nothing left to draw, which a blinking caret puts
// off until it rests.
func (w *caretWatch) rests(t *testing.T, what string, input func()) {
	t.Helper()
	since := time.Now()
	input()
	time.Sleep(time.Until(since.Add(restAfter)))
	for w.drawnWithin(5*blinkEvery) != 0 {
		if time.Since(since) > restAfter+5*time.Second {
			t.Fatalf("%s: the caret was still blinking %v after the input", what, time.Since(since).Round(time.Millisecond))
		}
	}
	if n, _ := w.drawnBetween(since, since.Add(restAfter)); n < int(restAfter/blinkEvery)/3 {
		t.Errorf("%s: the field was drawn %d times in the %v after the input, so its caret was not blinking", what, n, restAfter)
	}
	if _, last := w.drawnBetween(since, time.Now()); last.Before(since.Add(restAfter)) {
		t.Errorf("%s: the caret stopped blinking %v after the input, before the %v it should blink for",
			what, last.Sub(since).Round(time.Millisecond), restAfter)
	}
	if n := w.drawnWithin(12 * blinkEvery); n != 0 {
		t.Errorf("%s: with the caret at rest the field was drawn %d more times in %v", what, n, 12*blinkEvery)
	}
}

// A field's caret blinks while the field is in use and stops, shown, once
// nothing has been typed or pressed and the caret has not moved for a while,
// since every blink redraws the whole window. A key that moves nothing, a
// press, typing and a caret move made by the application each set it blinking
// again, and a rest changes neither the text nor the selection.
func TestAFieldCaretStopsBlinkingWhenLeftAlone(t *testing.T) {
	t.Cleanup(kvitui.SetCaretRestAfter(restAfter))
	var f *kvitui.Field
	var w *caretWatch
	changes := 0
	s := uitest.Open(t, uitest.Options{Width: 500, Height: 200}, func(ui *kvitui.UI) unison.Paneler {
		f = kvitui.NewField(ui)
		f.Label = "Name"
		f.OnChange = func(string) { changes++ }
		w = watchCaret(f.Edit())
		return kvitui.Width(ui, kvitui.Px(300), f)
	})
	screen, e := s.Screen, f.Edit()
	state := func() (st *unison.FieldState) {
		s.Do(func() { st = e.GetFieldState() })
		return st
	}

	w.rests(t, "pressed in an empty field", func() { screen.Click(screen.PanelCenter(e)) })
	caretShown(t, s, e)

	w.rests(t, "Left in an empty field", func() { screen.KeyPress(unison.KeyLeft, 0) })
	if st := state(); st.Text != "" || st.SelectionEnd != 0 {
		t.Errorf("resting changed an empty field to %q with the caret at %d", st.Text, st.SelectionEnd)
	}

	w.rests(t, "typing", func() { screen.Type("abc") })
	if st := state(); st.Text != "abc" || st.SelectionStart != 3 || st.SelectionEnd != 3 || changes != 3 {
		t.Errorf("after typing and resting the field holds %q, selection %d-%d, with %d changes reported, want \"abc\", 3-3 and 3",
			st.Text, st.SelectionStart, st.SelectionEnd, changes)
	}

	w.rests(t, "the caret moved by the application", func() { s.Do(e.SetSelectionToStart) })
	if st := state(); st.SelectionStart != 0 || st.SelectionEnd != 0 {
		t.Errorf("resting moved the caret from the start to %d-%d", st.SelectionStart, st.SelectionEnd)
	}

	// A press where the caret already is moves nothing.
	w.rests(t, "a press at the caret", func() { screen.Click(screen.PanelPoint(e, geom.NewPoint(0, 1))) })
	if st := state(); st.SelectionStart != 0 || st.SelectionEnd != 0 {
		t.Errorf("the press at the caret moved it to %d-%d", st.SelectionStart, st.SelectionEnd)
	}

	// A selection shows no caret, but unison blinks it all the same, and
	// redraws for every blink, until the field is left alone.
	w.rests(t, "a selection", func() { screen.KeyPress(unison.KeyRight, mod.Shift) })
	if st := state(); st.SelectionStart != 0 || st.SelectionEnd != 1 || st.SelectionAnchor != 0 || st.Text != "abc" || changes != 3 {
		t.Errorf("resting left %q with selection %d-%d from %d and %d changes reported, want \"abc\", 0-1 from 0 and 3",
			st.Text, st.SelectionStart, st.SelectionEnd, st.SelectionAnchor, changes)
	}
}

// caretShown fails the test unless an empty field's caret is drawn: a mark
// in the field's text box at the caret's place, and nothing else there but
// the field's ground.
func caretShown(t *testing.T, s *uitest.Session, e *unison.Field) {
	t.Helper()
	var box geom.Rect
	var at geom.Point
	s.Do(func() {
		box = e.RectToRoot(e.ContentRect(false))
		at = e.PointToRoot(e.FromSelectionIndex(0))
	})
	img := s.Capture()
	ground := img.At(int(box.Right())-2, int(box.CenterY()))
	marked, stray := 0, 0
	for y := int(box.Y); y < int(box.Bottom()); y++ {
		for x := int(box.X); x < int(box.Right()); x++ {
			if samePixel(img.At(x, y), ground) {
				continue
			}
			if x >= int(at.X)-2 && x <= int(at.X)+1 {
				marked++
			} else {
				stray++
			}
		}
	}
	if marked < int(box.Height)/2 || stray != 0 {
		t.Errorf("at rest the caret is not drawn: %d pixels marked at its place, %d elsewhere in the field's text box %v",
			marked, stray, image.Rect(int(box.X), int(box.Y), int(box.Right()), int(box.Bottom())))
	}
}

// A caret coming to rest does not bring itself into view: the source of a
// file, pressed in and then scrolled away from, stays where it was scrolled
// to.
func TestARestingCaretLeavesTheViewWhereItWasScrolled(t *testing.T) {
	t.Cleanup(kvitui.SetCaretRestAfter(restAfter))
	var region *kvitui.Region
	var area *kvitui.TextArea
	var w *caretWatch
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		area = kvitui.NewTextArea(ui)
		area.Label, area.ReadOnly, area.Plain, area.Fit = "Source", true, true, kvitui.FitText
		lines := ""
		for range 200 {
			lines += "a line of the file\n"
		}
		area.SetText(lines)
		w = watchCaret(area.Edit())
		region = kvitui.NewRegion(ui, area)
		region.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(480, 200)})
		return []unison.Paneler{region}
	})
	e := area.Edit()
	// The area has had the keyboard since the window opened, with its text
	// selected and in view to its end. The press is on its first line, and
	// then the region is scrolled to the middle of the text.
	screen.Do(func() { region.ScrollTo(0) })
	var top float32
	w.rests(t, "pressed in a text area", func() {
		screen.Click(screen.PanelPoint(e, geom.NewPoint(4, 4)))
		screen.Do(func() {
			region.ScrollTo(region.Children()[1].Self.(*unison.ScrollBar).MaxValue() / 2)
			_, top = region.Position()
		})
	})
	if top == 0 {
		t.Fatal("the region did not scroll")
	}
	var y float32
	var start, end int
	screen.Do(func() {
		_, y = region.Position()
		start, end = e.Selection()
	})
	if y != top || start != 0 || end != 0 {
		t.Errorf("at rest the region is scrolled to %.1f with the caret at %d-%d, want %.1f with the caret at the start", y, start, end, top)
	}
}

// A field whose caret came to rest blinks again when the keyboard comes back
// to it, by Tab or the application, with nothing typed or pressed in it.
func TestARestingCaretBlinksAgainWhenTheKeyboardReturns(t *testing.T) {
	t.Cleanup(kvitui.SetCaretRestAfter(restAfter))
	var name, place *kvitui.Field
	var w *caretWatch
	s := uitest.Open(t, uitest.Options{Width: 500, Height: 200}, func(ui *kvitui.UI) unison.Paneler {
		name, place = kvitui.NewField(ui), kvitui.NewField(ui)
		name.Label, place.Label = "Name", "Place"
		w = watchCaret(name.Edit())
		place.Edit().BlinkRate = blinkEvery
		return kvitui.Width(ui, kvitui.Px(300), kvitui.Column(ui, kvitui.SizeSpace, name, place))
	})
	s.Do(place.Focus)
	w.rests(t, "the keyboard moved to the first field", func() { s.Do(name.Focus) })
	s.Do(place.Focus)
	w.rests(t, "the keyboard back in the first field", func() { s.Do(name.Focus) })
}
