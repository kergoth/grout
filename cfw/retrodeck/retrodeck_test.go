package retrodeck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfigParsesRetrodeckJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("RETRODECK_CFG", "")

	configDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "retrodeck")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configJSON := `{
		"version": "0.10.7b",
		"paths": {
			"rd_home_path": "/run/media/deck/SD/retrodeck",
			"roms_path": "/run/media/deck/SD/retrodeck/roms",
			"saves_path": "/run/media/deck/SD/retrodeck/saves",
			"bios_path": "/run/media/deck/SD/retrodeck/bios",
			"downloaded_media_path": "/run/media/deck/SD/retrodeck/ES-DE/downloaded_media"
		}
	}`
	if err := os.WriteFile(filepath.Join(configDir, "retrodeck.json"), []byte(configJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := readConfig()
	if err != nil {
		t.Fatalf("readConfig failed: %v", err)
	}
	if cfg.Paths.RomsPath != "/run/media/deck/SD/retrodeck/roms" {
		t.Fatalf("roms_path = %q", cfg.Paths.RomsPath)
	}
	if cfg.Paths.BiosPath != "/run/media/deck/SD/retrodeck/bios" {
		t.Fatalf("bios_path = %q", cfg.Paths.BiosPath)
	}
}

func TestConfigPathEnvVarOverride(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	customDir := t.TempDir()
	customPath := filepath.Join(customDir, "retrodeck.json")
	configJSON := `{"version":"0.10.7b","paths":{"rd_home_path":"/custom","roms_path":"/custom/roms","saves_path":"/custom/saves","bios_path":"/custom/bios","downloaded_media_path":"/custom/media"}}`
	if err := os.WriteFile(customPath, []byte(configJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("RETRODECK_CFG", customPath)

	cfg, err := readConfig()
	if err != nil {
		t.Fatalf("readConfig with RETRODECK_CFG failed: %v", err)
	}
	if cfg.Paths.RdHomePath != "/custom" {
		t.Fatalf("expected /custom, got %q", cfg.Paths.RdHomePath)
	}
}

func TestParseConfig(t *testing.T) {
	configJSON := `{"version":"0.10.7b","paths":{"rd_home_path":"/r","roms_path":"/r/roms","saves_path":"/r/saves","bios_path":"/r/bios","downloaded_media_path":"/r/media"}}`
	f, err := os.CreateTemp("", "retrodeck-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	if _, err := f.WriteString(configJSON); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := ParseConfig(f.Name())
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}
	if cfg.Paths.RomsPath != "/r/roms" {
		t.Errorf("RomsPath = %q, want /r/roms", cfg.Paths.RomsPath)
	}
}

func TestParseConfig_FileNotFound(t *testing.T) {
	_, err := ParseConfig("/nonexistent/retrodeck.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseConfig_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "retrodeck-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	if _, err := f.WriteString(`{ not valid`); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = ParseConfig(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestValidateConfigFailsWithoutRetrodeckJSON(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("RETRODECK_CFG", "")
	if err := ValidateConfig(); err == nil {
		t.Fatal("expected error when retrodeck.json missing")
	}
}

func TestValidateConfigFailsWithEmptyRdHomePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("RETRODECK_CFG", "")

	configDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "retrodeck")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configJSON := `{"version": "0.10.7b", "paths": {"rd_home_path": ""}}`
	if err := os.WriteFile(filepath.Join(configDir, "retrodeck.json"), []byte(configJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateConfig(); err == nil {
		t.Fatal("expected error when rd_home_path empty")
	}
}

func TestValidateConfigFailsWithoutESDE(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("RETRODECK_CFG", "")

	configDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "retrodeck")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configJSON := `{"version": "0.10.7b", "paths": {"rd_home_path": "/run/media/deck/SD/retrodeck"}}`
	if err := os.WriteFile(filepath.Join(configDir, "retrodeck.json"), []byte(configJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	// ES-DE config directory not created
	if err := ValidateConfig(); err == nil {
		t.Fatal("expected error when ES-DE config missing")
	}
}

func TestValidateConfigSucceedsWithValidConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("RETRODECK_CFG", "")

	configDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "retrodeck")
	esdeDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "ES-DE", "settings")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(esdeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configJSON := `{"version": "0.10.7b", "paths": {"rd_home_path": "/run/media/deck/SD/retrodeck"}}`
	if err := os.WriteFile(filepath.Join(configDir, "retrodeck.json"), []byte(configJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(esdeDir, "es_settings.xml"), []byte(`<?xml version="1.0"?><config/>`), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateConfig(); err != nil {
		t.Fatalf("ValidateConfig failed: %v", err)
	}
}

func setupValidConfig(t *testing.T) string {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("RETRODECK_CFG", "")

	configDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "retrodeck")
	esdeDir := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "ES-DE", "settings")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(esdeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configJSON := `{
		"version": "0.10.7b",
		"paths": {
			"rd_home_path": "/run/media/deck/SD/retrodeck",
			"roms_path": "/run/media/deck/SD/retrodeck/roms",
			"saves_path": "/run/media/deck/SD/retrodeck/saves",
			"bios_path": "/run/media/deck/SD/retrodeck/bios",
			"downloaded_media_path": "/run/media/deck/SD/retrodeck/ES-DE/downloaded_media"
		}
	}`
	if err := os.WriteFile(filepath.Join(configDir, "retrodeck.json"), []byte(configJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(esdeDir, "es_settings.xml"), []byte(`<?xml version="1.0"?><config/>`), 0o644); err != nil {
		t.Fatal(err)
	}

	return home
}

func TestGetRomDirectory(t *testing.T) {
	setupValidConfig(t)
	if got := GetRomDirectory(); got != "/run/media/deck/SD/retrodeck/roms" {
		t.Fatalf("GetRomDirectory() = %q", got)
	}
}

func TestGetBIOSDirectory(t *testing.T) {
	setupValidConfig(t)
	if got := GetBIOSDirectory(); got != "/run/media/deck/SD/retrodeck/bios" {
		t.Fatalf("GetBIOSDirectory() = %q", got)
	}
}

func TestGetBaseSavePath(t *testing.T) {
	setupValidConfig(t)
	if got := GetBaseSavePath(); got != "/run/media/deck/SD/retrodeck/saves" {
		t.Fatalf("GetBaseSavePath() = %q", got)
	}
}

func TestGetArtDirectory(t *testing.T) {
	setupValidConfig(t)
	romDir := "/run/media/deck/SD/retrodeck/roms/psx"
	if got := GetArtDirectory(romDir); got != "/run/media/deck/SD/retrodeck/ES-DE/downloaded_media/psx/covers" {
		t.Fatalf("GetArtDirectory() = %q", got)
	}
}

func TestGetVideoDirectory(t *testing.T) {
	setupValidConfig(t)
	romDir := "/run/media/deck/SD/retrodeck/roms/snes"
	if got := GetVideoDirectory(romDir); got != "/run/media/deck/SD/retrodeck/ES-DE/downloaded_media/snes/videos" {
		t.Fatalf("GetVideoDirectory() = %q", got)
	}
}

func TestGetGroutGamelist(t *testing.T) {
	home := setupValidConfig(t)
	expected := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "ES-DE", "gamelists", "psx", "gamelist.xml")
	if got := GetGroutGamelist("psx"); got != expected {
		t.Fatalf("GetGroutGamelist() = %q, want %q", got, expected)
	}
}

// Stub discovery and short-ID validation are covered by
// internal/scummvm.FindStub's own tests. The cases below cover only
// RetroDECK's folder-plus-stub rename policy once a stub has already been
// found.

func TestPrepareScummVMLauncherRenamesDirectoryAndStub(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Beneath a Steel Sky")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "launcher.scummvm"), []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}

	launcherDir, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ready {
		t.Fatal("ready = false, want true")
	}
	wantDir := filepath.Join(root, "Beneath a Steel Sky.scummvm")
	if launcherDir != wantDir {
		t.Fatalf("launcher dir = %q, want %q", launcherDir, wantDir)
	}
	wantStub := filepath.Join(wantDir, "Beneath a Steel Sky.scummvm")
	content, err := os.ReadFile(wantStub)
	if err != nil {
		t.Fatalf("read renamed stub: %v", err)
	}
	if string(content) != "sky" {
		t.Fatalf("stub contents = %q, want %q", content, "sky")
	}
	if _, err := os.Stat(extractDir); !os.IsNotExist(err) {
		t.Fatalf("original directory still exists: %v", err)
	}
}

func TestPrepareScummVMLauncherKeepsAlreadyCanonicalLayout(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Beneath a Steel Sky")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "Beneath a Steel Sky.scummvm"), []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}

	launcherDir, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ready {
		t.Fatal("ready = false, want true")
	}
	wantDir := filepath.Join(root, "Beneath a Steel Sky.scummvm")
	if launcherDir != wantDir {
		t.Fatalf("launcher dir = %q, want %q", launcherDir, wantDir)
	}
}

func TestPrepareScummVMLauncherNotReadyWhenNoStub(t *testing.T) {
	extractDir := t.TempDir()

	launcherDir, ready, err := PrepareScummVMLauncher(extractDir)
	if err != nil {
		t.Fatal(err)
	}
	if ready {
		t.Fatal("ready = true, want false")
	}
	if launcherDir != "" {
		t.Fatalf("launcher dir = %q, want empty", launcherDir)
	}
	if _, err := os.Stat(extractDir); err != nil {
		t.Fatalf("original directory changed: %v", err)
	}
}

func TestPrepareScummVMLauncherLeavesSourceWhenStubRenameFails(t *testing.T) {
	root := t.TempDir()
	extractDir := filepath.Join(root, "Beneath a Steel Sky")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extractDir, "launcher.scummvm"), []byte("sky"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A directory already occupying the stub's rename destination makes the
	// outer directory rename succeed but the inner stub rename fail
	// (renaming a file onto an existing directory always fails), exercising
	// the rollback path.
	if err := os.MkdirAll(filepath.Join(extractDir, "Beneath a Steel Sky.scummvm"), 0o755); err != nil {
		t.Fatal(err)
	}

	launcherDir, ready, err := PrepareScummVMLauncher(extractDir)
	if err == nil {
		t.Fatal("expected stub rename to fail")
	}
	if ready {
		t.Fatal("launcher reported ready after rename failure")
	}
	if launcherDir != "" {
		t.Fatalf("launcher dir = %q, want empty", launcherDir)
	}
	if _, err := os.Stat(extractDir); err != nil {
		t.Fatalf("source directory not rolled back: %v", err)
	}
}
