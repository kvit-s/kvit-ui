package platform

import (
	"reflect"
	"slices"
	"testing"
)

// A folder picker asks the system for folders, with its title and no file
// filters: zenity takes --directory, kdialog --getexistingdirectory.
func TestAFolderPickerAsksForFoldersWithoutFilters(t *testing.T) {
	d := FileDialog{Title: "Add project", Folder: "/home/reader", Directories: true}
	if got, want := d.zenityArgs(), []string{
		"--file-selection", "--title=Add project", "--filename=/home/reader/", "--directory",
	}; !slices.Equal(got, want) {
		t.Errorf("zenity is given %q", got)
	}
	if got, want := d.kdialogArgs(), []string{
		"--title", "Add project", "--getexistingdirectory", "/home/reader/",
	}; !slices.Equal(got, want) {
		t.Errorf("kdialog is given %q", got)
	}
	several := FileDialog{Directories: true, Multiple: true}
	if got, want := several.zenityArgs(), []string{
		"--file-selection", "--multiple", "--separator=\n", "--directory",
	}; !slices.Equal(got, want) {
		t.Errorf("zenity is given %q for several folders", got)
	}
	if got, want := several.kdialogArgs(), []string{
		"--getexistingdirectory", "--multiple", "--separate-output", "./",
	}; !slices.Equal(got, want) {
		t.Errorf("kdialog is given %q for several folders", got)
	}
}

// The portal takes the filters as named globs, in order.
func TestThePortalIsGivenTheNamedFilters(t *testing.T) {
	got := portalFilters(workspacePicker.filters())
	want := []portalFilter{
		{"Kvit Cash workspace (*.sqlite)", []portalPattern{{0, "*.sqlite"}}},
		{"All files (*)", []portalPattern{{0, "*"}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the portal is given %+v", got)
	}
	if got := portalFilters(nil); len(got) != 0 {
		t.Errorf("with no filters the portal is given %+v", got)
	}
}

// The portal takes a folder suggestion NUL-terminated.
func TestThePortalIsGivenTheFolderNulTerminated(t *testing.T) {
	if got, want := portalFolder("/home/reader/finance"), []byte("/home/reader/finance\x00"); !reflect.DeepEqual(got, want) {
		t.Errorf("the portal is given %q", got)
	}
}

// Chosen URIs come back as paths: file URIs decoded, anything else dropped.
func TestChosenUrisAreReadAsPaths(t *testing.T) {
	uris := []string{
		"file:///home/reader/My%20Project",
		"file://localhost/home/reader/other",
		"https://example.com/remote",
		"file:///home/reader/with%2Fslash",
		"",
	}
	if got, want := portalURIPaths(uris), []string{
		"/home/reader/My Project", "/home/reader/other", "/home/reader/with/slash",
	}; !slices.Equal(got, want) {
		t.Errorf("chosen URIs read as %q", got)
	}
	if got := portalURIPaths(nil); len(got) != 0 {
		t.Errorf("nothing chosen reads as %q", got)
	}
	if _, ok := portalURIPath("file://remotehost/home/reader"); ok {
		t.Error("another machine's file counts as chosen")
	}
}
