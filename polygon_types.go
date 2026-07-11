package h3

// ContainmentMode mirrors C ContainmentMode for polygon flags.
type ContainmentMode int

const (
	ContainmentCenter          ContainmentMode = 0
	ContainmentFull            ContainmentMode = 1
	ContainmentOverlapping     ContainmentMode = 2
	ContainmentOverlappingBBox ContainmentMode = 3
	ContainmentInvalid         ContainmentMode = 4
)

// Polygon flag helpers.
const flagContainmentModeMask uint32 = 15

func flagGetContainmentMode(flags uint32) ContainmentMode {
	return ContainmentMode(flags & flagContainmentModeMask)
}
