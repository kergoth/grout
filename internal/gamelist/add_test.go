package gamelist

import (
	"os"
	"path/filepath"
	"testing"

	"grout/romm"
)

func TestAddRomGamesUsesExplicitGamelistPath(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "ES-DE", "gamelists", "psx", "gamelist.xml")
	entry := RomGameEntry{
		Game:         &romm.Rom{Name: "Ridge Racer"},
		GamePath:     "/run/media/deck/SD/Emulation/roms/psx/Ridge Racer.chd",
		RomDirectory: "/run/media/deck/SD/Emulation/roms/psx",
		GamelistPath: target,
		Platform:     &romm.Platform{FSSlug: "psx"},
	}
	if err := AddRomGamesToGamelist([]RomGameEntry{entry}, GameListFileName); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatal(err)
	}
}

func TestAddRomGamesWritesFolderLink(t *testing.T) {
	target := filepath.Join(t.TempDir(), "gamelist.xml")
	game := &romm.Rom{Name: "Beneath a Steel Sky"}
	if err := AddRomGamesToGamelist([]RomGameEntry{{
		Game:         game,
		GamePath:     "/roms/scummvm/Beneath a Steel Sky/launcher.scummvm",
		GamelistPath: target,
	}}, GameListFileName); err != nil {
		t.Fatal(err)
	}
	entry := RomGameEntry{
		Game:         game,
		GamePath:     "/roms/scummvm/Beneath a Steel Sky",
		FolderLink:   "sky.scummvm",
		GamelistPath: target,
	}
	if err := AddRomGamesToGamelist([]RomGameEntry{entry}, GameListFileName); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	gl := New()
	if err := gl.Parse(data); err != nil {
		t.Fatal(err)
	}
	folder := gl.root().SelectElement(FolderElement)
	if folder == nil {
		t.Fatal("expected folder entry")
	}
	if gl.GetGameElementByName(game.Name) != nil {
		t.Fatal("expected obsolete game entry to be removed")
	}
	if got := folder.FindElement(PathElement).Text(); got != entry.GamePath {
		t.Fatalf("folder path = %q, want %q", got, entry.GamePath)
	}
	if got := folder.FindElement(FolderLinkElement).Text(); got != entry.FolderLink {
		t.Fatalf("folderlink = %q, want %q", got, entry.FolderLink)
	}
}
