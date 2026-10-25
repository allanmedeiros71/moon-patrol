package input

import "testing"

func TestActionsString(t *testing.T) {
	if got := Actions(0).String(); got != "" {
		t.Errorf("empty = %q", got)
	}
	s := Actions(0).With(Fire).With(Right).With(Jump)
	if got := s.String(); got != "Right Jump Fire" {
		t.Errorf("got %q", got)
	}
	if got := s.With(Quit).With(ToggleDebug).With(ToggleFullscreen).With(Left).String(); got !=
		"Left Right Jump Fire ToggleFullscreen ToggleDebug Quit" {
		t.Errorf("got %q", got)
	}
	if s.Has(Left) || !s.Has(Fire) {
		t.Errorf("Has wrong for %v", s)
	}
}

func TestParse(t *testing.T) {
	for i, n := range keyNames {
		if k, err := ParseKey(n); err != nil || k != Key(i) {
			t.Errorf("ParseKey(%q) = %v, %v", n, k, err)
		}
	}
	for i, n := range padButtonNames {
		if b, err := ParsePadButton(n); err != nil || b != PadButton(i) {
			t.Errorf("ParsePadButton(%q) = %v, %v", n, b, err)
		}
	}
	if _, err := ParseKey("Nope"); err == nil {
		t.Error("ParseKey unknown: want error")
	}
	if _, err := ParsePadButton("Nope"); err == nil {
		t.Error("ParsePadButton unknown: want error")
	}
}

func stdPad(x float64, btns ...PadButton) Pad {
	p := Pad{Name: "pad", Standard: true, Buttons: map[PadButton]bool{}, LeftStickX: x}
	for _, b := range btns {
		p.Buttons[b] = true
	}
	return p
}

func TestDeadZone(t *testing.T) {
	tests := []struct {
		x    float64
		want Actions
	}{
		{-0.5, Actions(0).With(Left)},
		{-0.49, 0},
		{0.49, 0},
		{0.5, Actions(0).With(Right)},
	}
	for _, tc := range tests {
		m := NewMapper(DefaultBindings())
		if got := m.Map(Raw{Pads: []Pad{stdPad(tc.x)}}).Held; got != tc.want {
			t.Errorf("x=%v: got %q want %q", tc.x, got, tc.want)
		}
	}
}

func TestChord(t *testing.T) {
	for _, alt := range []Key{KeyAltLeft, KeyAltRight} {
		m := NewMapper(DefaultBindings())
		f := m.Map(Raw{Keys: map[Key]bool{alt: true, KeyEnter: true}})
		if !f.Pressed.Has(ToggleFullscreen) {
			t.Errorf("alt %v: not pressed", alt)
		}
	}
	m := NewMapper(DefaultBindings())
	if f := m.Map(Raw{Keys: map[Key]bool{KeyEnter: true}}); f.Held != 0 {
		t.Errorf("enter alone: %q", f.Held)
	}
	f := m.Map(Raw{Keys: map[Key]bool{KeyEnter: true, KeyAltLeft: true}})
	if !f.Pressed.Has(ToggleFullscreen) {
		t.Error("alt after enter: want pressed")
	}
	f = m.Map(Raw{Keys: map[Key]bool{KeyEnter: true, KeyAltLeft: true}})
	if f.Pressed != 0 || !f.Held.Has(ToggleFullscreen) {
		t.Errorf("held chord repeated: %+v", f)
	}
}

func TestTwoPadsAndNonStandard(t *testing.T) {
	m := NewMapper(DefaultBindings())
	raw := Raw{Pads: []Pad{
		stdPad(0, LeftLeft),
		stdPad(0.9, RightRight),
		{Name: "weird", RawButtons: []bool{true}, LeftStickX: 1},
	}}
	f := m.Map(raw)
	if want := Actions(0).With(Left).With(Right).With(Fire); f.Held != want {
		t.Errorf("held %q want %q", f.Held, want)
	}
	if u := Unsupported(raw); len(u) != 1 || u[0] != "weird" {
		t.Errorf("unsupported = %v", u)
	}
}

func TestPressedAfterRelease(t *testing.T) {
	m := NewMapper(DefaultBindings())
	down := Raw{Keys: map[Key]bool{KeyX: true}}
	seq := []struct {
		r    Raw
		want Actions
	}{
		{down, Actions(0).With(Fire)},
		{down, 0},
		{Raw{}, 0},
		{down, Actions(0).With(Fire)},
	}
	for i, s := range seq {
		if got := m.Map(s.r).Pressed; got != s.want {
			t.Errorf("tick %d: got %q want %q", i+1, got, s.want)
		}
	}
}
