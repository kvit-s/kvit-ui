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
	// sourceOnly shows the code without drawing it, for the one thing a page
	// cannot draw: a window, which has no place inside another window. The
	// gallery's tests still build and run it.
	sourceOnly bool
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
			{caption: "The seven roles", build: labelRoles},
			{caption: "Monospace and tabular numerals", build: labelMonoAndTabular},
		},
	},
	{
		name:    "KvitIcon",
		group:   "Foundation",
		summary: "A symbol asked for by what it means rather than by what it looks like. The font is embedded in the library, so nothing has to be added to an application.",
		specimens: []specimen{
			{caption: "At the sizes the chrome uses", build: iconSizes},
			{caption: "An unrecognised name is visible, not blank", build: iconUnknown},
		},
	},
	{
		name:    "KvitIconButton",
		group:   "Foundation",
		summary: "A button whose whole label is a symbol. A real button, so it takes tab focus and a screen reader is told it is there; `Label` fills both the accessible name and the tooltip shown on pointer hover, and `Explanation` is the second sentence beside it. `Dense` draws the symbol at the 13 every symbol beside words in this library is drawn at, rather than at 18.",
		specimens: []specimen{
			{caption: "Quiet, ordinary and checked", build: iconButtonForms},
			{caption: "Two symbol sizes", build: iconButtonSizes},
			{caption: "A second sentence, where the name is not enough", build: iconButtonExplanation},
		},
	},
	{
		name:    "KvitLink",
		group:   "Foundation",
		summary: "An inline destination or action with link semantics, natural width and optional symbol. Hover and keyboard focus both add accent and an underline; no chevron, button ground or border is invented. `Explanation` says in a sentence where following it goes, shown as the tooltip and announced as the accessible description.",
		specimens: []specimen{
			{caption: "Plain, symbolic and keyboard-focused", build: linkForms},
			{caption: "Elided in a narrow column", build: linkElided},
		},
	},
	{
		name:    "KvitHeader",
		group:   "Structure",
		summary: "The strip across the top of the window: the wordmark, navigation in the middle, actions on the right, in that fixed order on every screen.",
		specimens: []specimen{
			{caption: "With navigation and actions", build: headerWithNavigation},
		},
	},
	{
		name:    "KvitSidebar",
		group:   "Structure",
		summary: "The list of places down the left edge. Below the laptop breakpoint it collapses to its rail and every item becomes its symbol alone.",
		specimens: []specimen{
			{caption: "Expanded and collapsed", build: sidebarExpandedAndCollapsed},
		},
	},
	{
		name:    "KvitSidebarItem",
		group:   "Structure",
		summary: "One place in the sidebar. The selected item takes a bar down its leading edge as well as a tint, so the selection is not resting on colour. A count is hidden once the sidebar collapses to the rail unless `CountInRail` asks for it, and `CountMax` raises the badge's cap where the place the item points at states the true number.",
		specimens: []specimen{
			{caption: "Selected, hovered and counted", build: sidebarItemStates},
			{caption: "The rail, with and without its counts", build: sidebarItemRail},
		},
	},
	{
		name:    "KvitBreadcrumb",
		group:   "Structure",
		summary: "Where the reader is and the way back. The last crumb is the current place and is deliberately not a link; a long trail elides from the middle, keeping the section and the current place.",
		specimens: []specimen{
			{caption: "A short trail and a long one", build: breadcrumbTrails},
		},
	},
	{
		name:    "KvitRegion",
		group:   "Structure",
		summary: "A body that takes the height left over and scrolls what does not fit. The scroll bar sits beside the content rather than over it, so nothing is ever hidden behind it.",
		specimens: []specimen{
			{caption: "Scrolling a column of rows", build: regionScrolling},
		},
	},
	{
		name:    "KvitViewHead",
		group:   "Structure",
		summary: "The strip at the top of a view: what the view is, how much is in it, and the controls that act on all of it.",
		specimens: []specimen{
			{caption: "Title, count and controls", build: viewHeadWithControls},
		},
	},
	{
		name:    "KvitStatusBar",
		group:   "Structure",
		summary: "The strip along the bottom: what is happening on the left, standing facts on the right. Facts come in three shapes — plain strings, named groups whose facts open what they name, and whole controls at the end. What the bar has no room for goes into a menu behind a control saying how many there are, rather than being cut off the end of the list. `GroupsFirst` puts the groups before the activity, for a bar whose left end is the list of what is waiting.",
		specimens: []specimen{
			{caption: "Working, with two facts", build: statusBarWorking},
			{caption: "Grouped facts that open what they name", build: statusBarGroups},
			{caption: "Narrow: what does not fit is in the menu, not gone", build: statusBarNarrow},
			{caption: "Groups first, for a bar whose left end is the work", build: statusBarGroupsFirst},
		},
	},
	{
		name:    "KvitWindow",
		group:   "Structure",
		summary: "The application shell: header, optional sidebar, body and status bar, with the sidebar collapsing to its rail below the laptop breakpoint and opening over the body while the pointer rests on it. Not shown here because it is a window; see the code sample.",
		specimens: []specimen{
			{caption: "The whole shell", build: windowShell, sourceOnly: true},
		},
	},
	{
		name:    "KvitSectionHeading",
		group:   "Content",
		summary: "A group heading: a filled bar with a disclosure chevron, the name, what the group holds, the count with the word for what was counted, and the one action that applies to every row under it. A count that is not a number goes in `CountText` and is drawn beside the name as written; `ActionSymbol` draws the action as a symbol and keeps its words as the button's name and tooltip, with `ActionExplanation` saying in a sentence what running it does. A collapsible heading joins the tab order and opens on Return, Enter or Space; the hoisted action is a control of its own, so running it never also collapses the group.",
		specimens: []specimen{
			{caption: "Collapsible, counted, with an action", build: sectionHeadingForms},
			{caption: "Reached by tab, opened by Return", build: sectionHeadingKeyboard},
			{caption: "A written count, and an action drawn as a symbol", build: sectionHeadingWrittenCount},
			{caption: "What the action does, in a sentence", build: sectionHeadingExplanation},
		},
	},
	{
		name:    "KvitRow",
		group:   "Content",
		summary: "A list row at one of four heights, chosen by what the row carries rather than by how many rows a view wants to fit. Hover and keyboard focus are separate marks, because they are different rows. A row that does something when it is pressed says so first: the hover tint, the press and the chevron at its trailing edge all follow `Interactive`.",
		specimens: []specimen{
			{caption: "Pressable and static", build: rowPressableAndStatic},
			{caption: "Opened from the keyboard, without opening it twice", build: rowFromTheKeyboard},
			{caption: "The four heights, and the three states", build: rowHeights},
		},
	},
	{
		name:    "KvitSlimRow",
		group:   "Content",
		summary: "The row a reader sees most: a name, what it is, one phrase about it and one figure, right aligned, in that order on every screen.",
		specimens: []specimen{
			{caption: "Measured and unmeasured", build: slimRowMeasured},
		},
	},
	{
		name:    "KvitCard",
		group:   "Content",
		summary: "A bounded block of content on a surface. An interactive card takes the control-boundary token, the hover tint and the chevron a pressable row carries, so a card that responds to a click looks like it will before it is clicked.",
		specimens: []specimen{
			{caption: "Static and interactive", build: cardForms},
		},
	},
	{
		name:    "KvitPanel",
		group:   "Content",
		summary: "A region of the window with its own ground: a sidebar, a toolbar strip. Structural, where a card is content — which is why it has no radius.",
		specimens: []specimen{
			{caption: "With rules on two edges", build: panelWithRules},
		},
	},
	{
		name:    "KvitPane",
		group:   "Content",
		summary: "The side pane: detail about the one thing selected, the same width wherever it appears so the list beside it never reflows.",
		specimens: []specimen{
			{caption: "Open", build: paneOpen},
		},
	},
	{
		name:    "KvitDivider",
		group:   "Content",
		summary: "A rule between two things. One design pixel, in the decorative border token rather than the control-boundary one.",
		specimens: []specimen{
			{caption: "Horizontal and vertical", build: dividerBothWays},
		},
	},
	{
		name:    "KvitDisclosure",
		group:   "Content",
		summary: "A trigger with a body under it. Takes a `Group` so several sections can behave as an accordion, which is what most of the eighteen hand-rolled versions in the estate actually are.",
		specimens: []specimen{
			{caption: "Open and closed", build: disclosureOpenAndClosed},
		},
	},
	{
		name:    "KvitEmptyState",
		group:   "Content",
		summary: "What a view says when it has nothing to show: what would be here, why it is not, and the action that would fill it. Also the answer for a chart with no data, in place of an axis drawn around zeros. The compact form is the same sentence on one line, at the height of a slim row and starting where the rows start, for a section in a stack of sections that may each be empty.",
		specimens: []specimen{
			{caption: "With an action", build: emptyStateWithAction},
			{caption: "Dashed, as a drop target", build: emptyStateDropTarget},
			{caption: "Compact, one line per empty section", build: emptyStateCompact},
		},
	},
	{
		name:    "KvitBadge",
		group:   "Marks",
		summary: "A count attached to something else. Caps rather than growing wide, and hides at zero — a badge showing nought says look here about nothing. The number is drawn and announced with the reader's own digit grouping, and the noun beside it comes from the caller in two slots, a singular and a plural.",
		specimens: []specimen{
			{caption: "Counts, capped, and hidden at zero", build: badgeCounts},
			{caption: "What the count counts", build: badgeNouns},
		},
	},
	{
		name:    "KvitFigure",
		group:   "Quantities",
		summary: "A measured value: tabular numerals, the unit in muted colour at the smaller role, and an em dash where nothing was measured rather than a zero. A balance nobody computed is not a balance of zero.",
		specimens: []specimen{
			{caption: "Measured, unmeasured and bounded", build: figureForms},
		},
	},
	{
		name:    "KvitButton",
		group:   "Controls",
		summary: "A button with words on it, in three forms. One primary per screen region: it is the action the screen is for. `Danger` is separate from the form, because a destructive action can be any of the three. `Explanation` is one sentence saying what the words cannot — why it is disabled, or what pressing it opens — and it is read on a disabled button too.",
		specimens: []specimen{
			{caption: "Text, icon plus text, and busy text in every form", build: buttonForms},
			{caption: "Destructive and disabled", build: buttonDanger},
			{caption: "A button that stays on, beside the same button off", build: buttonChecked},
			{caption: "Why a button cannot be pressed", build: buttonExplanations},
		},
	},
	{
		name:    "KvitStepper",
		group:   "Controls",
		summary: "A number with a minus and a plus beside it, for a small range a reader adjusts by one or two. The range comes from whatever owns it rather than being repeated at the call site.",
		specimens: []specimen{
			{caption: "The interface-size row, pointed at the right setting", build: stepperInterfaceSize},
		},
	},
	{
		name:    "KvitField",
		group:   "Controls",
		summary: "A single line of text. The outline is the control-boundary token, so its edges are visible. An error is a message and a border together, never a border alone.",
		specimens: []specimen{
			{caption: "Resting, filled and in error", build: fieldStates},
		},
	},
	{
		name:    "KvitSearchField",
		group:   "Controls",
		summary: "A field that filters something. Escape clears rather than reverting, and it announces its result count — filtering is the one interaction whose whole outcome happens somewhere else on the screen.",
		specimens: []specimen{
			{caption: "Empty and filtering", build: searchFieldStates},
		},
	},
	{
		name:    "KvitTab",
		group:   "Controls",
		summary: "One tab in a row of them. The selected tab is marked by an underline as well as by colour and weight — selection shown by colour alone is the most common place the rule gets broken. `Explanation` is one sentence saying what the view behind the tab shows, which two or three words cannot.",
		specimens: []specimen{
			{caption: "Selected, counted and plain", build: tabForms},
			{caption: "Each tab saying what its view shows", build: tabExplanations},
		},
	},
	{
		name:    "KvitScrollBar",
		group:   "Flow",
		summary: "A scroll bar, occupying its own strip rather than floating over content. Thirty-two files in the estate have a private version; an overlay bar hides the right-hand column of a table and the last character of every elided label.",
		specimens: []specimen{
			{caption: "Beside a scrolling column", build: scrollBarBesideColumn},
		},
	},
	{
		name:    "KvitSegmented",
		group:   "Flow",
		summary: "One choice from two to five short options, all visible — the shape kvit-cash's dashboard period control needs, where a control governing seven widgets should say what they are showing without being opened.",
		specimens: []specimen{
			{caption: "A period control", build: segmentedPeriod},
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
