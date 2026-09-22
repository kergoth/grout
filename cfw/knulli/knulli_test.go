package knulli

import (
	"os"
	"path/filepath"
	"testing"
)

// Stub discovery and short-ID validation are covered by
// internal/scummvm.FindStub's own tests. Knulli performs no renaming, so
// these cases only confirm the stub's existing path is returned unchanged.

func TestPrepareScummVMLauncherReturnsStubPath(t *testing.T) {
	extractDir := t.TempDir()
	stubPath := filepath.Join(extractDir, "any-name.scummvm")
	if err := os.WriteFile(stubPath, []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}

	gotPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ready {
		t.Fatal("ready = false, want true")
	}
	if gotPath != stubPath {
		t.Fatalf("stub path = %q, want %q", gotPath, stubPath)
	}
	if _, err := os.Stat(stubPath); err != nil {
		t.Fatalf("stub file moved or removed: %v", err)
	}
}

func TestPrepareScummVMLauncherNotReadyWhenNoStub(t *testing.T) {
	extractDir := t.TempDir()

	gotPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if ready {
		t.Fatal("ready = true, want false")
	}
	if gotPath != "" {
		t.Fatalf("stub path = %q, want empty", gotPath)
	}
}

func TestPrepareScummVMLauncherNotReadyWithMultipleStubs(t *testing.T) {
	extractDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(extractDir, "one.scummvm"), []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "two.scummvm"), []byte("monkey"), 0o644); err != nil {
		t.Fatal(err)
	}

	gotPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if ready {
		t.Fatal("ready = true, want false")
	}
	if gotPath != "" {
		t.Fatalf("stub path = %q, want empty", gotPath)
	}
}
