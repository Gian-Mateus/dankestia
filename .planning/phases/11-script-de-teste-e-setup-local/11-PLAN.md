# Phase 11 Plan: Script de Teste e Setup Local

## Objective
Garantir que o script `run_dev.sh` isole o ambiente de testes do Dankestia, gerenciando serviços, recompilação limpa condicional e logs de depuração.

## Implementation Steps

### 1. Adicionar compilação com `--force` (run_dev.sh)
- Modificar o início do `run_dev.sh` para aceitar a flag `--force`.
- Se `--force` for passado, rodar `make clean` no backend Go e limpar o cache do `cmake` (ex: `rm -rf quickshell/plugin/build/*`) antes do fluxo normal de build.

### 2. Isolar o ambiente de testes (run_dev.sh)
- Expandir a linha de `killall -9` para incluir todos os daemons Wayland comuns que podem sobrepor ou conflitar com o Dankestia:
  - Painéis: `waybar`, `ironbar`, `ags`, `eww`, `polybar`
  - Notificações: `mako`, `dunst`, `swaync`
  - Wallpapers: `swaybg`, `hyprpaper`, `wpaperd`
- Manter o `killall -9 dankestia quickshell dms` e `systemctl --user stop dankestia.service`.
- Atualizar o `cleanup` (trap) para garantir que o `dankestia.service` seja reiniciado via `systemctl --user start dankestia.service`.

### 3. Implementar logs divididos (`tee`)
- Alterar a execução do Dankestia:
  ```bash
  ./core/bin/dankestia run -c $PWD/quickshell 2>&1 | tee .dankestia-dev.log
  ```
- Isso permitirá ver os logs no terminal E guardar um registro completo no arquivo `.dankestia-dev.log`.

## Verification
- Executar `./run_dev.sh` e verificar se o arquivo `.dankestia-dev.log` é criado e alimentado.
- Executar `./run_dev.sh --force` e verificar se `make clean` é invocado.
- Garantir que processos como `waybar` sejam encerrados ao rodar o script e que o trap volte a ligar o `dankestia.service` original após o encerramento com `Ctrl+C`.
