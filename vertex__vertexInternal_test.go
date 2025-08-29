package h3

import "testing"

// Tests ported from testVertexInternal.c

func TestVertexNumForDirection_hex(t *testing.T) {
	t.Parallel()
	origin := H3Index(0x823d6ffffffffff)
	vertexNums := make([]bool, NUM_HEX_VERTS)

	for dir := K_AXES_DIGIT; dir < NUM_DIGITS; dir++ {
		vertexNum := vertexNumForDirection(origin, dir)
		if vertexNum < 0 || vertexNum >= int32(NUM_HEX_VERTS) {
			t.Errorf("vertex number %d for direction %d is not valid (should be 0-%d)", vertexNum, dir, NUM_HEX_VERTS-1)
		}
		if vertexNums[vertexNum] {
			t.Errorf("vertex number %d for direction %d appears more than once", vertexNum, dir)
		}
		vertexNums[vertexNum] = true
	}
}

func TestVertexNumForDirection_pent(t *testing.T) {
	t.Parallel()
	pentagon := H3Index(0x823007fffffffff)
	vertexNums := make([]bool, NUM_PENT_VERTS)

	for dir := J_AXES_DIGIT; dir < NUM_DIGITS; dir++ {
		vertexNum := vertexNumForDirection(pentagon, dir)
		if vertexNum < 0 || vertexNum >= int32(NUM_PENT_VERTS) {
			t.Errorf("vertex number %d for direction %d is not valid (should be 0-%d)", vertexNum, dir, NUM_PENT_VERTS-1)
		}
		if vertexNums[vertexNum] {
			t.Errorf("vertex number %d for direction %d appears more than once", vertexNum, dir)
		}
		vertexNums[vertexNum] = true
	}
}

func TestVertexNumForDirection_badDirections(t *testing.T) {
	t.Parallel()
	origin := H3Index(0x823007fffffffff)

	// Test CENTER_DIGIT
	if vertexNumForDirection(origin, CENTER_DIGIT) != INVALID_VERTEX_NUM {
		t.Error("center digit should return invalid vertex")
	}

	// Test INVALID_DIGIT
	if vertexNumForDirection(origin, INVALID_DIGIT) != INVALID_VERTEX_NUM {
		t.Error("invalid digit should return invalid vertex")
	}

	// Test K direction on pentagon
	pentagon := H3Index(0x823007fffffffff)
	if vertexNumForDirection(pentagon, K_AXES_DIGIT) != INVALID_VERTEX_NUM {
		t.Error("K direction on pentagon should return invalid vertex")
	}
}

func TestDirectionForVertexNum_hex(t *testing.T) {
	t.Parallel()
	origin := H3Index(0x823d6ffffffffff)
	seenDirs := make([]bool, NUM_DIGITS)

	for vertexNum := int32(0); vertexNum < int32(NUM_HEX_VERTS); vertexNum++ {
		dir := directionForVertexNum(origin, vertexNum)
		if dir <= 0 || dir >= INVALID_DIGIT {
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
	origin := H3Index(0x823d6ffffffffff)

	// Test negative vertex
	if directionForVertexNum(origin, -1) != INVALID_DIGIT {
		t.Error("negative vertex should return invalid direction")
	}

	// Test invalid vertex (6 for hex)
	if directionForVertexNum(origin, 6) != INVALID_DIGIT {
		t.Error("invalid vertex should return invalid direction")
	}

	// Test invalid pent vertex (5 for pentagon)
	pentagon := H3Index(0x823007fffffffff)
	if directionForVertexNum(pentagon, 5) != INVALID_DIGIT {
		t.Error("invalid pent vertex should return invalid direction")
	}
}
