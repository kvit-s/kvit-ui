// Package settings keeps per-user application settings in one JSON file:
// theme, interface size, typography, option states. It writes the same
// format as the Qt library's SettingsStore, a single indented JSON object, so
// the Qt and Go versions of an app read each other's settings.
package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

// WriteDelay coalesces a burst of changes, such as a slider drag, into one
// write.
const WriteDelay = 500 * time.Millisecond

// Store holds the settings in memory and writes them to its file shortly
// after each change. Unknown keys survive a round trip, because the whole
// object is held and written back. A Store is safe for use from several
// goroutines.
type Store struct {
	mu          sync.Mutex
	path        string
	values      map[string]any
	dirty       bool
	timer       *time.Timer
	onChange    []func(key string)
	onWriteFail []func(path string, err error)
}

// New returns an empty store bound to no file; values set on it are held
// until Open binds it to one.
func New() *Store { return &Store{values: map[string]any{}} }

// Open binds the store to a file and loads it. A missing or corrupt file
// gives an empty store; the file is created on the first write.
//
// If a write to the previous file is pending and fails, Open refuses and the
// store stays bound to the previous file with its values, since those values
// are then the only copy; discardPending gives them up instead. Open also
// refuses when the file's directory cannot be created.
func (s *Store) Open(path string, discardPending bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dirty {
		if err := s.writeLocked(); err != nil && !discardPending {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	s.stopTimerLocked()
	s.path = path
	s.values = map[string]any{}
	s.dirty = false
	if data, err := os.ReadFile(path); err == nil {
		var v map[string]any
		if json.Unmarshal(data, &v) == nil && v != nil {
			s.values = v
		}
	}
	return nil
}

// Path is the file the store is bound to.
func (s *Store) Path() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.path
}

// Value returns a stored value as JSON decodes it (string, float64, bool,
// []any, map[string]any) and whether it was there.
func (s *Store) Value(key string) (any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	return v, ok
}

// Contains reports whether a key is stored.
func (s *Store) Contains(key string) bool {
	_, ok := s.Value(key)
	return ok
}

// SetValue stores a value and schedules a write. Setting the value a key
// already has does nothing. The value is normalised through JSON, so an int
// reads back as float64, as it would after a restart.
func (s *Store) SetValue(key string, value any) {
	norm := normalise(value)
	s.mu.Lock()
	if old, ok := s.values[key]; ok && reflect.DeepEqual(old, norm) {
		s.mu.Unlock()
		return
	}
	s.values[key] = norm
	s.scheduleLocked()
	fns := append([]func(string){}, s.onChange...)
	s.mu.Unlock()
	for _, fn := range fns {
		fn(key)
	}
}

// Remove deletes a key and schedules a write.
func (s *Store) Remove(key string) {
	s.mu.Lock()
	if _, ok := s.values[key]; !ok {
		s.mu.Unlock()
		return
	}
	delete(s.values, key)
	s.scheduleLocked()
	fns := append([]func(string){}, s.onChange...)
	s.mu.Unlock()
	for _, fn := range fns {
		fn(key)
	}
}

// Flush writes now if anything is pending, for quitting and for tests.
func (s *Store) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.dirty {
		return nil
	}
	return s.writeLocked()
}

// OnValueChanged calls fn with the key after every real change.
func (s *Store) OnValueChanged(fn func(key string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onChange = append(s.onChange, fn)
}

// OnWriteFailed calls fn when a write does not reach disk. The values stay
// pending and are written by the next flush or change, so this is a warning
// the user can act on rather than a report of lost data.
func (s *Store) OnWriteFailed(fn func(path string, err error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onWriteFail = append(s.onWriteFail, fn)
}

func (s *Store) scheduleLocked() {
	s.dirty = true
	s.stopTimerLocked()
	s.timer = time.AfterFunc(WriteDelay, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.dirty {
			_ = s.writeLocked()
		}
	})
}

func (s *Store) stopTimerLocked() {
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
}

// writeLocked writes the whole object through a temporary file renamed into
// place, so a crash leaves the old file or the new one. The dirty flag is
// cleared only once the bytes are there.
func (s *Store) writeLocked() error {
	s.stopTimerLocked()
	if s.path == "" {
		s.dirty = false // nowhere to write: nothing is being lost
		return nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	if err := enc.Encode(s.values); err != nil {
		return s.failLocked(err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".settings-*")
	if err != nil {
		return s.failLocked(err)
	}
	_, werr := tmp.Write(buf.Bytes())
	cerr := tmp.Close()
	if err = errors.Join(werr, cerr); err == nil {
		err = os.Rename(tmp.Name(), s.path)
	}
	if err != nil {
		os.Remove(tmp.Name())
		return s.failLocked(err)
	}
	s.dirty = false
	return nil
}

func (s *Store) failLocked(err error) error {
	// The change stays pending for the next flush or change. No timer is
	// restarted: a location that cannot be written usually stays that way.
	s.dirty = true
	path := s.path
	for _, fn := range s.onWriteFail {
		go fn(path, err)
	}
	return err
}

func normalise(v any) any {
	data, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var out any
	if json.Unmarshal(data, &out) != nil {
		return v
	}
	return out
}
