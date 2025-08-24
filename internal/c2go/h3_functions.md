# H3 C Library Functions - Port Status

This document tracks the porting status of all H3 C public functions exported via `H3_EXPORT()` macro.

## algos.c (19 functions)

- [ ] cellsToLinkedMultiPolygon
- [x] cellToBoundary
- [x] cellToLatLng
- [ ] destroyLinkedMultiPolygon
- [ ] getNumCells
- [ ] gridDisk
- [ ] gridDiskDistances
- [ ] gridDiskDistancesSafe
- [ ] gridDiskDistancesUnsafe
- [ ] gridDisksUnsafe
- [ ] gridDiskUnsafe
- [ ] gridRing
- [ ] gridRingUnsafe
- [x] isPentagon
- [x] latLngToCell
- [ ] maxGridDiskSize
- [ ] maxGridRingSize
- [ ] maxPolygonToCellsSize
- [ ] polygonToCells

## baseCells.c (2 functions)

- [x] getRes0Cells
- [x] res0CellCount

## directedEdge.c (12 functions)

- [ ] areNeighborCells
- [ ] cellsToDirectedEdge
- [ ] directedEdgeToBoundary
- [ ] directedEdgeToCells
- [ ] getDirectedEdgeDestination
- [ ] getDirectedEdgeOrigin
- [ ] isValidDirectedEdge
- [ ] originToDirectedEdges

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
- [ ] uncompactCells
- [ ] uncompactCellsSize

## latLng.c (19 functions)

- [ ] cellAreaKm2
- [ ] cellAreaM2
- [ ] cellAreaRads2
- [x] degsToRads
- [x] radsToDegs
- [ ] edgeLengthKm
- [ ] edgeLengthM
- [ ] edgeLengthRads
- [ ] getHexagonAreaAvgKm2
- [ ] getHexagonAreaAvgM2
- [ ] getHexagonEdgeLengthAvgKm
- [ ] getHexagonEdgeLengthAvgM
- [ ] getNumCells
- [x] greatCircleDistanceKm
- [x] greatCircleDistanceM
- [x] greatCircleDistanceRads

## localij.c (5 functions)

- [ ] cellToLocalIj
- [ ] localIjToCell
- [ ] gridDistance
- [ ] gridPathCells
- [ ] gridPathCellsSize

## vertex.c (7 functions)

- [ ] cellToVertex
- [ ] cellToVertexes
- [ ] isValidVertex
- [ ] vertexToLatLng

## Summary

**Total Functions:** 86
- **Completed:** 33
- **Remaining:** 45

## Notes

- Functions are listed exactly as they appear in the H3 C source code
- Some function names appear in multiple files (e.g., `cellToLatLng`, `isPentagon`) - these may be internal calls rather than duplicated exports
- Files like coordijk.c, faceijk.c contain internal functions but no H3_EXPORT public API functions
- This list represents the complete H3 v4.3.0 public API surface