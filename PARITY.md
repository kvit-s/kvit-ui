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
  keyboard focus. Its tooltip shows on pointer hover and on keyboard focus
  (`TestATooltipShowsForTheKeyboardBesideItsControl`).
- [x] `KvitLink` (`Link`). Evidence: `TestLinkFollowsAndIsALink`;
  `light-KvitLink.png` and `highContrast-KvitLink.png` match, including the
  link cut short with "…" in a narrow column.

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
  tooltip in the rail shows after the half-second every Kvit tooltip waits,
  where Qt's waits 400 ms.
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

- [x] `KvitChip` (`Chip`, with `Tone`). Evidence: `TestMarksSayWhatTheyMean`;
  `light-KvitChip.png` matches.
- [x] `KvitTag` (`Tag`). Evidence: `TestMarksSayWhatTheyMean`, which also
  checks the colour is announced by its name ("groceries, Red");
  `dark-KvitTag.png` matches.
- [x] `KvitBadge` (`Badge`). Evidence:
  `TestABadgeCapsHidesAtZeroAndSaysItsNoun`; `dark-KvitBadge.png` and
  `highContrast-24px-KvitBadge.png` match.
- [x] `KvitSlug` (`Slug`). Evidence: `TestMarksSayWhatTheyMean`,
  `TestAnIdentifierIsCutInTheMiddle` (the text layer's middle eliding);
  `sepia-KvitSlug.png` matches.
- [x] `KvitDot` (`Dot`, with `Shape`). Evidence: `TestMarksSayWhatTheyMean`;
  `highContrast-KvitDot.png` matches, the diamond's box a little wider:
  Qt draws the turned square past the dot's box, and unison clips to it, so
  the box grows and what follows a diamond sits 2 px further along.
- [x] `KvitSignal` (`Signal`). Evidence: `TestASignalCountsOnlyPastOne`;
  `light-KvitSignal.png` matches.
- [x] `KvitPip` (`Pip`). Evidence: `TestMarksSayWhatTheyMean`;
  `light-KvitPip.png` matches.

### Quantities (2)

- [x] `KvitFigure` (`Figure`). Evidence: `TestAFigureSaysWhatWasMeasured`;
  `sepia-KvitFigure.png` matches.
- [ ] `KvitBeforeAfter`

### Controls (9)

- [x] `KvitButton` (`Button`). Evidence:
  `TestAButtonSaysWhatItDoesInEveryState`; `light-KvitButton.png` matches
  for the primary and ordinary forms. The quiet form has kvit-cash's chip
  ground and faint border at rest, so it differs from the Qt shot.
- [x] `KvitChipButton` (`ChipButton`). Evidence:
  `TestAnUnavailableChipStaysReachableAndDoesNothing`;
  `light-KvitChipButton.png` matches. It has both versions' additions: the
  Qt library's `explanation` and `current`, and kvit-cash's quick-filter
  state (`Selectable`, `Selected`), a control's height and filled with the
  accent when on.
- [x] `KvitStepper` (`Stepper`). Evidence: `TestAStepperStaysInItsRange`;
  `light-KvitStepper.png`, and the gallery's own header. unison reads a spin
  button as text and leaves out its children, so the minus and plus reach a
  screen reader as the spin button's increment and decrement actions.
- [x] `KvitField` (`Field`). Evidence:
  `TestAFieldTakesTypingAndSaysWhatIsWrong`; `light-KvitField.png` matches,
  except that the error message takes room of its own under the field,
  where Qt's is drawn over the space below it. The editing is unison's field
  (caret, selection, clipboard, undo, screen reader text), and the typed
  text is drawn by unison's text engine rather than the Kvit text layer.
  kvit-cash's `shortcut` is `Shortcut`, shown with the label as the tooltip
  on pointer hover.
- [x] `KvitTextArea` (`TextArea`). Evidence:
  `TestAReadOnlyTextAreaTakesNoTyping`; `light-KvitTextArea.png` matches as
  far as the Qt shot reaches (it ends at the window's bottom edge). unison's
  field has no read-only mode, so `ReadOnly` (on `Field` too) drops typing,
  the keys that delete or break lines, and cut and paste.
- [x] `KvitSearchField` (`SearchField`). Evidence:
  `TestASearchFieldClearsWithEscapeAndCountsWhatItLeft`;
  `dark-KvitSearchField.png` matches, with kvit-cash's clear button a
  control's height square.
- [x] `KvitCheck` (`Check`). Evidence: `TestACheckHasThreeStates`;
  `dark-KvitCheck.png` matches, with the 6 px of padding on each side that
  Qt Quick's style gives a CheckBox and KvitCheck keeps; it is 6 px at
  every interface size, as the style's is. A disabled box that is checked draws its
  mark in the disabled text colour; the Qt one draws it in the colour made
  for the accent fill, which a disabled box does not have, so it vanishes.
- [x] `KvitSelect` (`Select`). Evidence:
  `TestASelectStepsWithTheArrowsAndOpensAMenu`; `sepia-KvitSelect.png`
  matches. The list is unison's menu with the chosen option ticked, as the
  owner chose for menus; Up and Down change the choice without opening it.
- [x] `KvitTab` (`Tab`). Evidence:
  `TestATabSaysWhatItHoldsAndWhetherItIsChosen`; `dark-KvitTab.png`,
  `dark-24px-KvitTab.png` and `light-10px-KvitTab.png` match.

### Feedback (7)

Everything that floats is shown in a layer the Kvit `Window` keeps above its
content (`Window.Show`), as the Qt pop-ups are drawn inside their window.
Menus are the exception: they are unison's, by the owner's choice.

- [x] `KvitTooltip` (every control's tooltip, and `UI.ShowTooltip`).
  Evidence: `TestATooltipShowsForTheKeyboardBesideItsControl`;
  `light-KvitTooltip.png` differs on purpose: the tooltip is placed as
  kvit-cash's copy places it, beside the control, then inside it, then
  below, then above, where the older Qt shot centres it above. It shows on
  pointer hover and on keyboard focus after half a second, and stays at
  least three seconds. A screen reader is given it as a tooltip named by
  its words. Outside a Kvit window unison's own hover tooltip is the
  fallback.
- [x] `KvitPopover` (`Popover`). Evidence:
  `TestAPopoverClosesOnEscapeAndGivesTheFocusBack`; `light-KvitPopover.png`
  matches. Escape and a press outside close it, and the focus goes back to
  where it was unless the reader has put it elsewhere, as the newer Qt
  popover does.
- [x] `KvitHint` (`Hint`). Evidence: the popover test;
  `light-KvitHint.png` matches.
- [x] `KvitHoverCard` (`HoverCard`). Evidence: `light-KvitHoverCard.png`
  matches. A screen reader passes over it.
- [x] `KvitToast` (`Toast`). Evidence:
  `TestAToastWithoutAnActionGoesAndOneWithAnActionStays`;
  `light-KvitToast.png` matches, the Undo in kvit-cash's quiet form. It is
  announced to a screen reader when shown.
- [x] `KvitNotice` (`Notice`). Evidence:
  `TestANoticeSaysItsConditionAndIsDismissedOnlyWhenItCanBe`;
  `light-KvitNotice.png` matches. It has kvit-cash's `ActionHint`.
- [x] `KvitDialog` (`Dialog`). Evidence: `TestADialogHasToBeAnsweredFirst`;
  `light-KvitDialog.png` matches. It has kvit-cash's close control in the
  title row and leaves out a button with no words, and a dialog with no
  buttons has no row for them.

### Data (11)

A Qt item may draw past its own box, and several charts do: a trend's top
and bottom axis values are centred on gridlines at its edges, and a gauge's
pace tick stands out of the bar. unison clips each panel's drawing to its
box, so a component can hand the part outside to the panel above it that
draws overflow (`spill.go`). The window does this for its body and its
sidebar: it draws each component's overflow after the page, under the rail
and under anything floating, clipped where a Qt item would be clipped, at a
scroll area's view and at a disclosure's body. Outside a Kvit window the
component draws it itself, cut off at its box.

- [x] `KvitBar` (`Bar`). Evidence: `TestChartsSayWhatTheyShow`;
  `light-KvitBar.png` matches.
- [x] `KvitStackedBar` (`StackedBar`). Evidence:
  `TestChartsSayWhatTheyShow`; `light-KvitStackedBar.png` matches.
- [x] `KvitSpark` (`Spark`). Evidence: `TestChartsSayWhatTheyShow`;
  `light-KvitSpark.png` matches.
- [x] `KvitTrend` (`Trend`). Evidence: `TestATrendSaysWhenItIsEmpty`;
  `light-KvitTrend.png` and `dark-24px-KvitTrend.png` match, including the
  axis values half a line outside the trend's box and the bottom gridline
  on the pixel row below it, which are drawn as overflow. The values are
  rounded half away from zero, as JavaScript's `toFixed` rounds them, so a
  scale of 0 to 1 is labelled 1, 1, 0 in both.
- [x] `KvitDistribution` (`Distribution`). Evidence:
  `TestChartsSayWhatTheyShow`; `light-KvitDistribution.png` matches.
- [x] `KvitGauge` (`Gauge`). Evidence: `TestChartsSayWhatTheyShow`,
  `TestAGaugesPaceTickStandsOutOfTheBar`; `light-KvitGauge.png` matches,
  the pace tick drawn as overflow.
- [x] `KvitDelta` (`Delta`). Evidence: `TestChartsSayWhatTheyShow`;
  `light-KvitDelta.png` matches.
- [x] `KvitStatTile` (`StatTile`). Evidence: `TestChartsSayWhatTheyShow`;
  `light-KvitStatTile.png` matches. unison's flex layout gives a hidden
  child its room and its gap, so the tile takes out a part it does not have
  rather than hiding it.
- [x] `KvitFigureBlock` (`FigureBlock`). Evidence:
  `TestChartsSayWhatTheyShow`; `light-KvitFigureBlock.png` matches.
- [x] `KvitCell` (`Cell`, with `CellKind` and `CellValue`). Evidence:
  `TestACellSaysWhatItHolds`, `TestAnAmountShowsWhatIsBehindItInAHoverCard`;
  `light-KvitCell.png` matches, except that a figure's unit sits one pixel
  lower, as it does on the `KvitFigure` page. It has kvit-cash's `leads`
  (the row's name in the link colour) and the hover card behind an amount.
  A date given as a date is drawn as 2006-01-02, which is what the C locale
  writes and what the Qt shots show; the Qt cell writes a date in the
  reader's locale, and Go has no locale date formats to do that with, so a
  model that wants another form passes the date as text. unison reads a
  cell as text and leaves out its children, so a Marks cell is a group of
  its marks, each named, and a Check cell leaves the naming to its box.
- [x] `KvitTable` (`Table`, over a `TableModel`, with `BenchmarkTableModel`).
  Evidence: `TestTheHeaderSortsAndResizesFromTheKeyboard`,
  `TestTheKeyboardWalksTheRowsAndTicksTheirBoxes`,
  `TestTickingABoxDoesNotAlsoOpenTheRow`,
  `TestACellCutShortIsShownWholeUnderThePointer`,
  `TestTheEmptyStateFollowsTheRowCount`, the two measurements under "Checks
  the Qt library runs"; `light-KvitTable.png` and `dark-24px-KvitTable.png`
  match, except that the header is a slim row tall, as kvit-cash's commit
  f4cb189 makes it, so the rows sit 6 px lower, and the second table has
  kvit-cash's strip of row marks down its right edge. The table draws only
  the cells in view, each with the component that draws that kind
  elsewhere, placed and drawn without being added to the panel tree, as
  unison's own table draws its cells; it holds no panel per cell or per row.
  It has everything in kvit-cash's copy: the header as one stop in the tab
  order with a cursor; sorting, resizing (a drag on an edge, or Ctrl+Left
  and Ctrl+Right) and moving (a drag, or the menu) columns, hiding and
  showing them; the column menu; a Check column with a box for every row
  shown in its header; the strip of row marks and its tooltip; the row's
  name in the link colour; one press opening a row; the ring round the
  keyboard's row; the full value of a cut-short cell under the pointer and
  under the keyboard; and the wheel distance. Differences, each on purpose:
  - The ring round the keyboard's row shows once the keyboard has been used
    in the rows, by Tab, `FocusRow` or a key, and not after a mouse press,
    which is the rule every Go control follows; the Qt ring shows whenever
    the table holds the focus.
  - Ctrl+C puts the selected rows on the clipboard. The Qt key handler
    calls `copySelection()` and drops the text it returns, so nothing is
    copied. The Go copy is in the column order shown, as the Qt comment
    says it should be, where the Qt code orders it by model column, and
    leaves out the column of boxes.
  - Moving a column moves it past the next column shown. The Qt table moves
    it one place in the full order, hidden columns included, so moving it
    past a hidden column changes nothing on the screen.
  - Widths the reader sets are kept in design pixels (`Widths`), so they
    follow the interface size; the column order (`Order`) and hidden set
    are public for an application to save with its view.
  - The model is told nothing about editing: the Qt table declares
    `editable` and never reads it.

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
- [x] `KvitSegmented` (`Segmented`). Evidence:
  `TestASegmentedControlChoosesOne`; `light-KvitSegmented.png` and the
  gallery's own header. It is kvit-cash's form, a tab's height with the
  chosen segment filled with the accent, so it differs from the Qt shots,
  where the chosen segment is outlined.
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
- [x] `KvitTable`: 250,000 rows and twelve columns, scrolling smoothly and
  filtering in under 100 ms. Evidence:
  `TestFilteringAQuarterMillionRowsStaysUnderTheBudget` (about 5 ms to
  filter to nothing and to 20,832 matches, against 100 ms) and
  `TestATableDrawsAPageOfAQuarterMillionRowsInsideAFrame`, which draws a
  1200 by 700 table forty pages down from the middle of 250,000 rows and
  requires the middle page to fit a 60 Hz frame (median 1.4 ms, worst
  2.8 ms, against 16 ms), and a page to ask the model for no more cells than
  are on screen (253). The time is the table's drawing into an offscreen
  canvas, without the canvas's setup and readback, as the Qt check times its
  page layout without the window grab.
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

- [ ] One page per component, with its states and a working sample: 58 of
  74 so far, each listed in the sidebar; the rest are listed faint.
- [x] The gallery is built from the library's own components, as the Qt
  gallery is: `Window`, a `Header` holding `Segmented` and `Stepper`, a
  sidebar of `SectionHeading`s and `ListRow`s under a `SearchField` that
  filters them, the page in a `Region`, and a `StatusBar` that says what the
  screenshot run is writing. Evidence: `TestTheSidebarFilters`, and every
  `compare/*.png`, whose frame now matches the Qt gallery's except for the
  extra Foundations page and kvit-cash's segmented control.
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
  search field's clear button at a control's height, the segmented control
  at a tab's height, the status bar's
  overflow link hidden when nothing is hidden, the card's chevron, the empty
  state's action reachable by the keyboard, the pane's floor (not needed
  in unison), the tooltip kept inside the window and beside its control,
  and in the table the hover card behind an amount, the strip of row marks,
  and the column menu's entries for moving, widening and narrowing.
- [x] `f4cb189` 2026-09-11: Raise a table's header to the height of a row somebody presses.
  Evidence: `light-KvitTable.png`, the header a slim row tall and its box
  and menu button a control's height.
- [ ] `1c237df` 2026-09-11: Say what a mark opens, and draw no control with nothing on it.
  Done so far: the row's, the card's and the table's `OpensLabel`, the
  pane's `closable`, the dialog's buttons left out when they have no words.
- [ ] `364c3dc` 2026-09-11: Float a view over a list, open the rail on hover, draw a choice as chosen.
  Done so far: the rail opening on hover and for the keyboard, the segmented
  control's filled choice, the chip button's quick-filter state, the
  dialog's close control in its title row.
- [ ] `d32c373` 2026-09-12: Give the wheel a distance, the card a fourth side, and a chart one baseline.
  Done so far: `ViewHead.Padding`, the view head's own side margin, and the
  region's and the table's wheel distance.
- [ ] `526b619` 2026-09-13: One press acts on one thing, and what acts says so before it is pressed.
  Done so far: the quiet button's ground at rest, the field's shortcut, one
  press acting on one thing in a card, the tooltip's placement, the
  notice's action hint, and in the table the row's name in the link colour
  (`OpensColumn`) and the ring round the keyboard's row.
- [x] `722906e` 2026-09-13: Name the shortest window the chrome holds.
  `Interface.HeightFloor`, and `Window`'s minimum height.
