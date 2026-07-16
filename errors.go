package h3

import "errors"

// Sentinel errors returned by the public API, one per H3 C error code
// (H3ErrorCodes in h3api.h). Match with errors.Is. The message text is the C
// library's describeH3Error text with an "h3: " prefix. The "returned by"
// notes below name representative operations; they are not exhaustive
// lists. Parse and unmarshal syntax errors are not sentinels — they wrap
// the underlying strconv error (see ParseCell).
var (
	// ErrFailed corresponds to eFailed: the operation failed but a more
	// specific error is not available. Returned by operations that cannot
	// complete for algorithmic reasons, e.g. GridDistance or CellToLocalIJ
	// across pentagon distortion or for very distant cells.
	ErrFailed = errors.New("h3: " + describeH3Error(eFailed))
	// ErrDomain corresponds to eDomain: argument was outside of acceptable
	// range. Returned by, among others, the k-taking traversal functions
	// for k outside 0..math.MaxInt32 and Cell.ChildAtPos for an
	// out-of-range position.
	ErrDomain = errors.New("h3: " + describeH3Error(eDomain))
	// ErrLatLngDomain corresponds to eLatlngDomain: latitude or longitude
	// arguments were outside of acceptable range. Returned by LatLngToCell
	// only for non-finite (NaN or infinite) coordinates; finite
	// out-of-range values are not validated.
	ErrLatLngDomain = errors.New("h3: " + describeH3Error(eLatlngDomain))
	// ErrResolutionDomain corresponds to eResDomain: resolution argument
	// was outside of acceptable range. Returned by resolution-taking
	// operations for res outside 0..MaxResolution (e.g. LatLngToCell,
	// Cell.Parent, Cell.IndexDigit) and by Cell.CenterChild for a
	// resolution coarser than the cell's.
	ErrResolutionDomain = errors.New("h3: " + describeH3Error(eResDomain))
	// ErrCellInvalid corresponds to eCellInvalid: cell argument was not
	// valid. Returned by ParseCell and Cell.UnmarshalText for a well-formed
	// string that is not a valid cell index, and by cell operations given
	// an invalid cell.
	ErrCellInvalid = errors.New("h3: " + describeH3Error(eCellInvalid))
	// ErrDirectedEdgeInvalid corresponds to eDirEdgeInvalid: directed edge
	// argument was not valid. Returned by ParseDirectedEdge and
	// DirectedEdge.UnmarshalText for a well-formed string that is not a
	// valid directed edge index, and by edge operations given an invalid
	// edge.
	ErrDirectedEdgeInvalid = errors.New("h3: " + describeH3Error(eDirEdgeInvalid))
	// ErrUndirectedEdgeInvalid corresponds to eUndirEdgeInvalid:
	// undirected edge argument was not valid. H3 4.x exposes no public
	// undirected-edge operations, so no operation in this package is
	// expected to return it.
	ErrUndirectedEdgeInvalid = errors.New("h3: " + describeH3Error(eUndirEdgeInvalid))
	// ErrVertexInvalid corresponds to eVertexInvalid: vertex argument was
	// not valid. Returned by ParseVertex and Vertex.UnmarshalText for a
	// well-formed string that is not a valid vertex index, and by vertex
	// operations given an invalid vertex.
	ErrVertexInvalid = errors.New("h3: " + describeH3Error(eVertexInvalid))
	// ErrPentagon corresponds to ePentagon: pentagon distortion was
	// encountered and the algorithm could not proceed. Returned by the
	// Unsafe traversal variants (GridDiskUnsafe, GridRingUnsafe,
	// GridDisksUnsafe, ...) when a pentagon or its distortion area is
	// encountered.
	ErrPentagon = errors.New("h3: " + describeH3Error(ePentagon))
	// ErrDuplicateInput corresponds to eDuplicateInput: duplicate input.
	// May be returned by CompactCells when duplicate input is detected;
	// detection is best-effort, not guaranteed (see CompactCells).
	ErrDuplicateInput = errors.New("h3: " + describeH3Error(eDuplicateInput))
	// ErrNotNeighbors corresponds to eNotNeighbors: cell arguments were not
	// neighbors. Returned by Cell.DirectedEdgeTo when the two cells are not
	// adjacent.
	ErrNotNeighbors = errors.New("h3: " + describeH3Error(eNotNeighbors))
	// ErrResolutionMismatch corresponds to eResMismatch: cell arguments had
	// incompatible resolutions. Returned by, among others, GridDistance,
	// GridPath, and CellToLocalIJ for cells of different resolutions,
	// UncompactCells for input cells finer than the target resolution, and
	// Cell.Parent for a requested resolution finer than the cell's.
	ErrResolutionMismatch = errors.New("h3: " + describeH3Error(eResMismatch))
	// ErrMemoryAlloc corresponds to eMemoryAlloc: memory allocation failed.
	// The Go port sizes its own buffers, so public operations are not
	// expected to return it.
	ErrMemoryAlloc = errors.New("h3: " + describeH3Error(eMemoryAlloc))
	// ErrMemoryBounds corresponds to eMemoryBounds: bounds of provided
	// memory were insufficient. The Go port sizes its own buffers, so
	// public operations are not expected to return it.
	ErrMemoryBounds = errors.New("h3: " + describeH3Error(eMemoryBounds))
	// ErrOptionInvalid corresponds to eOptionInvalid: mode or flags
	// argument was not valid. Returned by the PolygonToCellsExperimental
	// family for an invalid ContainmentMode.
	ErrOptionInvalid = errors.New("h3: " + describeH3Error(eOptionInvalid))
	// ErrIndexInvalid corresponds to eIndexInvalid: index argument was not
	// valid (H3 C 4.4.0).
	ErrIndexInvalid = errors.New("h3: " + describeH3Error(eIndexInvalid))
	// ErrBaseCellDomain corresponds to eBaseCellDomain: base cell number was
	// outside of acceptable range (H3 C 4.4.0).
	ErrBaseCellDomain = errors.New("h3: " + describeH3Error(eBaseCellDomain))
	// ErrDigitDomain corresponds to eDigitDomain: child digits invalid
	// (H3 C 4.4.0).
	ErrDigitDomain = errors.New("h3: " + describeH3Error(eDigitDomain))
	// ErrDeletedDigit corresponds to eDeletedDigit: deleted subsequence
	// indicates invalid index (H3 C 4.4.0).
	ErrDeletedDigit = errors.New("h3: " + describeH3Error(eDeletedDigit))
)

// errTable maps H3 C error codes to the public sentinel errors. Index 0
// (eSuccess) is nil.
var errTable = [h3ErrorEnd]error{
	eSuccess:          nil,
	eFailed:           ErrFailed,
	eDomain:           ErrDomain,
	eLatlngDomain:     ErrLatLngDomain,
	eResDomain:        ErrResolutionDomain,
	eCellInvalid:      ErrCellInvalid,
	eDirEdgeInvalid:   ErrDirectedEdgeInvalid,
	eUndirEdgeInvalid: ErrUndirectedEdgeInvalid,
	eVertexInvalid:    ErrVertexInvalid,
	ePentagon:         ErrPentagon,
	eDuplicateInput:   ErrDuplicateInput,
	eNotNeighbors:     ErrNotNeighbors,
	eResMismatch:      ErrResolutionMismatch,
	eMemoryAlloc:      ErrMemoryAlloc,
	eMemoryBounds:     ErrMemoryBounds,
	eOptionInvalid:    ErrOptionInvalid,
	eIndexInvalid:     ErrIndexInvalid,
	eBaseCellDomain:   ErrBaseCellDomain,
	eDigitDomain:      ErrDigitDomain,
	eDeletedDigit:     ErrDeletedDigit,
}

// toErr converts an internal H3 C error code to the public sentinel error.
// eSuccess maps to nil; codes outside the known range map to ErrFailed.
func toErr(errC h3Error) error {
	if int(errC) < len(errTable) {
		return errTable[errC]
	}
	return ErrFailed
}
