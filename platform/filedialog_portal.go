package platform

import (
	"net/url"
)

// The portal backend of the file dialog: the desktop's FileChooser portal,
// which draws the system's own GTK dialog where neither zenity nor kdialog
// is installed (unison's own dialog there shows neither the title nor the
// filters' names). The D-Bus exchange itself is in filedialog_portal_linux.go;
// this file holds what is portable: the filters, the folder suggestion and
// the chosen URIs as paths.

// portalPattern is one filter pattern as the portal takes it: a glob (0) or
// a MIME type (1). Every pattern here is a glob: "*" matches every file, and
// a bare name such as "Makefile" matches itself.
type portalPattern struct {
	Kind    uint32
	Pattern string
}

// portalFilter is one filter as the portal takes it: what the reader sees
// and the patterns it matches.
type portalFilter struct {
	Name     string
	Patterns []portalPattern
}

// portalFilters are the dialog's name filters as the portal takes them, in
// order.
func portalFilters(filters []nameFilter) []portalFilter {
	out := make([]portalFilter, 0, len(filters))
	for _, f := range filters {
		patterns := make([]portalPattern, 0, len(f.patterns))
		for _, p := range f.patterns {
			patterns = append(patterns, portalPattern{Kind: 0, Pattern: p})
		}
		out = append(out, portalFilter{Name: f.label, Patterns: patterns})
	}
	return out
}

// portalFolder is a folder suggestion as the portal takes it: the path in
// the file system's encoding, NUL-terminated.
func portalFolder(folder string) []byte {
	return append([]byte(folder), 0)
}

// portalURIPaths are chosen URIs as paths: file URIs decoded, anything else
// discarded, since backends normalize to file URIs and drop what they cannot.
func portalURIPaths(uris []string) []string {
	var paths []string
	for _, uri := range uris {
		if p, ok := portalURIPath(uri); ok {
			paths = append(paths, p)
		}
	}
	return paths
}

// portalURIPath is one chosen URI as a path: file URIs only, on this machine.
func portalURIPath(uri string) (string, bool) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme != "file" {
		return "", false
	}
	if u.Host != "" && u.Host != "localhost" {
		return "", false
	}
	if u.Path == "" {
		return "", false
	}
	return u.Path, true
}
