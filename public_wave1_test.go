package h3

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"
)

// Canonical vector from the H3 documentation: geoToH3(37.775938728915946,
// -122.41795063018799, 9) == 8928308280fffff.
const (
	sfLatDegs         = 37.775938728915946
	sfLngDegs         = -122.41795063018799
	sfCellRes9   Cell = 0x8928308280fffff
	sfCellRes15  Cell = 0x8f2830828052d25
	invalidIndex Cell = 0xffffffffffffffff
)

// mustCell returns a helper that unwraps (Cell, error) results, failing the
// test on error.
func mustCell(t *testing.T) func(Cell, error) Cell {
	return func(c Cell, err error) Cell {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
}

func TestLatLngToCellKnownVector(t *testing.T) {
	t.Parallel()

	must := mustCell(t)
	c := must(LatLngToCell(LatLngDegs(sfLatDegs, sfLngDegs), 9))
	if c != sfCellRes9 {
		t.Fatalf("LatLngToCell = %v, want %v", c, sfCellRes9)
	}
	if got := must(NewLatLng(Deg(sfLatDegs), Deg(sfLngDegs)).Cell(9)); got != c {
		t.Fatalf("LatLng.Cell = %v, want %v", got, c)
	}
}

func TestLatLngToCellDomainErrors(t *testing.T) {
	t.Parallel()

	if _, err := LatLngToCell(LatLngDegs(0, 0), -1); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("res -1: got %v, want ErrResolutionDomain", err)
	}
	if _, err := LatLngToCell(LatLngDegs(0, 0), MaxResolution+1); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("res 16: got %v, want ErrResolutionDomain", err)
	}
	// Enormous res values must not wrap through int32 narrowing.
	if _, err := LatLngToCell(LatLngDegs(0, 0), 1<<40+9); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("res 2^40+9: got %v, want ErrResolutionDomain", err)
	}
	if _, err := LatLngToCell(LatLng{Lat: Rad(math.NaN()), Lng: 0}, 9); !errors.Is(err, ErrLatLngDomain) {
		t.Errorf("NaN lat: got %v, want ErrLatLngDomain", err)
	}
}

func TestCellLatLngRoundTrip(t *testing.T) {
	t.Parallel()

	must := mustCell(t)
	for _, res := range []int{0, 5, 9, 15} {
		c := must(LatLngToCell(LatLngDegs(sfLatDegs, sfLngDegs), res))
		center, err := c.LatLng()
		if err != nil {
			t.Fatalf("res %d: LatLng: %v", res, err)
		}
		back := must(LatLngToCell(center, res))
		if back != c {
			t.Errorf("res %d: center round trip %v != %v", res, back, c)
		}
	}

	if _, err := invalidIndex.LatLng(); err == nil {
		t.Error("invalid index LatLng should fail")
	}
}

func TestCellBoundary(t *testing.T) {
	t.Parallel()

	b, err := sfCellRes9.Boundary()
	if err != nil {
		t.Fatal(err)
	}
	if b.Len() != 6 {
		t.Fatalf("hexagon boundary Len = %d, want 6", b.Len())
	}
	if len(b.Verts()) != 6 {
		t.Fatalf("Verts() len = %d, want 6", len(b.Verts()))
	}
	for i, v := range b.Verts() {
		if b.At(i) != v {
			t.Errorf("At(%d) != Verts()[%d]", i, i)
		}
		if math.Abs(v.Lat.Deg()) > 90 || math.Abs(v.Lng.Deg()) > 180 {
			t.Errorf("vertex %d out of range: %v", i, v)
		}
	}

	// Class II resolution (even): pentagon boundary has exactly 5 vertices.
	pentagonsII, err := Pentagons(8)
	if err != nil {
		t.Fatal(err)
	}
	pb, err := pentagonsII[0].Boundary()
	if err != nil {
		t.Fatal(err)
	}
	if pb.Len() != 5 {
		t.Errorf("Class II pentagon boundary Len = %d, want 5", pb.Len())
	}
	// Class III resolution (odd): icosahedron-edge distortion adds vertices,
	// up to MaxCellBoundaryVerts.
	pentagonsIII, err := Pentagons(9)
	if err != nil {
		t.Fatal(err)
	}
	pb3, err := pentagonsIII[0].Boundary()
	if err != nil {
		t.Fatal(err)
	}
	if pb3.Len() != MaxCellBoundaryVerts {
		t.Errorf("Class III pentagon boundary Len = %d, want %d", pb3.Len(), MaxCellBoundaryVerts)
	}

	defer func() {
		if recover() == nil {
			t.Error("At out of range should panic")
		}
	}()
	_ = b.At(6)
}

func TestInspection(t *testing.T) {
	t.Parallel()

	if got := sfCellRes9.Resolution(); got != 9 {
		t.Errorf("Resolution = %d, want 9", got)
	}
	if got := sfCellRes15.Resolution(); got != 15 {
		t.Errorf("Resolution = %d, want 15", got)
	}
	if got := sfCellRes9.BaseCellNumber(); got != 20 {
		t.Errorf("BaseCellNumber = %d, want 20", got)
	}
	if !sfCellRes9.IsValid() || Cell(0).IsValid() || invalidIndex.IsValid() {
		t.Error("IsValid misbehaves")
	}
	if sfCellRes9.IsPentagon() {
		t.Error("SF cell is not a pentagon")
	}
	if !sfCellRes9.IsResClassIII() {
		t.Error("res 9 is Class III")
	}
	c8 := mustCell(t)(sfCellRes9.Parent(8))
	if c8.IsResClassIII() {
		t.Error("res 8 is Class II")
	}
}

func TestIcosahedronFaces(t *testing.T) {
	t.Parallel()

	faces, err := sfCellRes9.IcosahedronFaces()
	if err != nil {
		t.Fatal(err)
	}
	if len(faces) != 1 {
		t.Errorf("hexagon far from edges should span 1 face, got %v", faces)
	}
	pentagons, err := Pentagons(0)
	if err != nil {
		t.Fatal(err)
	}
	pf, err := pentagons[0].IcosahedronFaces()
	if err != nil {
		t.Fatal(err)
	}
	if len(pf) != 5 {
		t.Errorf("res-0 pentagon should span 5 faces, got %v", pf)
	}
	for _, f := range append(faces, pf...) {
		if f < 0 || f >= numIcosaFaces {
			t.Errorf("face %d out of range", f)
		}
	}
}

func TestPentagonsAndRes0Cells(t *testing.T) {
	t.Parallel()

	for _, res := range []int{0, 7, 15} {
		ps, err := Pentagons(res)
		if err != nil {
			t.Fatal(err)
		}
		if len(ps) != NumPentagons {
			t.Fatalf("res %d: %d pentagons, want %d", res, len(ps), NumPentagons)
		}
		for _, p := range ps {
			if !p.IsPentagon() || p.Resolution() != res {
				t.Errorf("res %d: bad pentagon %v", res, p)
			}
		}
	}
	if _, err := Pentagons(-1); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("Pentagons(-1): got %v", err)
	}

	r0 := Res0Cells()
	if len(r0) != NumRes0Cells {
		t.Fatalf("Res0Cells len = %d, want %d", len(r0), NumRes0Cells)
	}
	for i, c := range r0 {
		if !c.IsValid() || c.Resolution() != 0 || c.BaseCellNumber() != i {
			t.Errorf("res0 cell %d invalid: %v", i, c)
		}
	}
}

func TestNumCells(t *testing.T) {
	t.Parallel()

	n0, err := NumCells(0)
	if err != nil || n0 != 122 {
		t.Errorf("NumCells(0) = %d, %v; want 122", n0, err)
	}
	n15, err := NumCells(15)
	if err != nil || n15 != 569707381193162 {
		t.Errorf("NumCells(15) = %d, %v; want 569707381193162", n15, err)
	}
	if _, err := NumCells(16); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("NumCells(16): got %v", err)
	}
}

func TestHierarchy(t *testing.T) {
	t.Parallel()

	parent, err := sfCellRes9.Parent(8)
	if err != nil {
		t.Fatal(err)
	}
	if parent.Resolution() != 8 {
		t.Fatalf("parent res = %d", parent.Resolution())
	}

	children, err := parent.Children(9)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) != 7 {
		t.Fatalf("children = %d, want 7", len(children))
	}
	n, err := parent.NumChildren(9)
	if err != nil || n != int64(len(children)) {
		t.Fatalf("NumChildren = %d, %v", n, err)
	}
	found := false
	for _, ch := range children {
		if ch == sfCellRes9 {
			found = true
		}
		p, err := ch.Parent(8)
		if err != nil || p != parent {
			t.Errorf("child %v parent mismatch", ch)
		}
	}
	if !found {
		t.Error("original cell not among parent's children")
	}

	cc, err := parent.CenterChild(9)
	if err != nil {
		t.Fatal(err)
	}
	if cc != children[0] {
		t.Errorf("center child %v != first child %v", cc, children[0])
	}

	// Pentagons have 6 children (1 pentagon + 5 hexagons).
	pentagons, err := Pentagons(4)
	if err != nil {
		t.Fatal(err)
	}
	pn, err := pentagons[0].NumChildren(5)
	if err != nil || pn != 6 {
		t.Errorf("pentagon NumChildren = %d, %v; want 6", pn, err)
	}

	// ChildPos round trip.
	pos, err := sfCellRes15.ChildPos(9)
	if err != nil {
		t.Fatal(err)
	}
	anc, err := sfCellRes15.Parent(9)
	if err != nil {
		t.Fatal(err)
	}
	back, err := anc.ChildAtPos(pos, 15)
	if err != nil {
		t.Fatal(err)
	}
	if back != sfCellRes15 {
		t.Errorf("ChildAtPos(ChildPos) = %v, want %v", back, sfCellRes15)
	}

	// Coarser-than-cell child resolution must error.
	if _, err := sfCellRes9.Children(3); err == nil {
		t.Error("Children at coarser res should fail")
	}
}

func TestAppendChildrenSemantics(t *testing.T) {
	t.Parallel()

	parent, _ := sfCellRes9.Parent(5)
	prefix := []Cell{42, 43}
	out, err := parent.AppendChildren(prefix, 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2+49 {
		t.Fatalf("append result len = %d, want 51", len(out))
	}
	if out[0] != 42 || out[1] != 43 {
		t.Error("prefix clobbered")
	}
	for _, c := range out[2:] {
		if !c.IsValid() || c.Resolution() != 7 {
			t.Errorf("bad child %v", c)
		}
	}
}

func TestParseFormat(t *testing.T) {
	t.Parallel()

	s := sfCellRes9.String()
	if s != "8928308280fffff" {
		t.Fatalf("String = %q", s)
	}
	for _, in := range []string{s, "0x" + s, "0X" + s} {
		c, err := ParseCell(in)
		if err != nil || c != sfCellRes9 {
			t.Errorf("ParseCell(%q) = %v, %v", in, c, err)
		}
	}
	if _, err := ParseCell("not-hex"); err == nil {
		t.Error("garbage should not parse")
	}
	if _, err := ParseCell(""); err == nil {
		t.Error("empty string should not parse")
	}
	// A valid *edge* index must be rejected as a cell (mode validation).
	var edges [6]h3Index
	if errC := originToDirectedEdges(sfCellRes9, edges[:]); errC != eSuccess {
		t.Fatalf("originToDirectedEdges: %v", errC)
	}
	edgeStr := DirectedEdge(edges[1]).String()
	if _, err := ParseCell(edgeStr); !errors.Is(err, ErrCellInvalid) {
		t.Errorf("ParseCell(edge) = %v, want ErrCellInvalid", err)
	}
	if e, err := ParseDirectedEdge(edgeStr); err != nil || e != DirectedEdge(edges[1]) {
		t.Errorf("ParseDirectedEdge(%q) = %v, %v", edgeStr, e, err)
	}
	if _, err := ParseDirectedEdge(s); !errors.Is(err, ErrDirectedEdgeInvalid) {
		t.Errorf("ParseDirectedEdge(cell) = %v, want ErrDirectedEdgeInvalid", err)
	}
}

func TestTextMarshaling(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(sfCellRes9)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"8928308280fffff"` {
		t.Fatalf("json.Marshal = %s", data)
	}
	var c Cell
	if err := json.Unmarshal(data, &c); err != nil || c != sfCellRes9 {
		t.Fatalf("json.Unmarshal = %v, %v", c, err)
	}
	if err := json.Unmarshal([]byte(`"ffffffffffffffff"`), &c); err == nil {
		t.Error("invalid cell must not unmarshal")
	}
}

func TestWave1Allocations(t *testing.T) {
	ll := LatLngDegs(sfLatDegs, sfLngDegs)

	assertAllocs := func(name string, want float64, f func()) {
		t.Helper()
		if got := testing.AllocsPerRun(200, f); got > want {
			t.Errorf("%s allocates %v/run, want <= %v", name, got, want)
		}
	}

	assertAllocs("LatLngToCell", 0, func() { _, _ = LatLngToCell(ll, 9) })
	assertAllocs("Cell.LatLng", 0, func() { _, _ = sfCellRes9.LatLng() })
	assertAllocs("Cell.Boundary", 0, func() { _, _ = sfCellRes9.Boundary() })
	assertAllocs("Cell.Resolution", 0, func() { _ = sfCellRes9.Resolution() })
	assertAllocs("Cell.IsValid", 0, func() { _ = sfCellRes9.IsValid() })
	assertAllocs("Cell.String", 1, func() { _ = sfCellRes9.String() })

	parent, _ := sfCellRes9.Parent(5)
	buf := make([]Cell, 0, 2500)
	assertAllocs("AppendChildren warm", 0, func() {
		out, err := parent.AppendChildren(buf, 7)
		if err != nil || len(out) != 49 {
			t.Fatal(err, len(out))
		}
	})

	b, _ := sfCellRes9.Boundary()
	assertAllocs("CellBoundary.Verts", 0, func() { _ = b.Verts() })
}

func TestBoundaryMatchesInternal(t *testing.T) {
	t.Parallel()

	// The public path must agree exactly with the parity-tested internal path.
	cells := []Cell{sfCellRes9, sfCellRes15}
	pentagons, err := Pentagons(3)
	if err != nil {
		t.Fatal(err)
	}
	cells = append(cells, pentagons...)
	for _, c := range cells {
		var want CellBoundary
		if errC := cellToBoundary(c, &want); errC != eSuccess {
			t.Fatalf("internal boundary failed for %v", c)
		}
		got, err := c.Boundary()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("boundary mismatch for %v", c)
		}
	}
}

func TestStringZeroAndInvalid(t *testing.T) {
	t.Parallel()

	if got := Cell(0).String(); got != "0" {
		t.Errorf("Cell(0).String() = %q, want \"0\" (C h3ToString behavior)", got)
	}
	if !strings.HasPrefix(invalidIndex.String(), "ffffffffffffffff") {
		t.Errorf("invalid String = %q", invalidIndex.String())
	}
}
