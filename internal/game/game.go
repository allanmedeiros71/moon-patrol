package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/allanmedeiros/moonpatrol/internal/input"
	"github.com/allanmedeiros/moonpatrol/internal/render"
	"github.com/allanmedeiros/moonpatrol/internal/sim"
)

// TPS is the number of logic ticks per second (R-02).
const TPS = sim.TPS

// Window is the part of the window API the game controls.
type Window interface {
	SetFullscreen(bool)
	SetWindowSize(w, h int)
	DeviceScaleFactor() float64
}

// Poller reads the raw state of the input devices once per tick.
type Poller interface {
	Poll() input.Raw
}

// Renderer draws a view onto the screen.
type Renderer interface {
	Draw(screen *ebiten.Image, v render.View)
}

// NopRenderer draws nothing; the screen stays cleared.
type NopRenderer struct{}

// Draw does nothing.
func (NopRenderer) Draw(*ebiten.Image, render.View) {}

// Config holds the game's dependencies.
type Config struct {
	Seed     int64
	Options  Options
	Window   Window
	Poller   Poller
	Renderer Renderer
	// Stats returns the measured TPS and FPS; nil reports zeros.
	Stats func() (tps, fps float64)
}

// Game implements ebiten.Game.
type Game struct {
	c       Config
	world   *sim.World
	mapper  *input.Mapper
	display *Display
	raw     input.Raw
	held    input.Actions
}

// New creates a game and applies the initial display mode.
func New(c Config) *Game {
	g := &Game{
		c:       c,
		world:   sim.New(c.Seed),
		mapper:  input.NewMapper(input.DefaultBindings()),
		display: NewDisplay(c.Options),
	}
	c.Window.SetFullscreen(g.display.Fullscreen())
	if !g.display.Fullscreen() {
		c.Window.SetWindowSize(g.display.WindowSize())
	}
	return g
}

// World returns the simulation; use it read-only.
func (g *Game) World() *sim.World { return g.world }

// Update runs one simulation tick.
func (g *Game) Update() error {
	g.raw = g.c.Poller.Poll()
	f := g.mapper.Map(g.raw)
	g.held = f.Held
	g.world.Step(f)
	if g.display.Apply(f).FullscreenChanged {
		g.c.Window.SetFullscreen(g.display.Fullscreen())
		if !g.display.Fullscreen() {
			g.c.Window.SetWindowSize(g.display.WindowSize())
		}
	}
	if g.display.QuitRequested() {
		return ebiten.Termination
	}
	return nil
}

// BuildView assembles what the renderer needs for a screen of the given size.
func (g *Game) BuildView(b image.Rectangle) render.View {
	var tps, fps float64
	if g.c.Stats != nil {
		tps, fps = g.c.Stats()
	}
	pads := make([]render.PadInfo, 0, len(g.raw.Pads))
	for _, p := range g.raw.Pads {
		pads = append(pads, render.PadInfo{Name: p.Name, Standard: p.Standard})
	}
	return render.View{
		MarkerX: g.world.MarkerX(),
		Overlay: render.OverlayInfo{
			TPS: tps, FPS: fps,
			Tick:    g.world.Tick(),
			ScreenW: b.Dx(), ScreenH: b.Dy(),
			Held:    g.held,
			MarkerX: g.world.MarkerX(),
			Pads:    pads,
		},
		OverlayVisible: g.display.DebugVisible(),
	}
}

// Draw renders the current state.
func (g *Game) Draw(screen *ebiten.Image) {
	g.c.Renderer.Draw(screen, g.BuildView(screen.Bounds()))
}

// LayoutF returns the outside size scaled to native resolution (D-10).
func (g *Game) LayoutF(outsideW, outsideH float64) (float64, float64) {
	s := g.c.Window.DeviceScaleFactor()
	return outsideW * s, outsideH * s
}

// Layout is required by ebiten.Game; LayoutF is used instead.
func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	w, h := g.LayoutF(float64(outsideW), float64(outsideH))
	return int(math.Round(w)), int(math.Round(h))
}
