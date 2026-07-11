package h3

// _gridRingInternal implements the safe but slow version of gridRing algorithm.
//
// This function uses _gridDiskDistancesInternal to get all cells up to distance
// k, then filters those results to include only cells exactly at distance k.
//
// Ported from H3 C: algos.c::_gridRingInternal.
func _gridRingInternal(origin h3Index, k int32, out []h3Index) h3Error {
	// Short-circuit on 'identity' ring
	if k == 0 {
		out[0] = origin
		return eSuccess
	}

	var maxIdx int64
	err := maxGridDiskSize(k, &maxIdx)
	if err != eSuccess {
		return err
	}

	// Allocate temporary buffers for disk and distances
	diskOut := make([]h3Index, maxIdx)
	diskDistances := make([]int32, maxIdx)

	err = _gridDiskDistancesInternal(origin, k, diskOut, diskDistances, maxIdx, 0)
	if err != eSuccess {
		return err
	}

	currentIdx := 0
	for i := int64(0); i < maxIdx; i++ {
		if diskOut[i] != 0 && diskDistances[i] == k {
			out[currentIdx] = diskOut[i]
			currentIdx++
		}
	}

	return eSuccess
}
