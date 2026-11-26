package jade_test

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"

	"github.com/adm87/jade/pkg/jade"
)

func ExampleDrawPass() {
	trace := func(name string) jade.DrawPass {
		return func(next jade.Drawer) jade.Drawer {
			return jade.DrawFunc(func(dst *ebiten.Image) {
				fmt.Println(name, "before")
				next.Draw(dst)
				fmt.Println(name, "after")
			})
		}
	}
	world := jade.DrawFunc(func(*ebiten.Image) { fmt.Println("world") })

	outer, inner := trace("outer"), trace("inner")
	outer(inner(world)).Draw(nil)
}

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
