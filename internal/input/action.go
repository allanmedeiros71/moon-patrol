// Package input maps raw keyboard and gamepad state to game actions. It is pure:
// it does not import the engine.
package input

import "strings"

// Action is a game or system action produced by the input layer.
type Action uint16

// Actions, in the order used by Actions.String.
const (
	Left Action = iota
	Right
	Jump
	Fire
	ToggleFullscreen
	ToggleDebug
	Quit
)

var actionNames = [...]string{"Left", "Right", "Jump", "Fire", "ToggleFullscreen", "ToggleDebug", "Quit"}

// Actions is a bitset of Action values.
type Actions uint16

// Has reports whether a is in the set.
func (s Actions) Has(a Action) bool { return s&(1<<a) != 0 }

// With returns the set with a added.
func (s Actions) With(a Action) Actions { return s | 1<<a }

// String lists the names of the set actions in constant order, separated by a space.
func (s Actions) String() string {
	var names []string
	for i, n := range actionNames {
		if s.Has(Action(i)) {
			names = append(names, n)
		}
	}
	return strings.Join(names, " ")
}
