package input

// Chord is a set of alternative modifier keys plus a key; it is active when any
// modifier and the key are held.
type Chord struct {
	Modifiers []Key
	Key       Key
	Actions   []Action
}

// Bindings maps physical inputs to actions. It is plain data so it can be rebound and persisted.
type Bindings struct {
	Keys          map[Key][]Action
	PadButtons    map[PadButton][]Action
	Chords        []Chord
	StickDeadZone float64
}

// DefaultBindings returns the default bindings (R-03).
func DefaultBindings() Bindings {
	return Bindings{
		Keys: map[Key][]Action{
			KeyArrowLeft:  {Left},
			KeyArrowRight: {Right},
			KeyArrowUp:    {Jump},
			KeyZ:          {Jump},
			KeyX:          {Fire},
			KeySpace:      {Fire},
			KeyF11:        {ToggleFullscreen},
			KeyF3:         {ToggleDebug},
			KeyEscape:     {Quit},
		},
		Chords: []Chord{{
			Modifiers: []Key{KeyAltLeft, KeyAltRight},
			Key:       KeyEnter,
			Actions:   []Action{ToggleFullscreen},
		}},
		PadButtons: map[PadButton][]Action{
			LeftLeft:    {Left},
			LeftRight:   {Right},
			LeftTop:     {Jump},
			RightBottom: {Jump},
			RightRight:  {Fire},
			RightLeft:   {Fire},
		},
		StickDeadZone: 0.5,
	}
}
