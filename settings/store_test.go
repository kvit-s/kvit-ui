package settings_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/kvit-s/kvit-ui/settings"
)

func open(t *testing.T, path string) *settings.Store {
	t.Helper()
	s := settings.New()
	if err := s.Open(path, false); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestUnknownKeysSurviveARoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"someone.else": [1, 2, 3], "theme.id": "dark"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	s := open(t, path)
	s.SetValue("theme.id", "sepia")
	if err := s.Flush(); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	data, _ := os.ReadFile(path)
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got["theme.id"] != "sepia" || len(got["someone.else"].([]any)) != 3 {
		t.Errorf("the file after a write: %s", data)
	}
}

func TestACorruptFileGivesDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{not json`), 0o644); err != nil {
		t.Fatal(err)
	}
	s := open(t, path)
	if s.Contains("theme.id") {
		t.Error("a corrupt file produced values")
	}
	s.SetValue("theme.id", "dark")
	if err := s.Flush(); err != nil {
		t.Fatal(err)
	}
	if v, _ := open(t, path).Value("theme.id"); v != "dark" {
		t.Error("the file was not rewritten after a corrupt load")
	}
}

func TestWritesAreDelayedAndCoalesced(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	s := open(t, path)
	for i := 10; i <= 24; i++ {
		s.SetValue("interface.fontSize", i)
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("the file was written before the delay")
	}
	deadline := time.Now().Add(5 * settings.WriteDelay)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(path); err == nil {
			var got map[string]any
			if json.Unmarshal(data, &got) == nil && got["interface.fontSize"] == 24.0 {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Error("the last value never reached the file")
}

func TestUnchangedValuesDoNotNotify(t *testing.T) {
	s := open(t, filepath.Join(t.TempDir(), "settings.json"))
	n := 0
	s.OnValueChanged(func(string) { n++ })
	s.SetValue("theme.id", "dark")
	s.SetValue("theme.id", "dark")
	s.Remove("missing")
	s.Remove("theme.id")
	if n != 2 {
		t.Errorf("%d notifications, want 2", n)
	}
}

func TestAFailedWriteStaysPending(t *testing.T) {
	if runtime.GOOS == "windows" || os.Getuid() == 0 {
		t.Skip("needs a directory the test user cannot write")
	}
	dir := filepath.Join(t.TempDir(), "locked")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	s := open(t, filepath.Join(dir, "settings.json"))
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0o755)
	s.SetValue("theme.id", "dark")
	if err := s.Flush(); err == nil {
		t.Fatal("a write into a read-only directory succeeded")
	}
	// Opening another file must refuse while the only copy is pending.
	if err := s.Open(filepath.Join(t.TempDir(), "other.json"), false); err == nil {
		t.Error("Open dropped a pending write")
	}
	if v, _ := s.Value("theme.id"); v != "dark" {
		t.Error("the pending value was lost")
	}
	os.Chmod(dir, 0o755)
	if err := s.Flush(); err != nil {
		t.Errorf("the retry failed: %v", err)
	}
}
