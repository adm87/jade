package display

import (
	"math"

	"github.com/adm87/jade/pkg/geom"
)

// ScaleMode returns a scale for target in window.
type ScaleMode func(window, target geom.Size) float64

// Fit returns the largest scale at which target fits inside window.
func Fit(window, target geom.Size) float64 {
	return min(window.Width/target.Width, window.Height/target.Height)
}

// Cover returns the smallest scale at which target fills window.
func Cover(window, target geom.Size) float64 {
	return max(window.Width/target.Width, window.Height/target.Height)
}

// Integer returns Fit rounded down to a whole number, and at least 1.
func Integer(window, target geom.Size) float64 {
	return max(1, math.Floor(Fit(window, target)))
}
