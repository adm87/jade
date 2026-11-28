package update

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/adm87/jade/pkg/jade"
)

// Pipeline returns the game's updater.
func Pipeline() jade.Updater {
	return jade.UpdateAll(
		jade.UpdateFunc(quitOnEscape),
	)
}

func quitOnEscape(time.Duration) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}
