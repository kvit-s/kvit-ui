// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_CATALOGWRITER_H
#define KVIT_UI_CATALOGWRITER_H

#include <QString>

class QQmlEngine;

namespace KvitUi {

// Write the vocabulary skill's catalogue.
//
// Generated rather than written by hand, because a hand-written catalogue goes
// stale the first time a property is added and nothing says so — and the whole
// value of the skill is that what an agent is told exists is what the
// repository actually contains.
//
// Two sources, both inside this binary. The prose and the code samples come
// from gallery/Catalog.qml, which is the same list the gallery renders and
// tests/test_gallery compiles. The property lists are read out of the
// component QML in the module's own resources, so a property added to a
// component appears in the catalogue on the next build.
//
// Returns false and prints why on failure.
bool writeCatalog(QQmlEngine *engine, const QString &path);

}   // namespace KvitUi

#endif // KVIT_UI_CATALOGWRITER_H
