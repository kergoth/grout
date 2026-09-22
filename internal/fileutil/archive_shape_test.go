package fileutil

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func writeZip(t *testing.T, entries ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := zip.NewWriter(f)
	for _, e := range entries {
		if e[len(e)-1] == '/' {
			_, err = w.Create(e)
		} else {
			fw, cerr := w.Create(e)
			if cerr != nil {
				t.Fatal(cerr)
			}
			fw.Write([]byte("x"))
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAnalyzeArchive(t *testing.T) {
	cases := []struct {
		name          string
		entries       []string
		wantFiles     int
		wantRoot      string
		wantStructure bool
	}{
		{"single flat file", []string{"game.sfc"}, 1, "", false},
		{"two flat files", []string{"game.sfc", "game.srm"}, 2, "", false},
		{"single file in toplevel folder", []string{"Game Name/game.sfc"}, 1, "Game Name", true},
		{"multi file in toplevel folder", []string{"Game/monkey.000", "Game/monkey.001"}, 2, "Game", true},
		{"toplevel folder with subdir", []string{"Game/data/monkey.000", "Game/monkey.001"}, 2, "Game", true},
		{"multi root not flattened", []string{"disc1/a.bin", "disc2/b.bin"}, 2, "", true},
		{"mixed root file and folder", []string{"readme.txt", "Game/monkey.000"}, 2, "", true},
		{"macos junk ignored", []string{"Game/monkey.000", "__MACOSX/Game/._monkey.000", ".DS_Store"}, 1, "Game", true},
		{"windows separators", []string{"Game\\monkey.000", "Game\\monkey.001"}, 2, "Game", true},
		{"all junk", []string{"__MACOSX/._x", ".DS_Store"}, 0, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			shape, err := AnalyzeArchive(writeZip(t, c.entries...))
			if err != nil {
				t.Fatal(err)
			}
			if shape.UsableFiles != c.wantFiles {
				t.Errorf("UsableFiles = %d, want %d", shape.UsableFiles, c.wantFiles)
			}
			if shape.SingleRoot != c.wantRoot {
				t.Errorf("SingleRoot = %q, want %q", shape.SingleRoot, c.wantRoot)
			}
			if shape.HasStructure != c.wantStructure {
				t.Errorf("HasStructure = %v, want %v", shape.HasStructure, c.wantStructure)
			}
		})
	}
}
