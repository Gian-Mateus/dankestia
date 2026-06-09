# Testing

**Date:** 2026-06-07

## Go Backend (`core/`)
- **Framework**: Standard `go test` library augmented with `github.com/stretchr/testify` for assertions.
- **Mocks**: Heavy usage of mocks. Generated via `mockery` (configured via `.mockery.yml`) and stored in `internal/mocks/`.
- **Coverage**: Tests exist across most internal modules (e.g., `windowrules/providers/`, `server/freedesktop/`, `server/sysupdate/`, etc.). They end with `_test.go`.

## Rust Backend (`core-rust/`)
- Standard cargo tests (`cargo test`), though currently in experimental rewrite phase.

## UI Frontend (`quickshell/`)
- Manual visual regression testing and Wayland compositor integration testing.
- `test_translator.qml` exists for verifying localization.
