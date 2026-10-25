package game

import (
	"testing"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

func pressed(as ...input.Action) input.Frame {
	var s input.Actions
	for _, a := range as {
		s = s.With(a)
	}
	return input.Frame{Held: s, Pressed: s}
}

func TestStart(t *testing.T) {
	d := NewDisplay(Options{})
	if !d.Fullscreen() || !d.DebugVisible() || d.QuitRequested() {
		t.Errorf("default start wrong: %+v", d)
	}
	if w, h := d.WindowSize(); w != 1280 || h != 720 {
		t.Errorf("size = %dx%d", w, h)
	}
	if NewDisplay(Options{Windowed: true}).Fullscreen() {
		t.Error("windowed option still fullscreen")
	}
}

func TestQuit(t *testing.T) {
	d := NewDisplay(Options{})
	if c := d.Apply(pressed(input.Quit)); c.FullscreenChanged {
		t.Error("quit changed fullscreen")
	}
	if !d.QuitRequested() {
		t.Error("quit not requested")
	}
}

func TestBothToggles(t *testing.T) {
	d := NewDisplay(Options{})
	c := d.Apply(pressed(input.ToggleFullscreen, input.ToggleDebug))
	if !c.FullscreenChanged || d.Fullscreen() || d.DebugVisible() {
		t.Errorf("after both: change=%v fs=%v dbg=%v", c, d.Fullscreen(), d.DebugVisible())
	}
	d.Apply(pressed(input.ToggleFullscreen, input.ToggleDebug))
	if !d.Fullscreen() || !d.DebugVisible() {
		t.Error("second press did not restore")
	}
}

func TestHeldOnlyDoesNothing(t *testing.T) {
	d := NewDisplay(Options{})
	var s input.Actions
	s = s.With(input.ToggleFullscreen).With(input.ToggleDebug).With(input.Quit)
	c := d.Apply(input.Frame{Held: s})
	if c.FullscreenChanged || !d.Fullscreen() || !d.DebugVisible() || d.QuitRequested() {
		t.Error("held without pressed changed state")
	}
}
