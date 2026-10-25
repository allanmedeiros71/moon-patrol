package sim

import (
	"math/rand"
	"testing"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

func held(as ...input.Action) input.Frame {
	var h input.Actions
	for _, a := range as {
		h = h.With(a)
	}
	return input.Frame{Held: h}
}

func stepN(w *World, n int, f input.Frame) {
	for range n {
		w.Step(f)
	}
}

func TestMove(t *testing.T) {
	w := New(1)
	stepN(w, 60, held(input.Right))
	if w.MarkerX() != 720 {
		t.Errorf("right: %v", w.MarkerX())
	}
	stepN(w, 10, held(input.Left))
	if w.MarkerX() != 680 {
		t.Errorf("left: %v", w.MarkerX())
	}
	stepN(w, 5, held(input.Jump, input.Fire))
	if w.MarkerX() != 680 {
		t.Errorf("other actions moved marker: %v", w.MarkerX())
	}
}

func TestLeftAndRight(t *testing.T) {
	w := New(1)
	stepN(w, 30, held(input.Left, input.Right))
	if w.MarkerX() != MarkerStartX {
		t.Errorf("got %v", w.MarkerX())
	}
}

func TestClamp(t *testing.T) {
	w := New(1)
	stepN(w, 200, held(input.Left))
	if w.MarkerX() != 0 {
		t.Errorf("left clamp: %v", w.MarkerX())
	}
	w.Step(held(input.Left))
	if w.MarkerX() != 0 {
		t.Errorf("left clamp after: %v", w.MarkerX())
	}
	stepN(w, 479, held(input.Right))
	if w.MarkerX() != 1916 {
		t.Errorf("before right clamp: %v", w.MarkerX())
	}
	w.Step(held(input.Right))
	if w.MarkerX() != WorldWidth {
		t.Errorf("exact right edge: %v", w.MarkerX())
	}
	w.Step(held(input.Right))
	if w.MarkerX() != WorldWidth {
		t.Errorf("right clamp: %v", w.MarkerX())
	}
}

func TestElapsed(t *testing.T) {
	w := New(1)
	if w.Tick() != 0 || w.Elapsed() != 0 {
		t.Errorf("initial tick %d elapsed %v", w.Tick(), w.Elapsed())
	}
	w.Step(input.Frame{})
	if w.Tick() != 1 || w.Elapsed() != 1.0/60 {
		t.Errorf("1 tick: %d %v", w.Tick(), w.Elapsed())
	}
	stepN(w, 60, input.Frame{})
	if w.Tick() != 61 || w.Elapsed() != 61.0/60 {
		t.Errorf("61 ticks: %d %v", w.Tick(), w.Elapsed())
	}
}

func randomRun(seed int64, inputSeed int64) Snapshot {
	w := New(seed)
	r := rand.New(rand.NewSource(inputSeed)) //nolint:gosec // test input generator
	for range 1000 {
		w.Step(input.Frame{Held: input.Actions(r.Intn(16))})
	}
	return w.Snapshot()
}

func TestDeterminism(t *testing.T) {
	if a, b := randomRun(7, 99), randomRun(7, 99); a != b {
		t.Errorf("same seed differ: %+v %+v", a, b)
	}
	a, b := randomRun(7, 99), randomRun(8, 99)
	if a == b || a.Noise == b.Noise {
		t.Errorf("different seeds equal: %+v %+v", a, b)
	}
	if a.Tick != b.Tick || a.MarkerX != b.MarkerX {
		t.Errorf("seed affected more than noise: %+v %+v", a, b)
	}
}

func TestNoiseChangesEachTick(t *testing.T) {
	w := New(3)
	if w.Snapshot().Noise != 0 {
		t.Error("noise before first step")
	}
	w.Step(input.Frame{})
	n1 := w.Snapshot().Noise
	w.Step(input.Frame{})
	if n1 == 0 || n1 == w.Snapshot().Noise {
		t.Errorf("noise not advancing: %d %d", n1, w.Snapshot().Noise)
	}
}
