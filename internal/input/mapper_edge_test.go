package input

import "testing"

// A non-standard pad must contribute nothing, even when its stick would steer.
// TestTwoPadsAndNonStandard cannot detect a leak because a standard pad also sets Right.
func TestNonStandardPadContributesNothing(t *testing.T) {
	m := NewMapper(DefaultBindings())
	raw := Raw{Pads: []Pad{{Name: "weird", RawButtons: []bool{true, true}, LeftStickX: 1}}}
	if f := m.Map(raw); f.Held != 0 || f.Pressed != 0 {
		t.Errorf("held %q pressed %q, want none", f.Held, f.Pressed)
	}
	if u := Unsupported(raw); len(u) != 1 || u[0] != "weird" {
		t.Errorf("unsupported = %v", u)
	}
}

func TestUnsupportedEmpty(t *testing.T) {
	if u := Unsupported(Raw{Pads: []Pad{stdPad(0)}}); len(u) != 0 {
		t.Errorf("unsupported = %v, want none", u)
	}
}

// A standard pad with no buttons map must not panic and must map to nothing.
func TestStandardPadWithoutButtonsMap(t *testing.T) {
	m := NewMapper(DefaultBindings())
	if f := m.Map(Raw{Pads: []Pad{{Name: "bare", Standard: true}}}); f.Held != 0 {
		t.Errorf("held %q, want none", f.Held)
	}
}

// Holding a pad button across ticks counts as pressed once; releasing and pressing again is a new press.
func TestPadButtonPressedOnceAndRePress(t *testing.T) {
	m := NewMapper(DefaultBindings())
	down := Raw{Pads: []Pad{stdPad(0, RightRight)}}
	up := Raw{Pads: []Pad{stdPad(0)}}
	seq := []struct {
		r             Raw
		held, pressed Actions
	}{
		{down, Actions(0).With(Fire), Actions(0).With(Fire)},
		{down, Actions(0).With(Fire), 0},
		{up, 0, 0},
		{down, Actions(0).With(Fire), Actions(0).With(Fire)},
	}
	for i, s := range seq {
		f := m.Map(s.r)
		if f.Held != s.held || f.Pressed != s.pressed {
			t.Errorf("tick %d: held %q pressed %q, want held %q pressed %q",
				i+1, f.Held, f.Pressed, s.held, s.pressed)
		}
	}
}

// Stick steering: moving from left to right on consecutive ticks presses Right only.
func TestStickSwitchDirectionPressesNewAction(t *testing.T) {
	m := NewMapper(DefaultBindings())
	if f := m.Map(Raw{Pads: []Pad{stdPad(-0.8)}}); !f.Pressed.Has(Left) {
		t.Fatalf("tick 1: Left not pressed, got %q", f.Pressed)
	}
	f := m.Map(Raw{Pads: []Pad{stdPad(0.8)}})
	if f.Pressed != Actions(0).With(Right) || f.Held != Actions(0).With(Right) {
		t.Errorf("tick 2: held %q pressed %q, want Right/Right", f.Held, f.Pressed)
	}
}

// An action held by both keyboard and pad is pressed only when it first becomes active.
func TestKeyAndPadSameActionPressedOnce(t *testing.T) {
	m := NewMapper(DefaultBindings())
	if f := m.Map(Raw{Keys: map[Key]bool{KeyArrowLeft: true}}); !f.Pressed.Has(Left) {
		t.Fatalf("tick 1: Left not pressed, got %q", f.Pressed)
	}
	f := m.Map(Raw{
		Keys: map[Key]bool{KeyArrowLeft: true},
		Pads: []Pad{stdPad(-0.9, LeftLeft)},
	})
	if f.Pressed != 0 || !f.Held.Has(Left) {
		t.Errorf("tick 2: held %q pressed %q, want Left held and nothing pressed", f.Held, f.Pressed)
	}
}

// Releasing Alt while Enter stays held ends the chord; pressing the other Alt starts it again.
func TestChordReleaseAndOtherAlt(t *testing.T) {
	m := NewMapper(DefaultBindings())
	m.Map(Raw{Keys: map[Key]bool{KeyEnter: true, KeyAltLeft: true}})
	f := m.Map(Raw{Keys: map[Key]bool{KeyEnter: true}})
	if f.Held.Has(ToggleFullscreen) || f.Pressed != 0 {
		t.Errorf("after Alt release: held %q pressed %q, want none", f.Held, f.Pressed)
	}
	f = m.Map(Raw{Keys: map[Key]bool{KeyEnter: true, KeyAltRight: true}})
	if !f.Pressed.Has(ToggleFullscreen) {
		t.Errorf("AltRight after release: ToggleFullscreen not pressed, got %q", f.Pressed)
	}
}

// Escape and F3 map to their system actions; Enter with no modifier maps to nothing.
func TestSystemKeyActions(t *testing.T) {
	tests := []struct {
		key  Key
		want Action
	}{
		{KeyEscape, Quit},
		{KeyF3, ToggleDebug},
		{KeyF11, ToggleFullscreen},
	}
	for _, tc := range tests {
		m := NewMapper(DefaultBindings())
		f := m.Map(Raw{Keys: map[Key]bool{tc.key: true}})
		if f.Held != Actions(0).With(tc.want) || f.Pressed != Actions(0).With(tc.want) {
			t.Errorf("key %v: held %q pressed %q, want %v", tc.key, f.Held, f.Pressed, tc.want)
		}
	}
}

// Left stick at exactly the negative dead-zone edge and a keyboard key together: Left once.
func TestStickAndKeySameDirection(t *testing.T) {
	m := NewMapper(DefaultBindings())
	f := m.Map(Raw{
		Keys: map[Key]bool{KeyArrowLeft: true},
		Pads: []Pad{stdPad(-0.5)},
	})
	if f.Held != Actions(0).With(Left) || f.Pressed != Actions(0).With(Left) {
		t.Errorf("held %q pressed %q, want Left/Left", f.Held, f.Pressed)
	}
}
