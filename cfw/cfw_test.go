package cfw

import (
	"testing"
)

func TestEmuDeckIsEmulationStationBased(t *testing.T) {
	if !CFW("EMUDECK").IsBasedOnEmulationStation() {
		t.Fatal("EMUDECK should be treated as ES-based")
	}
}
