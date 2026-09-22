package cfw

import (
	"grout/internal/gamelist"
	"grout/romm"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddGroutToGamelist(t *testing.T) {
	tests := []struct {
		name         string
		cfw          CFW
		expectedPath string
	}{
		{
			name:         "emudeck ports gamelist",
			cfw:          EmuDeck,
			expectedPath: filepath.Join("ES-DE", "gamelists", "ports", "gamelist.xml"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)

			target := filepath.Join(home, tc.expectedPath)
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				t.Fatalf("mkdir gamelist dir: %v", err)
			}

			AddGroutToGamelist(tc.cfw)
			data, err := os.ReadFile(target)
			if err != nil {
				t.Fatalf("read gamelist: %v", err)
			}
			content := string(data)
			if !strings.Contains(content, "<path>./Grout.sh</path>") {
				t.Fatalf("missing grout launcher path in gamelist: %s", content)
			}
		})
	}
}

// Grout installs on RetroDECK as a non-Steam game (see
// docs/getting-started/install-retrodeck.md), not as an ES-DE port, so it
// must not register itself in RetroDECK's ES-DE gamelist.xml.
func TestAddGroutToGamelistRetroDeckIsNoOp(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	AddGroutToGamelist(RetroDeck)

	gamelistPath := filepath.Join(home, ".var", "app", "net.retrodeck.retrodeck", "config", "ES-DE", "gamelists", "desktop", "gamelist.xml")
	if _, err := os.Stat(gamelistPath); !os.IsNotExist(err) {
		t.Fatalf("expected no gamelist to be written for RetroDeck, got err=%v", err)
	}
}

func TestFillGamesMetadataWritesESDEGamelistAndSchedulesRestart(t *testing.T) {
	for _, cfwName := range []string{"EMUDECK", "RETRODECK"} {
		t.Run(cfwName, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("CFW", cfwName)

			gamelistPath := filepath.Join(t.TempDir(), "gamelist.xml")
			entries := []gamelist.RomGameEntry{
				{
					Game:         &romm.Rom{Name: "Beneath a Steel Sky"},
					RomDirectory: filepath.Dir(gamelistPath),
					GamelistPath: gamelistPath,
				},
			}

			FillGamesMetadata(entries)

			data, err := os.ReadFile(gamelistPath)
			if err != nil {
				t.Fatalf("read gamelist: %v", err)
			}
			if !strings.Contains(string(data), "Beneath a Steel Sky") {
				t.Fatalf("missing game entry in gamelist: %s", data)
			}
			if _, err := os.Stat("es_restart_request"); err != nil {
				t.Fatalf("expected ES restart to be scheduled for %s: %v", cfwName, err)
			}
		})
	}
}
