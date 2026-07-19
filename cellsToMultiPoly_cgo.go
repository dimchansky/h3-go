//go:build cgo && c2go && h3v450

package h3

// Bridges to the H3 4.5.0 cellsToMultiPoly machinery: the public (in C)
// cellsToMultiPolygon/destroyGeoMultiPolygon, the extern linkedGeo
// conversions, and the file-statics via the same-TU h3goTest_* wrappers
// in h3lib_cellsToMultiPoly_c2go.c.

/*
#include <stdint.h>
#include <stdlib.h>
#include "h3api.h"

H3Error cellsToMultiPolygon(const H3Index *cells, const int64_t numCells,
                            GeoMultiPolygon *out);
void destroyGeoMultiPolygon(GeoMultiPolygon *mpoly);

H3Error h3goTest_validateCellSet(const H3Index *cells, int64_t numCells);
int64_t h3goTest_getNumEdges(const H3Index *cells, int64_t numCells);
uint64_t h3goTest_hashEdge(H3Index x, uint64_t n);
H3Error h3goTest_checkCellsToMultiPolyOverflow(int64_t numCells,
                                               int64_t hashMultiplier);
H3Error h3goTest_arcState(const H3Index *cells, int64_t numCells, int phase,
                          H3Index *ids, uint8_t *removed, int64_t *nextIdx,
                          int64_t *prevIdx, H3Index *rootId,
                          int64_t *numArcsOut);
int64_t h3goTest_findArcIndex(const H3Index *cells, int64_t numCells,
                              H3Index e);
int64_t h3goTest_countLoopsAfterCancel(const H3Index *cells,
                                       int64_t numCells);
H3Error h3goTest_loopSet(const H3Index *cells, int64_t numCells,
                         int64_t *numLoopsOut, int64_t *numPolysOut,
                         H3Index *roots, double *areas, int64_t *numVerts,
                         LatLng *verts);
H3Error h3goTest_globeMultiPolygon(int64_t *numPolysOut, int64_t *numVertsOut,
                                   LatLng *verts);

// linkedGeo.h externs (4.5.0 conversion helpers).
H3Error linkedGeoPolygonToGeoMultiPolygon(const LinkedGeoPolygon *linked,
                                          GeoMultiPolygon *out);
H3Error geoMultiPolygonToLinkedGeoPolygon(const GeoMultiPolygon *mpoly,
                                          LinkedGeoPolygon *out);

// Run the full C linked pipeline (cellsToLinkedMultiPolygon) and
// flatten its output through C's own linkedGeoPolygonToGeoMultiPolygon
// so both the linked algorithm and the conversion extern are exercised.
static H3Error linkedMultiPolyAsGeo_c(const H3Index *cells, int numCells,
                                      GeoMultiPolygon *out) {
    LinkedGeoPolygon linked;
    H3Error err = cellsToLinkedMultiPolygon(cells, numCells, &linked);
    if (err) return err;
    err = linkedGeoPolygonToGeoMultiPolygon(&linked, out);
    destroyLinkedMultiPolygon(&linked);
    return err;
}
*/
import "C"
import "unsafe"

func cellsPtr(cells []h3Index) *C.H3Index {
	if len(cells) == 0 {
		return nil
	}
	return (*C.H3Index)(unsafe.Pointer(&cells[0]))
}

// validateCellSetC calls the original C implementation (same-TU wrapper).
func validateCellSetC(cells []h3Index, numCells int64) h3Error {
	return h3Error(C.h3goTest_validateCellSet(cellsPtr(cells), C.int64_t(numCells)))
}

// getNumEdgesC calls the original C implementation (same-TU wrapper).
func getNumEdgesC(cells []h3Index, numCells int64) int64 {
	return int64(C.h3goTest_getNumEdges(cellsPtr(cells), C.int64_t(numCells)))
}

// hashEdgeC calls the original C implementation (same-TU wrapper).
func hashEdgeC(x h3Index, n uint64) uint64 {
	return uint64(C.h3goTest_hashEdge(C.H3Index(x), C.uint64_t(n)))
}

// checkCellsToMultiPolyOverflowC calls the original C implementation.
func checkCellsToMultiPolyOverflowC(numCells, hashMultiplier int64) h3Error {
	return h3Error(C.h3goTest_checkCellsToMultiPolyOverflow(
		C.int64_t(numCells), C.int64_t(hashMultiplier)))
}

// cArcState mirrors the serialized C arc state.
type cArcState struct {
	ids     []h3Index
	removed []bool
	nextIdx []int64
	prevIdx []int64
	rootID  []h3Index
}

// arcStateC serializes the C-side arc state after createArcSet
// (phase 0) or cancelArcPairs (phase 1).
func arcStateC(cells []h3Index, numCells int64, phase int32) (cArcState, h3Error) {
	n := getNumEdgesC(cells, numCells)
	ids := make([]h3Index, n)
	removed := make([]uint8, n)
	next := make([]int64, n)
	prev := make([]int64, n)
	root := make([]h3Index, n)
	var numArcs C.int64_t
	err := h3Error(C.h3goTest_arcState(cellsPtr(cells), C.int64_t(numCells), C.int(phase),
		(*C.H3Index)(unsafe.Pointer(&ids[0])), (*C.uint8_t)(&removed[0]),
		(*C.int64_t)(&next[0]), (*C.int64_t)(&prev[0]),
		(*C.H3Index)(unsafe.Pointer(&root[0])), &numArcs))
	if err != eSuccess {
		return cArcState{}, err
	}
	st := cArcState{ids: ids[:numArcs], nextIdx: next[:numArcs],
		prevIdx: prev[:numArcs], rootID: root[:numArcs]}
	for _, r := range removed[:numArcs] {
		st.removed = append(st.removed, r != 0)
	}
	return st, eSuccess
}

// findArcIndexC returns the C findArc result as an arcs-array index
// (-1 when the edge is not in the set).
func findArcIndexC(cells []h3Index, numCells int64, e h3Index) int64 {
	return int64(C.h3goTest_findArcIndex(cellsPtr(cells), C.int64_t(numCells), C.H3Index(e)))
}

// countLoopsAfterCancelC runs createArcSet+cancelArcPairs+countLoops in C.
func countLoopsAfterCancelC(cells []h3Index, numCells int64) int64 {
	return int64(C.h3goTest_countLoopsAfterCancel(cellsPtr(cells), C.int64_t(numCells)))
}

// cLoopSet mirrors the serialized C sortable loop set.
type cLoopSet struct {
	numPolys int64
	roots    []h3Index
	areas    []float64
	loops    []GeoLoop
}

// loopSetC serializes the C-side sorted loop set plus countPolys.
func loopSetC(cells []h3Index, numCells int64) (cLoopSet, h3Error) {
	n := getNumEdgesC(cells, numCells)
	roots := make([]h3Index, n)
	areas := make([]float64, n)
	numVerts := make([]int64, n)
	verts := make([]LatLng, 2*n)
	var numLoops, numPolys C.int64_t
	err := h3Error(C.h3goTest_loopSet(cellsPtr(cells), C.int64_t(numCells),
		&numLoops, &numPolys,
		(*C.H3Index)(unsafe.Pointer(&roots[0])), (*C.double)(&areas[0]),
		(*C.int64_t)(&numVerts[0]), (*C.LatLng)(unsafe.Pointer(&verts[0]))))
	if err != eSuccess {
		return cLoopSet{}, err
	}
	out := cLoopSet{numPolys: int64(numPolys),
		roots: roots[:numLoops], areas: areas[:numLoops]}
	v := int64(0)
	for i := int64(0); i < int64(numLoops); i++ {
		out.loops = append(out.loops, GeoLoop(verts[v:v+numVerts[i]]))
		v += numVerts[i]
	}
	return out, eSuccess
}

// globeMultiPolygonC serializes the C createGlobeMultiPolygon output.
func globeMultiPolygonC() (numPolys int64, verts []LatLng, err h3Error) {
	out := make([]LatLng, 24)
	var np, nv C.int64_t
	err = h3Error(C.h3goTest_globeMultiPolygon(&np, &nv,
		(*C.LatLng)(unsafe.Pointer(&out[0]))))
	return int64(np), out[:nv], err
}

// cellsToMultiPolygonC calls the original C implementation and marshals
// the GeoMultiPolygon into the Go representation.
func cellsToMultiPolygonC(cells []h3Index, numCells int64) (geoMultiPolygon, h3Error) {
	var cOut C.GeoMultiPolygon
	err := h3Error(C.cellsToMultiPolygon(cellsPtr(cells), C.int64_t(numCells), &cOut))
	if err != eSuccess {
		return geoMultiPolygon{}, err
	}
	var out geoMultiPolygon
	out.NumPolygons = int32(cOut.numPolygons)
	if cOut.numPolygons > 0 {
		cPolys := (*[1 << 20]C.GeoPolygon)(unsafe.Pointer(cOut.polygons))[:cOut.numPolygons:cOut.numPolygons]
		for i := range cPolys {
			out.Polygons = append(out.Polygons, geoPolygonFromC(&cPolys[i]))
		}
	}
	C.destroyGeoMultiPolygon(&cOut)
	return out, eSuccess
}

// linkedMultiPolyAsGeoC runs C cellsToLinkedMultiPolygon and flattens
// the linked output through C linkedGeoPolygonToGeoMultiPolygon.
func linkedMultiPolyAsGeoC(cells []h3Index, numCells int32) (geoMultiPolygon, h3Error) {
	var cOut C.GeoMultiPolygon
	err := h3Error(C.linkedMultiPolyAsGeo_c(cellsPtr(cells), C.int(numCells), &cOut))
	if err != eSuccess {
		return geoMultiPolygon{}, err
	}
	var out geoMultiPolygon
	out.NumPolygons = int32(cOut.numPolygons)
	if cOut.numPolygons > 0 {
		cPolys := (*[1 << 20]C.GeoPolygon)(unsafe.Pointer(cOut.polygons))[:cOut.numPolygons:cOut.numPolygons]
		for i := range cPolys {
			out.Polygons = append(out.Polygons, geoPolygonFromC(&cPolys[i]))
		}
	}
	C.destroyGeoMultiPolygon(&cOut)
	return out, eSuccess
}

func geoLoopFromC(cl C.GeoLoop) GeoLoop {
	n := int(cl.numVerts)
	if n == 0 {
		return nil
	}
	cVerts := (*[1 << 26]C.LatLng)(unsafe.Pointer(cl.verts))[:n:n]
	loop := make(GeoLoop, n)
	for i, v := range cVerts {
		loop[i] = LatLng{Lat: Rad(float64(v.lat)), Lng: Rad(float64(v.lng))}
	}
	return loop
}

func geoPolygonFromC(cp *C.GeoPolygon) GeoPolygon {
	poly := GeoPolygon{GeoLoop: geoLoopFromC(cp.geoloop)}
	n := int(cp.numHoles)
	if n > 0 {
		cHoles := (*[1 << 20]C.GeoLoop)(unsafe.Pointer(cp.holes))[:n:n]
		for i := range cHoles {
			poly.Holes = append(poly.Holes, geoLoopFromC(cHoles[i]))
		}
	}
	return poly
}
