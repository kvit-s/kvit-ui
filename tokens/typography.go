package tokens

import "math"

// FontRole is the entry of the document's type scale a run of text is set at.
type FontRole int

// The document's font roles.
const (
	Body FontRole = iota
	Heading1
	Heading2
	Heading3
	Heading4
	Mono // a code fence, and anything drawn in the monospace family at the code size
)

// Typography's clamps and defaults.
const (
	MinBaseSize             = 10
	MaxBaseSize             = 28
	DefaultBaseSize         = 14
	MinLineHeight           = 1.0
	MaxLineHeight           = 2.0
	DefaultLineHeight       = 1.3
	MinParagraphSpacing     = 0
	MaxParagraphSpacing     = 40
	DefaultParagraphSpacing = 4
	MinContentWidth         = 300 // when not 0
)

const (
	keyDocFamily        = "typography.fontFamily"
	keyBaseSize         = "typography.fontSize"
	keyLineHeight       = "typography.lineHeight"
	keyParagraphSpacing = "typography.paragraphSpacing"
	keyMaxContentWidth  = "typography.maxContentWidth"
	keyDocMonoFamily    = "typography.monoFamily"
)

// Typography is the document's text: family, base size, line height,
// paragraph spacing, maximum content width and the code family. Heading and
// code sizes derive from the base by fixed ratios, so one setting scales the
// whole document. It is separate from Interface because a reader who wants
// large body text in a dense list is asking for something coherent.
type Typography struct {
	settings         Settings
	loading          bool
	fontFamily       string
	baseSize         int
	lineHeight       float64
	paragraphSpacing int
	maxContentWidth  int
	monoFamily       string
	changed          signal
}

// NewTypography returns the defaults: 14 px, line height 1.3, 4 px between
// blocks, no width cap, the desktop's families.
func NewTypography() *Typography {
	return &Typography{
		baseSize:         DefaultBaseSize,
		lineHeight:       DefaultLineHeight,
		paragraphSpacing: DefaultParagraphSpacing,
		monoFamily:       "monospace",
	}
}

// OnChanged calls fn whenever any value may have changed.
func (t *Typography) OnChanged(fn func()) (disconnect func()) { return t.changed.connect(fn) }

// SetSettings reads the stored values, clamped as a live change would be, and
// saves every later change back.
func (t *Typography) SetSettings(s Settings) {
	t.settings = s
	if s == nil {
		return
	}
	t.loading = true
	t.SetFontFamily(stringValue(s, keyDocFamily, t.fontFamily))
	t.SetBaseSize(int(math.Round(numberValue(s, keyBaseSize, float64(t.baseSize)))))
	t.SetLineHeight(numberValue(s, keyLineHeight, t.lineHeight))
	t.SetParagraphSpacing(int(math.Round(numberValue(s, keyParagraphSpacing, float64(t.paragraphSpacing)))))
	t.SetMaxContentWidth(int(math.Round(numberValue(s, keyMaxContentWidth, float64(t.maxContentWidth)))))
	t.SetMonoFamily(stringValue(s, keyDocMonoFamily, t.monoFamily))
	t.loading = false
	t.changed.emit()
}

func (t *Typography) FontFamily() string    { return t.fontFamily }
func (t *Typography) BaseSize() int         { return t.baseSize }
func (t *Typography) LineHeight() float64   { return t.lineHeight }
func (t *Typography) ParagraphSpacing() int { return t.paragraphSpacing }
func (t *Typography) MaxContentWidth() int  { return t.maxContentWidth }
func (t *Typography) MonoFamily() string    { return t.monoFamily }

// SetFontFamily chooses the document family; "" is the application default.
func (t *Typography) SetFontFamily(family string) {
	if t.fontFamily == family {
		return
	}
	t.fontFamily = family
	t.save(keyDocFamily, family)
	t.changed.emit()
}

// SetBaseSize sets the body size, clamped to 10–28.
func (t *Typography) SetBaseSize(size int) {
	size = min(max(size, MinBaseSize), MaxBaseSize)
	if t.baseSize == size {
		return
	}
	t.baseSize = size
	t.save(keyBaseSize, size)
	t.changed.emit()
}

// SetLineHeight sets the line height as a multiple of the font's own, 1–2.
func (t *Typography) SetLineHeight(height float64) {
	height = min(max(height, MinLineHeight), MaxLineHeight)
	if math.Abs(t.lineHeight-height) < 1e-9 {
		return
	}
	t.lineHeight = height
	t.save(keyLineHeight, height)
	t.changed.emit()
}

// SetParagraphSpacing sets the pixels between blocks, 0–40.
func (t *Typography) SetParagraphSpacing(spacing int) {
	spacing = min(max(spacing, MinParagraphSpacing), MaxParagraphSpacing)
	if t.paragraphSpacing == spacing {
		return
	}
	t.paragraphSpacing = spacing
	t.save(keyParagraphSpacing, spacing)
	t.changed.emit()
}

// SetMaxContentWidth caps the block column at this many pixels, centred; 0
// fills the pane, and any other value is at least 300.
func (t *Typography) SetMaxContentWidth(width int) {
	if width != 0 {
		width = max(MinContentWidth, width)
	}
	if t.maxContentWidth == width {
		return
	}
	t.maxContentWidth = width
	t.save(keyMaxContentWidth, width)
	t.changed.emit()
}

// SetMonoFamily chooses the family code is set in.
func (t *Typography) SetMonoFamily(family string) {
	if t.monoFamily == family {
		return
	}
	t.monoFamily = family
	t.save(keyDocMonoFamily, family)
	t.changed.emit()
}

// SizeForRole is the pixel size a role is set at. The ratios are written
// against a 15 px base because that is what they were measured at: at the
// default of 14 the headings are 30, 22, 19 and 16 with code at 12, and at 15
// they are the original 32, 24, 20, 17 and 13.
func (t *Typography) SizeForRole(role FontRole) int {
	numerator := 15.0
	switch role {
	case Heading1:
		numerator = 32
	case Heading2:
		numerator = 24
	case Heading3:
		numerator = 20
	case Heading4:
		numerator = 17
	case Mono:
		numerator = 13
	}
	return int(math.Round(float64(t.baseSize) * numerator / 15))
}

// BodySize is the size ordinary document text is set at.
func (t *Typography) BodySize() int { return t.SizeForRole(Body) }

// MonoSize is the size code is set at.
func (t *Typography) MonoSize() int { return t.SizeForRole(Mono) }

// ResetToDefaults returns every value to the defaults.
func (t *Typography) ResetToDefaults() {
	t.SetFontFamily("")
	t.SetBaseSize(DefaultBaseSize)
	t.SetLineHeight(DefaultLineHeight)
	t.SetParagraphSpacing(DefaultParagraphSpacing)
	t.SetMaxContentWidth(0)
	t.SetMonoFamily("monospace")
}

func (t *Typography) save(key string, value any) {
	if t.settings != nil && !t.loading {
		t.settings.SetValue(key, value)
	}
}
