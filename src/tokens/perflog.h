// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef PERFLOG_H
#define PERFLOG_H

#include <QObject>
#include <QElapsedTimer>
#include <QFile>
#include <QHash>
#include <QJsonObject>
#include <QMutex>
#include <QVariantMap>

#include <atomic>
#include <optional>

class PerfLog : public QObject
{
    Q_OBJECT
    Q_PROPERTY(bool enabled READ enabled WRITE setEnabled NOTIFY enabledChanged)
    Q_PROPERTY(QString levelName READ levelName NOTIFY levelChanged)

public:
    enum Level {
        Off = 0,
        Major = 1,
        Verbose = 2,
    };
    Q_ENUM(Level)

    struct Sample {
        qint64 sinceStartMs = 0;
        QString operation;
        double durationMs = 0.0;
        QVariantMap context;
    };

    class ScopedTimer
    {
    public:
        explicit ScopedTimer(const QString &operation,
                             Level minLevel = Major,
                             double thresholdMs = 0.0);
        ScopedTimer(const QString &operation,
                    const QVariantMap &context,
                    Level minLevel = Major,
                    double thresholdMs = 0.0);
        ~ScopedTimer();

        ScopedTimer(const ScopedTimer &) = delete;
        ScopedTimer &operator=(const ScopedTimer &) = delete;

        void addContext(const QString &key, const QVariant &value);
        void setContext(const QVariantMap &context);

    private:
        void startIfEnabled();

        QString m_operation;
        QVariantMap m_context;
        Level m_minLevel = Major;
        double m_thresholdMs = 0.0;
        std::optional<QElapsedTimer> m_timer;
    };

    static PerfLog &instance();

    static bool isEnabled();
    static bool isVerbose();
    static bool levelAllows(Level level);

    bool enabled() const;
    Level level() const;
    QString levelName() const;

    void configureFromEnvironment();
    bool hasEnvironmentOverride() const { return m_hasEnvironmentOverride; }
    bool hasLogFilePath() const;
    void configureFromSetting(const QVariant &value);

    Q_INVOKABLE void setEnabled(bool enabled);
    void setLevel(Level level);
    void setHumanOutput(bool human);
    void setEmitToStderr(bool enabled);
    void setLogFilePath(const QString &path);

    Q_INVOKABLE void begin(const QString &id,
                           const QVariantMap &context = QVariantMap());
    Q_INVOKABLE void end(const QString &id,
                         const QVariantMap &context = QVariantMap());
    Q_INVOKABLE void mark(const QString &operation,
                          double durationMs,
                          const QVariantMap &context = QVariantMap());
    Q_INVOKABLE void markIfAbove(const QString &operation,
                                 double durationMs,
                                 double thresholdMs,
                                 const QVariantMap &context = QVariantMap());

    void record(const QString &operation,
                double durationMs,
                const QVariantMap &context = QVariantMap(),
                Level minLevel = Major,
                double thresholdMs = 0.0);

    QList<Sample> samples(const QString &operation = QString()) const;
    Q_INVOKABLE QVariantList sampleMaps(const QString &operation = QString()) const;
    Q_INVOKABLE void clear();

    static QJsonObject sampleToJsonObject(const Sample &sample);
    static QByteArray sampleToJsonLine(const Sample &sample);

    static int timerConstructionCount();
    static void resetTimerConstructionCount();

signals:
    void enabledChanged();
    void levelChanged();
    void sampleRecorded(const QString &operation,
                        double durationMs,
                        const QVariantMap &context);

private:
    explicit PerfLog(QObject *parent = nullptr);

    struct Pending {
        QElapsedTimer timer;
        QVariantMap context;
    };

    static Level parseLevel(const QByteArray &value, bool *recognized = nullptr,
                            bool *human = nullptr);
    static QString levelToString(Level level);

    // In-memory retention window. Samples are a debugging aid read back
    // through samples()/sampleMaps(); the complete record is the JSONL file.
    static constexpr int kMaxRetainedSamples = 5000;
    static constexpr int kRetentionSlack = 1000;

    void writeSample(const Sample &sample);
    void openFileIfNeeded();
    void rotateFileIfNeeded(const QString &path);

    mutable QMutex m_mutex;
    QElapsedTimer m_sinceStart;
    QList<Sample> m_samples;
    QHash<QString, Pending> m_pending;
    QFile m_file;
    QString m_logFilePath;
    bool m_hasEnvironmentOverride = false;
    bool m_emitToStderr = true;
    bool m_humanOutput = false;

    static std::atomic<int> s_level;
    static std::atomic<int> s_timerConstructionCount;
};

#define KVIT_PERF_CONCAT_INNER(a, b) a##b
#define KVIT_PERF_CONCAT(a, b) KVIT_PERF_CONCAT_INNER(a, b)
#define PERF_SCOPE(operation) \
    PerfLog::ScopedTimer KVIT_PERF_CONCAT(_kvitPerfScope, __LINE__)(operation)
#define PERF_SCOPE_CTX(operation, context) \
    PerfLog::ScopedTimer KVIT_PERF_CONCAT(_kvitPerfScope, __LINE__)(operation, context)
#define PERF_VERBOSE_SCOPE_THRESHOLD(operation, thresholdMs) \
    PerfLog::ScopedTimer KVIT_PERF_CONCAT(_kvitPerfScope, __LINE__)( \
        operation, PerfLog::Verbose, thresholdMs)

#endif // PERFLOG_H
