# Directory Structure

**Date:** 2026-06-07

## Root Directories

- `core/`: The production backend daemon written in Go.
  - `cmd/dankestia/`: Main application entry point for the Go daemon/CLI.
  - `internal/`: Internal Go packages (IPC server, bluetooth, wayland, clipboard, utils).
- `core-rust/`: The experimental rewrite of the backend in Rust.
  - `src/`: Rust source files.
- `quickshell/`: The frontend QML UI components and Quickshell application logic.
  - `dankestia/components/`: Reusable QML visual components.
  - `dankestia/config/`: Configuration definitions.
  - `dankestia/modules/`: High-level UI modules (e.g., control center, launcher).
  - `dankestia/services/`: QML controllers abstracting the IPC interface to the backend.
  - `dankestia/translations/`: `.ts` / `.qm` files for i18n (PT-BR).
  - `dankestia/utils/`: QML/JS utility scripts.
  - `nix/`: Nix package derivations for the shell.
- `assets/`: Media, icons, and systemd service files.
- `scripts/` (inside `quickshell`): Additional shell scripts for integrations.

## Key Files
- `README.md`: Primary project documentation.
- `Makefile`: Build and run instructions.
- `flake.nix`: Nix configuration for building and development environment.
- `install.sh`: Central installation script.
- `dankestia-start.sh`: Dev entry point to run the system.
