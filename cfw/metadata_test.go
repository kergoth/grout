package cfw

import (
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
		{
			name:         "retrodeck ports gamelist",
			cfw:          RetroDeck,
			expectedPath: filepath.Join(".var", "app", "net.retrodeck.retrodeck", "config", "ES-DE", "gamelists", "ports", "gamelist.xml"),
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
