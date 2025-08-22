package c2go

// ContainmentMode mirrors C ContainmentMode for polygon flags
type ContainmentMode int

const (
    CONTAINMENT_CENTER          ContainmentMode = 0
    CONTAINMENT_FULL            ContainmentMode = 1
    CONTAINMENT_OVERLAPPING     ContainmentMode = 2
    CONTAINMENT_OVERLAPPING_BBOX ContainmentMode = 3
    CONTAINMENT_INVALID         ContainmentMode = 4
)

// Polygon flag helpers
const FLAG_CONTAINMENT_MODE_MASK uint32 = 15

func FLAG_GET_CONTAINMENT_MODE(flags uint32) ContainmentMode {
    return ContainmentMode(flags & FLAG_CONTAINMENT_MODE_MASK)
}

