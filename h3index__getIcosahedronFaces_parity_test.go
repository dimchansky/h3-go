//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_getIcosahedronFaces_parity(t *testing.T) {
	tests := []struct {
		name string
		h3   h3Index
	}{
		{
			name: "regular hexagon res 0",
			h3:   0x8001fffffffffff, // res 0 hexagon
		},
		{
			name: "regular hexagon res 3",
			h3:   0x83080dfffffff, // res 3 hexagon
		},
		{
			name: "regular hexagon res 9",
			h3:   0x8930062838bffff, // res 9 hexagon
		},
		{
			name: "pentagon res 0",
			h3:   0x8007fffffffffff, // res 0 pentagon (base cell 4)
		},
		{
			name: "pentagon res 3",
			h3:   0x83026ffffffffff, // res 3 pentagon
		},
		{
			name: "pentagon res 9",
			h3:   0x89283080ddbffff, // res 9 pentagon
		},
		{
			name: "class II pentagon",
			h3:   0x8207fffffffffff, // res 2 pentagon (class II)
		},
		{
			name: "class III pentagon",
			h3:   0x8307fffffffffff, // res 3 pentagon (class III)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get expected face count
			var expectedFaceCount int32
			if err := maxFaceCountC(tt.h3, &expectedFaceCount); err != 0 {
				t.Fatalf("maxFaceCountC failed: %d", err)
			}

			var actualFaceCount int32
			if err := maxFaceCount(tt.h3, &actualFaceCount); err != eSuccess {
				t.Fatalf("maxFaceCount failed: %d", err)
			}

			if actualFaceCount != expectedFaceCount {
				t.Errorf("face count mismatch: Go=%d, C=%d", actualFaceCount, expectedFaceCount)
			}

			// Test getIcosahedronFaces
			goFaces := make([]int32, expectedFaceCount)
			cFaces := make([]int32, expectedFaceCount)

			// Call C implementation
			if err := getIcosahedronFacesC(tt.h3, cFaces); err != 0 {
				t.Fatalf("getIcosahedronFacesC failed: %d", err)
			}

			// Call Go implementation
			if err := getIcosahedronFaces(tt.h3, goFaces); err != eSuccess {
				t.Fatalf("getIcosahedronFaces failed: %d", err)
			}

			// Compare results - both arrays should have the same faces in potentially different order
			// Convert to sets for comparison
			goSet := make(map[int32]bool)
			cSet := make(map[int32]bool)

			for _, face := range goFaces {
				if face != invalidFace {
					goSet[face] = true
				}
			}

			for _, face := range cFaces {
				if face != invalidFace {
					cSet[face] = true
				}
			}

			// Compare sets
			if len(goSet) != len(cSet) {
				t.Errorf("face set size mismatch: Go has %d faces, C has %d faces", len(goSet), len(cSet))
			}

			for face := range goSet {
				if !cSet[face] {
					t.Errorf("Go implementation returned face %d which C implementation didn't", face)
				}
			}

			for face := range cSet {
				if !goSet[face] {
					t.Errorf("C implementation returned face %d which Go implementation didn't", face)
				}
			}

			t.Logf("h3Index %x: faces Go=%v, C=%v", tt.h3, goFaces, cFaces)
		})
	}
}
