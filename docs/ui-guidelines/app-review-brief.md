# Brief: reviewing the apps against the interface rules

This is the brief for step 8 of `guidelines-plan.md`: bringing kvit-cash, kvit-notes and kvit-works in line with the interface rules in `rules.md`. It is written for the session that works in each app's repository, one app at a time, and says what has to happen before the work can start, what the work is, in what order, and when it is done.

## Background

The rules came from Apple's and Microsoft's interface guidelines, condensed to 199 rules for Kvit's desktop apps. They cover wording, dialogs, menus, shortcuts and keyboard use, and each rule has a Windows form, which Linux also follows, and a macOS form; the app picks the form of the system it runs on. A survey of the three apps on 2026-10-22 found the problems the rules are meant to prevent: sentences on buttons, letters typed where symbols belong, one command under several names, OK on buttons that insert or import, permanent paragraphs of explanation, and menus that show no shortcuts.

Steps 1 to 7 and 9 of the plan built the rules into kvit-ui. The components follow the form by themselves (dialog button order, title case on the macOS form, access keys on the Windows form, standard shortcuts), and `uitest` has checks that read a running screen and report where it breaks a rule. As of 2026-10-25 kvit-ui fully keeps 77 of the rules, partly keeps 93 and breaks none; the other 29 are judged by a reviewer or wait until an app needs them. Step 8 is what is left: moving each app onto this version of kvit-ui, running the checks in the app's tests, and fixing what they and a reviewer find.

## Before any app can start

1. **kvit-ui's work has to reach GitHub's `main`.** Everything described here is on the kvit-ui branch `track/implement-guidelines-79f5f2`. The apps take kvit-ui only from GitHub: `./build.sh --update-gh` in an app runs `go get github.com/kvit-s/kvit-ui@main`, and each app's `CLAUDE.md` forbids pointing `go.mod` at a local folder. So the branch has to be merged into kvit-ui's `main` and pushed first.
2. **unison has to be one version.** kvit-ui, kvit-cash and kvit-notes use unison v0.108.0, and kvit-works uses v0.108.1. The Kvit repositories keep one unison version between them, and `tools/bump-unison.sh <version>` in kvit-ui moves them all and runs every repository's tests.
3. **The apps compile against this kvit-ui as they are.** On 2026-10-25 each app was built against this checkout through a temporary module file outside its repository, and all three built without errors. Their tests were not run. Moving an app is therefore a change to `go.mod`, and what changes is what its screens do: a dialog's buttons move into the form's order, labels are title-cased on the macOS form, access keys appear on the Windows form. An app test that pins a button order, a label's case or a menu's lines may fail after the move, and the failure shows what changed.

## The work in each app, in order

### 1. Move to this kvit-ui

Run `./build.sh --update-gh` and then `./build.sh --test`, and fix what fails. Where a test fails because kvit-ui now does what the rules ask, the screen is already right and the test is what changes.

### 2. Give the app's agents the kvit-ui skills

kvit-ui has two agent skills: `kvit-ui`, which teaches an agent how to build a Kvit screen and includes the 28 rules of its section "Words, symbols, commands and shortcuts", and `kvit-preview`, for drawing a proposed screen as an HTML page before writing it in Go. Today they reach no session working in the apps: Claude Code finds skills in a project's `.claude/skills/` folder and in the user's `~/.claude/skills/`, and neither holds them.

They are already on disk. When `go get` moves an app to a kvit-ui commit, Go downloads the whole kvit-ui repository at that commit, skills included, into its module cache (`~/go/pkg/mod/github.com/kvit-s/kvit-ui@<version>/agent/skills/`). Extend `./build.sh --update-gh` so that, after `go mod tidy`, it copies both skill folders from there into the app's `.claude/skills/`, replacing the copies from the last update, and commit them together with the `go.mod` change. The agents working in the app then read the skills written for the kvit-ui version the app builds against, on any machine and in cloud sessions, with nothing installed for each user; a skill from a newer kvit-ui could describe components the app cannot use yet.

```sh
# The kvit-ui skills, from the kvit-ui commit go.mod now names. Files in
# Go's module cache are read-only, so the copies are made writable for the
# next update to replace them.
kvitui=$(go list -m -f '{{.Dir}}' github.com/kvit-s/kvit-ui)
mkdir -p .claude/skills
for skill in kvit-ui kvit-preview; do
    rm -rf ".claude/skills/$skill"
    cp -R "$kvitui/agent/skills/$skill" ".claude/skills/$skill"
    chmod -R u+w ".claude/skills/$skill"
done
```

None of the three apps ignores `.claude/` in git (checked 2026-10-25). The skill's description, which decides when an agent uses it, names the apps it is for; on 2026-10-25 it named kvit-cash and kvit-notes but not kvit-works, which has to be added in kvit-ui before kvit-works takes the skill.

### 3. Add the checks to the app's tests

The apps already use `uitest`: kvit-cash in 3 test files, kvit-notes in 14 and kvit-works in 55, and each calls `uitest.CheckRules` on its source. Three checks are added:

- **One screen, in both forms.** `Session.CheckLabels` reads what the screen shows: command labels (a verb first, at most three words and 26 characters, no articles, the form's case), dialogs (a title, a safe button, a named action, never Confirm), menus, tooltips, access keys, help links, the Tab order, and Apply or OK where they do not belong. About 40 checks in all; `rules.md` names the rule behind each report.
- **The whole app.** `uitest.CheckApp` compares every screen collected with `Session.Collect`: one name, one symbol and one key per command, one meaning per symbol, standard commands on their standard keys, every command on a menu bar or toolbar, and every help topic present in `UI.HelpRoot`.
- **The source.** `uitest.CheckRules`, which the apps already call, now also reports "Oops", "click here", three full stops typed for an ellipsis, a key binding written with `mod.Control` or `mod.Command`, and a label starting with "+ ".

Each check takes a list of known problems. A problem not on the list fails the test, and so does a listed problem that no longer occurs, so the list can only get shorter. Start each list with every problem the checks report on the first run, so the tests pass, and then remove entries as screens are fixed. A known problem is matched by any part of its report, such as `WRD-2 (Windows form): button "Save to the workspace"`.

```go
func TestScreensFollowTheInterfaceRules(t *testing.T) {
	var screens []uitest.Screen
	for _, form := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
		s := uitest.Open(t, uitest.Options{Form: form}, buildTransactionsScreen)
		s.CheckLabels(knownTransactionsProblems...)
		screens = append(screens, s.Collect("transactions, "+form.String()+" form"))
		s.Screen.Stop() // one headless screen runs at a time
	}
	// ... the same for every other screen ...
	uitest.CheckApp(t, screens, knownAppProblems...)
}
```

### 4. Move the commands onto `Command` and a menu bar

The rules need every command to have one name, one symbol and one key, shown wherever the command is offered (CON-1, KEY-2, CMD-1). kvit-ui's `Command` carries all three, and a `MenuBar`, a `Toolbar`, `NewCommandButton` and menu lines draw from it; `Window.AddCommands` makes the window answer the command's key, and a command with `Std` set (StdSave, StdUndo, StdFind and the rest) gets the platform's standard name and key without choosing them. The apps still catch keys in `Window.OnKeyDown` (kvit-cash in 3 places, kvit-notes in 3, kvit-works in 1) and build menus from `MenuItem` literals (20, 62 and 34 of them), so each app needs its commands declared once and its menus and keys taken from them. `agent/skills/kvit-ui/SKILL.md`, rules 14 to 22, says how.

### 5. Fix the screens, starting with what the survey found

Fix each screen through kvit-ui's components rather than by drawing around them. These are the survey's findings, every one still in the code on 2026-10-25; the line numbers are where they were then.

| App | Where | Now | Fix | Rule |
|---|---|---|---|---|
| kvit-cash | `app/health/card.go:167` | button "The data is right" | "Mark correct" | WRD-1 |
| kvit-cash | `app/configuration/editor.go:116`, `:153` | "Take the new default", "Use the shipped default" | one command, "Reset to default…" | WRD-2, CON-1 |
| kvit-cash | `app/configuration/editor.go:157` | "Save to the workspace" | "Save" | WRD-2 |
| kvit-cash | `app/preview/preview.go:56` | "Close the preview" on a back arrow | "Back" | WRD-3 |
| kvit-cash | `app/transactions/transactions.go:112`, `views/proposal/proposal.go:39` | "Add the rows shown", which selects every row | "Select all" | WRD-10 |
| kvit-cash | its menus | none of 23 menu items shows a shortcut | commands with keys, shown in the menus | KEY-14 |
| kvit-notes | `app/toolbar.go:105-146`, `editor/formatbar.go:59-71` | buttons "B", "I", "U", "S", "<>", "H" | symbol buttons named Bold, Italic, Underline, Strikethrough, Code, Highlight | CMD-22 |
| kvit-notes | `app/links.go:260` | OK on Insert link | "Insert" | WRD-7 |
| kvit-notes | `app/exchange.go:299` | OK on Confirm import | "Import" | WRD-7 |
| kvit-notes | `app/tray.go:52`, `app/menus.go`, `app/capture.go` | "Quick Capture…", "Quick capture note…" and "Quick capture" for one command | one `Command`, "Quick capture…" | CON-1, WRD-6 |
| kvit-notes | `app/actions.go:108` | "Delete…" that moves a folder to the trash | "Move to Trash…" on the macOS form; "Delete…" saying it goes to the Recycle Bin on the Windows form | WRD-10 |
| kvit-notes | `app/settings.go:269` | a permanent 50-word paragraph under one setting | one line, the rest in a Hint | WRD-22 |
| kvit-notes | its menus | 4 of 109 menu items show a shortcut | every item with a key shows it | KEY-14 |
| kvit-works | `ui/shell/overviewcard.go:61` | the close symbol for Archive | the `archive` symbol | CON-2 |
| kvit-works | `ui/conversation/view.go:1242` | the close symbol for Withdraw | the `undo` symbol | CON-2 |
| kvit-works | `ui/shell/detailheader.go:184` | the close symbol named "Close agent", which opens a browser | the `external` symbol, "Open in browser" | CON-2, WRD-10 |

### 6. Review each screen with the checklist

Some rules need judgment a test cannot make: whether a verb is the specific one, whether a message says what to do next, whether a dialog was needed at all. `checklist.md` lists them for a reviewer agent, which looks at a screenshot of one screen together with the screen reader's list of its controls (`Session.Tree`) and reports each failure with the control, what is wrong and the fix. Review every screen in both forms, since some items differ between them.

### 7. Done

An app is done when its known-problem lists are empty and the reviewer's report on every screen has been worked through. Record what changed in rules.md's "kvit-ui now" lines only where kvit-ui itself changed; the app's own record is its commit history.

## Where things are

- `rules.md`: every rule, its Windows and macOS forms, how it is enforced, and what kvit-ui does about it now.
- `checklist.md`: the reviewer's checklist.
- `agent/skills/kvit-ui/SKILL.md`, section "Words, symbols, commands and shortcuts": 28 rules for whoever builds or changes a screen, with before-and-after examples; step 2 puts it in each app's `.claude/skills/`.
- `agent/skills/kvit-ui/catalog.md` and the gallery (`./build.sh --run` in kvit-ui): every component, with a code sample per example.
- `guidelines-plan.md`: the plan, and what is built in each step.
