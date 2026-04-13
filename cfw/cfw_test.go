package cfw

import (
	"testing"
)

func TestEmuDeckIsEmulationStationBased(t *testing.T) {
	if !CFW("EMUDECK").IsBasedOnEmulationStation() {
		t.Fatal("EMUDECK should be treated as ES-based")
	}
}

func TestArkOSAndKorikiAreEmulationStationBased(t *testing.T) {
	for _, platform := range []CFW{ArkOS, Koriki} {
		if !platform.IsBasedOnEmulationStation() {
			t.Fatalf("%s should be treated as ES-based", platform)
		}
	}
}
