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
	{
		name:    "KvitHeader",
		group:   "Structure",
		summary: "The strip across the top of the window: the wordmark, navigation in the middle, actions on the right, in that fixed order on every screen.",
		specimens: []specimen{
			{"With navigation and actions", headerWithNavigation},
		},
	},
	{
		name:    "KvitSidebar",
		group:   "Structure",
		summary: "The list of places down the left edge. Below the laptop breakpoint it collapses to its rail and every item becomes its symbol alone.",
		specimens: []specimen{
			{"Expanded and collapsed", sidebarExpandedAndCollapsed},
		},
	},
	{
		name:    "KvitSidebarItem",
		group:   "Structure",
		summary: "One place in the sidebar. The selected item takes a bar down its leading edge as well as a tint, so the selection is not resting on colour. A count is hidden once the sidebar collapses to the rail unless `CountInRail` asks for it, and `CountMax` raises the badge's cap where the place the item points at states the true number.",
		specimens: []specimen{
			{"Selected, hovered and counted", sidebarItemStates},
			{"The rail, with and without its counts", sidebarItemRail},
		},
	},
	{
		name:    "KvitBreadcrumb",
		group:   "Structure",
		summary: "Where the reader is and the way back. The last crumb is the current place and is deliberately not a link; a long trail elides from the middle, keeping the section and the current place.",
		specimens: []specimen{
			{"A short trail and a long one", breadcrumbTrails},
		},
	},
	{
		name:    "KvitPanel",
		group:   "Content",
		summary: "A region of the window with its own ground: a sidebar, a toolbar strip. Structural, where a card is content — which is why it has no radius.",
		specimens: []specimen{
			{"With rules on two edges", panelWithRules},
		},
	},
	{
		name:    "KvitDivider",
		group:   "Content",
		summary: "A rule between two things. One design pixel, in the decorative border token rather than the control-boundary one.",
		specimens: []specimen{
			{"Horizontal and vertical", dividerBothWays},
		},
	},
	{
		name:    "KvitBadge",
		group:   "Marks",
		summary: "A count attached to something else. Caps rather than growing wide, and hides at zero — a badge showing nought says look here about nothing. The number is drawn and announced with the reader's own digit grouping, and the noun beside it comes from the caller in two slots, a singular and a plural.",
		specimens: []specimen{
			{"Counts, capped, and hidden at zero", badgeCounts},
			{"What the count counts", badgeNouns},
		},
	},
	{
		name:    "KvitTab",
		group:   "Controls",
		summary: "One tab in a row of them. The selected tab is marked by an underline as well as by colour and weight — selection shown by colour alone is the most common place the rule gets broken. `Explanation` is one sentence saying what the view behind the tab shows, which two or three words cannot.",
		specimens: []specimen{
			{"Selected, counted and plain", tabForms},
			{"Each tab saying what its view shows", tabExplanations},
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
