package main

type scenario struct {
	name        string
	label       string
	description string
}

// scenarios is the public catalog for the complete benchmark report. Keeping
// it explicit makes an added, removed, or renamed benchmark a documentation
// change instead of silently dropping it from the published tables.
var scenarios = []scenario{
	{"LatLngToCell/res=9", "`LatLngToCell` · res 9", "Coordinate to cell index"},
	{"CellToLatLng/res=9", "`Cell.LatLng` · res 9", "Cell index to center coordinate"},
	{"CellToBoundary/res=9", "`Cell.Boundary` · res 9", "Cell boundary vertices"},
	{"CellToParent/res=9to7", "`Cell.Parent` · res 9→7", "Parent cell lookup"},
	{"GridDistance", "`GridDistance`", "Grid distance between two cells"},
	{"IsNeighbor", "`Cell.IsNeighbor`", "Neighbor test for two cells"},
	{"CellArea/unit=km2", "`Cell.AreaKm2`", "Cell area in square kilometers"},
	{"ParseCell", "`ParseCell`", "Text to cell index"},
	{"CellToString", "`Cell.String`", "Cell index to text"},
	{"IsValidCell", "`Cell.IsValid`", "Cell validity check"},
	{"Resolution", "`Cell.Resolution`", "Cell resolution lookup"},
	{"DirectedEdges", "`Cell.DirectedEdges`", "All directed edges from a cell"},
	{"DirectedEdgeBoundary", "`DirectedEdge.Boundary`", "Directed-edge boundary vertices"},
	{"Vertexes", "`Cell.Vertexes`", "All vertex indexes of a cell"},
	{"VertexLatLng", "`Vertex.LatLng`", "Vertex index to coordinate"},
	{"Children/depth=1", "`Cell.Children` · depth 1", "Immediate child cells"},
	{"Children/depth=3", "`Cell.Children` · depth 3", "Child cells three levels down"},
	{"Children/depth=5", "`Cell.Children` · depth 5", "Child cells five levels down"},
	{"GridDisk/k=1", "`GridDisk` · k 1", "Disk around one origin"},
	{"GridDisk/k=5", "`GridDisk` · k 5", "Disk around one origin"},
	{"GridDisk/k=20", "`GridDisk` · k 20", "Large disk around one origin"},
	{"GridDiskDistances/k=5", "`GridDiskDistances` · k 5", "Disk grouped by grid distance"},
	{"GridRing/k=5", "`GridRing` · k 5", "Hollow ring around one origin"},
	{"GridPath", "`GridPath`", "Shortest grid path between two cells"},
	{"GridDisksUnsafe/origins=64/k=2", "`GridDisksUnsafe` · 64 origins, k 2", "Batched disks for 64 origins"},
	{"Compact/set=sf9", "`CompactCells` · 1,253 cells", "Compact a San Francisco resolution-9 set"},
	{"Uncompact/res=4to9", "`UncompactCells` · res 4→9", "Expand one coarse cell to resolution 9"},
	{"PolygonToCells/poly=sf/res=7", "`PolygonToCells` · SF, res 7", "Fill a San Francisco polygon"},
	{"PolygonToCells/poly=sf/res=9", "`PolygonToCells` · SF, res 9", "Fill a San Francisco polygon"},
	{"PolygonToCells/poly=sf/res=11", "`PolygonToCells` · SF, res 11", "Fill a San Francisco polygon"},
	{"CellsToMultiPolygon/n=331", "`CellsToMultiPolygon` · 331 cells", "Convert a cell set to polygon loops"},
	{"BatchLatLngToCell/n=10000", "`LatLngToCell` batch · 10,000", "Convert a deterministic coordinate batch"},
	{"ServiceWorkload/pts=256", "service workload · 256 points", "For each point: `LatLngToCell` (res 9), `GridDisk` (k 1), then `Parent` (res 7) for every disk cell"},
}

var scenarioByName = func() map[string]scenario {
	byName := make(map[string]scenario, len(scenarios))
	for _, scenario := range scenarios {
		byName[scenario.name] = scenario
	}
	return byName
}()

var memoryScenarios = []scenario{
	{"polyfill-large", "large polygon fill", "`PolygonToCells` (SF, res 11), about 61k cells per call; last result retained"},
	{"children-deep", "deep children", "`Children` (res 2→8), 117,649 cells per call; last result retained"},
	{"uncompact-res0-to-5", "large uncompact", "`UncompactCells` (all 122 res-0 cells → res 5), 2,050,854 cells; last result retained"},
	{"compact-large", "large compact", "`CompactCells` over that 2,050,854-cell set; input and last result retained"},
	{"scalar-1m", "one million scalar calls", "1,000,000 `LatLngToCell` (res 9) calls; nothing retained"},
	{"multipolygon-sf9", "SF multipolygon", "`CellsToMultiPolygon` over 1,253 cells; last result retained"},
	{"retained-polyfill", "retained polygon fills", "200 × `PolygonToCells` (SF, res 9); every result retained"},
}
