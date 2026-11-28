package game

import (
	"github.com/adm87/jade/cmd/game/internal/game/draw"
	"github.com/adm87/jade/cmd/game/internal/game/update"
	"github.com/adm87/jade/pkg/display"
	"github.com/adm87/jade/pkg/jade"
)

// Boot initializes the game.
func Boot() error {
	d := display.New(800, 600, display.WithScaleMode(display.Cover))
	return jade.Run(
		jade.WithLayouter(d),
		jade.WithUpdater(update.Pipeline()),
		jade.WithDrawer(draw.Pipeline(d)),
	)
}
