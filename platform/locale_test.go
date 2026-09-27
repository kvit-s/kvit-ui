package platform

import "testing"

func TestPosixLocaleNamesBecomeTags(t *testing.T) {
	for in, want := range map[string]string{
		"en_US.UTF-8":     "en-US",
		"de_DE@euro":      "de-DE",
		"en_US@rg=gbzzzz": "en-US",
		"fr":              "fr",
		"C":               "",
		"C.UTF-8":         "",
		"POSIX":           "",
	} {
		if got := posixLocale(in); got != want {
			t.Errorf("posixLocale(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTheEnvironmentIsReadInTheCLibrarysOrder(t *testing.T) {
	t.Setenv("LANG", "en_GB.UTF-8")
	t.Setenv("LC_NUMERIC", "de_DE.UTF-8")
	t.Setenv("LC_ALL", "")
	if got := envLocale(); got != "de-DE" {
		t.Errorf("LC_NUMERIC should win over LANG, got %q", got)
	}
	t.Setenv("LC_ALL", "C")
	if got := envLocale(); got != "" {
		t.Errorf("LC_ALL=C should group nothing, got %q", got)
	}
}
