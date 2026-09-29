package platform

import (
	"slices"
	"testing"
)

// Kvit Cash's workspace picker, as its  FileDialog was written.
var workspacePicker = FileDialog{
	Title:       "Open a Kvit Cash workspace",
	NameFilters: []string{"Kvit Cash workspace (*.sqlite)", "All files (*)"},
	Folder:      "/home/reader/finance",
}

// A name filter is read as  reads it: the whole string is what the reader
// sees, and the patterns are what is in its last parentheses, or the whole
// string when it has none.
func TestNameFiltersAreReadAsQtReadsThem(t *testing.T) {
	d := FileDialog{NameFilters: []string{"Kvit Cash workspace (*.sqlite)", "All files (*)", "*.csv *.tsv",
		"Statements (bank (2024)) (*.ofx *.qfx)", " ", "Nothing ()"}}
	want := []nameFilter{
		{"Kvit Cash workspace (*.sqlite)", []string{"*.sqlite"}},
		{"All files (*)", []string{"*"}},
		{"*.csv *.tsv", []string{"*.csv", "*.tsv"}},
		{"Statements (bank (2024)) (*.ofx *.qfx)", []string{"*.ofx", "*.qfx"}},
		{"Nothing ()", []string{"*"}},
	}
	got := d.filters()
	if len(got) != len(want) {
		t.Fatalf("read %+v", got)
	}
	for i := range want {
		if got[i].label != want[i].label || !slices.Equal(got[i].patterns, want[i].patterns) {
			t.Errorf("filter %d is %+v, want %+v", i, got[i], want[i])
		}
	}
}

// Each system is handed the title and the filters under their own names,
// the first filter first.
func TestEachSystemIsGivenTheTitleAndTheNamedFilters(t *testing.T) {
	// Windows adds " (*)" to a name whose pattern is "*" itself.
	if got, want := workspacePicker.fileTypes(), [][2]string{
		{"Kvit Cash workspace (*.sqlite)", "*.sqlite"}, {"All files", "*"},
	}; !slices.Equal(got, want) {
		t.Errorf("Windows is given %q", got)
	}
	if got, want := workspacePicker.zenityArgs(), []string{
		"--file-selection", "--title=Open a Kvit Cash workspace", "--filename=/home/reader/finance/",
		"--file-filter=Kvit Cash workspace (*.sqlite) | *.sqlite", "--file-filter=All files (*) | *",
	}; !slices.Equal(got, want) {
		t.Errorf("zenity is given %q", got)
	}
	if got, want := workspacePicker.kdialogArgs(), []string{
		"--title", "Open a Kvit Cash workspace", "--getopenfilename", "/home/reader/finance/",
		"*.sqlite|Kvit Cash workspace (*.sqlite)\n*|All files (*)",
	}; !slices.Equal(got, want) {
		t.Errorf("kdialog is given %q", got)
	}
	several := FileDialog{Multiple: true, NameFilters: []string{"*.csv *.tsv"}}
	if got, want := several.zenityArgs(), []string{"--file-selection", "--multiple", "--separator=\n",
		"--file-filter=*.csv *.tsv | *.csv *.tsv"}; !slices.Equal(got, want) {
		t.Errorf("zenity is given %q for several files", got)
	}
	if got, want := several.kdialogArgs(), []string{"--getopenfilename", "--multiple", "--separate-output", "./",
		"*.csv *.tsv|*.csv *.tsv"}; !slices.Equal(got, want) {
		t.Errorf("kdialog is given %q for several files", got)
	}
	if got := (FileDialog{}).fileTypes(); !slices.Equal(got, [][2]string{{"All files", "*"}}) {
		t.Errorf("Windows is given %q with no filters", got)
	}
}

// unison's own dialog takes extensions, once each; a pattern that is not an
// extension becomes every file, so that its files can still be chosen.
func TestUnisonsDialogIsGivenTheFiltersExtensions(t *testing.T) {
	if got := workspacePicker.extensions(); !slices.Equal(got, []string{"sqlite", "*"}) {
		t.Errorf("the workspace picker's extensions are %q", got)
	}
	d := FileDialog{NameFilters: []string{"Tables (*.csv *.tsv)", "Comma-separated (*.csv)", "Makefiles (Makefile)"}}
	if got := d.extensions(); !slices.Equal(got, []string{"csv", "tsv", "*"}) {
		t.Errorf("the extensions are %q", got)
	}
	if got := (FileDialog{}).extensions(); got != nil {
		t.Errorf("with no filters the extensions are %q", got)
	}
}

// An external dialog writes the files chosen one to a line; with one file
// asked for, only the first counts.
func TestTheFilesAnExternalDialogWritesAreRead(t *testing.T) {
	out := []byte("/a/one.sqlite\n/a/two words.sqlite\r\n\n")
	if got := chosenPaths(out, true); !slices.Equal(got, []string{"/a/one.sqlite", "/a/two words.sqlite"}) {
		t.Errorf("several files read as %q", got)
	}
	if got := chosenPaths(out, false); !slices.Equal(got, []string{"/a/one.sqlite"}) {
		t.Errorf("one file read as %q", got)
	}
	if got := chosenPaths([]byte("\n"), false); got != nil {
		t.Errorf("nothing read as %q", got)
	}
}
