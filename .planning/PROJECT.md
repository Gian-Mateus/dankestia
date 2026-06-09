# Dankestia: Caelestia UI over DMS Backend

## What This Is

Dankestia is a Linux Wayland shell project that merges the robust backend architecture of Dank Material Shell (DMS) with the phenomenal visual design and animations of the Caelestia interface. The ultimate goal is to achieve complete feature parity with both reference projects, delivering a beautiful, highly-performant, and stable desktop experience.

## Core Value

To combine the technical stability and performance of the DMS backend (Go/Rust daemons, IPC, system integrations) with the rich aesthetics and UX of Caelestia (QML/Qt6 via Quickshell). It brings the best of both worlds into a single cohesive system for the open-source community.

## Context

- **Problem:** DMS has excellent architecture and robust backend daemons, but its UI might not be as visually striking as Caelestia. Caelestia has incredible aesthetics but might lack the robust, standalone backend daemon approach of DMS.
- **Audience:** The open-source community and Linux desktop enthusiasts who want a highly customizable, beautiful, and stable Wayland shell.
- **Current State:** Existing codebase mappings show a split-architecture design with a QML frontend and a Go/Rust daemon backend communicating via JSON-RPC over UNIX sockets. The reference source for DMS logic is in `./references/DankMaterialShell/` and for Caelestia aesthetics in `./references/shell/`.

## Requirements

### Validated

*(Inferred from existing architecture and codebase map)*
- ✓ Split-architecture design separating UI rendering (Frontend) from system-level state management (Backend).
- ✓ Backend daemon written in Go/Rust managing Wayland compositors (Niri), D-Bus (Bluez, NetworkManager), hardware sensors, and system utilities.
- ✓ Frontend written in QML (Qt6) via Quickshell for handling visual elements, animations, and layouts.
- ✓ IPC mechanism using JSON-RPC over a UNIX socket (`/tmp/dankestia.sock`) for bidirectional communication.

### Active

- [ ] Complete parity of all Caelestia UI features (panels, widgets, animations, layout) mapped to the Dankestia frontend.
- [ ] Complete parity of all DMS backend features (daemons, state management, system APIs) mapped to the Dankestia backend.
- [ ] Integration of the Caelestia UI components to consume data and send commands through the DMS-style backend IPC.

### Out of Scope

- [ ] Tight coupling with a single compositor. While Niri is first-class, the architecture must support multi-compositor modularity.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Merge DMS backend with Caelestia UI | Get the best stability and the best aesthetics | — Pending |
| Full parity for Version 1 | To ensure the community gets a complete and usable product | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-06-08 after initialization*
