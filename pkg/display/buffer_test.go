package display

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adm87/jade/pkg/geom"
	"github.com/adm87/jade/pkg/jade"
)

func TestBuffer(t *testing.T) {
	d := &Display{
		size:   geom.Size{Width: 32, Height: 18},
		bounds: rect(0, 0, 64, 36),
	}

	var got []*ebiten.Image
	next := jade.DrawFunc(func(dst *ebiten.Image) { got = append(got, dst) })

	screen := ebiten.NewImage(64, 36)
	drawer := Buffer(d, ebiten.FilterNearest)(next)

	drawer.Draw(screen)
	require.Len(t, got, 1)
	assert.NotSame(t, screen, got[0], "draws into its own image")
	assert.Equal(t, 32, got[0].Bounds().Dx())
	assert.Equal(t, 18, got[0].Bounds().Dy())

	drawer.Draw(screen)
	require.Len(t, got, 2)
	assert.Same(t, got[0], got[1], "reuses the image while the size holds")

	d.size = geom.Size{Width: 16, Height: 9}
	drawer.Draw(screen)
	require.Len(t, got, 3)
	assert.Equal(t, 16, got[2].Bounds().Dx(), "resizes with the display")
	assert.Equal(t, 9, got[2].Bounds().Dy())
}

func TestBufferSkipsEmpty(t *testing.T) {
	tests := []struct {
		name    string
		display *Display
	}{
		{"zero size", &Display{bounds: rect(0, 0, 64, 36)}},
		{"zero bounds", &Display{size: geom.Size{Width: 32, Height: 18}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := jade.DrawFunc(func(*ebiten.Image) { called = true })

			Buffer(tt.display, ebiten.FilterNearest)(next).Draw(ebiten.NewImage(1, 1))
			assert.False(t, called)
		})
	}
}
