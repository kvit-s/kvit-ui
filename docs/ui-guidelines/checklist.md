# Reviewer checklist

For a reviewer agent that looks at one screenshot of a Kvit screen together with the screen reader's list of its controls (role, name, description, state). Each item is a judgment a test cannot make, or cannot make yet. Each names its rules (rules.md) and says what a failure looks like.

Before starting, note which platform form the screenshot is in (Windows and Linux, or macOS): items marked **[form]** differ between the two, and rules.md gives both.

Report each failure as: item number, the control (its name from the list), what is wrong, and the fix in the screen's own words.

## Labels

1. **Every command label is a verb that names its result.** WRD-1. Fails: a button or menu item that is a statement or a description ("The data is right", "Rows shown"), or a generic verb where a specific one exists ("Manage", "Configure", "Do it"). Standard verbless labels (Settings, About, New, Next) are fine.
2. **Command labels are one to three words.** WRD-2. Fails: any button or menu item with four or more words ("Take the new default", "Save to the workspace"), or a sentence on a button. A Windows toolbar button may have four.
3. **No articles or possessives on labels, menu items or tooltips.** WRD-3. Fails: "Close the preview", "Add the rows shown", "Use the shipped default", "Your settings".
4. **Confirming buttons name the action.** WRD-7 **[form]**. Fails: OK, Confirm or Yes on a button that inserts, imports, deletes or sends. On the Windows form, OK is right on a settings window and Yes/No on a yes-or-no question.
5. **The abandoning button is Cancel, and Cancel undoes nothing that already happened.** WRD-8 **[form]**. Fails: "Never mind", "Go back", "Dismiss" on a dialog's abandoning button; Cancel shown after the work is done (it should read Close).
6. **Capitalisation follows the form and is the same for every label of a kind.** WRD-6 **[form]**. Fails on Windows and Linux: Title Case on a button, menu item or tab. Fails on macOS: sentence case on a menu item or button. Fails everywhere: a menu that mixes the two ("Quick Capture…" beside "Quick capture note…").
7. **Ellipsis exactly where the command asks for more.** WRD-4 **[form]**. Fails: a command that opens a dialog, chooser or confirmation without "…"; a command that acts at once with one; three full stops instead of the ellipsis character.
8. **The label fits what the control does.** WRD-10, WRD-1. Fails: "Delete…" that moves to the trash, "Close agent" that opens a browser, "Add the rows shown" that selects all.
9. **Tooltips say what the control does, briefly.** WRD-13, HLP-2 **[form]**. Fails: a tooltip that repeats a visible label, explains how to operate the control ("Click to…"), runs past one line on the macOS form, or ends with a full stop on the macOS form. On the Windows form the tooltip of a symbol-only control is its name, with the shortcut in parentheses.
10. **No implementation words.** WRD-23. Fails: "popover", "modal", "pane", "list view", "text box", "agent session object" in visible text.
11. **Plain, calm words.** WRD-15, WRD-28. Fails: "Oops", exclamation marks, "invalid", "illegal", "fatal", "abort", "failed to", jargon without a definition.

## Symbols and text

12. **Symbols for familiar actions, words for the rest.** CMD-22. Fails: a text button made of letters or punctuation standing in for a symbol ("B", "I", "U", "S", "<>", "H", "+ Insert"); an unfamiliar symbol with no words where one or two words would say it.
13. **Every symbol-only control has a name and a tooltip.** A11Y-1, HLP-1. Fails: in the controls list, a button whose name is empty, a symbol name ("x", "chevron"), or a different word from its tooltip.
14. **One symbol, one meaning.** CON-2. Fails: the same symbol on two controls on this screen that do different things (an X that archives beside an X that closes); a standard action with a non-standard symbol (a trash can that archives).

## Commands and menus

15. **Every command on the screen is also in the menu bar.** CMD-1 **[form]**. Fails: a button or context-menu command with no menu-bar item (check the app's menus when they are in the set).
16. **Menus show their shortcuts.** KEY-14, KEY-2. Fails: Save, New, Find, Close, Undo or Copy in a menu with no shortcut beside it; a context menu showing shortcuts on the macOS form (CMD-11).
17. **Standard menus in the standard order, standard items where people look.** CMD-2, CMD-3, CMD-4, CMD-5 **[form]**. Fails on Windows and Linux: Settings in the File menu, About outside Help, Exit not last in File, Quit instead of Exit. Fails on macOS: Settings or About outside the app menu, "Preferences" instead of "Settings…", Exit instead of Quit.
18. **Menus are grouped and short.** CMD-10, CMD-12, CMD-13. Fails: a context menu over about 15 lines or four groups; a submenu inside a submenu; unrelated commands in one group.
19. **A destructive item is last, after a separator.** CMD-23. Fails: Delete in the middle of a menu, or directly under an ordinary item with no separator.
20. **A toggle says what it does in each state.** WRD-12 **[form]**. Fails on macOS: a menu item reading "Sidebar" that hides it. Fails on Windows: "Show status bar" that becomes "Hide status bar" instead of keeping its name and a tick.
21. **The screen's main actions are visible, one of them primary.** CTL-1, CMD-17. Fails: two or more primary (accent) buttons in one region; the main action hidden in a "More" menu; a menu button with one or two items.

## Dialogs

22. **Buttons in the platform's order.** WIN-1 **[form]**. Fails on Windows and Linux: Cancel left of the confirming button. Fails on macOS: the default button not at the trailing end.
23. **A default where the platform has one, none where it must not.** KEY-6 **[form]**. Fails: a destructive confirming button drawn as the default; on the Windows form, a dialog with no default at all.
24. **A dialog asks one thing, in its title.** WRD-26, MSG-4. Fails: a title that is "Error", "Warning", "Confirm" or the app's name; a question buried in the detail; two questions in one dialog.
25. **Detail adds something and does not describe the buttons.** MSG-5. Fails: detail that restates the title, or says "Click OK to continue".
26. **A destructive dialog has Cancel, and says what will be lost.** WIN-3, MSG-3. Fails: a delete confirmation with no Cancel, or one that does not name what is deleted or say it cannot be undone.
27. **No dialog for something that can be undone.** MSG-3. Fails: "Are you sure?" before an action the app could undo with a toast or ConfirmInPlace.
28. **Informing dialogs have one way to close.** WRD-7, WRD-8 **[form]**. Fails: OK and Cancel on a dialog that only informs (Close on Windows, OK on macOS).

## Messages and feedback

29. **Errors sit next to their cause and say how to fix it.** MSG-6. Fails: an error at the top of a form about one field; "Invalid input" with no hint of what is valid; an error shown in a toast that disappears.
30. **Progress says what is happening.** WRD-25, FDB-2, FDB-10. Fails: "Loading…", "Working…" or "Please wait" alone; a percentage beside a bar that already shows it; a spinner where the length is known.
31. **Empty views say what would be there and how to fill it.** MSG-9. Fails: a blank area, or an axis drawn around zeros, where a list or chart has no data.
32. **Nothing needed vanishes on a timer.** MSG-8. Fails: an undo, an error or a result shown only in a toast that times out.
33. **Status is brief and next to what it describes.** FDB-5. Fails: sentences in the status bar; a selection count far from the list it counts.

## Keyboard (from the controls list)

34. **Every control can be reached and named.** KEY-1, A11Y-1. Fails: a control in the screenshot missing from the list, or listed without a role or with a role that does not match what it looks like (a "group" that is pressable).
35. **Access keys on the Windows form.** KEY-15 **[form]**. Fails (once Q-11 is decided): two controls on one screen with the same underlined letter; a menu or menu item with none.
36. **Focus is visible and sensible.** KEY-8. Fails: the screenshot's focused control (ring) is not the one people would use first; in a dialog with a field, focus on a button.

## Layout

37. **Fields have visible labels.** WRD-19, LAY-6 **[form]**. Fails: a field whose only label is its placeholder (except a search field); a label far from its field; on the Windows form, a label beside its field without a colon.
38. **Settings explain in a line, not a paragraph.** WRD-22. Fails: more than one short line of permanent text under a setting (the 50-word paragraph); explanations on options that need none.
39. **Option labels are parallel and say only what differs.** WRD-29. Fails: radio options of very different length or grammar; shared words repeated in every option; a label that is half a sentence finished by a field.
40. **Nothing needed is cut short where there is room.** LAY-9, PTR-12. Fails: an ellipsis on a label, column or button with space beside it; a cut-short name with no way to see the rest.
41. **The most important thing is at the top leading corner.** LAY-1. Fails: the main content or action below or behind secondary material.
42. **Related things together, at most about seven groups.** LAY-2. Fails: a control far from what it acts on; a group inside a group; a page of more than about seven sections without disclosure.
43. **Advanced options hidden, common ones shown.** CTL-16. Fails: a long page of rarely used settings in full view; a "More options" that hides something everyone needs.

## Consistency across screens

When the reviewer has several screens of one app:

44. **One name per command.** CON-1. Fails: the same command named differently on two screens, in the tray or in a tooltip.
45. **The same thing done the same way.** CON-4. Fails: two similar choices made with different controls (a switch here, a checkbox there); a field that moves or changes its label between screens.
