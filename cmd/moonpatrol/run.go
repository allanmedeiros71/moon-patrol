package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/allanmedeiros/moonpatrol/internal/buildinfo"
	"github.com/allanmedeiros/moonpatrol/internal/game"
	"github.com/allanmedeiros/moonpatrol/internal/input/inputebiten"
	"github.com/allanmedeiros/moonpatrol/internal/render"
	"github.com/hajimehoshi/ebiten/v2"
)

// options holds the parsed command-line flags.
type options struct {
	windowed bool
	seed     int64
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
	ebiten.SetWindowTitle("Moon Patrol")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(game.TPS)
	g := game.New(game.Config{
		Seed:     opts.seed,
		Options:  game.Options{Windowed: opts.windowed},
		Window:   game.EbitenWindow{},
		Poller:   &inputebiten.EbitenPoller{},
		Renderer: render.NewRenderer(),
		Stats:    func() (float64, float64) { return ebiten.ActualTPS(), ebiten.ActualFPS() },
	})
	if err := runGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		_, _ = fmt.Fprintln(os.Stderr, "moonpatrol:", err)
		return 1
	}
	return 0
}
