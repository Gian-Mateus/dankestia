import xml.etree.ElementTree as ET
from deep_translator import GoogleTranslator
import time

def translate_file(filepath):
    tree = ET.parse(filepath)
    root = tree.getroot()
    
    translator = GoogleTranslator(source='en', target='pt')
    
    count = 0
    total = len(root.findall('.//message'))
    
    print(f"Total messages to process: {total}")
    
    # Pre-caching common words to avoid API spam
    cache = {
        "Cancel": "Cancelar",
        "OK": "OK",
        "Save": "Salvar",
        "Apply": "Aplicar",
        "Close": "Fechar",
        "Settings": "Configurações",
        "System": "Sistema",
        "Network": "Rede",
        "Appearance": "Aparência",
        "Notifications": "Notificações",
        "Bluetooth": "Bluetooth",
        "Audio": "Áudio",
        "General": "Geral",
        "Time": "Hora",
        "Date": "Data",
        "Search": "Buscar",
        "Enable": "Habilitar",
        "Disable": "Desabilitar",
        "Enabled": "Habilitado",
        "Disabled": "Desabilitado",
        "Volume": "Volume",
        "Brightness": "Brilho",
        "Battery": "Bateria",
        "Charging": "Carregando",
        "Disconnect": "Desconectar",
        "Connect": "Conectar",
        "Connected": "Conectado",
        "Disconnected": "Desconectado"
    }
    
    for message in root.findall('.//message'):
        source_elem = message.find('source')
        translation_elem = message.find('translation')
        
        if source_elem is not None and translation_elem is not None:
            if translation_elem.get('type') == 'unfinished':
                original_text = source_elem.text
                if not original_text:
                    continue
                
                # Use cache or translate
                if original_text in cache:
                    translated_text = cache[original_text]
                elif original_text.strip() == "":
                    translated_text = original_text
                else:
                    try:
                        translated_text = translator.translate(original_text)
                        time.sleep(0.1) # small sleep to avoid rate limiting
                    except Exception as e:
                        print(f"Error translating '{original_text}': {e}")
                        translated_text = original_text # fallback
                
                translation_elem.text = translated_text
                del translation_elem.attrib['type']
                count += 1
                
                if count % 50 == 0:
                    print(f"Translated {count} strings...")
                    
    tree.write(filepath, encoding='utf-8', xml_declaration=True)
    print(f"Finished translating {count} strings!")

if __name__ == '__main__':
    translate_file('quickshell/dankestia/i18n/pt_BR.ts')
