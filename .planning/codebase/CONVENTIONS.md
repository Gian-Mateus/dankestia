# Conventions

**Date:** 2026-06-07

## Code Style
- **Go Backend**: 
  - Standard Go formatting (`gofmt`).
  - Linting enforced via `golangci-lint` (using `revive` and `errcheck`).
  - Strict error checking, except for explicitly ignored best-effort operations (like connection closes and teardowns).
- **QML/JS**:
  - Declarative UI descriptions with minimal inline JavaScript.
  - Usage of specific formatting tools like `.clang-format` or `qmlformat`.
- **Rust Backend**:
  - Standard Rust formatting (`rustfmt`).
  - Follows standard Cargo practices.

## Naming Patterns
- QML components are capitalized (e.g., `CUtils.qml`, `DeviceList.qml`).
- Go internal modules correspond to feature domains (e.g., `clipboard`, `colorpicker`, `wayland`).

## Error Handling
- **Go**: Explicit error handling. Errors are returned and propagated. Ignored errors are explicitly listed in `.golangci.yml`.
