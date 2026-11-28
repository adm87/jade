package draw

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/adm87/jade/pkg/display"
	"github.com/adm87/jade/pkg/jade"
)

// Pipeline returns the game's drawer.
func Pipeline(d *display.Display) jade.Drawer {
	return display.Buffer(d, ebiten.FilterNearest)(
		jade.DrawAll(
			world(),
			corners(d),
			fps(d),
		),
	)
}

func world() jade.Drawer {
	img := ebiten.NewImage(10, 10)
	img.Fill(color.RGBA{R: 255, B: 255, A: 255})
	return jade.DrawFunc(func(dst *ebiten.Image) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(100, 100)
		dst.DrawImage(img, op)
	})
}

// corners marks each corner of the display's visible area.
func corners(d *display.Display) jade.Drawer {
	return jade.DrawFunc(func(dst *ebiten.Image) {
		visible := d.Visible()
		minX, minY := visible.Min()
		maxX, maxY := visible.Max()

		mark(minX, minY, color.RGBA{R: 255, A: 255}, dst)
		mark(maxX-10, minY, color.RGBA{R: 255, G: 255, A: 255}, dst)
		mark(minX, maxY-10, color.RGBA{G: 255, A: 255}, dst)
		mark(maxX-10, maxY-10, color.RGBA{B: 255, A: 255}, dst)
	})
}

func fps(d *display.Display) jade.Drawer {
	return jade.DrawFunc(func(dst *ebiten.Image) {
		minX, minY := d.Visible().Min()
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("FPS: %0.2f", ebiten.ActualTPS()), int(minX), int(minY))
	})
}

func mark(x, y float64, c color.Color, dst *ebiten.Image) {
	vector.FillRect(dst, float32(x), float32(y), 10, 10, c, false)
}
