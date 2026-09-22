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

//go:embed input_mappings/*.json
var embeddedInputMappings embed.FS

var (
	Platforms = jsonutil.MustLoadJSONMap[string, []string](embeddedFiles, "data/platforms.json")
)

func GetInputMappingBytes() ([]byte, error) {
	overridePath := filepath.Join("overrides", "cfw", "emudeck", "input_mappings", "steamdeck.json")
	data, err := os.ReadFile(overridePath)
	if err != nil {
		data, err = embeddedInputMappings.ReadFile("input_mappings/steamdeck.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded input mapping: %w", err)
		}
	}
	return data, nil
}

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

// PrepareScummVMLauncher renames an extracted ScummVM game's populated launcher
// stub to <shortid>.scummvm, retaining the game-name directory for ES-DE.
// It reports false when the source directory does not contain exactly one
// usable stub, leaving the directory unchanged.
func PrepareScummVMLauncher(extractDir string) (string, bool, error) {
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return "", false, fmt.Errorf("read extracted ScummVM game: %w", err)
	}

	var stubPath string
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".scummvm") {
			continue
		}
		if stubPath != "" {
			return "", false, nil
		}
		stubPath = filepath.Join(extractDir, entry.Name())
	}
	if stubPath == "" {
		return "", false, nil
	}

	contents, err := os.ReadFile(stubPath)
	if err != nil {
		return "", false, fmt.Errorf("read ScummVM launcher stub: %w", err)
	}
	shortID := strings.TrimSpace(string(contents))
	if shortID == "" || strings.ContainsAny(shortID, "/\\\r\n") || filepath.Base(shortID) != shortID || shortID == "." {
		return "", false, nil
	}

	launcherName := shortID + ".scummvm"
	launcherPath := filepath.Join(extractDir, launcherName)
	if launcherPath != stubPath {
		if err := os.Rename(stubPath, launcherPath); err != nil {
			return "", false, fmt.Errorf("rename ScummVM launcher stub: %w", err)
		}
	}
	return launcherPath, true, nil
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
