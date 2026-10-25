package input

import "fmt"

// Key is a keyboard key, independent of the engine.
type Key uint8

// Keys known to the input layer.
const (
	KeyArrowLeft Key = iota
	KeyArrowRight
	KeyArrowUp
	KeyArrowDown
	KeyZ
	KeyX
	KeyA
	KeySpace
	KeyEnter
	KeyAltLeft
	KeyAltRight
	KeyF3
	KeyF11
	KeyEscape
)

var keyNames = [...]string{"ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown", "Z", "X", "A",
	"Space", "Enter", "AltLeft", "AltRight", "F3", "F11", "Escape"}

// ParseKey returns the key with the given name.
func ParseKey(name string) (Key, error) {
	for i, n := range keyNames {
		if n == name {
			return Key(i), nil
		}
	}
	return 0, fmt.Errorf("unknown key %q", name)
}

// PadButton is a button of a standard-layout gamepad.
type PadButton uint8

// Standard gamepad buttons, mirroring Ebitengine's StandardGamepadButton names.
const (
	RightBottom PadButton = iota
	RightRight
	RightLeft
	RightTop
	FrontTopLeft
	FrontTopRight
	FrontBottomLeft
	FrontBottomRight
	CenterLeft
	CenterRight
	LeftStick
	RightStick
	LeftTop
	LeftBottom
	LeftLeft
	LeftRight
	CenterCenter
)

var padButtonNames = [...]string{"RightBottom", "RightRight", "RightLeft", "RightTop",
	"FrontTopLeft", "FrontTopRight", "FrontBottomLeft", "FrontBottomRight", "CenterLeft",
	"CenterRight", "LeftStick", "RightStick", "LeftTop", "LeftBottom", "LeftLeft", "LeftRight",
	"CenterCenter"}

// ParsePadButton returns the pad button with the given name.
func ParsePadButton(name string) (PadButton, error) {
	for i, n := range padButtonNames {
		if n == name {
			return PadButton(i), nil
		}
	}
	return 0, fmt.Errorf("unknown gamepad button %q", name)
}

// Pad is the raw state of one gamepad.
type Pad struct {
	Name       string
	Standard   bool
	Buttons    map[PadButton]bool
	LeftStickX float64
	// RawButtons is only filled for non-standard pads.
	RawButtons []bool
}

// Raw is a snapshot of all input devices for one tick.
type Raw struct {
	Keys map[Key]bool
	Pads []Pad
}
