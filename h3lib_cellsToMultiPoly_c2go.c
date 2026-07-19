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

// Direct single-cell cellToEdgeArcs state: ids, linkage as element
// indices, parent index, rank, and the explicitly initialized
// isVisited/isRemoved flags. Buffers hold 6 entries.
H3Error h3goTest_cellToEdgeArcs(H3Index h, H3Index *ids, int64_t *nextIdx,
                                int64_t *prevIdx, int64_t *parentIdx,
                                int64_t *rank, uint8_t *visited,
                                uint8_t *removed, int64_t *numEdgesOut) {
    Arc arcs[6];
    // Poison the flags so the comparison proves cellToEdgeArcs itself
    // clears them (it explicitly initializes both fields).
    for (int i = 0; i < 6; i++) {
        arcs[i].isVisited = true;
        arcs[i].isRemoved = true;
    }
    H3Error err = cellToEdgeArcs(h, arcs, numEdgesOut);
    if (err) return err;
    for (int64_t i = 0; i < *numEdgesOut; i++) {
        ids[i] = arcs[i].id;
        nextIdx[i] = arcs[i].next - arcs;
        prevIdx[i] = arcs[i].prev - arcs;
        parentIdx[i] = arcs[i].parent - arcs;
        rank[i] = arcs[i].rank;
        visited[i] = arcs[i].isVisited ? 1 : 0;
        removed[i] = arcs[i].isRemoved ? 1 : 0;
    }
    return E_SUCCESS;
}

// Build a SortableLoopSet from caller-supplied synthetic loops
// (identical bytes on both sides) — shared by the FromLoops wrappers.
static void h3goTestBuildLoopSet(const H3Index *roots, const double *areas,
                                 const int64_t *numVerts, LatLng *verts,
                                 int64_t numLoops, SortableLoop *sloops,
                                 SortableLoopSet *loopset) {
    int64_t v = 0;
    for (int64_t i = 0; i < numLoops; i++) {
        sloops[i].root = roots[i];
        sloops[i].area = areas[i];
        sloops[i].loop.numVerts = numVerts[i];
        sloops[i].loop.verts = &verts[v];
        v += numVerts[i];
    }
    loopset->numLoops = numLoops;
    loopset->sloops = sloops;
}

// Direct createSortablePoly on a synthetic, caller-supplied loop set.
H3Error h3goTest_createSortablePolyFromLoops(
    const H3Index *roots, const double *areas, const int64_t *numVerts,
    LatLng *verts, int64_t numLoops, int64_t loopStart, int64_t numHoles,
    double *outerAreaOut, int64_t *outerNumVertsOut, LatLng *outerVerts,
    int64_t *holeNumVerts, LatLng *holeVerts) {
    SortableLoop sloops[16];
    SortableLoopSet loopset;
    h3goTestBuildLoopSet(roots, areas, numVerts, verts, numLoops, sloops,
                         &loopset);
    SortablePoly spoly;
    H3Error err =
        createSortablePoly(&loopset.sloops[loopStart], numHoles, &spoly);
    if (err) return err;
    *outerAreaOut = spoly.outerArea;
    *outerNumVertsOut = spoly.poly.geoloop.numVerts;
    for (int64_t i = 0; i < spoly.poly.geoloop.numVerts; i++) {
        outerVerts[i] = spoly.poly.geoloop.verts[i];
    }
    int64_t hv = 0;
    for (int64_t h = 0; h < numHoles; h++) {
        holeNumVerts[h] = spoly.poly.holes[h].numVerts;
        for (int64_t i = 0; i < spoly.poly.holes[h].numVerts; i++) {
            holeVerts[hv++] = spoly.poly.holes[h].verts[i];
        }
    }
    if (spoly.poly.holes) H3_MEMORY(free)(spoly.poly.holes);
    return E_SUCCESS;
}

// Direct createMultiPolygon on a synthetic, caller-supplied loop set,
// serialized as flattened per-polygon vertex/hole counts plus vertices.
// numLoops == 0 drives the createGlobeMultiPolygon branch.
H3Error h3goTest_createMultiPolygonFromLoops(
    const H3Index *roots, const double *areas, const int64_t *numVerts,
    LatLng *verts, int64_t numLoops, int64_t *numPolysOut,
    int64_t *polyNumVerts, int64_t *polyNumHoles, int64_t *holeNumVerts,
    LatLng *outVerts) {
    SortableLoop sloops[16];
    SortableLoopSet loopset;
    h3goTestBuildLoopSet(roots, areas, numVerts, verts, numLoops, sloops,
                         &loopset);
    GeoMultiPolygon mpoly;
    H3Error err = createMultiPolygon(loopset, &mpoly);
    if (err) return err;
    *numPolysOut = mpoly.numPolygons;
    int64_t v = 0, h = 0;
    for (int p = 0; p < mpoly.numPolygons; p++) {
        polyNumVerts[p] = mpoly.polygons[p].geoloop.numVerts;
        polyNumHoles[p] = mpoly.polygons[p].numHoles;
        for (int i = 0; i < mpoly.polygons[p].geoloop.numVerts; i++) {
            outVerts[v++] = mpoly.polygons[p].geoloop.verts[i];
        }
        for (int k = 0; k < mpoly.polygons[p].numHoles; k++) {
            holeNumVerts[h++] = mpoly.polygons[p].holes[k].numVerts;
            for (int i = 0; i < mpoly.polygons[p].holes[k].numVerts; i++) {
                outVerts[v++] = mpoly.polygons[p].holes[k].verts[i];
            }
        }
    }
    if (numLoops == 0) {
        // Globe branch: the octant vertices are C-allocated.
        H3_EXPORT(destroyGeoMultiPolygon)(&mpoly);
    } else {
        // The output polygons alias the synthetic loop set's caller
        // memory (createSortablePoly copies GeoLoop headers, not
        // vertices), so only the C-allocated arrays are freed here.
        for (int p = 0; p < mpoly.numPolygons; p++) {
            if (mpoly.polygons[p].holes) {
                H3_MEMORY(free)(mpoly.polygons[p].holes);
            }
        }
        H3_MEMORY(free)(mpoly.polygons);
    }
    return E_SUCCESS;
}

// Hash-bucket layout after createArcSet: for each bucket, the index of
// the arc it holds (-1 when empty). buckets must hold
// getNumEdges()*HASH_TABLE_MULTIPLIER entries.
H3Error h3goTest_bucketState(const H3Index *cells, int64_t numCells,
                             int64_t *bucketArcIdx, int64_t *numBucketsOut) {
    ArcSet arcset;
    H3Error err = createArcSet(cells, numCells, &arcset);
    if (err) return err;
    for (int64_t j = 0; j < arcset.numBuckets; j++) {
        bucketArcIdx[j] =
            arcset.buckets[j] ? (arcset.buckets[j] - arcset.arcs) : -1;
    }
    *numBucketsOut = arcset.numBuckets;
    destroyArcSet(&arcset);
    return E_SUCCESS;
}

// isVisited flags after countLoops (mode 0) or after countLoops +
// resetVisited (mode 1) — the direct resetVisited observation.
H3Error h3goTest_visitedState(const H3Index *cells, int64_t numCells,
                              int mode, uint8_t *visited,
                              int64_t *numArcsOut) {
    ArcSet arcset;
    H3Error err = createArcSet(cells, numCells, &arcset);
    if (err) return err;
    err = cancelArcPairs(arcset);
    if (err) {
        destroyArcSet(&arcset);
        return err;
    }
    countLoops(arcset);
    if (mode >= 1) resetVisited(arcset);
    for (int64_t i = 0; i < arcset.numArcs; i++) {
        visited[i] = arcset.arcs[i].isVisited ? 1 : 0;
    }
    *numArcsOut = arcset.numArcs;
    destroyArcSet(&arcset);
    return E_SUCCESS;
}

// Direct unionArcs/getRoot: union the given arc-index pairs on a fresh
// ArcSet, then serialize every arc's root id and rank.
H3Error h3goTest_unionSequence(const H3Index *cells, int64_t numCells,
                               const int64_t *pairA, const int64_t *pairB,
                               int64_t numPairs, H3Index *rootId,
                               int64_t *rank, int64_t *numArcsOut) {
    ArcSet arcset;
    H3Error err = createArcSet(cells, numCells, &arcset);
    if (err) return err;
    for (int64_t p = 0; p < numPairs; p++) {
        unionArcs(&arcset.arcs[pairA[p]], &arcset.arcs[pairB[p]]);
    }
    for (int64_t i = 0; i < arcset.numArcs; i++) {
        rootId[i] = getRoot(&arcset.arcs[i])->id;
        rank[i] = arcset.arcs[i].rank;
    }
    *numArcsOut = arcset.numArcs;
    destroyArcSet(&arcset);
    return E_SUCCESS;
}

// Direct createSortableLoop on the arc at arcIdx (after cancellation).
H3Error h3goTest_createSortableLoop(const H3Index *cells, int64_t numCells,
                                    int64_t arcIdx, H3Index *rootOut,
                                    double *areaOut, int64_t *numVertsOut,
                                    LatLng *verts) {
    ArcSet arcset;
    H3Error err = createArcSet(cells, numCells, &arcset);
    if (err) return err;
    err = cancelArcPairs(arcset);
    if (err) {
        destroyArcSet(&arcset);
        return err;
    }
    SortableLoop sloop;
    err = createSortableLoop(&arcset.arcs[arcIdx], &sloop);
    if (err) {
        destroyArcSet(&arcset);
        return err;
    }
    *rootOut = sloop.root;
    *areaOut = sloop.area;
    *numVertsOut = sloop.loop.numVerts;
    for (int64_t i = 0; i < sloop.loop.numVerts; i++) verts[i] = sloop.loop.verts[i];
    H3_MEMORY(free)(sloop.loop.verts);
    destroyArcSet(&arcset);
    return E_SUCCESS;
}

// The three comparators, directly on synthetic values.
int h3goTest_cmp_SortableLoop(H3Index rootA, double areaA, H3Index rootB,
                              double areaB) {
    SortableLoop a = {.root = rootA, .area = areaA};
    SortableLoop b = {.root = rootB, .area = areaB};
    return cmp_SortableLoop(&a, &b);
}

int h3goTest_cmp_SortablePoly(double areaA, double areaB) {
    SortablePoly a = {.outerArea = areaA};
    SortablePoly b = {.outerArea = areaB};
    return cmp_SortablePoly(&a, &b);
}

int h3goTest_cmp_uint64(H3Index a, H3Index b) { return cmp_uint64(&a, &b); }

// Destroy-helper observable state. Bit 1: fields nulled after the
// first call; bit 2: a second call is safe and leaves them nulled.
int h3goTest_destroyArcSet_state(const H3Index *cells, int64_t numCells) {
    ArcSet arcset;
    if (createArcSet(cells, numCells, &arcset)) return -1;
    destroyArcSet(&arcset);
    int bits = (arcset.arcs == NULL && arcset.buckets == NULL) ? 1 : 0;
    destroyArcSet(&arcset);
    if (arcset.arcs == NULL && arcset.buckets == NULL) bits |= 2;
    return bits;
}

int h3goTest_destroyLoopSet_state(const H3Index *cells, int64_t numCells,
                                  int shallow) {
    ArcSet arcset;
    if (createArcSet(cells, numCells, &arcset)) return -1;
    if (cancelArcPairs(arcset)) return -1;
    SortableLoopSet loopset;
    if (createSortableLoopSet(arcset, &loopset)) {
        destroyArcSet(&arcset);
        return -1;
    }
    if (shallow) {
        // Free the vertex arrays first so the shallow variant leaks
        // nothing in this harness.
        for (int64_t i = 0; i < loopset.numLoops; i++) {
            H3_MEMORY(free)(loopset.sloops[i].loop.verts);
        }
        destroySortableLoopSetShallow(&loopset);
    } else {
        destroySortableLoopSet(&loopset);
    }
    int bits = (loopset.sloops == NULL) ? 1 : 0;
    if (shallow) {
        destroySortableLoopSetShallow(&loopset);
    } else {
        destroySortableLoopSet(&loopset);
    }
    if (loopset.sloops == NULL) bits |= 2;
    destroyArcSet(&arcset);
    return bits;
}

// polygon.h externs defined in algos.c.
void destroyGeoLoop(GeoLoop *loop);
void destroyGeoPolygon(GeoPolygon *poly);

int h3goTest_destroyGeoLoop_state(void) {
    GeoLoop loop;
    loop.numVerts = 3;
    loop.verts = H3_MEMORY(malloc)(3 * sizeof(LatLng));
    if (!loop.verts) return -1;
    destroyGeoLoop(&loop);
    int bits = (loop.verts == NULL && loop.numVerts == 0) ? 1 : 0;
    destroyGeoLoop(&loop);
    if (loop.verts == NULL && loop.numVerts == 0) bits |= 2;
    return bits;
}

int h3goTest_destroyGeoPolygon_state(void) {
    GeoPolygon poly;
    poly.geoloop.numVerts = 3;
    poly.geoloop.verts = H3_MEMORY(malloc)(3 * sizeof(LatLng));
    poly.numHoles = 1;
    poly.holes = H3_MEMORY(malloc)(sizeof(GeoLoop));
    if (!poly.geoloop.verts || !poly.holes) return -1;
    poly.holes[0].numVerts = 3;
    poly.holes[0].verts = H3_MEMORY(malloc)(3 * sizeof(LatLng));
    if (!poly.holes[0].verts) return -1;
    destroyGeoPolygon(&poly);
    int bits = (poly.geoloop.verts == NULL && poly.geoloop.numVerts == 0 &&
                poly.holes == NULL && poly.numHoles == 0)
                   ? 1
                   : 0;
    destroyGeoPolygon(&poly);
    if (poly.geoloop.verts == NULL && poly.holes == NULL) bits |= 2;
    return bits;
}

int h3goTest_destroyGeoMultiPolygon_state(const H3Index *cells,
                                          int64_t numCells) {
    GeoMultiPolygon mpoly;
    if (H3_EXPORT(cellsToMultiPolygon)(cells, numCells, &mpoly)) return -1;
    H3_EXPORT(destroyGeoMultiPolygon)(&mpoly);
    int bits = (mpoly.polygons == NULL && mpoly.numPolygons == 0) ? 1 : 0;
    H3_EXPORT(destroyGeoMultiPolygon)(&mpoly);
    if (mpoly.polygons == NULL && mpoly.numPolygons == 0) bits |= 2;
    return bits;
}

// destroyLinkedMultiPolygon idempotence (H3 4.5.0 delta): destroy the
// linked output twice; bit 1 = head zeroed after the first call, bit 2
// = the second call was reached safely and the head stayed zeroed.
int h3goTest_destroyLinkedTwice(const H3Index *cells, int numCells) {
    LinkedGeoPolygon linked;
    if (H3_EXPORT(cellsToLinkedMultiPolygon)(cells, numCells, &linked)) {
        return -1;
    }
    H3_EXPORT(destroyLinkedMultiPolygon)(&linked);
    int bits = (linked.first == NULL && linked.last == NULL &&
                linked.next == NULL)
                   ? 1
                   : 0;
    H3_EXPORT(destroyLinkedMultiPolygon)(&linked);
    if (linked.first == NULL && linked.last == NULL && linked.next == NULL) {
        bits |= 2;
    }
    return bits;
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
