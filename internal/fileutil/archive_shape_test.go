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

func TestExtractArchiveFlat(t *testing.T) {
	t.Run("flat single file", func(t *testing.T) {
		src := writeZip(t, "game.sfc")
		shape, _ := AnalyzeArchive(src)
		dest := t.TempDir()
		out, err := ExtractArchiveFlat(src, shape, dest, nil)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Base(out) != "game.sfc" || filepath.Dir(out) != dest {
			t.Errorf("got %q, want %s/game.sfc", out, dest)
		}
	})

	t.Run("strips toplevel folder", func(t *testing.T) {
		src := writeZip(t, "Game Name/game.sfc")
		shape, _ := AnalyzeArchive(src)
		dest := t.TempDir()
		out, err := ExtractArchiveFlat(src, shape, dest, nil)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Base(out) != "game.sfc" || filepath.Dir(out) != dest {
			t.Errorf("got %q, want %s/game.sfc", out, dest)
		}
	})

	t.Run("rejects multi-file", func(t *testing.T) {
		src := writeZip(t, "a.sfc", "b.sfc")
		shape, _ := AnalyzeArchive(src)
		if _, err := ExtractArchiveFlat(src, shape, t.TempDir(), nil); err == nil {
			t.Error("expected error for multi-file archive")
		}
	})
}

func TestExtractArchiveToFolder(t *testing.T) {
	t.Run("strips single root", func(t *testing.T) {
		src := writeZip(t, "Game/monkey.000", "Game/monkey.001")
		shape, _ := AnalyzeArchive(src)
		dest := t.TempDir()
		files, err := ExtractArchiveToFolder(src, shape, dest, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 2 {
			t.Fatalf("got %d files, want 2", len(files))
		}
		for _, f := range files {
			if filepath.Dir(f) != dest {
				t.Errorf("file %q not directly in dest", f)
			}
		}
	})

	t.Run("preserves multi-root structure", func(t *testing.T) {
		src := writeZip(t, "disc1/a.bin", "disc2/b.bin")
		shape, _ := AnalyzeArchive(src)
		dest := t.TempDir()
		files, err := ExtractArchiveToFolder(src, shape, dest, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 2 {
			t.Fatalf("got %d files, want 2", len(files))
		}
		if !fileExists(filepath.Join(dest, "disc1", "a.bin")) || !fileExists(filepath.Join(dest, "disc2", "b.bin")) {
			t.Errorf("expected disc1/disc2 structure preserved, got %v", files)
		}
	})

	t.Run("skips junk", func(t *testing.T) {
		src := writeZip(t, "Game/monkey.000", "__MACOSX/Game/._monkey.000", ".DS_Store")
		shape, _ := AnalyzeArchive(src)
		dest := t.TempDir()
		files, err := ExtractArchiveToFolder(src, shape, dest, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 1 {
			t.Errorf("got %d files, want 1 (junk skipped): %v", len(files), files)
		}
	})

	t.Run("rejects path traversal", func(t *testing.T) {
		src := writeZip(t, "../evil.sfc")
		shape, _ := AnalyzeArchive(src)
		if _, err := ExtractArchiveToFolder(src, shape, t.TempDir(), nil); err == nil {
			t.Error("expected error for traversal entry")
		}
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
