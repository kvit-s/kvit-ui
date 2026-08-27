# The gallery: every component, in every state, in all four themes.
#
# An application of its own rather than a view inside one, because it is three
# things at once — the documentation, the screenshot run that makes a token
# change reviewable as a diff, and the harness the `kvit-preview` skill wraps
# so an agent can look at what it built.
#
# The QML is a library rather than living in the executable, so that
# tests/test_gallery.cpp can link it. That test is what holds the catalogue
# honest: it asserts every component in the module has a page, and that every
# code sample in the catalogue compiles and runs. Samples are strings rather
# than QML files, so nothing else in the build would notice a broken one.
add_library(kvit-ui-gallery-qml ${KVIT_UI_LIB_TYPE})

# Catalog.qml is `pragma Singleton`, and qt_add_qml_module only writes the
# `singleton` line into the qmldir when the source file is marked as one. This
# has to be set before the module is declared: without it the type resolves to
# a component, `Catalog.components` is undefined, and every page renders empty
# with no error saying why.
set_source_files_properties(Catalog.qml PROPERTIES QT_QML_SINGLETON_TYPE TRUE)

qt_add_qml_module(kvit-ui-gallery-qml
    URI Kvit.Gallery
    VERSION 1.0
    # Beside the Kvit.Ui module under one import root, so one -I reaches both.
    OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/qml-modules/Kvit/Gallery
    NO_PLUGIN
    DEPENDENCIES QtQuick
    SOURCES
        catalogwriter.h
        catalogwriter.cpp
    QML_FILES
        Gallery.qml
        GalleryPage.qml
        Specimen.qml
        Catalog.qml
)

target_include_directories(kvit-ui-gallery-qml PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(kvit-ui-gallery-qml
    PUBLIC kvit-ui-qml kvit-ui-tokens Qt6::Quick Qt6::Gui)

if(KVIT_UI_SHARED_LIBS)
    set_target_properties(kvit-ui-gallery-qml PROPERTIES
        WINDOWS_EXPORT_ALL_SYMBOLS ON)
endif()

qt_add_executable(kvit-ui-gallery main.cpp)
target_link_libraries(kvit-ui-gallery PRIVATE kvit-ui-gallery-qml)
