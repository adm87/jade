package jade

import "github.com/hajimehoshi/ebiten/v2"

// Run starts the game.
func Run(opts ...Option) error {
	o := applyOptions(opts)

	setupEbitengine()
	return ebiten.RunGame(&shell{
		drawer:   o.drawer,
		updater:  o.updater,
		layouter: o.layouter,
	})
}

func setupEbitengine() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// TODO: Set the window size, title, etc.
}
