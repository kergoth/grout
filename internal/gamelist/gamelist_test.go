package gamelist

import "testing"

func TestParseAcceptsAlternativeEmulatorBeforeGameList(t *testing.T) {
	raw := []byte(`<?xml version="1.0"?><alternativeEmulator label="DuckStation">duckstation</alternativeEmulator><gameList><game><name>Game</name></game></gameList>`)
	gl := New()
	if err := gl.Parse(raw); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !gl.Contains(NameElement, "Game") {
		t.Fatal("expected gameList entry to remain accessible")
	}
}
