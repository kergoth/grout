package ui

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"grout/internal"
	"grout/internal/fileutil"
	"grout/romm"
)

func TestMain(m *testing.M) {
	// cfw.GetCFW() (called from buildDownloads) requires the CFW env var to
	// be set to a recognised value, otherwise it terminates the process.
	if os.Getenv("CFW") == "" {
		os.Setenv("CFW", "ROCKNIX")
	}
	os.Exit(m.Run())
}

// TestBuildDownloads_EmptyFiles is a regression test for issue #223:
// a cached ROM with HasMultipleFiles == false and an empty Files slice
// must not panic. The entry should be skipped instead.
func TestBuildDownloads_EmptyFiles(t *testing.T) {
	s := NewDownloadScreen()
	config := internal.Config{}
	host := romm.Host{RootURI: "http://example.invalid"}
	platform := romm.Platform{ID: 1, FSSlug: "nds", Name: "Nintendo DS"}

	games := []romm.Rom{
		{
			ID:               60,
			Name:             "Professor Layton and the Curious Village",
			FsName:           "Professor Layton and the Curious Village (USA).nds",
			FsNameNoExt:      "Professor Layton and the Curious Village (USA)",
			HasMultipleFiles: false,
			Files:            nil,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("buildDownloads panicked on empty Files slice: %v", r)
		}
	}()

	downloads, artDownloads, gamelistEntries := s.buildDownloads(config, host, platform, games, 0)

	if len(downloads) != 0 {
		t.Errorf("expected 0 downloads when Files is empty, got %d", len(downloads))
	}
	if len(artDownloads) != 0 {
		t.Errorf("expected 0 art downloads when Files is empty, got %d", len(artDownloads))
	}
	if len(gamelistEntries) != 0 {
		t.Errorf("expected 0 gamelist entries when Files is empty, got %d", len(gamelistEntries))
	}
}

// TestBuildDownloads_SingleFile_HappyPath sanity-checks that a normal
// single-file ROM still produces a download URL after the empty-Files guard.
func TestBuildDownloads_SingleFile_HappyPath(t *testing.T) {
	s := NewDownloadScreen()
	config := internal.Config{}
	host := romm.Host{RootURI: "http://example.invalid"}
	platform := romm.Platform{ID: 1, FSSlug: "nds", Name: "Nintendo DS"}

	games := []romm.Rom{
		{
			ID:               42,
			Name:             "Test Game",
			FsName:           "test.nds",
			FsNameNoExt:      "test",
			HasMultipleFiles: false,
			Files: []romm.RomFile{
				{ID: 100, FileName: "test.nds"},
			},
		},
	}

	downloads, _, gamelistEntries := s.buildDownloads(config, host, platform, games, 0)

	if len(downloads) != 1 {
		t.Fatalf("expected 1 download, got %d", len(downloads))
	}
	if len(gamelistEntries) != 1 {
		t.Fatalf("expected 1 gamelist entry, got %d", len(gamelistEntries))
	}
	if downloads[0].URL == "" {
		t.Error("expected non-empty download URL")
	}
}

func TestShouldExtractSingleFileDownload(t *testing.T) {
	cases := []struct {
		name           string
		unzipDownloads bool
		fsSlug         string
		want           bool
	}{
		{"uncompress on, regular platform", true, "snes", true},
		{"uncompress on, dos", true, "dos", false},
		{"uncompress on, arcade", true, "arcade", false},
		{"uncompress on, neogeo alias", true, "neogeocd", false},
		{"uncompress on, scummvm still extracts", true, "scummvm", true},
		{"uncompress off, regular platform", false, "snes", false},
		{"uncompress off, dos", false, "dos", false},
		{"uncompress on, empty slug extracts", true, "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := shouldExtractSingleFileDownload(c.unzipDownloads, c.fsSlug); got != c.want {
				t.Errorf("shouldExtractSingleFileDownload(%v, %q) = %v, want %v",
					c.unzipDownloads, c.fsSlug, got, c.want)
			}
		})
	}
}

func TestRouteSingleFileExtraction(t *testing.T) {
	cases := []struct {
		name       string
		shape      fileutil.ArchiveShape
		wantFolder bool
	}{
		{"single flat file extracts as file", fileutil.ArchiveShape{UsableFiles: 1}, false},
		{"single file in folder extracts as file", fileutil.ArchiveShape{UsableFiles: 1, SingleRoot: "Game", HasStructure: true}, false},
		{"multi file extracts to folder", fileutil.ArchiveShape{UsableFiles: 3, SingleRoot: "Game", HasStructure: true}, true},
		{"multi root extracts to folder", fileutil.ArchiveShape{UsableFiles: 2, HasStructure: true}, true},
		{"two flat files extract to folder", fileutil.ArchiveShape{UsableFiles: 2}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := singleFileExtractsToFolder(c.shape); got != c.wantFolder {
				t.Errorf("singleFileExtractsToFolder(%+v) = %v, want %v", c.shape, got, c.wantFolder)
			}
		})
	}
}

func TestMultiFileExtractionStripsToplevel(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "server.zip")

	// Server-built zip carrying a toplevel folder
	f, _ := os.Create(zipPath)
	w := zip.NewWriter(f)
	for _, name := range []string{"My Game/disc1.chd", "My Game/disc2.chd"} {
		fw, _ := w.Create(name)
		fw.Write([]byte("x"))
	}
	w.Close()
	f.Close()

	shape, err := fileutil.AnalyzeArchive(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	extractDir := filepath.Join(dir, "My Game")
	files, err := fileutil.ExtractArchiveToFolder(zipPath, shape, extractDir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("got %d files, want 2", len(files))
	}
	for _, want := range []string{"disc1.chd", "disc2.chd"} {
		if !fileutil.FileExists(filepath.Join(extractDir, want)) {
			t.Errorf("expected %s directly in extract dir", want)
		}
	}
	if fileutil.FileExists(filepath.Join(extractDir, "My Game")) {
		t.Error("double-nested toplevel folder was not stripped")
	}
}

func TestResolvePostExtractionGamePathForEmuDeckScummVM(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Game Name")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "launcher.scummvm"), []byte("sky\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CFW", "EMUDECK")
	got := resolvePostExtractionGamePath("scummvm", root, extractDir, "Game Name")
	want := filepath.Join(root, "sky.scummvm", "sky.scummvm")
	if got != want {
		t.Fatalf("game path = %q, want %q", got, want)
	}
}

func TestResolvePostExtractionGamePathLeavesOtherCFWsUnchanged(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Game Name")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stubPath := filepath.Join(extractDir, "launcher.scummvm")
	if err := os.WriteFile(stubPath, []byte("sky\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CFW", "ROCKNIX")
	got := resolvePostExtractionGamePath("scummvm", root, extractDir, "Game Name")
	if got != stubPath {
		t.Fatalf("game path = %q, want %q", got, stubPath)
	}
	if _, err := os.Stat(extractDir); err != nil {
		t.Fatalf("original directory changed: %v", err)
	}
}
