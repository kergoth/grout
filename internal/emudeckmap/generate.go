// Package emudeckmap generates the EmuDeck platform map by merging upstream
// ES-DE system definitions with the EmuDeck custom-systems overlay.
package emudeckmap

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

// esdeToRomM maps ES-DE platform names to their canonical RomM slug(s).
//
// ES-DE's <platform> field uses its own naming conventions that diverge from
// RomM slugs for many well-known systems. Entries with multiple RomM slugs
// represent systems that RomM tracks under more than one slug (e.g. Neo Geo
// MVS and AES share the same ES-DE folder but are separate RomM platforms).
//
// Systems where the ES-DE platform name already matches the RomM slug (e.g.
// "psx", "snes", "arcade") are not listed here and pass through unchanged.
var esdeToRomM = map[string][]string{
	// Atari
	"atarijaguar":    {"jaguar"},
	"atarijaguarcd":  {"jaguar"},
	"atarilynx":      {"lynx"},
	"atarist":        {"atari-st"},
	// Sega
	"dreamcast":      {"dc"},
	"mastersystem":   {"sms"},
	"megadrive":      {"genesis"},
	"sega32x":        {"sega32"},
	"sg-1000":        {"sg1000"},
	// Nintendo
	"gc":             {"ngc"},
	"n3ds":           {"3ds"},
	// NEC
	"pcengine":       {"tg16"},
	"pcenginecd":     {"turbografx-cd"},
	"pcfx":           {"pc-fx"},
	// SNK
	"neogeo":         {"neogeomvs", "neogeoaes"},
	"neogeocd":       {"neo-geo-cd"},
	"ngp":            {"neo-geo-pocket"},
	"ngpc":           {"neo-geo-pocket-color"},
	// Handheld / LCD
	"gameandwatch":   {"g-and-w"},
	"lcdgames":       {"g-and-w"},
	// Computers
	"amstradcpc":     {"acpc"},
	"cdimono1":       {"philips-cd-i"},
	"channelf":       {"fairchild-channel-f"},
	"odyssey2":       {"odyssey"},
	"pc88":           {"pc-8000"},
	"pc98":           {"pc-9800-series"},
	"x68000":         {"sharp-x68000"},
	"zxspectrum":     {"zxs"},
	// PC / Windows
	"pc":             {"windows"},
	"pcwindows":      {"windows"},
	// Misc handhelds
	"megaduck":       {"mega-duck-slash-cougar-boy"},
	"pico8":          {"pico-8"},
	"pokemini":       {"pokemon-mini"},
	"tic80":          {"tic-80"},
	"wonderswancolor": {"wonderswan-color"},
}

type systemList struct {
	Systems []system `xml:"system"`
}

type system struct {
	Name     string `xml:"name"`
	Platform string `xml:"platform"`
	Path     string `xml:"path"`
}

func GeneratePlatforms(baseXML, overlayXML []byte) (map[string][]string, error) {
	base, err := parseSystems(baseXML)
	if err != nil {
		return nil, err
	}
	overlay, err := parseSystems(overlayXML)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]system)
	order := make([]string, 0, len(base)+len(overlay))
	// appendSystem merges an entry into the map. The overlay is applied after
	// the base, so overlay values win on lookup; base entries retain their
	// original position in output order.
	appendSystem := func(entry system) {
		if _, exists := merged[entry.Name]; !exists {
			order = append(order, entry.Name)
		}
		merged[entry.Name] = entry
	}
	for _, entry := range base {
		appendSystem(entry)
	}
	for _, entry := range overlay {
		appendSystem(entry)
	}

	result := make(map[string][]string)
	for _, name := range order {
		entry := merged[name]
		if entry.Platform == "" {
			continue
		}
		if !strings.HasPrefix(entry.Path, "%ROMPATH%/") {
			continue
		}
		rel := filepath.ToSlash(strings.TrimPrefix(entry.Path, "%ROMPATH%/"))
		if rel == "" || rel == "." {
			continue
		}
		// Split comma-separated ES-DE platform fields, normalize each part to
		// its canonical RomM slug(s), and accumulate folder paths per slug.
		for _, esPlatform := range splitPlatforms(entry.Platform) {
			for _, rommSlug := range toRomMSlugs(esPlatform) {
				result[rommSlug] = appendUnique(result[rommSlug], rel)
			}
		}
	}
	return result, nil
}

// splitPlatforms splits an ES-DE platform field on commas and trims whitespace.
func splitPlatforms(platform string) []string {
	parts := strings.Split(platform, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// toRomMSlugs maps a single ES-DE platform name to one or more canonical RomM
// slugs. Returns the input unchanged when no mapping is registered.
func toRomMSlugs(esPlatform string) []string {
	if slugs, ok := esdeToRomM[esPlatform]; ok {
		return slugs
	}
	return []string{esPlatform}
}

func parseSystems(raw []byte) ([]system, error) {
	var doc systemList
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("parse es_systems.xml: %w", err)
	}
	return doc.Systems, nil
}

func appendUnique(in []string, value string) []string {
	for _, existing := range in {
		if existing == value {
			return in
		}
	}
	return append(in, value)
}
