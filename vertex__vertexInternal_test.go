package h3

import "testing"

// Tests ported from testVertexInternal.c

func TestVertexNumForDirection_hex(t *testing.T) {
	t.Parallel()
	origin := h3Index(0x823d6ffffffffff)
	vertexNums := make([]bool, numHexVerts)

	for dir := kAxesDigit; dir < numDigits; dir++ {
		vertexNum := vertexNumForDirection(origin, dir)
		if vertexNum < 0 || vertexNum >= int32(numHexVerts) {
			t.Errorf("vertex number %d for direction %d is not valid (should be 0-%d)", vertexNum, dir, numHexVerts-1)
		}
		if vertexNums[vertexNum] {
			t.Errorf("vertex number %d for direction %d appears more than once", vertexNum, dir)
		}
		vertexNums[vertexNum] = true
	}
}

func TestVertexNumForDirection_pent(t *testing.T) {
	t.Parallel()
	pentagon := h3Index(0x823007fffffffff)
	vertexNums := make([]bool, numPentVerts)

	for dir := jAxesDigit; dir < numDigits; dir++ {
		vertexNum := vertexNumForDirection(pentagon, dir)
		if vertexNum < 0 || vertexNum >= int32(numPentVerts) {
			t.Errorf("vertex number %d for direction %d is not valid (should be 0-%d)", vertexNum, dir, numPentVerts-1)
		}
		if vertexNums[vertexNum] {
			t.Errorf("vertex number %d for direction %d appears more than once", vertexNum, dir)
		}
		vertexNums[vertexNum] = true
	}
}

func TestVertexNumForDirection_badDirections(t *testing.T) {
	t.Parallel()
	origin := h3Index(0x823007fffffffff)

	// Test centerDigit
	if vertexNumForDirection(origin, centerDigit) != invalidVertexNum {
		t.Error("center digit should return invalid vertex")
	}

	// Test invalidDigit
	if vertexNumForDirection(origin, invalidDigit) != invalidVertexNum {
		t.Error("invalid digit should return invalid vertex")
	}

	// Test K direction on pentagon
	pentagon := h3Index(0x823007fffffffff)
	if vertexNumForDirection(pentagon, kAxesDigit) != invalidVertexNum {
		t.Error("K direction on pentagon should return invalid vertex")
	}
}

func TestDirectionForVertexNum_hex(t *testing.T) {
	t.Parallel()
	origin := h3Index(0x823d6ffffffffff)
	seenDirs := make([]bool, numDigits)

	for vertexNum := int32(0); vertexNum < int32(numHexVerts); vertexNum++ {
		dir := directionForVertexNum(origin, vertexNum)
		if dir <= 0 || dir >= invalidDigit {
			t.Errorf("direction %d for vertex %d is not valid", dir, vertexNum)
		}
		if seenDirs[dir] {
			t.Errorf("direction %d for vertex %d appears more than once", dir, vertexNum)
		}
		seenDirs[dir] = true
	}
}

func TestDirectionForVertexNum_badVerts(t *testing.T) {
	t.Parallel()
	origin := h3Index(0x823d6ffffffffff)

	// Test negative vertex
	if directionForVertexNum(origin, -1) != invalidDigit {
		t.Error("negative vertex should return invalid direction")
	}

	// Test invalid vertex (6 for hex)
	if directionForVertexNum(origin, 6) != invalidDigit {
		t.Error("invalid vertex should return invalid direction")
	}

	// Test invalid pent vertex (5 for pentagon)
	pentagon := h3Index(0x823007fffffffff)
	if directionForVertexNum(pentagon, 5) != invalidDigit {
		t.Error("invalid pent vertex should return invalid direction")
	}
}
