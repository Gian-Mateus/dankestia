#pragma once

#include "configobject.hpp"

namespace dankestia::config {

class NexusConfig : public ConfigObject {
    Q_OBJECT
    QML_ANONYMOUS

    CONFIG_PROPERTY(int, wallpapersPerRow, 4)
    CONFIG_GLOBAL_PROPERTY(int, networkRescanInterval, 15000)

public:
    explicit NexusConfig(QObject* parent = nullptr)
        : ConfigObject(parent) {}
};

} // namespace dankestia::config
