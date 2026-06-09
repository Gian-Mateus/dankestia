package luaconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModuleToRelPath(t *testing.T) {
	tests := map[string]string{
		"dankestia.binds":       filepath.Join("dankestia", "binds.lua"),
		"dankestia/binds-user":  filepath.Join("dankestia", "binds-user.lua"),
		"awesome/anim":    filepath.Join("awesome", "anim.lua"),
		"awesome.colors":  filepath.Join("awesome", "colors.lua"),
		" awesome.binds ": filepath.Join("awesome", "binds.lua"),
	}

	for input, want := range tests {
		if got := ModuleToRelPath(input); got != want {
			t.Fatalf("ModuleToRelPath(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestRequiresSkipsComments(t *testing.T) {
	if modules := Requires(`-- require("dankestia.binds")`); len(modules) != 0 {
		t.Fatalf("expected commented require to be ignored, got %#v", modules)
	}

	modules := Requires(`print("-- not a comment") require("dankestia.binds") -- require("ignored")`)
	if len(modules) != 1 || modules[0] != "dankestia.binds" {
		t.Fatalf("unexpected modules: %#v", modules)
	}
}

func TestRequiresTargetRecurses(t *testing.T) {
	tmpDir := t.TempDir()
	dankestiaDir := filepath.Join(tmpDir, "dankestia")
	if err := os.MkdirAll(dankestiaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dankestiaDir, "windowrules.lua")
	if err := os.WriteFile(filepath.Join(tmpDir, "hyprland.lua"), []byte(`require("dankestia.extra")`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dankestiaDir, "extra.lua"), []byte(`require("dankestia.windowrules")`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte(`-- rules`), 0o644); err != nil {
		t.Fatal(err)
	}

	if !RequiresTarget(filepath.Join(tmpDir, "hyprland.lua"), target, make(map[string]bool)) {
		t.Fatal("expected recursive require lookup to find target")
	}
}
