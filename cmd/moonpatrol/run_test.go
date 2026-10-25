package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/allanmedeiros/moonpatrol/internal/game"
	"github.com/allanmedeiros/moonpatrol/internal/input"
)

func fakeRun(got *ebiten.Game, err error) func(ebiten.Game) error {
	return func(g ebiten.Game) error {
		*got = g
		return err
	}
}

func TestRunVersion(t *testing.T) {
	var out bytes.Buffer
	var g ebiten.Game
	if code := run([]string{"-version"}, &out, fakeRun(&g, nil)); code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if !strings.HasPrefix(out.String(), "moonpatrol ") {
		t.Errorf("output = %q", out.String())
	}
	if g != nil {
		t.Error("game started on -version")
	}
}

func TestRunOptions(t *testing.T) {
	seeds := map[string]uint64{}
	for _, args := range [][]string{nil, {"-windowed"}, {"-seed", "7"}} {
		var g ebiten.Game
		if code := run(args, &bytes.Buffer{}, fakeRun(&g, nil)); code != 0 {
			t.Fatalf("code = %d, want 0", code)
		}
		gg, ok := g.(*game.Game)
		if !ok {
			t.Fatalf("game type = %T", g)
		}
		gg.World().Step(input.Frame{})
		seeds[strings.Join(args, " ")] = gg.World().Snapshot().Noise
	}
	if seeds[""] != seeds["-windowed"] {
		t.Error("default seed differs between runs")
	}
	if seeds[""] == seeds["-seed 7"] {
		t.Error("-seed did not change the world")
	}
}

func TestRunBadFlag(t *testing.T) {
	var g ebiten.Game
	if code := run([]string{"-nope"}, &bytes.Buffer{}, fakeRun(&g, nil)); code != 2 {
		t.Errorf("code = %d, want 2", code)
	}
	if g != nil {
		t.Error("game started on bad flag")
	}
}

func TestRunHelp(t *testing.T) {
	var g ebiten.Game
	if code := run([]string{"-h"}, &bytes.Buffer{}, fakeRun(&g, nil)); code != 0 {
		t.Errorf("code = %d, want 0", code)
	}
}

func TestRunGameErrors(t *testing.T) {
	var g ebiten.Game
	if code := run([]string{"-windowed"}, &bytes.Buffer{}, fakeRun(&g, errors.New("boom"))); code != 1 {
		t.Errorf("error code = %d, want 1", code)
	}
	if code := run([]string{"-windowed"}, &bytes.Buffer{}, fakeRun(&g, ebiten.Termination)); code != 0 {
		t.Errorf("termination code = %d, want 0", code)
	}
}
