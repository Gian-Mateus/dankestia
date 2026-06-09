# ROADMAP: Dankestia (Caelestia + DMS Backend)

## Phase 1: Bootstrap do Projeto (Cópia Base)
Goal: Copiar todo o conteúdo do `references/DankMaterialShell` para a raiz do repositório.
Requirements: BACK-01
Success criteria:
1. Pasta core/ e scripts base estão no repositório.
2. Projeto compila e serviços base do DMS rodam.

## Phase 2: Transplante do Frontend (Caelestia)
Goal: Substituir a pasta `quickshell` nativa do DMS pela interface QML do `references/shell` (Caelestia).
Requirements: FRONT-01
Success criteria:
1. O repositório contém o backend DMS puro e o frontend Caelestia puro.

## Phase 3: Rebranding Global
Goal: Renomear todas as referências de "DMS" e "DankMaterialShell" para "Dankestia" em todos os módulos.
Requirements: BACK-01, FRONT-01
Success criteria:
1. Nenhuma string "DMS_" de ambiente ou caminhos de socket sobrou.
2. Compilação bem sucedida sob a nova marca.

## Phase 4: Mapeamento de Dependências Nativas do Caelestia
Goal: Catalogar e documentar como a pasta `quickshell/` do Caelestia busca dados do sistema nativamente.
Requirements: FRONT-02
Success criteria:
1. Documento de integração com todas as requisições (dbus, bash, sysfs) detalhado.

## Phase 5: Padrão Adapter/Bridge (QML -> Go)
Goal: Interceptar as requisições do Caelestia e redirecioná-las ao Backend Go via IPC/Socket.
Requirements: IPC-01, FRONT-03
Success criteria:
1. Serviços de rede, brilho, bateria e afins do Caelestia recebem dados do Go.
2. A interface não trava por falta de dados vitais de hardware.

## Phase 6: Auditoria de Lógica Externa (Scripts)
Goal: Encontrar e listar todos os scripts em Python/C++/Bash que o Caelestia usava para widgets pesados.
Requirements: BACK-02
Success criteria:
1. Lista documentada de lógicas órfãs ou que devem ir para o backend.

## Phase 7: Migração da Lógica Externa para o Backend Go
Goal: Desenvolver módulos no Go para realizar as tarefas descobertas na Fase 6.
Requirements: BACK-03
Success criteria:
1. Nenhuma dependência externa pesada do Caelestia é mantida no QML.
2. Dados complexos fluem do Go para os widgets.

## Phase 8: Internacionalização (i18n)
Goal: Mapear textos estáticos do layout QML para funções de tradução (`qsTr`).
Requirements: FRONT-04
Success criteria:
1. Todo o layout suporta múltiplos idiomas dinamicamente.

## Phase 9: Quality Assurance & Empacotamento
Goal: Testes E2E, verificação de memory leaks e empacotamento.
Requirements: SYS-01
Success criteria:
1. Shell inicia sem crashes.
2. Instalação e distribuição confirmadas.
