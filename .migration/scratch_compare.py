import os
import difflib

def get_files(base_dir):
    files_dict = {}
    for root, _, files in os.walk(base_dir):
        for f in files:
            if f.endswith('.qml'):
                full_path = os.path.join(root, f)
                rel_path = os.path.relpath(full_path, base_dir)
                try:
                    with open(full_path, 'r', encoding='utf-8') as file:
                        content = file.read()
                    files_dict[rel_path] = content
                except Exception as e:
                    pass
    return files_dict

cael_files = get_files('references/shell')
dank_files = get_files('quickshell/dankestia')

# Filter out nexus
cael_filtered = {k: v for k, v in cael_files.items() if not k.startswith('modules/nexus/')}

exact_matches = []
missing_in_dank = []
renamed_candidates = []

for c_rel, c_content in cael_filtered.items():
    if c_rel in dank_files:
        exact_matches.append(c_rel)
    else:
        # It's missing exactly, let's look for content similarity
        best_match = None
        best_ratio = 0
        c_lines = c_content.splitlines()
        for d_rel, d_content in dank_files.items():
            if d_rel not in cael_filtered: # Only check against files that aren't exact matches for something else
                # Fast heuristic: compare lengths
                if len(c_content) == 0 and len(d_content) == 0:
                    ratio = 1.0
                elif len(c_content) == 0 or len(d_content) == 0:
                    ratio = 0.0
                elif abs(len(c_content) - len(d_content)) / max(len(c_content), len(d_content)) > 0.5:
                    ratio = 0.0 # Too different in size
                else:
                    d_lines = d_content.splitlines()
                    # difflib for similarity
                    sm = difflib.SequenceMatcher(None, c_lines, d_lines)
                    ratio = sm.quick_ratio()
                
                if ratio > best_ratio:
                    best_ratio = ratio
                    best_match = d_rel
                    
        if best_ratio > 0.4:
            renamed_candidates.append((c_rel, best_match, best_ratio))
        else:
            missing_in_dank.append(c_rel)

print(f"--- EXACT MATCHES: {len(exact_matches)} ---")
print(f"--- RENAMED CANDIDATES (>70% similarity): {len(renamed_candidates)} ---")
for c, d, r in sorted(renamed_candidates, key=lambda x: x[2], reverse=True):
    print(f"{c} -> {d} ({r:.2f})")

print(f"\n--- TRULY MISSING: {len(missing_in_dank)} ---")
for m in sorted(missing_in_dank):
    print(m)

