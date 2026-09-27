package palette

import (
	"fmt"
	"math"
)

// RampRules are the thresholds a chart ramp is held to.
type RampRules struct {
	MinContrast            float64 // against every surface the ramp is drawn on
	MinSeparation          float64 // ΔE between any two categorical steps
	MinSeparationDeficient float64 // the same, under each deficiency
	MinChroma              float64 // below this a categorical step reads as a grey
	LightnessBand          bool    // hold categorical steps to a lightness band for the ground
	MinReservedDistance    float64 // ΔE from every colour that already means something
}

// DefaultRampRules are the thresholds kvit-ui's ramps are designed against.
var DefaultRampRules = RampRules{
	MinContrast:            3.0,
	MinSeparation:          8.0,
	MinSeparationDeficient: 8.0,
	MinChroma:              0.09,
	LightnessBand:          true,
	MinReservedDistance:    8.0,
}

// Finding is one broken rule: which, and a sentence saying where.
type Finding struct {
	Rule   string
	Detail string
}

func (f Finding) String() string { return f.Rule + ": " + f.Detail }

// darkGround reports whether a theme's marks are drawn over a dark ground. The
// first surface is the window's own, which is what "the ground" means.
func darkGround(surfaces []Color) bool {
	return len(surfaces) > 0 && RelativeLuminance(surfaces[0]) < 0.2
}

func checkContrast(ramp, surfaces []Color, floor float64, out *[]Finding) {
	for i, c := range ramp {
		for _, s := range surfaces {
			if r := ContrastRatio(c, s); r < floor {
				*out = append(*out, Finding{"contrast", fmt.Sprintf(
					"step %d (%s) reads at %.2f:1 on %s, below the %.1f:1 floor", i, c, r, s, floor)})
			}
		}
	}
}

func checkReserved(ramp, reserved []Color, min float64, out *[]Finding) {
	for i, c := range ramp {
		for _, h := range reserved {
			if d := DeltaE(c, h); d < min {
				*out = append(*out, Finding{"reserved", fmt.Sprintf(
					"step %d (%s) sits ΔE %.1f from the reserved hue %s, inside the %.1f keep-out", i, c, d, h, min)})
			}
		}
	}
}

// ValidateCategorical checks a ramp whose steps name series: every step
// contrasts with every surface, keeps away from the reserved colours, is
// coloured enough not to read as grey, sits in the ground's lightness band,
// and every pair of steps stays apart in ordinary vision and under each
// deficiency. Pairs rather than neighbours, because a legend is read by
// matching a swatch anywhere in the list to a mark anywhere in the chart.
func ValidateCategorical(ramp, surfaces, reserved []Color, rules RampRules) []Finding {
	var out []Finding
	checkContrast(ramp, surfaces, rules.MinContrast, &out)
	checkReserved(ramp, reserved, rules.MinReservedDistance, &out)

	lowest, highest := 0.30, 0.75
	if darkGround(surfaces) {
		lowest, highest = 0.50, 0.98
	}
	for i, c := range ramp {
		o := ToOklch(c)
		if o.C < rules.MinChroma {
			out = append(out, Finding{"chroma", fmt.Sprintf(
				"step %d (%s) has chroma %.3f, below the %.3f floor, so it reads as a grey", i, c, o.C, rules.MinChroma)})
		}
		if rules.LightnessBand && (o.L < lowest || o.L > highest) {
			out = append(out, Finding{"lightness", fmt.Sprintf(
				"step %d (%s) has lightness %.3f, outside the %.2f–%.2f band this ground allows", i, c, o.L, lowest, highest)})
		}
	}
	for i := range ramp {
		for j := i + 1; j < len(ramp); j++ {
			if d := DeltaE(ramp[i], ramp[j]); d < rules.MinSeparation {
				out = append(out, Finding{"separation", fmt.Sprintf(
					"steps %d and %d (%s, %s) are ΔE %.1f apart, under the %.1f floor", i, j, ramp[i], ramp[j], d, rules.MinSeparation)})
			}
			for _, def := range Deficiencies {
				if d := DeltaE(Simulate(ramp[i], def), Simulate(ramp[j], def)); d < rules.MinSeparationDeficient {
					out = append(out, Finding{"separation", fmt.Sprintf(
						"steps %d and %d (%s, %s) converge to ΔE %.1f under %s, under the %.1f floor",
						i, j, ramp[i], ramp[j], d, def, rules.MinSeparationDeficient)})
				}
			}
		}
	}
	return out
}

// ValidateSequential checks a ramp that encodes magnitude: the end away from
// the ground contrasts with every surface, lightness moves one way only, and
// adjacent steps are far enough apart to rank.
func ValidateSequential(ramp, surfaces []Color, rules RampRules) []Finding {
	var out []Finding
	if len(ramp) < 3 {
		return append(out, Finding{"shape", fmt.Sprintf("a sequential ramp of %d steps cannot show a progression", len(ramp))})
	}
	first, last := ramp[0], ramp[len(ramp)-1]
	lf, ll := ToOklch(first).L, ToOklch(last).L
	// Only the end away from the ground is held to the contrast floor: the
	// "least" end is supposed to sit close to the surface.
	far := first
	if darkGround(surfaces) {
		if ll > lf {
			far = last
		}
	} else if ll < lf {
		far = last
	}
	checkContrast([]Color{far}, surfaces, rules.MinContrast, &out)

	ascending := ll > lf
	for i := 1; i < len(ramp); i++ {
		prev, cur := ToOklch(ramp[i-1]).L, ToOklch(ramp[i]).L
		ordered := cur < prev
		if ascending {
			ordered = cur > prev
		}
		if !ordered {
			out = append(out, Finding{"monotonic", fmt.Sprintf(
				"step %d (%s, lightness %.3f) does not continue the progression from step %d (%s, lightness %.3f)",
				i, ramp[i], cur, i-1, ramp[i-1], prev)})
		}
		if step := DeltaE(ramp[i-1], ramp[i]); step < 3.0 {
			out = append(out, Finding{"separation", fmt.Sprintf(
				"steps %d and %d are ΔE %.1f apart, which is not a visible step", i-1, i, step)})
		}
	}
	return out
}

// ValidateDiverging checks a ramp that encodes distance from a reference and
// its direction: an odd number of steps with a neutral midpoint, two arms of
// different hues that stay apart under each deficiency, and each arm a
// sequential ramp of its own.
func ValidateDiverging(ramp, surfaces []Color, rules RampRules) []Finding {
	var out []Finding
	if len(ramp) < 5 || len(ramp)%2 == 0 {
		return append(out, Finding{"shape", fmt.Sprintf(
			"a diverging ramp needs an odd number of steps of at least five so the midpoint is a step; this has %d", len(ramp))})
	}
	middle := len(ramp) / 2
	first, last := ramp[0], ramp[len(ramp)-1]
	checkContrast([]Color{first, last}, surfaces, rules.MinContrast, &out)

	if c := ToOklch(ramp[middle]).C; c > 0.03 {
		out = append(out, Finding{"midpoint", fmt.Sprintf("the midpoint %s has chroma %.3f, so zero is drawn as a colour", ramp[middle], c)})
	}
	if d := DeltaE(first, last); d < rules.MinSeparation*2 {
		out = append(out, Finding{"separation", fmt.Sprintf(
			"the two ends (%s, %s) are ΔE %.1f apart, which does not read as two directions", first, last, d)})
	}
	for _, def := range Deficiencies {
		if d := DeltaE(Simulate(first, def), Simulate(last, def)); d < rules.MinSeparationDeficient {
			out = append(out, Finding{"separation", fmt.Sprintf(
				"the two ends converge to ΔE %.1f under %s, so the sign of the difference is lost", d, def)})
		}
	}
	var low, high []Color
	for i := middle; i >= 0; i-- {
		low = append(low, ramp[i])
	}
	high = append(high, ramp[middle:]...)
	for _, arm := range [][]Color{low, high} {
		for _, f := range ValidateSequential(arm, surfaces, rules) {
			if f.Rule != "contrast" { // already checked at the ends
				out = append(out, f)
			}
		}
	}
	return out
}

// HueDistance is the angle between two hues, from 0 to 180 degrees.
func HueDistance(a, b float64) float64 {
	d := math.Abs(a - b)
	return math.Min(d, 360-d)
}
