package muos

import (
	"fmt"
	"os"
	"path/filepath"

	"grout/scummvm"
)

// ScummVMLauncher sets up an extracted ScummVM game the way muOS launches it:
// a game folder hidden behind an underscore, with an empty "<name>.scummvm"
// beside it in the rom directory. muOS finds the game by autodetecting the
// folder and rewrites the stub on first launch, so any ID the archive's stub
// held is dropped.
//
// gameDir is where the game would have been extracted. OrganizeMultiFileRom
// has already renamed it to its underscore form by now.
func ScummVMLauncher(gameDir, romDir string) (string, bool, error) {
	name := filepath.Base(gameDir)

	stub, ok, err := scummvm.FindStub(filepath.Join(romDir, "_"+name))
	if err != nil || !ok {
		return "", false, err
	}

	launcher := filepath.Join(romDir, name+".scummvm")
	if err := os.WriteFile(launcher, nil, 0644); err != nil {
		return "", false, fmt.Errorf("writing ScummVM launcher: %w", err)
	}
	if err := os.Remove(stub); err != nil {
		return "", false, fmt.Errorf("removing the archive's ScummVM stub: %w", err)
	}
	return launcher, true, nil
}
