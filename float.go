package h3

//lint:file-ignore U1000 Ignore all unused code

import (
	"math"
)

// square of a number
// `x`: The input number.
// Returns the square of the input number.
func square(x float64) float64 { return x * x }

// isFinite reports whether f is neither NaN nor an infinity.
func isFinite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}
