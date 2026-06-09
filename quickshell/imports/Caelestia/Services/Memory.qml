pragma Singleton
import QtQuick
import "../../services"

QtObject {
    property real total: 0
    property real used: 0
    property real available: 0

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data.memory) {
                total = data.memory.totalMB || total;
                used = data.memory.usedMB || used;
                available = data.memory.availableMB || available;
            }
        }
    }
}
