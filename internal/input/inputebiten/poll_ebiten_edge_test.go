package inputebiten

import "testing"

// Without a running game and without a gamepad, Poll reports every mapped key as released
// and no pads; this exercises the Ebitengine calls headlessly.
func TestPollWithoutDevicesReportsReleasedKeys(t *testing.T) {
	var p EbitenPoller
	r := p.Poll()
	if len(r.Keys) != len(ebitenKeys) {
		t.Fatalf("Keys has %d entries, want %d", len(r.Keys), len(ebitenKeys))
	}
	for k := range ebitenKeys {
		if pressed, ok := r.Keys[k]; !ok || pressed {
			t.Errorf("key %d: present=%v pressed=%v", k, ok, pressed)
		}
	}
	if len(r.Pads) != 0 {
		t.Errorf("Pads = %d, want 0", len(r.Pads))
	}
}

// GamepadID 0 is not connected: the non-standard branch reports no raw buttons.
func TestPollPadNotConnected(t *testing.T) {
	pad := pollPad(0)
	if pad.Standard {
		t.Error("disconnected pad reported a standard layout")
	}
	if len(pad.RawButtons) != 0 || pad.Buttons != nil {
		t.Errorf("pad = %+v", pad)
	}
}
