package kvitui

import (
	"math"
	"runtime"
	"time"

	"github.com/kvit-s/kvit-ui/platform"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// wheelScroll moves a scrolling view the distance the wheel was turned: a
// notch is the desktop's lines per notch, of one slim row each, eased in over
// a few frames. Distance a finer wheel reports in smaller steps is added up
// rather than replaced, so turning twice as far travels twice as far however
// the turn arrives. Every turn moves its share at once, even one begun while
// the last is still easing in, so reversing out of a fast spin answers in
// its own direction rather than paying off what is still owed first. A view
// with nothing to scroll leaves the wheel to whatever scrolling area it sits
// in.
type wheelScroll struct {
	ui *UI
	// at is where the view is scrolled to, across and down.
	at func() (x, y float32)
	// travel is how far it can scroll, across and down.
	travel func() (x, y float32)
	// moveTo scrolls it, kept within its ends.
	moveTo func(x, y float32)

	owedX, owedY float32
	paying       bool
}

func (s *wheelScroll) wheel(delta geom.Point) bool {
	tx, ty := s.travel()
	if tx <= 0 && ty <= 0 {
		return false
	}
	x, y := s.at()
	if runtime.GOOS == "darwin" {
		// A Mac's trackpad and wheel report distance already, which unison
		// scales to pixels; those move the view at once.
		s.moveTo(x-delta.X*unison.MouseWheelMultiplier, y-delta.Y*unison.MouseWheelMultiplier)
		return true
	}
	notch := float32(platform.WheelScrollLines() * s.ui.Interface.RowHeightSlim())
	if s.paying {
		// A turn begun while the last is still easing in moves its share
		// at once, in its own direction, with the rest joining what is
		// owed. Without this the turn only pays off what is still owed,
		// so reversing out of a fast spin goes nowhere until enough turns
		// have been spun off. moveTo keeps the move within the ends; what
		// is left past either end is dropped below as before.
		share := s.share()
		x, y := s.at()
		s.moveTo(x-delta.X*notch*share, y-delta.Y*notch*share)
		nx, ny := s.at()
		s.owedY = clampOwed(s.owedY-(delta.Y*notch+(ny-y)), ny, ty)
		s.owedX = clampOwed(s.owedX-(delta.X*notch+(nx-x)), nx, tx)
		return true
	}
	// What is owed past either end is dropped as it is asked for, so a long
	// spin against the bottom does not build a debt that has to be spun off.
	s.owedY = clampOwed(s.owedY-delta.Y*notch, y, ty)
	s.owedX = clampOwed(s.owedX-delta.X*notch, x, tx)
	s.pay()
	return true
}

func clampOwed(owed, at, travel float32) float32 { return max(-at, min(travel-at, owed)) }

// share is the portion of what is owed paid at once: from the frame's
// length, so the motion is the same at any refresh rate, and all of it when
// motion is reduced.
func (s *wheelScroll) share() float32 {
	const frame = 16 * time.Millisecond
	settle := 35 * s.ui.Theme.MotionScale()
	if settle <= 0 {
		return 1
	}
	return float32(1 - math.Exp(-float64(frame.Milliseconds())/settle))
}

// pay moves the view a share of what the wheel is owed, once a frame, until
// it is paid. The share comes from the frame's length, so the motion is the
// same at any refresh rate, and it is all paid at once when motion is
// reduced.
func (s *wheelScroll) pay() {
	if s.paying {
		return
	}
	const frame = 16 * time.Millisecond
	share := s.share()
	step := func(owed float32) float32 {
		move := owed * share
		if d := owed - move; d < 0.5 && d > -0.5 {
			move = owed
		}
		return move
	}
	dx, dy := step(s.owedX), step(s.owedY)
	s.owedX -= dx
	s.owedY -= dy
	x, y := s.at()
	s.moveTo(x+dx, y+dy)
	if s.owedX != 0 || s.owedY != 0 {
		s.paying = true
		unison.InvokeTaskAfter(func() {
			s.paying = false
			s.pay()
		}, frame)
	}
}
