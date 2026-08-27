// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "systemappearance.h"

#include <QCoreApplication>
#include <QProcess>
#include <QString>

#ifdef Q_OS_WIN
#include <QAbstractNativeEventFilter>
#include <QByteArray>
#include <windows.h>
#endif

namespace {

#if defined(Q_OS_MACOS) || defined(Q_OS_LINUX)
// Run a short command and return its trimmed standard output, or an empty
// string if it is missing, slow or fails.
//
// A command rather than a library binding on purpose: `gsettings` and
// `defaults` are the documented way to read these preferences, and reading
// them this way keeps the module free of a GIO or Foundation dependency for
// two booleans. The timeout is what makes it safe to call during startup —
// a desktop without the tool answers by not being there, and a desktop with
// a wedged settings daemon answers by timing out rather than by hanging the
// application.
QString commandOutput(const QString &program, const QStringList &arguments)
{
    QProcess process;
    process.setProgram(program);
    process.setArguments(arguments);
    process.setProcessChannelMode(QProcess::SeparateChannels);
    process.start();
    if (!process.waitForStarted(300))
        return QString();
    if (!process.waitForFinished(500)) {
        process.kill();
        process.waitForFinished(200);
        return QString();
    }
    if (process.exitStatus() != QProcess::NormalExit || process.exitCode() != 0)
        return QString();
    return QString::fromUtf8(process.readAllStandardOutput()).trimmed();
}
#endif

#ifdef Q_OS_WIN
// Windows announces a change to either preference with WM_SETTINGCHANGE,
// which arrives as a native event on the application's own message loop.
// The filter forwards it and nothing else; deciding what changed is
// refresh()'s job, which re-reads both values and only notifies if one moved.
class SettingChangeFilter : public QObject, public QAbstractNativeEventFilter
{
public:
    explicit SettingChangeFilter(SystemAppearance *owner)
        : QObject(owner), m_owner(owner)
    {
        QCoreApplication::instance()->installNativeEventFilter(this);
    }
    ~SettingChangeFilter() override
    {
        if (QCoreApplication::instance())
            QCoreApplication::instance()->removeNativeEventFilter(this);
    }

    bool nativeEventFilter(const QByteArray &eventType, void *message,
                           qintptr *) override
    {
        if (eventType != "windows_generic_MSG")
            return false;
        const MSG *msg = static_cast<MSG *>(message);
        if (msg && msg->message == WM_SETTINGCHANGE)
            m_owner->refresh();
        return false;   // never consume it: other code watches it too
    }

private:
    SystemAppearance *m_owner;
};
#endif

} // namespace

SystemAppearance::SystemAppearance(QObject *parent)
    : QObject(parent)
{
#if defined(Q_OS_WIN) || defined(Q_OS_MACOS) || defined(Q_OS_LINUX)
    m_available = true;
#endif
    refresh();
    installPlatformWatch();
}

SystemAppearance::~SystemAppearance()
{
    removePlatformWatch();
}

void SystemAppearance::refresh()
{
    if (m_overridden)
        return;
    bool high = m_highContrast;
    bool reduced = m_reducedMotion;
    readPlatform(&high, &reduced);
    if (high == m_highContrast && reduced == m_reducedMotion)
        return;
    m_highContrast = high;
    m_reducedMotion = reduced;
    emit changed();
}

void SystemAppearance::setOverride(bool highContrast, bool reducedMotion)
{
    m_overridden = true;
    m_available = true;
    if (highContrast == m_highContrast && reducedMotion == m_reducedMotion)
        return;
    m_highContrast = highContrast;
    m_reducedMotion = reducedMotion;
    emit changed();
}

#if defined(Q_OS_WIN)

void SystemAppearance::readPlatform(bool *highContrast,
                                    bool *reducedMotion) const
{
    HIGHCONTRASTW hc {};
    hc.cbSize = sizeof(hc);
    if (SystemParametersInfoW(SPI_GETHIGHCONTRAST, sizeof(hc), &hc, 0))
        *highContrast = (hc.dwFlags & HCF_HIGHCONTRASTON) != 0;

    // SPI_GETCLIENTAREAANIMATION answers the question the other way round:
    // TRUE means animations are wanted.
    BOOL animations = TRUE;
    if (SystemParametersInfoW(SPI_GETCLIENTAREAANIMATION, 0, &animations, 0))
        *reducedMotion = !animations;
}

void SystemAppearance::installPlatformWatch()
{
    if (QCoreApplication::instance())
        m_watch = new SettingChangeFilter(this);
}

void SystemAppearance::removePlatformWatch()
{
    delete m_watch;
    m_watch = nullptr;
}

#elif defined(Q_OS_MACOS)

void SystemAppearance::readPlatform(bool *highContrast,
                                    bool *reducedMotion) const
{
    // The two preferences NSWorkspace publishes as
    // accessibilityDisplayShouldIncreaseContrast and
    // accessibilityDisplayShouldReduceMotion, read from the defaults they are
    // backed by. Reading them this way keeps the module in plain C++; the
    // cost is that a change made while the application is running is picked
    // up on the next refresh() rather than announced, which is why the
    // settings dialog refreshes when it opens.
    const QString contrast = commandOutput(
        QStringLiteral("defaults"),
        {QStringLiteral("read"), QStringLiteral("com.apple.universalaccess"),
         QStringLiteral("increaseContrast")});
    if (!contrast.isEmpty())
        *highContrast = (contrast == QLatin1String("1"));

    const QString motion = commandOutput(
        QStringLiteral("defaults"),
        {QStringLiteral("read"), QStringLiteral("com.apple.Accessibility"),
         QStringLiteral("ReduceMotionEnabled")});
    if (!motion.isEmpty())
        *reducedMotion = (motion == QLatin1String("1"));
}

void SystemAppearance::installPlatformWatch() {}
void SystemAppearance::removePlatformWatch() {}

#elif defined(Q_OS_LINUX)

void SystemAppearance::readPlatform(bool *highContrast,
                                    bool *reducedMotion) const
{
    // GNOME's two keys, which KDE and the other desktops that follow the
    // freedesktop settings conventions also honour. As on macOS, a change
    // made mid-session is picked up on the next refresh().
    const QString contrast = commandOutput(
        QStringLiteral("gsettings"),
        {QStringLiteral("get"),
         QStringLiteral("org.gnome.desktop.a11y.interface"),
         QStringLiteral("high-contrast")});
    if (!contrast.isEmpty())
        *highContrast = (contrast == QLatin1String("true"));

    // enable-animations reads the other way round, as SPI_GETCLIENTAREAANIMATION
    // does on Windows.
    const QString animations = commandOutput(
        QStringLiteral("gsettings"),
        {QStringLiteral("get"), QStringLiteral("org.gnome.desktop.interface"),
         QStringLiteral("enable-animations")});
    if (!animations.isEmpty())
        *reducedMotion = (animations == QLatin1String("false"));
}

void SystemAppearance::installPlatformWatch() {}
void SystemAppearance::removePlatformWatch() {}

#else

void SystemAppearance::readPlatform(bool *, bool *) const {}
void SystemAppearance::installPlatformWatch() {}
void SystemAppearance::removePlatformWatch() {}

#endif
