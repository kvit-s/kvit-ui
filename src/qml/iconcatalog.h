// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_ICONCATALOG_H
#define KVIT_UI_ICONCATALOG_H

#include <QObject>
#include <QString>
#include <QStringList>
#include <QtQml/qqmlregistration.h>

namespace KvitUi {

// What symbol a name draws.
//
// KvitIcon.qml asks this for the one character to put in a Text item. The
// lookup is in C++ rather than in a QML switch statement for two reasons: the
// table is 1,530 glyphs, which is not a switch anybody maintains by hand, and
// the same table is what the vocabulary skill's catalogue is generated from,
// so what an agent is told exists and what a call site can ask for are the
// same list.
//
// The table itself is generated — src/qml/iconcatalog.cpp is written by
// tools/generate-icon-catalog.py and CI fails if it is stale.
class IconCatalog : public QObject
{
    Q_OBJECT
    QML_ELEMENT
    QML_SINGLETON

public:
    explicit IconCatalog(QObject *parent = nullptr);

    // The character to draw for `name`, or an empty string when the name is
    // not a symbol.
    //
    // A meaning name resolves first (`chevron-right`), then a Phosphor name
    // (`caret-right`). Returning empty rather than a placeholder is
    // deliberate: KvitIcon turns it into something visible, and the gallery
    // check and the QML lint both fail on it, so a typo cannot ship as a
    // silently blank box the way kvit-notes-pro's version allowed.
    Q_INVOKABLE QString glyph(const QString &name) const;

    // The names a call site should be written in: what a symbol means, not
    // what it looks like. This is the list the vocabulary skill publishes.
    Q_INVOKABLE static QStringList meaningNames();

    // Every name the font itself knows. Much longer, and the escape hatch for
    // a symbol that has not been given a meaning yet.
    Q_INVOKABLE static QStringList glyphNames();
};

}   // namespace KvitUi

#endif // KVIT_UI_ICONCATALOG_H
