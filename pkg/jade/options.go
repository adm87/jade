package jade

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type options struct {
	drawer   Drawer
	updater  Updater
	layouter Layouter
}

// Option configures Run.
type Option func(*options)

// WithDrawer sets the Drawer that receives the screen each frame.
// A nil Drawer is ignored.
func WithDrawer(d Drawer) Option {
	return func(o *options) {
		if d != nil {
			o.drawer = d
		}
	}
}

// WithUpdater sets the Updater called each tick.
// A nil Updater is ignored.
func WithUpdater(u Updater) Option {
	return func(o *options) {
		if u != nil {
			o.updater = u
		}
	}
}

// WithLayouter sets the Layouter that decides the game's screen size from the window size.
// A nil Layouter is ignored.
func WithLayouter(l Layouter) Option {
	return func(o *options) {
		if l != nil {
			o.layouter = l
		}
	}
}

func defaultOptions() *options {
	return &options{
		drawer: DrawFunc(func(_ *ebiten.Image) {}),
		updater: UpdateFunc(func(_ time.Duration) error {
			return nil
		}),
		layouter: LayoutFunc(func(outsideWidth, outsideHeight float64) (float64, float64) {
			return outsideWidth, outsideHeight
		}),
	}
}

func applyOptions(opts []Option) *options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}
