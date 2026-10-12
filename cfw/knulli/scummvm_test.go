package knulli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScummVMLauncher(t *testing.T) {
	dir := t.TempDir()
	stub := filepath.Join(dir, "Game", "sky.scummvm")
	if err := os.MkdirAll(filepath.Dir(stub), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stub, []byte("sky"), 0644); err != nil {
		t.Fatal(err)
	}

	got, ok, err := ScummVMLauncher(dir, "")
	if err != nil || !ok || got != stub {
		t.Errorf("ScummVMLauncher = %q, %v, %v; want %q, true, nil", got, ok, err, stub)
	}

	if _, ok, _ := ScummVMLauncher(t.TempDir(), ""); ok {
		t.Error("expected no launcher for a folder without a stub")
	}
}
