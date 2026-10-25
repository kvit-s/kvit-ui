package kvitui

import (
	"strings"
	"time"
	"unicode"
)

// typeSelect is a list's jumping by typing (KEY-9): letters typed within a
// second of each other make up the start of a name, and the list moves to
// the first row whose name starts with it, ignoring case; one letter typed
// again and again moves on to the next row starting with it.
type typeSelect struct {
	typed string
	last  time.Time
}

// next is the row typing ch moves to, from the row at from, among n rows
// named by name; -1 when no row's name starts with what was typed.
func (s *typeSelect) next(ch rune, from, n int, name func(i int) string) int {
	if unicode.IsControl(ch) || n == 0 {
		return -1
	}
	now := time.Now()
	if now.Sub(s.last) > time.Second {
		s.typed = ""
	}
	s.last = now
	if s.typed == "" && ch == ' ' {
		return -1
	}
	s.typed += strings.ToLower(string(ch))
	if i := s.find(from, n, name); i >= 0 || len([]rune(s.typed)) == 1 {
		return i
	}
	// Nothing starts with what was typed: the letter starts a new name.
	s.typed = strings.ToLower(string(ch))
	return s.find(from, n, name)
}

// find is the row the typed letters mean, from the row at from.
func (s *typeSelect) find(from, n int, name func(i int) string) int {
	prefix, start := s.typed, from
	if rs := []rune(s.typed); len(rs) > 1 && strings.Count(s.typed, string(rs[0])) == len(rs) {
		// The same letter again: the next row starting with it.
		prefix, start = string(rs[0]), from+1
	} else if len([]rune(s.typed)) == 1 {
		start = from + 1
	}
	for k := range n {
		i := (max(0, start) + k) % n
		if strings.HasPrefix(strings.ToLower(name(i)), prefix) {
			return i
		}
	}
	return -1
}

// typing reports whether a rune is one a reader types to jump: printable,
// and not the Space that opens or chooses unless a name is being typed.
func (s *typeSelect) typing(ch rune) bool {
	if unicode.IsControl(ch) {
		return false
	}
	return ch != ' ' || (s.typed != "" && time.Since(s.last) <= time.Second)
}
