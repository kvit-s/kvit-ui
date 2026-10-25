package uitest

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

var (
	oops      = regexp.MustCompile(`(?i)\boops\b`)
	tap       = regexp.MustCompile(`(?i)\btap(ped|ping|s)?\b`)
	clickHere = regexp.MustCompile(`(?i)\bclick here\b`)
	// labelCalls are the constructors whose argument at the index given is
	// a command's words.
	labelCalls = map[string]int{"NewButton": 1, "NewMenuButton": 1, "NewIconButton": 2}
	// labelFields are the fields of the composites named that hold a
	// command's words.
	labelFields = map[string]map[string]bool{
		"MenuItem": {"Text": true},
		"Command":  {"Name": true, "TurnOn": true, "TurnOff": true},
	}
)

// WordingViolations lists where a file's source writes what the interface
// rules forbid: "Oops" (WRD-15), "tap" (WRD-17) or "click here" (WRD-18) in a
// string; three full stops for an ellipsis, or an ellipsis typed into a
// command's words, where Asks adds it in the form's way (WRD-4); a key
// binding written with mod.Control or mod.Command, which is right on one
// system only, where kvitui.Keys with ModPrimary is right on both (KEY-3);
// and a button's words starting with a sign typed for a symbol (CMD-22).
func WordingViolations(fset *token.FileSet, f *ast.File) []string {
	var out []string
	flag := func(at ast.Node, what string) {
		out = append(out, fmt.Sprintf("%s: %s", fset.Position(at.Pos()), what))
	}
	words := func(e ast.Expr) (string, bool) {
		lit, ok := e.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(lit.Value)
		return s, err == nil
	}
	label := func(e ast.Expr) {
		s, ok := words(e)
		if !ok {
			return
		}
		if strings.HasSuffix(s, "…") {
			flag(e, "an ellipsis typed into a command's words; set Asks and kvit-ui adds it in the form's way")
		}
		if strings.HasPrefix(s, "+ ") || strings.HasPrefix(s, "- ") {
			flag(e, "a sign typed before a command's words; give the button its symbol instead")
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.BasicLit:
			s, ok := words(v)
			if !ok {
				return true
			}
			switch {
			case oops.MatchString(s):
				flag(v, `"Oops" in a message; say what happened (WRD-15)`)
			case tap.MatchString(s):
				flag(v, `"tap" on a desktop; say select, choose or click (WRD-17)`)
			case clickHere.MatchString(s):
				flag(v, `"click here"; a link says where it goes (WRD-18)`)
			}
			if strings.HasSuffix(s, "...") {
				flag(v, "three full stops; write the ellipsis character, or set Asks (WRD-4)")
			}
		case *ast.CallExpr:
			name := ""
			switch fn := v.Fun.(type) {
			case *ast.SelectorExpr:
				name = fn.Sel.Name
			case *ast.Ident:
				name = fn.Name
			}
			if i, ok := labelCalls[name]; ok && i < len(v.Args) {
				label(v.Args[i])
			}
		case *ast.CompositeLit:
			own, _ := literalTypeNames(v.Type)
			if own == "KeyBinding" {
				for _, el := range v.Elts {
					if kv, ok := el.(*ast.KeyValueExpr); ok {
						ast.Inspect(kv.Value, func(m ast.Node) bool {
							if sel, ok := m.(*ast.SelectorExpr); ok && (selectorName(sel) == "mod.Control" || selectorName(sel) == "mod.Command") {
								flag(sel, "a key binding written with "+selectorName(sel)+", right on one system only; write kvitui.Keys with ModPrimary (KEY-3)")
							}
							return true
						})
					}
				}
			}
			if fields, ok := labelFields[own]; ok {
				for _, el := range v.Elts {
					if kv, ok := el.(*ast.KeyValueExpr); ok {
						if k, ok := kv.Key.(*ast.Ident); ok && fields[k.Name] {
							label(kv.Value)
						}
					}
				}
			}
		}
		return true
	})
	return out
}
