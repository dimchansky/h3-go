//go:build cgo && c2go && h3v450

// cellsToMultiPoly.c is new in H3 4.5.0 (docs/sync/4.4.0-to-4.5.0.md
// §5.2): the arc-based cellsToMultiPolygon machinery that the rewritten
// algos.c::cellsToLinkedMultiPolygon delegates to. Compiled only in the
// h3v450 harness configuration.
#include "cellsToMultiPoly.c"

// Test-only wrappers exposing the file-statics to the parity harness,
// in the same translation unit as the statics they call. Pointer-based
// arc state is serialized as element indices into the ArcSet's arcs
// array (identical construction order on both sides makes the indices
// comparable).

H3Error h3goTest_validateCellSet(const H3Index *cells, int64_t numCells) {
    return validateCellSet(cells, numCells);
}

int64_t h3goTest_getNumEdges(const H3Index *cells, int64_t numCells) {
    return getNumEdges(cells, numCells);
}

uint64_t h3goTest_hashEdge(H3Index x, uint64_t n) { return hashEdge(x, n); }

H3Error h3goTest_checkCellsToMultiPolyOverflow(int64_t numCells,
                                               int64_t hashMultiplier) {
    return checkCellsToMultiPolyOverflow(numCells, hashMultiplier);
}

// Serialize arc state after createArcSet (phase 0) or after
// cancelArcPairs (phase 1). Buffers must hold getNumEdges() entries.
H3Error h3goTest_arcState(const H3Index *cells, int64_t numCells, int phase,
                          H3Index *ids, uint8_t *removed, int64_t *nextIdx,
                          int64_t *prevIdx, H3Index *rootId,
                          int64_t *numArcsOut) {
    ArcSet arcset;
    H3Error err = createArcSet(cells, numCells, &arcset);
    if (err) return err;
    if (phase >= 1) {
        err = cancelArcPairs(arcset);
        if (err) {
            destroyArcSet(&arcset);
            return err;
        }
    }
    for (int64_t i = 0; i < arcset.numArcs; i++) {
        ids[i] = arcset.arcs[i].id;
        removed[i] = arcset.arcs[i].isRemoved ? 1 : 0;
        nextIdx[i] = arcset.arcs[i].next - arcset.arcs;
        prevIdx[i] = arcset.arcs[i].prev - arcset.arcs;
        rootId[i] = getRoot(&arcset.arcs[i])->id;
    }
    *numArcsOut = arcset.numArcs;
    destroyArcSet(&arcset);
    return E_SUCCESS;
}

// findArc lookup result as an arcs-array index; -1 when not found.
int64_t h3goTest_findArcIndex(const H3Index *cells, int64_t numCells,
                              H3Index e) {
    ArcSet arcset;
    if (createArcSet(cells, numCells, &arcset)) return -2;
    Arc *a = findArc(arcset, e);
    int64_t idx = a ? (a - arcset.arcs) : -1;
    destroyArcSet(&arcset);
    return idx;
}

int64_t h3goTest_countLoopsAfterCancel(const H3Index *cells,
                                       int64_t numCells) {
    ArcSet arcset;
    if (createArcSet(cells, numCells, &arcset)) return -1;
    if (cancelArcPairs(arcset)) {
        destroyArcSet(&arcset);
        return -1;
    }
    int64_t n = countLoops(arcset);
    destroyArcSet(&arcset);
    return n;
}

// Serialize the sorted loop set: per-loop root/area/vertex count plus
// the flattened vertices, and the resulting countPolys value. roots/
// areas/numVerts must hold getNumEdges() entries; verts must hold
// 2*getNumEdges() entries.
H3Error h3goTest_loopSet(const H3Index *cells, int64_t numCells,
                         int64_t *numLoopsOut, int64_t *numPolysOut,
                         H3Index *roots, double *areas, int64_t *numVerts,
                         LatLng *verts) {
    ArcSet arcset;
    H3Error err = createArcSet(cells, numCells, &arcset);
    if (err) return err;
    err = cancelArcPairs(arcset);
    if (err) {
        destroyArcSet(&arcset);
        return err;
    }
    SortableLoopSet loopset;
    err = createSortableLoopSet(arcset, &loopset);
    if (err) {
        destroyArcSet(&arcset);
        return err;
    }
    *numLoopsOut = loopset.numLoops;
    *numPolysOut = countPolys(loopset);
    int64_t v = 0;
    for (int64_t i = 0; i < loopset.numLoops; i++) {
        roots[i] = loopset.sloops[i].root;
        areas[i] = loopset.sloops[i].area;
        numVerts[i] = loopset.sloops[i].loop.numVerts;
        for (int64_t j = 0; j < loopset.sloops[i].loop.numVerts; j++) {
            verts[v++] = loopset.sloops[i].loop.verts[j];
        }
    }
    destroySortableLoopSet(&loopset);
    destroyArcSet(&arcset);
    return E_SUCCESS;
}

// Serialize createGlobeMultiPolygon: 8 triangles, flattened vertices.
// verts must hold 24 entries.
H3Error h3goTest_globeMultiPolygon(int64_t *numPolysOut, int64_t *numVertsOut,
                                   LatLng *verts) {
    GeoMultiPolygon mpoly;
    H3Error err = createGlobeMultiPolygon(&mpoly);
    if (err) return err;
    *numPolysOut = mpoly.numPolygons;
    int64_t v = 0;
    for (int i = 0; i < mpoly.numPolygons; i++) {
        if (mpoly.polygons[i].numHoles != 0) return E_FAILED;
        for (int j = 0; j < mpoly.polygons[i].geoloop.numVerts; j++) {
            verts[v++] = mpoly.polygons[i].geoloop.verts[j];
        }
    }
    *numVertsOut = v;
    H3_EXPORT(destroyGeoMultiPolygon)(&mpoly);
    return E_SUCCESS;
}
