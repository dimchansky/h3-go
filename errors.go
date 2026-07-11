package h3

import "errors"

// Sentinel errors returned by the public API, one per H3 C error code
// (H3ErrorCodes in h3api.h). Match with errors.Is. The message text is the C
// library's describeH3Error text with an "h3: " prefix.
var (
	// ErrFailed corresponds to E_FAILED: the operation failed but a more
	// specific error is not available.
	ErrFailed = errors.New("h3: " + describeH3Error(E_FAILED))
	// ErrDomain corresponds to E_DOMAIN: argument was outside of acceptable range.
	ErrDomain = errors.New("h3: " + describeH3Error(E_DOMAIN))
	// ErrLatLngDomain corresponds to E_LATLNG_DOMAIN: latitude or longitude
	// arguments were outside of acceptable range.
	ErrLatLngDomain = errors.New("h3: " + describeH3Error(E_LATLNG_DOMAIN))
	// ErrResolutionDomain corresponds to E_RES_DOMAIN: resolution argument
	// was outside of acceptable range.
	ErrResolutionDomain = errors.New("h3: " + describeH3Error(E_RES_DOMAIN))
	// ErrCellInvalid corresponds to E_CELL_INVALID: cell argument was not valid.
	ErrCellInvalid = errors.New("h3: " + describeH3Error(E_CELL_INVALID))
	// ErrDirectedEdgeInvalid corresponds to E_DIR_EDGE_INVALID: directed edge
	// argument was not valid.
	ErrDirectedEdgeInvalid = errors.New("h3: " + describeH3Error(E_DIR_EDGE_INVALID))
	// ErrUndirectedEdgeInvalid corresponds to E_UNDIR_EDGE_INVALID:
	// undirected edge argument was not valid.
	ErrUndirectedEdgeInvalid = errors.New("h3: " + describeH3Error(E_UNDIR_EDGE_INVALID))
	// ErrVertexInvalid corresponds to E_VERTEX_INVALID: vertex argument was
	// not valid.
	ErrVertexInvalid = errors.New("h3: " + describeH3Error(E_VERTEX_INVALID))
	// ErrPentagon corresponds to E_PENTAGON: pentagon distortion was
	// encountered and the algorithm could not proceed.
	ErrPentagon = errors.New("h3: " + describeH3Error(E_PENTAGON))
	// ErrDuplicateInput corresponds to E_DUPLICATE_INPUT: duplicate input.
	ErrDuplicateInput = errors.New("h3: " + describeH3Error(E_DUPLICATE_INPUT))
	// ErrNotNeighbors corresponds to E_NOT_NEIGHBORS: cell arguments were not
	// neighbors.
	ErrNotNeighbors = errors.New("h3: " + describeH3Error(E_NOT_NEIGHBORS))
	// ErrResolutionMismatch corresponds to E_RES_MISMATCH: cell arguments had
	// incompatible resolutions.
	ErrResolutionMismatch = errors.New("h3: " + describeH3Error(E_RES_MISMATCH))
	// ErrMemoryAlloc corresponds to E_MEMORY_ALLOC: memory allocation failed.
	ErrMemoryAlloc = errors.New("h3: " + describeH3Error(E_MEMORY_ALLOC))
	// ErrMemoryBounds corresponds to E_MEMORY_BOUNDS: bounds of provided
	// memory were insufficient.
	ErrMemoryBounds = errors.New("h3: " + describeH3Error(E_MEMORY_BOUNDS))
	// ErrOptionInvalid corresponds to E_OPTION_INVALID: mode or flags
	// argument was not valid.
	ErrOptionInvalid = errors.New("h3: " + describeH3Error(E_OPTION_INVALID))
)

// errTable maps H3 C error codes to the public sentinel errors. Index 0
// (E_SUCCESS) is nil.
var errTable = [16]error{
	E_SUCCESS:            nil,
	E_FAILED:             ErrFailed,
	E_DOMAIN:             ErrDomain,
	E_LATLNG_DOMAIN:      ErrLatLngDomain,
	E_RES_DOMAIN:         ErrResolutionDomain,
	E_CELL_INVALID:       ErrCellInvalid,
	E_DIR_EDGE_INVALID:   ErrDirectedEdgeInvalid,
	E_UNDIR_EDGE_INVALID: ErrUndirectedEdgeInvalid,
	E_VERTEX_INVALID:     ErrVertexInvalid,
	E_PENTAGON:           ErrPentagon,
	E_DUPLICATE_INPUT:    ErrDuplicateInput,
	E_NOT_NEIGHBORS:      ErrNotNeighbors,
	E_RES_MISMATCH:       ErrResolutionMismatch,
	E_MEMORY_ALLOC:       ErrMemoryAlloc,
	E_MEMORY_BOUNDS:      ErrMemoryBounds,
	E_OPTION_INVALID:     ErrOptionInvalid,
}

// toErr converts an internal H3 C error code to the public sentinel error.
// E_SUCCESS maps to nil; codes outside the known range map to ErrFailed.
func toErr(errC H3Error) error {
	if int(errC) < len(errTable) {
		return errTable[errC]
	}
	return ErrFailed
}
