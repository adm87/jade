package jade

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

type fakeDrawer struct{}

func (*fakeDrawer) Draw(*ebiten.Image) {}

type fakeUpdater struct{}

func (*fakeUpdater) Update(time.Duration) error { return nil }

type fakeLayouter struct{}

func (*fakeLayouter) Layout(w, h int) (int, int) { return w, h }

func TestApplyOptions(t *testing.T) {
	d, u, l := &fakeDrawer{}, &fakeUpdater{}, &fakeLayouter{}

	t.Run("defaults are safe no-ops", func(t *testing.T) {
		o := applyOptions(nil)
		assert.NotPanics(t, func() { o.drawer.Draw(nil) })
		assert.NoError(t, o.updater.Update(time.Second))
	})

	t.Run("default layout matches the window", func(t *testing.T) {
		w, h := applyOptions(nil).layouter.Layout(640, 360)
		assert.Equal(t, 640, w)
		assert.Equal(t, 360, h)
	})

	t.Run("sets drawer, updater and layouter", func(t *testing.T) {
		o := applyOptions([]Option{WithDrawer(d), WithUpdater(u), WithLayouter(l)})
		assert.Same(t, d, o.drawer)
		assert.Same(t, u, o.updater)
		assert.Same(t, l, o.layouter)
	})

	t.Run("ignores nil", func(t *testing.T) {
		o := applyOptions([]Option{WithDrawer(nil), WithUpdater(nil), WithLayouter(nil)})
		assert.NotNil(t, o.drawer)
		assert.NotNil(t, o.updater)
		assert.NotNil(t, o.layouter)
	})

	t.Run("last option wins", func(t *testing.T) {
		o := applyOptions([]Option{WithDrawer(&fakeDrawer{}), WithDrawer(d)})
		assert.Same(t, d, o.drawer)
	})
}
