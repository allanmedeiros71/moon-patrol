package input

import "github.com/hajimehoshi/ebiten/v2"

var ebitenKeys = map[Key]ebiten.Key{
	KeyArrowLeft: ebiten.KeyArrowLeft, KeyArrowRight: ebiten.KeyArrowRight,
	KeyArrowUp: ebiten.KeyArrowUp, KeyArrowDown: ebiten.KeyArrowDown,
	KeyZ: ebiten.KeyZ, KeyX: ebiten.KeyX, KeyA: ebiten.KeyA, KeySpace: ebiten.KeySpace,
	KeyEnter: ebiten.KeyEnter, KeyAltLeft: ebiten.KeyAltLeft, KeyAltRight: ebiten.KeyAltRight,
	KeyF3: ebiten.KeyF3, KeyF11: ebiten.KeyF11, KeyEscape: ebiten.KeyEscape,
}

var ebitenPadButtons = map[PadButton]ebiten.StandardGamepadButton{
	RightBottom: ebiten.StandardGamepadButtonRightBottom, RightRight: ebiten.StandardGamepadButtonRightRight,
	RightLeft: ebiten.StandardGamepadButtonRightLeft, RightTop: ebiten.StandardGamepadButtonRightTop,
	FrontTopLeft: ebiten.StandardGamepadButtonFrontTopLeft, FrontTopRight: ebiten.StandardGamepadButtonFrontTopRight,
	FrontBottomLeft: ebiten.StandardGamepadButtonFrontBottomLeft, FrontBottomRight: ebiten.StandardGamepadButtonFrontBottomRight,
	CenterLeft: ebiten.StandardGamepadButtonCenterLeft, CenterRight: ebiten.StandardGamepadButtonCenterRight,
	LeftStick: ebiten.StandardGamepadButtonLeftStick, RightStick: ebiten.StandardGamepadButtonRightStick,
	LeftTop: ebiten.StandardGamepadButtonLeftTop, LeftBottom: ebiten.StandardGamepadButtonLeftBottom,
	LeftLeft: ebiten.StandardGamepadButtonLeftLeft, LeftRight: ebiten.StandardGamepadButtonLeftRight,
	CenterCenter: ebiten.StandardGamepadButtonCenterCenter,
}

// EbitenPoller reads the keyboard and gamepads through Ebitengine.
type EbitenPoller struct{ ids []ebiten.GamepadID }

// Poll returns the raw input state for this tick.
func (p *EbitenPoller) Poll() Raw {
	r := Raw{Keys: make(map[Key]bool, len(ebitenKeys))}
	for k, ek := range ebitenKeys {
		r.Keys[k] = ebiten.IsKeyPressed(ek)
	}
	p.ids = ebiten.AppendGamepadIDs(p.ids[:0])
	for _, id := range p.ids {
		r.Pads = append(r.Pads, pollPad(id))
	}
	return r
}

func pollPad(id ebiten.GamepadID) Pad {
	pad := Pad{Name: ebiten.GamepadName(id), Standard: ebiten.IsStandardGamepadLayoutAvailable(id)}
	if !pad.Standard {
		n := ebiten.GamepadButtonCount(id)
		pad.RawButtons = make([]bool, n)
		for i := range n {
			pad.RawButtons[i] = ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(i))
		}
		return pad
	}
	pad.Buttons = make(map[PadButton]bool, len(ebitenPadButtons))
	for b, eb := range ebitenPadButtons {
		pad.Buttons[b] = ebiten.IsStandardGamepadButtonPressed(id, eb)
	}
	pad.LeftStickX = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
	return pad
}
