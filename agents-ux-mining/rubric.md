# Failures the owner has flagged in the Kvit apps

This file lists the problems the owner of the Kvit apps has pointed out in their screens, collapsed into 31 recurring failures. Each entry is a rule a reviewer agent can hold a new screen against, with how often the owner raised it, the owner's own words, the cases it comes from and, where one exists, a before/after pair of pictures that the owner accepted. It is the "rubric of observed failures" that `agents-ux.md` (sections 1.6 and 3.5) proposes: published studies found that a list of failures already seen in the same product improves an agent's screens, whereas a list of general best practice does not.

The cases come from the session transcripts of seven Kvit projects, read by `scan.py` and checked by hand; `agents-ux-mining.md` says how. Every case is one line of `cases.jsonl` in this folder, and a case is named by its project, session and line, such as `works/20261022-f46b4c/L18`: line 18 of `~/kvit-works/.kvit/sessions/20261022-f46b4c/chat.md`. Kvit Works ("Works" below) runs coding agents; Kvit Hub ("Hub") is the owner's dashboard over many projects.

## The failures at a glance

223 of the 293 owner critiques found concern a screen; the rest are bug reports and questions. Each screen critique is counted under one failure. The last column says what could catch the failure before the owner does: a reviewer agent; a program (a test or check in the build); or one component that, fixed once, fixes every screen using it.

| Entry | Failure | Cases | Works | Hub | Other | Caught by |
|---|---|---:|---:|---:|---:|---|
| R01 | Result before progress | 9 | 9 | 0 | 0 | reviewer |
| R02 | Only what needs the user | 5 | 1 | 4 | 0 | reviewer |
| R03 | Each thing once | 8 | 6 | 2 | 0 | partly a program |
| R04 | Only the choices that apply | 6 | 6 | 0 | 0 | reviewer |
| R05 | Separate subjects, separate places | 7 | 4 | 3 | 0 | reviewer |
| R06 | Fits the desktop window | 5 | 2 | 3 | 0 | a program |
| R07 | Panes and selection | 5 | 5 | 0 | 0 | a program |
| R08 | Like items together | 7 | 3 | 4 | 0 | reviewer |
| R09 | Who said what | 4 | 4 | 0 | 0 | one component |
| R10 | No internal names | 10 | 6 | 4 | 0 | a program |
| R11 | Errors that say what to do | 15 | 12 | 2 | 1 | partly a program |
| R12 | Short plain labels | 15 | 6 | 8 | 1 | partly a program |
| R13 | Labels that agree with the screen | 9 | 5 | 4 | 0 | reviewer |
| R14 | Names from content | 2 | 2 | 0 | 0 | reviewer |
| R15 | Background work in view | 7 | 1 | 6 | 0 | partly a program |
| R16 | Visible result of every action | 12 | 8 | 4 | 0 | a program |
| R17 | Marks that lead somewhere | 6 | 6 | 0 | 0 | a program |
| R18 | Alarms only when true | 6 | 6 | 0 | 0 | reviewer |
| R19 | Disabled controls explain themselves | 4 | 3 | 0 | 1 | a program |
| R20 | A way forward from every state | 8 | 2 | 6 | 0 | partly a program |
| R21 | The thing being decided is on screen | 4 | 3 | 1 | 0 | reviewer |
| R22 | The app does not decide for the user | 4 | 2 | 2 | 0 | reviewer |
| R23 | Links open everywhere | 8 | 8 | 0 | 0 | one component |
| R24 | Desktop basics | 8 | 5 | 1 | 2 | one component |
| R25 | Actions keep working | 3 | 3 | 0 | 0 | partly a program |
| R26 | One text editor | 6 | 6 | 0 | 0 | one component |
| R27 | Rendered content | 10 | 8 | 2 | 0 | one component |
| R28 | Same thing, same look | 9 | 4 | 5 | 0 | partly a program |
| R29 | Like the product the user knows | 3 | 3 | 0 | 0 | reviewer |
| R30 | Size, spacing and alignment | 8 | 4 | 4 | 0 | partly a program |
| R31 | Drawn cleanly | 10 | 6 | 1 | 3 | one component |

By group:

| Group | Entries | Cases | Share |
|---|---|---:|---:|
| What a screen leads with | R01–R04 | 28 | 13% |
| How a screen is organised | R05–R09 | 28 | 13% |
| Words on screen | R10–R14 | 51 | 23% |
| State and response | R15–R19 | 35 | 16% |
| Reaching the next step | R20–R22 | 16 | 7% |
| Controls | R23–R26 | 25 | 11% |
| Drawing and likeness | R27–R31 | 40 | 18% |

By what could catch them: reviewer 62 cases (28%), partly a program 73 cases (33%), a program 42 cases (19%), one component 46 cases (21%).

## What a screen leads with

### R01 Result before progress

When a piece of work finishes, its result is what the screen shows first. The steps and history that led to it fold behind one summary row, and anything written to the user (a report, a summary before a question) is never folded away with that activity.

> "when I click on ship it opens up whole history - I don't care about it!" (`works/20261112-1ed198/L21`)
>
> "When an agent… produces a question… Write a report. And the report doesn't show… somehow the report gets swallowed. Here is an example." (`works/20261105-a111c5/L21`)

**Cases (9):** `works/20261009-f16de8/L21`, `works/20261012-b4e8ce/L21`, `works/20261022-f46b4c/L18`, `works/20261112-1ed198/L21`, `works/20261008-534a13/L21`, `works/20261102-2a9e06/L86`, `works/20261105-a111c5/L21`, `works/20261109-cd7739/L21`, `hub/20261109-f0c7ad/L32`.

**Example pair** (`works/20261022-f46b4c/L18`, owner's verdict: approved, "Show report start (Recommended)"): The agent's intermediate work before each report now folds into one row such as 'Worked for 1m 5s +1 −1 ›' with the line under it; a click opens the work.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261022-173832-f1577371.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261022-174700-b28c87fc.png`
- after: `~/kvit-works/.kvit/sessions/20261022-f46b4c/pictures/compare-dark-df05c8761342.png`
- after: `~/kvit-works/.kvit/sessions/20261022-f46b4c/pictures/running-turn-folds-at-its-report-14aa99dbb6df.png`

**Caught by:** A reviewer. A replayed workflow can capture the end of each task for the reviewer to look at.

### R02 Only what needs the user

A list of things to act on holds only items that need the user now. Handled, finished, outdated or unactionable items leave it or fold away at the end, while failures stay visible.

> ""the program that applies changes" - this is bad labeling… If they can be ignored, why are they shown to me at all?" (`hub/20261030-01364d/L75`)
>
> "job list has lots of stuff… lots of finished jobs. I usually don't care about finished jobs, so there probably should be some collapsed panel. Where the finished jobs are" (`hub/20261110-d93fa1/L21`)

**Cases (5):** `works/20261029-32cdd5/L77`, `hub/20261030-01364d/L75`, `hub/20261030-01364d/L120`, `hub/20261030-01364d/L240`, `hub/20261110-d93fa1/L21`.

**Example pair** (`hub/20261110-d93fa1/L21`, owner's verdict: approved, "Stop here"): Finished jobs now sit under an "Ended" heading that starts folded, and failed jobs that still need Run again or Dismiss have their own "Failed" group above it that never folds.

- before: `~/kvit-hub/.kvit/pasted/assets/message-20261110-152730-56b33691.png`
- after: `~/kvit-hub/.kvit/sessions/20261110-d93fa1/pictures/jobs-ended-folded-8bc008a66d3a.png`
- after: `~/kvit-hub/.kvit/sessions/20261110-d93fa1/pictures/jobs-ended-open-21ee62fcc2b6.png`

**Caught by:** A reviewer, helped by a replayed workflow that checks an item leaves the list once the user has acted on it.

### R03 Each thing once

Each piece of information appears once on a screen and each action has one control. There is no second view of what is already shown, no repeated button per row when every row's button opens the same thing, and a question answered in one place is cleared everywhere.

> "there is already stop button within composer box, there is no need for [Terminate agent] button" (`works/20261012-892b23/L77`)
>
> "This bar shows multiple machines buttons, but when I open each, that shows already multiple machines." (`works/20261105-fe5e31/L21`)

**Cases (8):** `works/20261012-892b23/L77`, `works/20261012-95b172/L21`, `works/20261022-3c414b/L21`, `works/20261112-fdf76f/L21`, `works/20261120-1eb872/L21`, `hub/20261112-d75530/L21`, `works/20261105-fe5e31/L21`, `hub/20261112-06bd47/L21`.

**Example pair** (`works/20261022-3c414b/L21`, owner's verdict: approved, "While a Push is running or stopped on conflicts, what should [x] do? — Close the view only"): Each step now starts with a status icon, the warning is a small orange pill with details on hover and press, the title and the state's buttons moved next to the × in the header, and a finished Push shows a summary with a collapsed file list and a centred primary Close button.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261022-165749-1746fb08.png`
- after: `~/kvit-works/.kvit/sessions/20261022-3c414b/pictures/push-conflict-dark-77ccb2fb293a.png`
- after: `~/kvit-works/.kvit/sessions/20261022-3c414b/pictures/push-ready-light-1a2af2bc641c.png`
- after: `~/kvit-works/.kvit/sessions/20261022-3c414b/pictures/push-running-dark-95a3bb7975d7.png`
- after: `~/kvit-works/.kvit/sessions/20261022-3c414b/pictures/operation-warning-light-b77e82d88d11.png`
- after: `~/kvit-works/.kvit/sessions/20261022-3c414b/pictures/push-complete-dark-5cd5b85f3841.png`
- after: `~/kvit-works/.kvit/sessions/20261022-3c414b/pictures/operation-receipt-open-light-40fb171da878.png`

**Caught by:** Partly a program: two visible controls on one screen that run the same command can be found from the screen-reader tree. A reviewer for repeated information.

### R04 Only the choices that apply

A card, dialog or page offers only the actions and choices that mean something for the item in its current state. Actions the user rarely takes are folded away. A screen put to the user for approval is shown in the state they will normally see.

> "What this keep this track? … I always let agents to do the work. And all I see is hinge buttons, and I don't see information that agents do…" (`works/20261117-1cb519/L78`)
>
> "home track is now actually main, there is no point in push/pull button on it anymore:" (`works/20261001-019af4/L20`)

**Cases (6):** `works/20261117-1cb519/L21`, `works/20261117-1cb519/L78`, `works/20260921-4a956a/L19`, `works/20260921-9e9e27/L20`, `works/20261001-019af4/L20`, `works/20261107-50170f/L21`.

**Example pair** (`works/20261117-1cb519/L21`, owner's verdict: corrected, "Shorten Keep label"): Added a turning mark and 'Nothing to press while it works' under the agent's sentence, and folded the file list and Keep/Cancel buttons under a closed 'Other choices'.

- before: `~/kvit-works/.kvit/pasted/assets/phone-20261117-182026-3a6b6f1b.png`
- after: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/after-05d553e6a96f.png`
- after: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/after-open-61b55697a61c.png`

**Caught by:** A reviewer.

## How a screen is organised

### R05 Separate subjects, separate places

Settings for a separate subject get their own tab, a settings page shows only settings, and the output of running work appears where that work is started or followed. A setting lives where its scope is: per project with the project, once for everything in one general place.

> "what is going on with log being displayed in settings view? … that is setting! it should show settings! whatever is running does not go there!" (`hub/20261102-ef8d24/L21`)
>
> "should these be on project settings page, not on general?" (`works/20261013-13d012/L21`)

**Cases (7):** `works/20261031-21d23f/L47`, `works/20261110-6b1d27/L97`, `hub/20261031-ad2899/L21`, `hub/20261102-ef8d24/L21`, `works/20261013-13d012/L21`, `works/20261105-7096bc/L39`, `hub/20261103-5ab638/L21`.

**Example pair** (`hub/20261102-ef8d24/L21`, owner's verdict: approved, "Write jobs design (Recommended)"): The running block and its output were taken off the Agents tab, which now shows only settings, plus a notice when an agent is logged out or out of credits.

- before: `~/kvit-hub/.kvit/pasted/assets/message-20261102-122928-b8f62430.png`
- after: `~/kvit-hub/.kvit/sessions/20261102-ef8d24/pictures/settings-agents-after-bee9f6f0fdfb.png`

**Caught by:** A reviewer.

### R06 Fits the desktop window

A screen bigger than its window is divided into panes that each scroll with a visible scroll bar. Content wider than its pane can be reached. A dialog with a list of unknown length keeps its buttons on screen, and a long menu has a maximum height and scrolls.

> "This view feels like a webpage… there is no scroll bar on the side… doesn't feel like. Desktop app… maybe splitting into multiple views makes more sense." (`hub/20261123-9d0b02/L21`)
>
> "when I on main press commit I see: … I suppose there is button at the bottom that I don't see and cannot press because it is clipped." (`works/20261120-20f68e/L44`)

**Cases (5):** `works/20261110-c9ea0a/L226`, `works/20261120-20f68e/L44`, `hub/20261030-01364d/L466`, `hub/20261110-99bad7/L21`, `hub/20261123-9d0b02/L21`.

**Example pair** (`hub/20261123-9d0b02/L21`, owner's verdict: approved, "Build C (Recommended)"): The agent found that the Projects tabs' scroll bars were drawn outside the visible area and proposed three layouts; the recommended one puts the board, Your work today and the dated list in separate panes that each scroll with a visible bar.

- before: `~/kvit-hub/.kvit/pasted/assets/capture-20261124-043045-06b009b6.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/check-projects-board-343c7e345d7b.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/c-panes-aorus-62ccbb73ce64.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/c-panes-7d4621ca1511.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/a-board-77a509741966.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/check-projects-board-343c7e345d7b.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/c-panes-aorus-62ccbb73ce64.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/c-panes-7d4621ca1511.png`
- after: `~/kvit-hub/.kvit/sessions/20261123-9d0b02/pictures/a-board-77a509741966.png`

**Caught by:** A program, built: `Session.CheckReach` in kvit-ui's `uitest` shrinks the window to its smallest size and fails a control, scroll view, dialog or menu cut off where nothing scrolls it into view (rule LAY-3). It found that a kvit-ui dialog holding a long list grew past the window; a dialog now keeps inside the window and scrolls its content, with its title and buttons in place. kvit-ui menus stop at fifteen lines and scroll. A page scrolled as a whole instead of in panes needs a reviewer.

### R07 Panes and selection

Whatever fills the main pane has a visibly selected row in the navigation. A pane is named once, in its own title bar, with its actions and its close button at the right end of that bar, and closing it clears the selection.

> "whatever is shown in right view should be selected in left detail view - I don't see it selected anywhere. Also when I switch to another track, history view… stays" (`works/20261012-10817c/L21`)
>
> "also - close button needs to be in the title panel:" (`works/20260921-4a956a/L40`)

**Cases (5):** `works/20261012-10817c/L21`, `works/20261012-c69de3/L21`, `works/20260921-4a956a/L40`, `works/20260921-4a956a/L62`, `works/20261013-fd890a/L21`.

**Example pair** (`works/20261012-c69de3/L21`, owner's verdict: approved, "The × is in and its tests pass. What next? — Stop here"): A × at the right end of the right side's title bar now closes that side's views for the track and leaves nothing selected on the left.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261012-210900-1b56a681.png`
- after: `~/kvit-works/.kvit/sessions/20261012-c69de3/pictures/stage-close-after-2d7693ed2865.png`

**Caught by:** Mostly a program: from the screen-reader tree, every pane showing an item has a title with a close button, and some row in the navigation is selected.

### R08 Like items together

Items of one kind appear in one list, each with its state and its own action, rather than split between a list and a dropdown. Stages that run in order are numbered steps. A long list is grouped (by month, by folder) and a list that mixes kinds marks each kind. A group heading is visibly tied to the items under it.

> "this whole thing is confusing: … What is the difference between the ones that down below, with the ones that in… drop down. It's not clear and it's confusing." (`hub/20261110-7caecc/L21`)
>
> "this list is important… it is hard to parse… more reading than viewing… maybe we should group by month… I'm not clear what what is shown in the third column." (`hub/20261123-841d3f/L21`)

**Cases (7):** `works/20261107-56a8dc/L63`, `hub/20261102-ef8d24/L40`, `hub/20261110-7caecc/L21`, `hub/20261123-841d3f/L21`, `hub/20261124-d6506d/L61`, `works/20261008-93f42a/L21`, `works/20261105-3b294f/L21`.

**Example pair** (`hub/20261110-7caecc/L21`, owner's verdict: approved, "Stop here"): The dropdown was removed and the card shows one row per project still to plan, each with its state (Not started, Working, Waits for you, Ready on Decisions, Failed) and its own Plan, Plan tasks or Open button.

- before: `~/kvit-hub/.kvit/pasted/assets/message-20261110-152303-9188751b.png`
- after: `~/kvit-hub/.kvit/sessions/20261110-7caecc/pictures/planning-card-six-dd6dba6ff5ba.png`
- after: `~/kvit-hub/.kvit/sessions/20261110-7caecc/pictures/recap-planning-frozen-copy-8886b894f82c.png`

**Caught by:** A reviewer.

### R09 Who said what

Messages the app sends to an agent by itself are drawn as part of the agent's work, never as if the user wrote them. Everything the user said, including answers through a question card, looks like the user's own messages, and an agent's progress looks the same whichever agent produced it.

> "instead of showing original report, then something like looks like user input (which actually automated kvit-works message to copilot) and then updated report, instead should be within chat tool calls…" (`works/20261120-20f68e/L21`)
>
> "when model sends me question and I respond this shows on history of chat like: … this is barely visible - it should show with background similar to regular request" (`works/20261003-2ca34c/L21`)

**Cases (4):** `works/20261003-2ca34c/L21`, `works/20261120-20f68e/L21`, `works/20261023-3a7881/L21`, `hub/20261120-b73cf0/L18`.

**Example pair** (`works/20261003-2ca34c/L21`, owner's verdict: approved, "Answers now have your background; the text is still the quotation's dimmer grey. What next? — Stop here (Recommended)"): An answer to a question is now drawn on the same grey background, with the owner's mark, as the owner's own messages, in finished and running turns; its text stays the quotation's dimmer grey.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261003-011505.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261003-011606.png`
- after: `~/kvit-works/.kvit/sessions/20261003-2ca34c/pictures/transcript-answer-dark-bb7ba7c2aea2.png`
- after: `~/kvit-works/.kvit/sessions/20261003-2ca34c/pictures/reply-so-far-answer-93abafec5370.png`

**Caught by:** Fixed once in the conversation component of Kvit Works, which draws every agent's transcript.

## Words on screen

### R10 No internal names

On-screen text uses the words the user would use. The code's own names (surface, target, stage, records, modes), raw markers such as [dictated:, model IDs, and descriptions of how the system works inside never appear.

> "I just got a message: … surface is your jargon that user has no clue about. 'View' is what user understands." (`works/20261022-691b03/L21`)
>
> "what the heck that phrase even means? Do we really need it? It looks like some internal jargon mambo-jumbo" (`hub/20261030-01364d/L193`)

**Cases (10):** `works/20261022-691b03/L21`, `works/20261022-691b03/L65`, `hub/20261030-01364d/L109`, `hub/20261030-01364d/L193`, `hub/20261031-4756e1/L21`, `hub/20261108-8e2513/L41`, `works/20261012-d73f2b/L21`, `works/20261023-76239e/L21`, `works/20261023-ad6bae/L21`, `works/20261123-84f488/L21`.

**Example pair** (`works/20261022-691b03/L65`, owner's verdict: approved, "Stop here"): About 100 on-screen strings in 62 files were reworded (for example 'Put the work on main', 'Pause after this step', 'project history'), and notices now also appear on the status bar for a few seconds.

- before: `~/kvit-works/.kvit/sessions/20261022-691b03/pictures/ambiguous-0-unavailable-f8e5f18ad323.png`
- after: `~/kvit-works/.kvit/sessions/20261022-691b03/pictures/push-running-dark-95a3bb7975d7.png`
- after: `~/kvit-works/.kvit/sessions/20261022-691b03/pictures/push-running-dark-fc0a6a712988.png`

**Caught by:** A program: each app keeps a list of its internal words, and the existing forbidden-wording check (CheckRules in source, CheckLabels on screens) fails when one appears in visible text.

### R11 Errors that say what to do

A message that refuses an action or reports a failure names the actual cause in the user's terms and offers a button that fixes it. It never shows a raw network error or a bare exit code, never reports a failure that did not happen, names the one cause that applies rather than a list of possible ones, and stays explained after a restart.

> "what does it mean 'cannot be answered here'. 1. This is bad message, it should have suggestion what can be done to resolve, instead of just complaining." (`hub/20261030-01364d/L21`)
>
> "Something happened when I pushed, and it's not clear how to resolve. And it's message completely not clear. What what does it mean?" (`works/20261105-3e40d7/L21`)

**Cases (15):** `works/20261012-b84545/L21`, `works/20261112-10baab/L21`, `works/20261116-fa808e/L21`, `works/20261123-9c0915/L21`, `hub/20261030-01364d/L21`, `works/20261001-0a7353/L21`, `works/20261022-12ac19/L21`, `works/20261103-2b00f8/L21`, `works/20261103-2b00f8/L56`, `works/20261104-64c2e5/L21`, `works/20261105-3e40d7/L21`, `works/20261107-86a63b/L21`, `works/20261117-cb703d/L21`, `hub/20261108-c77622/L21`, `dictate/20261120-954120/L195`.

**Example pair** (`works/20261116-fa808e/L21`, owner's verdict: corrected, "Find the root cause"): The refusal now names the blocking conversation, marks it if closed, and offers a Stop or End turn button; pressing Push again re-checks, and Stop clears a stale running mark.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261116-150714-903de849.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261116-150737-7c36ddf5.png`
- after: `~/kvit-works/.kvit/sessions/20261116-fa808e/pictures/operation-blockers-light-6aad21aa8efd.png`

**Caught by:** Partly a program: visible text checked for raw error patterns (exit code, connection refused, stack traces), and every error notice in the screen-reader tree has an action next to it. A reviewer for whether the cause is the right one.

### R12 Short plain labels

Labels are short ordinary words. A button is a verb for the outcome the user wants, without counts, dates or names. A tooltip is a few words, a confirmation one sentence, and on-screen text is an instruction to the user rather than a narration of what the program does. Empty sections and obvious empty states are left out, and figures are short numbers in the largest sensible unit (18 h, 4.2M), with the explanation in a hover tip.

> "improve robotic comments: Counts ach day's activity. No agent; about a minute a day. want to puke to this! …Should have been: Collect each day's activity." (`hub/20261102-ef8d24/L160`)
>
> "I don't like this. Lots of reading. 18 hours of yours. Six weeks left. Can you make it short? Something like 18h slash six weeks. Maybe when I hover…" (`hub/20261123-3f4dc0/L110`)

**Cases (15):** `works/20261115-67f1ff/L21`, `hub/20261031-7729b2/L525`, `hub/20261102-ef8d24/L160`, `hub/20261122-d13081/L21`, `hub/20261123-3f4dc0/L110`, `hub/20261123-8e010f/L21`, `hub/20261124-d6506d/L21`, `works/20260920-50df12/L19`, `works/20260930-a6bd7f/L20`, `works/20261022-776ce8/L21`, `works/20261023-98fc02/L21`, `works/20261023-e5f54d/L21`, `hub/20261031-7729b2/L611`, `hub/20261107-13ab3d/L46`, `dictate/20261123-8ffe92/L127`.

**Example pair** (`hub/20261031-7729b2/L525`, owner's verdict: corrected, "this looks like improperly alligned. Measure (and other relevant) buttons shouls have dropdown attached with [All], [5 days]"): The agent drew a new panel with one verb-only button per stage (Measure all, Write up, Gather), the count beside the stage name, a short line under each name instead of an information icon, and a progress bar with Stop while a batch runs.

- before: `~/kvit-hub/.kvit/sessions/20261031-7729b2/pictures/message-20261101-133934-d2a44400-d2a44400397e.png`
- after: `~/kvit-hub/.kvit/sessions/20261031-7729b2/pictures/mock-catch-up-cd4a0bf68609.png`
- after: `~/kvit-hub/.kvit/sessions/20261031-7729b2/pictures/mock-project-describe-9b061356556b.png`

**Caught by:** Partly a program: CheckLabels can limit the words in a button or tooltip and reject digits in button labels and raw large numbers. A reviewer for tone; the owner's own phrase for it is "LLM looking words".

### R13 Labels that agree with the screen

A button says exactly what it will do and matches what the screen shows. A heading agrees with the states listed under it, a count of the same thing is the same everywhere, a name means one thing, two things that share a name can be told apart, and a status word gives the specific outcome (accepted, not decided).

> "labels should not be: 18 are outdated, 24 are routine… as a user I'm not sure what it means… what does 'Close these 42' means?" (`hub/20261030-01364d/L367`)
>
> "What's going on here? It shows some files were changed, but then it says the agent didn't run. What what's happening?" (`works/20261106-9f871c/L21`)

**Cases (9):** `works/20261011-819701/L174`, `hub/20261030-01364d/L367`, `hub/20261102-ef8d24/L34`, `works/20261012-19feaf/L21`, `works/20261105-888c08/L21`, `works/20261106-9f871c/L21`, `works/20261108-e1f200/L21`, `works/20261112-bdb4ed/L18`, `hub/20261030-01364d/L530`.

**Example pair** (`hub/20261030-01364d/L367`, owner's verdict: approved, "Ship it now (Recommended)"): The bar now says '42 of 44 old decisions can be closed without you', shows the agent's summary and what Close does, links to the full review, and the group is renamed 'Can be closed without you'.

- before: `~/kvit-hub/.kvit/sessions/20261030-01364d/pictures/message-20261030-111716-9b7c1565-9b7c15657de5.png`
- after: `~/kvit-hub/.kvit/sessions/20261030-01364d/pictures/decisions-checked-2-0a1e6348c249.png`

**Caught by:** A reviewer, given the whole workflow rather than one screen; counts that must match can be asserted in a replayed workflow.

### R14 Names from content

An item the app names by itself is named after what it contains, with a short summary, not "Untitled 3" or the first characters of the user's text.

> "when I create an agent it is given untitled name unless it is main agent of a track. Can it be automatically renamed once chat starts" (`works/20261012-396f18/L21`)
>
> "it gets renamed to the initial part of my input. Does it run summarizer to get better name, shorter and better name? … it doesn't look like it's doing that." (`works/20261031-21d23f/L21`)

**Cases (2):** `works/20261012-396f18/L21`, `works/20261031-21d23f/L21`.

**Example pair** (`works/20261012-396f18/L21`, owner's verdict: corrected, "Agent naming works and is off until you pick kc… What next? — Name old agents (Recommended)"): An unnamed agent is now shown under the opening words of its first message and then renamed with a short title written by a summarising agent or model chosen in a new Summaries setting.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261012-161037-73b41e8b.png`
- after: `~/kvit-works/.kvit/sessions/20261012-396f18/pictures/after-first-message-0eacec7632b9.png`
- after: `~/kvit-works/.kvit/sessions/20261012-396f18/pictures/2-first-message-2b8ad1e12a21.png`
- after: `~/kvit-works/.kvit/sessions/20261012-396f18/pictures/3-titled-0043ae23f5cd.png`
- after: `~/kvit-works/.kvit/sessions/20261012-396f18/pictures/summaries-agent-dark-f1bb1cbb0600.png`
- after: `~/kvit-works/.kvit/sessions/20261012-396f18/pictures/summaries-api-light-efa2655da890.png`

**Caught by:** A reviewer.

## State and response

### R15 Background work in view

While work runs in the background, the screen says what is running, how far it has got, what it is waiting for and whether the user needs to do anything, as its most prominent line. Starting long work does not block the screen, and closing the window while work runs warns first.

> "This is bad UI. Waiting to save answers when I already pressed save. Waiting for whom? … Do I, as a user, need to do anything about it?" (`hub/20261111-733183/L21`)
>
> "it shows that's resolving with agent, but small font somewhere in the middle of the text and a wall of text that I'm not going to read. All these buttons…" (`works/20261117-1cb519/L104`)

**Cases (7):** `works/20261117-1cb519/L104`, `hub/20261031-7729b2/L21`, `hub/20261109-409af9/L21`, `hub/20261111-733183/L21`, `hub/20261111-733183/L28`, `hub/20261030-01364d/L548`, `hub/20261108-661157/L21`.

**Example pair** (`works/20261117-1cb519/L104`, owner's verdict: approved, "Stop here"): The page now opens with one bold, larger line 'Push conflict is resolving the conflict' with a turning mark and 'Nothing to press while it works'; the repeated grey status line is gone and the rest stays folded.

- before: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/agent-working-05d553e6a96f.png`
- before: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/conflict-before-42ee4181c947.png`
- before: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/conflict-folded-3d034da48d9a.png`
- after: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/phone-agent-v4-74defb46d2ee.png`
- after: `~/kvit-works/.kvit/sessions/20261117-1cb519/pictures/desktop-agent-5e2e7ba30d66.png`

**Caught by:** Partly a program: a replayed workflow can check the screen stays usable while a job runs. A reviewer for what the status line says.

### R16 Visible result of every action

After the user acts, the screen shows the result where they are looking: an answered item leaves the to-do list at once, a sent message shows whether it is waiting, sent or read, a toggle shows its state, Copy confirms, a step shows what it produced, and an action whose effect happens elsewhere says plainly whether it arrived.

> "there is [Copy] button on top right. When I click it has no feedback at all (should show short popup - 'copied' or something)…" (`works/20261022-35671b/L21`)
>
> "when I run ship to windows button it does something but it doesn't show if actual shipping happened. Or am I doing something wrong?" (`works/20261012-d1df08/L21`)

**Cases (12):** `works/20261031-895b37/L21`, `hub/20261111-733183/L32`, `hub/20261112-c2063a/L21`, `hub/20261122-8be1e0/L111`, `works/20261008-99fa3a/L21`, `works/20261012-d1df08/L21`, `works/20261012-d1df08/L55`, `works/20261022-35671b/L21`, `works/20261028-e1366c/L21`, `works/20261117-9dadb0/L21`, `hub/20261110-aa218f/L21`, `cash/20261012-6b628a/L21`.

**Example pair** (`hub/20261112-c2063a/L21`, owner's verdict: approved, "Stop here"): While the pane is open the Ask button is drawn as pressed (outlined and tinted), its tooltip reads 'Hide the Ask pane', and pressing it again hides the pane.

- before: `~/kvit-hub/.kvit/pasted/assets/message-20261112-032343-11b4b939.png`
- after: `~/kvit-hub/.kvit/sessions/20261112-c2063a/pictures/ask-button-before-after-e72768cd5f8f.png`
- after: `~/kvit-hub/.kvit/sessions/20261112-c2063a/pictures/ask-pane-open-after-9ad003452817.png`

**Caught by:** A program, built: `Session.CheckResponds` in kvit-ui's `uitest` does one step of a scripted task and fails it when nothing a screen reader is told about the window changed. Whether the change is the right one, and in the right place, needs a reviewer.

### R17 Marks that lead somewhere

A summary mark (a dot on a project tile, a badge on a tab) uses the same rule as the marks on the items it summarises, appears only for something that waits for the user, shows on hover what it stands for, and following it leads step by step to the exact view that needs the user.

> "is this kind of dot considered a warning for archived … Only something that is waiting explicitely for user should be shown with that dot" (`works/20261008-583e14/L51`)
>
> "Why does it show error on project's tile in rail, even though there is no error shown within project view:" (`works/20261012-506c84/L21`)

**Cases (6):** `works/20261001-bf469d/L20`, `works/20261008-583e14/L93`, `works/20261117-014fec/L21`, `works/20261008-583e14/L51`, `works/20261012-506c84/L21`, `works/20261023-dd4ae8/L21`.

**Example pair** (`works/20261008-583e14/L93`, owner's verdict: corrected, "The rail's dot now matches the track dots… What next? — Count claude.ai prompts (Recommended)"): The rail tile now shows the blue waiting dot whenever any of the project's tracks shows one, using the same rule as the track list, and pressing the dot opens the first waiting track's conversation.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261008-012047-d59b5820.png`
- after: `~/kvit-works/.kvit/sessions/20261008-583e14/pictures/rail-waiting-dark-c489ad81202b.png`

**Caught by:** Mostly a program: a test in each app that the summary mark is computed from the item marks; a replayed workflow that follows a mark to its item.

### R18 Alarms only when true

Warning colours, error marks and warnings about losing work appear only when they are true. An outcome the user chose (closing, stopping) is never an error, "nothing to do" is not a warning, a reading that is missing or out of date says so, and a label beside a setting describes the current choice.

> "When I close an agent (terminate it) it shows like it is an error. It is not, I did it intentionally and it is not an error!" (`works/20261013-4f7d62/L21`)
>
> "why 1% limit is shown in yellow/orange? Feels it is leftover from before credit reset:" (`works/20261027-0b402c/L21`)

**Cases (6):** `works/20260923-9bf494/L19`, `works/20261013-4f7d62/L21`, `works/20261027-0b402c/L21`, `works/20261106-0f3ce9/L32`, `works/20261106-a71a8b/L21`, `ui/20261112-2260b1/L21`.

**Example pair:** none; no case of this failure has an accepted before/after pair.

**Caught by:** A reviewer.

### R19 Disabled controls explain themselves

A disabled control says on screen, without hovering, why it is disabled, and when the main action does not apply the screen offers the action that does (Archive when there is nothing to push). A fresh install does not leave a feature greyed out.

> "some tracks end up not changing any files so [Push + Archive] is disabled. We should have button [Archive] to be able to close the track…" (`works/20261022-ee7090/L21`)
>
> "the Pusher guy was not available … push is not available because another track being pushed probably should have a message there under when it's grayed out" (`works/20261112-38e82d/L30`)

**Cases (4):** `works/20261022-ee7090/L21`, `works/20261112-38e82d/L30`, `works/20261106-0f3ce9/L21`, `dictate/20261120-954120/L21`.

**Example pair** (`works/20261022-ee7090/L21`, owner's verdict: corrected, "Archive without asking"): When a track has nothing to push the button reads 'Archive' and archives it (after the same confirmation as the other archive control), switching back to 'Push + Archive' once the track has changes.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261022-234015-bc88fd53.png`
- after: `~/kvit-works/.kvit/sessions/20261022-ee7090/pictures/nothing-to-push-dark-42ed1db00c11.png`
- after: `~/kvit-works/.kvit/sessions/20261022-ee7090/pictures/something-to-push-dark-22c2e4ce310d.png`

**Caught by:** A program, existing for buttons: a disabled kvit-ui Button shows its Explanation, and `uitest.Labels` (Session.CheckLabels) fails one without (rule CTL-8). The cases here come from screens that did not use it or did not call the check; other disabled controls are not checked.

## Reaching the next step

### R20 A way forward from every state

Every held, failed or skipped item has a button on its own row that resumes it or runs it again, and its error stays visible until it succeeds. The action that starts the whole job stays beside a failure, a conversation pane has a way to start a new conversation, finished work offers the next steps, and a dashboard answers how much is done and what to do next.

> "…all those buttons on the right… I just want to continue that job. And it says held, and I have no idea how to get it back running." (`hub/20261112-ab6264/L21`)
>
> "I'm always behind … I need to understand like how much is done, how far it is, what is the potential of the project … what should I do next?" (`hub/20261122-710388/L114`)

**Cases (8):** `hub/20261112-2f9a54/L106`, `hub/20261112-768202/L21`, `hub/20261112-ab6264/L21`, `hub/20261122-710388/L114`, `works/20261025-7b9c12/L21`, `works/20261107-194b4d/L21`, `hub/20261106-2ef185/L21`, `hub/20261109-f0c7ad/L21`.

**Example pair** (`hub/20261112-ab6264/L21`, owner's verdict: corrected, "Build the review (Recommended) — the agent had said the rest of the page was still confusing"): Each Held row got Continue and Remove, Start all and Remove all left the page head, the step's Start became Continue, and a job cut off by a restart is listed as cut off and resumed first instead of shown as failed.

- before: `~/kvit-hub/.kvit/pasted/assets/message-20261112-092752-d530cb35.png`
- after: `~/kvit-hub/.kvit/sessions/20261112-ab6264/pictures/jobs-after-restart-fixed-c9cdc1503f5f.png`
- after: `~/kvit-hub/.kvit/sessions/20261112-ab6264/pictures/jobs-after-restart-fixed-c9cdc1503f5f.png`

**Caught by:** Partly a program: in a replayed workflow that makes a job fail, the failed row has an enabled action. A reviewer for the rest.

### R21 The thing being decided is on screen

When the user is asked to approve or choose, the thing itself (the document, the pictures, the command) is on screen beside the question, and a message that refers to a picture, file or command shows it.

> "I don't see a button to take a look at description … we need a link or something so that I can open description before I can accept" (`hub/20261108-58a596/L21`)
>
> "can you give me report where I can see the pictures? I don't see them in your responses" (`works/20261106-802499/L36`)

**Cases (4):** `works/20261106-802499/L36`, `works/20261106-f95ceb/L21`, `hub/20261108-58a596/L21`, `dictate/20261120-954120/L53`.

**Example pair:** none; no case of this failure has an accepted before/after pair.

**Caught by:** A reviewer.

### R22 The app does not decide for the user

The app does not apply a change on the user's behalf without the user confirming it, and does not ask again what the user has already decided. A form fills in what the app can work out, and an empty field's hint says what leaving it empty does rather than showing an example that looks like a default.

> "I intentionally kept them separate. I already kind of made the decision about this. And if algorithmically constantly gonna be suggesting me… this is useless to me." (`hub/20261122-185466/L21`)
>
> "I'm trying to set it up and I see: [picture] what should I use for sessions folder? why is it not populated automatically?" (`works/20261109-bdb78e/L41`)

**Cases (4):** `works/20261110-0be235/L21`, `hub/20261122-185466/L21`, `works/20261109-bdb78e/L41`, `hub/20261122-8be1e0/L21`.

**Example pair** (`hub/20261122-185466/L21`, owner's verdict: approved, "Asked whether connections between projects are recorded anywhere, then after the answer: "Stop here""): Projects already in the owner's list are never proposed as part of another; a pair worked on together now gets an information line such as 'Used with home-auto on 7 days' on the project's card, and nothing waits for an answer.

- before: `~/kvit-hub/.kvit/pasted/assets/message-20261122-113543-50e6accf.png`
- after: `~/kvit-hub/.kvit/sessions/20261122-185466/pictures/map-suggestions-after-00383de3276a.png`
- after: `~/kvit-hub/.kvit/sessions/20261122-185466/pictures/map-card-eero2beryl-after-aceb883cf960.png`

**Caught by:** A reviewer.

## Controls

### R23 Links open everywhere

Every file name, path, picture reference or link shown as text opens the same way (Ctrl+click on the desktop, a tap on the phone) wherever it appears: replies, question cards, message boxes, settings, Markdown files. It shows that it can be clicked, and a name that refers to something is itself the link.

> "I has file name in question tool, but I cannot Ctrl+click on it to open - can you fix this?" (`works/20261022-ef9a1f/L22`)
>
> "crap - can you just include file path within your response how you usally do - not full path, relative to your root! … none of these are clickable!" (`works/20261102-2a9e06/L103`)

**Cases (8):** `works/20261002-b946e0/L21`, `works/20261022-ef9a1f/L22`, `works/20261110-b115b0/L21`, `works/20261112-fcb7b2/L21`, `works/20261117-3edb7c/L21`, `works/20261120-a35501/L21`, `works/20261023-85c9ee/L21`, `works/20261102-2a9e06/L103`.

**Example pair** (`works/20261112-fcb7b2/L21`, owner's verdict: approved, "Stop here"): Ctrl+clicking a picture path in a chat or Markdown file now underlines it on hover and opens it in the picture viewer.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261112-083239-0b94d5b3.png`
- after: `~/kvit-works/.kvit/sessions/20261112-fcb7b2/pictures/named-picture-hover-c614eea5979d.png`
- after: `~/kvit-works/.kvit/sessions/20261112-fcb7b2/pictures/named-picture-viewer-c2915c6d44f8.png`

**Caught by:** Fixed once in the component that draws text, so that every place showing text finds paths and links by the same rule.

### R24 Desktop basics

Text the user may need, such as logs and output, can be selected and copied. Pressing the control that opened a dropdown closes it, Enter presses the primary button, a hover card opens anywhere over its control, a tooltip never covers the control it describes, every overlay and panel has a close button and Esc, and a picture viewer opens at 100% with a way to fit.

> "when I press the chevron… I press chevron again. It closes and opens again. My expectation is that when I press the chevron, it would close." (`works/20261008-587061/L21`)
>
> "this text should be selectable: … this is details in jobs view" (`hub/20261112-36dd4a/L21`)

**Cases (8):** `works/20261008-4e1fc6/L21`, `works/20261008-587061/L21`, `works/20261022-5c4614/L21`, `works/20261123-fab1a6/L21`, `hub/20261112-36dd4a/L21`, `cash/20261011-00c960/L21`, `works/20261029-aaca4c/L21`, `cash/20261012-6b628a/L55`.

**Example pair** (`works/20261008-587061/L21`, owner's verdict: approved, "How should I proceed? — Push and update app (Recommended)"): In the component library, a press on the control that opened a menu or dropdown now only closes it, so the chevron toggles the list; this applies to every dropdown and menu button.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261008-222547-2b649424.png`
- after: `~/kvit-works/.kvit/sessions/20261008-587061/pictures/second-press-fixed-1fd6c57f6f6e.png`

**Caught by:** Fixed once in kvit-ui's components, each with a behaviour test: dropdown, dialog, tooltip, hover card, text view, picture viewer.

### R25 Actions keep working

An action stays available while background work runs unless it depends on that work's result, and a button that is shown always does what it says. Choosing "stop" stops at once, and a failed action does not put text the user already sent back into the input.

> "Clicking 'ask' does nothing (maybe related to the fact that agent still running, but there is no reason for blocking that action)" (`works/20261012-fc04a2/L21`)
>
> "all it does is send stop here to agent and agent responds - stopping. But this is useless … if I close agent … it says that there was an error." (`works/20261022-dc518b/L26`)

**Cases (3):** `works/20261012-fc04a2/L21`, `works/20261022-dc518b/L26`, `works/20261012-201016/L21`.

**Example pair** (`works/20261012-fc04a2/L21`, owner's verdict: approved, "Ask now quotes the running turn's words; the new test passes. What next? — Stop here"): Ask now works over a running turn's text: it puts the passage in the message box as a quote, numbered as the turn the reply will become.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261012-191709-5429de17.png`
- after: `~/kvit-works/.kvit/sessions/20261012-fc04a2/pictures/transcript-selection-actions-running-turn-e2a8ccdc2c2f.png`
- after: `~/kvit-works/.kvit/sessions/20261012-fc04a2/pictures/transcript-running-turn-quoted-52448d600e86.png`

**Caught by:** Partly a program: a replayed workflow can press each visible control while a job runs and check something happens. A reviewer for the rest.

### R26 One text editor

Every box where the user writes free text uses the same editor as the main message box, with the same formatting, pasted pictures (as thumbnails) and keys. It grows up to a limit, then scrolls with the caret kept in view, and free text typed while several questions wait is one message.

> "'type your own answer' input should be similar to composer input box - allow markdown, copy paste pictures, etc." (`works/20261012-d46400/L74`)
>
> "when I type into agent's input box it expands… there is a limit… it needs to become scrollable and maintain visibility of caret after paste" (`works/20261001-66c0e5/L21`)

**Cases (6):** `works/20261001-66c0e5/L21`, `works/20261002-2ceabe/L21`, `works/20261012-d46400/L74`, `works/20260921-6e6921/L19`, `works/20261123-c2d9ec/L21`, `works/in-sessions-i-see-input-box-but-when-i-type/L17`.

**Example pair** (`works/20261001-66c0e5/L21`, owner's verdict: corrected, "Should the box also show that it has more text than fits? — Thin scroll bar (Recommended)"): After a paste the box now scrolls itself so the line with the caret is visible; this applies to every message box in the app.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261001-204044.png`
- after: `~/kvit-works/.kvit/sessions/20261001-66c0e5/pictures/growing-input-tall-picture-1e3390ff29a9.png`

**Caught by:** Fixed once by giving the apps one editor component for free text, with tests for growth, scrolling and the caret.

## Drawing and likeness

### R27 Rendered content

Markdown is rendered everywhere it appears (replies, documents, settings, decision sheets): tables as tables, formulas typeset, diagrams drawn, pictures shown at their own size up to a maximum. Raw #, | or TeX, a "!" in place of a picture, or a tab drawn as an empty box are defects.

> "I thoght we fixed markdown rendering… However I still see not rendered table in agent chat. Investigate and fix:" (`works/20260924-889c45/L16`)
>
> "I see you putting out the. Instructions as unrendered Markdown. Why do we need unrendered? We have all the editor, right? … render and allow editing within." (`hub/20261031-ad2899/L81`)

**Cases (10):** `works/20261031-48e09d/L21`, `works/20261110-9c9549/L21`, `hub/20261031-ad2899/L81`, `hub/20261122-b8db41/L21`, `works/20260921-7d4818/L19`, `works/20260924-581a33/L16`, `works/20260924-889c45/L16`, `works/20261024-cc4bb5/L21`, `works/20261106-9f9962/L21`, `works/20261107-e402fe/L21`.

**Example pair** (`works/20261110-9c9549/L21`, owner's verdict: approved, "Stop here"): The Windows package now includes the missing math library, the file view runs the same table cleanup as replies, and agents are told to write diagrams as mermaid and math between dollar signs.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261110-161049-c67a1a0c.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261110-161111-8c41eb14.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261110-161207-0f49e115.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261110-161223-67dd8ae2.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261110-161258-1249845a.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261110-161407-874d8da0.png`
- after: `~/kvit-works/.kvit/sessions/20261110-9c9549/pictures/file-view-after-1d183082df85.png`

**Caught by:** Fixed once in the block editor that draws Markdown, with rendering tests; a program can also scan the visible text of captured screens for raw Markdown patterns.

### R28 Same thing, same look

One kind of thing looks and behaves alike everywhere: annotations, charts side by side, status marks, icons for row actions, sidebar items (each opens a page, none a dialog), documents (each opens in the same side panel). Where kvit-ui has a control, the app uses it rather than drawing its own.

> "each of the items in the left panel opens the view on the right panel. Other than ketchup… weird dialog box that cannot move, cannot expand." (`hub/20261102-ef8d24/L133`)
>
> "this looks like improperly alligned. Measure (and other relevant) buttons shouls have dropdown attached with [All], [5 days]" (`hub/20261031-7729b2/L563`)

**Cases (9):** `works/20261002-47752f/L21`, `works/20261022-9a482b/L21`, `hub/20261031-7729b2/L563`, `hub/20261031-ad2899/L30`, `hub/20261102-ef8d24/L133`, `hub/20261123-1ec2ef/L21`, `works/20260921-18100c/L19`, `works/20260921-4a956a/L87`, `hub/20261103-8b0b3a/L21`.

**Example pair** (`hub/20261102-ef8d24/L133`, owner's verdict: approved, "Catch-up as steps"): Catch up became a page in the main view like the other sidebar items, and then a page of numbered steps.

- before: `~/kvit-hub/.kvit/sessions/20261102-ef8d24/pictures/catch-up-after-64a421d48275.png`
- after: `~/kvit-hub/.kvit/sessions/20261102-ef8d24/pictures/catch-up-running-ac00aaa439b2.png`

**Caught by:** Partly a program: a source check that reports panels with their own drawing or pointer handling outside an allow-list (proposed in agents-ux.md, section 3.4). A reviewer for likeness across screens.

### R29 Like the product the user knows

Where the user already knows a product for the same job and points to it (Claude's web and phone apps for agent chat), the screen follows that product's layout: the question card, the reading column of about 90 characters.

> "question tool produces ui like this: … I want to improve to look more like what I see in claude web or iphone app" (`works/20261012-d46400/L21`)
>
> "[earlier] chat windows for each agent … kind of mostly works, but it's still quite clunky … iPhone app often looks better. [here] claude pic:" (`works/20261110-c9ea0a/L102`)

**Cases (3):** `works/20261012-d46400/L21`, `works/20261012-d46400/L37`, `works/20261110-c9ea0a/L102`.

**Example pair** (`works/20261012-d46400/L21`, owner's verdict: corrected, "'type your own answer' input should be similar to composer input box - allow markdown, copy paste pictures, etc."): The question is now one card with the question as its title, each option a full-width tile with a bold label, description and number key, an 'Other' tile holding the text field, a × for ending the turn, and several questions shown one at a time.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261012-140932-13535bbe.png`
- before: `~/kvit-works/.kvit/pasted/assets/message-20261012-141212-e64f536f.png`
- after: `~/kvit-works/.kvit/sessions/20261012-d46400/pictures/question-box-end-turn-dark-244a02c38494.png`
- after: `~/kvit-works/.kvit/sessions/20261012-d46400/pictures/question-box-first-of-two-dark-4d7e154efd4f.png`
- after: `~/kvit-works/.kvit/sessions/20261012-d46400/pictures/agent-asks-two-light-bb2fe49fd72a.png`

**Caught by:** A reviewer given screenshots of the reference product, which the proposed UX.md file for each app would list.

### R30 Size, spacing and alignment

Body text has one size and colour across the window, including embedded documents and summaries. Row height fits the content, list items sit closer than paragraphs, icons are no larger than the text they sit in, things in a row share a line, and a status mark sits inside its tile's corner.

> "the font is much heavier. I think it's even larger, or I don't know. It feels much heavier than the rest of the app. Can you see how to improve?" (`hub/20261108-333e95/L21`)
>
> "there is no need to have so much padding around each item, unless we want to have like in email - title and line or two of body." (`hub/20261030-01364d/L236`)

**Cases (8):** `works/20261022-20a28e/L21`, `hub/20261030-01364d/L236`, `hub/20261108-333e95/L21`, `hub/20261123-ac6aed/L21`, `hub/20261124-dac3cc/L21`, `works/20260921-c430f3/L19`, `works/20261002-8a248a/L21`, `works/in-sessions-view-when-i-select-a-session-i-see/L18`.

**Example pair** (`hub/20261030-01364d/L236`, owner's verdict: approved, "Ship it now"): Rows are now slim one-line rows with a small space above and below the question.

- before: `~/kvit-hub/.kvit/sessions/20261030-01364d/pictures/decisions-inbox-14f04858610a.png`
- after: `~/kvit-hub/.kvit/sessions/20261030-01364d/pictures/decisions-folded-a82ba2598ea7.png`
- after: `~/kvit-hub/.kvit/sessions/20261030-01364d/pictures/decisions-checked-91a41fa71eee.png`

**Caught by:** Partly a program: the existing font-size and geometry checks keep sizes to the design values; alignment within a row can be tested from layout positions. A reviewer for the rest.

### R31 Drawn cleanly

Controls are drawn whole and sharp: lines inside their outlines and on whole pixels, round marks round at every scale, nothing clipped, a dropdown's arrow compact, a divider in a filled button soft. A page appears complete at once, and a popup for a selection stays next to that selection.

> "split button rendering looks off - there is this vertical line that goes beyon horizontal, and is word alignment correct?" (`hub/20261102-c37560/L21`)
>
> "When I'm in MD file and select something, I see a pop-up… But it's way misaligned. Showing somewhere else, very far away from selection." (`works/20261105-38265f/L21`)

**Cases (10):** `works/20261012-892b23/L21`, `works/20261026-3e4351/L21`, `works/20261120-56c428/L21`, `hub/20261102-c37560/L21`, `ui/20261101-82667a/L59`, `ui/20261101-82667a/L158`, `ui/20261101-82667a/L234`, `works/20261001-095012/L21`, `works/20261105-319b12/L21`, `works/20261105-38265f/L21`.

**Example pair** (`works/20261120-56c428/L21`, owner's verdict: approved, "Stop here"): Each dot is now placed on whole pixels before drawing, so the circle fits inside its clip at every screen scale.

- before: `~/kvit-works/.kvit/pasted/assets/message-20261120-091934-d37bd147.png`
- after: `~/kvit-works/.kvit/sessions/20261120-56c428/pictures/marks-before-after-f16440013c11.png`
- after: `~/kvit-works/.kvit/sessions/20261120-56c428/pictures/window-marks-before-after-291c670715c4.png`

**Caught by:** Fixed once in kvit-ui's drawing, with picture tests at several interface sizes.
