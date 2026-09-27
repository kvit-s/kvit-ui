package tokens

import "math"

// The interface size's clamps: 10 px is the smallest the chrome stays legible
// at, and 24 px is where a list row stops fitting a title and a date on one
// line.
const (
	DefaultInterfaceSize = 12
	MinInterfaceSize     = 10
	MaxInterfaceSize     = 24
)

const (
	keyInterfaceSize = "interface.fontSize"
	keyFontFamily    = "interface.fontFamily"
	keyMonoFamily    = "interface.monoFamily"
)

// Interface is the chrome's type scale and geometry: every size, gap, height,
// radius and reflow width the library's components use, all derived from one
// interface size through Px. Type and geometry travel together because a
// 24 px label inside a 28 px button clips. Document text is Typography's.
type Interface struct {
	settings   Settings
	loading    bool
	fontSize   int
	fontFamily string
	monoFamily string
	changed    signal
}

// NewInterface returns the metrics at the default size of 12 px.
func NewInterface() *Interface {
	return &Interface{fontSize: DefaultInterfaceSize, monoFamily: "monospace"}
}

// OnChanged calls fn whenever any value may have changed.
func (m *Interface) OnChanged(fn func()) (disconnect func()) { return m.changed.connect(fn) }

// SetSettings reads the stored values, clamped as a live change would be,
// and saves every later change back.
func (m *Interface) SetSettings(s Settings) {
	m.settings = s
	if s == nil {
		return
	}
	m.loading = true
	m.SetFontSize(int(math.Round(numberValue(s, keyInterfaceSize, float64(m.fontSize)))))
	m.SetFontFamily(stringValue(s, keyFontFamily, m.fontFamily))
	m.SetMonoFamily(stringValue(s, keyMonoFamily, m.monoFamily))
	m.loading = false
	m.changed.emit()
}

// FontSize is the interface size in pixels, 10 to 24.
func (m *Interface) FontSize() int { return m.fontSize }

// SetFontSize sets the interface size, clamped to 10–24.
func (m *Interface) SetFontSize(size int) {
	size = min(max(size, MinInterfaceSize), MaxInterfaceSize)
	if m.fontSize == size {
		return
	}
	m.fontSize = size
	m.save(keyInterfaceSize, size)
	m.changed.emit()
}

// FontFamily is the chrome's chosen family; "" means the desktop's own.
func (m *Interface) FontFamily() string { return m.fontFamily }

// SetFontFamily chooses the chrome's family; "" follows the desktop.
func (m *Interface) SetFontFamily(family string) {
	if m.fontFamily == family {
		return
	}
	m.fontFamily = family
	m.save(keyFontFamily, family)
	m.changed.emit()
}

// MonoFamily is the family identifiers are set in; "monospace" means the
// platform's own fixed-pitch face.
func (m *Interface) MonoFamily() string { return m.monoFamily }

// SetMonoFamily chooses the identifier family; "" means "monospace", since
// an identifier in a proportional face no longer lines up down a column.
func (m *Interface) SetMonoFamily(family string) {
	if family == "" {
		family = "monospace"
	}
	if m.monoFamily == family {
		return
	}
	m.monoFamily = family
	m.save(keyMonoFamily, family)
	m.changed.emit()
}

// ResetToDefaults returns to 12 px and the desktop's families.
func (m *Interface) ResetToDefaults() {
	m.SetFontSize(DefaultInterfaceSize)
	m.SetFontFamily("")
	m.SetMonoFamily("monospace")
}

// Scale is FontSize / 12, what Px multiplies by.
func (m *Interface) Scale() float64 { return float64(m.fontSize) / DefaultInterfaceSize }

// Px scales a design-pixel value to the interface size and rounds it. A
// non-zero value never rounds to zero, so a one-pixel rule stays visible at
// the smallest size. Reach for a named value first: Px with a literal in it
// is right for a one-off and wrong for anything a second view also needs.
func (m *Interface) Px(design int) int {
	if design == 0 {
		return 0
	}
	scaled := int(math.Round(float64(design) * m.Scale()))
	if design > 0 {
		return max(1, scaled)
	}
	return min(-1, scaled)
}

func (m *Interface) save(key string, value any) {
	if m.settings != nil && !m.loading {
		m.settings.SetValue(key, value)
	}
}

// The seven type roles, at their pixel sizes for the default 12 px.
func (m *Interface) Caption() int  { return m.Px(10) } // kind tags, counts
func (m *Interface) Small() int    { return m.Px(11) } // chip labels, sub-lines
func (m *Interface) Body() int     { return m.Px(12) } // row text, prose
func (m *Interface) Strong() int   { return m.Px(13) } // a name, an emphasised row
func (m *Interface) Title() int    { return m.Px(15) } // a section heading
func (m *Interface) Headline() int { return m.Px(17) } // a pane title, the wordmark
func (m *Interface) Display() int  { return m.Px(20) } // a page title

// The spacing scale. 3 and 5 are deliberately absent: a scale with every
// integer in it is not a scale.
func (m *Interface) SpaceTight() int { return m.Px(2) }
func (m *Interface) SpaceSnug() int  { return m.Px(4) }
func (m *Interface) SpaceNear() int  { return m.Px(6) }
func (m *Interface) Space() int      { return m.Px(8) } // the ordinary gap
func (m *Interface) SpaceWide() int  { return m.Px(10) }
func (m *Interface) SpaceLoose() int { return m.Px(12) }

// Layout.
func (m *Interface) ViewMargin() int       { return m.Px(16) }
func (m *Interface) ColumnGap() int        { return m.Px(14) }
func (m *Interface) StackGap() int         { return m.Px(7) }
func (m *Interface) SidebarWidth() int     { return m.Px(232) }
func (m *Interface) RailWidth() int        { return m.Px(48) } // the sidebar collapsed
func (m *Interface) PaneWidth() int        { return m.Px(392) }
func (m *Interface) HeaderHeight() int     { return m.Px(52) }
func (m *Interface) BreadcrumbHeight() int { return m.Px(34) }
func (m *Interface) StatusBarHeight() int  { return m.Px(22) }

// The row, in its four heights. A view picks the one its content needs.
func (m *Interface) RowHeight() int        { return m.Px(56) }
func (m *Interface) RowHeightSub() int     { return m.Px(48) }
func (m *Interface) RowHeightSlim() int    { return m.Px(30) }
func (m *Interface) RowHeightCompact() int { return m.Px(24) }

// Control and mark heights.
func (m *Interface) ControlHeight() int { return m.Px(28) }
func (m *Interface) TabHeight() int     { return m.Px(30) }
func (m *Interface) ChipHeight() int    { return m.Px(17) }
func (m *Interface) TagHeight() int     { return m.Px(16) }
func (m *Interface) PillHeight() int    { return m.Px(15) }
func (m *Interface) BarHeight() int     { return m.Px(7) }
func (m *Interface) BarHeightWide() int { return m.Px(9) }
func (m *Interface) IconSize() int      { return m.Px(18) }
func (m *Interface) IconSizeSmall() int { return m.Px(13) }

// Corner radii and rules. A focus ring is two pixels because one is invisible
// against a border.
func (m *Interface) RadiusBar() int      { return m.Px(2) }
func (m *Interface) RadiusChip() int     { return m.Px(3) }
func (m *Interface) RadiusControl() int  { return m.Px(4) }
func (m *Interface) RadiusCard() int     { return m.Px(6) }
func (m *Interface) RadiusPill() int     { return m.Px(8) }
func (m *Interface) Hairline() int       { return m.Px(1) }
func (m *Interface) FocusRingWidth() int { return m.Px(2) }

// The widths a view reflows at. They scale too: at twice the interface size a
// screen needs twice the width before it stops being cramped.
func (m *Interface) WidthFloor() int  { return m.Px(880) }
func (m *Interface) WidthLaptop() int { return m.Px(1100) }
func (m *Interface) WidthDrawn() int  { return m.Px(1440) }

// HeightFloor is the shortest window the chrome holds: a header, a status
// bar and enough body for a list to be a list. It is the window's minimum
// height, named because a check of what a surface does at the smallest
// window has to know what that is (kvit-cash's copy of kvit-ui, 722906e).
func (m *Interface) HeightFloor() int { return m.Px(600) }
