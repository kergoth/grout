package esde

import (
	"embed"
	"fmt"
)

//go:embed input_mappings/*.json
var embeddedInputMappings embed.FS

// GetInputMappingBytes returns the controller mapping for the Steam Deck.
//
// Steam Input presents every pad to Grout as one Xbox-style controller, so a
// single mapping covers the Deck's own controls and external controllers. The
// toolkit's default does not suit it: it swaps A and B, which feels backwards
// on this layout, and it binds nothing but buttons. The toolkit drops any axis
// that has no binding, so the triggers are dead, and the screen that lets a
// user remap buttons can never see L2. Here the triggers are L2 and R2, and a
// click of the left stick is a second L2 for the BIOS shortcut. The sticks stay
// unbound on purpose: the toolkit passes every motion event of a bound axis, so
// a stick's constant jitter would push button presses out of the event queue.
func GetInputMappingBytes() ([]byte, error) {
	const filename = "input_mappings/steam-input.json"

	data, err := embeddedInputMappings.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded input mapping %s: %w", filename, err)
	}
	return data, nil
}
