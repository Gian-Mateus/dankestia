import xml.etree.ElementTree as ET

ts_file = "quickshell/i18n/dankestia_pt_BR.ts"
tree = ET.parse(ts_file)
root = tree.getroot()

# Dicionário de traduções essenciais
translations = {
    "Settings": "Configurações",
    "Language & region": "Idioma e região",
    "Language": "Idioma",
    "System language": "Idioma do sistema",
    "Weather": "Clima",
    "Units": "Unidades",
    "Temperature": "Temperatura",
    "System temperatures": "Temperaturas do sistema",
    "Time & date": "Hora e data",
    "Clock format": "Formato da hora",
    "24-hour": "24 horas",
    "12-hour": "12 horas",
    "Desktop": "Área de Trabalho",
    "Displays": "Telas",
    "Display": "Tela",
    "Monitors": "Monitores",
    "Keyboard": "Teclado",
    "Shortcuts": "Atalhos",
    "Resolution": "Resolução",
    "Refresh Rate": "Taxa de Atualização",
    "Scale": "Escala",
    "Position": "Posição",
    "Audio": "Áudio",
    "Bluetooth": "Bluetooth",
    "Network": "Rede",
    "Wi-Fi": "Wi-Fi",
    "Power": "Energia",
    "Restart": "Reiniciar",
    "Log out": "Encerrar Sessão",
    "Shut down": "Desligar",
    "Cancel": "Cancelar"
}

count = 0
for context in root.findall('context'):
    for message in context.findall('message'):
        translation_elem = message.find('translation')
        if translation_elem is not None and translation_elem.get('type') == 'unfinished':
            source_elem = message.find('source')
            if source_elem is not None and source_elem.text:
                source_text = source_elem.text
                if source_text in translations:
                    translation_elem.text = translations[source_text]
                    del translation_elem.attrib['type']
                    count += 1

tree.write(ts_file, encoding='utf-8', xml_declaration=True)
print(f"Tradução base concluída. {count} strings traduzidas localmente.")
