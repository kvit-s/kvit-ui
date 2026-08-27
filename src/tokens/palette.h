// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_PALETTE_H
#define KVIT_UI_PALETTE_H

#include <QColor>
#include <QList>
#include <QString>

// The colour arithmetic a chart ramp is checked with.
//
// prd.md §5.5 calls the absence of chart colour the largest single gap in the
// estate, and the reason it is a gap rather than an oversight is that the five
// saturated hues the applications have are all spoken for: the amber that
// means "my time" is the same value as the amber that means "warning", the
// violet that means "discovered scope" is the same as the violet that means
// "a hygiene signal", and the red that means "a hard signal" is the same as
// the red that means "danger". A chart drawn in any of them says something it
// does not mean.
//
// So the three ramps are new, and a ramp asserted by eye is a ramp that fails
// the first time somebody edits one step. What is here is the arithmetic that
// lets a ramp be *checked*: perceptual coordinates, a distance that means
// something, the three common colour-vision deficiencies, and WCAG contrast.
// tests/test_palette.cpp runs the checks and tools/check-ramps is the same
// checks from a shell.
namespace KvitPalette {

// A colour in OKLCh: perceptual lightness 0–1, chroma 0 to about 0.4, and hue
// in degrees.
//
// OKLab rather than CIELAB because the questions asked of a ramp are exactly
// the ones CIELAB answers badly: whether two hues are equally light (CIELAB's
// blue is famously wrong), and whether a hue is saturated enough to read as a
// colour at all rather than as a grey.
struct Oklch
{
    qreal l = 0.0;
    qreal c = 0.0;
    qreal h = 0.0;
};

Oklch toOklch(const QColor &color);
QColor fromOklch(const Oklch &value);
// Whether fromOklch(value) had to be clipped to fit in sRGB. A ramp step that
// only exists outside the gamut is a step that renders as something else.
bool inGamut(const Oklch &value);

// Perceptual distance between two colours: the Euclidean distance in OKLab,
// multiplied by 100 so the numbers sit in the same range a reader knows from
// CIELAB ΔE. Roughly: 1 is the smallest difference anyone notices side by
// side, 3 is a difference visible without comparing, and 8 is the separation a
// series needs to be identifiable in a legend across a chart.
qreal deltaE(const QColor &a, const QColor &b);

// The three colour-vision deficiencies worth checking. Together they cover
// about 8% of men and 0.5% of women; the anomalous trichromacies they
// approximate (protanomaly and so on) are more common still and are less
// severe, so a ramp that clears these clears those.
enum class Deficiency {
    Protanopia,     // no long-wavelength cone: red and green converge
    Deuteranopia,   // no medium-wavelength cone: the common one
    Tritanopia,     // no short-wavelength cone: blue and green converge
};

// What `color` looks like to someone with `deficiency`, using the Machado,
// Oliveira and Fernandes (2009) matrices at full severity, applied in linear
// light.
QColor simulate(const QColor &color, Deficiency deficiency);

// WCAG 2.1 relative contrast, 1.0 (identical) to 21.0 (black on white). The
// same function Theme::contrastRatio exposes; it is here as well so the
// validator does not have to construct a Theme to ask.
qreal contrastRatio(const QColor &a, const QColor &b);

// ── What a ramp has to satisfy ─────────────────────────────────────────────
//
// Every threshold is stated rather than implied, because the interesting
// output of a failing check is which rule failed and by how much.
struct RampRules
{
    // Against every surface the ramp can be drawn on. 3.0 is WCAG's floor for
    // a graphical object, which is what a bar or a line is; a chart mark is
    // not text and 4.5 would rule out most usable hues.
    qreal minContrast = 3.0;

    // Between any two steps, in ordinary vision and under each of the three
    // deficiencies. A legend is read by matching a swatch to a mark, so two
    // steps that converge for a deuteranope are two series that cannot be told
    // apart by 6% of the men looking at the screen.
    qreal minSeparation = 8.0;
    qreal minSeparationDeficient = 8.0;

    // Below this chroma a hue reads as a grey rather than as a colour, which
    // is what happened to the teal already in the estate: at 0.085 it is
    // legible and not identifiable.
    qreal minChroma = 0.09;

    // The usable lightness band, which depends on the ground and so is
    // derived from `surfaces` rather than stated here: [0.30, 0.75] over a
    // light ground and [0.50, 0.98] over a dark one.
    //
    // It is not about contrast — the rule above covers that. It is that sRGB
    // has almost no chroma left at either extreme, so a step pushed past the
    // band is a step that arrives as a near-grey whatever hue it was asked
    // for.
    //
    // Set `lightnessBand` to override it for a ramp with a reason to.
    bool lightnessBand = true;

    // How far a step must stay from a hue that already means something.
    //
    // A chart series that lands on the value of `danger` says "stalled" to
    // anyone who has read the conventions, whatever the legend claims.
    //
    // The same distance as the separation floor, because it is the same
    // question: can a reader tell these two apart. A larger keep-out was tried
    // and it closes off enough of the hue wheel that no categorical ramp of
    // eight steps exists at all on the light and sepia grounds.
    qreal minReservedDistance = 8.0;
};

// One thing wrong with a ramp: which rule, where, and by how much.
struct Finding
{
    QString rule;
    QString detail;
};

// Check `ramp` against `surfaces` (every ground the mark can be drawn on in
// this theme) and `reserved` (the hues that already mean something in this
// theme). An empty result is a ramp that passes.
//
// `sequential` and `diverging` ramps are deliberately exempted from the
// separation and lightness rules, because both encode magnitude *through*
// lightness: adjacent steps of a sequential ramp are supposed to be close, and
// a diverging ramp's midpoint is supposed to be neutral. A sequential ramp
// holds only its far end to the contrast floor, since its near end sitting
// close to the surface is what "least" looks like; a diverging ramp holds both
// of its ends.
QList<Finding> validateCategorical(const QList<QColor> &ramp,
                                   const QList<QColor> &surfaces,
                                   const QList<QColor> &reserved,
                                   const RampRules &rules = RampRules());

QList<Finding> validateSequential(const QList<QColor> &ramp,
                                  const QList<QColor> &surfaces,
                                  const RampRules &rules = RampRules());

QList<Finding> validateDiverging(const QList<QColor> &ramp,
                                 const QList<QColor> &surfaces,
                                 const RampRules &rules = RampRules());

}   // namespace KvitPalette

#endif // KVIT_UI_PALETTE_H
