package render

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestRendererDrawAllocatesAndReusesOffscreen(t *testing.T) {
	r := NewRenderer()
	screen := ebiten.NewImage(640, 400)
	v := View{MarkerX: 480, OverlayVisible: true, Overlay: OverlayInfo{ScreenW: 640, ScreenH: 400}}
	r.Draw(screen, v)
	first := r.scene
	if first == nil || first.Bounds().Dx() != 640 || first.Bounds().Dy() != 360 {
		t.Fatalf("scene = %v", first)
	}
	if r.text == nil {
		t.Error("overlay text image not allocated")
	}
	r.Draw(screen, v)
	if r.scene != first {
		t.Error("offscreen reallocated for an unchanged size")
	}
	r.Draw(ebiten.NewImage(1280, 720), View{})
	if r.scene == first || r.scene.Bounds().Dx() != 1280 {
		t.Error("offscreen not reallocated for a new size")
	}
}

func TestRendererDrawTinyFramebuffer(t *testing.T) {
	r := NewRenderer()
	r.Draw(ebiten.NewImage(1, 1), View{OverlayVisible: true})
	if r.scene == nil || r.scene.Bounds().Dx() != 1 {
		t.Error("tiny framebuffer should still get a 1x1 scene")
	}
}

func TestRendererOverlayScalesAtHighHeight(t *testing.T) {
	r := NewRenderer()
	r.Draw(ebiten.NewImage(1920, 1080), View{OverlayVisible: true, Overlay: OverlayInfo{ScreenW: 1920, ScreenH: 1080}})
	if r.scene.Bounds().Dy() != 1080 || r.text == nil {
		t.Error("full-size scene not drawn")
	}
}
