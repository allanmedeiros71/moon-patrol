package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Debug font metrics (ebitenutil.DebugPrint): 6x16 px per glyph at factor 1.
const (
	glyphW        = 6
	lineH         = 16
	overlayMargin = 16
	overlayPad    = 8
)

var overlayPanel = color.RGBA{0, 0, 0, 0x99} // #000000 at 60% alpha

// Renderer draws the placeholder scene and the debug overlay.
type Renderer struct {
	scene *ebiten.Image // offscreen of the scene size; the phase 4 glow pass runs here
	text  *ebiten.Image // overlay text at factor 1
}

// NewRenderer returns a Renderer.
func NewRenderer() *Renderer { return &Renderer{} }

// sized returns img when it already has the wanted size, otherwise a new image.
func sized(img *ebiten.Image, w, h int) *ebiten.Image {
	if img != nil && img.Bounds().Dx() == w && img.Bounds().Dy() == h {
		return img
	}
	return ebiten.NewImage(w, h)
}

// Draw renders v into screen.
func (r *Renderer) Draw(screen *ebiten.Image, v View) {
	screen.Fill(ColorBars)
	b := screen.Bounds()
	l := NewLetterbox(b.Dx(), b.Dy(), WorldW, WorldH)
	if l.Scene.Empty() {
		return
	}
	r.scene = sized(r.scene, l.Scene.W, l.Scene.H)
	for _, s := range SceneShapes(l, v.MarkerX) {
		vector.FillRect(r.scene, float32(s.Rect.X-float64(l.Scene.X)), float32(s.Rect.Y-float64(l.Scene.Y)),
			float32(s.Rect.W), float32(s.Rect.H), s.Color, false)
	}
	if v.OverlayVisible {
		r.drawOverlay(v.Overlay, l.Scene.H)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Min.X+l.Scene.X), float64(b.Min.Y+l.Scene.Y))
	screen.DrawImage(r.scene, op)
}

// drawOverlay draws the overlay text on a translucent panel at the scene's top-left.
func (r *Renderer) drawOverlay(info OverlayInfo, sceneH int) {
	lines := OverlayLines(info)
	cols := 0
	for _, ln := range lines {
		cols = max(cols, len(ln))
	}
	tw, th := cols*glyphW, len(lines)*lineH
	r.text = sized(r.text, tw, th)
	r.text.Clear()
	for i, ln := range lines {
		ebitenutil.DebugPrintAt(r.text, ln, 0, i*lineH)
	}
	f := float32(max(1, sceneH/540))
	vector.FillRect(r.scene, overlayMargin*f, overlayMargin*f,
		float32(tw)*f+2*overlayPad*f, float32(th)*f+2*overlayPad*f, overlayPanel, false)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(float64(f), float64(f))
	op.GeoM.Translate(float64((overlayMargin+overlayPad)*f), float64((overlayMargin+overlayPad)*f))
	r.scene.DrawImage(r.text, op)
}
