package ui

import (
	"testing"

	"grout/cfw"
	"grout/romm"
)

// A first run must map every platform that has a matching folder, not just the
// first one it finds.
func TestDetectMappings_MapsEveryPlatformWithAFolder(t *testing.T) {
	platforms := []romm.Platform{
		{FSSlug: "gba"}, {FSSlug: "snes"}, {FSSlug: "psx"}, {FSSlug: "dos"}, {FSSlug: "saturn"},
	}
	directories := []string{"gba", "snes", "psx", "dos"}
	input := PlatformMappingInput{CFW: cfw.Knulli}

	got := (&PlatformMappingScreen{}).detectMappings(platforms, directories, input)

	for _, slug := range []string{"gba", "snes", "psx", "dos"} {
		if got[slug].RelativePath != slug {
			t.Errorf("%s maps to %q, want %q", slug, got[slug].RelativePath, slug)
		}
	}
	if m, ok := got["saturn"]; ok {
		t.Errorf("saturn has no folder but maps to %q", m.RelativePath)
	}
}
