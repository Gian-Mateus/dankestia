pragma Singleton
import QtQuick
import "../../services"

QtObject {
    property string name: "CPU"
    property real percentage: 0
    property real temperature: 0

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data.cpu) {
                name = data.cpu.name || name;
                percentage = data.cpu.percentage || percentage;
                temperature = data.cpu.temperature || temperature;
            }
        }
    }
}
