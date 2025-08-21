package coordijk

import "testing"

// TestCoordIJKBasicOperations tests Add, Sub, Scale operations with C-derived expected results.
func TestCoordIJKBasicOperations(t *testing.T) {
	tests := []struct {
		name          string
		a, b          CoordIJK
		expectedAdd   CoordIJK
		expectedSub   CoordIJK
		scaleFactor   int
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
			aCopy := tt.a
			resultAdd := aCopy.Add(tt.b)
			if *resultAdd != tt.expectedAdd {
				t.Errorf("Add() = %v, want %v", *resultAdd, tt.expectedAdd)
			}

			aCopy = tt.a
			resultSub := aCopy.Sub(tt.b)
			if *resultSub != tt.expectedSub {
				t.Errorf("Sub() = %v, want %v", *resultSub, tt.expectedSub)
			}

			aCopy = tt.a
			resultScale := aCopy.Scale(tt.scaleFactor)
			if *resultScale != tt.expectedScale {
				t.Errorf("Scale() = %v, want %v", *resultScale, tt.expectedScale)
			}
		})
	}
}
