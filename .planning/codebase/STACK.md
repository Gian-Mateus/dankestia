# Tech Stack

**Date:** 2026-06-07

## Languages & Frameworks
- **Frontend**: QML (Qt6) via Quickshell
- **Backend (Daemon)**: Go (v1.26.1)
- **Backend (Experimental/Re-write)**: Rust (edition 2021)
- **Scripting**: Bash

## Runtimes & Environments
- **Wayland Compositor**: Niri
- **Package Manager**: Nix (`flake.nix`), Go Modules, Cargo

## Key Dependencies
- **Go**: 
  - `github.com/Wifx/gonetworkmanager`
  - `github.com/godbus/dbus/v5`
  - `github.com/holoplot/go-evdev`
  - `github.com/pilebones/go-udev`
  - `github.com/sblinch/kdl-go`
  - `github.com/spf13/cobra`
  - `go.etcd.io/bbolt`
  - `tailscale.com`
- **Rust**:
  - `tokio` (async)
  - `sysinfo` (telemetry)
  - `nix` (socket bindings)
  - `clap`, `serde`
- **QML**:
  - `quickshell` (Wayland/X11 dynamic shell framework)
  - `dankestia-cli`

## System Integrations
- DBus (Bluez, NetworkManager, logind)
- PipeWire / WirePlumber
- `sysfs` / DDC/CI (Hardware metrics & brightness)
- IPC via UNIX Socket JSON-RPC
