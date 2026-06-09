# Concerns

**Date:** 2026-06-07

## Technical Debt & Ongoing Refactoring
- **Backend Port**: The project is currently maintaining a Go backend (`core/`) while actively porting it to Rust (`core-rust/`). This dual-maintenance creates potential feature divergence until the Rust version reaches parity.
- **Dependencies**: The Go backend has many dependencies for DBus, networking, and Wayland protocols that are complex to maintain and mock.
- **Wayland Fragility**: Deep integration with Wayland and Niri IPC can be brittle when compositor APIs change.
- **Hardware Integration**: Integrating directly with `sysfs` and DDC/CI (I2C) for brightness/hardware metrics usually requires correct user permissions/privileges (`udev` rules, `video`/`i2c` groups), which can lead to friction during installation.

## Known Issues
- Translating the QML frontend is managed via scripts (`translate_ts.py`), which might become unwieldy.
