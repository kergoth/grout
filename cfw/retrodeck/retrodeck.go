package retrodeck

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"grout/internal/jsonutil"
)

//go:embed data/*.json
var embeddedFiles embed.FS

//go:embed input_mappings/*.json
var embeddedInputMappings embed.FS

var (
	Platforms = jsonutil.MustLoadJSONMap[string, []string](embeddedFiles, "data/platforms.json")
)

func GetInputMappingBytes() ([]byte, error) {
	overridePath := filepath.Join("overrides", "cfw", "retrodeck", "input_mappings", "steamdeck.json")
	data, err := os.ReadFile(overridePath)
	if err != nil {
		data, err = embeddedInputMappings.ReadFile("input_mappings/steamdeck.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded input mapping: %w", err)
		}
	}
	return data, nil
}

type retroDeckConfig struct {
	Version string         `json:"version"`
	Paths   retroDeckPaths `json:"paths"`
}

type retroDeckPaths struct {
	RdHomePath          string `json:"rd_home_path"`
	RomsPath            string `json:"roms_path"`
	SavesPath           string `json:"saves_path"`
	BiosPath            string `json:"bios_path"`
	DownloadedMediaPath string `json:"downloaded_media_path"`
}

const configPathEnv = "RETRODECK_CFG"

func configPath() string {
	if p := os.Getenv(configPathEnv); p != "" {
		return p
	}
	return filepath.Join(os.Getenv("HOME"), ".var", "app", "net.retrodeck.retrodeck", "config", "retrodeck", "retrodeck.json")
}

func esdeConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".var", "app", "net.retrodeck.retrodeck", "config", "ES-DE")
}

// ParseConfig reads and parses the RetroDECK config file at path.
func ParseConfig(path string) (*retroDeckConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read retrodeck.json: %w", err)
	}
	var cfg retroDeckConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse retrodeck.json: %w", err)
	}
	return &cfg, nil
}

func readConfig() (*retroDeckConfig, error) {
	return ParseConfig(configPath())
}

// ValidateConfig checks that RetroDECK and ES-DE config files required for
// path resolution are present and parseable. Call this at startup when CFW is
// RetroDECK; fail fast so the user sees a clear message before any directory
// resolution is attempted.
func ValidateConfig() error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	if cfg.Paths.RdHomePath == "" {
		return fmt.Errorf("retrodeck.json: missing rd_home_path")
	}
	esSettingsPath := filepath.Join(esdeConfigDir(), "settings", "es_settings.xml")
	if _, err := os.Stat(esSettingsPath); err != nil {
		return fmt.Errorf("ES-DE settings not found at %s: %w", esSettingsPath, err)
	}
	return nil
}

func mustReadConfig() *retroDeckConfig {
	cfg, err := readConfig()
	if err != nil {
		log.Fatalf("retrodeck: %v", err)
	}
	return cfg
}

func GetRomDirectory() string {
	return mustReadConfig().Paths.RomsPath
}

func GetBIOSDirectory() string {
	return mustReadConfig().Paths.BiosPath
}

func GetBaseSavePath() string {
	return mustReadConfig().Paths.SavesPath
}

func GetArtDirectory(romDir string) string {
	return filepath.Join(mustReadConfig().Paths.DownloadedMediaPath, mediaSystemDir(romDir), "covers")
}

func GetVideoDirectory(romDir string) string {
	return filepath.Join(mustReadConfig().Paths.DownloadedMediaPath, mediaSystemDir(romDir), "videos")
}

func GetManualDirectory(romDir string) string {
	return filepath.Join(mustReadConfig().Paths.DownloadedMediaPath, mediaSystemDir(romDir), "manuals")
}

func GetMarqueeDirectory(romDir string) string {
	return filepath.Join(mustReadConfig().Paths.DownloadedMediaPath, mediaSystemDir(romDir), "marquees")
}

func GetBoxbackDirectory(romDir string) string {
	return filepath.Join(mustReadConfig().Paths.DownloadedMediaPath, mediaSystemDir(romDir), "backcovers")
}

func GetFanartDirectory(romDir string) string {
	return filepath.Join(mustReadConfig().Paths.DownloadedMediaPath, mediaSystemDir(romDir), "fanart")
}

func GetGroutGamelist(system string) string {
	return filepath.Join(esdeConfigDir(), "gamelists", system, "gamelist.xml")
}

// mediaSystemDir extracts the ES-DE system name from a ROM directory path.
// RetroDECK ROM paths are flat (e.g. "psx", "snes"), so the first path
// component relative to the ROM root is the system name.
func mediaSystemDir(romDir string) string {
	rel, _ := filepath.Rel(GetRomDirectory(), romDir)
	rel = filepath.ToSlash(rel)
	return strings.Split(rel, "/")[0]
}
