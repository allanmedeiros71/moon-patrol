// Package inputebiten reads raw input through Ebitengine.
package inputebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

var ebitenKeys = map[input.Key]ebiten.Key{
	input.KeyArrowLeft: ebiten.KeyArrowLeft, input.KeyArrowRight: ebiten.KeyArrowRight,
	input.KeyArrowUp: ebiten.KeyArrowUp, input.KeyArrowDown: ebiten.KeyArrowDown,
	input.KeyZ: ebiten.KeyZ, input.KeyX: ebiten.KeyX, input.KeyA: ebiten.KeyA, input.KeySpace: ebiten.KeySpace,
	input.KeyEnter: ebiten.KeyEnter, input.KeyAltLeft: ebiten.KeyAltLeft, input.KeyAltRight: ebiten.KeyAltRight,
	input.KeyF3: ebiten.KeyF3, input.KeyF11: ebiten.KeyF11, input.KeyEscape: ebiten.KeyEscape,
}

var ebitenPadButtons = map[input.PadButton]ebiten.StandardGamepadButton{
	input.RightBottom: ebiten.StandardGamepadButtonRightBottom, input.RightRight: ebiten.StandardGamepadButtonRightRight,
	input.RightLeft: ebiten.StandardGamepadButtonRightLeft, input.RightTop: ebiten.StandardGamepadButtonRightTop,
	input.FrontTopLeft: ebiten.StandardGamepadButtonFrontTopLeft, input.FrontTopRight: ebiten.StandardGamepadButtonFrontTopRight,
	input.FrontBottomLeft: ebiten.StandardGamepadButtonFrontBottomLeft, input.FrontBottomRight: ebiten.StandardGamepadButtonFrontBottomRight,
	input.CenterLeft: ebiten.StandardGamepadButtonCenterLeft, input.CenterRight: ebiten.StandardGamepadButtonCenterRight,
	input.LeftStick: ebiten.StandardGamepadButtonLeftStick, input.RightStick: ebiten.StandardGamepadButtonRightStick,
	input.LeftTop: ebiten.StandardGamepadButtonLeftTop, input.LeftBottom: ebiten.StandardGamepadButtonLeftBottom,
	input.LeftLeft: ebiten.StandardGamepadButtonLeftLeft, input.LeftRight: ebiten.StandardGamepadButtonLeftRight,
	input.CenterCenter: ebiten.StandardGamepadButtonCenterCenter,
}

// EbitenPoller reads the keyboard and gamepads through Ebitengine.
type EbitenPoller struct{ ids []ebiten.GamepadID }

// Poll returns the raw input state for this tick.
func (p *EbitenPoller) Poll() input.Raw {
	r := input.Raw{Keys: make(map[input.Key]bool, len(ebitenKeys))}
	for k, ek := range ebitenKeys {
		r.Keys[k] = ebiten.IsKeyPressed(ek)
	}
	p.ids = ebiten.AppendGamepadIDs(p.ids[:0])
	for _, id := range p.ids {
		r.Pads = append(r.Pads, pollPad(id))
	}
	return r
}

func pollPad(id ebiten.GamepadID) input.Pad {
	pad := input.Pad{Name: ebiten.GamepadName(id), Standard: ebiten.IsStandardGamepadLayoutAvailable(id)}
	if !pad.Standard {
		n := ebiten.GamepadButtonCount(id)
		pad.RawButtons = make([]bool, n)
		for i := range n {
			pad.RawButtons[i] = ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(i))
		}
		return pad
	}
	pad.Buttons = make(map[input.PadButton]bool, len(ebitenPadButtons))
	for b, eb := range ebitenPadButtons {
		pad.Buttons[b] = ebiten.IsStandardGamepadButtonPressed(id, eb)
	}
	pad.LeftStickX = ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
	return pad
}
