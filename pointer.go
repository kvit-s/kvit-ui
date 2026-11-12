package kvitui

import (
	"weak"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// pointerWatchedKey is the client-data key under which a window keeps the
// listeners OnPointer gave it.
const pointerWatchedKey = "kvitui.pointerWatched"

// pointerListeners are a window's pointer listeners by their owners, held
// weakly, so neither a closed window nor a dropped panel is kept alive.
type pointerListeners map[weak.Pointer[unison.Panel]]func(where geom.Point, inside bool)

// OnPointer has a window tell f where the pointer is after every move, in the
// window's root coordinates, for as long as owner is alive; inside is false
// once the pointer has left the window. unison tells only the innermost panel
// under the pointer that the pointer arrived or left, so a panel that has to
// know the pointer is on anything inside it listens here, and asks PointerOn
// or PointerOver.
//
// A second call for the same owner replaces its listener. The window holds
// the owner weakly and f strongly, so f must not hold the owner, or a
// component the owner is part of, except through a weak pointer: otherwise
// the owner is never dropped. The pointer is also reported as gone and back
// when it crosses from one panel to the next, so a listener that acts on a
// change does so once the event is over.
func (u *UI) OnPointer(w *unison.Window, owner *unison.Panel, f func(where geom.Point, inside bool)) {
	if w == nil || owner == nil {
		return
	}
	listeners, ok := w.ClientData()[pointerWatchedKey].(pointerListeners)
	if !ok {
		listeners = pointerListeners{}
		w.ClientData()[pointerWatchedKey] = listeners
		tell := func(where geom.Point, inside bool) {
			for p, f := range listeners {
				if p.Value() == nil {
					delete(listeners, p)
					continue
				}
				f(where, inside)
			}
		}
		move, enter, exit := w.MouseMoveCallback, w.MouseEnterCallback, w.MouseExitCallback
		w.MouseMoveCallback = func(where geom.Point, mods mod.Modifiers) bool {
			tell(where, true)
			return move != nil && move(where, mods)
		}
		w.MouseEnterCallback = func(where geom.Point, mods mod.Modifiers) bool {
			tell(where, true)
			return enter != nil && enter(where, mods)
		}
		w.MouseExitCallback = func() bool {
			tell(geom.Point{}, false)
			return exit != nil && exit()
		}
	}
	listeners[weak.Make(owner)] = f
}

// PointerOver reports whether a root point is inside p's box, p being in a
// window and it and everything it sits in shown. The point and inside are
// what an OnPointer listener is given.
func PointerOver(p *unison.Panel, where geom.Point, inside bool) bool {
	if !inside || p == nil || p.Window() == nil {
		return false
	}
	for q := p; q != nil; q = q.Parent() {
		if q.Hidden {
			return false
		}
	}
	return p.PointFromRoot(where).In(p.ContentRect(true))
}

// PointerOn reports whether the innermost panel under a root point is p or
// inside p, so a part of p scrolled out of view or under a popup does not
// count.
func PointerOn(p *unison.Panel, where geom.Point, inside bool) bool {
	if !PointerOver(p, where, inside) {
		return false
	}
	content := p.Window().Content()
	for q := content.PanelAt(content.PointFromRoot(where)); q != nil; q = q.Parent() {
		if q == p {
			return true
		}
	}
	return false
}
