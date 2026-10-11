package esde

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestScummVMLauncher(t *testing.T) {
	romDir := t.TempDir()
	gameDir := filepath.Join(romDir, "Beneath a Steel Sky")
	writeFile(t, filepath.Join(gameDir, "sky.scummvm"), "sky:sky")
	writeFile(t, filepath.Join(gameDir, "SKY.DNR"), "data")

	got, ok, err := ScummVMLauncher(gameDir, romDir)
	folder := filepath.Join(romDir, "Beneath a Steel Sky.scummvm")
	if err != nil || !ok || got != folder {
		t.Fatalf("ScummVMLauncher = %q, %v, %v; want %q, true, nil", got, ok, err, folder)
	}

	stub := filepath.Join(folder, "Beneath a Steel Sky.scummvm")
	if b, err := os.ReadFile(stub); err != nil || string(b) != "sky:sky" {
		t.Errorf("stub %q = %q, %v; want the ID kept", stub, b, err)
	}
	if _, err := os.Stat(filepath.Join(folder, "SKY.DNR")); err != nil {
		t.Error("game data should move with the folder")
	}
	if _, err := os.Stat(gameDir); !os.IsNotExist(err) {
		t.Error("the original folder should be gone")
	}
	if _, err := os.Stat(filepath.Join(folder, "sky.scummvm")); !os.IsNotExist(err) {
		t.Error("the original stub name should be gone")
	}
}

func TestScummVMLauncher_AlreadyArranged(t *testing.T) {
	romDir := t.TempDir()
	gameDir := filepath.Join(romDir, "Sky.scummvm")
	writeFile(t, filepath.Join(gameDir, "Sky.scummvm"), "sky")

	got, ok, err := ScummVMLauncher(gameDir, romDir)
	if err != nil || !ok || got != gameDir {
		t.Errorf("ScummVMLauncher = %q, %v, %v; want %q, true, nil", got, ok, err, gameDir)
	}
}

func TestScummVMLauncher_LeavesUnusableGamesAlone(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
	}{
		{"no stub", map[string]string{"SKY.DNR": "data"}},
		{"empty stub", map[string]string{"sky.scummvm": ""}},
		{"two stubs", map[string]string{"a.scummvm": "a", "b.scummvm": "b"}},
		{"stub below the game data", map[string]string{"Sub/sky.scummvm": "sky", "SKY.DNR": "data"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			romDir := t.TempDir()
			gameDir := filepath.Join(romDir, "Game")
			for name, content := range tt.files {
				writeFile(t, filepath.Join(gameDir, name), content)
			}

			if _, ok, err := ScummVMLauncher(gameDir, romDir); ok || err != nil {
				t.Errorf("got ok=%v err=%v, want no launcher", ok, err)
			}
			for name := range tt.files {
				if _, err := os.Stat(filepath.Join(gameDir, name)); err != nil {
					t.Errorf("%s should be untouched: %v", name, err)
				}
			}
		})
	}
}

func TestScummVMLauncher_WontOverwrite(t *testing.T) {
	romDir := t.TempDir()
	gameDir := filepath.Join(romDir, "Sky")
	writeFile(t, filepath.Join(gameDir, "sky.scummvm"), "sky")
	writeFile(t, filepath.Join(romDir, "Sky.scummvm", "keep"), "x")

	if _, ok, err := ScummVMLauncher(gameDir, romDir); ok || err == nil {
		t.Errorf("got ok=%v err=%v, want an error and no launcher", ok, err)
	}
	if _, err := os.Stat(filepath.Join(gameDir, "sky.scummvm")); err != nil {
		t.Error("the game folder should be untouched")
	}
}
