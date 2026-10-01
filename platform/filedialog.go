package platform

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/richardwilkes/unison"
)

// FileDialog is the system's own dialog for choosing files to open, with a
// title and named filters. unison's OpenDialog has neither: it cannot be
// given a title, and it names its filters itself from bare extensions ("All
// Readable Files", "sqlite Files", "All Files" on Windows). This one is shown
// with the system's own API on each desktop, as unison's is:
//   - Windows: the common item dialog (IFileOpenDialog), with SetTitle and
//     SetFileTypes;
//   - macOS: NSOpenPanel, with the title as its message, since a panel shows
//     no title bar, and the filters in a pop-up below the file list;
//   - Linux: kdialog in a KDE session and zenity elsewhere, as unison
//     chooses; with neither installed, the desktop's portal FileChooser,
//     which draws the same GTK dialog, except under Windows Subsystem for
//     Linux, where the portal is skipped; with no portal either, unison's
//     own dialog, which shows neither the title nor the filters' names.
//
// In a headless session (tests) it is unison's own dialog, drawn on the
// headless screen, so a test never opens a window on the desktop.
type FileDialog struct {
	// Title is what the dialog says it is for, such as "Open a Kvit Cash
	// workspace".
	Title string
	// NameFilters are the filters the reader chooses among, the first in
	// force when the dialog opens. Each is a name with its patterns in
	// parentheses, "Kvit Cash workspace (*.sqlite)", or patterns alone,
	// "*.csv *.tsv". "*" matches every file.
	// None offers every file.
	NameFilters []string
	// Folder is where the dialog opens; unset, the system's choice, which is
	// usually the folder a file was last chosen from.
	Folder string
	// Multiple lets the reader choose several files.
	Multiple bool
	// Directories picks folders instead of files, such as the Add project
	// chooser. The filters do not apply to it.
	Directories bool
}

// Open shows the dialog over the active window and waits for the reader,
// with the application's windows taking no input meanwhile. It returns the
// files chosen, or none when the reader cancelled.
func (d FileDialog) Open() []string {
	if unison.ActiveHeadlessScreen() != nil {
		return d.openWithUnison()
	}
	return d.openNative()
}

// openWithUnison is unison's own dialog, which takes the filters'
// extensions and neither their names nor the title.
func (d FileDialog) openWithUnison() []string {
	picker := unison.NewOpenDialog()
	picker.SetAllowsMultipleSelection(d.Multiple)
	picker.SetCanChooseFiles(!d.Directories)
	picker.SetCanChooseDirectories(d.Directories)
	if d.Folder != "" {
		picker.SetInitialDirectory(d.Folder)
	}
	if !d.Directories {
		picker.SetAllowedExtensions(d.extensions()...)
	}
	if !picker.RunModal() {
		return nil
	}
	return picker.Paths()
}

// nameFilter is one of the reader's choices: what it is called and the
// patterns it matches.
type nameFilter struct {
	label    string
	patterns []string
}

// nameFilterPattern reads a name filter: a name, then the patterns in the last
// parentheses.
var nameFilterPattern = regexp.MustCompile(`^(.*)\(([a-zA-Z0-9_.,*? +;#\-\[\]@\{\}/!<>\$%&=^~:\|]*)\)$`)

// filters are the dialog's name filters as labels and patterns. The label
// is the whole string, parentheses and all, and is what every system shows.
// A filter with no patterns matches every file.
func (d FileDialog) filters() []nameFilter {
	out := make([]nameFilter, 0, len(d.NameFilters))
	for _, f := range d.NameFilters {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		patterns := f
		if m := nameFilterPattern.FindStringSubmatch(f); m != nil {
			patterns = m[2]
		}
		fields := strings.Fields(patterns)
		if len(fields) == 0 {
			fields = []string{"*"}
		}
		out = append(out, nameFilter{label: f, patterns: fields})
	}
	return out
}

// matchesEverything reports a filter that lets every file through.
func (f nameFilter) matchesEverything() bool {
	return slices.ContainsFunc(f.patterns, func(p string) bool { return p == "*" || p == "*.*" })
}

// extensions are the filters as unison takes them, in order and once each:
// "sqlite" for "*.sqlite", and "*" for a filter that matches every file or
// has a pattern that is not an extension, so that such a file can still be
// chosen.
func (d FileDialog) extensions() []string {
	var out []string
	add := func(ext string) {
		if !slices.Contains(out, ext) {
			out = append(out, ext)
		}
	}
	for _, f := range d.filters() {
		for _, p := range f.patterns {
			if ext, ok := strings.CutPrefix(p, "*."); ok && ext != "*" && !strings.ContainsAny(ext, "*?[") {
				add(ext)
			} else {
				add("*")
			}
		}
	}
	return out
}

// fileTypes are the filters as Windows takes them: a name and the patterns
// joined by semicolons. Windows adds " (patterns)" to a name that does not
// hold the first pattern's extension, such as ".sqlite" or ".*" (seen on
// Windows 11, 2026-09-28), and a bare "*" has none, so "All files (*)"
// would read "All files (*) (*)". A filter whose one pattern is "*" is
// handed over without its " (*)", which Windows then puts back. With no
// filters, one that matches every file.
func (d FileDialog) fileTypes() [][2]string {
	filters := d.filters()
	if len(filters) == 0 {
		return [][2]string{{"All files", "*"}}
	}
	out := make([][2]string, len(filters))
	for i, f := range filters {
		name := f.label
		if slices.Equal(f.patterns, []string{"*"}) {
			if trimmed := strings.TrimSuffix(name, " (*)"); trimmed != "" {
				name = trimmed
			}
		}
		out[i] = [2]string{name, strings.Join(f.patterns, ";")}
	}
	return out
}

// zenityArgs are zenity's arguments for the dialog; with several files
// chosen, zenity writes one to a line. A folder picker takes --directory
// and no filters.
func (d FileDialog) zenityArgs() []string {
	args := []string{"--file-selection"}
	if d.Title != "" {
		args = append(args, "--title="+d.Title)
	}
	if d.Folder != "" {
		args = append(args, "--filename="+withSeparator(d.Folder))
	}
	if d.Multiple {
		args = append(args, "--multiple", "--separator=\n")
	}
	if d.Directories {
		args = append(args, "--directory")
		return args
	}
	for _, f := range d.filters() {
		args = append(args, "--file-filter="+f.label+" | "+strings.Join(f.patterns, " "))
	}
	return args
}

// kdialogArgs are kdialog's arguments for the dialog. kdialog takes the
// filters as KDE writes them, one to a line, "*.sqlite|Kvit Cash workspace
// (*.sqlite)", and with several files chosen writes one to a line.
func (d FileDialog) kdialogArgs() []string {
	var args []string
	if d.Title != "" {
		args = append(args, "--title", d.Title)
	}
	if d.Directories {
		args = append(args, "--getexistingdirectory")
	} else {
		args = append(args, "--getopenfilename")
	}
	if d.Multiple {
		args = append(args, "--multiple", "--separate-output")
	}
	folder := d.Folder
	if folder == "" {
		folder = "."
	}
	args = append(args, withSeparator(folder))
	if d.Directories {
		return args
	}
	if filters := d.filters(); len(filters) > 0 {
		lines := make([]string, len(filters))
		for i, f := range filters {
			lines[i] = strings.Join(f.patterns, " ") + "|" + f.label
		}
		args = append(args, strings.Join(lines, "\n"))
	}
	return args
}

// withSeparator ends a folder's path with a separator, which is how zenity
// and kdialog tell a folder to open in from a file to choose.
func withSeparator(dir string) string {
	if strings.HasSuffix(dir, string(filepath.Separator)) {
		return dir
	}
	return dir + string(filepath.Separator)
}

// chosenPaths are the lines an external dialog wrote, one file to a line.
func chosenPaths(out []byte, multiple bool) []string {
	var paths []string
	for line := range strings.SplitSeq(strings.TrimRight(string(out), "\r\n"), "\n") {
		if line = strings.TrimRight(line, "\r"); line != "" {
			paths = append(paths, line)
		}
	}
	if !multiple && len(paths) > 1 {
		paths = paths[:1]
	}
	return paths
}
