// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include <QCommandLineParser>
#include <QDir>
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQuickWindow>
#include <QStandardPaths>

#include "catalogwriter.h"
#include "interfacemetrics.h"
#include "theme.h"
#include "uiservices.h"

// The gallery: every component, in every state, in all four themes.
//
// kvit-hub's GalleryView.qml is where this idea comes from, and its header
// says why it exists: it is where the rule that no distinction rests on colour
// alone gets checked once, so every view built from verified components
// inherits the check rather than repeating it.
//
// Here it is an application of its own rather than a view inside one, for
// three reasons. It is the documentation — one page per component, its states
// and a working snippet, following the arrangement HuskarUI uses. It is the
// screenshot run, which is what makes a token change reviewable as a diff
// against the previous set rather than as sixty images to look at again. And
// it is the harness the `kvit-preview` skill wraps, so an agent can render one
// component against sample data and look at what it built.
//
//   kvit-ui-gallery                       # browse
//   kvit-ui-gallery --shots <directory>   # write the screenshot set and exit
//   kvit-ui-gallery --page KvitRow --theme dark
//                                         # one component, one theme
//   kvit-ui-gallery --catalog <file>      # write the skill's catalogue
int main(int argc, char *argv[])
{
    QGuiApplication app(argc, argv);
    app.setApplicationName(QStringLiteral("kvit-ui-gallery"));
    app.setOrganizationName(QStringLiteral("kvit"));

    QCommandLineParser parser;
    parser.setApplicationDescription(
        QStringLiteral("Every kvit-ui component, in every state and theme."));
    parser.addHelpOption();
    const QCommandLineOption shots(
        QStringLiteral("shots"),
        QStringLiteral("Write one PNG per page per theme into <directory> and "
                       "exit."),
        QStringLiteral("directory"));
    const QCommandLineOption page(
        QStringLiteral("page"),
        QStringLiteral("Open one component's page (its type name)."),
        QStringLiteral("name"));
    const QCommandLineOption theme(
        QStringLiteral("theme"),
        QStringLiteral("light | dark | sepia | highContrast."),
        QStringLiteral("id"));
    const QCommandLineOption catalog(
        QStringLiteral("catalog"),
        QStringLiteral("Write the vocabulary skill's catalogue to <file> and "
                       "exit."),
        QStringLiteral("file"));
    const QCommandLineOption size(
        QStringLiteral("interface-size"),
        QStringLiteral("The interface size to render at (10 to 24)."),
        QStringLiteral("px"));
    parser.addOption(shots);
    parser.addOption(page);
    parser.addOption(theme);
    parser.addOption(catalog);
    parser.addOption(size);
    parser.process(app);

    // The gallery persists its own settings — which theme and size it was left
    // at — beside the applications rather than in their file, so browsing it
    // never moves a reader's real theme.
    const QString settingsPath =
        QStandardPaths::writableLocation(QStandardPaths::AppConfigLocation)
        + QStringLiteral("/gallery.json");
    QDir().mkpath(QFileInfo(settingsPath).absolutePath());
    KvitUi::DefaultServices::openSettings(settingsPath);

    if (parser.isSet(theme))
        KvitUi::DefaultServices::theme()->setThemeId(parser.value(theme));
    if (parser.isSet(size)) {
        KvitUi::DefaultServices::interfaceMetrics()->setFontSize(
            parser.value(size).toInt());
    } else if (parser.isSet(shots)) {
        // A screenshot run starts from the default interface size unless it
        // is told otherwise.
        //
        // The gallery remembers the size it was last left at, which is right
        // when browsing and wrong here: a set written after somebody looked
        // at one page at 20 px is a set that cannot be diffed against the
        // last one, and nothing about the images says why they all moved.
        KvitUi::DefaultServices::interfaceMetrics()->resetToDefaults();
    }

    QQmlApplicationEngine engine;

    // The catalogue is written without showing a window: it reads the same
    // list the gallery renders, and building a window to read a list would
    // make the generator need a display.
    if (parser.isSet(catalog))
        return KvitUi::writeCatalog(&engine, parser.value(catalog)) ? 0 : 1;

    // Initial properties rather than context properties. A context property
    // is a name a static analyser cannot see, so every use of one is reported
    // by qmllint as unqualified access, and the category that catches a real
    // typo stops being worth reading.
    engine.setInitialProperties({
        { QStringLiteral("shotDirectory"),
          parser.isSet(shots) ? parser.value(shots) : QString() },
        { QStringLiteral("startingPage"),
          parser.isSet(page) ? parser.value(page) : QString() },
    });

    engine.loadFromModule("Kvit.Gallery", "Gallery");
    if (engine.rootObjects().isEmpty())
        return 1;

    return app.exec();
}
