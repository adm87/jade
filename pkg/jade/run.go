package jade

import "github.com/hajimehoshi/ebiten/v2"

// Run starts the game.
func Run() error {
	return ebiten.RunGame(&shell{})
}
