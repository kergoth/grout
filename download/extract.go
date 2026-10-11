package download

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/atomic"

	"grout/archive"
	"grout/cfw"
	"grout/files"
	"grout/romm"
)

// IsArchive reports whether a rom arrived as an archive worth unpacking.
func IsArchive(fileName string) bool {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".zip", ".7z":
		return true
	}
	return false
}

// Treatment is what to do with a downloaded archive.
type Treatment int

const (
	// KeepArchive leaves the archive as downloaded.
	KeepArchive Treatment = iota
	// ExtractToRomDirectory unpacks the one file inside next to the other roms.
	ExtractToRomDirectory
	// ExtractToGameFolder unpacks into a folder named for the game.
	ExtractToGameFolder
)

// ArchiveTreatment decides what happens to a downloaded archive from what is
// inside it, not from how the uploader named or arranged it.
//
// A platform that plays archives keeps them. Past that, content made of
// several files needs a folder of its own and is always unpacked, since the
// frontend cannot launch it any other way. A lone rom in a 7z is unpacked too,
// because the emulators that read a zip directly do not read 7z, and so is one
// on a platform whose emulators never read a zip. Only a lone rom in a zip on
// the rest is a matter of preference.
func ArchiveTreatment(fsSlug, archivePath string, shape archive.Shape, preferExtracted bool) Treatment {
	switch {
	case KeepsArchive(fsSlug), shape.UsableFiles == 0:
		return KeepArchive
	case shape.FolderShaped():
		return ExtractToGameFolder
	case preferExtracted, AlwaysExtracts(fsSlug), strings.EqualFold(filepath.Ext(archivePath), ".7z"):
		return ExtractToRomDirectory
	}
	return KeepArchive
}

// ExtractMultiFile unpacks a game the server delivered as one zip holding
// several files, and returns the path the frontend should point at.
//
// The archive lands in a temp directory rather than the rom directory, so a
// failure here leaves nothing half-unpacked where the frontend would find it.
func ExtractMultiFile(game romm.Rom, romDirectory string, progress *atomic.Float64) (string, error) {
	archivePath := MultiFileArchivePath(game)

	shape, err := archive.Analyze(archivePath)
	if err != nil {
		os.Remove(archivePath)
		return "", fmt.Errorf("reading %s: %w", game.Name, err)
	}

	path, err := extractToGameFolder(archivePath, romDirectory, game.FsNameNoExt, shape, progress)
	if err != nil {
		os.Remove(archivePath)
		return "", fmt.Errorf("unpacking %s: %w", game.Name, err)
	}
	return path, nil
}

// ExtractToFolder unpacks a downloaded archive into a folder named for the game
// inside its rom directory, and returns the path the frontend should point at.
// A failure leaves the archive alone.
func ExtractToFolder(archivePath, romDirectory, baseName string, shape archive.Shape, progress *atomic.Float64) (string, error) {
	path, err := extractToGameFolder(archivePath, romDirectory, baseName, shape, progress)
	if err != nil {
		return "", fmt.Errorf("unpacking %s: %w", filepath.Base(archivePath), err)
	}
	return path, nil
}

func extractToGameFolder(archivePath, romDirectory, baseName string, shape archive.Shape, progress *atomic.Float64) (string, error) {
	extractDir := filepath.Join(romDirectory, baseName)

	if _, err := archive.ExtractToFolder(archivePath, shape, extractDir, progress); err != nil {
		return "", err
	}

	if err := cfw.OrganizeExtractedRom(extractDir, romDirectory, baseName); err != nil {
		os.RemoveAll(extractDir)
		return "", fmt.Errorf("arranging for this firmware: %w", err)
	}

	if err := os.Remove(archivePath); err != nil {
		slog.Default().Warn("Could not remove the downloaded archive", "path", archivePath, "error", err)
	}

	return multiFileGamePath(romDirectory, extractDir, baseName), nil
}

// ExtractArchive unpacks the one rom inside a .zip or .7z straight into its
// rom directory, and returns the path the frontend should point at.
//
// The archive is removed once it is unpacked. A failure leaves it alone, so
// the game is still playable from the archive if the frontend can read one.
func ExtractArchive(archivePath, romDirectory string, shape archive.Shape, progress *atomic.Float64) (string, error) {
	path, err := archive.ExtractFlat(archivePath, shape, romDirectory, progress)
	if err != nil {
		return "", fmt.Errorf("unpacking %s: %w", filepath.Base(archivePath), err)
	}

	if err := os.Remove(archivePath); err != nil {
		slog.Default().Warn("Could not remove the archive after unpacking", "path", archivePath, "error", err)
	}
	return path, nil
}

// multiFileGamePath finds what to launch after a multi-file game is unpacked.
//
// A firmware that rearranges the files may have lifted a playlist up into the
// rom directory, so that is looked for first. Failing everything the extract
// directory itself is named, which at least points somewhere real.
func multiFileGamePath(romDirectory, extractDir, baseName string) string {
	logger := slog.Default()

	if playlist := filepath.Join(romDirectory, baseName+".m3u"); files.FileExists(playlist) {
		return playlist
	}

	entries, err := os.ReadDir(extractDir)
	if err != nil {
		logger.Debug("Cannot read the extract directory, naming it instead", "path", extractDir, "error", err)
		return extractDir
	}

	firstFile := ""
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".m3u") {
			return filepath.Join(extractDir, entry.Name())
		}
		if firstFile == "" {
			firstFile = filepath.Join(extractDir, entry.Name())
		}
	}

	if firstFile == "" {
		logger.Debug("Nothing was unpacked, naming the directory instead", "path", extractDir)
		return extractDir
	}
	return firstFile
}
