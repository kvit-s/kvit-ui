# Second set of additions for Kvit Works: what was built and where it differs from the request

*Written 2026-11-12 in kvit-ui, against the request
`docs/proposals/architecture-review-kvit-ui-request-2.md` of the Kvit Works
repository, which asked kvit-ui at commit `5faae05` for options that let Kvit
Works delete the private copies it still keeps of rows, focus tracking, a
popover's focus, a dialog's drag, two layouts, a line-limited text and two
layout helpers.*

Kvit Works keeps those copies because the kvit-ui piece that would replace
each one differs from it in what is drawn or in what a press, a key or the
focus does. The request named eleven additions, one of them a defect fix, and
all eleven are built. Every new field is off at its zero value, so an existing
caller sees no change, with two exceptions, both defect fixes. The first is
the one the request marked, a `Label` with `MaxLines`; its fix has two parts,
the height it reports, as asked, and an empty line it drew where the limit
falls just after a line break, found while writing the requested test. The
second is a `Mover` whose drag never received the release of its left
button, found while comparing the Mover with the Kvit Works copy and then
confirmed by a test. This file is the report the request asks
for. It gives the names as built, each place where a quoted Kvit Works
behaviour is not reproduced exactly, and the change for existing callers with
the result of running the other applications' tests against it.

## State

- `./build.sh --test` passes (gofmt, `go vet`, every headless test).
- Each new test was checked to fail with its option switched off in the
  source, one option at a time, and the source restored afterwards. The test
  of the `Mover` fix failed before the fix was written.
- `agent/skills/kvit-ui/catalog.md` is rebuilt. `names_gen.go` and
  `ux/tokens.css` did not need rebuilding, since no design value changed.
- The gallery's screenshot set was written from `5faae05` and from this
  working copy, 1,209 pictures each, and compared pixel by pixel. 39 differ,
  the 13 pictures (four themes at three interface sizes, and the macOS form)
  of each of three pages: `KvitRow` and `KvitLine`, which gained the
  specimens and name the new options in their summaries, and `KvitPopover`,
  whose summary names `NoFocusBack` and `AnyPressEnds` and where only the
  summary's lines differ. No other page changed, the `KvitLabel` page with
  its line-limited label among them. The set was written before the `Mover`
  fix and before the Browse button of the new line specimen took its "…" from
  `Asks` instead of its words; neither changes what is drawn.
- The tests of kvit-works, kvit-hub, kvit-cash and kvit-notes were run against
  kvit-ui at `5faae05` and against this working copy; "App tests" below has
  the result.

## Where the code is

| Addition | File | Test |
|---|---|---|
| `ListRow.RingKeptOnPress` | `listrow.go` | `TestAListRowKeepsItsRingOnAPressWhenTold` in `listrow_test.go` |
| `KeyFocus.KeptOnPress` | `keyfocus.go` | `TestKeyFocusKeptOnAPressWhenTold` in `pieces_test.go` |
| `ListRow.NoFocusOnPress` | `listrow.go` | `TestAListRowPressedWithoutTakingTheFocus` |
| `ListRow.RingHalfIn` | `listrow.go` | `TestAListRowsMarks` |
| `ListRow.OnKeyDown` | `listrow.go` | `TestAListRowsKeyAndFocusHooks` |
| `ListRow.FitRoom`, `FitAtLeastForm`, `Arrange` | `listrow.go` | `TestAListRowsFitRoomAndArrange` |
| `ListRow.HoverYields` | `listrow.go` | `TestAListRowYieldsItsHoverToAButton` |
| `Popover.NoFocusBack` | `popover.go` | `TestAPopoverOpenedModal` in `pieces_test.go` |
| `Mover.AnyPressEnds` | `mover.go` | `TestAMoverEndsADragAtASecondPressWhenTold` in `pieces_test.go` |
| `Mover` drag that loses its release (fix) | `mover.go`, `window.go`, `windowkeys.go` | `TestAMoverDragEndsWhenASecondPressGoesElsewhere` in `pieces_test.go` |
| `LineLayout.AlignTop` | `line.go` | `TestALinesOptions` in `layouts_test.go` |
| `NewRowLayout` | `layout.go` | `TestTheSmallerLayouts` in `layouts_test.go` |
| `Label` height with `MaxLines` | `label.go` | `TestALabelKeepsToItsLineLimit` in `label_test.go` |
| `ShownIn`, `InsetsOf` | `stack.go` | `TestTheSmallerLayouts` |

The `FitRoom` test is a function of its own rather than part of
`TestAListRowsHeightAndPlacement`, because it presses the rows.

Gallery specimens: the `KvitRow` page's "A selection bar, rounded tints, and
the ring at the edge or half its width in" has a fourth row with `RingHalfIn`,
and the `KvitLine` page has a new specimen, "A button level with the top of a
field that grows", a path field showing an error with a Browse button beside
it on a line with `AlignTop`. The catalogue summaries of `KvitRow`,
`KvitLine` and `KvitPopover` name the new options.

## Names

Every exported name is as the request wrote it, and nothing was renamed.
`NewRowLayout` returns `unison.Layout`, as asked; its concrete type is
unexported. `KeyFocus.KeptOnPress` is the struct's first field, so
`kvitui.KeyFocus{KeptOnPress: true}` makes one.

Nothing exported was added beyond the request; the `Mover` fix adds only
unexported code.

## Behaviour to compare

### Keyboard ring on a press

- **`RingKeptOnPress`** leaves `keyboardFocus` as it was at a press; the row
  still sets the marker that says a press is giving it the focus. A press on
  a row that does not hold the focus therefore gives it the focus with no
  ring, and `OnFocus` is told `byPointer` true, with or without the option.
- **`KeyFocus.KeptOnPress`** behaves the same way for a panel of the
  caller's own.

### Press without taking the focus

- **`NoFocusOnPress`** skips the three lines the request names. A press then
  does not call `OnFocus`, since the row gains no focus. `OnPress`,
  `OnDoubleClick`, `PressOn` and the release work as without it. Tab still
  reaches the row while it is `Interactive`; the test checks this.

### Ring half a width inside the edge

- **`RingHalfIn`** strokes the ring with its centre line one ring width
  inside the row's box, so at a ring width of 2 pixels it covers the second
  and third pixel rows. The Kvit Works copy insets the box by half a width and
  then strokes half a width further in, which puts the centre line at the same
  place. `RingAtEdge` wins when both are set. The radius, when set, is
  `Radius` as for the other two positions; the copy's `drawRing` takes a
  radius of 0 in the quoted call.

### Keys with the repeat flag

- **`OnKeyDown`** is offered every key that reaches the row, before `OnKey`.
  That includes keys that come up from a control inside the row and keys on a
  row that is not `Interactive`, as `OnKey` already does. The row's own
  handling after it is unchanged: Return, Enter and Space only while the row
  itself holds the focus, only on a first press, and not with Command, Option
  or Control held. A Kvit Works row whose rule acts on a key from inside the
  row, or on a repeat, answers that key in `OnKeyDown` and returns true.

### Row height from content

- **`FitRoom`** is added to the height with `FitContent` whatever arranges
  the children: the row's own line, the slots, `Place` or `Arrange`. The
  content's box grows by it and the release box does not change, so with the
  row's default margins a release in the left or right margin still does not
  count and one in the added room does (the test checks both).
- **`FitAtLeastForm`** applies after `FitRoom`: the row's height is the
  larger of the fitted height and the form's.
- **`Arrange`** is given the row's own panel and the row's hint unchanged.
  The row's border is its margins, chevron room included, so a `LineLayout`
  counts them in the sizes it reports and places the children inside them.
  A layout that leaves the insets out of its sizes (a `LineLayout` with
  `InsetsUncounted`) leaves them out of the row's fitted height too. With
  `Arrange`, `FitContent` and `FitAtLeastForm`, the least, preferred and
  greatest heights are all the taller of the form and the arranged
  content, and both widths come from `Arrange`, which is what `rowLine`
  reports. The greatest width is unison's default maximum, as in `rowLine`.

### Hover over parts with a hover of their own

- **`HoverYields`** walks from the innermost panel under the pointer up to
  the row, not counting the row, and gives up the hover at the first panel
  with a `MouseEnterCallback` or a `MouseExitCallback`. Over a label nested
  inside another plain panel the row stays tinted, which the Kvit Works copy
  does not do, as the request expected; the field's comment says so.
- **The row's own tooltip.** The copy hands the row's enter and exit
  callbacks to its plain parts, so the row's tooltip also shows over its
  words. A `ListRow` with `HoverWhole` hides its own tooltip when the pointer
  moves from its padding onto a child, as before. A tooltip that should
  follow the row's hover, `HoverYields` included, is an `AddTooltip` with
  `HoverTold`, told from `OnHover` with `SetHovered`.

### Focus after a popover closes

- **`NoFocusBack`** leaves the focus where closing the popover leaves it. In
  the test, after a modal popover holding the focus closes, the button that
  held the focus before it opened does not hold it again.

### End of a drag at a second press

- **`AnyPressEnds`** ends the drag at any press that arrives while one is
  under way, and with `EscapeAbandons` ends the window's registered drag, so a
  later Escape puts nothing back. The press then goes on as any press does:
  a left press in the band starts a new drag.
- **Which panel unison gives the second press to.** unison sends the drags
  and the release to the panel whose `MouseDownCallback` takes the press. A
  right press reaches the Mover's panel, and stays there, only when the
  callback the Mover replaced takes it; the test gives the panel such a
  callback, as a dialog's body has. Where nothing on the panel takes it, the
  rest of the gesture goes elsewhere, which is the defect described under
  "A `Mover` drag that loses its release" in "Changes for existing callers".

### Line with children at the top

- **`AlignTop`** frames each child at the top of the content at the height
  it asks for at the width it gets, without cutting it to the line's height,
  and wins over `Top`. The line's reported sizes are unchanged; as the request
  says, `fieldAndButton` reports different least and preferred widths, which
  this does not reproduce.

### Row layout that counts its children

- **`NewRowLayout`** sets the column count to the number of children, hidden
  ones included, and gives a child with no layout data the middle alignment,
  before every size question and every layout. The gap is read from its
  `Measure` each time, so it follows the interface size.

### Exported helpers for layouts

- **`ShownIn` and `InsetsOf`** have the bodies `shownIn` and `insetsOf` had;
  the unexported names are one-line calls to them.

## Changes for existing callers

### `Label` with `MaxLines`

Two defects were fixed, and the second was found while writing the test the
request asked for.

1. **The height asked for.** A wrapping label with `MaxLines` above zero,
   offered a width, now reports the height of the text as it is drawn: the
   lines kept and the cut last line. Before, it reported that height only when
   the text's widest line, laid out without wrapping, was wider than the room;
   a text of short lines of its own, more of them than `MaxLines`, reported
   the height of all of them. Offered no width, a label still reports the
   height of all its own lines.
2. **An empty line after a line break.** Where `MaxLines` cuts the text just
   after a line break of the text's own, the part kept whole ended with that
   line break, which the text package lays out as an empty line of its own.
   The label then drew, and reported, one line more than `MaxLines`: for
   `"a\nb\nc\nd\ne"` with `MaxLines = 2`, the lines "a", an empty line, and
   "b c d e" cut short. The line break is now trimmed from the part kept
   whole, so it draws "a" over "b c d e", two lines. With the first fix alone
   that label reported three lines' height, 51 pixels at the default
   interface size against 17 for a line, which the requested test (between
   1.5 and 2.5 lines) does not accept.

The second fix also changes a long wrapped text whose kept lines end at a
line break of its own, which took the cut-short path before this work and drew
the empty line. A label with `MaxLines` of 0 or 1 is unchanged by either fix.

In the local checkouts of the four applications, no application's own code
sets `MaxLines` on a kvit-ui `Label`: Kvit Works uses its `para` type and
kvit-hub its `parts.Words` for line-limited text. The gallery's one label
with `MaxLines` holds a single sentence with no line break.

### A `Mover` drag that loses its release

**The defect.** A `Mover` ended its drag only in its own `MouseUpCallback`.
When a second button was pressed during a drag and the Mover's panel did not
take that press (a panel whose own callback takes left presses only), unison
sent the rest of the gesture, the left button's release included, to
another panel or to none. The Mover then went on believing the drag was
under way. `TestAMoverDragEndsWhenASecondPressGoesElsewhere` shows what
followed, without `AnyPressEnds`: after both buttons were let go, a drag of
50 pixels left and 20 up below the band, which the panel's own callback
took, moved the panel 50 pixels left and 64 down, because the Mover measured
it from the old drag's grab point; and in a first version of the test,
Escape pressed after the buttons were let go put the panel back where the
old drag began. With `AnyPressEnds` the right press ended the drag, so
that setting passed before the fix.

**The fix.** A Kvit window now has a `MouseUpCallback` on its unison window,
which sees every release before it reaches a panel. When a Mover's drag
starts it asks the window, through `watchRelease`, to be told when the left
button is next let go. The window tells it once that release has been
handled, wherever it went; if the drag is still under way then, the Mover
ends it as `AnyPressEnds` does, leaving the panel where it is and ending the
window's registered drag for Escape. A drag that receives its own release
ends there as before, and the later word from the window finds nothing to
do. A drag counter keeps the word about one drag from ending the next.

**What changes for existing callers.** Without `AnyPressEnds` a second
button pressed during a drag is still passed on, and the drag still goes on
while the panel's own callback takes that press. The only change is that the
drag ends when the left button is let go, where before it could stay under
way. The window's new callback returns false, so every release still reaches
its panel. No application in the local checkouts uses `NewMover`. An
application that sets its own `MouseUpCallback` on a Kvit window's unison
window would replace the window's, as one that sets `MouseDownCallback`
already replaces the popups' handling; none of the four does. The other
drags that register with the window for Escape (a dialog's move, a split
view's divider, a table header's column) were not examined for the same
defect.

## App tests

Each application's `go test ./...` was run on the Ubuntu development machine
against kvit-ui at `5faae05` and against this working copy as it was before
the `Mover` fix, through a temporary module file that points both kvit-ui and
its patched go-text copy at the version under test; `go list -m` confirmed
the replacement in each run, and the applications' own files were not
changed. The local checkouts
pin kvit-ui at `bb47829` (Kvit Works and kvit-hub), `2a4ecab` (kvit-cash)
and `5faae05` (kvit-notes); the Kvit Works checkout is therefore older than
the `2800d10` the request names.

| App | Tests failing in both runs | Failing only with this work | Failing only without it |
|---|---|---|---|
| kvit-works | 7 | `TestAStartNamingATrackMakesAnAgentOnIt` | `TestAnAgentsRequestForAProjectIsAnsweredAsItsTurn` |
| kvit-hub | 1 | none | none |
| kvit-cash | 4 | none | none |
| kvit-notes | 2 | none | none |

The two Kvit Works tests that failed in one run only were each run alone
three more times against both versions. `TestAStartNamingATrackMakesAnAgentOnIt`
(`views/reader`) failed once in three against this working copy and passed
three times against `5faae05`; its failure in the full run was a temporary
folder that could not be removed ("unlinkat …/project/.git: directory not
empty"), which does not involve drawing or layout.
`TestAnAgentsRequestForAProjectIsAnsweredAsItsTurn` (`agent/cluster`, which
does not use kvit-ui) failed twice in three against `5faae05` and passed three
times against this working copy. Both fail or pass by chance in either
version, so neither result is caused by this work.

The tests failing in both runs, for the record: in kvit-works
`TestTheWindowsStatesFollowTheInterfaceRules`,
`TestTheLineUnderTheFieldsSaysWhatWouldFail`, `TestWorksJourney`,
`TestTheWorksShellLoadsAsARealAlternateComposition`,
`TestAConversationComesBackWhereItWasLeft`,
`TestALongConversationInAJoinedClientLoadsItsEarlierTurns` and
`TestSeveralFormsWorkOverTheJoint`, with the subtests
`TestTheWindowsStatesFollowTheInterfaceRules/ask-menu-question` and
`TestWorksJourney/TheNewTrackButtonStartsATrack`; in kvit-hub
`TestRunAgainIsOnTheNewestAttemptOnly`; in kvit-cash
`TestFiltersFullScaleWorkspaceWithinBudget`,
`TestGeneratesFullScaleIndexedWorkspace`,
`TestHiddenCountAtFullScaleStaysWithinBudget` and
`TestEachPeriodIsTheCalendarPeriodHoldingToday`; in kvit-notes
`TestAMarginItemSitsBesideItsLine` and `TestTallInlineMathGrowsItsLine`.

After the `Mover` fix, which gives every Kvit window a release handler, each
application's tests were run once more against this working copy. Every test
that failed in both earlier runs failed again, except kvit-works'
`TestTheLineUnderTheFieldsSaysWhatWouldFail`, which passed. Two tests failed
that had passed in both earlier runs, and each was then run again against
both versions, alone five times and in its whole package three times:

- kvit-notes `TestDiagramResetSceneDropsLastGood` ("the reset kept the
  scene") failed once in three package runs against this working copy and
  once in three against `5faae05`, and passed every time alone; it fails by
  chance in either version.
- kvit-cash `TestAnEntryThatCannotBeReversedSaysSoAsAMark` failed because
  the popup it read said "Opens transactions this change wrote", the tooltip
  of something other than the mark it rests the pointer on. It passed all
  eight further runs against each version. The test only moves the pointer
  and reads a tooltip, with no button pressed, and the window's new handler
  does anything only when a button is released during a Mover's drag, so
  this failure is not attributed to the fix; the cause was not found.

## Left out

Nothing the request asked for was left out.
