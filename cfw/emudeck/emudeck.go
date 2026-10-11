// Package emudeck finds where an EmuDeck install on a Steam Deck keeps things.
//
// EmuDeck records where it lives in ~/.config/EmuDeck/settings.json, and ES-DE
// (its frontend) records the rom and media folders in its own settings, so
// both are read rather than assumed: the install may be on the SD card.
package emudeck

import (
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"grout/cfw/esde"
	"grout/tables"
)

//go:embed data/*.json
var embeddedFiles embed.FS

// Platforms is the ES-DE table with the folders an EmuDeck install adds or
// lays out differently.
var Platforms = esde.With(tables.MustLoad[string, []string](embeddedFiles, "data/platforms.json"))

// Paths are the locations grout needs.
type Paths struct {
	Roms  string
	BIOS  string
	Saves string
	Media string
	// Gamelists is ES-DE's gamelist tree, with a folder per system.
	Gamelists string
}

// LoadPaths reads the install's locations for the user whose home is home.
//
// A settings file that cannot be read still yields usable paths, from an
// install in the home directory, along with the error to report: a wrong guess
// is better than a launcher that will not start.
func LoadPaths(home string) (Paths, error) {
	emulation := filepath.Join(home, "Emulation")
	var loadErr error

	storage, err := readStoragePath(filepath.Join(home, ".config", "EmuDeck", "settings.json"))
	if err != nil {
		loadErr = err
	} else {
		emulation = filepath.Join(storage, "Emulation")
	}

	es := readESSettings(filepath.Join(home, "ES-DE", "settings", "es_settings.xml"))

	return Paths{
		Roms:      orDefault(es["ROMDirectory"], filepath.Join(emulation, "roms")),
		BIOS:      filepath.Join(emulation, "bios"),
		Saves:     filepath.Join(emulation, "saves"),
		Media:     orDefault(es["MediaDirectory"], filepath.Join(emulation, "tools", "downloaded_media")),
		Gamelists: filepath.Join(home, "ES-DE", "gamelists"),
	}, loadErr
}

func orDefault(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

var paths = sync.OnceValue(func() Paths {
	p, err := LoadPaths(os.Getenv("HOME"))
	if err != nil {
		slog.Default().Error("Could not read the EmuDeck settings, assuming an install in the home folder", "error", err)
	}
	return p
})

func GetRomDirectory() string                { return paths().Roms }
func GetBIOSDirectory() string               { return paths().BIOS }
func GetBaseSavePath() string                { return paths().Saves }
func GetMediaDirectory() string              { return paths().Media }
func GetBezelDirectory(romDir string) string { return filepath.Join(romDir, "bezels") }

// GetGroutGamelist is the gamelist the launcher shortcut goes in: grout is
// installed as a port.
func GetGroutGamelist() string {
	return filepath.Join(paths().Gamelists, "ports", "gamelist.xml")
}

// GetGamelistPath is where ES-DE reads a system's gamelist from.
func GetGamelistPath(romDir, fileName string) string {
	return filepath.Join(paths().Gamelists, System(romDir), fileName)
}

// System is the ES-DE system a rom directory belongs to, which names its
// gamelist and media folders. Some systems keep their roms a level down
// ("wiiu/roms"), so it is the first folder under the rom root, not the last.
func System(romDir string) string {
	rel, err := filepath.Rel(GetRomDirectory(), romDir)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return filepath.Base(romDir)
	}
	return strings.Split(filepath.ToSlash(rel), "/")[0]
}

func readStoragePath(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading EmuDeck settings: %w", err)
	}
	var settings struct {
		StoragePath string `json:"storagePath"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return "", fmt.Errorf("parsing EmuDeck settings: %w", err)
	}
	if settings.StoragePath == "" {
		return "", fmt.Errorf("EmuDeck settings have no storagePath")
	}
	return settings.StoragePath, nil
}

// readESSettings returns ES-DE's string settings, or nothing if the file is
// missing or unreadable, in which case the defaults stand.
//
// The file has no root element, just one element per setting, so it is read as
// a stream of elements rather than unmarshalled as a document.
func readESSettings(path string) map[string]string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	values := map[string]string{}
	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		element, ok := token.(xml.StartElement)
		if !ok || element.Name.Local != "string" {
			continue
		}

		var name, value string
		for _, attr := range element.Attr {
			switch attr.Name.Local {
			case "name":
				name = attr.Value
			case "value":
				value = attr.Value
			}
		}
		values[name] = value
	}
	return values
}
