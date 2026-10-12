package archive

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"go.uber.org/atomic"
)

// Shape describes how an archive is packed, independent of how the uploader
// chose to arrange it. Layout decisions are made from the shape, never from
// the archive's filename.
type Shape struct {
	// UsableFiles counts the entries that are real content, leaving out the
	// metadata macOS and Windows tools tuck into archives.
	UsableFiles int
	// SingleRoot is the one top-level folder every usable entry sits under,
	// which extraction strips so the game folder is not nested twice. It is
	// left empty for a folder that is part of the content's own layout.
	SingleRoot string
	// HasStructure is set when any usable entry sits below the archive root.
	HasStructure bool
}

// FolderShaped reports whether the content needs a folder of its own, as
// opposed to being one file that can sit directly among the other roms.
func (s Shape) FolderShaped() bool {
	return s.UsableFiles > 1 || (s.UsableFiles == 1 && s.HasStructure && s.SingleRoot == "")
}

// layoutFolders are top-level folders that belong to the content's own layout,
// so stripping one would damage the game: console disc dumps whose emulators
// look for them by name, and the folders updates and DLC are filed under.
var layoutFolders = map[string]bool{
	"ps3_game": true, "psp_game": true, "usrdir": true,
	"update": true, "updates": true, "dlc": true, "dlcs": true, "extcontent": true,
}

func isJunkEntry(name string) bool {
	norm := strings.ReplaceAll(name, "\\", "/")
	base := norm
	if i := strings.LastIndex(norm, "/"); i >= 0 {
		base = norm[i+1:]
	}
	return strings.HasPrefix(norm, "__MACOSX/") ||
		base == ".DS_Store" ||
		strings.HasPrefix(base, "._")
}

func normalizeEntryName(name string) string {
	return strings.TrimSuffix(strings.ReplaceAll(name, "\\", "/"), "/")
}

func isSevenZip(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".7z")
}

// Analyze reads an archive's entry names and works out its shape.
func Analyze(path string) (Shape, error) {
	var names []string
	var err error
	if isSevenZip(path) {
		names, err = SevenZipFileNames(path)
	} else {
		names, err = ZipFileNames(path)
	}
	if err != nil {
		return Shape{}, err
	}

	shape := Shape{}
	root := ""
	rootSet := false

	for _, raw := range names {
		if isJunkEntry(raw) {
			continue
		}
		norm := normalizeEntryName(raw)
		if norm == "" {
			continue
		}

		shape.UsableFiles++
		first, rest, _ := strings.Cut(norm, "/")
		if rest != "" {
			shape.HasStructure = true
		}

		if !rootSet {
			root = first
			rootSet = true
		} else if root != first {
			root = ""
		}
	}

	// A shared first component only counts as a folder to strip when
	// something actually lives beneath it.
	if rootSet && root != "" && shape.HasStructure && !layoutFolders[strings.ToLower(root)] {
		shape.SingleRoot = root
	}
	return shape, nil
}

// ExtractFlat unpacks the archive's one usable file straight into destDir,
// whatever folders it sat in, and returns where it landed.
func ExtractFlat(archivePath string, shape Shape, destDir string, progress *atomic.Float64) (string, error) {
	if shape.UsableFiles != 1 {
		return "", errors.New("ExtractFlat needs an archive holding exactly one file")
	}
	paths, err := extractShaped(archivePath, shape, destDir, progress, true)
	if err != nil {
		return "", err
	}
	if len(paths) != 1 {
		return "", fmt.Errorf("expected 1 extracted file, got %d", len(paths))
	}
	return paths[0], nil
}

// ExtractToFolder unpacks every usable entry into destDir, dropping the shared
// top-level folder when there is one, and returns the paths it wrote.
func ExtractToFolder(archivePath string, shape Shape, destDir string, progress *atomic.Float64) ([]string, error) {
	return extractShaped(archivePath, shape, destDir, progress, false)
}

// shapedEntry is one archive entry with where it lands relative to the
// destination, and how to copy it there.
type shapedEntry struct {
	name    string
	rel     string
	size    uint64
	extract func(destPath string, buffer []byte, total uint64, done *uint64, progress *atomic.Float64) error
}

func extractShaped(archivePath string, shape Shape, destDir string, progress *atomic.Float64, flatten bool) ([]string, error) {
	var entries []shapedEntry
	var closeReader func() error

	if isSevenZip(archivePath) {
		reader, err := sevenzip.OpenReader(archivePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open 7z file: %w", err)
		}
		closeReader = reader.Close
		for _, f := range reader.File {
			if f.FileInfo().IsDir() {
				continue
			}
			entries = append(entries, shapedEntry{
				name: f.Name,
				size: f.UncompressedSize,
				extract: func(destPath string, buffer []byte, total uint64, done *uint64, progress *atomic.Float64) error {
					return extract7zFile(f, destPath, buffer, total, done, progress)
				},
			})
		}
	} else {
		reader, err := zip.OpenReader(archivePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open zip file: %w", err)
		}
		closeReader = reader.Close
		for _, f := range reader.File {
			if f.FileInfo().IsDir() {
				continue
			}
			entries = append(entries, shapedEntry{
				name: f.Name,
				size: f.UncompressedSize64,
				extract: func(destPath string, buffer []byte, total uint64, done *uint64, progress *atomic.Float64) error {
					return extractFile(f, destPath, buffer, total, done, progress)
				},
			})
		}
	}
	defer closeReader()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	var usable []shapedEntry
	var totalBytes uint64
	for _, e := range entries {
		if isJunkEntry(e.name) {
			continue
		}
		norm := normalizeEntryName(e.name)
		// Checked before the root is stripped, so a ".." component cannot pass
		// itself off as a shared root and vanish.
		if err := validateEntryPath(norm); err != nil {
			return nil, err
		}
		rel := norm
		if shape.SingleRoot != "" {
			rel = strings.TrimPrefix(norm, shape.SingleRoot+"/")
		}
		if flatten {
			rel = filepath.Base(rel)
		}
		e.rel = filepath.FromSlash(rel)
		usable = append(usable, e)
		totalBytes += e.size
	}

	buffer := make([]byte, DefaultBufferSize)
	var extractedBytes uint64
	createdDirs := map[string]bool{destDir: true}
	written := make([]string, 0, len(usable))

	for _, e := range usable {
		destPath := filepath.Join(destDir, e.rel)
		if parent := filepath.Dir(destPath); !createdDirs[parent] {
			if err := os.MkdirAll(parent, 0755); err != nil {
				return nil, fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
			}
			createdDirs[parent] = true
		}
		if err := e.extract(destPath, buffer, totalBytes, &extractedBytes, progress); err != nil {
			return nil, fmt.Errorf("failed to extract file %s: %w", e.name, err)
		}
		written = append(written, destPath)
	}
	return written, nil
}

// validateEntryPath rejects entries that would land outside the destination.
func validateEntryPath(norm string) error {
	for _, part := range strings.Split(norm, "/") {
		if part == ".." || part == "" {
			return fmt.Errorf("archive entry escapes destination: %q", norm)
		}
	}
	return nil
}
