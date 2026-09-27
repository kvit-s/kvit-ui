package kvitui_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// The library's rules, as the Qt version enforces them with qmllint and its
// tests: outside the packages that define design values, no colour is
// written as a literal and no font size as a number. A literal is right in
// the theme it was picked in and wrong in the other three, and a numeric size
// stands still when the reader changes the interface size.
//
// Checked: the root package (the components) and every command. Exempt: the
// packages whose job is to hold or apply design values.
var exempt = []string{"palette", "tokens", "text", "settings", "platform", "icons", "tools"}

func checkedFiles(t *testing.T) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			top := strings.Split(filepath.ToSlash(path), "/")[0]
			for _, e := range exempt {
				if top == e {
					return filepath.SkipDir
				}
			}
			if strings.HasPrefix(d.Name(), ".") || d.Name() == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") && !strings.HasSuffix(path, "_gen.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func isLiteral(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.BasicLit:
		return true
	case *ast.UnaryExpr:
		return isLiteral(v.X)
	case *ast.ParenExpr:
		return isLiteral(v.X)
	}
	return false
}

func selectorName(e ast.Expr) string {
	if s, ok := e.(*ast.SelectorExpr); ok {
		if x, ok := s.X.(*ast.Ident); ok {
			return x.Name + "." + s.Sel.Name
		}
		return "." + s.Sel.Name
	}
	return ""
}

var (
	colourCalls = map[string]bool{
		"palette.Hex": true, "palette.ParseHex": true, "palette.RGB8": true,
		"unison.RGB": true, "unison.ARGB": true, "unison.ARGBfloat": true,
	}
	// The UI methods whose first argument is a font size.
	sizeMethods = map[string]bool{"Chrome": true, "Mono": true, "Icon": true}
)

// violations lists where a file breaks the two rules.
func violations(fset *token.FileSet, f *ast.File) []string {
	var out []string
	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.CallExpr:
			if name := selectorName(v.Fun); colourCalls[name] {
				for _, a := range v.Args {
					if isLiteral(a) {
						out = append(out, fmt.Sprintf("%s: a colour written as a literal (%s); take it from the theme's tokens", fset.Position(v.Pos()), name))
						break
					}
				}
			}
			if sel, ok := v.Fun.(*ast.SelectorExpr); ok && sizeMethods[sel.Sel.Name] && len(v.Args) > 0 && isLiteral(v.Args[0]) {
				out = append(out, fmt.Sprintf("%s: a font size written as a number; take it from Interface's type roles", fset.Position(v.Pos())))
			}
		case *ast.CompositeLit:
			switch selectorName(v.Type) {
			case "text.Color", "palette.Color":
				for _, el := range v.Elts {
					if kv, ok := el.(*ast.KeyValueExpr); ok && isLiteral(kv.Value) || isLiteral(el) {
						out = append(out, fmt.Sprintf("%s: a colour written as a literal; take it from the theme's tokens", fset.Position(v.Pos())))
						break
					}
				}
			case "text.Style":
				for _, el := range v.Elts {
					if kv, ok := el.(*ast.KeyValueExpr); ok {
						if k, ok := kv.Key.(*ast.Ident); ok && k.Name == "Size" && isLiteral(kv.Value) {
							out = append(out, fmt.Sprintf("%s: a font size written as a number; take it from Interface's type roles", fset.Position(v.Pos())))
						}
					}
				}
			}
		}
		return true
	})
	return out
}

func TestNoColourLiteralsOrNumericFontSizes(t *testing.T) {
	fset := token.NewFileSet()
	for _, path := range checkedFiles(t) {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range violations(fset, f) {
			t.Error(v)
		}
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
	if got := violations(fset, f); len(got) != 5 {
		t.Errorf("found %d of the 5 planted literals:\n%s", len(got), strings.Join(got, "\n"))
	}
}
