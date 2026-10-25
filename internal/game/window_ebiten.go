package game

import "github.com/hajimehoshi/ebiten/v2"

// EbitenWindow implements Window with the Ebitengine window API.
type EbitenWindow struct{}

// SetFullscreen switches fullscreen mode.
func (EbitenWindow) SetFullscreen(on bool) { ebiten.SetFullscreen(on) }

// SetWindowSize sets the window size in device-independent pixels.
func (EbitenWindow) SetWindowSize(w, h int) { ebiten.SetWindowSize(w, h) }

// DeviceScaleFactor returns the monitor scale factor.
func (EbitenWindow) DeviceScaleFactor() float64 { return ebiten.Monitor().DeviceScaleFactor() }
