pragma Singleton

import QtQuick
import Quickshell

pragma ComponentBehavior: Bound

QtObject {
    id: root

    property var monitors: []
    property var workspaces: []
    property var windows: []
    property int activeWorkspaceId: 0
    property bool hasCompositor: DankestiaIPC.capabilities.indexOf("compositor") !== -1

    readonly property var focusedWindow: {
        for (let i = 0; i < windows.length; i++) {
            if (windows[i].focused) return windows[i]
        }
        return null
    }

    Connections {
        target: DankestiaIPC
        function onEventReceived(service, data) {
            if (service === "compositor") {
                root.monitors = data.monitors || []
                root.workspaces = data.workspaces || []
                root.windows = data.windows || []
                root.activeWorkspaceId = data.activeWorkspaceId || 0
            }
        }
    }

    function dispatch(command) {
        DankestiaIPC.sendRequest("compositor.dispatch", { command: command })
    }

    function getWorkspace(id) {
        for (let i = 0; i < root.workspaces.length; i++) {
            if (root.workspaces[i].id === id) {
                return root.workspaces[i]
            }
        }
        return null
    }

    function getMonitor(id) {
        for (let i = 0; i < root.monitors.length; i++) {
            if (root.monitors[i].id === id) {
                return root.monitors[i]
            }
        }
        return null
    }

    function getFocusedWindow() {
        return root.focusedWindow
    }
}
