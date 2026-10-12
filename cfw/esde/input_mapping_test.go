package esde

import (
	"encoding/json"
	"testing"
)

// The toolkit drops events from any axis the mapping does not bind, so a
// missing trigger binding silently kills it. It passes every motion event of a
// bound axis, whatever its value, so binding a stick lets its constant jitter
// crowd button presses out of the event queue.
func TestInputMappingBindsTheAxes(t *testing.T) {
	data, err := GetInputMappingBytes()
	if err != nil {
		t.Fatal(err)
	}

	var mapping struct {
		Axes map[string]json.RawMessage `json:"joystick_axis_map"`
	}
	if err := json.Unmarshal(data, &mapping); err != nil {
		t.Fatalf("mapping is not valid JSON: %v", err)
	}

	for _, axis := range []string{"4", "5"} {
		if _, ok := mapping.Axes[axis]; !ok {
			t.Errorf("axis %s is not bound, so the toolkit would drop it", axis)
		}
	}
	for axis := range mapping.Axes {
		switch axis {
		case "4", "5":
		default:
			t.Errorf("axis %s is bound but should not be: stray axes flood the event queue", axis)
		}
	}
}
