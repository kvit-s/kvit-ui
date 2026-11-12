package text

// KeptLayouts is how many layouts the store behind Fonts.Layout holds, for
// tests outside the package.
func (f *Fonts) KeptLayouts() int {
	f.cacheMu.Lock()
	defer f.cacheMu.Unlock()
	return len(f.cache)
}
