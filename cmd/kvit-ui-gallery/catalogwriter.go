package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kvit-s/kvit-ui/icons"
)

// catalogPath is where the vocabulary skill's catalogue lives, from the
// module's root.
const catalogPath = "agent/skills/kvit-ui/catalog.md"

// writeCatalog writes the vocabulary skill's component catalogue to path,
// from the same list the gallery draws and the library's own source.
func writeCatalog(path string) error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	text, err := renderCatalog(root)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s: %d components\n", path, len(catalog))
	return nil
}

// moduleRoot is the directory of the kvit-ui module's go.mod, found by
// walking up from the working directory: the catalogue is written from a
// checkout of the library, since it reads the library's source.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if data, err := os.ReadFile(filepath.Join(dir, "go.mod")); err == nil &&
			strings.Contains(string(data), "module github.com/kvit-s/kvit-ui\n") {
			return dir, nil
		}
		up := filepath.Dir(dir)
		if up == dir {
			return "", errors.New("catalog: not inside a checkout of github.com/kvit-s/kvit-ui")
		}
		dir = up
	}
}

// library is the kvitui package's documentation, read from its source.
type library struct {
	fset *token.FileSet
	pkg  *doc.Package
}

func readLibrary(root string) (*library, error) {
	fset := token.NewFileSet()
	files := map[string]*ast.File{}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var parsed []*ast.File
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(root, name), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[name] = f
		parsed = append(parsed, f)
	}
	pkg, err := doc.NewFromFiles(fset, parsed, "github.com/kvit-s/kvit-ui")
	if err != nil {
		return nil, err
	}
	return &library{fset: fset, pkg: pkg}, nil
}

func (l *library) print(node any) string {
	var b bytes.Buffer
	if err := printer.Fprint(&b, l.fset, node); err != nil {
		return ""
	}
	return b.String()
}

// signature is a function's declaration line.
func (l *library) signature(f *doc.Func) string {
	d := *f.Decl
	d.Doc, d.Body = nil, nil
	return l.print(&d)
}

func (l *library) typeNamed(name string) *doc.Type {
	for _, t := range l.pkg.Types {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// firstSentence is a comment's first sentence: a catalogue is a reference,
// and the reasoning belongs in the source where somebody changing it reads
// it.
func firstSentence(comment string) string {
	s := strings.Join(strings.Fields(comment), " ")
	if i := strings.Index(s, ". "); i > 0 {
		return s[:i+1]
	}
	return s
}

// otherEntries are the components the library has as functions rather than
// types.
var otherEntries = map[string][]string{
	"KvitScrollBar": {"NewScrollBar"},
	"KvitMenu":      {"UI.ShowMenu", "UI.ShowMenuAt"},
	"KvitTooltip":   {"UI.ShowTooltip"},
}

func (l *library) funcNamed(name string) *doc.Func {
	if recv, method, ok := strings.Cut(name, "."); ok {
		if t := l.typeNamed(recv); t != nil {
			for _, m := range t.Methods {
				if m.Name == method {
					return m
				}
			}
		}
		return nil
	}
	for _, f := range l.pkg.Funcs {
		if f.Name == name {
			return f
		}
	}
	for _, t := range l.pkg.Types {
		for _, f := range t.Funcs {
			if f.Name == name {
				return f
			}
		}
	}
	return nil
}

// describe writes what a caller needs to use a component: how to make one,
// its fields and its methods.
func (l *library) describe(b *strings.Builder, component string) {
	if names, ok := otherEntries[component]; ok {
		for _, name := range names {
			if f := l.funcNamed(name); f != nil {
				fmt.Fprintf(b, "```go\n%s\n```\n\n%s\n\n", l.signature(f), firstSentence(f.Doc))
			}
		}
		return
	}
	t := l.typeNamed(strings.TrimPrefix(component, "Kvit"))
	if t == nil {
		return
	}
	var made []string
	for _, f := range t.Funcs {
		made = append(made, l.signature(f))
	}
	if len(made) > 0 {
		fmt.Fprintf(b, "```go\n%s\n```\n\n", strings.Join(made, "\n"))
	}
	spec, _ := t.Decl.Specs[0].(*ast.TypeSpec)
	if st, ok := spec.Type.(*ast.StructType); ok {
		var rows []string
		for _, field := range st.Fields.List {
			comment := ""
			if field.Doc != nil {
				comment = field.Doc.Text()
			} else if field.Comment != nil {
				comment = field.Comment.Text()
			}
			note := strings.ReplaceAll(firstSentence(comment), "|", "\\|")
			kind := strings.ReplaceAll(l.print(field.Type), "|", "\\|")
			if len(field.Names) == 0 {
				if ast.IsExported(strings.TrimPrefix(kind, "*")) && !strings.Contains(kind, ".") {
					rows = append(rows, fmt.Sprintf("| `%s` | embedded | Everything a %s has. |", kind, kind))
				}
				continue
			}
			for _, n := range field.Names {
				if n.IsExported() {
					rows = append(rows, fmt.Sprintf("| `%s` | `%s` | %s |", n.Name, kind, note))
				}
			}
		}
		if len(rows) > 0 {
			b.WriteString("| Field | Type | |\n|---|---|---|\n")
			b.WriteString(strings.Join(rows, "\n"))
			b.WriteString("\n\n")
		}
	}
	var methods []string
	for _, m := range t.Methods {
		if ast.IsExported(m.Name) && !strings.HasPrefix(m.Name, "ProvideAccessibility") && m.Name != "PerformAccessibilityAction" {
			methods = append(methods, "`"+m.Name+"`")
		}
	}
	sort.Strings(methods)
	if len(methods) > 0 {
		fmt.Fprintf(b, "Methods: %s.\n\n", strings.Join(methods, ", "))
	}
}

// renderCatalog is the catalogue's text.
func renderCatalog(root string) (string, error) {
	lib, err := readLibrary(root)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString(`# The kvit-ui component catalogue

GENERATED FILE — do not edit. Written by ` + "`kvit-ui-gallery --catalog`" + `
from the same list the gallery draws and the library's own source, so what is
described here is what the repository contains.

Every component is in the ` + "`kvitui`" + ` package: one
` + "`import kvitui \"github.com/kvit-s/kvit-ui\"`" + ` reaches all of them, the ` + "`UI`" + ` that
holds the theme, the interface size and the fonts, and every colour and size
by name.

Each entry gives what the component is for, how to make one, its fields and
methods, and a working sample. The samples are the functions the gallery
draws, compiled with it and run by its tests, so one that does not work stops
the build.

`)
	fmt.Fprintf(&b, "## What there is\n\n%d components, grouped by what they are for.\n", len(catalog))
	for _, g := range groups {
		var names []string
		for _, p := range g.pages {
			if strings.HasPrefix(p, "Kvit") {
				names = append(names, "`"+p+"`")
			}
		}
		if len(names) > 0 {
			fmt.Fprintf(&b, "\n**%s** — %s", g.name, strings.Join(names, ", "))
		}
	}
	b.WriteString("\n\n## Symbols\n\n" +
		"`NewIcon` and every component that takes a `Symbol` accept these names.\n" +
		"They say what a symbol means rather than what it looks like, so the drawing\n" +
		"can change without touching a call site. A name outside this list draws a\n" +
		"marked placeholder and logs a warning, which fails the gallery's tests.\n\n")
	var symbols []string
	for _, s := range icons.MeaningNames() {
		symbols = append(symbols, "`"+s+"`")
	}
	b.WriteString(strings.Join(symbols, ", ") + "\n\n## The components\n\n")
	for _, e := range catalog {
		fmt.Fprintf(&b, "### %s\n\n%s\n\n", e.name, e.summary)
		lib.describe(&b, e.name)
		for _, s := range e.specimens {
			fmt.Fprintf(&b, "*%s*\n\n```go\n%s\n```\n\n", s.caption, sourceOf(s.build))
		}
	}
	return strings.TrimRight(b.String(), "\n") + "\n", nil
}
