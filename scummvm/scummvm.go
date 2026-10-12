// Package scummvm finds the launcher stubs that ship inside ScummVM game
// archives. A stub is a file with a ".scummvm" extension; for most frontends
// its content is the ScummVM game ID and its name is arbitrary.
package scummvm

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const stubExt = ".scummvm"

// FindStub locates the launcher stub inside an extracted game folder.
//
// The stub is usually at the top of the folder, but an archive that wraps the
// game in a further folder puts it deeper, so the shallowest level holding any
// stub wins. ok is false when that level holds none or more than one, because
// guessing between two stubs would launch the wrong game. Names starting with a
// dot are archive metadata and are skipped. FindStub never changes the
// filesystem.
func FindStub(dir string) (path string, ok bool, err error) {
	var found []string
	foundDepth := -1

	err = filepath.WalkDir(dir, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") || !strings.EqualFold(filepath.Ext(d.Name()), stubExt) {
			return nil
		}

		rel, relErr := filepath.Rel(dir, p)
		if relErr != nil {
			return relErr
		}
		depth := strings.Count(rel, string(filepath.Separator))

		switch {
		case foundDepth == -1 || depth < foundDepth:
			found, foundDepth = []string{p}, depth
		case depth == foundDepth:
			found = append(found, p)
		}
		return nil
	})
	if err != nil {
		return "", false, fmt.Errorf("searching %s for a ScummVM stub: %w", dir, err)
	}

	if len(found) != 1 {
		return "", false, nil
	}
	return found[0], true, nil
}
