package h3

// Index fits within a 64-bit unsigned integer.
type Index = uint64

const (
	// MAX_CELL_BNDRY_VERTS is a maximum number of cell boundary vertices; worst case is pentagon:
	// 5 original verts + 5 edge crossings
	MAX_CELL_BNDRY_VERTS = 10
)

// GeoCoord is a struct for geographic coordinates, holds latitude/longitude in radians.
type GeoCoord struct {
	Lat float64 // latitude in radians
	Lon float64 // longitude in radians
}

// GeoBoundary is a slice of `GeoCoord`.  Note, `len(GeoBoundary)` will never
// exceed `MAX_CELL_BNDRY_VERTS`. Vertices are in ccw order.
type GeoBoundary []GeoCoord

// GeoPolygon is a geofence with 0 or more geofence holes
type GeoPolygon struct {
	// GeoFence is the exterior boundary of the polygon
	GeoFence []GeoCoord

	// Holes is a slice of interior boundary (holes) in the polygon
	Holes [][]GeoCoord
}
