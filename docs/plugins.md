# Plugins: building Kvit apps from replaceable parts

This note describes how a Kvit desktop app could be put together from parts
that its users can rearrange, replace and extend, while the app itself stays a
single executable. It is a design proposal written on 2026-10-02, and none of
it is implemented yet. It sets out the model, the rules that let systems of
this kind keep working as the app changes, the existing systems the model
borrows from, and the decisions that are still open.

The owner's goal is composability: an app is assembled from parts, a default
assembly ships with it, and users adjust that assembly or add their own parts
without replacing the app. Loading and unloading code while the app runs is
not a goal, and neither is letting outside code paint.

## Summary

Kvit Notes, Kvit Works, Kvit Cash, Kvit Hub and kvit-term are Go programs
drawn with kvit-ui on top of the unison toolkit. Each of them builds its
screens in Go code, so changing what a screen contains means changing the app
and building it again. The proposal changes four things:

1. **A window is a tree of named areas.** An area is a named place in the
   window, such as the sidebar, the row of a list, or the status bar. Each area
   declares a Go interface, and it is filled by a **component** that
   implements that interface: a kvit-ui widget or one of the app's own,
   registered under a name.
2. **The app's own arrangement is a description file.** This file, the
   **default composition**, says which component fills which area and with
   which settings. It is embedded in the executable.
3. **Users change the app with snippets.** A snippet is a small file in the
   user's folder that changes the default composition. It can replace the
   component in an area, add to an area that holds a list, wrap a component,
   remove one, or change its settings. A snippet can also define a new
   component assembled from existing ones.
4. **Events can be handled outside the app's Go code.** When a snippet routes
   an event such as a click, a menu command or a confirmed value to a
   **handler**, the app sends the event to it. A handler is a script run by an
   engine compiled into the app, a WebAssembly module, or a separate program.
   The handler answers with changes to the screen.

Users cannot add new drawing or new kinds of interaction. Every component that
draws is compiled into the executable. Snippets arrange, configure and combine
those components, and handlers decide what happens when something is clicked
or chosen.

```
kvit-notes (one executable)
├── compiled components    kvit-ui's widgets and the app's own
├── default composition    a description file, embedded
└── handler runtimes       script engine, WebAssembly runtime, process launcher

the app's folder in the user's configuration directory, next to ui.json
├── snippets/              changes applied on top of the default composition
└── handlers/              scripts, .wasm modules, programs
```

## How a Kvit screen is built today

unison has no layout file format, so it has nothing like HTML, Qt's `.ui`
files or QML. A screen is a tree of panels (the unit unison lays out and
draws), and Go code creates that tree when a window or view opens. kvit-ui's
components are constructors that build panels. For example,
`NewCard(ui, content...)` creates a card and adds the content panels to it as
children (`card.go`).

Five properties of unison and kvit-ui shape the design below:

- **Panels keep their behaviour in function fields.** `DrawCallback` draws
  the panel, `DrawOverCallback` draws on top of its children, and
  `InputCallbacks` holds mouse and keyboard handling (unison 0.108.0,
  `panel.go`). Any code that holds a panel can replace these, or keep the old
  function and call it from the new one. Wrapping a component therefore needs
  no change to unison.
- **Panels have no names.** A panel has no ID, and unison has no lookup by
  name of the kind `getElementById` gives a web page. The only name a panel
  has is the one it gives to screen readers. A panel can hold any data in a
  key/value map (`ClientData`). The loader proposed here keeps its own table
  from area name to panel.
- **All input and drawing happen on one thread, the UI thread.**
  `unison.InvokeTask(f)` schedules a function to run on it (`task.go`), and
  results computed on other goroutines reach the screen that way.
- **Some components draw many parts without a panel for each.** A table draws
  its cells with detached components used as stamps (`stampAt` in `cell.go`).
  An area inside such a component, such as a table's cell, is filled by a
  stamp rather than a panel, and that area's interface has to say so.
- **kvit-ui's own rules apply.** cgo is always off and every build must
  cross-compile, and unison stays unmodified. Everything proposed here would
  be written in kvit-ui or in the apps. Settings are kept in
  `ui.json` in the app's folder under the user's configuration directory
  (`DefaultSettingsPath` in `ui.go`).

## The composition model

### Areas, interfaces and components

Each area states four things:

| What it states | Example |
|---|---|
| a name, unique among its siblings; its path is the chain of names from the window down | `window/sidebar/row` |
| the Go interface a component in it must implement | a row renderer that draws one note in the tree |
| whether it holds one component or a list | the editor holds one; the status bar holds a list |
| the values it passes down to the components beneath it | the open note, the current selection |

A component is registered under a name together with a constructor. The
constructor takes the component's settings from the description and returns
something that implements one or more area interfaces. Registration happens
in Go when the app starts, the way Caddy's modules call
`caddy.RegisterModule` from an `init` function. A sketch, in which every name
is a placeholder:

```go
// The interface of window/sidebar/row: draws one note of the tree.
type NoteRow interface {
	Row(note NoteItem) unison.Paneler
}

func init() {
	compose.Register("note-row", func(s compose.Settings) (any, error) {
		return newNoteRow(s.Bool("show-dates")), nil
	})
}
```

### The description

The default composition of a notes window might read:

```yaml
version: 1
window:
  component: notes-window
  areas:
    sidebar:
      component: note-tree
      areas:
        row: { component: note-row, show-dates: true }
    editor:
      component: rich-editor
    statusbar:
      items:
        - { id: word-count, component: word-count }
        - { id: sync-state, component: sync-state }
```

### What a snippet can do

Each change in a snippet addresses an area by its path and applies one of
five operations:

| Operation | Meaning | Example |
|---|---|---|
| replace | put a different component in an area that holds one | `window/editor` becomes `markdown-editor` |
| add | insert into an area that holds a list, placed relative to a named item | add a label after `word-count` |
| remove | take a component out | remove `sync-state` |
| wrap | keep the component and put another around it; the wrapper can draw on top of it or see its events first | wrap `note-row` with `tag-badge` |
| configure | change a component's settings | `word-count` counts characters instead of words |

```yaml
version: 1
changes:
  - wrap: window/sidebar/row
    with: { component: tag-badge }
  - add: window/statusbar
    after: word-count
    item: { id: drafts, component: label, text: "Drafts", on-click: drafts }
  - configure: window/statusbar/word-count
    set: { count: characters }
```

Wrapping is the operation that makes small changes cheap. Without it, a user
who wants one extra mark on each row has to copy the whole row component and
keep the copy in step with the app's own version as the app changes.

Positions in a list are given relative to named items, never as indexes. An
index points at a different item as soon as the default composition gains or
loses one.

### Components defined in snippets

A snippet can define a component as a named piece of description with
parameters. Once defined, it can be used like a compiled component:

```yaml
components:
  account-summary:
    params: [account]
    body:
      component: card
      children:
        - { component: label, text: "{account.name}" }
        - { component: amount, value: "{account.balance}" }
        - { component: button, text: "Refresh", on-click: refresh-account }
```

This is how users make new parts without writing code. How much they can make
this way depends on how broad kvit-ui's set of components is, since every
defined component is a combination of compiled ones.

### Values passed down the tree

Components need shared information: the open note, the current selection,
the account being shown. Each area can provide named values to everything
beneath it. A component asks for a value by name and receives the one from
the nearest area above it that provides that name. A snippet can override a
value for one subtree only, for example giving one pane its own data source,
without affecting the rest of the window. When a value changes, the
components that asked for it are told, so a pane follows the selection
without subscribing to anything itself. Eclipse 4's contexts and Cordis's
isolated services resolve values in the same way (see "Existing systems this
borrows from").

### Loading

When a window opens, the app:

1. reads the default composition from the executable;
2. reads the snippets in the app's `snippets/` folder in file-name order, so
   `10-tags.yaml` applies before `20-statusbar.yaml`, and applies their
   changes in that order;
3. checks the result: every component name is registered, every component
   implements its area's interface, every area that must be filled is
   filled, and every handler a snippet names exists. Each error names the
   snippet file, the line and the area;
4. builds the panel tree.

Snippets stay stored as changes and are never merged into a saved copy of
the composition. When a new version of the app changes its default
composition, the user's snippets apply on top of the new default. Eclipse 4
shows what happens otherwise. Its workbench saves the whole model when it
closes and restores it when it starts, so, in the words of the Eclipse 4 FAQ,
"any changes in the Application.e4xmi and previously-loaded fragments will be
ignored." The way out is starting with `-clearPersistedState`, which also
discards the user's layout.

What a user changes by using the app, such as pane sizes or open tabs, is
not a snippet. It stays in the settings file, so a user's snippets contain
only changes they made on purpose.

## Handling events outside the app

### Messages

A component's events are handled in Go as they are now. When the description
names a handler for an event (`on-click: drafts`), the component sends a
message to that handler and returns at once. The handler answers with
changes:

```
app → handler   {"seq": 41, "area": "window/statusbar/drafts", "event": "click"}
handler → app   {"seq": 41, "set": {"area": "window/statusbar/drafts", "text": "Drafts (3)"}}
handler → app   {"seq": 41, "replace": {"area": "window/editor", "with": { ... }}}
handler → app   {"seq": 41, "open": {"window": "drafts", "with": { ... }}}
```

The `with` values use the description format, so the loader that builds a
window at startup also applies the changes a handler sends. htmx works this
way on the web: an element names the server address its click goes to, and
the server answers with an HTML fragment that replaces part of the page.

Handlers also need to ask the app for things, such as the text of the open
note, the selection, or a search. The app offers handlers a fixed list of
functions, the same list whichever kind of handler calls them. That list and
the area interfaces together are the app's public interface for extensions.

### Three kinds of handler

| Kind | What the author ships | What keeps it apart from the app | Libraries |
|---|---|---|---|
| a script run by an engine compiled into the app | a text file that runs on every system | the engine; a script reaches only the functions the app gives it | goja (JavaScript: ECMAScript 5.1 and most of ES6), starlark-go (Starlark, a Python dialect made for Bazel), gopher-lua (Lua 5.1) |
| a WebAssembly module | one `.wasm` file for every system | the runtime: no file or network access unless the app grants it | wazero; Extism and knqyf263/go-plugin add plugin conventions on top of it |
| a separate program | an executable for each system and processor, or a script plus its interpreter | the process boundary; a crash in it leaves the app running | JSON lines over standard input and output; hashicorp/go-plugin for gRPC between Go programs |

All five libraries are written in Go without cgo, so the app remains a single
executable that cross-compiles.

Scripts suit snippets best, because the author writes a text file next to the
snippet, nobody builds anything per system, and the app decides exactly what
a script can reach. WebAssembly suits handlers written in a compiled language
or handlers that need more speed. Separate programs suit long-running work,
existing tools in other languages, and handlers that need free access to the
operating system.

Two properties of the script engines affect how the app runs them. A goja
`Runtime` and a gopher-lua `LState` may each be used by only one goroutine at
a time, so each handler gets its own runtime on its own goroutine, and
messages to it wait in a queue. All three engines can stop a script that runs
too long: goja through `Interrupt`, starlark-go through `Thread.Cancel` or a
step limit (`SetMaxExecutionSteps`), and gopher-lua through a context.
starlark-go's documentation notes that it cannot bound a script's memory.

### Rules for handlers

1. **Only events that are not tied to a frame go out.** Painting, hover,
   dragging and typing into a field happen on every frame, which is 16.7 ms
   at 60 Hz, so they stay in compiled components. Clicks, choices, confirmed
   values, menu commands and timers can go to handlers. A round trip to a
   process on the same machine takes well under a millisecond, which is
   unnoticeable for those events. VS Code draws the same line: its
   extensions supply data and commands, and the editor does all the
   painting.
2. **The UI thread never waits for a handler.** The component sends the
   message and returns. Answers arrive on another goroutine and are applied
   with `unison.InvokeTask`. If an answer takes longer than a short delay,
   the area shows that it is waiting, for example with a disabled button.
   Every answer has the sequence number of the event it answers. An answer
   to an event that a newer event has replaced is dropped.
3. **Each piece of state has one owner.** Either the handler owns the data
   and the screen shows only what it was sent, or the app owns the data and
   tells the handler about changes. When both hold the same value, they
   disagree after the first lost or late message.
4. **A failed handler affects only its own areas.** A handler that crashes,
   stops answering or sends a change that fails the check puts its areas into
   a failed state. That state says which handler failed and offers a restart.
   The rest of the window keeps working. After a restart, the app sends the
   handler the current state again.
5. **Each handler has listed permissions.** A handler can change only the
   areas its snippet lists and call only the app functions its snippet asks
   for, and the app shows both lists when a user installs it. For the owner's
   own handlers this is bookkeeping. For handlers written by others, it
   protects against things like a dialog made to look like the app's own.

## What keeps the system working over time

1. **The app is built from the same description.** The default composition
   uses the same areas, components and loader as users' snippets, and nothing
   in a window is wired in Go code where a snippet could not reach it. If an
   area is used only by outside code, the app's own development never
   exercises it, so it can break without anyone noticing. Eclipse and
   DeepSeek Harness describe themselves as "everything is a plugin", and the
   practical meaning of that phrase is this rule.
2. **Area names and interfaces are a public contract.** Once snippets exist,
   renaming an area or changing its interface breaks them. Snippets and area
   interfaces both have version numbers. When a snippet was written for an
   older version, the app reports which file and which area. The app also
   marks which areas are public. The others use the same mechanism but may
   change between versions, and snippets that touch them are flagged. The
   Firefox history below shows the cost of leaving everything public.
3. **Snippets are checked, and there is a safe mode.** A snippet that fails
   the check at loading is skipped and reported, and the rest still apply.
   The app can also start with every snippet and handler switched off, as
   Firefox and VS Code can.
4. **User-made parts are built from kvit-ui's components.** kvit-ui's
   components follow the four themes, the 10–24 px interface size and the
   screen-reader rules by themselves. Parts built from them keep those
   properties, and this is one reason the design leaves all painting to
   compiled components.
5. **There is a view of the assembled window.** A window lists each area,
   its component, the snippet that put the component there, and the
   handlers attached to it. With this much indirection, the code alone cannot
   tell you why a pane looks the way it does. Eclipse has a tool for exactly
   this question, Plug-in Spy (Alt+Shift+F1).

## Existing systems this borrows from

**Eclipse 3: extension points.** A plugin declares a named place, called an
extension point, and other plugins add to it in their `plugin.xml`: views,
editors, menu items. A menu item names its place with a string such as
`menu:org.eclipse.ui.main.menu?after=additions`, meaning the main menu after
the marker called `additions`. Eclipse reads these contributions as data,
without loading the contributing plugin's code. The mechanism is mostly
additive: a plugin can add to another plugin's menu, but cannot replace the
class that draws its editor unless that class was itself made an extension
point. The `add … after` operation in this proposal comes from here. Eclipse
plugins are JAR files loaded into the one Java process Eclipse runs in. This
proposal loads no code into the app, so that part of Eclipse has no
counterpart here.

**Eclipse 4: the modeled workbench.** The workbench is described by a model
file, `Application.e4xmi`, which lists windows, stacks of parts and parts.
Each part names the class that draws it. Plugins add elements with model
fragments, and with processors, which are code that changes the model at
startup. Each window and part has its own context for dependency injection,
and a lookup goes from the part's context up through its parents to the
application's. Of the systems listed here, this is the closest to the
proposal. Its persistence problem, described under "Loading", is why
snippets stay stored as changes.

**Cordis and DeepSeek Harness.** DeepSeek Harness is an AI agent app written
in TypeScript on Node.js. It is built on Cordis, a plugin framework described
in a paper from Peking University and DeepSeek. In Harness the tool registry,
the language-model adapters and the agent runner are all plugins, listed in a
file called `cordis.yml`. A plugin receives a context object and registers
everything through it. Services are named and requested with `inject`. A
branch of the plugin tree can be given its own provider of a service. Each
plugin's settings are checked against a schema when it loads. During
development, `--patch` files are layered on top of the base `cordis.yml`.
Three parts of this proposal correspond to these: snippets stored as changes
on top of a base file, values overridden for one subtree, and checking at
loading. Cordis's main feature, undoing everything a plugin did so it can be
unloaded while the app runs, is not needed here.

**VS Code.** Extensions run in a separate process and contribute through a
fixed set of contribution points declared in their `package.json`. Although
the editor is built on web technology, extensions cannot touch its page. The
fixed surface keeps extensions working across editor versions, and the price
is that many changes cannot be made at all. Handlers in this proposal follow
VS Code's split: data and commands go out, and painting stays in.

**The web.** A page is a tree of elements, and any script can find elements
by ID or selector and then add, replace or wrap them. Userscripts and browser
extensions reshape sites this way without the sites' cooperation. Two
weaknesses show up. Scripts change the generated page, which the site's own
code regenerates on its next re-render, so extension authors watch for
changes (`MutationObserver`) and apply their edits again. And a site never
promised to keep its IDs and class names, many of which are generated anew
by each build. Snippets in this proposal change the composition itself, and
area names are a declared contract. Web components later added encapsulation
with explicit openings. A component's shadow DOM hides its internals, and it
exposes chosen places: named `<slot>`s for content, CSS custom properties for
settings, and `::part()` for styling named internals. Those openings are the
web's version of public areas. htmx and Phoenix LiveView send events to a
server and replace parts of the page with what it answers, which is the model
for handlers here.

**Firefox before version 57.** Firefox's own window was written in XUL, an
XML interface language, plus JavaScript. Extensions used overlays to change
any part of the window and could replace any internal function. In August
2015 Mozilla announced that it would end this model. From Firefox 57, in
November 2017, only WebExtensions run, which is a fixed API largely based on
Chrome's. Mozilla gave three reasons: extensions had unrestricted access to
Firefox's internals, which was a security problem; they interfered with the
move to a multi-process architecture; and the tight coupling between Firefox
and its extensions caused delays and crashes. Once many extensions depend on
an app's internals, those internals can no longer change, and this is the
reason for marking areas public or internal.

**Caddy.** Caddy is a Go web server whose JSON configuration is a tree of
modules. Each module has an ID in a namespace, such as
`http.handlers.file_server`, and the namespace fixes which interface the
module implements. A module loads its child modules from its own part of the
configuration (`ctx.LoadModule`). Modules are compiled in and register
themselves from `init` functions. When the configuration is reloaded, "new
modules are started before the old ones are stopped." Of the Go systems
looked at, Caddy is the closest to the registry and loader proposed here.

**Go dependency-injection libraries.** Uber's Fx nests modules, can keep a
provided value private to its module (`fx.Private`), and lets many modules
each add an item to one shared slice (value groups), which resembles an area
that holds a list. Its graph is fixed once the app is constructed. samber/do
v2 has a tree of scopes in which a child scope sees its parent's services.
Either library could wire an app's services internally. Neither reads a
description from a file, so neither replaces the loader.

## Decisions still open

1. **The description format.** YAML allows comments and is easier to read.
   JSON is what `ui.json` already uses and needs no extra library.
2. **Which app goes first, and which of its areas become public.**
3. **How data reaches defined components.** This covers which values each
   area provides, and the syntax that snippets use to refer to them, such as
   `{account.balance}`.
4. **The script engine.** Supporting one engine keeps the documentation and
   the users' learning to one language. JavaScript has the most users.
   Starlark is the closest to Python and runs deterministically. Lua is the
   smallest.
5. **Where snippets and handlers live.** One option is a folder per app next
   to `ui.json`. Another is a shared folder used by every Kvit app as well.
6. **The protocol for separate programs.** The candidates are JSON lines over
   standard input and output, or gRPC through hashicorp/go-plugin.
7. **Installing others' snippets.** How a user installs a snippet written by
   someone else, and how they review the areas and functions it asks for.

## Sources

- Yifan Shi, Wei Zhang and Tianyi Cui, "A Programming Paradigm for
  Spatiotemporal Composability", [arXiv 2608.25512](https://arxiv.org/abs/2608.25512)
  (Cordis; sections 1.2.1, 3.2.3, 6.2 and 7 compare it with VS Code, OSGi and
  Eclipse).
- DeepSeek Harness: [product page](https://www.deepseek.com/en/harness/),
  developer documentation on
  [plugins](https://deepseek-harness.github.io/deepseek-harness/en/develop/basic/),
  [services](https://deepseek-harness.github.io/deepseek-harness/en/develop/framework/service)
  and [configuration](https://deepseek-harness.github.io/deepseek-harness/en/develop/basic/config).
- Eclipse: [Eclipse4/RCP/FAQ](https://wiki.eclipse.org/Eclipse4/RCP/FAQ)
  (persisted model state and `-clearPersistedState`);
  [EclipseSource's Eclipse 4 tutorial](https://eclipsesource.com/blogs/2012/05/10/eclipse-4-final-sprint-part-1-the-e4-application-model/)
  (the application model).
- Firefox: [Add-on (Mozilla)](https://en.wikipedia.org/wiki/Add-on_(Mozilla));
  [Why did Mozilla remove XUL add-ons?](https://yoric.github.io/post/why-did-mozilla-remove-xul-addons/).
- [Extending Caddy](https://caddyserver.com/docs/extending-caddy).
- [Uber Fx](https://pkg.go.dev/go.uber.org/fx), [samber/do](https://github.com/samber/do).
- Handler libraries: [goja](https://github.com/dop251/goja),
  [starlark-go](https://github.com/google/starlark-go),
  [gopher-lua](https://github.com/yuin/gopher-lua),
  [wazero](https://github.com/tetratelabs/wazero),
  [Extism Go SDK](https://github.com/extism/go-sdk),
  [knqyf263/go-plugin](https://github.com/knqyf263/go-plugin),
  [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin).
- unison 0.108.0: `panel.go` (callback fields and `ClientData`), `task.go`
  (`InvokeTask`). kvit-ui: `card.go`, `cell.go` (`stampAt`), `ui.go`
  (`DefaultSettingsPath`).

The descriptions of Eclipse 3's menu locations, Eclipse 4's contexts, Plug-in
Spy, VS Code, web components, htmx and Phoenix LiveView come from general
knowledge of those systems and have no source listed above.
