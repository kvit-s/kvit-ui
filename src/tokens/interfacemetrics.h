// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef INTERFACEMETRICS_H
#define INTERFACEMETRICS_H

#include <QObject>
#include <QString>

class SettingsStore;

// The chrome's type scale and its geometry: the sidebar, the lists, the
// headers, the rows, the chips, the bars, the dialogs and the space between
// all of it. Everything, that is, except document text, which Typography owns
// (accessibility.md Finding 4).
//
// Two settings rather than one, because the needs differ and both directions
// are coherent: a person who wants large body text with a small, dense note
// list is asking for something reasonable, and so is the reverse. Operating
// system display scaling — which Qt honours on all three platforms —
// enlarges everything at once and cannot express either.
//
// The shape: state in the settings store under an `interface.` key prefix,
// clamped setters, one change signal, and every value below derived from one
// base through px(). Nothing here is a literal that stands still while the
// rest of the interface grows.
//
// ── What is in here, and why geometry sits beside type ──────────────────────
//
// Font size alone is not enough. A 24-pixel label inside a 28-pixel button
// clips, so the geometry has to travel with the type scale. kvit-hub learned
// the same thing from the other end: its density block was 141 lines of bare
// integers — row heights, chip heights, radii, pane widths — that stood still
// while its type setting moved, so raising the interface size made rows too
// tight for the text in them. Every one of those values is a property here.
//
// ── The type scale ─────────────────────────────────────────────────────────
//
// Seven roles, which is the merge of two scales that had five and nine.
//
// kvit-notes' chrome had five roles in pixels off a base of 12. kvit-hub had
// nine written as integer offsets from a base of 15 in points, which at 96 dpi
// is 12 to 25 pixels. The nine collapse onto seven because three of them —
// `typeSecondary`, `typeRow` and `typeBody` — were three names for the two
// sizes ordinary row text is set at, and `typeName` is what `strong` already
// meant. What kvit-hub had above `title` and kvit-notes did not is real, so
// `headline` and `display` are new.
//
// The five original names keep their meanings and their pixel values, so
// kvit-notes' chrome is pixel-identical after Wave 2 and its 3,702 call sites
// need no rename.
class InterfaceMetrics : public QObject
{
    Q_OBJECT

    Q_PROPERTY(int fontSize READ fontSize WRITE setFontSize NOTIFY changed)
    // fontSize / 12.0: what px() multiplies by, exposed so a call site that
    // needs a real rather than a rounded integer can do its own arithmetic.
    Q_PROPERTY(qreal scale READ scale NOTIFY changed)

    // The clamps, so a settings stepper takes its range from the one place
    // that enforces it rather than repeating the numbers. kvit-hub's
    // "Interface size" row is the caller that most needs this: its label has
    // always said interface size while the value it wrote was the document
    // font size, which is kvit-notes' setting (divergence 3.4).
    Q_PROPERTY(int minFontSize READ minFontSize CONSTANT)
    Q_PROPERTY(int maxFontSize READ maxFontSize CONSTANT)

    // The two families the chrome draws in. Empty means the platform's own
    // default, which is what a desktop application should look like unless
    // somebody has asked otherwise.
    //
    // These are here rather than on Typography for the same reason the sizes
    // are: Typography is the document, and a person who sets their notes in a
    // serif has not asked for a serif sidebar. kvit-hub reached for
    // `typography.fontFamily` because it had nowhere else to look, and drawing
    // its chrome in the document font is the visible half of the settings
    // collision described in divergence 3.4.
    Q_PROPERTY(QString fontFamily READ fontFamily WRITE setFontFamily NOTIFY changed)
    // For an identifier: a reference, a hash, a key. A monospace family is
    // what makes two of them comparable down the column.
    Q_PROPERTY(QString monoFamily READ monoFamily WRITE setMonoFamily NOTIFY changed)

    // The seven roles the chrome asks for by name. The comment on each is the
    // pixel size it answers at the default base, and what asks for it.
    Q_PROPERTY(int caption READ caption NOTIFY changed)    // 10 — kind tags, counts
    Q_PROPERTY(int small READ small NOTIFY changed)        // 11 — chip labels, sub-lines
    Q_PROPERTY(int body READ body NOTIFY changed)          // 12 — row text, prose
    Q_PROPERTY(int strong READ strong NOTIFY changed)      // 13 — a name, an emphasised row
    Q_PROPERTY(int title READ title NOTIFY changed)        // 15 — a section heading
    Q_PROPERTY(int headline READ headline NOTIFY changed)  // 17 — a pane title, the wordmark
    Q_PROPERTY(int display READ display NOTIFY changed)    // 20 — a page title

    // ── Spacing ────────────────────────────────────────────────────────────
    //
    // The gap prd.md §5.3 names: kvit-hub writes 146 bare `spacing:` literals
    // and kvit-notes 59 more, concentrated at 8, 6 and 2 and running down to
    // 1 — every one of them below `stackGap`, which was the smallest value
    // either had a name for.
    //
    // Seven steps. 3 and 5 are deliberately not among them: a scale with
    // every integer in it is not a scale, and a call site writing 5 is
    // choosing between 4 and 6 without having thought about which. The two
    // migrations round each 3 and each 5 onto a neighbour one call site at a
    // time, which is a judgement about that row rather than a global rule.
    Q_PROPERTY(int spaceTight READ spaceTight NOTIFY changed)  // 2
    Q_PROPERTY(int spaceSnug READ spaceSnug NOTIFY changed)    // 4
    Q_PROPERTY(int spaceNear READ spaceNear NOTIFY changed)    // 6
    Q_PROPERTY(int space READ space NOTIFY changed)            // 8 — the ordinary gap
    Q_PROPERTY(int spaceWide READ spaceWide NOTIFY changed)    // 10
    Q_PROPERTY(int spaceLoose READ spaceLoose NOTIFY changed)  // 12

    // ── Layout ─────────────────────────────────────────────────────────────
    Q_PROPERTY(int viewMargin READ viewMargin NOTIFY changed)   // 16
    Q_PROPERTY(int columnGap READ columnGap NOTIFY changed)     // 14
    Q_PROPERTY(int stackGap READ stackGap NOTIFY changed)       // 7
    Q_PROPERTY(int sidebarWidth READ sidebarWidth NOTIFY changed)  // 232
    Q_PROPERTY(int railWidth READ railWidth NOTIFY changed)     // 48 — the sidebar collapsed
    Q_PROPERTY(int paneWidth READ paneWidth NOTIFY changed)     // 392 — the side pane
    Q_PROPERTY(int headerHeight READ headerHeight NOTIFY changed)      // 52
    Q_PROPERTY(int breadcrumbHeight READ breadcrumbHeight NOTIFY changed)  // 34
    Q_PROPERTY(int statusBarHeight READ statusBarHeight NOTIFY changed)    // 22

    // ── The row, in its four heights ───────────────────────────────────────
    //
    // A view picks the one its content needs rather than shrinking a row to
    // fit more in. kvit-hub's note on why: 26 projects at the full height need
    // 1,456 pixels against the 874 a stage has, which is why its list at rest
    // is slim.
    Q_PROPERTY(int rowHeight READ rowHeight NOTIFY changed)              // 56
    Q_PROPERTY(int rowHeightSub READ rowHeightSub NOTIFY changed)        // 48
    Q_PROPERTY(int rowHeightSlim READ rowHeightSlim NOTIFY changed)      // 30
    Q_PROPERTY(int rowHeightCompact READ rowHeightCompact NOTIFY changed)  // 24

    // ── Control and mark heights ───────────────────────────────────────────
    Q_PROPERTY(int controlHeight READ controlHeight NOTIFY changed)  // 28
    Q_PROPERTY(int tabHeight READ tabHeight NOTIFY changed)          // 30
    Q_PROPERTY(int chipHeight READ chipHeight NOTIFY changed)        // 17
    Q_PROPERTY(int tagHeight READ tagHeight NOTIFY changed)          // 16
    Q_PROPERTY(int pillHeight READ pillHeight NOTIFY changed)        // 15
    Q_PROPERTY(int barHeight READ barHeight NOTIFY changed)          // 7
    Q_PROPERTY(int barHeightWide READ barHeightWide NOTIFY changed)  // 9
    Q_PROPERTY(int iconSize READ iconSize NOTIFY changed)            // 18
    Q_PROPERTY(int iconSizeSmall READ iconSizeSmall NOTIFY changed)  // 13

    // ── Corner radii and rules ─────────────────────────────────────────────
    Q_PROPERTY(int radiusBar READ radiusBar NOTIFY changed)          // 2
    Q_PROPERTY(int radiusChip READ radiusChip NOTIFY changed)        // 3
    Q_PROPERTY(int radiusControl READ radiusControl NOTIFY changed)  // 4
    Q_PROPERTY(int radiusCard READ radiusCard NOTIFY changed)        // 6
    Q_PROPERTY(int radiusPill READ radiusPill NOTIFY changed)        // 8
    // One device pixel at the default, and never less than one at any size:
    // a separator that scales to zero is a line that vanishes at the smallest
    // interface size, which reads as a layout bug rather than as a smaller
    // interface.
    Q_PROPERTY(int hairline READ hairline NOTIFY changed)            // 1
    // A focus ring is two pixels because one is invisible against a border
    // (accessibility.md Finding 3).
    Q_PROPERTY(int focusRingWidth READ focusRingWidth NOTIFY changed)  // 2

    // ── The widths a view reflows at ───────────────────────────────────────
    //
    // These scale too, which is the point of putting them here: at twice the
    // interface size the same screen needs twice the width before it stops
    // being cramped, so a fixed 880 would keep a doubled interface in its
    // widest layout while its content collided.
    Q_PROPERTY(int widthFloor READ widthFloor NOTIFY changed)    // 880
    Q_PROPERTY(int widthLaptop READ widthLaptop NOTIFY changed)  // 1100
    Q_PROPERTY(int widthDrawn READ widthDrawn NOTIFY changed)    // 1440

public:
    // Clamps: 10 is the smallest size the chrome stays legible at, and 24 is
    // where a list row stops fitting a title and a date on one line.
    static constexpr int DefaultFontSize = 12;
    static constexpr int MinFontSize = 10;
    static constexpr int MaxFontSize = 24;

    explicit InterfaceMetrics(QObject *parent = nullptr);

    void setSettings(SettingsStore *settings);

    int fontSize() const { return m_fontSize; }
    void setFontSize(int size);
    QString fontFamily() const { return m_fontFamily; }
    void setFontFamily(const QString &family);
    QString monoFamily() const { return m_monoFamily; }
    void setMonoFamily(const QString &family);
    qreal scale() const { return m_fontSize / qreal(DefaultFontSize); }
    int minFontSize() const { return MinFontSize; }
    int maxFontSize() const { return MaxFontSize; }

    int caption() const { return px(10); }
    int small() const { return px(11); }
    int body() const { return px(12); }
    int strong() const { return px(13); }
    int title() const { return px(15); }
    int headline() const { return px(17); }
    int display() const { return px(20); }

    int spaceTight() const { return px(2); }
    int spaceSnug() const { return px(4); }
    int spaceNear() const { return px(6); }
    int space() const { return px(8); }
    int spaceWide() const { return px(10); }
    int spaceLoose() const { return px(12); }

    int viewMargin() const { return px(16); }
    int columnGap() const { return px(14); }
    int stackGap() const { return px(7); }
    int sidebarWidth() const { return px(232); }
    int railWidth() const { return px(48); }
    int paneWidth() const { return px(392); }
    int headerHeight() const { return px(52); }
    int breadcrumbHeight() const { return px(34); }
    int statusBarHeight() const { return px(22); }

    int rowHeight() const { return px(56); }
    int rowHeightSub() const { return px(48); }
    int rowHeightSlim() const { return px(30); }
    int rowHeightCompact() const { return px(24); }

    int controlHeight() const { return px(28); }
    int tabHeight() const { return px(30); }
    int chipHeight() const { return px(17); }
    int tagHeight() const { return px(16); }
    int pillHeight() const { return px(15); }
    int barHeight() const { return px(7); }
    int barHeightWide() const { return px(9); }
    int iconSize() const { return px(18); }
    int iconSizeSmall() const { return px(13); }

    int radiusBar() const { return px(2); }
    int radiusChip() const { return px(3); }
    int radiusControl() const { return px(4); }
    int radiusCard() const { return px(6); }
    int radiusPill() const { return px(8); }
    int hairline() const { return px(1); }
    int focusRingWidth() const { return px(2); }

    int widthFloor() const { return px(880); }
    int widthLaptop() const { return px(1100); }
    int widthDrawn() const { return px(1440); }

    // A design-pixel value scaled and rounded, for a measurement that has no
    // name yet: `implicitHeight: 28` becomes `implicitHeight: Interface.px(28)`.
    //
    // Reach for a named property first. px() with a literal in it is the same
    // unnamed value the migration is removing, only scaled; it is right for a
    // one-off and wrong for anything a second view will also need.
    Q_INVOKABLE int px(int designPx) const;

    // Back to the built-in default (a settings dialog's "Reset interface size").
    Q_INVOKABLE void resetToDefaults();

signals:
    void changed();

private:
    SettingsStore *m_settings = nullptr;
    bool m_loading = false;
    int m_fontSize = DefaultFontSize;
    QString m_fontFamily;
    QString m_monoFamily = QStringLiteral("monospace");
};

#endif // INTERFACEMETRICS_H
