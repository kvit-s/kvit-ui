package platform

import "testing"

// On Windows a folder chooser picks folders, which no file need be in, and
// a file chooser asks for files that are there, several where Multiple
// asks; both take items of the file system only.
func TestTheWindowsDialogPicksFoldersForAFolderChooser(t *testing.T) {
	folders := FileDialog{Directories: true}.dialogOptions(0)
	if folders&fosPickFolders == 0 || folders&fosFileMustExist != 0 || folders&fosForceFileSystem == 0 {
		t.Errorf("a folder chooser's options are %#x", folders)
	}
	files := FileDialog{Multiple: true}.dialogOptions(0)
	if files&fosPickFolders != 0 || files&fosFileMustExist == 0 || files&fosAllowMultiSelect == 0 {
		t.Errorf("a file chooser's options are %#x", files)
	}
}
