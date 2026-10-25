package game

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/allanmedeiros/moonpatrol/internal/input"
	"github.com/allanmedeiros/moonpatrol/internal/render"
)

type fakeWindow struct {
	fullscreenCalls []bool
	sizes           [][2]int
	scale           float64
}

func (w *fakeWindow) SetFullscreen(b bool)       { w.fullscreenCalls = append(w.fullscreenCalls, b) }
func (w *fakeWindow) SetWindowSize(x, y int)     { w.sizes = append(w.sizes, [2]int{x, y}) }
func (w *fakeWindow) DeviceScaleFactor() float64 { return w.scale }

type fakePoller struct{ raw input.Raw }

func (p *fakePoller) Poll() input.Raw { return p.raw }

func keys(ks ...input.Key) input.Raw {
	m := map[input.Key]bool{}
	for _, k := range ks {
		m[k] = true
	}
	return input.Raw{Keys: m}
}

func newTest(o Options) (*Game, *fakeWindow, *fakePoller) {
	w, p := &fakeWindow{scale: 1}, &fakePoller{}
	return New(Config{Seed: 1, Options: o, Window: w, Poller: p, Renderer: NopRenderer{}}), w, p
}

func TestInitialModes(t *testing.T) {
	_, w, _ := newTest(Options{Windowed: true})
	if len(w.fullscreenCalls) != 1 || w.fullscreenCalls[0] || len(w.sizes) != 1 || w.sizes[0] != [2]int{1280, 720} {
		t.Errorf("windowed start: %+v", w)
	}
	_, w, _ = newTest(Options{})
	if len(w.fullscreenCalls) != 1 || !w.fullscreenCalls[0] || len(w.sizes) != 0 {
		t.Errorf("fullscreen start: %+v", w)
	}
}

func TestAltEnterToggles(t *testing.T) {
	g, w, p := newTest(Options{Windowed: true})
	p.raw = keys(input.KeyAltLeft, input.KeyEnter)
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	if n := len(w.fullscreenCalls); n != 2 || !w.fullscreenCalls[1] {
		t.Errorf("calls = %v, want fullscreen on", w.fullscreenCalls)
	}
	if len(w.sizes) != 1 {
		t.Errorf("entering fullscreen resized the window: %v", w.sizes)
	}
}

func TestHeldF11TogglesOnce(t *testing.T) {
	g, w, p := newTest(Options{})
	p.raw = keys(input.KeyF11)
	for range 10 {
		if err := g.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if len(w.fullscreenCalls) != 2 || w.fullscreenCalls[1] {
		t.Errorf("calls = %v, want one switch to windowed", w.fullscreenCalls)
	}
	if len(w.sizes) != 1 || w.sizes[0] != [2]int{1280, 720} {
		t.Errorf("sizes = %v", w.sizes)
	}
	if g.World().Tick() != 10 {
		t.Errorf("tick = %d", g.World().Tick())
	}
}

func TestUpdateQuit(t *testing.T) {
	g, _, p := newTest(Options{})
	p.raw = keys(input.KeyEscape)
	if err := g.Update(); err != ebiten.Termination {
		t.Errorf("Update = %v, want Termination", err)
	}
}

func TestLayout(t *testing.T) {
	g, w, _ := newTest(Options{})
	w.scale = 1.5
	if x, y := g.LayoutF(1000, 500); x != 1500 || y != 750 {
		t.Errorf("LayoutF = %v,%v", x, y)
	}
	if x, y := g.Layout(1001, 501); x != 1502 || y != 752 {
		t.Errorf("Layout = %v,%v", x, y)
	}
}

func TestBuildView(t *testing.T) {
	w, p := &fakeWindow{scale: 1}, &fakePoller{}
	g := New(Config{Seed: 1, Window: w, Poller: p, Renderer: NopRenderer{},
		Stats: func() (float64, float64) { return 59.5, 30 }})
	p.raw = keys(input.KeyArrowRight, input.KeyX)
	p.raw.Pads = []input.Pad{{Name: "Pad A", Standard: true}, {Name: "Odd"}}
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	v := g.BuildView(image.Rect(0, 0, 800, 600))
	o := v.Overlay
	checkViewHeldPads(t, o)
	if o.TPS != 59.5 || o.FPS != 30 || o.Tick != 1 || o.ScreenW != 800 || o.ScreenH != 600 {
		t.Errorf("overlay = %+v", o)
	}
	if v.MarkerX != 484 || o.MarkerX != 484 || !v.OverlayVisible {
		t.Errorf("view = %+v", v)
	}
}

func TestBuildViewNoStats(t *testing.T) {
	g, _, _ := newTest(Options{})
	if o := g.BuildView(image.Rect(0, 0, 1, 1)).Overlay; o.TPS != 0 || o.FPS != 0 {
		t.Errorf("overlay = %+v", o)
	}
}

func checkViewHeldPads(t *testing.T, o render.OverlayInfo) {
	t.Helper()
	if !o.Held.Has(input.Right) || !o.Held.Has(input.Fire) || o.Held.Has(input.Left) {
		t.Errorf("held = %q", o.Held)
	}
	if len(o.Pads) != 2 || o.Pads[0].Name != "Pad A" || !o.Pads[0].Standard || o.Pads[1].Standard {
		t.Errorf("pads = %+v", o.Pads)
	}
}
