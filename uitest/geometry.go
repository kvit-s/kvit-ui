package uitest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// CheckGeometry fails the test for every place in the Go files under dir that
// writes a position, a size, a spacing, a radius or a margin as a bare
// number, as the library's build refused an "unnamed geometry value" in
// view markup. A number there stands still when the reader changes the interface
// size, while everything taken from Interface, or given in design pixels to
// Px, moves with it. exempt, the files read and the return value are as for
// CheckRules. It is a check of its own rather than part of CheckRules so that
// an application takes it on when its screens are ready for it.
func CheckGeometry(t testing.TB, dir string, exempt ...string) int {
	t.Helper()
	files, err := goFiles(dir, exempt)
	if err != nil {
		t.Fatal(err)
		return 0
	}
	if len(files) == 0 {
		t.Fatalf("the geometry check found no Go files under %s", dir)
		return 0
	}
	fset := token.NewFileSet()
	for _, path := range files {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range GeometryViolations(fset, f) {
			t.Error(v)
		}
	}
	return len(files)
}

var (
	// geometryConstructors are the geom functions that take a position, a
	// size or insets.
	geometryConstructors = map[string]bool{"NewRect": true, "NewSize": true, "NewPoint": true, "NewInsets": true,
		"NewUniformInsets": true, "NewHorizontalInsets": true, "NewVerticalInsets": true}
	// geometryFields are the struct fields that hold a position, a size, a
	// spacing, a radius or a margin: the properties the gate named in
	// view markup (x, width, spacing, radius, the anchors' margins and the rest).
	geometryFields = map[string]bool{"X": true, "Y": true, "Width": true, "Height": true, "HSpacing": true,
		"VSpacing": true, "Radius": true, "Top": true, "Left": true, "Bottom": true, "Right": true}
	// designPixelTypes are the library's structs whose sizes are design
	// pixels, which the component passes through Px itself: a number in a
	// TableColumn's Width is already a measure that moves with the
	// interface size.
	designPixelTypes = map[string]bool{"TableColumn": true}
)

// GeometryViolations lists where a file writes a position, a size, a
// spacing, a radius or a margin as a bare number: a number other than zero
// given to one of geom's constructors, or given to one of those fields in a
// struct literal or an assignment. Zero is not a measure, and arithmetic on a
// measure, such as halving it, is not a value of its own.
func GeometryViolations(fset *token.FileSet, f *ast.File) []string {
	var out []string
	flag := func(at ast.Node, what string) {
		out = append(out, fmt.Sprintf("%s: a geometry value written as a number (%s); take it from Interface, or give design pixels to Px", fset.Position(at.Pos()), what))
	}
	inDesignPixels := designPixelLiterals(f)
	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.CallExpr:
			if s, ok := v.Fun.(*ast.SelectorExpr); ok && geometryConstructors[s.Sel.Name] {
				if pkg, ok := s.X.(*ast.Ident); ok && pkg.Name == "geom" {
					for _, a := range v.Args {
						if numberNotZero(a) {
							flag(a, "geom."+s.Sel.Name)
						}
					}
				}
			}
		case *ast.CompositeLit:
			if inDesignPixels[v] {
				return true
			}
			for _, el := range v.Elts {
				if kv, ok := el.(*ast.KeyValueExpr); ok {
					if k, ok := kv.Key.(*ast.Ident); ok && geometryFields[k.Name] && numberNotZero(kv.Value) {
						flag(kv.Value, k.Name)
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range v.Lhs {
				if s, ok := lhs.(*ast.SelectorExpr); ok && geometryFields[s.Sel.Name] && i < len(v.Rhs) && numberNotZero(v.Rhs[i]) {
					flag(v.Rhs[i], s.Sel.Name)
				}
			}
		}
		return true
	})
	return out
}

// designPixelLiterals are the struct literals in a file whose type is one of
// designPixelTypes, written out or left to the slice, array or map that
// holds them.
func designPixelLiterals(f *ast.File) map[*ast.CompositeLit]bool {
	found := map[*ast.CompositeLit]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		own, elem := literalTypeNames(lit.Type)
		if designPixelTypes[own] {
			found[lit] = true
		}
		if designPixelTypes[elem] {
			for _, el := range lit.Elts {
				if kv, ok := el.(*ast.KeyValueExpr); ok {
					el = kv.Value
				}
				if u, ok := el.(*ast.UnaryExpr); ok {
					el = u.X
				}
				if c, ok := el.(*ast.CompositeLit); ok && c.Type == nil {
					found[c] = true
				}
			}
		}
		return true
	})
	return found
}

// literalTypeNames is the name of a literal's type without its package, and,
// for a slice, an array or a map, the name of the type it holds.
func literalTypeNames(e ast.Expr) (own, elem string) {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name, ""
	case *ast.SelectorExpr:
		return v.Sel.Name, ""
	case *ast.StarExpr:
		return literalTypeNames(v.X)
	case *ast.ArrayType:
		elem, _ = literalTypeNames(v.Elt)
		return "", elem
	case *ast.MapType:
		elem, _ = literalTypeNames(v.Value)
		return "", elem
	}
	return "", ""
}

// numberNotZero reports a bare numeric literal other than zero, signed or
// in parentheses.
func numberNotZero(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind != token.INT && v.Kind != token.FLOAT {
			return false
		}
		return strings.Trim(v.Value, "0._") != ""
	case *ast.UnaryExpr:
		return numberNotZero(v.X)
	case *ast.ParenExpr:
		return numberNotZero(v.X)
	}
	return false
}
