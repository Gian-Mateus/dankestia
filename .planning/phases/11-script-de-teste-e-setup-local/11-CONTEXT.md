# Phase 11: Script de Teste e Setup Local

## Domain
Esta fase entregará as atualizações no script `run_dev.sh` (ou script equivalente de setup local) para garantir que ele gerencia de forma segura o build e o lifecycle de execução do Dankestia, isolando o ambiente de teste de outros daemons em execução.

## Locked Decisions (Implementation)

1. **Lógica de Compilação Condicional**:
   - Manter a delegação para o `make`, que já faz o tracking de alterações nativamente e só recompila o necessário.
   - Adicionar uma flag `--force` ao script que vai executar um "clean build" quando acionada, forçando a recompilação total caso algo corrompa.

2. **Serviços Interrompidos**:
   - Para evitar conflito, o script deve interromper *todos os painéis, widgets e daemons* comuns (que não o próprio compositor Niri/Hyprland) antes de rodar o ambiente de desenvolvimento.
   - Serviços alvo incluem: `dankestia.service`, `quickshell`, `dms`, além de painéis de terceiros clássicos como `waybar`, `mako`, `dunst`, `swaync`, `swaybg`, `hyprpaper`, `wpaperd`, etc.
   - O trap (cleanup) deve tentar restaurar os serviços primários do usuário, ou pelo menos o próprio `dankestia.service`, ao ser encerrado.

3. **Gestão de Logs**:
   - Os logs de saída (stdout e stderr) tanto do Go quanto do Quickshell deverão ser duplicados (usando ferramentas como `tee`).
   - Eles devem aparecer em tempo real no terminal E ser salvos automaticamente no arquivo `.dankestia-dev.log` na raiz do projeto para depuração posterior.

## Canonical References
- `run_dev.sh` (Arquivo a ser modificado)
