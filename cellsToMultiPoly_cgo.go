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
H3Error h3goTest_cellToEdgeArcs(H3Index h, H3Index *ids, int64_t *nextIdx,
                                int64_t *prevIdx, int64_t *parentIdx,
                                int64_t *rank, uint8_t *visited,
                                uint8_t *removed, int64_t *numEdgesOut);
H3Error h3goTest_createSortablePolyFromLoops(
    const H3Index *roots, const double *areas, const int64_t *numVerts,
    LatLng *verts, int64_t numLoops, int64_t loopStart, int64_t numHoles,
    double *outerAreaOut, int64_t *outerNumVertsOut, LatLng *outerVerts,
    int64_t *holeNumVerts, LatLng *holeVerts);
H3Error h3goTest_createMultiPolygonFromLoops(
    const H3Index *roots, const double *areas, const int64_t *numVerts,
    LatLng *verts, int64_t numLoops, int64_t *numPolysOut,
    int64_t *polyNumVerts, int64_t *polyNumHoles, int64_t *holeNumVerts,
    LatLng *outVerts);
H3Error h3goTest_bucketState(const H3Index *cells, int64_t numCells,
                             int64_t *bucketArcIdx, int64_t *numBucketsOut);
H3Error h3goTest_visitedState(const H3Index *cells, int64_t numCells,
                              int mode, uint8_t *visited, int64_t *numArcsOut);
H3Error h3goTest_unionSequence(const H3Index *cells, int64_t numCells,
                               const int64_t *pairA, const int64_t *pairB,
                               int64_t numPairs, H3Index *rootId,
                               int64_t *rank, int64_t *numArcsOut);
H3Error h3goTest_createSortableLoop(const H3Index *cells, int64_t numCells,
                                    int64_t arcIdx, H3Index *rootOut,
                                    double *areaOut, int64_t *numVertsOut,
                                    LatLng *verts);
int h3goTest_cmp_SortableLoop(H3Index rootA, double areaA, H3Index rootB,
                              double areaB);
int h3goTest_cmp_SortablePoly(double areaA, double areaB);
int h3goTest_cmp_uint64(H3Index a, H3Index b);
int h3goTest_destroyArcSet_state(const H3Index *cells, int64_t numCells);
int h3goTest_destroyLoopSet_state(const H3Index *cells, int64_t numCells,
                                  int shallow);
int h3goTest_destroyGeoLoop_state(void);
int h3goTest_destroyGeoPolygon_state(void);
int h3goTest_destroyGeoMultiPolygon_state(const H3Index *cells,
                                          int64_t numCells);
int h3goTest_destroyLinkedTwice(const H3Index *cells, int numCells);

// linkedGeo.h externs (4.5.0 conversion helpers).
H3Error linkedGeoPolygonToGeoMultiPolygon(const LinkedGeoPolygon *linked,
                                          GeoMultiPolygon *out);
H3Error geoMultiPolygonToLinkedGeoPolygon(const GeoMultiPolygon *mpoly,
                                          LinkedGeoPolygon *out);

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

// cEdgeArcs mirrors the serialized single-cell cellToEdgeArcs state.
type cEdgeArcs struct {
	ids                []h3Index
	nextIdx, prevIdx   []int64
	parentIdx, rank    []int64
	isVisited, removed []bool
}

// cellToEdgeArcsC serializes the C cellToEdgeArcs output for one cell.
func cellToEdgeArcsC(h h3Index) (cEdgeArcs, h3Error) {
	ids := make([]h3Index, 6)
	next := make([]int64, 6)
	prev := make([]int64, 6)
	parent := make([]int64, 6)
	rank := make([]int64, 6)
	visited := make([]uint8, 6)
	removed := make([]uint8, 6)
	var n C.int64_t
	err := h3Error(C.h3goTest_cellToEdgeArcs(C.H3Index(h),
		(*C.H3Index)(unsafe.Pointer(&ids[0])), (*C.int64_t)(&next[0]),
		(*C.int64_t)(&prev[0]), (*C.int64_t)(&parent[0]),
		(*C.int64_t)(&rank[0]), (*C.uint8_t)(&visited[0]),
		(*C.uint8_t)(&removed[0]), &n))
	if err != eSuccess {
		return cEdgeArcs{}, err
	}
	out := cEdgeArcs{ids: ids[:n], nextIdx: next[:n], prevIdx: prev[:n],
		parentIdx: parent[:n], rank: rank[:n]}
	for i := int64(0); i < int64(n); i++ {
		out.isVisited = append(out.isVisited, visited[i] != 0)
		out.removed = append(out.removed, removed[i] != 0)
	}
	return out, eSuccess
}

// bucketStateC serializes the hash-bucket layout after C createArcSet.
func bucketStateC(cells []h3Index, numCells int64) ([]int64, h3Error) {
	n := getNumEdgesC(cells, numCells) * hashTableMultiplier
	buckets := make([]int64, n)
	var nb C.int64_t
	err := h3Error(C.h3goTest_bucketState(cellsPtr(cells), C.int64_t(numCells),
		(*C.int64_t)(&buckets[0]), &nb))
	if err != eSuccess {
		return nil, err
	}
	return buckets[:nb], eSuccess
}

// visitedStateC serializes isVisited after countLoops (mode 0) or
// countLoops+resetVisited (mode 1).
func visitedStateC(cells []h3Index, numCells int64, mode int32) ([]bool, h3Error) {
	n := getNumEdgesC(cells, numCells)
	visited := make([]uint8, n)
	var na C.int64_t
	err := h3Error(C.h3goTest_visitedState(cellsPtr(cells), C.int64_t(numCells),
		C.int(mode), (*C.uint8_t)(&visited[0]), &na))
	if err != eSuccess {
		return nil, err
	}
	out := make([]bool, na)
	for i := range out {
		out[i] = visited[i] != 0
	}
	return out, eSuccess
}

// unionSequenceC unions the given arc-index pairs on the C side and
// serializes every arc's root id and rank.
func unionSequenceC(cells []h3Index, numCells int64, pairs [][2]int64) ([]h3Index, []int64, h3Error) {
	n := getNumEdgesC(cells, numCells)
	pa := make([]int64, len(pairs))
	pb := make([]int64, len(pairs))
	for i, p := range pairs {
		pa[i], pb[i] = p[0], p[1]
	}
	roots := make([]h3Index, n)
	ranks := make([]int64, n)
	var na C.int64_t
	err := h3Error(C.h3goTest_unionSequence(cellsPtr(cells), C.int64_t(numCells),
		(*C.int64_t)(&pa[0]), (*C.int64_t)(&pb[0]), C.int64_t(len(pairs)),
		(*C.H3Index)(unsafe.Pointer(&roots[0])), (*C.int64_t)(&ranks[0]), &na))
	if err != eSuccess {
		return nil, nil, err
	}
	return roots[:na], ranks[:na], eSuccess
}

// createSortableLoopC calls C createSortableLoop on the arc at arcIdx.
func createSortableLoopC(cells []h3Index, numCells, arcIdx int64) (sortableLoop, h3Error) {
	n := getNumEdgesC(cells, numCells)
	verts := make([]LatLng, 2*n)
	var root C.H3Index
	var area C.double
	var nv C.int64_t
	err := h3Error(C.h3goTest_createSortableLoop(cellsPtr(cells), C.int64_t(numCells),
		C.int64_t(arcIdx), &root, &area, &nv,
		(*C.LatLng)(unsafe.Pointer(&verts[0]))))
	if err != eSuccess {
		return sortableLoop{}, err
	}
	return sortableLoop{root: h3Index(root), area: float64(area),
		loop: GeoLoop(verts[:nv])}, eSuccess
}

func cmp_SortableLoopC(rootA h3Index, areaA float64, rootB h3Index, areaB float64) int32 {
	return int32(C.h3goTest_cmp_SortableLoop(C.H3Index(rootA), C.double(areaA),
		C.H3Index(rootB), C.double(areaB)))
}

func cmp_SortablePolyC(areaA, areaB float64) int32 {
	return int32(C.h3goTest_cmp_SortablePoly(C.double(areaA), C.double(areaB)))
}

func cmp_uint64C(a, b h3Index) int32 {
	return int32(C.h3goTest_cmp_uint64(C.H3Index(a), C.H3Index(b)))
}

// flattenLoops marshals a synthetic loop set for the FromLoops
// wrappers (identical bytes fed to both sides).
func flattenLoops(loops []sortableLoop) (roots []h3Index, areas []float64, numVerts []int64, verts []LatLng) {
	for _, l := range loops {
		roots = append(roots, l.root)
		areas = append(areas, l.area)
		numVerts = append(numVerts, int64(len(l.loop)))
		verts = append(verts, l.loop...)
	}
	return
}

// createSortablePolyFromLoopsC calls C createSortablePoly on a
// synthetic caller-supplied loop set.
func createSortablePolyFromLoopsC(loops []sortableLoop, loopStart, numHoles int64) (sortablePoly, h3Error) {
	roots, areas, numVerts, verts := flattenLoops(loops)
	outerVerts := make([]LatLng, len(verts)+1)
	holeVerts := make([]LatLng, len(verts)+1)
	holeNum := make([]int64, len(loops)+1)
	var outerArea C.double
	var outerNV C.int64_t
	err := h3Error(C.h3goTest_createSortablePolyFromLoops(
		(*C.H3Index)(unsafe.Pointer(&roots[0])), (*C.double)(&areas[0]),
		(*C.int64_t)(&numVerts[0]), (*C.LatLng)(unsafe.Pointer(&verts[0])),
		C.int64_t(len(loops)), C.int64_t(loopStart), C.int64_t(numHoles),
		&outerArea, &outerNV, (*C.LatLng)(unsafe.Pointer(&outerVerts[0])),
		(*C.int64_t)(&holeNum[0]), (*C.LatLng)(unsafe.Pointer(&holeVerts[0]))))
	if err != eSuccess {
		return sortablePoly{}, err
	}
	out := sortablePoly{outerArea: float64(outerArea),
		poly: GeoPolygon{GeoLoop: GeoLoop(outerVerts[:outerNV])}}
	v := int64(0)
	for h := int64(0); h < numHoles; h++ {
		out.poly.Holes = append(out.poly.Holes, GeoLoop(holeVerts[v:v+holeNum[h]]))
		v += holeNum[h]
	}
	return out, eSuccess
}

// createMultiPolygonFromLoopsC calls C createMultiPolygon on a
// synthetic caller-supplied loop set (empty = globe branch).
func createMultiPolygonFromLoopsC(loops []sortableLoop) (geoMultiPolygon, h3Error) {
	roots, areas, numVerts, verts := flattenLoops(loops)
	// Padded one-element inputs keep the C pointers valid for the
	// empty (globe) case.
	if len(loops) == 0 {
		roots, areas, numVerts, verts = make([]h3Index, 1), make([]float64, 1), make([]int64, 1), make([]LatLng, 1)
	}
	bound := int64(len(verts)) + 30
	polyNV := make([]int64, bound)
	polyNH := make([]int64, bound)
	holeNV := make([]int64, bound)
	outVerts := make([]LatLng, 2*bound)
	var np C.int64_t
	err := h3Error(C.h3goTest_createMultiPolygonFromLoops(
		(*C.H3Index)(unsafe.Pointer(&roots[0])), (*C.double)(&areas[0]),
		(*C.int64_t)(&numVerts[0]), (*C.LatLng)(unsafe.Pointer(&verts[0])),
		C.int64_t(len(loops)), &np, (*C.int64_t)(&polyNV[0]),
		(*C.int64_t)(&polyNH[0]), (*C.int64_t)(&holeNV[0]),
		(*C.LatLng)(unsafe.Pointer(&outVerts[0]))))
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

// Destroy-helper state probes (bit 1: nulled after first call; bit 2:
// second call safe and still nulled).
func destroyArcSetStateC(cells []h3Index, numCells int64) int32 {
	return int32(C.h3goTest_destroyArcSet_state(cellsPtr(cells), C.int64_t(numCells)))
}

func destroyLoopSetStateC(cells []h3Index, numCells int64, shallow int32) int32 {
	return int32(C.h3goTest_destroyLoopSet_state(cellsPtr(cells), C.int64_t(numCells), C.int(shallow)))
}

func destroyGeoLoopStateC() int32    { return int32(C.h3goTest_destroyGeoLoop_state()) }
func destroyGeoPolygonStateC() int32 { return int32(C.h3goTest_destroyGeoPolygon_state()) }
func destroyGeoMultiPolygonStateC(cells []h3Index, numCells int64) int32 {
	return int32(C.h3goTest_destroyGeoMultiPolygon_state(cellsPtr(cells), C.int64_t(numCells)))
}

// destroyLinkedTwiceC exercises the 4.5.0 destroyLinkedMultiPolygon
// idempotence on the C side (bits as above).
func destroyLinkedTwiceC(cells []h3Index, numCells int32) int32 {
	return int32(C.h3goTest_destroyLinkedTwice(cellsPtr(cells), C.int(numCells)))
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
