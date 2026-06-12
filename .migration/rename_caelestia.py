import os

# Directories to ignore
IGNORE_DIRS = {'.git', 'references', 'docs', '.agent', '.planning', '.gemini'}

def process_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except UnicodeDecodeError:
        return False

    new_content = content.replace('Dankestia', 'Dankestia')
    new_content = new_content.replace('dankestia', 'dankestia')
    new_content = new_content.replace('DANKESTIA', 'DANKESTIA')

    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        return True
    return False

def rename_in_directory(root_dir):
    modified_count = 0
    for root, dirs, files in os.walk(root_dir):
        # Modify dirs in-place to avoid walking ignored directories
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
        
        for file in files:
            filepath = os.path.join(root, file)
            # only process text files, avoid binaries
            if file.endswith(('.go', '.qml', '.js', '.py', '.md', '.sh', '.service', '.conf', '.desktop', 'Makefile', '.cpp', '.hpp', '.c', '.h')):
                if process_file(filepath):
                    modified_count += 1
                    print(f"Modificado: {filepath}")

    print(f"\nTotal de arquivos modificados: {modified_count}")

if __name__ == "__main__":
    rename_in_directory('.')
