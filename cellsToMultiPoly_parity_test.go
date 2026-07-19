//go:build cgo && c2go && h3v450

package h3

import (
	"sort"
	"testing"
)

// Direct parity for the H3 4.5.0 cellsToMultiPoly machinery against the
// exact 4.5.0 C implementation: the file-statics via same-TU h3goTest_*
// wrappers (validateCellSet, getNumEdges, hashEdge,
// checkCellsToMultiPolyOverflow, createArcSet/cellToEdgeArcs/findArc/
// cancelArcPairs/unionArcs/getRoot via serialized arc state, countLoops,
// createSortableLoop/createSortableLoopSet/cmp_SortableLoop/countPolys
// via the serialized loop set, createGlobeMultiPolygon) and the
// C-public entry points (cellsToMultiPolygon, and
// cellsToLinkedMultiPolygon with its linked output serialized directly
// — no conversion pipeline in between).
//
// Discrete outputs (error codes, counts, arc ids, linkage indices,
// component roots, vertex counts) compare exactly. Loop areas compare
// with areaClose (accumulated Cagnoli trig; 1e-14, upstream's own
// bound) and boundary vertices with vec3UlpClose (libm-dependent).

// Deterministic input sets drawn from the ported upstream suites.
func multiPolyParitySets(t *testing.T) map[string][]h3Index {
	t.Helper()
	sets := map[string][]h3Index{
		"singleHex": {0x890dab6220bffff},
		"contiguous2": {
			0x8928308291bffff, 0x89283082957ffff},
		"nonContiguous2": {
			0x8928308291bffff, 0x89283082943ffff},
		"hole": {
			0x892830828c7ffff, 0x892830828d7ffff, 0x8928308289bffff,
			0x89283082813ffff, 0x8928308288fffff, 0x89283082883ffff},
		"pentagon": {0x851c0003fffffff},
		"twoRing": {
			0x8930062838bffff, 0x8930062838fffff, 0x89300628383ffff,
			0x8930062839bffff, 0x893006283d7ffff, 0x893006283c7ffff,
			0x89300628313ffff, 0x89300628317ffff, 0x893006283bbffff,
			0x89300628387ffff, 0x89300628397ffff, 0x89300628393ffff,
			0x89300628067ffff, 0x8930062806fffff, 0x893006283d3ffff,
			0x893006283c3ffff, 0x893006283cfffff, 0x8930062831bffff,
			0x89300628303ffff},
		"nestedDonut": {
			0x89283082813ffff, 0x8928308281bffff, 0x8928308280bffff,
			0x8928308280fffff, 0x89283082807ffff, 0x89283082817ffff,
			0x8928308289bffff, 0x892830828d7ffff, 0x892830828c3ffff,
			0x892830828cbffff, 0x89283082853ffff, 0x89283082843ffff,
			0x8928308284fffff, 0x8928308287bffff, 0x89283082863ffff,
			0x89283082867ffff, 0x8928308282bffff, 0x89283082823ffff,
			0x89283082837ffff, 0x892830828afffff, 0x892830828a3ffff,
			0x892830828b3ffff, 0x89283082887ffff, 0x89283082883ffff},
		"threePolygonsRes0": {
			0x8027fffffffffff, 0x802bfffffffffff, 0x804dfffffffffff,
			0x8067fffffffffff, 0x806dfffffffffff, 0x8049fffffffffff,
			0x805ffffffffffff, 0x8057fffffffffff, 0x807dfffffffffff,
			0x80a5fffffffffff, 0x80a9fffffffffff, 0x808bfffffffffff,
			0x801bfffffffffff, 0x8035fffffffffff, 0x803ffffffffffff,
			0x8053fffffffffff, 0x8043fffffffffff, 0x8021fffffffffff,
			0x8011fffffffffff, 0x801ffffffffffff, 0x8097fffffffffff},
		// The exact input of upstream 4.5.0's new globalEquatorCells
		// regression (testCellsToLinkedMultiPolygon.c): a continuous
		// band of 60 res-1 cells around the equator.
		"globalEquatorCells": {
			0x81807ffffffffff, 0x817efffffffffff, 0x81723ffffffffff,
			0x817ebffffffffff, 0x817c3ffffffffff, 0x817e3ffffffffff,
			0x817a3ffffffffff, 0x8166fffffffffff, 0x8172bffffffffff,
			0x816afffffffffff, 0x81933ffffffffff, 0x8168fffffffffff,
			0x8188fffffffffff, 0x81853ffffffffff, 0x817f7ffffffffff,
			0x8180bffffffffff, 0x81783ffffffffff, 0x81743ffffffffff,
			0x8170bffffffffff, 0x8173bffffffffff, 0x8179bffffffffff,
			0x817cbffffffffff, 0x8188bffffffffff, 0x81857ffffffffff,
			0x816f7ffffffffff, 0x8177bffffffffff, 0x81617ffffffffff,
			0x816f3ffffffffff, 0x8174bffffffffff, 0x8180fffffffffff,
			0x817a7ffffffffff, 0x81767ffffffffff, 0x81757ffffffffff,
			0x81957ffffffffff, 0x81787ffffffffff, 0x81847ffffffffff,
			0x81653ffffffffff, 0x817bbffffffffff, 0x816cfffffffffff,
			0x816abffffffffff, 0x815f3ffffffffff, 0x817c7ffffffffff,
			0x8168bffffffffff, 0x818cbffffffffff, 0x818cfffffffffff,
			0x818afffffffffff, 0x8174fffffffffff, 0x8172fffffffffff,
			0x8170fffffffffff, 0x816fbffffffffff, 0x81657ffffffffff,
			0x816c7ffffffffff, 0x8186bffffffffff, 0x81763ffffffffff,
			0x818a7ffffffffff, 0x8186fffffffffff, 0x81707ffffffffff,
			0x8182bffffffffff, 0x818f3ffffffffff, 0x8182fffffffffff},
	}
	var res0 [122]h3Index
	if err := getRes0Cells(res0[:]); err != eSuccess {
		t.Fatalf("getRes0Cells: %v", err)
	}
	sets["globeAllRes0"] = res0[:]
	return sets
}

func Test_validateCellSet_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		goErr := validateCellSet(cells, int64(len(cells)))
		cErr := validateCellSetC(cells, int64(len(cells)))
		if goErr != cErr {
			t.Errorf("%s: Go=%v C=%v", name, goErr, cErr)
		}
	}

	// Error paths: negative count, invalid cell, resolution mismatch,
	// duplicates — all exact.
	errSets := map[string]struct {
		cells []h3Index
		n     int64
	}{
		"negative":    {nil, -1},
		"empty":       {nil, 0},
		"invalid":     {[]h3Index{0x8027fffffffffff, 0x81efbffffffffff + 1}, 2},
		"resMismatch": {[]h3Index{0x8027fffffffffff, 0x81efbffffffffff}, 2},
		"duplicates":  {[]h3Index{0x81efbffffffffff, 0x81efbffffffffff, 0x81efbffffffffff}, 3},
		// Invalid checked before resolution: an index that is both
		// invalid and res-mismatched reports E_CELL_INVALID.
		"invalidBeforeRes": {[]h3Index{0x8027fffffffffff, 0x8928308291bffff + 1}, 2},
	}
	for name, tc := range errSets {
		goErr := validateCellSet(tc.cells, tc.n)
		cErr := validateCellSetC(tc.cells, tc.n)
		if goErr != cErr {
			t.Errorf("%s: Go=%v C=%v", name, goErr, cErr)
		}
	}
}

func Test_hashEdge_getNumEdges_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		if got, want := getNumEdges(cells, int64(len(cells))), getNumEdgesC(cells, int64(len(cells))); got != want {
			t.Errorf("getNumEdges(%s): Go=%d C=%d", name, got, want)
		}
		for _, c := range cells {
			var edges [6]h3Index
			if err := originToDirectedEdges(c, edges[:]); err != eSuccess {
				t.Fatalf("originToDirectedEdges: %v", err)
			}
			for _, e := range edges {
				if e == h3Null {
					continue
				}
				for _, n := range []uint64{1, 7, 60, 1260, 7320} {
					if got, want := hashEdge(e, n), hashEdgeC(e, n); got != want {
						t.Fatalf("hashEdge(%x, %d): Go=%d C=%d", uint64(e), n, got, want)
					}
				}
			}
		}
	}
}

func Test_checkCellsToMultiPolyOverflow_parity(t *testing.T) {
	// Exercise both sides of the threshold with the pinned C sizes,
	// several multipliers, and degenerate counts — all exact. This
	// validates the Go port's pinned sizeof(Arc)/sizeof(Arc*) against
	// the C oracle.
	maxSafe := int64(cSizeMax / uint64(6*hashTableMultiplier*cSizeofArcPtr))
	cases := []struct{ numCells, mult int64 }{
		{-1, hashTableMultiplier},
		{0, hashTableMultiplier},
		{1000000, hashTableMultiplier},
		{maxSafe, hashTableMultiplier},
		{maxSafe + 1, hashTableMultiplier},
		{int64(^uint64(0) >> 1), hashTableMultiplier},
		{1000000, 1},
		{1000000, 100},
		{int64(cSizeMax / (6 * cSizeofArc)), 1},
		{int64(cSizeMax/(6*cSizeofArc)) + 1, 1},
	}
	for _, tc := range cases {
		goErr := checkCellsToMultiPolyOverflow(tc.numCells, tc.mult)
		cErr := checkCellsToMultiPolyOverflowC(tc.numCells, tc.mult)
		if goErr != cErr {
			t.Errorf("(%d, %d): Go=%v C=%v", tc.numCells, tc.mult, goErr, cErr)
		}
	}
}

// goArcState serializes the Go-side arc state the same way the C
// wrapper does (linkage as element indices, roots as cell ids).
func goArcState(t *testing.T, cells []h3Index, phase int32) cArcState {
	t.Helper()
	var arcset arcSet
	if err := createArcSet(cells, int64(len(cells)), &arcset); err != eSuccess {
		t.Fatalf("createArcSet: %v", err)
	}
	if phase >= 1 {
		if err := cancelArcPairs(arcset); err != eSuccess {
			t.Fatalf("cancelArcPairs: %v", err)
		}
	}
	idx := make(map[*arc]int64, arcset.numArcs)
	for i := range arcset.arcs {
		idx[&arcset.arcs[i]] = int64(i)
	}
	st := cArcState{}
	for i := range arcset.arcs {
		a := &arcset.arcs[i]
		st.ids = append(st.ids, a.id)
		st.removed = append(st.removed, a.isRemoved)
		st.nextIdx = append(st.nextIdx, idx[a.next])
		st.prevIdx = append(st.prevIdx, idx[a.prev])
		st.rootID = append(st.rootID, getRoot(a).id)
	}
	return st
}

func Test_arcState_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		for _, phase := range []int32{0, 1} {
			goSt := goArcState(t, cells, phase)
			cSt, err := arcStateC(cells, int64(len(cells)), phase)
			if err != eSuccess {
				t.Fatalf("arcStateC(%s, phase %d): %v", name, phase, err)
			}
			if len(goSt.ids) != len(cSt.ids) {
				t.Fatalf("%s phase %d: numArcs Go=%d C=%d", name, phase, len(goSt.ids), len(cSt.ids))
			}
			for i := range goSt.ids {
				if goSt.ids[i] != cSt.ids[i] || goSt.removed[i] != cSt.removed[i] ||
					goSt.nextIdx[i] != cSt.nextIdx[i] || goSt.prevIdx[i] != cSt.prevIdx[i] ||
					goSt.rootID[i] != cSt.rootID[i] {
					t.Fatalf("%s phase %d arc %d: Go={%x %v %d %d %x} C={%x %v %d %d %x}",
						name, phase, i,
						uint64(goSt.ids[i]), goSt.removed[i], goSt.nextIdx[i], goSt.prevIdx[i], uint64(goSt.rootID[i]),
						uint64(cSt.ids[i]), cSt.removed[i], cSt.nextIdx[i], cSt.prevIdx[i], uint64(cSt.rootID[i]))
				}
			}
		}
	}
}

func Test_findArc_parity(t *testing.T) {
	cells := multiPolyParitySets(t)["hole"]
	n := int64(len(cells))

	var arcset arcSet
	if err := createArcSet(cells, n, &arcset); err != eSuccess {
		t.Fatalf("createArcSet: %v", err)
	}
	idx := make(map[*arc]int64, arcset.numArcs)
	for i := range arcset.arcs {
		idx[&arcset.arcs[i]] = int64(i)
	}

	// Every edge in the set is found at the same index on both sides.
	for i := range arcset.arcs {
		e := arcset.arcs[i].id
		goIdx := idx[findArc(arcset, e)]
		cIdx := findArcIndexC(cells, n, e)
		if goIdx != cIdx {
			t.Fatalf("findArc(%x): Go idx=%d C idx=%d", uint64(e), goIdx, cIdx)
		}
	}

	// An edge not in the set is not found on either side.
	var outside h3Index
	outsideGeo := LatLng{Lat: Rad(0.9), Lng: Rad(0.9)}
	if err := latLngToCell(&outsideGeo, 9, &outside); err != eSuccess {
		t.Fatalf("latLngToCell: %v", err)
	}
	var outsideEdges [6]h3Index
	if err := originToDirectedEdges(outside, outsideEdges[:]); err != eSuccess {
		t.Fatalf("originToDirectedEdges: %v", err)
	}
	if got := findArc(arcset, outsideEdges[0]); got != nil {
		t.Errorf("Go findArc(missing) = %v, want nil", got)
	}
	if got := findArcIndexC(cells, n, outsideEdges[0]); got != -1 {
		t.Errorf("C findArc(missing) idx = %d, want -1", got)
	}
}

func Test_countLoops_countPolys_loopSet_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		n := int64(len(cells))

		// countLoops after cancellation.
		var arcset arcSet
		if err := createArcSet(cells, n, &arcset); err != eSuccess {
			t.Fatalf("createArcSet: %v", err)
		}
		if err := cancelArcPairs(arcset); err != eSuccess {
			t.Fatalf("cancelArcPairs: %v", err)
		}
		if got, want := countLoops(arcset), countLoopsAfterCancelC(cells, n); got != want {
			t.Fatalf("countLoops(%s): Go=%d C=%d", name, got, want)
		}

		// Full sorted loop set: roots and vertex counts exact, areas
		// within areaClose, vertices within vec3UlpClose per component.
		var loopset sortableLoopSet
		if err := createSortableLoopSet(arcset, &loopset); err != eSuccess {
			t.Fatalf("createSortableLoopSet: %v", err)
		}
		cSet, err := loopSetC(cells, n)
		if err != eSuccess {
			t.Fatalf("loopSetC(%s): %v", name, err)
		}
		if loopset.numLoops != int64(len(cSet.roots)) {
			t.Fatalf("%s: numLoops Go=%d C=%d", name, loopset.numLoops, len(cSet.roots))
		}
		if got, want := countPolys(loopset), cSet.numPolys; got != want {
			t.Fatalf("countPolys(%s): Go=%d C=%d", name, got, want)
		}
		for i := int64(0); i < loopset.numLoops; i++ {
			gl := loopset.sloops[i]
			if gl.root != cSet.roots[i] {
				t.Fatalf("%s loop %d: root Go=%x C=%x", name, i, uint64(gl.root), uint64(cSet.roots[i]))
			}
			if !areaClose(gl.area, cSet.areas[i]) {
				t.Fatalf("%s loop %d: area Go=%v C=%v", name, i, gl.area, cSet.areas[i])
			}
			if len(gl.loop) != len(cSet.loops[i]) {
				t.Fatalf("%s loop %d: numVerts Go=%d C=%d", name, i, len(gl.loop), len(cSet.loops[i]))
			}
			for j := range gl.loop {
				if !vec3UlpClose(gl.loop[j].Lat.Rad(), cSet.loops[i][j].Lat.Rad()) ||
					!vec3UlpClose(gl.loop[j].Lng.Rad(), cSet.loops[i][j].Lng.Rad()) {
					t.Fatalf("%s loop %d vert %d: Go=%v C=%v", name, i, j, gl.loop[j], cSet.loops[i][j])
				}
			}
		}
	}
}

// triangleKey serializes a 3-vertex loop for order-independent globe
// comparison (the 8 octant triangles tie on area, so the final
// qsort/sort.Slice order of equal keys is implementation-defined on
// both sides; the triangle *contents* are exact literals).
func triangleKey(loop GeoLoop) [6]float64 {
	return [6]float64{
		loop[0].Lat.Rad(), loop[0].Lng.Rad(),
		loop[1].Lat.Rad(), loop[1].Lng.Rad(),
		loop[2].Lat.Rad(), loop[2].Lng.Rad(),
	}
}

func sortedTriangles(t *testing.T, loops []GeoLoop) [][6]float64 {
	t.Helper()
	keys := make([][6]float64, 0, len(loops))
	for _, l := range loops {
		if len(l) != 3 {
			t.Fatalf("expected triangle, got %d verts", len(l))
		}
		keys = append(keys, triangleKey(l))
	}
	sort.Slice(keys, func(a, b int) bool {
		for k := 0; k < 6; k++ {
			if keys[a][k] != keys[b][k] {
				return keys[a][k] < keys[b][k]
			}
		}
		return false
	})
	return keys
}

func Test_createGlobeMultiPolygon_parity(t *testing.T) {
	var goOut geoMultiPolygon
	if err := createGlobeMultiPolygon(&goOut); err != eSuccess {
		t.Fatalf("createGlobeMultiPolygon: %v", err)
	}
	cPolys, cVerts, err := globeMultiPolygonC()
	if err != eSuccess {
		t.Fatalf("globeMultiPolygonC: %v", err)
	}
	if goOut.NumPolygons != 8 || cPolys != 8 || len(cVerts) != 24 {
		t.Fatalf("polygon counts: Go=%d C=%d (C verts %d)", goOut.NumPolygons, cPolys, len(cVerts))
	}
	goLoops := make([]GeoLoop, 0, 8)
	for _, p := range goOut.Polygons {
		if len(p.Holes) != 0 {
			t.Fatal("globe polygons must have no holes")
		}
		goLoops = append(goLoops, p.GeoLoop)
	}
	cLoops := make([]GeoLoop, 0, 8)
	for i := 0; i < 24; i += 3 {
		cLoops = append(cLoops, GeoLoop(cVerts[i:i+3]))
	}
	goKeys := sortedTriangles(t, goLoops)
	cKeys := sortedTriangles(t, cLoops)
	for i := range goKeys {
		if goKeys[i] != cKeys[i] {
			// The octant vertices are exact literal constants on both
			// sides, so this comparison is bit-exact.
			t.Fatalf("octant %d: Go=%v C=%v", i, goKeys[i], cKeys[i])
		}
	}
}

// assertMultiPolyParity compares a Go and C geoMultiPolygon: counts
// exact; vertices within vec3UlpClose. When orderNormalize is true (the
// globe case: all outer areas tie, so the polygon order is
// implementation-defined), the polygon records are order-normalized by
// sortLinkedPolyRecords independently on each side — each polygon's
// internal outer/hole loop order and vertex sequence untouched — and
// then the complete structure is compared: outer/hole counts, every
// loop's vertex count, and every vertex.
func assertMultiPolyParity(t *testing.T, name string, goOut, cOut geoMultiPolygon, orderNormalize bool) {
	t.Helper()
	if goOut.NumPolygons != cOut.NumPolygons {
		t.Fatalf("%s: NumPolygons Go=%d C=%d", name, goOut.NumPolygons, cOut.NumPolygons)
	}
	if orderNormalize {
		goRecs, cRecs := multiPolyRecords(goOut), multiPolyRecords(cOut)
		sortLinkedPolyRecords(goRecs)
		sortLinkedPolyRecords(cRecs)
		assertPolyRecordsEqual(t, name, goRecs, cRecs)
		return
	}
	for i := int32(0); i < goOut.NumPolygons; i++ {
		gp, cp := goOut.Polygons[i], cOut.Polygons[i]
		if len(gp.GeoLoop) != len(cp.GeoLoop) || len(gp.Holes) != len(cp.Holes) {
			t.Fatalf("%s poly %d: shape Go=(%d verts, %d holes) C=(%d verts, %d holes)",
				name, i, len(gp.GeoLoop), len(gp.Holes), len(cp.GeoLoop), len(cp.Holes))
		}
		assertLoopUlpClose(t, name, i, -1, gp.GeoLoop, cp.GeoLoop)
		for h := range gp.Holes {
			if len(gp.Holes[h]) != len(cp.Holes[h]) {
				t.Fatalf("%s poly %d hole %d: verts Go=%d C=%d", name, i, h, len(gp.Holes[h]), len(cp.Holes[h]))
			}
			assertLoopUlpClose(t, name, i, h, gp.Holes[h], cp.Holes[h])
		}
	}
}

func assertLoopUlpClose(t *testing.T, name string, poly int32, hole int, a, b GeoLoop) {
	t.Helper()
	for j := range a {
		if !vec3UlpClose(a[j].Lat.Rad(), b[j].Lat.Rad()) ||
			!vec3UlpClose(a[j].Lng.Rad(), b[j].Lng.Rad()) {
			t.Fatalf("%s poly %d hole %d vert %d: Go=%v C=%v", name, poly, hole, j, a[j], b[j])
		}
	}
}

func Test_cellsToMultiPolygon_parity(t *testing.T) {
	for name, cells := range multiPolyParitySets(t) {
		n := int64(len(cells))
		var goOut geoMultiPolygon
		goErr := cellsToMultiPolygon(cells, n, &goOut)
		cOut, cErr := cellsToMultiPolygonC(cells, n)
		if goErr != cErr {
			t.Fatalf("%s: Go err=%v C err=%v", name, goErr, cErr)
		}
		assertMultiPolyParity(t, name, goOut, cOut, name == "globeAllRes0")
	}

	// Error paths — exact.
	errCases := map[string]struct {
		cells []h3Index
		n     int64
	}{
		"negative":    {nil, -1},
		"duplicates":  {[]h3Index{0x81efbffffffffff, 0x81efbffffffffff}, 2},
		"resMismatch": {[]h3Index{0x8027fffffffffff, 0x81efbffffffffff}, 2},
		"invalid":     {[]h3Index{0x8027fffffffffff, 0x81efbffffffffff + 1}, 2},
		"overflow":    {nil, int64(^uint64(0)>>1) / 100},
	}
	for name, tc := range errCases {
		var goOut geoMultiPolygon
		goErr := cellsToMultiPolygon(tc.cells, tc.n, &goOut)
		_, cErr := cellsToMultiPolygonC(tc.cells, tc.n)
		if goErr != cErr {
			t.Errorf("%s: Go err=%v C err=%v", name, goErr, cErr)
		}
	}
}

// linkedPolyRecord is one polygon node of a serialized linked chain:
// per-loop coordinate counts plus the flattened vertices, preserving
// the node's internal loop order and vertex sequence.
type linkedPolyRecord struct {
	coordsPerLoop []int32
	verts         []LatLng
}

// splitLinkedShape splits a flattened cLinkedShape into per-polygon
// records.
func splitLinkedShape(t *testing.T, shape cLinkedShape) []linkedPolyRecord {
	t.Helper()
	recs := make([]linkedPolyRecord, 0, len(shape.loopsPerPoly))
	loopIdx, vertIdx := 0, 0
	for _, numLoops := range shape.loopsPerPoly {
		var rec linkedPolyRecord
		for l := int32(0); l < numLoops; l++ {
			n := int(shape.coordsPerLoop[loopIdx])
			loopIdx++
			rec.coordsPerLoop = append(rec.coordsPerLoop, int32(n))
			rec.verts = append(rec.verts, shape.verts[vertIdx:vertIdx+n]...)
			vertIdx += n
		}
		recs = append(recs, rec)
	}
	if loopIdx != len(shape.coordsPerLoop) || vertIdx != len(shape.verts) {
		t.Fatalf("shape not fully consumed: %d/%d loops, %d/%d verts",
			loopIdx, len(shape.coordsPerLoop), vertIdx, len(shape.verts))
	}
	return recs
}

// sortLinkedPolyRecords orders polygon records by a canonical key —
// loop count, then per-loop coordinate counts, then raw vertex values
// — without touching any record's internal loop/vertex sequence. Used
// only where the chain order of the polygon nodes themselves is
// implementation-defined (the globe tiling's tied octant areas); the
// octant vertices are exact literal constants on both sides, so
// exact-value sorting yields the same canonical order on each side.
func sortLinkedPolyRecords(recs []linkedPolyRecord) {
	sort.Slice(recs, func(x, y int) bool {
		a, b := recs[x], recs[y]
		if len(a.coordsPerLoop) != len(b.coordsPerLoop) {
			return len(a.coordsPerLoop) < len(b.coordsPerLoop)
		}
		for i := range a.coordsPerLoop {
			if a.coordsPerLoop[i] != b.coordsPerLoop[i] {
				return a.coordsPerLoop[i] < b.coordsPerLoop[i]
			}
		}
		for i := range a.verts {
			if a.verts[i].Lat.Rad() != b.verts[i].Lat.Rad() {
				return a.verts[i].Lat.Rad() < b.verts[i].Lat.Rad()
			}
			if a.verts[i].Lng.Rad() != b.verts[i].Lng.Rad() {
				return a.verts[i].Lng.Rad() < b.verts[i].Lng.Rad()
			}
		}
		return false
	})
}

// multiPolyRecords converts a geoMultiPolygon into per-polygon records
// (outer loop first, then holes, matching the linked-chain layout),
// preserving each polygon's internal loop order and vertex sequence.
func multiPolyRecords(m geoMultiPolygon) []linkedPolyRecord {
	recs := make([]linkedPolyRecord, 0, m.NumPolygons)
	for _, p := range m.Polygons {
		rec := linkedPolyRecord{coordsPerLoop: []int32{int32(len(p.GeoLoop))}}
		rec.verts = append(rec.verts, p.GeoLoop...)
		for _, h := range p.Holes {
			rec.coordsPerLoop = append(rec.coordsPerLoop, int32(len(h)))
			rec.verts = append(rec.verts, h...)
		}
		recs = append(recs, rec)
	}
	return recs
}

// assertPolyRecordsEqual compares two per-polygon record slices in
// full: loop counts and per-loop coordinate counts exact, every vertex
// within vec3UlpClose.
func assertPolyRecordsEqual(t *testing.T, name string, goRecs, cRecs []linkedPolyRecord) {
	t.Helper()
	if len(goRecs) != len(cRecs) {
		t.Fatalf("%s: NumPolygons Go=%d C=%d", name, len(goRecs), len(cRecs))
	}
	for p := range goRecs {
		gr, cr := goRecs[p], cRecs[p]
		if len(gr.coordsPerLoop) != len(cr.coordsPerLoop) {
			t.Fatalf("%s poly %d: loops Go=%d C=%d", name, p, len(gr.coordsPerLoop), len(cr.coordsPerLoop))
		}
		for l := range gr.coordsPerLoop {
			if gr.coordsPerLoop[l] != cr.coordsPerLoop[l] {
				t.Fatalf("%s poly %d loop %d: coords Go=%d C=%d", name, p, l, gr.coordsPerLoop[l], cr.coordsPerLoop[l])
			}
		}
		if len(gr.verts) != len(cr.verts) {
			t.Fatalf("%s poly %d: verts Go=%d C=%d", name, p, len(gr.verts), len(cr.verts))
		}
		for i := range gr.verts {
			if !vec3UlpClose(gr.verts[i].Lat.Rad(), cr.verts[i].Lat.Rad()) ||
				!vec3UlpClose(gr.verts[i].Lng.Rad(), cr.verts[i].Lng.Rad()) {
				t.Fatalf("%s poly %d vert %d: Go=%v C=%v", name, p, i, gr.verts[i], cr.verts[i])
			}
		}
	}
}

func Test_cellsToLinkedMultiPolygon_450_parity(t *testing.T) {
	// Isolated cellsToLinkedMultiPolygon: the C side calls ONLY the
	// public function and serializes its linked output directly —
	// polygon/loop/coordinate counts, every vertex, and the
	// First/Last + tail-Next linkage invariants — with no conversion
	// pipeline in between; the Go side serializes its own linked
	// output the same way. Vertices carry the boundary pipeline's
	// vec3UlpClose discipline. For the globe tiling the eight octant
	// polygons tie on outer area, so the chain order of the polygon
	// nodes is implementation-defined; the records are order-normalized
	// by sortLinkedPolyRecords independently on each side and then
	// compared in full like every other set.
	for name, cells := range multiPolyParitySets(t) {
		var linked linkedGeoPolygon
		goErr := cellsToLinkedMultiPolygon(cells, int32(len(cells)), &linked)
		if goErr != eSuccess {
			t.Fatalf("%s: cellsToLinkedMultiPolygon: %v", name, goErr)
		}
		var goRecs []linkedPolyRecord
		goInv := true
		for poly := &linked; poly != nil; poly = poly.Next {
			st := goLinkedPolyState(poly)
			goInv = goInv && st.invariantsOK
			goRecs = append(goRecs, linkedPolyRecord{coordsPerLoop: st.coordsPerLoop, verts: st.verts})
		}
		maxLoops := 4 * len(cells)
		if maxLoops < 64 {
			maxLoops = 64
		}
		cShape, cInv, cErr := cellsToLinkedMultiPolygonSerializedC(cells, int32(len(cells)), maxLoops, 12*len(cells)+64)
		if cErr != eSuccess {
			t.Fatalf("%s: C cellsToLinkedMultiPolygon: %v", name, cErr)
		}
		if !goInv || !cInv {
			t.Fatalf("%s: linkage invariants Go=%v C=%v", name, goInv, cInv)
		}
		cRecs := splitLinkedShape(t, cShape)
		if name == "globeAllRes0" {
			sortLinkedPolyRecords(goRecs)
			sortLinkedPolyRecords(cRecs)
		}
		assertPolyRecordsEqual(t, name, goRecs, cRecs)
	}

	// The 4.5.0 behavioral change: invalid cells now fail with
	// E_CELL_INVALID (was E_FAILED at 4.4.0) — exact on both sides.
	invalid := []h3Index{0xd60006d60000f100, 0x3c3c403c1300d668}
	var linked linkedGeoPolygon
	goErr := cellsToLinkedMultiPolygon(invalid, 2, &linked)
	_, _, cErr := cellsToLinkedMultiPolygonSerializedC(invalid, 2, 8, 8)
	if goErr != cErr || goErr != eCellInvalid {
		t.Errorf("invalid: Go=%v C=%v (want eCellInvalid from both)", goErr, cErr)
	}
}
