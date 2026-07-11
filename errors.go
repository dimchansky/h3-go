package h3

import "errors"

// Sentinel errors returned by the public API, one per H3 C error code
// (H3ErrorCodes in h3api.h). Match with errors.Is. The message text is the C
// library's describeH3Error text with an "h3: " prefix.
var (
	// ErrFailed corresponds to eFailed: the operation failed but a more
	// specific error is not available.
	ErrFailed = errors.New("h3: " + describeH3Error(eFailed))
	// ErrDomain corresponds to eDomain: argument was outside of acceptable range.
	ErrDomain = errors.New("h3: " + describeH3Error(eDomain))
	// ErrLatLngDomain corresponds to eLatlngDomain: latitude or longitude
	// arguments were outside of acceptable range.
	ErrLatLngDomain = errors.New("h3: " + describeH3Error(eLatlngDomain))
	// ErrResolutionDomain corresponds to eResDomain: resolution argument
	// was outside of acceptable range.
	ErrResolutionDomain = errors.New("h3: " + describeH3Error(eResDomain))
	// ErrCellInvalid corresponds to eCellInvalid: cell argument was not valid.
	ErrCellInvalid = errors.New("h3: " + describeH3Error(eCellInvalid))
	// ErrDirectedEdgeInvalid corresponds to eDirEdgeInvalid: directed edge
	// argument was not valid.
	ErrDirectedEdgeInvalid = errors.New("h3: " + describeH3Error(eDirEdgeInvalid))
	// ErrUndirectedEdgeInvalid corresponds to eUndirEdgeInvalid:
	// undirected edge argument was not valid.
	ErrUndirectedEdgeInvalid = errors.New("h3: " + describeH3Error(eUndirEdgeInvalid))
	// ErrVertexInvalid corresponds to eVertexInvalid: vertex argument was
	// not valid.
	ErrVertexInvalid = errors.New("h3: " + describeH3Error(eVertexInvalid))
	// ErrPentagon corresponds to ePentagon: pentagon distortion was
	// encountered and the algorithm could not proceed.
	ErrPentagon = errors.New("h3: " + describeH3Error(ePentagon))
	// ErrDuplicateInput corresponds to eDuplicateInput: duplicate input.
	ErrDuplicateInput = errors.New("h3: " + describeH3Error(eDuplicateInput))
	// ErrNotNeighbors corresponds to eNotNeighbors: cell arguments were not
	// neighbors.
	ErrNotNeighbors = errors.New("h3: " + describeH3Error(eNotNeighbors))
	// ErrResolutionMismatch corresponds to eResMismatch: cell arguments had
	// incompatible resolutions.
	ErrResolutionMismatch = errors.New("h3: " + describeH3Error(eResMismatch))
	// ErrMemoryAlloc corresponds to eMemoryAlloc: memory allocation failed.
	ErrMemoryAlloc = errors.New("h3: " + describeH3Error(eMemoryAlloc))
	// ErrMemoryBounds corresponds to eMemoryBounds: bounds of provided
	// memory were insufficient.
	ErrMemoryBounds = errors.New("h3: " + describeH3Error(eMemoryBounds))
	// ErrOptionInvalid corresponds to eOptionInvalid: mode or flags
	// argument was not valid.
	ErrOptionInvalid = errors.New("h3: " + describeH3Error(eOptionInvalid))
)

// errTable maps H3 C error codes to the public sentinel errors. Index 0
// (eSuccess) is nil.
var errTable = [16]error{
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
}

// toErr converts an internal H3 C error code to the public sentinel error.
// eSuccess maps to nil; codes outside the known range map to ErrFailed.
func toErr(errC h3Error) error {
	if int(errC) < len(errTable) {
		return errTable[errC]
	}
	return ErrFailed
}
