qml_module(dankestia-internal
    URI Dankestia.Internal
    SOURCES
        circularbuffer.cpp
        circularindicatormanager.cpp
        linearindicatormanager.cpp
        hyprdevices.cpp
        hyprextras.cpp
        logindmanager.cpp
        sparklineitem.cpp
        visualiserbars.cpp
    LIBRARIES
        Qt::Gui
        Qt::Quick
        Qt::Network
        Qt::DBus
)
