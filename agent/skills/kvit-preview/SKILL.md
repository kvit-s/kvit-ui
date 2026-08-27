---
name: kvit-preview
description: Look at what you built. Renders a kvit HTML mockup to a PNG, runs an application's or the gallery's screenshot pass across all four themes, or opens one component against sample data. Use after changing anything that draws, before saying a screen is done, and whenever asking whether a screen looks right.
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

A proposal for a screen, before any QML is written. Static HTML against the
shared stylesheets, rendered headless at the size the design assumes.

```
kvit-ui/ux/render.sh                  # every mockup in the directory
kvit-ui/ux/render.sh today-a.html     # just these
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
Fastest way to see whether a change to a component did what was intended.

```
kvit-ui-gallery --page KvitSlimRow --theme dark
kvit-ui-gallery --page KvitBar --theme highContrast --interface-size 20
```

The gallery is also the whole reference: launched with no arguments it lists
every component with its states and a working code sample, and the theme and
interface size are controls along the top so the same component can be
compared across all four themes without leaving the page.

## Every component, or every screen

The screenshot pass. One image per component per theme, under fixed names.

```
kvit-ui-gallery --shots <directory>       # the component set: 68 x 4 images
./build.sh --shots                        # build, test, then write the set
./build.sh --shots-only                   # just the set
```

In an application, its own screenshot run does the same for its screens
(`build.sh --shots` in kvit-notes and kvit-hub).

The names are stable on purpose. What gets reviewed after a token change is
the diff against the previous run, not the whole set looked at again — sixty
components in four themes is more images than anyone judges carefully twice.
So: keep the previous run, write the new one beside it, and compare.

```
mv build/screenshots build/screenshots-before
./build.sh --shots-only
# then compare the two directories and look only at what moved
```

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

## When it will not run

The gallery needs a display or an offscreen platform. Under a headless shell,
`QT_QPA_PLATFORM=offscreen` is enough for the screenshot pass and for a single
component; it is not enough to interact with the gallery.

The mockup renderer needs a headless Chromium and says so if it cannot find
one. That failure is worth reporting rather than working around: a screen
described but not rendered has not been reviewed.
