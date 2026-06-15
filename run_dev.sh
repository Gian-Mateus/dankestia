#!/bin/bash
set -e

# Cores para o terminal
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

FORCE_BUILD=0
if [ "$1" == "--force" ]; then
    FORCE_BUILD=1
    echo -e "${BLUE}=== Modo --force ativado: Limpando builds anteriores ===${NC}"
    make -C core clean
    rm -rf quickshell/plugin/build/*
fi

echo -e "${BLUE}=== Construindo o Backend Go ===${NC}"
make -C core build

echo -e "${BLUE}=== Construindo o Plugin C++ do Quickshell ===${NC}"
mkdir -p quickshell/plugin/build
cd quickshell/plugin/build
cmake .. -DFETCHCONTENT_UPDATES_DISCONNECTED=ON
make -j$(nproc)
cd ../../../

echo -e "${BLUE}Parando a interface original e serviços de terceiros para evitar conflitos...${NC}"
# Usa systemctl para evitar que o Linux reinicie o painel original automaticamente
systemctl --user stop dankestia.service 2>/dev/null || true

# Derruba o ambiente atual do Dankestia e outros serviços gráficos comuns no Wayland
killall -9 dankestia quickshell dms 2>/dev/null || true
killall -9 waybar ironbar ags eww polybar 2>/dev/null || true
killall -9 mako dunst swaync 2>/dev/null || true
killall -9 swaybg hyprpaper wpaperd 2>/dev/null || true

# Prepara uma armadilha para ligar a interface original de volta quando fecharmos o teste
cleanup() {
    echo -e "${GREEN}Restaurando sua interface original do sistema...${NC}"
    systemctl --user start dankestia.service || true
}
trap cleanup EXIT

# Exporta a variável para que o backend repasse pro quickshell filho
export QML2_IMPORT_PATH=$PWD/quickshell/plugin/build/qml

echo -e "${GREEN}Iniciando servidor Dankestia... Pressione Ctrl+C para desligar tudo.${NC}"
echo -e "${GREEN}Os logs estão sendo salvos em .dankestia-dev.log${NC}"
# Inicia em foreground. O backend já gerencia o quickshell por padrão.
# O uso de stdbuf garante que o output não fique preso no buffer e vá para o tee em tempo real
stdbuf -o0 -e0 ./core/bin/dankestia run -c $PWD/quickshell 2>&1 | tee .dankestia-dev.log
