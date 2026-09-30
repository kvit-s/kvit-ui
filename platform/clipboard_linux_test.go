package platform

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"image"
	"image/png"
	"testing"
	"unicode/utf16"

	"github.com/richardwilkes/unison"
)

// fakeWindows makes the program run under WSL or not, with PowerShell
// answering out, and counts the times PowerShell is asked.
func fakeWindows(t *testing.T, wsl bool, out string) *int {
	t.Helper()
	asked := 0
	savedWSL, savedRun := underWSL, runPowerShell
	t.Cleanup(func() { underWSL, runPowerShell = savedWSL, savedRun })
	underWSL = func() bool { return wsl }
	runPowerShell = func(script string) ([]byte, error) {
		asked++
		if script != windowsPictureScript {
			t.Errorf("PowerShell is given %q", script)
		}
		return []byte(out), nil
	}
	return &asked
}

func samplePNG(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// Under WSL the picture PowerShell writes, base64 with a line end after it,
// is the PNG answered.
func TestAPictureOnWindowsClipboardIsReadUnderWSL(t *testing.T) {
	want := samplePNG(t)
	asked := fakeWindows(t, true, base64.StdEncoding.EncodeToString(want)+"\r\n")
	got, ok := WindowsClipboardPicture()
	if !ok || !bytes.Equal(got, want) {
		t.Fatalf("got %d bytes, %v; want the %d-byte PNG", len(got), ok, len(want))
	}
	if *asked != 1 {
		t.Errorf("PowerShell was asked %d times", *asked)
	}
}

// Outside WSL, Windows is never asked.
func TestWindowsIsNotAskedOutsideWSL(t *testing.T) {
	asked := fakeWindows(t, false, base64.StdEncoding.EncodeToString(samplePNG(t)))
	if _, ok := WindowsClipboardPicture(); ok {
		t.Error("a picture was answered outside WSL")
	}
	if *asked != 0 {
		t.Errorf("PowerShell was asked %d times", *asked)
	}
}

// A headless session has a clipboard of its own, so Windows' is not asked
// even under WSL.
func TestWindowsIsNotAskedFromAHeadlessSession(t *testing.T) {
	asked := fakeWindows(t, true, base64.StdEncoding.EncodeToString(samplePNG(t)))
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 200, Height: 100})
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.Sync()
	var ok bool
	screen.Do(func() { _, ok = WindowsClipboardPicture() })
	if ok || *asked != 0 {
		t.Errorf("answered a picture: %v, PowerShell asked %d times", ok, *asked)
	}
}

// No output, output that is not base64, and base64 of something other than
// a PNG are all no picture.
func TestOnlyAPNGIsAPicture(t *testing.T) {
	for _, out := range []string{"", "\r\n", "not base64!", base64.StdEncoding.EncodeToString([]byte("text"))} {
		fakeWindows(t, true, out)
		if _, ok := WindowsClipboardPicture(); ok {
			t.Errorf("output %q was taken as a picture", out)
		}
	}
}

// The script is handed over as base64 of its UTF-16 little-endian text,
// characters outside the basic plane included.
func TestTheScriptGoesAsUTF16InBase64(t *testing.T) {
	script := "Write-Output 'añ€𝄞'"
	raw, err := base64.StdEncoding.DecodeString(encodePowerShell(script))
	if err != nil || len(raw)%2 != 0 {
		t.Fatalf("not base64 of UTF-16: %v", err)
	}
	units := make([]uint16, len(raw)/2)
	for i := range units {
		units[i] = binary.LittleEndian.Uint16(raw[2*i:])
	}
	if got := string(utf16.Decode(units)); got != script {
		t.Errorf("decodes to %q", got)
	}
}
