package muos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScummVMLauncher(t *testing.T) {
	romDir := t.TempDir()
	hidden := filepath.Join(romDir, "_Sky")
	if err := os.MkdirAll(hidden, 0755); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{"sky.scummvm": "sky:sky", "SKY.DNR": "data"} {
		if err := os.WriteFile(filepath.Join(hidden, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	got, ok, err := ScummVMLauncher(filepath.Join(romDir, "Sky"), romDir)
	want := filepath.Join(romDir, "Sky.scummvm")
	if err != nil || !ok || got != want {
		t.Fatalf("ScummVMLauncher = %q, %v, %v; want %q, true, nil", got, ok, err, want)
	}

	if info, err := os.Stat(want); err != nil || info.Size() != 0 {
		t.Errorf("launcher should exist and be empty, got %v, %v", info, err)
	}
	if _, err := os.Stat(filepath.Join(hidden, "sky.scummvm")); !os.IsNotExist(err) {
		t.Error("the archive's stub should be gone from the hidden folder")
	}
	if _, err := os.Stat(filepath.Join(hidden, "SKY.DNR")); err != nil {
		t.Error("game data must be left alone")
	}
}

func TestScummVMLauncher_NoStub(t *testing.T) {
	romDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(romDir, "_Sky"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := ScummVMLauncher(filepath.Join(romDir, "Sky"), romDir); ok || err != nil {
		t.Errorf("got ok=%v err=%v, want no launcher", ok, err)
	}
	if _, err := os.Stat(filepath.Join(romDir, "Sky.scummvm")); !os.IsNotExist(err) {
		t.Error("no launcher should be written when the archive has no stub")
	}
}
