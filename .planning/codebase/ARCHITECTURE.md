# Architecture

**Date:** 2026-06-07

## Projeto: Dankestia (Objetivo e Escopo)
O objetivo principal do Dankestia é criar uma shell Linux unificada que extraia a robustez técnica do **Dank Material Shell (DMS)** e o design estético do **Caelestia**.

### Regras de Escopo (Crucial para a IA)
1. **Backend e Lógica (Inspirado no DMS):** Toda a arquitetura de backend, scripts de gerenciamento de estado, deamons de escuta de eventos e manipulação de sockets devem seguir estritamente o padrão estrutural e robusto do Dank Material Shell.
2. **Frontend e UI (Inspirado no Caelestia):** Toda a identidade visual, layout de painéis, estilos CSS/SCSS, animações e o design conceitual dos widgets devem seguir estritamente a estética fenomenal do Caelestia.
3. **Multi-compositor:** O core do backend herdado do DMS deve suportar a modularidade multi-compositora inspirada pelo Caelestia, mantendo o Niri como primeira classe, mas sem acoplamento rígido que impeça outros.

### Arquivos de Referência Mapeados
- Fontes originais do DMS para cópia lógica: `./references/DankMaterialShell/...`
- Fontes originais do Caelestia para cópia estética: `./references/shell/...`

## System Design
The Dankestia shell is built on a split-architecture design separating the heavy UI rendering from the system-level state management and communication:
- **Frontend Layer**: A pure QML/Qt6 application built on top of the Quickshell Wayland shell framework. It handles all visual elements, animations, layout, and user interactions.
- **Backend Layer**: A headless Go/Rust daemon (`core`/`core-rust`) that runs in the background. It interfaces directly with Wayland compositors (Niri), D-Bus (Bluez, NetworkManager), hardware sensors (sysfs), and provides utilities (OCR, screenshots).

## Data Flow
- **IPC Mechanism**: The QML frontend communicates with the backend daemon via a JSON-RPC interface over a UNIX socket (`/tmp/dankestia.sock`).
- **Telemetry**: The backend polls or subscribes to hardware and system metrics and broadcasts updates to the frontend over the socket.
- **Commands**: The frontend issues commands (e.g., toggle launcher, open clipboard, set brightness) via the IPC to the backend, which executes the corresponding system calls or shell commands.

## Key Abstractions
- **Services (Frontend)**: QML modules in `quickshell/dankestia/services/` that encapsulate IPC calls and state logic for specific domains (e.g., NiriIPC, Bluetooth).
- **Internal Modules (Backend)**: Go packages in `core/internal/` (e.g., `clipboard`, `server`, `desktop`, `config`, `wayland`) abstracting specific system interactions.

## Entry Points
- `install.sh` / `install_dependencies.sh`: Shell scripts to bootstrap the environment.
- `dankestia-start.sh`: Shell script used to start the daemon and QML UI locally.
- `core/cmd/dankestia`: The main Go entry point for the backend daemon and CLI commands.
- `core-rust/src/main.rs`: The experimental Rust rewrite entry point.
- `quickshell/shell.qml`: The root entry point for the QML frontend.
