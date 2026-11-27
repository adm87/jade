package jade

import "github.com/hajimehoshi/ebiten/v2"

// Drawer draws onto an image.
type Drawer interface {
	Draw(dst *ebiten.Image)
}

// DrawFunc lets an ordinary function be used as a Drawer.
type DrawFunc func(dst *ebiten.Image)

// Draw calls f(dst).
func (f DrawFunc) Draw(dst *ebiten.Image) {
	f(dst)
}

// DrawPass wraps a Drawer to add behavior before or after it.
type DrawPass func(Drawer) Drawer

// Chain applies passes in order, with the first one outermost.
// Nil passes are skipped.
func Chain(passes ...DrawPass) DrawPass {
	return func(next Drawer) Drawer {
		for i := len(passes) - 1; i >= 0; i-- {
			if passes[i] != nil {
				next = passes[i](next)
			}
		}
		return next
	}
}

// DrawAll draws each drawer in order onto the same image.
// Nil drawers are skipped.
func DrawAll(drawers ...Drawer) Drawer {
	return DrawFunc(func(dst *ebiten.Image) {
		for _, d := range drawers {
			if d != nil {
				d.Draw(dst)
			}
		}
	})
}
