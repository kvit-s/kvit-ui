// Package text lays out and draws styled text for the Kvit components: it
// finds fonts on the system and falls back to another font for any character
// the chosen one lacks, shapes each run with HarfBuzz (kerning, ligatures,
// emoji sequences, other scripts), breaks lines by the Unicode rules, orders
// mixed-direction runs for display, answers where a caret goes and which
// character a point is over, and draws the result on a unison canvas.
//
// unison's own text support draws one glyph per character without shaping,
// so this package does the layout with go-text and draws through the canvas
// library unison is built on: it builds glyph runs with canvas's textblob
// package and hands them to unison.Canvas.DrawTextBlob.
package text

import (
	"bytes"
	"fmt"
	"github.com/go-text/typesetting/shaping"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	gfont "github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/language"
	cfont "github.com/richardwilkes/canvas/font"
)

// Fonts is the set of fonts text is drawn with: the system's, found once and
// indexed in a cache directory, plus any added from memory, such as the icon
// font. One Fonts serves a whole application.
type Fonts struct {
	mu     sync.Mutex
	fm     *fontscan.FontMap
	data   map[string][]byte                // font bytes by file or identifier
	canvas map[gfont.FontID]*cfont.Typeface // canvas typefaces by file and face index
	sized  map[fontKey]*cfont.Font          // canvas fonts by face and size
	// shaper is kept for the life of the fonts, because it caches what
	// HarfBuzz builds from each font; a shaper made per paragraph rebuilt
	// the kerning tables of every font for every piece of text.
	shaper shaping.HarfbuzzShaper
	// cache holds recent layouts of one span; see Layout.
	cacheMu sync.Mutex
	cache   map[layoutKey]*Layout
}

// NewFonts scans the system's fonts, keeping the index in cacheDir so later
// starts are fast. An empty cacheDir uses the user cache directory.
func NewFonts(cacheDir string) (*Fonts, error) {
	if cacheDir == "" {
		base, err := os.UserCacheDir()
		if err != nil {
			base = os.TempDir()
		}
		cacheDir = filepath.Join(base, "kvit", "fontscan")
	}
	fm := fontscan.NewFontMap(log.New(io.Discard, "", 0))
	if err := fm.UseSystemFonts(cacheDir); err != nil {
		return nil, fmt.Errorf("text: scanning system fonts: %w", err)
	}
	f := &Fonts{
		fm:     fm,
		data:   map[string][]byte{},
		canvas: map[gfont.FontID]*cfont.Typeface{},
	}
	fm.SetFaceLoader(f.shapingFace)
	return f, nil
}

// AddFont makes a font held in memory available under a family name, such as
// the embedded icon font.
func (f *Fonts) AddFont(data []byte, family string) error {
	// A new font can change what an earlier layout would draw with.
	f.cacheMu.Lock()
	f.cache = nil
	f.cacheMu.Unlock()
	f.mu.Lock()
	defer f.mu.Unlock()
	id := "memory:" + family
	face, err := strippedFace(data, 0)
	if err != nil {
		return err
	}
	f.fm.AddFace(face, fontscan.Location{File: id}, gfont.Description{Family: family, Aspect: face.Describe().Aspect})
	f.data[id] = data
	return nil
}

// Generic family names every platform resolves: the desktop's own
// proportional face, and its fixed-pitch one.
const (
	SansSerif = "sans-serif"
	Monospace = "monospace"
)

// familiesFor turns a requested family into the list fontscan tries in order.
// "" and "sans-serif" mean the desktop's own interface face, and "monospace"
// its fixed-pitch face, named explicitly where the generic name would not
// find it (Windows treats "monospace" as an ordinary face name).
func familiesFor(family string) []string {
	switch strings.ToLower(strings.TrimSpace(family)) {
	case "", SansSerif:
		switch runtime.GOOS {
		case "windows":
			return []string{"Segoe UI", SansSerif}
		case "darwin":
			return []string{"SF Pro Text", "Helvetica Neue", SansSerif}
		}
		return []string{SansSerif}
	case Monospace:
		switch runtime.GOOS {
		case "windows":
			return []string{"Cascadia Mono", "Consolas", Monospace}
		case "darwin":
			return []string{"SF Mono", "Menlo", Monospace}
		}
		return []string{Monospace}
	}
	return []string{family, SansSerif}
}

// queryFamilies is what a style asks fontscan for: its own families, then
// the colour emoji fonts. fontscan only falls back to fonts that share a
// family with the request or cover the character's script, and emoji have no
// script of their own, so without these an emoji would find no font.
func queryFamilies(family string) []string {
	return append(familiesFor(family), emojiFamilies...)
}

// scriptFontmap sets each character's script before choosing its font, so
// fontscan's fallback finds fonts made for that script (Devanagari, Han and
// the rest) where they are installed.
type scriptFontmap struct{ fm *fontscan.FontMap }

func (s scriptFontmap) ResolveFace(r rune) *gfont.Face {
	s.fm.SetScript(language.LookupScript(r))
	return s.fm.ResolveFace(r)
}

// ResolveFamily is the family a request actually draws with: never empty,
// so a component can show or store it.
func (f *Fonts) ResolveFamily(family string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fm.SetQuery(fontscan.Query{Families: familiesFor(family)})
	face := f.fm.ResolveFace('a')
	if face == nil {
		return family
	}
	// The font's own name, as its name table spells it; fontscan's metadata
	// holds a normalised form ("dejavusans") meant for matching.
	if fam := face.Font.Describe().Family; fam != "" {
		return fam
	}
	return family
}

// typeface returns the canvas typeface for a go-text face, loading the same
// file and face index so glyph numbers agree.
func (f *Fonts) typeface(face *gfont.Face) (*cfont.Typeface, error) {
	loc := f.fm.FontLocation(face.Font)
	if tf, ok := f.canvas[loc]; ok {
		return tf, nil
	}
	data, err := f.fileData(loc.File)
	if err != nil {
		return nil, err
	}
	tf, err := cfont.NewTypefaceFromData(data, int(loc.Index))
	if err != nil {
		return nil, err
	}
	f.canvas[loc] = tf
	return tf, nil
}

// fileData returns a font file's bytes, mapped from disk once and kept.
func (f *Fonts) fileData(path string) ([]byte, error) {
	if data, ok := f.data[path]; ok {
		return data, nil
	}
	data, err := mapFile(path)
	if err != nil {
		return nil, err
	}
	f.data[path] = data
	return data, nil
}

// shapingFace is the loader fontscan calls once it has chosen a system font.
// Shaping reads the character map, the substitution and positioning tables
// and the metrics, never the glyphs themselves; canvas draws those from its
// own parse of the full file. So the face shaping uses is parsed from a copy
// without the outline, bitmap and colour tables, which go-text would
// otherwise parse and hold for the life of the face.
func (f *Fonts) shapingFace(loc fontscan.Location) (*gfont.Face, error) {
	data, err := f.fileData(loc.File)
	if err != nil {
		return nil, err
	}
	return strippedFace(data, int(loc.Index))
}

// withoutGlyphs are the tables shaping never reads: outlines (TrueType,
// CFF), bitmap strikes, and colour glyphs.
var withoutGlyphs = map[ot.Tag]bool{}

func init() {
	for _, t := range []string{"glyf", "loca", "CFF ", "CFF2", "COLR", "CPAL", "CBDT", "CBLC", "sbix", "SVG ", "EBDT", "EBLC", "EBSC"} {
		withoutGlyphs[ot.MustNewTag(t)] = true
	}
}

// strippedFace parses face index of a font file or collection from a copy
// holding every table but the glyph tables.
func strippedFace(data []byte, index int) (*gfont.Face, error) {
	lds, err := ot.NewLoaders(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if index >= len(lds) {
		return nil, fmt.Errorf("text: face %d of a collection of %d", index, len(lds))
	}
	ld := lds[index]
	var tables []ot.Table
	for _, tag := range ld.Tables() {
		if withoutGlyphs[tag] {
			continue
		}
		raw, err := ld.RawTable(tag)
		if err != nil {
			return nil, err
		}
		tables = append(tables, ot.Table{Tag: tag, Content: raw})
	}
	sort.Slice(tables, func(a, b int) bool { return tables[a].Tag < tables[b].Tag })
	small, err := ot.NewLoader(bytes.NewReader(ot.WriteTTF(tables)))
	if err != nil {
		return nil, err
	}
	ft, err := gfont.NewFont(small)
	if err != nil {
		return nil, err
	}
	return gfont.NewFace(ft), nil
}
