# Requirements

## v1 Requirements

### Architecture (ARCH)
- [ ] **ARCH-01**: Establish Quickshell QML frontend connected to Go/Rust backend daemon via JSON-RPC UNIX socket.
- [ ] **ARCH-02**: Support multi-compositor design with primary support for Niri.

### Backend (BACK)
- [ ] **BACK-01**: Implement Niri IPC daemon in backend to handle workspace and window events.
- [ ] **BACK-02**: Implement D-Bus services (Bluez, NetworkManager, logind) in backend.
- [ ] **BACK-03**: Implement hardware polling (sysfs, DDC/CI) for battery and brightness.

### Frontend UI (FRONT)
- [ ] **FRONT-01**: Implement Caelestia-styled Top Bar with workspaces, clock, and system tray.
- [ ] **FRONT-02**: Implement Caelestia-styled App Launcher (dynamic grid/list of installed apps).
- [ ] **FRONT-03**: Implement Caelestia-styled Control Center/Dashboard (quick settings, media controls).
- [ ] **FRONT-04**: Implement Caelestia-styled Notification Daemon and OSDs (volume, brightness).

## v2 Requirements
- Multi-monitor full support
- External plugin system

## Out of Scope
- Tight coupling to Niri (must be modular).

## Traceability
