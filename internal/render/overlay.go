package render

import (
	"fmt"
	"strconv"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

// PadInfo describes a connected gamepad for the overlay.
type PadInfo struct {
	Name     string
	Standard bool
}

// OverlayInfo is the data shown by the debug overlay.
type OverlayInfo struct {
	TPS, FPS         float64
	Tick             uint64
	ScreenW, ScreenH int
	Held             input.Actions
	MarkerX          float64
	Pads             []PadInfo
}

// View is what the renderer needs per frame.
type View struct {
	MarkerX        float64
	Overlay        OverlayInfo
	OverlayVisible bool
}

const (
	gameActions = 1<<input.Left | 1<<input.Right | 1<<input.Jump | 1<<input.Fire
	overlayHelp = "F3: hide overlay  F11/Alt+Enter: fullscreen  Esc: quit"
)

// OverlayLines builds the debug overlay text, one string per line.
func OverlayLines(i OverlayInfo) []string {
	sc := NewLetterbox(i.ScreenW, i.ScreenH, 1920, 1080).Scene
	actions := (i.Held & gameActions).String()
	if actions == "" {
		actions = "-"
	}
	lines := []string{
		fmt.Sprintf("TPS %.1f  FPS %.1f", i.TPS, i.FPS),
		"Tick " + strconv.FormatUint(i.Tick, 10),
		fmt.Sprintf("Window %dx%d  Scene %dx%d at %d,%d", i.ScreenW, i.ScreenH, sc.W, sc.H, sc.X, sc.Y),
		"Actions: " + actions,
		fmt.Sprintf("Marker X %.1f", i.MarkerX),
		fmt.Sprintf("Gamepads: %d", len(i.Pads)),
	}
	for n, p := range i.Pads {
		kind := "unsupported"
		if p.Standard {
			kind = "standard"
		}
		lines = append(lines, fmt.Sprintf("Pad %d: %s (%s)", n, p.Name, kind))
	}
	return append(lines, overlayHelp)
}
