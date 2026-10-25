package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
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
	tests := []struct {
		name string
		args []string
		want options
	}{
		{"defaults", nil, options{windowed: false, seed: 1}},
		{"windowed", []string{"-windowed"}, options{windowed: true, seed: 1}},
		{"seed", []string{"-seed", "7"}, options{windowed: false, seed: 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g ebiten.Game
			if code := run(tt.args, &bytes.Buffer{}, fakeRun(&g, nil)); code != 0 {
				t.Fatalf("code = %d, want 0", code)
			}
			bg, ok := g.(blankGame)
			if !ok {
				t.Fatalf("game type = %T", g)
			}
			if bg.opts != tt.want {
				t.Errorf("opts = %+v, want %+v", bg.opts, tt.want)
			}
		})
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

func TestBlankGameLayout(t *testing.T) {
	w, h := blankGame{}.Layout(640, 360)
	if w != 640 || h != 360 {
		t.Errorf("Layout = %d,%d", w, h)
	}
}

func TestBlankGameUpdateDraw(t *testing.T) {
	g := blankGame{}
	if err := g.Update(); err != nil {
		t.Errorf("Update() = %v, want nil with no key pressed", err)
	}
	g.Draw(nil)
}
