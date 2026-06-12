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

echo -e "${BLUE}Parando a interface original do sistema para não dar conflito...${NC}"
# Usa systemctl para evitar que o Linux reinicie o painel original automaticamente
systemctl --user stop dankestia.service || true
killall -9 dankestia quickshell dms 2>/dev/null || true

# Prepara uma armadilha para ligar a interface original de volta quando fecharmos o teste
cleanup() {
    echo -e "${GREEN}Restaurando sua interface original do sistema...${NC}"
    systemctl --user start dankestia.service || true
}
trap cleanup EXIT

# Exporta a variável para que o backend repasse pro quickshell filho
export QML2_IMPORT_PATH=$PWD/quickshell/plugin/build/qml

echo -e "${GREEN}Iniciando servidor Dankestia... Pressione Ctrl+C para desligar tudo.${NC}"
# Inicia em foreground. O backend já gerencia o quickshell por padrão.
./core/bin/dankestia run -c $PWD/quickshell
