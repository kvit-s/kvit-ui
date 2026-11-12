package kvitui

import (
	"testing"

	"github.com/kvit-s/kvit-ui/text"
)

// The style a label's fields ask for is the style its text is laid out in.
func TestALabelsFieldsReachItsStyle(t *testing.T) {
	ui, err := New(Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	l := NewLabel(ui, "run 12 of 40")
	if st := l.style(); st.Italic || st.Weight != text.Regular || st.Size != float32(ui.Size(RoleBody)) {
		t.Fatalf("a plain label's style is %+v", st)
	}
	l.Italic = true
	if !l.style().Italic {
		t.Error("an italic label's chrome style is not italic")
	}
	l.Mono, l.Weight = true, text.Semibold
	if st := l.style(); !st.Italic || st.Weight == text.Semibold {
		t.Errorf("an italic monospace label with a weight and no WeightedMono has style %+v", st)
	}
	l.WeightedMono = true
	if st := l.style(); st.Weight != text.Semibold || st.Family != ui.Interface.MonoFamily() {
		t.Errorf("a weighted monospace label has style %+v", st)
	}
	l.Weight = 0
	if st := l.style(); st.Weight != text.Regular {
		t.Errorf("a weighted monospace label with no weight has weight %v, want regular", st.Weight)
	}
	// A size of its own replaces the role's, and follows the interface size.
	l.Size = Px(30)
	if st := l.style(); st.Size != float32(ui.Interface.Px(30)) {
		t.Errorf("a label with a size of 30 design pixels is set at %v", st.Size)
	}
	ui.Interface.SetFontSize(24)
	if st := l.style(); st.Size != float32(ui.Interface.Px(30)) || st.Size == 30 {
		t.Errorf("at the largest interface size the label is set at %v", st.Size)
	}
}
