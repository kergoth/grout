package fileutil

import "strings"

// archiveAsRomPlatforms are RomM platform FS slugs whose emulators consume
// the archive file itself as the ROM: DOSBox Pure mounts zips directly, and
// MAME/FBNeo romsets are zip/7z files by definition. ES-DE's es_systems.xml
// lists .zip/.7z as valid ROM extensions for these systems. Extracting their
// archives would break launching, so downloads for these platforms are kept
// as-is.
var archiveAsRomPlatforms = map[string]bool{
	"dos": true,

	"arcade": true,
	"mame":   true,
	"fbneo":  true,
	"fba":    true,
	"neogeo": true,

	"cps1": true,
	"cps2": true,
	"cps3": true,

	"neogeoaes": true,
	"neogeomvs": true,
	"neogeocd":  true,

	"naomi":      true,
	"naomi2":     true,
	"atomiswave": true,
}

// UsesArchiveAsRom reports whether the platform treats an archive file as the
// playable ROM artifact rather than as packaging to extract.
func UsesArchiveAsRom(fsSlug string) bool {
	return archiveAsRomPlatforms[strings.ToLower(fsSlug)]
}
