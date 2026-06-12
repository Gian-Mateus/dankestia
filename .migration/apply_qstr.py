import os
import re

search_dir = "quickshell/"

# Match: property: "String" or property: 'String'
# properties: text, title, label, description, placeholderText, tooltip, displayName
pattern = re.compile(r'(text|title|label|description|placeholderText|tooltip|displayName)\s*:\s*("([^"]+)"|\'([^\']+)\')')

files_to_update = {}
total_found = 0

def is_icon_or_symbol(s):
    if len(s) == 0:
        return True
    # Ignore single chars like ":" or symbols
    if s in ["dankestiafetch.sh", "Dankestia", "Dankestia", "°C", "°F", "•••", "•", ">", "/", ":", " "]:
        return True
    # Ignore material icons (all lowercase, underscores)
    if re.match(r'^[a-z0-9_]+$', s):
        return True
    return False

for root, dirs, files in os.walk(search_dir):
    for f in files:
        if f.endswith(".qml"):
            filepath = os.path.join(root, f)
            with open(filepath, 'r') as file:
                lines = file.readlines()
                
            modified = False
            for i, line in enumerate(lines):
                if "qsTr(" in line:
                    continue
                
                # Check for matches
                matches = list(pattern.finditer(line))
                if not matches:
                    continue
                    
                # We need to replace in reverse order to not mess up indices
                new_line = line
                for match in reversed(matches):
                    prop = match.group(1)
                    val = match.group(3) if match.group(3) is not None else match.group(4)
                    full_str = match.group(2) # "..." or '...'
                    
                    if not is_icon_or_symbol(val):
                        # Construct replacement: qsTr("val")
                        # Always use double quotes for qsTr
                        replacement = f'qsTr("{val}")'
                        start, end = match.span(2)
                        new_line = new_line[:start] + replacement + new_line[end:]
                        modified = True
                        total_found += 1
                        print(f"[{filepath}:{i+1}] {prop}: {full_str} -> {replacement}")
                
                if modified:
                    lines[i] = new_line
                    
            if modified:
                files_to_update[filepath] = lines

print(f"\nTotal items to translate: {total_found}")
print(f"Files to modify: {len(files_to_update)}")

# Uncomment to actually write the changes:
for filepath, lines in files_to_update.items():
    with open(filepath, 'w') as f:
        f.writelines(lines)
print("Changes applied!")
