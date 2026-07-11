package h3

import (
	"errors"
	"slices"
	"testing"
)

func TestGridDisk(t *testing.T) {
	t.Parallel()

	disk, err := sfCellRes9.GridDisk(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(disk) != 7 {
		t.Fatalf("hexagon 1-disk = %d cells, want 7", len(disk))
	}
	if !slices.Contains(disk, sfCellRes9) {
		t.Error("disk must contain the origin")
	}
	for _, c := range disk {
		if !c.IsValid() {
			t.Errorf("invalid cell %v in disk", c)
		}
		d, err := sfCellRes9.GridDistance(c)
		if err != nil || d > 1 {
			t.Errorf("cell %v at distance %d (%v)", c, d, err)
		}
	}

	if _, err := sfCellRes9.GridDisk(-1); !errors.Is(err, ErrDomain) {
		t.Errorf("negative k: got %v, want ErrDomain", err)
	}

	// Pentagon disks exercise the hole-pruning path: 1-disk of a pentagon has
	// 6 cells (5 neighbors + origin), not 7.
	pentagons, err := Pentagons(4)
	if err != nil {
		t.Fatal(err)
	}
	pdisk, err := pentagons[0].GridDisk(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(pdisk) != 6 {
		t.Fatalf("pentagon 1-disk = %d cells, want 6 (holes pruned)", len(pdisk))
	}
	for _, c := range pdisk {
		if c == h3Null {
			t.Fatal("pruning left a null cell")
		}
	}
}

func TestGridDiskUnsafeVariants(t *testing.T) {
	t.Parallel()

	// Hexagon-only neighborhoods succeed and are ring-walk ordered.
	disk, err := sfCellRes9.GridDiskUnsafe(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(disk) != 19 || disk[0] != sfCellRes9 {
		t.Fatalf("unsafe 2-disk: len %d, first %v", len(disk), disk[0])
	}

	// Pentagons must fail with ErrPentagon.
	pentagons, _ := Pentagons(4)
	if _, err := pentagons[0].GridDiskUnsafe(1); !errors.Is(err, ErrPentagon) {
		t.Errorf("pentagon GridDiskUnsafe: got %v, want ErrPentagon", err)
	}
	if _, _, err := pentagons[0].GridDiskDistancesUnsafe(1); !errors.Is(err, ErrPentagon) {
		t.Errorf("pentagon GridDiskDistancesUnsafe: got %v, want ErrPentagon", err)
	}
	if _, err := GridDisksUnsafe([]Cell{sfCellRes9, pentagons[0]}, 1); !errors.Is(err, ErrPentagon) {
		t.Errorf("GridDisksUnsafe with pentagon: got %v, want ErrPentagon", err)
	}

	multi, err := GridDisksUnsafe([]Cell{sfCellRes9, sfCellRes9}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(multi) != 14 {
		t.Fatalf("GridDisksUnsafe len = %d, want 14", len(multi))
	}
}

func TestGridDiskDistances(t *testing.T) {
	t.Parallel()

	cells, dists, err := sfCellRes9.GridDiskDistances(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(cells) != 19 || len(dists) != 19 {
		t.Fatalf("lens = %d/%d, want 19/19", len(cells), len(dists))
	}
	counts := map[int32]int{}
	for i, c := range cells {
		counts[dists[i]]++
		want, err := sfCellRes9.GridDistance(c)
		if err != nil || int32(want) != dists[i] {
			t.Errorf("cell %v: distance %d, GridDistance says %d (%v)", c, dists[i], want, err)
		}
	}
	if counts[0] != 1 || counts[1] != 6 || counts[2] != 12 {
		t.Errorf("distance histogram = %v, want 1/6/12", counts)
	}

	// Safe variant must agree as a set.
	sCells, sDists, err := sfCellRes9.GridDiskDistancesSafe(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(sCells) != 19 || len(sDists) != 19 {
		t.Fatalf("safe lens = %d/%d", len(sCells), len(sDists))
	}

	// Pentagon: pruned in tandem.
	pentagons, _ := Pentagons(4)
	pCells, pDists, err := pentagons[0].GridDiskDistances(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(pCells) != 6 || len(pDists) != 6 {
		t.Fatalf("pentagon lens = %d/%d, want 6/6", len(pCells), len(pDists))
	}
}

func TestGridRing(t *testing.T) {
	t.Parallel()

	ring, err := sfCellRes9.GridRing(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(ring) != 12 {
		t.Fatalf("2-ring = %d cells, want 12", len(ring))
	}
	for _, c := range ring {
		d, err := sfCellRes9.GridDistance(c)
		if err != nil || d != 2 {
			t.Errorf("ring cell %v at distance %d (%v)", c, d, err)
		}
	}

	ring0, err := sfCellRes9.GridRing(0)
	if err != nil || len(ring0) != 1 || ring0[0] != sfCellRes9 {
		t.Errorf("0-ring = %v (%v)", ring0, err)
	}

	// Unsafe ordered variant on hexagons agrees as a set.
	ringU, err := sfCellRes9.GridRingUnsafe(2)
	if err != nil {
		t.Fatal(err)
	}
	a, b := slices.Clone(ring), slices.Clone(ringU)
	slices.Sort(a)
	slices.Sort(b)
	if !slices.Equal(a, b) {
		t.Error("GridRing and GridRingUnsafe disagree")
	}

	// Pentagon-adjacent rings prune holes: ring around a pentagon still only
	// contains valid cells.
	pentagons, _ := Pentagons(4)
	pring, err := pentagons[0].GridRing(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(pring) != 5 {
		t.Fatalf("pentagon 1-ring = %d, want 5", len(pring))
	}
	if _, err := pentagons[0].GridRingUnsafe(1); !errors.Is(err, ErrPentagon) {
		t.Errorf("pentagon GridRingUnsafe: got %v, want ErrPentagon", err)
	}
}

func TestGridPathAndDistance(t *testing.T) {
	t.Parallel()

	other := mustCell(t)(LatLngToCell(LatLngDegs(37.8, -122.5), 9))
	d, err := sfCellRes9.GridDistance(other)
	if err != nil {
		t.Fatal(err)
	}
	if d <= 0 {
		t.Fatalf("distance = %d", d)
	}

	n, err := sfCellRes9.GridPathLen(other)
	if err != nil || n != d+1 {
		t.Fatalf("GridPathLen = %d (%v), want %d", n, err, d+1)
	}

	path, err := sfCellRes9.GridPath(other)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) != n || path[0] != sfCellRes9 || path[len(path)-1] != other {
		t.Fatalf("path len %d endpoints %v..%v", len(path), path[0], path[len(path)-1])
	}
	for i := 1; i < len(path); i++ {
		ok, err := path[i-1].IsNeighbor(path[i])
		if err != nil || !ok {
			t.Errorf("path step %d not adjacent (%v)", i, err)
		}
	}
}

func TestLocalIJRoundTrip(t *testing.T) {
	t.Parallel()

	disk, err := sfCellRes9.GridDisk(3)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range disk {
		ij, err := CellToLocalIJ(sfCellRes9, c)
		if err != nil {
			t.Fatalf("CellToLocalIJ(%v): %v", c, err)
		}
		back, err := LocalIJToCell(sfCellRes9, ij)
		if err != nil || back != c {
			t.Errorf("LocalIJ round trip %v -> %+v -> %v (%v)", c, ij, back, err)
		}
	}
}

func TestDirectedEdges(t *testing.T) {
	t.Parallel()

	edges, err := sfCellRes9.DirectedEdges()
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 6 {
		t.Fatalf("hexagon edges = %d, want 6", len(edges))
	}
	for _, e := range edges {
		if !e.IsValid() {
			t.Fatalf("invalid edge %v", e)
		}
		if e.Resolution() != 9 {
			t.Errorf("edge resolution = %d", e.Resolution())
		}
		origin, err := e.Origin()
		if err != nil || origin != sfCellRes9 {
			t.Errorf("origin = %v (%v)", origin, err)
		}
		dest, err := e.Destination()
		if err != nil {
			t.Fatal(err)
		}
		o2, d2, err := e.Cells()
		if err != nil || o2 != origin || d2 != dest {
			t.Errorf("Cells() = %v,%v (%v)", o2, d2, err)
		}
		ok, err := origin.IsNeighbor(dest)
		if err != nil || !ok {
			t.Errorf("edge endpoints not neighbors (%v)", err)
		}
		back, err := origin.DirectedEdgeTo(dest)
		if err != nil || back != e {
			t.Errorf("DirectedEdgeTo = %v, want %v (%v)", back, e, err)
		}
		b, err := e.Boundary()
		if err != nil || b.Len() < 2 {
			t.Errorf("edge boundary len %d (%v)", b.Len(), err)
		}
	}

	// Pentagons have 5 edges (deleted k-axis hole pruned).
	pentagons, _ := Pentagons(4)
	pEdges, err := pentagons[0].DirectedEdges()
	if err != nil {
		t.Fatal(err)
	}
	if len(pEdges) != 5 {
		t.Fatalf("pentagon edges = %d, want 5", len(pEdges))
	}

	// Non-neighbors.
	far := mustCell(t)(LatLngToCell(LatLngDegs(0, 0), 9))
	if _, err := sfCellRes9.DirectedEdgeTo(far); !errors.Is(err, ErrNotNeighbors) {
		t.Errorf("far DirectedEdgeTo: got %v, want ErrNotNeighbors", err)
	}
	ok, err := sfCellRes9.IsNeighbor(far)
	if err != nil || ok {
		t.Errorf("far IsNeighbor = %v (%v)", ok, err)
	}
	if _, err := sfCellRes9.IsNeighbor(sfCellRes15); !errors.Is(err, ErrResolutionMismatch) {
		t.Errorf("cross-res IsNeighbor: got %v, want ErrResolutionMismatch", err)
	}

	// A cell index is not a valid edge.
	if DirectedEdge(sfCellRes9).IsValid() {
		t.Error("cell index must not validate as edge")
	}
}

func TestVertexes(t *testing.T) {
	t.Parallel()

	verts, err := sfCellRes9.Vertexes()
	if err != nil {
		t.Fatal(err)
	}
	if len(verts) != 6 {
		t.Fatalf("hexagon vertexes = %d, want 6", len(verts))
	}
	seen := map[Vertex]bool{}
	for i, v := range verts {
		if !v.IsValid() || seen[v] {
			t.Fatalf("vertex %d invalid or duplicate: %v", i, v)
		}
		seen[v] = true
		if v.Resolution() != 9 {
			t.Errorf("vertex resolution = %d", v.Resolution())
		}
		single, err := sfCellRes9.Vertex(i)
		if err != nil || single != v {
			t.Errorf("Vertex(%d) = %v, want %v (%v)", i, single, v, err)
		}
		ll, err := v.LatLng()
		if err != nil {
			t.Fatal(err)
		}
		// Vertex coordinates must appear in the cell boundary.
		b, _ := sfCellRes9.Boundary()
		found := false
		for _, bv := range b.Verts() {
			if bv.Lat.EqualApprox(ll.Lat, 1e-9) && bv.Lng.EqualApprox(ll.Lng, 1e-9) {
				found = true
			}
		}
		if !found {
			t.Errorf("vertex %v coords %v not on boundary", v, ll)
		}
	}

	pentagons, _ := Pentagons(4)
	pVerts, err := pentagons[0].Vertexes()
	if err != nil {
		t.Fatal(err)
	}
	if len(pVerts) != 5 {
		t.Fatalf("pentagon vertexes = %d, want 5", len(pVerts))
	}

	if _, err := sfCellRes9.Vertex(6); !errors.Is(err, ErrDomain) {
		t.Errorf("Vertex(6): got %v, want ErrDomain", err)
	}
	if Vertex(sfCellRes9).IsValid() {
		t.Error("cell index must not validate as vertex")
	}
}

func TestWave2Allocations(t *testing.T) {
	assertAllocs := func(name string, want float64, f func()) {
		t.Helper()
		if got := testing.AllocsPerRun(200, f); got > want {
			t.Errorf("%s allocates %v/run, want <= %v", name, got, want)
		}
	}

	// AppendGridDisk warm path: 1 internal distance-scratch alloc remains
	// (the C implementation heap-allocates the same scratch).
	buf := make([]Cell, 0, 64)
	assertAllocs("AppendGridDisk warm", 1, func() {
		out, err := sfCellRes9.AppendGridDisk(buf, 2)
		if err != nil || len(out) != 19 {
			t.Fatal(err, len(out))
		}
	})

	// AppendGridDiskDistances with caller buffers: 0 allocs.
	distBuf := make([]int32, 0, 64)
	assertAllocs("AppendGridDiskDistances warm", 0, func() {
		c, d, err := sfCellRes9.AppendGridDiskDistances(buf, distBuf, 2)
		if err != nil || len(c) != 19 || len(d) != 19 {
			t.Fatal(err)
		}
	})

	assertAllocs("AppendGridRingUnsafe warm", 0, func() {
		out, err := sfCellRes9.AppendGridRingUnsafe(buf, 2)
		if err != nil || len(out) != 12 {
			t.Fatal(err, len(out))
		}
	})

	ring, err := sfCellRes9.GridRing(1)
	if err != nil || len(ring) == 0 {
		t.Fatal(err)
	}
	other := ring[0]
	assertAllocs("GridDistance", 0, func() {
		if _, err := sfCellRes9.GridDistance(other); err != nil {
			t.Fatal(err)
		}
	})
	assertAllocs("CellToLocalIJ", 0, func() {
		if _, err := CellToLocalIJ(sfCellRes9, other); err != nil {
			t.Fatal(err)
		}
	})

	edge, err := sfCellRes9.DirectedEdgeTo(other)
	if err != nil {
		t.Fatal(err)
	}
	assertAllocs("DirectedEdgeTo", 0, func() {
		if _, err := sfCellRes9.DirectedEdgeTo(other); err != nil {
			t.Fatal(err)
		}
	})
	assertAllocs("DirectedEdge.Boundary", 0, func() {
		if _, err := edge.Boundary(); err != nil {
			t.Fatal(err)
		}
	})
	assertAllocs("DirectedEdge.Cells", 0, func() {
		if _, _, err := edge.Cells(); err != nil {
			t.Fatal(err)
		}
	})
}
