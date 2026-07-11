//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_square_ParityWithC(t *testing.T) {
	testCases := []float64{
		// Basic cases
		0.0,
		1.0,
		-1.0,
		2.0,
		-2.0,

		// Fractional values
		0.5,
		-0.5,
		0.25,
		-0.75,

		// Small values
		1e-10,
		-1e-10,
		1e-15,
		-1e-15,

		// Large values
		1000.0,
		-1000.0,
		1e6,
		-1e6,

		// Mathematical constants
		math.Pi,
		-math.Pi,
		math.E,
		-math.E,
		math.Sqrt2,
		-math.Sqrt2,

		// Edge cases
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),

		// Values around 1
		0.99999,
		1.00001,
		-0.99999,
		-1.00001,

		// Random values for good coverage
		3.14159,
		-2.71828,
		123.456,
		-987.654,
		0.123456789,
		-0.987654321,
	}

	for _, x := range testCases {
		goResult := _square(x)
		cResult := _squareC(x)

		// Handle special values
		if math.IsNaN(x) {
			if !math.IsNaN(goResult) || !math.IsNaN(cResult) {
				t.Fatalf("_square NaN handling mismatch: x=%.15g go=%.15g c=%.15g", x, goResult, cResult)
			}
			continue
		}

		if math.IsInf(x, 0) {
			if !math.IsInf(goResult, 1) || !math.IsInf(cResult, 1) {
				t.Fatalf("_square Inf handling mismatch: x=%.15g go=%.15g c=%.15g", x, goResult, cResult)
			}
			continue
		}

		// For normal values, check exact equality or very tight tolerance
		if goResult != cResult && math.Abs(goResult-cResult) > 1e-15 {
			t.Fatalf("_square mismatch: x=%.15g go=%.15g c=%.15g diff=%.15g", x, goResult, cResult, math.Abs(goResult-cResult))
		}
	}
}
