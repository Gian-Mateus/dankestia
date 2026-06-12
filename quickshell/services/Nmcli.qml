pragma Singleton
pragma ComponentBehavior: Bound

import QtQuick
import Quickshell
import Quickshell.Io
import Caelestia.Services

Singleton {
    id: root

    // Exposed State
    property bool isConnected: false
    property bool connecting: false
    property string activeInterface: ""
    property string activeConnection: ""
    property bool wifiEnabled: true
    property bool scanning: false

    property list<AccessPoint> networks: []
    property AccessPoint active: null
    property list<string> savedConnectionSsids: []
    property list<string> savedConnections: []

    // Ethernet
    property list<var> ethernetDevices: []
    property var activeEthernet: null
    property var ethernetDeviceDetails: null
    property var wirelessDeviceDetails: null

    signal connectionFailed(string ssid)

    Connections {
        target: DankestiaIPC
        function onNetworkStateUpdate(data) {
            if (!data) return;
            root.isConnected = data.wifiConnected || (data.ethernetStatus === "connected");
            root.connecting = data.isConnecting;
            root.wifiEnabled = data.wifiEnabled;
            root.activeConnection = data.wifiSSID;

            if (data.isScanning !== undefined) {
                root.scanning = data.isScanning;
            }

            if (data.wifiDevice) {
                root.activeInterface = data.wifiDevice;
            }
        }
    }

    Component.onCompleted: {
        getNetworks(() => {});
        DankestiaIPC.sendRequest("network.getState", {}, result => {
            root.wifiEnabled = result.wifiEnabled;
            root.isConnected = result.wifiConnected;
        });
    }

    function rescanWifi() {
        scanning = true;
        DankestiaIPC.sendRequest("network.wifi.scan", {}, result => {
            scanning = false;
            getNetworks(() => {});
        });
    }

    function toggleWifi(callback) {
        DankestiaIPC.sendRequest("network.wifi.toggle", {}, result => {
            if (callback) callback(result);
            if (result && result.enabled !== undefined) {
                root.wifiEnabled = result.enabled;
            }
        });
    }

    function enableWifi(enabled, callback) {
        let method = enabled ? "network.wifi.enable" : "network.wifi.disable";
        DankestiaIPC.sendRequest(method, {}, result => {
            if (callback) callback(result);
        });
    }

    function getWifiStatus(callback) {
        DankestiaIPC.sendRequest("network.getState", {}, result => {
            if (callback) callback(result.wifiEnabled);
        });
    }

    function disconnectFromNetwork() {
        DankestiaIPC.sendRequest("network.wifi.disconnect", {}, result => {
            getNetworks(() => {});
        });
    }

    function getNetworks(callback) {
        DankestiaIPC.sendRequest("network.wifi.networks", {}, result => {
            if (!result || !Array.isArray(result)) {
                if (callback) callback([]);
                return;
            }

            // Sync the arrays
            let currentMap = new Map();
            for (let n of result) {
                currentMap.set(`${n.frequency}:${n.ssid}:${n.bssid}`, n);
            }

            let toRemove = [];
            let rNetworks = root.networks;
            let existingMap = new Map();
            for (let i = 0; i < rNetworks.length; i++) {
                let rn = rNetworks[i];
                let key = `${rn.frequency}:${rn.ssid}:${rn.bssid}`;
                existingMap.set(key, rn);
                if (!currentMap.has(key)) {
                    toRemove.push(i);
                }
            }

            for (let i = toRemove.length - 1; i >= 0; i--) {
                let idx = toRemove[i];
                let ap = rNetworks[idx];
                rNetworks.splice(idx, 1);
                ap.destroy();
            }

            let savedSsids = [];
            for (let [key, network] of currentMap) {
                let match = existingMap.get(key);
                if (match) {
                    match.lastIpcObject = network;
                } else {
                    rNetworks.push(apComp.createObject(root, { lastIpcObject: network }));
                }

                if (network.saved) {
                    savedSsids.push(network.ssid);
                }
            }
            root.savedConnectionSsids = savedSsids;

            // Find active
            root.active = root.networks.find(n => n.active) ?? null;

            if (callback) callback(root.networks);
        });
    }

    function hasSavedProfile(ssid) {
        return root.savedConnectionSsids.includes(ssid);
    }

    function forgetNetwork(ssid, callback) {
        DankestiaIPC.sendRequest("network.wifi.forget", {ssid: ssid}, result => {
            if (callback) callback({success: !result.error, error: result.error});
            getNetworks(() => {});
        });
    }

    function connectToNetwork(ssid, password, bssid, callback) {
        DankestiaIPC.sendRequest("network.wifi.connect", {
            ssid: ssid,
            password: password,
            interactive: false
        }, result => {
            if (result && result.error && (result.error.includes("secrets") || result.error.includes("password"))) {
                if (callback) callback({success: false, needsPassword: true, error: result.error});
            } else {
                if (callback) callback({success: !result.error, error: result.error});
            }
            getNetworks(() => {});
        });
    }

    function connectToNetworkWithPasswordCheck(ssid, isSecure, callback, bssid) {
        // Just call connect, if it needs password the backend will fail with a secrets error since interactive: false.
        connectToNetwork(ssid, "", bssid, result => {
            if (result && result.error && (result.error.includes("secrets") || result.error.includes("password"))) {
                if (callback) callback({success: false, needsPassword: true, error: result.error});
            } else {
                if (callback) callback(result);
            }
        });
    }

    // Dummy / Unsupported methods that shouldn't crash
    function getDeviceStatus(callback) { if(callback) callback("") }
    function getWirelessInterfaces(callback) { if(callback) callback([]) }
    function getEthernetInterfaces(callback) { if(callback) callback([]) }
    function connectEthernet(connName, iface, cb) { if(cb) cb({success:false}) }
    function disconnectEthernet(connName, cb) { if(cb) cb({success:false}) }
    function isInterfaceConnected(iface, cb) { if(cb) cb(false) }
    function loadSavedConnections(cb) { getNetworks(cb) }
    function getDeviceDetails(iface, cb) { if(cb) cb("") }
    function refreshStatus(cb) { if(cb) cb({connected: root.isConnected, interface: root.activeInterface}) }

    Component {
        id: apComp
        AccessPoint {}
    }

    component AccessPoint: QtObject {
        required property var lastIpcObject
        readonly property string ssid: lastIpcObject.ssid || ""
        readonly property string bssid: lastIpcObject.bssid || ""
        readonly property int strength: lastIpcObject.signal || 0
        readonly property int frequency: lastIpcObject.frequency || 0
        readonly property bool active: lastIpcObject.connected || false
        readonly property string security: lastIpcObject.secured ? "secured" : ""
        readonly property bool isSecure: lastIpcObject.secured || false
    }

    // Unused Timers for compatibility
    Timer { id: connectionCheckTimer }
    Timer { id: immediateCheckTimer; property int checkCount: 0 }
    property var pendingConnection: null
}
