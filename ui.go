package kvitui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/platform"
	"github.com/kvit-s/kvit-ui/settings"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/thememode"
	"golang.org/x/text/language"
)

// UI is one application's design values and fonts: the theme, the
// interface metrics and the document typography, read from and saved to the
// settings file, following the desktop's appearance, and applied to unison so
// its own widgets draw in the same colours.
type UI struct {
	Theme      *tokens.Theme
	Interface  *tokens.Interface
	Typography *tokens.Typography
	Fonts      *text.Fonts
	Settings   *settings.Store      // nil when Options.SettingsPath is empty
	Appearance *platform.Appearance // what the desktop says
	// Locale is the language and region numbers are written for; see
	// Number. language.Und groups no digits, as the C locale does.
	Locale   language.Tag
	onChange []func()
	// windows are the open Kvit windows, by the unison window they wrap.
	windows map[*unison.Window]*Window
	// keyTurn is set while a window is handling a key press, so focus that
	// moves then came from the keyboard.
	keyTurn bool
}

// Options configure New.
type Options struct {
	// SettingsPath is the settings file; "" keeps every choice for this run
	// only, which is what a test wants.
	SettingsPath string
	// FontCacheDir is where the index of system fonts is kept; "" uses the
	// user cache directory.
	FontCacheDir string
	// IgnoreDesktop leaves the desktop's appearance and locale out, so
	// "system" is light, motion is on and numbers are written as under the C
	// locale, with no digit grouping. Screenshots and tests want the same
	// result on every machine.
	IgnoreDesktop bool
}

// DefaultSettingsPath is where an application keeps kvit-ui's settings:
// ui.json in the app's folder under the user's configuration directory, the
// file the Kvit versions built with Qt saved them in.
func DefaultSettingsPath(app string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, app, "ui.json")
}

// New sets up the design values and fonts. It can be called before
// unison.Start; the theme is applied to unison's colours straight away and
// again whenever it changes.
func New(opt Options) (*UI, error) {
	fonts, err := text.NewFonts(opt.FontCacheDir)
	if err != nil {
		return nil, err
	}
	if err := fonts.AddFont(icons.Font, icons.FontFamily); err != nil {
		return nil, fmt.Errorf("kvitui: loading the icon font: %w", err)
	}
	u := &UI{
		Theme:      tokens.NewTheme(),
		Interface:  tokens.NewInterface(),
		Typography: tokens.NewTypography(),
		Fonts:      fonts,
		Appearance: platform.NewAppearance(),
	}
	if opt.IgnoreDesktop {
		u.Appearance.Set(false, false, false)
		u.Locale = language.Und
	} else if tag, err := language.Parse(platform.Locale()); err == nil {
		u.Locale = tag
	}
	u.Theme.SetAppearance(u.Appearance)
	if opt.SettingsPath != "" {
		u.Settings = settings.New()
		if err := u.Settings.Open(opt.SettingsPath, false); err != nil {
			return nil, err
		}
		u.Theme.SetSettings(u.Settings)
		u.Interface.SetSettings(u.Settings)
		u.Typography.SetSettings(u.Settings)
	}
	u.Theme.OnChanged(u.changed)
	u.Interface.OnChanged(u.changed)
	u.Typography.OnChanged(u.changed)
	u.applyColors()
	return u, nil
}

// OnChanged calls fn after the theme, the interface size or the typography
// changed and unison has been told; a component that caches anything derived
// from them drops it here.
func (u *UI) OnChanged(fn func()) { u.onChange = append(u.onChange, fn) }

func (u *UI) changed() {
	u.applyColors()
	for _, fn := range u.onChange {
		fn()
	}
	// Relayout and redraw every window: sizes as well as colours may have moved.
	for _, w := range unison.Windows() {
		w.Content().MarkForLayoutRecursively()
		w.MarkForRedraw()
	}
	unison.ThemeChanged()
}

// applyColors writes the theme into unison's theme colours, so unison's own
// fields, menus, dialogs and scroll bars draw in the Kvit theme. unison keeps
// a light and a dark value for each; both get the current theme's value, and
// unison's mode is set to the side whose derived shades suit the ground.
func (u *UI) applyColors() {
	t := u.Theme.Tokens()
	set := func(c *unison.ThemeColor, v palette.Color) { c.Light, c.Dark = Color(v), Color(v) }
	set(unison.ThemeSurface, t.WindowBackground)
	set(unison.ThemeBanding, t.ListBackground)
	set(unison.ThemeFocus, t.Accent)
	set(unison.ThemeTooltip, t.PopupBackground)
	set(unison.ThemeError, t.Danger)
	set(unison.ThemeWarning, t.Warning)
	u.applyMenuTheme()
	switch u.Theme.ResolvedTheme() {
	case tokens.Dark, tokens.HighContrast:
		unison.SetThemeMode(thememode.Dark)
	default:
		unison.SetThemeMode(thememode.Light)
	}
}

// Color converts a design colour for unison.
func Color(c palette.Color) unison.Color {
	r, g, b := c.RGBA8()
	return unison.RGB(int(r), int(g), int(b))
}

// opaque is a colour's alpha with nothing showing through.
const opaque = 255

// TextColor converts a design colour for the text package.
func TextColor(c palette.Color) text.Color {
	r, g, b := c.RGBA8()
	return text.Color{R: r, G: g, B: b, A: opaque}
}

// Chrome is the text style of one of the interface's type roles, at the
// current interface size, in the chrome's family and the given colour.
func (u *UI) Chrome(size int, weight int, color palette.Color) text.Style {
	return text.Style{Family: u.Interface.FontFamily(), Size: float32(size), Weight: weight, Color: TextColor(color)}
}

// Mono is the chrome's fixed-pitch style for identifiers.
func (u *UI) Mono(size int, color palette.Color) text.Style {
	return text.Style{Family: u.Interface.MonoFamily(), Size: float32(size), Color: TextColor(color)}
}

// Icon is the style an icon glyph is drawn in, at a size in pixels.
func (u *UI) Icon(size int, color palette.Color) text.Style {
	return text.Style{Family: icons.FontFamily, Size: float32(size), Color: TextColor(color)}
}
