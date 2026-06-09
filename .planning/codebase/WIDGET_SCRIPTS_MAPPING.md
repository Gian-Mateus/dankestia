# Mapeamento de Scripts, Funções e Processos (Fase 6)

Durante a auditoria da Fase 6, mapeamos as funções nativas C++, possíveis scripts (Python/Bash) e chamadas explícitas no layout (`Process { ... }`) que devem ser substituídas por integrações via IPC no backend Go.

## 1. Funções Nativas C++ (Caelestia Plugin)
Muitos recursos que antes necessitavam de C++ nativo agora podem ser resolvidos pelos módulos do Dankestia:

| Função Original (C++) | Descrição | Substituto Ideal no Backend Go |
| :--- | :--- | :--- |
| `cpu.cpp` / `memory.cpp` | Uso de CPU e Memória RAM | Serviço `sysinfo` recém criado na Fase 5 |
| `gpu.cpp` | Carga e temperatura da GPU | Serviço `sysinfo` (A expandir para suportar GPU) |
| `storage.cpp` / `diskinfo.cpp`| Uso de partições de disco | Serviço `sysinfo` (A expandir para leitura do Filesystem) |
| `sensorslib.cpp` | Leitura de Fans e LM-Sensors | Serviço `sysinfo` (Usar APIs Linux nativas) |
| `appdb.cpp` | Indexação de `.desktop` para Launcher | Serviço `AppPickerManager` do Dankestia (já existente) |
| `logindmanager.cpp` | Reboot, Shutdown, Lock (DBus) | Serviço `loginctl` do Dankestia (já existente) |
| `audiocollector.cpp` (Cava) | Visualizador de Áudio (FFT) | Serviço IPC de Áudio/Pipewire a criar ou adaptar via Go |
| `lyrics.cpp` | Letras de Música do MPRIS | Backend MPRIS em Go (fazer fetch de letras HTTP/DBus) |
| `hyprextras.cpp` | Dados do Wayland/Hyprland | Serviço `wlroutput` e `dwl` (ou Hyprland específico) do Dankestia |
| `qalculator.cpp` | Cálculos matemáticos inline | Criar serviço simples `calculator` no Backend Go com Go-Math |

## 2. Processos Disparados pelo Layout via QML (`Process { ... }`)
O layout atual utiliza componentes `Process` e `CommandProcess` para invocar binários de terminal. A conversão destes para o Backend Go deixa o QML mais limpo, eficiente e sem bloqueios de subshell.

| Arquivo / Módulo (`quickshell/`) | Comando / Binário Invocado | Substituto Ideal no Backend Go |
| :--- | :--- | :--- |
| `services/Nmcli.qml` | `nmcli` | Serviço IPC nativo `network` (NetworkManager via DBus) |
| `services/Network.qml` | `ping` ou checagens de rede | Serviço IPC `network` |
| `services/Brightness.qml` | `brightnessctl` | Serviço IPC nativo `brightness` (já implementado no Dankestia) |
| `services/VPN.qml` | `tailscale` e `openvpn` | Expandir o serviço IPC `network` ou usar `tailscale` manager em Go |
| `services/Recorder.qml` | `wf-recorder` | Implementar `recorder` ou gerenciar via Systemd/Go |
| `modules/nexus/pages/AboutPage.qml` | `uname -r`, bash scripts | Serviço IPC `sysinfo` (fornecer nome do OS e Kernel) |
| `modules/launcher/services/Schemes.qml`| `matugen` | Integrar as chamadas do matugen nas rotinas Go do Dankestia |
| `modules/bar/popouts/kblayout/KbLayoutModel.qml` | `hyprctl` (layouts de teclado) | Serviço `wlroutput` / `hyprland` via Socket Wayland no Go |
| `modules/areapicker/Picker.qml` | `slurp` / `grim` | Utilitários de screenshot (Dankestia já possui `commands_screenshot.go`) |
| `modules/lock/Pam.qml` | `pamtester` (ou similar) | Interface PAM via CGO ou utilitário setuid standalone do Dankestia |

## Conclusão e Próximos Passos (Fase 7)
Para a Fase 7, os focos de criação dos novos scripts ou adaptações no backend Go devem ser:
1. **Concluir as pontes dos Módulos Faltantes**: Aproveitar os serviços Go existentes (`loginctl`, `brightness`, `network`) para os Singletons que ainda estão vazios (`Power.qml`, `DBus.qml`, etc).
2. **Expandir o `SysInfo`**: Ler Temperatura, Uso de Discos e Dados do Kernel (`AboutPage`).
3. **Migrar `Nmcli.qml` e `Brightness.qml`**: Fazer com que utilizem os canais do `DankestiaIPC` (Go) em vez de executarem `Process` localmente.
