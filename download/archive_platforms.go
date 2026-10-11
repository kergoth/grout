package download

import "strings"

// archiveAsRom lists the platform slugs whose emulators take the archive
// itself as the rom, so unpacking one leaves nothing they can load. It is
// kept in step with the list Argosy Launcher (also a RomM client) uses,
// including the alternate spellings of a slug it resolves to the same
// platform.
//
// Exact slugs only: Argosy also resolves a slug with a suffix (such as
// "mame-2003") to its base, which needs the whole platform table to avoid
// matching a different platform, "neo-geo-pocket" for one.
var archiveAsRom = map[string]bool{
	// DOS.
	"dos":    true,
	"ibm-pc": true,
	"ibm_pc": true,
	"ms-dos": true,
	"ms_dos": true,
	"msdos":  true,
	"pcxt":   true,

	// Arcade romsets.
	"arcade":   true,
	"mame":     true,
	"hbmame":   true,
	"fbneo":    true,
	"fba":      true,
	"fb_alpha": true,

	// Neo Geo.
	"neogeo":           true,
	"neo-geo":          true,
	"neo_geo":          true,
	"neogeoaes":        true,
	"neogeo-aes":       true,
	"neogeo_aes":       true,
	"neogeomvs":        true,
	"neo-geo-mvs":      true,
	"neogeo-mvs":       true,
	"neogeo_mvs":       true,
	"neogeocd":         true,
	"neo-geo-cd":       true,
	"neogeo-cd":        true,
	"neogeo_cd":        true,
	"neo_geo_cd":       true,
	"neocd":            true,
	"hyperneogeo64":    true,
	"hyper-neo-geo-64": true,

	// Capcom, Sega and Sammy arcade boards.
	"cps1":       true,
	"cps2":       true,
	"cps3":       true,
	"cps-1":      true,
	"cps-2":      true,
	"cps-3":      true,
	"naomi":      true,
	"naomi2":     true,
	"hikaru":     true,
	"atomiswave": true,
	"model2":     true,
	"model3":     true,

	// Handhelds and consoles whose packages are zips.
	"gameandwatch":     true,
	"game-and-watch":   true,
	"game_and_watch":   true,
	"g-and-w":          true,
	"vita":             true,
	"psvita":           true,
	"ps-vita":          true,
	"ps_vita":          true,
	"playstation-vita": true,
	"playstation_vita": true,
}

// alwaysExtract lists the platform slugs whose emulators never read a zip,
// so even an archive holding one flat rom is unpacked whatever the user
// prefers. Modern consoles take disc images, which are already compressed, or
// their own container formats. Mirrors Argosy Launcher's list, with the same
// alternate spellings.
var alwaysExtract = map[string]bool{
	// Nintendo.
	"3ds":               true,
	"ds":                true,
	"gamecube":          true,
	"gc":                true,
	"n3ds":              true,
	"nds":               true,
	"new-nintendo-3ds":  true,
	"new_nintendo3ds":   true,
	"ngc":               true,
	"nintendo-3ds":      true,
	"nintendo-ds":       true,
	"nintendo-gamecube": true,
	"nintendo-switch":   true,
	"nintendo_3ds":      true,
	"nintendo_ds":       true,
	"nintendoswitch":    true,
	"nsw":               true,
	"nswitch":           true,
	"switch":            true,
	"wii":               true,
	"wii-u":             true,
	"wiiu":              true,

	// Sony.
	"playstation":                true,
	"playstation-1":              true,
	"playstation-2":              true,
	"playstation-3":              true,
	"playstation-portable":       true,
	"playstation-portable-minis": true,
	"playstation1":               true,
	"playstation2":               true,
	"playstation3":               true,
	"playstation_1":              true,
	"playstation_2":              true,
	"playstation_3":              true,
	"playstation_portable":       true,
	"playstation_portable_minis": true,
	"playstationportable":        true,
	"ps":                         true,
	"ps1":                        true,
	"ps2":                        true,
	"ps3":                        true,
	"psp":                        true,
	"psp-minis":                  true,
	"psp_minis":                  true,
	"pspminis":                   true,
	"psx":                        true,
	"sony-playstation":           true,
	"sony-playstation-2":         true,
	"sony-playstation-3":         true,
	"sony_playstation":           true,
	"sony_playstation_2":         true,
	"sony_playstation_3":         true,

	// Sega.
	"dc":             true,
	"dreamcast":      true,
	"saturn":         true,
	"sega-dreamcast": true,
	"sega-saturn":    true,
	"sega_dreamcast": true,
	"sega_saturn":    true,
}

// AlwaysExtracts reports whether the platform cannot use a zipped rom.
func AlwaysExtracts(fsSlug string) bool {
	return alwaysExtract[strings.ToLower(fsSlug)]
}

// KeepsArchive reports whether the platform plays archives as they are.
func KeepsArchive(fsSlug string) bool {
	return archiveAsRom[strings.ToLower(fsSlug)]
}
