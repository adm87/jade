package jade

import (
	"errors"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestShellDt(t *testing.T) {
	t.Run("zero on first call", func(t *testing.T) {
		var s shell
		assert.Zero(t, s.dt())
	})

	t.Run("time since previous call", func(t *testing.T) {
		s := shell{last: time.Now().Add(-10 * time.Millisecond)}
		dt := s.dt()
		assert.GreaterOrEqual(t, dt, 10*time.Millisecond)
		assert.Less(t, dt, maxDelta)
	})

	t.Run("capped at maxDelta", func(t *testing.T) {
		s := shell{last: time.Now().Add(-time.Hour)}
		assert.Equal(t, maxDelta, s.dt())
	})
}

func TestShellUpdate(t *testing.T) {
	errBoom := errors.New("boom")
	var got []time.Duration
	s := shell{updater: UpdateFunc(func(dt time.Duration) error {
		got = append(got, dt)
		return errBoom
	})}

	assert.ErrorIs(t, s.Update(), errBoom)
	assert.Equal(t, []time.Duration{0}, got, "first tick gets dt 0")
}

func TestShellDraw(t *testing.T) {
	target := new(ebiten.Image)
	var got *ebiten.Image
	s := shell{drawer: DrawFunc(func(dst *ebiten.Image) { got = dst })}

	s.Draw(target)
	assert.Same(t, target, got)
}

func TestShellLayout(t *testing.T) {
	var gotW, gotH float64
	s := shell{layouter: LayoutFunc(func(w, h float64) (float64, float64) {
		gotW, gotH = w, h
		return 320.5, 180.5
	})}

	w, h := s.LayoutF(1280, 720)
	assert.Equal(t, []float64{1280, 720}, []float64{gotW, gotH}, "window size passed in")
	assert.Equal(t, []float64{320.5, 180.5}, []float64{w, h}, "layouter's size returned")

	iw, ih := s.Layout(1280, 720)
	assert.Equal(t, []int{320, 180}, []int{iw, ih}, "Layout truncates LayoutF")
}
