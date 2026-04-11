package internal

import (
	"os"
	"path/filepath"
	"testing"

	"grout/romm"
)

func TestSlotPreference_DefaultsToAutosave(t *testing.T) {
	c := Config{}
	if got := c.GetSlotPreference(1); got != "autosave" {
		t.Errorf("default slot = %q, want autosave", got)
	}
}

// Picking "autosave" must persist as an EXPLICIT preference (not be discarded), so it
// can override a sticky recorded slot. Otherwise a user can never switch a ROM back to
// autosave once another slot has been recorded (issue #250).
func TestSetSlotPreference_AutosavePersistsAsExplicit(t *testing.T) {
	c := Config{}
	c.SetSlotPreference(1, "quicksave")
	c.SetSlotPreference(1, "autosave") // user explicitly chooses autosave

	slot, ok := c.SlotPreferenceExplicit(1)
	if !ok || slot != "autosave" {
		t.Errorf("explicit autosave should persist: got (%q, %v), want (\"autosave\", true)", slot, ok)
	}
	if got := c.GetSlotPreference(1); got != "autosave" {
		t.Errorf("GetSlotPreference = %q, want autosave", got)
	}
}

func TestGetPlatformGamelistPathUsesPlatformRomDirectorySystemName(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CFW", "EMUDECK")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".config", "EmuDeck"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, "ES-DE", "settings"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".config", "EmuDeck", "settings.json"), []byte(`{"storagePath":"/run/media/deck/SD"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	settingsXML := `<?xml version="1.0"?><config><string name="ROMDirectory" value="/run/media/deck/SD/Emulation/roms"/></config>`
	if err := os.WriteFile(filepath.Join(home, "ES-DE", "settings", "es_settings.xml"), []byte(settingsXML), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		DirectoryMappings: map[string]DirectoryMapping{
			"arcade": {
				RomMSlug:     "arcade",
				RelativePath: "model2/roms",
			},
		},
	}
	platform := romm.Platform{FSSlug: "arcade", Name: "Arcade"}

	got := cfg.GetPlatformGamelistPath(platform)
	want := filepath.Join(home, "ES-DE", "gamelists", "model2", "gamelist.xml")
	if got != want {
		t.Fatalf("gamelist path = %q, want %q", got, want)
	}
}

func TestGetPlatformRomDirectoryUsesEmuDeckAliasForModel2(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CFW", "EMUDECK")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".config", "EmuDeck"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, "ES-DE", "settings"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".config", "EmuDeck", "settings.json"), []byte(`{"storagePath":"/run/media/deck/SD"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	settingsXML := `<?xml version="1.0"?><config><string name="ROMDirectory" value="/run/media/deck/SD/Emulation/roms"/></config>`
	if err := os.WriteFile(filepath.Join(home, "ES-DE", "settings", "es_settings.xml"), []byte(settingsXML), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{}
	platform := romm.Platform{FSSlug: "model2", Name: "Sega Model 2"}

	got := cfg.GetPlatformRomDirectory(platform)
	want := "/run/media/deck/SD/Emulation/roms/model2/roms"
	if got != want {
		t.Fatalf("rom dir = %q, want %q", got, want)
	}
}
