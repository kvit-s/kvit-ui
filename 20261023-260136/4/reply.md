I've recorded all 41 decisions and carried them into the documents.

**`open-questions.md`** now starts with a table of every decision. Each question also has a "Decided (2026-10-23)" line, with its options kept for the record.

**The rules now state the decided behaviour.** About 50 rules changed. The most visible:
- **Linux** takes the Windows form, with nothing Linux-specific.
- **Capitalisation:** labels are written once in sentence case, and kvit-ui converts them to title case on macOS. Mac alert buttons are title case too.
- **Dialogs:**
  - Enter presses a focused button on Windows only.
  - Every dialog has a title, and macOS dialogs have no close symbol.
  - A message that only informs has Close on Windows and OK on macOS.
  - Flows end with Close on Windows and Done on macOS.
  - The click that closes a popover goes no further, unless it lands on a control that opens another popover.
- **Words:**
  - No "we" on either system.
  - Move to Trash… on macOS; Delete… on Windows.
  - No command links anywhere.
  - Four former Windows-only details now apply everywhere: button text within 26 characters, no italics, short option labels, and an asterisk on required fields.
- **Keys:**
  - Full access keys on Windows and Linux.
  - Redo on both Ctrl+Y and Ctrl+Shift+Z, and Ctrl+Q never bound to an app command.
  - Ctrl+F and Ctrl+E both go to the search field, and Ctrl+Tab switches tabs everywhere.
  - On a Mac, Tab focus follows the system's full-keyboard-access setting.
  - The Windows 7 forms stand for the keyboard rules.
- **Window:** kvit-ui keeps the system frame and its Px(880) × Px(600) minimum.
- **Saving:** kvit-notes and kvit-works save automatically; kvit-cash asks on close and gets its Ctrl+S.

`library-gaps.md`, `checklist.md`, `skill-section-draft.md` and the README match the decisions. The platform table in `library-gaps.md` lists what still differs between Windows and macOS, and what now behaves the same on both.

**kvit-ui against the decided rules:** it enforces 18, partly covers 105, contradicts 9 and doesn't cover 67, and 94 rules keep separate Windows and macOS forms. One rule newly conflicts with the library:
- **PTR-13:** a press outside a popover still reaches whatever is under it (`layer.go:31-34`).

The other eight conflicts are as before:
- **WIN-1:** dialog button order.
- **WRD-7:** the "Confirm" default label.
- **WRD-3:** the library's own labels carry articles.
- **WRD-25:** "Working…" as a busy button's default text.
- **FDB-10:** the progress percentage shown by default.
- **CMD-11:** shortcuts in the text field's context menu.
- **WIN-8:** typed input lost when a dialog is dismissed.
- **MSG-10:** notifications posted while the app is in front.

The pipeline will commit the changes when this turn ends.
