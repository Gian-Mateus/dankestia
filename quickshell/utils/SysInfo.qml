pragma Singleton

import QtQuick
import Quickshell
import qs.Services
import Quickshell.Io
import Caelestia.Config
import qs.utils

Singleton {
    id: root

    property string osName: "Linux"
    property string osPrettyName: "Linux"
    property string osId: "linux"
    property list<string> osIdLike: []
    property string osLogo: Qt.resolvedUrl(`${Quickshell.shellDir}/assets/logo.svg`)
    property bool isDefaultLogo: true

    property string uptime: "0 minutes"
    readonly property string user: Quickshell.env("USER")
    readonly property string wm: Quickshell.env("XDG_CURRENT_DESKTOP") || Quickshell.env("XDG_SESSION_DESKTOP")
    readonly property string shell: Quickshell.env("SHELL").split("/").pop()

    property string kernel: "Unknown"
    property string hostname: Quickshell.env("HOSTNAME") || "Dankestia"
    property string firmware: ""

    // DMI vendor/model, combined into a single human-readable device name
    property string boardVendor: ""
    property string boardName: ""
    readonly property string device: {
        if (!boardName)
            return boardVendor;
        if (!boardVendor || boardName.toLowerCase().startsWith(boardVendor.toLowerCase()))
            return boardName;
        return `${boardVendor} ${boardName}`;
    }

    // Strips the placeholder strings OEMs commonly leave in DMI fields
    function sanitiseDmi(s: string): string {
        const t = s.trim();
        const junk = ["to be filled by o.e.m.", "system product name", "system manufacturer", "system version", "default string", "o.e.m.", "not specified", "not applicable", "unknown", "none", ""];
        return junk.includes(t.toLowerCase()) ? "" : t;
    }

    Connections {
        target: GlobalConfig.general
        function onLogoChanged(): void {
            root.updateLogo()
        }
    }

    function updateLogo() {
        if (GlobalConfig.general.logo === "caelestia") {
            root.osLogo = Qt.resolvedUrl(`${Quickshell.shellDir}/assets/logo.svg`);
            root.isDefaultLogo = true;
        } else if (GlobalConfig.general.logo) {
            root.osLogo = Quickshell.iconPath(GlobalConfig.general.logo, true) || "file://" + Paths.absolutePath(GlobalConfig.general.logo);
            root.isDefaultLogo = false;
        }
    }

    Connections {
        target: DankestiaIPC
        function onSysinfoStateUpdate(data) {
            if (data.osName && data.osName !== root.osPrettyName) {
                root.osPrettyName = data.osName;
                root.osName = data.osName;
                root.updateLogo();
            }
            if (data.kernelVersion && data.kernelVersion !== root.kernel) {
                root.kernel = data.kernelVersion;
            }
            if (data.uptimeSeconds !== undefined) {
                const up = data.uptimeSeconds;
                const days = Math.floor(up / 86400);
                const hours = Math.floor((up % 86400) / 3600);
                const minutes = Math.floor((up % 3600) / 60);

                let str = "";
                if (days > 0)
                    str += `${days} day${days === 1 ? "" : "s"}`;
                if (hours > 0)
                    str += `${str ? ", " : ""}${hours} hour${hours === 1 ? "" : "s"}`;
                if (minutes > 0 || !str)
                    str += `${str ? ", " : ""}${minutes} minute${minutes === 1 ? "" : "s"}`;
                root.uptime = str;
            }
        }
    }

    // Fallback info via fileview just for DMI strings since they are one-shot
    FileView {
        path: "/sys/class/dmi/id/sys_vendor"
        printErrors: false
        onLoaded: root.boardVendor = root.sanitiseDmi(text())
    }

    FileView {
        path: "/sys/class/dmi/id/product_name"
        printErrors: false
        onLoaded: root.boardName = root.sanitiseDmi(text())
    }

    FileView {
        path: "/sys/class/dmi/id/bios_version"
        printErrors: false
        onLoaded: root.firmware = root.sanitiseDmi(text())
    }

    FileView {
        path: "/proc/sys/kernel/hostname"
        printErrors: false
        onLoaded: root.hostname = text().trim()
    }
}
