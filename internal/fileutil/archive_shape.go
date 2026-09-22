package fileutil

import (
	"archive/zip"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
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
