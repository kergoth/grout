package scummvm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindStub(t *testing.T) {
	tests := []struct {
		name         string
		stubs        map[string]string
		wantStubName string
		wantShortID  string
		wantOK       bool
	}{
		{
			name:         "finds a single stub",
			stubs:        map[string]string{"launcher.scummvm": "sky\n"},
			wantStubName: "launcher.scummvm",
			wantShortID:  "sky",
			wantOK:       true,
		},
		{
			name:         "accepts an already-canonical stub name",
			stubs:        map[string]string{"sky.scummvm": "sky"},
			wantStubName: "sky.scummvm",
			wantShortID:  "sky",
			wantOK:       true,
		},
		{
			name:   "does not derive an ID when stub is missing",
			stubs:  nil,
			wantOK: false,
		},
		{
			name:   "does not use an empty stub",
			stubs:  map[string]string{"launcher.scummvm": " \n"},
			wantOK: false,
		},
		{
			name:   "does not choose among multiple stubs",
			stubs:  map[string]string{"one.scummvm": "sky", "two.scummvm": "monkey"},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, contents := range tt.stubs {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			stubPath, shortID, ok, err := FindStub(dir)
			if err != nil {
				t.Fatal(err)
			}
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				if stubPath != "" || shortID != "" {
					t.Fatalf("stubPath = %q, shortID = %q, want both empty", stubPath, shortID)
				}
				return
			}
			wantStubPath := filepath.Join(dir, tt.wantStubName)
			if stubPath != wantStubPath {
				t.Fatalf("stubPath = %q, want %q", stubPath, wantStubPath)
			}
			if shortID != tt.wantShortID {
				t.Fatalf("shortID = %q, want %q", shortID, tt.wantShortID)
			}
		})
	}
}
