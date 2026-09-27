package text_test

import (
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

var (
	fontsOnce sync.Once
	fonts     *text.Fonts
	fontsErr  error
)

func sharedFonts(t *testing.T) *text.Fonts {
	t.Helper()
	fontsOnce.Do(func() {
		fonts, fontsErr = text.NewFonts("")
		if fontsErr == nil {
			fontsErr = fonts.AddFont(icons.Font, icons.FontFamily)
		}
	})
	if fontsErr != nil {
		t.Fatal(fontsErr)
	}
	return fonts
}

var black = text.Color{R: 0x1a, G: 0x1a, B: 0x1a, A: 255}

func plain(s string, size float32) []text.Span {
	return []text.Span{{Text: s, Style: text.Style{Size: size, Color: black}}}
}

func TestKerningNarrowsAPair(t *testing.T) {
	f := sharedFonts(t)
	w := func(s string) float32 { w, _ := f.Layout(plain(s, 40), text.Options{}).Size(); return w }
	// With kerning, "AV" is narrower than "A" and "V" set apart.
	if av, a, v := w("AV"), w("A"), w("V"); av >= a+v-0.5 {
		t.Errorf("AV is %.1f wide, A and V %.1f and %.1f: no kerning applied", av, a, v)
	}
}

func TestEmojiSequencesShapeToOneGlyph(t *testing.T) {
	f := sharedFonts(t)
	for _, s := range []string{"👍🏽", "🇺🇸", "👨‍👩‍👧", "1️⃣"} {
		l := f.Layout(plain(s, 24), text.Options{})
		w, _ := l.Size()
		single, _ := f.Layout(plain("😀", 24), text.Options{}).Size()
		// One emoji wide, give or take, rather than two or three.
		if w > single*1.5 {
			t.Errorf("%q is %.1f wide against %.1f for one emoji: it was not joined", s, w, single)
		}
	}
}

func TestWrapsAtTheWidthAndAtLineBreaks(t *testing.T) {
	f := sharedFonts(t)
	long := "The quick brown fox jumps over the lazy dog, and then it does it again."
	l := f.Layout(plain(long, 16), text.Options{MaxWidth: 200})
	if l.LineCount() < 2 {
		t.Errorf("a long sentence at 200 px made %d line", l.LineCount())
	}
	if w, _ := l.Size(); w > 200 {
		t.Errorf("a wrapped line is %.1f wide, over 200", w)
	}
	if n := f.Layout(plain("one\ntwo\nthree", 16), text.Options{}).LineCount(); n != 3 {
		t.Errorf("three lines separated by line breaks made %d lines", n)
	}
	// Chinese has no spaces; it must still wrap between characters.
	cjk := "这是一个没有空格的很长的中文句子它必须在字符之间换行才能放进窄的列里"
	if n := f.Layout(plain(cjk, 16), text.Options{MaxWidth: 120}).LineCount(); n < 3 {
		t.Errorf("a long Chinese sentence at 120 px made %d lines", n)
	}
}

func TestLineHeightMultiplies(t *testing.T) {
	f := sharedFonts(t)
	_, h1 := f.Layout(plain("x", 20), text.Options{}).Size()
	_, h2 := f.Layout(plain("x", 20), text.Options{LineHeight: 1.5}).Size()
	if d := h2 - h1*1.5; d > 0.01 || d < -0.01 {
		t.Errorf("line height 1.5 gave %.2f against %.2f natural", h2, h1)
	}
	if _, h := f.Layout(nil, text.Options{}).Size(); h <= 0 {
		t.Error("empty text must still take a line")
	}
}

func TestCaretAndHitTestAgree(t *testing.T) {
	f := sharedFonts(t)
	s := "office hello"
	l := f.Layout(plain(s, 20), text.Options{})
	prev := float32(-1)
	for i := 0; i <= l.Len(); i++ {
		x, _, _ := l.CaretAt(i)
		if x < prev {
			t.Errorf("the caret moves left from index %d to %d", i-1, i)
		}
		prev = x
		if got := l.IndexAt(x, 1); got != i {
			t.Errorf("the caret at index %d is at x %.1f, which hits index %d", i, x, got)
		}
	}
	// A point past the end lands at the end.
	if got := l.IndexAt(10000, 1); got != l.Len() {
		t.Errorf("past the end hits %d, want %d", got, l.Len())
	}
}

func TestStylesKeepTheirOwnSize(t *testing.T) {
	f := sharedFonts(t)
	l := f.Layout([]text.Span{
		{Text: "Big ", Style: text.Style{Size: 32, Weight: text.Bold, Color: black}},
		{Text: "small", Style: text.Style{Size: 12, Italic: true, Color: black}},
	}, text.Options{})
	_, h := l.Size()
	_, hSmall := f.Layout(plain("small", 12), text.Options{}).Size()
	if h <= hSmall*1.5 {
		t.Errorf("a line mixing 32 px and 12 px is only %.1f tall", h)
	}
}

func TestResolveFamilyNamesAFace(t *testing.T) {
	f := sharedFonts(t)
	for _, fam := range []string{"", text.SansSerif, text.Monospace} {
		if got := f.ResolveFamily(fam); got == "" {
			t.Errorf("%q resolved to nothing", fam)
		}
	}
}

// TestDrawSample draws the kinds of text the checks above cover and saves the
// picture next to the test's output, to look at.
func TestDrawSample(t *testing.T) {
	f := sharedFonts(t)
	lines := [][]text.Span{
		plain("Kerning: AVATAR To Wa Ty — office ﬁ", 22),
		plain("Emoji: 👍 👍🏽 🇺🇸 🇩🇪 👨‍👩‍👧 1️⃣", 22),
		{
			{Text: "Styled: ", Style: text.Style{Size: 18, Color: black}},
			{Text: "bold", Style: text.Style{Size: 18, Weight: text.Bold, Color: black}},
			{Text: ", ", Style: text.Style{Size: 18, Color: black}},
			{Text: "italic", Style: text.Style{Size: 18, Italic: true, Color: black}},
			{Text: ", ", Style: text.Style{Size: 18, Color: black}},
			{Text: "code", Style: text.Style{Size: 16, Family: text.Monospace, Color: black, Background: text.Color{0xf0, 0xf0, 0xee, 255}}},
			{Text: ", ", Style: text.Style{Size: 18, Color: black}},
			{Text: "link", Style: text.Style{Size: 18, Color: text.Color{0x29, 0x70, 0xc8, 255}, Underline: true}},
			{Text: ", ", Style: text.Style{Size: 18, Color: black}},
			{Text: "struck", Style: text.Style{Size: 18, Color: black, Strike: true}},
		},
		plain("Scripts: Ελληνικά, Русский, العربية, עברית, हिन्दी, 中文, 日本語", 20),
		{{Text: string([]rune{mustGlyph(t, "search"), ' ', mustGlyph(t, "settings"), ' ', mustGlyph(t, "trend-up"), ' ', mustGlyph(t, "folder")}),
			Style: text.Style{Family: icons.FontFamily, Size: 26, Color: black}}},
	}
	var layouts []*text.Layout
	for _, l := range lines {
		layouts = append(layouts, f.Layout(l, text.Options{MaxWidth: 760}))
	}
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 800, Height: 300},
		unison.StartupFinishedCallback(func() {
			wnd, err := unison.NewWindow("text")
			if err != nil {
				t.Error(err)
				return
			}
			p := unison.NewPanel()
			p.DrawCallback = func(gc *unison.Canvas, r geom.Rect) {
				gc.DrawRect(r, unison.White.Paint(gc, r, paintstyle.Fill))
				y := float32(12)
				for _, l := range layouts {
					l.Draw(gc, 12, y)
					_, h := l.Size()
					y += h + 8
				}
			}
			p.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
			wnd.Content().SetLayout(&unison.FlexLayout{Columns: 1})
			wnd.Content().AddChild(p)
			wnd.SetContentRect(geom.NewRect(0, 0, 800, 300))
			wnd.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()
	dir := os.Getenv("KVIT_SHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	out, err := os.Create(filepath.Join(dir, "text-sample.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if err := png.Encode(out, screen.Capture()); err != nil {
		t.Fatal(err)
	}
	for _, e := range screen.Errors() {
		t.Error(e)
	}
}

func mustGlyph(t *testing.T, name string) rune {
	r, ok := icons.Glyph(name)
	if !ok {
		t.Fatalf("no icon %q", name)
	}
	return r
}

// The last glyph of a line keeps its width. go-text's wrapper decides what
// is a trailing space by ink width, and shaping faces carry no ink (see
// markInk), so a line of several runs, or a wrapped line, lost its final
// glyph's width before markInk.
func TestTheLastGlyphOfALineKeepsItsWidth(t *testing.T) {
	f := sharedFonts(t)
	bold := text.Style{Size: 20, Weight: text.Bold, Color: black}
	two := f.Layout([]text.Span{{Text: "link, ", Style: bold}, {Text: "struck", Style: text.Style{Size: 20, Color: black}}}, text.Options{})
	w, _ := two.Size()
	a, _ := f.Layout([]text.Span{{Text: "link, ", Style: bold}}, text.Options{KeepTrailingSpace: true}).Size()
	b, _ := f.Layout(plain("struck", 20), text.Options{}).Size()
	if d := w - (a + b); d < -0.5 || d > 0.5 {
		t.Errorf("two runs are %.2f wide, their parts %.2f: the last glyph lost its width", w, a+b)
	}
	if x, _, _ := two.CaretAt(two.Len()); x < w-0.5 {
		t.Errorf("the caret after the last letter is at %.2f, inside the %.2f wide line", x, w)
	}
	// Each line of a wrapped paragraph is as wide as its own text alone.
	long := "abcdefghij klmnopqrst uvwxyzabcd efghijklmn opqrstuvwx"
	rs := []rune(long)
	for _, width := range []float32{90, 140, 200} {
		l := f.Layout(plain(long, 20), text.Options{MaxWidth: width})
		for i := 0; i < l.LineCount(); i++ {
			start, end, _, _ := l.LineBounds(i)
			for end > start && rs[end-1] == ' ' {
				end--
			}
			alone, _ := f.Layout(plain(string(rs[start:end]), 20), text.Options{}).Size()
			x0, _, _ := l.CaretAt(start)
			x1, _, _ := l.CaretAt(end)
			if i < l.LineCount()-1 {
				// the caret at a wrapped line's end is placed on the next line
				last := []rune{rs[end-1]}
				lw, _ := f.Layout(plain(string(last), 20), text.Options{}).Size()
				xl, _, _ := l.CaretAt(end - 1)
				x1 = xl + lw
			}
			if d := (x1 - x0) - alone; d < -0.5 || d > 0.5 {
				t.Errorf("at width %.0f line %d is %.2f wide, its text alone %.2f", width, i, x1-x0, alone)
			}
		}
	}
}
