// Package h3 provides a pure-Go implementation of Uber's H3 hexagonal hierarchical
// geospatial index, targeting behavioral equivalence with H3 C v4.3.0.
package h3

import (
	"github.com/dimchansky/h3-go/internal/indexbits"
	"github.com/dimchansky/h3-go/internal/tables"
)

// Input validation helpers
func validateResolution(res int) error {
	if res < 0 || res > MaxResolution {
		return ErrResolutionDomain
	}
	return nil
}

func validateLatLng(lat, lng float64) error {
	if lat < -90.0 || lat > 90.0 {
		return ErrLatLngDomain
	}
	
	if lng <= -180.0 || lng > 180.0 {
		return ErrLatLngDomain
	}
	
	return nil
}

func validateCell(c Cell) error {
	if !indexbits.IsValidCell(uint64(c)) {
		return ErrCellInvalid
	}
	return nil
}

func validateKValue(k int) error {
	if k < 0 {
		return ErrDomain
	}
	return nil
}

func validateCellPair(a, b Cell) error {
	if err := validateCell(a); err != nil {
		return err
	}
	if err := validateCell(b); err != nil {
		return err
	}
	
	// Check resolution match
	resA := indexbits.GetResolution(uint64(a))
	resB := indexbits.GetResolution(uint64(b))
	if resA != resB {
		return ErrResolutionMismatch
	}
	
	return nil
}

// IsValidCell checks if the given Cell represents a valid H3 index.
// This performs bit-level validation of the index structure.
func IsValidCell(c Cell) bool {
	return indexbits.IsValidCell(uint64(c))
}

// Resolution returns the resolution of the given cell.
// Resolution ranges from 0 (coarsest) to 15 (finest).
// Returns ErrCellInvalid if the cell is not valid.
func Resolution(c Cell) (int, error) {
	if err := validateCell(c); err != nil {
		return 0, err
	}
	return indexbits.GetResolution(uint64(c)), nil
}

// BaseCell returns the base cell number (0-121) for the given cell.
// Returns ErrCellInvalid if the cell is not valid.
func BaseCell(c Cell) (int, error) {
	if err := validateCell(c); err != nil {
		return 0, err
	}
	return indexbits.GetBaseCell(uint64(c)), nil
}

// IsPentagon returns true if the given cell is a pentagon.
// Pentagons occur at each resolution due to the spherical geometry.
// Returns ErrCellInvalid if the cell is not valid.
func IsPentagon(c Cell) (bool, error) {
	if err := validateCell(c); err != nil {
		return false, err
	}
	
	baseCell := indexbits.GetBaseCell(uint64(c))
	
	// Check if base cell is a pentagon
	if tables.IsPentagonBaseCell(baseCell) {
		// At resolution 0, pentagon base cells are pentagons
		res := indexbits.GetResolution(uint64(c))
		if res == 0 {
			return true, nil
		}
		
		// At higher resolutions, check if all digits are 0
		// (center child of pentagon is also a pentagon)
		isPent := true
		for i := 0; i < res; i++ {
			if indexbits.GetDigit(uint64(c), i) != 0 {
				isPent = false
				break
			}
		}
		return isPent, nil
	}
	
	return false, nil
}

// Additional helper type for range-with-distance results.
type CellDistance struct {
	Cell     Cell
	Distance int
}

// Stub functions that return errors for now

// LatLngToCell converts a geographic coordinate to an H3 cell at the specified resolution.
// Resolution must be between 0 and 15 inclusive.
// Latitude must be in [-90, 90], longitude in (-180, 180].
func LatLngToCell(p LatLng, res int) (Cell, error) {
	if err := validateResolution(res); err != nil {
		return 0, err
	}
	if err := validateLatLng(p.Lat, p.Lng); err != nil {
		return 0, err
	}
	
	// TODO: Implement actual conversion
	return 0, ErrOptionInvalid
}

// CellToLatLng returns the center point of the given cell.
func CellToLatLng(c Cell) (LatLng, error) {
	if err := validateCell(c); err != nil {
		return LatLng{}, err
	}
	
	// TODO: Implement actual conversion
	return LatLng{}, ErrOptionInvalid
}

// CellToBoundary returns the boundary vertices of the given cell.
// The returned vertices are in counterclockwise order starting from
// a canonical vertex. Hexagons have 6 vertices, pentagons have 5.
// The dst buffer is reused if it has sufficient capacity.
func CellToBoundary(dst []LatLng, c Cell) ([]LatLng, error) {
	if err := validateCell(c); err != nil {
		return nil, err
	}
	
	// TODO: Implement actual boundary calculation
	return nil, ErrOptionInvalid
}

// AreNeighbors returns true if the two cells are neighbors.
// Cells must be at the same resolution.
func AreNeighbors(a, b Cell) (bool, error) {
	if err := validateCellPair(a, b); err != nil {
		return false, err
	}
	
	// TODO: Implement actual neighbor check
	return false, ErrOptionInvalid
}

// GridDistance returns the grid distance between two cells.
// Cells must be at the same resolution.
func GridDistance(a, b Cell) (int, error) {
	if err := validateCellPair(a, b); err != nil {
		return 0, err
	}
	
	// TODO: Implement actual distance calculation
	return 0, ErrOptionInvalid
}

// KRing returns all cells within k grid steps of the origin cell.
// Results are returned in ascending order by cell index.
// The dst buffer is reused if it has sufficient capacity.
func KRing(dst []Cell, origin Cell, k int) ([]Cell, error) {
	if err := validateCell(origin); err != nil {
		return nil, err
	}
	if err := validateKValue(k); err != nil {
		return nil, err
	}
	
	// TODO: Implement actual k-ring calculation
	return nil, ErrOptionInvalid
}

// HexRange is a synonym for KRing with bounded traversal semantics.
func HexRange(dst []Cell, origin Cell, k int) ([]Cell, error) {
	return KRing(dst, origin, k)
}

// HexRangeDistances returns all cells within k grid steps of the origin,
// annotated with their distance from the origin.
// Results are ordered by (distance, cell index).
// The dst buffer is reused if it has sufficient capacity.
func HexRangeDistances(dst []CellDistance, origin Cell, k int) ([]CellDistance, error) {
	if err := validateCell(origin); err != nil {
		return nil, err
	}
	if err := validateKValue(k); err != nil {
		return nil, err
	}
	
	// TODO: Implement actual range with distances
	return nil, ErrOptionInvalid
}

// HexRing returns all cells exactly k grid steps from the origin.
// Results are returned in ring-walk order (documented canonical order).
// The dst buffer is reused if it has sufficient capacity.
func HexRing(dst []Cell, origin Cell, k int) ([]Cell, error) {
	if err := validateCell(origin); err != nil {
		return nil, err
	}
	if err := validateKValue(k); err != nil {
		return nil, err
	}
	
	// TODO: Implement actual hex ring
	return nil, ErrOptionInvalid
}

// MaxKRingSize returns the maximum number of cells in a k-ring.
// Formula: 3*k*(k+1) + 1
func MaxKRingSize(k int) int {
	if k < 0 {
		return 0
	}
	return 3*k*(k+1) + 1
}