pragma Singleton

import QtQuick
import Dankestia.Services

Item {
    id: root

    property real percentage: 0
    property var primaryDisk: ({ name: "/", total: 1, used: 0, available: 1 })
    property list<var> disks: [primaryDisk]
    property string manualPrimaryDisk: ""

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data && data.storage) {
                const tot = data.storage.totalMB || 1;
                const usd = data.storage.usedMB || 0;
                const free = data.storage.freeMB || 0;
                
                root.percentage = (usd / tot) * 100.0;
                root.primaryDisk = {
                    name: "/",
                    total: tot * 1024 * 1024,
                    used: usd * 1024 * 1024,
                    available: free * 1024 * 1024
                };
                root.disks = [root.primaryDisk];
            }
        }
    }
}
