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
		name:    "KvitChip",
		group:   "Marks",
		summary: "A small labelled mark saying what kind of thing this is or what state it is in. The tone names a meaning rather than a colour, and every tone carries its outline as well as its tint so the distinction is not resting on hue. It is a mark and not a control: a chip that opens something when it is pressed is KvitChipButton, drawn from this same tone table.",
		specimens: []specimen{
			{caption: "Every tone, tinted and filled", build: chipTones},
			{caption: "A state whose word is not enough to act on", build: chipExplanation},
		},
	},
	{
		name:    "KvitTag",
		group:   "Marks",
		summary: "A label a person put there, in a colour they chose. Theme-independent, because the reader picked that red and would not expect it to change; the label colour is derived from the fill rather than taken from a token.",
		specimens: []specimen{
			{caption: "Tinted, plain and removable", build: tagForms},
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
		name:    "KvitSlug",
		group:   "Marks",
		summary: "An identifier: a reference, a hash, a key. Monospace, because the task is comparison rather than reading, and elided from the middle because the end is what distinguishes one from its neighbours.",
		specimens: []specimen{
			{caption: "With and without a ground", build: slugGrounds},
		},
	},
	{
		name:    "KvitDot",
		group:   "Marks",
		summary: "A small filled circle standing for one thing's state. The shape is the second channel: a level that differs only by hue says nothing to a reader who cannot separate red from amber.",
		specimens: []specimen{
			{caption: "Three shapes, three levels", build: dotLevels},
		},
	},
	{
		name:    "KvitSignal",
		group:   "Marks",
		summary: "A mark saying what state something is in and how many things are in it, for a list where one row may have several of each. The number is drawn only past one, because a column of marks all reading `1` says nothing the mark did not already say. The colour is the caller’s, since which states exist is an application’s own question; the shape and the hollow form are the second channel, so states told apart by hue alone are not.",
		specimens: []specimen{
			{caption: "One of each, and several of one", build: signalCounts},
			{caption: "Beside the words it marks", build: signalBesideWords},
		},
	},
	{
		name:    "KvitPip",
		group:   "Marks",
		summary: "A row of dots standing for a small count — three of five days recorded. For counts a reader takes in without counting; past about seven, the figure is faster.",
		specimens: []specimen{
			{caption: "Three of five, and none of four", build: pipCounts},
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
		name:    "KvitBeforeAfter",
		group:   "Quantities",
		summary: "One value as it stands and the value something proposes to replace it with. Position, colour weight and the arrow all say which is which, so no reader depends on separating the two colours. A record being added has no before, and the em dash says so.",
		specimens: []specimen{
			{caption: "Changed, added, removed and unchanged", build: beforeAfterForms},
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
		name:    "KvitChipButton",
		group:   "Controls",
		summary: "KvitChip's twin for a fact that opens something. Drawn from the same tone table, so a row mixing facts that act with facts that do not reads as one row; what separates them is what a control has anyway — a ground that changes under the pointer, a hand cursor, a focus ring and a button role. A chip that cannot be pressed keeps its place in the tab order and says why, a chip that can may say where it goes, and the chip whose destination is already open is drawn as the current one.",
		specimens: []specimen{
			{caption: "Every tone, tinted, filled and keyboard-focused", build: chipButtonTones},
			{caption: "What pressing it leads to, where the caller knows", build: chipButtonLeadsTo},
			{caption: "Unavailable, with the reason attached", build: chipButtonUnavailable},
			{caption: "Elided in a narrow column, at any interface size", build: chipButtonElided},
			{caption: "The one that is already open", build: chipButtonCurrent},
			{caption: "Where pressing one goes", build: chipButtonExplanation},
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
		name:    "KvitTextArea",
		group:   "Controls",
		summary: "Several lines of text, typed or read. KvitField's outline and error rule in the field form, no ground at all in the plain one, which is for a document filling a pane. `Underlay` draws behind the words in the text's own coordinates, so a wash over a marked passage does not mean replacing the background.",
		specimens: []specimen{
			{caption: "Resting, filled, monospace and in error", build: textAreaStates},
			{caption: "Plain: a document filling a pane, with no ground of its own", build: textAreaPlain},
			{caption: "A wash behind a marked passage, drawn through `Underlay`", build: textAreaUnderlay},
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
		name:    "KvitCheck",
		group:   "Controls",
		summary: "A checkbox in three states. The third — partial — is what a parent row shows when some of its children are checked; drawing that as unchecked loses the information and drawing it as checked is a lie.",
		specimens: []specimen{
			{caption: "Off, on, partial and disabled", build: checkStates},
		},
	},
	{
		name:    "KvitSelect",
		group:   "Controls",
		summary: "A choice from a list too long to lay out. For two to four self-evident options, KvitSegmented is right instead — a three-item dropdown hides two of the three answers for no reason.",
		specimens: []specimen{
			{caption: "A currency picker", build: selectCurrency},
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
		name:    "KvitTooltip",
		group:   "Feedback",
		summary: "A short label next to a control after a pause. Never the only place a control's meaning lives — a control explained only by its tooltip is unusable on a keyboard, a screen reader and a touch screen.",
		specimens: []specimen{
			{caption: "On a button", build: tooltipOnButton},
		},
	},
	{
		name:    "KvitPopover",
		group:   "Feedback",
		summary: "A small surface anchored to a control holding something to act on. Takes focus, closes on Escape and on a click outside — which is what separates it from a hover card.",
		specimens: []specimen{
			{caption: "Open, holding a form", build: popoverWithForm},
		},
	},
	{
		name:    "KvitHint",
		group:   "Feedback",
		summary: "An information-icon trigger for an explanation too long for a tooltip. Click or keyboard activation keeps its KvitPopover open for reading, and Escape or an outside click dismisses it.",
		specimens: []specimen{
			{caption: "Open for a longer explanation", build: hintOpen},
		},
	},
	{
		name:    "KvitHoverCard",
		group:   "Feedback",
		summary: "More about the thing under the pointer. Read-only, always: a surface that appears on hover is unreachable by keyboard and by touch, so nothing inside one may be the only route to an action.",
		specimens: []specimen{
			{caption: "Showing what a row could not fit", build: hoverCardRow},
		},
	},
	{
		name:    "KvitToast",
		group:   "Feedback",
		summary: "A message that appears, says one thing and goes away. For confirming what already happened. A toast with an action stays until dismissed — an undo that times out mid-read is an undo the reader cannot use.",
		specimens: []specimen{
			{caption: "Four tones, and one with an undo", build: toastTones},
		},
	},
	{
		name:    "KvitNotice",
		group:   "Feedback",
		summary: "A message that stays: a condition that is still true, where a toast is an acknowledgement of something that finished. Dismissible only when dismissing it is meaningful.",
		specimens: []specimen{
			{caption: "A warning that can be dismissed, and an error that cannot", build: noticeForms},
		},
	},
	{
		name:    "KvitDialog",
		group:   "Feedback",
		summary: "A modal surface that has to be answered first. Expensive, and worth it only where the answer really does come first; the confirming button is on the right, and a destructive dialog takes no keyboard default. There are three ways out, because different readers find different ones: the close symbol at the right of the title, the cancel button, and Escape.",
		specimens: []specimen{
			{caption: "A destructive confirmation", build: dialogDestructive},
		},
	},
	{
		name:    "KvitFloatingView",
		group:   "Feedback",
		summary: "Detail about the one thing selected, held over the list rather than beside it: a centred card with a name and a close control, and the list dimmed behind. Unlike KvitPane it takes no width out of the layout, so the list keeps the whole window at every size; unlike KvitDialog it is closed rather than answered, and it floats over one region rather than over the window. From kvit-cash's copy of the Qt library.",
		specimens: []specimen{
			{caption: "A record held over the list it came from", build: floatingRecord},
		},
	},
	{
		name:    "KvitBar",
		group:   "Data",
		summary: "One quantity against a stated scale. A bar with no value draws a tick rather than a zero-width fill; a bounded figure is hatched, and the hatch survives grayscale; and the scale is required, because a bar drawn against its own list rescales invisibly.",
		specimens: []specimen{
			{caption: "Measured, bounded and unmeasured", build: barForms},
		},
	},
	{
		name:    "KvitStackedBar",
		group:   "Data",
		summary: "Several quantities adding to one total. The two-pixel gap is the reason this is a component: adjacent fills that touch read as one fill with a colour change in it.",
		specimens: []specimen{
			{caption: "A four-way breakdown", build: stackedBarBreakdown},
		},
	},
	{
		name:    "KvitSpark",
		group:   "Data",
		summary: "The shape of a series, small enough to sit in a row. A period that was never measured is a baseline tick rather than a zero-height bar: a missing week and a quiet week are different facts.",
		specimens: []specimen{
			{caption: "With a hole in the middle", build: sparkWithHole},
		},
	},
	{
		name:    "KvitNetFlow",
		group:   "Data",
		summary: "Two opposed series over the same periods, on one baseline and one scale, with what they come to drawn as a line. What comes in is drawn up from the baseline and what goes out is drawn down, so which of the two was larger in a period is a glance rather than a comparison between two strips an inch apart. Each column names its own period underneath, and the names thin out to whatever spacing they actually fit in. From kvit-cash's copy of the Qt library.",
		specimens: []specimen{
			{caption: "A year of money in and money out, netted off", build: netFlowYear},
			{caption: "Daily columns, where the names thin out, and a period nobody measured", build: netFlowDays},
		},
	},
	{
		name:    "KvitTrend",
		group:   "Data",
		summary: "A series with a value axis and a hover crosshair — read for values, where a spark is read for shape. A gap in the data draws as a gap: interpolating over a hole asserts values nobody measured. A second series is dashed as well as differently coloured, and both are named in the key and in the crosshair.",
		specimens: []specimen{
			{caption: "A series, and the empty state", build: trendAndEmpty},
			{caption: "Two series, with an annotation drawn over the plot", build: trendTwoSeries},
		},
	},
	{
		name:    "KvitDistribution",
		group:   "Data",
		summary: "How a set of values is spread. An average of four days and an average made of one twenty-day outlier are the same number and different situations.",
		specimens: []specimen{
			{caption: "Two rows on one scale", build: distributionRows},
		},
	},
	{
		name:    "KvitGauge",
		group:   "Data",
		summary: "How much of an allowance is used, with the target marked and a pace mark saying where an even rate would have reached. Sixty percent spent is fine on day eighteen and a problem on day six.",
		specimens: []specimen{
			{caption: "Ahead of pace, behind it, and over", build: gaugePace},
		},
	},
	{
		name:    "KvitDelta",
		group:   "Data",
		summary: "How much something changed and in which direction, with an arrow as well as a colour. Whether up is good is the caller's to say: a rise in spending and a rise in savings are the same arrow and opposite colours.",
		specimens: []specimen{
			{caption: "Up, down and unchanged", build: deltaDirections},
		},
	},
	{
		name:    "KvitStatTile",
		group:   "Data",
		summary: "A card carrying one figure, what it is, how it changed and its recent shape — the shape of six of kvit-cash's seven dashboard widgets. The order is fixed so a row of tiles can be scanned one part at a time.",
		specimens: []specimen{
			{caption: "A dashboard row", build: statTileRow},
		},
	},
	{
		name:    "KvitFigureBlock",
		group:   "Data",
		summary: "A figure with its name under it: one number a reader is meant to take away. The number is above and larger, because a row of these is read across the numbers.",
		specimens: []specimen{
			{caption: "A row of three", build: figureBlockRow},
		},
	},
	{
		name:    "KvitCell",
		group:   "Data",
		summary: "One cell of a table, drawn according to what kind of value its column holds. The kind comes from the column, so every cell in it aligns the same way and says the same thing about a missing value. A cell formats nothing: Money and Figure are given the string they draw.",
		specimens: []specimen{
			{caption: "The six kinds a column of values can be", build: cellKinds},
			{caption: "Several states at once, and a row the reader picks", build: cellStatesAndPick},
			{caption: "A value the column cut short", build: cellCutShort},
		},
	},
	{
		name:    "KvitTable",
		group:   "Data",
		summary: "A dense, configurable table over a model. Columns sort, resize, move and open a menu from their own header, by keyboard as well as pointer; a row opens its record on one press; and a Check column draws a box per row with a box in its header for every row shown. The model says what each column is and answers each cell with its value, its unit, its marks, whether its box is ticked and the whole of a value the column cut short, which the cell then discloses on hover and under the keyboard cursor. Holds smooth scrolling and sub-100 ms filtering at 250,000 rows with twelve columns of value.",
		specimens: []specimen{
			{caption: "Two hundred and fifty thousand rows", build: tableQuarterMillion},
			{caption: "Sorting, a column menu and a column of boxes", build: tableSortingAndBoxes},
			{caption: "Nothing to show", build: tableEmpty},
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
		name:    "KvitMenu",
		group:   "Flow",
		summary: "A list of commands. Eighteen private versions in the estate, and what they mostly get wrong is the same two things: no keyboard route in, and no separator before the destructive item. Opened by a button, a right-click, the Menu key or Shift+F10; drawn by Kvit on Windows and Linux and by the system on macOS.",
		specimens: []specimen{
			{caption: "Opened, with shortcuts and a destructive item", build: menuOpened},
		},
	},
	{
		name:    "KvitMenuItem",
		group:   "Flow",
		summary: "One line of a menu, carrying its shortcut on the right — which is how a menu teaches a faster route to a reader who keeps using it.",
		specimens: []specimen{
			{caption: "Ordinary, disabled and destructive", build: menuItemForms},
		},
	},
	{
		name:    "KvitTree",
		group:   "Flow",
		summary: "A nested list the reader can open and close. Twelve private versions; what a shared one has to get right is the keyboard, because depth without Left and Right is a wall.",
		specimens: []specimen{
			{caption: "A small hierarchy", build: treeAccounts},
		},
	},
	{
		name:    "KvitSwitch",
		group:   "Flow",
		summary: "An option that takes effect the moment it moves — where a checkbox is a value in a form that takes effect on submit. The knob moves and the track fills, so the state is not resting on hue.",
		specimens: []specimen{
			{caption: "On, off and disabled", build: switchForms},
		},
	},
	{
		name:    "KvitRadioGroup",
		group:   "Flow",
		summary: "One choice from a handful where the choice needs explaining. Arrow keys move within the group and Tab leaves it, which a column of separate controls does not do.",
		specimens: []specimen{
			{caption: "Three options with detail", build: radioAppearance},
		},
	},
	{
		name:    "KvitProgress",
		group:   "Flow",
		summary: "How far through something the application is. An unknown total is a moving band rather than a bar creeping toward the end, and it carries a label saying how far through what.",
		specimens: []specimen{
			{caption: "Determinate and indeterminate", build: progressForms},
		},
	},
	{
		name:    "KvitSlider",
		group:   "Flow",
		summary: "A value chosen by dragging, for a feel rather than a number. Where the exact value matters, KvitStepper or KvitNumberField are right — hitting a particular number on a slider is hard and it cannot be typed.",
		specimens: []specimen{
			{caption: "With its value shown", build: sliderWithValue},
		},
	},
	{
		name:    "KvitSplitView",
		group:   "Flow",
		summary: "Two regions the reader can resize. The handle is a wide invisible strip with a hairline down the middle, so the target is comfortable and the rule is still thin; it also moves with the arrow keys.",
		specimens: []specimen{
			{caption: "Two panes", build: splitTwoPanes},
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
	{
		name:    "KvitTypeAhead",
		group:   "Flow",
		summary: "A field offering matches as the reader types, for a list too long to read. Whether the reader may make a new value has no default: a tag picker should let the reader invent one, a category picker should not.",
		specimens: []specimen{
			{caption: "A category picker that will not invent categories", build: typeAheadCategory},
		},
	},
	{
		name:    "KvitConfirmInPlace",
		group:   "Flow",
		summary: "A strip saying what just happened, with the undo inside it. Cheaper than a dialog before every action when the reader is doing the same thing forty times — and right only when the action can really be undone.",
		specimens: []specimen{
			{caption: "After a bulk edit", build: confirmAfterBulkEdit},
		},
	},
	{
		name:    "KvitTimeline",
		group:   "Flow",
		summary: "What happened to something, newest first, with who did it. In an estate where an agent and a person change the same things, who is the column that makes a history worth reading.",
		specimens: []specimen{
			{caption: "An account's recent history", build: timelineHistory},
		},
	},
	{
		name:    "KvitNumberField",
		group:   "Flow",
		summary: "A number the reader types, right-aligned in tabular numerals. It validates and says why rather than refusing keystrokes — a field that ignores a key gives no reason, and the usual cause is a decimal separator the reader's locale writes differently.",
		specimens: []specimen{
			{caption: "Integer and decimal, valid and out of range", build: numberFieldForms},
		},
	},
	{
		name:    "KvitMoneyField",
		group:   "Flow",
		summary: "An amount of money, in minor units. The reader types 12.34 and the field reports 1234, an integer — money in floating point drifts by a penny somewhere nobody can find. The decimal count comes from the currency.",
		specimens: []specimen{
			{caption: "Sterling, yen and dinar", build: moneyFieldCurrencies},
		},
	},
	{
		name:    "KvitDualList",
		group:   "Flow",
		summary: "Two lists with items moving between them: what is available, and what is chosen and in what order. Everything works from the keyboard, which most implementations of this shape do not.",
		specimens: []specimen{
			{caption: "Choosing table columns", build: dualListColumns},
		},
	},
	{
		name:    "KvitSpotlight",
		group:   "Flow",
		summary: "Darken everything except one region and say something about it — the primitive under a guided tour, and deliberately only the primitive: what drives the stepping differs between the two applications that want one.",
		specimens: []specimen{
			{caption: "Focusing a button", build: spotlightOnButton},
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
