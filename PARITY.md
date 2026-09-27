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

- [x] `Theme`: the four colour tables (light, dark, sepia, high contrast),
  switched while running, plus a `system` setting that follows the
  desktop's light or dark preference. Evidence: `tokens/theme_test.go`
  (23 tests ported from `test_theme.cpp`); tables generated from the Qt
  source by `tools/import-qt-theme`.
- [x] `Theme`: the reduced-motion setting and the motion scale. Evidence:
  `TestReducedMotionScale`, `TestReducedMotionFollowsTheSystem`,
  `TestAnExistingMotionChoiceSurvivesTheUpgrade`.
- [x] `Interface`: the interface size from 10 to 24 px, and the seven type
  roles, spacing scale, row heights, control heights, radii, hairline,
  focus-ring width and three reflow breakpoints derived from it. Evidence:
  `tokens/interface_test.go` (ported from `test_interfacemetrics.cpp` and
  `test_density.cpp`).
- [x] `Interface.px(n)` for one-off measurements (`Interface.Px`).
  Evidence: `TestDefaultReproducesTheOldLiterals`,
  `TestNoDensityValueCollapsesAtTheSmallestSize`.
- [x] `Typography`: document text (family, base size, line height,
  paragraph spacing, maximum content width), set separately from the
  chrome. Evidence: `TestTypographyScale`, `TestTypographyClampsAndSignals`,
  `TestTheDocumentTypeScaleIsUntouched`.
- [x] The theme and interface size saved between runs (`ui.json`), with the
  Qt library's keys and file format, so both read one file. Evidence:
  `settings/store_test.go`, `TestChoicesPersistAcrossRuns`.

## Components

Taken from `agent/skills/kvit-ui/catalog.md` in `~/kvit-ui` at commit 0a0b210 (74 components), in the groups `README.md` there uses.

The screenshot evidence below is `./build.sh --shots`, which writes each Go
page under the Qt file name and stacks it above the Qt shot in
`build/shots/compare`. As in the Qt run, the window is as wide as the design
width at each interface size (1200, 1440 and 2880 px at 10, 12 and 24 px),
and numbers are written under the C locale, with no digit grouping. "Matches" means the pair was compared by eye region by
region, and positions agree to the pixel unless a difference is listed.

### Foundation (4)

Each of the four themes was compared on at least one of these pages, at
12 px, the only size the Qt set has. Two differences are
expected and not defects: each code sample is the Go source of the specimen
rather than QML, and the Go gallery has no filter field above its sidebar
yet.

- [x] `KvitLabel` (`Label`). Evidence: `TestLabelSizesByRoleAndElides`,
  `TestTheLastGlyphOfALineKeepsItsWidth`, `TestPagesFollowTheInterfaceSize`;
  `dark-KvitLabel.png` matches, every role line and the frame to the pixel.
- [x] `KvitIcon` (`Icon`). Evidence: `TestAnUnknownIconIsMarked`;
  `light-KvitIcon.png` matches except that the small search symbol is centred
  on the row rather than at its top, because Go's `Row` centres what it
  holds and Qt's `Row` does not.
- [x] `KvitIconButton` (`IconButton`). Evidence:
  `TestIconButtonActivatesAndDescribesItself`,
  `TestTheFocusRingShowsForTheKeyboardOnly`; `dark-KvitIconButton.png` and
  `sepia-KvitIconButton.png` match, including the ring on the button with
  keyboard focus. Not yet done: the Qt button also shows its tooltip when it
  gains keyboard focus, and the Go one shows it only on pointer hover.
- [x] `KvitLink` (`Link`). Evidence: `TestLinkFollowsAndIsALink`;
  `light-KvitLink.png` and `highContrast-KvitLink.png` match, including the
  link cut short with "…" in a narrow column. The same tooltip gap as
  `KvitIconButton` applies.

### Structure (8)

- [x] `KvitHeader` (`Header`). Evidence: `TestTheHeaderKeepsItsThreePlaces`;
  `light-KvitHeader.png` matches.
- [x] `KvitSidebar` (`Sidebar`). Evidence:
  `TestASidebarItemSaysWhereItGoesAndHowMuchIsThere`; `light-KvitSidebar.png`
  matches for the expanded sidebar. The rail differs on purpose: in Qt the
  items keep the full sidebar width (232 px at 12 px) inside a 48 px rail,
  because a QML `Column` does not give its children its width, so the
  selection tint and the centred symbols reach far past the rail. Go's items
  take the sidebar's width.
- [x] `KvitSidebarItem` (`SidebarItem`). Evidence: the same test;
  `light-KvitSidebarItem.png` matches for the expanded items, and its rail
  specimen differs from Qt for the reason given under `KvitSidebar`. The
  tooltip in the rail shows on pointer hover after unison's own delay
  rather than Qt's 400 ms.
- [x] `KvitBreadcrumb` (`Breadcrumb`). Evidence:
  `TestTheBreadcrumbLinksBackAndCutsTheMiddle`; `light-KvitBreadcrumb.png`
  matches. The crumbs that can be followed are `Link`s, so unlike Qt's they
  can also be reached with Tab and followed with the keyboard.
- [ ] `KvitRegion`
- [ ] `KvitViewHead` (`ViewHead`): built and used for the head of every
  gallery page, with kvit-cash's side padding. Evidence so far:
  `TestTheViewHeadCountsInTheReadersLocale` (ported from
  `test_components.cpp`), `TestTheViewHeadIsARowTallWithItsControlsAtTheRight`.
  Not ticked until its own gallery page exists, which needs `KvitSearchField`
  and `KvitButton` for its specimen.
- [ ] `KvitStatusBar`
- [ ] `KvitWindow`

### Content (9)

- [ ] `KvitSectionHeading`
- [ ] `KvitRow`
- [ ] `KvitSlimRow`
- [ ] `KvitCard`
- [x] `KvitPanel` (`Panel`). Evidence: `TestAPanelDrawsTheRulesItIsAskedFor`;
  `light-KvitPanel.png` matches.
- [ ] `KvitPane`
- [x] `KvitDivider` (`Divider`). Evidence: `TestADividerIsOneHairline`;
  `sepia-KvitDivider.png` matches.
- [ ] `KvitDisclosure`
- [ ] `KvitEmptyState`

### Marks (7)

- [ ] `KvitChip`
- [ ] `KvitTag`
- [x] `KvitBadge` (`Badge`). Evidence:
  `TestABadgeCapsHidesAtZeroAndSaysItsNoun`; `dark-KvitBadge.png` and
  `highContrast-24px-KvitBadge.png` match.
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
- [x] `KvitTab` (`Tab`). Evidence:
  `TestATabSaysWhatItHoldsAndWhetherItIsChosen`; `dark-KvitTab.png`,
  `dark-24px-KvitTab.png` and `light-10px-KvitTab.png` match. The Qt tab also
  shows its explanation as a tooltip when it gains keyboard focus; the Go
  one shows it only on pointer hover so far, as with `KvitIconButton`.

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

- [x] The Phosphor icons, named by meaning: 83 meaning names over the
  font's 1,530 glyphs, generated from the Qt catalogue by
  `tools/import-qt-icons`. Evidence: `icons/icons_test.go` (ported from
  `test_icons.cpp`).
- [x] An unknown icon name draws a marked placeholder: a hairline box and
  "?" in the danger colour, and a warning in the log. Evidence:
  `TestAnUnknownIconIsMarked`; `light-KvitIcon.png`.

## Checks the Qt library runs, and their Go equivalents

- [x] Colour tables and contrast floors, per theme
  (`TestEveryTokenPairMeetsItsFloor`, `TestHighContrastMeetsStricterFloor`)
- [x] The type scale (`TestTheScalesAreOrderedAtEverySize`,
  `TestTypographyScale`)
- [x] The density and spacing scale (`TestEveryValueScalesWithOneSetting`)
- [x] The three chart ramps in all four themes, checked in OKLab against that
  theme's surfaces and every hue that already means something, under normal
  vision and three colour-vision deficiencies (`palette/palette_test.go`)
- [ ] Every component opening in all four themes with no warning
- [ ] `KvitTable`: 250,000 rows and twelve columns, scrolling smoothly and
  filtering in under 100 ms
- [ ] Every gallery sample compiling
- [x] No colour literal and no numeric font size outside the design values
  (`rules_test.go`; in Qt, qmllint at full strength)
- [ ] No unnamed geometry value outside the design values: not checked by a
  test yet, since a literal offset cannot be told from ordinary arithmetic
  by reading the source
- [ ] Every interactive component giving screen readers a role and a name
- [x] Generated files matching their sources: the theme tables and the icon
  catalogue (`tools/import-qt-*/main_test.go`)
- [ ] The same for the drawing stylesheet and the agent skill's component
  catalogue, which do not exist in Go yet

## Gallery

- [ ] One page per component, with its states and a working sample
- [x] `--page`, `--theme`, `--interface-size`, and in the window Ctrl+1 to
  Ctrl+4, Ctrl+plus and Ctrl+minus
- [ ] `--shots`: the fixed screenshot set (74 components × 4 themes, plus
  control-size variants, about 300 images), compared with the Qt set. The
  mechanism exists: `./build.sh --shots` writes the set under the Qt file
  names and stacks each image above the Qt one from
  `~/kvit-qt-reference/kvit-ui-0a0b210`. The foundations page and the four
  Foundation components have shots so far.
- [ ] `--catalog`: writes the agent skill's component catalogue

## Drawing a screen before building it

- [ ] `ux/tokens.css` generated from the Go design values
- [ ] `ux/frame.css`, `ux/render.sh`, `visual-language.md` and
  `PATTERN.md` copied across

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
- [ ] `d32c373` 2026-09-12: Give the wheel a distance, the card a fourth side, and a chart one baseline.
  Done so far: `ViewHead.Padding`, the view head's own side margin.
- [ ] `526b619` 2026-09-13: One press acts on one thing, and what acts says so before it is pressed
- [ ] `722906e` 2026-09-13: Name the shortest window the chrome holds
