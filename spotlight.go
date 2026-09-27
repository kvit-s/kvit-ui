package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// Spotlight darkens a region except for one part of it and says something
// about that part. It is the primitive under a guided tour and deliberately
// only that: kvit-cash's first-run flow reads a fixed script while
// kvit-notes-pro's decides the next step from what the agent is doing, so a
// shared stepper would be wrong for one of them; what is shared is lighting a
// part, and each application brings its own reason to move on. Escape always
// closes it, and so does a press anywhere on the shade, since a reader who
// cannot get out of a tour force-quits the application.
type Spotlight struct {
	unison.Panel
	ui *UI
	// Target is the part left lit; nil darkens everything, the state before
	// the first step.
	Target unison.Paneler
	// Title and Detail are what is said about it.
	Title, Detail string
	// Padding is the room around the target inside the lit hole; a space
	// unless set.
	Padding Measure
	// OnDismiss runs when the reader asks to leave, by Escape or a press.
	OnDismiss func()

	over  unison.Paneler // the region darkened
	card  *Card
	title *Label
	words *Label
	hide  func()
}

// NewSpotlight returns a spotlight on target, closed.
func NewSpotlight(ui *UI, target unison.Paneler, title, detail string) *Spotlight {
	s := &Spotlight{ui: ui, Target: target, Title: title, Detail: detail, Padding: SizeSpace}
	s.Self = s
	s.title = NewLabel(ui, "")
	s.title.Role, s.title.Weight = RoleStrong, text.Bold
	s.words = NewLabel(ui, "")
	s.words.Role, s.words.Ink, s.words.Wrap = RoleSmall, InkTextSecondary, true
	// The Qt card's column is the card's padding narrower again than the
	// room inside the card, so its words wrap that much sooner.
	twice := func(m *tokens.Interface) int { return 2 * m.SpaceLoose() }
	s.title.SetBorder(Insets(ui, nil, nil, nil, twice))
	s.words.SetBorder(Insets(ui, nil, nil, nil, twice))
	s.card = NewCard(ui, s.title, s.words)
	s.card.SetLayout(&spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceNear})
	s.AddChild(s.card)
	s.SetLayout(syncing{Layout: spotlightLayout{s}, sync: func() {
		s.title.Text, s.words.Text = s.Title, s.Detail
		s.words.Hidden = s.Detail == ""
		s.card.Hidden = s.Title == "" && s.Detail == ""
		showOnly(s.card.AsPanel(), s.title, s.words)
	}})
	s.DrawCallback = s.draw
	s.MouseDownCallback = func(geom.Point, int, int, mod.Modifiers) bool {
		s.dismiss()
		return true
	}
	return s
}

// Open darkens over, which must be in a Kvit window, lighting the target.
func (s *Spotlight) Open(over unison.Paneler) {
	w := s.ui.windowOf(over)
	if w == nil || s.hide != nil {
		return
	}
	s.over = over
	s.hide = w.Show(&Popup{Panel: s, Place: func(geom.Rect, geom.Size) geom.Rect { return anchorIn(over) },
		OnEscape: s.dismiss, Anchor: over})
	unison.AnnounceForAccessibility(joinLines(s.Title, s.Detail))
}

// Opened reports whether the spotlight is showing.
func (s *Spotlight) Opened() bool { return s.hide != nil }

// Close takes the spotlight away.
func (s *Spotlight) Close() {
	if s.hide != nil {
		s.hide()
		s.hide = nil
	}
}

func (s *Spotlight) dismiss() {
	if s.OnDismiss != nil {
		s.OnDismiss()
	}
}

// hole is the lit part, in the spotlight's coordinates, or an empty box.
func (s *Spotlight) hole() geom.Rect {
	if s.Target == nil || s.Target.AsPanel().Window() == nil {
		return geom.Rect{}
	}
	t := s.Target.AsPanel()
	r := s.RectFromRoot(t.RectToRoot(t.ContentRect(true)))
	pad := float32(s.Padding.Of(s.ui))
	return r.Inset(geom.NewUniformInsets(-pad))
}

type spotlightLayout struct{ s *Spotlight }

func (l spotlightLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	return geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l spotlightLayout) PerformLayout(target *unison.Panel) {
	s, m := l.s, l.s.ui.Interface
	r := target.ContentRect(false)
	hole := s.hole()
	w := float32(m.Px(280))
	_, p, _ := s.card.Sizes(geom.NewSize(w, 0))
	margin, gap := float32(m.ViewMargin()), float32(m.Space())
	x := max(margin, min(r.Width-w-margin, hole.X))
	// Below the lit part where it fits, and above it where it does not.
	y := hole.Bottom() + gap
	if y+p.Height >= r.Height {
		y = max(margin, hole.Y-p.Height-gap)
	}
	s.card.SetFrameRect(geom.NewRect(r.X+x, r.Y+y, w, p.Height))
}

func (s *Spotlight) draw(gc *unison.Canvas, _ geom.Rect) {
	t, m := s.ui.Theme.Tokens(), s.ui.Interface
	r := s.ContentRect(false)
	hole := s.hole()
	// The shade as four boxes around the hole, which leave the lit part
	// untouched.
	shade := unison.RGB(0, 0, 0).SetAlphaIntensity(0.55)
	for _, b := range []geom.Rect{
		geom.NewRect(r.X, r.Y, r.Width, max(0, hole.Y-r.Y)),
		geom.NewRect(hole.Right(), hole.Y, max(0, r.Right()-hole.Right()), hole.Height),
		geom.NewRect(r.X, hole.Bottom(), r.Width, max(0, r.Bottom()-hole.Bottom())),
		geom.NewRect(r.X, hole.Y, max(0, hole.X-r.X), hole.Height),
	} {
		if hole.Empty() {
			b = r
		}
		gc.DrawRect(b, shade.Paint(gc, b, paintstyle.Fill))
		if hole.Empty() {
			break
		}
	}
	if !hole.Empty() {
		// Outlined, so the lit part reads as chosen rather than as a gap in
		// the shade.
		painterFor(gc, s.ui).outline(hole, float32(m.RadiusCard()), float32(m.FocusRingWidth()), t.FocusRing)
	}
}

// ProvideAccessibility describes the spotlight as a dialog saying its title
// and detail.
func (s *Spotlight) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Dialog
	}
	n.Name = s.Title
	n.Description = s.Detail
}
