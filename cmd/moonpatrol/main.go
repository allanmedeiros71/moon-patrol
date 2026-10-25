// Command moonpatrol runs the Moon Patrol game.
package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, ebiten.RunGame))
}
