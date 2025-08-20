package coordijk

import (
	"fmt"
	"math"
	"testing"
)

// TestCoordIJKNormalizeCValidation validates Go implementation against H3 C test cases
func TestCoordIJKNormalizeCValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    CoordIJK
		expected CoordIJK
	}{
		{"normalize (0,0,0)", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"normalize (1,0,0)", CoordIJK{1, 0, 0}, CoordIJK{1, 0, 0}},
		{"normalize (0,1,0)", CoordIJK{0, 1, 0}, CoordIJK{0, 1, 0}},
		{"normalize (0,0,1)", CoordIJK{0, 0, 1}, CoordIJK{0, 0, 1}},
		{"normalize (1,1,0)", CoordIJK{1, 1, 0}, CoordIJK{1, 1, 0}},
		{"normalize (1,0,1)", CoordIJK{1, 0, 1}, CoordIJK{1, 0, 1}},
		{"normalize (0,1,1)", CoordIJK{0, 1, 1}, CoordIJK{0, 1, 1}},
		{"normalize (1,1,1)", CoordIJK{1, 1, 1}, CoordIJK{0, 0, 0}},
		{"normalize (2,2,2)", CoordIJK{2, 2, 2}, CoordIJK{0, 0, 0}},
		{"normalize (3,1,2)", CoordIJK{3, 1, 2}, CoordIJK{2, 0, 1}},
		{"normalize (-1,0,0)", CoordIJK{-1, 0, 0}, CoordIJK{0, 1, 1}},
		{"normalize (0,-1,0)", CoordIJK{0, -1, 0}, CoordIJK{1, 0, 1}},
		{"normalize (0,0,-1)", CoordIJK{0, 0, -1}, CoordIJK{1, 1, 0}},
		{"normalize (-1,-1,-1)", CoordIJK{-1, -1, -1}, CoordIJK{0, 0, 0}},
		{"normalize (5,-2,3)", CoordIJK{5, -2, 3}, CoordIJK{7, 0, 5}},
		{"normalize (-3,4,1)", CoordIJK{-3, 4, 1}, CoordIJK{0, 7, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input
			result.Normalize()
			if result != tt.expected {
				t.Errorf("CoordIJK.Normalize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestUnitIJKToDigitCValidation validates Go implementation against H3 C test cases
func TestUnitIJKToDigitCValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    CoordIJK
		expected Direction
	}{
		{"CENTER_DIGIT", CoordIJK{0, 0, 0}, 0},
		{"K_AXES_DIGIT", CoordIJK{0, 0, 1}, 1},
		{"J_AXES_DIGIT", CoordIJK{0, 1, 0}, 2},
		{"JK_AXES_DIGIT", CoordIJK{0, 1, 1}, 3},
		{"I_AXES_DIGIT", CoordIJK{1, 0, 0}, 4},
		{"IK_AXES_DIGIT", CoordIJK{1, 0, 1}, 5},
		{"IJ_AXES_DIGIT", CoordIJK{1, 1, 0}, 6},
		{"INVALID_DIGIT (out of range)", CoordIJK{2, 0, 0}, 7},
		{"CENTER_DIGIT (unnormalized zero)", CoordIJK{2, 2, 2}, 0},
		{"CENTER_DIGIT (unnormalized zero)", CoordIJK{3, 3, 3}, 0},
		{"INVALID_DIGIT (no match)", CoordIJK{-1, 1, 0}, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnitIJKToDigit(tt.input)
			if result != tt.expected {
				t.Errorf("UnitIJKToDigit(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCoordIJKBasicOperations tests Add, Sub, Scale operations with C-derived expected results
func TestCoordIJKBasicOperations(t *testing.T) {
	tests := []struct {
		name string
		a, b CoordIJK
		expectedAdd CoordIJK
		expectedSub CoordIJK
		scaleFactor int
		expectedScale CoordIJK
	}{
		{"zero coordinates", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, 2, CoordIJK{0, 0, 0}},
		{"simple addition", CoordIJK{1, 0, 0}, CoordIJK{0, 1, 0}, CoordIJK{1, 1, 0}, CoordIJK{1, -1, 0}, 3, CoordIJK{3, 0, 0}},
		{"negative coordinates", CoordIJK{-1, 2, 0}, CoordIJK{1, -1, 1}, CoordIJK{0, 1, 1}, CoordIJK{-2, 3, -1}, -2, CoordIJK{2, -4, 0}},
		{"large coordinates", CoordIJK{5, -3, 2}, CoordIJK{-2, 4, -1}, CoordIJK{3, 1, 1}, CoordIJK{7, -7, 3}, 4, CoordIJK{20, -12, 8}},
		{"unit vectors", CoordIJK{1, 1, 0}, CoordIJK{0, 0, 1}, CoordIJK{1, 1, 1}, CoordIJK{1, 1, -1}, 5, CoordIJK{5, 5, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultAdd := tt.a.Add(tt.b)
			if resultAdd != tt.expectedAdd {
				t.Errorf("Add() = %v, want %v", resultAdd, tt.expectedAdd)
			}

			resultSub := tt.a.Sub(tt.b)
			if resultSub != tt.expectedSub {
				t.Errorf("Sub() = %v, want %v", resultSub, tt.expectedSub)
			}

			resultScale := tt.a.Scale(tt.scaleFactor)
			if resultScale != tt.expectedScale {
				t.Errorf("Scale() = %v, want %v", resultScale, tt.expectedScale)
			}
		})
	}
}

// TestCoordIJKDistance tests Distance function based on H3 C testGridDistanceInternal.c  
func TestCoordIJKDistance(t *testing.T) {
	tests := []struct {
		name string
		a, b CoordIJK
		expected int
	}{
		// Identity tests - same coordinate should have distance 0
		{"identity distance 0,0,0", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, 0},
		{"identity distance 1,0,0", CoordIJK{1, 0, 0}, CoordIJK{1, 0, 0}, 0},
		{"identity distance 1,0,1", CoordIJK{1, 0, 1}, CoordIJK{1, 0, 1}, 0},
		{"identity distance 1,1,0", CoordIJK{1, 1, 0}, CoordIJK{1, 1, 0}, 0},
		{"identity distance 0,2,0", CoordIJK{0, 2, 0}, CoordIJK{0, 2, 0}, 0},
		// Test cases directly from H3 C testGridDistanceInternal.c expected results
		{"0,0,0 to 1,0,0", CoordIJK{0, 0, 0}, CoordIJK{1, 0, 0}, 1},
		{"0,0,0 to 0,2,0", CoordIJK{0, 0, 0}, CoordIJK{0, 2, 0}, 2},
		{"0,0,0 to 1,0,1", CoordIJK{0, 0, 0}, CoordIJK{1, 0, 1}, 0},
		{"1,0,0 to 1,0,1", CoordIJK{1, 0, 0}, CoordIJK{1, 0, 1}, 1},
		{"1,0,1 to 0,2,0", CoordIJK{1, 0, 1}, CoordIJK{0, 2, 0}, 3},
		{"1,0,1 to 1,1,0", CoordIJK{1, 0, 1}, CoordIJK{1, 1, 0}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Distance(tt.a, tt.b)
			if result != tt.expected {
				// Output actual vs expected for debugging
				t.Logf("Distance calculation debug: (%v) - (%v) = %v", tt.a, tt.b, tt.a.Sub(tt.b))
				t.Errorf("Distance(%v, %v) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestCoordIJKNeighbor tests Neighbor function based on H3 C tests
func TestCoordIJKNeighbor(t *testing.T) {
	tests := []struct {
		name string
		input CoordIJK
		direction Direction
		expected CoordIJK
	}{
		{"center neighbor from origin", CoordIJK{0, 0, 0}, 0, CoordIJK{0, 0, 0}},
		{"i-axis neighbor from origin", CoordIJK{0, 0, 0}, 4, CoordIJK{1, 0, 0}},
		{"j-axis neighbor from origin", CoordIJK{0, 0, 0}, 2, CoordIJK{0, 1, 0}},
		{"k-axis neighbor from origin", CoordIJK{0, 0, 0}, 1, CoordIJK{0, 0, 1}},
		{"invalid direction", CoordIJK{1, 0, 0}, 7, CoordIJK{1, 0, 0}},
		{"center from non-origin", CoordIJK{2, 1, -1}, 0, CoordIJK{3, 2, 0}},
		{"multiple steps", CoordIJK{1, 0, 0}, 1, CoordIJK{1, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input
			result.Neighbor(tt.direction)
			if result != tt.expected {
				t.Errorf("Neighbor(%d) = %v, want %v", tt.direction, result, tt.expected)
			}
		})
	}
}

// TestCoordIJKRotation tests Rotate60CCW and Rotate60CW functions
func TestCoordIJKRotation(t *testing.T) {
	tests := []struct {
		name string
		input CoordIJK
		expectedCCW CoordIJK
		expectedCW CoordIJK
	}{
		{"rotate (0,0,0)", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"rotate (1,0,0)", CoordIJK{1, 0, 0}, CoordIJK{0, -1, 0}, CoordIJK{0, 0, -1}},
		{"rotate (0,1,0)", CoordIJK{0, 1, 0}, CoordIJK{0, 0, -1}, CoordIJK{-1, 0, 0}},
		{"rotate (0,0,1)", CoordIJK{0, 0, 1}, CoordIJK{-1, 0, 0}, CoordIJK{0, -1, 0}},
		{"rotate (1,1,0)", CoordIJK{1, 1, 0}, CoordIJK{0, -1, -1}, CoordIJK{-1, 0, -1}},
		{"rotate (1,0,1)", CoordIJK{1, 0, 1}, CoordIJK{-1, -1, 0}, CoordIJK{0, -1, -1}},
		{"rotate (0,1,1)", CoordIJK{0, 1, 1}, CoordIJK{-1, 0, -1}, CoordIJK{-1, -1, 0}},
		{"rotate (2,-1,1)", CoordIJK{2, -1, 1}, CoordIJK{-1, -2, 1}, CoordIJK{1, -1, -2}},
		{"rotate (-1,2,-1)", CoordIJK{-1, 2, -1}, CoordIJK{1, 1, -2}, CoordIJK{-2, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultCCW := tt.input
			resultCCW.Rotate60CCW()
			if resultCCW != tt.expectedCCW {
				t.Errorf("Rotate60CCW() = %v, want %v", resultCCW, tt.expectedCCW)
			}

			resultCW := tt.input
			resultCW.Rotate60CW()
			if resultCW != tt.expectedCW {
				t.Errorf("Rotate60CW() = %v, want %v", resultCW, tt.expectedCW)
			}
		})
	}
}

// TestCoordIJKApertureTransforms tests aperture transformation functions with C-derived expected results
func TestCoordIJKApertureTransforms(t *testing.T) {
	tests := []struct {
		name string
		input CoordIJK
		expectedUpAp7 CoordIJK
		expectedUpAp7r CoordIJK
		expectedDownAp7 CoordIJK
		expectedDownAp7r CoordIJK
		expectedDownAp3 CoordIJK
		expectedDownAp3r CoordIJK
	}{
		{"aperture (0,0,0)", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"aperture (7,0,0)", CoordIJK{7, 0, 0}, CoordIJK{3, 1, 0}, CoordIJK{3, 0, 1}, CoordIJK{21, 0, 7}, CoordIJK{21, 7, 0}, CoordIJK{14, 0, 7}, CoordIJK{14, 7, 0}},
		{"aperture (0,7,0)", CoordIJK{0, 7, 0}, CoordIJK{0, 3, 1}, CoordIJK{1, 3, 0}, CoordIJK{7, 21, 0}, CoordIJK{0, 21, 7}, CoordIJK{7, 14, 0}, CoordIJK{0, 14, 7}},
		{"aperture (1,0,0)", CoordIJK{1, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{3, 0, 1}, CoordIJK{3, 1, 0}, CoordIJK{2, 0, 1}, CoordIJK{2, 1, 0}},
		{"aperture (0,1,0)", CoordIJK{0, 1, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{1, 3, 0}, CoordIJK{0, 3, 1}, CoordIJK{1, 2, 0}, CoordIJK{0, 2, 1}},
		{"aperture (0,0,1)", CoordIJK{0, 0, 1}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 1, 3}, CoordIJK{1, 0, 3}, CoordIJK{0, 1, 2}, CoordIJK{1, 0, 2}},
		{"aperture (2,1,0)", CoordIJK{2, 1, 0}, CoordIJK{1, 1, 0}, CoordIJK{1, 0, 0}, CoordIJK{5, 1, 0}, CoordIJK{5, 4, 0}, CoordIJK{3, 0, 0}, CoordIJK{3, 3, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upAp7 := tt.input
			upAp7.UpAp7()
			if upAp7 != tt.expectedUpAp7 {
				t.Errorf("UpAp7() = %v, want %v", upAp7, tt.expectedUpAp7)
			}

			upAp7r := tt.input
			upAp7r.UpAp7r()
			if upAp7r != tt.expectedUpAp7r {
				t.Errorf("UpAp7r() = %v, want %v", upAp7r, tt.expectedUpAp7r)
			}

			downAp7 := tt.input
			downAp7.DownAp7()
			if downAp7 != tt.expectedDownAp7 {
				t.Errorf("DownAp7() = %v, want %v", downAp7, tt.expectedDownAp7)
			}

			downAp7r := tt.input
			downAp7r.DownAp7r()
			if downAp7r != tt.expectedDownAp7r {
				t.Errorf("DownAp7r() = %v, want %v", downAp7r, tt.expectedDownAp7r)
			}

			downAp3 := tt.input
			downAp3.DownAp3()
			if downAp3 != tt.expectedDownAp3 {
				t.Errorf("DownAp3() = %v, want %v", downAp3, tt.expectedDownAp3)
			}

			downAp3r := tt.input
			downAp3r.DownAp3r()
			if downAp3r != tt.expectedDownAp3r {
				t.Errorf("DownAp3r() = %v, want %v", downAp3r, tt.expectedDownAp3r)
			}
		})
	}
}

// TestCoordIJKHex2dConversion tests 2D hex coordinate conversion functions with C-derived expected results
func TestCoordIJKHex2dConversion(t *testing.T) {
	// Test IJK to Hex2d conversion
	ijkToHex2dTests := []struct {
		name string
		input CoordIJK
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
			if math.Abs(result.X - tt.expectedX) > 1e-10 {
				t.Errorf("ToHex2d().X = %f, want %f", result.X, tt.expectedX)
			}
			if math.Abs(result.Y - tt.expectedY) > 1e-10 {
				t.Errorf("ToHex2d().Y = %f, want %f", result.Y, tt.expectedY)
			}
		})
	}

	// Test Hex2d to IJK conversion
	hex2dToIJKTests := []struct {
		name string
		inputX, inputY float64
		expected CoordIJK
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
			v := Vec2d{X: tt.inputX, Y: tt.inputY}
			result := Hex2dToCoordIJK(v)
			if result != tt.expected {
				t.Errorf("Hex2dToCoordIJK(%v) = %v, want %v", v, result, tt.expected)
			}
		})
	}

	// Test coordinate equivalence through hex2d conversion
	// Note: Multiple IJK coordinates can represent the same hexagonal position.
	// The hex2d conversion acts as a normalizing function, so IJK->Hex2d->IJK 
	// may not preserve the original IJK, but it MUST preserve the hex2d position.
	equivalenceTests := []CoordIJK{
		{0, 0, 0},   // Should round-trip exactly (canonical form)
		{1, 0, 0},   // Should round-trip exactly (canonical form)  
		{0, 1, 0},   // Should round-trip exactly (canonical form)
		{0, 0, 1},   // Should round-trip exactly (canonical form)
		{1, 1, 0},   // Should round-trip exactly (canonical form)
		{-1, 0, 0},  // May normalize to equivalent form {0, 1, 1}
		{0, -1, 0},  // May normalize to equivalent form {1, 0, 1}
		{2, 1, -1},  // May normalize to equivalent form {3, 2, 0}
		{3, -2, 1},  // May normalize to equivalent form {5, 0, 3}
		{10, 5, -3}, // May normalize to equivalent form {13, 8, 0}
	}

	for i, original := range equivalenceTests {
		t.Run(fmt.Sprintf("hex2dEquivalence_%d", i), func(t *testing.T) {
			hex2d1 := original.ToHex2d()
			recovered := Hex2dToCoordIJK(hex2d1)
			hex2d2 := recovered.ToHex2d()
			
			// The key test: both coordinates must map to the same hex2d position
			if math.Abs(hex2d1.X-hex2d2.X) > 1e-10 || math.Abs(hex2d1.Y-hex2d2.Y) > 1e-10 {
				t.Errorf("Hex2d position not preserved: %v -> %v -> %v (hex2d: %v -> %v)", 
					original, hex2d1, recovered, hex2d1, hex2d2)
			}
			
			// Log when coordinates are normalized (this is expected, not an error)
			if recovered != original {
				t.Logf("Coordinate normalized: %v -> %v (same hex position)", original, recovered)
			}
		})
	}
}
