# Phase 1: Architecture Foundation - Discussion Log

**Date:** 2026-06-08

## Area: Inicialização do Daemon e UI
**Options Presented:**
- Via script de shell (ex: dankestia-start.sh) que inicia ambos
- Via systemd user services
- O backend (daemon) inicia a UI como um processo filho

**User Selection:**
O backend (daemon) inicia a UI como um processo filho (Abordagem do DMS).

**Notes:**
Usuário prefere o padrão do DMS de usar o backend como supervisor da UI.

## Area: Linguagem do Backend
**Options Presented:**
- Manter o backend em Go (Recomendado)
- Migrar totalmente para o Rust

**User Selection:**
Manter o backend em Go.

**Notes:**
Garante estabilidade para a versão 1 do Dankestia.

## Area: Detecção de Compositor
**Options Presented:**
- Automático via variáveis de ambiente
- Explícito via arquivo de configuração

**User Selection:**
Automático via variáveis de ambiente.

**Notes:**
Usuário notou que era a forma mais inteligente baseada na arquitetura de Wayland.
