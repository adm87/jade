package jade

import "github.com/hajimehoshi/ebiten/v2"

type options struct {
	drawer Drawer
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

func defaultOptions() *options {
	return &options{
		drawer: DrawFunc(func(_ *ebiten.Image) {}),
	}
}

func applyOptions(opts []Option) *options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}
