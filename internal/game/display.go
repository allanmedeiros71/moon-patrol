// Package game holds the game states; this file is the pure display-mode controller.
package game

import "github.com/allanmedeiros/moonpatrol/internal/input"

// Windowed size used when not fullscreen (D-14).
const (
	WindowedWidth  = 1280
	WindowedHeight = 720
)

// Options are the start-up options that affect the display.
type Options struct {
	// Windowed starts in a window instead of fullscreen (-windowed).
	Windowed bool
}

// Change tells the caller which window API calls to make after Apply.
type Change struct {
	FullscreenChanged bool
}

// Display decides the display mode, overlay visibility and quit request from input frames.
type Display struct {
	fullscreen bool
	debug      bool
	quit       bool
}

// NewDisplay returns a controller that starts fullscreen unless o.Windowed, with the overlay visible.
func NewDisplay(o Options) *Display {
	return &Display{fullscreen: !o.Windowed, debug: true}
}

// Fullscreen reports whether the game should be fullscreen.
func (d *Display) Fullscreen() bool { return d.fullscreen }

// WindowSize returns the window size used when windowed.
func (d *Display) WindowSize() (int, int) { return WindowedWidth, WindowedHeight }

// DebugVisible reports whether the debug overlay is shown.
func (d *Display) DebugVisible() bool { return d.debug }

// QuitRequested reports whether the player asked to quit.
func (d *Display) QuitRequested() bool { return d.quit }

// Apply reacts to the actions pressed this tick; held actions alone do nothing.
func (d *Display) Apply(f input.Frame) Change {
	var c Change
	if f.Pressed.Has(input.ToggleFullscreen) {
		d.fullscreen = !d.fullscreen
		c.FullscreenChanged = true
	}
	if f.Pressed.Has(input.ToggleDebug) {
		d.debug = !d.debug
	}
	if f.Pressed.Has(input.Quit) {
		d.quit = true
	}
	return c
}
