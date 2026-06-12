#!/bin/bash
set -e

# Cores para o terminal
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Construindo o Backend Go ===${NC}"
make -C core build

echo -e "${BLUE}=== Construindo o Plugin C++ do Quickshell ===${NC}"
mkdir -p quickshell/plugin/build
cd quickshell/plugin/build
cmake ..
make
cd ../../../

echo -e "${GREEN}=== Build concluído! Iniciando o Dankestia ===${NC}"

# Mata processos anteriores do dankestia se existirem
pkill -x dankestia || true

# Inicia o backend em segundo plano (background)
./core/bin/dankestia run &
BACKEND_PID=$!

echo -e "${GREEN}Backend Go rodando no PID: $BACKEND_PID${NC}"
echo -e "${GREEN}Abrindo Interface Gráfica...${NC}"

# Inicia o Quickshell bloqueando o terminal principal
QML2_IMPORT_PATH=$PWD/quickshell/plugin/build/qml quickshell -p quickshell/

# Assim que a janela do Quickshell for fechada pelo usuário, o script continua e mata o backend
echo -e "${BLUE}Quickshell encerrado. Desligando o servidor backend...${NC}"
kill $BACKEND_PID 2>/dev/null || true
echo -e "${GREEN}Ambiente de desenvolvimento encerrado com sucesso!${NC}"
