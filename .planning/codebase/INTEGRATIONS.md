# Integrations

**Date:** 2026-06-07

## External Services & APIs
- **Google Lens**: Visual search for screen contents.
- **OCR (Tesseract)**: Text extraction from screen region.

## System APIs
- **DBus**:
  - Bluez (Bluetooth management)
  - NetworkManager (Network configuration)
  - logind (Session management and brightness)
- **Wayland / Niri**:
  - IPC with Niri compositor
  - `xdg-desktop-portal` for native screen sharing
- **Hardware APIs**:
  - DDC/CI (External monitor brightness via I2C)
  - `sysfs` (Internal screen brightness, hardware metrics)
  - `udev` / `evdev` (Device input / events)
- **Audio**: PipeWire / WirePlumber

## Data Storage
- `bbolt` (Go backend internal key-value store)

## Inter-Process Communication
- UNIX Socket JSON-RPC (`/tmp/dankestia.sock` or via `XDG_RUNTIME_DIR`) connecting QML frontend to Go/Rust backend.
