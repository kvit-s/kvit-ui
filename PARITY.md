# Parity with the Qt kvit-ui

This checklist tracks how far kvit-ui-go is from doing everything the Qt/QML
library in `~/kvit-ui` does. When every box is ticked, kvit-ui-go is ready
to replace it; `~/kvit-shirei/go-ui-plan.md` (sections 6 and 7) describes
what happens then.

Tick an item only with evidence, and write the evidence after it: the name of
a test, or a screenshot pair in which the Go version is shown beside the Qt
one.

Changes made to the Qt library after the Go port started have to be ported
too. `git -C ~/kvit-ui log qt-port-start..` lists them.

## Design values

- [ ] `Theme`: the four colour tables (light, dark, sepia, high contrast),
  switched while running, plus a `system` setting that follows the
  desktop's light or dark preference
- [ ] `Theme`: the reduced-motion setting and the motion scale
- [ ] `Interface`: the interface size from 10 to 24 px, and the seven type
  roles, spacing scale, row heights, control heights, radii, hairline,
  focus-ring width and three reflow breakpoints derived from it
- [ ] `Interface.px(n)` for one-off measurements
- [ ] `Typography`: document text (family, base size, line height,
  paragraph spacing, maximum content width), set separately from the chrome
- [ ] The theme and interface size saved between runs (`ui.json`)

## Components

Taken from `agent/skills/kvit-ui/catalog.md` in `~/kvit-ui` at commit 0a0b210 (74 components), in the groups `README.md` there uses.

### Foundation (4)

- [ ] `KvitLabel`
- [ ] `KvitIcon`
- [ ] `KvitIconButton`
- [ ] `KvitLink`

### Structure (8)

- [ ] `KvitHeader`
- [ ] `KvitSidebar`
- [ ] `KvitSidebarItem`
- [ ] `KvitBreadcrumb`
- [ ] `KvitRegion`
- [ ] `KvitViewHead`
- [ ] `KvitStatusBar`
- [ ] `KvitWindow`

### Content (9)

- [ ] `KvitSectionHeading`
- [ ] `KvitRow`
- [ ] `KvitSlimRow`
- [ ] `KvitCard`
- [ ] `KvitPanel`
- [ ] `KvitPane`
- [ ] `KvitDivider`
- [ ] `KvitDisclosure`
- [ ] `KvitEmptyState`

### Marks (7)

- [ ] `KvitChip`
- [ ] `KvitTag`
- [ ] `KvitBadge`
- [ ] `KvitSlug`
- [ ] `KvitDot`
- [ ] `KvitSignal`
- [ ] `KvitPip`

### Quantities (2)

- [ ] `KvitFigure`
- [ ] `KvitBeforeAfter`

### Controls (9)

- [ ] `KvitButton`
- [ ] `KvitChipButton`
- [ ] `KvitStepper`
- [ ] `KvitField`
- [ ] `KvitTextArea`
- [ ] `KvitSearchField`
- [ ] `KvitCheck`
- [ ] `KvitSelect`
- [ ] `KvitTab`

### Feedback (7)

- [ ] `KvitTooltip`
- [ ] `KvitPopover`
- [ ] `KvitHint`
- [ ] `KvitHoverCard`
- [ ] `KvitToast`
- [ ] `KvitNotice`
- [ ] `KvitDialog`

### Data (11)

- [ ] `KvitBar`
- [ ] `KvitStackedBar`
- [ ] `KvitSpark`
- [ ] `KvitTrend`
- [ ] `KvitDistribution`
- [ ] `KvitGauge`
- [ ] `KvitDelta`
- [ ] `KvitStatTile`
- [ ] `KvitFigureBlock`
- [ ] `KvitCell`
- [ ] `KvitTable`

### Flow (17)

- [ ] `KvitScrollBar`
- [ ] `KvitMenu`
- [ ] `KvitMenuItem`
- [ ] `KvitTree`
- [ ] `KvitSwitch`
- [ ] `KvitRadioGroup`
- [ ] `KvitProgress`
- [ ] `KvitSlider`
- [ ] `KvitSplitView`
- [ ] `KvitSegmented`
- [ ] `KvitTypeAhead`
- [ ] `KvitConfirmInPlace`
- [ ] `KvitTimeline`
- [ ] `KvitNumberField`
- [ ] `KvitMoneyField`
- [ ] `KvitDualList`
- [ ] `KvitSpotlight`

## Icons

- [ ] The 74 Phosphor icons, named by meaning; an unknown name draws a
  marked placeholder and fails a test

## Checks the Qt library runs, and their Go equivalents

- [ ] Colour tables and contrast floors, per theme
- [ ] The type scale
- [ ] The density and spacing scale
- [ ] The three chart ramps in all four themes, checked in OKLab against that
  theme's surfaces and every hue that already means something, under normal
  vision and three colour-vision deficiencies
- [ ] Every component opening in all four themes with no warning
- [ ] `KvitTable`: 250,000 rows and twelve columns, scrolling smoothly and
  filtering in under 100 ms
- [ ] Every gallery sample compiling
- [ ] No colour literal, numeric font size or unnamed geometry value outside
  the design values (in Qt: qmllint at full strength)
- [ ] Every interactive component giving screen readers a role and a name
- [ ] Generated files matching their sources: the drawing stylesheet, the
  icon catalogue and the agent skill's component catalogue

## Gallery

- [ ] One page per component, with its states and a working sample
- [ ] `--page`, `--theme`, `--interface-size`
- [ ] `--shots`: the fixed screenshot set (71 components × 4 themes, plus
  control-size variants, about 300 images), compared with the Qt set
- [ ] `--catalog`: writes the agent skill's component catalogue

## Drawing a screen before building it

- [ ] `ux/tokens.css` generated from the Go design values
- [ ] `ux/frame.css`, `ux/render.sh`, `visual-language.md` and
  `PATTERN.md` carried over

## Agent skills

- [ ] `kvit-ui`: the vocabulary and workflow, rewritten for Go
- [ ] `kvit-preview`: rendering, screenshot sets and single components,
  rewritten for Go

## Changes that exist only in kvit-cash's copy of kvit-ui

kvit-cash's kvit-ui checkout (`~/kvit-cash/third-party/kvit-ui`, pinned at
`722906e`) has 7 commits that `~/kvit-ui` lacks. The two split after `46a67c0`
(2026-09-10). kvit-ui-go has to do what these commits do as well.

- [ ] `97575ee` 2026-09-11: Say what opens something, and keep a popup inside the window
- [ ] `f4cb189` 2026-09-11: Raise a table's header to the height of a row somebody presses
- [ ] `1c237df` 2026-09-11: Say what a mark opens, and draw no control with nothing on it
- [ ] `364c3dc` 2026-09-11: Float a view over a list, open the rail on hover, draw a choice as chosen
- [ ] `d32c373` 2026-09-12: Give the wheel a distance, the card a fourth side, and a chart one baseline
- [ ] `526b619` 2026-09-13: One press acts on one thing, and what acts says so before it is pressed
- [ ] `722906e` 2026-09-13: Name the shortest window the chrome holds
