// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "uiservices.h"

#include <QCoreApplication>

#include "interfacemetrics.h"
#include "settingsstore.h"
#include "systemappearance.h"
#include "theme.h"
#include "typography.h"
#include "uisingletons.h"

namespace {
// The engine property the table hangs on. A property rather than a static
// map because the lifetime that matters is the engine's.
const char *kServicesProperty = "_kvitUiServices";
}   // namespace

namespace KvitUi {

void attachServices(QQmlEngine *engine, ServiceTable *table)
{
    if (engine)
        engine->setProperty(kServicesProperty, QVariant::fromValue(
                                static_cast<void *>(table)));
}

const ServiceTable *services(const QQmlEngine *engine)
{
    if (!engine)
        return nullptr;
    const QVariant value = engine->property(kServicesProperty);
    if (!value.isValid())
        return nullptr;
    return static_cast<const ServiceTable *>(value.value<void *>());
}

namespace {

// The process-wide fallback composition, built on first use.
//
// Constructed in dependency order and wired the way an application would wire
// it: the desktop's appearance drives Theme's "system" resolution, and all
// three token objects share one settings store so an interface size and a
// theme choice persist together.
struct Default
{
    Default()
    {
        systemAppearance = new SystemAppearance;
        theme = new Theme;
        typography = new Typography;
        interfaceMetrics = new InterfaceMetrics;
        settings = new SettingsStore;
        theme->setSystemAppearance(systemAppearance);
    }

    ~Default()
    {
        // Ordered against construction: Theme holds a pointer to both of the
        // objects it was given, so it goes first.
        delete theme;
        delete typography;
        delete interfaceMetrics;
        delete systemAppearance;
        delete settings;
    }

    SystemAppearance *systemAppearance = nullptr;
    Theme *theme = nullptr;
    Typography *typography = nullptr;
    InterfaceMetrics *interfaceMetrics = nullptr;
    SettingsStore *settings = nullptr;
};

Default &defaults()
{
    // Function-local static: constructed on the first singleton lookup, which
    // is after QGuiApplication exists, which is what SystemAppearance needs.
    static Default instance;
    return instance;
}

}   // namespace

Theme *DefaultServices::theme() { return defaults().theme; }
Typography *DefaultServices::typography() { return defaults().typography; }
InterfaceMetrics *DefaultServices::interfaceMetrics()
{
    return defaults().interfaceMetrics;
}
SystemAppearance *DefaultServices::systemAppearance()
{
    return defaults().systemAppearance;
}
SettingsStore *DefaultServices::settings() { return defaults().settings; }

void DefaultServices::openSettings(const QString &path)
{
    Default &d = defaults();
    if (!d.settings->open(path))
        return;
    d.theme->setSettings(d.settings);
    d.typography->setSettings(d.settings);
    d.interfaceMetrics->setSettings(d.settings);
}

QStringList singletonNames()
{
#define KVIT_UI_SINGLETON_NAME(Type, Name, Accessor) QStringLiteral(#Name),
    return QStringList{ KVIT_UI_SINGLETONS(KVIT_UI_SINGLETON_NAME) };
#undef KVIT_UI_SINGLETON_NAME
}

}   // namespace KvitUi
