pragma Singleton
import QtQuick
import "../../services"

QtObject {
    property string name: "Unknown GPU"
    property real percentage: 0
    property real temperature: 0
    property string type: "NONE"

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data.gpu) {
                name = data.gpu.name || name;
                percentage = data.gpu.percentage || percentage;
                temperature = data.gpu.temperature || temperature;
                type = data.gpu.type || type;
            }
        }
    }
}
