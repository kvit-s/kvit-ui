package uitest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// CheckRules fails the test for every place in the Go files under dir that
// writes a colour as a literal or a font size as a number, as the Qt library
// fails its build on them. A literal colour is right in the theme it was
// picked in and wrong in the other three, and a numeric size stands still
// when the reader changes the interface size. exempt names top-level
// directories under dir whose job is to hold or apply design values; tests,
// generated files, hidden directories and build output are not read. It
// returns how many files it read, which a caller can hold to a floor: a
// check that reads nothing passes.
func CheckRules(t testing.TB, dir string, exempt ...string) int {
	t.Helper()
	files, err := goFiles(dir, exempt)
	if err != nil {
		t.Fatal(err)
		return 0
	}
	if len(files) == 0 {
		t.Fatalf("the rule check found no Go files under %s", dir)
		return 0
	}
	fset := token.NewFileSet()
	for _, path := range files {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range Violations(fset, f) {
			t.Error(v)
		}
	}
	return len(files)
}

// goFiles are the Go files the rules apply to.
func goFiles(dir string, exempt []string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == dir {
				return nil
			}
			rel, _ := filepath.Rel(dir, path)
			if slices.Contains(exempt, strings.Split(filepath.ToSlash(rel), "/")[0]) ||
				strings.HasPrefix(d.Name(), ".") || d.Name() == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") && !strings.HasSuffix(path, "_gen.go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
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

// Violations lists where a file writes a colour as a literal or a font size
// as a number.
func Violations(fset *token.FileSet, f *ast.File) []string {
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
