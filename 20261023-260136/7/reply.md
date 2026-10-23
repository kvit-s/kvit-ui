Not as a catalogue to copy. Add a component when a rule can't be kept without one, or when an app screen needs it. Most of what the guidelines describe either already exists under a Kvit name, or is a field or two on an existing component.

Each new component costs a gallery page, tests, catalogue entries and maintenance across four themes and two platform forms. The project's own rule (CLAUDE.md, and the skill's workflow step 3) is that a missing control becomes a kvit-ui component when a screen needs it, not before.

**Already covered under a Kvit name** (nothing to add):

| Guideline control | Kvit component |
|---|---|
| InfoBar | Notice |
| InfoBadge | Badge |
| Expander, disclosure triangle | Disclosure |
| AutoSuggestBox, combo box | TypeAhead, Select |
| NumberBox | NumberField |
| NavigationView, sidebar | Sidebar |
| SelectorBar, segmented control | Segmented |
| BreadcrumbBar | Breadcrumb |
| Flyout, popover | Popover |
| Teaching tip | Spotlight, which covers tours but not single tips (see below) |

**Missing, and the rules need them now.** These would close the gaps the survey found:

1. **MenuBar**, Kvit-drawn on Windows and Linux and the system's on macOS (CMD-1 to CMD-8).
2. **Toolbar** with an overflow More menu (CMD-1, CMD-18, CMD-19). On Windows it is the command surface for apps without a menu bar.
3. **Split button and menu button** (CMD-17). kvit-notes' "+ Insert" is one.
4. **Settings pieces:**
   - a setting row (label and one-line description, with the control at the trailing side);
   - a Settings item pinned to the bottom of the sidebar;
   - an About section (WIN-16, WRD-22).
5. **Form row:** a visible label tied to its field, plus a required-field mark (WRD-19, LAY-6, CON-4).
6. **Tip:** anchored and non-blocking, for teaching tips (HLP-4). You chose this over reusing Spotlight.
7. **Back history** for moving between pages inside a window (LAY-10).

The Command type and shortcut registry come before all of these. They aren't components, but the menu bar, toolbar and buttons are all built from them.

**Extensions to existing components instead of new ones:**
- Field: a masked mode for secrets, a clear button on Windows, a length limit.
- Switch: the On or Off word beside it on Windows.
- Button and IconButton: the `Asks` flag that adds the ellipsis, Enter on Windows, the shortcut shown in the tooltip.
- MenuItem: both toggle words, a radio (bullet) item, the default item drawn bold.
- Progress: Cancel or Stop, and time remaining.
- Dialog: the platform form, and a Close-only form for messages.
- Label: selectable text, and the full text as a tooltip when it is cut short.
- RadioGroup: arrow keys stop at the ends on Windows.

**Only when an app asks:**
- A date picker or calendar.
- A colour chooser limited to the chart ramp's steps, as decided.
- A token field.
- Document tabs, which kvit-notes might want.
- A rating control.
- An annotated scroll bar.
- A list-and-details layout helper.

**Don't add:**
- Command links (decided against).
- A column browser (Windows has none).
- Flip view, pips pager, pivot, semantic zoom.
- Media playback, person picture.
- A Kvit-drawn title bar (decided against).

So the next components are the seven in the "missing" list, in roughly that order, all built on the Command type. Everything else waits for a screen that needs it. `library-gaps.md` already has these in its build order, under steps 3 and 10.
