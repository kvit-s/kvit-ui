package kvitui_test

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kvit-s/kvit-ui/uitest"
)

// The library's rules, as the version enforces them with viewlint and its
// tests: outside the packages that define design values, no colour is
// written as a literal and no font size as a number. Checked: the root
// package (the components) and every command. Exempt: the packages whose job
// is to hold or apply design values.
func TestNoColourLiteralsOrNumericFontSizes(t *testing.T) {
	// The components and the gallery at the least.
	if n := uitest.CheckRules(t, ".", "palette", "tokens", "text", "settings", "platform", "icons", "tools"); n < 50 {
		t.Errorf("the rule check read only %d files", n)
	}
}

// The third rule: outside the same packages, no position, size, spacing,
// radius or margin is written as a bare number. A table column's width is in
// design pixels, which the table passes through Px, so it is not one.
func TestNoUnnamedGeometryValues(t *testing.T) {
	if n := uitest.CheckGeometry(t, ".", "palette", "tokens", "text", "settings", "platform", "icons", "tools"); n < 50 {
		t.Errorf("the geometry check read only %d files", n)
	}
}

// The geometry check catches each kind of number planted here, and passes
// zero, a measure, arithmetic on a measure and a column's design pixels.
func TestTheGeometryCheckCatchesNumbers(t *testing.T) {
	src := `package x
func f() {
	_ = geom.NewRect(4, 0, w, 12.5)
	_ = unison.FlexLayout{HSpacing: 6, VSpacing: m.Space()}
	_ = geom.Insets{Top: -2}
	r.Width = (40)
	_ = geom.NewPoint(0, 0.0)
	_ = geom.NewSize(float32(m.Px(400)), h/2)
	_ = []kvitui.TableColumn{{Title: "Date", Width: 96}}
	_ = TableColumn{Width: 40}
	_ = map[int]*TableColumn{1: {Width: 30}}
}`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got := uitest.GeometryViolations(fset, f); len(got) != 5 {
		t.Errorf("found %d of the 5 planted numbers:\n%s", len(got), strings.Join(got, "\n"))
	}
}

// The check itself works: it catches each kind of literal planted here.
func TestTheRuleCheckCatchesLiterals(t *testing.T) {
	src := `package x
func f() {
	_ = palette.Hex("#ff0000")
	_ = ui.Chrome(12, 400, c)
	_ = text.Style{Size: 14}
	_ = text.Color{R: 1, A: 255}
	_ = unison.RGB(1, 2, 3)
	_ = ui.Chrome(m.Body(), 400, c)
}`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got := uitest.Violations(fset, f); len(got) != 5 {
		t.Errorf("found %d of the 5 planted literals:\n%s", len(got), strings.Join(got, "\n"))
	}
}
