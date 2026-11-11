package kvitui_test

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	"math"
	"slices"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// leadColour is what the pictures in these tests draw, a colour no theme uses.
var leadColour = unison.RGB(250, 20, 200)

// leadPanel is a picture held to a size in design pixels the way an application
// holds one, with Sized: outer is what is handed to SetLead and inner, inside
// it, is what the pointer lands on.
func leadPanel(ui *kvitui.UI, px int) (outer, inner *unison.Panel) {
	inner = unison.NewPanel()
	inner.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		r := inner.ContentRect(false)
		gc.DrawRect(r, leadColour.Paint(gc, r, paintstyle.Fill))
	}
	return kvitui.Sized(ui, kvitui.Px(px), kvitui.Px(px), inner), inner
}

func rootRect(p unison.Paneler) geom.Rect {
	return p.AsPanel().RectToRoot(p.AsPanel().ContentRect(true))
}

func samePixel(a, b color.Color) bool {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	return ar == br && ag == bg && ab == bb
}

func isLeadColour(c color.Color) bool {
	return samePixel(c, color.NRGBA{R: 250, G: 20, B: 200, A: 255})
}

// sameStretch compares two boxes of pixels of one size, and reports how many
// pixels differ and how many of a's are not its corner's colour, so a test can
// tell two matching stretches of words from two matching stretches of nothing.
func sameStretch(img image.Image, a, b image.Point, size image.Point) (differ, drawn int) {
	ground := img.At(a.X, a.Y)
	for y := range size.Y {
		for x := range size.X {
			pa := img.At(a.X+x, a.Y+y)
			if !samePixel(pa, img.At(b.X+x, b.Y+y)) {
				differ++
			}
			if !samePixel(pa, ground) {
				drawn++
			}
		}
	}
	return differ, drawn
}

func wholePixels(r geom.Rect) bool {
	return r.X == float32(math.Round(float64(r.X))) && r.Y == float32(math.Round(float64(r.Y)))
}

var leadSizes = []int{tokens.MinInterfaceSize, tokens.DefaultInterfaceSize, tokens.MaxInterfaceSize}

// A lead stands centred on where the symbol stands, a larger one no less, and
// the label stays where it is, in the sidebar at the smallest and largest
// interface sizes. In the rail it is centred as the symbol is.
func TestASidebarItemsLeadStandsWhereItsSymbolDoes(t *testing.T) {
	for _, size := range leadSizes {
		t.Run(fmt.Sprintf("%dpx", size), func(t *testing.T) {
			var picture, symbol, tiny, tinyAlone, huge, railed *kvitui.SidebarItem
			var lead, hugeLead, railLead *unison.Panel
			s := uitest.Open(t, uitest.Options{InterfaceSize: size}, func(ui *kvitui.UI) unison.Paneler {
				picture = kvitui.NewSidebarItem(ui, "The same words", "folder")
				outer, inner := leadPanel(ui, 16)
				picture.SetLead(outer)
				lead = inner
				symbol = kvitui.NewSidebarItem(ui, "The same words", "folder")
				// A small lead leaves the symbol's place in view, to show that the
				// symbol is not drawn behind it.
				tiny = kvitui.NewSidebarItem(ui, "Tiny", "folder")
				outer, _ = leadPanel(ui, 4)
				tiny.SetLead(outer)
				tinyAlone = kvitui.NewSidebarItem(ui, "Tiny", "")
				outer, _ = leadPanel(ui, 4)
				tinyAlone.SetLead(outer)
				huge = kvitui.NewSidebarItem(ui, "The same words", "folder")
				outer, inner = leadPanel(ui, 200)
				huge.SetLead(outer)
				hugeLead = inner
				railed = kvitui.NewSidebarItem(ui, "In the rail", "folder")
				outer, inner = leadPanel(ui, 16)
				railed.SetLead(outer)
				railLead = inner
				rail := kvitui.NewSidebar(ui, railed)
				rail.Collapsed = true
				return kvitui.Row(ui, kvitui.SizeColumnGap, kvitui.NewSidebar(ui, picture, symbol, tiny, tinyAlone, huge), rail)
			})
			img := s.Capture()
			s.Do(func() {
				m := s.UI.Interface
				iconSize := float32(m.IconSizeSmall())
				item := rootRect(picture)
				box := rootRect(lead)
				if box.Width != float32(m.Px(16)) || box.Height != float32(m.Px(16)) {
					t.Errorf("%d px: the lead is %v, want %d square", size, box.Size, m.Px(16))
				}
				if !wholePixels(box) {
					t.Errorf("%d px: the lead is at %v, between pixels", size, box.Point)
				}
				if d := box.CenterX() - (item.X + float32(m.SpaceLoose()) + iconSize/2); d < -0.5 || d > 0.5 {
					t.Errorf("%d px: the lead's middle is %v off the symbol's", size, d)
				}
				if d := box.CenterY() - item.CenterY(); d < -0.5 || d > 0.5 {
					t.Errorf("%d px: the lead's middle is %v off the row's", size, d)
				}
				if !isLeadColour(img.At(int(box.CenterX()), int(box.CenterY()))) {
					t.Errorf("%d px: the lead is not drawn", size)
				}
				// From where the label starts to the end of the row, the two items
				// match pixel for pixel, and there are words there to match.
				left := int(item.X) + m.SpaceLoose() + m.IconSizeSmall() + m.Space()
				stretch := image.Pt(int(item.Right())-left, int(item.Height))
				at := func(p *kvitui.SidebarItem) image.Point { return image.Pt(left, int(rootRect(p).Y)) }
				if differ, drawn := sameStretch(img, at(picture), at(symbol), stretch); differ != 0 || drawn == 0 {
					t.Errorf("%d px: the labels differ in %d pixels, with %d drawn", size, differ, drawn)
				}
				if differ, drawn := sameStretch(img, at(huge), at(symbol), stretch); differ != 0 || drawn == 0 {
					t.Errorf("%d px: a huge lead moved the label: %d pixels differ, with %d drawn", size, differ, drawn)
				}
				// A lead takes the symbol's place: with one, the symbol set or not
				// makes no difference to the whole row.
				if differ, drawn := sameStretch(img, image.Pt(int(item.X), int(rootRect(tiny).Y)), image.Pt(int(item.X), int(rootRect(tinyAlone).Y)),
					image.Pt(int(item.Width), int(item.Height))); differ != 0 || drawn == 0 {
					t.Errorf("%d px: the symbol is drawn behind the lead: %d pixels differ", size, differ)
				}
				// A lead too large for the place is drawn smaller, square still,
				// short of the label and inside the row.
				big := rootRect(hugeLead)
				if big.Right() >= float32(left) || big.Height >= item.Height || big.Width != big.Height {
					t.Errorf("%d px: a huge lead is %v, with the label at %d in a row %v tall", size, big, left, item.Height)
				}
				// In the rail, in the middle.
				rr := rootRect(railed)
				rb := rootRect(railLead)
				if d := rb.CenterX() - rr.CenterX(); d < -0.5 || d > 0.5 {
					t.Errorf("%d px: the lead's middle is %v off the rail's", size, d)
				}
				if d := rb.CenterY() - rr.CenterY(); d < -0.5 || d > 0.5 {
					t.Errorf("%d px: the lead's middle is %v off the row's in the rail", size, d)
				}
			})
		})
	}
}

// A lead given, replaced and taken away leaves the rest of the item alone, and
// the badge draws over it.
func TestASidebarItemsLeadComesAndGoes(t *testing.T) {
	var item *kvitui.SidebarItem
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		item = kvitui.NewSidebarItem(ui, "Inbox", "note")
		item.Count = 3
		return kvitui.NewSidebar(ui, item)
	})
	s.Do(func() {
		badge := item.Children()[0]
		first, _ := leadPanel(s.UI, 16)
		item.SetLead(first)
		if item.Lead() != unison.Paneler(first) || item.Children()[0] != first || item.Children()[1] != badge {
			t.Errorf("with a lead the children are %v", item.Children())
		}
		second, _ := leadPanel(s.UI, 16)
		item.SetLead(second)
		if first.Parent() != nil || len(item.Children()) != 2 || item.Children()[0] != second {
			t.Errorf("replacing the lead left %v", item.Children())
		}
		item.SetLead(nil)
		if item.Lead() != nil || len(item.Children()) != 1 || item.Children()[0] != badge {
			t.Errorf("taking the lead away left %v", item.Children())
		}
	})
}

// A press on a lead presses what it stands in, the pointer resting on it rests
// on what it stands in, and the panel inside a lead held with Sized behaves the
// same as the lead.
func TestALeadIsPartOfWhatItStandsIn(t *testing.T) {
	var item, railed *kvitui.SidebarItem
	var heading *kvitui.SectionHeading
	var link *kvitui.Link
	var itemLead, railLead, headingLead, linkLead *unison.Panel
	var pressed, railPressed, followed, toggled int
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var outer *unison.Panel
		item = kvitui.NewSidebarItem(ui, "An item", "folder")
		item.OnClick = func() { pressed++ }
		outer, itemLead = leadPanel(ui, 16)
		item.SetLead(outer)
		railed = kvitui.NewSidebarItem(ui, "In the rail", "folder")
		railed.OnClick = func() { railPressed++ }
		outer, railLead = leadPanel(ui, 16)
		railed.SetLead(outer)
		rail := kvitui.NewSidebar(ui, railed)
		rail.Collapsed = true
		heading = kvitui.NewSectionHeading(ui, "A heading")
		heading.Collapsible = true
		heading.OnToggle = func(bool) { toggled++ }
		outer, headingLead = leadPanel(ui, 16)
		heading.SetLead(outer)
		link = kvitui.NewLink(ui, "A link")
		link.OnActivate = func() { followed++ }
		outer, linkLead = leadPanel(ui, 16)
		link.SetLead(outer)
		return kvitui.Column(ui, kvitui.SizeSpace,
			kvitui.Row(ui, kvitui.SizeColumnGap, kvitui.NewSidebar(ui, item), rail),
			kvitui.FullWidth(heading), kvitui.Left(link))
	})
	for _, c := range []struct {
		name  string
		lead  *unison.Panel
		hover func() bool
		count func() int
	}{
		{"a sidebar item", itemLead, item.Hovered, func() int { return pressed }},
		{"a sidebar item in the rail", railLead, railed.Hovered, func() int { return railPressed }},
		{"a section heading", headingLead, heading.Hovered, func() int { return toggled }},
		{"a link", linkLead, link.Hovered, func() int { return followed }},
	} {
		s.Screen.MouseMove(s.Screen.PanelCenter(c.lead), 0)
		s.Sync()
		if !c.hover() {
			t.Errorf("the pointer on the lead of %s does not rest on it", c.name)
		}
		s.Screen.MouseMove(geom.NewPoint(1, 1), 0)
		s.Sync()
		if c.hover() {
			t.Errorf("the pointer off the lead of %s still rests on it", c.name)
		}
		s.Screen.Click(s.Screen.PanelCenter(c.lead))
		s.Sync()
		if c.count() != 1 {
			t.Errorf("a press on the lead of %s pressed it %d times", c.name, c.count())
		}
	}
	// In the rail the tooltip is the only thing that says where an item goes,
	// and it shows for the pointer resting on the lead.
	s.Screen.MouseMove(geom.NewPoint(1, 1), 0)
	s.Sync()
	s.Screen.MouseMove(s.Screen.PanelCenter(railLead), 0)
	s.WaitFor("the rail item's tooltip", func() bool {
		for _, n := range s.Nodes(role.Tooltip) {
			if n.Name == "In the rail" {
				return true
			}
		}
		return false
	})
}

// What a screen reader hears is the name, whatever the lead calls itself, in
// each form a heading takes, and in every theme CheckNamed accepts them.
func TestAScreenReaderHearsTheNameAndNothingOfALead(t *testing.T) {
	for _, theme := range tokens.BuiltInThemes() {
		t.Run(theme, func(t *testing.T) {
			var item *kvitui.SidebarItem
			var collapsible, withAction, fixed *kvitui.SectionHeading
			s := uitest.Open(t, uitest.Options{Theme: theme}, func(ui *kvitui.UI) unison.Paneler {
				named := func() unison.Paneler {
					outer, inner := leadPanel(ui, 16)
					outer.Accessibility.Role, outer.Accessibility.Name = role.Image, "Not heard"
					inner.Accessibility.Role, inner.Accessibility.Name = role.Image, "Not heard"
					return outer
				}
				item = kvitui.NewSidebarItem(ui, "Inbox", "note")
				item.Count, item.Counted = 12, "message"
				item.SetLead(named())
				collapsible = kvitui.NewSectionHeading(ui, "Opens")
				collapsible.Collapsible = true
				collapsible.SetLead(named())
				withAction = kvitui.NewSectionHeading(ui, "Holds an action")
				withAction.Action = "Run it"
				withAction.SetLead(named())
				fixed = kvitui.NewSectionHeading(ui, "Fixed")
				fixed.SetLead(named())
				return kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewSidebar(ui, item),
					kvitui.FullWidth(collapsible), kvitui.FullWidth(withAction), kvitui.FullWidth(fixed))
			})
			for _, c := range []struct {
				p    unison.Paneler
				role role.Enum
				name string
			}{
				{item, role.ListItem, "Inbox, 12 messages"},
				{collapsible, role.Button, "Opens"},
				{withAction, role.Group, "Holds an action"},
				{fixed, role.Heading, "Fixed"},
			} {
				if n := s.Node(c.p); n == nil || n.Role != c.role || n.Name != c.name {
					t.Errorf("%s: the node for %q is %+v", theme, c.name, n)
				}
			}
			for _, n := range s.Tree().Nodes {
				if n.Name == "Not heard" {
					t.Errorf("%s: a screen reader is told the lead's own name: %+v", theme, n)
				}
			}
			s.CheckNamed()
		})
	}
}

// A heading's lead stands after the chevron and before the name, at a stack gap
// from each, at the smallest and largest interface sizes, and the name is the
// same pixels moved along by the lead.
func TestASectionHeadingsLeadStandsBeforeItsName(t *testing.T) {
	for _, size := range leadSizes {
		t.Run(fmt.Sprintf("%dpx", size), func(t *testing.T) {
			var open, fixed, plain, tall *kvitui.SectionHeading
			var openLead, fixedLead, tallLead *unison.Panel
			s := uitest.Open(t, uitest.Options{InterfaceSize: size}, func(ui *kvitui.UI) unison.Paneler {
				var outer *unison.Panel
				open = kvitui.NewSectionHeading(ui, "The same words")
				open.Collapsible = true
				outer, openLead = leadPanel(ui, 16)
				open.SetLead(outer)
				fixed = kvitui.NewSectionHeading(ui, "The same words")
				outer, fixedLead = leadPanel(ui, 16)
				fixed.SetLead(outer)
				plain = kvitui.NewSectionHeading(ui, "The same words")
				tall = kvitui.NewSectionHeading(ui, "The same words")
				outer, tallLead = leadPanel(ui, 200)
				tall.SetLead(outer)
				return kvitui.Column(ui, kvitui.SizeSpace, kvitui.FullWidth(open), kvitui.FullWidth(fixed), kvitui.FullWidth(plain), kvitui.FullWidth(tall))
			})
			img := s.Capture()
			s.Do(func() {
				m := s.UI.Interface
				px16 := m.Px(16)
				box := rootRect(fixedLead)
				bar := rootRect(fixed)
				if want := bar.X + float32(m.Space()); box.X != want {
					t.Errorf("%d px: the lead of a heading with no chevron is at %v, want %v", size, box.X, want)
				}
				if d := box.CenterY() - bar.CenterY(); d < -0.5 || d > 0.5 || !wholePixels(box) {
					t.Errorf("%d px: the lead is at %v in a bar at %v", size, box, bar)
				}
				if want := rootRect(open).X + float32(m.Space()+m.IconSizeSmall()+m.StackGap()); rootRect(openLead).X != want {
					t.Errorf("%d px: the lead after a chevron is at %v, want %v", size, rootRect(openLead).X, want)
				}
				if !isLeadColour(img.At(int(box.CenterX()), int(box.CenterY()))) {
					t.Errorf("%d px: the lead is not drawn", size)
				}
				// The name is the plain heading's, moved along by the lead and a gap.
				words := image.Pt(int(bar.Width)-m.Space()-px16-m.StackGap()-m.Space(), int(bar.Height)-2)
				with := image.Pt(int(bar.X)+m.Space()+px16+m.StackGap(), int(bar.Y)+1)
				without := image.Pt(int(bar.X)+m.Space(), int(rootRect(plain).Y)+1)
				if differ, drawn := sameStretch(img, with, without, words); differ != 0 || drawn == 0 {
					t.Errorf("%d px: the name differs in %d pixels, with %d drawn", size, differ, drawn)
				}
				// One taller than the bar is drawn smaller, square still.
				big := rootRect(tallLead)
				if big.Height > bar.Height || big.Width != big.Height {
					t.Errorf("%d px: a tall lead is %v in a bar %v tall", size, big.Size, bar.Height)
				}
			})
		})
	}
}

// recordedPicture is a picture that says where it was drawn, and fills what it
// was given in one colour.
type recordedPicture struct {
	natural geom.Size
	at      []geom.Rect
	painted bool
}

func (p *recordedPicture) LogicalSize() geom.Size { return p.natural }

func (p *recordedPicture) DrawInRect(gc *unison.Canvas, r geom.Rect, _ *unison.SamplingOptions, paint *unison.Paint) {
	p.at = append(p.at, r)
	p.painted = p.painted || paint != nil
	gc.DrawRect(r, leadColour.Paint(gc, r, paintstyle.Fill))
}

// A menu line's picture goes in the symbol's column, fitted inside it with its
// shape kept, in place of the symbol, and the words stay where they are, at the
// smallest and largest interface sizes. A line's state takes the column from
// it: a tick, and the bullet of the chosen one of a set of options, while the
// others of the set have no mark and draw their pictures. A destructive line
// leaves its colours alone, a disabled one dims it, and a nil image counts as
// no picture.
func TestAMenuLinesPictureStandsInTheSymbolsColumn(t *testing.T) {
	for _, size := range leadSizes {
		t.Run(fmt.Sprintf("%dpx", size), func(t *testing.T) {
			wide := &recordedPicture{natural: geom.NewSize(40, 20)}
			tall := &recordedPicture{natural: geom.NewSize(20, 40)}
			large := &recordedPicture{natural: geom.NewSize(256, 256)}
			same := &recordedPicture{natural: geom.NewSize(16, 16)}
			both := &recordedPicture{natural: geom.NewSize(16, 16)}
			ticked := &recordedPicture{natural: geom.NewSize(16, 16)}
			option := &recordedPicture{natural: geom.NewSize(16, 16)}
			bullet := &recordedPicture{natural: geom.NewSize(16, 16)}
			danger := &recordedPicture{natural: geom.NewSize(16, 16)}
			dimmed := &recordedPicture{natural: geom.NewSize(16, 16)}
			var nilImage *unison.Image
			var button *kvitui.Button
			s := uitest.Open(t, uitest.Options{InterfaceSize: size}, func(ui *kvitui.UI) unison.Paneler {
				button = kvitui.NewButton(ui, "Add to project")
				return kvitui.Left(button)
			})
			lines := []kvitui.MenuItem{
				{Text: "Wide", Verbatim: true, Picture: wide},
				{Text: "Tall", Verbatim: true, Picture: tall},
				{Text: "Large", Verbatim: true, Picture: large},
				{Text: "The same words", Verbatim: true, Picture: same},
				{Text: "The same words", Verbatim: true, Symbol: "folder"},
				{Text: "The same words", Verbatim: true, Symbol: "folder", Picture: both},
				{Text: "Ticked", Verbatim: true, Picture: ticked, Checked: true},
				{Text: "Not chosen", Verbatim: true, Picture: option, Radio: true},
				{Text: "No image", Verbatim: true, Picture: nilImage, Symbol: "folder"},
				{Text: "Dimmed", Verbatim: true, Picture: dimmed, Disabled: true},
				{Text: "Chosen", Verbatim: true, Picture: bullet, Radio: true, Checked: true},
				{Separator: true},
				{Text: "Danger", Verbatim: true, Picture: danger, Danger: true},
			}
			s.Do(func() { s.UI.ShowMenu(button, "Add to project", lines) })
			s.WaitFor("the menu's lines", func() bool { return len(s.Nodes(role.MenuItem)) == len(lines)-1 })
			items := s.Nodes(role.MenuItem)
			slices.SortFunc(items, func(a, b *accessibility.Node) int { return cmp.Compare(a.Bounds.Y, b.Bounds.Y) })
			img := s.Capture()
			s.Do(func() {
				m := s.UI.Interface
				iconSize := float32(m.IconSizeSmall())
				if len(wide.at) == 0 {
					t.Fatalf("%d px: the picture was not drawn", size)
				}
				// The column, from the first line's left edge in the menu's own
				// coordinates: the recorded rectangles are in those.
				columnX := wide.at[len(wide.at)-1].X
				for _, c := range []struct {
					name string
					p    *recordedPicture
					w, h float32
				}{
					{"wide", wide, iconSize, iconSize / 2},
					{"tall", tall, iconSize / 2, iconSize},
					{"large", large, iconSize, iconSize},
					{"the same size as a symbol's", same, iconSize, iconSize},
					{"with a symbol", both, iconSize, iconSize},
					{"disabled", dimmed, iconSize, iconSize},
					{"unchosen option's", option, iconSize, iconSize},
					{"destructive", danger, iconSize, iconSize},
				} {
					if len(c.p.at) == 0 {
						t.Errorf("%d px: the %s picture was not drawn", size, c.name)
						continue
					}
					r := c.p.at[len(c.p.at)-1]
					if d := r.Width - c.w; d < -0.01 || d > 0.01 {
						t.Errorf("%d px: the %s picture is %v wide, want %v", size, c.name, r.Width, c.w)
					}
					if d := r.Height - c.h; d < -0.01 || d > 0.01 {
						t.Errorf("%d px: the %s picture is %v tall, want %v", size, c.name, r.Height, c.h)
					}
					if c.p.painted {
						t.Errorf("%d px: the %s picture was given a paint, which would replace an SVG's colours", size, c.name)
					}
					if !wholePixels(r) {
						t.Errorf("%d px: the %s picture is at %v, between pixels", size, c.name, r.Point)
					}
				}
				if d := tall.at[len(tall.at)-1].X - (columnX + iconSize/4); d < -1 || d > 1 {
					t.Errorf("%d px: a tall picture is %v off the column's middle", size, d)
				}
				if len(ticked.at) != 0 || len(bullet.at) != 0 {
					t.Errorf("%d px: a line's state left its picture drawn: ticked %d, chosen option %d", size, len(ticked.at), len(bullet.at))
				}
				// The words of the line with a picture and the line with a symbol
				// match pixel for pixel.
				line := func(i int) geom.Rect { return items[i].Bounds }
				lineOf := func(i int) image.Point {
					return image.Pt(int(line(i).X)+m.Space()+m.IconSizeSmall()+1, int(line(i).Y))
				}
				stretch := image.Pt(int(line(3).Width)-m.Space()-m.IconSizeSmall()-1, int(line(3).Height))
				if differ, drawn := sameStretch(img, lineOf(3), lineOf(4), stretch); differ != 0 || drawn == 0 {
					t.Errorf("%d px: the words of a picture's line and a symbol's differ in %d pixels, with %d drawn", size, differ, drawn)
				}
				// With both set, the line is the picture's alone.
				if differ, drawn := sameStretch(img, image.Pt(int(line(3).X), int(line(3).Y)), image.Pt(int(line(5).X), int(line(5).Y)),
					image.Pt(int(line(3).Width), int(line(3).Height))); differ != 0 || drawn == 0 {
					t.Errorf("%d px: a line with a symbol and a picture differs from one with the picture in %d pixels", size, differ)
				}
				// A destructive line leaves the picture's colours alone; a disabled
				// one dims them.
				middle := func(i int) color.Color {
					return img.At(int(line(i).X)+m.Space()+int(iconSize)/2, int(line(i).Y)+int(line(i).Height)/2)
				}
				if c := middle(11); !isLeadColour(c) {
					t.Errorf("%d px: the destructive line's picture is %v", size, c)
				}
				if c := middle(7); !isLeadColour(c) {
					t.Errorf("%d px: the unchosen option's picture is %v", size, c)
				}
				// The chosen option draws its mark there, which is not the tick
				// of the line above it.
				column := func(i int) image.Point { return image.Pt(int(line(i).X)+m.Space(), int(line(i).Y)) }
				mark := image.Pt(m.IconSizeSmall(), int(line(10).Height))
				if differ, drawn := sameStretch(img, column(10), column(6), mark); differ == 0 || drawn == 0 {
					t.Errorf("%d px: the chosen option's bullet and a tick differ in %d pixels, with %d drawn", size, differ, drawn)
				}
				if c := middle(9); isLeadColour(c) {
					t.Errorf("%d px: the disabled line's picture is not dimmed: %v", size, c)
				}
				// A nil image leaves the symbol.
				r := line(8)
				noImage := 0
				for y := int(r.Y); y < int(r.Bottom()); y++ {
					for x := int(r.X) + m.Space(); x < int(r.X)+m.Space()+m.IconSizeSmall(); x++ {
						if !samePixel(img.At(x, y), img.At(int(r.X), y)) {
							noImage++
						}
					}
				}
				if noImage == 0 {
					t.Errorf("%d px: a line with a nil image drew no symbol", size)
				}
				// The screen reader hears the words.
				if items[3].Name != "The same words" {
					t.Errorf("%d px: a picture's line is read as %q", size, items[3].Name)
				}
			})
			s.CheckNamed()
		})
	}
}

// An option of a set draws its own picture or symbol until it is chosen, when
// its mark takes the column: a bullet on the Windows form and a tick on the
// macOS form. It takes no symbol from its command, so a set of options made
// from commands stays as plain as its neighbours, whereas an ordinary line
// does. A line on for part of the selection draws no picture or symbol on
// either form, and a screen reader is told the same in every case.
func TestAnOptionLineDrawsItsPictureUntilItIsChosen(t *testing.T) {
	for _, form := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
		t.Run(form.String(), func(t *testing.T) {
			option := &recordedPicture{natural: geom.NewSize(16, 16)}
			chosen := &recordedPicture{natural: geom.NewSize(16, 16)}
			mixed := &recordedPicture{natural: geom.NewSize(16, 16)}
			run := func() {}
			var button *kvitui.Button
			s := uitest.Open(t, uitest.Options{Form: form}, func(ui *kvitui.UI) unison.Paneler {
				button = kvitui.NewButton(ui, "Move into")
				return kvitui.Left(button)
			})
			const (
				plain = iota
				ordinarySymbol
				optionPicture
				optionSymbol
				chosenPicture
				chosenSymbol
				chosenBare
				ticked
				mixedPicture
				mixedBare
				optionCommand
				optionCommandOwn
				ordinaryCommand
			)
			lines := []kvitui.MenuItem{
				plain:            {Text: "Words", Verbatim: true},
				ordinarySymbol:   {Text: "Words", Verbatim: true, Symbol: "folder"},
				optionPicture:    {Text: "Words", Verbatim: true, Radio: true, Picture: option},
				optionSymbol:     {Text: "Words", Verbatim: true, Radio: true, Symbol: "folder"},
				chosenPicture:    {Text: "Words", Verbatim: true, Radio: true, Checked: true, Picture: chosen},
				chosenSymbol:     {Text: "Words", Verbatim: true, Radio: true, Checked: true, Symbol: "folder"},
				chosenBare:       {Text: "Words", Verbatim: true, Radio: true, Checked: true},
				ticked:           {Text: "Words", Verbatim: true, Checked: true},
				mixedPicture:     {Text: "Words", Verbatim: true, Mixed: true, Picture: mixed, Symbol: "folder"},
				mixedBare:        {Text: "Words", Verbatim: true, Mixed: true},
				optionCommand:    {Radio: true, Command: &kvitui.Command{Name: "Words", Symbol: "folder", Run: run}},
				optionCommandOwn: {Radio: true, Symbol: "folder", Command: &kvitui.Command{Name: "Words", Symbol: "search", Run: run}},
				ordinaryCommand:  {Command: &kvitui.Command{Name: "Words", Symbol: "folder", Run: run}},
			}
			s.Do(func() { s.UI.ShowMenu(button, "Move into", lines) })
			s.WaitFor("the menu's lines", func() bool { return len(s.Nodes(role.MenuItem)) == len(lines) })
			items := s.Nodes(role.MenuItem)
			slices.SortFunc(items, func(a, b *accessibility.Node) int { return cmp.Compare(a.Bounds.Y, b.Bounds.Y) })
			img := s.Capture()
			s.Do(func() {
				m := s.UI.Interface
				column := func(i int) image.Point {
					return image.Pt(int(items[i].Bounds.X)+m.Space(), int(items[i].Bounds.Y))
				}
				size := image.Pt(m.IconSizeSmall(), int(items[plain].Bounds.Height))
				same := func(what string, a, b int) {
					t.Helper()
					if differ, drawn := sameStretch(img, column(a), column(b), size); differ != 0 || drawn == 0 {
						t.Errorf("%s: %d pixels differ from what it should match, with %d drawn", what, differ, drawn)
					}
				}
				blank := func(what string, i int) {
					t.Helper()
					if _, drawn := sameStretch(img, column(i), column(plain), size); drawn != 0 {
						t.Errorf("%s: the column has %d pixels drawn, want none", what, drawn)
					}
				}
				blank("a line with nothing to draw", plain)
				if len(option.at) == 0 {
					t.Error("an unchosen option did not draw its picture")
				}
				same("an unchosen option's symbol", optionSymbol, ordinarySymbol)
				if len(chosen.at) != 0 {
					t.Error("the chosen option drew its picture")
				}
				same("the chosen option in place of a picture", chosenPicture, chosenBare)
				same("the chosen option in place of a symbol", chosenSymbol, chosenBare)
				differ, drawn := sameStretch(img, column(chosenBare), column(ticked), size)
				if form == kvitui.FormMac && (differ != 0 || drawn == 0) {
					t.Errorf("the chosen option differs from a tick in %d pixels on the macOS form, with %d drawn", differ, drawn)
				}
				if form == kvitui.FormWindows && (differ == 0 || drawn == 0) {
					t.Errorf("the chosen option differs from a tick in %d pixels on the Windows form, with %d drawn", differ, drawn)
				}
				if len(mixed.at) != 0 {
					t.Error("a line on for part of the selection drew its picture")
				}
				if form == kvitui.FormMac {
					same("a dash in place of a mixed line's picture", mixedPicture, mixedBare)
				} else {
					blank("a mixed line on the Windows form", mixedPicture)
				}
				blank("an option made from a command with a symbol", optionCommand)
				same("an option's own symbol over its command's", optionCommandOwn, ordinarySymbol)
				same("an ordinary line made from a command with a symbol", ordinaryCommand, ordinarySymbol)
				for _, c := range []struct {
					name     string
					i        int
					hasCheck bool
					checked  check.Enum
				}{
					{"an unchosen option", optionPicture, true, check.Off},
					{"the chosen option", chosenPicture, true, check.On},
					{"a mixed line", mixedPicture, true, check.Mixed},
					{"a plain line", plain, false, check.Off},
				} {
					if n := items[c.i]; n.Name != "Words" || n.HasCheck != c.hasCheck || n.Checked != c.checked {
						t.Errorf("%s is told to a screen reader as %+v", c.name, n)
					}
				}
			})
			s.CheckNamed()
		})
	}
}

// On the Windows form the chosen option's bullet is a solid dot in the accent,
// large enough to find while scanning the column of marks: at least half the
// column's width at every interface size, round, and in the column's middle.
func TestTheChosenOptionsBulletIsLargeEnoughToSee(t *testing.T) {
	for _, theme := range tokens.BuiltInThemes() {
		for _, size := range []int{tokens.MinInterfaceSize, tokens.DefaultInterfaceSize, tokens.MaxInterfaceSize} {
			t.Run(fmt.Sprintf("%s-%d", theme, size), func(t *testing.T) {
				var button *kvitui.Button
				s := uitest.Open(t, uitest.Options{Theme: theme, InterfaceSize: size}, func(ui *kvitui.UI) unison.Paneler {
					button = kvitui.NewButton(ui, "Move into")
					return kvitui.Left(button)
				})
				lines := []kvitui.MenuItem{
					{Text: "Chosen", Verbatim: true, Radio: true, Checked: true},
					{Text: "Other", Verbatim: true, Radio: true},
				}
				s.Do(func() { s.UI.ShowMenu(button, "Move into", lines) })
				s.WaitFor("the menu's lines", func() bool { return len(s.Nodes(role.MenuItem)) == len(lines) })
				items := s.Nodes(role.MenuItem)
				slices.SortFunc(items, func(a, b *accessibility.Node) int { return cmp.Compare(a.Bounds.Y, b.Bounds.Y) })
				img := s.Capture()
				s.Do(func() {
					m, line := s.UI.Interface, items[0].Bounds
					column := m.IconSizeSmall()
					x0 := int(line.X) + m.Space()
					ground := img.At(x0, int(line.Y)+1)
					left, right, top, bottom := math.MaxInt, -1, math.MaxInt, -1
					// Across the column and past it, up to where the words start,
					// and the line's whole height.
					for y := int(line.Y); y < int(line.Bottom()); y++ {
						for x := int(line.X); x < x0+column+m.Space(); x++ {
							if !samePixel(img.At(x, y), ground) {
								left, right = min(left, x), max(right, x)
								top, bottom = min(top, y), max(bottom, y)
							}
						}
					}
					if right < 0 {
						t.Fatal("the chosen line has no bullet")
					}
					wide, high := right-left+1, bottom-top+1
					if wide*2 < column || high*2 < column {
						t.Errorf("the bullet is %d by %d in a column %d wide", wide, high, column)
					}
					if wide != high {
						t.Errorf("the bullet is %d by %d, not round", wide, high)
					}
					middleX, middleY := float64(left+right+1)/2, float64(top+bottom+1)/2
					if want := float64(x0) + float64(column)/2; math.Abs(middleX-want) > 1 {
						t.Errorf("the bullet's middle is at x %v, the column's at %v", middleX, want)
					}
					if want := float64(line.Y) + float64(line.Height)/2; math.Abs(middleY-want) > 1 {
						t.Errorf("the bullet's middle is at y %v, the line's at %v", middleY, want)
					}
					accent := kvitui.Color(s.UI.Theme.Tokens().Accent)
					if c := img.At((left+right)/2, (top+bottom)/2); !samePixel(c, accent) {
						t.Errorf("the bullet's middle is %v, the accent is %v", c, accent)
					}
				})
				s.CheckNamed()
			})
		}
	}
}

// A select's options draw their pictures and symbols on their lines in its
// list, as a set of options in a menu does: an unchosen option its picture or
// its symbol, one with neither nothing, and the chosen option, whatever it has,
// its bullet on the Windows form and its tick on the macOS form. A screen
// reader is told each line's words alone. The closed select, in its field,
// compact, plain and fused forms, draws the same pixels as one whose options
// have neither, so its words stand where a plain select's always have.
func TestASelectsOptionsDrawTheirPicturesInItsList(t *testing.T) {
	for _, form := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
		t.Run(form.String(), func(t *testing.T) {
			picture := &recordedPicture{natural: geom.NewSize(16, 16)}
			chosen := &recordedPicture{natural: geom.NewSize(16, 16)}
			const (
				hasPicture = iota
				hasSymbol
				hasNeither
				isChosen
			)
			decorated := []kvitui.Option{
				hasPicture: {Value: "picture", Label: "Has a picture", Picture: picture},
				hasSymbol:  {Value: "symbol", Label: "Has a symbol", Symbol: "folder"},
				hasNeither: {Value: "neither", Label: "Has neither"},
				isChosen:   {Value: "chosen", Label: "Chosen", Picture: chosen, Symbol: "folder"},
			}
			bare := make([]kvitui.Option, len(decorated))
			for i, o := range decorated {
				bare[i] = kvitui.Option{Value: o.Value, Label: o.Label}
			}
			type pair struct {
				form          string
				with, without unison.Paneler
			}
			var pairs []pair
			var field *kvitui.Select
			s := uitest.Open(t, uitest.Options{Form: form}, func(ui *kvitui.UI) unison.Paneler {
				choice := func(options []kvitui.Option, how func(*kvitui.Select)) *kvitui.Select {
					c := kvitui.NewSelect(ui, "Move to", options...)
					c.Current = "chosen"
					how(c)
					return c
				}
				fused := func(options []kvitui.Option) *kvitui.ChoiceButton {
					c := kvitui.NewChoiceButton(ui, "Move", "Where to", options...)
					c.Current = "chosen"
					return c
				}
				asField := func(*kvitui.Select) {}
				compact := func(c *kvitui.Select) { c.Compact, c.Symbol = true, "calendar" }
				plain := func(c *kvitui.Select) { c.Plain = true }
				field = choice(decorated, asField)
				pairs = []pair{
					{"field", kvitui.Width(ui, kvitui.Px(200), field), kvitui.Width(ui, kvitui.Px(200), choice(bare, asField))},
					{"compact", kvitui.Left(choice(decorated, compact)), kvitui.Left(choice(bare, compact))},
					{"plain", kvitui.Width(ui, kvitui.Px(200), choice(decorated, plain)), kvitui.Width(ui, kvitui.Px(200), choice(bare, plain))},
					{"fused", kvitui.Left(fused(decorated)), kvitui.Left(fused(bare))},
				}
				// A button first takes the focus the window gives its first
				// control, which a field select would draw as its edge.
				children := []unison.Paneler{kvitui.Left(kvitui.NewButton(ui, "First"))}
				for _, p := range pairs {
					children = append(children, p.with, p.without)
				}
				return kvitui.Column(ui, kvitui.Px(12), children...)
			})
			s.Screen.MouseMove(geom.NewPoint(1400, 900), 0)
			s.Sync()

			img := s.Capture()
			for _, p := range pairs {
				a, b := rootRect(p.with), rootRect(p.without)
				if a.Size != b.Size {
					t.Errorf("the closed %s select is %v with pictures and %v without", p.form, a.Size, b.Size)
					continue
				}
				size := image.Pt(int(a.Width), int(a.Height))
				if differ, drawn := sameStretch(img, image.Pt(int(a.X), int(a.Y)), image.Pt(int(b.X), int(b.Y)), size); differ != 0 || drawn == 0 {
					t.Errorf("the closed %s select differs in %d pixels from one with no pictures, with %d drawn", p.form, differ, drawn)
				}
			}
			if len(picture.at) != 0 || len(chosen.at) != 0 {
				t.Errorf("the closed selects drew their options' pictures %d and %d times", len(picture.at), len(chosen.at))
			}
			if n := s.Node(field); n == nil || n.Name != "Move to" || n.Value != "Chosen" {
				t.Errorf("the closed select is told to a screen reader as %+v", n)
			}

			// The list, and what the menu was given for each line.
			box := rootRect(field)
			s.Screen.Click(geom.NewPoint(box.X+box.Width/2, box.Y+box.Height/2))
			s.WaitFor("the select's list", func() bool { return len(s.Nodes(role.MenuItem)) == len(decorated) })
			s.Do(func() {
				menus := s.Window.ShownMenus()
				if len(menus) != 1 || len(menus[0].Items) != len(decorated) {
					t.Fatalf("the select opened %+v", menus)
				}
				for i, it := range menus[0].Items {
					o := decorated[i]
					if it.Symbol != o.Symbol || it.Picture != o.Picture || !it.Radio || it.Checked != (i == isChosen) {
						t.Errorf("the line for %q is %+v", o.Label, it)
					}
				}
			})
			lines := s.Nodes(role.MenuItem)
			slices.SortFunc(lines, func(a, b *accessibility.Node) int { return cmp.Compare(a.Bounds.Y, b.Bounds.Y) })
			for i, n := range lines {
				want := check.Off
				if i == isChosen {
					want = check.On
				}
				if n.Name != decorated[i].Label || !n.HasCheck || n.Checked != want {
					t.Errorf("the line for %q is told to a screen reader as %+v", decorated[i].Label, n)
				}
			}
			list := s.Capture()

			// What the same lines draw in a menu of MenuItems, to hold the
			// select's lines against.
			s.Screen.KeyPress(unison.KeyEscape, 0)
			s.WaitFor("the list to close", func() bool { return len(s.Nodes(role.MenuItem)) == 0 })
			const (
				refSymbol = iota
				refChosen
				refTick
			)
			reference := []kvitui.MenuItem{
				refSymbol: {Text: "Has a symbol", Verbatim: true, Symbol: "folder"},
				refChosen: {Text: "Chosen", Verbatim: true, Radio: true, Checked: true},
				refTick:   {Text: "Chosen", Verbatim: true, Checked: true},
			}
			s.Do(func() { s.UI.ShowMenu(field, "Move to", reference) })
			s.WaitFor("the reference menu", func() bool { return len(s.Nodes(role.MenuItem)) == len(reference) })
			refLines := s.Nodes(role.MenuItem)
			slices.SortFunc(refLines, func(a, b *accessibility.Node) int { return cmp.Compare(a.Bounds.Y, b.Bounds.Y) })
			ref := s.Capture()

			s.Do(func() {
				m := s.UI.Interface
				column := func(n *accessibility.Node) image.Point {
					return image.Pt(int(n.Bounds.X)+m.Space(), int(n.Bounds.Y))
				}
				size := image.Pt(m.IconSizeSmall(), int(lines[0].Bounds.Height))
				// across compares a column of the select's list with one of the
				// reference menu.
				across := func(what string, line int, refLine int) {
					t.Helper()
					a, b := column(lines[line]), column(refLines[refLine])
					differ, drawn := 0, 0
					ground := list.At(a.X, a.Y)
					for y := range size.Y {
						for x := range size.X {
							pa := list.At(a.X+x, a.Y+y)
							if !samePixel(pa, ref.At(b.X+x, b.Y+y)) {
								differ++
							}
							if !samePixel(pa, ground) {
								drawn++
							}
						}
					}
					if differ != 0 || drawn == 0 {
						t.Errorf("%s: %d pixels differ from what it should match, with %d drawn", what, differ, drawn)
					}
				}
				if len(picture.at) == 0 {
					t.Error("the unchosen option with a picture did not draw it")
				}
				middle := column(lines[hasPicture]).Add(image.Pt(m.IconSizeSmall()/2, int(lines[hasPicture].Bounds.Height)/2))
				if c := list.At(middle.X, middle.Y); !isLeadColour(c) {
					t.Errorf("the middle of the picture's column is %v", c)
				}
				across("the unchosen option's symbol", hasSymbol, refSymbol)
				if _, drawn := sameStretch(list, column(lines[hasNeither]), column(lines[hasNeither]), size); drawn != 0 {
					t.Errorf("the option with neither has %d pixels drawn in its column, want none", drawn)
				}
				if len(chosen.at) != 0 {
					t.Error("the chosen option drew its picture")
				}
				across("the chosen option's mark", isChosen, refChosen)
				differ, drawn := sameStretch(ref, column(refLines[refChosen]), column(refLines[refTick]), size)
				if form == kvitui.FormWindows && (differ == 0 || drawn == 0) {
					t.Errorf("the chosen option's mark differs from a tick in %d pixels on the Windows form, with %d drawn", differ, drawn)
				}
				if form == kvitui.FormMac && (differ != 0 || drawn == 0) {
					t.Errorf("the chosen option's mark differs from a tick in %d pixels on the macOS form, with %d drawn", differ, drawn)
				}
			})
			s.CheckNamed()
		})
	}
}
