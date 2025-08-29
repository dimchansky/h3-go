package h3

// IterCellsChildren mirrors the C IterCellsChildren struct for iterator state.
type IterCellsChildren struct {
	H         H3Index // Current H3 index
	ParentRes int32   // Parent resolution
	SkipDigit int32   // Skip digit for pentagons
}

// IterCellsPolygonCompact mirrors the C IterCellsPolygonCompact struct for iterating
// through all cells within a given polygon, outputting a compact set.
type IterCellsPolygonCompact struct {
	Cell    H3Index     // current value
	Error   H3Error     // error, if any
	res     int32       // target resolution
	flags   uint32      // Mode flags for the polygonToCells operation
	polygon *GeoPolygon // the polygon we're filling
	bboxes  []BBox      // Bounding box(es) for the polygon and its holes
	started bool        // Whether iteration has started
}

// IterCellsPolygon mirrors the C IterCellsPolygon struct for iterating through
// all cells within a given polygon at a fixed resolution.
type IterCellsPolygon struct {
	Cell      H3Index                 // current value
	Error     H3Error                 // error, if any
	cellIter  IterCellsPolygonCompact // sub-iterator for compact cells
	childIter IterCellsChildren       // sub-iterator for cell children
}

// IterCellsResolution mirrors the C IterCellsResolution struct for iterating
// through all cells at a given resolution.
type IterCellsResolution struct {
	H           H3Index           // current value
	baseCellNum int32             // current base cell number
	res         int32             // target resolution
	itC         IterCellsChildren // child iterator for current base cell
}
