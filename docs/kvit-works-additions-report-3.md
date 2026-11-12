# Third set of additions for Kvit Works

*Written 2026-11-12 in kvit-ui, against the request
`docs/proposals/architecture-review-kvit-ui-request-3.md` of the Kvit Works
repository, which asked kvit-ui at commit `c28b8b5` for two options on
`ListRow` and one test, so that Kvit Works can delete the last private
copies it keeps of a list row's focus and key handling.*

Kvit Works is a desktop application built on kvit-ui. After two earlier sets
of additions (`docs/kvit-works-additions-report.md` and
`docs/kvit-works-additions-report-2.md`), it still kept four private copies,
each because the kvit-ui piece that would replace it behaves differently. The
request took the four in turn and asked for three things: an option that
makes a row forget a waiting request for the keyboard focus when the row
loses the focus (`FocusRequestEndsOnBlur`), an option that leaves Return,
Enter and Space pressed with Shift to what the row sits in
(`NoShiftActivate`), and a test, moved into kvit-ui, that a row's
registration with the window's pointer watch does not keep the row alive. It
also named three differences it did not ask to change, and asked whether
they read in the source as it describes them.

All three are built, under the names the request used. Both options are off
at their zero value, and nothing else changed for an existing caller: the four
applications' tests fail in the same places against this work as against
`c28b8b5`, and of the gallery's 1,209 screenshots only the 13 of the `KvitRow`
page differ, in its summary text. This file says where the code is, how each
addition behaves where it could be read more than one way, what was found
about the three cases not requested, and the two comparisons. Of the three
cases, two read exactly as the request describes, and the first of them, the
ring shown while a pointer button is held on a row, holds only when the press
is what gave the row the focus.

## State

- `./build.sh --test` passes (gofmt, `go vet`, every headless test).
- Each new test was checked to fail with its addition switched off in the
  source, one at a time, and the source restored afterwards; the failures
  are quoted below under each addition.
- `agent/skills/kvit-ui/catalog.md` is rebuilt from the `KvitRow` summary,
  which now names both options. `names_gen.go` and `ux/tokens.css` did not
  need rebuilding, since no design value changed.
- No gallery specimen was added, as the request expected: neither option
  changes what a still picture of a row shows.
- The gallery's screenshot set was written from `c28b8b5` and from this
  working copy, 1,209 pictures each, and compared pixel by pixel on every
  colour channel. 13 differ: the `KvitRow` page in each of its four themes at
  three interface sizes, and in the macOS form. On each, the differing pixels
  lie inside the band of the page's summary (at the default interface size,
  rows 148 to 205 of a page 2,691 pixels tall), and every page kept its
  height, because the summary that names the two options still takes six
  lines. No other page changed.
- The tests of kvit-works, kvit-hub, kvit-cash and kvit-notes were run
  against kvit-ui at `c28b8b5` and against this working copy; "App tests"
  below has the result.
- The work is on this session's branch, whose base is `main` at `c28b8b5`,
  with no merge in between, so it moves onto `main` as a linear history. It is
  not pushed.

## Where the code is

| Addition | File | Test |
|---|---|---|
| `ListRow.FocusRequestEndsOnBlur` | `listrow.go` | `TestAListRowEndsAFocusRequestOnBlurWhenTold` in `listrow_test.go` |
| `ListRow.NoShiftActivate` | `listrow.go` | `TestAListRowLeavesShiftedKeysWhenTold` in `listrow_test.go` |
| The row's pointer watch registration, split so a test can give it a window | `listrow.go` (`watchPointerIn`) | `TestAListRowIsNotKeptByThePointerWatch` in `listrow_internal_test.go` |

The catalogue entry is in `cmd/kvit-ui-gallery/catalog.go`.

## Names

Both exported names are as the request wrote them. `FocusRequestEndsOnBlur`
sits beside `RingKeptOnPress` in the struct, since both decide when the ring
shows; `NoShiftActivate` sits after `OnKeyDown`, with the other key options.
The unexported function the test calls is `watchPointerIn(w *unison.Window)`.
Nothing exported was added beyond the two fields, and no `export_test.go` was
needed.

## Behaviour to compare

### Focus request ended by a loss of the focus

When code calls `Focus()` on a row that already holds the focus, unison
changes nothing, because `Window.SetFocus` does nothing when the target is
the current focus. The row records the request (`askedFocus`) and counts its
next gain of the focus as one from the keyboard, with the ring, however that
gain comes. The request calls such a request a waiting request.

- **What the option does.** In the `LostFocusCallback` that `initListRow`
  sets, after the control's own handling and before `OnBlur`, a row with
  `FocusRequestEndsOnBlur` clears `askedFocus`. Nothing else reads the
  option.
- **What it leaves alone.** A `Focus()` call that moves the focus to the row
  gains it at once, and the gain uses the request then, before any loss; the
  test checks this. A `Focus()` call made while the row could not take the
  focus leaves a waiting request on a row that does not hold the focus, which
  therefore cannot lose it, so that request still waits for the row's next
  gain. The case the request accepts, such a request followed by its
  sequence 3, still shows the ring on `ListRow`.
- **The sequences.** With the option on, a row that held the focus after a
  press, was asked for the keyboard, lost the focus and got it back from code
  shows no ring, whether or not it was pressed again before losing it
  (sequences 1 and 2 of the request); with the option off it shows the ring
  in both. Sequence 3 (the row's `OnActivate` calling `RequestFocus` after
  the focus left) is a gain from code while the window is not handling a key,
  the same path the test drives with `RequestFocus`, so it shows no ring with
  the option on.
- **`OnFocus`.** `OnFocus` is told `byPointer` as the opposite of what
  `KeyboardFocus` then answers. A gain from code after a forgotten request
  therefore reports `byPointer` true, as any gain from code without a request
  already does.

With the option's line removed from the source, the test failed with "with
FocusRequestEndsOnBlur true and pressed again true, the focus given back by
code shows the ring true", and the same for "pressed again false".

### Shift with the activation keys

- **What the option does.** In the `KeyDownCallback` that `initListRow`
  sets, after `OnKeyDown`, after `OnKey`, and after the check that the row is
  `Interactive` and holds the focus itself, a row with `NoShiftActivate`
  answers false for a key pressed with Shift, before the control's own
  handler. That handler acts only on Return, Enter and Space, so these are
  the only keys whose handling changes. A key answered false goes on to the
  panels the row sits in, as unison's window offers a key to the focused
  panel and then to each parent until one answers true.
- **The test's numbers.** Thirty interactive rows in a `Region` 200 pixels
  tall, scrolled to 366 pixels, with a row in view given the focus. A page is
  the larger of a `RowFull` row's height and the view's height less that
  height, 144 pixels here. With the option off, Shift+Space presses the row
  once and the region stays at 366; with it on, the row is not pressed and the
  region moves to 222. With it on, Shift+Return presses nothing, and Space
  without Shift presses the row once. With the option's lines removed from the
  source, the test failed with "with NoShiftActivate, Shift+Space pressed the
  row 1 times and moved the region from 366.0 to 366.0, want to 222.0", and
  Shift+Return pressed the row once.
- **Kvit Works' rule moved into `OnKeyDown`.** The quoted rule activates on
  Return, Enter or Space when no modifier in `mod.NonSticky` is held; in
  unison that mask is Shift, Control, Option and Command. With the rule in
  `OnKeyDown` and `NoShiftActivate` set, a key the rule leaves reaches the
  row's own handling only with one of those modifiers held. The row then
  answers false for Shift (the new option) and for Control, Option and
  Command (the control's handler has always refused them), so every key the
  rule leaves goes on to the parents, as with the copy. A held key and a key
  from a link or button inside the row are the rule's to take in
  `OnKeyDown`, which is offered both.

### Pointer watch

- **The split.** `watchPointer()` is now one line,
  `r.watchPointerIn(r.Window())`. `watchPointerIn` holds the checks and the
  registration unchanged: it does nothing without `HoverWhole`, without a
  window, or for the window the row is registered with already, and the
  listener holds the row only through `weak.Make`.
- **The test.** It is in a new internal test file, `listrow_internal_test.go`
  (package `kvitui`), because it calls the unexported `watchPointerIn`. That
  package cannot import `uitest`, which imports kvit-ui, so the test starts
  unison's headless screen and a window itself, as the tests in `ui_test.go`
  do. It makes a `ListRow` with `HoverWhole` holding a `Label`, never added to
  the window, and registers it with `watchPointerIn`. Before dropping the
  row it checks that the row's panel is among the window's pointer listeners,
  so the test cannot pass by registering nothing. It then moves the pointer
  once over the window and makes up to ten `runtime.GC()` calls; a
  `weak.Pointer` to the row must be nil.
- **Checked to fail.** With the listener changed to hold the row directly
  (the closure capturing the row rather than the weak pointer), the test
  failed with "a row nothing else holds is kept by the window's pointer
  watch". With the listener restored it passes.

## Cases not requested

Each was read in the source at `c28b8b5` and then run once in a throwaway
headless test, deleted afterwards.

### Ring while the button is held

This reads as the request describes when the press is what gives row A the
focus, and differently when A already holds the focus at the press.

`ListRow.mouseDown` sets `pointerFocus` and calls `RequestFocus`. The
control's `GainedFocusCallback` reads the marker and clears it, but unison
calls that callback only when the focus moves. When A already holds the
focus, for example because the window gave it to the first row when it
opened, or because A was pressed before, `RequestFocus` changes nothing and
`pointerFocus` stays set until the release. The gain at Shift+Tab in step 3
then reads it, and A shows no ring, as Kvit Works' `pressRow` does.

The throwaway test ran the request's sequence (press A and hold, Tab, then
Shift+Tab with the button still down) from two starting states:

| Before the press | A after the press | B after Tab | A after Shift+Tab |
|---|---|---|---|
| A holds the focus (the window gave it) | focus, no ring | focus, ring | focus, **no ring** |
| B holds the focus | focus, no ring | focus, ring | focus, **ring** |

So `ListRow` and the copy differ in step 3 only when the press moved the
focus to A. The request accepts `ListRow`'s behaviour in that case already,
and in the other case the two agree.

### Role of a row that cannot be pressed

This reads as the request describes. `ListRow.ProvideAccessibility` gives a
row that is not `Interactive` the role `role.Label` when it has a `Label` and
no control inside it, and `role.Group` otherwise; it decides the role only
while the node's role is still `role.Auto`. unison's
`Panel.Accessibility.Callback` "runs last, after everything else about the
node has been decided" (`accessibility.go` in unison v0.108.1). The throwaway
test read `role.Label` for a row with `Label` set, `Interactive` false and a
`Label` panel inside, and `role.Group` for the same row with the callback the
request quotes.

### Zero-width line in a two-line row

This reads as the request describes. `LineLayout.widths` in `line.go` gives
a growing child its own width when the content width it is given is 0, and
`PerformLayout` passes the content box's width. The throwaway test made a
`ListRow` 4 pixels wide at scale 1, holding a two-line content marked with
`Grow`, with a `LineLayout` as `Arrange`, `FitContent`, and `FitRoom` of
`SizeSpaceTight`: the content was framed at x = 2 with its own width of 111.7
pixels and a height of 34, in a row 36 pixels tall.

## Changes for existing callers

None. Both new fields are false unless a caller sets them, and the code they
guard is two `if` tests that do nothing while they are false. The pointer
watch's registration runs the same statements as before, through one more
function call.

The applications' tests agree; see "App tests".

## App tests

Each application's `go test -count=1 ./...` was run on the Ubuntu
development machine against kvit-ui at `c28b8b5` and against this working
copy. Each run used a temporary module file, a copy of the application's own
`go.mod` and `go.sum` with `replace` lines pointing kvit-ui and its patched
go-text copy at the version under test; `go list -m` confirmed the
replacement for both versions in all four applications, and `git status` in
each application's folder was the same before and after the runs. The four
applications ran at the same time, each running the `c28b8b5` version first.

The local checkouts pin kvit-ui at `bb47829` (kvit-works), `c28b8b5`
(kvit-hub), `2a4ecab` (kvit-cash) and `5faae05` (kvit-notes). The Kvit Works
checkout is therefore older than the `c28b8b5` the request names, and its
code does not set either new option, so this comparison tests what an
application that leaves them off sees.

| App | Tests failing in both runs | Failing only with this work | Failing only without it |
|---|---|---|---|
| kvit-works | 6 (and 4 of their subtests) | none | `TestStartPressedAgainAfterTheLinkWentMakesOneTrack`, `TestAnAgentsRequestForAProjectIsAnsweredAsItsTurn` |
| kvit-hub | none | none | none |
| kvit-cash | 4 | none | none |
| kvit-notes | 2 | `TestDiagramResetSceneDropsLastGood` | none |

Each test that failed in one run only was then run against both versions,
alone five times and in its whole package three times, and the kvit-notes
test alone 15 more times against each:

| Test | Against `c28b8b5` | Against this working copy |
|---|---|---|
| kvit-works `TestStartPressedAgainAfterTheLinkWentMakesOneTrack` (`ui/window/joinedtest`) | failed in 1 of 3 package runs, 0 of 5 alone | failed in 0 of 3 package runs, 1 of 5 alone |
| kvit-works `TestAnAgentsRequestForAProjectIsAnsweredAsItsTurn` (`agent/cluster`) | passed all 8 | passed all 8 |
| kvit-notes `TestDiagramResetSceneDropsLastGood` (`editor`) | failed in 0 of 3 package runs, 5 of 20 alone | failed in 0 of 3 package runs, 3 of 20 alone |

The first and third fail by chance against either version. The kvit-notes
test failed every time with the message "the reset kept the scene", and the
second report found it failing by chance in the same way. The second does not
depend on kvit-ui at all: `go list -deps ./agent/cluster` does not list it,
and the second report also found it failing by chance. None of the three is
caused by this work.

The tests failing in both runs, for the record: in kvit-works
`TestTheWindowsStatesFollowTheInterfaceRules` (with the subtests
`ask-menu-question`, `narrow` and `narrow-empty`), `TestWorksJourney` (with
the subtest `TheNewTrackButtonStartsATrack`),
`TestTheWorksShellLoadsAsARealAlternateComposition`,
`TestAConversationComesBackWhereItWasLeft`,
`TestALongConversationInAJoinedClientLoadsItsEarlierTurns` and
`TestSeveralFormsWorkOverTheJoint`; in kvit-cash
`TestFiltersFullScaleWorkspaceWithinBudget`,
`TestGeneratesFullScaleIndexedWorkspace`,
`TestHiddenCountAtFullScaleStaysWithinBudget` and
`TestEachPeriodIsTheCalendarPeriodHoldingToday`; in kvit-notes
`TestAMarginItemSitsBesideItsLine` and `TestTallInlineMathGrowsItsLine`.

## Left out

Nothing the request asked for was left out.
