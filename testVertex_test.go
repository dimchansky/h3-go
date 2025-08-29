// Tests ported from testVertex.c
package h3

import "testing"

func TestCellToVertex_badVerts(t *testing.T) {
	t.Parallel()
	
	origin := H3Index(0x823d6ffffffffff)

	var vert H3Index
	if err := cellToVertex(origin, -1, &vert); err != E_DOMAIN {
		t.Errorf("negative vertex should return E_DOMAIN, got %v", err)
	}
	if err := cellToVertex(origin, 6, &vert); err != E_DOMAIN {
		t.Errorf("invalid vertex should return E_DOMAIN, got %v", err)
	}

	pentagon := H3Index(0x823007fffffffff)
	if err := cellToVertex(pentagon, 5, &vert); err != E_DOMAIN {
		t.Errorf("invalid pent vertex should return E_DOMAIN, got %v", err)
	}
}

func TestCellToVertex_invalid(t *testing.T) {
	t.Parallel()
	
	invalid := H3Index(0xFFFFFFFFFFFFFFFF)
	var vert H3Index
	if err := cellToVertex(invalid, 3, &vert); err != E_FAILED {
		t.Errorf("Invalid cell should return E_FAILED, got %v", err)
	}
}

func TestCellToVertex_invalid2(t *testing.T) {
	t.Parallel()
	
	index := H3Index(0x685b2396e900fff9)
	var vert H3Index
	if err := cellToVertex(index, 2, &vert); err != E_CELL_INVALID {
		t.Errorf("Invalid cell should return E_CELL_INVALID, got %v", err)
	}
}

func TestCellToVertex_invalid3(t *testing.T) {
	t.Parallel()
	
	index := H3Index(0x20ff20202020ff35)
	var vert H3Index
	if err := cellToVertex(index, 0, &vert); err != E_CELL_INVALID {
		t.Errorf("Invalid cell should return E_CELL_INVALID, got %v", err)
	}
}

func TestIsValidVertex_hex(t *testing.T) {
	t.Parallel()
	
	origin := H3Index(0x823d6ffffffffff)
	vert := H3Index(0x2222597fffffffff)

	if !isValidVertex(vert) {
		t.Error("known vertex should be valid")
	}

	for i := int32(0); i < NUM_HEX_VERTS; i++ {
		if err := cellToVertex(origin, i, &vert); err != E_SUCCESS {
			t.Fatalf("cellToVertex should succeed, got %v", err)
		}
		if !isValidVertex(vert) {
			t.Errorf("vertex %d should be valid", i)
		}
	}
}

func TestIsValidVertex_invalidOwner(t *testing.T) {
	t.Parallel()
	
	origin := H3Index(0x823d6ffffffffff)
	vertexNum := int32(0)
	var vert H3Index
	if err := cellToVertex(origin, vertexNum, &vert); err != E_SUCCESS {
		t.Fatalf("cellToVertex should succeed, got %v", err)
	}

	// Set a bit for an unused digit to something else.
	vert ^= 1

	if isValidVertex(vert) {
		t.Error("vertex with invalid owner should not be valid")
	}
}

func TestIsValidVertex_wrongOwner(t *testing.T) {
	t.Parallel()
	
	origin := H3Index(0x823d6ffffffffff)
	vertexNum := int32(0)
	var vert H3Index
	if err := cellToVertex(origin, vertexNum, &vert); err != E_SUCCESS {
		t.Fatalf("cellToVertex should succeed, got %v", err)
	}

	// Assert that origin does not own the vertex
	owner := vert
	owner = setMode(owner, H3_CELL_MODE)
	owner = setReservedBits(owner, 0)

	if origin == owner {
		t.Error("origin should not own the canonical vertex")
	}

	nonCanonicalVertex := origin
	nonCanonicalVertex = setMode(nonCanonicalVertex, H3_VERTEX_MODE)
	nonCanonicalVertex = setReservedBits(nonCanonicalVertex, vertexNum)

	if isValidVertex(nonCanonicalVertex) {
		t.Error("vertex with incorrect owner should not be valid")
	}
}

func TestIsValidVertex_badVerts(t *testing.T) {
	t.Parallel()
	
	origin := H3Index(0x823d6ffffffffff)
	if isValidVertex(origin) {
		t.Error("cell should not be valid")
	}

	fakeEdge := origin
	fakeEdge = setMode(fakeEdge, H3_DIRECTEDEDGE_MODE)
	if isValidVertex(fakeEdge) {
		t.Error("edge mode should not be valid")
	}

	var vert H3Index
	if err := cellToVertex(origin, 0, &vert); err != E_SUCCESS {
		t.Fatalf("cellToVertex should succeed, got %v", err)
	}
	vert = setReservedBits(vert, 6)
	if isValidVertex(vert) {
		t.Error("invalid vertexNum should not be valid")
	}

	pentagon := H3Index(0x823007fffffffff)
	var vert2 H3Index
	if err := cellToVertex(pentagon, 0, &vert2); err != E_SUCCESS {
		t.Fatalf("cellToVertex should succeed, got %v", err)
	}
	vert2 = setReservedBits(vert2, 5)
	if isValidVertex(vert2) {
		t.Error("invalid pentagon vertexNum should not be valid")
	}
}

func TestVertexToLatLng_invalid(t *testing.T) {
	t.Parallel()
	
	invalid := H3Index(0xFFFFFFFFFFFFFFFF)
	var latLng LatLng
	if err := vertexToLatLng(invalid, &latLng); err == E_SUCCESS {
		t.Error("Invalid vertex should return error")
	}
}

func TestCellToVertexes_invalid(t *testing.T) {
	t.Parallel()
	
	invalid := H3Index(0xFFFFFFFFFFFFFFFF)
	var verts [6]H3Index
	if err := cellToVertexes(invalid, &verts); err != E_FAILED {
		t.Errorf("cellToVertexes should fail for invalid cell with E_FAILED, got %v", err)
	}
}