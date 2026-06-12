import os

def is_text_file(filepath):
    try:
        with open(filepath, 'rt', encoding='utf-8') as f:
            f.read(1024)
        return True
    except Exception:
        return False

def replace_in_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content
    # Replace any remaining DANKESTIA with DANKESTIA
    content = content.replace('DANKESTIA', 'DANKESTIA')
    content = content.replace('dankestia', 'dankestia')

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    return False

def main():
    skip_dirs = {'.git', 'references', '.planning', '.agent', '.gemini', 'docs'}
    for root, dirs, files in os.walk('.', topdown=True):
        dirs[:] = [d for d in dirs if d not in skip_dirs]
        for name in files:
            filepath = os.path.join(root, name)
            if is_text_file(filepath):
                replace_in_file(filepath)

if __name__ == '__main__':
    main()
