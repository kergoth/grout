package emudeckmap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePlatformsMergesOverlayPaths(t *testing.T) {
	baseXML, err := os.ReadFile(filepath.Join("testdata", "base_es_systems.xml"))
	if err != nil {
		t.Fatal(err)
	}
	overlayXML, err := os.ReadFile(filepath.Join("testdata", "emudeck_overlay.xml"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := GeneratePlatforms(baseXML, overlayXML)
	if err != nil {
		t.Fatal(err)
	}

	if len(got["psx"]) == 0 {
		t.Fatalf("psx: no paths in result")
	}
	if got["psx"][0] != "psx" {
		t.Fatalf("psx path = %q, want %q", got["psx"][0], "psx")
	}
	if len(got["arcade"]) == 0 {
		t.Fatalf("arcade: no paths in result")
	}
	if got["arcade"][0] != "model2/roms" {
		t.Fatalf("arcade path = %q, want %q", got["arcade"][0], "model2/roms")
	}
	if len(got["wiiu"]) == 0 {
		t.Fatalf("wiiu: no paths in result")
	}
	if got["wiiu"][0] != "wiiu/roms" {
		t.Fatalf("wiiu path = %q, want %q", got["wiiu"][0], "wiiu/roms")
	}
}

func TestGeneratePlatformsNormalizesESDESlug(t *testing.T) {
	baseXML := []byte(`<?xml version="1.0"?><systemList>
		<system><name>dreamcast</name><platform>dreamcast</platform><path>%ROMPATH%/dreamcast</path></system>
		<system><name>mastersystem</name><platform>mastersystem</platform><path>%ROMPATH%/mastersystem</path></system>
	</systemList>`)
	overlayXML := []byte(`<?xml version="1.0"?><systemList></systemList>`)

	got, err := GeneratePlatforms(baseXML, overlayXML)
	if err != nil {
		t.Fatal(err)
	}

	if len(got["dreamcast"]) != 0 {
		t.Fatalf("expected no 'dreamcast' key, got %v", got["dreamcast"])
	}
	if len(got["dc"]) == 0 || got["dc"][0] != "dreamcast" {
		t.Fatalf("dc paths = %v, want [dreamcast]", got["dc"])
	}
	if len(got["sms"]) == 0 || got["sms"][0] != "mastersystem" {
		t.Fatalf("sms paths = %v, want [mastersystem]", got["sms"])
	}
}

func TestGeneratePlatformsSplitsCommaSeparatedPlatforms(t *testing.T) {
	// daphne has no normalization entry; arcade passes through unchanged.
	// Both should point to the same folder.
	baseXML := []byte(`<?xml version="1.0"?><systemList>
		<system><name>daphne</name><platform>daphne, arcade</platform><path>%ROMPATH%/daphne</path></system>
	</systemList>`)
	overlayXML := []byte(`<?xml version="1.0"?><systemList></systemList>`)

	got, err := GeneratePlatforms(baseXML, overlayXML)
	if err != nil {
		t.Fatal(err)
	}

	if len(got["daphne"]) == 0 || got["daphne"][0] != "daphne" {
		t.Fatalf("daphne paths = %v, want [daphne]", got["daphne"])
	}
	if len(got["arcade"]) == 0 || got["arcade"][0] != "daphne" {
		t.Fatalf("arcade paths = %v, want [daphne]", got["arcade"])
	}
	if _, ok := got["daphne, arcade"]; ok {
		t.Fatal("raw comma key should not appear in result")
	}
}

func TestGeneratePlatformsExpandsMultiSlugMapping(t *testing.T) {
	// neogeo → both neogeomvs and neogeoaes, same folder
	baseXML := []byte(`<?xml version="1.0"?><systemList>
		<system><name>neogeo</name><platform>neogeo</platform><path>%ROMPATH%/neogeo</path></system>
	</systemList>`)
	overlayXML := []byte(`<?xml version="1.0"?><systemList></systemList>`)

	got, err := GeneratePlatforms(baseXML, overlayXML)
	if err != nil {
		t.Fatal(err)
	}

	if len(got["neogeomvs"]) == 0 || got["neogeomvs"][0] != "neogeo" {
		t.Fatalf("neogeomvs paths = %v, want [neogeo]", got["neogeomvs"])
	}
	if len(got["neogeoaes"]) == 0 || got["neogeoaes"][0] != "neogeo" {
		t.Fatalf("neogeoaes paths = %v, want [neogeo]", got["neogeoaes"])
	}
	if _, ok := got["neogeo"]; ok {
		t.Fatal("raw ES-DE slug 'neogeo' should not appear as a key")
	}
}

func TestGeneratePlatformsSkipsEntriesWithoutPlatform(t *testing.T) {
	baseXML := []byte(`<?xml version="1.0"?><systemList><system><name>ports</name><path>%ROMPATH%/ports</path></system></systemList>`)
	overlayXML := []byte(`<?xml version="1.0"?><systemList></systemList>`)

	got, err := GeneratePlatforms(baseXML, overlayXML)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 0 {
		t.Fatalf("unexpected platforms: %#v", got)
	}
}
