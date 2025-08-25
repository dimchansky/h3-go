# Ported H3 Internal Functions Report

This report lists all H3 C functions that have been ported to Go in the internal/c2go package.
These are primarily internal/static functions used to build the public API.

## algos.c (15 functions)

- DIRECTIONS
- K_ALL_CELLS_AT_RES_15
- NEXT_RING_DIRECTION
- _gridDiskDistancesInternal
- _gridRingInternal
- directionForNeighbor
- gridDisk
- gridDiskDistances
- gridDiskDistancesSafe
- gridDiskDistancesUnsafe
- gridRing
- gridRingUnsafe
- h3NeighborRotations
- maxGridDiskSize
- maxGridRingSize

## baseCells.c (15 functions)

- _baseCellIsCwOffset
- _baseCellToCCWrot60
- _baseCellToFaceIjk
- _faceIjkToBaseCell
- _faceIjkToBaseCellCCWrot60
- _getBaseCellDirection
- _getBaseCellNeighbor
- _isBaseCellPentagon
- _isBaseCellPolarPentagon
- baseCellData
- baseCellNeighbor60CCWRots
- baseCellNeighbors
- faceIjkBaseCells
- getRes0Cells
- res0CellCount

## bbox.c (10 functions)

- bboxCenter
- bboxContains
- bboxContainsBBox
- bboxEquals
- bboxHeightRads
- bboxIsTransmeridian
- bboxNormalization
- bboxOverlapsBBox
- bboxWidthRads
- scaleBBox

## coordijk.c (28 functions)

- _downAp3
- _downAp3r
- _downAp7
- _downAp7r
- _hex2dToCoordIJK
- _ijkAdd
- _ijkMatches
- _ijkNormalize
- _ijkNormalizeCouldOverflow
- _ijkRotate60ccw
- _ijkRotate60cw
- _ijkScale
- _ijkSub
- _ijkToHex2d
- _neighbor
- _rotate60ccw
- _rotate60cw
- _setIJK
- _unitIjkToDigit
- _upAp7
- _upAp7Checked
- _upAp7r
- _upAp7rChecked
- cubeToIjk
- ijToIjk
- ijkDistance
- ijkToCube
- ijkToIj

## directedEdge.c (8 functions)

- areNeighborCells
- cellsToDirectedEdge
- directedEdgeToBoundary
- directedEdgeToCells
- getDirectedEdgeDestination
- getDirectedEdgeOrigin
- isValidDirectedEdge
- originToDirectedEdges

## faceijk.c (15 functions)

- _adjustOverageClassII
- _adjustPentVertOverage
- _faceIjkPentToCellBoundary
- _faceIjkPentToVerts
- _faceIjkToCellBoundary
- _faceIjkToGeo
- _faceIjkToVerts
- _geoToClosestFace
- _geoToFaceIjk
- _geoToHex2d
- _hex2dToGeo
- adjacentFaceDir
- faceNeighbors
- maxDimByCIIres
- unitScaleByCIIres

## h3Index.c (44 functions)

- H3_GET_BASE_CELL
- _faceIjkToH3
- _firstOneIndex
- _h3LeadingNonZeroDigit
- _h3Rotate60ccw
- _h3Rotate60cw
- _h3RotatePent60ccw
- _h3RotatePent60cw
- _h3ToFaceIjk
- _h3ToFaceIjkWithInitializedFijk
- _hasAll7AfterRes
- _hasAny7UptoRes
- _hasChildAtRes
- _hasDeletedSubsequence
- _hasGoodTopBits
- _zeroIndexDigits
- cellToBoundary
- cellToCenterChild
- cellToChildPos
- cellToChildren
- cellToChildrenSize
- cellToLatLng
- cellToParent
- childPosToCell
- compactCells
- describeH3Error
- getBaseCellNumber
- getIcosahedronFaces
- getPentagons
- getResolution
- h3ToString
- isPentagon
- isResClassIII
- isResolutionClassIII
- isValidCell
- latLngToCell
- makeDirectChild
- maxFaceCount
- pentagonCount
- setH3Index
- stringToH3
- uncompactCells
- uncompactCellsSize
- validateChildPos

## iterators.c (4 functions)

- _incrementResDigit
- _iterInitParent
- _null_iter
- iterStepChild

## latLng.c (28 functions)

- _geoAzDistanceRads
- _geoAzimuthRads
- _posAngleRads
- _setGeoRads
- cellAreaKm2
- cellAreaM2
- cellAreaRads2
- constrainLat
- constrainLng
- degsToRads
- edgeLengthKm
- edgeLengthM
- edgeLengthRads
- geoAlmostEqual
- geoAlmostEqualThreshold
- getHexagonAreaAvgKm2
- getHexagonAreaAvgM2
- getHexagonEdgeLengthAvgKm
- getHexagonEdgeLengthAvgM
- getNumCells
- greatCircleDistanceKm
- greatCircleDistanceM
- greatCircleDistanceRads
- normalizeLng
- radsToDegs
- setGeoDegs
- triangleArea
- triangleEdgeLengthsToArea

## linkedGeo.c (11 functions)

- addLinkedCoord
- addLinkedLoop
- addNewLinkedLoop
- addNewLinkedPolygon
- countContainers
- countLinkedCoords
- countLinkedLoops
- countLinkedPolygons
- findDeepestContainer
- findPolygonForHole
- normalizeMultiPolygon

## localij.c (13 functions)

- FAILED_DIRECTIONS
- PENTAGON_ROTATIONS
- PENTAGON_ROTATIONS_REVERSE
- PENTAGON_ROTATIONS_REVERSE_NONPOLAR
- PENTAGON_ROTATIONS_REVERSE_POLAR
- cellToLocalIj
- cellToLocalIjk
- cubeRound
- gridDistance
- gridPathCells
- gridPathCellsSize
- localIjToCell
- localIjkToCell

## mathExtensions.c (1 functions)

- _ipow

## polyfill.c (1 functions)

- baseCellNumToCell

## polygon.c (9 functions)

- bboxFromGeoLoop
- bboxesFromGeoPolygon
- cellBoundaryCrossesGeoLoop
- cellBoundaryCrossesPolygon
- cellBoundaryInsidePolygon
- lineCrossesLine
- pointInsideGeoLoop
- pointInsidePolygon
- validatePolygonFlags

## vec2d.c (3 functions)

- _v2dAlmostEquals
- _v2dIntersect
- _v2dMag

## vec3d.c (3 functions)

- _geoToVec3d
- _pointSquareDist
- _square

## vertex.c (15 functions)

- DIRECTIONS
- DIRECTION_INDEX_OFFSET
- cellToVertex
- cellToVertexes
- directionForVertexNum
- directionToVertexNumHex
- directionToVertexNumPent
- isValidVertex
- pentagonDirectionFaces
- revNeighborDirectionsHex
- vertexNumForDirection
- vertexNumToDirectionHex
- vertexNumToDirectionPent
- vertexRotations
- vertexToLatLng

## vertexGraph.c (2 functions)

- _hashVertex
- initVertexGraph

**Total ported internal functions: 225**

---

*This report is automatically generated. Run `./scripts/update-h3-status.sh` to refresh.*
