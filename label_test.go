package kvitui_test

import (
	"image"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

const longSentence = "The quick brown fox jumps over the lazy dog, and then goes round the field once more to do it again, " +
	"because one jump was never going to be enough for a fox with this much energy and a dog this lazy, " +
	"and the afternoon is long, and nobody has anywhere else to be until the light goes."

func TestALabelCutShortInItsMiddleKeepsBothEnds(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		l := kvitui.NewLabel(ui, "conversation-transcript-2026-11-11-final.md")
		whole, _ := l.Layout(0).Size()
		end := l.Layout(whole / 2).Text()
		if !strings.HasPrefix(end, "conv") || !strings.HasSuffix(end, "…") {
			t.Errorf("a label cut short at its end reads %q", end)
		}
		l.Middle = true
		middle := l.Layout(whole / 2).Text()
		if !strings.HasPrefix(middle, "conv") || !strings.HasSuffix(middle, ".md") || !strings.Contains(middle, "…") {
			t.Errorf("a label cut short in its middle reads %q", middle)
		}
		// Wrapped text is not cut short at all.
		l.Wrap = true
		if wrapped := l.Layout(whole / 2).Text(); strings.Contains(wrapped, "…") {
			t.Errorf("a wrapped label with Middle set reads %q", wrapped)
		}
	})
}

func TestALabelAlignsItsLines(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		for _, wrap := range []bool{false, true} {
			l := kvitui.NewLabel(ui, "short")
			l.Wrap = wrap
			w, _ := l.Layout(0).Size()
			if x, _, _ := l.Layout(300).CaretAt(0); x != 0 {
				t.Errorf("wrap %v: a label's line starts at %v", wrap, x)
			}
			l.Align = text.AlignMiddle
			if x, _, _ := l.Layout(300).CaretAt(0); abs(x-(300-w)/2) > 1 {
				t.Errorf("wrap %v: a centred line %v wide starts at %v in 300", wrap, w, x)
			}
			l.Align = text.AlignEnd
			if x, _, _ := l.Layout(300).CaretAt(0); abs(x-(300-w)) > 1 {
				t.Errorf("wrap %v: a line %v wide at the end starts at %v in 300", wrap, w, x)
			}
		}
	})
}

func TestWhatALabelTellsAScreenReader(t *testing.T) {
	var plain, named, heading, ignored, empty *kvitui.Label
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		plain = kvitui.NewLabel(ui, "3 of 12")
		named = kvitui.NewLabel(ui, "3/12")
		named.Name = "3 of 12 files"
		heading = kvitui.NewLabel(ui, "Changes")
		heading.Heading = true
		ignored = kvitui.NewLabel(ui, "·")
		ignored.Ignored = true
		empty = kvitui.NewLabel(ui, "")
		empty.EmptyTakesNoRoom = true
		return []unison.Paneler{plain, named, heading, ignored, empty}
	})
	if n := screen.AccessibilityNodeFor(plain); n == nil || n.Role != role.Label || n.Name != "3 of 12" || n.Ignored {
		t.Errorf("a plain label: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(named); n == nil || n.Role != role.Label || n.Name != "3 of 12 files" {
		t.Errorf("a label with a name: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(heading); n == nil || n.Role != role.Heading || n.Name != "Changes" {
		t.Errorf("a heading: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(ignored); n != nil && !n.Ignored {
		t.Errorf("an ignored label: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(empty); n != nil && !n.Ignored {
		t.Errorf("an empty label that takes no room: %+v", n)
	}
	screen.Do(func() { empty.SetText("Saved") })
	screen.Sync()
	if n := screen.AccessibilityNodeFor(empty); n == nil || n.Ignored || n.Name != "Saved" {
		t.Errorf("the same label once it says something: %+v", n)
	}
}

func TestALabelCutShortShowsItsTextAsATooltipUnlessToldNot(t *testing.T) {
	var plain, quiet *kvitui.Label
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		plain = kvitui.NewLabel(ui, longSentence)
		quiet = kvitui.NewLabel(ui, longSentence)
		quiet.NoTooltip = true
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(120), plain), kvitui.Width(ui, kvitui.Px(120), quiet)}
	})
	screen.Do(func() {
		if !plain.Elided() || !quiet.Elided() {
			t.Fatal("the labels are not cut short, so the test shows nothing")
		}
		plain.UpdateTooltipCallback(geom.Point{}, geom.Rect{})
		quiet.UpdateTooltipCallback(geom.Point{}, geom.Rect{})
		if plain.Tooltip == nil {
			t.Error("a label cut short has no tooltip")
		}
		if quiet.Tooltip != nil {
			t.Error("a label told to show no tooltip has one")
		}
	})
}

func TestTheRoomALabelAsksFor(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		// An empty label is a line tall, unless it is told to take no room.
		empty := kvitui.NewLabel(ui, "")
		_, line, _ := kvitui.NewLabel(ui, "x").Sizes(geom.Size{})
		if _, pref, _ := empty.Sizes(geom.Size{}); pref.Height != line.Height {
			t.Errorf("a plain empty label is %v tall, want a line's %v", pref.Height, line.Height)
		}
		empty.EmptyTakesNoRoom = true
		wantSizes(t, "an empty label that takes no room", empty.AsPanel(), geom.NewSize(200, 0),
			geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, 0))
		empty.Text = "x"
		if _, pref, _ := empty.Sizes(geom.Size{}); pref != line {
			t.Errorf("the same label with text asks for %v, want %v", pref, line)
		}

		// A wrapping label asks for no width of its own, unless told to ask
		// for its text's, up to the width offered.
		words := kvitui.NewLabel(ui, "a few words")
		words.Wrap = true
		own, _ := words.Layout(0).Size()
		if _, pref, _ := words.Sizes(geom.Size{}); pref.Width != 0 {
			t.Errorf("a wrapping label asks for %v wide", pref.Width)
		}
		words.WrapOwnWidth = true
		if _, pref, _ := words.Sizes(geom.Size{}); pref.Width != own {
			t.Errorf("a wrapping label asking for its own width asks for %v, want %v", pref.Width, own)
		}
		if _, pref, _ := words.Sizes(geom.NewSize(own+100, 0)); pref.Width != own {
			t.Errorf("offered more than its text's width it asks for %v, want %v", pref.Width, own)
		}
		_, wrapped := words.Layout(own / 2).Size()
		if _, pref, _ := words.Sizes(geom.NewSize(own/2, 0)); pref.Width != own/2 || pref.Height != wrapped {
			t.Errorf("offered half its text's width it asks for %v, want %v by %v", pref, own/2, wrapped)
		}

		// A size of its own makes the text larger.
		sized := kvitui.NewLabel(ui, "x")
		sized.Size = kvitui.Px(30)
		if _, pref, _ := sized.Sizes(geom.Size{}); pref.Height <= line.Height {
			t.Errorf("a label at 30 design pixels is %v tall, a body label %v", pref.Height, line.Height)
		}
	})
}

// As Kvit Works tests it: in a column 300 design pixels wide, a long text
// limited to two lines is about two lines tall, the same text without the
// limit is over five, and a limit of one is one line.
func TestALabelKeepsToItsLineLimit(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		width := float32(ui.Interface.Px(300))
		_, line, _ := kvitui.NewLabel(ui, "x").Sizes(geom.Size{})
		l := kvitui.NewLabel(ui, longSentence)
		l.Wrap = true
		if _, pref, _ := l.Sizes(geom.NewSize(width, 0)); pref.Height <= 5*line.Height {
			t.Fatalf("the text with no limit is %v tall, a line %v", pref.Height, line.Height)
		}
		l.MaxLines = 2
		if _, pref, _ := l.Sizes(geom.NewSize(width, 0)); pref.Height < 1.5*line.Height || pref.Height > 2.5*line.Height {
			t.Errorf("the text limited to two lines is %v tall, a line %v", pref.Height, line.Height)
		}
		l.MaxLines = 1
		if _, pref, _ := l.Sizes(geom.NewSize(width, 0)); pref.Height != line.Height {
			t.Errorf("the text limited to one line is %v tall, a line %v", pref.Height, line.Height)
		}
		// A text that fits within the limit is left alone.
		short := kvitui.NewLabel(ui, "Two words")
		short.Wrap, short.MaxLines = true, 2
		if _, pref, _ := short.Sizes(geom.NewSize(width, 0)); pref.Height != line.Height {
			t.Errorf("a short text with a limit is %v tall, a line %v", pref.Height, line.Height)
		}
	})
}

func TestSetTextLaysOutOnlyOnAChange(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		l := kvitui.NewLabel(ui, "Saved")
		holder := unison.NewPanel()
		holder.AddChild(l)
		holder.NeedsLayout, l.NeedsLayout = false, false
		l.SetText("Saved")
		if holder.NeedsLayout || l.NeedsLayout {
			t.Error("setting the text a label already has asked for a layout")
		}
		l.SetText("Saving…")
		if l.Text != "Saving…" || !holder.NeedsLayout || !l.NeedsLayout {
			t.Errorf("setting new text left the label saying %q, with layout asked of the holder %v and the label %v", l.Text, holder.NeedsLayout, l.NeedsLayout)
		}
	})
}

// inkRows are the first and last rows of a picture, within a box, that hold
// a pixel much darker than the window's ground.
func inkRows(img image.Image, box geom.Rect) (first, last int) {
	first, last = -1, -1
	for y := int(box.Y); y < int(box.Bottom()); y++ {
		for x := int(box.X); x < int(box.Right()); x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if (r+g+b)/3 < 0x6000 {
				if first < 0 {
					first = y
				}
				last = y
				break
			}
		}
	}
	return first, last
}

// A wrapping label in a box taller than its text is centred on the box,
// unless it is told to start at the top. One unwrapped line is centred
// either way.
func TestAWrappingLabelDrawsFromTheTopWhenTold(t *testing.T) {
	t.Setenv(uitest.ScaleVariable, "1")
	var centred, top, line *kvitui.Label
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		make := func(wrap, atTop bool) (*kvitui.Label, *unison.Panel) {
			l := kvitui.NewLabel(ui, "Two words")
			l.Wrap, l.Top = wrap, atTop
			return l, kvitui.Sized(ui, kvitui.Px(160), kvitui.Px(120), l)
		}
		var a, b, c *unison.Panel
		centred, a = make(true, false)
		top, b = make(true, true)
		line, c = make(false, true)
		return kvitui.Row(ui, kvitui.SizeSpace, a, b, c)
	})
	img := s.Capture()
	where := func(l *kvitui.Label) (offset, height float32) {
		var r geom.Rect
		s.Do(func() { r = l.RectToRoot(l.ContentRect(true)) })
		first, _ := inkRows(img, r)
		if first < 0 {
			t.Fatal("no text was drawn in the label's box")
		}
		return float32(first) - r.Y, r.Height
	}
	if at, h := where(centred); at < h/4 {
		t.Errorf("a wrapping label's text starts %v down a box %v tall, so it is not centred", at, h)
	}
	if at, h := where(top); at > h/6 {
		t.Errorf("a wrapping label told to start at the top starts %v down a box %v tall", at, h)
	}
	if at, h := where(line); at < h/4 {
		t.Errorf("an unwrapped label starts %v down a box %v tall, so it is not centred", at, h)
	}
}
