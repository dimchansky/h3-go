package h3

// Cell is an H3 cell index: a single hexagon (or pentagon) in the H3 grid
// system at a particular resolution.
//
// The zero value is not a valid cell; IsValid reports false for it. A Cell
// may be constructed by conversion from a raw uint64 index (unchecked — use
// IsValid to verify), parsed from its hexadecimal string form with ParseCell,
// or produced by operations such as LatLngToCell.
type Cell uint64

// Curated public constants. Each is defined in terms of the corresponding
// ported H3 C constant so upstream changes propagate mechanically.
const (
	// MaxResolution is the finest (highest) H3 resolution. Resolutions run
	// from 0 (coarsest) to MaxResolution (finest).
	// H3 C API: MAX_H3_RES.
	MaxResolution = maxH3Res

	// NumBaseCells is the number of resolution-0 base cells.
	// H3 C API: NUM_BASE_CELLS.
	NumBaseCells = numBaseCells

	// NumRes0Cells is the number of resolution-0 cells, equal to NumBaseCells.
	// H3 C API: res0CellCount.
	NumRes0Cells = numBaseCells

	// NumPentagons is the number of pentagonal cells per resolution.
	// H3 C API: pentagonCount.
	NumPentagons = numPentagons

	// MaxCellBoundaryVerts is the maximum number of vertices a cell boundary
	// can have (pentagon cells crossing icosahedron edges produce up to 10).
	// H3 C API: MAX_CELL_BNDRY_VERTS.
	MaxCellBoundaryVerts = 10

	// VersionMajor, VersionMinor and VersionPatch identify the H3 C release
	// this library is behaviorally equivalent to.
	// H3 C API: H3_VERSION_MAJOR, H3_VERSION_MINOR, H3_VERSION_PATCH.
	VersionMajor = h3VersionMajor
	VersionMinor = h3VersionMinor
	VersionPatch = h3VersionPatch
)
