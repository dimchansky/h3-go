package h3

import "slices"

// Vertex is an H3 vertex index: a single topological vertex of an H3 cell,
// shared by three cells. One of the three neighboring cells is arbitrarily
// designated the vertex's "owner" and determines its canonical index — so
// the same topological corner yields the identical Vertex value regardless
// of which of its three adjacent cells it is derived from.
//
// The zero value is not a valid vertex; IsValid reports false for it. A
// Vertex may be constructed by conversion from a raw uint64 index (unchecked
// — use IsValid to verify), parsed with ParseVertex, or produced by
// operations such as Cell.Vertex.
type Vertex uint64

// Vertex returns the cell's topological vertex with the given number
// (0..5 for hexagons, 0..4 for pentagons).
//
// H3 C API: cellToVertex.
func (c Cell) Vertex(vertexNum int) (Vertex, error) {
	if vertexNum < 0 || vertexNum >= numHexVerts {
		return 0, ErrDomain
	}
	var out h3Index
	if errC := cellToVertex(c, int32(vertexNum), &out); errC != eSuccess {
		return 0, toErr(errC)
	}
	return Vertex(out), nil
}

// Vertexes returns all topological vertexes of the cell: 6 for hexagons, 5
// for pentagons. Unlike Cell.Boundary, which returns geographic vertices
// and may additionally include distortion vertices, the result covers the
// cell's topological corners only — distortion vertices have no Vertex
// index.
//
// H3 C API: cellToVertexes.
func (c Cell) Vertexes() ([]Vertex, error) { return c.AppendVertexes(nil) }

// AppendVertexes appends all topological vertexes of the cell to dst and
// returns the extended slice: 6 vertexes for a hexagon, 5 for a pentagon
// (see Vertexes). Pass dst[:0] (or nil) to reuse dst's capacity; a capacity
// of 6 is always sufficient and makes the call allocation-free. On error
// the returned slice has dst's original length and elements.
//
// H3 C API: cellToVertexes.
func (c Cell) AppendVertexes(dst []Vertex) ([]Vertex, error) {
	var raw [6]h3Index
	if errC := cellToVertexes(c, &raw); errC != eSuccess {
		return dst, toErr(errC)
	}
	n := numHexVerts
	if raw[numHexVerts-1] == h3Null { // pentagons leave the last slot empty
		n = numPentVerts
	}
	start := len(dst)
	dst = slices.Grow(dst, n)[:start+n]
	for i := range n {
		dst[start+i] = Vertex(raw[i])
	}
	return dst, nil
}

// IsValid reports whether the index is a valid H3 vertex index.
//
// H3 C API: isValidVertex.
func (v Vertex) IsValid() bool { return isValidVertex(h3Index(v)) }

// Resolution returns the vertex's resolution. It is a pure bit accessor and
// does not validate the index.
//
// H3 C API: getResolution.
func (v Vertex) Resolution() int { return int(getResolution(h3Index(v))) }

// IndexDigit returns the indexing digit of the vertex's owner cell at the
// given resolution (1..MaxResolution; resolution 0 is the base cell number,
// not a digit). res may exceed the vertex's actual resolution, in which case
// the stored digit (7 for valid vertexes) is returned.
//
// H3 C API: getIndexDigit.
func (v Vertex) IndexDigit(res int) (int, error) {
	return indexDigit(v, res)
}

// LatLng returns the geographic coordinates of the vertex.
//
// H3 C API: vertexToLatLng.
func (v Vertex) LatLng() (LatLng, error) {
	var g LatLng
	if errC := vertexToLatLng(h3Index(v), &g); errC != eSuccess {
		return LatLng{}, toErr(errC)
	}
	return g, nil
}
