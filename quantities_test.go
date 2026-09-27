package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

func TestAFigureSaysWhatWasMeasured(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	plain := kvitui.NewFigure(ui, "1,284.50", "GBP")
	bounded := kvitui.NewFigure(ui, "12", "d")
	bounded.Bounded = true
	unmeasured := kvitui.NewFigure(ui, "0", "d")
	unmeasured.Measured = false
	for f, want := range map[*kvitui.Figure]string{plain: "1,284.50 GBP", bounded: "at most 12 d", unmeasured: "not measured"} {
		if got := f.Phrase(); got != want {
			t.Errorf("%q, want %q", got, want)
		}
	}
	// The unit is dropped with the value: "— d" would say a unit of nothing.
	_, withUnit, _ := unmeasured.Sizes(geom.Size{})
	unmeasured.Unit = ""
	_, without, _ := unmeasured.Sizes(geom.Size{})
	if withUnit != without {
		t.Errorf("an unmeasured figure is %v with a unit and %v without", withUnit, without)
	}
	_ = unison.DefaultMaxSize
}
