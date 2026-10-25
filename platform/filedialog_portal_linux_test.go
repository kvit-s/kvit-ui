package platform

import "testing"

// Under Windows Subsystem for Linux the portal is skipped, since it shows
// the Windows files while the program works with the Linux ones: the Add
// project picker falls back to the dialog that opens where the program is.
func TestThePortalIsSkippedOnWsl(t *testing.T) {
	t.Setenv("WSL_DISTRO_NAME", "ubuntu")
	t.Setenv("WSL_INTEROP", "/run/WSL/interop")
	if !isWSL() {
		t.Fatal("a WSL session is not recognized")
	}
	if _, ok := portalConnection(true); ok {
		t.Error("the portal is used on WSL")
	}
}
