package jade

// Layouter handles the layout of the game window.
type Layouter interface {
	Layout(outsideWidth, outsideHeight int) (int, int)
}

// LayoutFunc lets an ordinary function be used as a Layouter.
type LayoutFunc func(outsideWidth, outsideHeight int) (int, int)

// Layout calls f(outsideWidth, outsideHeight).
func (f LayoutFunc) Layout(outsideWidth, outsideHeight int) (int, int) {
	return f(outsideWidth, outsideHeight)
}
