package features

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/allanmedeiros/moonpatrol/internal/game"
	"github.com/allanmedeiros/moonpatrol/internal/input"
	"github.com/allanmedeiros/moonpatrol/internal/render"
)

type displayWorld struct {
	d     *game.Display
	info  render.OverlayInfo
	lines []string
}

type displayKey struct{}

func dw(ctx context.Context) *displayWorld { return ctx.Value(displayKey{}).(*displayWorld) }

func parseOverlayInfo(t *godog.Table) (render.OverlayInfo, error) {
	var i render.OverlayInfo
	for _, row := range t.Rows {
		k, v := row.Cells[0].Value, strings.TrimSpace(row.Cells[1].Value)
		var err error
		switch k {
		case "tps":
			i.TPS, err = strconv.ParseFloat(v, 64)
		case "fps":
			i.FPS, err = strconv.ParseFloat(v, 64)
		case "tick":
			i.Tick, err = strconv.ParseUint(v, 10, 64)
		case "marker":
			i.MarkerX, err = strconv.ParseFloat(v, 64)
		case "screen":
			_, err = fmt.Sscanf(v, "%dx%d", &i.ScreenW, &i.ScreenH)
		case "held":
			i.Held, err = parseActionNames(v)
		case "gamepads":
			i.Pads = parsePads(v)
		default:
			err = fmt.Errorf("unknown overlay key %q", k)
		}
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func parseActionNames(v string) (input.Actions, error) {
	var s input.Actions
	for _, n := range strings.Fields(v) {
		found := false
		for a := input.Left; a <= input.Quit; a++ {
			if input.Actions(0).With(a).String() == n {
				s, found = s.With(a), true
			}
		}
		if !found {
			return 0, fmt.Errorf("unknown action %q", n)
		}
	}
	return s, nil
}

func parsePads(v string) []render.PadInfo {
	if v == "" {
		return nil
	}
	var pads []render.PadInfo
	for _, p := range strings.Split(v, ";") {
		p = strings.TrimSpace(p)
		name, kind, _ := strings.Cut(p, " (")
		pads = append(pads, render.PadInfo{Name: name, Standard: kind == "standard)"})
	}
	return pads
}

func (w *displayWorld) checkMode(fullscreen bool, sw, sh int) error {
	if w.d.Fullscreen() != fullscreen {
		return fmt.Errorf("fullscreen = %v, want %v", w.d.Fullscreen(), fullscreen)
	}
	if fullscreen {
		return nil
	}
	if gw, gh := w.d.WindowSize(); gw != sw || gh != sh {
		return fmt.Errorf("window size = %dx%d, want %dx%d", gw, gh, sw, sh)
	}
	return nil
}

func (w *displayWorld) checkDebug(want bool) error {
	if w.d.DebugVisible() != want {
		return fmt.Errorf("overlay visible = %v, want %v", w.d.DebugVisible(), want)
	}
	return nil
}

func (w *displayWorld) checkLine(n int, want string) error {
	if n < 1 || n > len(w.lines) {
		return fmt.Errorf("no overlay line %d (have %d)", n, len(w.lines))
	}
	if w.lines[n-1] != want {
		return fmt.Errorf("overlay line %d = %q, want %q", n, w.lines[n-1], want)
	}
	return nil
}

func initDisplay(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, displayKey{}, &displayWorld{}), nil
	})
	sc.Step(`^the display controller is created with default options$`, func(ctx context.Context) {
		dw(ctx).d = game.NewDisplay(game.Options{})
	})
	sc.Step(`^the display controller is created with the windowed option$`, func(ctx context.Context) {
		dw(ctx).d = game.NewDisplay(game.Options{Windowed: true})
	})
	sc.Step(`^"([^"]*)" is pressed$`, func(ctx context.Context, name string) error {
		a, err := parseActionNames(name)
		if err != nil {
			return err
		}
		dw(ctx).d.Apply(input.Frame{Held: a, Pressed: a})
		return nil
	})
	sc.Step(`^the keys "([^"]*)" are held for (\d+) ticks$`, func(ctx context.Context, list string, n int) error {
		keys, err := parseKeys(list)
		if err != nil {
			return err
		}
		m := input.NewMapper(input.DefaultBindings())
		for range n {
			dw(ctx).d.Apply(m.Map(input.Raw{Keys: keys}))
		}
		return nil
	})
	sc.Step(`^the display mode is fullscreen$`, func(ctx context.Context) error {
		return dw(ctx).checkMode(true, 0, 0)
	})
	sc.Step(`^the display mode is windowed with size (\d+)x(\d+)$`, func(ctx context.Context, w, h int) error {
		return dw(ctx).checkMode(false, w, h)
	})
	sc.Step(`^the debug overlay is visible$`, func(ctx context.Context) error { return dw(ctx).checkDebug(true) })
	sc.Step(`^the debug overlay is hidden$`, func(ctx context.Context) error { return dw(ctx).checkDebug(false) })
	initOverlay(sc)
}

func initOverlay(sc *godog.ScenarioContext) {
	sc.Step(`^the overlay info:$`, func(ctx context.Context, t *godog.Table) error {
		i, err := parseOverlayInfo(t)
		dw(ctx).info = i
		return err
	})
	sc.Step(`^the overlay lines are built$`, func(ctx context.Context) {
		w := dw(ctx)
		w.lines = render.OverlayLines(w.info)
	})
	sc.Step(`^the overlay lines are:$`, func(ctx context.Context, t *godog.Table) error {
		w := dw(ctx)
		if len(t.Rows) != len(w.lines) {
			return fmt.Errorf("got %d lines %q, want %d", len(w.lines), w.lines, len(t.Rows))
		}
		for n, row := range t.Rows {
			// Gherkin trims cell padding, so exact inner spacing is compared.
			if err := w.checkLine(n+1, row.Cells[0].Value); err != nil {
				return err
			}
		}
		return nil
	})
	sc.Step(`^(?:And )?overlay line (\d+) is "([^"]*)"$`, func(ctx context.Context, n int, want string) error {
		return dw(ctx).checkLine(n, want)
	})
}
