package fileutil

import "testing"

func TestUsesArchiveAsRom(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		// Core set
		{"dos", true},
		{"arcade", true},
		{"mame", true},
		{"fbneo", true},
		{"neogeo", true},
		// Aliases RomM may report
		{"fba", true},
		{"cps1", true},
		{"cps2", true},
		{"cps3", true},
		{"neogeoaes", true},
		{"neogeomvs", true},
		{"neogeocd", true},
		{"naomi", true},
		{"naomi2", true},
		{"atomiswave", true},
		// Case-insensitivity
		{"DOS", true},
		{"MAME", true},
		// Non-members
		{"scummvm", false},
		{"snes", false},
		{"psx", false},
		{"switch", false},
		// Empty must never skip extraction
		{"", false},
	}
	for _, c := range cases {
		if got := UsesArchiveAsRom(c.slug); got != c.want {
			t.Errorf("UsesArchiveAsRom(%q) = %v, want %v", c.slug, got, c.want)
		}
	}
}
