pragma Singleton
pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import Quickshell.Io
import Dankestia.Services
import qs.components.misc

Singleton {
    id: root

    readonly property list<Monitor> monitors: variants.instances // qmllint disable incompatible-type

    function getMonitorForScreen(screen: ShellScreen): var {
        return monitors.find(m => m.modelData === screen); // qmllint disable missing-property
    }

    function getMonitor(query: string): var {
        if (query === "active") {
            // Try Hyprland first, fallback to first monitor
            if (typeof Hypr !== "undefined" && Hypr.monitorFor) {
                var found = monitors.find(m => Hypr.monitorFor(m.modelData)?.focused);
                if (found) return found;
            }
            return monitors.length > 0 ? monitors[0] : null;
        }
        return monitors.find(m => m.modelData.name === query); // qmllint disable missing-property
    }

    function increaseBrightness(): void {
        DankestiaIPC.sendRequest("brightness.increment", { "device": "active", "step": 5 })
    }

    function decreaseBrightness(): void {
        DankestiaIPC.sendRequest("brightness.decrement", { "device": "active", "step": 5 })
    }

    Variants {
        id: variants
        model: Quickshell.screens
        Monitor {}
    }

    // qmllint disable unresolved-type
    CustomShortcut {
        // qmllint enable unresolved-type
        name: "brightnessUp"
        description: qsTr("Increase brightness")
        onPressed: root.increaseBrightness()
    }

    // qmllint disable unresolved-type
    CustomShortcut {
        // qmllint enable unresolved-type
        name: "brightnessDown"
        description: qsTr("Decrease brightness")
        onPressed: root.decreaseBrightness()
    }

    component Monitor: Item {
        id: monitor

        required property ShellScreen modelData
        property real brightness: 0
        property string deviceName: modelData.name

        Connections {
            target: DankestiaIPC
            function onBrightnessDeviceUpdate(device) {
                if (device.name === monitor.deviceName || device.name === "default") {
                    monitor.brightness = device.percentage / 100.0;
                }
            }
            function onBrightnessStateUpdate(data) {
                if (data.devices) {
                    for (let i = 0; i < data.devices.length; i++) {
                        let d = data.devices[i];
                        if (d.name === monitor.deviceName || d.name === "default") {
                            monitor.brightness = d.percentage / 100.0;
                            break;
                        }
                    }
                }
            }
        }

        function setBrightness(value: real): void {
            value = Math.max(0, Math.min(1, value));
            brightness = value;
            DankestiaIPC.sendRequest("brightness.setBrightness", { "device": deviceName, "percent": value * 100 });
        }
    }
}
