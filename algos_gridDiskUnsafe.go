package h3

// gridDiskUnsafe produces indexes within k distance of the origin index.
// Output behavior is undefined when one of the indexes returned by this
// function is a pentagon or is in the pentagon distortion area.
// k-ring 0 is defined as the origin index, k-ring 1 is defined as k-ring 0 and
// all neighboring indexes, and so on.
// Output is placed in the provided array in order of increasing distance from
// the origin.
// Ported from H3 C: algos.c::gridDiskUnsafe.
func gridDiskUnsafe(origin H3Index, k int32, out []H3Index) H3Error {
	return gridDiskDistancesUnsafe(origin, k, out, nil)
}
