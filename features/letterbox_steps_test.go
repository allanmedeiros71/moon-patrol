package features

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/allanmedeiros/moonpatrol/internal/render"
)

type boxWorld struct {
	sw, sh int
	box    render.Letterbox
	x, y   float64
}

type boxKey struct{}

func bw(ctx context.Context) *boxWorld { return ctx.Value(boxKey{}).(*boxWorld) }

func parsePair(s string) (float64, float64, error) {
	a, b, ok := strings.Cut(s, ",")
	if !ok {
		return 0, 0, fmt.Errorf("bad pair %q", s)
	}
	x, err := strconv.ParseFloat(strings.TrimSpace(a), 64)
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(b), 64)
	return x, y, err
}

func (w *boxWorld) compute() { w.box = render.NewLetterbox(w.sw, w.sh, 1920, 1080) }

func checkPoint(ctx context.Context, pair string) error {
	wx, wy, err := parsePair(pair)
	if err != nil {
		return err
	}
	w := bw(ctx)
	if math.Abs(w.x-wx) > 1e-6 || math.Abs(w.y-wy) > 1e-6 {
		return fmt.Errorf("point = %v,%v, want %v,%v", w.x, w.y, wx, wy)
	}
	return nil
}

func initLetterbox(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, boxKey{}, &boxWorld{}), nil
	})
	sc.Step(`^a screen of (\d+)x(\d+) pixels$`, func(ctx context.Context, sw, sh int) {
		w := bw(ctx)
		w.sw, w.sh = sw, sh
		w.compute()
	})
	sc.Step(`^the letterbox is computed$`, func(ctx context.Context) { bw(ctx).compute() })
	sc.Step(`^the scene rectangle is at (\d+),(\d+) with size (\d+)x(\d+)$`,
		func(ctx context.Context, x, y, wd, h int) error {
			if got := bw(ctx).box.Scene; got != (render.Rect{X: x, Y: y, W: wd, H: h}) {
				return fmt.Errorf("scene = %+v", got)
			}
			return nil
		})
	sc.Step(`^the scene rectangle is empty$`, func(ctx context.Context) error {
		if !bw(ctx).box.Scene.Empty() {
			return fmt.Errorf("scene = %+v, want empty", bw(ctx).box.Scene)
		}
		return nil
	})
	sc.Step(`^the world point ([\d.]+),([\d.]+) is mapped to the screen$`,
		func(ctx context.Context, wx, wy float64) {
			w := bw(ctx)
			w.x, w.y = w.box.ToScreen(wx, wy)
		})
	sc.Step(`^the screen point is ([\d.]+,[\d.]+)$`, checkPoint)
	sc.Step(`^the screen point ([\d.]+),([\d.]+) is mapped to the world$`,
		func(ctx context.Context, sx, sy float64) {
			w := bw(ctx)
			w.x, w.y = w.box.ToWorld(sx, sy)
		})
	sc.Step(`^the world point is ([\d.]+,[\d.]+)$`, checkPoint)
}
