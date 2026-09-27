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
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	gfont "github.com/go-text/typesetting/font"
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
	return &Fonts{
		fm:     fm,
		data:   map[string][]byte{},
		canvas: map[gfont.FontID]*cfont.Typeface{},
	}, nil
}

// AddFont makes a font held in memory available under a family name, such as
// the embedded icon font.
func (f *Fonts) AddFont(data []byte, family string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	id := "memory:" + family
	if err := f.fm.AddFont(bytes.NewReader(data), id, family); err != nil {
		return err
	}
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
	data, ok := f.data[loc.File]
	if !ok {
		var err error
		if data, err = os.ReadFile(loc.File); err != nil {
			return nil, err
		}
		f.data[loc.File] = data
	}
	tf, err := cfont.NewTypefaceFromData(data, int(loc.Index))
	if err != nil {
		return nil, err
	}
	f.canvas[loc] = tf
	return tf, nil
}
