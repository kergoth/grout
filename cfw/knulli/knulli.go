package knulli

import (
	"embed"
	"grout/internal/jsonutil"
	"grout/internal/scummvm"
	"os"
	"path/filepath"
)

//go:embed data/*.json
var embeddedFiles embed.FS

var (
	Platforms       = jsonutil.MustLoadJSONMap[string, []string](embeddedFiles, "data/platforms.json")
	SaveDirectories = jsonutil.MustLoadJSONMap[string, []string](embeddedFiles, "data/save_directories.json")
)

func GetBasePath() string {
	if basePath := os.Getenv("BASE_PATH"); basePath != "" {
		return basePath
	}
	return "/userdata"
}

func GetRomDirectory() string {
	return filepath.Join(GetBasePath(), "roms")
}

func GetBIOSDirectory() string {
	return filepath.Join(GetBasePath(), "bios")
}

func GetBaseSavePath() string {
	return filepath.Join(GetBasePath(), "saves")
}

func GetArtDirectory(romDir string) string {
	return filepath.Join(romDir, "images")
}

func GetGroutGamelist() string {
	return filepath.Join(GetRomDirectory(), "tools", "gamelist.xml")
}

func GetVideoDirectory(romDir string) string {
	return filepath.Join(romDir, "videos")
}

func GetManualDirectory(romDir string) string {
	return filepath.Join(romDir, "manuals")
}

func GetBezelDirectory(romDir string) string {
	return filepath.Join(romDir, "bezels")
}

// PrepareScummVMLauncher locates an extracted ScummVM game's launcher stub
// for Knulli's manual-add convention: a top-level game folder containing a
// single ".scummvm" file whose contents are the ScummVM short ID. Knulli
// accepts any filename for the stub, so unlike EmuDeck and RetroDECK, no
// renaming is needed; the gamelist entry can point directly at the stub
// wherever it already is. It reports false, leaving extractDir unchanged,
// when extractDir does not contain exactly one usable stub.
func PrepareScummVMLauncher(extractDir string) (string, bool, error) {
	stubPath, _, ok, err := scummvm.FindStub(extractDir)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return stubPath, true, nil
}
