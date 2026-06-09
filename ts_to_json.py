import xml.etree.ElementTree as ET
import json
import os

def ts_to_json(ts_file, json_file):
    tree = ET.parse(ts_file)
    root = tree.getroot()
    
    translations = {}
    
    for context in root.findall('context'):
        ctx_name = context.find('name').text if context.find('name') is not None else ""
        if ctx_name not in translations:
            translations[ctx_name] = {}
            
        for message in context.findall('message'):
            source = message.find('source')
            translation = message.find('translation')
            
            if source is not None and translation is not None and translation.text:
                translations[ctx_name][source.text] = translation.text
                # Also add it to the global context if not there
                if "" not in translations:
                    translations[""] = {}
                translations[""][source.text] = translation.text
                
    os.makedirs(os.path.dirname(json_file), exist_ok=True)
    with open(json_file, 'w', encoding='utf-8') as f:
        json.dump(translations, f, ensure_ascii=False, indent=2)

if __name__ == '__main__':
    ts_to_json('quickshell/dankestia/i18n/pt_BR.ts', 'quickshell/dankestia/translations/poexports/pt_BR.json')
    ts_to_json('quickshell/dankestia/i18n/pt_BR.ts', 'quickshell/dankestia/translations/poexports/pt-BR.json')
    ts_to_json('quickshell/dankestia/i18n/pt_BR.ts', 'quickshell/dankestia/translations/poexports/pt.json')
