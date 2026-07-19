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

func Test_addLinkedGeoLoop_direct_parity(t *testing.T) {
	tri := convertParityLoops(t)[1]
	// times=2 covers both the first-loop and the append branch.
	for _, times := range []int32{1, 2} {
		var poly linkedGeoPolygon
		goErr := eSuccess
		for i := int32(0); i < times && goErr == eSuccess; i++ {
			goErr = addLinkedGeoLoop(tri, &poly)
		}
		var goCoords []int32
		for loop := poly.First; loop != nil; loop = loop.Next {
			goCoords = append(goCoords, countLinkedCoords(loop))
		}
		cLoops, cCoords, cErr := addLinkedGeoLoopC(tri, times)
		if goErr != cErr || int32(len(goCoords)) != cLoops {
			t.Fatalf("times %d: Go=(%d loops, %v) C=(%d loops, %v)", times, len(goCoords), goErr, cLoops, cErr)
		}
		for i := range goCoords {
			if goCoords[i] != cCoords[i] {
				t.Fatalf("times %d loop %d: coords Go=%d C=%d", times, i, goCoords[i], cCoords[i])
			}
		}
	}

	// < 3 verts propagates eFailed from geoLoopToLinkedGeoLoop.
	short := GeoLoop{{Lat: Rad(0.1)}, {Lat: Rad(0.2)}}
	var poly linkedGeoPolygon
	goErr := addLinkedGeoLoop(short, &poly)
	_, _, cErr := addLinkedGeoLoopC(short, 1)
	if goErr != cErr || goErr != eFailed {
		t.Errorf("short loop: Go=%v C=%v, want eFailed from both", goErr, cErr)
	}
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
		var goCoords []int32
		for loop := goPoly.First; loop != nil; loop = loop.Next {
			goCoords = append(goCoords, countLinkedCoords(loop))
		}
		cLoops, cCoords, cErr := geoPolygonToLinkedGeoLoopsC(in)
		if goErr != cErr || int32(len(goCoords)) != cLoops {
			t.Fatalf("input %d: Go=(%d loops, %v) C=(%d loops, %v)", i, len(goCoords), goErr, cLoops, cErr)
		}
		for l := range goCoords {
			if goCoords[l] != cCoords[l] {
				t.Fatalf("input %d loop %d: coords Go=%d C=%d", i, l, goCoords[l], cCoords[l])
			}
		}
	}
}

func Test_geoMultiPolygonToLinked_direct_parity(t *testing.T) {
	// Isolated success path of geoMultiPolygonToLinkedGeoPolygon: both
	// sides convert cellsToMultiPolygon output (independently
	// parity-verified) and compare the linked chain's shape and
	// vertices. Vertices are copied verbatim by the conversion, so
	// they carry the pipeline's vec3UlpClose discipline.
	sets := map[string][]h3Index{
		"singleHex": {0x890dab6220bffff},
		"hole": {
			0x892830828c7ffff, 0x892830828d7ffff, 0x8928308289bffff,
			0x89283082813ffff, 0x8928308288fffff, 0x89283082883ffff},
		"nonContiguous2": {0x8928308291bffff, 0x89283082943ffff},
	}
	for name, cells := range sets {
		n := int64(len(cells))
		var mpoly geoMultiPolygon
		if err := cellsToMultiPolygon(cells, n, &mpoly); err != eSuccess {
			t.Fatalf("cellsToMultiPolygon(%s): %v", name, err)
		}
		var linked linkedGeoPolygon
		if err := geoMultiPolygonToLinkedGeoPolygon(&mpoly, &linked); err != eSuccess {
			t.Fatalf("geoMultiPolygonToLinkedGeoPolygon(%s): %v", name, err)
		}
		var goShape cLinkedShape
		for poly := &linked; poly != nil; poly = poly.Next {
			loops := int32(0)
			for loop := poly.First; loop != nil; loop = loop.Next {
				coords := int32(0)
				for c := loop.First; c != nil; c = c.Next {
					goShape.verts = append(goShape.verts, c.Vertex)
					coords++
				}
				goShape.coordsPerLoop = append(goShape.coordsPerLoop, coords)
				loops++
			}
			goShape.loopsPerPoly = append(goShape.loopsPerPoly, loops)
		}
		cShape, err := geoMultiPolygonToLinkedC(cells, n, 64, 4096)
		if err != eSuccess {
			t.Fatalf("geoMultiPolygonToLinkedC(%s): %v", name, err)
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
			if !vec3UlpClose(goShape.verts[i].Lat.Rad(), cShape.verts[i].Lat.Rad()) ||
				!vec3UlpClose(goShape.verts[i].Lng.Rad(), cShape.verts[i].Lng.Rad()) {
				t.Fatalf("%s vert %d: Go=%v C=%v", name, i, goShape.verts[i], cShape.verts[i])
			}
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
