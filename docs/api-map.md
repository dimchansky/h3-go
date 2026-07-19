# H3 C → Go API map

Every public function of **H3 C v4.5.0**, mapped to this library's public
Go API. Use this page to find the Go name for a C function you already
know; full signatures and documentation are on
[pkg.go.dev](https://pkg.go.dev/github.com/dimchansky/h3-go).

How to read the table:

- **Go API** is the idiomatic entry point. Operations on a single index are
  methods on the typed indexes (`Cell`, `DirectedEdge`, `Vertex`);
  free-standing operations are package functions.
- **Additional Go forms** are additive variants this library layers on top
  of the C surface: caller-owned-buffer `Append*` functions, streaming
  `*Seq` iterators, and grouped-result variants.
- A long dash (—) marks a C function that has no direct Go counterpart on
  purpose because Go semantics absorb its job.

Related documents: the [comparison matrix](comparison-uber-h3-go.md) maps
the same 79 functions against the official uber/h3-go cgo binding with
status, semantics, allocation, and migration columns; the
[migration guide](migration-from-uber-h3-go.md) covers call-site changes
when switching from the binding.

The tables on this page are generated from
[comparison-uber-h3-go.csv](comparison-uber-h3-go.csv) by
[tools/ubercompare](../tools/README.md#ubercompare); this framing text is
maintained by hand. Regenerate with
`make gen-ubercompare`; `make check-ubercompare` fails in CI when it
drifts from the matrix, the C-API inventory, or the locked API surface.

<!-- BEGIN GENERATED: ubercompare api-map (edit docs/comparison-uber-h3-go.csv and run `make gen-ubercompare`) -->

All **79** public functions of the pinned H3 C release are mapped. A
long dash (—) in the Go API column marks a C function whose job is
absorbed by Go semantics — sentinel errors carry the message text, the
garbage collector replaces destructors, and one sizing helper stays
internal; the annotation says where that behavior lives.

### Indexing

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `latLngToCell` | `LatLngToCell`, `LatLng.Cell` | — |
| `cellToLatLng` | `Cell.LatLng` | — |
| `cellToBoundary` | `Cell.Boundary` | — |

### Index inspection

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `getResolution` | `Cell.Resolution`, `DirectedEdge.Resolution`, `Vertex.Resolution` | — |
| `getBaseCellNumber` | `Cell.BaseCellNumber` | — |
| `getIndexDigit` | `Cell.IndexDigit`, `DirectedEdge.IndexDigit`, `Vertex.IndexDigit` | — |
| `stringToH3` | `ParseCell`, `ParseDirectedEdge`, `ParseVertex`, `Cell.UnmarshalText` | — |
| `h3ToString` | `Cell.String`, `DirectedEdge.String`, `Vertex.String`, `Cell.MarshalText` | — |
| `isValidCell` | `Cell.IsValid` | — |
| `isValidIndex` | `IsValidIndex` | — |
| `isResClassIII` | `Cell.IsResClassIII` | — |
| `isPentagon` | `Cell.IsPentagon` | — |
| `maxFaceCount` | — (sizes IcosahedronFaces internally) | — |
| `getIcosahedronFaces` | `Cell.IcosahedronFaces` | — |
| `constructCell` | `ConstructCell` | — |

### Grid traversal

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `maxGridDiskSize` | `MaxGridDiskSize` | — |
| `gridDisk` | `Cell.GridDisk` | `Cell.AppendGridDisk` |
| `gridDiskUnsafe` | `Cell.GridDiskUnsafe` | `Cell.AppendGridDiskUnsafe` |
| `gridDiskDistances` | `Cell.GridDiskDistances` | `Cell.AppendGridDiskDistances`, `Cell.GridDiskDistancesGrouped` |
| `gridDiskDistancesSafe` | `Cell.GridDiskDistancesSafe` | — |
| `gridDiskDistancesUnsafe` | `Cell.GridDiskDistancesUnsafe` | — |
| `gridDisksUnsafe` | `GridDisksUnsafe` | — |
| `maxGridRingSize` | `MaxGridRingSize` | — |
| `gridRing` | `Cell.GridRing` | `Cell.AppendGridRing` |
| `gridRingUnsafe` | `Cell.GridRingUnsafe` | `Cell.AppendGridRingUnsafe` |
| `gridDistance` | `Cell.GridDistance` | — |
| `gridPathCellsSize` | `Cell.GridPathLen` | — |
| `gridPathCells` | `Cell.GridPath` | `Cell.AppendGridPath` |
| `cellToLocalIj` | `CellToLocalIJ` | — |
| `localIjToCell` | `LocalIJToCell` | — |

### Hierarchy and compaction

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `cellToParent` | `Cell.Parent`, `Cell.ImmediateParent` | — |
| `cellToChildrenSize` | `Cell.NumChildren` | — |
| `cellToChildren` | `Cell.Children`, `Cell.ImmediateChildren` | `Cell.AppendChildren`, `Cell.ChildrenSeq`, `Cell.AppendImmediateChildren` |
| `cellToCenterChild` | `Cell.CenterChild` | — |
| `cellToChildPos` | `Cell.ChildPos` | — |
| `childPosToCell` | `Cell.ChildAtPos` | — |
| `compactCells` | `CompactCells` | `AppendCompactCells` |
| `uncompactCellsSize` | `UncompactCellsSize` | — |
| `uncompactCells` | `UncompactCells` | `AppendUncompactCells` |

### Directed edges

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `areNeighborCells` | `Cell.IsNeighbor` | — |
| `cellsToDirectedEdge` | `Cell.DirectedEdgeTo` | — |
| `isValidDirectedEdge` | `DirectedEdge.IsValid` | — |
| `getDirectedEdgeOrigin` | `DirectedEdge.Origin` | — |
| `getDirectedEdgeDestination` | `DirectedEdge.Destination` | — |
| `directedEdgeToCells` | `DirectedEdge.Cells` | — |
| `originToDirectedEdges` | `Cell.DirectedEdges` | — |
| `directedEdgeToBoundary` | `DirectedEdge.Boundary` | — |
| `reverseDirectedEdge` | `DirectedEdge.Reverse` | — |

### Vertexes

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `cellToVertex` | `Cell.Vertex` | — |
| `cellToVertexes` | `Cell.Vertexes` | `Cell.AppendVertexes` |
| `vertexToLatLng` | `Vertex.LatLng` | — |
| `isValidVertex` | `Vertex.IsValid` | — |

### Regions (polyfill and multi-polygon)

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `maxPolygonToCellsSize` | `MaxPolygonToCellsSize` | — |
| `polygonToCells` | `PolygonToCells` | `AppendPolygonToCells` |
| `maxPolygonToCellsSizeExperimental` | `MaxPolygonToCellsSizeExperimental` | — |
| `polygonToCellsExperimental` | `PolygonToCellsExperimental` | `AppendPolygonToCellsExperimental`, `PolygonToCellsExperimentalSeq` |
| `cellsToLinkedMultiPolygon` | `CellsToMultiPolygon` | — |
| `destroyLinkedMultiPolygon` | — (garbage collected) | — |

### Measurement

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `cellAreaRads2` | `Cell.AreaRads2` | — |
| `cellAreaKm2` | `Cell.AreaKm2` | — |
| `cellAreaM2` | `Cell.AreaM2` | — |
| `edgeLengthRads` | `DirectedEdge.LengthRads` | — |
| `edgeLengthKm` | `DirectedEdge.LengthKm` | — |
| `edgeLengthM` | `DirectedEdge.LengthM` | — |
| `getHexagonAreaAvgKm2` | `HexagonAreaAvgKm2` | — |
| `getHexagonAreaAvgM2` | `HexagonAreaAvgM2` | — |
| `getHexagonEdgeLengthAvgKm` | `HexagonEdgeLengthAvgKm` | — |
| `getHexagonEdgeLengthAvgM` | `HexagonEdgeLengthAvgM` | — |
| `greatCircleDistanceRads` | `GreatCircleDistanceRads` | — |
| `greatCircleDistanceKm` | `GreatCircleDistanceKm` | — |
| `greatCircleDistanceM` | `GreatCircleDistanceM` | — |

### Constants, conversions, and error description

| H3 C API | Go API | Additional Go forms |
|---|---|---|
| `describeH3Error` | — (sentinel error messages) | — |
| `degsToRads` | `Deg`, `Angle.Rad` | — |
| `radsToDegs` | `Rad`, `Angle.Deg` | — |
| `getNumCells` | `NumCells` | — |
| `res0CellCount` | `NumRes0Cells` | — |
| `getRes0Cells` | `Res0Cells` | — |
| `pentagonCount` | `NumPentagons` | — |
| `getPentagons` | `Pentagons` | — |

<!-- END GENERATED: ubercompare api-map -->
