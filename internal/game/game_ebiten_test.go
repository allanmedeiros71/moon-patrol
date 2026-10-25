package game

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/allanmedeiros/moonpatrol/internal/input"
	"github.com/allanmedeiros/moonpatrol/internal/render"
)

// recordingRenderer keeps the last view and screen size it was asked to draw.
type recordingRenderer struct {
	calls int
	w, h  int
	view  render.View
}

func (r *recordingRenderer) Draw(screen *ebiten.Image, v render.View) {
	r.calls++
	r.w, r.h = screen.Bounds().Dx(), screen.Bounds().Dy()
	r.view = v
}

func TestDrawPassesHeldActionsAndPadsToRenderer(t *testing.T) {
	w, p := &fakeWindow{scale: 1}, &fakePoller{}
	rr := &recordingRenderer{}
	g := New(Config{Seed: 1, Window: w, Poller: p, Renderer: rr})
	p.raw = keys(input.KeyArrowRight, input.KeyX)
	p.raw.Pads = []input.Pad{{Name: "Pad A", Standard: true}, {Name: "Odd"}}
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	g.Draw(ebiten.NewImage(640, 360))
	if rr.calls != 1 || rr.w != 640 || rr.h != 360 {
		t.Fatalf("renderer calls=%d size=%dx%d", rr.calls, rr.w, rr.h)
	}
	checkViewHeldPads(t, rr.view.Overlay)
	if rr.view.Overlay.ScreenW != 640 || rr.view.Overlay.ScreenH != 360 {
		t.Errorf("overlay screen = %dx%d", rr.view.Overlay.ScreenW, rr.view.Overlay.ScreenH)
	}
}

func TestNopRendererDrawIsSafe(_ *testing.T) {
	img := ebiten.NewImage(2, 2)
	NopRenderer{}.Draw(img, render.View{})
}

func TestEbitenWindowSetters(t *testing.T) {
	var w Window = EbitenWindow{}
	w.SetFullscreen(false)
	w.SetWindowSize(WindowedWidth, WindowedHeight)
	if s := w.DeviceScaleFactor(); s <= 0 {
		t.Errorf("DeviceScaleFactor = %v, want > 0", s)
	}
}
