package h3

// iterCellsChildren mirrors the C iterCellsChildren struct for iterator state.
type iterCellsChildren struct {
	H         h3Index // Current H3 index
	ParentRes int32   // Parent resolution
	SkipDigit int32   // Skip digit for pentagons
}

// iterCellsPolygonCompact mirrors the C iterCellsPolygonCompact struct for iterating
// through all cells within a given polygon, outputting a compact set.
type iterCellsPolygonCompact struct {
	Cell    h3Index     // current value
	Error   h3Error     // error, if any
	res     int32       // target resolution
	flags   uint32      // Mode flags for the polygonToCells operation
	polygon *GeoPolygon // the polygon we're filling
	bboxes  []bbox      // Bounding box(es) for the polygon and its holes
	started bool        // Whether iteration has started
}

// iterCellsPolygon mirrors the C iterCellsPolygon struct for iterating through
// all cells within a given polygon at a fixed resolution.
type iterCellsPolygon struct {
	Cell      h3Index                 // current value
	Error     h3Error                 // error, if any
	cellIter  iterCellsPolygonCompact // sub-iterator for compact cells
	childIter iterCellsChildren       // sub-iterator for cell children
}

// iterCellsResolution mirrors the C iterCellsResolution struct for iterating
// through all cells at a given resolution.
type iterCellsResolution struct {
	H           h3Index           // current value
	baseCellNum int32             // current base cell number
	res         int32             // target resolution
	itC         iterCellsChildren // child iterator for current base cell
}
