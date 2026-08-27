// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef SETTINGSSTORE_H
#define SETTINGSSTORE_H

#include <QObject>
#include <QJsonObject>
#include <QTimer>
#include <QVariant>

// Per-user application settings: one flat
// JSON object in settings.json under the app config location, written
// atomically through a debounce so slider drags do not grind the disk.
// The path is injected (open()), which keeps tests hermetic. Unknown
// keys survive round-trips — the same tolerance rule collection.json
// follows — because the whole object is held and written back.
//
// App settings are per-user, not per-collection: theme, typography,
// toolbar layout, option states. Collection-shaped state (folder
// expansion, manual note order) stays in collection.json.
class SettingsStore : public QObject
{
    Q_OBJECT
    Q_PROPERTY(int revision READ revision NOTIFY revisionChanged)

public:
    // Write debounce: coalesces bursts (a slider drag) into one write.
    static constexpr int WriteDelayMs = 500;

    explicit SettingsStore(QObject *parent = nullptr);
    ~SettingsStore() override;  // flushes a pending write

    // Binds the store to its file and loads it. An absent or corrupt
    // file yields defaults (empty store); the file is (re)created on
    // the first write.
    //
    // Returns false, and leaves the store bound to its previous file with
    // its values and pending write intact, when the path's directory cannot
    // be created or when a pending write to the previous file fails. The
    // second case is the one worth stating: rebinding would drop values the
    // store has already told the caller it is holding, which contradicts the
    // write-failure contract everywhere else. A caller that would rather
    // lose them than stay bound passes discardPendingWrite.
    Q_INVOKABLE bool open(const QString &filePath,
                          bool discardPendingWrite = false);
    Q_INVOKABLE QString filePath() const { return m_filePath; }

    Q_INVOKABLE QVariant value(const QString &key,
                               const QVariant &fallback = QVariant()) const;
    Q_INVOKABLE void setValue(const QString &key, const QVariant &value);
    Q_INVOKABLE void remove(const QString &key);
    Q_INVOKABLE bool contains(const QString &key) const;

    // Write now if anything is pending (quit paths and tests).
    Q_INVOKABLE void flush();

    int revision() const { return m_revision; }

signals:
    // Emitted per real change; consumers watching one key filter here.
    void valueChanged(const QString &key);
    void revisionChanged();
    // A write did not reach disk (read-only location, full disk). The
    // values stay pending and are retried on the next flush or change, so
    // this is a warning the user can act on, not a report of lost data.
    void writeFailed(const QString &filePath, const QString &error);

private:
    void bumpRevision();
    void scheduleWrite();
    // False when the bytes did not reach disk; the values then stay pending.
    bool writeFile();
    void reportWriteFailure(const QString &error);

    QString m_filePath;
    QJsonObject m_values;
    QTimer m_writeTimer;
    bool m_dirty = false;
    int m_revision = 0;
};

#endif // SETTINGSSTORE_H
