// Command tokens-to-css writes ux/tokens.css: the design values as CSS
// variables, for drawing a screen as a static HTML page before it is built,
// so a drawing is judged in the colours and sizes the application will use.
// Its test fails when ux/tokens.css differs from what it writes.
//
//	go run ./tools/tokens-to-css > ux/tokens.css
package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"

	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/tokens"
)

func main() {
	css, err := render()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tokens-to-css:", err)
		os.Exit(1)
	}
	os.Stdout.WriteString(css)
}

// The themes, with the class a drawing puts on <body> for each.
// ".theme-contrast" rather than ".theme-highContrast": it is the name the
// existing drawings write.
var themes = []struct{ id, selector string }{
	{tokens.Dark, ":root,\n.theme-dark"},
	{tokens.Light, ".theme-light"},
	{tokens.Sepia, ".theme-sepia"},
	{tokens.HighContrast, ".theme-contrast"},
}

// The colours in the groups the stylesheet lists them under.
var groups = []struct {
	heading string
	members []string
}{
	{"surfaces", []string{"WindowBackground", "PanelBackground", "ListBackground", "FooterBackground",
		"PopupBackground", "ChipBackground", "BannerBackground", "CodePanelBackground"}},
	{"text", []string{"TextPrimary", "TextSecondary", "TextMuted", "TextFaint", "TextDisabled", "BannerText", "OnAccent"}},
	{"lines and glyphs", []string{"Border", "BorderStrong", "QuoteBar", "MutedGlyph"}},
	{"interactive tints", []string{"HoverTint", "BlockHoverTint", "FocusTint", "FocusRing", "SelectionTint",
		"SelectionActiveTint", "BlockSelectionTint"}},
	{"accent and semantic", []string{"Accent", "Danger", "DangerBright", "Success", "Warning", "PinColor", "Link",
		"Marker", "InlineCodeBackground", "HighlightBackground", "SearchMatchBackground", "SearchCurrentBackground",
		"ChangedTextBackground", "AddedTextBackground", "RemovedTextBackground", "CalloutTip"}},
	{"code, for a mockup that shows a document", []string{"CodeKeyword", "CodeType", "CodeString", "CodeComment", "CodeNumber"}},
	{"portfolio vocabulary", []string{"AxisAttention", "AxisAttentionText", "AxisAgent", "AxisAgentText",
		"ScopeDiscovered", "SignalHard", "SignalSoft", "SignalHygiene", "HatchAlt"}},
}

// cssName is a field's CSS variable: WindowBackground is --window-background.
func cssName(field string) string {
	var b strings.Builder
	b.WriteString("--")
	for i, r := range field {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('-')
			}
			r = unicode.ToLower(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}

func themeBlock(b *strings.Builder, id, selector string) error {
	th := tokens.NewTheme()
	th.SetThemeID(id)
	t := reflect.ValueOf(th.Tokens())
	fmt.Fprintf(b, "%s {\n", selector)
	grouped := map[string]bool{}
	for i, g := range groups {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(b, "  /* %s */\n", g.heading)
		for _, name := range g.members {
			f := t.FieldByName(name)
			if !f.IsValid() {
				return fmt.Errorf("the token table has no %s", name)
			}
			fmt.Fprintf(b, "  %s: %s;\n", cssName(name), f.Interface().(palette.Color).Hex())
			grouped[name] = true
		}
	}
	// A colour in no group would be missing from every drawing, which would
	// still look right because the browser falls back to inheriting.
	for i := range t.NumField() {
		f := t.Type().Field(i)
		if f.Type == reflect.TypeOf(palette.Color{}) && !grouped[f.Name] {
			return fmt.Errorf("the colour %s belongs to no group in tools/tokens-to-css", f.Name)
		}
	}
	tk := th.Tokens()
	for _, ramp := range []struct {
		name  string
		steps []palette.Color
	}{{"categorical", tk.CategoricalRamp}, {"sequential", tk.SequentialRamp}, {"diverging", tk.DivergingRamp}} {
		fmt.Fprintf(b, "\n  /* %s ramp */\n", ramp.name)
		for i, c := range ramp.steps {
			fmt.Fprintf(b, "  --%s-%d: %s;\n", ramp.name, i, c.Hex())
		}
	}
	b.WriteString("}\n")
	return nil
}

type entry struct {
	name  string
	value int
	note  string
}

func scaleBlock(b *strings.Builder) {
	m := tokens.NewInterface()
	b.WriteString(":root {\n")
	b.WriteString("  /* Local families only — no @font-face, no network, so a mockup\n" +
		"   * renders identically here and on any machine. */\n")
	b.WriteString("  --font-ui: Ubuntu, \"Ubuntu Sans\", \"DejaVu Sans\", system-ui, sans-serif;\n")
	b.WriteString("  --font-mono: \"DejaVu Sans Mono\", \"Ubuntu Mono\", ui-monospace, monospace;\n")
	fmt.Fprintf(b, "\n  /* The seven type roles, in pixels at the default interface\n   * size of %d. */\n", m.FontSize())
	for _, e := range []entry{
		{"type-caption", m.Caption(), "kind tags, counts"},
		{"type-small", m.Small(), "chip labels, sub-lines"},
		{"type-body", m.Body(), "row text, prose"},
		{"type-strong", m.Strong(), "a name, an emphasised row"},
		{"type-title", m.Title(), "a section heading"},
		{"type-headline", m.Headline(), "a pane title, the wordmark"},
		{"type-display", m.Display(), "a page title"},
	} {
		pad := max(1, 22-len(e.name)-len(fmt.Sprint(e.value)))
		fmt.Fprintf(b, "  --%s: %dpx;%s/* %s */\n", e.name, e.value, strings.Repeat(" ", pad), e.note)
	}
	list := func(heading string, entries []entry) {
		b.WriteString(heading)
		for _, e := range entries {
			fmt.Fprintf(b, "  --%s: %dpx;", e.name, e.value)
			if e.note != "" {
				fmt.Fprintf(b, "  /* %s */", e.note)
			}
			b.WriteString("\n")
		}
	}
	list("\n  /* The spacing scale. 3 and 5 are deliberately not steps. */\n", []entry{
		{"space-tight", m.SpaceTight(), ""}, {"space-snug", m.SpaceSnug(), ""}, {"space-near", m.SpaceNear(), ""},
		{"space", m.Space(), "the ordinary gap"}, {"space-wide", m.SpaceWide(), ""}, {"space-loose", m.SpaceLoose(), ""},
	})
	// kvit-hub's older type names on the merged scale, which the drawings
	// made before the merge still write.
	list("\n  /* kvit-hub's older type names, on the merged scale. New work\n"+
		"   * uses the seven roles above; these keep the existing mockups\n"+
		"   * rendering until Wave 3 moves them over. */\n", []entry{
		{"type-micro", m.Caption(), "was baseSize - 6"}, {"type-small", m.Small(), "was baseSize - 5"},
		{"type-secondary", m.Small(), "was baseSize - 4"}, {"type-row", m.Body(), "was baseSize - 3"},
		{"type-name", m.Strong(), "was baseSize"}, {"type-heading", m.Title(), "was baseSize + 1"},
		{"type-page", m.Headline(), "was baseSize + 2"},
	})
	list("\n  /* Density. One reader, desktop, no mobile target. */\n", []entry{
		{"view-margin", m.ViewMargin(), "the outer margin of a view"},
		{"column-gap", m.ColumnGap(), ""},
		{"stack-gap", m.StackGap(), "between stacked blocks"},
		{"sidebar-width", m.SidebarWidth(), ""},
		{"rail-width", m.RailWidth(), "the sidebar collapsed"},
		{"pane-width", m.PaneWidth(), "the side pane"},
		{"floating-view-width", m.FloatingViewWidth(), "a view held over a list"},
		{"header-height", m.HeaderHeight(), ""},
		{"breadcrumb-height", m.BreadcrumbHeight(), ""},
		{"status-bar-height", m.StatusBarHeight(), ""},
		{"row-height", m.RowHeight(), "name over description"},
		{"row-height-sub", m.RowHeightSub(), "an expanded line"},
		{"row-height-slim", m.RowHeightSlim(), "a list at rest"},
		{"row-height-compact", m.RowHeightCompact(), "a disclosure"},
		{"control-height", m.ControlHeight(), ""},
		{"tab-height", m.TabHeight(), ""},
		{"chip-height", m.ChipHeight(), ""},
		{"tag-height", m.TagHeight(), ""},
		{"pill-height", m.PillHeight(), ""},
		{"bar-height", m.BarHeight(), "a compact bar"},
		{"bar-height-wide", m.BarHeightWide(), "a full-size bar"},
		{"icon-size", m.IconSize(), ""},
		{"icon-size-small", m.IconSizeSmall(), ""},
		{"radius-bar", m.RadiusBar(), ""},
		{"radius-chip", m.RadiusChip(), ""},
		{"radius-control", m.RadiusControl(), ""},
		{"radius-card", m.RadiusCard(), ""},
		{"radius-pill", m.RadiusPill(), ""},
		{"hairline", m.Hairline(), ""},
		{"focus-ring-width", m.FocusRingWidth(), ""},
		{"width-floor", m.WidthFloor(), "below this, effort dots collide"},
		{"width-laptop", m.WidthLaptop(), "the sidebar collapses"},
		{"width-drawn", m.WidthDrawn(), "every drawing is this wide"},
	})
	b.WriteString("}\n")
}

func render() (string, error) {
	var b strings.Builder
	b.WriteString(`/* kvit-ui design tokens, for the drawing layer.
 *
 * GENERATED FILE — do not edit.
 *
 * Written by tools/tokens-to-css from the same design values the
 * applications draw with (the tokens package: tokens.Tokens and
 * tokens.Interface), so a drawing is judged in the colours the
 * application will actually use. Editing this file by hand fails
 * the generator's test.
 *
 * Values only. The window, the base elements and the component
 * classes are in frame.css, which is hand-written.
 *
 * Themes: :root carries dark. Put .theme-light, .theme-dark,
 * .theme-sepia or .theme-contrast on <body> to render a mockup in
 * another one.
 */

`)
	b.WriteString("/* ══════════════════════════════════════════════ theme tokens ══════ */\n\n")
	for _, th := range themes {
		if err := themeBlock(&b, th.id, th.selector); err != nil {
			return "", err
		}
		b.WriteString("\n")
	}
	b.WriteString("/* ═══════════════════════════════════ type, density, geometry ══════ */\n\n")
	scaleBlock(&b)
	return b.String(), nil
}
