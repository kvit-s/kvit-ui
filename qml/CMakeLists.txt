# The `Kvit.Ui` QML module: the token singletons and the component vocabulary.
#
# One module rather than two, because a component and the token it draws with
# are the same import for a consumer. `import Kvit.Ui` is the whole surface.

# The component QML, in the order prd.md §5.4 builds it: structure first,
# because everything else is placed inside it, then content, then marks, then
# quantities, then controls, then feedback. Each group is a list of its own so
# a reader can see the vocabulary from the build file.
set(KVIT_UI_QML_FOUNDATION
    KvitLabel.qml
    KvitIcon.qml
    KvitIconButton.qml
    KvitLink.qml
)

set(KVIT_UI_QML_STRUCTURE
    KvitWindow.qml
    KvitSidebar.qml
    KvitSidebarItem.qml
    KvitHeader.qml
    KvitBreadcrumb.qml
    KvitRegion.qml
    KvitViewHead.qml
    KvitStatusBar.qml
)

set(KVIT_UI_QML_CONTENT
    KvitSectionHeading.qml
    KvitRow.qml
    KvitSlimRow.qml
    KvitCard.qml
    KvitPanel.qml
    KvitPane.qml
    KvitDivider.qml
    KvitDisclosure.qml
    KvitEmptyState.qml
)

set(KVIT_UI_QML_MARKS
    KvitChip.qml
    KvitTag.qml
    KvitBadge.qml
    KvitSlug.qml
    KvitPip.qml
    KvitDot.qml
    KvitSignal.qml
)

set(KVIT_UI_QML_QUANTITIES
    KvitFigure.qml
    KvitBeforeAfter.qml
)

set(KVIT_UI_QML_CONTROLS
    KvitButton.qml
    KvitChipButton.qml
    KvitStepper.qml
    KvitField.qml
    KvitSearchField.qml
    KvitCheck.qml
    KvitSelect.qml
    KvitTab.qml
)

set(KVIT_UI_QML_FEEDBACK
    KvitToast.qml
    KvitPopover.qml
    KvitHint.qml
    KvitHoverCard.qml
    KvitDialog.qml
    KvitTooltip.qml
    KvitNotice.qml
)

# The data-display vocabulary (prd.md §5.5): the marks, the tiles and the
# table. These are the components each application had independently invented,
# and the nine common rules travel with them rather than being restated per
# caller.
set(KVIT_UI_QML_DATA
    KvitBar.qml
    KvitStackedBar.qml
    KvitSpark.qml
    KvitTrend.qml
    KvitDistribution.qml
    KvitGauge.qml
    KvitDelta.qml
    KvitStatTile.qml
    KvitFigureBlock.qml
    KvitTable.qml
    KvitCell.qml
)

# The components a written flow already names (prd.md §5.6). Group A is what
# the estate has built by hand between three and thirty-two times; Group B is
# absent everywhere and named in a flow that exists.
set(KVIT_UI_QML_FLOW
    KvitScrollBar.qml
    KvitMenu.qml
    KvitMenuItem.qml
    KvitTree.qml
    KvitSwitch.qml
    KvitRadioGroup.qml
    KvitProgress.qml
    KvitSlider.qml
    KvitSplitView.qml
    KvitSegmented.qml
    KvitTypeAhead.qml
    KvitConfirmInPlace.qml
    KvitTimeline.qml
    KvitNumberField.qml
    KvitMoneyField.qml
    KvitDualList.qml
    KvitSpotlight.qml
)

# Named before the module so the aliases below can be set on the same paths.
set(KVIT_UI_FONTS
    ../resources/fonts/Phosphor.ttf
    ../resources/fonts/Phosphor-LICENSE.txt
)

# The aliases, so the font lands at `fonts/Phosphor.ttf` inside the module
# rather than at `../resources/fonts/Phosphor.ttf`. KvitIcon then names it
# relative to its own URL, which resolves the same whether the module was
# loaded from a resource or from a directory on disk.
set_source_files_properties(../resources/fonts/Phosphor.ttf
    PROPERTIES QT_RESOURCE_ALIAS fonts/Phosphor.ttf)
set_source_files_properties(../resources/fonts/Phosphor-LICENSE.txt
    PROPERTIES QT_RESOURCE_ALIAS fonts/Phosphor-LICENSE.txt)

qt_add_qml_module(kvit-ui-qml
    ${KVIT_UI_LIB_TYPE}
    URI Kvit.Ui
    VERSION 1.0
    # Where the generated qmldir and .qmltypes go.
    #
    # Under one root shared with the gallery module, laid out as the URI —
    # <root>/Kvit/Ui — because that is the shape an import path has to have.
    # Left to itself Qt writes the qmldir straight into this directory, which
    # the engine finds through the module's own resources but qmllint cannot,
    # so the lint gate reports every component as "not found".
    OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/qml-modules/Kvit/Ui
    # The types are registered from C++ through the QML_FOREIGN wrappers in
    # ../src/qml/uisingletons.h, so there is no plugin to load: a consumer
    # links this target and the module is there.
    NO_PLUGIN
    DEPENDENCIES QtQuick
    SOURCES
        ../src/qml/uiservices.h
        ../src/qml/uiservices.cpp
        ../src/qml/uisingletons.h
        ../src/qml/iconcatalog.h
        ../src/qml/iconcatalog.cpp
        ../src/qml/tablemodelbase.h
        ../src/qml/tablemodelbase.cpp
        ../src/qml/benchmarktablemodel.h
        ../src/qml/benchmarktablemodel.cpp
        ../src/qml/treenodemodel.h
        ../src/qml/treenodemodel.cpp
    QML_FILES
        ${KVIT_UI_QML_FOUNDATION}
        ${KVIT_UI_QML_STRUCTURE}
        ${KVIT_UI_QML_CONTENT}
        ${KVIT_UI_QML_MARKS}
        ${KVIT_UI_QML_QUANTITIES}
        ${KVIT_UI_QML_CONTROLS}
        ${KVIT_UI_QML_FEEDBACK}
        ${KVIT_UI_QML_DATA}
        ${KVIT_UI_QML_FLOW}
    RESOURCES
        # The icon font ships with the module rather than with an application.
        # kvit-notes-pro adds a `qrc:/fonts/Phosphor.ttf` alias to its own
        # .qrc today, which means a second consumer of KvitIcon would have to
        # know to do the same and would draw nothing if it did not.
        ${KVIT_UI_FONTS}
)

target_include_directories(kvit-ui-qml PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/../src/qml)
target_link_libraries(kvit-ui-qml PUBLIC kvit-ui-tokens Qt6::Quick Qt6::QuickControls2)

if(KVIT_UI_SHARED_LIBS)
    set_target_properties(kvit-ui-qml PROPERTIES WINDOWS_EXPORT_ALL_SYMBOLS ON)
endif()

# The module's type description, always current.
#
# qmllint reads .qmltypes to resolve `Kvit.Ui`, and qt_add_qml_module only
# generates it as a by-product of building something that needs it. kvit-notes
# learned what that costs: its lint job named a generated target that had since
# moved, so the blocking check failed before running a single check. Naming the
# generated target in ALL means the description exists after any build.
if(TARGET kvit-ui-qml_qmltyperegistration)
    add_custom_target(kvit-ui-qmltypes ALL)
    # The backing module owns both AUTOMOC metadata production and type
    # registration. Depending on the registrar helper alone can run it before
    # its JSON inputs exist under Visual Studio.
    add_dependencies(kvit-ui-qmltypes kvit-ui-qml)
endif()
