package jade

import "github.com/hajimehoshi/ebiten/v2"

type shell struct{}

func (s shell) Update() error {
	return nil
}

func (s shell) Draw(_ *ebiten.Image) {

}

func (s shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
