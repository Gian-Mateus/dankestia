import xml.etree.ElementTree as ET
from deep_translator import GoogleTranslator
import time
import os

ts_file = "quickshell/i18n/dankestia_pt_BR.ts"
tree = ET.parse(ts_file)
root = tree.getroot()

translator = GoogleTranslator(source='auto', target='pt')

cache = {}

for context in root.findall('context'):
    for message in context.findall('message'):
        translation_elem = message.find('translation')
        if translation_elem is not None and translation_elem.get('type') == 'unfinished':
            source_elem = message.find('source')
            if source_elem is not None and source_elem.text:
                source_text = source_elem.text
                if source_text in cache:
                    translated_text = cache[source_text]
                else:
                    try:
                        translated_text = translator.translate(source_text)
                        # Consertar interpolação Qt (ex: %1)
                        translated_text = translated_text.replace('% 1', '%1').replace('% 2', '%2')
                        cache[source_text] = translated_text
                        print(f"Translating: '{source_text}' -> '{translated_text}'")
                        time.sleep(0.05) # Rate limit
                    except Exception as e:
                        print(f"Error translating '{source_text}': {e}")
                        translated_text = source_text
                
                translation_elem.text = translated_text
                del translation_elem.attrib['type'] # Remove o atributo unfinished para dar como traduzido

tree.write(ts_file, encoding='utf-8', xml_declaration=True)
print("Tradução finalizada.")
