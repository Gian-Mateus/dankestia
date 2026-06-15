pragma Singleton

import QtQuick
import Dankestia.Services

Item {
    id: root

    // Enum-like constants (QML requires lowercase property names)
    readonly property int typeNone: 0
    readonly property int typeIntegrated: 1
    readonly property int typeDiscrete: 2

    property int type: root.typeNone
    property string name: "GPU"
    property real percentage: 0
    property real temperature: 0

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data && data.gpu) {
                if (data.gpu.name !== undefined) root.name = data.gpu.name;
                if (data.gpu.percentage !== undefined) root.percentage = data.gpu.percentage;
                if (data.gpu.temperature !== undefined) root.temperature = data.gpu.temperature;
                if (data.gpu.type !== undefined) {
                    if (data.gpu.type === "discrete") root.type = root.typeDiscrete;
                    else if (data.gpu.type === "integrated") root.type = root.typeIntegrated;
                    else root.type = root.typeNone;
                }
            }
        }
    }
}
