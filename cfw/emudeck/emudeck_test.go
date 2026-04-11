package emudeck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRuntimeDerivesRomMediaAndGamelistRoots(t *testing.T) {
	home := t.TempDir()
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
	settingsXML := `<?xml version="1.0"?><config><string name="ROMDirectory" value="/run/media/deck/SD/Emulation/roms"/><string name="MediaDirectory" value="/run/media/deck/SD/Emulation/tools/downloaded_media"/></config>`
	if err := os.WriteFile(filepath.Join(home, "ES-DE", "settings", "es_settings.xml"), []byte(settingsXML), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := GetRomDirectory(); got != "/run/media/deck/SD/Emulation/roms" {
		t.Fatalf("rom dir = %q", got)
	}
	if got := GetArtDirectory("/run/media/deck/SD/Emulation/roms/psx"); got != "/run/media/deck/SD/Emulation/tools/downloaded_media/psx/covers" {
		t.Fatalf("cover dir = %q", got)
	}
	if got := GetGroutGamelist("psx"); got != filepath.Join(home, "ES-DE", "gamelists", "psx", "gamelist.xml") {
		t.Fatalf("gamelist path = %q", got)
	}
}

func TestValidateConfigFailsWithoutSettings(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := ValidateConfig(); err == nil {
		t.Fatal("expected error when settings.json missing")
	}
}

func TestValidateConfigFailsWithoutESDE(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".config", "EmuDeck"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".config", "EmuDeck", "settings.json"), []byte(`{"storagePath":"/run/media/deck/SD"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// ES-DE settings directory not created — ValidateConfig must fail
	if err := ValidateConfig(); err == nil {
		t.Fatal("expected error when es_settings.xml missing")
	}
}
