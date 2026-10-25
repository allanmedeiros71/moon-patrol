package input

// Frame is the result of mapping one tick of raw input.
type Frame struct {
	// Held holds the actions active this tick.
	Held Actions
	// Pressed holds the actions that became active this tick.
	Pressed Actions
}

// Mapper turns raw snapshots into frames, remembering the previous tick.
type Mapper struct {
	b    Bindings
	prev Actions
}

// NewMapper returns a Mapper using the given bindings.
func NewMapper(b Bindings) *Mapper { return &Mapper{b: b} }

// Map computes the held and pressed actions for r.
func (m *Mapper) Map(r Raw) Frame {
	var held Actions
	add := func(as []Action) {
		for _, a := range as {
			held = held.With(a)
		}
	}
	for k, as := range m.b.Keys {
		if r.Keys[k] {
			add(as)
		}
	}
	for _, c := range m.b.Chords {
		if r.Keys[c.Key] && anyKey(r.Keys, c.Modifiers) {
			add(c.Actions)
		}
	}
	for _, p := range r.Pads {
		if p.Standard {
			held |= m.mapPad(p)
		}
	}
	f := Frame{Held: held, Pressed: held &^ m.prev}
	m.prev = held
	return f
}

func (m *Mapper) mapPad(p Pad) Actions {
	var held Actions
	for b, as := range m.b.PadButtons {
		if p.Buttons[b] {
			for _, a := range as {
				held = held.With(a)
			}
		}
	}
	if p.LeftStickX <= -m.b.StickDeadZone {
		held = held.With(Left)
	}
	if p.LeftStickX >= m.b.StickDeadZone {
		held = held.With(Right)
	}
	return held
}

func anyKey(keys map[Key]bool, ks []Key) bool {
	for _, k := range ks {
		if keys[k] {
			return true
		}
	}
	return false
}

// Unsupported returns the names of the pads without the standard layout.
func Unsupported(r Raw) []string {
	var names []string
	for _, p := range r.Pads {
		if !p.Standard {
			names = append(names, p.Name)
		}
	}
	return names
}
