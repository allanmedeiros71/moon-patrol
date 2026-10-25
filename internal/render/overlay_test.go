package render

import (
	"slices"
	"testing"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

func TestOverlayTwoPads(t *testing.T) {
	held := input.Actions(0).With(input.Fire).With(input.Left).With(input.Quit).With(input.ToggleDebug)
	got := OverlayLines(OverlayInfo{
		TPS: 0, FPS: 59.95, Tick: 7, ScreenW: 2560, ScreenH: 1080, Held: held, MarkerX: 12.34,
		Pads: []PadInfo{{"A", true}, {"B", false}},
	})
	want := []string{
		"TPS 0.0  FPS 60.0",
		"Tick 7",
		"Window 2560x1080  Scene 1920x1080 at 320,0",
		"Actions: Left Fire",
		"Marker X 12.3",
		"Gamepads: 2",
		"Pad 0: A (standard)",
		"Pad 1: B (unsupported)",
		"F3: hide overlay  F11/Alt+Enter: fullscreen  Esc: quit",
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestOverlayNoActionsNoPads(t *testing.T) {
	got := OverlayLines(OverlayInfo{ScreenW: 1920, ScreenH: 1080, Held: input.Actions(0).With(input.Quit)})
	if got[3] != "Actions: -" || got[5] != "Gamepads: 0" || len(got) != 7 {
		t.Errorf("got %q", got)
	}
}

func TestOverlayZeroFPS(t *testing.T) {
	got := OverlayLines(OverlayInfo{FPS: 0, ScreenW: 1920, ScreenH: 1080})
	if got[0] != "TPS 0.0  FPS 0.0" {
		t.Errorf("got %q", got[0])
	}
}
