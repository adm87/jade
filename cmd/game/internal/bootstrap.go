package internal

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/adm87/jade/pkg/jade"
)

// Boot initializes the sandbox game.
func Boot() error {
	return jade.Run(
		jade.WithUpdater(jade.UpdateFunc(func(_ time.Duration) error {
			if ebiten.IsKeyPressed(ebiten.KeyEscape) {
				return ebiten.Termination
			}
			return nil
		})),
		jade.WithDrawer(jade.DrawFunc(func(target *ebiten.Image) {
			ebitenutil.DebugPrint(target, "Hello, World!")
		})),
	)
}
