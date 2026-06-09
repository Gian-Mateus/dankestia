# Phase 1: Architecture Foundation - Context

**Gathered:** 2026-06-08
**Status:** Ready for planning

<domain>
## Phase Boundary

Estabelecer a camada base de comunicação IPC (JSON-RPC via UNIX socket) e a inicialização entre a interface QML e o backend em Go.

</domain>

<decisions>
## Implementation Decisions

### Inicialização do Daemon e UI
- **D-01:** O backend atuará como o processo "mestre". O sistema rodará o backend (dms daemon), e o backend fará o spawn da interface (Caelestia/Quickshell) como um processo filho, monitorando e controlando o ciclo de vida da UI.

### Linguagem do Backend
- **D-02:** O backend adotado será em **Go**, garantindo estabilidade por usar a base exaustivamente testada do DankMaterialShell. A reescrita experimental em Rust será ignorada por enquanto.

### Detecção de Compositor
- **D-03:** A detecção do compositor (foco no Niri) será feita de forma automática em tempo de execução via variáveis de ambiente (como `NIRI_SOCKET` ou `WAYLAND_DISPLAY`), eliminando a necessidade de configuração manual pelo usuário.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Backend & Inicialização
- `references/DankMaterialShell/core/cmd/dms/shell.go` — Lógica original de inicialização de processos e IPC do DMS.
- `references/DankMaterialShell/core/internal/server/server.go` — Servidor de socket JSON-RPC.

### Frontend
- `references/shell/shell.qml` — Arquivo raiz da interface Caelestia.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `core/cmd/dms/shell.go`: Pode ser adaptado/reutilizado para gerenciar a inicialização do Quickshell via `exec.Command`.
- Servidor IPC em `core/internal/server`: Lógica já pronta para expor o UNIX socket.

### Established Patterns
- **Ciclo de vida atrelado:** O `shell.go` do DMS já escuta por sinais de sistema para fechar o processo filho.

### Integration Points
- O Quickshell recebe a variável de ambiente `DMS_SOCKET` apontando para o socket de comunicação ao ser invocado pelo Go.

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-architecture-foundation*
*Context gathered: 2026-06-08*
