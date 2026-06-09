# Integração e Dependências Nativas (Caelestia QML)

Este documento cataloga de forma detalhada como o frontend puro do Caelestia (`quickshell/`) busca dados do sistema.

## 1. Visão Geral
Durante a análise do código (`grep` por `bash`, `sh -c`, `dbus-send`, `qdbus` e `/sys/`), constatou-se que o Caelestia **não realiza chamadas sujas de terminal ou DBus via QML**. 

Todo o acesso ao sistema operacional foi perfeitamente encapsulado em um **Plugin C++ Nativo (Caelestia Plugin)**, localizado em `quickshell/plugin/src/Caelestia/`.

Isso significa que, para o QML, as informações do sistema aparecem magicamente através de imports como `import Caelestia.Services` e `import Caelestia.Config`.

## 2. Mapa do Plugin C++ (O que precisaremos recriar/adaptar no Go)

Para a Fase 5 (Padrão Adapter/Bridge), nós substituiremos ou faremos uma ponte das seguintes capacidades do C++ para o backend Dankestia (Go):

### 2.1. Serviços de Hardware e Monitoramento (`Services/`)
- **`cpu.cpp` / `memory.cpp`:** Leitura de carga de CPU e uso de RAM.
- **`gpu.cpp`:** Monitoramento de uso e temperatura da placa de vídeo.
- **`storage.cpp` / `diskinfo.cpp`:** Leitura das partições do sistema de arquivos e espaço livre.
- **`sensorslib.cpp`:** Integração com o `lm-sensors` para leitura de ventoinhas e temperatura de componentes.

### 2.2. Áudio e Visualização de Mídia
- **`audiocollector.cpp` / `cavaprovider.cpp`:** Coleta dados do áudio do sistema (PulseAudio/Pipewire) para alimentar os visualizadores do Caelestia.
- **`lyrics.cpp` / `lyriccandidate.cpp`:** Busca e sincroniza letras de músicas ativas no MPRIS.
- **`beattracker.cpp`:** Algoritmo que detecta as batidas da música para pulsar a interface.

### 2.3. Lógica do Desktop e Sistema (`Internal/` e Raiz do Plugin)
- **`logindmanager.cpp`:** Única interface que faz chamadas DBus (`qdbus`) diretamente para o `org.freedesktop.login1` (usado para suspender, reiniciar e desligar a máquina).
- **`appdb.cpp`:** Lê e indexa arquivos `.desktop` em `/usr/share/applications` para a gaveta de aplicativos (Launcher).
- **`qalculator.cpp`:** Integra a biblioteca libqalculate para resolver equações matemáticas diretamente pela barra de pesquisa.
- **`requests.cpp`:** Realiza chamadas HTTP GET nativas (muito possivelmente para buscar clima ou updates).
- **`hyprextras.cpp` / `hyprdevices.cpp`:** Interação nativa com o Wayland/Hyprland (ex: lista de teclados conectados e mapeamento de periféricos).

### 2.4. Custom Qt Quick Items (Frontend)
Itens que são renderizados em C++ por motivos de performance (não precisaremos transferir isso para o Go, apenas manter no C++ se necessário, ou usar os equivalentes do Quickshell):
- `visualiserbars.cpp` e `sparklineitem.cpp`: Gráficos de barra desenhados no Canvas nativo do Qt.
- `circularindicatormanager.cpp`: Gerenciamento otimizado de anéis de progresso.

## 3. Conclusão da Fase 4
A arquitetura purista do Caelestia facilita muito a migração. O nosso Backend Go (`core/`) já possui 90% dessas funcionalidades implementadas via IPC. 

**O desafio da Fase 5 será:** Criar adaptadores no QML que finjam ser os objetos C++ `Caelestia.Services`, mas que nos bastidores assinem (subscribe) os dados do nosso Socket IPC (`DANKESTIA_SOCKET`) alimentado pelo Go.
