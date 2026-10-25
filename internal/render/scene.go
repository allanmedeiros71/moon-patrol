package render

import "image/color"

// World size and placeholder-scene geometry in world units (contract S-01-01).
const (
	WorldW = 1920.0
	WorldH = 1080.0

	groundTopY    = 900.0 // world Y of the ground's top edge
	groundEdgeH   = 4.0   // thickness of the ground edge line
	markerW       = 96.0
	markerH       = 48.0
	markerBottomY = groundTopY
)

// Colour tokens of the placeholder scene (.project/DESIGN.md).
var (
	ColorBars       = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	ColorSky        = color.RGBA{0x0B, 0x08, 0x20, 0xFF}
	ColorGround     = color.RGBA{0x3B, 0x2A, 0x4F, 0xFF}
	ColorGroundEdge = color.RGBA{0xFF, 0x3E, 0xA5, 0xFF}
	ColorMarker     = color.RGBA{0xFF, 0xD2, 0x3F, 0xFF}
)

// FRect is a rectangle in screen pixels with fractional coordinates.
type FRect struct{ X, Y, W, H float64 }

// Shape is a filled rectangle with its colour.
type Shape struct {
	Rect  FRect
	Color color.RGBA
}

// SceneShapes returns the placeholder scene in draw order: sky, ground, ground edge, marker.
// An empty scene yields no shapes.
func SceneShapes(l Letterbox, markerX float64) []Shape {
	if l.Scene.Empty() {
		return nil
	}
	rect := func(wx, wy, ww, wh float64) FRect {
		x0, y0 := l.ToScreen(wx, wy)
		x1, y1 := l.ToScreen(wx+ww, wy+wh)
		return FRect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
	}
	return []Shape{
		{rect(0, 0, l.WorldW, l.WorldH), ColorSky},
		{rect(0, groundTopY, l.WorldW, l.WorldH-groundTopY), ColorGround},
		{rect(0, groundTopY, l.WorldW, groundEdgeH), ColorGroundEdge},
		{rect(markerX-markerW/2, markerBottomY-markerH, markerW, markerH), ColorMarker},
	}
}
