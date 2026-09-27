//go:build !darwin && !windows

package platform

// readLocale reads the locale variables, which is where Linux desktops put
// the region settings for the programs they start.
func readLocale() string { return envLocale() }
