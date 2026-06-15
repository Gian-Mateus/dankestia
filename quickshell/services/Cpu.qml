pragma Singleton

import QtQuick
import Dankestia.Services

Item {
    id: root

    property string name: "CPU"
    property real percentage: 0
    property real temperature: 0

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data && data.cpu) {
                if (data.cpu.name !== undefined) root.name = data.cpu.name;
                if (data.cpu.percentage !== undefined) root.percentage = data.cpu.percentage;
                if (data.cpu.temperature !== undefined) root.temperature = data.cpu.temperature;
            }
        }
    }
}
