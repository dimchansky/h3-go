package c2go

// gridDiskDistances produces cells and their distances from the given origin cell, up to distance k.
// k-ring 0 is defined as the origin cell, k-ring 1 is defined as k-ring 0 and all neighboring cells, and so on.
// Output is placed in the provided array in no particular order. Elements of the output array may be left zero,
// as can happen when crossing a pentagon.
// Ported from H3 C: algos.c::gridDiskDistances
func gridDiskDistances(origin H3Index, k int32, out []H3Index, distances []int32) H3Error {
	// Optimistically try the faster gridDiskUnsafe algorithm first
	failed := gridDiskDistancesUnsafe(origin, k, out, distances)
	if failed != E_SUCCESS {
		var maxIdx int64
		err := maxGridDiskSize(k, &maxIdx)
		if err != E_SUCCESS {
			return err
		}
		// Fast algo failed, fall back to slower, correct algo
		// and also wipe out array because contents untrustworthy
		for i := range out[:maxIdx] {
			out[i] = 0
		}

		if distances == nil {
			// Allocate temporary distances array
			tempDistances := make([]int32, maxIdx)
			result := _gridDiskDistancesInternal(origin, k, out, tempDistances, maxIdx, 0)
			return result
		} else {
			// Clear distances array
			for i := range distances[:maxIdx] {
				distances[i] = 0
			}
			return _gridDiskDistancesInternal(origin, k, out, distances, maxIdx, 0)
		}
	}
	return E_SUCCESS
}
