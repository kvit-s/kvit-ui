package kvitui

import (
	"strconv"
	"strings"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// TreeModel is what a Tree shows, for a tree whose size is not known in
// advance, such as the folders of a notes vault: the tree asks only about
// the nodes it shows. A node is named by its path, the index of each node on
// the way down from the top; the empty path is the root above the top level.
type TreeModel interface {
	// Children is how many children the node at path has.
	Children(path []int) int
	// Label names the node at path.
	Label(path []int) string
}

// TreeNode is one node of a tree written out in full, for a tree whose shape
// is fixed: account groups, a category hierarchy, the sections of a settings
// page.
type TreeNode struct {
	Label    string
	Children []TreeNode
}

// TreeNodes is a TreeModel over nodes written out in full.
type TreeNodes []TreeNode

func (n TreeNodes) at(path []int) ([]TreeNode, *TreeNode) {
	level := []TreeNode(n)
	var node *TreeNode
	for _, i := range path {
		if i < 0 || i >= len(level) {
			return nil, nil
		}
		node = &level[i]
		level = node.Children
	}
	return level, node
}

// Children is how many children the node at path has.
func (n TreeNodes) Children(path []int) int {
	level, _ := n.at(path)
	return len(level)
}

// Label names the node at path.
func (n TreeNodes) Label(path []int) string {
	if _, node := n.at(path); node != nil {
		return node.Label
	}
	return ""
}

// Tree is a nested list the reader opens and closes: folders, an outline, a
// hierarchy of accounts. What a shared tree has to get right is the keyboard:
// Right opens a node and then steps into it, Left closes it or steps out to
// its parent, and Home and End go to the ends, so depth is never a wall a
// keyboard user has to arrow through a row at a time. Indentation is a loose
// space a level, so a deep tree stays in proportion at every interface size.
// It draws only the rows in view, and scrolls with the wheel.
type Tree struct {
	unison.Panel
	ui *UI
	// Label names the tree for a screen reader.
	Label string
	// Model is what the tree shows.
	Model TreeModel
	// OnChoose runs when the reader moves to a node.
	OnChoose func(path []int)
	// OnActivate runs when the reader presses Return on a node with no
	// children, or presses one twice.
	OnActivate func(path []int)

	expanded map[string]bool
	rows     []treeRow // the rows shown, top to bottom
	stale    bool      // the rows have to be worked out again
	current  string    // the key of the node the reader is on
	hovered  int
	scrollY  float32
	ease     wheelScroll
}

type treeRow struct {
	path     []int
	key      string
	label    string
	children int
}

// NewTree returns a tree of nodes written out in full, every node closed.
func NewTree(ui *UI, label string, nodes ...TreeNode) *Tree {
	t := &Tree{ui: ui, Label: label, Model: TreeNodes(nodes), expanded: map[string]bool{}, stale: true, hovered: -1}
	t.Self = t
	t.SetFocusable(true)
	t.ease = wheelScroll{ui: ui,
		at:     func() (x, y float32) { return 0, t.scrollY },
		travel: func() (x, y float32) { return 0, t.travel() },
		moveTo: func(_, y float32) { t.scrollTo(y) }}
	t.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		m := ui.Interface
		h := float32(len(t.visible())) * t.rowHeight()
		return geom.NewSize(float32(m.Px(80)), t.rowHeight()), geom.NewSize(float32(m.SidebarWidth()), h),
			geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
	})
	t.DrawCallback = t.draw
	t.KeyDownCallback = t.keyDown
	t.MouseDownCallback = t.mouseDown
	t.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { t.hover(t.rowAt(where)); return false }
	t.MouseEnterCallback = t.MouseMoveCallback
	t.MouseExitCallback = func() bool { t.hover(-1); return false }
	t.MouseWheelCallback = func(_, delta geom.Point, _ mod.Modifiers) bool { return t.ease.wheel(delta) }
	t.GainedFocusCallback = func() {
		if t.current == "" && len(t.visible()) > 0 {
			t.current = t.visible()[0].key
		}
		t.MarkForRedraw()
	}
	t.LostFocusCallback = t.MarkForRedraw
	t.FrameChangeCallback = func() {
		if w := t.Window(); w != nil {
			ui.watchWindow(w)
		}
	}
	return t
}

func pathKey(path []int) string {
	parts := make([]string, len(path))
	for i, p := range path {
		parts[i] = strconv.Itoa(p)
	}
	return strings.Join(parts, "/")
}

func (t *Tree) rowHeight() float32 { return float32(t.ui.Interface.RowHeightSlim()) }

// Refresh works out the rows again after the model changed.
func (t *Tree) Refresh() {
	t.stale = true
	t.MarkForLayoutAndRedraw()
}

// visible are the rows shown: each node whose ancestors are all open.
func (t *Tree) visible() []treeRow {
	if !t.stale || t.Model == nil {
		return t.rows
	}
	t.rows = t.rows[:0]
	var walk func(path []int)
	walk = func(path []int) {
		for i := range t.Model.Children(path) {
			p := append(append([]int(nil), path...), i)
			key := pathKey(p)
			n := t.Model.Children(p)
			t.rows = append(t.rows, treeRow{path: p, key: key, label: t.Model.Label(p), children: n})
			if n > 0 && t.expanded[key] {
				walk(p)
			}
		}
	}
	walk(nil)
	t.stale = false
	return t.rows
}

// SetExpanded opens or closes the node at path.
func (t *Tree) SetExpanded(path []int, open bool) {
	key := pathKey(path)
	if t.expanded[key] == open {
		return
	}
	if open {
		t.expanded[key] = true
	} else {
		delete(t.expanded, key)
		// A reader inside a node that closes is moved to the node.
		if strings.HasPrefix(t.current, key+"/") {
			t.current = key
		}
	}
	t.Refresh()
}

// IsExpanded reports whether the node at path is open.
func (t *Tree) IsExpanded(path []int) bool { return t.expanded[pathKey(path)] }

// ExpandTo opens every node shallower than depth: 1 opens the top level.
func (t *Tree) ExpandTo(depth int) {
	var walk func(path []int)
	walk = func(path []int) {
		if len(path) >= depth {
			return
		}
		for i := range t.Model.Children(path) {
			p := append(append([]int(nil), path...), i)
			if t.Model.Children(p) > 0 {
				t.expanded[pathKey(p)] = true
				walk(p)
			}
		}
	}
	walk(nil)
	t.Refresh()
}

// Current is the path of the node the reader is on, or nil.
func (t *Tree) Current() []int {
	for _, r := range t.visible() {
		if r.key == t.current {
			return r.path
		}
	}
	return nil
}

func (t *Tree) index(key string) int {
	for i, r := range t.visible() {
		if r.key == key {
			return i
		}
	}
	return -1
}

// choose moves the reader to a row, keeping it in view.
func (t *Tree) choose(i int) {
	rows := t.visible()
	if i < 0 || i >= len(rows) {
		return
	}
	if rows[i].key != t.current {
		t.current = rows[i].key
		if t.OnChoose != nil {
			t.OnChoose(rows[i].path)
		}
	}
	top := float32(i) * t.rowHeight()
	view := t.ContentRect(false).Height
	switch {
	case top < t.scrollY:
		t.scrollTo(top)
	case top+t.rowHeight() > t.scrollY+view:
		t.scrollTo(top + t.rowHeight() - view)
	}
	t.MarkForRedraw()
}

func (t *Tree) travel() float32 {
	return max(0, float32(len(t.visible()))*t.rowHeight()-t.ContentRect(false).Height)
}

func (t *Tree) scrollTo(y float32) {
	y = max(0, min(t.travel(), y))
	if y != t.scrollY {
		t.scrollY = y
		t.MarkForRedraw()
	}
}

func (t *Tree) rowAt(where geom.Point) int {
	r := t.ContentRect(false)
	if !where.In(r) {
		return -1
	}
	i := int((where.Y - r.Y + t.scrollY) / t.rowHeight())
	if i < 0 || i >= len(t.visible()) {
		return -1
	}
	return i
}

func (t *Tree) hover(i int) {
	if i != t.hovered {
		t.hovered = i
		t.MarkForRedraw()
	}
}

// chevron is where a row's open-and-close mark is.
func (t *Tree) chevron(i int, depth int) geom.Rect {
	m := t.ui.Interface
	r := t.ContentRect(false)
	s := float32(m.IconSizeSmall())
	y := r.Y + float32(i)*t.rowHeight() - t.scrollY
	return geom.NewRect(r.X+float32(m.SpaceSnug()+depth*m.SpaceLoose()), y+(t.rowHeight()-s)/2, s, s)
}

func (t *Tree) mouseDown(where geom.Point, button, clicks int, _ mod.Modifiers) bool {
	t.RequestFocus()
	i := t.rowAt(where)
	if button != unison.ButtonLeft || i < 0 {
		return true
	}
	row := t.visible()[i]
	if row.children > 0 && where.In(t.chevron(i, len(row.path)-1).Inset(geom.NewUniformInsets(-float32(t.ui.Interface.SpaceTight())))) {
		t.SetExpanded(row.path, !t.expanded[row.key])
		return true
	}
	t.choose(i)
	if clicks == 2 {
		if row.children > 0 {
			t.SetExpanded(row.path, !t.expanded[row.key])
		} else if t.OnActivate != nil {
			t.OnActivate(row.path)
		}
	}
	return true
}

func (t *Tree) keyDown(key unison.KeyCode, _ mod.Modifiers, _ bool) bool {
	rows := t.visible()
	if len(rows) == 0 {
		return false
	}
	i := t.index(t.current)
	if i < 0 {
		t.choose(0)
		return true
	}
	row := rows[i]
	switch key {
	case unison.KeyUp:
		t.choose(max(0, i-1))
	case unison.KeyDown:
		t.choose(min(len(rows)-1, i+1))
	case unison.KeyHome:
		t.choose(0)
	case unison.KeyEnd:
		t.choose(len(rows) - 1)
	case unison.KeyRight:
		switch {
		case row.children > 0 && !t.expanded[row.key]:
			t.SetExpanded(row.path, true)
		case row.children > 0:
			t.choose(i + 1)
		}
	case unison.KeyLeft:
		switch {
		case row.children > 0 && t.expanded[row.key]:
			t.SetExpanded(row.path, false)
		case len(row.path) > 1:
			t.choose(t.index(pathKey(row.path[:len(row.path)-1])))
		}
	case unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace:
		if row.children > 0 {
			t.SetExpanded(row.path, !t.expanded[row.key])
		} else if t.OnActivate != nil {
			t.OnActivate(row.path)
		}
	default:
		return false
	}
	t.MarkForRedraw()
	return true
}

func (t *Tree) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	p := painterFor(gc, ui)
	r := t.ContentRect(false)
	rows := t.visible()
	rowH := t.rowHeight()
	for i := int(t.scrollY / rowH); i < len(rows); i++ {
		y := r.Y + float32(i)*rowH - t.scrollY
		if y >= r.Bottom() {
			break
		}
		row := rows[i]
		box := geom.NewRect(r.X, y, r.Width, rowH)
		current := row.key == t.current
		switch {
		case current:
			p.fill(box, tk.SelectionTint)
		case i == t.hovered:
			p.fill(box, tk.HoverTint)
		}
		depth := len(row.path) - 1
		if row.children > 0 {
			name := "chevron-right"
			if t.expanded[row.key] {
				name = "chevron-down"
			}
			if g, ok := icons.Glyph(name); ok {
				c := t.chevron(i, depth)
				drawGlyph(gc, ui, g, c.Width, tk.TextMuted, c)
			}
		}
		ink := tk.TextSecondary
		if current {
			ink = tk.TextPrimary
		}
		x := r.X + float32(m.Px(24)+depth*m.SpaceLoose())
		l := ui.Fonts.Layout([]text.Span{{Text: row.label, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}},
			text.Options{MaxWidth: max(1, r.Right()-float32(m.SpaceNear())-x), Elide: true})
		_, h := l.Size()
		l.Draw(gc, x, y+(rowH-h)/2)
	}
}

// ProvideAccessibility describes the tree and the rows in view, each with its
// level and whether it is open, the focus on the reader's node.
func (t *Tree) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Tree
	}
	n.Name = t.Label
	rows := t.visible()
	n.RowCount = len(rows)
	r := t.ContentRect(false)
	rowH := t.rowHeight()
	for i, row := range rows {
		y := r.Y + float32(i)*rowH - t.scrollY
		if (y+rowH < r.Y || y > r.Bottom()) && row.key != t.current {
			continue
		}
		id := b.AddVirtualChild(row.key, func(v *accessibility.Node) {
			v.Role = role.Row
			v.Name = row.label
			v.Level = len(row.path)
			v.RowIndex = i
			v.Bounds = geom.NewRect(r.X, y, r.Width, rowH)
			v.Selectable, v.Selected = true, row.key == t.current
			v.Focusable = true
			v.Expandable = row.children > 0
			v.Expanded = t.expanded[row.key]
			v.Actions = v.Actions.With(accessibility.Select, accessibility.Focus, accessibility.ScrollIntoView)
			if row.children > 0 {
				v.Actions = v.Actions.With(accessibility.Expand, accessibility.Collapse)
			}
		})
		if row.key == t.current {
			b.FocusChild(id)
		}
	}
}

// PerformAccessibilityAction moves to, opens or closes a node for a screen
// reader.
func (t *Tree) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	key, ok := req.Key.(string)
	if !ok {
		return false
	}
	i := t.index(key)
	if i < 0 {
		return false
	}
	row := t.visible()[i]
	switch req.Action {
	case accessibility.Select, accessibility.Focus, accessibility.ScrollIntoView:
		t.choose(i)
	case accessibility.Expand:
		t.SetExpanded(row.path, true)
	case accessibility.Collapse:
		t.SetExpanded(row.path, false)
	default:
		return false
	}
	return true
}
