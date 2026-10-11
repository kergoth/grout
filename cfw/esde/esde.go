// Package esde holds what the firmwares that run ES-DE as their frontend
// (RetroDECK, EmuDeck) have in common.
package esde

import (
	"embed"
	"maps"

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
