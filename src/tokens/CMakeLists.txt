# The token layer: every named value the four applications draw with, as
# plain C++ QObjects with no QML machinery on them.
#
# Six classes rather than the three prd.md §2 names, because the three do not
# stand alone. Theme, Typography and InterfaceMetrics each persist through
# SettingsStore and each time their own reflow through PerfLog, and Theme
# follows the desktop's light/dark and high-contrast settings through
# SystemAppearance. Stripping those out would have meant rewriting the three
# rather than copying them, and the comments in these headers are the design
# rationale worth keeping intact.
#
# Nothing here declares a QML type. The registrations are one file in
# src/qml/, which is the layering kvit-notes arrived at: the list of what QML
# can reach is one place rather than a macro in every header.
# qt_add_library rather than add_library, for one reason: it makes Qt generate
# this target's metatypes. Without them, qmltyperegistrar writes the QML_FOREIGN
# wrappers in src/qml/uisingletons.h as bare components with no properties, so
# qmllint reports every use of `Theme` or `Interface` as unqualified access and
# the two most valuable lint categories become noise.
qt_add_library(kvit-ui-tokens ${KVIT_UI_LIB_TYPE}
    theme.h theme.cpp
    typography.h typography.cpp
    interfacemetrics.h interfacemetrics.cpp
    palette.h palette.cpp
    settingsstore.h settingsstore.cpp
    systemappearance.h systemappearance.cpp
    perflog.h perflog.cpp
)

target_include_directories(kvit-ui-tokens PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

# Gui rather than Core alone: Theme is a table of QColor, and it reads the
# desktop's colour scheme through QGuiApplication::styleHints().
#
# Qml is PRIVATE and is one function's worth of dependency. SettingsStore takes
# values from QML, and QML hands a JavaScript array or object across as a
# QVariant wrapping a QJSValue; QJsonValue::fromVariant turns one of those into
# null without complaining, so the store unwraps it at the boundary. That
# unwrapping is the only thing in this library that needs QtQml, it lives
# entirely in settingsstore.cpp, and PRIVATE keeps it off a consumer's link
# line. Qt6::Quick is not linked at all.
target_link_libraries(kvit-ui-tokens PUBLIC Qt6::Core Qt6::Gui)
target_link_libraries(kvit-ui-tokens PRIVATE Qt6::Qml)

if(KVIT_UI_SHARED_LIBS)
    # A shared library on Windows exports nothing by default, and the whole
    # point of building shared here is that the same source tree links both
    # ways without a per-symbol export macro.
    set_target_properties(kvit-ui-tokens PROPERTIES
        WINDOWS_EXPORT_ALL_SYMBOLS ON)
endif()
