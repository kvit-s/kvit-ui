// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef SYSTEMAPPEARANCE_H
#define SYSTEMAPPEARANCE_H

#include <QObject>

// The two accessibility preferences the operating system knows about and Qt
// does not pass on (accessibility.md Finding 7).
//
// Qt already follows the light/dark preference through
// QStyleHints::colorScheme, which is what Theme's "system" setting reads.
// High contrast and reduced motion have no Qt equivalent, and each desktop
// exposes them its own way:
//
//   Windows  SystemParametersInfo(SPI_GETHIGHCONTRAST) and
//            SPI_GETCLIENTAREAANIMATION, with WM_SETTINGCHANGE announcing a
//            change.
//   macOS    NSWorkspace's accessibilityDisplayShouldIncreaseContrast and
//            accessibilityDisplayShouldReduceMotion, each with a notification.
//   Linux    the GSettings keys org.gnome.desktop.a11y.interface
//            high-contrast and org.gnome.desktop.interface enable-animations,
//            read through gsettings.
//
// Both properties are false where the platform gives no answer, which is the
// safe direction: a person who has turned neither on sees no change, and the
// application's own settings stay in charge. Neither property is ever written
// — this reports what the desktop says and nothing else.
//
// Of the four, Windows high contrast is the one that matters most: someone who
// turns it on system-wide expects every application to follow, and Kvit has a
// high-contrast theme already.
class SystemAppearance : public QObject
{
    Q_OBJECT

    // Whether the desktop is in a high-contrast mode.
    Q_PROPERTY(bool highContrast READ highContrast NOTIFY changed)
    // Whether the desktop asks applications to still their animations.
    Q_PROPERTY(bool reducedMotion READ reducedMotion NOTIFY changed)
    // Whether this platform answers at all. False on a desktop with no
    // readable settings, which is how the "follow the system" choices know to
    // present themselves as having nothing to follow.
    Q_PROPERTY(bool available READ available CONSTANT)

public:
    explicit SystemAppearance(QObject *parent = nullptr);
    ~SystemAppearance() override;

    bool highContrast() const { return m_highContrast; }
    bool reducedMotion() const { return m_reducedMotion; }
    bool available() const { return m_available; }

    // Re-read both values from the platform and notify if either moved.
    // Called by the platform's own change notification, and by tests.
    Q_INVOKABLE void refresh();

    // Force both values, for the tests: no CI machine has a high-contrast
    // desktop to turn on, and the behaviour that matters — the theme
    // resolving to highContrast, motion stilling — is on this side of the
    // platform boundary. Setting either marks the object as overridden, so a
    // later refresh() from the platform cannot undo the test's setup.
    void setOverride(bool highContrast, bool reducedMotion);

signals:
    void changed();

private:
    // Reads the platform, or leaves both alone where there is nothing to
    // read. Defined once per platform in the .cpp behind #ifdef.
    void readPlatform(bool *highContrast, bool *reducedMotion) const;
    void installPlatformWatch();
    void removePlatformWatch();

    bool m_highContrast = false;
    bool m_reducedMotion = false;
    bool m_available = false;
    bool m_overridden = false;
    // The platform's change notification, where it needs a live object: the
    // Windows message window, or the GSettings monitor process.
    QObject *m_watch = nullptr;
};

#endif // SYSTEMAPPEARANCE_H
