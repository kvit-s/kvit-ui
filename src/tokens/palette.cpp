// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "palette.h"

#include <QStringList>

#include <algorithm>
#include <cmath>

namespace {

// sRGB's transfer function and its inverse. Every piece of arithmetic below
// happens in linear light, because averaging, mixing and the colour-vision
// matrices are all linear operations and doing them on encoded values gives
// answers that are wrong in a way that looks plausible.
qreal toLinear(qreal encoded)
{
    return encoded <= 0.04045 ? encoded / 12.92
                              : std::pow((encoded + 0.055) / 1.055, 2.4);
}

qreal toEncoded(qreal linear)
{
    return linear <= 0.0031308 ? linear * 12.92
                               : 1.055 * std::pow(linear, 1.0 / 2.4) - 0.055;
}

struct Linear
{
    qreal r = 0.0, g = 0.0, b = 0.0;
};

Linear linearOf(const QColor &color)
{
    const QColor rgb = color.toRgb();
    return { toLinear(rgb.redF()), toLinear(rgb.greenF()),
             toLinear(rgb.blueF()) };
}

QColor colorOf(const Linear &l)
{
    auto clamp = [](qreal v) { return std::clamp(v, 0.0, 1.0); };
    return QColor::fromRgbF(clamp(toEncoded(std::clamp(l.r, 0.0, 1.0))),
                            clamp(toEncoded(std::clamp(l.g, 0.0, 1.0))),
                            clamp(toEncoded(std::clamp(l.b, 0.0, 1.0))));
}

struct Oklab
{
    qreal l = 0.0, a = 0.0, b = 0.0;
};

// Björn Ottosson's OKLab, verbatim. The cube root in the middle is what makes
// it perceptual: the cone responses are compressed the way the visual system
// compresses them, rather than by CIELAB's single lightness curve applied
// after the fact.
Oklab toOklab(const Linear &c)
{
    const qreal l = 0.4122214708 * c.r + 0.5363325363 * c.g + 0.0514459929 * c.b;
    const qreal m = 0.2119034982 * c.r + 0.6806995451 * c.g + 0.1073969566 * c.b;
    const qreal s = 0.0883024619 * c.r + 0.2817188376 * c.g + 0.6299787005 * c.b;

    const qreal l_ = std::cbrt(l);
    const qreal m_ = std::cbrt(m);
    const qreal s_ = std::cbrt(s);

    return { 0.2104542553 * l_ + 0.7936177850 * m_ - 0.0040720468 * s_,
             1.9779984951 * l_ - 2.4285922050 * m_ + 0.4505937099 * s_,
             0.0259040371 * l_ + 0.7827717662 * m_ - 0.8086757660 * s_ };
}

Linear fromOklab(const Oklab &v)
{
    const qreal l_ = v.l + 0.3963377774 * v.a + 0.2158037573 * v.b;
    const qreal m_ = v.l - 0.1055613458 * v.a - 0.0638541728 * v.b;
    const qreal s_ = v.l - 0.0894841775 * v.a - 1.2914855480 * v.b;

    const qreal l = l_ * l_ * l_;
    const qreal m = m_ * m_ * m_;
    const qreal s = s_ * s_ * s_;

    return { +4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s,
             -1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s,
             -0.0041960863 * l - 0.7034186147 * m + 1.7076147010 * s };
}

// The Machado, Oliveira and Fernandes (2009) matrices at full severity,
// applied in linear light. Rows are the red, green and blue outputs.
const qreal kProtanopia[9] = {
    0.152286,  1.052583, -0.204868,
    0.114503,  0.786281,  0.099216,
   -0.003882, -0.048116,  1.051998,
};
const qreal kDeuteranopia[9] = {
    0.367322,  0.860646, -0.227968,
    0.280085,  0.672501,  0.047413,
   -0.011820,  0.042940,  0.968881,
};
const qreal kTritanopia[9] = {
    1.255528, -0.076749, -0.178779,
   -0.078411,  0.930809,  0.147602,
    0.004733,  0.691367,  0.303900,
};

QString hexOf(const QColor &c)
{
    return c.name(QColor::HexRgb);
}

// WCAG 2.1 relative luminance: each channel linearised, then weighted by how
// much the eye takes from it.
qreal relativeLuminance(const QColor &c)
{
    const Linear l = linearOf(c);
    return 0.2126 * l.r + 0.7152 * l.g + 0.0722 * l.b;
}

}   // namespace

namespace KvitPalette {

Oklch toOklch(const QColor &color)
{
    const Oklab lab = toOklab(linearOf(color));
    const qreal chroma = std::hypot(lab.a, lab.b);
    qreal hue = std::atan2(lab.b, lab.a) * 180.0 / M_PI;
    if (hue < 0.0)
        hue += 360.0;
    return { lab.l, chroma, hue };
}

QColor fromOklch(const Oklch &value)
{
    const qreal radians = value.h * M_PI / 180.0;
    return colorOf(fromOklab({ value.l,
                               value.c * std::cos(radians),
                               value.c * std::sin(radians) }));
}

bool inGamut(const Oklch &value)
{
    const qreal radians = value.h * M_PI / 180.0;
    const Linear l = fromOklab({ value.l,
                                 value.c * std::cos(radians),
                                 value.c * std::sin(radians) });
    // A small tolerance, because a value that lands a thousandth outside the
    // cube is a rounding artefact rather than an unrenderable colour.
    const qreal slack = 0.001;
    for (qreal channel : { l.r, l.g, l.b }) {
        if (channel < -slack || channel > 1.0 + slack)
            return false;
    }
    return true;
}

qreal deltaE(const QColor &a, const QColor &b)
{
    const Oklab x = toOklab(linearOf(a));
    const Oklab y = toOklab(linearOf(b));
    return 100.0 * std::sqrt((x.l - y.l) * (x.l - y.l)
                             + (x.a - y.a) * (x.a - y.a)
                             + (x.b - y.b) * (x.b - y.b));
}

QColor simulate(const QColor &color, Deficiency deficiency)
{
    const qreal *m = nullptr;
    switch (deficiency) {
    case Deficiency::Protanopia:   m = kProtanopia; break;
    case Deficiency::Deuteranopia: m = kDeuteranopia; break;
    case Deficiency::Tritanopia:   m = kTritanopia; break;
    }
    const Linear c = linearOf(color);
    return colorOf({ m[0] * c.r + m[1] * c.g + m[2] * c.b,
                     m[3] * c.r + m[4] * c.g + m[5] * c.b,
                     m[6] * c.r + m[7] * c.g + m[8] * c.b });
}

qreal contrastRatio(const QColor &a, const QColor &b)
{
    const qreal la = relativeLuminance(a);
    const qreal lb = relativeLuminance(b);
    const qreal lighter = qMax(la, lb);
    const qreal darker = qMin(la, lb);
    return (lighter + 0.05) / (darker + 0.05);
}

namespace {

const char *nameOf(Deficiency deficiency)
{
    switch (deficiency) {
    case Deficiency::Protanopia:   return "protanopia";
    case Deficiency::Deuteranopia: return "deuteranopia";
    case Deficiency::Tritanopia:   return "tritanopia";
    }
    return "";
}

// Whether a theme's marks are drawn over a dark ground. The first surface is
// the window's own, which is what "the ground" means.
bool darkGround(const QList<QColor> &surfaces)
{
    return !surfaces.isEmpty() && relativeLuminance(surfaces.first()) < 0.2;
}

void checkContrast(const QList<QColor> &ramp, const QList<QColor> &surfaces,
                   qreal minContrast, QList<Finding> *out)
{
    for (int i = 0; i < ramp.size(); ++i) {
        for (const QColor &surface : surfaces) {
            const qreal ratio = contrastRatio(ramp.at(i), surface);
            if (ratio < minContrast) {
                out->append({ QStringLiteral("contrast"),
                              QStringLiteral("step %1 (%2) reads at %3:1 on %4, "
                                             "below the %5:1 floor")
                                  .arg(i)
                                  .arg(hexOf(ramp.at(i)))
                                  .arg(ratio, 0, 'f', 2)
                                  .arg(hexOf(surface))
                                  .arg(minContrast, 0, 'f', 1) });
            }
        }
    }
}

void checkReserved(const QList<QColor> &ramp, const QList<QColor> &reserved,
                   qreal minDistance, QList<Finding> *out)
{
    for (int i = 0; i < ramp.size(); ++i) {
        for (const QColor &hue : reserved) {
            const qreal distance = deltaE(ramp.at(i), hue);
            if (distance < minDistance) {
                out->append({ QStringLiteral("reserved"),
                              QStringLiteral("step %1 (%2) sits ΔE %3 from the "
                                             "reserved hue %4, inside the %5 "
                                             "keep-out")
                                  .arg(i)
                                  .arg(hexOf(ramp.at(i)))
                                  .arg(distance, 0, 'f', 1)
                                  .arg(hexOf(hue))
                                  .arg(minDistance, 0, 'f', 1) });
            }
        }
    }
}

}   // namespace

QList<Finding> validateCategorical(const QList<QColor> &ramp,
                                   const QList<QColor> &surfaces,
                                   const QList<QColor> &reserved,
                                   const RampRules &rules)
{
    QList<Finding> findings;
    checkContrast(ramp, surfaces, rules.minContrast, &findings);
    checkReserved(ramp, reserved, rules.minReservedDistance, &findings);

    const bool dark = darkGround(surfaces);
    const qreal lowest = dark ? 0.50 : 0.30;
    const qreal highest = dark ? 0.98 : 0.75;

    for (int i = 0; i < ramp.size(); ++i) {
        const Oklch coords = toOklch(ramp.at(i));
        if (coords.c < rules.minChroma) {
            findings.append({ QStringLiteral("chroma"),
                              QStringLiteral("step %1 (%2) has chroma %3, below "
                                             "the %4 floor — it reads as a grey")
                                  .arg(i)
                                  .arg(hexOf(ramp.at(i)))
                                  .arg(coords.c, 0, 'f', 3)
                                  .arg(rules.minChroma, 0, 'f', 3) });
        }
        if (rules.lightnessBand
            && (coords.l < lowest || coords.l > highest)) {
            findings.append({ QStringLiteral("lightness"),
                              QStringLiteral("step %1 (%2) has lightness %3, "
                                             "outside the %4–%5 band this "
                                             "ground allows")
                                  .arg(i)
                                  .arg(hexOf(ramp.at(i)))
                                  .arg(coords.l, 0, 'f', 3)
                                  .arg(lowest, 0, 'f', 2)
                                  .arg(highest, 0, 'f', 2) });
        }
    }

    // Every pair, in ordinary vision and under each deficiency. Pairs rather
    // than neighbours: a legend is read by matching a swatch anywhere in the
    // list to a mark anywhere in the chart, so two steps six apart converging
    // is the same failure as two adjacent ones converging.
    for (int i = 0; i < ramp.size(); ++i) {
        for (int j = i + 1; j < ramp.size(); ++j) {
            const qreal plain = deltaE(ramp.at(i), ramp.at(j));
            if (plain < rules.minSeparation) {
                findings.append({ QStringLiteral("separation"),
                                  QStringLiteral("steps %1 and %2 (%3, %4) are "
                                                 "ΔE %5 apart, under the %6 floor")
                                      .arg(i).arg(j)
                                      .arg(hexOf(ramp.at(i)))
                                      .arg(hexOf(ramp.at(j)))
                                      .arg(plain, 0, 'f', 1)
                                      .arg(rules.minSeparation, 0, 'f', 1) });
            }
            for (Deficiency deficiency : { Deficiency::Protanopia,
                                           Deficiency::Deuteranopia,
                                           Deficiency::Tritanopia }) {
                const qreal seen = deltaE(simulate(ramp.at(i), deficiency),
                                          simulate(ramp.at(j), deficiency));
                if (seen < rules.minSeparationDeficient) {
                    findings.append(
                        { QStringLiteral("separation"),
                          QStringLiteral("steps %1 and %2 (%3, %4) converge to "
                                         "ΔE %5 under %6, under the %7 floor")
                              .arg(i).arg(j)
                              .arg(hexOf(ramp.at(i)))
                              .arg(hexOf(ramp.at(j)))
                              .arg(seen, 0, 'f', 1)
                              .arg(QLatin1String(nameOf(deficiency)))
                              .arg(rules.minSeparationDeficient, 0, 'f', 1) });
                }
            }
        }
    }
    return findings;
}

QList<Finding> validateSequential(const QList<QColor> &ramp,
                                  const QList<QColor> &surfaces,
                                  const RampRules &rules)
{
    QList<Finding> findings;
    if (ramp.size() < 3) {
        findings.append({ QStringLiteral("shape"),
                          QStringLiteral("a sequential ramp of %1 steps cannot "
                                         "show a progression")
                              .arg(ramp.size()) });
        return findings;
    }

    // One end is held to the contrast floor, and it is the end away from the
    // ground. A sequential ramp encodes magnitude as distance from the
    // surface, so its "least" end is *supposed* to sit close to the surface —
    // holding that end to 3:1 would be asking the ramp not to do the one thing
    // it is for. What must not happen is the far end disappearing too.
    const QColor far = darkGround(surfaces)
        ? (toOklch(ramp.last()).l > toOklch(ramp.first()).l ? ramp.last()
                                                            : ramp.first())
        : (toOklch(ramp.last()).l < toOklch(ramp.first()).l ? ramp.last()
                                                            : ramp.first());
    checkContrast({ far }, surfaces, rules.minContrast, &findings);

    // The ordering is the whole meaning. A sequential ramp that is not
    // monotonic in lightness encodes magnitude as something the reader cannot
    // rank, and it stops working entirely in grayscale.
    const bool ascending = toOklch(ramp.last()).l > toOklch(ramp.first()).l;
    for (int i = 1; i < ramp.size(); ++i) {
        const qreal previous = toOklch(ramp.at(i - 1)).l;
        const qreal current = toOklch(ramp.at(i)).l;
        const bool ordered = ascending ? current > previous : current < previous;
        if (!ordered) {
            findings.append({ QStringLiteral("monotonic"),
                              QStringLiteral("step %1 (%2, lightness %3) does "
                                             "not continue the progression from "
                                             "step %4 (%5, lightness %6)")
                                  .arg(i)
                                  .arg(hexOf(ramp.at(i)))
                                  .arg(current, 0, 'f', 3)
                                  .arg(i - 1)
                                  .arg(hexOf(ramp.at(i - 1)))
                                  .arg(previous, 0, 'f', 3) });
        }
        // Adjacent steps must be far enough apart to be told apart at all,
        // which is a much smaller distance than a categorical ramp needs
        // because the reader is ranking rather than identifying.
        const qreal step = deltaE(ramp.at(i - 1), ramp.at(i));
        if (step < 3.0) {
            findings.append({ QStringLiteral("separation"),
                              QStringLiteral("steps %1 and %2 are ΔE %3 apart, "
                                             "which is not a visible step")
                                  .arg(i - 1).arg(i)
                                  .arg(step, 0, 'f', 1) });
        }
    }
    return findings;
}

QList<Finding> validateDiverging(const QList<QColor> &ramp,
                                 const QList<QColor> &surfaces,
                                 const RampRules &rules)
{
    QList<Finding> findings;
    if (ramp.size() < 5 || ramp.size() % 2 == 0) {
        findings.append({ QStringLiteral("shape"),
                          QStringLiteral("a diverging ramp needs an odd number "
                                         "of steps of at least five so the "
                                         "midpoint is a step; this has %1")
                              .arg(ramp.size()) });
        return findings;
    }

    const int middle = ramp.size() / 2;
    checkContrast({ ramp.first(), ramp.last() }, surfaces, rules.minContrast,
                  &findings);

    // The midpoint is the "no difference" value and must read as neutral: a
    // midpoint with a hue in it puts a direction on zero.
    const Oklch centre = toOklch(ramp.at(middle));
    if (centre.c > 0.03) {
        findings.append({ QStringLiteral("midpoint"),
                          QStringLiteral("the midpoint %1 has chroma %2, so "
                                         "zero is drawn as a colour")
                              .arg(hexOf(ramp.at(middle)))
                              .arg(centre.c, 0, 'f', 3) });
    }

    // The two arms must be different hues, and each must be far enough from
    // the midpoint that the direction is visible at one step out.
    const qreal armDistance = deltaE(ramp.first(), ramp.last());
    if (armDistance < rules.minSeparation * 2) {
        findings.append({ QStringLiteral("separation"),
                          QStringLiteral("the two ends (%1, %2) are ΔE %3 apart, "
                                         "which does not read as two directions")
                              .arg(hexOf(ramp.first()))
                              .arg(hexOf(ramp.last()))
                              .arg(armDistance, 0, 'f', 1) });
    }
    for (Deficiency deficiency : { Deficiency::Protanopia,
                                   Deficiency::Deuteranopia,
                                   Deficiency::Tritanopia }) {
        const qreal seen = deltaE(simulate(ramp.first(), deficiency),
                                  simulate(ramp.last(), deficiency));
        if (seen < rules.minSeparationDeficient) {
            findings.append({ QStringLiteral("separation"),
                              QStringLiteral("the two ends converge to ΔE %1 "
                                             "under %2, so the sign of the "
                                             "difference is lost")
                                  .arg(seen, 0, 'f', 1)
                                  .arg(QLatin1String(nameOf(deficiency))) });
        }
    }

    // Each arm is a sequential ramp in its own right, and the same reasoning
    // about lightness ordering applies to it.
    QList<QColor> low, high;
    for (int i = middle; i >= 0; --i)
        low.append(ramp.at(i));
    for (int i = middle; i < ramp.size(); ++i)
        high.append(ramp.at(i));
    for (const QList<QColor> &arm : { low, high }) {
        const QList<Finding> armFindings = validateSequential(arm, surfaces, rules);
        for (const Finding &finding : armFindings) {
            if (finding.rule == QLatin1String("contrast"))
                continue;   // already checked at the ends
            findings.append(finding);
        }
    }
    return findings;
}

}   // namespace KvitPalette
