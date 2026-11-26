package jade

import "github.com/hajimehoshi/ebiten/v2"

type shell struct {
	drawer Drawer
}

func (s *shell) Update() error {
	return nil
}

func (s *shell) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen)
}

func (s *shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
