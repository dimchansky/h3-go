package h3

// gridDiskDistancesSafe is the safe but slow version of gridDiskDistances.
//
// This function computes cells and their distances from the given origin cell, up to distance k.
// It is the public API equivalent of _gridDiskDistancesInternal but calculates the required
// buffer size automatically.
//
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and
// all neighboring cells, and so on.
//
// Output is placed in the provided array in no particular order. Elements of
// the output array may be left zero, as can happen when crossing a pentagon.
//
// Ported from H3 C: algos.c::gridDiskDistancesSafe.
func gridDiskDistancesSafe(origin h3Index, k int32, out []h3Index, distances []int32) h3Error {
	var maxIdx int64
	err := maxGridDiskSize(k, &maxIdx)
	if err != eSuccess {
		return err
	}
	return _gridDiskDistancesInternal(origin, k, out, distances, maxIdx, 0)
}
