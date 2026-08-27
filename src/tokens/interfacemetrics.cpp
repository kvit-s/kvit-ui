// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "interfacemetrics.h"

#include <QFont>
#include <QGuiApplication>
#include <QVariant>

#include "settingsstore.h"
#include "perflog.h"

namespace {
const QString kFontSize = QStringLiteral("interface.fontSize");
const QString kFontFamily = QStringLiteral("interface.fontFamily");
const QString kMonoFamily = QStringLiteral("interface.monoFamily");
} // namespace

InterfaceMetrics::InterfaceMetrics(QObject *parent)
    : QObject(parent)
{
}

void InterfaceMetrics::setSettings(SettingsStore *settings)
{
    m_settings = settings;
    if (!m_settings)
        return;
    // Loaded through the setter so a persisted value clamps exactly as a live
    // one does; m_loading suppresses writing back what was just read.
    m_loading = true;
    setFontSize(m_settings->value(kFontSize, m_fontSize).toInt());
    setFontFamily(m_settings->value(kFontFamily, m_fontFamily).toString());
    setMonoFamily(m_settings->value(kMonoFamily, m_monoFamily).toString());
    m_loading = false;
    emit changed();
}

void InterfaceMetrics::setFontSize(int size)
{
    size = qBound(MinFontSize, size, MaxFontSize);
    if (m_fontSize == size)
        return;
    // A chrome resize relays out every pane at once, which is the same class
    // of event as a typography reflow and is measured the same way.
    PerfLog::ScopedTimer perf(
        QStringLiteral("interface.reflow"),
        QVariantMap{{QStringLiteral("value"), size}},
        m_loading ? PerfLog::Verbose : PerfLog::Major);
    m_fontSize = size;
    if (m_settings && !m_loading)
        m_settings->setValue(kFontSize, m_fontSize);
    emit changed();
}

void InterfaceMetrics::setFontFamily(const QString &family)
{
    if (m_fontFamily == family)
        return;
    m_fontFamily = family;
    if (m_settings && !m_loading)
        m_settings->setValue(kFontFamily, m_fontFamily);
    emit changed();
}

void InterfaceMetrics::setMonoFamily(const QString &family)
{
    // Empty would leave an identifier in the proportional family, where two
    // references down a column no longer line up, so it falls back to the
    // generic name every platform resolves to something fixed-pitch.
    const QString wanted = family.isEmpty() ? QStringLiteral("monospace")
                                            : family;
    if (m_monoFamily == wanted)
        return;
    m_monoFamily = wanted;
    if (m_settings && !m_loading)
        m_settings->setValue(kMonoFamily, m_monoFamily);
    emit changed();
}

QString InterfaceMetrics::resolvedFontFamily() const
{
    if (!m_fontFamily.isEmpty())
        return m_fontFamily;
    // The desktop's own default, asked for by name rather than left empty.
    // QGuiApplication::font() is what the platform theme decided, and on a
    // Linux desktop that is the generic "Sans Serif", which fontconfig
    // resolves; an empty family is resolved against nothing.
    return QGuiApplication::font().family();
}

QString InterfaceMetrics::resolvedMonoFamily() const
{
    // setMonoFamily already refuses to store an empty value, so this is the
    // stored one unless something bypassed the setter — a hand-edited
    // settings file loaded before setSettings runs, for instance.
    return m_monoFamily.isEmpty() ? QStringLiteral("monospace") : m_monoFamily;
}

int InterfaceMetrics::px(int designPx) const
{
    if (designPx == 0)
        return 0;
    const int scaled = qRound(designPx * scale());
    // A one-pixel rule must not round away to nothing: a separator scaled to
    // zero is a line that vanishes at the smallest interface size, which
    // reads as a layout bug rather than as a smaller interface.
    if (designPx > 0)
        return qMax(1, scaled);
    return qMin(-1, scaled);
}

void InterfaceMetrics::resetToDefaults()
{
    setFontSize(DefaultFontSize);
    setFontFamily(QString());
    setMonoFamily(QStringLiteral("monospace"));
}
