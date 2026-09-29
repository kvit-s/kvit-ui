package tokens_test

import (
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/settings"
	"github.com/kvit-s/kvit-ui/tokens"
)

// Ported from kvit-ui's tests/test_theme.cpp.

// fakeAppearance stands in for the desktop: no machine running the tests has
// high contrast turned on, and what matters is on this side of the platform.
type fakeAppearance struct {
	dark, high, reduced bool
	listeners           []func()
}

func (f *fakeAppearance) DarkMode() bool      { return f.dark }
func (f *fakeAppearance) HighContrast() bool  { return f.high }
func (f *fakeAppearance) ReducedMotion() bool { return f.reduced }
func (f *fakeAppearance) OnChanged(fn func()) func() {
	f.listeners = append(f.listeners, fn)
	return func() {}
}
func (f *fakeAppearance) set(high, reduced bool) {
	f.high, f.reduced = high, reduced
	for _, fn := range f.listeners {
		fn()
	}
}

func openStore(t *testing.T, name string) *settings.Store {
	t.Helper()
	s := settings.New()
	if err := s.Open(filepath.Join(t.TempDir(), name), false); err != nil {
		t.Fatal(err)
	}
	return s
}

func counter(connect func(func()) func()) *int {
	n := 0
	connect(func() { n++ })
	return &n
}

func TestDefaultIsLightTable(t *testing.T) {
	th := tokens.NewTheme()
	if th.ThemeID() != tokens.Light || th.ResolvedTheme() != tokens.Light {
		t.Fatalf("a new theme is %q resolving to %q", th.ThemeID(), th.ResolvedTheme())
	}
	light := tokens.TokensFor(tokens.Light)
	got := th.Tokens()
	if !got.WindowBackground.Equal(light.WindowBackground) || !got.Accent.Equal(light.Accent) {
		t.Error("a new theme does not draw with the light table")
	}
	for name, want := range map[string]string{
		"Marker": "#949494", "Link": "#2970c8", "HighlightBackground": "#fdf3a9",
		"InlineCodeBackground": "#f0f0ee", "SearchMatchBackground": "#b5dcff", "SearchCurrentBackground": "#ffb454",
	} {
		if c := reflect.ValueOf(got).FieldByName(name).Interface().(palette.Color); c.Hex() != want {
			t.Errorf("light %s is %s, want %s", name, c, want)
		}
	}
}

func TestThemeSwitchSwapsTokens(t *testing.T) {
	th := tokens.NewTheme()
	n := counter(th.OnChanged)
	th.SetThemeID(tokens.Dark)
	if *n != 1 || th.ResolvedTheme() != tokens.Dark || !th.Tokens().WindowBackground.Equal(tokens.TokensFor(tokens.Dark).WindowBackground) {
		t.Errorf("switching to dark: %d notifications, resolved %q", *n, th.ResolvedTheme())
	}
	th.SetThemeID(tokens.Sepia)
	if *n != 2 || !th.Tokens().WindowBackground.Equal(tokens.TokensFor(tokens.Sepia).WindowBackground) {
		t.Error("switching to sepia did not swap the tokens")
	}
	th.SetThemeID(tokens.Sepia)
	if *n != 2 {
		t.Error("choosing the current theme again notified")
	}
}

func TestInvalidThemeIdRejected(t *testing.T) {
	th := tokens.NewTheme()
	n := counter(th.OnChanged)
	th.SetThemeID("neon")
	if th.ThemeID() != tokens.Light || *n != 0 {
		t.Errorf("an unknown theme id was taken: %q, %d notifications", th.ThemeID(), *n)
	}
}

func TestSystemResolvesToLightOrDark(t *testing.T) {
	th := tokens.NewTheme()
	th.SetThemeID(tokens.System)
	if th.ThemeID() != tokens.System || th.ResolvedTheme() != tokens.Light {
		t.Errorf("system with no desktop answer resolved to %q", th.ResolvedTheme())
	}
	th.SetAppearance(&fakeAppearance{dark: true})
	if th.ResolvedTheme() != tokens.Dark || !th.Tokens().WindowBackground.Equal(tokens.TokensFor(tokens.Dark).WindowBackground) {
		t.Errorf("system on a dark desktop resolved to %q", th.ResolvedTheme())
	}
}

// Every colour field of Tokens is set by every table: the generator lists the
// names the tables set, and a field missing from them would be black.
func TestTablesAreCompleteAndDistinct(t *testing.T) {
	colorType := reflect.TypeOf(palette.Color{})
	typ := reflect.TypeOf(tokens.Tokens{})
	set := tokens.TokenNames()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type == colorType && f.Name != "OnAccent" && !slices.Contains(set, f.Name) {
			t.Errorf("Tokens.%s is not set by the tables", f.Name)
		}
	}
	for _, id := range tokens.BuiltInThemes() {
		tk := tokens.TokensFor(id)
		code := []palette.Color{tk.CodeKeyword, tk.CodeType, tk.CodeString, tk.CodeComment, tk.CodeNumber}
		for i := range code {
			for j := i + 1; j < len(code); j++ {
				if code[i].Equal(code[j]) {
					t.Errorf("%s: code tokens %d and %d are the same colour", id, i, j)
				}
			}
		}
		// Added, removed and changed mark three different claims about a line.
		v := []palette.Color{tk.AddedTextBackground, tk.RemovedTextBackground, tk.ChangedTextBackground}
		for i := range v {
			for j := i + 1; j < len(v); j++ {
				if v[i].Equal(v[j]) {
					t.Errorf("%s: version tints %d and %d are the same colour", id, i, j)
				}
			}
		}
	}
	bg := func(id string) palette.Color { return tokens.TokensFor(id).WindowBackground }
	if bg(tokens.Light).Equal(bg(tokens.Dark)) || bg(tokens.Light).Equal(bg(tokens.Sepia)) || bg(tokens.Dark).Equal(bg(tokens.Sepia)) {
		t.Error("the light, dark and sepia themes should have different backgrounds")
	}
}

// A coarse legibility floor on plain (not linearised) luminance, as the
// test states it.
func TestDarkAndSepiaKeepContrast(t *testing.T) {
	lum := func(c palette.Color) float64 { return 0.2126*c.R + 0.7152*c.G + 0.0722*c.B }
	abs := func(v float64) float64 { return max(v, -v) }
	for _, id := range []string{tokens.Light, tokens.Dark, tokens.Sepia} {
		tk := tokens.TokensFor(id)
		if abs(lum(tk.TextPrimary)-lum(tk.WindowBackground)) <= 0.55 {
			t.Errorf("%s: body text vs background", id)
		}
		if abs(lum(tk.TextPrimary)-lum(tk.SelectionTint)) <= 0.35 {
			t.Errorf("%s: text vs selection tint", id)
		}
		if abs(lum(tk.TextPrimary)-lum(tk.SearchMatchBackground)) <= 0.3 {
			t.Errorf("%s: text vs search tint", id)
		}
	}
}

// Every foreground/background pair that appears together on screen, in every
// theme, against WCAG 2.1 AA: 4.5:1 for text, 3:1 for a control's boundary or
// state mark. Border is absent on purpose: it is decorative.
func TestEveryTokenPairMeetsItsFloor(t *testing.T) {
	type pair struct {
		name   string
		fg, bg func(tokens.Tokens) palette.Color
		floor  float64
	}
	pairs := []pair{
		{"onAccent on accent", func(k tokens.Tokens) palette.Color { return k.OnAccent }, func(k tokens.Tokens) palette.Color { return k.Accent }, 4.5},
		{"warning on window", func(k tokens.Tokens) palette.Color { return k.Warning }, func(k tokens.Tokens) palette.Color { return k.WindowBackground }, 4.5},
		{"success on window", func(k tokens.Tokens) palette.Color { return k.Success }, func(k tokens.Tokens) palette.Color { return k.WindowBackground }, 4.5},
		{"textFaint on window", func(k tokens.Tokens) palette.Color { return k.TextFaint }, func(k tokens.Tokens) palette.Color { return k.WindowBackground }, 4.5},
		{"codeComment on code panel", func(k tokens.Tokens) palette.Color { return k.CodeComment }, func(k tokens.Tokens) palette.Color { return k.CodePanelBackground }, 4.5},
		{"codeString on code panel", func(k tokens.Tokens) palette.Color { return k.CodeString }, func(k tokens.Tokens) palette.Color { return k.CodePanelBackground }, 4.5},
		{"codeType on code panel", func(k tokens.Tokens) palette.Color { return k.CodeType }, func(k tokens.Tokens) palette.Color { return k.CodePanelBackground }, 4.5},
		{"bannerText on banner", func(k tokens.Tokens) palette.Color { return k.BannerText }, func(k tokens.Tokens) palette.Color { return k.BannerBackground }, 4.5},
		{"textPrimary on the current match", func(k tokens.Tokens) palette.Color { return k.TextPrimary }, func(k tokens.Tokens) palette.Color { return k.SearchCurrentBackground }, 4.5},
		{"textPrimary on active selection", func(k tokens.Tokens) palette.Color { return k.TextPrimary }, func(k tokens.Tokens) palette.Color { return k.SelectionActiveTint }, 4.5},
		{"textPrimary on changed text", func(k tokens.Tokens) palette.Color { return k.TextPrimary }, func(k tokens.Tokens) palette.Color { return k.ChangedTextBackground }, 4.5},
		{"textPrimary on an added line", func(k tokens.Tokens) palette.Color { return k.TextPrimary }, func(k tokens.Tokens) palette.Color { return k.AddedTextBackground }, 4.5},
		{"textPrimary on a removed line", func(k tokens.Tokens) palette.Color { return k.TextPrimary }, func(k tokens.Tokens) palette.Color { return k.RemovedTextBackground }, 4.5},
		{"borderStrong on window", func(k tokens.Tokens) palette.Color { return k.BorderStrong }, func(k tokens.Tokens) palette.Color { return k.WindowBackground }, 3},
		{"mutedGlyph on panel", func(k tokens.Tokens) palette.Color { return k.MutedGlyph }, func(k tokens.Tokens) palette.Color { return k.PanelBackground }, 3},
		{"quoteBar on window", func(k tokens.Tokens) palette.Color { return k.QuoteBar }, func(k tokens.Tokens) palette.Color { return k.WindowBackground }, 3},
		{"marker on window", func(k tokens.Tokens) palette.Color { return k.Marker }, func(k tokens.Tokens) palette.Color { return k.WindowBackground }, 3},
	}
	for _, id := range tokens.BuiltInThemes() {
		tk := tokens.TokensFor(id)
		for _, p := range pairs {
			if r := palette.ContrastRatio(p.fg(tk), p.bg(tk)); r < p.floor {
				t.Errorf("%s: %s is only %.2f:1 (need %.1f:1)", id, p.name, r, p.floor)
			}
		}
	}
}

func TestAccentLabelIsDerivedFromTheAccent(t *testing.T) {
	th := tokens.NewTheme()
	th.SetAccentOverride("#ffe9a0")
	if r := palette.ContrastRatio(th.Tokens().OnAccent, th.Tokens().Accent); r < 4.5 {
		t.Errorf("the label on a pale accent is %.2f:1", r)
	}
	onPale := th.Tokens().OnAccent
	th.SetAccentOverride("#102030")
	if r := palette.ContrastRatio(th.Tokens().OnAccent, th.Tokens().Accent); r < 4.5 {
		t.Errorf("the label on a dark accent is %.2f:1", r)
	}
	if th.Tokens().OnAccent.Equal(onPale) {
		t.Error("the label has to move when the accent moves from pale to dark")
	}
	th.SetAccentOverride("")
	if !th.Tokens().OnAccent.Equal(tokens.TokensFor(tokens.Light).OnAccent) {
		t.Error("clearing the override did not restore the theme's own label")
	}
}

func TestPaletteNamesCoverEveryColor(t *testing.T) {
	if len(tokens.ColorPaletteNames()) != len(tokens.ColorPalette()) || len(tokens.HighlightPaletteNames()) != len(tokens.HighlightPalette()) {
		t.Fatal("the name lists and the palettes differ in length")
	}
	for i, c := range tokens.ColorPalette() {
		want := tokens.ColorPaletteNames()[i]
		if got := tokens.ColorName(c.Hex()); got != want {
			t.Errorf("%s is named %q, want %q", c, got, want)
		}
		// Spelling must not decide the answer.
		upper := "#" + string([]byte(c.Hex()[1:]))
		if got := tokens.ColorName(toUpper(upper)); got != want {
			t.Errorf("%s in upper case is named %q, want %q", c, got, want)
		}
	}
	if tokens.ColorName("") == "" || tokens.ColorName("#123456") == "" {
		t.Error("the default and a custom colour must have names too")
	}
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'f' {
			b[i] = c - 32
		}
	}
	return string(b)
}

// The keyboard focus ring is held to WCAG's 3:1 non-text floor. The  test
// linearises with WCAG 2.0's 0.03928 threshold, reproduced here.
func wcag20Ratio(a, b palette.Color) float64 {
	lin := func(c float64) float64 {
		if c <= 0.03928 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	rel := func(c palette.Color) float64 { return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B) }
	hi, lo := max(rel(a), rel(b)), min(rel(a), rel(b))
	return (hi + 0.05) / (lo + 0.05)
}

func TestFocusRingIsVisible(t *testing.T) {
	for _, id := range []string{tokens.Light, tokens.Dark, tokens.Sepia} {
		tk := tokens.TokensFor(id)
		if r := wcag20Ratio(tk.FocusRing, tk.WindowBackground); r < 3 {
			t.Errorf("%s: focus ring vs background is only %.2f:1 (need 3:1)", id, r)
		}
	}
}

func TestHighContrastMeetsStricterFloor(t *testing.T) {
	tk := tokens.TokensFor(tokens.HighContrast)
	if r := wcag20Ratio(tk.TextPrimary, tk.WindowBackground); r < 7 {
		t.Errorf("high contrast body text is only %.2f:1 (need 7:1)", r)
	}
	if r := wcag20Ratio(tk.FocusRing, tk.WindowBackground); r < 4.5 {
		t.Errorf("high contrast focus ring is only %.2f:1 (need 4.5:1)", r)
	}
	th := tokens.NewTheme()
	th.SetThemeID(tokens.HighContrast)
	if th.ResolvedTheme() != tokens.HighContrast || th.Tokens().WindowBackground.Hex() != "#000000" || tokens.DisplayName(tokens.HighContrast) != "High contrast" {
		t.Error("the high-contrast theme is not selectable as itself")
	}
}

func TestReducedMotionScale(t *testing.T) {
	th := tokens.NewTheme()
	if th.ReducedMotion() || th.MotionScale() != 1 {
		t.Fatal("a new theme should have motion")
	}
	n := counter(th.OnReducedMotionChanged)
	th.SetReducedMotion(true)
	if *n != 1 || !th.ReducedMotion() || th.MotionScale() != 0 {
		t.Error("turning reduced motion on did not still motion")
	}
	th.SetReducedMotion(true)
	if *n != 1 {
		t.Error("setting the same value again notified")
	}
	store := openStore(t, "s.json")
	a := tokens.NewTheme()
	a.SetSettings(store)
	a.SetReducedMotion(true)
	b := tokens.NewTheme()
	b.SetSettings(store)
	if !b.ReducedMotion() || b.MotionScale() != 0 {
		t.Error("reduced motion did not persist")
	}
}

func TestReducedMotionFollowsTheSystem(t *testing.T) {
	sys := &fakeAppearance{reduced: true}
	th := tokens.NewTheme()
	if th.ReducedMotionSetting() != "system" || th.ReducedMotion() {
		t.Fatal("a new theme should follow the system and, with no system, have motion")
	}
	n := counter(th.OnReducedMotionChanged)
	th.SetAppearance(sys)
	if *n < 1 || !th.ReducedMotion() || th.MotionScale() != 0 {
		t.Error("the system's reduced motion was not followed")
	}
	th.SetReducedMotionSetting("off")
	if th.ReducedMotion() {
		t.Error("an explicit off lost to the system")
	}
	th.SetReducedMotionSetting("on")
	if !th.ReducedMotion() {
		t.Error("an explicit on was not taken")
	}
	th.SetReducedMotionSetting("system")
	if !th.ReducedMotion() {
		t.Error("back to system: should follow the system again")
	}
	sys.set(false, false)
	if th.ReducedMotion() {
		t.Error("the system turning motion back on was not followed")
	}
	th.SetReducedMotionSetting("sometimes")
	if th.ReducedMotionSetting() != "system" {
		t.Error("an invalid setting was taken")
	}
}

func TestHighContrastFollowsTheSystem(t *testing.T) {
	sys := &fakeAppearance{high: true}
	th := tokens.NewTheme()
	th.SetThemeID(tokens.System)
	th.SetAppearance(sys)
	if th.ResolvedTheme() != tokens.HighContrast {
		t.Errorf("system with high contrast on resolved to %q", th.ResolvedTheme())
	}
	th.SetThemeID(tokens.Sepia)
	if th.ResolvedTheme() != tokens.Sepia {
		t.Error("an explicit theme lost to the system's high contrast")
	}
	th.SetThemeID(tokens.System)
	if th.ResolvedTheme() != tokens.HighContrast {
		t.Error("back to system: should be high contrast again")
	}
	sys.set(false, false)
	if th.ResolvedTheme() != tokens.Light {
		t.Errorf("high contrast turned off: resolved to %q", th.ResolvedTheme())
	}
}

func TestAnExistingMotionChoiceSurvivesTheUpgrade(t *testing.T) {
	store := openStore(t, "s.json")
	store.SetValue("view.reducedMotion", true)
	upgraded := tokens.NewTheme()
	upgraded.SetSettings(store)
	if upgraded.ReducedMotionSetting() != "on" || !upgraded.ReducedMotion() {
		t.Error("an old explicit on was not carried over")
	}
	off := openStore(t, "off.json")
	off.SetValue("view.reducedMotion", false)
	sys := &fakeAppearance{reduced: true}
	keptOff := tokens.NewTheme()
	keptOff.SetAppearance(sys)
	keptOff.SetSettings(off)
	if keptOff.ReducedMotionSetting() != "off" || keptOff.ReducedMotion() {
		t.Error("an old explicit off was not carried over")
	}
	fresh := openStore(t, "fresh.json")
	brandNew := tokens.NewTheme()
	brandNew.SetAppearance(sys)
	brandNew.SetSettings(fresh)
	if brandNew.ReducedMotionSetting() != "system" || !brandNew.ReducedMotion() {
		t.Error("a new installation should follow the system")
	}
}

func TestAccentOverride(t *testing.T) {
	th := tokens.NewTheme()
	n := counter(th.OnChanged)
	th.SetAccentOverride("#aa3366")
	tk := th.Tokens()
	if *n != 1 || tk.Accent.Hex() != "#aa3366" || tk.Link.Hex() != "#aa3366" || !tk.TextPrimary.Equal(tokens.TokensFor(tokens.Light).TextPrimary) {
		t.Error("the accent override did not ride on the light table")
	}
	th.SetThemeID(tokens.Dark)
	if th.Tokens().Accent.Hex() != "#aa3366" {
		t.Error("the accent override did not survive a theme switch")
	}
	th.SetAccentOverride("")
	dark := tokens.TokensFor(tokens.Dark)
	if !th.Tokens().Accent.Equal(dark.Accent) || !th.Tokens().Link.Equal(dark.Link) {
		t.Error("clearing the override did not restore dark's accent and link")
	}
}

func TestHighlightOverride(t *testing.T) {
	th := tokens.NewTheme()
	th.SetHighlightOverride("#c2f0c2")
	if th.Tokens().HighlightBackground.Hex() != "#c2f0c2" || !th.Tokens().Accent.Equal(tokens.TokensFor(tokens.Light).Accent) {
		t.Error("the highlight override was not applied on its own")
	}
	th.SetHighlightOverride("")
	if !th.Tokens().HighlightBackground.Equal(tokens.TokensFor(tokens.Light).HighlightBackground) {
		t.Error("clearing the highlight override did not restore the theme's own")
	}
}

func TestInvalidOverrideClears(t *testing.T) {
	th := tokens.NewTheme()
	th.SetAccentOverride("#aa3366")
	th.SetAccentOverride("not-a-color")
	if th.AccentOverride() != "" || !th.Tokens().Accent.Equal(tokens.TokensFor(tokens.Light).Accent) {
		t.Error("an invalid override should clear it")
	}
}

func TestInvalidPersistedOverrideIgnored(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"theme.accent": "not-a-color", "theme.highlight": "#zzzzzz"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	store := settings.New()
	if err := store.Open(path, false); err != nil {
		t.Fatal(err)
	}
	th := tokens.NewTheme()
	th.SetSettings(store)
	if th.AccentOverride() != "" || th.HighlightOverride() != "" {
		t.Error("invalid stored overrides were taken")
	}
	r := tokens.TokensFor(th.ResolvedTheme())
	if !th.Tokens().Accent.Equal(r.Accent) || !th.Tokens().HighlightBackground.Equal(r.HighlightBackground) {
		t.Error("the theme's own accent and highlight should be in force")
	}
}

func TestPersistsThroughSettings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	{
		store := settings.New()
		if err := store.Open(path, false); err != nil {
			t.Fatal(err)
		}
		th := tokens.NewTheme()
		th.SetSettings(store)
		th.SetThemeID(tokens.Sepia)
		th.SetAccentOverride("#aa3366")
		th.SetHighlightOverride("#c2f0c2")
		if err := store.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	reopened := settings.New()
	if err := reopened.Open(path, false); err != nil {
		t.Fatal(err)
	}
	th := tokens.NewTheme()
	th.SetSettings(reopened)
	if th.ThemeID() != tokens.Sepia || th.Tokens().Accent.Hex() != "#aa3366" || th.Tokens().HighlightBackground.Hex() != "#c2f0c2" {
		t.Error("the choices did not survive a reopen")
	}
}

func TestFirstStartDefaultsToSystem(t *testing.T) {
	th := tokens.NewTheme()
	th.SetSettings(openStore(t, "settings.json"))
	if th.ThemeID() != tokens.System || th.ResolvedTheme() != tokens.Light {
		t.Errorf("a first start is %q resolving to %q", th.ThemeID(), th.ResolvedTheme())
	}
}

func TestStaleSettingsValueFallsBack(t *testing.T) {
	store := openStore(t, "settings.json")
	store.SetValue("theme.id", "neon")
	th := tokens.NewTheme()
	th.SetSettings(store)
	if th.ThemeID() != tokens.System {
		t.Errorf("a stale stored theme id gave %q", th.ThemeID())
	}
}
