package features

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/cucumber/godog"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/allanmedeiros/moonpatrol/internal/game"
	"github.com/allanmedeiros/moonpatrol/internal/input"
)

type loopWindow struct {
	scale      float64
	fullscreen bool
	sizeSet    [2]int
}

func (w *loopWindow) SetFullscreen(b bool)       { w.fullscreen = b }
func (w *loopWindow) SetWindowSize(x, y int)     { w.sizeSet = [2]int{x, y} }
func (w *loopWindow) DeviceScaleFactor() float64 { return w.scale }

type loopPoller struct{ raw input.Raw }

func (p *loopPoller) Poll() input.Raw { return p.raw }

type loopGame struct {
	g *game.Game
	w *loopWindow
	p *loopPoller
}

type loopWorld struct {
	games   []*loopGame
	lastErr error
	lw, lh  int
}

type loopKey struct{}

func lw(ctx context.Context) *loopWorld { return ctx.Value(loopKey{}).(*loopWorld) }

func (w *loopWorld) last() *loopGame { return w.games[len(w.games)-1] }

func (w *loopWorld) update(lg *loopGame, n int) {
	for range n {
		if err := lg.g.Update(); err != nil {
			w.lastErr = err
		}
	}
}

func initGameLoop(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, loopKey{}, &loopWorld{}), nil
	})
	sc.Step(`^a headless game with seed (\d+) and a fake window of (\d+)x(\d+)$`, func(ctx context.Context, seed int64, _, _ int) {
		w := lw(ctx)
		lg := &loopGame{w: &loopWindow{scale: 1}, p: &loopPoller{}}
		lg.g = game.New(game.Config{Seed: seed, Window: lg.w, Poller: lg.p, Renderer: game.NopRenderer{}})
		w.games = append(w.games, lg)
	})
	sc.Step(`^the fake window is fullscreen$`, func(ctx context.Context) error {
		if !lw(ctx).last().w.fullscreen {
			return errors.New("fake window is not fullscreen")
		}
		return nil
	})
	sc.Step(`^the fake window device scale factor is ([0-9.]+)$`, func(ctx context.Context, f float64) {
		lw(ctx).last().w.scale = f
	})
	sc.Step(`^Update is called (\d+) times$`, func(ctx context.Context, n int) {
		w := lw(ctx)
		w.update(w.last(), n)
	})
	sc.Step(`^the simulation tick count is (\d+)$`, func(ctx context.Context, n uint64) error {
		if got := lw(ctx).last().g.World().Tick(); got != n {
			return fmt.Errorf("tick = %d, want %d", got, n)
		}
		return nil
	})
	sc.Step(`^the configured ticks per second is (\d+)$`, func(n int) error {
		if game.TPS != n {
			return fmt.Errorf("TPS = %d, want %d", game.TPS, n)
		}
		return nil
	})
	initGameLoopInput(sc)
}

func initGameLoopInput(sc *godog.ScenarioContext) {
	sc.Step(`^the fake input holds the key "([^"]*)" for (\d+) Updates?$`, func(ctx context.Context, name string, n int) error {
		keys, err := parseKeys(name)
		if err != nil {
			return err
		}
		w := lw(ctx)
		w.last().p.raw = input.Raw{Keys: keys}
		w.update(w.last(), n)
		return nil
	})
	sc.Step(`^the fake window was set to windowed with size (\d+)x(\d+)$`, func(ctx context.Context, x, y int) error {
		fw := lw(ctx).last().w
		if fw.fullscreen || fw.sizeSet != [2]int{x, y} {
			return fmt.Errorf("fullscreen=%v size=%v, want windowed %dx%d", fw.fullscreen, fw.sizeSet, x, y)
		}
		return nil
	})
	sc.Step(`^Update returns the termination signal$`, func(ctx context.Context) error {
		if !errors.Is(lw(ctx).lastErr, ebiten.Termination) {
			return fmt.Errorf("Update error = %v, want Termination", lw(ctx).lastErr)
		}
		return nil
	})
	sc.Step(`^both games receive the key "([^"]*)" for (\d+) Updates$`, func(ctx context.Context, name string, n int) error {
		keys, err := parseKeys(name)
		if err != nil {
			return err
		}
		w := lw(ctx)
		for _, lg := range w.games {
			lg.p.raw = input.Raw{Keys: keys}
			w.update(lg, n)
		}
		return nil
	})
	sc.Step(`^both simulations have marker X ([0-9.]+)$`, func(ctx context.Context, want float64) error {
		for i, lg := range lw(ctx).games {
			if got := lg.g.World().MarkerX(); math.Abs(got-want) > 1e-9 {
				return fmt.Errorf("game %d marker X = %v, want %v", i, got, want)
			}
		}
		return nil
	})
	sc.Step(`^Layout is called with an outside size of (\d+)x(\d+)$`, func(ctx context.Context, x, y int) {
		w := lw(ctx)
		w.lw, w.lh = w.last().g.Layout(x, y)
	})
	sc.Step(`^the layout size is (\d+)x(\d+)$`, func(ctx context.Context, x, y int) error {
		if w := lw(ctx); w.lw != x || w.lh != y {
			return fmt.Errorf("layout = %dx%d, want %dx%d", w.lw, w.lh, x, y)
		}
		return nil
	})
}
