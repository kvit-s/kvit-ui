package palette_test

import (
	"math"
	"strings"
	"testing"

	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/tokens"
)

// Ported from kvit-ui's tests/test_palette.cpp.

func TestOklchRoundTrips(t *testing.T) {
	for _, h := range []string{"#000000", "#ffffff", "#c0392b", "#2970c8", "#e0a34c", "#6f7178", "#971c4b"} {
		c := palette.Hex(h)
		back := palette.FromOklch(palette.ToOklch(c))
		if !back.Equal(c) {
			t.Errorf("%s came back from OKLCH as %s", h, back)
		}
	}
	// A value outside sRGB is reported as such.
	if palette.InGamut(palette.Oklch{L: 0.7, C: 0.4, H: 150}) {
		t.Error("chroma 0.4 at hue 150 should be outside sRGB")
	}
	if !palette.InGamut(palette.ToOklch(palette.Hex("#2970c8"))) {
		t.Error("an sRGB colour should be inside sRGB")
	}
}

func TestDeltaEIsZeroForIdenticalColours(t *testing.T) {
	c := palette.Hex("#4a90d9")
	if d := palette.DeltaE(c, c); d != 0 {
		t.Errorf("ΔE of a colour with itself is %v", d)
	}
	if d := palette.DeltaE(palette.Hex("#000000"), palette.Hex("#ffffff")); math.Abs(d-100) > 0.01 {
		t.Errorf("black to white should be ΔE 100, got %v", d)
	}
}

// A deuteranope loses the red-green axis, so a saturated red and green must
// come out much closer than they went in; without this a validator could be
// running a near-identity matrix and reporting that everything is fine.
func TestDeficiencySimulationCollapsesTheAxisItShould(t *testing.T) {
	red, green := palette.Hex("#c0392b"), palette.Hex("#1e874b")
	ordinary := palette.DeltaE(red, green)
	seen := palette.DeltaE(palette.Simulate(red, palette.Deuteranopia), palette.Simulate(green, palette.Deuteranopia))
	if seen >= ordinary*0.6 {
		t.Errorf("red/green were ΔE %.1f apart and are ΔE %.1f under deuteranopia", ordinary, seen)
	}
	// The blue-yellow axis survives it.
	by := palette.DeltaE(palette.Simulate(palette.Hex("#2970c8"), palette.Deuteranopia),
		palette.Simulate(palette.Hex("#e0a34c"), palette.Deuteranopia))
	if by <= 30 {
		t.Errorf("blue/yellow should stay apart under deuteranopia, got ΔE %.1f", by)
	}
}

func describe(fs []palette.Finding) string {
	var b strings.Builder
	for _, f := range fs {
		b.WriteString("  " + f.String() + "\n")
	}
	return b.String()
}

func TestCategoricalRampHolds(t *testing.T) {
	for _, theme := range tokens.BuiltInThemes() {
		tk := tokens.TokensFor(theme)
		if len(tk.CategoricalRamp) != 8 {
			t.Errorf("%s: categorical ramp has %d steps, want 8", theme, len(tk.CategoricalRamp))
		}
		if fs := palette.ValidateCategorical(tk.CategoricalRamp, tk.Surfaces(), tk.ReservedHues(), palette.DefaultRampRules); len(fs) > 0 {
			t.Errorf("the %s categorical ramp fails:\n%s", theme, describe(fs))
		}
	}
}

func TestSequentialRampHolds(t *testing.T) {
	for _, theme := range tokens.BuiltInThemes() {
		tk := tokens.TokensFor(theme)
		if len(tk.SequentialRamp) != 7 {
			t.Errorf("%s: sequential ramp has %d steps, want 7", theme, len(tk.SequentialRamp))
		}
		if fs := palette.ValidateSequential(tk.SequentialRamp, tk.Surfaces(), palette.DefaultRampRules); len(fs) > 0 {
			t.Errorf("the %s sequential ramp fails:\n%s", theme, describe(fs))
		}
	}
}

func TestDivergingRampHolds(t *testing.T) {
	for _, theme := range tokens.BuiltInThemes() {
		tk := tokens.TokensFor(theme)
		if len(tk.DivergingRamp) != 7 {
			t.Errorf("%s: diverging ramp has %d steps, want 7", theme, len(tk.DivergingRamp))
		}
		if fs := palette.ValidateDiverging(tk.DivergingRamp, tk.Surfaces(), palette.DefaultRampRules); len(fs) > 0 {
			t.Errorf("the %s diverging ramp fails:\n%s", theme, describe(fs))
		}
	}
}

// Step n is the same colour in every theme, allowing for how much lighter and
// more saturated it has to be on a dark ground: each theme may sit up to 14
// degrees either side of a shared spoke, so two themes can be about 30 apart.
func TestRampsShareTheirHuesAcrossThemes(t *testing.T) {
	for step := 0; step < 8; step++ {
		ref := -1.0
		for _, theme := range tokens.BuiltInThemes() {
			hue := palette.ToOklch(tokens.TokensFor(theme).CategoricalRamp[step]).H
			if ref < 0 {
				ref = hue
				continue
			}
			if d := palette.HueDistance(hue, ref); d > 30 {
				t.Errorf("categorical step %d is hue %.0f in %s against %.0f elsewhere, a drift of %.0f degrees", step, hue, theme, ref, d)
			}
		}
	}
}

// The validators are not vacuous: each rejects a ramp built to break it.
func TestTheValidatorsRejectBrokenRamps(t *testing.T) {
	light := tokens.TokensFor(tokens.Light)
	grey := []palette.Color{palette.Hex("#777777"), palette.Hex("#787878")}
	if fs := palette.ValidateCategorical(grey, light.Surfaces(), nil, palette.DefaultRampRules); len(fs) == 0 {
		t.Error("two near-identical greys passed as a categorical ramp")
	}
	backwards := []palette.Color{palette.Hex("#5e2142"), palette.Hex("#f5aacd"), palette.Hex("#9b3f70")}
	if fs := palette.ValidateSequential(backwards, light.Surfaces(), palette.DefaultRampRules); len(fs) == 0 {
		t.Error("a non-monotonic ramp passed as sequential")
	}
	even := light.DivergingRamp[:6]
	if fs := palette.ValidateDiverging(even, light.Surfaces(), palette.DefaultRampRules); len(fs) == 0 {
		t.Error("a six-step ramp passed as diverging")
	}
}
