package features

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

type inputWorld struct {
	mapper *input.Mapper
	raw    input.Raw
	frames []input.Frame
}

type inputKey struct{}

func iw(ctx context.Context) *inputWorld { return ctx.Value(inputKey{}).(*inputWorld) }

func parseKeys(list string) (map[input.Key]bool, error) {
	keys := map[input.Key]bool{}
	for _, n := range strings.Split(list, ",") {
		k, err := input.ParseKey(strings.TrimSpace(n))
		if err != nil {
			return nil, err
		}
		keys[k] = true
	}
	return keys, nil
}

func checkActions(kind string, got input.Actions, want string) error {
	if got.String() != want {
		return fmt.Errorf("%s actions = %q, want %q", kind, got, want)
	}
	return nil
}

func (w *inputWorld) frame(tick int) (input.Frame, error) {
	if tick < 1 || tick > len(w.frames) {
		return input.Frame{}, fmt.Errorf("no tick %d (have %d)", tick, len(w.frames))
	}
	return w.frames[tick-1], nil
}

func (w *inputWorld) lastPad() *input.Pad { return &w.raw.Pads[len(w.raw.Pads)-1] }

func (w *inputWorld) mapOnce() { w.frames = append(w.frames, w.mapper.Map(w.raw)) }

func initInput(sc *godog.ScenarioContext) {
	initInputKeys(sc)
	initInputPads(sc)
	initInputThen(sc)
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, inputKey{}, &inputWorld{raw: input.Raw{Keys: map[input.Key]bool{}}}), nil
	})
	sc.Step(`^the default bindings$`, func(ctx context.Context) {
		iw(ctx).mapper = input.NewMapper(input.DefaultBindings())
	})
	sc.Step(`^a standard gamepad "([^"]*)" is connected$`, func(ctx context.Context, name string) {
		w := iw(ctx)
		w.raw.Pads = append(w.raw.Pads, input.Pad{Name: name, Standard: true, Buttons: map[input.PadButton]bool{}})
	})
	sc.Step(`^a non-standard gamepad "([^"]*)" is connected$`, func(ctx context.Context, name string) {
		w := iw(ctx)
		w.raw.Pads = append(w.raw.Pads, input.Pad{Name: name})
	})
}

func initInputKeys(sc *godog.ScenarioContext) {
	sc.Step(`^the keys? "([^"]*)" (?:is|are) held$`, func(ctx context.Context, list string) error {
		w := iw(ctx)
		keys, err := parseKeys(list)
		if err != nil {
			return err
		}
		for k := range keys {
			w.raw.Keys[k] = true
		}
		w.mapOnce()
		return nil
	})
}

func initInputPads(sc *godog.ScenarioContext) {
	sc.Step(`^the gamepad button "([^"]*)" is held$`, func(ctx context.Context, name string) error {
		w := iw(ctx)
		b, err := input.ParsePadButton(name)
		if err != nil {
			return err
		}
		w.lastPad().Buttons[b] = true
		w.mapOnce()
		return nil
	})
	sc.Step(`^the gamepad left stick horizontal axis is (-?[0-9.]+)$`, func(ctx context.Context, v string) error {
		w := iw(ctx)
		x, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		w.lastPad().LeftStickX = x
		w.mapOnce()
		return nil
	})
	sc.Step(`^its raw button (\d+) is held$`, func(ctx context.Context, n int) {
		w := iw(ctx)
		p := w.lastPad()
		p.RawButtons = make([]bool, n+1)
		p.RawButtons[n] = true
		w.mapOnce()
	})
	sc.Step(`^the key "([^"]*)" is held for (\d+) ticks?$`, func(ctx context.Context, name string, n int) error {
		w := iw(ctx)
		k, err := input.ParseKey(name)
		if err != nil {
			return err
		}
		w.raw.Keys[k] = true
		for range n {
			w.mapOnce()
		}
		return nil
	})
	sc.Step(`^the key "([^"]*)" is held for 1 tick, released for 1 tick and held for 1 tick$`,
		func(ctx context.Context, name string) error {
			w := iw(ctx)
			k, err := input.ParseKey(name)
			if err != nil {
				return err
			}
			for _, down := range []bool{true, false, true} {
				w.raw.Keys[k] = down
				w.mapOnce()
			}
			return nil
		})
}

func initInputThen(sc *godog.ScenarioContext) {
	sc.Step(`^the held actions are exactly "([^"]*)"$`, func(ctx context.Context, want string) error {
		w := iw(ctx)
		f, err := w.frame(len(w.frames))
		if err != nil {
			return err
		}
		return checkActions("held", f.Held, want)
	})
	sc.Step(`^the pressed actions are exactly "([^"]*)"$`, func(ctx context.Context, want string) error {
		w := iw(ctx)
		f, err := w.frame(len(w.frames))
		if err != nil {
			return err
		}
		return checkActions("pressed", f.Pressed, want)
	})
	sc.Step(`^the pressed actions on tick (\d+) are "([^"]*)"$`, func(ctx context.Context, n int, want string) error {
		f, err := iw(ctx).frame(n)
		if err != nil {
			return err
		}
		return checkActions("pressed", f.Pressed, want)
	})
	sc.Step(`^the held actions on tick (\d+) are "([^"]*)"$`, func(ctx context.Context, n int, want string) error {
		f, err := iw(ctx).frame(n)
		if err != nil {
			return err
		}
		return checkActions("held", f.Held, want)
	})
	sc.Step(`^the gamepad "([^"]*)" is reported as unsupported$`, func(ctx context.Context, name string) error {
		for _, n := range input.Unsupported(iw(ctx).raw) {
			if n == name {
				return nil
			}
		}
		return fmt.Errorf("%q not reported unsupported", name)
	})
}
