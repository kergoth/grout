package fileutil

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

// ArchiveShape describes how an archive is packed, independent of how the
// uploader chose to arrange it. Extraction layout decisions are made from
// this shape, never from the archive's filename or extension.
type ArchiveShape struct {
	UsableFiles  int
	SingleRoot   string
	HasStructure bool
}

// isJunkArchiveEntry reports whether an entry is macOS/Windows metadata
// noise that must not influence shape detection or be extracted.
func isJunkArchiveEntry(name string) bool {
	norm := strings.ReplaceAll(name, "\\", "/")
	base := norm
	if i := strings.LastIndex(norm, "/"); i >= 0 {
		base = norm[i+1:]
	}
	return strings.HasPrefix(norm, "__MACOSX/") ||
		base == ".DS_Store" ||
		strings.HasPrefix(base, "._")
}

// AnalyzeArchive sniffs an archive's entry layout. SingleRoot is set when
// every usable entry shares one first path component (a redundant toplevel
// folder the extractor should strip). HasStructure is set when any usable
// entry lives below the root, meaning the content is folder-shaped.
func AnalyzeArchive(path string) (ArchiveShape, error) {
	names, err := archiveEntryNames(path)
	if err != nil {
		return ArchiveShape{}, err
	}

	shape := ArchiveShape{}
	root := ""
	rootSet := false

	for _, raw := range names {
		if isJunkArchiveEntry(raw) {
			continue
		}
		norm := strings.ReplaceAll(raw, "\\", "/")
		norm = strings.TrimSuffix(norm, "/")
		if norm == "" {
			continue
		}

		shape.UsableFiles++
		first := norm
		rest := ""
		if i := strings.Index(norm, "/"); i >= 0 {
			first = norm[:i]
			rest = norm[i+1:]
		}
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

	// A common root only counts as a strippable toplevel folder when it is
	// actually a folder, i.e. at least one entry lives beneath it.
	if rootSet && root != "" && shape.HasStructure {
		shape.SingleRoot = root
	}
	return shape, nil
}

// normalizeEntryName converts an archive entry name to forward-slash form
// with any trailing slash removed.
func normalizeEntryName(name string) string {
	norm := strings.ReplaceAll(name, "\\", "/")
	return strings.TrimSuffix(norm, "/")
}

// stripRoot removes a shared toplevel folder component, when present.
func stripRoot(norm, root string) string {
	if root == "" {
		return norm
	}
	prefix := root + "/"
	if strings.HasPrefix(norm, prefix) {
		return norm[len(prefix):]
	}
	return norm
}

// validateEntryPath rejects archive entries that could extract outside the
// destination directory. Checked on the raw normalized name, before any
// toplevel-root stripping, so a malicious ".." component can't be hidden by
// being mistaken for a shared root and stripped away.
func validateEntryPath(norm string) error {
	for _, part := range strings.Split(norm, "/") {
		if part == ".." || part == "" {
			return fmt.Errorf("archive entry escapes destination: %q", norm)
		}
	}
	return nil
}

// ExtractArchiveFlat extracts the archive's single usable file to destDir,
// stripping any SingleRoot prefix. Returns the extracted path. Errors when
// the archive is not single-file.
func ExtractArchiveFlat(archivePath string, shape ArchiveShape, destDir string, progress *atomic.Float64) (string, error) {
	if shape.UsableFiles != 1 {
		return "", errors.New("ExtractArchiveFlat requires a single-file archive")
	}
	files, err := extractArchiveTo(archivePath, shape, destDir, progress, true)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", fmt.Errorf("expected 1 extracted file, got %d", len(files))
	}
	return files[0], nil
}

// ExtractArchiveToFolder extracts all usable entries into destDir, stripping
// the SingleRoot prefix when set, rejecting path traversal, and skipping
// junk entries. Returns extracted file paths.
func ExtractArchiveToFolder(archivePath string, shape ArchiveShape, destDir string, progress *atomic.Float64) ([]string, error) {
	return extractArchiveTo(archivePath, shape, destDir, progress, false)
}

// extractArchiveTo extracts an archive's usable entries to destDir. When
// flatten is true (single-file archives only), entries land directly under
// destDir by base name instead of preserving their relative path.
func extractArchiveTo(archivePath string, shape ArchiveShape, destDir string, progress *atomic.Float64, flatten bool) ([]string, error) {
	if strings.EqualFold(filepath.Ext(archivePath), ".7z") {
		return extract7zArchiveTo(archivePath, shape, destDir, progress, flatten)
	}
	return extractZipArchiveTo(archivePath, shape, destDir, progress, flatten)
}

func extractZipArchiveTo(archivePath string, shape ArchiveShape, destDir string, progress *atomic.Float64, flatten bool) ([]string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	type entry struct {
		file *zip.File
		rel  string
	}
	var entries []entry
	var totalBytes uint64
	for _, f := range reader.File {
		if f.FileInfo().IsDir() || isJunkArchiveEntry(f.Name) {
			continue
		}
		norm := normalizeEntryName(f.Name)
		if err := validateEntryPath(norm); err != nil {
			return nil, err
		}
		rel := stripRoot(norm, shape.SingleRoot)
		if flatten {
			rel = filepath.Base(rel)
		}
		entries = append(entries, entry{f, filepath.FromSlash(rel)})
		totalBytes += f.UncompressedSize64
	}

	buffer := make([]byte, DefaultBufferSize)
	var extractedBytes uint64
	createdDirs := map[string]bool{destDir: true}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		destPath := filepath.Join(destDir, e.rel)
		parentDir := filepath.Dir(destPath)
		if !createdDirs[parentDir] {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
			}
			createdDirs[parentDir] = true
		}
		if err := extractFile(e.file, destPath, buffer, totalBytes, &extractedBytes, progress); err != nil {
			return nil, fmt.Errorf("failed to extract file %s: %w", e.file.Name, err)
		}
		out = append(out, destPath)
	}
	return out, nil
}

func extract7zArchiveTo(archivePath string, shape ArchiveShape, destDir string, progress *atomic.Float64, flatten bool) ([]string, error) {
	reader, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open 7z file: %w", err)
	}
	defer reader.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	type entry struct {
		file *sevenzip.File
		rel  string
	}
	var entries []entry
	var totalBytes uint64
	for _, f := range reader.File {
		if f.FileInfo().IsDir() || isJunkArchiveEntry(f.Name) {
			continue
		}
		norm := normalizeEntryName(f.Name)
		if err := validateEntryPath(norm); err != nil {
			return nil, err
		}
		rel := stripRoot(norm, shape.SingleRoot)
		if flatten {
			rel = filepath.Base(rel)
		}
		entries = append(entries, entry{f, filepath.FromSlash(rel)})
		totalBytes += f.UncompressedSize
	}

	buffer := make([]byte, DefaultBufferSize)
	var extractedBytes uint64
	createdDirs := map[string]bool{destDir: true}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		destPath := filepath.Join(destDir, e.rel)
		parentDir := filepath.Dir(destPath)
		if !createdDirs[parentDir] {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
			}
			createdDirs[parentDir] = true
		}
		if err := extract7zFile(e.file, destPath, buffer, totalBytes, &extractedBytes, progress); err != nil {
			return nil, fmt.Errorf("failed to extract file %s: %w", e.file.Name, err)
		}
		out = append(out, destPath)
	}
	return out, nil
}

// archiveEntryNames returns the non-directory entry names of a zip or 7z,
// detected by extension (callers only hand it archives of known type).
func archiveEntryNames(path string) ([]string, error) {
	if strings.EqualFold(filepath.Ext(path), ".7z") {
		r, err := sevenzip.OpenReader(path)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		var names []string
		for _, f := range r.File {
			if !f.FileInfo().IsDir() {
				names = append(names, f.Name)
			}
		}
		return names, nil
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var names []string
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			names = append(names, f.Name)
		}
	}
	return names, nil
}
