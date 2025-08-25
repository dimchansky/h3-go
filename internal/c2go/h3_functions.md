# H3 C Library Functions - Port Status

This document tracks the porting status of all H3 C public functions exported via `H3_EXPORT()` macro.

## algos.c (19 functions)

- [ ] cellsToLinkedMultiPolygon
- [x] cellToBoundary
- [x] cellToLatLng
- [ ] destroyLinkedMultiPolygon
- [x] getNumCells
- [x] gridDisk
- [x] gridDiskDistances
- [x] gridDiskDistancesSafe
- [x] gridDiskDistancesUnsafe
- [ ] gridDisksUnsafe
- [ ] gridDiskUnsafe
- [ ] gridRing
- [ ] gridRingUnsafe
- [x] isPentagon
- [x] latLngToCell
- [x] maxGridDiskSize
- [x] maxGridRingSize
- [ ] maxPolygonToCellsSize
- [ ] polygonToCells

## baseCells.c (2 functions)

- [x] getRes0Cells
- [x] res0CellCount

## directedEdge.c (8 functions)

- [x] areNeighborCells
- [x] cellsToDirectedEdge
- [x] directedEdgeToBoundary
- [x] directedEdgeToCells
- [x] getDirectedEdgeDestination
- [x] getDirectedEdgeOrigin
- [x] isValidDirectedEdge
- [x] originToDirectedEdges

## h3Index.c (22+ functions)

- [x] cellToBoundary
- [x] cellToCenterChild
- [x] cellToChildPos
- [x] cellToChildren
- [x] cellToChildrenSize
- [x] cellToLatLng
- [x] cellToParent
- [x] childPosToCell
- [x] compactCells
- [x] describeH3Error
- [x] getBaseCellNumber
- [x] getIcosahedronFaces
- [x] getPentagons
- [x] getResolution
- [x] h3ToString
- [x] isPentagon
- [x] isResClassIII
- [x] isValidCell
- [x] latLngToCell
- [x] maxFaceCount
- [x] pentagonCount
- [x] stringToH3
- [x] uncompactCells
- [x] uncompactCellsSize

## latLng.c (19 functions)

- [x] cellAreaKm2
- [x] cellAreaM2
- [x] cellAreaRads2
- [x] degsToRads
- [x] radsToDegs
- [x] edgeLengthKm
- [x] edgeLengthM
- [x] edgeLengthRads
- [x] getHexagonAreaAvgKm2
- [x] getHexagonAreaAvgM2
- [x] getHexagonEdgeLengthAvgKm
- [x] getHexagonEdgeLengthAvgM
- [x] getNumCells
- [x] greatCircleDistanceKm
- [x] greatCircleDistanceM
- [x] greatCircleDistanceRads

## localij.c (5 functions)

- [x] cellToLocalIj
- [x] localIjToCell
- [x] gridDistance
- [x] gridPathCells
- [x] gridPathCellsSize

## vertex.c (7 functions)

- [x] cellToVertex
- [x] cellToVertexes
- [x] isValidVertex
- [x] vertexToLatLng

## Summary

**Total Functions:** 82
- **Completed:** 70
- **Remaining:** 8

## Notes

- Functions are listed exactly as they appear in the H3 C source code
- Some function names appear in multiple files (e.g., `cellToLatLng`, `isPentagon`) - these may be internal calls rather than duplicated exports
- Files like coordijk.c, faceijk.c contain internal functions but no H3_EXPORT public API functions
- This list represents the complete H3 v4.3.0 public API surface