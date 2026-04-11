package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"grout/internal/emudeckmap"
)

func main() {
	esSystems := flag.String("es-systems", "", "path to upstream ES-DE es_systems.xml")
	overlay := flag.String("overlay", "", "path to EmuDeck custom_systems es_systems.xml")
	output := flag.String("output", "cfw/emudeck/data/platforms.json", "output path")
	flag.Parse()

	baseXML, err := os.ReadFile(*esSystems)
	if err != nil {
		log.Fatal(err)
	}
	overlayXML, err := os.ReadFile(*overlay)
	if err != nil {
		log.Fatal(err)
	}

	platforms, err := emudeckmap.GeneratePlatforms(baseXML, overlayXML)
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.MarshalIndent(platforms, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*output, append(body, '\n'), 0644); err != nil {
		log.Fatal(err)
	}
}
