package emudeck

import (
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"grout/internal/jsonutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed data/*.json
var embeddedFiles embed.FS

var (
	Platforms = jsonutil.MustLoadJSONMap[string, []string](embeddedFiles, "data/platforms.json")
)

type emuDeckSettings struct {
	StoragePath string `json:"storagePath"`
}

type esSettings struct {
	Strings []esString `xml:"string"`
}

type esString struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// ValidateConfig checks that the EmuDeck and ES-DE config files required for
// path resolution are present and parseable. Call this at startup when CFW is
// EmuDeck; fail fast so the user sees a clear message before any directory
// resolution is attempted.
func ValidateConfig() error {
	settingsPath := filepath.Join(os.Getenv("HOME"), ".config", "EmuDeck", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("read EmuDeck settings.json: %w", err)
	}
	var settings emuDeckSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("parse EmuDeck settings.json: %w", err)
	}
	if settings.StoragePath == "" {
		return fmt.Errorf("EmuDeck settings.json: missing storagePath")
	}
	esSettingsPath := filepath.Join(os.Getenv("HOME"), "ES-DE", "settings", "es_settings.xml")
	if _, err := os.Stat(esSettingsPath); err != nil {
		return fmt.Errorf("ES-DE settings not found at %s: %w", esSettingsPath, err)
	}
	return nil
}

func GetRomDirectory() string {
	if rom := readESSetting("ROMDirectory"); rom != "" {
		return rom
	}
	return filepath.Join(readStoragePath(), "Emulation", "roms")
}

func GetBIOSDirectory() string {
	return filepath.Join(readStoragePath(), "Emulation", "bios")
}

func GetBaseSavePath() string {
	return filepath.Join(readStoragePath(), "Emulation", "saves")
}

func GetArtDirectory(romDir string) string {
	return filepath.Join(readMediaDirectory(), mediaSystemDir(romDir), "covers")
}

func GetVideoDirectory(romDir string) string {
	return filepath.Join(readMediaDirectory(), mediaSystemDir(romDir), "videos")
}

func GetManualDirectory(romDir string) string {
	return filepath.Join(readMediaDirectory(), mediaSystemDir(romDir), "manuals")
}

func GetMarqueeDirectory(romDir string) string {
	return filepath.Join(readMediaDirectory(), mediaSystemDir(romDir), "marquees")
}

func GetBoxbackDirectory(romDir string) string {
	return filepath.Join(readMediaDirectory(), mediaSystemDir(romDir), "backcovers")
}

func GetFanartDirectory(romDir string) string {
	return filepath.Join(readMediaDirectory(), mediaSystemDir(romDir), "fanart")
}

func GetGroutGamelist(system string) string {
	return filepath.Join(os.Getenv("HOME"), "ES-DE", "gamelists", system, "gamelist.xml")
}

// mediaSystemDir extracts the ES-DE system name from a ROM directory path.
// EmuDeck ROM paths are at most two levels deep relative to the ROM root
// (e.g. "psx", "model2/roms", "wiiu/roms"), so the first path component is
// always the system name used by ES-DE's media tree.
func mediaSystemDir(romDir string) string {
	rel, _ := filepath.Rel(GetRomDirectory(), romDir)
	rel = filepath.ToSlash(rel)
	return strings.Split(rel, "/")[0]
}

func readStoragePath() string {
	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "EmuDeck", "settings.json"))
	if err != nil {
		log.Fatalf("emudeck: read settings.json: %v", err)
	}
	var settings emuDeckSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Fatalf("emudeck: parse settings.json: %v", err)
	}
	if settings.StoragePath == "" {
		log.Fatalf("emudeck: settings.json missing storagePath")
	}
	return settings.StoragePath
}

func readMediaDirectory() string {
	if media := readESSetting("MediaDirectory"); media != "" {
		return media
	}
	return filepath.Join(readStoragePath(), "Emulation", "tools", "downloaded_media")
}

func readESSetting(name string) string {
	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), "ES-DE", "settings", "es_settings.xml"))
	if err != nil {
		return ""
	}
	var settings esSettings
	if err := xml.Unmarshal(data, &settings); err != nil {
		return ""
	}
	for _, entry := range settings.Strings {
		if entry.Name == name {
			return entry.Value
		}
	}
	return ""
}
