package kvitui

import (
	"strings"
	"unicode"
)

// An app writes every label once, in sentence case, which is the Windows
// form for all text, and kvit-ui turns the labels the macOS form writes in
// title case into title case there: menu titles and lines, buttons, tabs,
// segments and column headings (WRD-6, decided Q-3). Tooltips, field labels
// and messages stay in sentence case on both forms.

// TitleCaseKept lists words title case leaves as they are written, by their
// lower-case spelling: a name written lower case on purpose. An app adds its
// own; a word already holding a capital, such as iCloud, macOS or PDF, is
// kept as written without being listed.
var TitleCaseKept = map[string]bool{
	"kvit-ui": true,
}

// titleSmallWords are the words Apple's title case leaves lower case unless
// they come first or last: articles, coordinating conjunctions, and
// prepositions of three letters or fewer.
var titleSmallWords = map[string]bool{
	"a": true, "an": true, "the": true,
	"and": true, "but": true, "or": true, "nor": true, "for": true, "yet": true, "so": true,
	"as": true, "at": true, "by": true, "in": true, "of": true, "off": true, "on": true,
	"per": true, "to": true, "up": true, "via": true, "out": true,
}

// TitleCase is a sentence-case label in Apple's title case: every word
// capitalised except articles, coordinating conjunctions and prepositions
// of three letters or fewer, which stay lower case unless first or last.
// "Move to trash…" is "Move to Trash…", "Set up" is "Set Up". A word with a
// capital, a digit or a dot already in it is kept as written, and so is one
// in TitleCaseKept. An access-key marker ("Save &as") stays where it is.
func TitleCase(s string) string {
	words := strings.Split(s, " ")
	last := len(words) - 1
	for last > 0 && bareWord(words[last]) == "" {
		last--
	}
	first := 0
	for first < last && bareWord(words[first]) == "" {
		first++
	}
	for i, w := range words {
		bare := bareWord(w)
		if bare == "" || keptAsWritten(w, bare) {
			continue
		}
		if i != first && i != last && titleSmallWords[bare] {
			continue
		}
		words[i] = capitaliseParts(w)
	}
	return strings.Join(words, " ")
}

// bareWord is a word's letters in lower case, without the access-key
// marker and the punctuation around it: "&Save…" is "save".
func bareWord(w string) string {
	w = strings.ReplaceAll(w, "&", "")
	w = strings.TrimFunc(w, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	return strings.ToLower(w)
}

// keptAsWritten reports whether a word is left as the app wrote it: one in
// TitleCaseKept, or one with a capital after its first letter, a digit, a
// dot, a slash or an at sign in it.
func keptAsWritten(w, bare string) bool {
	if TitleCaseKept[bare] {
		return true
	}
	firstLetter := true
	for _, r := range strings.ReplaceAll(w, "&", "") {
		if unicode.IsDigit(r) || r == '.' || r == '/' || r == '@' {
			return true
		}
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) && !firstLetter {
				return true
			}
			firstLetter = false
		}
	}
	return false
}

// capitaliseParts capitalises the first letter of a word and of each part
// after a hyphen: "built-in" is "Built-In".
func capitaliseParts(w string) string {
	rs := []rune(w)
	start := true
	for i, r := range rs {
		switch {
		case r == '-':
			start = true
		case start && unicode.IsLetter(r):
			rs[i] = unicode.ToUpper(r)
			start = false
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			start = false
		}
	}
	return string(rs)
}

// lowerFirst writes a sentence-case label's first word in lower case, for a
// name built around it ("Decrease interface size"), unless the word keeps its
// capital: a name such as Kvit or an acronym such as PDF.
func lowerFirst(s string) string {
	first, _, _ := strings.Cut(s, " ")
	rs := []rune(first)
	if len(rs) < 2 || !unicode.IsUpper(rs[0]) || TitleCaseKept[strings.ToLower(first)] {
		return s
	}
	for _, r := range rs[1:] {
		if unicode.IsUpper(r) || unicode.IsDigit(r) {
			return s
		}
	}
	return string(unicode.ToLower(rs[0])) + string([]rune(s)[1:])
}

// CommandCase is a label of the kind the macOS form writes in title case,
// a menu line, a button, a tab, a segment or a column heading, in the app's
// form: title case on the macOS form, as written on the Windows form.
func (u *UI) CommandCase(s string) string {
	if u.form == FormMac {
		return TitleCase(s)
	}
	return s
}
