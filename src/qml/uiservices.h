// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_UISERVICES_H
#define KVIT_UI_UISERVICES_H

#include <QHash>
#include <QObject>
#include <QQmlEngine>
#include <QStringList>

class InterfaceMetrics;
class SettingsStore;
class SystemAppearance;
class Theme;
class Typography;

// How a `Kvit.Ui` QML singleton finds the objects it should be.
//
// The obvious arrangement — qmlRegisterSingletonInstance with one pointer —
// binds a single object for the whole process, and kvit-notes cannot use it:
// its services are members of an AppContext, and a process holds more than
// one, because the shell has its own and each test builds its own to stay
// isolated. This repository has to fit that, since kvit-notes is the first
// consumer.
//
// Qt resolves a QML_SINGLETON through a static create(QQmlEngine *,
// QJSEngine *) once per engine, which is the seam. A composition fills a
// table with its own token objects and hangs it on the engine; each
// singleton's create() reads its instance back out. One engine, one
// composition.
//
// The table is keyed by QMetaObject rather than by name, so the registering
// side and the reading side are tied together by the type and a mismatch
// cannot compile.
namespace KvitUi {

class ServiceTable
{
public:
    template <typename T>
    void add(T *instance)
    {
        m_services.insert(&T::staticMetaObject, instance);
    }

    QObject *lookup(const QMetaObject *type) const
    {
        return m_services.value(type, nullptr);
    }

private:
    QHash<const QMetaObject *, QObject *> m_services;
};

// Publish `table` on `engine`. The table must outlive the engine; whatever
// owns the composition owns both ends of that.
void attachServices(QQmlEngine *engine, ServiceTable *table);

// The table attached to `engine`, or nullptr when nothing attached one.
const ServiceTable *services(const QQmlEngine *engine);

// The token objects an engine gets when nothing attached a composition.
//
// kvit-notes and kvit-hub both attach one, because both hold their own Theme
// and their own settings file. The gallery, a single-component preview and a
// small application that wants the library's defaults do not, and requiring
// them to build a composition before `import Kvit.Ui` resolves would make the
// simplest use of this repository the most ceremonious. So there is a
// process-wide fallback, constructed on first use and bound to no settings
// file, which means it is session-scoped: a theme chosen through it is
// forgotten when the process exits.
//
// Call setSettings() on it to give it somewhere to persist to. An application
// that wants more than one composition should attach tables instead and never
// touch this.
class DefaultServices
{
public:
    static Theme *theme();
    static Typography *typography();
    static InterfaceMetrics *interfaceMetrics();
    static SystemAppearance *systemAppearance();
    static SettingsStore *settings();

    // Bind every token object above to `path` and load what is already in it.
    // Safe to call once, at startup, before the first window is shown.
    static void openSettings(const QString &path);
};

// The instance of T this engine's composition owns, falling back to the
// process default.
//
// A null return would reach QML as a singleton whose every member is
// undefined, which shows up as a wall of load warnings rather than as a
// crash. The fallback is what keeps that from being the out-of-the-box
// experience.
template <typename T>
T *singleton(QQmlEngine *engine, T *(*fallback)())
{
    T *instance = nullptr;
    if (const ServiceTable *table = services(engine))
        instance = qobject_cast<T *>(table->lookup(&T::staticMetaObject));
    if (!instance)
        instance = fallback();
    // The composition owns these, not the engine. Without this the garbage
    // collector would take a singleton it believes it created and leave the
    // owner holding a dangling member.
    QQmlEngine::setObjectOwnership(instance, QQmlEngine::CppOwnership);
    return instance;
}

// The QML names of every singleton the `Kvit.Ui` module registers, generated
// from the same list that declares them (KVIT_UI_SINGLETONS in
// uisingletons.h). A consumer with an extension namespace of its own reserves
// these so a module name cannot be confusable with one of them.
QStringList singletonNames();

}   // namespace KvitUi

#endif // KVIT_UI_UISERVICES_H
