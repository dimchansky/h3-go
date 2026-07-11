// Tests ported from testGetIcosahedronFaces.c
package h3

import (
	"testing"
)

func countFaces(t *testing.T, h3 h3Index, expectedMax int32) int32 {
	t.Helper()
	var sz int32
	err := maxFaceCount(h3, &sz)
	if err != eSuccess {
		t.Fatalf("maxFaceCount failed: %v", err)
	}
	if sz != expectedMax {
		t.Fatalf("expected max face count %d, got %d", expectedMax, sz)
	}

	faces := make([]int32, sz)
	err = getIcosahedronFaces(h3, faces)
	if err != eSuccess {
		t.Fatalf("getIcosahedronFaces failed: %v", err)
	}

	var validCount int32
	for i := int32(0); i < sz; i++ {
		if faces[i] >= 0 && faces[i] <= 19 {
			validCount++
		}
	}

	return validCount
}

func assertSingleHexFace(t *testing.T, h3 h3Index) {
	t.Helper()
	validCount := countFaces(t, h3, 2)
	if validCount != 1 {
		t.Errorf("expected single valid face, got %d", validCount)
	}
}

func assertMultipleHexFaces(t *testing.T, h3 h3Index) {
	t.Helper()
	validCount := countFaces(t, h3, 2)
	if validCount != 2 {
		t.Errorf("expected multiple valid faces for hexagon, got %d", validCount)
	}
}

func assertPentagonFaces(t *testing.T, h3 h3Index) {
	t.Helper()
	if !isPentagon(h3) {
		t.Fatal("expected a pentagon")
	}
	validCount := countFaces(t, h3, 5)
	if validCount != 5 {
		t.Errorf("expected 5 valid faces for pentagon, got %d", validCount)
	}
}

func TestGetIcosahedronFaces_SingleFaceHexes(t *testing.T) {
	t.Parallel()

	// base cell 16 is at the center of an icosahedron face,
	// so all children should have the same face
	iterateBaseCellIndexesAtRes(t, 2, assertSingleHexFace, 16)
	iterateBaseCellIndexesAtRes(t, 3, assertSingleHexFace, 16)
}

func TestGetIcosahedronFaces_HexagonWithEdgeVertices(t *testing.T) {
	t.Parallel()

	// Class II pentagon neighbor - one face, two adjacent vertices on edge
	h3 := h3Index(0x821c37fffffffff)
	assertSingleHexFace(t, h3)
}

func TestGetIcosahedronFaces_HexagonWithDistortion(t *testing.T) {
	t.Parallel()

	// Class III pentagon neighbor, distortion across faces
	h3 := h3Index(0x831c06fffffffff)
	assertMultipleHexFaces(t, h3)
}

func TestGetIcosahedronFaces_HexagonCrossingFaces(t *testing.T) {
	t.Parallel()

	// Class II hex with two vertices on edge
	h3 := h3Index(0x821ce7fffffffff)
	assertMultipleHexFaces(t, h3)
}

func TestGetIcosahedronFaces_ClassIIIPentagon(t *testing.T) {
	t.Parallel()

	var pentagon h3Index
	setH3Index(&pentagon, 1, 4, 0)
	assertPentagonFaces(t, pentagon)
}

func TestGetIcosahedronFaces_ClassIIPentagon(t *testing.T) {
	t.Parallel()

	var pentagon h3Index
	setH3Index(&pentagon, 2, 4, 0)
	assertPentagonFaces(t, pentagon)
}

func TestGetIcosahedronFaces_Res15Pentagon(t *testing.T) {
	t.Parallel()

	var pentagon h3Index
	setH3Index(&pentagon, 15, 4, 0)
	assertPentagonFaces(t, pentagon)
}

func TestGetIcosahedronFaces_BaseCellHexagons(t *testing.T) {
	t.Parallel()

	var singleCount int32
	var multipleCount int32
	for i := int32(0); i < numBaseCells; i++ {
		if !_isBaseCellPentagon(i) {
			// Make the base cell index
			var baseCell h3Index
			setH3Index(&baseCell, 0, i, 0)
			validCount := countFaces(t, baseCell, 2)
			if validCount < 1 {
				t.Errorf("expected at least one face for base cell %d", i)
			}
			if validCount == 1 {
				singleCount++
			} else {
				multipleCount++
			}
		}
	}
	if singleCount != 4*20 {
		t.Errorf("expected single face for 4 aligned hex base cells per face (80), got %d", singleCount)
	}
	if multipleCount != int32(1.5*20) {
		t.Errorf("expected multiple faces for non-aligned hex base cells (30), got %d", multipleCount)
	}
}

func TestGetIcosahedronFaces_BaseCellPentagons(t *testing.T) {
	t.Parallel()

	for i := int32(0); i < numBaseCells; i++ {
		if _isBaseCellPentagon(i) {
			// Make the base cell index
			var baseCell h3Index
			setH3Index(&baseCell, 0, i, 0)
			assertPentagonFaces(t, baseCell)
		}
	}
}

func TestGetIcosahedronFaces_Invalid(t *testing.T) {
	t.Parallel()

	invalid := h3Index(0xFFFFFFFFFFFFFFFF)
	out := make([]int32, 1)
	err := getIcosahedronFaces(invalid, out)
	if err != eCellInvalid {
		t.Errorf("expected eCellInvalid for invalid cell, got %v", err)
	}
}

func TestGetIcosahedronFaces_Invalid2(t *testing.T) {
	t.Parallel()

	invalid := h3Index(0x71330073003f004e)
	var sz int32
	err := maxFaceCount(invalid, &sz)
	if err != eSuccess {
		t.Fatalf("maxFaceCount failed: %v", err)
	}
	faces := make([]int32, sz)
	err = getIcosahedronFaces(invalid, faces)
	if err != eFailed {
		t.Errorf("expected eFailed for invalid cell, got %v", err)
	}
}
