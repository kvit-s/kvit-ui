package platform

// readLocale reads the region set in System Settings. An application started
// from the Finder is given no LANG, so the environment is only the fallback.
func readLocale() string {
	if v := commandOutput("defaults", "read", "-g", "AppleLocale"); v != "" {
		return posixLocale(v)
	}
	return envLocale()
}
