import QtQuick
import Quickshell
import Quickshell.Wayland

ShellWindow {
    Component.onCompleted: {
        console.log("Toplevels: " + ToplevelManager.toplevels.values.length);
        for (let t of ToplevelManager.toplevels.values) {
            console.log(t.title);
        }
        Quickshell.exit(0);
    }
}
