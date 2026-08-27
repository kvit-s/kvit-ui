# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Search for chart ramps that pass the checks in src/tokens/palette.cpp.
#
# This is design work rather than a build step: it was run once per theme to
# produce the values now written into theme.cpp, and it is kept so that a later
# change to a rule — a raised contrast floor, a new reserved hue — can be
# answered with a new search rather than by hand-nudging hex values until the
# test goes quiet.
#
#   python3 tools/design-ramps.py            # search every theme
#   python3 tools/design-ramps.py dark       # just one
#
# The arithmetic here is the same as palette.cpp's, and tests/test_palette.cpp
# is what holds the two in agreement: the ramps this prints are checked there
# by the C++ implementation, so a divergence between the two shows up as a
# failing test rather than as a ramp that was never really checked.

import math
import sys

# ── colour arithmetic, matching src/tokens/palette.cpp ─────────────────────


def to_linear(v):
    return v / 12.92 if v <= 0.04045 else ((v + 0.055) / 1.055) ** 2.4


def to_encoded(v):
    return v * 12.92 if v <= 0.0031308 else 1.055 * (v ** (1 / 2.4)) - 0.055


def hex_to_linear(h):
    h = h.lstrip("#")
    return tuple(to_linear(int(h[i:i + 2], 16) / 255.0) for i in (0, 2, 4))


def linear_to_hex(c):
    out = ""
    for v in c:
        e = to_encoded(min(1.0, max(0.0, v)))
        out += "%02x" % max(0, min(255, round(e * 255)))
    return "#" + out


def oklab(c):
    r, g, b = c
    l = 0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b
    m = 0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b
    s = 0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b
    l_, m_, s_ = (math.copysign(abs(v) ** (1 / 3), v) for v in (l, m, s))
    return (0.2104542553 * l_ + 0.7936177850 * m_ - 0.0040720468 * s_,
            1.9779984951 * l_ - 2.4285922050 * m_ + 0.4505937099 * s_,
            0.0259040371 * l_ + 0.7827717662 * m_ - 0.8086757660 * s_)


def from_oklab(v):
    L, a, b = v
    l_ = L + 0.3963377774 * a + 0.2158037573 * b
    m_ = L - 0.1055613458 * a - 0.0638541728 * b
    s_ = L - 0.0894841775 * a - 1.2914855480 * b
    l, m, s = (x ** 3 for x in (l_, m_, s_))
    return (4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s,
            -1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s,
            -0.0041960863 * l - 0.7034186147 * m + 1.7076147010 * s)


def oklch(c):
    L, a, b = oklab(c)
    return (L, math.hypot(a, b), math.degrees(math.atan2(b, a)) % 360.0)


def from_oklch(L, C, H):
    r = math.radians(H)
    return from_oklab((L, C * math.cos(r), C * math.sin(r)))


def in_gamut(lin, slack=0.001):
    return all(-slack <= v <= 1.0 + slack for v in lin)


def delta_e(a, b):
    x, y = oklab(a), oklab(b)
    return 100.0 * math.dist(x, y)


def luminance(c):
    return 0.2126 * c[0] + 0.7152 * c[1] + 0.0722 * c[2]


def contrast(a, b):
    la, lb = luminance(a), luminance(b)
    hi, lo = max(la, lb), min(la, lb)
    return (hi + 0.05) / (lo + 0.05)


CVD = {
    "protanopia": (0.152286, 1.052583, -0.204868,
                   0.114503, 0.786281, 0.099216,
                   -0.003882, -0.048116, 1.051998),
    "deuteranopia": (0.367322, 0.860646, -0.227968,
                     0.280085, 0.672501, 0.047413,
                     -0.011820, 0.042940, 0.968881),
    "tritanopia": (1.255528, -0.076749, -0.178779,
                   -0.078411, 0.930809, 0.147602,
                   0.004733, 0.691367, 0.303900),
}


def simulate(c, kind):
    m = CVD[kind]
    r, g, b = c
    out = (m[0] * r + m[1] * g + m[2] * b,
           m[3] * r + m[4] * g + m[5] * b,
           m[6] * r + m[7] * g + m[8] * b)
    return tuple(min(1.0, max(0.0, v)) for v in out)


# What each theme draws on, and what its hues already mean.
#
# Surfaces are every ground a chart mark can land on. Reserved hues are what
# Theme::reservedHues() returns — the values theme.cpp already gives a meaning
# to — and the two lists are kept the same by tests/test_palette.cpp, which
# runs the C++ validator against the C++ list.
#
# `axisAttentionText` and `axisAgentText` are deliberately not reserved. They
# are the legible text partners of the two axis bar hues, used for small type
# rather than for marks, and they sit close enough to their bar hues that
# reserving both closes off enough of the wheel to leave no categorical ramp at
# all. Nothing is confused by a bar landing near the colour some words are set
# in.
THEMES = {
    "light": {
        "surfaces": ["#ffffff", "#f4f4f4", "#fafafa", "#f2f2f0"],
        "reserved": ["#4a90d9", "#b3261e", "#e05c5c", "#1e874b", "#a66908",
                     "#2970c8", "#e0a04c", "#2a9d8f", "#1f6feb", "#d99a3d",
                     "#4aa3a3", "#8a5cc0", "#c0392b"],
        "band": (0.30, 0.75),
    },
    "dark": {
        "surfaces": ["#1e1e1e", "#252526", "#212122", "#37373a"],
        "reserved": ["#5c9fe0", "#e06c60", "#e05c5c", "#5abd82", "#e0a34c",
                     "#6fb1ff", "#e0a04c", "#4db6ac", "#58a6ff", "#d9a04c",
                     "#4aa3a3", "#a37fd4"],
        "band": (0.50, 0.95),
    },
    "sepia": {
        "surfaces": ["#f6efdf", "#eee5d0", "#f2ead8", "#ece3cb"],
        "reserved": ["#9a6b2f", "#a33b30", "#c25a4e", "#4d7839", "#91641a",
                     "#8a5a20", "#b07a2a", "#3a8a7a", "#9a6a2b", "#c08a2e",
                     "#3f8f8a", "#7a4fa8", "#b07a1f"],
        "band": (0.30, 0.75),
    },
    "highContrast": {
        "surfaces": ["#000000", "#050505", "#1a1a1a"],
        "reserved": ["#33ccff", "#ff6666", "#ff8080", "#44ff99", "#ffbb33",
                     "#66ddff", "#33ffdd", "#ffff00", "#33ddcc", "#cc99ff"],
        "band": (0.55, 0.98),
    },
}

RULES = dict(min_contrast=3.0, min_separation=8.0, min_chroma=0.09,
             min_reserved=8.0)

# The most saturated a step is allowed to be.
#
# Not a correctness rule — it is a ceiling on how loud the chart may be, and it
# exists because without one the search runs straight to the edge of sRGB and
# returns eight neon values that clear every check and look nothing like the
# rest of the interface. The hues already in theme.cpp run from chroma 0.085
# (the teal) to 0.20 (the light theme's focus ring), clustered around 0.13, so
# a ceiling just above that cluster keeps a chart looking like it belongs on
# the same screen as the chrome around it.
MAX_CHROMA = 0.16


def findings(ramp, surfaces, reserved, band, rules=RULES):
    """The categorical checks from palette.cpp, as a list of strings."""
    out = []
    lin = [hex_to_linear(h) for h in ramp]
    for i, c in enumerate(lin):
        for s in surfaces:
            r = contrast(c, hex_to_linear(s))
            if r < rules["min_contrast"]:
                out.append("contrast %d on %s: %.2f" % (i, s, r))
        L, C, _ = oklch(c)
        if C < rules["min_chroma"]:
            out.append("chroma %d: %.3f" % (i, C))
        if not (band[0] <= L <= band[1]):
            out.append("lightness %d: %.3f" % (i, L))
        for h in reserved:
            d = delta_e(c, hex_to_linear(h))
            if d < rules["min_reserved"]:
                out.append("reserved %d vs %s: %.1f" % (i, h, d))
    for i in range(len(lin)):
        for j in range(i + 1, len(lin)):
            d = delta_e(lin[i], lin[j])
            if d < rules["min_separation"]:
                out.append("separation %d/%d: %.1f" % (i, j, d))
            for kind in CVD:
                d = delta_e(simulate(lin[i], kind), simulate(lin[j], kind))
                if d < rules["min_separation"]:
                    out.append("separation %d/%d under %s: %.1f"
                               % (i, j, kind, d))
    return out


def worst_separation(ramp):
    """The tightest pair in the ramp, in any vision. This is what a search
    maximises: a ramp is only as identifiable as its closest two steps."""
    lin = [hex_to_linear(h) for h in ramp]
    worst = 1e9
    for i in range(len(lin)):
        for j in range(i + 1, len(lin)):
            worst = min(worst, delta_e(lin[i], lin[j]))
            for kind in CVD:
                worst = min(worst, delta_e(simulate(lin[i], kind),
                                           simulate(lin[j], kind)))
    return worst


def candidates(theme):
    """Every colour this theme could use for a chart mark: inside the gamut,
    clearing contrast on every surface, saturated enough to read as a colour,
    inside the lightness band, and outside every reserved hue's keep-out.

    Returned as hue -> list of (lightness, chroma, hex), because what the
    categorical search needs is the set of ways one hue can be placed in one
    theme."""
    spec = THEMES[theme]
    surfaces = [hex_to_linear(s) for s in spec["surfaces"]]
    reserved = [hex_to_linear(h) for h in spec["reserved"]]
    # A margin inside every rule.
    #
    # The search works in continuous OKLCh and the answer is written out as a
    # six-digit hex value, which moves it. A step designed to sit exactly on
    # the chroma floor comes back from that rounding at 0.089 against a floor
    # of 0.090, and the C++ validator — which reads the hex, because that is
    # what ships — reports a ramp the search called passing. The margin is what
    # keeps the two agreeing.
    margin_l, margin_c = 0.01, 0.005
    lo, hi = spec["band"][0] + margin_l, spec["band"][1] - margin_l
    out = {}
    for H in range(0, 360, 2):
        placements = []
        L = lo
        while L <= hi + 1e-9:
            C = MAX_CHROMA
            while C > 0.0 and not in_gamut(from_oklch(L, C, H)):
                C -= 0.004
            if C >= RULES["min_chroma"] + margin_c:
                lin = from_oklch(L, C, H)
                ok = all(contrast(lin, s) >= RULES["min_contrast"]
                         for s in surfaces)
                if ok and all(delta_e(lin, r) >= RULES["min_reserved"]
                              for r in reserved):
                    placements.append((L, C, linear_to_hex(lin)))
            L += 0.02
        if placements:
            out[H] = placements
    return out


def _views(hex_color):
    """One colour as ordinary vision plus the three deficiencies. Separation is
    the worst of the four, because a pair that is far apart for most readers
    and converges for a deuteranope is a pair that fails."""
    lin = hex_to_linear(hex_color)
    return [lin] + [simulate(lin, kind) for kind in CVD]


def _worst_pair(views):
    worst = 1e9
    for i in range(len(views)):
        for j in range(i + 1, len(views)):
            worst = min(worst, min(delta_e(a, b)
                                   for a, b in zip(views[i], views[j])))
    return worst


# How far a theme may move a hue off the shared wheel.
#
# Zero would be the pure version of "one wheel everywhere" and it does not
# reach the separation floor: the best even wheel placeable in all four themes
# leaves its tightest pair at ΔE 7.4 in sepia, against a floor of 8. Fourteen
# degrees is enough to clear it, and it is small enough that a category keeps
# its identity — the same step is an orange in every theme, a slightly
# different orange.
HUE_TOLERANCE = 14.0


def place(hues, theme, grid, tolerance=HUE_TOLERANCE):
    """Put each of `hues` somewhere in `theme`, as far from the others as the
    theme allows.

    Lightness does part of the work, which is the point. Strip the red-green
    axis out for a deuteranope and eight hues at one lightness collapse onto a
    single blue-yellow line, where eight points cannot all be ΔE 8 apart. It is
    also what makes the ramp survive a grayscale print."""
    options = []
    for want in hues:
        near = []
        for H in grid:
            d = abs(H - want) % 360.0
            if min(d, 360.0 - d) <= tolerance:
                near.extend(grid[H])
        if not near:
            return None
        options.append(near)

    # Precompute the four vision models for every placement once; the loop
    # below asks for them thousands of times.
    seen = [[_views(hexed) for _L, _C, hexed in opts] for opts in options]

    best_ramp, best_worst = None, -1.0
    # Three starts — every hue at its lightest allowed placement, at its
    # darkest, and alternating — because one greedy pass from a single start
    # settles into an arrangement where no single step can improve while the
    # set as a whole is poor.
    starts = [[0] * len(hues),
              [len(opts) - 1 for opts in options],
              [0 if k % 2 == 0 else len(options[k]) - 1
               for k in range(len(hues))]]
    for start in starts:
        chosen = list(start)
        improved = True
        rounds = 0
        while improved and rounds < 20:
            improved = False
            rounds += 1
            for k in range(len(hues)):
                others = [seen[j][chosen[j]] for j in range(len(hues)) if j != k]
                best_at, best_score = chosen[k], -1.0
                for index in range(len(options[k])):
                    mine = seen[k][index]
                    score = min(min(delta_e(a, b) for a, b in zip(mine, other))
                                for other in others)
                    if score > best_score:
                        best_at, best_score = index, score
                if best_at != chosen[k]:
                    chosen[k] = best_at
                    improved = True
        ramp = [options[k][chosen[k]][2] for k in range(len(hues))]
        worst = _worst_pair([seen[k][chosen[k]] for k in range(len(hues))])
        if worst > best_worst:
            best_ramp, best_worst = ramp, worst
    return best_ramp


def search_categorical(steps=8):
    """One hue wheel for the whole estate, placed differently in each theme.

    The alternative — searching each theme for its own eight best hues —
    produces four unrelated palettes, so switching from light to dark recolours
    every category on the screen and the reader has to learn the legend again.
    A category should be the same colour everywhere; what changes between
    themes is how light and how saturated that colour has to be to sit on that
    theme's ground.

    An even wheel is not available. Some hues have nowhere to go in some
    themes: teal has no placement at all on the light ground, because every
    lightness it could take at a chroma the display can reach lands inside the
    keep-out around `axisAgent` and `calloutTip`, which are both teal. So the
    wheel is chosen from the hues that are placeable in all four themes, seeded
    as evenly as that set allows and then improved."""
    grids = {theme: candidates(theme) for theme in THEMES}
    feasible = sorted(set.intersection(*(set(g) for g in grids.values())))
    if len(feasible) < steps:
        return None

    cache = {}

    def score(hues):
        key = tuple(hues)
        if key in cache:
            return cache[key]
        placed = {}
        worst_overall = 1e9
        for theme in THEMES:
            ramp = place(hues, theme, grids[theme])
            if ramp is None:
                worst_overall = -1.0
                break
            placed[theme] = ramp
            worst_overall = min(worst_overall,
                                _worst_pair([_views(c) for c in ramp]))
        cache[key] = (worst_overall, placed)
        return cache[key]

    def nearest(want):
        return min(feasible,
                   key=lambda h: min(abs(h - want) % 360.0,
                                     360.0 - abs(h - want) % 360.0))

    best = None
    # Sweep the wheel's rotation a degree at a time, snapping each spoke to the
    # nearest hue that is placeable everywhere. Rotation is the only free
    # parameter: letting each spoke wander independently was tried and is not
    # worth what it costs — it searches for an hour and returns a set whose
    # hues are no longer evenly spaced, which is the property that makes a
    # categorical palette read as one thing rather than as eight colours.
    for offset in range(0, 360 // steps):
        hues = []
        for k in range(steps):
            pick = nearest((offset + 360.0 * k / steps) % 360.0)
            while pick in hues:
                # Two spokes landing on the same feasible hue: nudge the second
                # onto its neighbour rather than dropping a category.
                pick = feasible[(feasible.index(pick) + 1) % len(feasible)]
            hues.append(pick)
        hues.sort()

        worst, placed = score(hues)
        if worst < RULES["min_separation"]:
            continue
        if best is None or worst > best[0]:
            best = (worst, list(hues), placed)
    return best


def _sequential_for(theme, hue, chroma, steps):
    """Build one theme's sequential ramp at a given hue, and say whether it
    holds. Returns (ramp, smallest visible step) or None."""
    spec = THEMES[theme]
    surfaces = [hex_to_linear(s) for s in spec["surfaces"]]
    dark_ground = luminance(surfaces[0]) < 0.2
    # On a dark ground the ramp runs from near the surface up to the brightest
    # end; on a light one it runs the other way. Either way the last step is
    # the one furthest from the ground.
    lo, hi = (0.45, 0.92) if dark_ground else (0.35, 0.82)
    ramp = []
    for k in range(steps):
        f = k / (steps - 1)
        L = lo + (hi - lo) * f if dark_ground else hi - (hi - lo) * f
        # Chroma peaks in the middle: an end near white or near black cannot
        # hold saturation without leaving the gamut, and forcing it there just
        # clips to a colour nobody chose.
        c = chroma * (0.45 + 0.55 * math.sin(math.pi * (0.15 + 0.7 * f)))
        while c > 0.01 and not in_gamut(from_oklch(L, c, hue)):
            c -= 0.005
        ramp.append(linear_to_hex(from_oklch(L, c, hue)))

    lin = [hex_to_linear(h) for h in ramp]
    if any(contrast(lin[-1], s) < RULES["min_contrast"] for s in surfaces):
        return None
    smallest = min(delta_e(lin[k - 1], lin[k]) for k in range(1, steps))
    if smallest < 3.0:
        return None
    # A smaller keep-out than a categorical ramp gets. A sequential ramp is
    # read as a scale — the reader sees the whole gradient at once and ranks
    # within it — rather than as eight identities to match against a legend, so
    # a step passing near a reserved hue is not mistaken for it.
    for h in spec["reserved"]:
        for c in lin:
            if delta_e(c, hex_to_linear(h)) < 6.0:
                return None
    return ramp, smallest


def search_sequential(steps=7):
    """One hue, light to dark, in every theme.

    Shared for the same reason the categorical wheel is shared: a heat map that
    is blue in the light theme and pink in the dark one is two heat maps."""
    best = None
    for hue in range(0, 360, 5):
        for chroma in (0.08, 0.10, 0.12, 0.14):
            built = {}
            worst = 1e9
            for theme in THEMES:
                made = _sequential_for(theme, hue, chroma, steps)
                if made is None:
                    worst = -1.0
                    break
                built[theme], step = made
                worst = min(worst, step)
            if worst < 0:
                continue
            if best is None or worst > best[0]:
                best = (worst, hue, chroma, built)
    return best


def _diverging_for(theme, low_hue, high_hue, chroma, steps):
    """Build one theme's diverging ramp, and say whether it holds."""
    spec = THEMES[theme]
    surfaces = [hex_to_linear(s) for s in spec["surfaces"]]
    dark_ground = luminance(surfaces[0]) < 0.2
    mid_l = 0.62 if dark_ground else 0.70
    half = steps // 2

    ramp = []
    for k in range(steps):
        d = (k - half) / half          # −1 … 0 … +1
        hue = low_hue if d < 0 else high_hue
        c = chroma * abs(d)
        # Lightness moves away from the midpoint in whichever direction leaves
        # the ground, so distance from zero reads in grayscale too.
        L = mid_l + (0.20 * abs(d) if dark_ground else -0.20 * abs(d))
        while c > 0.0 and not in_gamut(from_oklch(L, c, hue)):
            c -= 0.005
        ramp.append(linear_to_hex(from_oklch(L, c, hue)))

    lin = [hex_to_linear(h) for h in ramp]
    if any(contrast(e, s) < RULES["min_contrast"]
           for e in (lin[0], lin[-1]) for s in surfaces):
        return None
    # The midpoint is "no difference" and must be neutral: a midpoint with a
    # hue in it puts a direction on zero.
    if oklch(lin[half])[1] > 0.03:
        return None
    if delta_e(lin[0], lin[-1]) < 2 * RULES["min_separation"]:
        return None
    ends = min(delta_e(simulate(lin[0], kind), simulate(lin[-1], kind))
               for kind in CVD)
    if ends < RULES["min_separation"]:
        return None
    if min(delta_e(lin[k - 1], lin[k]) for k in range(1, steps)) < 3.0:
        return None
    # The keep-out applies to the two ends only.
    #
    # The steps between an end and the midpoint are progressively less
    # saturated versions of it, and on the sepia ground — whose whole reserved
    # set is warm and mid-lightness — one of them lands within ΔE 6 of a
    # reserved hue for every hue pair there is, so checking them all leaves no
    # diverging ramp at all. It is also the wrong thing to check: those steps
    # are read as positions on a gradient, and it is the ends that carry the
    # two meanings a reader could confuse with a semantic colour.
    for h in spec["reserved"]:
        for c in (lin[0], lin[-1]):
            if delta_e(c, hex_to_linear(h)) < 6.0:
                return None
    return ramp, ends


def search_diverging(steps=7):
    """Two hues meeting at a neutral midpoint, the same two in every theme."""
    best = None
    for low_hue in range(0, 360, 5):
        for gap in range(90, 271, 10):
            high_hue = (low_hue + gap) % 360
            for chroma in (0.10, 0.13, 0.16):
                built = {}
                worst = 1e9
                for theme in THEMES:
                    made = _diverging_for(theme, low_hue, high_hue, chroma, steps)
                    if made is None:
                        worst = -1.0
                        break
                    built[theme], ends = made
                    worst = min(worst, ends)
                if worst < 0:
                    continue
                if best is None or worst > best[0]:
                    best = (worst, low_hue, high_hue, chroma, built)
    return best


def main():
    print("── categorical: one hue wheel, placed per theme ──")
    cat = search_categorical()
    if cat is None:
        print("NO PASSING RAMP FOUND")
    else:
        worst, hues, placed = cat
        print("tightest pair in any vision, in the worst theme: ΔE %.1f" % worst)
        print("hues: " + " ".join("%.0f°" % h for h in hues))
        for theme in THEMES:
            print("  %-13s %s" % (theme, " ".join(placed[theme])))
    print()

    print("── sequential: one hue, light to dark ──")
    seq = search_sequential()
    if seq is None:
        print("NO PASSING RAMP FOUND")
    else:
        worst, hue, chroma, built = seq
        print("hue %d°, chroma %.2f; smallest visible step ΔE %.1f"
              % (hue, chroma, worst))
        for theme in THEMES:
            print("  %-13s %s" % (theme, " ".join(built[theme])))
    print()

    print("── diverging: two hues, neutral midpoint ──")
    div = search_diverging()
    if div is None:
        print("NO PASSING RAMP FOUND")
    else:
        worst, low, high, chroma, built = div
        print("hues %d° and %d°, chroma %.2f; ends stay ΔE %.1f apart under "
              "the worst deficiency" % (low, high, chroma, worst))
        for theme in THEMES:
            print("  %-13s %s" % (theme, " ".join(built[theme])))


if __name__ == "__main__":
    main()
