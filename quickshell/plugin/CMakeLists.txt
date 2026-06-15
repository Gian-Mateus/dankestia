cmake_minimum_required(VERSION 3.16)
project(DankestiaPlugins)

find_package(Qt6 REQUIRED COMPONENTS ShaderTools Core Qml Gui Quick QuickControls2 Concurrent Sql Network DBus)
find_package(PkgConfig REQUIRED)
pkg_check_modules(Qalculate IMPORTED_TARGET libqalculate REQUIRED)
pkg_check_modules(Pipewire IMPORTED_TARGET libpipewire-0.3 REQUIRED)
pkg_check_modules(Aubio IMPORTED_TARGET aubio REQUIRED)
pkg_check_modules(Cava IMPORTED_TARGET libcava QUIET)
if(NOT Cava_FOUND)
    pkg_check_modules(Cava IMPORTED_TARGET cava QUIET)
endif()
include(cmake/sensorslib.cmake)

set(QT_QML_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/qml")
qt_standard_project_setup(REQUIRES 6.9)

include(cmake/pch.cmake)
include(cmake/qml-module.cmake)

add_subdirectory(src/Dankestia)

message(STATUS "Fetching M3Shapes module")
include(FetchContent)
FetchContent_Declare(
    m3shapes_external
    GIT_REPOSITORY https://github.com/soramanew/m3shapes.git
    GIT_TAG        main
)
set(FETCHCONTENT_UPDATES_DISCONNECTED ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(m3shapes_external)
message(STATUS "Done fetching M3Shapes module")

# Fix m3shapes wrong rpath
if(TARGET m3shapesplugin)
    set_target_properties(m3shapesplugin PROPERTIES INSTALL_RPATH "$ORIGIN")
endif()
