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
- [x] `KvitRegion` (`Region`). Evidence:
  `TestARegionScrollsBesideItsBarAndAnswersTheWheelAndKeys`;
  `light-KvitRegion.png` matches, except that Qt's specimen rows are a fixed
  460 px in a 438 px column and run under the bar, while Go's rows take the
  column's width. It includes kvit-cash's keyboard scrolling, bringing the
  focused control into view, and the wheel travelling the desktop's lines
  per notch of one slim row each, eased over a few frames.
- [x] `KvitViewHead` (`ViewHead`). Evidence:
  `TestTheViewHeadCountsInTheReadersLocale` (ported from
  `test_components.cpp`), `TestTheViewHeadIsARowTallWithItsControlsAtTheRight`;
  `light-KvitViewHead.png` matches except for kvit-cash's side padding, which
  insets the title and the controls by the view margin.
- [x] `KvitStatusBar` (`StatusBar`). Evidence:
  `TestAStatusBarKeepsWhatDoesNotFitBehindACount`,
  `TestTheStatusBarMenuHoldsWhatDidNotFit`; `light-KvitStatusBar.png`
  matches, the overflow link 6 px to the left of Qt's. The overflow menu is
  unison's menu (see `KvitMenu`), and the overflow link is hidden, not drawn
  at no width, when nothing is hidden (kvit-cash).
- [x] `KvitWindow` (`Window`). Evidence:
  `TestTheWindowCollapsesItsSidebarAndOpensTheRailOverTheBody`,
  `TestTheSourceOnlySamplesRun`; `dark-KvitWindow.png` matches (the page
  shows the sample's code only, as in Qt). It includes kvit-cash's rail that
  opens over the body under the pointer or the keyboard, and its height
  floor. The rail opens for focus moved by a key, not for the focus unison
  gives the first control when the window becomes active.

### Content (9)

- [x] `KvitSectionHeading` (`SectionHeading`). Evidence:
  `TestASectionHeadingOpensAndItsActionStaysItsOwn`;
  `light-KvitSectionHeading.png` matches. Two differences: it opens and
  closes by itself and reports the new state, where the Qt heading only
  signals and leaves the state to its caller; and unison treats a heading as
  text and leaves its children out, so a heading that opens and closes is
  announced as a button with its expanded state, and a fixed heading with an
  action as a group named by the heading, to keep the action reachable.
- [x] `KvitRow` (`ListRow`, since `Row` is the Go layout helper). Evidence:
  `TestAListRowActsOnlyWhenItSaysItDoes`; `dark-KvitRow.png` matches, with
  one intended difference: rows that can be pressed carry kvit-cash's
  chevron at their trailing edge, so their content ends 16 px sooner.
- [x] `KvitSlimRow` (`SlimRow`). Evidence: `TestASlimRowSaysItsPartsInOrder`;
  `dark-KvitSlimRow.png` matches. A slim row nobody presses is announced as
  a named group rather than as text, because unison drops the children of a
  text node and the figure would not be read.
- [x] `KvitCard` (`Card`). Evidence: `TestACardOpensOnceWhenARowInsideActs`;
  `light-KvitCard.png` matches, with kvit-cash's chevron on the pressable
  card. It also does kvit-cash's other two card changes: `OpensLabel`, and a
  press on something inside that acts on its own not also opening the card.
  Unlike the Qt card, a pressable card takes the keyboard focus and opens on
  Return, Enter or Space.
- [x] `KvitPanel` (`Panel`). Evidence: `TestAPanelDrawsTheRulesItIsAskedFor`;
  `light-KvitPanel.png` matches.
- [x] `KvitPane` (`Pane`, placed with `WithPane`). Evidence:
  `TestAPaneSlidesOverTheRightEdge`; `light-KvitPane.png` matches. kvit-cash's
  `closable` is `Closable`. Its floor, which stops a press or a wheel turn on
  the pane reaching the list under it, is not needed: unison delivers both
  to the topmost panel and its own ancestors only.
- [x] `KvitDivider` (`Divider`). Evidence: `TestADividerIsOneHairline`;
  `sepia-KvitDivider.png` matches.
- [x] `KvitDisclosure` (`Disclosure`). Evidence:
  `TestDisclosuresInAGroupOpenOneAtATime`; `sepia-KvitDisclosure.png`
  matches. It opens and closes by itself, and sections sharing a `Group`
  close each other, which the Qt disclosure's `group` property promises but
  leaves to its callers.
- [x] `KvitEmptyState` (`EmptyState`). Evidence:
  `TestAnEmptyStateSaysWhatWouldBeHereAndOffersTheWayIn`;
  `dark-KvitEmptyState.png` matches, except that the Go drop target shows its
  dashed edge and the Qt shot does not: the Qt edge is drawn on a Canvas,
  which draws nothing under the software renderer the Qt screenshots were
  taken with. Its action button takes the keyboard focus, as kvit-cash made
  it do.

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

- [x] `KvitFigure` (`Figure`). Evidence: `TestAFigureSaysWhatWasMeasured`;
  `sepia-KvitFigure.png` matches.
- [ ] `KvitBeforeAfter`

### Controls (9)

- [x] `KvitButton` (`Button`). Evidence:
  `TestAButtonSaysWhatItDoesInEveryState`; `light-KvitButton.png` matches
  for the primary and ordinary forms. The quiet form has kvit-cash's chip
  ground and faint border at rest, so it differs from the Qt shot. Its
  explanation shows as a tooltip on pointer hover only, as with
  `KvitIconButton`.
- [ ] `KvitChipButton`
- [ ] `KvitStepper`
- [x] `KvitField` (`Field`). Evidence:
  `TestAFieldTakesTypingAndSaysWhatIsWrong`; `light-KvitField.png` matches,
  except that the error message takes room of its own under the field,
  where Qt's is drawn over the space below it. The editing is unison's field
  (caret, selection, clipboard, undo, screen reader text), and the typed
  text is drawn by unison's text engine rather than the Kvit text layer.
  kvit-cash's `shortcut` is `Shortcut`, shown with the label as the tooltip
  on pointer hover.
- [ ] `KvitTextArea`
- [x] `KvitSearchField` (`SearchField`). Evidence:
  `TestASearchFieldClearsWithEscapeAndCountsWhatItLeft`;
  `dark-KvitSearchField.png` matches, with kvit-cash's clear button a
  control's height square.
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

- [x] `KvitScrollBar` (`NewScrollBar`, unison's scroll bar drawn in the
  Kvit colours in a strip of its own). Evidence: the region test;
  `light-KvitScrollBar.png` matches, with the same row-width difference as
  `KvitRegion`.
- [ ] `KvitMenu`, as `UI.ShowMenu`: unison's menu, drawn in the Kvit colours
  and type on Windows and Linux and native on macOS, as the owner chose on
  2026-09-26. It has text, a shortcut, a check and separators, and answers
  the arrow keys, Return and Escape. Not ticked: it cannot show the Qt
  menu's symbols, explanations or danger colour, and its minimum width is
  not the Qt menu's 200 px.
- [ ] `KvitMenuItem`: the `MenuItem` lines of `ShowMenu`, with the same gaps.
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

- [ ] `97575ee` 2026-09-11: Say what opens something, and keep a popup inside the window.
  Done so far: the region's keyboard scrolling and bringing focus into view,
  the slim row's kind capped at a quarter of the row, the row's chevron, the
  search field's clear button at a control's height, the status bar's
  overflow link hidden when nothing is hidden, the card's chevron, the empty
  state's action reachable by the keyboard, and the pane's floor (not needed
  in unison).
- [ ] `f4cb189` 2026-09-11: Raise a table's header to the height of a row somebody presses
- [ ] `1c237df` 2026-09-11: Say what a mark opens, and draw no control with nothing on it.
  Done so far: the row's and the card's `OpensLabel`, the pane's `closable`.
- [ ] `364c3dc` 2026-09-11: Float a view over a list, open the rail on hover, draw a choice as chosen.
  Done so far: the rail opening on hover and for the keyboard.
- [ ] `d32c373` 2026-09-12: Give the wheel a distance, the card a fourth side, and a chart one baseline.
  Done so far: `ViewHead.Padding`, the view head's own side margin, and the
  region's wheel distance.
- [ ] `526b619` 2026-09-13: One press acts on one thing, and what acts says so before it is pressed.
  Done so far: the quiet button's ground at rest, the field's shortcut, one
  press acting on one thing in a card.
- [x] `722906e` 2026-09-13: Name the shortest window the chrome holds.
  `Interface.HeightFloor`, and `Window`'s minimum height.
