package jade

// Layouter handles the layout of the game window.
// Sizes are float64.
type Layouter interface {
	Layout(outsideWidth, outsideHeight float64) (width, height float64)
}

// LayoutFunc lets an ordinary function be used as a Layouter.
type LayoutFunc func(outsideWidth, outsideHeight float64) (width, height float64)

// Layout calls f(outsideWidth, outsideHeight).
func (f LayoutFunc) Layout(outsideWidth, outsideHeight float64) (width, height float64) {
	return f(outsideWidth, outsideHeight)
}
