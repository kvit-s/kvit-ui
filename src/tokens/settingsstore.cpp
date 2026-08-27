// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "settingsstore.h"

#include <QDir>
#include <QFile>
#include <QFileInfo>
#include <QJsonDocument>
#include <QJsonParseError>
#include <QJsonValue>
#include <QJSValue>
#include <QSaveFile>

// QML passes JS arrays/objects into QVariant parameters as wrapped
// QJSValues; QJsonValue::fromVariant would silently store those as
// null. Unwrap to plain variants at the boundary.
static QVariant unwrapped(const QVariant &value)
{
    if (value.userType() == qMetaTypeId<QJSValue>())
        return value.value<QJSValue>().toVariant();
    return value;
}

SettingsStore::SettingsStore(QObject *parent)
    : QObject(parent)
{
    m_writeTimer.setSingleShot(true);
    m_writeTimer.setInterval(WriteDelayMs);
    connect(&m_writeTimer, &QTimer::timeout, this, &SettingsStore::writeFile);
}

SettingsStore::~SettingsStore()
{
    if (m_dirty)
        writeFile();
}

bool SettingsStore::open(const QString &filePath, bool discardPendingWrite)
{
    // A pending write belongs to the previous path; land it there first. If
    // it does not land, the values are still the only copy that exists, and
    // replacing the path here would discard them -- the outcome the write
    // path goes out of its way to avoid. Stay where we are and report the
    // refusal; the caller retries once the location is writable, or asks
    // explicitly to drop the values.
    if (m_dirty && !writeFile() && !discardPendingWrite)
        return false;

    // Nothing is given up until the new location is known to be usable. The
    // binding, the values and the dirty flag used to be cleared first and this
    // check made afterwards, so a directory that could not be created left the
    // store holding nothing and bound to a path it can neither read nor write
    // — the outcome the second half of this function's contract exists to
    // rule out, and the same rule the pending-write refusal above follows.
    const QDir dir = QFileInfo(filePath).absoluteDir();
    if (!dir.exists() && !dir.mkpath(QStringLiteral(".")))
        return false;

    m_writeTimer.stop();
    m_filePath = filePath;
    m_values = QJsonObject();
    m_dirty = false;

    QFile file(filePath);
    if (file.open(QIODevice::ReadOnly)) {
        QJsonParseError error;
        const QJsonDocument doc =
            QJsonDocument::fromJson(file.readAll(), &error);
        if (error.error == QJsonParseError::NoError && doc.isObject())
            m_values = doc.object();
        // Corrupt or non-object: defaults; the file is rewritten on the
        // first setValue.
    }

    bumpRevision();
    return true;
}

QVariant SettingsStore::value(const QString &key, const QVariant &fallback) const
{
    const auto it = m_values.constFind(key);
    if (it == m_values.constEnd())
        return unwrapped(fallback);
    return it->toVariant();
}

void SettingsStore::setValue(const QString &key, const QVariant &value)
{
    const QJsonValue jsonValue = QJsonValue::fromVariant(unwrapped(value));
    const auto it = m_values.constFind(key);
    if (it != m_values.constEnd() && *it == jsonValue)
        return;

    m_values.insert(key, jsonValue);
    scheduleWrite();
    emit valueChanged(key);
    bumpRevision();
}

void SettingsStore::remove(const QString &key)
{
    if (!m_values.contains(key))
        return;

    m_values.remove(key);
    scheduleWrite();
    emit valueChanged(key);
    bumpRevision();
}

bool SettingsStore::contains(const QString &key) const
{
    return m_values.contains(key);
}

void SettingsStore::flush()
{
    if (m_dirty)
        writeFile();
}

void SettingsStore::bumpRevision()
{
    ++m_revision;
    emit revisionChanged();
}

void SettingsStore::scheduleWrite()
{
    m_dirty = true;
    m_writeTimer.start();  // restarts on every change: bursts coalesce
}

bool SettingsStore::writeFile()
{
    m_writeTimer.stop();
    if (m_filePath.isEmpty()) {
        m_dirty = false; // nowhere to write: nothing is being lost
        return true;
    }

    // The dirty flag is only cleared once the bytes are committed. Clearing
    // it up front loses the change outright when the disk is full or the
    // location is read-only: the next flush sees nothing pending and the
    // user's preference is gone with no trace.
    const QByteArray bytes =
        QJsonDocument(m_values).toJson(QJsonDocument::Indented);
    QSaveFile file(m_filePath);
    if (!file.open(QIODevice::WriteOnly)) {
        reportWriteFailure(file.errorString());
        return false;
    }
    if (file.write(bytes) != bytes.size()) {
        const QString error = file.errorString();
        file.cancelWriting();
        reportWriteFailure(error);
        return false;
    }
    if (!file.commit()) {
        reportWriteFailure(file.errorString());
        return false;
    }
    m_dirty = false;
    return true;
}

void SettingsStore::reportWriteFailure(const QString &error)
{
    // The change stays pending so the next flush, the next setting change,
    // or the destructor retries it. No timer is restarted here: a location
    // that cannot be written usually stays that way, and a self-restarting
    // timer would spin failing writes for the life of the process.
    m_dirty = true;
    emit writeFailed(m_filePath, error);
}
