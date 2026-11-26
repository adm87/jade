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
