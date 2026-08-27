// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_UISINGLETONS_H
#define KVIT_UI_UISINGLETONS_H

#include <QtQml/qqmlregistration.h>

#include "interfacemetrics.h"
#include "settingsstore.h"
#include "systemappearance.h"
#include "theme.h"
#include "typography.h"
#include "uiservices.h"

// Which token objects QML reaches as `Kvit.Ui` module singletons.
//
// Every declaration here is a QML_FOREIGN wrapper rather than a macro on the
// class itself, for two reasons.
//
// The first is layering: kvit-ui-tokens stays free of QML registration, so
// the list of what a shell can reach is this one file, and the token classes
// link nothing from QtQml.
//
// The second is a trap. Qt chooses how to construct a singleton in
// QQmlPrivate::singletonConstructionMode(), and it tests
// std::is_default_constructible<T> BEFORE it looks for a create() factory.
// Every class here takes `QObject *parent = nullptr`, so putting
// QML_SINGLETON and a create() on the class itself gets the factory silently
// ignored: Qt default-constructs its own instance, QML gets a valid object
// wired to nothing, and there is no warning anywhere. A view bound to such a
// Theme simply draws the built-in light table for ever, whatever the
// application's setting says.
//
// The FactoryWrapper branch is tested first, before default-constructibility,
// so a foreign wrapper carrying the create() is honoured. That is what these
// are.

// Each entry is (C++ class, QML name, DefaultServices accessor).
//
// The QML names are Theme, Interface, Typography and AppSettings, unchanged
// from what kvit-notes registers in its own `Kvit` module. 3,702 call sites
// across the two editors already spell them that way, so the migration in
// Waves 2 and 5 is an import line rather than a rename.
#define KVIT_UI_SINGLETONS(X)                                     \
    X(Theme, Theme, theme)                                        \
    X(Typography, Typography, typography)                         \
    X(InterfaceMetrics, Interface, interfaceMetrics)              \
    X(SystemAppearance, SystemAppearance, systemAppearance)       \
    X(SettingsStore, AppSettings, settings)

#define KVIT_UI_SINGLETON_WRAPPER(Type, Name, Accessor)                       \
    struct Name##Foreign                                                      \
    {                                                                         \
        Q_GADGET                                                              \
        QML_FOREIGN(Type)                                                     \
        QML_NAMED_ELEMENT(Name)                                               \
        QML_SINGLETON                                                         \
    public:                                                                   \
        static Type *create(QQmlEngine *engine, QJSEngine *)                  \
        {                                                                     \
            return KvitUi::singleton<Type>(                                   \
                engine, &KvitUi::DefaultServices::Accessor);                  \
        }                                                                     \
    };

KVIT_UI_SINGLETONS(KVIT_UI_SINGLETON_WRAPPER)

#endif // KVIT_UI_UISINGLETONS_H
