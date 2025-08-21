package coordijk

import (
	"fmt"
	"math"
	"testing"

	"github.com/dimchansky/h3-go/internal/v2d"
)

// TestCoordIJKHex2dConversion tests 2D hex coordinate conversion functions.
func TestCoordIJKHex2dConversion(t *testing.T) {
	// IJK to Hex2d
	ijkToHex2dTests := []struct {
		name                 string
		input                CoordIJK
		expectedX, expectedY float64
	}{
		{"ijkToHex2d (0,0,0)", CoordIJK{0, 0, 0}, 0.0000000000, 0.0000000000},
		{"ijkToHex2d (1,0,0)", CoordIJK{1, 0, 0}, 1.0000000000, 0.0000000000},
		{"ijkToHex2d (0,1,0)", CoordIJK{0, 1, 0}, -0.5000000000, 0.8660254038},
		{"ijkToHex2d (0,0,1)", CoordIJK{0, 0, 1}, -0.5000000000, -0.8660254038},
		{"ijkToHex2d (1,1,0)", CoordIJK{1, 1, 0}, 0.5000000000, 0.8660254038},
		{"ijkToHex2d (-1,0,0)", CoordIJK{-1, 0, 0}, -1.0000000000, 0.0000000000},
		{"ijkToHex2d (0,-1,0)", CoordIJK{0, -1, 0}, 0.5000000000, -0.8660254038},
		{"ijkToHex2d (2,1,-1)", CoordIJK{2, 1, -1}, 2.0000000000, 1.7320508076},
		{"ijkToHex2d (3,-2,1)", CoordIJK{3, -2, 1}, 3.5000000000, -2.5980762114},
		{"ijkToHex2d (10,5,-3)", CoordIJK{10, 5, -3}, 9.0000000000, 6.9282032303},
	}

	for _, tt := range ijkToHex2dTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToHex2d()
			if math.Abs(result.X-tt.expectedX) > 1e-10 {
				t.Errorf("ToHex2d().X = %f, want %f", result.X, tt.expectedX)
			}
			if math.Abs(result.Y-tt.expectedY) > 1e-10 {
				t.Errorf("ToHex2d().Y = %f, want %f", result.Y, tt.expectedY)
			}
		})
	}

	// Hex2d to IJK
	hex2dToIJKTests := []struct {
		name           string
		inputX, inputY float64
		expected       CoordIJK
	}{
		{"hex2dToCoordIJK (0.000000,0.000000)", 0.0000000000, 0.0000000000, CoordIJK{0, 0, 0}},
		{"hex2dToCoordIJK (1.000000,0.000000)", 1.0000000000, 0.0000000000, CoordIJK{1, 0, 0}},
		{"hex2dToCoordIJK (0.000000,1.000000)", 0.0000000000, 1.0000000000, CoordIJK{1, 1, 0}},
		{"hex2dToCoordIJK (1.000000,1.000000)", 1.0000000000, 1.0000000000, CoordIJK{2, 1, 0}},
		{"hex2dToCoordIJK (-1.000000,0.000000)", -1.0000000000, 0.0000000000, CoordIJK{0, 1, 1}},
		{"hex2dToCoordIJK (0.000000,-1.000000)", 0.0000000000, -1.0000000000, CoordIJK{1, 0, 1}},
		{"hex2dToCoordIJK (-1.000000,-1.000000)", -1.0000000000, -1.0000000000, CoordIJK{0, 1, 2}},
		{"hex2dToCoordIJK (0.500000,0.866025)", 0.5000000000, 0.8660250000, CoordIJK{1, 1, 0}},
		{"hex2dToCoordIJK (2.500000,1.299038)", 2.5000000000, 1.2990380000, CoordIJK{3, 1, 0}},
		{"hex2dToCoordIJK (-2.500000,-1.299038)", -2.5000000000, -1.2990380000, CoordIJK{0, 2, 3}},
	}

	for _, tt := range hex2dToIJKTests {
		t.Run(tt.name, func(t *testing.T) {
			v := v2d.Vec2d{X: tt.inputX, Y: tt.inputY}
			result := Hex2dToCoordIJK(v)
			if result != tt.expected {
				t.Errorf("Hex2dToCoordIJK(%v) = %v, want %v", v, result, tt.expected)
			}
		})
	}

	// Equivalence through hex2d conversion
	equivalenceTests := []CoordIJK{
		{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0},
		{-1, 0, 0}, {0, -1, 0}, {2, 1, -1}, {3, -2, 1}, {10, 5, -3},
	}
	for i, original := range equivalenceTests {
		t.Run(fmt.Sprintf("hex2dEquivalence_%d", i), func(t *testing.T) {
			hex2d1 := original.ToHex2d()
			recovered := Hex2dToCoordIJK(hex2d1)
			hex2d2 := recovered.ToHex2d()
			if math.Abs(hex2d1.X-hex2d2.X) > 1e-10 || math.Abs(hex2d1.Y-hex2d2.Y) > 1e-10 {
				t.Errorf("Hex2d position not preserved: %v -> %v -> %v (hex2d: %v -> %v)", original, hex2d1, recovered, hex2d1, hex2d2)
			}
			if recovered != original {
				t.Logf("Coordinate normalized: %v -> %v (same hex position)", original, recovered)
			}
		})
	}
}
