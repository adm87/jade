package game

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/adm87/jade/pkg/display"
	"github.com/adm87/jade/pkg/jade"
)

// Boot initializes the sandbox game.
func Boot() error {
	d := display.New(800, 600, display.WithScaleMode(display.Cover))
	world := corners(d)
	return jade.Run(
		jade.WithLayouter(d),
		jade.WithUpdater(jade.UpdateFunc(quitOnEscape)),
		jade.WithDrawer(jade.DrawAll(
			jade.Chain(
				display.Buffer(d, ebiten.FilterNearest),
			)(world),
		)),
	)
}

func quitOnEscape(time.Duration) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

// corners marks each corner of the display's visible area.
func corners(d *display.Display) jade.Drawer {
	return jade.DrawFunc(func(dst *ebiten.Image) {
		visible := d.Visible()
		minX, minY := visible.Min()
		maxX, maxY := visible.Max()

		mark := func(x, y float64, c color.Color) {
			vector.FillRect(dst, float32(x), float32(y), 10, 10, c, false)
		}

		mark(minX, minY, color.RGBA{R: 255, A: 255})
		mark(maxX-10, minY, color.RGBA{R: 255, G: 255, A: 255})
		mark(minX, maxY-10, color.RGBA{G: 255, A: 255})
		mark(maxX-10, maxY-10, color.RGBA{B: 255, A: 255})
	})
}
