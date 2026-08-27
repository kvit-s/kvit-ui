find_package(Qt6 REQUIRED COMPONENTS Test QuickTest)

# One executable per test class, which looks wasteful and is not: QTest runs
# one QObject per binary, and a suite squeezed into fewer binaries has to
# either give up QTEST_MAIN or run unrelated cases in one process, where a
# crash in one takes the rest of the report with it.
function(kvit_ui_test name)
    qt_add_executable(${name} ${name}.cpp)
    target_link_libraries(${name} PRIVATE kvit-ui-tokens Qt6::Test ${ARGN})
    add_test(NAME ${name} COMMAND ${name})
endfunction()

# The three ported tests. test_theme and test_interfacemetrics are kvit-notes'
# files; test_typography is its file with the three BlockEditorEngine cases
# removed, because that class is the editor's document renderer and stays
# there.
#
# Qt6::Gui on the theme test: Theme resolves "system" through
# QGuiApplication::styleHints(), so without it QTEST_MAIN builds a
# QCoreApplication and the system-theme cases assert against a null pointer.
kvit_ui_test(test_theme Qt6::Gui)
set_tests_properties(test_theme PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")
kvit_ui_test(test_typography Qt6::Gui)
set_tests_properties(test_typography PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")
kvit_ui_test(test_interfacemetrics)

# The density and spacing tokens, and the one claim the whole arrangement
# rests on: one setting moves every type role and every geometry value, and
# leaves the document type scale alone.
kvit_ui_test(test_density)

# The chart ramps, checked rather than asserted. This is what makes an edit to
# a ramp fail the build instead of the eye: every step of all three ramps, in
# all four themes, against that theme's surfaces and its reserved hues.
kvit_ui_test(test_palette Qt6::Gui)
set_tests_properties(test_palette PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")

# The icon catalogue: that every meaning name resolves to a glyph, that an
# unknown name resolves to nothing rather than to a placeholder, and that the
# nine characters kvit-notes draws as literals all have somewhere to go.
qt_add_executable(test_icons test_icons.cpp)
target_link_libraries(test_icons PRIVATE kvit-ui-qml Qt6::Test Qt6::Gui)
add_test(NAME test_icons COMMAND test_icons)
set_tests_properties(test_icons PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")

# KvitTable's benchmark model, and the 250,000-row measurement prd.md
# Decision 3 turns on. The model half runs everywhere; the view half needs a
# window, so it is a Qt Quick Test and is labelled so a headless CI run can
# see which one it skipped.
qt_add_executable(test_tablemodel test_tablemodel.cpp)
target_link_libraries(test_tablemodel PRIVATE kvit-ui-qml Qt6::Test Qt6::Gui)
add_test(NAME test_tablemodel COMMAND test_tablemodel)
set_tests_properties(test_tablemodel PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")

# Every component in the module loads, in all four themes, with no QML
# warning. This is the check that a component referencing a token that does not
# exist, or an icon name with a typo in it, fails the build rather than drawing
# an empty box.
qt_add_executable(test_components test_components.cpp)
target_link_libraries(test_components PRIVATE kvit-ui-qml Qt6::Test Qt6::Quick Qt6::Gui)
add_test(NAME test_components COMMAND test_components)
set_tests_properties(test_components PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")

# The gallery catalogue: that every component has a page, and that every code
# sample in it compiles and runs. The samples are strings rather than QML
# files, so without this a broken one reaches a reader who copies it.
if(KVIT_UI_BUILD_GALLERY)
    qt_add_executable(test_gallery test_gallery.cpp)
    # kvit-ui-qml is named as well as reached through the gallery module.
    # With static libraries the linker keeps only the object files something
    # references, and a QML module's type registration is a file nothing
    # references by name — so `import Kvit.Ui` inside a specimen fails with
    # "module is not installed" unless the module is on this target's own link
    # line.
    target_link_libraries(test_gallery
        PRIVATE kvit-ui-gallery-qml kvit-ui-qml Qt6::Test Qt6::Quick Qt6::Gui)
    add_test(NAME test_gallery COMMAND test_gallery)
    set_tests_properties(test_gallery
        PROPERTIES ENVIRONMENT "QT_QPA_PLATFORM=offscreen")
endif()

# The QML lint gate, as a test.
#
# It catches what nothing else can: a binding to a property that does not
# exist. QML does not throw on one — it evaluates to `undefined`, which a
# colour property renders as transparent and a numeric one as zero, so a
# component with a typo in a token name renders as nothing and the build says
# nothing about it.
find_program(KVIT_UI_BASH bash)
if(KVIT_UI_BASH)
    add_test(NAME qmllint
        COMMAND ${KVIT_UI_BASH} ${CMAKE_SOURCE_DIR}/tools/run-qmllint.sh)
    set_tests_properties(qmllint PROPERTIES
        ENVIRONMENT "KVIT_UI_BUILD_DIR=${CMAKE_BINARY_DIR}")
endif()

# The catalogue the vocabulary skill publishes, which is generated from the
# same list the gallery renders. A stale one tells an agent a component exists
# that does not, or hides one that does.
if(KVIT_UI_BUILD_GALLERY AND KVIT_UI_BASH)
    add_test(NAME generator_skill_catalog
        COMMAND ${KVIT_UI_BASH} ${CMAKE_SOURCE_DIR}/tools/check-catalog.sh
                $<TARGET_FILE:kvit-ui-gallery>
                ${CMAKE_SOURCE_DIR}/agent/skills/kvit-ui/catalog.md)
    set_tests_properties(generator_skill_catalog PROPERTIES
        ENVIRONMENT "QT_QPA_PLATFORM=offscreen")
endif()

# The two generated files, checked so a stale one fails the build in the same
# place everything else does. CI runs these too; having them here as well means
# a developer who never opens CI still finds out.
#
# ux/tokens.css is what a drawing is judged in. When it and the token table
# disagree, a screen is approved in colours the application does not draw with,
# and the drawing workflow quietly stops meaning anything.
add_test(NAME generator_tokens_css
    COMMAND tokens-to-css --check ${CMAKE_SOURCE_DIR}/ux/tokens.css)

# src/qml/iconcatalog.cpp is what every named symbol resolves through.
find_package(Python3 COMPONENTS Interpreter)
if(Python3_Interpreter_FOUND)
    add_test(NAME generator_icon_catalog
        COMMAND ${Python3_EXECUTABLE}
                ${CMAKE_SOURCE_DIR}/tools/generate-icon-catalog.py --check)
else()
    message(STATUS "No Python 3 found: the icon catalogue check will not run")
endif()
