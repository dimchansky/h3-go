//go:build cgo && c2go && h3v450

package h3

import "testing"

// Direct parity for the H3 4.5.0 linkedGeo.c conversion helpers: the
// file-statics (linkedGeoLoopToGeoLoop, geoLoopToLinkedGeoLoop,
// linkedGeoPolygonToGeoPolygon, addLinkedGeoLoop,
// geoPolygonToLinkedGeoLoops) via same-TU wrappers in
// h3lib_linkedGeo_c2go.c, and the error branches of the two conversion
// externs (linkedGeoPolygonToGeoMultiPolygon,
// geoMultiPolygonToLinkedGeoPolygon). Vertex values are copied
// verbatim by both implementations, so all comparisons are exact.

func convertParityLoops(t *testing.T) []GeoLoop {
	t.Helper()
	// A real cell boundary plus a synthetic triangle.
	var gb CellBoundary
	if err := cellToBoundary(0x890dab6220bffff, &gb); err != eSuccess {
		t.Fatalf("cellToBoundary: %v", err)
	}
	real := make(GeoLoop, gb.numVerts)
	copy(real, gb.verts[:gb.numVerts])
	tri := GeoLoop{{Lat: Rad(0.1), Lng: Rad(0.2)}, {Lat: Rad(0.3), Lng: Rad(0.1)}, {Lat: Rad(0.2), Lng: Rad(0.4)}}
	return []GeoLoop{real, tri}
}

func loopsEqual(a, b GeoLoop) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Test_linkedGeoLoopToGeoLoop_direct_parity(t *testing.T) {
	for i, verts := range convertParityLoops(t) {
		// Go side: build the linked loop the same way the C wrapper does.
		var poly linkedGeoPolygon
		loop := addNewLinkedLoop(&poly)
		for j := range verts {
			addLinkedCoord(loop, &verts[j])
		}
		var goOut GeoLoop
		goErr := linkedGeoLoopToGeoLoop(loop, &goOut)
		cOut, cErr := linkedGeoLoopToGeoLoopC(verts)
		if goErr != cErr || !loopsEqual(goOut, cOut) {
			t.Fatalf("loop %d: Go=(%d verts, %v) C=(%d verts, %v)", i, len(goOut), goErr, len(cOut), cErr)
		}
	}

	// < 3 verts fails with eFailed on both sides.
	short := GeoLoop{{Lat: Rad(0.1)}, {Lat: Rad(0.2)}}
	var poly linkedGeoPolygon
	loop := addNewLinkedLoop(&poly)
	for j := range short {
		addLinkedCoord(loop, &short[j])
	}
	var goOut GeoLoop
	goErr := linkedGeoLoopToGeoLoop(loop, &goOut)
	_, cErr := linkedGeoLoopToGeoLoopC(short)
	if goErr != cErr || goErr != eFailed {
		t.Errorf("short loop: Go=%v C=%v, want eFailed from both", goErr, cErr)
	}
}

func Test_geoLoopToLinkedGeoLoop_direct_parity(t *testing.T) {
	for i, verts := range convertParityLoops(t) {
		var goLoop linkedGeoLoop
		goErr := geoLoopToLinkedGeoLoop(verts, &goLoop)
		var goOut GeoLoop
		for c := goLoop.First; c != nil; c = c.Next {
			goOut = append(goOut, c.Vertex)
		}
		cOut, cErr := geoLoopToLinkedGeoLoopC(verts)
		if goErr != cErr || !loopsEqual(goOut, cOut) {
			t.Fatalf("loop %d: Go=(%d verts, %v) C=(%d verts, %v)", i, len(goOut), goErr, len(cOut), cErr)
		}
	}

	short := GeoLoop{{Lat: Rad(0.1)}, {Lat: Rad(0.2)}}
	var goLoop linkedGeoLoop
	goErr := geoLoopToLinkedGeoLoop(short, &goLoop)
	_, cErr := geoLoopToLinkedGeoLoopC(short)
	if goErr != cErr || goErr != eFailed {
		t.Errorf("short loop: Go=%v C=%v, want eFailed from both", goErr, cErr)
	}
}

func Test_linkedGeoPolygonToGeoPolygon_direct_parity(t *testing.T) {
	loops := convertParityLoops(t)
	inputs := [][]GeoLoop{
		{loops[0]},           // outer only
		{loops[0], loops[1]}, // outer + one hole
	}
	for i, in := range inputs {
		// Go side: build the linked polygon like the C wrapper does.
		var poly linkedGeoPolygon
		for _, l := range in {
			loop := addNewLinkedLoop(&poly)
			for j := range l {
				addLinkedCoord(loop, &l[j])
			}
		}
		var goOut GeoPolygon
		goErr := linkedGeoPolygonToGeoPolygon(&poly, &goOut)
		cLoops, cErr := linkedGeoPolygonToGeoPolygonC(in)
		if goErr != cErr {
			t.Fatalf("input %d: Go err=%v C err=%v", i, goErr, cErr)
		}
		goLoops := append([]GeoLoop{goOut.GeoLoop}, goOut.Holes...)
		if len(goLoops) != len(cLoops) {
			t.Fatalf("input %d: loops Go=%d C=%d", i, len(goLoops), len(cLoops))
		}
		for l := range goLoops {
			if !loopsEqual(goLoops[l], cLoops[l]) {
				t.Fatalf("input %d loop %d differs", i, l)
			}
		}
	}
}

// goLinkedPolyState serializes one Go linked polygon node exactly as
// the C serializer does, checking the same First/Last linkage
// invariants.
func goLinkedPolyState(poly *linkedGeoPolygon) cLinkedPolyState {
	st := cLinkedPolyState{invariantsOK: true}
	var lastLoop *linkedGeoLoop
	for loop := poly.First; loop != nil; loop = loop.Next {
		c := int32(0)
		var lastCoord *linkedLatLng
		for coord := loop.First; coord != nil; coord = coord.Next {
			st.verts = append(st.verts, coord.Vertex)
			lastCoord = coord
			c++
		}
		if loop.Last != lastCoord || (lastCoord != nil && lastCoord.Next != nil) {
			st.invariantsOK = false
		}
		st.coordsPerLoop = append(st.coordsPerLoop, c)
		lastLoop = loop
	}
	if poly.Last != lastLoop {
		st.invariantsOK = false
	}
	return st
}

func assertLinkedStateExact(t *testing.T, label string, goSt, cSt cLinkedPolyState) {
	t.Helper()
	if !goSt.invariantsOK || !cSt.invariantsOK {
		t.Fatalf("%s: linkage invariants Go=%v C=%v", label, goSt.invariantsOK, cSt.invariantsOK)
	}
	if len(goSt.coordsPerLoop) != len(cSt.coordsPerLoop) || len(goSt.verts) != len(cSt.verts) {
		t.Fatalf("%s: shape Go=(%d loops, %d verts) C=(%d loops, %d verts)",
			label, len(goSt.coordsPerLoop), len(goSt.verts), len(cSt.coordsPerLoop), len(cSt.verts))
	}
	for i := range goSt.coordsPerLoop {
		if goSt.coordsPerLoop[i] != cSt.coordsPerLoop[i] {
			t.Fatalf("%s loop %d: coords Go=%d C=%d", label, i, goSt.coordsPerLoop[i], cSt.coordsPerLoop[i])
		}
	}
	for i := range goSt.verts {
		if goSt.verts[i] != cSt.verts[i] {
			t.Fatalf("%s vert %d: Go=%v C=%v", label, i, goSt.verts[i], cSt.verts[i])
		}
	}
}

func Test_addLinkedGeoLoop_direct_parity(t *testing.T) {
	tri := convertParityLoops(t)[1]
	// times=2 covers both the first-loop and the append branch; the
	// full state — every vertex plus the First/Last linkage
	// invariants — compares exactly.
	for _, times := range []int32{1, 2} {
		var poly linkedGeoPolygon
		goErr := eSuccess
		for i := int32(0); i < times && goErr == eSuccess; i++ {
			goErr = addLinkedGeoLoop(tri, &poly)
		}
		goSt := goLinkedPolyState(&poly)
		cSt, cErr := addLinkedGeoLoopC(tri, times)
		if goErr != cErr {
			t.Fatalf("times %d: Go err=%v C err=%v", times, goErr, cErr)
		}
		assertLinkedStateExact(t, "addLinkedGeoLoop", goSt, cSt)
	}

	// < 3 verts fails, leaving identical partial state on both sides
	// (the empty loop node was appended before the vertex check).
	short := GeoLoop{{Lat: Rad(0.1)}, {Lat: Rad(0.2)}}
	var poly linkedGeoPolygon
	goErr := addLinkedGeoLoop(short, &poly)
	goSt := goLinkedPolyState(&poly)
	cSt, cErr := addLinkedGeoLoopC(short, 1)
	if goErr != cErr || goErr != eFailed {
		t.Fatalf("short loop: Go=%v C=%v, want eFailed from both", goErr, cErr)
	}
	if len(goSt.coordsPerLoop) != 1 || goSt.coordsPerLoop[0] != 0 {
		t.Fatalf("short loop partial state: Go %v, want one empty loop node", goSt.coordsPerLoop)
	}
	assertLinkedStateExact(t, "addLinkedGeoLoop partial", goSt, cSt)
}

func Test_geoPolygonToLinkedGeoLoops_direct_parity(t *testing.T) {
	loops := convertParityLoops(t)
	inputs := [][]GeoLoop{
		{loops[0]},
		{loops[0], loops[1]},
	}
	for i, in := range inputs {
		gp := GeoPolygon{GeoLoop: in[0]}
		if len(in) > 1 {
			gp.Holes = in[1:]
		}
		var goPoly linkedGeoPolygon
		goErr := geoPolygonToLinkedGeoLoops(&gp, &goPoly)
		goSt := goLinkedPolyState(&goPoly)
		cSt, cErr := geoPolygonToLinkedGeoLoopsC(in)
		if goErr != cErr {
			t.Fatalf("input %d: Go err=%v C err=%v", i, goErr, cErr)
		}
		assertLinkedStateExact(t, "geoPolygonToLinkedGeoLoops", goSt, cSt)
	}

	// Valid outer + short hole: fails after the outer loop was fully
	// converted; the partial linked state (outer loop + empty hole
	// node) must match C before the owning extern cleans it.
	short := GeoLoop{{Lat: Rad(0.1)}, {Lat: Rad(0.2)}}
	gp := GeoPolygon{GeoLoop: loops[1], Holes: []GeoLoop{short}}
	var goPoly linkedGeoPolygon
	goErr := geoPolygonToLinkedGeoLoops(&gp, &goPoly)
	goSt := goLinkedPolyState(&goPoly)
	cSt, cErr := geoPolygonToLinkedGeoLoopsC([]GeoLoop{loops[1], short})
	if goErr != cErr || goErr != eFailed {
		t.Fatalf("outer+short hole: Go=%v C=%v, want eFailed from both", goErr, cErr)
	}
	if len(goSt.coordsPerLoop) != 2 || goSt.coordsPerLoop[0] != int32(len(loops[1])) || goSt.coordsPerLoop[1] != 0 {
		t.Fatalf("partial state: Go %v, want [outer, empty hole node]", goSt.coordsPerLoop)
	}
	assertLinkedStateExact(t, "geoPolygonToLinkedGeoLoops partial", goSt, cSt)
}

func Test_geoMultiPolygonToLinked_direct_parity(t *testing.T) {
	// Isolated success path of geoMultiPolygonToLinkedGeoPolygon on a
	// synthetic, identical GeoMultiPolygon fed to both sides: the
	// linked chain's shape, every vertex (bit-exact — verbatim
	// copies), and the linkage invariants compare.
	tri := GeoLoop{{Lat: Rad(0.1), Lng: Rad(0.2)}, {Lat: Rad(0.3), Lng: Rad(0.1)}, {Lat: Rad(0.2), Lng: Rad(0.4)}}
	quad := GeoLoop{{Lat: Rad(0.5)}, {Lat: Rad(0.6), Lng: Rad(0.1)}, {Lat: Rad(0.55), Lng: Rad(0.2)}, {Lat: Rad(0.45), Lng: Rad(0.1)}}
	penta := GeoLoop{{Lat: Rad(-0.1)}, {Lat: Rad(-0.2), Lng: Rad(0.1)}, {Lat: Rad(-0.3)}, {Lat: Rad(-0.2), Lng: Rad(-0.1)}, {Lat: Rad(-0.15), Lng: Rad(-0.05)}}
	inputs := map[string]geoMultiPolygon{
		"onePolyNoHoles": {NumPolygons: 1, Polygons: []GeoPolygon{{GeoLoop: tri}}},
		"onePolyOneHole": {NumPolygons: 1, Polygons: []GeoPolygon{{GeoLoop: quad, Holes: []GeoLoop{tri}}}},
		"twoPolys":       {NumPolygons: 2, Polygons: []GeoPolygon{{GeoLoop: quad, Holes: []GeoLoop{tri}}, {GeoLoop: penta}}},
	}
	for name, mp := range inputs {
		var linked linkedGeoPolygon
		if err := geoMultiPolygonToLinkedGeoPolygon(&mp, &linked); err != eSuccess {
			t.Fatalf("geoMultiPolygonToLinkedGeoPolygon(%s): %v", name, err)
		}
		var goShape cLinkedShape
		goInv := true
		for poly := &linked; poly != nil; poly = poly.Next {
			st := goLinkedPolyState(poly)
			goInv = goInv && st.invariantsOK
			goShape.loopsPerPoly = append(goShape.loopsPerPoly, int32(len(st.coordsPerLoop)))
			goShape.coordsPerLoop = append(goShape.coordsPerLoop, st.coordsPerLoop...)
			goShape.verts = append(goShape.verts, st.verts...)
		}
		cShape, cInv, err := geoMultiPolygonToLinkedSyntheticC(mp)
		if err != eSuccess {
			t.Fatalf("geoMultiPolygonToLinkedSyntheticC(%s): %v", name, err)
		}
		if !goInv || !cInv {
			t.Fatalf("%s: linkage invariants Go=%v C=%v", name, goInv, cInv)
		}
		if len(goShape.loopsPerPoly) != len(cShape.loopsPerPoly) ||
			len(goShape.coordsPerLoop) != len(cShape.coordsPerLoop) ||
			len(goShape.verts) != len(cShape.verts) {
			t.Fatalf("%s: shape Go=(%d polys, %d loops, %d verts) C=(%d polys, %d loops, %d verts)",
				name, len(goShape.loopsPerPoly), len(goShape.coordsPerLoop), len(goShape.verts),
				len(cShape.loopsPerPoly), len(cShape.coordsPerLoop), len(cShape.verts))
		}
		for i := range goShape.loopsPerPoly {
			if goShape.loopsPerPoly[i] != cShape.loopsPerPoly[i] {
				t.Fatalf("%s poly %d: loops Go=%d C=%d", name, i, goShape.loopsPerPoly[i], cShape.loopsPerPoly[i])
			}
		}
		for i := range goShape.coordsPerLoop {
			if goShape.coordsPerLoop[i] != cShape.coordsPerLoop[i] {
				t.Fatalf("%s loop %d: coords Go=%d C=%d", name, i, goShape.coordsPerLoop[i], cShape.coordsPerLoop[i])
			}
		}
		for i := range goShape.verts {
			if goShape.verts[i] != cShape.verts[i] {
				t.Fatalf("%s vert %d: Go=%v C=%v", name, i, goShape.verts[i], cShape.verts[i])
			}
		}
	}
}

func Test_linkedConvertCleanup_parity(t *testing.T) {
	// Valid-outer + invalid-hole (and invalid-second-element) cases for
	// both conversion directions, comparing the promised cleanup state
	// as well as the error code. Constructions mirror
	// h3goTest_linkedConvertCleanup exactly.
	tri := GeoLoop{{Lat: Rad(0.1), Lng: Rad(0.2)}, {Lat: Rad(0.3), Lng: Rad(0.1)}, {Lat: Rad(0.2), Lng: Rad(0.4)}}
	duo := GeoLoop{{}, {Lat: Rad(1)}}
	goErrs := make([]h3Error, 4)
	goClean := make([]bool, 4)

	{ // 0: linkedGeoPolygonToGeoPolygon leaves out zeroed after a hole fails.
		var poly linkedGeoPolygon
		outer := addNewLinkedLoop(&poly)
		for i := range tri {
			addLinkedCoord(outer, &tri[i])
		}
		hole := addNewLinkedLoop(&poly)
		for i := range duo {
			addLinkedCoord(hole, &duo[i])
		}
		var out GeoPolygon
		goErrs[0] = linkedGeoPolygonToGeoPolygon(&poly, &out)
		goClean[0] = out.GeoLoop == nil && out.Holes == nil
	}
	{ // 1: linkedGeoPolygonToGeoMultiPolygon cleans the partial multipolygon.
		var poly linkedGeoPolygon
		outer := addNewLinkedLoop(&poly)
		for i := range tri {
			addLinkedCoord(outer, &tri[i])
		}
		second := addNewLinkedPolygon(&poly)
		bad := addNewLinkedLoop(second)
		for i := range duo {
			addLinkedCoord(bad, &duo[i])
		}
		var mpoly geoMultiPolygon
		goErrs[1] = linkedGeoPolygonToGeoMultiPolygon(&poly, &mpoly)
		goClean[1] = mpoly.Polygons == nil && mpoly.NumPolygons == 0
	}
	{ // 2: geoMultiPolygonToLinkedGeoPolygon zeroes the head (short hole).
		mp := geoMultiPolygon{NumPolygons: 1, Polygons: []GeoPolygon{{GeoLoop: tri, Holes: []GeoLoop{duo}}}}
		var out linkedGeoPolygon
		goErrs[2] = geoMultiPolygonToLinkedGeoPolygon(&mp, &out)
		goClean[2] = out.First == nil && out.Last == nil && out.Next == nil
	}
	{ // 3: geoMultiPolygonToLinkedGeoPolygon zeroes the head (short second poly).
		mp := geoMultiPolygon{NumPolygons: 2, Polygons: []GeoPolygon{{GeoLoop: tri}, {GeoLoop: duo}}}
		var out linkedGeoPolygon
		goErrs[3] = geoMultiPolygonToLinkedGeoPolygon(&mp, &out)
		goClean[3] = out.First == nil && out.Last == nil && out.Next == nil
	}

	for which := int32(0); which < 4; which++ {
		cErr, cClean := linkedConvertCleanupC(which)
		if goErrs[which] != cErr || goErrs[which] != eFailed {
			t.Errorf("case %d: err Go=%v C=%v, want eFailed", which, goErrs[which], cErr)
		}
		if !goClean[which] || !cClean {
			t.Errorf("case %d: cleanup state Go=%v C=%v, want clean on both sides", which, goClean[which], cClean)
		}
	}
}

func Test_linkedConvertErrorBranches_parity(t *testing.T) {
	// The upstream testLinkedGeoConvert error constructions, executed
	// identically on both sides (see h3goTest_linkedConvertError).
	v1 := LatLng{}
	v2 := LatLng{Lat: Rad(1)}
	goErrs := make([]h3Error, 5)

	{ // 0: linked loop with 2 verts -> linked->flat eFailed
		var poly linkedGeoPolygon
		loop := addNewLinkedLoop(&poly)
		addLinkedCoord(loop, &v1)
		addLinkedCoord(loop, &v2)
		var mpoly geoMultiPolygon
		goErrs[0] = linkedGeoPolygonToGeoMultiPolygon(&poly, &mpoly)
	}
	{ // 1: empty polygon node -> linked->flat eFailed
		var poly linkedGeoPolygon
		addNewLinkedPolygon(&poly)
		var mpoly geoMultiPolygon
		goErrs[1] = linkedGeoPolygonToGeoMultiPolygon(&poly, &mpoly)
	}
	{ // 2: geoloop with 2 verts -> flat->linked eFailed
		mpoly := geoMultiPolygon{NumPolygons: 1, Polygons: []GeoPolygon{{GeoLoop: GeoLoop{v1, v2}}}}
		var out linkedGeoPolygon
		goErrs[2] = geoMultiPolygonToLinkedGeoPolygon(&mpoly, &out)
	}
	{ // 3: empty chain -> eSuccess
		var empty linkedGeoPolygon
		var mpoly geoMultiPolygon
		goErrs[3] = linkedGeoPolygonToGeoMultiPolygon(&empty, &mpoly)
	}
	{ // 4: empty mpoly -> eSuccess
		var mpoly geoMultiPolygon
		var out linkedGeoPolygon
		goErrs[4] = geoMultiPolygonToLinkedGeoPolygon(&mpoly, &out)
	}

	want := []h3Error{eFailed, eFailed, eFailed, eSuccess, eSuccess}
	for which := int32(0); which < 5; which++ {
		cErr := linkedConvertErrorC(which)
		if goErrs[which] != cErr || goErrs[which] != want[which] {
			t.Errorf("branch %d: Go=%v C=%v, want %v", which, goErrs[which], cErr, want[which])
		}
	}
}
