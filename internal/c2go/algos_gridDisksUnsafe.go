package c2go

// gridDisksUnsafe takes an array of input hex IDs and a max k-ring and returns
// an array of hexagon IDs sorted first by the original hex IDs and then by the
// k-ring (0 to max), with no guaranteed sorting within each k-ring group.
// The memory block should be equal to maxGridDiskSize(k) * length
// Ported from H3 C: algos.c::gridDisksUnsafe
func gridDisksUnsafe(h3Set []H3Index, k int32, out []H3Index) H3Error {
	if len(h3Set) == 0 || len(out) == 0 {
		return E_FAILED
	}
	
	length := int32(len(h3Set))
	var segmentSize int64
	err := maxGridDiskSize(k, &segmentSize)
	if err != E_SUCCESS {
		return err
	}

	for i := int32(0); i < length; i++ {
		// Determine the appropriate segment of the output array to operate on
		segmentStart := i * int32(segmentSize)
		segment := out[segmentStart : segmentStart+int32(segmentSize)]
		failed := gridDiskUnsafe(h3Set[i], k, segment)
		if failed != E_SUCCESS {
			return failed
		}
	}
	return E_SUCCESS
}
