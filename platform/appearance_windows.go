package platform

import (
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	user32                   = windows.NewLazySystemDLL("user32.dll")
	procSystemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

const (
	spiGetHighContrast        = 0x0042
	spiGetClientAreaAnimation = 0x1042
	hcfHighContrastOn         = 0x00000001
	personalizeKey            = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	personalizeAppsLightValue = "AppsUseLightTheme"
)

type highContrast struct {
	cbSize            uint32
	dwFlags           uint32
	lpszDefaultScheme uintptr
}

// readPlatform asks Windows through SystemParametersInfo and the registry.
func readPlatform() (dark, high, reduced, available bool) {
	if k, err := registry.OpenKey(registry.CURRENT_USER, personalizeKey, registry.QUERY_VALUE); err == nil {
		if v, _, err := k.GetIntegerValue(personalizeAppsLightValue); err == nil {
			dark = v == 0
		}
		k.Close()
	}
	hc := highContrast{cbSize: uint32(unsafe.Sizeof(highContrast{}))}
	if r, _, _ := procSystemParametersInfo.Call(spiGetHighContrast, uintptr(hc.cbSize), uintptr(unsafe.Pointer(&hc)), 0); r != 0 {
		high = hc.dwFlags&hcfHighContrastOn != 0
	}
	// SPI_GETCLIENTAREAANIMATION answers the other way round: true means
	// animations are wanted.
	var animations int32 = 1
	if r, _, _ := procSystemParametersInfo.Call(spiGetClientAreaAnimation, 0, uintptr(unsafe.Pointer(&animations)), 0); r != 0 {
		reduced = animations == 0
	}
	return dark, high, reduced, true
}
