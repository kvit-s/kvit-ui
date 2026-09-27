package platform

import (
	"os"
	"strings"
)

// Locale is the language and region the desktop writes numbers for, as a
// BCP 47 tag such as "en-US" or "de-DE". It is "" for the C and POSIX
// locales, which group no digits, and where the desktop gives no answer.
func Locale() string { return readLocale() }

// posixLocale turns a POSIX locale name, such as "de_DE.UTF-8@euro", into a
// BCP 47 tag.
func posixLocale(s string) string {
	if i := strings.IndexAny(s, ".@"); i >= 0 {
		s = s[:i]
	}
	if s == "C" || s == "POSIX" {
		return ""
	}
	return strings.ReplaceAll(s, "_", "-")
}

// envLocale reads the locale variables in the order the C library consults
// them for numbers.
func envLocale() string {
	for _, k := range []string{"LC_ALL", "LC_NUMERIC", "LANG"} {
		if v := os.Getenv(k); v != "" {
			return posixLocale(v)
		}
	}
	return ""
}
