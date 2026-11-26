package jade_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"

	"github.com/adm87/jade/pkg/jade"
)

func TestDrawPassRedirectsTarget(t *testing.T) {
	screen, buf := new(ebiten.Image), new(ebiten.Image)

	redirect := func(next jade.Drawer) jade.Drawer {
		return jade.DrawFunc(func(*ebiten.Image) { next.Draw(buf) })
	}

	var got *ebiten.Image
	world := jade.DrawFunc(func(dst *ebiten.Image) { got = dst })

	jade.DrawPass(redirect)(world).Draw(screen)
	assert.Same(t, buf, got)
}
