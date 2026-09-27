package platform

import "strings"

// readPlatform reads GNOME's keys, which KDE and the other desktops following
// the freedesktop settings conventions also honour.
func readPlatform() (dark, high, reduced, available bool) {
	scheme := commandOutput("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	if scheme == "" {
		// Older desktops name the theme instead of a preference.
		scheme = commandOutput("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	}
	dark = strings.Contains(strings.ToLower(scheme), "dark")
	high = commandOutput("gsettings", "get", "org.gnome.desktop.a11y.interface", "high-contrast") == "true"
	// enable-animations reads the other way round.
	reduced = commandOutput("gsettings", "get", "org.gnome.desktop.interface", "enable-animations") == "false"
	return dark, high, reduced, true
}
