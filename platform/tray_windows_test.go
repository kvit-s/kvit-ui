package platform

import (
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

// uxIsDarkModeAllowedForApp is uxtheme.dll's IsDarkModeAllowedForApp, from
// Windows 10 1903 (build 18362), which says whether SetPreferredAppMode let
// the process be drawn dark.
const uxIsDarkModeAllowedForApp = 139

// On Windows 10 1809 and later the tray finds uxtheme's functions for dark
// menus, and from 1903 Windows reports that the process may be drawn dark
// once the tray has asked, so its menu follows Windows' app mode.
func TestTheTrayMenuIsAllowedToBeDark(t *testing.T) {
	v := windows.RtlGetVersion()
	if v.MajorVersion < 10 || v.BuildNumber < darkMenuBuild {
		t.Skipf("Windows %d.%d build %d has no dark menus", v.MajorVersion, v.MinorVersion, v.BuildNumber)
	}
	p := menuTheme()
	if p.refresh == 0 || p.flush == 0 {
		t.Fatalf("uxtheme's functions were not found: %+v", p)
	}
	followMenuTheme()
	if v.BuildNumber < 18362 {
		return
	}
	allowed, err := windows.GetProcAddressByOrdinal(windows.Handle(uxtheme.Handle()), uxIsDarkModeAllowedForApp)
	if allowed == 0 {
		t.Fatalf("IsDarkModeAllowedForApp was not found: %v", err)
	}
	if r, _, _ := syscall.SyscallN(allowed); r&0xFF == 0 {
		t.Error("Windows does not let the process be drawn dark after the tray asked")
	}
}
