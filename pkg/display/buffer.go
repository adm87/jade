package display

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/adm87/jade/pkg/geom"
	"github.com/adm87/jade/pkg/jade"
)

// Buffer returns a pass that draws into an image of the display's logical size,
// and then draws that image into the destination with the given filter.
func Buffer(d *Display, filter ebiten.Filter) jade.DrawPass {
	return func(next jade.Drawer) jade.Drawer {
		var buffer *ebiten.Image
		return jade.DrawFunc(func(dst *ebiten.Image) {
			w, h := pixels(d.size)

			if w <= 0 || h <= 0 || d.bounds.Size.Width <= 0 || d.bounds.Size.Height <= 0 {
				return
			}

			if buffer == nil || buffer.Bounds().Dx() != w || buffer.Bounds().Dy() != h {
				if buffer != nil {
					buffer.Deallocate()
				}
				buffer = ebiten.NewImage(w, h)
			} else {
				buffer.Clear()
			}

			next.Draw(buffer)

			op := &ebiten.DrawImageOptions{Filter: filter}
			op.GeoM.Scale(d.bounds.Size.Width/float64(w), d.bounds.Size.Height/float64(h))
			op.GeoM.Translate(d.bounds.Position.X, d.bounds.Position.Y)

			dst.DrawImage(buffer, op)
		})
	}
}

func pixels(size geom.Size) (int, int) {
	return int(size.Width), int(size.Height)
}
