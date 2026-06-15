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

## Phase 10: Traduções pt_BR
Goal: Criar as traduções para pt_BR e aplicar ao projeto (complementando a internacionalização da Fase 8).
Requirements: FRONT-04
Success criteria:
1. A interface consegue ser exibida totalmente em Português do Brasil de forma funcional.

## Phase 11: Script de Teste e Setup Local
Goal: Verificar e ajustar o script de setup (`run_dev.sh` ou similar) para testar o ambiente localmente, parando serviços rodando na máquina, compilando se necessário, e restaurando os serviços anteriores ao fechar.
Requirements: SYS-02
Success criteria:
1. O script permite rodar e testar o Dankestia sem quebrar ou interferir permanentemente no desktop atual do usuário.

## Phase 12: Screenshooter do Caelestia
Goal: Verificar se o screenshooter original do Caelestia foi importado para o projeto, substituindo o atual do DMS.
Requirements: FRONT-05
Success criteria:
1. Funcionalidades superiores de screenshot do Caelestia estão totalmente integradas e operantes no Dankestia.

## Phase 13: Configurações de Displays
Goal: Adicionar controle de configurações de displays na interface.
Requirements: FRONT-06, BACK-04
Success criteria:
1. Gerenciamento de monitores/telas funcional via GUI.

## Phase 14: Layout de Teclado
Goal: Adicionar configuração de troca de layout de teclado na interface.
Requirements: FRONT-07, BACK-05
Success criteria:
1. O usuário consegue alterar o idioma/layout do teclado diretamente pela GUI.

## Phase 15: Gerenciamento de Atalhos via GUI
Goal: Adicionar configurações de atalhos (listagem, edição e criação de novos) totalmente gerenciáveis via interface gráfica.
Requirements: FRONT-08, BACK-06
Success criteria:
1. O usuário visualiza, edita e cria atalhos do compositor/sistema pela GUI de forma persistente.
