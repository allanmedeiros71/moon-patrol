// Package render draws the game; this file holds the pure letterbox geometry.
package render

import "math"

// Rect is an integer rectangle in screen pixels.
type Rect struct{ X, Y, W, H int }

// Empty reports whether the rectangle has no area.
func (r Rect) Empty() bool { return r.W <= 0 || r.H <= 0 }

// Letterbox maps between world units and screen pixels, keeping the world aspect ratio.
type Letterbox struct {
	// Screen is the full framebuffer, Scene the centred rectangle the world is drawn into.
	Screen, Scene  Rect
	WorldW, WorldH float64
}

// NewLetterbox computes the largest centred rectangle with the world's aspect ratio that
// fits a screenW x screenH framebuffer. A non-positive screen size yields an empty Scene.
func NewLetterbox(screenW, screenH int, worldW, worldH float64) Letterbox {
	l := Letterbox{Screen: Rect{0, 0, screenW, screenH}, WorldW: worldW, WorldH: worldH}
	if screenW <= 0 || screenH <= 0 || worldW <= 0 || worldH <= 0 {
		return l
	}
	s := math.Min(float64(screenW)/worldW, float64(screenH)/worldH)
	w := min(int(math.Round(worldW*s)), screenW)
	h := min(int(math.Round(worldH*s)), screenH)
	l.Scene = Rect{X: (screenW - w) / 2, Y: (screenH - h) / 2, W: w, H: h}
	return l
}

// ToScreen maps a world point to screen pixels.
func (l Letterbox) ToScreen(wx, wy float64) (float64, float64) {
	return float64(l.Scene.X) + wx*float64(l.Scene.W)/l.WorldW,
		float64(l.Scene.Y) + wy*float64(l.Scene.H)/l.WorldH
}

// ToWorld maps a screen point to world units (the inverse of ToScreen).
func (l Letterbox) ToWorld(sx, sy float64) (float64, float64) {
	return (sx - float64(l.Scene.X)) * l.WorldW / float64(l.Scene.W),
		(sy - float64(l.Scene.Y)) * l.WorldH / float64(l.Scene.H)
}

// Scale returns screen pixels per world unit, for line widths and fonts.
func (l Letterbox) Scale() float64 { return float64(l.Scene.W) / l.WorldW }
