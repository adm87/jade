package geom

// Rectangle is an area with its top-left corner at Position.
type Rectangle struct {
	Position Point
	Size     Size
}

// Min is the top-left corner.
func (r Rectangle) Min() (x, y float64) {
	return r.Position.X, r.Position.Y
}

// Max is the bottom-right corner.
func (r Rectangle) Max() (x, y float64) {
	return r.Position.X + r.Size.Width, r.Position.Y + r.Size.Height
}
