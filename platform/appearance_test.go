package platform_test

import (
	"testing"

	"github.com/kvit-s/kvit-ui/platform"
)

func TestSetNotifiesOnlyOnChange(t *testing.T) {
	a := platform.NewAppearance()
	a.Set(false, false, false)
	n := 0
	a.OnChanged(func() { n++ })
	a.Set(false, true, false)
	a.Set(false, true, false)
	a.Set(true, true, true)
	if n != 2 {
		t.Errorf("%d notifications for two changes", n)
	}
	if !a.DarkMode() || !a.HighContrast() || !a.ReducedMotion() || !a.Available() {
		t.Error("the answers set were not kept")
	}
}

// Reading the desktop never fails or hangs, whatever it says.
func TestReadingTheDesktopReturns(t *testing.T) {
	a := platform.NewAppearance()
	a.Refresh()
	t.Logf("dark %v, high contrast %v, reduced motion %v", a.DarkMode(), a.HighContrast(), a.ReducedMotion())
}
