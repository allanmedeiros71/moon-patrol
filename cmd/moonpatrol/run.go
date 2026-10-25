package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/allanmedeiros/moonpatrol/internal/buildinfo"
	"github.com/hajimehoshi/ebiten/v2"
)

// options holds the parsed command-line flags.
type options struct {
	windowed bool
	seed     int64
}

// blankGame is a minimal placeholder game: a black screen that quits on Esc.
type blankGame struct{ opts options }

// Update ends the game when Esc is pressed.
func (blankGame) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

// Draw leaves the screen black.
func (blankGame) Draw(*ebiten.Image) {}

// Layout returns the outside size unchanged.
func (blankGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// run parses args, starts the game through runGame and returns the exit code.
func run(args []string, stdout io.Writer, runGame func(ebiten.Game) error) int {
	var opts options
	var showVersion bool
	fs := flag.NewFlagSet("moonpatrol", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.BoolVar(&opts.windowed, "windowed", false, "run in a window instead of fullscreen")
	fs.Int64Var(&opts.seed, "seed", 1, "random seed")
	fs.BoolVar(&showVersion, "version", false, "print version and exit")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if showVersion {
		_, _ = fmt.Fprintln(stdout, buildinfo.String())
		return 0
	}
	if !opts.windowed {
		ebiten.SetFullscreen(true)
	}
	if err := runGame(blankGame{opts: opts}); err != nil && !errors.Is(err, ebiten.Termination) {
		_, _ = fmt.Fprintln(os.Stderr, "moonpatrol:", err)
		return 1
	}
	return 0
}
