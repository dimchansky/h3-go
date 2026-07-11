package h3

import (
	"errors"
	"math"
	"slices"
	"testing"
)

// sfTestPolygon is the San Francisco test polygon used throughout the
// upstream H3 test suite (testPolygonToCells.c).
func sfTestPolygon() GeoPolygon {
	return GeoPolygon{GeoLoop: GeoLoop{
		{Lat: Rad(0.659966917655), Lng: Rad(-2.1364398519396)},
		{Lat: Rad(0.6595011102219), Lng: Rad(-2.1359434279405)},
		{Lat: Rad(0.6583348114025), Lng: Rad(-2.1354884206045)},
		{Lat: Rad(0.6581220034068), Lng: Rad(-2.1382437718946)},
		{Lat: Rad(0.6594479998527), Lng: Rad(-2.1384597563896)},
		{Lat: Rad(0.6599990002976), Lng: Rad(-2.1376771158464)},
	}}
}

func TestPublicPolygonToCells(t *testing.T) {
	t.Parallel()

	cells, err := PolygonToCells(sfTestPolygon(), 9)
	if err != nil {
		t.Fatal(err)
	}
	// Canonical count from upstream testPolygonToCells.c.
	if len(cells) != 1253 {
		t.Fatalf("PolygonToCells = %d cells, want 1253", len(cells))
	}
	seen := map[Cell]bool{}
	for _, c := range cells {
		if !c.IsValid() || c.Resolution() != 9 {
			t.Fatalf("bad cell %v", c)
		}
		if seen[c] {
			t.Fatalf("duplicate cell %v", c)
		}
		seen[c] = true
	}
	if !seen[sfCellRes9] {
		t.Error("expected SF cell in polygon fill")
	}

	if _, err := PolygonToCells(sfTestPolygon(), -1); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("res -1: got %v", err)
	}

	// Append form reuses capacity.
	buf := make([]Cell, 0, 2048)
	out, err := AppendPolygonToCells(buf, sfTestPolygon(), 9)
	if err != nil || len(out) != 1253 {
		t.Fatalf("Append form: %d cells (%v)", len(out), err)
	}
}

func TestPolygonToCellsExperimental(t *testing.T) {
	t.Parallel()

	center, err := PolygonToCellsExperimental(sfTestPolygon(), 9, ContainmentCenter)
	if err != nil {
		t.Fatal(err)
	}
	if len(center) != 1253 {
		t.Fatalf("experimental center = %d cells, want 1253", len(center))
	}

	full, err := PolygonToCellsExperimental(sfTestPolygon(), 9, ContainmentFull)
	if err != nil {
		t.Fatal(err)
	}
	overlapping, err := PolygonToCellsExperimental(sfTestPolygon(), 9, ContainmentOverlapping)
	if err != nil {
		t.Fatal(err)
	}
	if len(full) >= len(center) || len(center) >= len(overlapping) {
		t.Errorf("containment ordering violated: full %d, center %d, overlapping %d",
			len(full), len(center), len(overlapping))
	}

	if _, err := PolygonToCellsExperimental(sfTestPolygon(), 9, ContainmentInvalid); !errors.Is(err, ErrOptionInvalid) {
		t.Errorf("invalid mode: got %v, want ErrOptionInvalid", err)
	}
}

func TestPolygonToCellsExperimentalSeq(t *testing.T) {
	t.Parallel()

	seq, err := PolygonToCellsExperimentalSeq(sfTestPolygon(), 9, ContainmentCenter)
	if err != nil {
		t.Fatal(err)
	}
	var got []Cell
	for c := range seq {
		got = append(got, c)
	}
	want, err := PolygonToCellsExperimental(sfTestPolygon(), 9, ContainmentCenter)
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("Seq result (%d) differs from slice result (%d)", len(got), len(want))
	}

	// Early break must not panic or leak.
	n := 0
	for range seq {
		n++
		if n == 10 {
			break
		}
	}
	if n != 10 {
		t.Errorf("early break iterated %d", n)
	}

	if _, err := PolygonToCellsExperimentalSeq(sfTestPolygon(), 99, ContainmentCenter); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("bad res: got %v", err)
	}
	if _, err := PolygonToCellsExperimentalSeq(sfTestPolygon(), 9, ContainmentInvalid); !errors.Is(err, ErrOptionInvalid) {
		t.Errorf("bad mode: got %v", err)
	}
}

func TestChildrenSeqAndCellsAtRes(t *testing.T) {
	t.Parallel()

	parent, _ := sfCellRes9.Parent(5)
	var fromSeq []Cell
	for c := range parent.ChildrenSeq(8) {
		fromSeq = append(fromSeq, c)
	}
	want, err := parent.Children(8)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(fromSeq, want) {
		t.Errorf("ChildrenSeq (%d) differs from Children (%d)", len(fromSeq), len(want))
	}

	// Empty on invalid input.
	for range parent.ChildrenSeq(-1) {
		t.Fatal("invalid res must yield empty sequence")
	}
	for range parent.ChildrenSeq(3) {
		t.Fatal("coarser res must yield empty sequence")
	}

	// CellsAtRes(0) must equal Res0Cells.
	var res0 []Cell
	for c := range CellsAtRes(0) {
		res0 = append(res0, c)
	}
	if !slices.Equal(res0, Res0Cells()) {
		t.Error("CellsAtRes(0) != Res0Cells()")
	}

	// Count cells at res 1: 122*7 - 12*1 = 842.
	n := 0
	for range CellsAtRes(1) {
		n++
	}
	if want, err := NumCells(1); err != nil || int64(n) != want {
		t.Errorf("CellsAtRes(1) count = %d, want %d (%v)", n, want, err)
	}

	for range CellsAtRes(-1) {
		t.Fatal("invalid res must yield empty sequence")
	}
}

func TestIterAllocations(t *testing.T) {
	parent, _ := sfCellRes9.Parent(5)
	allocs := testing.AllocsPerRun(50, func() {
		n := 0
		for range parent.ChildrenSeq(8) {
			n++
		}
		if n != 343 {
			t.Fatal(n)
		}
	})
	if allocs != 0 {
		t.Errorf("ChildrenSeq iteration allocates %v/run, want 0", allocs)
	}
}

func TestCompactUncompactRoundTrip(t *testing.T) {
	t.Parallel()

	parent, _ := sfCellRes9.Parent(4)
	cells, err := parent.Children(6)
	if err != nil {
		t.Fatal(err)
	}

	compacted, err := CompactCells(cells)
	if err != nil {
		t.Fatal(err)
	}
	if len(compacted) != 1 || compacted[0] != parent {
		t.Fatalf("compacting all children should yield the parent, got %v", compacted)
	}

	expanded, err := UncompactCells(compacted, 6)
	if err != nil {
		t.Fatal(err)
	}
	a, b := slices.Clone(cells), slices.Clone(expanded)
	slices.Sort(a)
	slices.Sort(b)
	if !slices.Equal(a, b) {
		t.Error("uncompact(compact(children)) != children")
	}

	// Partial set stays partially compacted.
	partial := cells[:len(cells)-1]
	pc, err := CompactCells(partial)
	if err != nil {
		t.Fatal(err)
	}
	if len(pc) >= len(partial) || len(pc) <= 1 {
		t.Errorf("partial compact = %d cells (from %d)", len(pc), len(partial))
	}

	// Duplicate detection matches C: a full child set plus one duplicate
	// overflows the parent bucket (upstream compactCells_duplicateMinimum).
	dupParent, _ := sfCellRes9.Parent(8)
	dupChildren, err := dupParent.Children(9)
	if err != nil {
		t.Fatal(err)
	}
	dupChildren = append(dupChildren, dupChildren[0])
	if _, err := CompactCells(dupChildren); !errors.Is(err, ErrDuplicateInput) {
		t.Errorf("duplicate compact: got %v, want ErrDuplicateInput", err)
	}
	// Null entries in uncompact input are skipped (C contract).
	un, err := UncompactCells([]Cell{0, parent, 0}, 5)
	if err != nil {
		t.Fatal(err)
	}
	if int64(len(un)) != 7 {
		t.Errorf("uncompact with nulls = %d cells, want 7", len(un))
	}
	// UncompactCellsSize agrees.
	sz, err := UncompactCellsSize([]Cell{0, parent, 0}, 5)
	if err != nil || sz != 7 {
		t.Errorf("UncompactCellsSize = %d (%v), want 7", sz, err)
	}
}

func TestCellsToMultiPolygon(t *testing.T) {
	t.Parallel()

	// Empty input.
	polys, err := CellsToMultiPolygon(nil)
	if err != nil || polys != nil {
		t.Fatalf("empty input: %v, %v", polys, err)
	}

	// A single cell produces one polygon with one 6-vertex loop, no holes.
	polys, err = CellsToMultiPolygon([]Cell{sfCellRes9})
	if err != nil {
		t.Fatal(err)
	}
	if len(polys) != 1 || len(polys[0].GeoLoop) != 6 || len(polys[0].Holes) != 0 {
		t.Fatalf("single cell: %d polys, %d verts, %d holes",
			len(polys), len(polys[0].GeoLoop), len(polys[0].Holes))
	}

	// A 1-ring (hollow: origin excluded) produces one polygon with a hole.
	ring, err := sfCellRes9.GridRing(1)
	if err != nil {
		t.Fatal(err)
	}
	polys, err = CellsToMultiPolygon(ring)
	if err != nil {
		t.Fatal(err)
	}
	if len(polys) != 1 || len(polys[0].Holes) != 1 {
		t.Fatalf("donut: %d polys, %d holes; want 1, 1", len(polys), len(polys[0].Holes))
	}
	if len(polys[0].GeoLoop) < 6 || len(polys[0].Holes[0]) != 6 {
		t.Errorf("donut loops: outer %d, hole %d", len(polys[0].GeoLoop), len(polys[0].Holes[0]))
	}

	// Two disjoint cells produce two polygons.
	far := mustCell(t)(LatLngToCell(LatLngDegs(0, 0), 9))
	polys, err = CellsToMultiPolygon([]Cell{sfCellRes9, far})
	if err != nil {
		t.Fatal(err)
	}
	if len(polys) != 2 {
		t.Fatalf("disjoint cells: %d polys, want 2", len(polys))
	}
}

func TestMetrics(t *testing.T) {
	t.Parallel()

	// Exact cell area: res-9 hexagon ~0.1 km^2; unit consistency.
	km2, err := sfCellRes9.AreaKm2()
	if err != nil {
		t.Fatal(err)
	}
	if km2 < 0.05 || km2 > 0.2 {
		t.Errorf("res-9 cell area = %f km2", km2)
	}
	m2, err := sfCellRes9.AreaM2()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(m2-km2*1e6)/m2 > 1e-9 {
		t.Errorf("m2/km2 inconsistent: %f vs %f", m2, km2*1e6)
	}
	rads2, err := sfCellRes9.AreaRads2()
	if err != nil || rads2 <= 0 {
		t.Errorf("rads2 = %f (%v)", rads2, err)
	}

	// Edge lengths: consistent units, plausible magnitude (res-9 edge ~200m).
	edges, err := sfCellRes9.DirectedEdges()
	if err != nil {
		t.Fatal(err)
	}
	lkm, err := edges[0].LengthKm()
	if err != nil {
		t.Fatal(err)
	}
	lm, err := edges[0].LengthM()
	if err != nil {
		t.Fatal(err)
	}
	if lkm < 0.1 || lkm > 0.3 || math.Abs(lm-lkm*1000)/lm > 1e-9 {
		t.Errorf("edge lengths: %f km, %f m", lkm, lm)
	}

	// Average metrics against known upstream values.
	a0, err := HexagonAreaAvgKm2(0)
	if err != nil || math.Abs(a0-4357449.416078381)/a0 > 1e-9 {
		t.Errorf("HexagonAreaAvgKm2(0) = %f (%v)", a0, err)
	}
	e0, err := HexagonEdgeLengthAvgKm(0)
	if err != nil || e0 < 1000 || e0 > 1300 {
		t.Errorf("HexagonEdgeLengthAvgKm(0) = %f (%v)", e0, err)
	}
	if _, err := HexagonAreaAvgM2(16); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("res 16: got %v", err)
	}
	if _, err := HexagonEdgeLengthAvgM(-1); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("res -1: got %v", err)
	}

	// Great-circle distance: SF to LA ~559 km.
	sf := LatLngDegs(37.775938728915946, -122.41795063018799)
	la := LatLngDegs(34.052235, -118.243683)
	dkm := GreatCircleDistanceKm(sf, la)
	if dkm < 540 || dkm > 580 {
		t.Errorf("SF-LA distance = %f km", dkm)
	}
	dm := GreatCircleDistanceM(sf, la)
	if math.Abs(dm-dkm*1000)/dm > 1e-9 {
		t.Errorf("m/km inconsistent: %f vs %f", dm, dkm*1000)
	}
	drads := GreatCircleDistanceRads(sf, la)
	if math.Abs(drads*earthRadiusKm-dkm)/dkm > 1e-9 {
		t.Errorf("rads/km inconsistent")
	}
}

func TestWave3Allocations(t *testing.T) {
	assertAllocs := func(name string, want float64, f func()) {
		t.Helper()
		if got := testing.AllocsPerRun(100, f); got > want {
			t.Errorf("%s allocates %v/run, want <= %v", name, got, want)
		}
	}

	assertAllocs("Cell.AreaKm2", 0, func() { _, _ = sfCellRes9.AreaKm2() })
	sf := LatLngDegs(37.7, -122.4)
	la := LatLngDegs(34.0, -118.2)
	assertAllocs("GreatCircleDistanceKm", 0, func() { _ = GreatCircleDistanceKm(sf, la) })

	parent, _ := sfCellRes9.Parent(4)
	children, err := parent.Children(6)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]Cell, 0, len(children))
	// compactCells allocates internal working arrays per call, exactly as the
	// C implementation mallocs them; the destination itself is reused. The
	// count is 3 or 4 depending on platform (escape analysis keeps one array
	// on the stack on darwin/arm64 but not on linux/amd64).
	assertAllocs("AppendCompactCells warm", 4, func() {
		out, err := AppendCompactCells(buf, children)
		if err != nil || len(out) != 1 {
			t.Fatal(err, len(out))
		}
	})
	compacted := []Cell{parent}
	big := make([]Cell, 0, 2500)
	assertAllocs("AppendUncompactCells warm", 0, func() {
		out, err := AppendUncompactCells(big, compacted, 6)
		if err != nil || len(out) != len(children) {
			t.Fatal(err, len(out))
		}
	})
}
