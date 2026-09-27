package main

import (
	"embed"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

// specimen is one state or use of a component worth looking at: a caption,
// and the function that builds it. The function's own source is the code
// sample shown beside it, so a sample is always code that compiles.
type specimen struct {
	caption string
	build   func(ui *kvitui.UI) unison.Paneler
}

// entry is one component's page.
type entry struct {
	name      string
	group     string
	summary   string
	specimens []specimen
}

// catalog is every component that has a page, in the Qt catalogue's order.
var catalog = []entry{
	{
		name:    "KvitLabel",
		group:   "Foundation",
		summary: "A run of chrome text at one of the seven type roles. Every other component uses it, which is what keeps the chrome family, the colour and the eliding rule in one place.",
		specimens: []specimen{
			{"The seven roles", labelRoles},
			{"Monospace and tabular numerals", labelMonoAndTabular},
		},
	},
	{
		name:    "KvitIcon",
		group:   "Foundation",
		summary: "A symbol asked for by what it means rather than by what it looks like. The font is embedded in the library, so nothing has to be added to an application.",
		specimens: []specimen{
			{"At the sizes the chrome uses", iconSizes},
			{"An unrecognised name is visible, not blank", iconUnknown},
		},
	},
	{
		name:    "KvitIconButton",
		group:   "Foundation",
		summary: "A button whose whole label is a symbol. A real button, so it takes tab focus and a screen reader is told it is there; `Label` fills both the accessible name and the tooltip shown on pointer hover, and `Explanation` is the second sentence beside it. `Dense` draws the symbol at the 13 every symbol beside words in this library is drawn at, rather than at 18.",
		specimens: []specimen{
			{"Quiet, ordinary and checked", iconButtonForms},
			{"Two symbol sizes", iconButtonSizes},
			{"A second sentence, where the name is not enough", iconButtonExplanation},
		},
	},
	{
		name:    "KvitLink",
		group:   "Foundation",
		summary: "An inline destination or action with link semantics, natural width and optional symbol. Hover and keyboard focus both add accent and an underline; no chevron, button ground or border is invented. `Explanation` says in a sentence where following it goes, shown as the tooltip and announced as the accessible description.",
		specimens: []specimen{
			{"Plain, symbolic and keyboard-focused", linkForms},
			{"Elided in a narrow column", linkElided},
		},
	},
}

func entryNamed(name string) (entry, bool) {
	for _, e := range catalog {
		if e.name == name {
			return e, true
		}
	}
	return entry{}, false
}

//go:embed specimens_*.go
var specimenSources embed.FS

// sourceOf returns the source of a specimen function as written.
func sourceOf(fn func(ui *kvitui.UI) unison.Paneler) string {
	full := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	name := full[strings.LastIndex(full, ".")+1:]
	files, _ := specimenSources.ReadDir(".")
	for _, f := range files {
		data, err := specimenSources.ReadFile(f.Name())
		if err != nil {
			continue
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, f.Name(), data, parser.ParseComments)
		if err != nil {
			continue
		}
		for _, d := range file.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == name {
				start, end := fset.Position(fd.Pos()).Offset, fset.Position(fd.End()).Offset
				return strings.ReplaceAll(string(data[start:end]), "\t", "    ")
			}
		}
	}
	return "// source not found for " + name
}
