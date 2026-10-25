package uitest

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

// Screen is what the whole-app checks read from one screen of an app, taken
// with Collect while the screen is open: the commands it offers, the words
// and symbols on its controls, and the commands it runs that are on no
// command surface. Screens taken in different sessions are compared by
// their words, keys and symbols, since their commands are different objects.
type Screen struct {
	// Name says which screen it is, in what the checks report.
	Name string
	// Form is the form the screen was read in; screens are compared only
	// with others read in the same form.
	Form kvitui.Form
	// Commands are the screen's commands.
	Commands []CommandFact
	// Names are the words of its buttons, symbol buttons and menu lines.
	Names []NameUse
	// Symbols are the symbols its controls draw and what each is called.
	Symbols []NameUse
	// OffSurface are the names of the commands it runs that are on no menu
	// bar or toolbar (CMD-1).
	OffSurface []string
	// Menus are its menu bar's menus, by title, and how many lines each has;
	// the app's own menus only.
	Menus map[string]int
	// MenuCount is how many menus its menu bar has.
	MenuCount int
	// Help are the help topics its help links open.
	Help []HelpUse
}

// HelpUse is a help link's topic, and whether the app ships it.
type HelpUse struct {
	Topic, Question string
	// Missing says why the topic cannot be opened, "" when its page is in
	// UI.HelpRoot with its anchor.
	Missing string
}

// CommandFact is one command as the whole-app checks read it.
type CommandFact struct {
	Name   string
	Std    kvitui.StdCommand
	Symbol string
	// Keys are its shortcuts as each form writes them.
	Keys map[kvitui.Form][]string
}

// NameUse is a word or a symbol and where it is used.
type NameUse struct {
	Key, Name, Where string
}

// Collect reads a screen for App. Call it on the interface's thread.
func Collect(name string, w *kvitui.Window) Screen {
	ui := w.UI()
	s := Screen{Name: name, Form: ui.Form(), Menus: map[string]int{}}
	var cmds []*kvitui.Command
	addCmd := func(c *kvitui.Command) {
		if c != nil && !slices.Contains(cmds, c) {
			cmds = append(cmds, c)
		}
	}
	for _, c := range w.Commands() {
		addCmd(c)
	}
	if b := w.MenuBar(); b != nil {
		for _, c := range b.Commands() {
			addCmd(c)
		}
		menus := b.Menus()
		s.MenuCount = len(menus)
		for _, m := range menus {
			lines := ui.MenuLines(m.Items, false)
			title, _, _ := kvitui.AccessText(m.Title)
			if isAppMenu(m.Title) {
				s.Menus[title] = countLines(lines)
			}
			var walk func(items []kvitui.MenuItem)
			walk = func(items []kvitui.MenuItem) {
				for _, it := range items {
					if it.Separator || it.Verbatim {
						continue
					}
					words, _, _ := kvitui.AccessText(it.Text)
					s.Names = append(s.Names, NameUse{Name: words, Where: "menu " + quote(title)})
					if it.Symbol != "" {
						s.Symbols = append(s.Symbols, NameUse{Key: it.Symbol, Name: words, Where: "menu " + quote(title)})
					}
					walk(it.Items)
				}
			}
			walk(lines)
		}
	}
	var walk func(p *unison.Panel)
	walk = func(p *unison.Panel) {
		if p.Hidden {
			return
		}
		switch v := p.Self.(type) {
		case *kvitui.Toolbar:
			for _, c := range v.Commands() {
				addCmd(c)
			}
		case *kvitui.Button:
			addCmd(v.Command)
			if !v.Busy {
				s.Names = append(s.Names, NameUse{Name: v.Words(), Where: "button"})
				if sym := v.Symbol; sym != "" {
					s.Symbols = append(s.Symbols, NameUse{Key: sym, Name: v.Words(), Where: "button"})
				}
			}
		case *kvitui.HelpLink:
			s.Help = append(s.Help, HelpUse{Topic: v.Topic, Question: v.Question, Missing: missingTopic(ui.HelpRoot, v.Topic)})
		case *kvitui.IconButton:
			addCmd(v.Command)
			sym := v.Symbol
			if sym == "" && v.Command != nil {
				sym = v.Command.SymbolName()
			}
			s.Names = append(s.Names, NameUse{Name: v.Name(), Where: "symbol button"})
			s.Symbols = append(s.Symbols, NameUse{Key: sym, Name: v.Name(), Where: "symbol button"})
		case *kvitui.Select:
			// A compact select's symbol names what the choice is about, its
			// Label, rather than the value in force: the calendar on a
			// period filter means period whatever the period is.
			if v.Symbol != "" {
				s.Symbols = append(s.Symbols, NameUse{Key: v.Symbol, Name: v.Label, Where: "select"})
			}
		}
		for _, c := range p.Children() {
			walk(c)
		}
	}
	walk(w.Content())
	for _, c := range cmds {
		f := CommandFact{Name: c.Label(ui), Std: c.Std, Symbol: c.SymbolName(), Keys: map[kvitui.Form][]string{}}
		for _, form := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
			for _, k := range c.Shortcuts(form) {
				f.Keys[form] = append(f.Keys[form], kvitui.ShortcutTextIn(form, k))
			}
		}
		s.Commands = append(s.Commands, f)
		if c.Symbol != "" {
			s.Symbols = append(s.Symbols, NameUse{Key: c.Symbol, Name: c.Label(ui), Where: "command"})
		}
	}
	for _, c := range w.OffSurface() {
		s.OffSurface = append(s.OffSurface, c.Label(ui))
	}
	return s
}

// isAppMenu reports whether a menu bar title is one of the app's own menus
// rather than a standard one.
func isAppMenu(title string) bool {
	switch title {
	case "&File", "&Edit", "F&ormat", "&View", "&Tools", "&Window", "&Help":
		return false
	}
	return true
}

// sameName is a label reduced to what makes two labels the same command:
// lower case, without the access-key marker, the ellipsis and spaces at the
// ends.
func sameName(s string) string {
	s = strings.ReplaceAll(s, "&", "")
	s = strings.TrimSuffix(strings.TrimSpace(s), "…")
	return strings.ToLower(strings.TrimSpace(s))
}

// symbolVerbs are the verbs that name one meaning of a symbol: adding for
// plus, taking away for minus, taking away for close. A symbol's uses are
// compared by their verb, the first word of the name, so two Add … buttons
// carrying the plus symbol are one meaning; two names whose verbs differ
// outside such a table are still two meanings, as next period and show fields
// on one chevron are.
var symbolVerbs = map[string]map[string]string{
	"close": {"close": "close", "dismiss": "close", "clear": "close"},
	"plus":  {"add": "add", "new": "add", "create": "add", "increase": "add"},
	"minus": {"remove": "remove", "decrease": "remove"},
}

// symbolMeaning reduces a use of a symbol to the meaning CON-2 compares: the
// verb, the first word of the name, with the verbs in the symbol's table
// standing for the one meaning.
func symbolMeaning(symbol, name string) string {
	verb := strings.Fields(sameName(name) + " x")[0]
	if table, ok := symbolVerbs[symbol]; ok {
		if canon, ok := table[verb]; ok {
			return canon
		}
	}
	return verb
}

// App lists where an app's screens, read together, break the rules that
// need every screen at once (rules.md, the "test, all screens of an app"
// entries): one command written two ways or with two symbols or keys
// (CON-1), one symbol with two meanings (CON-2), one key for two commands
// (KEY-5), a standard command made by hand without its standard key (KEY-2),
// a command on no command surface (CMD-1), and a menu bar of more than ten
// menus or an app menu of one or two lines (CMD-2).
func App(screens ...Screen) []Problem {
	var out []Problem
	add := func(rule, where, what string) {
		out = append(out, Problem{Rule: rule, Where: where, What: what})
	}
	// CON-1: the same command written two ways, anywhere it is shown in one
	// form.
	written := map[string]map[string]string{}
	for _, s := range screens {
		for _, n := range s.Names {
			k := s.Form.String() + " form: " + sameName(n.Name)
			if k == "" {
				continue
			}
			if written[k] == nil {
				written[k] = map[string]string{}
			}
			plain, _, _ := kvitui.AccessText(n.Name)
			written[k][strings.TrimSuffix(plain, "…")] = s.Name + ", " + n.Where
		}
	}
	for _, k := range sortedKeys(written) {
		if forms := written[k]; len(forms) > 1 {
			var ways []string
			for w, where := range forms {
				ways = append(ways, quote(w)+" ("+where+")")
			}
			sort.Strings(ways)
			add("CON-1", quote(k), "one command written "+strings.Join(ways, ", "))
		}
	}
	// CON-1, KEY-5 and KEY-2 on the commands themselves.
	byName := map[string][]CommandFact{}
	byKey := map[string][]CommandFact{}
	for _, s := range screens {
		for _, c := range s.Commands {
			name := s.Form.String() + " form: " + sameName(c.Name)
			byName[name] = append(byName[name], c)
			for _, k := range c.Keys[s.Form] {
				byKey[s.Form.String()+" "+k] = append(byKey[s.Form.String()+" "+k], c)
			}
			if c.Std == kvitui.StdNone {
				if stdNamed(sameName(c.Name)) {
					add("KEY-2", "command "+quote(c.Name), "a standard command made by hand; give it its Std so it has the standard name and keys")
				}
			}
		}
	}
	for _, name := range sortedKeys(byName) {
		cs := byName[name]
		for _, c := range cs[1:] {
			if c.Symbol != cs[0].Symbol {
				add("CON-1", "command "+quote(cs[0].Name), "two symbols: "+quote(cs[0].Symbol)+" and "+quote(c.Symbol))
				break
			}
		}
		for _, c := range cs[1:] {
			if fmt.Sprint(c.Keys) != fmt.Sprint(cs[0].Keys) {
				add("CON-1", "command "+quote(cs[0].Name), "two sets of keys")
				break
			}
		}
	}
	for _, key := range sortedKeys(byKey) {
		names := map[string]bool{}
		for _, c := range byKey[key] {
			names[sameName(c.Name)] = true
		}
		if len(names) > 1 {
			add("KEY-5", key, "one key for "+strings.Join(sortedKeys(names), " and "))
		}
	}
	// CON-2: one symbol, one meaning.
	meanings := map[string]map[string]bool{}
	for _, s := range screens {
		for _, u := range s.Symbols {
			if u.Key == "" {
				continue
			}
			n := symbolMeaning(u.Key, u.Name)
			key := s.Form.String() + " form: symbol " + quote(u.Key)
			if meanings[key] == nil {
				meanings[key] = map[string]bool{}
			}
			meanings[key][n] = true
		}
	}
	for _, key := range sortedKeys(meanings) {
		if names := meanings[key]; len(names) > 1 {
			add("CON-2", key, "used for "+strings.Join(sortedKeys(names), ", "))
		}
	}
	for _, s := range screens {
		for _, name := range s.OffSurface {
			add("CMD-1", s.Name+": command "+quote(name), "on no menu bar or toolbar; a command reached only by a key, a button or a context menu cannot be found")
		}
		for _, h := range s.Help {
			if h.Missing != "" {
				add("HLP-6", s.Name+": help link "+quote(h.Question), "its topic "+quote(h.Topic)+" "+h.Missing)
			}
		}
		if s.MenuCount > 10 {
			add("CMD-2", s.Name+": menu bar", fmt.Sprintf("%d menus; keep to ten", s.MenuCount))
		}
		for _, title := range sortedKeys(s.Menus) {
			if n := s.Menus[title]; n <= 2 {
				add("CMD-2", s.Name+": menu "+quote(title), fmt.Sprintf("%d lines; a menu of one or two belongs in another", n))
			}
		}
	}
	return out
}

// missingTopic says why a help topic, "budgets#rollover", cannot be opened
// from the help folder: its page is not there, or has no element with the
// anchor as its id or name. "" means it can.
func missingTopic(root, topic string) string {
	if root == "" {
		return "cannot open: UI.HelpRoot is not set"
	}
	page, anchor, _ := strings.Cut(topic, "#")
	b, err := os.ReadFile(filepath.Join(root, page+".html"))
	if err != nil {
		return "has no page: " + page + ".html is not in " + root
	}
	if anchor == "" {
		return ""
	}
	html := string(b)
	for _, attr := range []string{"id", "name"} {
		for _, q := range []string{`"`, `'`} {
			if strings.Contains(html, attr+"="+q+anchor+q) {
				return ""
			}
		}
	}
	return "has no anchor: " + page + ".html has nothing named " + quote(anchor)
}

// stdNames are the names of the standard commands, which an app makes with
// Std rather than by hand.
var stdNames = wordSet("new open save print undo redo cut copy paste find replace rename refresh settings help exit quit")

// stdNamed reports whether a name is a standard command's.
func stdNamed(name string) bool {
	switch name {
	case "save as", "select all", "find next", "full screen":
		return true
	}
	return stdNames[name]
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// CheckApp fails the test for every problem App finds in the screens that
// known does not list, as Report does.
func CheckApp(t testing.TB, screens []Screen, known ...string) {
	t.Helper()
	Report(t, App(screens...), known...)
}

// Collect reads the session's window for App, as a screen named name.
func (s *Session) Collect(name string) Screen {
	var out Screen
	s.Do(func() { out = Collect(name, s.Window) })
	return out
}
