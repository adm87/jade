package display

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/adm87/jade/pkg/geom"
)

// Display calculates the bounds and visible area of a fixed logical size in the window.
// The zero value has a zero size and scales with Fit.
type Display struct {
	options  *options
	visible  geom.Rectangle
	bounds   geom.Rectangle
	size     geom.Size
	physical geom.Size
	scale    float64
}

// New returns a Display with the given logical size.
func New(width, height int, opts ...Option) *Display {
	return &Display{
		options: applyOptions(opts),
		size: geom.Size{
			Width:  float64(width),
			Height: float64(height),
		},
		scale: 1,
	}
}

// Visible is the part of the logical area inside the window, in logical pixels.
func (d *Display) Visible() geom.Rectangle {
	return d.visible
}

// ToLogical converts screen coordinates in physical pixels to logical coordinates.
func (d *Display) ToLogical(screenX, screenY float64) (x, y float64) {
	if d.scale == 0 {
		return screenX, screenY
	}
	return (screenX - d.bounds.Position.X) / d.scale, (screenY - d.bounds.Position.Y) / d.scale
}

// ToScreen converts logical coordinates to screen coordinates in physical pixels.
func (d *Display) ToScreen(x, y float64) (screenX, screenY float64) {
	return x*d.scale + d.bounds.Position.X, y*d.scale + d.bounds.Position.Y
}

// Layout updates the scale and placement for the window size, and returns that size in physical pixels.
func (d *Display) Layout(outsideWidth, outsideHeight float64) (width, height float64) {
	scale := deviceScale()

	physical := geom.Size{
		Width:  outsideWidth * scale,
		Height: outsideHeight * scale,
	}

	if physical != d.physical {
		d.physical = physical
		d.bounds, d.visible, d.scale = layout(d.physical, d.size, d.scaleMode())
	}

	return d.physical.Width, d.physical.Height
}

func (d *Display) scaleMode() ScaleMode {
	if d.options == nil {
		return Fit
	}
	return d.options.scaleMode
}

func deviceScale() float64 {
	if mon := ebiten.Monitor(); mon != nil {
		return mon.DeviceScaleFactor()
	}
	return 1
}

func layout(window, size geom.Size, mode ScaleMode) (bounds, visible geom.Rectangle, scale float64) {
	if window.Width <= 0 || window.Height <= 0 || size.Width <= 0 || size.Height <= 0 {
		return geom.Rectangle{}, geom.Rectangle{}, 1
	}

	scale = mode(window, size)

	scaled := geom.Size{
		Width:  size.Width * scale,
		Height: size.Height * scale,
	}

	position := geom.Point{
		X: math.Floor((window.Width - scaled.Width) / 2.0),
		Y: math.Floor((window.Height - scaled.Height) / 2.0),
	}

	visible = geom.Rectangle{
		Position: geom.Point{
			X: max(0, -position.X) / scale,
			Y: max(0, -position.Y) / scale,
		},
		Size: geom.Size{
			Width:  min(scaled.Width, window.Width) / scale,
			Height: min(scaled.Height, window.Height) / scale,
		},
	}

	bounds = geom.Rectangle{
		Position: position,
		Size:     scaled,
	}

	return bounds, visible, scale
}
