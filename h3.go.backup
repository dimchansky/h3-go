package h3

import (
	"errors"
	"strconv"
	"strings"
	"unsafe"
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

const (
	MaxResolution = MAX_H3_RES
)

type (
	// Cell is an Index that identifies a single hexagon cell at a resolution.
	Cell H3Index

	// DirectedEdge is an Index that identifies a directed edge between two cells.
	DirectedEdge H3Index

	// Vertex is an Index that identifies a single topological vertex, shared by three cells.
	// A vertex is arbitrarily assigned one of the three neighboring cells as its "owner", which is used to calculate
	// the canonical index and geographic coordinates for the vertex.
	Vertex H3Index
)

// NewLatLng is a helper function to create a LatLng.
func NewLatLng(lat, lng Angle) LatLng {
	return LatLng{Lat: lat, Lng: lng}
}

// LatLngToCell returns the Cell at resolution for a geographic coordinate.
func LatLngToCell(latLng LatLng, resolution int) (Cell, error) {
	var out H3Index

	errC := latLngToCell(&latLng, int32(resolution), &out)

	return Cell(out), toErr(errC)
}

// Cell returns the Cell at resolution for a geographic coordinate.
func (g LatLng) Cell(resolution int) (Cell, error) {
	return LatLngToCell(g, resolution)
}

// CellToLatLng returns the geographic centerpoint of a Cell.
func CellToLatLng(c Cell) (LatLng, error) {
	var g LatLng

	errC := cellToLatLng(H3Index(c), &g)

	return g, toErr(errC)
}

// LatLng returns the Cell at resolution for a geographic coordinate.
func (c Cell) LatLng() (LatLng, error) {
	return CellToLatLng(c)
}

// CellToBoundary returns a CellBoundary of the Cell.
func CellToBoundary(c Cell) ([]LatLng, error) {
	var cb CellBoundary

	errC := cellToBoundary(H3Index(c), &cb)

	return cb.Verts, toErr(errC)
}

// Boundary returns a CellBoundary of the Cell.
func (c Cell) Boundary() ([]LatLng, error) {
	return CellToBoundary(c)
}

// GridDisk produces cells within grid distance k of the origin cell.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
//
// Output is placed in an array in no particular order. Elements of the output
// array may be left zero, as can happen when crossing a pentagon.
func GridDisk(origin Cell, k int) ([]Cell, error) {
	var outSize int64
	if err := toErr(maxGridDiskSize(int32(k), &outSize)); err != nil {
		return nil, err
	}

	out := make([]Cell, outSize)
	errC := gridDisk(H3Index(origin), int32(k), castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// GridDisk produces cells within grid distance k of the origin cell.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
//
// Output is placed in an array in no particular order. Elements of the output
// array may be left zero, as can happen when crossing a pentagon.
func (c Cell) GridDisk(k int) ([]Cell, error) {
	return GridDisk(c, k)
}

// GridDisksUnsafe produces cells within grid distance k of all provided origin
// cells.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func GridDisksUnsafe(origins []Cell, k int) ([]Cell, error) {
	var gridDiskSize int64
	if err := toErr(maxGridDiskSize(int32(k), &gridDiskSize)); err != nil {
		return nil, err
	}

	out := make([]Cell, int64(len(origins))*gridDiskSize)
	errC := gridDisksUnsafe(castSlice[Cell, H3Index](origins), int32(k), castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// GridDiskDistances produces cells within grid distance k of the origin cell.
// This method optimistically tries the faster GridDiskDistancesUnsafe first.
// If a cell was a pentagon or was in the pentagon distortion area, it falls
// back to GridDiskDistancesSafe.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func GridDiskDistances(origin Cell, k int) (outCells []Cell, outDists []int32, err error) {
	var rsz int64
	if err := toErr(maxGridDiskSize(int32(k), &rsz)); err != nil {
		return nil, nil, err
	}

	outCells = make([]Cell, rsz)
	outDists = make([]int32, rsz)
	err = toErr(gridDiskDistances(H3Index(origin), int32(k), castSlice[Cell, H3Index](outCells), outDists))

	return
}

// GridDiskDistances produces cells within grid distance k of the origin cell.
// This method optimistically tries the faster GridDiskDistancesUnsafe first.
// If a cell was a pentagon or was in the pentagon distortion area, it falls
// back to GridDiskDistancesSafe.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func (c Cell) GridDiskDistances(k int) (outCells []Cell, outDists []int32, err error) {
	return GridDiskDistances(c, k)
}

// GridDiskDistancesUnsafe produces cells within grid distance k of the origin cell.
// Output behavior is undefined when one of the cells returned by this
// function is a pentagon or is in the pentagon distortion area.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func GridDiskDistancesUnsafe(origin Cell, k int) (outCells []Cell, outDists []int32, err error) {
	var rsz int64
	if err := toErr(maxGridDiskSize(int32(k), &rsz)); err != nil {
		return nil, nil, err
	}

	outCells = make([]Cell, rsz)
	outDists = make([]int32, rsz)
	err = toErr(gridDiskDistancesUnsafe(H3Index(origin), int32(k), castSlice[Cell, H3Index](outCells), outDists))

	return
}

// GridDiskDistancesUnsafe produces cells within grid distance k of the origin cell.
// Output behavior is undefined when one of the cells returned by this
// function is a pentagon or is in the pentagon distortion area.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func (c Cell) GridDiskDistancesUnsafe(k int) (outCells []Cell, outDists []int32, err error) {
	return GridDiskDistancesUnsafe(c, k)
}

// GridDiskDistancesSafe produces cells within grid distance k of the origin cell.
// This is the safe, but slow version of GridDiskDistances.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func GridDiskDistancesSafe(origin Cell, k int) (outCells []Cell, outDists []int32, err error) {
	var rsz int64
	if err := toErr(maxGridDiskSize(int32(k), &rsz)); err != nil {
		return nil, nil, err
	}

	outCells = make([]Cell, rsz)
	outDists = make([]int32, rsz)
	err = toErr(gridDiskDistancesSafe(H3Index(origin), int32(k), castSlice[Cell, H3Index](outCells), outDists))

	return
}

// GridDiskDistancesSafe produces cells within grid distance k of the origin cell.
// This is the safe, but slow version of GridDiskDistances.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
func (c Cell) GridDiskDistancesSafe(k int) (outCells []Cell, outDists []int32, err error) {
	return GridDiskDistancesSafe(c, k)
}

// GridRing produces the "hollow" ring of cells at exactly grid distance k from the origin cell.
//
// k-ring 0 returns just the origin hexagon.
//
// Elements of the output array may be left zero, as can happen when crossing a pentagon.
func GridRing(origin Cell, k int) ([]Cell, error) {
	if k < 0 {
		return nil, ErrDomain
	}

	out := make([]Cell, ringSize(int32(k)))
	errC := gridRing(H3Index(origin), int32(k), castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// GridRing produces the "hollow" ring of cells at exactly grid distance k from the origin cell.
//
// k-ring 0 returns just the origin hexagon.
//
// Elements of the output array may be left zero, as can happen when crossing a pentagon.
func (c Cell) GridRing(k int) ([]Cell, error) {
	return GridRing(c, k)
}

// GridRingUnsafe produces the "hollow" ring of cells at exactly grid distance k from the origin cell.
//
// k-ring 0 returns just the origin hexagon.
func GridRingUnsafe(origin Cell, k int) ([]Cell, error) {
	if k < 0 {
		return nil, ErrDomain
	}

	out := make([]Cell, ringSize(int32(k)))
	errC := gridRingUnsafe(H3Index(origin), int32(k), castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// GridRingUnsafe produces the "hollow" ring of cells at exactly grid distance k from the origin cell.
//
// k-ring 0 returns just the origin hexagon.
func (c Cell) GridRingUnsafe(k int) ([]Cell, error) {
	return GridRingUnsafe(c, k)
}

// PolygonToCells takes a given GeoJSON-like data structure fills it with the
// hexagon cells that are contained by the GeoJSON-like data structure.
//
// This implementation traces the GeoJSON geoloop(s) in cartesian space with
// hexagons, tests them and their neighbors to be contained by the geoloop(s),
// and then any newly found hexagons are used to test again until no new
// hexagons are found.
func PolygonToCells(polygon GeoPolygon, resolution int) ([]Cell, error) {
	if len(polygon.GeoLoop) == 0 {
		return nil, nil
	}

	var maxLen int64
	if err := toErr(maxPolygonToCellsSize(&polygon, int32(resolution), 0, &maxLen)); err != nil {
		return nil, err
	}

	out := make([]Cell, maxLen)
	errC := polygonToCells(&polygon, int32(resolution), 0, castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// Cells takes a given GeoJSON-like data structure fills it with the
// hexagon cells that are contained by the GeoJSON-like data structure.
//
// This implementation traces the GeoJSON geoloop(s) in cartesian space with
// hexagons, tests them and their neighbors to be contained by the geoloop(s),
// and then any newly found hexagons are used to test again until no new
// hexagons are found.
func (p GeoPolygon) Cells(resolution int) ([]Cell, error) {
	return PolygonToCells(p, resolution)
}

// CellsToLinkedMultiPolygon creates a LinkedGeoPolygon describing the outline(s) of a set of hexagons.
// Polygon outlines will follow GeoJSON MultiPolygon order: Each polygon will
// have one outer loop, which is first in the list, followed by any holes.
//
// It is expected that all hexagons in the set have the same resolution and
// that the set contains no duplicates. Behavior is undefined if duplicates
// or multiple resolutions are present, and the algorithm may produce
// unexpected or invalid output.
func CellsToLinkedMultiPolygon(cells []Cell) (out LinkedGeoPolygon, err error) {
	if len(cells) == 0 {
		return
	}

	errC := cellsToLinkedMultiPolygon(castSlice[Cell, H3Index](cells), int32(len(cells)), &out)

	return out, toErr(errC)
}

// GreatCircleDistanceRads returns the "great circle" or "haversine" distance between
// pairs of LatLng points (lat/lng pairs) in radians.
func GreatCircleDistanceRads(a, b LatLng) float64 {
	return greatCircleDistanceRads(&a, &b)
}

// GreatCircleDistanceKm returns the "great circle" or "haversine" distance between pairs
// of LatLng points (lat/lng pairs) in kilometers.
func GreatCircleDistanceKm(a, b LatLng) float64 {
	return greatCircleDistanceKm(&a, &b)
}

// GreatCircleDistanceM returns the "great circle" or "haversine" distance between pairs
// of LatLng points (lat/lng pairs) in meters.
func GreatCircleDistanceM(a, b LatLng) float64 {
	return greatCircleDistanceM(&a, &b)
}

// HexagonAreaAvgKm2 returns the average hexagon area in square kilometers at the given
// resolution.
func HexagonAreaAvgKm2(resolution int) (float64, error) {
	var out float64

	errC := getHexagonAreaAvgKm2(int32(resolution), &out)

	return out, toErr(errC)
}

// HexagonAreaAvgM2 returns the average hexagon area in square meters at the given
// resolution.
func HexagonAreaAvgM2(resolution int) (float64, error) {
	var out float64

	errC := getHexagonAreaAvgM2(int32(resolution), &out)

	return out, toErr(errC)
}

// CellAreaRads2 returns the exact area of specific cell in square radians.
func CellAreaRads2(c Cell) (float64, error) {
	out, errC := cellAreaRads2(H3Index(c))
	return out, toErr(errC)
}

// CellAreaKm2 returns the exact area of specific cell in square kilometers.
func CellAreaKm2(c Cell) (float64, error) {
	out, errC := cellAreaKm2(H3Index(c))

	return out, toErr(errC)
}

// CellAreaM2 returns the exact area of specific cell in square meters.
func CellAreaM2(c Cell) (float64, error) {
	out, errC := cellAreaM2(H3Index(c))

	return out, toErr(errC)
}

// HexagonEdgeLengthAvgKm returns the average hexagon edge length in kilometers
// at the given resolution.
func HexagonEdgeLengthAvgKm(resolution int) (float64, error) {
	var out float64
	errC := getHexagonEdgeLengthAvgKm(int32(resolution), &out)

	return out, toErr(errC)
}

// HexagonEdgeLengthAvgM returns the average hexagon edge length in meters at
// the given resolution.
func HexagonEdgeLengthAvgM(resolution int) (float64, error) {
	var out float64

	errC := getHexagonEdgeLengthAvgM(int32(resolution), &out)

	return out, toErr(errC)
}

// EdgeLengthRads returns the exact edge length of specific unidirectional edge
// in radians.
func EdgeLengthRads(e DirectedEdge) (float64, error) {
	var out float64

	errC := edgeLengthRads(H3Index(e), &out)

	return out, toErr(errC)
}

// EdgeLengthKm returns the exact edge length of specific unidirectional
// edge in kilometers.
func EdgeLengthKm(e DirectedEdge) (float64, error) {
	var out float64

	errC := edgeLengthKm(H3Index(e), &out)

	return out, toErr(errC)
}

// EdgeLengthM returns the exact edge length of specific unidirectional
// edge in meters.
func EdgeLengthM(e DirectedEdge) (float64, error) {
	var out float64

	errC := edgeLengthM(H3Index(e), &out)

	return out, toErr(errC)
}

// NumCells returns the number of cells at the given resolution.
func NumCells(resolution int) (int64, error) {
	out, errC := getNumCells(int32(resolution))
	return out, toErr(errC)
}

// Res0Cells returns all the cells at resolution 0.
func Res0Cells() ([]Cell, error) {
	out := make([]Cell, res0CellCount())
	errC := getRes0Cells(castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// Pentagons returns all the pentagons at resolution.
func Pentagons(resolution int) ([]Cell, error) {
	out := make([]Cell, NUM_PENTAGONS)
	errC := getPentagons(int32(resolution), castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// Resolution returns the resolution of the cell.
func (c Cell) Resolution() int {
	return int(getResolution(H3Index(c)))
}

// Resolution returns the resolution of the edge.
func (e DirectedEdge) Resolution() int {
	return int(getResolution(H3Index(e)))
}

// BaseCellNumber returns the integer ID (0-121) of the base cell the H3Index h
// belongs to.
func BaseCellNumber(h Cell) int32 {
	return getBaseCellNumber(H3Index(h))
}

// BaseCellNumber returns the integer ID (0-121) of the base cell the H3Index h
// belongs to.
func (c Cell) BaseCellNumber() int32 {
	return BaseCellNumber(c)
}

const (
	base16  = 16
	bitSize = 64
)

// IndexFromString returns an uint64 from a string. Should call c.IsValid() to check
// if the Cell is valid before using it.
func IndexFromString(s string) uint64 {
	if len(s) > 2 && strings.ToLower(s[:2]) == "0x" {
		s = s[2:]
	}
	c, _ := strconv.ParseUint(s, base16, bitSize)

	return c
}

// IndexToString returns a string from a Cell.
func IndexToString(i uint64) string {
	return strconv.FormatUint(i, base16)
}

// CellFromString returns a Cell from a string. Should call c.IsValid() to check
// if the Cell is valid before using it.
func CellFromString(s string) Cell {
	return Cell(IndexFromString(s))
}

// CellToString returns a string from a Cell.
func CellToString(c Cell) string {
	return IndexToString(uint64(c))
}

// VertexFromString returns a Vertex from a string. Should call v.IsValid() to check
// if the Vertex is valid before using it.
func VertexFromString(s string) Vertex {
	return Vertex(IndexFromString(s))
}

// String returns the string representation of the H3Index h.
func (c Cell) String() string {
	return CellToString(c)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c Cell) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (c *Cell) UnmarshalText(text []byte) error {
	*c = Cell(IndexFromString(string(text)))
	if !c.IsValid() {
		return errors.New("invalid cell index")
	}

	return nil
}

// IsValid returns if a Cell is a valid cell (hexagon or pentagon).
func (c Cell) IsValid() bool {
	return c != 0 && isValidCell(H3Index(c))
}

// Parent returns the parent or grandparent Cell of this Cell.
func (c Cell) Parent(resolution int) (Cell, error) {
	out, errC := cellToParent(H3Index(c), int32(resolution))

	return Cell(out), toErr(errC)
}

// ImmediateParent returns the immediate parent of the cell.
func (c Cell) ImmediateParent() (Cell, error) {
	return c.Parent(c.Resolution() - 1)
}

// Children returns the children or grandchildren cells of this Cell.
func (c Cell) Children(resolution int) ([]Cell, error) {
	outSz, errC := cellToChildrenSize(H3Index(c), int32(resolution))
	if err := toErr(errC); err != nil {
		return nil, err
	}

	out := make([]Cell, outSz)
	errC = cellToChildren(H3Index(c), int32(resolution), castSlice[Cell, H3Index](out))

	return out, toErr(errC)
}

// ImmediateChildren returns the children or grandchildren cells of this Cell.
func (c Cell) ImmediateChildren() ([]Cell, error) {
	return c.Children(c.Resolution() + 1)
}

// CenterChild returns the center child Cell of this Cell.
func (c Cell) CenterChild(resolution int) (Cell, error) {
	out, errC := cellToCenterChild(H3Index(c), int32(resolution))

	return Cell(out), toErr(errC)
}

// IsResClassIII returns true if this is a class III index. If false, this is a
// class II index.
func (c Cell) IsResClassIII() bool {
	return isResClassIII(H3Index(c))
}

// IsPentagon returns true if this is a pentagon.
func (c Cell) IsPentagon() bool {
	return isPentagon(H3Index(c))
}

// IcosahedronFaces finds all icosahedron faces (0-19) intersected by this Cell.
func (c Cell) IcosahedronFaces() ([]int32, error) {
	var outsz int32
	if err := toErr(maxFaceCount(H3Index(c), &outsz)); err != nil {
		return nil, err
	}

	out := make([]int32, outsz)
	errC := getIcosahedronFaces(H3Index(c), out)

	return out, toErr(errC)
}

// IsNeighbor returns true if this Cell is a neighbor of the other Cell.
func (c Cell) IsNeighbor(other Cell) (bool, error) {
	out, errC := areNeighborCells(H3Index(c), H3Index(other))

	return out, toErr(errC)
}

// DirectedEdge returns a DirectedEdge from this Cell to other.
func (c Cell) DirectedEdge(other Cell) (DirectedEdge, error) {
	out, errC := cellsToDirectedEdge(H3Index(c), H3Index(other))

	return DirectedEdge(out), toErr(errC)
}

const (
	numCellEdges    = 6
	numEdgeCells    = 2
	numCellVertexes = 6
)

// DirectedEdges returns 6 directed edges with h as the origin.
func (c Cell) DirectedEdges() ([]DirectedEdge, error) {
	out := make([]DirectedEdge, numCellEdges) // always 6 directed edges

	// Seems like this function always returns E_SUCCESS.
	errC := originToDirectedEdges(H3Index(c), castSlice[DirectedEdge, H3Index](out))

	return out, toErr(errC)
}

// IsValid determines if the directed edge is valid.
func (e DirectedEdge) IsValid() bool {
	return isValidDirectedEdge(H3Index(e))
}

// Origin returns the origin cell of this directed edge.
func (e DirectedEdge) Origin() (Cell, error) {
	out, errC := getDirectedEdgeOrigin(H3Index(e))

	return Cell(out), toErr(errC)
}

// Destination returns the destination cell of this directed edge.
func (e DirectedEdge) Destination() (Cell, error) {
	var out H3Index
	errC := getDirectedEdgeDestination(H3Index(e), &out)

	return Cell(out), toErr(errC)
}

// Cells returns the origin and destination cells in that order.
func (e DirectedEdge) Cells() ([]Cell, error) {
	out := make([]Cell, numEdgeCells)
	if err := toErr(directedEdgeToCells(H3Index(e), castSlice[Cell, H3Index](out))); err != nil {
		return nil, err
	}

	return out, nil
}

// Boundary provides the coordinates of the boundary of the directed edge. Note,
// the type returned is CellBoundary, but the coordinates will be from the
// center of the origin to the center of the destination. There may be more than
// 2 coordinates to account for crossing faces.
func (e DirectedEdge) Boundary() ([]LatLng, error) {
	var out CellBoundary
	err := toErr(directedEdgeToBoundary(H3Index(e), &out))

	return out.Verts, err
}

// CompactCells merges full sets of children into their parent H3Index
// recursively, until no more merges are possible.
func CompactCells(in []Cell) ([]Cell, error) {
	csz := len(in)
	out := make([]Cell, csz)
	errC := compactCells(castSlice[Cell, H3Index](in), castSlice[Cell, H3Index](out), int64(csz))

	return out, toErr(errC)
}

// UncompactCells splits every H3Index in in if its resolution is greater
// than resolution recursively. Returns all the H3Indexes at resolution resolution.
func UncompactCells(in []Cell, resolution int) ([]Cell, error) {
	csz, cErr := uncompactCellsSize(castSlice[Cell, H3Index](in), int64(len(in)), int32(resolution))
	if err := toErr(cErr); err != nil {
		return nil, err
	}

	out := make([]Cell, csz)
	errC := uncompactCells(castSlice[Cell, H3Index](in), int64(len(in)), castSlice[Cell, H3Index](out), csz, int32(resolution))

	return out, toErr(errC)
}

// ChildPosToCell returns the child of cell a at a given position within an ordered list of all
// children at the specified resolution.
func ChildPosToCell(position int64, a Cell, resolution int) (Cell, error) {
	out, errC := childPosToCell(position, H3Index(a), int32(resolution))

	return Cell(out), toErr(errC)
}

// ChildPosToCell returns the child cell at a given position within an ordered list of all
// children at the specified resolution.
func (c Cell) ChildPosToCell(position int64, resolution int) (Cell, error) {
	return ChildPosToCell(position, c, resolution)
}

// CellToChildPos returns the position of the cell a within an ordered list of all children of the cell's parent
// at the specified resolution.
func CellToChildPos(a Cell, resolution int) (int64, error) {
	out, errC := cellToChildPos(H3Index(a), int32(resolution))

	return out, toErr(errC)
}

// ChildPos returns the position of the cell within an ordered list of all children of the cell's parent
// at the specified resolution.
func (c Cell) ChildPos(resolution int) (int64, error) {
	return CellToChildPos(c, resolution)
}

// GridDistance returns grid distance between two cells.
//
// This function may fail to find the distance between two indexes, for example if they are very far apart. It may also
// fail when finding distances for indexes on opposite sides of a pentagon.
func GridDistance(a, b Cell) (int64, error) {
	var out int64
	errC := gridDistance(H3Index(a), H3Index(b), &out)

	return out, toErr(errC)
}

// GridDistance returns grid distance between two cells.
//
// This function may fail to find the distance between two indexes, for example if they are very far apart. It may also
// fail when finding distances for indexes on opposite sides of a pentagon.
func (c Cell) GridDistance(other Cell) (int64, error) {
	return GridDistance(c, other)
}

// GridPath returns the line of cells between the two cells (inclusive).
//
// This function may fail to find the line between two indexes, for example if they are very far apart. It may also fail
// when finding distances for indexes on opposite sides of a pentagon.
func GridPath(a, b Cell) ([]Cell, error) {
	var outSz int64
	if err := toErr(gridPathCellsSize(H3Index(a), H3Index(b), &outSz)); err != nil {
		return nil, err
	}

	out := make([]Cell, outSz)
	if err := toErr(gridPathCells(castSlice[Cell, H3Index](out), H3Index(a), H3Index(b))); err != nil {
		return nil, err
	}

	return out, nil
}

// GridPath returns the line of cells between the two cells (inclusive).
//
// This function may fail to find the line between two indexes, for example if they are very far apart. It may also fail
// when finding distances for indexes on opposite sides of a pentagon.
func (c Cell) GridPath(other Cell) ([]Cell, error) {
	return GridPath(c, other)
}

// CellToLocalIJ produces ij coordinates for cell anchored by an origin.
//
// The coordinate space used by this function may have deleted regions or warping due to pentagonal distortion.
//
// Coordinates are only comparable if they come from the same origin index.
//
// Failure may occur if the index is too far away from the origin or if the index is on the other side of a pentagon.
func CellToLocalIJ(origin, cell Cell) (CoordIJ, error) {
	var out CoordIJ
	errC := cellToLocalIj(H3Index(origin), H3Index(cell), 0, &out)

	return out, toErr(errC)
}

// LocalIJToCell produces a cell for ij coordinates anchored by an origin.
//
// The coordinate space used by this function may have deleted regions or warping due to pentagonal distortion.
//
// Failure may occur if the index is too far away from the origin or if the index is on the other side of a pentagon.
func LocalIJToCell(origin Cell, ij CoordIJ) (Cell, error) {
	var out H3Index
	errC := localIjToCell(H3Index(origin), &ij, 0, &out)

	return Cell(out), toErr(errC)
}

// CellToVertex returns a single vertex for a given cell, or InvalidH3Index if the vertex is invalid.
func CellToVertex(c Cell, vertexNum int) (Vertex, error) {
	var out H3Index
	errC := cellToVertex(H3Index(c), int32(vertexNum), &out)

	return Vertex(out), toErr(errC)
}

// Vertex returns a single vertex for a given cell, or InvalidH3Index if the vertex is invalid.
func (c Cell) Vertex(vertexNum int) (Vertex, error) {
	return CellToVertex(c, vertexNum)
}

// CellToVertexes returns all vertexes for the given cell.
func CellToVertexes(c Cell) ([]Vertex, error) {
	var out [6]Vertex
	h3Arr := asH3Array6(&out)
	if err := toErr(cellToVertexes(H3Index(c), h3Arr)); err != nil {
		return nil, err
	}
	return out[:], nil
}

// Vertexes returns all vertexes for the given cell.
func (c Cell) Vertexes() ([]Vertex, error) {
	return CellToVertexes(c)
}

// VertexToLatLng returns the geographic coordinates of the vertex.
func VertexToLatLng(vertex Vertex) (LatLng, error) {
	var out LatLng
	errC := vertexToLatLng(H3Index(vertex), &out)
	return out, toErr(errC)
}

// LatLng returns the geographic coordinates of the vertex.
func (v Vertex) LatLng() (LatLng, error) {
	return VertexToLatLng(v)
}

// IsValidVertex returns whether the cell is a valid vertex.
func IsValidVertex(v Vertex) bool {
	return isValidVertex(H3Index(v))
}

// IsValid returns whether the cell is a valid vertex.
func (v Vertex) IsValid() bool {
	return IsValidVertex(v)
}

// String returns a string from a Vertex.
func (v Vertex) String() string {
	return IndexToString(uint64(v))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (v Vertex) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (v *Vertex) UnmarshalText(text []byte) error {
	*v = VertexFromString(string(text))
	if !v.IsValid() {
		return errors.New("invalid cell index")
	}

	return nil
}

const (
	latLngFloatPrecision = 5
	// latLngStringSize is the size to pre-allocate the buffer for.
	// Given latLngFloatPrecision, a typical string is "(DD.DDDDD, -DDD.DDDDD)"
	// which is ~25-30 bytes. 32 is a safe and efficient capacity to start with
	// to avoid re-allocation.
	latLngStringSize = 32
)

func (g LatLng) String() string {
	buf := make([]byte, 0, latLngStringSize)
	buf = append(buf, '(')
	buf = strconv.AppendFloat(buf, g.Lat.Deg(), 'f', latLngFloatPrecision, 64) //nolint:mnd // float bit size
	buf = append(buf, ',', ' ')
	buf = strconv.AppendFloat(buf, g.Lng.Deg(), 'f', latLngFloatPrecision, 64) //nolint:mnd // float bit size
	buf = append(buf, ')')
	return string(buf)
}

func ringSize(k int32) int64 {
	if k == 0 {
		return 1
	}

	return int64(6 * k)
}

func toErr(errC H3Error) error {
	if int(errC) < len(errMap) {
		return errMap[errC]
	}
	return ErrUnknown
}

// castSlice reinterprets a []From as []To without copying.
func castSlice[From ~uint64, To ~uint64](in []From) []To {
	return unsafe.Slice((*To)(unsafe.Pointer(unsafe.SliceData(in))), len(in))
}

func asH3Array6[T ~uint64](p *[6]T) *[6]H3Index {
	return (*[6]H3Index)(unsafe.Pointer(p))
}
