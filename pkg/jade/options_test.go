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

func TestApplyOptions(t *testing.T) {
	d, u := &fakeDrawer{}, &fakeUpdater{}

	t.Run("defaults are safe no-ops", func(t *testing.T) {
		o := applyOptions(nil)
		assert.NotPanics(t, func() { o.drawer.Draw(nil) })
		assert.NoError(t, o.updater.Update(time.Second))
	})

	t.Run("sets drawer and updater", func(t *testing.T) {
		o := applyOptions([]Option{WithDrawer(d), WithUpdater(u)})
		assert.Same(t, d, o.drawer)
		assert.Same(t, u, o.updater)
	})

	t.Run("ignores nil", func(t *testing.T) {
		o := applyOptions([]Option{WithDrawer(nil), WithUpdater(nil)})
		assert.NotNil(t, o.drawer)
		assert.NotNil(t, o.updater)
	})

	t.Run("last option wins", func(t *testing.T) {
		o := applyOptions([]Option{WithDrawer(&fakeDrawer{}), WithDrawer(d)})
		assert.Same(t, d, o.drawer)
	})
}
