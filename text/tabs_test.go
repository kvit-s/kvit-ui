package text_test

import (
	"testing"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

// A tab reaches the next tab stop, four of its font's spaces apart, counted
// from the start of the paragraph: two tabs indent a line of code by eight
// spaces, and a tab after two characters reaches the fourth space.
func TestATabReachesTheNextStop(t *testing.T) {
	f := sharedFonts(t)
	for _, family := range []string{text.Monospace, text.SansSerif} {
		st := text.Style{Family: family, Size: 14, Color: black}
		lay := func(s string) *text.Layout { return f.Layout([]text.Span{{Text: s, Style: st}}, text.Options{}) }
		space, _ := lay("    ").Size()
		width := func(s string) float32 { w, _ := lay(s).Size(); return w }
		// next is the first stop past x.
		next := func(x float32) float32 { return float32(int(x/space)+1) * space }
		near := func(got, want float32) bool { return got > want-0.05 && got < want+0.05 }
		cases := []struct {
			text  string
			index int
			want  float32
		}{
			{"\tx", 1, space},
			{"\t\tx", 2, 2 * space},
			{"ab\tx", 3, next(width("ab"))},
			{"\tx\ty", 3, next(space + width("x"))},
		}
		for _, c := range cases {
			if x, _, _ := lay(c.text).CaretAt(c.index); !near(x, c.want) {
				t.Errorf("%s: in %q the character after the tab starts at %.2f, want %.2f", family, c.text, x, c.want)
			}
		}
		// A tab alone in its span, as an indent before a keyword coloured apart
		// from it is: it is shaped in whichever font has a glyph for it, and
		// still reaches four of its own span's spaces.
		keyword := st
		keyword.Color = text.Color{R: 0x9a, G: 0x2f, B: 0x8a, A: 255}
		alone := f.Layout([]text.Span{{Text: "\t", Style: st}, {Text: "if", Style: keyword}}, text.Options{})
		if x, _, _ := alone.CaretAt(1); !near(x, space) {
			t.Errorf("%s: after a tab alone in its span the keyword starts at %.2f, want %.2f", family, x, space)
		}
		// Two lines of a code block, one of them indented, as the editor lays
		// it out: the second paragraph counts its stops from its own start.
		l := f.Layout([]text.Span{{Text: "if x {\n\t\treturn\n}", Style: st}}, text.Options{})
		if x, _, _ := l.CaretAt(9); !near(x, 2*space) {
			t.Errorf("%s: the indented line's word starts at %.2f, want %.2f", family, x, 2*space)
		}
	}
}

// A tab draws nothing: it had been the font's glyph for U+0009, which most
// fonts lack, so it drew as the empty box of a missing character.
func TestATabDrawsNothing(t *testing.T) {
	f := sharedFonts(t)
	st := text.Style{Family: text.Monospace, Size: 20, Color: black}
	l := f.Layout([]text.Span{{Text: "\t\tx", Style: st}}, text.Options{})
	tabs, _, _ := l.CaretAt(2)
	const left, top, width, height = 10, 10, 300, 60
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: width, Height: height},
		unison.StartupFinishedCallback(func() {
			wnd, err := unison.NewWindow("tabs")
			if err != nil {
				t.Error(err)
				return
			}
			p := unison.NewPanel()
			p.DrawCallback = func(gc *unison.Canvas, r geom.Rect) {
				gc.DrawRect(r, unison.White.Paint(gc, r, paintstyle.Fill))
				l.Draw(gc, left, top)
			}
			p.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
			wnd.Content().SetLayout(&unison.FlexLayout{Columns: 1})
			wnd.Content().AddChild(p)
			wnd.SetContentRect(geom.NewRect(0, 0, width, height))
			wnd.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()
	pic := screen.Capture()
	inked := func(x0, x1 int) int {
		n := 0
		for y := 0; y < height; y++ {
			for x := x0; x < x1; x++ {
				if r, g, b, _ := pic.At(x, y).RGBA(); (r+g+b)/3 < 0xc000 {
					n++
				}
			}
		}
		return n
	}
	if n := inked(0, left+int(tabs)-1); n > 0 {
		t.Errorf("%d dark pixels where the two tabs are, left of %.0f", n, left+tabs)
	}
	if n := inked(left+int(tabs), width); n == 0 {
		t.Error("the x after the tabs was not drawn")
	}
	for _, e := range screen.Errors() {
		t.Error(e)
	}
}
