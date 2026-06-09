# Phase 1: Architecture Foundation - Research

## Overview
This document captures the technical research necessary to plan Phase 1 of Dankestia, which focuses on setting up the split-architecture and IPC communication layer between the QML frontend (Quickshell) and the Go backend daemon (dankestia).

## Context Recap
- **Daemon Initialization:** The Go backend (`dankestia`) acts as the master process, initializing its internal server and spawning Quickshell as a child process.
- **Backend Language:** Go (`core/`).
- **Compositor Detection:** Automatic detection via environment variables (e.g., `NIRI_SOCKET` for Niri).

## Current Codebase Insights
The `core/` directory contains the Go backend migrated from DankMaterialShell, now under the `dankestia` binary name.
The `quickshell/` directory contains the QML frontend migrated from Caelestia.

### 1. Quickshell Invocation (`shell.go`)
- **Location:** `core/cmd/dankestia/shell.go`
- **Current Behavior:** It spawns `qs -p <configPath>`. In DMS, `configPath` pointed to `DankMaterialShell`'s `shell.qml`.
- **Required Change:** It must be updated to resolve the path to `quickshell/shell.qml` from the Dankestia project. The environment variable `DMS_SOCKET` (or `DANKESTIA_SOCKET`) needs to be passed correctly to the QML process.

### 2. IPC Server (`server.go`)
- **Location:** `core/internal/server/server.go`
- **Current Behavior:** Sets up a UNIX socket for JSON-RPC communication.
- **Required Change:** Ensure the socket path name is updated to reflect Dankestia (`/tmp/dankestia.sock` or similar, depending on XDG_RUNTIME_DIR). Ensure the QML frontend `NetworkConnection.qml` (or equivalent in Caelestia) connects to this exact socket.

### 3. Compositor Detection
- **Required Change:** Add logic in the backend initialization (possibly in `shell.go` or `server.go`) to detect the current compositor.
- **Detection Logic:**
  - If `NIRI_SOCKET` is present in the environment -> Niri
  - If `HYPRLAND_INSTANCE_SIGNATURE` is present -> Hyprland
  - Else fallback to generic Wayland.
- This detection should be logged and made available to internal server modules to toggle compositor-specific IPC.

## Validation Architecture
- **Verification strategy:** 
  1. Build the `dankestia` Go binary.
  2. Run the `dankestia run` command.
  3. Verify that the Quickshell UI launches successfully.
  4. Verify that the UI successfully connects to the UNIX socket (no connection error logs).

## Conclusion
The architectural changes are highly contained. We need to point the Go daemon to the new QML entry point, rename socket variables to reflect the new project name, and introduce the compositor detection logic based on environment variables.
