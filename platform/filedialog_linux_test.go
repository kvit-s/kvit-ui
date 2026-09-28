package platform

import (
	"os/exec"
	"slices"
	"testing"
	"time"

	"github.com/richardwilkes/unison"
)

// While zenity or kdialog runs, the application goes on running its event
// loop behind a modal window of its own, which ends when the program exits:
// with the files the program wrote, or none when it was cancelled.
func TestAnExternalDialogRunsModallyAndReportsWhatWasChosen(t *testing.T) {
	var owner *unison.Window
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 800, Height: 600},
		unison.StartupFinishedCallback(func() {
			var err error
			if owner, err = unison.NewWindow("Owner"); err == nil {
				owner.SetFrameRect(owner.FrameRectForContentRect(unison.PrimaryDisplay().Usable))
				owner.ToFront()
			}
		}))
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.Sync()
	var ticks, during int
	var chosen, cancelled []string
	screen.Do(func() {
		// Work handed to the event loop runs while the program does.
		unison.InvokeTaskAfter(func() { ticks++ }, 50*time.Millisecond)
		chosen = runExternalDialog(exec.Command("sh", "-c", "sleep 0.2; printf '/w/Household.sqlite\\n'"), false)
		during = ticks
		cancelled = runExternalDialog(exec.Command("sh", "-c", "exit 1"), false)
	})
	if !slices.Equal(chosen, []string{"/w/Household.sqlite"}) || cancelled != nil {
		t.Errorf("chosen %q, cancelled %q", chosen, cancelled)
	}
	if during != 1 {
		t.Errorf("the event loop ran %d tasks while the program ran", during)
	}
	var left []string
	screen.Do(func() {
		for _, w := range unison.Windows() {
			left = append(left, w.Title())
		}
	})
	if !slices.Equal(left, []string{"Owner"}) {
		t.Errorf("the windows left are %q", left)
	}
}
