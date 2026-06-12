import os
import re

search_dir = "quickshell/"
output_file = ".planning/codebase/TRANSLATIONS_MAPPING.md"

# Regex to find text: "...", label: "...", title: "...", description: "..."
pattern = re.compile(r'(text|label|title|description|placeholderText|tooltip)\s*:\s*"([^"]+)"')

candidates = []

def is_icon_or_symbol(s):
    if len(s) == 0:
        return True
    if s in ["dankestiafetch.sh", "Dankestia", "Dankestia", "°C", "°F", "•••", "•", ">"]:
        return True
    # If all lowercase with underscores, highly likely a material icon
    if re.match(r'^[a-z_]+$', s):
        return True
    return False

for root, dirs, files in os.walk(search_dir):
    for f in files:
        if f.endswith(".qml"):
            filepath = os.path.join(root, f)
            with open(filepath, 'r') as file:
                lines = file.readlines()
                for i, line in enumerate(lines):
                    # Skip if qsTr is already present
                    if "qsTr(" in line:
                        continue
                    
                    matches = pattern.finditer(line)
                    for match in matches:
                        prop = match.group(1)
                        val = match.group(2)
                        
                        if not is_icon_or_symbol(val):
                            candidates.append({
                                'file': filepath,
                                'line': i + 1,
                                'prop': prop,
                                'val': val
                            })

with open(output_file, 'w') as out:
    out.write("# Mapeamento de Textos Hardcoded para qsTr()\n\n")
    out.write("| Arquivo | Linha | Propriedade | Texto Original | Sugestão | Ação |\n")
    out.write("| :--- | :--- | :--- | :--- | :--- | :--- |\n")
    
    for c in candidates:
        sug = f'qsTr("{c["val"]}")'
        out.write(f"| `{c['file']}` | {c['line']} | `{c['prop']}` | `{c['val']}` | `{sug}` | [ ] Traduzir |\n")

print(f"Encontrados {len(candidates)} textos hardcoded. Documento salvo em {output_file}")
