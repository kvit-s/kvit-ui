# The stylesheet generator.
#
# A tool rather than a script, because a script that reads theme.cpp would be a
# second implementation of the token table — which is exactly the arrangement
# that let kvit-hub's copy drift four values out of five. This links the real
# table and asks it.
qt_add_executable(tokens-to-css tokens-to-css.cpp)
target_link_libraries(tokens-to-css PRIVATE kvit-ui-tokens Qt6::Core Qt6::Gui)

# `cmake --build <dir> --target tokens-css` regenerates the committed file.
# The matching check runs in ctest and in CI.
add_custom_target(tokens-css
    COMMAND tokens-to-css ${CMAKE_SOURCE_DIR}/ux/tokens.css
    DEPENDS tokens-to-css
    COMMENT "Writing ux/tokens.css from the token table"
    VERBATIM)
