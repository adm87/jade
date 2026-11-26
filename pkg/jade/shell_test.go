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
	var gotW, gotH int
	s := shell{layouter: LayoutFunc(func(w, h int) (int, int) {
		gotW, gotH = w, h
		return 320, 180
	})}

	w, h := s.Layout(1280, 720)
	assert.Equal(t, []int{1280, 720}, []int{gotW, gotH}, "window size passed in")
	assert.Equal(t, []int{320, 180}, []int{w, h}, "layouter's size returned")
}
