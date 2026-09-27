package kvitui_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/geom"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

func TestChartsSayWhatTheyShow(t *testing.T) {
	var bar, unmeasured *kvitui.Bar
	var stack *kvitui.StackedBar
	var spark *kvitui.Spark
	var dist *kvitui.Distribution
	var gauge *kvitui.Gauge
	var up, down, flat, unknown *kvitui.Delta
	var tile *kvitui.StatTile
	var block *kvitui.FigureBlock
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		bar = kvitui.NewBar(ui, 62, 100)
		bar.Label, bar.Unit = "Attention", "h"
		unmeasured = kvitui.NewBar(ui, 0, 100)
		unmeasured.Measured = false
		stack = kvitui.NewStackedBar(ui, kvitui.Segment{Value: 420, Label: "rent"}, kvitui.Segment{Value: 180, Label: "groceries"})
		spark = kvitui.NewSpark(ui, 3, kvitui.NotMeasured, 9)
		dist = kvitui.NewDistribution(ui)
		dist.Minimum, dist.LowerQuartile, dist.Median, dist.UpperQuartile, dist.Maximum, dist.Unit = 1, 2, 3, 5, 21, "d"
		gauge = kvitui.NewGauge(ui, 720, 600)
		gauge.Unit = "GBP"
		up, down, flat, unknown = kvitui.NewDelta(ui, 12.4), kvitui.NewDelta(ui, -3.2), kvitui.NewDelta(ui, 0), kvitui.NewDelta(ui, 5)
		up.Unit, up.Precision, down.Unit, down.Precision = "%", 1, "%", 1
		unknown.Measured = false
		tile = kvitui.NewStatTile(ui, "Balance", "4,182.30", "GBP")
		block = kvitui.NewFigureBlock(ui, "1,284", "", "Transactions")
		return []unison.Paneler{kvitui.Column(ui, kvitui.SizeSpace, bar, unmeasured, stack, spark, dist, gauge,
			kvitui.Row(ui, kvitui.SizeSpace, up, down, flat, unknown), kvitui.Left(tile), kvitui.Left(block))}
	})
	for _, c := range []struct {
		p                 unison.Paneler
		role              role.Enum
		name, description string
	}{
		{bar, role.ProgressBar, "Attention", "62 of 100 h"},
		{unmeasured, role.ProgressBar, "", "not measured"},
		{stack, role.Image, "", "420 rent, 180 groceries"},
		{spark, role.Image, "", "3 periods"},
		{dist, role.Image, "", "median 3 d, middle half 2 to 5, range 1 to 21"},
		{gauge, role.ProgressBar, "", "720 of 600 GBP, over by 120"},
		{up, role.Label, "up 12.4 %", ""},
		{down, role.Label, "down 3.2 %", ""},
		{flat, role.Label, "unchanged", ""},
		{unknown, role.Label, "not measured", ""},
		{block, role.Label, "Transactions: 1,284", ""},
	} {
		n := screen.AccessibilityNodeFor(c.p)
		if n == nil || n.Role != c.role || n.Name != c.name || n.Description != c.description {
			t.Errorf("want %v %q %q, got %+v", c.role, c.name, c.description, n)
		}
	}
	if n := screen.AccessibilityNodeFor(tile); n == nil || n.Name != "Balance: 4,182.30 GBP" {
		t.Errorf("the tile's node: %+v", n)
	}
	screen.Do(func() {
		if w := tile.FrameRect().Width; w < 200 {
			t.Errorf("a stat tile is %.1f wide, under 200", w)
		}
	})
}

// A trend with nothing to plot says so, rather than drawing an empty axis.
func TestATrendSaysWhenItIsEmpty(t *testing.T) {
	var full, empty *kvitui.Trend
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		full = kvitui.NewTrend(ui, 400, 620, 580)
		full.Label, full.MaximumY = "Balance", 2000
		empty = kvitui.NewTrend(ui)
		empty.Label = "Savings"
		return []unison.Paneler{full, empty}
	})
	screen.Do(func() {
		if full.Children()[0].Hidden != true || empty.Children()[0].Hidden {
			t.Error("the empty state shows on the full trend, or not on the empty one")
		}
	})
	if n := screen.AccessibilityNodeFor(full); n == nil || n.Description != "3 points, 0 to 2000" {
		t.Errorf("the trend's node: %+v", n)
	}
}

// A Qt item draws past its own box and unison clips each panel to its box, so
// what a Kvit component draws outside it, like a gauge's pace tick standing
// out of the bar, is drawn by the window over the page.
func TestAGaugesPaceTickStandsOutOfTheBar(t *testing.T) {
	var gauge *kvitui.Gauge
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		gauge = kvitui.NewGauge(ui, 100, 600)
		gauge.Pace = 0.5
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(300), gauge)}
	})
	var box geom.Rect
	screen.Do(func() { box = gauge.RectToRoot(gauge.ContentRect(true)) })
	img := screen.Capture()
	want := kvitui.Color(ui.Theme.Tokens().TextPrimary)
	x := int(box.X + box.Width/2)
	for _, y := range []int{int(box.Y) - 1, int(box.Bottom())} {
		found := false
		for dx := -1; dx <= 1; dx++ {
			r, g, b, _ := img.At(x+dx, y).RGBA()
			found = found || (r>>8 == uint32(want.Red()) && g>>8 == uint32(want.Green()) && b>>8 == uint32(want.Blue()))
		}
		if !found {
			t.Errorf("no pace tick at row %d, outside the gauge's box %v", y, box)
		}
	}
}

func TestANetFlowSaysWhatItShows(t *testing.T) {
	var flow *kvitui.NetFlow
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		flow = kvitui.NewNetFlow(ui, []float64{0.6, kvitui.NotMeasured, 0.2}, []float64{0.4, kvitui.NotMeasured, 0.9}, "Jan", "Feb", "Mar")
		flow.Label = "Cash flow by month"
		flow.RisingLabel, flow.FallingLabel = "Money in", "Money out"
		return []unison.Paneler{kvitui.FullWidth(flow)}
	})
	if n := screen.AccessibilityNodeFor(flow); n == nil || n.Role != role.Image || n.Name != "Cash flow by month" ||
		n.Description != "3 periods, what came in above the line and what went out below it" {
		t.Errorf("the chart's node is %+v", n)
	}
}
