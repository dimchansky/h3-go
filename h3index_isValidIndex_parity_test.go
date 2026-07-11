//go:build cgo && c2go

package h3

import "testing"

func Test_isValidIndex_parity(t *testing.T) {
	var edges [6]h3Index
	cell := h3Index(0x8928308280fffff)
	if err := originToDirectedEdges(cell, edges[:]); err != eSuccess {
		t.Fatal(err)
	}
	var vertexes [6]h3Index
	if err := cellToVertexes(cell, &vertexes); err != eSuccess {
		t.Fatal(err)
	}
	inputs := []h3Index{
		0, ^h3Index(0), cell, edges[0], vertexes[0],
		0x8009fffffffffff, // pentagon
		cell | (1 << 63),  // high bit set
		cell ^ (7 << 59),  // mangled mode
	}
	for _, h := range inputs {
		if got, want := isValidIndex(h), isValidIndexC(h); got != want {
			t.Fatalf("isValidIndex(%#x): Go %v != C %v", uint64(h), got, want)
		}
	}
}
