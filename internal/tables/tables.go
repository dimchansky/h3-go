// Package tables provides static lookup tables for H3 operations.
// These tables are stubbed with correct types and sizes but zero values.
// TODO: Populate with actual H3 v4.3.0 data.
package tables

// BaseCellData contains metadata for each base cell.
type BaseCellData struct {
	// Face is the icosahedron face this base cell belongs to (0-19).
	Face int
	
	// IsPentagon indicates if this base cell is a pentagon.
	IsPentagon bool
	
	// IJK0 is the canonical IJK coordinates on the face.
	IJK0 [3]int
	
	// CWOffsetPent is the rotation offset for pentagon cells.
	CWOffsetPent [2]int
}

// NeighborRotation describes a rotation to reach a neighbor.
type NeighborRotation struct {
	// Direction is the neighbor direction (0-6).
	Direction int
	
	// RotationCCW is the counter-clockwise rotation count.
	RotationCCW int
}

// FaceIJK represents a face and IJK coordinate pair.
type FaceIJK struct {
	Face int
	IJK  [3]int
}

// Constants for table sizes.
const (
	// NumBaseCells is the total number of base cells (122 in H3).
	NumBaseCells = 122
	
	// NumIcosahedronFaces is the number of icosahedron faces.
	NumIcosahedronFaces = 20
	
	// MaxCellNeighbors is the maximum number of neighbors a cell can have.
	MaxCellNeighbors = 6
	
	// NumPentagons is the number of pentagon base cells.
	NumPentagons = 12
)

// BaseCells contains metadata for all base cells.
// TODO: Populate with actual base cell data from H3 v4.3.0.
var BaseCells = [NumBaseCells]BaseCellData{
	// Stubbed with zero values for now
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	{}, {},
}

// PentagonBaseCells lists the indices of all pentagon base cells.
// TODO: Populate with actual pentagon base cell indices.
var PentagonBaseCells = [NumPentagons]int{
	// Pentagon base cells (12 total in H3)
	// These are placeholders - actual values: 4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117
	4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117,
}

// IsPentagonBaseCell returns true if the given base cell is a pentagon.
func IsPentagonBaseCell(baseCell int) bool {
	if baseCell < 0 || baseCell >= NumBaseCells {
		return false
	}
	
	// Check against known pentagon base cells
	for _, p := range PentagonBaseCells {
		if baseCell == p {
			return true
		}
	}
	return false
}

// BaseCellNeighbors contains the neighbor relationships for base cells.
// TODO: Populate with actual neighbor data.
// Each entry is [7]int where index 0-5 are the neighbor base cells
// (or -1 if no neighbor in that direction), and index 6 is unused.
var BaseCellNeighbors = [NumBaseCells][7]int{
	// Stubbed with -1 (no neighbor) for now
	{-1, -1, -1, -1, -1, -1, -1},
	// ... repeat for all 122 base cells
}

// BaseCellNeighborRotations contains rotation data for base cell neighbors.
// TODO: Populate with actual rotation data.
var BaseCellNeighborRotations = [NumBaseCells][7]NeighborRotation{
	// Stubbed with zero values for now
}

// FaceNeighbors maps each icosahedron face to its neighboring faces.
// TODO: Populate with actual face neighbor data.
var FaceNeighbors = [NumIcosahedronFaces][3]FaceIJK{
	// Stubbed with zero values for now
}

// ResolutionAreaKm2 contains the average hexagon area in km² at each resolution.
// TODO: Verify these values against H3 v4.3.0.
var ResolutionAreaKm2 = [16]float64{
	4250546.848, 607220.978, 86745.854, 12392.265,
	1770.348, 252.907, 36.129, 5.161,
	0.737, 0.105, 0.015, 0.0021,
	0.0003, 0.00004, 0.000006, 0.0000009,
}

// ResolutionEdgeLengthKm contains the average edge length in km at each resolution.
// TODO: Verify these values against H3 v4.3.0.
var ResolutionEdgeLengthKm = [16]float64{
	1107.712, 418.677, 158.245, 59.811,
	22.607, 8.544, 3.229, 1.220,
	0.461, 0.174, 0.066, 0.025,
	0.009, 0.0035, 0.0013, 0.0005,
}

// DirectionToNeighbor maps a direction (0-6) to IJK deltas.
// Direction 0 is the center (self), 1-6 are the six neighbors.
// TODO: Verify IJK coordinate system and directions.
var DirectionToNeighbor = [7][3]int{
	{0, 0, 0},   // CENTER
	{0, 1, 0},   // Direction 1
	{1, 0, 0},   // Direction 2
	{1, 1, 0},   // Direction 3
	{0, 0, 1},   // Direction 4
	{-1, 0, 1},  // Direction 5
	{0, -1, 1},  // Direction 6
}