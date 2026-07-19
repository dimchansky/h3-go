//go:build cgo && c2go

package h3

// Bridges to the H3 4.5.0 linkedGeo.c conversion file-statics via the
// same-TU h3goTest_* wrappers in h3lib_linkedGeo_c2go.c (guarded to the
// 4.5.0 tree), plus the error branches of the two conversion externs.

/*
#include <stdint.h>
#include "h3api.h"

H3Error h3goTest_linkedGeoLoopToGeoLoop(const LatLng *verts, int n,
                                        LatLng *outVerts, int *outN);
H3Error h3goTest_geoLoopToLinkedGeoLoop(const LatLng *verts, int n,
                                        LatLng *outVerts, int *outN);
H3Error h3goTest_linkedGeoPolygonToGeoPolygon(
    const LatLng *verts, const int *loopNumVerts, int numLoops,
    int *outLoopNumVerts, LatLng *outVerts, int *outNumLoops);
H3Error h3goTest_addLinkedGeoLoop(const LatLng *verts, int n, int times,
                                  int *outNumLoops, int *outCoordsPerLoop,
                                  LatLng *outVerts, int *invariantsOut);
H3Error h3goTest_geoPolygonToLinkedGeoLoops(const LatLng *verts,
                                            const int *loopNumVerts,
                                            int numLoops, int *outNumLoops,
                                            int *outCoordsPerLoop,
                                            LatLng *outVerts,
                                            int *invariantsOut);
H3Error h3goTest_geoMultiPolygonToLinkedSynthetic(
    const int *polyNumVerts, const int *polyNumHoles, const int *holeNumVerts,
    LatLng *verts, int numPolys, int *numPolysOut, int *loopsPerPoly,
    int *coordsPerLoop, LatLng *outVerts, int *invariantsOut);
H3Error h3goTest_linkedConvertCleanup(int which, int *cleanOut);
H3Error h3goTest_linkedConvertError(int which);
H3Error h3goTest_cellsToLinkedMultiPolygonSerialized(
    const H3Index *cells, int numCells, int *numPolysOut, int *loopsPerPoly,
    int *coordsPerLoop, LatLng *outVerts, int *invariantsOut);
H3Error h3goTest_linkedToGeoMultiPolygonSynthetic(
    const int *loopsPerPoly, const int *coordsPerLoop, LatLng *verts,
    int numPolys, int64_t *numPolysOut, int64_t *polyNumVerts,
    int64_t *polyNumHoles, int64_t *holeNumVerts, LatLng *outVerts);
*/
import "C"
import "unsafe"

func latLngPtr(verts []LatLng) *C.LatLng {
	if len(verts) == 0 {
		return nil
	}
	return (*C.LatLng)(unsafe.Pointer(&verts[0]))
}

// linkedGeoLoopToGeoLoopC builds a C linked loop from verts and runs
// the C static, returning the resulting flat loop.
func linkedGeoLoopToGeoLoopC(verts []LatLng) (GeoLoop, h3Error) {
	out := make([]LatLng, len(verts)+1)
	var n C.int
	err := h3Error(C.h3goTest_linkedGeoLoopToGeoLoop(latLngPtr(verts),
		C.int(len(verts)), latLngPtr(out), &n))
	if err != eSuccess {
		return nil, err
	}
	return GeoLoop(out[:n]), eSuccess
}

// geoLoopToLinkedGeoLoopC runs the C static and returns the linked
// loop's coordinates in order.
func geoLoopToLinkedGeoLoopC(verts []LatLng) (GeoLoop, h3Error) {
	out := make([]LatLng, len(verts)+1)
	var n C.int
	err := h3Error(C.h3goTest_geoLoopToLinkedGeoLoop(latLngPtr(verts),
		C.int(len(verts)), latLngPtr(out), &n))
	if err != eSuccess {
		return nil, err
	}
	return GeoLoop(out[:n]), eSuccess
}

// linkedGeoPolygonToGeoPolygonC builds a C linked polygon from the
// given loops (loops[0] is the outer loop) and runs the C static.
func linkedGeoPolygonToGeoPolygonC(loops []GeoLoop) ([]GeoLoop, h3Error) {
	var flat []LatLng
	nums := make([]C.int, len(loops))
	total := 0
	for i, l := range loops {
		flat = append(flat, l...)
		nums[i] = C.int(len(l))
		total += len(l)
	}
	outNums := make([]C.int, len(loops)+1)
	outVerts := make([]LatLng, total+1)
	var outLoops C.int
	err := h3Error(C.h3goTest_linkedGeoPolygonToGeoPolygon(latLngPtr(flat),
		&nums[0], C.int(len(loops)), &outNums[0], latLngPtr(outVerts), &outLoops))
	if err != eSuccess {
		return nil, err
	}
	var out []GeoLoop
	v := 0
	for l := 0; l < int(outLoops); l++ {
		n := int(outNums[l])
		out = append(out, GeoLoop(outVerts[v:v+n]))
		v += n
	}
	return out, eSuccess
}

// cLinkedPolyState is the full serialized state of one linked polygon
// node: coords per loop, flattened vertices, and whether the C-side
// First/Last linkage invariants held.
type cLinkedPolyState struct {
	coordsPerLoop []int32
	verts         []LatLng
	invariantsOK  bool
}

// addLinkedGeoLoopC calls C addLinkedGeoLoop `times` times on one
// polygon and returns the full serialized state — including the
// partial state left behind on failure.
func addLinkedGeoLoopC(verts []LatLng, times int32) (cLinkedPolyState, h3Error) {
	outCoords := make([]C.int, times+1)
	outVerts := make([]LatLng, int(times)*len(verts)+1)
	var outLoops, inv C.int
	err := h3Error(C.h3goTest_addLinkedGeoLoop(latLngPtr(verts),
		C.int(len(verts)), C.int(times), &outLoops, &outCoords[0],
		latLngPtr(outVerts), &inv))
	st := cLinkedPolyState{invariantsOK: inv != 0}
	total := 0
	for i := 0; i < int(outLoops); i++ {
		st.coordsPerLoop = append(st.coordsPerLoop, int32(outCoords[i]))
		total += int(outCoords[i])
	}
	st.verts = outVerts[:total]
	return st, err
}

// geoPolygonToLinkedGeoLoopsC runs the C static on outer+holes and
// returns the full serialized state — including the partial state left
// behind when a hole loop fails.
func geoPolygonToLinkedGeoLoopsC(loops []GeoLoop) (cLinkedPolyState, h3Error) {
	var flat []LatLng
	nums := make([]C.int, len(loops))
	for i, l := range loops {
		flat = append(flat, l...)
		nums[i] = C.int(len(l))
	}
	outCoords := make([]C.int, len(loops)+1)
	outVerts := make([]LatLng, len(flat)+1)
	var outLoops, inv C.int
	err := h3Error(C.h3goTest_geoPolygonToLinkedGeoLoops(latLngPtr(flat),
		&nums[0], C.int(len(loops)), &outLoops, &outCoords[0],
		latLngPtr(outVerts), &inv))
	st := cLinkedPolyState{invariantsOK: inv != 0}
	total := 0
	for i := 0; i < int(outLoops); i++ {
		st.coordsPerLoop = append(st.coordsPerLoop, int32(outCoords[i]))
		total += int(outCoords[i])
	}
	st.verts = outVerts[:total]
	return st, err
}

// geoMultiPolygonToLinkedSyntheticC runs the flat->linked extern on a
// synthetic caller-supplied GeoMultiPolygon shape (identical bytes on
// both sides).
func geoMultiPolygonToLinkedSyntheticC(mp geoMultiPolygon) (cLinkedShape, bool, h3Error) {
	var polyNV, polyNH, holeNV []C.int
	var flat []LatLng
	for _, p := range mp.Polygons {
		polyNV = append(polyNV, C.int(len(p.GeoLoop)))
		polyNH = append(polyNH, C.int(len(p.Holes)))
		flat = append(flat, p.GeoLoop...)
		for _, h := range p.Holes {
			holeNV = append(holeNV, C.int(len(h)))
			flat = append(flat, h...)
		}
	}
	if len(holeNV) == 0 {
		holeNV = make([]C.int, 1)
	}
	loops := make([]C.int, 64)
	coords := make([]C.int, 64)
	outVerts := make([]LatLng, len(flat)+1)
	var np, inv C.int
	err := h3Error(C.h3goTest_geoMultiPolygonToLinkedSynthetic(
		&polyNV[0], &polyNH[0], &holeNV[0], latLngPtr(flat),
		C.int(len(mp.Polygons)), &np, &loops[0], &coords[0],
		latLngPtr(outVerts), &inv))
	if err != eSuccess {
		return cLinkedShape{}, false, err
	}
	var out cLinkedShape
	totalLoops := 0
	for p := 0; p < int(np); p++ {
		out.loopsPerPoly = append(out.loopsPerPoly, int32(loops[p]))
		totalLoops += int(loops[p])
	}
	totalCoords := 0
	for l := 0; l < totalLoops; l++ {
		out.coordsPerLoop = append(out.coordsPerLoop, int32(coords[l]))
		totalCoords += int(coords[l])
	}
	out.verts = outVerts[:totalCoords]
	return out, inv != 0, eSuccess
}

// cellsToLinkedMultiPolygonSerializedC calls ONLY the public C
// cellsToLinkedMultiPolygon and serializes its linked output directly
// (no conversion pipeline): chain shape, every vertex, and the
// linkage invariants.
func cellsToLinkedMultiPolygonSerializedC(cells []h3Index, numCells int32, maxLoops, maxVerts int) (cLinkedShape, bool, h3Error) {
	loops := make([]C.int, maxLoops)
	coords := make([]C.int, maxLoops)
	outVerts := make([]LatLng, maxVerts)
	var np, inv C.int
	err := h3Error(C.h3goTest_cellsToLinkedMultiPolygonSerialized(
		(*C.H3Index)(unsafe.Pointer(&cells[0])), C.int(numCells), &np,
		&loops[0], &coords[0], latLngPtr(outVerts), &inv))
	if err != eSuccess {
		return cLinkedShape{}, false, err
	}
	var out cLinkedShape
	totalLoops := 0
	for p := 0; p < int(np); p++ {
		out.loopsPerPoly = append(out.loopsPerPoly, int32(loops[p]))
		totalLoops += int(loops[p])
	}
	totalCoords := 0
	for l := 0; l < totalLoops; l++ {
		out.coordsPerLoop = append(out.coordsPerLoop, int32(coords[l]))
		totalCoords += int(coords[l])
	}
	out.verts = outVerts[:totalCoords]
	return out, inv != 0, eSuccess
}

// linkedToGeoMultiPolygonSyntheticC builds a synthetic linked chain
// from the flattened description (identical bytes on both sides),
// calls ONLY linkedGeoPolygonToGeoMultiPolygon, and returns the
// complete resulting GeoMultiPolygon.
func linkedToGeoMultiPolygonSyntheticC(loopsPerPoly []int32, coordsPerLoop []int32, verts []LatLng) (geoMultiPolygon, h3Error) {
	cLoops := make([]C.int, len(loopsPerPoly))
	for i, v := range loopsPerPoly {
		cLoops[i] = C.int(v)
	}
	cCoords := make([]C.int, len(coordsPerLoop))
	for i, v := range coordsPerLoop {
		cCoords[i] = C.int(v)
	}
	bound := int64(len(verts)) + 1
	polyNV := make([]int64, bound)
	polyNH := make([]int64, bound)
	holeNV := make([]int64, bound)
	outVerts := make([]LatLng, bound)
	var np C.int64_t
	err := h3Error(C.h3goTest_linkedToGeoMultiPolygonSynthetic(
		&cLoops[0], &cCoords[0], latLngPtr(verts), C.int(len(loopsPerPoly)),
		&np, (*C.int64_t)(&polyNV[0]), (*C.int64_t)(&polyNH[0]),
		(*C.int64_t)(&holeNV[0]), latLngPtr(outVerts)))
	if err != eSuccess {
		return geoMultiPolygon{}, err
	}
	out := geoMultiPolygon{NumPolygons: int32(np)}
	v, h := int64(0), int64(0)
	for p := int64(0); p < int64(np); p++ {
		poly := GeoPolygon{GeoLoop: GeoLoop(outVerts[v : v+polyNV[p]])}
		v += polyNV[p]
		for k := int64(0); k < polyNH[p]; k++ {
			poly.Holes = append(poly.Holes, GeoLoop(outVerts[v:v+holeNV[h]]))
			v += holeNV[h]
			h++
		}
		out.Polygons = append(out.Polygons, poly)
	}
	return out, eSuccess
}

// linkedConvertCleanupC runs a valid-outer/invalid-hole (or invalid
// second element) construction and reports the error plus whether the
// promised cleanup state held.
func linkedConvertCleanupC(which int32) (h3Error, bool) {
	var clean C.int
	err := h3Error(C.h3goTest_linkedConvertCleanup(C.int(which), &clean))
	return err, clean != 0
}

// linkedConvertErrorC exercises the conversion externs' error branches
// on the upstream testLinkedGeoConvert constructions.
func linkedConvertErrorC(which int32) h3Error {
	return h3Error(C.h3goTest_linkedConvertError(C.int(which)))
}

// cLinkedShape is the serialized shape of a linked polygon chain.
type cLinkedShape struct {
	loopsPerPoly  []int32
	coordsPerLoop []int32
	verts         []LatLng
}
