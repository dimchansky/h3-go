package h3

import (
	"errors"
)

// Error codes.
var (
	ErrFailed                = errors.New(describeH3Error(1))
	ErrDomain                = errors.New(describeH3Error(2))
	ErrLatLngDomain          = errors.New(describeH3Error(3))
	ErrResolutionDomain      = errors.New(describeH3Error(4))
	ErrCellInvalid           = errors.New(describeH3Error(5))
	ErrDirectedEdgeInvalid   = errors.New(describeH3Error(6))
	ErrUndirectedEdgeInvalid = errors.New(describeH3Error(7))
	ErrVertexInvalid         = errors.New(describeH3Error(8))
	ErrPentagon              = errors.New(describeH3Error(9))
	ErrDuplicateInput        = errors.New(describeH3Error(10))
	ErrNotNeighbors          = errors.New(describeH3Error(11))
	ErrResolutionMismatch    = errors.New(describeH3Error(12))
	ErrMemoryAlloc           = errors.New(describeH3Error(13))
	ErrMemoryBounds          = errors.New(describeH3Error(14))
	ErrOptionInvalid         = errors.New(describeH3Error(15))

	ErrUnknown = errors.New("unknown error code returned by H3")

	errMap = [16]error{
		0:  nil, // Success error code.
		1:  ErrFailed,
		2:  ErrDomain,
		3:  ErrLatLngDomain,
		4:  ErrResolutionDomain,
		5:  ErrCellInvalid,
		6:  ErrDirectedEdgeInvalid,
		7:  ErrUndirectedEdgeInvalid,
		8:  ErrVertexInvalid,
		9:  ErrPentagon,
		10: ErrDuplicateInput,
		11: ErrNotNeighbors,
		12: ErrResolutionMismatch,
		13: ErrMemoryAlloc,
		14: ErrMemoryBounds,
		15: ErrOptionInvalid,
	}
)

type (
	// Cell is an Index that identifies a single hexagon cell at a resolution.
	Cell = H3Index
)

// NewLatLng is a helper function to create a LatLng.
func NewLatLng(lat, lng Angle) LatLng {
	return LatLng{Lat: lat, Lng: lng}
}

// LatLngToCell returns the Cell at resolution for a geographic coordinate.
func LatLngToCell(latLng LatLng, resolution int) (Cell, error) {
	var out H3Index

	errC := latLngToCell(&latLng, int32(resolution), &out)

	return out, toErr(errC)
}

// Cell returns the Cell at resolution for a geographic coordinate.
func (g LatLng) Cell(resolution int) (Cell, error) {
	return LatLngToCell(g, resolution)
}

func toErr(errC H3Error) error {
	if int(errC) < len(errMap) {
		return errMap[errC]
	}
	return ErrUnknown
}
