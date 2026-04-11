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

func TestKnownEmuDeckMediaDirectories(t *testing.T) {
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
	// No ROMDirectory or MediaDirectory overrides — paths derive from storagePath.
	if err := os.WriteFile(filepath.Join(home, "ES-DE", "settings", "es_settings.xml"), []byte(`<?xml version="1.0"?><config></config>`), 0o644); err != nil {
		t.Fatal(err)
	}

	romDir := "/run/media/deck/SD/Emulation/roms/model2/roms"
	if got := GetArtDirectory(romDir); got != "/run/media/deck/SD/Emulation/tools/downloaded_media/model2/covers" {
		t.Fatalf("cover dir = %q", got)
	}
	if got := GetManualDirectory(romDir); got != "/run/media/deck/SD/Emulation/tools/downloaded_media/model2/manuals" {
		t.Fatalf("manual dir = %q", got)
	}
	if got := GetFanartDirectory(romDir); got != "/run/media/deck/SD/Emulation/tools/downloaded_media/model2/fanart" {
		t.Fatalf("fanart dir = %q", got)
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
