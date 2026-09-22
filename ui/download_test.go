package ui

import (
	"os"
	"testing"

	"grout/internal"
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
