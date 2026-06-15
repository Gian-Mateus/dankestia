pragma Singleton

import QtQuick
import Dankestia.Services

Item {
    id: root

    property real percentage: 0
    property real used: 0
    property real total: 1

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data && data.memory) {
                root.total = data.memory.totalMB || 1;
                root.used = data.memory.usedMB || 0;
                root.percentage = (root.used / root.total) * 100.0;
            }
        }
    }
}
