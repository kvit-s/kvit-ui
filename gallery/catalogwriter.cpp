// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "catalogwriter.h"

#include <QFile>
#include <QQmlComponent>
#include <QQmlContext>
#include <QQmlEngine>
#include <QRegularExpression>
#include <QScopedPointer>
#include <QTextStream>
#include <QVariantList>
#include <QVariantMap>

#include <cstdio>

#include "iconcatalog.h"

namespace {

struct Property
{
    QString name;
    QString type;
    bool required = false;
    // The line comment directly above the declaration, which in this tree is
    // where the reason for the property is written.
    QString note;
};

// Read a component's declared properties out of its QML source.
//
// The source is in the binary's own resources, which is what makes this exact
// rather than a guess: it is the file that was compiled, not a description of
// it. A regular expression is enough because the shape being matched is one
// line of QML with a fixed grammar, and anything it fails to match simply does
// not appear rather than appearing wrong.
//
// Only declarations on the root object count, which here means exactly four
// spaces of indentation. Without that rule the catalogue advertises every
// `required property var modelData` inside a Repeater delegate as part of the
// component's interface, and an agent reading it writes
// `KvitSpark { modelData: ... }`. Fifteen of the sixty-eight components have a
// delegate, so this is most of the ones that draw a list.
QList<Property> propertiesOf(const QString &component)
{
    QFile file(QStringLiteral(":/qt/qml/Kvit/Ui/%1.qml").arg(component));
    if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
        return {};
    const QStringList lines =
        QString::fromUtf8(file.readAll()).split(QLatin1Char('\n'));

    static const QRegularExpression declaration(
        QStringLiteral("^    (?:(readonly)\\s+)?(?:(required)\\s+)?property\\s+"
                       "([A-Za-z_][\\w.<>]*)\\s+([A-Za-z_]\\w*)"));

    QList<Property> found;
    for (int i = 0; i < lines.size(); ++i) {
        const QRegularExpressionMatch match = declaration.match(lines.at(i));
        if (!match.hasMatch())
            continue;

        Property property;
        property.required = !match.captured(2).isEmpty();
        property.type = match.captured(3);
        property.name = match.captured(4);
        if (!match.captured(1).isEmpty())
            property.type = QStringLiteral("readonly ") + property.type;

        // The comment block immediately above, up to a blank line. Only the
        // first sentence: a catalogue is a reference, and the reasoning
        // belongs in the source where somebody changing it will read it.
        QStringList above;
        for (int j = i - 1; j >= 0; --j) {
            const QString line = lines.at(j).trimmed();
            if (!line.startsWith(QLatin1String("//")))
                break;
            above.prepend(line.mid(2).trimmed());
        }
        if (!above.isEmpty()) {
            const QString joined = above.join(QLatin1Char(' '));
            const int stop = joined.indexOf(QLatin1String(". "));
            property.note = stop > 0 ? joined.left(stop + 1) : joined;
        }
        found.append(property);
    }
    return found;
}

}   // namespace

namespace KvitUi {

bool writeCatalog(QQmlEngine *engine, const QString &path)
{
    QQmlComponent reader(engine);
    reader.setData("import QtQml\nimport Kvit.Gallery\n"
                   "QtObject { property var all: Catalog.components }",
                   QUrl(QStringLiteral("qrc:/kvit-ui-gallery/catalog-reader.qml")));
    if (!reader.isReady()) {
        std::fprintf(stderr, "catalog: %s\n", qPrintable(reader.errorString()));
        return false;
    }
    QScopedPointer<QObject> holder(reader.create());
    if (holder.isNull()) {
        std::fprintf(stderr, "catalog: %s\n", qPrintable(reader.errorString()));
        return false;
    }
    const QVariantList entries = holder->property("all").toList();

    QString text;
    QTextStream out(&text);

    out << "# The kvit-ui component catalogue\n\n"
           "GENERATED FILE — do not edit. Written by `kvit-ui-gallery "
           "--catalog`\nfrom the same list the gallery renders and the test "
           "suite compiles, so\nwhat is described here is what the repository "
           "contains.\n\n"
           "Every component is in the `Kvit.Ui` module: one `import Kvit.Ui` "
           "reaches\nall of them, the token singletons `Theme`, `Interface` "
           "and `Typography`,\nand the icon font.\n\n"
           "Each entry gives what the component is for, the properties it "
           "declares, and\na working sample. The samples are compiled by "
           "`tests/test_gallery`, so one\nthat does not work stops the "
           "build.\n\n";

    // A contents list, grouped, because a catalogue of sixty-eight components
    // is otherwise something an agent reads linearly to find one name.
    QString lastGroup;
    out << "## What there is\n\n";
    for (const QVariant &entry : entries) {
        const QVariantMap map = entry.toMap();
        const QString group = map.value(QStringLiteral("group")).toString();
        if (group != lastGroup) {
            out << "\n**" << group << "** — ";
            lastGroup = group;
        } else {
            out << ", ";
        }
        out << "`" << map.value(QStringLiteral("name")).toString() << "`";
    }
    out << "\n\n";

    // The icon names, which are the other half of what a call site needs: a
    // component takes a symbol by meaning, and a name that is not in this list
    // draws a marked placeholder and fails the build.
    out << "## Symbols\n\n"
           "`KvitIcon` and every component that takes a `symbol` accept these "
           "names.\nThey say what a symbol means rather than what it looks "
           "like, so the drawing\ncan change without touching a call site. A "
           "name outside this list draws a\nmarked placeholder and fails "
           "`tests/test_components`.\n\n";
    const QStringList symbols = IconCatalog::meaningNames();
    for (int i = 0; i < symbols.size(); ++i) {
        out << "`" << symbols.at(i) << "`";
        out << (i == symbols.size() - 1 ? "\n\n" : ", ");
    }

    out << "## The components\n\n";
    for (const QVariant &entry : entries) {
        const QVariantMap map = entry.toMap();
        const QString name = map.value(QStringLiteral("name")).toString();

        out << "### " << name << "\n\n";
        out << map.value(QStringLiteral("summary")).toString() << "\n\n";

        const QList<Property> properties = propertiesOf(name);
        if (!properties.isEmpty()) {
            out << "| Property | Type | |\n|---|---|---|\n";
            for (const Property &property : properties) {
                // A pipe in a note would end the table cell early. These
                // notes are often a list of the values a string may take,
                // written with pipes between them, so it is the common case
                // rather than an edge one.
                QString note = property.note;
                note.replace(QLatin1Char('|'), QLatin1String("\\|"));
                out << "| `" << property.name << "` | " << property.type
                    << " | " << (property.required ? "**required.** " : "")
                    << note << " |\n";
            }
            out << "\n";
        }

        const QVariantList specimens =
            map.value(QStringLiteral("specimens")).toList();
        for (const QVariant &specimen : specimens) {
            const QVariantMap one = specimen.toMap();
            out << "*" << one.value(QStringLiteral("caption")).toString()
                << "*\n\n```qml\n"
                << one.value(QStringLiteral("snippet")).toString()
                << "\n```\n\n";
        }
    }

    out.flush();

    QFile file(path);
    if (!file.open(QIODevice::WriteOnly | QIODevice::Truncate | QIODevice::Text)) {
        std::fprintf(stderr, "catalog: cannot write %s\n", qPrintable(path));
        return false;
    }
    file.write(text.toUtf8());
    std::printf("wrote %s: %d components\n", qPrintable(path),
                int(entries.size()));
    return true;
}

}   // namespace KvitUi
