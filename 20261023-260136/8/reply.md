I've written `docs/ui-guidelines/guidelines-plan.md` and added it to the folder's README.

It lists the next work in order, and for each step says what gets built, which rules it serves, what "done" means, and roughly how big it is:

1. **Platform form:** kvit-ui knows whether it is drawing the Windows or the macOS form, and the tests can render both.
2. **Commands and shortcuts:** a Command type, a shortcut registry, and the standard commands with their keys.
3. **Menu bar and toolbar:** the menu bar, the toolbar with its More menu, and split and menu buttons. This is the biggest step and splits into two pieces of work.
4. **Dialog by platform.**
5. **Wording in the components:** the title-case converter, the ellipsis flag, the library's own strings, and the formatting symbols.
6. **Tests:** one that reads a single screen, one that reads every screen of an app, and checks on the Go source, all run in both forms.
7. **Keyboard.**
8. **Review the three apps:** one app per piece of work. This needs a session where the app repositories are attached.
9. **The rest of the library:** settings pieces, form rows, field and window behaviour, Tip, back history, undo stack, help links.

The plan also records what gets built and what doesn't: guideline controls that already exist under a Kvit name, the seven new components the rules need, things to add only when an app asks, and things never to build. A last section covers ongoing upkeep: run `checklist.md` with a reviewer agent on every new screen now, and update SKILL.md and rules.md as each step lands.
