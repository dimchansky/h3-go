//go:build cgo

package c2go

import "testing"

func Test_h3index_maxFaceCount_ParityWithC(t *testing.T) {
	testCases := []H3Index{
		// Hexagons at various resolutions
		0x8928308280fffff, // res 9 hexagon
		0x821c07fffffffff, // res 2 hexagon
		0x8c2a1072b9dffff, // res 12 hexagon

		// Pentagons at various resolutions (base cell 4, 14, 24, 38, 49, etc.)
		0x8009fffffffffff, // res 0 pentagon (base cell 4)
		0x801dfffffffffff, // res 0 pentagon (base cell 14)
		0x8031fffffffffff, // res 0 pentagon (base cell 24)
		0x804dfffffffffff, // res 0 pentagon (base cell 38)
		0x8063fffffffffff, // res 0 pentagon (base cell 49)
		0x8077fffffffffff, // res 0 pentagon (base cell 58)
		0x807ffffffffffff, // res 0 pentagon (base cell 63)
		0x8091fffffffffff, // res 0 pentagon (base cell 72)
		0x80a7fffffffffff, // res 0 pentagon (base cell 83)
		0x80c3fffffffffff, // res 0 pentagon (base cell 97)
		0x80d7fffffffffff, // res 0 pentagon (base cell 107)
		0x80ebfffffffffff, // res 0 pentagon (base cell 117)

		// Higher resolution pentagons
		0x811d7ffffffffff, // res 1 pentagon
		0x821d7ffffffffff, // res 2 pentagon
		0x831d7ffffffffff, // res 3 pentagon
	}

	for _, h3 := range testCases {
		var goOut int32
		var cOut int32

		goErr := maxFaceCount(h3, &goOut)
		cErr := maxFaceCountC(h3, &cOut)

		if uint32(goErr) != cErr {
			t.Fatalf("maxFaceCount error mismatch for h3=%x: go=%d c=%d", uint64(h3), uint32(goErr), cErr)
		}

		if goOut != cOut {
			t.Fatalf("maxFaceCount result mismatch for h3=%x: go=%d c=%d", uint64(h3), goOut, cOut)
		}

		// Basic sanity checks
		if goErr == E_SUCCESS {
			if goOut != 2 && goOut != 5 {
				t.Fatalf("maxFaceCount returned unexpected value %d for h3=%x, expected 2 or 5", goOut, uint64(h3))
			}

			// Check consistency with isPentagon
			isPent := isPentagon(h3)
			expectedCount := int32(2)
			if isPent {
				expectedCount = 5
			}
			if goOut != expectedCount {
				t.Fatalf("maxFaceCount result %d inconsistent with isPentagon=%t for h3=%x", goOut, isPent, uint64(h3))
			}
		}
	}
}

func Test_h3index_maxFaceCount_InvalidIndex_ParityWithC(t *testing.T) {
	// Test with invalid H3 indices
	invalidIndices := []H3Index{
		0x0000000000000000, // H3_NULL
		0xffffffffffffffff, // Invalid bits
		0x1234567890abcdef, // Random invalid index
	}

	for _, h3 := range invalidIndices {
		var goOut int32
		var cOut int32

		goErr := maxFaceCount(h3, &goOut)
		cErr := maxFaceCountC(h3, &cOut)

		if uint32(goErr) != cErr {
			t.Fatalf("maxFaceCount error mismatch for invalid h3=%x: go=%d c=%d", uint64(h3), uint32(goErr), cErr)
		}

		if goOut != cOut {
			t.Fatalf("maxFaceCount result mismatch for invalid h3=%x: go=%d c=%d", uint64(h3), goOut, cOut)
		}
	}
}
