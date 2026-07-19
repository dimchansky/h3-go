//go:build cgo && c2go && h3v450

package h3

import "testing"

// Isolated (non-composite) parity for the cellsToMultiPoly functions
// that the pipeline-level tests in cellsToMultiPoly_parity_test.go
// exercise only in combination: cellToEdgeArcs, the createArcSet hash
// buckets, resetVisited, unionArcs/getRoot (incl. rank state),
// createSortableLoop, the three comparators, createSortablePoly,
// createMultiPolygon, the destroy helpers with observable post-state,
// and the H3 4.5.0 destroyLinkedMultiPolygon idempotence delta.

func Test_cellToEdgeArcs_direct_parity(t *testing.T) {
	// A hexagon and a pentagon cover both idx tables and edge counts.
	cells := []h3Index{0x890dab6220bffff, 0x851c0003fffffff}
	for _, cell := range cells {
		var goArcs [6]arc
		var goNum int64
		if err := cellToEdgeArcs(cell, goArcs[:], &goNum); err != eSuccess {
			t.Fatalf("cellToEdgeArcs(%x): %v", uint64(cell), err)
		}
		cSt, err := cellToEdgeArcsC(cell)
		if err != eSuccess {
			t.Fatalf("cellToEdgeArcsC(%x): %v", uint64(cell), err)
		}
		if int64(len(cSt.ids)) != goNum {
			t.Fatalf("%x: numEdges Go=%d C=%d", uint64(cell), goNum, len(cSt.ids))
		}
		idx := map[*arc]int64{}
		for i := range goArcs {
			idx[&goArcs[i]] = int64(i)
		}
		for i := int64(0); i < goNum; i++ {
			a := &goArcs[i]
			if a.id != cSt.ids[i] || idx[a.next] != cSt.nextIdx[i] ||
				idx[a.prev] != cSt.prevIdx[i] || idx[a.parent] != cSt.parentIdx[i] ||
				a.rank != cSt.rank[i] {
				t.Fatalf("%x arc %d: Go={%x n%d p%d par%d r%d} C={%x n%d p%d par%d r%d}",
					uint64(cell), i, uint64(a.id), idx[a.next], idx[a.prev], idx[a.parent], a.rank,
					uint64(cSt.ids[i]), cSt.nextIdx[i], cSt.prevIdx[i], cSt.parentIdx[i], cSt.rank[i])
			}
		}
	}
}

func Test_createArcSet_buckets_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		n := int64(len(cells))
		var arcset arcSet
		if err := createArcSet(cells, n, &arcset); err != eSuccess {
			t.Fatalf("createArcSet(%s): %v", name, err)
		}
		idx := map[*arc]int64{}
		for i := range arcset.arcs {
			idx[&arcset.arcs[i]] = int64(i)
		}
		cBuckets, err := bucketStateC(cells, n)
		if err != eSuccess {
			t.Fatalf("bucketStateC(%s): %v", name, err)
		}
		if int64(len(cBuckets)) != arcset.numBuckets {
			t.Fatalf("%s: numBuckets Go=%d C=%d", name, arcset.numBuckets, len(cBuckets))
		}
		for j := range arcset.buckets {
			goIdx := int64(-1)
			if arcset.buckets[j] != nil {
				goIdx = idx[arcset.buckets[j]]
			}
			if goIdx != cBuckets[j] {
				t.Fatalf("%s bucket %d: Go arc %d, C arc %d", name, j, goIdx, cBuckets[j])
			}
		}
	}
}

func Test_resetVisited_direct_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		n := int64(len(cells))
		for _, mode := range []int32{0, 1} {
			var arcset arcSet
			if err := createArcSet(cells, n, &arcset); err != eSuccess {
				t.Fatalf("createArcSet: %v", err)
			}
			if err := cancelArcPairs(arcset); err != eSuccess {
				t.Fatalf("cancelArcPairs: %v", err)
			}
			countLoops(arcset)
			if mode >= 1 {
				resetVisited(arcset)
			}
			cVisited, err := visitedStateC(cells, n, mode)
			if err != eSuccess {
				t.Fatalf("visitedStateC(%s, %d): %v", name, mode, err)
			}
			for i := int64(0); i < arcset.numArcs; i++ {
				if arcset.arcs[i].isVisited != cVisited[i] {
					t.Fatalf("%s mode %d arc %d: Go visited=%v C=%v",
						name, mode, i, arcset.arcs[i].isVisited, cVisited[i])
				}
			}
		}
	}
}

func Test_unionArcs_getRoot_direct_parity(t *testing.T) {
	cells := multiPolyParitySets(t)["hole"]
	n := int64(len(cells))
	// Pair sequences exercising the rank-swap branch, tree merging
	// through getRoot path compression, and the a == b (already same
	// component) no-op branch.
	sequences := [][][2]int64{
		{{0, 6}},
		{{0, 6}, {12, 18}, {0, 12}},
		{{0, 6}, {6, 12}, {0, 12}, {18, 24}, {24, 30}, {0, 18}},
		{{0, 1}}, // same per-cell component already (a == b branch)
	}
	for si, pairs := range sequences {
		var arcset arcSet
		if err := createArcSet(cells, n, &arcset); err != eSuccess {
			t.Fatalf("createArcSet: %v", err)
		}
		for _, p := range pairs {
			unionArcs(&arcset.arcs[p[0]], &arcset.arcs[p[1]])
		}
		cRoots, cRanks, err := unionSequenceC(cells, n, pairs)
		if err != eSuccess {
			t.Fatalf("unionSequenceC(seq %d): %v", si, err)
		}
		for i := int64(0); i < arcset.numArcs; i++ {
			goRoot := getRoot(&arcset.arcs[i]).id
			if goRoot != cRoots[i] || arcset.arcs[i].rank != cRanks[i] {
				t.Fatalf("seq %d arc %d: Go root=%x rank=%d, C root=%x rank=%d",
					si, i, uint64(goRoot), arcset.arcs[i].rank,
					uint64(cRoots[i]), cRanks[i])
			}
		}
	}
}

func Test_createSortableLoop_direct_parity(t *testing.T) {
	for _, name := range []string{"singleHex", "hole", "pentagon"} {
		cells := multiPolyParitySets(t)[name]
		n := int64(len(cells))
		var arcset arcSet
		if err := createArcSet(cells, n, &arcset); err != eSuccess {
			t.Fatalf("createArcSet: %v", err)
		}
		if err := cancelArcPairs(arcset); err != eSuccess {
			t.Fatalf("cancelArcPairs: %v", err)
		}
		// Both sides have identical arc state, so surviving indices
		// coincide; call the function directly on each surviving arc.
		for i := int64(0); i < arcset.numArcs; i++ {
			if arcset.arcs[i].isRemoved {
				continue
			}
			resetVisited(arcset)
			var goLoop sortableLoop
			if err := createSortableLoop(&arcset.arcs[i], &goLoop); err != eSuccess {
				t.Fatalf("createSortableLoop(%s, %d): %v", name, i, err)
			}
			cLoop, err := createSortableLoopC(cells, n, i)
			if err != eSuccess {
				t.Fatalf("createSortableLoopC(%s, %d): %v", name, i, err)
			}
			if goLoop.root != cLoop.root || len(goLoop.loop) != len(cLoop.loop) {
				t.Fatalf("%s arc %d: Go(root %x, %d verts) C(root %x, %d verts)",
					name, i, uint64(goLoop.root), len(goLoop.loop),
					uint64(cLoop.root), len(cLoop.loop))
			}
			if !areaClose(goLoop.area, cLoop.area) {
				t.Fatalf("%s arc %d: area Go=%v C=%v", name, i, goLoop.area, cLoop.area)
			}
			for j := range goLoop.loop {
				if !vec3UlpClose(goLoop.loop[j].Lat.Rad(), cLoop.loop[j].Lat.Rad()) ||
					!vec3UlpClose(goLoop.loop[j].Lng.Rad(), cLoop.loop[j].Lng.Rad()) {
					t.Fatalf("%s arc %d vert %d: Go=%v C=%v", name, i, j, goLoop.loop[j], cLoop.loop[j])
				}
			}
		}
	}
}

func Test_comparators_direct_parity(t *testing.T) {
	// cmp_SortableLoop: every (root, area) branch combination.
	loopCases := []struct {
		rootA h3Index
		areaA float64
		rootB h3Index
		areaB float64
	}{
		{1, 1.0, 2, 1.0}, {2, 1.0, 1, 1.0}, // root <, >
		{1, 1.0, 1, 2.0}, {1, 2.0, 1, 1.0}, // same root; area <, >
		{1, 1.0, 1, 1.0}, // fully equal
	}
	for _, tc := range loopCases {
		goA := sortableLoop{root: tc.rootA, area: tc.areaA}
		goB := sortableLoop{root: tc.rootB, area: tc.areaB}
		if got, want := cmp_SortableLoop(&goA, &goB), cmp_SortableLoopC(tc.rootA, tc.areaA, tc.rootB, tc.areaB); got != want {
			t.Errorf("cmp_SortableLoop(%v): Go=%d C=%d", tc, got, want)
		}
	}

	// cmp_SortablePoly: descending, ascending, equal.
	polyCases := [][2]float64{{200, 100}, {100, 200}, {100, 100}}
	for _, tc := range polyCases {
		goA := sortablePoly{outerArea: tc[0]}
		goB := sortablePoly{outerArea: tc[1]}
		if got, want := cmp_SortablePoly(&goA, &goB), cmp_SortablePolyC(tc[0], tc[1]); got != want {
			t.Errorf("cmp_SortablePoly(%v): Go=%d C=%d", tc, got, want)
		}
	}

	// cmp_uint64: <, >, ==, including the high-bit range.
	u64Cases := [][2]h3Index{
		{1, 2}, {2, 1}, {5, 5},
		{0x8000000000000000, 1}, {1, 0x8000000000000000},
	}
	for _, tc := range u64Cases {
		if got, want := cmp_uint64(tc[0], tc[1]), cmp_uint64C(tc[0], tc[1]); got != want {
			t.Errorf("cmp_uint64(%x, %x): Go=%d C=%d", uint64(tc[0]), uint64(tc[1]), got, want)
		}
	}
}

func Test_createSortablePoly_direct_parity(t *testing.T) {
	// The "hole" set sorts into one polygon: outer loop first, one
	// hole after it.
	cells := multiPolyParitySets(t)["hole"]
	n := int64(len(cells))
	var arcset arcSet
	if err := createArcSet(cells, n, &arcset); err != eSuccess {
		t.Fatalf("createArcSet: %v", err)
	}
	if err := cancelArcPairs(arcset); err != eSuccess {
		t.Fatalf("cancelArcPairs: %v", err)
	}
	var loopset sortableLoopSet
	if err := createSortableLoopSet(arcset, &loopset); err != eSuccess {
		t.Fatalf("createSortableLoopSet: %v", err)
	}
	if loopset.numLoops != 2 {
		t.Fatalf("hole set: %d loops, want 2", loopset.numLoops)
	}
	for _, numHoles := range []int64{0, 1} {
		var goPoly sortablePoly
		if err := createSortablePoly(loopset.sloops, numHoles, &goPoly); err != eSuccess {
			t.Fatalf("createSortablePoly(holes %d): %v", numHoles, err)
		}
		cPoly, err := createSortablePolyC(cells, n, 0, numHoles)
		if err != eSuccess {
			t.Fatalf("createSortablePolyC(holes %d): %v", numHoles, err)
		}
		if !areaClose(goPoly.outerArea, cPoly.outerArea) ||
			len(goPoly.poly.GeoLoop) != len(cPoly.poly.GeoLoop) ||
			len(goPoly.poly.Holes) != len(cPoly.poly.Holes) {
			t.Fatalf("holes %d: Go(area %v, %d verts, %d holes) C(area %v, %d verts, %d holes)",
				numHoles, goPoly.outerArea, len(goPoly.poly.GeoLoop), len(goPoly.poly.Holes),
				cPoly.outerArea, len(cPoly.poly.GeoLoop), len(cPoly.poly.Holes))
		}
		for h := range goPoly.poly.Holes {
			if len(goPoly.poly.Holes[h]) != len(cPoly.poly.Holes[h]) {
				t.Fatalf("holes %d hole %d: verts Go=%d C=%d",
					numHoles, h, len(goPoly.poly.Holes[h]), len(cPoly.poly.Holes[h]))
			}
		}
	}
}

func Test_createMultiPolygon_direct_parity(t *testing.T) {
	// Isolated from validate/overflow: the pipeline runs to the sorted
	// loop set, then createMultiPolygon only. The globe set drives its
	// createGlobeMultiPolygon branch (empty loop set).
	for name, cells := range multiPolyParitySets(t) {
		n := int64(len(cells))
		var arcset arcSet
		if err := createArcSet(cells, n, &arcset); err != eSuccess {
			t.Fatalf("createArcSet(%s): %v", name, err)
		}
		if err := cancelArcPairs(arcset); err != eSuccess {
			t.Fatalf("cancelArcPairs(%s): %v", name, err)
		}
		var loopset sortableLoopSet
		if err := createSortableLoopSet(arcset, &loopset); err != eSuccess {
			t.Fatalf("createSortableLoopSet(%s): %v", name, err)
		}
		var goOut geoMultiPolygon
		if err := createMultiPolygon(loopset, &goOut); err != eSuccess {
			t.Fatalf("createMultiPolygon(%s): %v", name, err)
		}
		cOut, err := createMultiPolygonOnlyC(cells, n)
		if err != eSuccess {
			t.Fatalf("createMultiPolygonOnlyC(%s): %v", name, err)
		}
		assertMultiPolyParity(t, name, goOut, cOut, name == "globeAllRes0")
	}
}

func Test_destroyHelpers_state_parity(t *testing.T) {
	cells := multiPolyParitySets(t)["hole"]
	n := int64(len(cells))

	// destroyArcSet: Go side.
	var arcset arcSet
	if err := createArcSet(cells, n, &arcset); err != eSuccess {
		t.Fatalf("createArcSet: %v", err)
	}
	destroyArcSet(&arcset)
	goBits := int32(0)
	if arcset.arcs == nil && arcset.buckets == nil {
		goBits |= 1
	}
	destroyArcSet(&arcset)
	if arcset.arcs == nil && arcset.buckets == nil {
		goBits |= 2
	}
	if cBits := destroyArcSetStateC(cells, n); goBits != 3 || cBits != 3 {
		t.Errorf("destroyArcSet state: Go bits=%d C bits=%d, want 3/3", goBits, cBits)
	}

	// destroySortableLoopSet / ...Shallow: Go side.
	for _, shallow := range []int32{0, 1} {
		var as arcSet
		if err := createArcSet(cells, n, &as); err != eSuccess {
			t.Fatalf("createArcSet: %v", err)
		}
		if err := cancelArcPairs(as); err != eSuccess {
			t.Fatalf("cancelArcPairs: %v", err)
		}
		var loopset sortableLoopSet
		if err := createSortableLoopSet(as, &loopset); err != eSuccess {
			t.Fatalf("createSortableLoopSet: %v", err)
		}
		destroy := destroySortableLoopSet
		if shallow == 1 {
			destroy = destroySortableLoopSetShallow
		}
		destroy(&loopset)
		bits := int32(0)
		if loopset.sloops == nil {
			bits |= 1
		}
		destroy(&loopset)
		if loopset.sloops == nil {
			bits |= 2
		}
		if cBits := destroyLoopSetStateC(cells, n, shallow); bits != 3 || cBits != 3 {
			t.Errorf("destroyLoopSet(shallow=%d) state: Go bits=%d C bits=%d, want 3/3", shallow, bits, cBits)
		}
	}

	// destroyGeoLoop: Go side.
	loop := make(GeoLoop, 3)
	destroyGeoLoop(&loop)
	goBits = 0
	if loop == nil {
		goBits |= 1
	}
	destroyGeoLoop(&loop)
	if loop == nil {
		goBits |= 2
	}
	if cBits := destroyGeoLoopStateC(); goBits != 3 || cBits != 3 {
		t.Errorf("destroyGeoLoop state: Go bits=%d C bits=%d, want 3/3", goBits, cBits)
	}

	// destroyGeoPolygon: Go side.
	poly := GeoPolygon{GeoLoop: make(GeoLoop, 3), Holes: []GeoLoop{make(GeoLoop, 3)}}
	destroyGeoPolygon(&poly)
	goBits = 0
	if poly.GeoLoop == nil && poly.Holes == nil {
		goBits |= 1
	}
	destroyGeoPolygon(&poly)
	if poly.GeoLoop == nil && poly.Holes == nil {
		goBits |= 2
	}
	if cBits := destroyGeoPolygonStateC(); goBits != 3 || cBits != 3 {
		t.Errorf("destroyGeoPolygon state: Go bits=%d C bits=%d, want 3/3", goBits, cBits)
	}

	// destroyGeoMultiPolygon: Go side.
	var mpoly geoMultiPolygon
	if err := cellsToMultiPolygon(cells, n, &mpoly); err != eSuccess {
		t.Fatalf("cellsToMultiPolygon: %v", err)
	}
	destroyGeoMultiPolygon(&mpoly)
	goBits = 0
	if mpoly.Polygons == nil && mpoly.NumPolygons == 0 {
		goBits |= 1
	}
	destroyGeoMultiPolygon(&mpoly)
	if mpoly.Polygons == nil && mpoly.NumPolygons == 0 {
		goBits |= 2
	}
	if cBits := destroyGeoMultiPolygonStateC(cells, n); goBits != 3 || cBits != 3 {
		t.Errorf("destroyGeoMultiPolygon state: Go bits=%d C bits=%d, want 3/3", goBits, cBits)
	}
}

func Test_destroyLinkedMultiPolygon_idempotence_parity(t *testing.T) {
	// H3 4.5.0 delta (record §7 item 7): destroyLinkedMultiPolygon
	// zeroes the caller-owned head node, so a second call is a no-op.
	// Both sides destroy twice; neither second call may fail (the C
	// side would crash the harness), and the head must be zeroed after
	// both calls. C side gated to 4.5.0 (this file's h3v450 build tag):
	// the 4.4.0 implementation did not provide this contract.
	for name, cells := range multiPolyParitySets(t) {
		var linked linkedGeoPolygon
		if err := cellsToLinkedMultiPolygon(cells, int32(len(cells)), &linked); err != eSuccess {
			t.Fatalf("cellsToLinkedMultiPolygon(%s): %v", name, err)
		}
		destroyLinkedMultiPolygon(&linked)
		goBits := int32(0)
		if linked.First == nil && linked.Last == nil && linked.Next == nil {
			goBits |= 1
		}
		destroyLinkedMultiPolygon(&linked)
		if linked.First == nil && linked.Last == nil && linked.Next == nil {
			goBits |= 2
		}
		cBits := destroyLinkedTwiceC(cells, int32(len(cells)))
		if goBits != 3 || cBits != 3 {
			t.Errorf("%s: idempotence Go bits=%d C bits=%d, want 3/3", name, goBits, cBits)
		}
	}
}
