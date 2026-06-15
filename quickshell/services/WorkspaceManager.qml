pragma Singleton
pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import Quickshell.Io
import "../Common"

Singleton {
    id: root

    readonly property bool isHyprland: typeof Hypr !== "undefined"
    readonly property bool isNiri: Quickshell.env("XDG_SESSION_DESKTOP") === "niri" || Quickshell.env("NIRI_SOCKET") !== ""

    property var workspaces: []
    property int activeWsId: 1
    property bool onSpecial: false
    property var occupied: ({})
    property var niriIdToIdx: ({})
    property string activeWindowTitle: ""
    property string activeWindowAppId: ""
    property var niriWindows: ({})

    Component.onCompleted: {
        if (isHyprland && !isNiri) {
            hyprlandTimer.start();
        } else if (isNiri) {
            niriSocket.connected = true;
        }
    }

    // --- Hyprland Fallback ---
    Timer {
        id: hyprlandTimer
        interval: 100
        repeat: true
        onTriggered: {
            if (typeof Hypr !== "undefined" && Hypr.workspaces) {
                root.activeWsId = Hypr.activeWsId;
                root.onSpecial = Hypr.focusedMonitor?.lastIpcObject?.specialWorkspace?.name !== "";
                
                const occ = {};
                for (const ws of Hypr.workspaces.values) {
                    occ[ws.id] = ws.lastIpcObject.windows > 0;
                }
                root.occupied = occ;
                root.activeWindowTitle = Hypr.activeToplevel?.title ?? "";
                root.activeWindowAppId = Hypr.activeToplevel?.lastIpcObject?.class ?? "";
            }
        }
    }

    // --- Niri Socket ---
    DankSocket {
        id: niriSocket
        path: Quickshell.env("NIRI_SOCKET")
        connected: false

        onConnectionStateChanged: {
            if (connected) {
                send('"EventStream"');
            }
        }

        parser: SplitParser {
            onRead: line => {
                if (!line) return;
                try {
                    const event = JSON.parse(line);
                    root.handleNiriEvent(event);
                } catch (e) {}
            }
        }
    }

    function handleNiriEvent(event) {
        if (event.WorkspacesChanged) {
            const occ = {};
            const map = {};
            for (const ws of event.WorkspacesChanged.workspaces) {
                occ[ws.idx] = true;
                map[ws.id] = ws.idx;
            }
            root.occupied = occ;
            root.niriIdToIdx = map;
        } else if (event.WorkspaceActivated) {
            const idx = root.niriIdToIdx[event.WorkspaceActivated.id];
            if (idx !== undefined) {
                root.activeWsId = idx;
            }
        } else if (event.WindowOpened) {
            const w = event.WindowOpened.window;
            const wins = root.niriWindows;
            wins[w.id] = { title: w.title, appId: w.app_id };
            root.niriWindows = wins;
        } else if (event.WindowClosed) {
            const wins = root.niriWindows;
            delete wins[event.WindowClosed.id];
            root.niriWindows = wins;
        } else if (event.WindowFocusChanged) {
            const id = event.WindowFocusChanged.id;
            if (id !== null && root.niriWindows[id]) {
                root.activeWindowTitle = root.niriWindows[id].title || "";
                root.activeWindowAppId = root.niriWindows[id].appId || "";
            } else {
                root.activeWindowTitle = "";
                root.activeWindowAppId = "";
            }
        }
    }

    function dispatch(command) {
        if (isNiri) {
            if (command.startsWith("workspace ")) {
                const wsIdx = parseInt(command.replace("workspace ", ""));
                if (!isNaN(wsIdx)) {
                    Quickshell.execute("niri msg action focus-workspace " + wsIdx);
                }
            } else if (command === "togglespecialworkspace special") {
                Quickshell.execute("niri msg action toggle-overview");
            }
        } else if (isHyprland) {
            Hypr.dispatch(command);
        }
    }
}
