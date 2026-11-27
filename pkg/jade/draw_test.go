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

func TestChain(t *testing.T) {
	var got []string
	pass := func(name string) jade.DrawPass {
		return func(next jade.Drawer) jade.Drawer {
			return jade.DrawFunc(func(dst *ebiten.Image) {
				got = append(got, name)
				next.Draw(dst)
			})
		}
	}
	world := jade.DrawFunc(func(*ebiten.Image) { got = append(got, "world") })

	tests := []struct {
		name   string
		passes []jade.DrawPass
		want   []string
	}{
		{"first is outermost", []jade.DrawPass{pass("a"), pass("b")}, []string{"a", "b", "world"}},
		{"skips nil", []jade.DrawPass{pass("a"), nil, pass("b")}, []string{"a", "b", "world"}},
		{"empty", nil, []string{"world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got = nil
			jade.Chain(tt.passes...)(world).Draw(new(ebiten.Image))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDrawAll(t *testing.T) {
	screen := new(ebiten.Image)

	var got []string
	record := func(name string) jade.Drawer {
		return jade.DrawFunc(func(dst *ebiten.Image) {
			assert.Same(t, screen, dst)
			got = append(got, name)
		})
	}

	tests := []struct {
		name    string
		drawers []jade.Drawer
		want    []string
	}{
		{"draws in order", []jade.Drawer{record("a"), record("b"), record("c")}, []string{"a", "b", "c"}},
		{"skips nil", []jade.Drawer{record("a"), nil, record("b")}, []string{"a", "b"}},
		{"empty", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got = nil
			jade.DrawAll(tt.drawers...).Draw(screen)
			assert.Equal(t, tt.want, got)
		})
	}
}
