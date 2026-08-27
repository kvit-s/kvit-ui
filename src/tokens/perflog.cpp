// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "perflog.h"

#include <QCoreApplication>
#include <QDateTime>
#include <QDir>
#include <QFileInfo>
#include <QJsonDocument>
#include <QJsonValue>
#include <QTextStream>

#include <cstdio>

std::atomic<int> PerfLog::s_level{PerfLog::Off};
std::atomic<int> PerfLog::s_timerConstructionCount{0};

PerfLog::ScopedTimer::ScopedTimer(const QString &operation,
                                  Level minLevel,
                                  double thresholdMs)
    : m_operation(operation)
    , m_minLevel(minLevel)
    , m_thresholdMs(thresholdMs)
{
    startIfEnabled();
}

PerfLog::ScopedTimer::ScopedTimer(const QString &operation,
                                  const QVariantMap &context,
                                  Level minLevel,
                                  double thresholdMs)
    : m_operation(operation)
    , m_context(context)
    , m_minLevel(minLevel)
    , m_thresholdMs(thresholdMs)
{
    startIfEnabled();
}

PerfLog::ScopedTimer::~ScopedTimer()
{
    if (!m_timer)
        return;

    const double elapsedMs = double(m_timer->nsecsElapsed()) / 1000000.0;
    PerfLog::instance().record(m_operation, elapsedMs, m_context,
                               m_minLevel, m_thresholdMs);
}

void PerfLog::ScopedTimer::addContext(const QString &key, const QVariant &value)
{
    if (m_timer)
        m_context.insert(key, value);
}

void PerfLog::ScopedTimer::setContext(const QVariantMap &context)
{
    if (m_timer)
        m_context = context;
}

void PerfLog::ScopedTimer::startIfEnabled()
{
    if (!PerfLog::levelAllows(m_minLevel))
        return;
    m_timer.emplace();
    s_timerConstructionCount.fetch_add(1, std::memory_order_relaxed);
    m_timer->start();
}

PerfLog &PerfLog::instance()
{
    static PerfLog log;
    return log;
}

bool PerfLog::isEnabled()
{
    return s_level.load(std::memory_order_relaxed) != Off;
}

bool PerfLog::isVerbose()
{
    return s_level.load(std::memory_order_relaxed) >= Verbose;
}

bool PerfLog::levelAllows(Level level)
{
    return s_level.load(std::memory_order_relaxed) >= int(level);
}

bool PerfLog::enabled() const
{
    return isEnabled();
}

PerfLog::Level PerfLog::level() const
{
    return static_cast<Level>(s_level.load(std::memory_order_relaxed));
}

QString PerfLog::levelName() const
{
    return levelToString(level());
}

void PerfLog::configureFromEnvironment()
{
    const QByteArray value = qgetenv("KVIT_PERF");
    const QByteArray path = qgetenv("KVIT_PERF_LOG");
    if (!path.isEmpty())
        setLogFilePath(QString::fromLocal8Bit(path));

    if (value.isEmpty()) {
        m_hasEnvironmentOverride = false;
        setLevel(Off);
        return;
    }

    bool recognized = false;
    bool human = false;
    const Level parsed = parseLevel(value, &recognized, &human);
    m_hasEnvironmentOverride = true;
    setHumanOutput(human);
    setLevel(recognized ? parsed : Major);
}

void PerfLog::configureFromSetting(const QVariant &value)
{
    if (!value.isValid() || value.isNull())
        return;

    if (value.typeId() == QMetaType::Bool) {
        setLevel(value.toBool() ? Major : Off);
        return;
    }

    bool recognized = false;
    bool human = false;
    const Level parsed = parseLevel(value.toString().toUtf8(),
                                    &recognized, &human);
    if (!recognized)
        return;
    setHumanOutput(human);
    setLevel(parsed);
}

void PerfLog::setEnabled(bool enabled)
{
    setLevel(enabled ? Major : Off);
}

void PerfLog::setLevel(Level level)
{
    const int old = s_level.exchange(int(level), std::memory_order_relaxed);
    if (old == int(level))
        return;
    emit enabledChanged();
    emit levelChanged();
}

void PerfLog::setHumanOutput(bool human)
{
    m_humanOutput = human;
}

void PerfLog::setEmitToStderr(bool enabled)
{
    QMutexLocker locker(&m_mutex);
    m_emitToStderr = enabled;
}

void PerfLog::setLogFilePath(const QString &path)
{
    QMutexLocker locker(&m_mutex);
    if (m_logFilePath == path)
        return;
    if (m_file.isOpen())
        m_file.close();
    m_logFilePath = path;
}

bool PerfLog::hasLogFilePath() const
{
    QMutexLocker locker(&m_mutex);
    return !m_logFilePath.isEmpty();
}

void PerfLog::begin(const QString &id, const QVariantMap &context)
{
    if (!isEnabled())
        return;
    Pending pending;
    pending.context = context;
    pending.timer.start();
    s_timerConstructionCount.fetch_add(1, std::memory_order_relaxed);

    QMutexLocker locker(&m_mutex);
    m_pending.insert(id, pending);
}

void PerfLog::end(const QString &id, const QVariantMap &context)
{
    if (!isEnabled())
        return;

    Pending pending;
    {
        QMutexLocker locker(&m_mutex);
        auto it = m_pending.find(id);
        if (it == m_pending.end())
            return;
        pending = it.value();
        m_pending.erase(it);
    }

    QVariantMap merged = pending.context;
    for (auto it = context.constBegin(); it != context.constEnd(); ++it)
        merged.insert(it.key(), it.value());

    mark(id, double(pending.timer.nsecsElapsed()) / 1000000.0, merged);
}

void PerfLog::mark(const QString &operation,
                   double durationMs,
                   const QVariantMap &context)
{
    record(operation, durationMs, context);
}

void PerfLog::markIfAbove(const QString &operation,
                          double durationMs,
                          double thresholdMs,
                          const QVariantMap &context)
{
    record(operation, durationMs, context, Major, thresholdMs);
}

void PerfLog::record(const QString &operation,
                     double durationMs,
                     const QVariantMap &context,
                     Level minLevel,
                     double thresholdMs)
{
    if (!levelAllows(minLevel))
        return;
    if (thresholdMs > 0.0 && durationMs < thresholdMs)
        return;

    Sample sample;
    sample.sinceStartMs = m_sinceStart.elapsed();
    sample.operation = operation;
    sample.durationMs = durationMs;
    sample.context = context;

    {
        QMutexLocker locker(&m_mutex);
        m_samples.append(sample);
        // Retention is a bounded window of the most recent samples. The JSONL
        // file rotates, so without this a long session with logging enabled
        // grows in memory without limit while the on-disk log stays capped.
        // Trimming only once the slack is used up keeps the cost amortized:
        // one bulk move per kRetentionSlack records rather than shifting the
        // whole list on every record.
        if (m_samples.size() > kMaxRetainedSamples + kRetentionSlack)
            m_samples.remove(0, m_samples.size() - kMaxRetainedSamples);
    }

    writeSample(sample);
    emit sampleRecorded(operation, durationMs, context);
}

QList<PerfLog::Sample> PerfLog::samples(const QString &operation) const
{
    QMutexLocker locker(&m_mutex);
    if (operation.isEmpty())
        return m_samples;

    QList<Sample> filtered;
    for (const Sample &sample : m_samples) {
        if (sample.operation == operation)
            filtered.append(sample);
    }
    return filtered;
}

QVariantList PerfLog::sampleMaps(const QString &operation) const
{
    QVariantList maps;
    for (const Sample &sample : samples(operation)) {
        maps.append(QVariantMap{
            {QStringLiteral("t"), sample.sinceStartMs},
            {QStringLiteral("op"), sample.operation},
            {QStringLiteral("ms"), sample.durationMs},
            {QStringLiteral("ctx"), sample.context},
        });
    }
    return maps;
}

void PerfLog::clear()
{
    QMutexLocker locker(&m_mutex);
    m_samples.clear();
    m_pending.clear();
}

QJsonObject PerfLog::sampleToJsonObject(const Sample &sample)
{
    QJsonObject object;
    object.insert(QStringLiteral("t"),
                  QJsonValue::fromVariant(QVariant::fromValue(sample.sinceStartMs)));
    object.insert(QStringLiteral("op"), sample.operation);
    object.insert(QStringLiteral("ms"), sample.durationMs);
    object.insert(QStringLiteral("ctx"),
                  QJsonObject::fromVariantMap(sample.context));
    return object;
}

QByteArray PerfLog::sampleToJsonLine(const Sample &sample)
{
    QByteArray line =
        QJsonDocument(sampleToJsonObject(sample)).toJson(QJsonDocument::Compact);
    line.append('\n');
    return line;
}

int PerfLog::timerConstructionCount()
{
    return s_timerConstructionCount.load(std::memory_order_relaxed);
}

void PerfLog::resetTimerConstructionCount()
{
    s_timerConstructionCount.store(0, std::memory_order_relaxed);
}

PerfLog::PerfLog(QObject *parent)
    : QObject(parent)
{
    m_sinceStart.start();
}

PerfLog::Level PerfLog::parseLevel(const QByteArray &value,
                                   bool *recognized,
                                   bool *human)
{
    const QByteArray normalized = value.trimmed().toLower();
    if (recognized)
        *recognized = true;
    if (human)
        *human = false;

    if (normalized.isEmpty() || normalized == "0" || normalized == "false"
        || normalized == "off" || normalized == "none") {
        return Off;
    }
    if (normalized == "1" || normalized == "true" || normalized == "on"
        || normalized == "major") {
        return Major;
    }
    if (normalized == "verbose" || normalized == "2")
        return Verbose;
    if (normalized == "human") {
        if (human)
            *human = true;
        return Major;
    }

    if (recognized)
        *recognized = false;
    return Off;
}

QString PerfLog::levelToString(Level level)
{
    switch (level) {
    case Off: return QStringLiteral("off");
    case Major: return QStringLiteral("major");
    case Verbose: return QStringLiteral("verbose");
    }
    return QStringLiteral("off");
}

void PerfLog::writeSample(const Sample &sample)
{
    const QByteArray jsonLine = sampleToJsonLine(sample);
    QByteArray humanLine;
    if (m_humanOutput) {
        humanLine = QStringLiteral("%1 ms  %2  %3 ms  %4\n")
                        .arg(sample.sinceStartMs, 8)
                        .arg(sample.operation, -28)
                        .arg(sample.durationMs, 8, 'f', 2)
                        .arg(QString::fromUtf8(QJsonDocument(
                                 QJsonObject::fromVariantMap(sample.context))
                                 .toJson(QJsonDocument::Compact)))
                        .toUtf8();
    }

    QMutexLocker locker(&m_mutex);
    if (m_emitToStderr) {
        const QByteArray &line = m_humanOutput ? humanLine : jsonLine;
        fwrite(line.constData(), 1, size_t(line.size()), stderr);
        fflush(stderr);
    }

    if (m_logFilePath.isEmpty())
        return;
    openFileIfNeeded();
    if (!m_file.isOpen())
        return;
    m_file.write(jsonLine);
    m_file.flush();
}

void PerfLog::openFileIfNeeded()
{
    if (m_file.isOpen())
        return;
    if (m_logFilePath.isEmpty())
        return;

    QFileInfo info(m_logFilePath);
    QDir().mkpath(info.absolutePath());
    rotateFileIfNeeded(m_logFilePath);
    m_file.setFileName(m_logFilePath);
    const bool opened =
        m_file.open(QIODevice::WriteOnly | QIODevice::Append | QIODevice::Text);
    Q_UNUSED(opened);
}

void PerfLog::rotateFileIfNeeded(const QString &path)
{
    static constexpr qint64 maxBytes = 1024 * 1024;
    QFileInfo info(path);
    if (!info.exists() || info.size() < maxBytes)
        return;
    const QString rotated = path + QStringLiteral(".1");
    QFile::remove(rotated);
    QFile::rename(path, rotated);
}
