//go:build cgo && c2go && h3v450

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
                                  int *outNumLoops, int *outCoordsPerLoop);
H3Error h3goTest_geoPolygonToLinkedGeoLoops(const LatLng *verts,
                                            const int *loopNumVerts,
                                            int numLoops, int *outNumLoops,
                                            int *outCoordsPerLoop);
H3Error h3goTest_linkedConvertError(int which);
H3Error h3goTest_geoMultiPolygonToLinked(const H3Index *cells, int64_t n,
                                         int *numPolysOut, int *loopsPerPoly,
                                         int *coordsPerLoop, LatLng *verts);
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

// addLinkedGeoLoopC calls C addLinkedGeoLoop `times` times on one
// polygon and returns the loop count and coords per loop.
func addLinkedGeoLoopC(verts []LatLng, times int32) (int32, []int32, h3Error) {
	outCoords := make([]C.int, times+1)
	var outLoops C.int
	err := h3Error(C.h3goTest_addLinkedGeoLoop(latLngPtr(verts),
		C.int(len(verts)), C.int(times), &outLoops, &outCoords[0]))
	coords := make([]int32, outLoops)
	for i := range coords {
		coords[i] = int32(outCoords[i])
	}
	return int32(outLoops), coords, err
}

// geoPolygonToLinkedGeoLoopsC runs the C static on outer+holes and
// returns the serialized loop/coord counts.
func geoPolygonToLinkedGeoLoopsC(loops []GeoLoop) (int32, []int32, h3Error) {
	var flat []LatLng
	nums := make([]C.int, len(loops))
	for i, l := range loops {
		flat = append(flat, l...)
		nums[i] = C.int(len(l))
	}
	outCoords := make([]C.int, len(loops)+1)
	var outLoops C.int
	err := h3Error(C.h3goTest_geoPolygonToLinkedGeoLoops(latLngPtr(flat),
		&nums[0], C.int(len(loops)), &outLoops, &outCoords[0]))
	if err != eSuccess {
		return 0, nil, err
	}
	coords := make([]int32, outLoops)
	for i := range coords {
		coords[i] = int32(outCoords[i])
	}
	return int32(outLoops), coords, eSuccess
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

// geoMultiPolygonToLinkedC runs the isolated C success path of
// geoMultiPolygonToLinkedGeoPolygon on cellsToMultiPolygon(cells).
func geoMultiPolygonToLinkedC(cells []h3Index, n int64, maxLoops, maxVerts int) (cLinkedShape, h3Error) {
	loops := make([]C.int, maxLoops)
	coords := make([]C.int, maxLoops)
	verts := make([]LatLng, maxVerts)
	var np C.int
	err := h3Error(C.h3goTest_geoMultiPolygonToLinked(
		(*C.H3Index)(unsafe.Pointer(&cells[0])), C.int64_t(n), &np,
		&loops[0], &coords[0], latLngPtr(verts)))
	if err != eSuccess {
		return cLinkedShape{}, err
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
	out.verts = verts[:totalCoords]
	return out, eSuccess
}
