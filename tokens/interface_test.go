package tokens_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kvit-s/kvit-ui/settings"
	"github.com/kvit-s/kvit-ui/tokens"
)

// derived lists every value Interface derives from the interface size: each
// method that takes nothing and returns an int, apart from the size itself.
func derived(m *tokens.Interface) map[string]int {
	out := map[string]int{}
	v := reflect.ValueOf(m)
	intType := reflect.TypeOf(0)
	for i := 0; i < v.NumMethod(); i++ {
		name := v.Type().Method(i).Name
		mt := v.Method(i).Type()
		if name == "FontSize" || mt.NumIn() != 0 || mt.NumOut() != 1 || mt.Out(0) != intType {
			continue
		}
		out[name] = int(v.Method(i).Call(nil)[0].Int())
	}
	return out
}

func TestDefaultReproducesTheOldLiterals(t *testing.T) {
	m := tokens.NewInterface()
	if m.FontSize() != 12 || m.Scale() != 1 {
		t.Fatalf("default size %d, scale %v", m.FontSize(), m.Scale())
	}
	for _, d := range []int{1, 2, 4, 6, 8, 10, 14, 16, 22, 24, 26, 28, 30, 60, 96, 148, 170, 480, 560} {
		if m.Px(d) != d {
			t.Errorf("Px(%d) = %d at the default size", d, m.Px(d))
		}
	}
}

func TestEveryValueIsTheIdentityAtTheDefault(t *testing.T) {
	want := map[string]int{
		"Caption": 10, "Small": 11, "Body": 12, "Strong": 13, "Title": 15, "Headline": 17, "Display": 20,
		"SpaceTight": 2, "SpaceSnug": 4, "SpaceNear": 6, "Space": 8, "SpaceWide": 10, "SpaceLoose": 12,
		"ViewMargin": 16, "ColumnGap": 14, "StackGap": 7, "SidebarWidth": 232, "RailWidth": 48, "PaneWidth": 392,
		"HeaderHeight": 52, "BreadcrumbHeight": 34, "StatusBarHeight": 22,
		"RowHeight": 56, "RowHeightSub": 48, "RowHeightSlim": 30, "RowHeightCompact": 24,
		"ControlHeight": 28, "TabHeight": 30, "ChipHeight": 17, "TagHeight": 16, "PillHeight": 15,
		"BarHeight": 7, "BarHeightWide": 9, "IconSize": 18, "IconSizeSmall": 13,
		"RadiusBar": 2, "RadiusChip": 3, "RadiusControl": 4, "RadiusCard": 6, "RadiusPill": 8,
		"Hairline": 1, "FocusRingWidth": 2, "WidthFloor": 880, "WidthLaptop": 1100, "WidthDrawn": 1440,
		"HeightFloor": 600, "FloatingViewWidth": 720,
	}
	got := derived(tokens.NewInterface())
	if len(got) != len(want) {
		t.Errorf("Interface derives %d values, want %d", len(got), len(want))
	}
	for name, w := range want {
		if g, ok := got[name]; !ok || g != w {
			t.Errorf("%s = %d, want %d", name, g, w)
		}
	}
}

func TestEveryValueScalesWithOneSetting(t *testing.T) {
	m := tokens.NewInterface()
	before := derived(m)
	if len(before) <= 30 {
		t.Fatalf("only %d derived values found", len(before))
	}
	n := counter(m.OnChanged)
	m.SetFontSize(24)
	if *n != 1 || m.Scale() != 2 {
		t.Fatalf("doubling the size: %d notifications, scale %v", *n, m.Scale())
	}
	for name, after := range derived(m) {
		b := before[name]
		if after <= b {
			t.Errorf("%s stayed at %d when the interface size doubled", name, b)
		}
		if d := after - 2*b; d < -1 || d > 1 {
			t.Errorf("%s went from %d to %d, not a doubling", name, b, after)
		}
	}
	if m.Body() != 24 || m.Caption() != 20 || m.Title() != 30 || m.Px(28) != 56 {
		t.Error("the roles did not scale with the geometry")
	}
	m.SetFontSize(24)
	if *n != 1 {
		t.Error("setting the same size again notified")
	}
}

func TestClamps(t *testing.T) {
	m := tokens.NewInterface()
	m.SetFontSize(2)
	if m.FontSize() != tokens.MinInterfaceSize {
		t.Errorf("size 2 clamped to %d", m.FontSize())
	}
	m.SetFontSize(400)
	if m.FontSize() != tokens.MaxInterfaceSize {
		t.Errorf("size 400 clamped to %d", m.FontSize())
	}
}

func TestNoDensityValueCollapsesAtTheSmallestSize(t *testing.T) {
	m := tokens.NewInterface()
	m.SetFontSize(tokens.MinInterfaceSize)
	for name, v := range derived(m) {
		if v < 1 {
			t.Errorf("%s collapsed to %d at the smallest size", name, v)
		}
	}
	if m.Px(0) != 0 {
		t.Error("Px(0) should stay 0")
	}
}

func TestTheScalesAreOrderedAtEverySize(t *testing.T) {
	m := tokens.NewInterface()
	for size := tokens.MinInterfaceSize; size <= tokens.MaxInterfaceSize; size++ {
		m.SetFontSize(size)
		typeScale := []int{m.Caption(), m.Small(), m.Body(), m.Strong(), m.Title(), m.Headline(), m.Display()}
		for i := 1; i < len(typeScale); i++ {
			if typeScale[i] <= typeScale[i-1] {
				t.Errorf("at size %d the type scale %v is not increasing", size, typeScale)
				break
			}
		}
		space := []int{m.Hairline(), m.SpaceTight(), m.SpaceSnug(), m.SpaceNear(), m.Space(), m.SpaceWide(), m.SpaceLoose()}
		for i := 1; i < len(space); i++ {
			if space[i] < space[i-1] {
				t.Errorf("at size %d spacing step %d is %d, below %d", size, i, space[i], space[i-1])
			}
		}
	}
}

// The density tokens of the Kvit Hub version built with Qt (viewMargin,
// typeTitle and the rest) each have a value here.
func TestEveryKvitHubTokenHasAHome(t *testing.T) {
	got := derived(tokens.NewInterface())
	for hub, name := range map[string]string{
		"viewMargin": "ViewMargin", "columnGap": "ColumnGap", "stackGap": "StackGap", "sidebarWidth": "SidebarWidth",
		"railWidth": "RailWidth", "paneWidth": "PaneWidth", "headerHeight": "HeaderHeight", "breadcrumbHeight": "BreadcrumbHeight",
		"rowHeight": "RowHeight", "rowHeightSub": "RowHeightSub", "rowHeightSlim": "RowHeightSlim", "rowHeightCompact": "RowHeightCompact",
		"tabHeight": "TabHeight", "chipHeight": "ChipHeight", "tagHeight": "TagHeight", "pillHeight": "PillHeight",
		"barHeight": "BarHeight", "barHeightWide": "BarHeightWide", "radiusBar": "RadiusBar", "radiusChip": "RadiusChip",
		"radiusControl": "RadiusControl", "radiusCard": "RadiusCard", "radiusPill": "RadiusPill", "hairline": "Hairline",
		"widthFloor": "WidthFloor", "widthLaptop": "WidthLaptop", "widthDrawn": "WidthDrawn",
		"typeTitle": "Display", "typePage": "Headline", "typeHeading": "Title", "typeName": "Strong",
		"typeBody": "Body", "typeRow": "Body", "typeSecondary": "Small", "typeSmall": "Small", "typeMicro": "Caption",
	} {
		if v, ok := got[name]; !ok || v <= 0 {
			t.Errorf("kvit-hub's %s has nowhere to go: Interface.%s", hub, name)
		}
	}
}

func TestInterfacePersistsThroughSettings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	{
		s := settings.New()
		if err := s.Open(path, false); err != nil {
			t.Fatal(err)
		}
		m := tokens.NewInterface()
		m.SetSettings(s)
		m.SetFontSize(16)
		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	s := settings.New()
	if err := s.Open(path, false); err != nil {
		t.Fatal(err)
	}
	m := tokens.NewInterface()
	m.SetSettings(s)
	if m.FontSize() != 16 {
		t.Errorf("the size came back as %d", m.FontSize())
	}
	corrupt := openStore(t, "corrupt.json")
	corrupt.SetValue("interface.fontSize", 400)
	clamped := tokens.NewInterface()
	clamped.SetSettings(corrupt)
	if clamped.FontSize() != tokens.MaxInterfaceSize {
		t.Errorf("a stored 400 loaded as %d", clamped.FontSize())
	}
}

func TestInterfaceResetToDefaults(t *testing.T) {
	m := tokens.NewInterface()
	m.SetFontSize(20)
	m.SetFontFamily("Some Sans")
	m.ResetToDefaults()
	if m.FontSize() != tokens.DefaultInterfaceSize || m.Body() != 12 || m.FontFamily() != "" {
		t.Error("reset did not restore the defaults")
	}
}

// The chrome and the document are two settings: moving one leaves the other.
func TestTheDocumentTypeScaleIsUntouched(t *testing.T) {
	s := openStore(t, "settings.json")
	m := tokens.NewInterface()
	ty := tokens.NewTypography()
	m.SetSettings(s)
	ty.SetSettings(s)
	docBase, docBody := ty.BaseSize(), ty.BodySize()
	m.SetFontSize(24)
	if m.Body() != 24 || ty.BaseSize() != docBase || ty.BodySize() != docBody {
		t.Error("the interface size moved the document's")
	}
	chrome := m.Body()
	ty.SetBaseSize(20)
	if ty.BaseSize() != 20 || m.Body() != chrome {
		t.Error("the document size moved the interface's")
	}
	if v, _ := s.Value("interface.fontSize"); v != 24.0 {
		t.Errorf("interface.fontSize stored as %v", v)
	}
	if v, _ := s.Value("typography.fontSize"); v != 20.0 {
		t.Errorf("typography.fontSize stored as %v", v)
	}
}

func TestTheFontFamiliesAreSeparateFromTheDocument(t *testing.T) {
	s := openStore(t, "settings.json")
	m := tokens.NewInterface()
	ty := tokens.NewTypography()
	m.SetSettings(s)
	ty.SetSettings(s)
	ty.SetFontFamily("Some Serif")
	if m.FontFamily() != "" {
		t.Error("the document family moved the chrome's")
	}
	m.SetFontFamily("Some Sans")
	if ty.FontFamily() != "Some Serif" {
		t.Error("the chrome family moved the document's")
	}
	m.SetMonoFamily("")
	if m.MonoFamily() != "monospace" {
		t.Errorf("an empty mono family became %q", m.MonoFamily())
	}
}

func TestTypographyScale(t *testing.T) {
	ty := tokens.NewTypography()
	check := func(base int, want [6]int) {
		t.Helper()
		ty.SetBaseSize(base)
		got := [6]int{ty.SizeForRole(tokens.Body), ty.SizeForRole(tokens.Heading1), ty.SizeForRole(tokens.Heading2),
			ty.SizeForRole(tokens.Heading3), ty.SizeForRole(tokens.Heading4), ty.SizeForRole(tokens.Mono)}
		if got != want {
			t.Errorf("at base %d the scale is %v, want %v", base, got, want)
		}
	}
	check(14, [6]int{14, 30, 22, 19, 16, 12}) // the default
	check(15, [6]int{15, 32, 24, 20, 17, 13}) // the base the ratios were measured at
	check(10, [6]int{10, 21, 16, 13, 11, 9})
	check(28, [6]int{28, 60, 45, 37, 32, 24})
	if ty.BodySize() != 28 || ty.MonoSize() != 24 {
		t.Error("BodySize and MonoSize disagree with SizeForRole")
	}
}

func TestTypographyClampsAndSignals(t *testing.T) {
	ty := tokens.NewTypography()
	n := counter(ty.OnChanged)
	ty.SetBaseSize(2)
	ty.SetLineHeight(9)
	ty.SetParagraphSpacing(-3)
	ty.SetMaxContentWidth(40)
	if ty.BaseSize() != tokens.MinBaseSize || ty.LineHeight() != tokens.MaxLineHeight ||
		ty.ParagraphSpacing() != tokens.MinParagraphSpacing || ty.MaxContentWidth() != tokens.MinContentWidth {
		t.Errorf("clamped to %d, %v, %d, %d", ty.BaseSize(), ty.LineHeight(), ty.ParagraphSpacing(), ty.MaxContentWidth())
	}
	if *n != 4 {
		t.Errorf("four changes gave %d notifications", *n)
	}
	ty.SetBaseSize(ty.BaseSize())
	ty.SetLineHeight(ty.LineHeight())
	if *n != 4 {
		t.Error("setting unchanged values notified")
	}
	ty.SetMaxContentWidth(0)
	if ty.MaxContentWidth() != 0 {
		t.Error("0 means no cap and must stay 0")
	}
}

func TestTypographyPersistsAndClampsOnLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	{
		s := settings.New()
		if err := s.Open(path, false); err != nil {
			t.Fatal(err)
		}
		ty := tokens.NewTypography()
		ty.SetSettings(s)
		ty.SetFontFamily("Some Serif")
		ty.SetBaseSize(18)
		ty.SetLineHeight(1.6)
		ty.SetParagraphSpacing(10)
		ty.SetMaxContentWidth(700)
		ty.SetMonoFamily("Some Mono")
		if err := s.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	s := settings.New()
	if err := s.Open(path, false); err != nil {
		t.Fatal(err)
	}
	ty := tokens.NewTypography()
	ty.SetSettings(s)
	if ty.FontFamily() != "Some Serif" || ty.BaseSize() != 18 || ty.LineHeight() != 1.6 ||
		ty.ParagraphSpacing() != 10 || ty.MaxContentWidth() != 700 || ty.MonoFamily() != "Some Mono" {
		t.Error("the typography did not survive a reopen")
	}
	corrupt := openStore(t, "corrupt.json")
	corrupt.SetValue("typography.fontSize", 999)
	corrupt.SetValue("typography.lineHeight", 0.1)
	c := tokens.NewTypography()
	c.SetSettings(corrupt)
	if c.BaseSize() != tokens.MaxBaseSize || c.LineHeight() != tokens.MinLineHeight {
		t.Errorf("corrupt values loaded as %d, %v", c.BaseSize(), c.LineHeight())
	}
	c.ResetToDefaults()
	if c.BaseSize() != tokens.DefaultBaseSize || c.LineHeight() != tokens.DefaultLineHeight || c.ParagraphSpacing() != tokens.DefaultParagraphSpacing {
		t.Error("reset did not restore the defaults")
	}
}
