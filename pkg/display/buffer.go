package display

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/adm87/jade/pkg/geom"
	"github.com/adm87/jade/pkg/jade"
)

// Viewport is a logical size and where it's placed in the window.
type Viewport interface {
	Size() geom.Size
	Bounds() geom.Rectangle
}

// Buffer returns a pass that draws into an image of the viewport's size, then scales it to the viewport's bounds.
func Buffer(viewport Viewport, filter ebiten.Filter) jade.DrawPass {
	return func(next jade.Drawer) jade.Drawer {
		var buffer *ebiten.Image
		return jade.DrawFunc(func(dst *ebiten.Image) {
			w, h := pixels(viewport.Size())

			bounds := viewport.Bounds()
			if w <= 0 || h <= 0 || bounds.Size.Width <= 0 || bounds.Size.Height <= 0 {
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
			op.GeoM.Scale(bounds.Size.Width/float64(w), bounds.Size.Height/float64(h))
			op.GeoM.Translate(bounds.Position.X, bounds.Position.Y)

			dst.DrawImage(buffer, op)
		})
	}
}

func pixels(size geom.Size) (int, int) {
	return int(size.Width), int(size.Height)
}
