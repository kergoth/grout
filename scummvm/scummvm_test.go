package scummvm

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, dir, rel, content string) string {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFindStub(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  string // relative to the game folder; "" means none found
	}{
		{"at the top", []string{"sky.scummvm", "SKY.DNR"}, "sky.scummvm"},
		{"wrapped one level down", []string{"Game/sky.scummvm", "Game/SKY.DNR"}, "Game/sky.scummvm"},
		{"extension case does not matter", []string{"SKY.SCUMMVM"}, "SKY.SCUMMVM"},
		{"shallowest level wins over a deeper one", []string{"sky.scummvm", "extras/other.scummvm"}, "sky.scummvm"},
		{"two at the same level is ambiguous", []string{"a.scummvm", "b.scummvm"}, ""},
		{"two at the shallowest level is ambiguous even with a deeper one", []string{"a.scummvm", "b.scummvm", "x/c.scummvm"}, ""},
		{"macOS metadata is skipped", []string{"._sky.scummvm", "sky.scummvm"}, "sky.scummvm"},
		{"no stub", []string{"SKY.DNR"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				write(t, dir, f, "sky")
			}

			got, ok, err := FindStub(dir)
			if err != nil {
				t.Fatalf("FindStub: %v", err)
			}
			if tt.want == "" {
				if ok {
					t.Fatalf("expected no stub, got %q", got)
				}
				return
			}
			if !ok || got != filepath.Join(dir, tt.want) {
				t.Errorf("FindStub = %q, %v; want %q", got, ok, filepath.Join(dir, tt.want))
			}
		})
	}
}

func TestFindStub_MissingDirectory(t *testing.T) {
	if _, _, err := FindStub(filepath.Join(t.TempDir(), "absent")); err == nil {
		t.Error("expected an error for a folder that does not exist")
	}
}
