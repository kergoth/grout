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

// Stub discovery and short-ID validation are covered by
// internal/scummvm.FindStub's own tests. The cases below cover only
// EmuDeck's rename policy once a stub has already been found.

func TestPrepareScummVMLauncherRenamesLauncherFromStub(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Game Name")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "launcher.scummvm"), []byte("sky\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	launcherPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ready {
		t.Fatal("ready = false, want true")
	}
	wantPath := filepath.Join(extractDir, "sky.scummvm")
	if launcherPath != wantPath {
		t.Fatalf("launcher path = %q, want %q", launcherPath, wantPath)
	}
	content, err := os.ReadFile(launcherPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "sky\n" {
		t.Fatalf("launcher contents = %q, want %q", content, "sky\n")
	}
}

func TestPrepareScummVMLauncherKeepsAlreadyCanonicalLayout(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Game Name")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "sky.scummvm"), []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}

	launcherPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ready {
		t.Fatal("ready = false, want true")
	}
	wantPath := filepath.Join(extractDir, "sky.scummvm")
	if launcherPath != wantPath {
		t.Fatalf("launcher path = %q, want %q", launcherPath, wantPath)
	}
}

func TestPrepareScummVMLauncherNotReadyWhenNoStub(t *testing.T) {
	extractDir := t.TempDir()

	launcherPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if ready {
		t.Fatal("ready = true, want false")
	}
	if launcherPath != "" {
		t.Fatalf("launcher path = %q, want empty", launcherPath)
	}
}

func TestPrepareScummVMLauncherLeavesSourceWhenLauncherDestinationExists(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Game Name")
	stubPath := filepath.Join(extractDir, "launcher.scummvm")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stubPath, []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(extractDir, "sky.scummvm"), 0o755); err != nil {
		t.Fatal(err)
	}

	launcherPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err == nil {
		t.Fatal("expected existing destination to fail")
	}
	if ready {
		t.Fatal("launcher reported ready after destination conflict")
	}
	if launcherPath != "" {
		t.Fatalf("launcher path = %q, want empty", launcherPath)
	}
	if _, err := os.Stat(stubPath); err != nil {
		t.Fatalf("source launcher changed: %v", err)
	}
}

func TestPrepareScummVMLauncherLeavesSourceWhenLauncherRenameFails(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Game Name")
	stubPath := filepath.Join(extractDir, "launcher.scummvm")
	if err := os.MkdirAll(filepath.Join(extractDir, "sky.scummvm"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stubPath, []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}

	launcherPath, ready, err := PrepareScummVMLauncher(extractDir)
	if err == nil {
		t.Fatal("expected launcher rename to fail")
	}
	if ready {
		t.Fatal("launcher reported ready after rename failure")
	}
	if launcherPath != "" {
		t.Fatalf("launcher path = %q, want empty", launcherPath)
	}
	if _, err := os.Stat(stubPath); err != nil {
		t.Fatalf("source launcher changed: %v", err)
	}
}
