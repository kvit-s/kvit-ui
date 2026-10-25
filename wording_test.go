package kvitui

import "testing"

// Apple's title case, with the cases the converter has to get right: small
// words inside but not at the ends, words with capitals kept, the access-key
// marker and the ellipsis left where they are.
func TestTitleCase(t *testing.T) {
	for in, want := range map[string]string{
		"Select all":         "Select All",
		"Move to trash…":     "Move to Trash…",
		"Set up":             "Set Up",
		"Save as…":           "Save As…",
		"Don't save":         "Don't Save",
		"Show all columns":   "Show All Columns",
		"Open in browser":    "Open in Browser",
		"Enter full screen":  "Enter Full Screen",
		"Export to PDF":      "Export to PDF",
		"Sync with iCloud":   "Sync With iCloud",
		"Open &folder…":      "Open &Folder…",
		"&save":              "&Save",
		"Built-in themes":    "Built-In Themes",
		"Quit kvit-ui":       "Quit kvit-ui",
		"Search the web":     "Search the Web",
		"A walk in the park": "A Walk in the Park",
		"Copy ui.json":       "Copy ui.json",
		"":                   "",
	} {
		if got := TitleCase(in); got != want {
			t.Errorf("TitleCase(%q) = %q, want %q", in, got, want)
		}
	}
}
