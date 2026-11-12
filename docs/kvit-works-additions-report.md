# Additions for Kvit Works: what was built and where it differs from the request

*Written 2026-11-11 in kvit-ui, against the request
`docs/proposals/architecture-review-kvit-ui-request.md` of the Kvit Works
repository (commit `ba649218`), which asked kvit-ui at commit `bb47829` for
additions that let Kvit Works delete its private copies of layouts, text,
rows, tooltips and five smaller pieces.*

All five parts of the request are built as additions. This file is the report
the request asks for: the exported names where they differ, what was added
beyond the request, and each place where a quoted behaviour is not reproduced
exactly or where the request left a choice to kvit-ui. The reader on the Kvit
Works side should check each item under "Behaviour to compare" against its own
screenshots before deleting the matching copy.

## State

- `./build.sh --test` passes (gofmt, `go vet`, every headless test).
- The gallery's screenshot set was written before and after the change. Of the
  1,170 pictures that existed before, 1,105 differ only in the page count in
  the status bar ("92 components", formerly 89); the other 65 are the five
  pages that gained specimens (`KvitLabel`, `KvitRow`, `KvitTooltip`,
  `KvitPopover`, `KvitTextArea`). Three pages are new: `KvitStack`,
  `KvitLine`, `KvitFlowLayout`.
- The tests of kvit-works, kvit-hub, kvit-cash and kvit-notes were run twice
  each through a temporary module file, against kvit-ui at `bb47829` and
  against this working copy; "App tests" below has the result. No Kvit Works
  screenshot was compared.
- `agent/skills/kvit-ui/catalog.md` is rebuilt. `names_gen.go` and
  `ux/tokens.css` did not need rebuilding, since no design value changed.

## App tests

Each app's `go test ./...` was run against kvit-ui at `bb47829` and against
this working copy, on the Ubuntu development machine. The failures are the
same tests in both runs, with one exception in kvit-works that is not caused
by this work.

| App | Tests failing in both runs | Failing only with this work | Failing only without it |
|---|---|---|---|
| kvit-works | 7 | `TestAnAgentsRequestForAProjectIsAnsweredAsItsTurn` | `TestACommandEndsAtItsLimitAndWithItsTurn` |
| kvit-hub | 1 | none | none |
| kvit-cash | 4 | none | none |
| kvit-notes | 2 | none | none |

`TestAnAgentsRequestForAProjectIsAnsweredAsItsTurn` is in kvit-works'
`agent/cluster` package, which does not use kvit-ui. Run alone it failed
three times against this working copy and also failed against `bb47829`, so
it fails or passes for a reason outside kvit-ui. The test that failed only
without this work was not looked into.

## Where the code is

| Part | Files | Tests |
|---|---|---|
| Layouts (part 1) | `stack.go`, `line.go`, `morelayouts.go`, `textarea.go`, `field.go` | `layouts_test.go` |
| Text (part 2) | `label.go` | `label_test.go`, `label_internal_test.go` |
| Rows (part 3) | `listrow.go` | `listrow_test.go` |
| Tooltips (part 4) | `tooltip.go`, `control.go` | `tooltip_test.go` |
| Exported pieces (part 5) | `keyfocus.go`, `popover.go`, `roving.go`, `mover.go`, `pointer.go` | `pieces_test.go` |

## Names

Every exported name is as the request wrote it. Nothing was renamed.

Added beyond the request:

- `ListRow.NoHoverTintDisabled`, off by default: a disabled row is not tinted
  under the pointer. See "Candidates the request named".
- `Popover.ModalPress`, with `ModalPressStops` (the default) and
  `ModalPressPasses`. See "Exported pieces" below.
- `Tooltip.HoverTold` and `Tooltip.SetHovered(on bool)`. `AddTooltip` learns
  of the pointer from the panel's own enter and exit callbacks, which unison
  calls only for the innermost panel. A `ListRow` with `HoverWhole` is left
  and entered again as the pointer crosses onto its words, so a tooltip on
  such a row would close there. With `HoverTold` set the caller reports the
  pointer itself, typically from the row's `OnHover`.
- `Label.Layout(width float32) *text.Layout`, which the request asked for if
  cheap. It is the text before `MaxLines` cuts it short.

## Candidates the request named

- **`Label.Size`**: added. It is a `Measure`, so it follows the interface
  size, and the rule tests pass. A caller passing `Px(n)` is writing a size
  outside the seven roles, which is the caller's to justify.
- **`Label.Layout`**: added, as above.
- **The `Enabled` test on a row's hover tint**: not added to the row's default
  drawing, by the owner's decision on 2026-11-11, so that no existing caller
  changes: a disabled pressable row is tinted under the pointer as at
  `bb47829`. A Kvit Works row, which tints only an enabled row, sets
  `NoHoverTintDisabled`.
- **`TextArea.Fit`**: a field on `TextArea`, as asked. No `FitToText` function
  was needed. An error message under the area adds its height.

## Behaviour to compare

### Layouts

- **Layout data by value or pointer.** `OwnWidth`, `Grow` and the other
  markers set `StackData` and `LineData` by value; both layouts also read a
  pointer to either.
- **`StackLayout.SkipEmpty`.** Children left out are framed at no height at
  the bottom edge of the last child placed, without the gap after it. The
  request said "just past the last of them", which may mean after the gap.
- **`StackLayout.MaxWidth` with `Left` and `Right`.** The request quotes them
  from different copies. Together, `Left` and `Right` are taken off the width
  first and the limit applies to what remains.
- **`LineLayout.ClipFixed`.** The request quotes the clip only in the layout
  pass. Here the reported height also uses the clipped width, when a width is
  offered.
- **`LineLayout.GrowLast`.** The last shown child grows even if it is marked
  with `FixedWidth`.
- **`AtMost`.** The wrapped layout is asked at the offered width with no
  height, as the quoted copy asks its child, and minimum, preferred and
  maximum height are all set to the capped height.
- **`FlowLayout`.** With `FillsHint` or `Unasked` on a panel with a border,
  the insets are added to the reported width; the quoted copies appear to
  have no border. With `ItemWidth`, a line breaks only when a width is known.

### Text

- **`MaxLines`.** The two layouts also take the label's `LineHeight` and
  `Align`, which the quoted copy does not have; at their zero values the
  result is the copy's.
- **`Elided`** still answers false for a wrapped label, `MaxLines` included.

### Rows

- **`OnFocus(byPointer)`** is the opposite of `KeyboardFocus()`, as the
  request defines it, so it is also true for the focus a window hands out
  when it opens.
- **Slots.** Metrics chosen here: `Leading` and `Trailing` at their own size,
  cut to the row's height; a space (`SizeSpace`) between the parts; `Title`
  directly over `Subtitle` with no gap, the pair centred.
- **`FitContent` with no width offered** measures the children unconstrained.
- **A right-button press on a row with `AnyButton`** follows the quoted code;
  as the request notes, Kvit Works has no test of its own for it.

### Tooltips

- **A press with `Stays` false** hides the tooltip until the pointer leaves
  and returns, or the keyboard arrives. A control's own tooltip comes back on
  the next pointer move over the control; this difference only matters to a
  caller that leaves `Stays` off.
- **`TipFocusWatched` with no `Keyboard` function** asks the owner's
  `KeyboardFocus()` when the owner is a Kvit control, and `UI.KeyTurn()`
  otherwise.

### Exported pieces

- **`Popover.OpenModal`, a press outside.** `Popover.ModalPress` decides what
  becomes of the press. With `ModalPressStops`, the default, the popover
  closes (or `OnPressOutside` runs) once the press is over, so the scrim is
  still there to take it and nothing behind is pressed. With
  `ModalPressPasses` it closes at once and the same press goes on to what is
  behind. The second is what the quoted Kvit Works code does today: it shows
  a `Popup` with `Modal` and an `OnPressOutside` that closes at once, and
  closing removes the scrim before unison looks for the press's target. Moved
  callers set `ModalPressPasses` to keep that. In either mode, while a
  caller's `OnPressOutside` leaves the popover open, the scrim takes the
  press.
- **`Popover.PressOutsideIgnored` with `Open`** (not modal) leaves the press
  to whatever it lands on.
- **`DialogButtonOrder`** is the order the request quotes. `Dialog` itself
  orders differently on the macOS form (second answer, Cancel, Apply, then
  the action) and is unchanged.
- **`Mover.EscapeAbandons`** works only in a Kvit window, which is where the
  drag is registered for Escape.
- **`DrawSurface`** strokes the edge as kvit-ui does, with the radius less
  half the hairline, as the request expected.
- **`RoveStep`** is `roveStep` renamed, with `roveStep` kept as a one-line
  call to it.

## Left out

Nothing the request asked for was left out.
