package h3

// ContainmentMode selects how PolygonToCellsExperimental decides whether a
// cell belongs to the polygon.
//
// H3 C API: ContainmentMode (h3api.h).
type ContainmentMode int

// Containment modes for PolygonToCellsExperimental.
const (
	// ContainmentCenter includes a cell when its center point is contained
	// in the polygon.
	ContainmentCenter ContainmentMode = 0
	// ContainmentFull includes a cell only when it is fully contained in
	// the polygon.
	ContainmentFull ContainmentMode = 1
	// ContainmentOverlapping includes a cell when any part of it overlaps
	// the polygon.
	ContainmentOverlapping ContainmentMode = 2
	// ContainmentOverlappingBBox includes a cell when any part of its
	// bounding box overlaps the polygon (faster, may include false
	// positives).
	ContainmentOverlappingBBox ContainmentMode = 3
	// ContainmentInvalid marks the end of the valid mode range; it is not a
	// usable mode.
	ContainmentInvalid ContainmentMode = 4
)

// Polygon flag helpers.
const flagContainmentModeMask uint32 = 15

func flagGetContainmentMode(flags uint32) ContainmentMode {
	return ContainmentMode(flags & flagContainmentModeMask)
}
