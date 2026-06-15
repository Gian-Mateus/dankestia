pragma Singleton

import Quickshell
import qs.components
import qs.services

Singleton {
    property var screens: new Map()
    property var bars: new Map()

    function load(screen: ShellScreen, visibilities: DrawerVisibilities): void {
        var key = screen.name;
        if (typeof Hypr !== "undefined" && Hypr.monitorFor) {
            var m = Hypr.monitorFor(screen);
            if (m) key = m;
        }
        screens.set(key, visibilities);
    }

    function getForActive(): DrawerVisibilities {
        var key = null;
        if (typeof Hypr !== "undefined" && Hypr.focusedMonitor) {
            key = Hypr.focusedMonitor;
        }
        if (!key && Quickshell.screens && Quickshell.screens.length > 0) {
            key = Quickshell.screens[0].name;
        }
        
        // Se ainda não achar, pega o primeiro mapa
        var val = screens.get(key);
        if (!val) {
            const iterator = screens.values();
            const first = iterator.next();
            if (!first.done) val = first.value;
        }
        return val;
    }
}
