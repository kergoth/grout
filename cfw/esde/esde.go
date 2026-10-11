// Package esde holds what the firmwares that run ES-DE as their frontend
// (RetroDECK, EmuDeck) have in common.
package esde

import (
	"embed"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"grout/scummvm"
	"grout/tables"
)

//go:embed data/*.json
var embeddedFiles embed.FS

// Platforms maps a RomM platform slug to the ES-DE system folders that hold
// its roms, the first being where downloads go. It is the table for a stock
// ES-DE install; a firmware whose folders differ layers its own on top.
var Platforms = tables.MustLoad[string, []string](embeddedFiles, "data/platforms.json")

// With returns Platforms with overrides applied: an override replaces the
// whole entry for its slug.
func With(overrides map[string][]string) map[string][]string {
	merged := maps.Clone(Platforms)
	maps.Copy(merged, overrides)
	return merged
}

// ScummVMLauncher arranges an extracted ScummVM game the way ES-DE launches it:
// a folder named "<name>.scummvm" holding a stub of the same name. ES-DE treats
// such a folder as one game file, so the game shows in its list without being
// entered, and it hands the stub's content to ScummVM as the game to run.
//
// gameDir is the extracted folder and must sit directly in romDir. ok is false,
// leaving the folder untouched, when there is no single usable stub at the top
// of it: moving a stub away from the game data would stop ScummVM finding it.
func ScummVMLauncher(gameDir, romDir string) (string, bool, error) {
	stub, ok, err := scummvm.FindStub(gameDir)
	if err != nil || !ok || filepath.Dir(stub) != gameDir {
		return "", false, err
	}
	if _, ok, err := scummvm.ReadID(stub); err != nil || !ok {
		return "", false, err
	}

	name := filepath.Base(gameDir)
	if !strings.HasSuffix(strings.ToLower(name), scummvm.StubExt) {
		name += scummvm.StubExt
	}
	folder := filepath.Join(romDir, name)

	if folder != gameDir {
		if _, err := os.Lstat(folder); err == nil {
			return "", false, fmt.Errorf("%s already exists", folder)
		}
		if err := os.Rename(gameDir, folder); err != nil {
			return "", false, fmt.Errorf("renaming the game folder: %w", err)
		}
		stub = filepath.Join(folder, filepath.Base(stub))
	}

	launcher := filepath.Join(folder, name)
	if launcher != stub {
		if err := os.Rename(stub, launcher); err != nil {
			if folder != gameDir {
				os.Rename(folder, gameDir)
			}
			return "", false, fmt.Errorf("renaming the ScummVM stub: %w", err)
		}
	}
	return folder, true, nil
}
