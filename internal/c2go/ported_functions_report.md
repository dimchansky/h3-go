# Ported H3 Internal Functions Report

This report lists all H3 C functions that have been ported to Go in the internal/c2go package.
These are primarily internal/static functions used to build the public API.

## algos.c (2 functions)

- directionForNeighbor
- h3NeighborRotations

## baseCells.c (15 functions)

- _baseCellIsCwOffset
- _baseCellToCCWrot60
- _baseCellToFaceIjk
- _faceIjkToBaseCell
- _faceIjkToBaseCellCCWrot60
- _getBaseCellDirection
- _getBaseCellNeighbor
- _isBaseCellPolarPentagon
- baseCellData
- baseCellNeighbor60CCWRots
- baseCellNeighbors
- faceIjkBaseCells
- getRes0Cells
- isBaseCellPentagon
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

- _rotate60cw
- cubeToIjk
- downAp3
- downAp3r
- downAp7
- downAp7r
- hex2dToCoordIJK
- ijToIjk
- ijkAdd
- ijkDistance
- ijkMatches
- ijkNormalize
- ijkNormalizeCouldOverflow
- ijkRotate60ccw
- ijkRotate60cw
- ijkScale
- ijkSub
- ijkToCube
- ijkToHex2d
- ijkToIj
- neighbor
- rotate60ccw
- setIJK
- unitIjkToDigit
- upAp7
- upAp7Checked
- upAp7r
- upAp7rChecked

## directedEdge.c (4 functions)

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
- _geoToFaceIjk
- adjacentFaceDir
- faceNeighbors
- geoToClosestFace
- geoToHex2d
- hex2dToGeo
- maxDimByCIIres
- unitScaleByCIIres

## h3Index.c (52 functions)

- _faceIjkToH3
- _firstOneIndex
- _getResDigit
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
- getHighBit
- getIcosahedronFaces
- getIndexDigit
- getMode
- getPentagons
- getReservedBits
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
- setHighBit
- setIndexDigit
- setMode
- setReservedBits
- stringToH3
- uncompactCells
- uncompactCellsSize
- validateChildPos

## iterators.c (4 functions)

- _incrementResDigit
- _iterInitParent
- _null_iter
- iterStepChild

## latLng.c (21 functions)

- _geoAzDistanceRads
- _geoAzimuthRads
- _setGeoRads
- cellAreaKm2
- cellAreaM2
- cellAreaRads2
- constrainLat
- constrainLng
- degsToRads
- geoAlmostEqual
- geoAlmostEqualThreshold
- getNumCells
- greatCircleDistanceKm
- greatCircleDistanceM
- greatCircleDistanceRads
- normalizeLng
- posAngleRads
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

## localij.c (1 functions)

- cubeRound

## mathExtensions.c (1 functions)

- ipow

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
- _v2dMag
- v2dIntersect

## vec3d.c (3 functions)

- _geoToVec3d
- _square
- pointSquareDist

## vertex.c (11 functions)

- DIRECTIONS
- cellToVertex
- directionForVertexNum
- directionToVertexNumHex
- directionToVertexNumPent
- pentagonDirectionFaces
- revNeighborDirectionsHex
- vertexNumForDirection
- vertexNumToDirectionHex
- vertexNumToDirectionPent
- vertexRotations

## vertexGraph.c (2 functions)

- _hashVertex
- initVertexGraph

**Total ported internal functions: 193**

---

*This report is automatically generated. Run `./scripts/update-h3-status.sh` to refresh.*
