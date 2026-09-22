// Package scummvm locates and validates ScummVM launcher stubs shipped
// inside extracted ScummVM game archives. A stub is a file with a
// ".scummvm" extension whose contents are the ScummVM short game ID (e.g.
// "sky"); the file name itself is arbitrary.
package scummvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindStub locates the single ".scummvm" launcher stub directly inside dir
// and returns its path and validated short ID. ok is false, with an empty
// path and ID, when dir contains zero or more than one top-level stub, or
// when the stub's contents are not a usable short ID. FindStub never
// modifies the filesystem.
func FindStub(dir string) (stubPath, shortID string, ok bool, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", "", false, fmt.Errorf("read ScummVM game directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".scummvm") {
			continue
		}
		if stubPath != "" {
			return "", "", false, nil
		}
		stubPath = filepath.Join(dir, entry.Name())
	}
	if stubPath == "" {
		return "", "", false, nil
	}

	contents, err := os.ReadFile(stubPath)
	if err != nil {
		return "", "", false, fmt.Errorf("read ScummVM launcher stub: %w", err)
	}
	shortID = strings.TrimSpace(string(contents))
	if shortID == "" || strings.ContainsAny(shortID, "/\\\r\n") || filepath.Base(shortID) != shortID || shortID == "." {
		return "", "", false, nil
	}

	return stubPath, shortID, true, nil
}
