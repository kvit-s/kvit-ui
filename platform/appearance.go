// Package platform reports what the desktop says about appearance: dark or
// light, high contrast, and reduced motion. Each answer is false where the
// desktop gives none, which is the safe direction: a person who has turned
// neither on sees no change, and the application's own settings stay in
// charge. Nothing here is ever written back to the desktop.
package platform

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Appearance is the desktop's current answers. It implements
// tokens.Appearance, and announces changes to anything connected with
// OnChanged once Refresh finds one.
type Appearance struct {
	mu        sync.Mutex
	dark      bool
	high      bool
	reduced   bool
	available bool
	listeners []func()
}

// NewAppearance reads the desktop once.
func NewAppearance() *Appearance {
	a := &Appearance{}
	a.dark, a.high, a.reduced, a.available = readPlatform()
	return a
}

// DarkMode reports whether the desktop prefers dark.
func (a *Appearance) DarkMode() bool { a.mu.Lock(); defer a.mu.Unlock(); return a.dark }

// HighContrast reports whether the desktop is in a high-contrast mode.
func (a *Appearance) HighContrast() bool { a.mu.Lock(); defer a.mu.Unlock(); return a.high }

// ReducedMotion reports whether the desktop asks applications to still their
// animations.
func (a *Appearance) ReducedMotion() bool { a.mu.Lock(); defer a.mu.Unlock(); return a.reduced }

// Available reports whether this platform answers at all, so a "follow the
// system" choice can say when there is nothing to follow.
func (a *Appearance) Available() bool { a.mu.Lock(); defer a.mu.Unlock(); return a.available }

// OnChanged calls fn after Refresh finds any answer changed.
func (a *Appearance) OnChanged(fn func()) (disconnect func()) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.listeners = append(a.listeners, fn)
	i := len(a.listeners) - 1
	return func() {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.listeners[i] = nil
	}
}

// Refresh reads the desktop again and notifies if anything moved. Call it
// when the desktop announces a change (unison's theme-changed callback does),
// and when a settings dialog opens.
func (a *Appearance) Refresh() {
	dark, high, reduced, available := readPlatform()
	a.mu.Lock()
	changed := dark != a.dark || high != a.high || reduced != a.reduced
	a.dark, a.high, a.reduced, a.available = dark, high, reduced, available
	fns := append([]func(){}, a.listeners...)
	a.mu.Unlock()
	if changed {
		for _, fn := range fns {
			if fn != nil {
				fn()
			}
		}
	}
}

// Set forces all three answers, for tests: no machine running them has high
// contrast turned on, and what matters is on this side of the platform.
func (a *Appearance) Set(dark, highContrast, reducedMotion bool) {
	a.mu.Lock()
	changed := dark != a.dark || highContrast != a.high || reducedMotion != a.reduced
	a.dark, a.high, a.reduced, a.available = dark, highContrast, reducedMotion, true
	fns := append([]func(){}, a.listeners...)
	a.mu.Unlock()
	if changed {
		for _, fn := range fns {
			if fn != nil {
				fn()
			}
		}
	}
}

// commandOutput runs a short command and returns its trimmed output, or ""
// if it is missing, slow or fails. The timeout makes it safe at startup: a
// desktop without the tool answers by not having it, and one with a wedged
// settings daemon by timing out rather than hanging the application.
func commandOutput(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
