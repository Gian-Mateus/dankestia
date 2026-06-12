qml_module(dankestia-images
    URI Dankestia.Images
    SOURCES
        cachingimageprovider.cpp
        imagecacher.cpp
        iutils.cpp
    LIBRARIES
        Qt::Gui
        Qt::Quick
        Qt::Concurrent
)
