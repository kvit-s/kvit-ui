package kvitui_test

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kvit-s/kvit-ui/uitest"
)

// The library's rules, as the Qt version enforces them with qmllint and its
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
