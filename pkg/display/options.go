package display

type options struct {
	scaleMode ScaleMode
}

// Option configures a Display.
type Option func(*options)

// WithScaleMode sets how the logical size is scaled to the window. A nil ScaleMode is ignored.
func WithScaleMode(scaleMode ScaleMode) Option {
	return func(o *options) {
		if scaleMode != nil {
			o.scaleMode = scaleMode
		}
	}
}

func defaultOptions() *options {
	return &options{
		scaleMode: Fit,
	}
}

func applyOptions(opts []Option) *options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}
