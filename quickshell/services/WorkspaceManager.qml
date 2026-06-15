pragma Singleton
pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import Dankestia.Services

Singleton {
    id: root

    readonly property bool hasCompositor: Compositor.hasCompositor

    property var workspaces: Compositor.workspaces
    property int activeWsId: Compositor.activeWorkspaceId
    
    // Calculate occupied property for workspaces backwards compatibility
    property var occupied: {
        let occ = {}
        for (let i = 0; i < Compositor.workspaces.length; i++) {
            let ws = Compositor.workspaces[i]
            occ[ws.id] = ws.windows > 0
        }
        return occ
    }

    property bool onSpecial: {
        for (let i = 0; i < Compositor.workspaces.length; i++) {
            let ws = Compositor.workspaces[i]
            if (ws.id === activeWsId && ws.isSpecial) return true
        }
        return false
    }

    property string activeWindowTitle: {
        let w = Compositor.focusedWindow
        return w ? w.title : ""
    }

    property string activeWindowAppId: {
        let w = Compositor.focusedWindow
        return w ? w.appId : ""
    }

    function dispatch(command) {
        Compositor.dispatch(command)
    }
}
