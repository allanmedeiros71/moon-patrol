package features

import (
	"context"
	"fmt"
	"math"

	"github.com/cucumber/godog"

	"github.com/allanmedeiros/moonpatrol/internal/render"
)

type sceneState struct {
	markerX float64
	shapes  []render.Shape
}

type sceneKey struct{}

func sceneOf(ctx context.Context) *sceneState { return ctx.Value(sceneKey{}).(*sceneState) }

func checkShape(ctx context.Context, idx int, name string, x, y, w, h float64) error {
	s := sceneOf(ctx).shapes
	if idx >= len(s) {
		return fmt.Errorf("no %s shape (have %d shapes)", name, len(s))
	}
	r := s[idx].Rect
	if math.Abs(r.X-x) > 1e-6 || math.Abs(r.Y-y) > 1e-6 || math.Abs(r.W-w) > 1e-6 || math.Abs(r.H-h) > 1e-6 {
		return fmt.Errorf("%s = %+v, want %v,%v %vx%v", name, r, x, y, w, h)
	}
	return nil
}

func initScene(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, sceneKey{}, &sceneState{}), nil
	})
	sc.Step(`^the marker is at world X (\d+)$`, func(ctx context.Context, x float64) {
		sceneOf(ctx).markerX = x
	})
	sc.Step(`^the scene shapes are computed$`, func(ctx context.Context) {
		s := sceneOf(ctx)
		s.shapes = render.SceneShapes(bw(ctx).box, s.markerX)
	})
	sc.Step(`^the marker shape is at screen (\d+),(\d+) with size (\d+)x(\d+)$`,
		func(ctx context.Context, x, y, w, h float64) error { return checkShape(ctx, 3, "marker", x, y, w, h) })
	sc.Step(`^the ground shape is at screen (\d+),(\d+) with size (\d+)x(\d+)$`,
		func(ctx context.Context, x, y, w, h float64) error { return checkShape(ctx, 1, "ground", x, y, w, h) })
}
