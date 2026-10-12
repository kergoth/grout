package emudeck

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadPaths_FollowsTheInstallAndES(t *testing.T) {
	home := t.TempDir()
	write(t, filepath.Join(home, ".config", "EmuDeck", "settings.json"), `{"storagePath": "/run/media/deck/SD"}`)
	write(t, filepath.Join(home, "ES-DE", "settings", "es_settings.xml"), `<?xml version="1.0"?>
<bool name="DebugMode" value="false" />
<string name="MediaDirectory" value="/run/media/deck/SD/Emulation/tools/downloaded_media" />
<string name="ROMDirectory" value="/run/media/deck/SD/Emulation/roms" />`)

	got, err := LoadPaths(home)
	if err != nil {
		t.Fatalf("LoadPaths: %v", err)
	}
	want := Paths{
		Roms:      "/run/media/deck/SD/Emulation/roms",
		BIOS:      "/run/media/deck/SD/Emulation/bios",
		Saves:     "/run/media/deck/SD/Emulation/saves",
		Media:     "/run/media/deck/SD/Emulation/tools/downloaded_media",
		Gamelists: filepath.Join(home, "ES-DE", "gamelists"),
	}
	if got != want {
		t.Errorf("LoadPaths = %+v, want %+v", got, want)
	}
}

// ES-DE lets the user move media and roms; the settings win over where
// EmuDeck would have put them.
func TestLoadPaths_ESSettingsOverrideDefaults(t *testing.T) {
	home := t.TempDir()
	write(t, filepath.Join(home, ".config", "EmuDeck", "settings.json"), `{"storagePath": "/mnt/sd"}`)
	write(t, filepath.Join(home, "ES-DE", "settings", "es_settings.xml"),
		`<?xml version="1.0"?>
<bool name="DebugMode" value="false" />
<string name="MediaDirectory" value="/elsewhere/media" />`)

	got, err := LoadPaths(home)
	if err != nil {
		t.Fatal(err)
	}
	if got.Media != "/elsewhere/media" {
		t.Errorf("Media = %q, want the ES-DE setting", got.Media)
	}
	if got.Roms != "/mnt/sd/Emulation/roms" {
		t.Errorf("Roms = %q, want the default under the install", got.Roms)
	}
}

func TestLoadPaths_MissingSettingsFallBackToHome(t *testing.T) {
	home := t.TempDir()

	got, err := LoadPaths(home)
	if err == nil {
		t.Error("expected the unreadable settings to be reported")
	}
	if want := filepath.Join(home, "Emulation", "roms"); got.Roms != want {
		t.Errorf("Roms = %q, want %q", got.Roms, want)
	}
}

func TestLoadPaths_EmptyStoragePathIsAnError(t *testing.T) {
	home := t.TempDir()
	write(t, filepath.Join(home, ".config", "EmuDeck", "settings.json"), `{"storagePath": ""}`)

	if _, err := LoadPaths(home); err == nil {
		t.Error("expected an error for a settings file with no storagePath")
	}
}

func TestPlatforms_EmuDeckLayersOverESDE(t *testing.T) {
	cases := map[string][]string{
		"wiiu":  {"wiiu/roms"}, // EmuDeck keeps wiiu roms a level down
		"psx":   {"psx"},       // untouched ES-DE entry
		"win9x": {"windows9x"}, // only EmuDeck separates these
	}
	for slug, want := range cases {
		got := Platforms[slug]
		if len(got) != len(want) || got[0] != want[0] {
			t.Errorf("Platforms[%q] = %v, want %v", slug, got, want)
		}
	}
}
