package tokens

import (
	"slices"

	"github.com/kvit-s/kvit-ui/palette"
)

// Settings is where Theme, Interface and Typography keep the reader's
// choices. The settings package's Store implements it; values are what JSON
// decodes to (string, float64, bool).
type Settings interface {
	Value(key string) (any, bool)
	SetValue(key string, value any)
}

// Appearance is what the desktop says about light or dark, high contrast and
// animation. The platform package implements it; each answer is false where
// the desktop says nothing.
type Appearance interface {
	DarkMode() bool
	HighContrast() bool
	ReducedMotion() bool
}

// The setting keys. They are the keys the Kvit versions built with Qt
// wrote, so a settings file those versions saved is read.
const (
	keyThemeID           = "theme.id"
	keyAccent            = "theme.accent"
	keyHighlight         = "theme.highlight"
	keyReducedMotion     = "view.reducedMotion" // the old boolean, read once to carry a choice over
	keyReducedMotionMode = "view.reducedMotionMode"
)

// signal is a list of functions to call when something changes.
type signal struct {
	next int
	fns  map[int]func()
}

func (s *signal) connect(fn func()) (disconnect func()) {
	if s.fns == nil {
		s.fns = map[int]func(){}
	}
	id := s.next
	s.next++
	s.fns[id] = fn
	return func() { delete(s.fns, id) }
}

func (s *signal) emit() {
	ids := make([]int, 0, len(s.fns))
	for id := range s.fns {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		if fn, ok := s.fns[id]; ok {
			fn()
		}
	}
}

// Theme is the chosen theme and the colours it resolves to. It starts on the
// light theme with no settings, which is what a test gets; an application
// attaches its settings and the desktop's appearance.
type Theme struct {
	settings             Settings
	appearance           Appearance
	stopFollowing        func()
	themeID              string
	reducedMotionSetting string
	resolved             string
	accentOverride       string
	highlightOverride    string
	tokens               Tokens
	changed              signal
	motionChanged        signal
}

// NewTheme returns a theme on the light table.
func NewTheme() *Theme {
	t := &Theme{themeID: Light, reducedMotionSetting: "system", resolved: Light}
	t.refresh()
	return t
}

// OnChanged calls fn whenever any colour may have changed.
func (t *Theme) OnChanged(fn func()) (disconnect func()) { return t.changed.connect(fn) }

// OnReducedMotionChanged calls fn whenever the effective reduced-motion value
// or its setting may have changed.
func (t *Theme) OnReducedMotionChanged(fn func()) (disconnect func()) {
	return t.motionChanged.connect(fn)
}

// SetSettings reads the stored choices and saves every later change back.
// A first start, with nothing stored, follows the desktop.
func (t *Theme) SetSettings(s Settings) {
	t.settings = s
	if s == nil {
		return
	}
	t.themeID = stringValue(s, keyThemeID, System)
	if !slices.Contains(AvailableThemes(), t.themeID) {
		t.themeID = System // stale or hand-edited
	}
	t.accentOverride = validColorOrEmpty(stringValue(s, keyAccent, ""))
	t.highlightOverride = validColorOrEmpty(stringValue(s, keyHighlight, ""))
	// Reduced motion: the three-way setting if it has been written, else the
	// old boolean, whose presence at all means somebody chose, else "system".
	previous := t.reducedMotionSetting
	if v, ok := s.Value(keyReducedMotionMode); ok {
		t.reducedMotionSetting, _ = v.(string)
	} else if v, ok := s.Value(keyReducedMotion); ok {
		if b, _ := v.(bool); b {
			t.reducedMotionSetting = "on"
		} else {
			t.reducedMotionSetting = "off"
		}
	} else {
		t.reducedMotionSetting = "system"
	}
	if !slices.Contains(ReducedMotionSettings(), t.reducedMotionSetting) {
		t.reducedMotionSetting = "system"
	}
	if previous != t.reducedMotionSetting {
		t.motionChanged.emit()
	}
	t.refresh()
}

// SetAppearance attaches what the desktop says. Without one, "system"
// resolves to light and following the system for motion means off. An
// Appearance that can announce its own changes (an OnChanged method) is
// followed from then on; otherwise call AppearanceChanged after a change.
func (t *Theme) SetAppearance(a Appearance) {
	if t.stopFollowing != nil {
		t.stopFollowing()
		t.stopFollowing = nil
	}
	t.appearance = a
	if n, ok := a.(interface{ OnChanged(func()) func() }); ok {
		t.stopFollowing = n.OnChanged(t.AppearanceChanged)
	}
	t.AppearanceChanged()
}

// AppearanceChanged re-resolves after the desktop's appearance changed.
func (t *Theme) AppearanceChanged() {
	t.motionChanged.emit()
	t.refresh()
}

// ThemeID is the chosen theme, possibly "system".
func (t *Theme) ThemeID() string { return t.themeID }

// SetThemeID chooses a theme. An id AvailableThemes does not list is ignored.
func (t *Theme) SetThemeID(id string) {
	if t.themeID == id || !slices.Contains(AvailableThemes(), id) {
		return
	}
	t.themeID = id
	t.save(keyThemeID, id)
	t.refresh()
}

// ResolvedTheme is what the theme currently renders as; never "system".
func (t *Theme) ResolvedTheme() string { return t.resolved }

// Tokens returns the effective colours: the resolved table with the
// overrides applied.
func (t *Theme) Tokens() Tokens { return t.tokens }

// ReducedMotionSetting is "on", "off" or "system".
func (t *Theme) ReducedMotionSetting() string { return t.reducedMotionSetting }

// SetReducedMotionSetting sets "on", "off" or "system"; anything else is ignored.
func (t *Theme) SetReducedMotionSetting(mode string) {
	if t.reducedMotionSetting == mode || !slices.Contains(ReducedMotionSettings(), mode) {
		return
	}
	t.reducedMotionSetting = mode
	t.save(keyReducedMotionMode, mode)
	t.motionChanged.emit()
}

// ReducedMotion is the effective value: the explicit choice, or under
// "system" what the desktop says, and off where it says nothing.
func (t *Theme) ReducedMotion() bool {
	switch t.reducedMotionSetting {
	case "on":
		return true
	case "off":
		return false
	}
	return t.appearance != nil && t.appearance.ReducedMotion()
}

// SetReducedMotion makes an explicit choice, as a checkable menu item does.
func (t *Theme) SetReducedMotion(reduced bool) {
	if reduced {
		t.SetReducedMotionSetting("on")
	} else {
		t.SetReducedMotionSetting("off")
	}
}

// MotionScale is 0 under reduced motion and 1 otherwise; an animation
// multiplies its duration by it.
func (t *Theme) MotionScale() float64 {
	if t.ReducedMotion() {
		return 0
	}
	return 1
}

// AccentOverride is the reader's accent colour, or "" for the theme's own.
func (t *Theme) AccentOverride() string { return t.accentOverride }

// SetAccentOverride sets the accent colour; "" or an unparsable value clears it.
func (t *Theme) SetAccentOverride(hex string) {
	v := validColorOrEmpty(hex)
	if t.accentOverride == v {
		return
	}
	t.accentOverride = v
	t.save(keyAccent, v)
	t.refresh()
}

// HighlightOverride is the reader's highlight colour, or "" for the theme's own.
func (t *Theme) HighlightOverride() string { return t.highlightOverride }

// SetHighlightOverride sets the highlight colour; "" or an unparsable value clears it.
func (t *Theme) SetHighlightOverride(hex string) {
	v := validColorOrEmpty(hex)
	if t.highlightOverride == v {
		return
	}
	t.highlightOverride = v
	t.save(keyHighlight, v)
	t.refresh()
}

func (t *Theme) resolveSystem() string {
	// High contrast first: turned on system-wide, it outranks light or dark,
	// and the high-contrast theme answers both at once. Only under "system":
	// an explicit theme choice wins.
	if t.appearance != nil && t.appearance.HighContrast() {
		return HighContrast
	}
	if t.appearance != nil && t.appearance.DarkMode() {
		return Dark
	}
	return Light
}

func (t *Theme) refresh() {
	t.resolved = t.themeID
	if t.themeID == System {
		t.resolved = t.resolveSystem()
	}
	t.tokens = TokensFor(t.resolved)
	// The overrides ride on top of the table; components read plain tokens
	// and never know overrides exist.
	if t.accentOverride != "" {
		a := palette.Hex(t.accentOverride)
		t.tokens.Accent = a
		t.tokens.Link = a
		t.tokens.OnAccent = LabelOn(a)
	}
	if t.highlightOverride != "" {
		t.tokens.HighlightBackground = palette.Hex(t.highlightOverride)
	}
	t.changed.emit()
}

func (t *Theme) save(key string, value any) {
	if t.settings != nil {
		t.settings.SetValue(key, value)
	}
}

func validColorOrEmpty(v string) string {
	if _, err := palette.ParseHex(v); err != nil {
		return ""
	}
	return v
}

// stringValue reads a stored string, or fallback when there is none.
func stringValue(s Settings, key, fallback string) string {
	if v, ok := s.Value(key); ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return fallback
}

// numberValue reads a stored number, or fallback when there is none. JSON
// numbers decode to float64; an int is accepted too, as a test may store one.
func numberValue(s Settings, key string, fallback float64) float64 {
	if v, ok := s.Value(key); ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		}
	}
	return fallback
}
