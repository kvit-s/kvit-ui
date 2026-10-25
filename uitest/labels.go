package uitest

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"unicode"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// Problem is one place a screen or an app breaks one of the interface rules
// in docs/ui-guidelines/rules.md.
type Problem struct {
	// Rule is the rule broken, "WRD-3".
	Rule string
	// Where is the thing that breaks it, by kind and words: `button "Delete
	// the file"`.
	Where string
	// What says what is wrong.
	What string
	// Form names the form the screen was read in, "" for a problem of the
	// whole app.
	Form string
}

// String is the problem on one line: "WRD-3 (Windows form): button "Delete
// the file": an article in a command's words".
func (p Problem) String() string {
	if p.Form == "" {
		return fmt.Sprintf("%s: %s: %s", p.Rule, p.Where, p.What)
	}
	return fmt.Sprintf("%s (%s form): %s: %s", p.Rule, p.Form, p.Where, p.What)
}

// Verbs are the words a command's label may start with (WRD-1): specific
// verbs naming what happens. An app adds its own after review; "configure",
// "manage" and "set" are left out on purpose.
var Verbs = wordSet(`accept add allow decrease increase pull push renew apply approve archive ask assign attach back block browse build calculate
cancel capture change check choose clear close collapse collect compare connect continue copy create cut decline
delete deny deselect detach disable discard dismiss done download draw duplicate edit email empty enable end enter
erase exit expand explain export fetch filter find finish fix flag follow format forward generate get go group hide
highlight hold ignore import insert install invite join keep label launch leave link list load lock log mark merge
minimize move mute narrow new next note notify open order paste pause pick pin play post preview print publish
quit quote rate read reconcile reconnect record redo refresh refuse reject reload remove rename reopen reorder
repeat replace reply report request reset resize restore resume retry return reveal review revert rotate run save
scan schedule scroll search select send share show shrink sign skip snooze sort split star start step stop store
submit subscribe suspend swap switch sync tag take test track trash trim try turn undo unarchive unlink unlock
unmute unpin unstar unsubscribe update upgrade upload use validate verify view wake watch widen withdraw wrap
write zoom`)

// Verbless are the labels that may stand without a verb (WRD-1): the
// standard ones on the Windows form, and the formatting commands.
var Verbless = wordSet(`ok yes no about advanced back details forward full less more new next personalize previous
properties settings help done bold italic underline strikethrough code heading`)

// ProperNouns are words a sentence-case label may capitalise inside it, such
// as the app's own name; an app adds its own.
var ProperNouns = wordSet(`Kvit Windows Linux Mac macOS Apple Microsoft`)

func wordSet(words string) map[string]bool {
	out := map[string]bool{}
	for _, w := range strings.Fields(words) {
		out[w] = true
	}
	return out
}

// Labels lists where a window's screen breaks the interface rules on words,
// dialogs, menus, tooltips, help links and the Tab order (rules.md, the
// "test, one screen" entries), in the form its UI follows now. It reads
// what is shown: the controls in the window, its menu bar's menus and any
// menu open. Call it on the interface's thread, as Session.Do does.
func Labels(w *kvitui.Window) []Problem {
	c := &labelCheck{ui: w.UI(), form: w.UI().Form(), window: w, primaries: map[*unison.Panel]int{},
		tabRows: map[*unison.Panel]int{}, helpButtons: map[*unison.Panel]int{}, settingRegions: map[*unison.Panel]bool{},
		commitButtons: map[*unison.Panel][]string{}}
	c.walk(w.Content(), w.Content(), false)
	for region, words := range c.commitButtons {
		if c.settingRegions[region] {
			for _, word := range words {
				c.add("WIN-23", "button "+quote(word), "on a page of settings, which apply each change at once")
			}
		}
	}
	for _, n := range c.helpButtons {
		if n > 1 {
			c.add("HLP-6", "a window", fmt.Sprintf("%d help buttons; the macOS form has one per window, for the window's topic", n))
		}
	}
	for _, n := range c.primaries {
		if n > 1 {
			c.add("CTL-1", "a region", fmt.Sprintf("%d primary buttons; one per region tells the reader which to press", n))
		}
	}
	for _, n := range c.tabRows {
		limit := 7
		if c.form == kvitui.FormMac {
			limit = 6
		}
		switch {
		case n == 1:
			c.add("CTL-10", "a tab row", "one tab alone")
		case n > limit:
			c.add("CTL-10", "a tab row", fmt.Sprintf("%d tabs, more than the %d the form allows", n, limit))
		}
	}
	if b := w.MenuBar(); b != nil {
		hasEdit := false
		for _, m := range b.Menus() {
			title, _, _ := kvitui.AccessText(m.Title)
			hasEdit = hasEdit || title == "Edit"
			if title != c.ui.AppName && title != "App" {
				c.caseOf("menu title", c.ui.CommandCase(title), true)
			}
			c.menu("menu "+quote(title), c.ui.MenuLines(m.Items, false), false, 0)
		}
		c.hasEdit = hasEdit
	}
	for _, m := range w.ShownMenus() {
		c.menu("menu "+quote(m.Title), m.Items, m.Context, 0)
	}
	// KEY-15: on the Windows form every control with words has an access
	// key, and no two enabled controls share one the app chose; kvit-ui
	// shares a letter only when every letter of a control's words is taken.
	chosen := map[rune]string{}
	for _, k := range w.AccessKeys() {
		where := "control " + quote(k.Words)
		switch {
		case k.Key == 0:
			c.add("KEY-15", where, "no letter left for an access key")
		case k.Chosen && k.Panel.Enabled():
			if other, ok := chosen[k.Key]; ok {
				c.add("KEY-15", where, fmt.Sprintf("its access key %c is also %s's", unicode.ToUpper(k.Key), quote(other)))
			}
			chosen[k.Key] = k.Words
		}
	}
	c.tabOrder(w)
	if c.hasEdit {
		for _, name := range c.clipboardButtons {
			c.add("CMD-21", "button "+quote(name), "a Cut, Copy or Paste button on a screen with an Edit menu")
		}
	}
	return c.out
}

type labelCheck struct {
	ui               *kvitui.UI
	form             kvitui.Form
	window           *kvitui.Window
	out              []Problem
	primaries        map[*unison.Panel]int
	tabRows          map[*unison.Panel]int
	helpButtons      map[*unison.Panel]int // the macOS form's help buttons, by window or dialog
	scope            *unison.Panel         // the dialog being read, nil for the window
	inProperties     bool                  // the dialog being read is a properties dialog
	settingRegions   map[*unison.Panel]bool
	commitButtons    map[*unison.Panel][]string // OK, Cancel and Apply outside dialogs, by region
	clipboardButtons []string
	hasEdit          bool
	menuVerbed       bool // reading a submenu whose title starts with a verb
}

func (c *labelCheck) add(rule, where, what string) {
	c.out = append(c.out, Problem{Rule: rule, Where: where, What: what, Form: c.form.String()})
}

// walk reads a panel and what it holds; region is the panel the nearest
// region, dialog, popover or pane makes, for one primary button per region.
func (c *labelCheck) walk(p *unison.Panel, region *unison.Panel, inDialog bool) {
	if p.Hidden {
		return
	}
	switch v := p.Self.(type) {
	case *kvitui.Region, *kvitui.Popover, *kvitui.Pane, *kvitui.Card:
		region = p
	case *kvitui.Dialog:
		region, inDialog = p, true
		c.dialog(v)
		// A dialog's help button is counted apart from the window's.
		outer, props := c.scope, c.inProperties
		c.scope, c.inProperties = p, v.Kind == kvitui.DialogProperties
		defer func() { c.scope, c.inProperties = outer, props }()
	case *kvitui.SettingRow:
		c.settingRegions[region] = true
	case *kvitui.Button:
		c.button(v, region, inDialog)
	case *kvitui.MenuButton:
		c.button(&v.Button, region, inDialog)
		lines := c.ui.MenuLines(v.Items, false)
		if n := countLines(lines); n < 3 {
			c.add("CMD-17", "menu button "+quote(v.Words()), fmt.Sprintf("a menu of %d; one or two commands are buttons of their own", n))
		}
		c.menu("menu button "+quote(v.Words()), lines, false, 0)
	case *kvitui.IconButton:
		c.iconButton(v)
	case *kvitui.Tab:
		c.tabRows[p.Parent()]++
		words := c.ui.CommandCase(v.Text)
		c.heading("tab", words)
		c.caseOf("tab", words, true)
	case *kvitui.Segmented:
		opts := v.Options()
		if len(opts) > 7 {
			c.add("CTL-3", "segmented "+quote(v.Label), fmt.Sprintf("%d segments; past about five a select is the control", len(opts)))
		}
		for _, o := range opts {
			words := c.ui.CommandCase(o.Label)
			c.heading("segment", words)
			c.caseOf("segment", words, true)
		}
	case *kvitui.RadioGroup:
		opts := v.Options()
		if len(opts) < 2 || len(opts) > 5 {
			c.add("CTL-2", "radio group "+quote(v.Label), fmt.Sprintf("%d options; a group holds two to five", len(opts)))
		}
		seen := map[string]bool{}
		for _, o := range opts {
			if seen[o.Label] {
				c.add("WRD-29", "radio group "+quote(v.Label), "two options named "+quote(o.Label))
			}
			seen[o.Label] = true
			if o.Detail != "" && strings.EqualFold(o.Detail, o.Label) {
				c.add("WRD-22", "radio option "+quote(o.Label), "its detail repeats its label")
			}
		}
	case *kvitui.SearchField:
		if v.Placeholder == "Filter…" {
			c.add("CTL-15", "search field "+quote(v.Label), "the bare default placeholder; say what it searches")
		}
		c.fieldLabel(&v.Field)
	case *kvitui.Field:
		c.fieldLabel(v)
		if c.form == kvitui.FormWindows && v.Placeholder != "" && !dashPlaceholder(v.Placeholder) && endsInPunctuation(v.Placeholder) {
			c.add("WRD-19", "field "+quote(v.Label), "its placeholder ends in punctuation")
		}
	case *kvitui.Check:
		c.optionLabel("check box", v.Text)
	case *kvitui.Switch:
		c.optionLabel("switch", v.Text)
	case *kvitui.Select:
		// A select's words are the value in force, read as a value and never
		// as a command: no WRD-1, in either form.
	case *kvitui.Link:
		c.link(v)
	case *kvitui.HelpLink:
		c.helpLink(v)
	case *kvitui.Notice:
		c.title("notice", v.Text)
	case *kvitui.Label:
		c.label(v, p)
	case *kvitui.Toolbar:
		if inDialog {
			c.add("CMD-18", "toolbar "+quote(v.Label), "a toolbar inside a dialog")
		}
		cmds := v.Commands()
		for i, cmd := range cmds {
			if slices.Index(cmds, cmd) != i {
				c.add("CMD-18", "toolbar "+quote(v.Label), "a command twice")
			}
		}
	case *kvitui.StatusBar:
		if inDialog {
			c.add("CMD-18", "status bar", "a status bar inside a dialog")
		}
	}
	for _, child := range p.Children() {
		c.walk(child, region, inDialog)
	}
}

// command reads a command's words, as on a button or a menu line.
func (c *labelCheck) command(where, words string, verbChecked bool) {
	plain := strings.TrimSuffix(words, "…")
	ws := strings.Fields(plain)
	if len(ws) == 0 {
		return
	}
	if verbChecked {
		first := strings.ToLower(strings.Trim(ws[0], "&"))
		first = strings.ReplaceAll(first, "&", "")
		if !Verbs[first] && !Verbless[first] && !Verbless[strings.ToLower(strings.ReplaceAll(plain, "&", ""))] {
			c.add("WRD-1", where, "does not start with a verb naming what happens")
		}
	}
	if len(ws) > 3 {
		c.add("WRD-2", where, fmt.Sprintf("%d words; a command is one to three", len(ws)))
	}
	for _, w := range ws {
		switch strings.ToLower(strings.Trim(w, ".,;:")) {
		case "a", "an", "the", "my", "your":
			c.add("WRD-3", where, "an article or possessive in a command's words")
		}
	}
	if strings.HasSuffix(plain, "...") {
		c.add("WRD-4", where, "three full stops where the ellipsis character belongs")
	}
	if strings.HasSuffix(plain, ":") {
		c.add("WRD-19", where, "ends in a colon")
	}
	if strings.HasSuffix(plain, ".") {
		c.add("WRD-27", where, "ends in a full stop")
	}
	if strings.Contains(words, "  ") {
		c.add("WRD-27", where, "two spaces")
	}
	if strings.HasPrefix(strings.ToLower(words), "click") {
		c.add("A11Y-1", where, "starts with Click")
	}
	c.caseOf(where, words, true)
}

// caseOf checks a label of a kind the macOS form writes in title case is in
// the form's case (WRD-6).
func (c *labelCheck) caseOf(where, words string, titled bool) {
	if !titled || words == "" {
		return
	}
	if c.form == kvitui.FormMac {
		if kvitui.TitleCase(words) != words {
			c.add("WRD-6", where, "not in title case on the macOS form; write it in sentence case and let kvit-ui convert it")
		}
		return
	}
	ws := strings.Fields(strings.ReplaceAll(words, "&", ""))
	for _, w := range ws[1:] {
		rs := []rune(strings.Trim(w, "…().,"))
		if len(rs) < 2 || !unicode.IsUpper(rs[0]) || ProperNouns[string(rs)] || kvitui.TitleCaseKept[strings.ToLower(string(rs))] {
			continue
		}
		rest := string(rs[1:])
		if strings.ToLower(rest) == rest {
			c.add("WRD-6", where, "title case on the Windows form, which writes sentence case")
			return
		}
	}
}

func (c *labelCheck) button(b *kvitui.Button, region *unison.Panel, inDialog bool) {
	words := b.Words()
	where := "button " + quote(words)
	if b.Busy {
		switch strings.TrimSpace(b.BusyText) {
		case "", "Working…", "Working...", "Please wait", "Please wait…", "Loading…", "Loading...":
			c.add("WRD-25", where, "busy words that do not name the work")
		}
		return
	}
	if b.Form == kvitui.ButtonPrimary && b.Enabled() {
		c.primaries[region]++
	}
	switch strings.ToLower(strings.TrimSuffix(words, "…")) {
	case "help", "troubleshoot":
		c.add("HLP-6", where, "a generic help button; give the place a HelpLink to its own topic")
	}
	// WIN-23: Apply only in a properties dialog, and only on the Windows
	// form; OK and Cancel never on a page of settings, which apply at once.
	switch {
	case words == "Apply" && !(inDialog && c.inProperties):
		c.add("WIN-23", where, "Apply outside a properties dialog; elsewhere a change applies at once or on the confirming button")
	case words == "Apply" && c.form == kvitui.FormMac:
		c.add("WIN-23", where, "Apply on the macOS form, which has none")
	}
	if !inDialog && (words == "OK" || words == "Cancel" || words == "Apply") {
		c.commitButtons[region] = append(c.commitButtons[region], words)
	}
	rs := []rune(strings.TrimSpace(words))
	onlyPunct := len(rs) > 0 && !slices.ContainsFunc(rs, unicode.IsLetter)
	if len(rs) > 0 && (len(rs) <= 2 || onlyPunct) && b.Symbol == "" && words != "OK" && words != "No" {
		c.add("CMD-22", where, "letters or punctuation standing in for a symbol; use an IconButton with a Label")
	}
	if strings.HasPrefix(words, "+ ") || strings.HasPrefix(words, "- ") {
		c.add("CMD-22", where, "a sign typed before the words; give the button the plus symbol instead")
	}
	if len([]rune(strings.TrimSuffix(words, "…"))) > 26 {
		c.add("WRD-2", where, "over 26 characters")
	}
	// A dialog's own buttons are read by the dialog rules; a standard
	// command's words are the platform's; a button that turns something on
	// and off names what it turns on, as Bold does, the way a menu's switch
	// line does.
	std := b.Command != nil && b.Command.Std != kvitui.StdNone
	toggle := b.Checkable || b.Checked || (b.Command != nil && b.Command.Checked != nil)
	if !inDialog {
		c.command(where, words, !std && !toggle)
	} else {
		c.caseOf(where, words, true)
	}
	explanation := b.Explanation
	if explanation == "" && b.Command != nil {
		explanation = b.Command.Explanation
	}
	if explanation != "" && strings.EqualFold(strings.TrimSuffix(explanation, "."), strings.TrimSuffix(words, "…")) {
		c.add("WRD-13", where, "its explanation repeats its words")
	}
	if !b.Enabled() && explanation == "" {
		c.add("CTL-8", where, "disabled without an explanation saying why")
	}
	switch strings.ToLower(strings.TrimSuffix(words, "…")) {
	case "cut", "copy", "paste":
		c.clipboardButtons = append(c.clipboardButtons, words)
	}
}

func (c *labelCheck) iconButton(b *kvitui.IconButton) {
	name := b.Name()
	where := "symbol button " + quote(name)
	if name == "" {
		return // CheckNamed reports it
	}
	std := b.Command != nil && b.Command.Std != kvitui.StdNone
	plain := strings.TrimSuffix(name, "…")
	if !std {
		first := strings.ToLower(strings.Fields(plain)[0])
		if !Verbs[first] && !Verbless[first] && !Verbless[strings.ToLower(plain)] {
			c.add("WRD-1", where, "does not start with a verb naming what happens")
		}
	}
	if n := len(strings.Fields(plain)); n > 3 {
		c.add("WRD-2", where, fmt.Sprintf("%d words; a command is one to three", n))
	}
	for _, w := range strings.Fields(plain) {
		switch strings.ToLower(w) {
		case "a", "an", "the", "my", "your":
			c.add("WRD-3", where, "an article or possessive in a command's words")
		}
	}
	tip := b.TooltipWords()
	if b.TooltipEnabled && tip == "" {
		c.add("HLP-1", where, "no tooltip")
	}
	if c.form == kvitui.FormMac {
		first, _, _ := strings.Cut(tip, "\n")
		if len([]rune(first)) > 60 {
			c.add("WRD-13", where, "a tooltip over 60 characters on the macOS form")
		}
		if strings.HasSuffix(first, ".") {
			c.add("WRD-13", where, "a tooltip ending in a full stop on the macOS form")
		}
	}
	least := float32(c.ui.Interface.Px(20))
	if r := b.FrameRect(); r.Width > 0 && (r.Width < least-0.5 || r.Height < least-0.5) {
		c.add("PTR-14", where, "smaller than 20 design pixels square")
	}
	switch strings.ToLower(plain) {
	case "cut", "copy", "paste":
		c.clipboardButtons = append(c.clipboardButtons, name)
	}
}

func (c *labelCheck) dialog(d *kvitui.Dialog) {
	where := "dialog " + quote(d.Title)
	if strings.TrimSpace(d.Title) == "" {
		c.add("WIN-6", where, "no title")
	}
	c.title("dialog", d.Title)
	if strings.HasSuffix(d.Title, ".") {
		c.add("WRD-26", where, "its title ends in a full stop")
	}
	if d.Detail != "" && d.Title != "" && strings.HasPrefix(strings.ToLower(d.Detail), strings.ToLower(d.Title)) {
		c.add("WRD-26", where, "its detail begins with its title's words")
	}
	if !d.SafeButton() {
		c.add("WIN-3", where, "no safe button that leaves without acting")
	}
	words := d.ButtonWords()
	if len(words) > 3 && d.Kind != kvitui.DialogProperties {
		c.add("WIN-2", where, fmt.Sprintf("%d buttons; an alert has at most three", len(words)))
	}
	has := func(w string) bool { return slices.Contains(words, w) }
	for _, w := range words {
		switch w {
		case "Confirm", "Ok", "Okay":
			c.add("WRD-7", where, quote(w)+" in place of the action's name")
		}
	}
	if c.form == kvitui.FormMac {
		for _, w := range []string{"OK", "Yes", "No"} {
			if has(w) && !(w == "OK" && d.Kind == kvitui.DialogInform && len(words) == 1) && d.Kind != kvitui.DialogProperties {
				c.add("WRD-7", where, quote(w)+" on the macOS form, which names the action")
			}
		}
	} else {
		if has("OK") && d.Kind != kvitui.DialogProperties {
			c.add("WRD-7", where, "OK outside a properties dialog")
		}
		if has("Yes") != has("No") {
			c.add("WRD-7", where, "Yes and No only as a pair")
		}
	}
	if has("Yes") && has("Cancel") && !has("No") {
		c.add("WRD-8", where, "Yes beside Cancel without No")
	}
	if d.Kind == kvitui.DialogAsk && d.CancelText != "" {
		switch {
		case d.CancelText == "Cancel", d.CancelText == "Close", d.CancelText == "Stop":
		case c.form == kvitui.FormWindows && strings.HasPrefix(d.CancelText, "Don't "):
		default:
			c.add("WRD-8", where, "the abandoning button is "+quote(d.CancelText)+", not Cancel, Close or Stop")
		}
	}
}

// title checks the title of a dialog or a notice (MSG-4).
func (c *labelCheck) title(kind, title string) {
	where := kind + " " + quote(title)
	switch strings.ToLower(strings.TrimSpace(title)) {
	case "error", "warning", "caution", "alert":
		c.add("MSG-4", where, "a title that names no problem")
	}
	if len([]rune(title)) > 100 {
		c.add("MSG-4", where, "a title over 100 characters")
	}
}

// heading checks a tab, segment or column heading (WRD-20, WRD-27).
func (c *labelCheck) heading(kind, words string) {
	where := kind + " " + quote(words)
	if endsInPunctuation(words) {
		c.add("WRD-20", where, "ends in punctuation")
	}
	if n := len(strings.Fields(words)); n > 3 {
		c.add("WRD-20", where, fmt.Sprintf("%d words; a %s is one to three", n, kind))
	}
}

func (c *labelCheck) fieldLabel(f *kvitui.Field) {
	if strings.HasSuffix(strings.TrimSpace(f.Label), ":") {
		c.add("WRD-19", "field "+quote(f.Label), "its label ends in a colon")
	}
}

func (c *labelCheck) optionLabel(kind, words string) {
	where := kind + " " + quote(words)
	if strings.HasSuffix(strings.TrimSpace(words), ":") {
		c.add("WRD-19", where, "ends in a colon")
	}
	if strings.HasSuffix(words, ".") {
		c.add("WRD-27", where, "ends in a full stop")
	}
	if (strings.HasPrefix(words, "Don't") || strings.HasPrefix(words, "Do not")) && !strings.Contains(words, "show") {
		c.add("WRD-29", where, "a negative option; name what it turns on")
	}
}

func (c *labelCheck) link(l *kvitui.Link) {
	where := "link " + quote(l.Text)
	switch strings.ToLower(strings.TrimSpace(l.Text)) {
	case "here", "click here", "more", "link", "learn more", "read more":
		c.add("WRD-18", where, "words that do not say where it goes")
	}
	if strings.HasSuffix(l.Text, ".") {
		c.add("WRD-18", where, "ends in a full stop")
	}
}

// tabOrder reads the window's Tab stops in the order Tab visits them, as
// unison does, inside the topmost dialog or other popup that keeps the
// keyboard while one is open and leaving out the other popups, and reports a step that goes back against the
// reading order (KEY-9): up within the same column, or to the left within
// the same row. A step to a column further up, as from the sidebar's foot
// to the top of the body, starts the next column and is reading order; a
// help link may come after the buttons it sits above.
func (c *labelCheck) tabOrder(w *kvitui.Window) {
	root := w.Content()
	popups := map[*unison.Panel]bool{}
	for _, p := range slices.Backward(w.Popups()) {
		popups[p.Panel.AsPanel()] = true
		if root == w.Content() && !p.Passive && !p.Ambient && (p.Modal || p.HoldsKeyboard) {
			root = p.Panel.AsPanel()
		}
	}
	type stop struct {
		p   *unison.Panel
		box geom.Rect
	}
	var stops []stop
	var walk func(p *unison.Panel, help bool)
	walk = func(p *unison.Panel, help bool) {
		// A menu, tooltip or card shown over the page is no part of its Tab
		// order.
		if p.Hidden || !p.Enabled() || (popups[p] && p != root) {
			return
		}
		if _, ok := p.Self.(*kvitui.HelpLink); ok {
			help = true
		}
		if p.Focusable() && !help {
			stops = append(stops, stop{p, p.RectToRoot(p.ContentRect(true))})
		}
		for _, ch := range p.Children() {
			walk(ch, help)
		}
	}
	walk(root, false)
	const slack = 1
	for i := 1; i < len(stops); i++ {
		a, b := stops[i-1].box, stops[i].box
		if a.Empty() || b.Empty() {
			continue
		}
		sameRow := b.Y < a.Bottom()-slack && b.Bottom() > a.Y+slack
		sameColumn := b.X < a.Right()-slack && b.Right() > a.X+slack
		var what string
		switch {
		case sameRow && b.Right() <= a.X+slack:
			what = "Tab goes back to the left, from " + stopName(stops[i-1].p)
		case !sameRow && sameColumn && b.Bottom() <= a.Y+slack:
			what = "Tab goes back up, from " + stopName(stops[i-1].p)
		default:
			continue
		}
		c.add("KEY-9", stopName(stops[i].p), what+"; Tab follows the reading order")
	}
}

// stopName names a Tab stop in a problem: its kind and words.
func stopName(p *unison.Panel) string {
	switch v := p.Self.(type) {
	case *kvitui.Button:
		return "button " + quote(v.Words())
	case *kvitui.IconButton:
		return "symbol button " + quote(v.Name())
	case *kvitui.Link:
		return "link " + quote(v.Text)
	case *kvitui.Check:
		return "check box " + quote(v.Text)
	case *kvitui.Switch:
		return "switch " + quote(v.Text)
	case *kvitui.Tab:
		return "tab " + quote(v.Text)
	}
	for q := p; q != nil; q = q.Parent() {
		switch f := q.Self.(type) {
		case *kvitui.Field:
			return "field " + quote(f.Label)
		case *kvitui.SearchField:
			return "search field " + quote(f.Label)
		}
	}
	return strings.TrimPrefix(fmt.Sprintf("%T", p.Self), "*kvitui.")
}

// helpLink checks a link to a help topic: worded as the whole question the
// topic answers, ending in a question mark, and never a generic "Learn
// more"; one help button to a window or dialog on the macOS form (HLP-6).
func (c *labelCheck) helpLink(h *kvitui.HelpLink) {
	q := strings.TrimSpace(h.Question)
	where := "help link " + quote(q)
	rs := []rune(q)
	switch {
	case len(rs) == 0 || rs[len(rs)-1] != '?':
		c.add("HLP-6", where, "not worded as the question its topic answers, ending in a question mark")
	case len(rs) > 1 && unicode.IsPunct(rs[len(rs)-2]):
		c.add("HLP-6", where, "other punctuation before its question mark")
	}
	lower := strings.ToLower(q)
	for _, generic := range []string{"learn more", "get help", "troubleshoot", "help"} {
		if lower == generic || strings.HasPrefix(lower, generic+" ") {
			c.add("HLP-6", where, "generic words; ask the question the topic answers")
			break
		}
	}
	if c.form == kvitui.FormMac {
		c.helpButtons[c.scope]++
	}
}

// label checks static text: a paragraph under a control (WRD-22).
func (c *labelCheck) label(l *kvitui.Label, p *unison.Panel) {
	if n := len(strings.Fields(l.Text)); n > 20 {
		if parent := p.Parent(); parent != nil {
			kids := parent.Children()
			if i := slices.Index(kids, p); i > 0 && interactive(kids[i-1]) {
				c.add("WRD-22", "text under a control", fmt.Sprintf("%d words; a setting gets one short line, the rest in a Hint", n))
			}
		}
	}
}

// interactive reports whether a panel is a control a setting is made with.
func interactive(p *unison.Panel) bool {
	switch p.Self.(type) {
	case *kvitui.Switch, *kvitui.Check, *kvitui.Field, *kvitui.Select, *kvitui.RadioGroup, *kvitui.Slider,
		*kvitui.Stepper, *kvitui.Segmented:
		return true
	}
	return false
}

// menu checks a menu's lines (CMD-10, CMD-12, CMD-13, CMD-23, WRD-1 to
// WRD-6), depth 0 for a menu and 1 for a submenu.
func (c *labelCheck) menu(where string, lines []kvitui.MenuItem, context bool, depth int) {
	group, level := 0, 0
	groups := 1
	for i, it := range lines {
		if it.Separator {
			groups++
			group = 0
			continue
		}
		group++
		level++
		if group == 8 {
			c.add("CMD-13", where, "a group of more than seven lines")
		}
		words, _, _ := kvitui.AccessText(it.Text)
		line := where + ", line " + quote(words)
		if it.Danger {
			last := i == len(lines)-1 || lines[len(lines)-1].Danger
			if !last {
				c.add("CMD-23", line, "a destructive line that is not last")
			} else if i > 0 && !lines[i-1].Separator && !lines[i-1].Danger {
				c.add("CMD-23", line, "no separator before the destructive line")
			}
		}
		if len(it.Items) > 0 {
			if strings.HasSuffix(words, "…") {
				c.add("WRD-4", line, "a line that opens a submenu ends in an ellipsis")
			}
			if depth > 0 {
				c.add("CMD-12", line, "a submenu inside a submenu")
			}
			if context && c.form == kvitui.FormWindows {
				c.add("CMD-12", line, "a submenu in a context menu on the Windows form")
			}
			if n := countLines(it.Items); n > 7 {
				c.add("CMD-12", line, fmt.Sprintf("a submenu of %d lines", n))
			}
			c.caseOf(line, words, !it.Verbatim)
			// The lines of a submenu named with a verb drop it: Copy as >
			// Markdown (WRD-1).
			first := strings.ToLower(strings.Fields(words + " x")[0])
			c.menuVerbed = Verbs[first]
			c.menu(line, it.Items, context, depth+1)
			c.menuVerbed = false
			continue
		}
		if it.Verbatim {
			continue
		}
		std := it.Command != nil && it.Command.Std != kvitui.StdNone
		toggle := it.Radio || it.Mixed || (it.Command != nil && it.Command.Checked != nil) || it.TurnOn != ""
		c.command(line, words, !std && !toggle && !(depth > 0 && c.menuVerbed))
		if context && c.form == kvitui.FormWindows && it.Disabled {
			c.add("CMD-10", line, "a disabled line in a context menu on the Windows form; leave it out")
		}
	}
	if level > 25 {
		c.add("CMD-13", where, fmt.Sprintf("%d lines in one level", level))
	}
	if context && (level > 15 || groups > 4) {
		c.add("CMD-10", where, fmt.Sprintf("a context menu of %d lines in %d groups", level, groups))
	}
}

func countLines(lines []kvitui.MenuItem) int {
	n := 0
	for _, it := range lines {
		if !it.Separator {
			n++
		}
	}
	return n
}

func endsInPunctuation(s string) bool {
	rs := []rune(strings.TrimSpace(s))
	return len(rs) > 0 && unicode.IsPunct(rs[len(rs)-1]) && rs[len(rs)-1] != ')'
}

// dashPlaceholder reports whether a placeholder is only a dash standing for
// a value that is not set, as a record pane's fields and Figure do for a
// value nobody measured, rather than words ending in punctuation (WRD-19).
func dashPlaceholder(s string) bool {
	switch strings.Join(strings.Fields(s), "") {
	case "—", "–", "-":
		return true
	}
	return false
}

// Report fails the test for every problem not listed in known, and for every
// entry of known that no problem matches any more, so a list of known
// problems only shrinks as they are fixed. An entry matches a problem whose
// String contains it.
func Report(t testing.TB, problems []Problem, known ...string) {
	t.Helper()
	used := make([]bool, len(known))
	for _, p := range problems {
		s := p.String()
		matched := false
		for i, k := range known {
			if strings.Contains(s, k) {
				used[i], matched = true, true
			}
		}
		if !matched {
			t.Error(s)
		}
	}
	for i, k := range known {
		if !used[i] {
			t.Errorf("a known problem is gone; take %q off the list", k)
		}
	}
}

// CheckLabels reads the session's window in both forms, the Windows form
// (also Linux's) and the macOS form, and fails the test for every problem
// Labels finds that known does not list, as Report does. The window is left
// in the form it was in.
func (s *Session) CheckLabels(known ...string) {
	s.t.Helper()
	var problems []Problem
	for _, f := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
		var was kvitui.Form
		s.Do(func() {
			was = s.UI.Form()
			s.UI.SetForm(f)
		})
		s.Sync()
		s.Do(func() {
			problems = append(problems, Labels(s.Window)...)
			s.UI.SetForm(was)
		})
		s.Sync()
	}
	Report(s.t, problems, known...)
}
