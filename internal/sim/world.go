// Package sim holds the headless, deterministic simulation. It advances one fixed tick per
// Step, never reads the clock and does not import the engine.
package sim

import (
	"math/rand"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

// World units: origin top-left, Y down. Durations are ticks (R-02, D-04).
const (
	// TPS is the number of simulation ticks per second.
	TPS = 60
	// WorldWidth is the logical world width in world units.
	WorldWidth = 1920.0
	// WorldHeight is the logical world height in world units.
	WorldHeight = 1080.0
	// GroundY is the Y coordinate of the ground line in world units.
	GroundY = 900.0

	// MarkerStartX is the initial X of the marker. The marker is a phase 1 placeholder
	// replaced by the buggy in phase 2.
	MarkerStartX = 480.0
	// MarkerSpeed is the marker speed in world units per tick (phase 1 placeholder,
	// replaced by the buggy in phase 2).
	MarkerSpeed = 4.0
)

// World is the simulation state.
type World struct {
	rng     *rand.Rand
	tick    uint64
	markerX float64
	noise   uint64
}

// Snapshot is a comparable copy of the observable world state.
type Snapshot struct {
	Tick    uint64
	MarkerX float64
	Noise   uint64
}

// New creates a world whose random generator is seeded with seed.
func New(seed int64) *World {
	return &World{rng: rand.New(rand.NewSource(seed)), markerX: MarkerStartX} //nolint:gosec // deterministic simulation, not security
}

// Step advances the world by one tick using the input frame f.
func (w *World) Step(f input.Frame) {
	w.tick++
	if f.Held.Has(input.Right) {
		w.markerX += MarkerSpeed
	}
	if f.Held.Has(input.Left) {
		w.markerX -= MarkerSpeed
	}
	w.markerX = min(max(w.markerX, 0), WorldWidth)
	w.noise = w.rng.Uint64()
}

// Tick returns the number of ticks stepped so far.
func (w *World) Tick() uint64 { return w.tick }

// Elapsed returns the simulated time in seconds.
func (w *World) Elapsed() float64 { return float64(w.tick) / TPS }

// MarkerX returns the X of the placeholder marker in world units.
func (w *World) MarkerX() float64 { return w.markerX }

// Snapshot returns the current observable state.
func (w *World) Snapshot() Snapshot {
	return Snapshot{Tick: w.tick, MarkerX: w.markerX, Noise: w.noise}
}
