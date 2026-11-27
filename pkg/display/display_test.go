package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adm87/jade/pkg/geom"
)

func rect(x, y, w, h float64) geom.Rectangle {
	return geom.Rectangle{
		Position: geom.Point{X: x, Y: y},
		Size:     geom.Size{Width: w, Height: h},
	}
}

func TestLayout(t *testing.T) {
	tests := []struct {
		name         string
		window, size geom.Size
		mode         ScaleMode
		bounds       geom.Rectangle
		visible      geom.Rectangle
		scale        float64
	}{
		{
			name:    "fit exact",
			window:  geom.Size{Width: 640, Height: 360},
			size:    geom.Size{Width: 320, Height: 180},
			mode:    Fit,
			bounds:  rect(0, 0, 640, 360),
			visible: rect(0, 0, 320, 180),
			scale:   2,
		},
		{
			name:    "fit letterbox",
			window:  geom.Size{Width: 1000, Height: 700},
			size:    geom.Size{Width: 320, Height: 180},
			mode:    Fit,
			bounds:  rect(0, 68, 1000, 562.5),
			visible: rect(0, 0, 320, 180),
			scale:   3.125,
		},
		{
			name:    "fit pillarbox",
			window:  geom.Size{Width: 800, Height: 360},
			size:    geom.Size{Width: 320, Height: 180},
			mode:    Fit,
			bounds:  rect(80, 0, 640, 360),
			visible: rect(0, 0, 320, 180),
			scale:   2,
		},
		{
			name:    "cover crops the sides",
			window:  geom.Size{Width: 320, Height: 200},
			size:    geom.Size{Width: 100, Height: 50},
			mode:    Cover,
			bounds:  rect(-40, 0, 400, 200),
			visible: rect(10, 0, 80, 50),
			scale:   4,
		},
		{
			name:    "integer rounds down",
			window:  geom.Size{Width: 1000, Height: 700},
			size:    geom.Size{Width: 320, Height: 180},
			mode:    Integer,
			bounds:  rect(20, 80, 960, 540),
			visible: rect(0, 0, 320, 180),
			scale:   3,
		},
		{
			name:    "integer never goes below 1",
			window:  geom.Size{Width: 160, Height: 90},
			size:    geom.Size{Width: 320, Height: 180},
			mode:    Integer,
			bounds:  rect(-80, -45, 320, 180),
			visible: rect(80, 45, 160, 90),
			scale:   1,
		},
		{
			name:   "zero window",
			window: geom.Size{},
			size:   geom.Size{Width: 320, Height: 180},
			mode:   Fit,
			scale:  1,
		},
		{
			name:   "zero size",
			window: geom.Size{Width: 640, Height: 360},
			size:   geom.Size{},
			mode:   Fit,
			scale:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bounds, visible, scale := layout(tt.window, tt.size, tt.mode)
			assert.Equal(t, tt.bounds, bounds, "bounds")
			assert.Equal(t, tt.visible, visible, "visible")
			assert.Equal(t, tt.scale, scale, "scale")
		})
	}
}

func TestScaleModes(t *testing.T) {
	target := geom.Size{Width: 100, Height: 50}
	tests := []struct {
		name   string
		mode   ScaleMode
		window geom.Size
		want   float64
	}{
		{"fit uses the smaller ratio", Fit, geom.Size{Width: 300, Height: 100}, 2},
		{"cover uses the larger ratio", Cover, geom.Size{Width: 300, Height: 100}, 3},
		{"integer floors fit", Integer, geom.Size{Width: 250, Height: 150}, 2},
		{"integer is at least 1", Integer, geom.Size{Width: 50, Height: 25}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mode(tt.window, target))
		})
	}
}

func TestDisplayZeroValueFits(t *testing.T) {
	var d Display
	mode := d.scaleMode()
	require.NotNil(t, mode)
	assert.Equal(t, 2.0, mode(geom.Size{Width: 300, Height: 100}, geom.Size{Width: 100, Height: 50}))
}
