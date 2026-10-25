package features

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/allanmedeiros/moonpatrol/internal/input"
	"github.com/allanmedeiros/moonpatrol/internal/sim"
)

type timingWorld struct {
	worlds []*sim.World
}

type timingKey struct{}

func tw(ctx context.Context) *timingWorld { return ctx.Value(timingKey{}).(*timingWorld) }

func parseHeld(names string) (input.Frame, error) {
	var h input.Actions
	if names == "none" {
		return input.Frame{}, nil
	}
	all := map[string]input.Action{
		"Left": input.Left, "Right": input.Right, "Jump": input.Jump, "Fire": input.Fire,
	}
	for _, n := range strings.Fields(names) {
		a, ok := all[n]
		if !ok {
			return input.Frame{}, fmt.Errorf("unknown action %q", n)
		}
		h = h.With(a)
	}
	return input.Frame{Held: h}, nil
}

// stepScript feeds a script like "30 Right, 10 none, 20 Left Fire" to w.
func stepScript(w *sim.World, script string) error {
	for _, part := range strings.Split(script, ",") {
		count, names, _ := strings.Cut(strings.TrimSpace(part), " ")
		n, err := strconv.Atoi(count)
		if err != nil {
			return fmt.Errorf("bad script part %q: %w", part, err)
		}
		f, err := parseHeld(strings.TrimSpace(names))
		if err != nil {
			return err
		}
		for range n {
			w.Step(f)
		}
	}
	return nil
}

func (t *timingWorld) newWorlds(seeds ...int64) {
	t.worlds = nil
	for _, s := range seeds {
		t.worlds = append(t.worlds, sim.New(s))
	}
}

func (t *timingWorld) stepHolding(n int, names string) error {
	f, err := parseHeld(names)
	if err != nil {
		return err
	}
	for range n {
		t.worlds[0].Step(f)
	}
	return nil
}

func (t *timingWorld) stepAll(script string) error {
	for _, w := range t.worlds {
		if err := stepScript(w, script); err != nil {
			return err
		}
	}
	return nil
}

func (t *timingWorld) markerIs(want float64) error {
	if got := t.worlds[0].MarkerX(); got != want {
		return fmt.Errorf("marker X = %v, want %v", got, want)
	}
	return nil
}

func (t *timingWorld) tickIs(want int) error {
	if got := t.worlds[0].Tick(); got != uint64(want) {
		return fmt.Errorf("tick = %d, want %d", got, want)
	}
	return nil
}

func (t *timingWorld) timeIs(want float64) error {
	if got := t.worlds[0].Elapsed(); got != want {
		return fmt.Errorf("elapsed = %v, want %v", got, want)
	}
	return nil
}

func (t *timingWorld) identical() error {
	a, b := t.worlds[0].Snapshot(), t.worlds[1].Snapshot()
	if a != b {
		return fmt.Errorf("snapshots differ: %+v vs %+v", a, b)
	}
	return nil
}

func (t *timingWorld) differInNoise() error {
	a, b := t.worlds[0].Snapshot(), t.worlds[1].Snapshot()
	if a.Noise == b.Noise {
		return fmt.Errorf("noise equal: %d", a.Noise)
	}
	a.Noise, b.Noise = 0, 0
	if a != b {
		return fmt.Errorf("snapshots differ beyond noise: %+v vs %+v", a, b)
	}
	return nil
}

func initTiming(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, timingKey{}, &timingWorld{}), nil
	})
	sc.Step(`^a new simulation with seed (\d+)$`, func(ctx context.Context, seed int) {
		tw(ctx).newWorlds(int64(seed))
	})
	sc.Step(`^two new simulations with seed (\d+)$`, func(ctx context.Context, seed int) {
		tw(ctx).newWorlds(int64(seed), int64(seed))
	})
	sc.Step(`^a new simulation with seed (\d+) and a new simulation with seed (\d+)$`,
		func(ctx context.Context, a, b int) { tw(ctx).newWorlds(int64(a), int64(b)) })
	sc.Step(`^it is stepped (\d+) ticks with no input$`, func(ctx context.Context, n int) error {
		return tw(ctx).stepHolding(n, "none")
	})
	sc.Step(`^it is stepped (\d+) ticks holding "([^"]*)"$`, func(ctx context.Context, n int, names string) error {
		return tw(ctx).stepHolding(n, names)
	})
	sc.Step(`^both are stepped with the script "([^"]*)"$`, func(ctx context.Context, s string) error {
		return tw(ctx).stepAll(s)
	})
	sc.Step(`^the tick count is (\d+)$`, func(ctx context.Context, n int) error { return tw(ctx).tickIs(n) })
	sc.Step(`^the simulated time is (\d+\.\d+) seconds$`, func(ctx context.Context, s float64) error {
		return tw(ctx).timeIs(s)
	})
	sc.Step(`^the marker X is (\d+\.\d+)$`, func(ctx context.Context, x float64) error { return tw(ctx).markerIs(x) })
	sc.Step(`^both world snapshots are identical$`, func(ctx context.Context) error { return tw(ctx).identical() })
	sc.Step(`^the world snapshots differ in the noise value only$`, func(ctx context.Context) error {
		return tw(ctx).differInNoise()
	})
}
