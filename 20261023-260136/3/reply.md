I read both epubs: 83 pages, 1,045 rules extracted, 854 kept. Where the current guidance disagrees with the Windows 7 guide on behaviour, the current guidance now wins for Windows (and so for Linux). All seven files in `docs/ui-guidelines/` are updated. There are now 199 rules, 99 of them with separate Windows and macOS forms, and every kept rule from all three sources is in at least one of them.

**Six pages the brief named aren't in the epubs:** keyboard accessibility, keyboard interactions, access keys, keyboard accelerators, focus navigation and visual feedback. The access-key, shortcut-modifier and focus rules still rest on the Windows 7 guide (Q-39). Send those six if you want them checked too.

**Changed on Windows**
- **Settings:** a page in the window that applies each change at once, with no OK or Cancel, and named Settings, not Options. It is the last item of the sidebar or the toolbar's More menu, and About moves to the bottom of the Settings page. OK, Cancel and Apply stay only for a properties dialog of one object.
- **Menu bar:** now optional. Apps with many commands (editing apps) keep one; others may use a toolbar with a More menu, as long as every command is on one of them.
- **Dialogs:** a default button only when one answer is clearly likely and safe. The title is the question, with no caption naming the command. Order is main action, optional second action, then the safe button at the right end, and every dialog has a safe button.
- **Wording:** sentence case for everything, window and dialog titles included. The problem message is the app's ("We couldn't save your file"). "Select" replaces "click" in help text. Toolbar labels are preferably one word.
- **Context menus:** every command for the item, common ones first, and commands that don't apply are left out rather than disabled.
- **Labels and links:** field labels sit above the field. Links only navigate; actions are buttons, and every link's tooltip says where it leads.

**Now the same on both platforms:** no colon after labels; placeholders stay while the field is empty; radio buttons for two to five options; a mixed checkbox goes all on, then all off; a submit button may stay disabled if the required fields are marked; steppers stop at their limits. kvit-ui's Check and Stepper already behave this way.

**Removed from the Windows 7 forms:** the four-word toolbar labels, wrapping steppers, "a link may run a light command", disabled Cut and Copy in text menus, the Options dialog, About in the Help menu, and the colon on labels.

**Added:** one new rule, LAY-10, for in-app navigation (back history, breadcrumb, sidebar or top row). Badges inside the window, in-window notice bars and document tabs were folded into existing rules.

**Open questions:** 11 are now settled (Q-1, Q-2, Q-9, Q-22, Q-26, Q-27, Q-28, Q-29, Q-30, Q-32, Q-37) and five narrowed. The current guidance raised 12 new ones (Q-41 to Q-52), so 41 remain open. The one with the biggest cost is Q-41: Windows 11 lets apps put Back, search and tabs in the title bar, which would mean kvit-ui drawing its own title bar; I recommend keeping the system frame.

kvit-ui now enforces 20 rules, partly covers 101, contradicts 9 and doesn't cover 69. KEY-6 (the default button and Enter) moved from contradicts to partly, because a dialog without a default now suits Windows. The order of library work in `library-gaps.md` is unchanged, except that the toolbar is now built alongside the menu bar, and a Settings page and a teaching-tip component are added.
