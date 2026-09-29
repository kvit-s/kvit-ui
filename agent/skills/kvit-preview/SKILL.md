---
name: kvit-preview
description: Look at what you built. Renders a kvit HTML mockup to a PNG, runs the gallery's or an application's screenshot pass across all four themes, or opens one component in the gallery. Use after changing anything that draws, before saying a screen is done, and whenever asking whether a screen looks right.
---

# Seeing what you built

An interface application's only output is what it draws. Its test suite says
whether the code runs; it says nothing about whether a view is right. An agent
that cannot see its own output is guessing, and this is what stops that.

Use it after any change that draws, before reporting a screen as finished, and
whenever the question is whether something looks right rather than whether it
works.

There are three things to look at, and which one depends on what changed.

## A mockup

A proposal for a screen, before any Go is written. Static HTML against the
shared stylesheets, rendered headless at the size the design assumes.

```
kvit-ui-go/ux/render.sh                  # every mockup in the directory
kvit-ui-go/ux/render.sh today-a.html     # just these
```

PNGs land in `renders/` beside the source. The viewport is 1440×960, which is
the desktop size every drawing in the set uses; `WIDTH` and `HEIGHT` override
it.

It needs a headless Chromium. The script finds one under
`~/.cache/ms-playwright/`, or takes `HEADLESS_SHELL` pointing at one.

A mockup loads nothing over the network: no web fonts, no remote images, no
fetched scripts. That is what makes it render the same here and on any other
machine, and a mockup that breaks the rule renders differently for the person
reviewing it than it did for whoever drew it.

To see the same mockup in another theme, put `theme-light`, `theme-dark`,
`theme-sepia` or `theme-contrast` on its `<body>`. A drawing with no theme
class renders dark.

## A component

One component, on its own, in a chosen theme and at a chosen interface size.
The fastest way to see whether a change to a component did what was
intended.

```
go run ./cmd/kvit-ui-gallery --page KvitSlimRow --theme dark
go run ./cmd/kvit-ui-gallery --page KvitBar --theme highContrast --interface-size 20
```

The gallery is also the whole reference: every component with its states and
a working code sample, which is the Go function that builds the sample, and
the theme and interface size are controls along the top (Ctrl+1 to Ctrl+4,
Ctrl+plus and Ctrl+minus) so one component can be compared across all four
themes without leaving the page.

The gallery opens a window, which needs a desktop. Under WSL,
`./build.sh --win` builds it for Windows and starts it on the Windows desktop.

## Every component, or every screen

The screenshot pass. One image per page per theme, at three interface sizes,
under the gallery's file names.

```
./build.sh --shots                                  # into build/shots
go run ./cmd/kvit-ui-gallery --shots <directory>    # anywhere
```

It runs headless, with no window and no desktop. `./build.sh --shots` also
stacks each image above the reference image of the same page, from
`~/kvit-reference/kvit-ui-0a0b210`, in `build/shots/compare/`.

The names are stable on purpose. What gets reviewed after a change to a
design value is the difference from the previous run, not the whole set
looked at again — seventy-odd pages in four themes is more images than anyone
judges carefully twice. So: keep the previous run, write the new one beside
it, and compare.

```
mv build/shots build/shots-before
./build.sh --shots
# then compare the two directories and look only at what moved
```

The gallery's `TestShots` writes the same set into a temporary directory and
fails on any warning logged while it does, a panic in a component's drawing
included.

## What to look at

**The high-contrast pass first.** It is where a distinction resting on colour
alone becomes obvious: everything that was carrying meaning through a subtle
tint goes flat, and what survives is what also had a shape, a weight, an
outline or a position.

**Then the light theme against the dark one.** A hardcoded colour looks
correct in exactly one of them.

**Then anything with no data in it.** An empty list, a chart with nothing to
plot, a figure that was never measured. Those states are the ones that get
built last and looked at least, and they are where a zero appears in place of
an em dash.

**Then the 24 px shots.** A size read once when a panel was built stays at
the size it had, and at twice the interface size it is obvious.

## When it will not run

The screenshot pass needs nothing but Go: it draws with unison's headless
screen. The interactive gallery needs a desktop.

The mockup renderer needs a headless Chromium and says so if it cannot find
one. That failure is worth reporting rather than working around: a screen
described but not rendered has not been reviewed.
