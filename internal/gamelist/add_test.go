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
