package platform

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/unison"
)

// windowsPictureScript writes the picture on Windows' clipboard as PNG in
// base64, and nothing when there is none. Base64 keeps the bytes clear of the
// text encoding PowerShell gives its output.
const windowsPictureScript = `Add-Type -AssemblyName System.Windows.Forms, System.Drawing
$picture = [System.Windows.Forms.Clipboard]::GetImage()
if ($picture -ne $null) {
    $stream = New-Object System.IO.MemoryStream
    $picture.Save($stream, [System.Drawing.Imaging.ImageFormat]::Png)
    [Console]::Out.Write([Convert]::ToBase64String($stream.ToArray()))
}`

// windowsClipboardTimeout bounds the wait for PowerShell, which answers in
// under a second.
const windowsClipboardTimeout = 10 * time.Second

// windowsPowerShell is Windows PowerShell where WSL mounts the C: drive,
// for when PATH does not include Windows' own folders.
const windowsPowerShell = "/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe"

var pngSignature = []byte("\x89PNG\r\n\x1a\n")

// underWSL and runPowerShell are replaced by the tests.
var (
	underWSL      = isWSL
	runPowerShell = powerShellOutput
)

// windowsClipboardPicture skips a headless session, whose clipboard is its
// own in memory, so a test never reads the desktop's clipboard.
func windowsClipboardPicture() ([]byte, bool) {
	if unison.ActiveHeadlessScreen() != nil || !underWSL() {
		return nil, false
	}
	out, err := runPowerShell(windowsPictureScript)
	if err != nil {
		errs.Log(err)
		return nil, false
	}
	return decodeWindowsPicture(out)
}

// decodeWindowsPicture turns windowsPictureScript's output into the PNG it
// holds, false for no output or for anything that is not a PNG.
func decodeWindowsPicture(out []byte) ([]byte, bool) {
	text := strings.TrimSpace(string(out))
	if text == "" {
		return nil, false
	}
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil || !bytes.HasPrefix(data, pngSignature) {
		return nil, false
	}
	return data, true
}

// powerShellOutput runs a script in Windows PowerShell and answers what it
// wrote; no output when there is no PowerShell. The script goes as
// -EncodedCommand so that no quoting has to survive WSL's passing of
// arguments to a Windows program. Windows PowerShell rather than PowerShell 7
// because every Windows has it, and -STA because the clipboard is read from a
// single-threaded apartment.
func powerShellOutput(script string) ([]byte, error) {
	program, err := exec.LookPath("powershell.exe")
	if err != nil {
		if _, err = os.Stat(windowsPowerShell); err != nil {
			return nil, nil
		}
		program = windowsPowerShell
	}
	ctx, cancel := context.WithTimeout(context.Background(), windowsClipboardTimeout)
	defer cancel()
	return exec.CommandContext(ctx, program, "-NoProfile", "-NonInteractive", "-STA",
		"-EncodedCommand", encodePowerShell(script)).Output()
}

// encodePowerShell is a script as -EncodedCommand takes it: its UTF-16
// little-endian text in base64.
func encodePowerShell(script string) string {
	units := utf16.Encode([]rune(script))
	raw := make([]byte, 2*len(units))
	for i, unit := range units {
		binary.LittleEndian.PutUint16(raw[2*i:], unit)
	}
	return base64.StdEncoding.EncodeToString(raw)
}
