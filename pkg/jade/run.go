package jade

import "github.com/hajimehoshi/ebiten/v2"

// Run starts the game.
func Run(opts ...Option) error {
	o := applyOptions(opts)

	return ebiten.RunGame(&shell{
		drawer: o.drawer,
	})
}
